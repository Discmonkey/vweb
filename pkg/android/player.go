package android

import (
	"context"
	"errors"
	"github.com/discmonkey/vweb/pkg/video"
	"net"
	"time"
)

type Player struct {
	sps         []byte
	pps         []byte
	source      chan video.Frame
	subscribers map[chan video.Frame]*receiverState
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

func (f Frame) Free() {}

type receiverState struct {
	idrSent bool
}

func safeSend(channel chan video.Frame, bytes []byte) {
	select {
	case channel <- Frame{bytes: bytes}:
	default:
	}
}
func NewPlayer(port int) (video.Player, context.CancelFunc, error) {
	newFrameSource := make(chan video.Frame)

	p := Player{
		sps:         nil,
		pps:         nil,
		source:      newFrameSource,
		subscribers: make(map[chan video.Frame]*receiverState),
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
				p.subscribers[s] = &receiverState{}
			case s := <-p.unsubscribe:
				delete(p.subscribers, s)
			case f := <-out:
				for channel := range p.subscribers {
					safeSend(channel, f)
				}
			}
		}
	}()

	return &p, cancel, nil
}

func (p *Player) Listen(ctxt context.Context, port int) (chan []byte, error) {
	addr := net.UDPAddr{
		Port: port,
		IP:   net.ParseIP("0.0.0.0"),
	}
	conn, err := net.ListenUDP("udp", &addr) // code does not block here
	if err != nil {
		return nil, err
	}

	input := make(chan []byte)
	go func() {
		h264Parser{}.parse(ctxt, conn, input)
	}()

	return input, nil
}

var _ video.Player = &Player{}
var _ video.Frame = &Frame{}
