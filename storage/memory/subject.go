package memory

import (
	"sort"

	"github.com/adlerhurst/eventstore"
)

type subject struct {
	topic  eventstore.TextSubject
	events events
	subs   []*subject
}

func (s *subject) push(subjects []eventstore.TextSubject, e *event) {
	if len(subjects) == 0 {
		s.events = append(s.events, e)
		return
	}
	for _, sub := range s.subs {
		if sub.topic == subjects[0] {
			sub.push(subjects[1:], e)
			return
		}
	}

	sub := &subject{topic: subjects[0]}
	s.subs = append(s.subs, sub)

	sub.push(subjects[1:], e)
}

func (s *subject) find(subjects []eventstore.Subject) (res events) {
	if len(subjects) == 0 {
		return nil
	}

	defer sort.Sort(res)

	if sub, ok := subjects[0].(eventstore.TextSubject); ok {
		if s.topic != sub {
			return nil
		} else if len(subjects) == 1 {
			return s.events
		} else {
			for _, sub := range s.subs {
				res = append(res, sub.find(subjects[1:])...)
			}
			return res
		}
	} else if subjects[0] == eventstore.MultiToken {
		res = s.getAll()
		return res
	} else if subjects[0] == eventstore.SingleToken {
		for _, sub := range s.subs {
			res = append(res, sub.find(subjects[1:])...)
		}
	}
	return res
}

func (n *subject) getAll() (res events) {
	res = n.events
	for _, sub := range n.subs {
		res = append(res, sub.getAll()...)
	}
	return res
}
