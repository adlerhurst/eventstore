package memory

import (
	"github.com/adlerhurst/eventstore"
	"github.com/adlerhurst/eventstore/storage"
)

type subject struct {
	topic  eventstore.TextSubject
	events []*eventstore.Event
	subs   []*subject
}

func (n *subject) push(subjects []eventstore.TextSubject, cmd eventstore.Command, seq uint64) (*eventstore.Event, error) {
	if len(subjects) == 0 {
		payload, err := storage.PayloadToBytes(cmd.Payload())
		if err != nil {
			return nil, err
		}
		e := &eventstore.Event{
			EditorService: cmd.EditorService(),
			EditorUser:    cmd.EditorService(),
			Subjects:      cmd.Subjects(),
			Payload:       payload,
			Sequence:      seq,
			ResourceOwner: cmd.ResourceOwner(),
		}
		n.events = append(n.events, e)
		return e, nil
	}
	for _, sub := range n.subs {
		if sub.topic == subjects[0] {
			return sub.push(subjects[1:], cmd, seq)
		}
	}

	sub := &subject{topic: subjects[0]}
	n.subs = append(n.subs, sub)

	return sub.push(subjects[1:], cmd, seq)
}

func (n *subject) find(subjects []eventstore.Subject) (events []*eventstore.Event) {
	if len(subjects) == 0 {
		return nil
	}
	if s, ok := subjects[0].(eventstore.TextSubject); ok {
		if n.topic != s {
			return nil
		} else if len(subjects) == 1 {
			return n.events
		} else {
			for _, sub := range n.subs {
				events = append(events, sub.find(subjects[1:])...)
			}
			return events
		}
	} else if subjects[0] == eventstore.MultiToken {
		return n.getAll()
	} else if subjects[0] == eventstore.SingleToken {
		for _, sub := range n.subs {
			events = append(events, sub.find(subjects[1:])...)
		}
	}
	return events
}

func (n *subject) getAll() (events []*eventstore.Event) {
	events = n.events
	for _, sub := range n.subs {
		events = append(events, sub.getAll()...)
	}
	return events
}
