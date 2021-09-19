package video

import (
	"github.com/discmonkey/vweb/internal/image"
	"github.com/pion/webrtc/v3"
	"time"
)

type Encoding = string

const (
	H264 Encoding = webrtc.MimeTypeH264
	JPEG          = "jpeg"
)

// Count is used for synchronization between video source and metadata
type Count = uint64

// Frame is a single quantum of a some video stream/container
type Frame interface {
	Aspect() (image.Aspect, error)
	Bytes() []byte
	Timestamp(Count) (time.Time, error)
}

// Player is a handle on some video source
type Player interface {
	Next() (Frame, Count, error)
	Encoding() (Encoding, error)
}
