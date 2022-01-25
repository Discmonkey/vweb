package android

import (
	"context"
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

type parseState = byte

const (
	NA   parseState = 0
	O___ parseState = 1
	OO__ parseState = 2
	OOO_ parseState = 3
	OOO1 parseState = 4
)

func (p h264Parser) parse(cancel context.Context, con net.Conn, out chan []byte) {
	// reader
	input := make([]byte, 2048)
	output := make([]byte, 0, 2048)
	sentSps := false
	sentPps := false

	var state = NA
	for n, err := con.Read(input); err == nil; {
		var c byte
		for i := 0; i < n; i++ {
			c = input[i]
			switch c {
			case 0:
				switch state {
				case NA:
					state = O___
				case O___:
					state = OO__
				case OO__:
					state = OOO_
				case OOO_:
					output = append(output, 0)
				}
			default:
				switch state {
				case NA:
					output = append(output, c)
				case O___:
					output = append(output, 0, c)
					state = NA
				case OO__:
					output = append(output, 0, 0, c)
					state = NA
				case OOO_:
					if input[i] == 1 {
						if len(output) > 0 {
							isPPS := output[0] == PPS
							isSPS := output[0] == SPS
							if sentPps && sentSps || isSPS || isPPS {
								out <- output
								//select {
								//case out <- output:
								//case <-cancel.Done():
								//	return
								//}
								sentSps = sentSps || isSPS
								sentPps = sentPps || isPPS
							}
						}
						output = output[:0]
						state = NA
					} else {
						output = append(output, 0, 0, 0, input[i])
					}
				}
			}
		}
	}
	close(out)
}
