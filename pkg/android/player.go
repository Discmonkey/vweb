package android

import (
	"context"
	"errors"
	"fmt"
	"github.com/discmonkey/vweb/pkg/video"
	"net"
	"time"
)

type Player struct {
	sps         []byte
	pps         []byte
	source      chan video.Frame
	subscribers map[chan video.Frame]bool
	subscribe   chan chan video.Frame
	unsubscribe chan chan video.Frame
	cancel      context.CancelFunc
}

func (p *Player) Stop() {
	if p.cancel != nil {
		p.cancel()
	}
}

func (p Player) Type() video.Type {
	return video.H264
}

func (p Player) Play() (chan video.Frame, context.CancelFunc, error) {
	out := make(chan video.Frame)

	select {
	case <-time.After(time.Second * 5):
		return nil, nil, errors.New("could not subscribe to stream")
	case p.subscribe <- out:
	}

	cancel := func() {
		select {
		case <-time.After(time.Second * 5):
		case p.unsubscribe <- out:
		}
	}

	return out, cancel, nil
}

type Frame struct {
	count int
	bytes []byte
}

func (f Frame) Bytes() ([]byte, error) {
	//TODO implement me
	return f.bytes, nil
}

func (f Frame) Count() (int, error) {
	//TODO implement me
	return 0, nil
}

func NewPlayer(port int) (video.Player, context.CancelFunc, error) {
	newFrameSource := make(chan video.Frame)

	p := Player{
		sps:         nil,
		pps:         nil,
		source:      newFrameSource,
		subscribers: make(map[chan video.Frame]bool),
		subscribe:   make(chan chan video.Frame),
		unsubscribe: make(chan chan video.Frame),
		cancel:      nil,
	}
	ctxt, cancel := context.WithCancel(context.Background())
	out, err := p.Listen(ctxt, port)
	if err != nil {
		cancel()
		return nil, nil, err
	}

	p.cancel = cancel
	go func() {
		for {
			select {
			case <-ctxt.Done():
				return
			case s := <-p.subscribe:
				p.subscribers[s] = true
				s <- Frame{bytes: p.pps}
				s <- Frame{bytes: p.sps}
			case s := <-p.unsubscribe:
				delete(p.subscribers, s)
			case f := <-out:
				for channel := range p.subscribers {
					select {
					case channel <- Frame{bytes: f}:
					default:
					}
				}
			}
		}
	}()

	return &p, cancel, nil
}

func (p *Player) Listen(ctxt context.Context, port int) (chan []byte, error) {
	addr := net.TCPAddr{
		Port: port,
		IP:   net.ParseIP("0.0.0.0"),
	}
	listener, err := net.ListenTCP("tcp", &addr) // code does not block here
	if err != nil {
		return nil, err
	}
	conn, err := listener.AcceptTCP()
	if err != nil {
		return nil, err
	}

	input := make(chan []byte)
	go h264Parser{}.parse(ctxt, conn, input)
	// wait to get the sps and pps and out
	timeout := time.After(time.Minute * 3)

	for p.sps == nil || p.pps == nil {
		var next []byte
		select {
		case next = <-input:
		case <-timeout:
			return nil, errors.New("could not find sps and pps")
		}
		if next[0] == PPS {
			fmt.Println("found pps")
			p.pps = next
		} else if next[0] == SPS {
			fmt.Println("found sps")
			p.sps = next
		}
	}

	return input, nil
}

var _ video.Player = &Player{}
var _ video.Frame = &Frame{}
