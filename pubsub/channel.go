package pubsub

import "github.com/adlerhurst/eventstore"

type Channel struct {
	m *eventstore.Mapping
}

func (c *Channel) Subscribe(ch chan<- eventstore.Event, subs ...eventstore.Subject) {
	eventstore.Register(c.m, subs, func(events ...eventstore.Event) {
		for _, event := range events {
			ch <- event
		}
	})
}

func (c *Channel) Publish(events ...eventstore.Event) {}
