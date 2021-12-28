package filestramer

import (
	"fmt"
	"github.com/discmonkey/vweb/internal/image"
	"github.com/discmonkey/vweb/pkg/video"
	"os"
	"path"
)

type Player struct {
	loop bool
	directory string
	index int
}

func NewPlayer(dir string, loop bool) *Player {
	return &Player {
		loop: loop,
		directory: dir,
		index: 0,
	}
}

type Frame struct {
	bytes []byte
	aspect image.Aspect
	isKey bool
}

func (f* Frame) Aspect() (image.Aspect, error) {
	return f.aspect, nil
}

func (f* Frame) Bytes() []byte {
	return f.bytes
}

func (f* Frame) IsKey() bool {
	return f.isKey
}

func (f *Player) Next() (video.Frame, video.Count, error) {
	currentIndex := f.index
	f.index += 1

	next := path.Join(f.directory, fmt.Sprintf("%d", currentIndex))

	contents, err := os.ReadFile(next)

	if err != nil && currentIndex > 0 && f.loop {
		f.index = 0
		return f.Next()
	}

	return &Frame {
		bytes: contents,
		aspect: image.Aspect{},
		isKey: true,
	}, video.Count(currentIndex), err
}

func (f *Player) Encoding() (video.Encoding, error) {
	return video.H264, nil
}

var _ video.Player = &Player{}
var _ video.Frame = &Frame{}
