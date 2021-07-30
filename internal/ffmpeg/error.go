package ffmpeg

import "fmt"

type Error struct {
	code    int
	message string
}

func (e Error) Error() string {
	return fmt.Sprintf("FFMPEG ERROR %d: %s", e.code, e.message)
}

var _ error = Error{}
