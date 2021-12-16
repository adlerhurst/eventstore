package memory

import (
	"context"
	"sync"

	"github.com/adlerhurst/eventstore"
)

type Storage struct {
	sequence uint64
	mu       sync.Mutex
	base     *subject
}

func New() *Storage {
	return &Storage{
		base: &subject{},
	}
}

//Health checks if the storage is available
func (s *Storage) Ready(context.Context) error { return nil }

//Push stores the command's and returns the resulting Event's
// the command's should be stored in a single transaction
func (s *Storage) Push(_ context.Context, cmds []eventstore.Command) (events []*eventstore.Event, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	events = make([]*eventstore.Event, len(cmds))

	for i, cmd := range cmds {
		s.sequence += 1
		events[i], err = s.base.push(cmd.Subjects(), cmd, s.sequence)
	}

	return events, nil
}

//Filter returns the events matching the subject
func (s *Storage) Filter(_ context.Context, filter eventstore.Filter) ([]*eventstore.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	filter.Subjects = append([]eventstore.Subject{eventstore.SingleToken}, filter.Subjects...)
	return s.base.find(filter.Subjects), nil
}
