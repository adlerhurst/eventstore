package memory

import (
	"context"
	"encoding/json"
	"time"

	"github.com/adlerhurst/eventstore/v0"
)

// New is the constructor
func New() *Storage {
	return &Storage{
		events:            make([]*Event, 0),
		aggregateSequence: make(map[string]uint64),
	}
}

var _ eventstore.Eventstore = (*Storage)(nil)

// Push implements [eventstore.Eventstore]
type Storage struct {
	events            []*Event
	aggregateSequence map[string]uint64
}

// Ready implements [eventstore.Eventstore]
func (s *Storage) Ready(context.Context) error { return nil }

// Push implements [eventstore.Eventstore]
func (s *Storage) Push(ctx context.Context, commands ...eventstore.Command) (_ []eventstore.Event, err error) {
	events := make([]*Event, len(commands))
	response := make([]eventstore.Event, len(commands))
	for i, command := range commands {
		var payload []byte
		if command.Payload() != nil {
			payload, err = json.Marshal(command.Payload())
			if err != nil {
				return nil, err
			}
		}
		aggregate := command.Aggregate().Join(".")
		if _, ok := s.aggregateSequence[aggregate]; !ok {
			s.aggregateSequence[aggregate] = 0
		}
		s.aggregateSequence[aggregate]++

		events[i] = &Event{
			Command:      command,
			sequence:     s.aggregateSequence[aggregate],
			creationDate: time.Now(),
			payload:      payload,
		}

		response[i] = events[i]
	}

	s.events = append(s.events, events...)

	return response, nil
}

// Filter implements [eventstore.Eventstore]
func (s *Storage) Filter(ctx context.Context, filter *eventstore.Filter) ([]eventstore.Event, error) {
	events := make([]eventstore.Event, 0, filter.Limit)

	for _, event := range s.events {
		if filter.Limit > 0 && len(events) == int(filter.Limit) {
			break
		}

		if !checkFrom(filter.From, event) ||
			!checkTo(filter.To, event) ||
			!checkAction(filter.Action, event) {
			continue
		}

		events = append(events, event)
	}

	return events, nil
}

func checkFrom(from time.Time, event *Event) bool {
	if from.IsZero() {
		return true
	}
	return event.creationDate.After(from)
}

func checkTo(to time.Time, event *Event) bool {
	if to.IsZero() {
		return true
	}
	return to.After(event.creationDate)
}

func checkAction(action []eventstore.Subject, event *Event) bool {
	// if len(action) != len(event.Action()) {
	// 	return action[len(action)-1] == eventstore.MultiToken
	// }
	for i, a := range action {
		switch a {
		case eventstore.SingleToken:
			if i > len(event.Action())-1 {
				return false
			}
			continue
		case eventstore.MultiToken:
			if i > len(event.Action())-1 {
				return false
			}
			return true
		default:
			if i > len(event.Action())-1 {
				return false
			}
			// must have type [eventstore.TextSubject] and therefore match
			if a.(eventstore.TextSubject) != event.Action()[i] {
				return false
			}
		}
	}

	return true
}
