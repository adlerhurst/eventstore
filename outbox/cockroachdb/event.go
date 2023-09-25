package cockroachdb

import (
	"encoding/json"
	"time"

	"github.com/adlerhurst/eventstore/v0"
)

var _ eventstore.Event = (*Event)(nil)

type Event struct {
	action       eventstore.TextSubjects
	aggregate    eventstore.TextSubjects
	revision     uint16
	creationDate time.Time
	position     float64
	sequence     uint32
	payload      []byte
}

// Action implements [eventstore.Event]
func (e *Event) Action() eventstore.TextSubjects {
	return e.action
}

// Aggregate implements [eventstore.Event]
func (e *Event) Aggregate() eventstore.TextSubjects {
	return e.aggregate
}

// Revision implements [eventstore.Event]
func (e *Event) Revision() uint16 {
	return e.revision
}

// CreationDate implements [eventstore.Event]
func (e *Event) CreationDate() time.Time {
	return e.creationDate
}

// Sequence implements [eventstore.Event]
func (e *Event) Sequence() uint64 {
	return uint64(e.sequence)
}

// UnmarshalPayload implements [eventstore.Event]
func (e *Event) UnmarshalPayload(object any) error {
	if len(e.payload) == 0 {
		return nil
	}
	return json.Unmarshal(e.payload, object)
}

func eventsFromAggregates(aggregates []eventstore.Aggregate) (events []*Event, err error) {
	events = make([]*Event, 0, len(aggregates))
	for _, aggregate := range aggregates {
		aggregateEvents, err := eventsFromAggregate(aggregate)
		if err != nil {
			return nil, err
		}
		events = append(events, aggregateEvents...)
	}

	return events, nil
}

func eventsFromAggregate(aggregate eventstore.Aggregate) ([]*Event, error) {
	events := make([]*Event, len(aggregate.Commands()))
	for i, command := range aggregate.Commands() {
		events[i] = &Event{
			aggregate: aggregate.ID(),
			action:    command.Action(),
			revision:  command.Revision(),
		}

		if command.Payload() != nil {
			payload, err := json.Marshal(command.Payload())
			if err != nil {
				return nil, err
			}
			if len(payload) > 0 {
				events[i].payload = payload
			}
		}
	}

	return events, nil
}
