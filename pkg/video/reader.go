package video

import (
	"github.com/dismonkey/vweb/internal/image"
	"time"
)

type Encoding = string
const (
	H264 Encoding = "h264"
	JPEG = "jpeg"
	VP8 = "vp8"
)

// Count is used for synchronization between video source and metadata
type Count = uint64

// Frame is a single quantum of a some video stream/container
type Frame interface {
	Aspect() (image.Aspect, error)
	Bytes() ([]byte, error)
	Timestamp() (time.Time, error)
}

// Reader is the main interface for the video package, it can be used to interop with different video sinks
type Reader interface {
	Next() (Frame, Count, error)
	Encoding() (Encoding, error)
}
