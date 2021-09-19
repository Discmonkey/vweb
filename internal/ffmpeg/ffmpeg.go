package ffmpeg

import "github.com/discmonkey/vweb/internal/image"

type ParameterSet struct {
	Raw    []byte
	Parsed map[string]string
	Aspect image.Aspect
}

type Frame struct {
	ParameterSet *ParameterSet
	IsKey        bool
	Bytes        []byte
}
