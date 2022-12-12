package utils

import "time"

type FPS = int

type Timer struct {
	logEveryN int
	lastTick  time.Time
	logFps    bool
}

func NewTimer() Timer {
	return Timer{
		logEveryN: 1,
		lastTick:  time.Now(),
		logFps:    true,
	}
}

func (t *Timer) SetFps(fps bool) {
	t.logFps = fps
}

func (t *Timer) SetLogEverN(logEvery int) {
	if logEvery < 1 {
		panic("negative log every passed to timer")
	}

	t.logEveryN = logEvery
}

func (t *Timer) Tick() FPS {
	now := time.Now()
	duration := now.Sub(t.lastTick)
	t.lastTick = now
	// without the cast to float64 secs will be set to zero
	secs := float64(duration/time.Millisecond) / 1000.0

	if secs <= 0 {
		return -1
	}

	return FPS(1 / float64(secs))
}
