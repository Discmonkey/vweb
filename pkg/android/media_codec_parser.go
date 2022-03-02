package android

import (
	"context"
	"github.com/discmonkey/vweb/internal/nal"
	"net"
)

// mediaCodecParser parses the buffers sent by an android media codec h264 encoder
type h264Parser struct {
}

type parserState byte

const (
	NA   parserState = 0
	O    parserState = 1
	OO   parserState = 2
	OOO  parserState = 3
	OOO1 parserState = 4 // never used
)

func (p h264Parser) parse(cancel context.Context, con net.Conn, out chan []byte) {
	// reader
	input := make([]byte, 256)
	output := make([]byte, 0, 0)

	defer func() {
		close(out)
	}()

	var last byte = 0
	shouldWriteLast := false
	clear := false

	state := NA
	for {
		n, err := con.Read(input)
		if err != nil {
			return
		}

		for i := 0; i < n; i++ {
			if input[i] == 0 {
				if state == NA || state == OO || state == O {
					state += 1
				}
			} else if input[i] == 1 && state == OOO {
				shouldWriteLast = true
				if len(output) > 3 {

					clear = !(nal.IsPPS(last) || nal.IsSPS(last))
					if clear {
						// shorten output by 3 since it would have 000 appended to it,
						// which technically belongs to the next buffer
						length := len(output) - 3
						outSend := make([]byte, length)
						copy(outSend, output[:length])

						select {
						case out <- outSend:
						case <-cancel.Done():
							return
						}
					}
				}

				if clear {
					output = output[:0]
					output = append(output, 0, 0, 0)
				}

				state = NA

			} else {
				if shouldWriteLast {
					last = input[i]
					shouldWriteLast = false
				}
				state = NA
			}
			output = append(output, input[i])
		}
	}
}
