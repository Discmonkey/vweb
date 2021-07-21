package ffmpeg

type ParameterSet struct {
	raw    []byte
	parsed map[string]string
}

type Frame struct {
	ParameterSet *ParameterSet
}
