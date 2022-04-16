package video

import (
	"context"
	"github.com/pion/webrtc/v3"
)

type Type = string

const (
	H264 Type = webrtc.MimeTypeH264
	VP9  Type = webrtc.MimeTypeVP9
)

// Count is used for synchronization between video source and metadata
type Count = uint64

// Frame is a single quantum of a some video stream/container
type Frame interface {
	Bytes() ([]byte, error)
	Count() (int, error)
	Free()
}

type Unsubscribe = func()

// Player is a handle on some video source
type Player interface {
	Play() (chan Frame, context.CancelFunc, error)
	Type() Type
}
