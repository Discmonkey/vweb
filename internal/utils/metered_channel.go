package utils

import (
	"github.com/discmonkey/vweb/pkg/video"
	"time"
)

func MeterChannel(input chan video.Frame, frequency time.Duration) chan video.Frame {
	output := make(chan video.Frame, 20)
	go func() {
		ticker := time.NewTicker(time.Millisecond * 1000 / frequency)

		for {
			<-ticker.C
			f, open := <-input
			if !open {
				close(output)
				break
			}

			output <- f
		}
	}()

	return output
}
