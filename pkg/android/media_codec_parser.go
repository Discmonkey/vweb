package android

import (
	"fmt"
	"net"
)

// mediaCodecParser parses the buffers sent by an android media codec h264 encoder
type h264Parser struct {
}

const (
	SPS   byte = 103
	PPS        = 104
	OTHER      = 128 // the largest valid NAL byte value is 127
	UNSET      = 129
)

func (p h264Parser) parse(con net.Conn, spsChan chan []byte, ppsChan chan []byte, frameChan chan []byte) {
	// reader
	next := make([]byte, 1024)
	prev := make([]byte, 1024)
	ready := make(chan []byte)
	go func() {
		for n := -1; n >= 0; {
			var err error
			n, err = con.Read(prev)
			if err != nil {
				fmt.Println(err)
				close(ready)
			}

			if n > 0 {
				prev = prev[:n]
				ready <- prev
				prev, next = next, prev
			}
		}
	}()
	go func() {
		markerPosition := 0
		var dest []byte

		send := func() {
			switch dest[0] {
			case SPS:
				spsChan <- dest
			case PPS:
				ppsChan <- dest
			default:
				frameChan <- dest
			}
		}
		for bytes := range ready {
			for i := range bytes {
				switch bytes[i] {
				case 0:
					if markerPosition < 3 {
						markerPosition += 1
					} else {
						dest = append(dest, 0)
					}
				case 1:
					if markerPosition == 3 {
						markerPosition = 0
						if len(dest) > 0 {
							send()
						}
						dest = dest[:0]
					}
				default:
					for markerPosition > 0 {
						dest = append(dest, 0)
						markerPosition--
					}
					dest = append(dest, bytes[i])
				}
			}
		}
	}()
}
