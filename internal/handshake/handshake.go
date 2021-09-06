package handshake

import (
	"context"
)

type Handshake struct {
	value *uint64
}

func (h *Handshake) Shake(ctx context.Context, callback func()) {
	atomic.Add
}
