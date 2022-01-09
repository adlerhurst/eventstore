package memory

import (
	"context"
	"sync"

	"github.com/adlerhurst/eventstore"
)

type Storage struct {
	sequence uint64
	mu       sync.Mutex
	root     *subject
}

func New() *Storage {
	return &Storage{
		root: &subject{},
	}
}

//Health checks if the storage is available
func (s *Storage) Ready(context.Context) error { return nil }

// Push stores the command's and returns the resulting Event's
// the command's should be stored in a single transaction
func (s *Storage) Push(_ context.Context, cmds []eventstore.Command) (storedEvents []eventstore.EventBase, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	toPush := make(events, len(cmds))
	storedEvents = make([]eventstore.EventBase, len(cmds))
	seq := s.sequence
	for i, cmd := range cmds {
		seq++
		toPush[i], err = NewEvent(cmd, seq)
		if err != nil {
			return nil, err
		}
		storedEvents[i] = toPush[i].toEventstore()
	}

	for _, e := range toPush {
		s.root.push(e.Subjects, e)
	}

	s.sequence = seq
	return storedEvents, nil
}

//Filter returns the events matching the subject
func (s *Storage) Filter(_ context.Context, filter eventstore.Filter) (res []eventstore.EventBase, _ error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	//needed because in the current implementation the base subject has no subject
	filter.Subjects = append([]eventstore.Subject{eventstore.SingleToken}, filter.Subjects...)

	res = make([]eventstore.EventBase, 0, filter.Limit)

	found := s.root.find(filter.Subjects)
	for _, e := range found {

		if filter.From > 0 && e.Sequence < filter.From {
			continue
		}
		if filter.To > 0 && e.Sequence > filter.To {
			break
		}

		res = append(res, e.toEventstore())

		if uint64(len(res)) == filter.Limit {
			break
		}
	}

	return res, nil
}
