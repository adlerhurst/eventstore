package eventstore

import "sync"

type eventPool[E Event] struct {
	pool sync.Pool
}

func NewEventPool[E Event]() *eventPool[E] {
	return &eventPool[E]{
		pool: sync.Pool{
			New: func() any { return new(E) },
		},
	}
}

func (p *eventPool[E]) Get() E {
	return p.pool.Get().(E)
}

func (p *eventPool[E]) Put(event E) {
	p.pool.Put(event)
}
