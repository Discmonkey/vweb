package android

import (
	"github.com/discmonkey/vweb/pkg/video"
	"net"
)

type Player struct {
	router *video.Router
	sps    []byte
	pps    []byte
}

func (p Player) Type() video.Type {
	return video.H264
}

func (p Player) Play() (chan video.Frame, video.Unsubscribe, error) {
	out := make(chan video.Frame)

	go func() {
		out <- Frame{
			bytes: p.sps,
		}

		out <- Frame{
			bytes: p.pps,
		}
	}()

	cancel := p.router.AddOut(out)

	return out, cancel, nil
}

type Frame struct {
	count int
	bytes []byte
}

func (f Frame) Bytes() ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (f Frame) Count() (int, error) {
	//TODO implement me
	panic("implement me")
}

func (f Frame) Type() video.Type {
	//TODO implement me
	panic("implement me")
}

func NewPlayer(port int) (video.Player, error) {
	spsC, ppsC, frameC := make(chan []byte), make(chan []byte), make(chan []byte)
	newFrameSource := make(chan video.Frame)
	addr := net.UDPAddr{
		Port: port,
		IP:   net.ParseIP("127.0.0.1"),
	}
	conn, err := net.ListenUDP("udp", &addr) // code does not block here
	if err != nil {
		return nil, err
	}
	h264Parser{}.parse(conn, spsC, ppsC, frameC)
	p := Player{
		router: video.NewRouter(newFrameSource),
	}

	go func() {
		{
			f := Frame{}
			select {
			case p.pps = <-spsC:
				f.bytes = p.pps
			case p.sps = <-ppsC:
				f.bytes = p.sps
			case f.bytes = <-frameC:
			}

			newFrameSource <- f
		}
	}()

	return &p, nil
}

var _ video.Player = &Player{}
var _ video.Frame = &Frame{}
