package android

import (
	"bytes"
	"context"
	"net"
)

// mediaCodecParser parses the buffers sent by an android media codec h264 encoder
type h264Parser struct {
}

type parseState = byte

const (
	NA   parseState = 0
	O___ parseState = 1
	OO__ parseState = 2
	OOO_ parseState = 3
	OOO1 parseState = 4
)

var sep = []byte{0, 0, 0, 1}

func (p h264Parser) parse(cancel context.Context, con net.Conn, out chan []byte) {
	// reader
	input := make([]byte, 16384)

	for {
		n, err := con.Read(input)
		if err != nil {
			close(out)
			return
		}

		outputs := bytes.Split(input[:n], sep)

		for _, slice := range outputs {
			if len(slice) == 0 {
				continue
			}
			output := make([]byte, len(slice)+4)
			copy(output, sep)
			copy(output[4:], slice)

			select {
			case out <- output:
			case <-cancel.Done():
				close(out)
				return
			}
		}

	}
}
