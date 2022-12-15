package android

import (
	"context"
	"github.com/discmonkey/vweb/internal/nal"
	"log"
	"net"
)

// mediaCodecParser parses the buffers sent by an android media codec h264 encoder
type h264Parser struct {
}

type parserState byte

const (
	READING parserState = 0
	O       parserState = 1
	OO      parserState = 2
	OOO     parserState = 3
	OOO1    parserState = 4
)

const extraAppendLength = 4

func copyUnit(in []byte, sps, pps []byte, nalUnit byte) []byte {
	// skip 0,0,0 being appended before this is called
	end := len(in) - extraAppendLength
	var out []byte
	if nal.IsIDR(nalUnit) {
		out = make([]byte, len(sps)+len(pps)+len(in[:end]))
		copy(out, sps)
		copy(out[len(sps):], pps)
		copy(out[len(pps)+len(sps):], in[:end])
	} else {
		out = make([]byte, len(in[:end]))
		copy(out, in[:end])
	}

	return out
}

func (p h264Parser) parse(cancel context.Context, con net.Conn, out chan []byte) {
	// reader
	input := make([]byte, 2048)
	readBuffer := make([]byte, 0, 0)

	defer func() {
		close(out)
		cancel.Done()
	}()

	var sps, pps []byte
	state := READING
	var current byte = 0
	for {
		n, err := con.Read(input)
		if err != nil {
			log.Println("error encountered reading udp connection", err)
			return
		}

		for i := 0; i < n; i++ {
			if input[i] == 0 && (state == READING || state == OO || state == O) {
				state += 1
			} else if state == OOO && input[i] == 1 {
				state = OOO1
			} else if state == OOO1 {
				if len(readBuffer) > extraAppendLength {
					lastUnit := copyUnit(readBuffer, sps, pps, current)
					if nal.IsSPS(current) {
						sps = lastUnit
					} else if nal.IsPPS(current) {
						pps = lastUnit
					} else {
						out <- lastUnit
					}
					readBuffer = readBuffer[:extraAppendLength]
					current = input[i]
				} else {
					current = input[i]
				}
				state = READING
			} else {
				state = READING
			}
			readBuffer = append(readBuffer, input[i])
		}
	}
}
