package x

import (
	"sync"
)

type pool[T any] struct {
	pool sync.Pool
}

func NewPool[T any]() *pool[T] {
	return &pool[T]{
		pool: sync.Pool{
			New: func() any { return new(T) },
		},
	}
}

func (p *pool[T]) Get() *T {
	return p.pool.Get().(*T)
}

func (p *pool[T]) Put(object *T) {
	p.pool.Put(object)
}
