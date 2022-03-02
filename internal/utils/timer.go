package utils

import "time"

type FpsTimer struct {
	last time.Time
}

func (f *FpsTimer) Start() {
	f.last = time.Now()
}

func (f *FpsTimer) Tick() float64 {
	next := time.Now()
	dif := next.Sub(f.last)
	f.last = next

	return 1 * 1000.0 / float64(dif.Milliseconds())
}
