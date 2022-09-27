package memory

import (
	"context"
	"sync"

	"github.com/adlerhurst/eventstore"
)

type Memory struct {
	events          []*eventstore.Event
	currentSequence uint64
	seqMu           *sync.Mutex
}

func New() *Memory {
	return &Memory{
		seqMu: &sync.Mutex{},
	}
}

func (m *Memory) Ready(context.Context) error {
	return nil
}

func (m *Memory) Push(ctx context.Context, cmds []eventstore.Command) ([]*eventstore.Event, error) {
	return m.saveEvents(cmds...)
}

func (m *Memory) Filter(ctx context.Context, filter *eventstore.Filter) (res []*eventstore.Event, err error) {
	for _, event := range m.events {
		for i, sub := range filter.Subjects {
			if event.Sequence > sub.To {
				//no events will be found anymore for this subject
				copy(filter.Subjects[i:], filter.Subjects[i+1:])
				filter.Subjects[len(filter.Subjects)-1] = nil
				filter.Subjects = filter.Subjects[:len(filter.Subjects)-1]
				continue
			}
			if sub.From > event.Sequence {
				continue
			}
			if !matchSubjects(event, sub.Subjects) {
				continue
			}

			res = append(res, event)
			break
		}
		if filter.Limit == len(res) {
			return res, nil
		}
	}
	return nil, nil
}

func matchSubjects(event *eventstore.Event, filter []eventstore.Subject) bool {
	for i := 0; i < len(filter); i++ {
		if filter[i] == eventstore.MultiToken {
			return true
		}

		if filter[i] == eventstore.SingleToken {
			continue
		}

		if filter[i] != event.Subjects[i] {
			return false
		}

		//more subjects to compare
		if len(event.Subjects) > i+1 && len(filter) > i+1 {
			continue
		}

		return len(event.Subjects) == i+1 && len(filter) == i+1
	}
	return len(event.Subjects) == len(filter)
}
