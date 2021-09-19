package ffmpeg

import "github.com/discmonkey/vweb/pkg/video"

// Player implements the videoPlayer interface with a ffmpeg backend
type Player struct {
}

// NewPlayer attempts to open url
func NewPlayer(url string) (video.Player, error) {
	return &Player{}, nil
}

// Next returns the next frame in the stream
func (p *Player) Next() (video.Frame, video.Count, error) {
	panic("implement me")
}

// Encoding for the ffmpeg player only supports H264
func (p *Player) Encoding() (video.Encoding, error) {
	return video.H264, nil
}
