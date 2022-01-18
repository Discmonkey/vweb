package ffmpeg

import (
	"errors"
	"github.com/discmonkey/vweb/internal/image"
	"github.com/discmonkey/vweb/pkg/video"
	"reflect"
	"unsafe"
)

// #cgo LDFLAGS: -L${SRCDIR}/c/build -lpacket_reader -lavcodec -lavformat -lavutil
// #include "c/src/packet_reader.h"
import "C"

// Player implements the videoPlayer interface with a ffmpeg backend
type Player struct {
	stream *C.Stream
	count  video.Count
}

func (p *Player) Play() (chan video.Frame, error) {
	//TODO implement me
	panic("implement me")
}

// Frame implements the video.Frame interface with a ffmpeg backend
type Frame struct {
	data       []byte
	aspect     image.Aspect
	isKeyFrame bool
}

func (f *Frame) Bytes() ([]byte, error) {
	//TODO implement me
	return f.data, nil
}

func (f *Frame) Type() video.Type {
	//TODO implement me
	panic("implement me")
}

func (f *Frame) Aspect() (image.Aspect, error) {
	return f.aspect, nil
}

func (f *Frame) IsKey() bool {
	return f.isKeyFrame
}

// NewPlayer attempts to open url using ffmpeg c bindings
func NewPlayer(url string) (video.Player, error) {
	cUrl := C.CString(url)
	defer C.free(unsafe.Pointer(cUrl))

	streamOrError := C.open_stream(cUrl)
	if streamOrError.error != nil {
		defer C.free_error(streamOrError.error)
		return nil, errors.New(C.GoString(streamOrError.error.reason))
	}

	return &Player{
		stream: streamOrError.stream,
		count:  0,
	}, nil
}

// Next returns the next frame in the stream
func (p *Player) Next() (video.Frame, video.Count, error) {
	packetOrError := C.next_packet(p.stream)
	if packetOrError.error != nil {
		defer C.free_error(packetOrError.error)
		return nil, 0, errors.New(C.GoString(packetOrError.error.reason))
	}

	p.count += 1
	data := wrapRawBytes(packetOrError.packet.data, int(packetOrError.packet.size))
	aspect := image.Aspect{}
	isKeyFrame := bool(packetOrError.packet.is_key_frame)
	return &Frame{
		data,
		aspect,
		isKeyFrame,
	}, p.count, nil
}

// Encoding for the ffmpeg player only supports H264
func (p *Player) Encoding() (video.Encoding, error) {
	return video.H264, nil
}

// wrapRawBytes converts a raw pointer to a go slice
// note that the memory for the byte array is not managed by the go runtime
// function body is inspired by https://stackoverflow.com/questions/61961793/wrapping-allocated-byte-buffer-in-c-as-go-slice-byte
func wrapRawBytes(ptr unsafe.Pointer, sz int) []byte {
	h := reflect.SliceHeader{
		Data: uintptr(ptr),
		Len:  sz,
		Cap:  sz,
	}
	buf := *(*[]byte)(unsafe.Pointer(&h))
	return buf
}

var _ video.Player = &Player{}
var _ video.Frame = &Frame{}
