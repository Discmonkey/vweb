package video

import "sync/atomic"

type Router struct {
	out map[uint64]chan Frame
	add chan subscriber
	del chan uint64
	id  int
}

type subscriber struct {
	c  chan Frame
	id uint64
}

var id uint64 = 0

func NewRouter(source chan Frame) *Router {
	r := &Router{
		out: make(map[uint64]chan Frame),
		add: make(chan subscriber),
	}

	go func() {
		for {
			next := <-source
			select {
			case subscriber := <-r.add:
				r.out[subscriber.id] = subscriber.c
			case id := <-r.del:
				delete(r.out, id)
			default:
				for _, out := range r.out {
					out <- next
				}
			}
		}
	}()

	return r
}

func (r *Router) AddOut(c chan Frame) Unsubscribe {
	id := atomic.AddUint64(&id, 1)
	r.add <- subscriber{
		c:  c,
		id: id,
	}

	return func() { r.del <- id }
}
