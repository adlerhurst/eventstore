package cockroachdb

import (
	"context"
	"encoding/json"
	"time"

	"github.com/adlerhurst/eventstore/v0"
)

var _ eventstore.Event = (*event)(nil)

type event struct {
	action       eventstore.TextSubjects
	aggregate    eventstore.TextSubjects
	revision     uint16
	creationDate time.Time
	position     float64
	sequence     uint32
	payload      []byte
}

// Action implements [eventstore.Event]
func (e *event) Action() eventstore.TextSubjects {
	return e.action
}

// Aggregate implements [eventstore.Event]
func (e *event) Aggregate() eventstore.TextSubjects {
	return e.aggregate
}

// Revision implements [eventstore.Event]
func (e *event) Revision() uint16 {
	return e.revision
}

// CreationDate implements [eventstore.Event]
func (e *event) CreationDate() time.Time {
	return e.creationDate
}

// Sequence implements [eventstore.Event]
func (e *event) Sequence() uint64 {
	return uint64(e.sequence)
}

// UnmarshalPayload implements [eventstore.Event]
func (e *event) UnmarshalPayload(object any) error {
	if len(e.payload) == 0 {
		return nil
	}
	return json.Unmarshal(e.payload, object)
}

func eventsFromAggregates(ctx context.Context, aggregates []eventstore.Aggregate) (events []*event, close func(), err error) {
	events = make([]*event, 0, len(aggregates))
	for _, aggregate := range aggregates {
		aggregateEvents, err := eventsFromAggregate(ctx, aggregate)
		if err != nil {
			return nil, func() {}, err
		}
		events = append(events, aggregateEvents...)
	}

	return events,
		func() {
			for _, e := range events {
				e.payload = nil
				eventPool.Put(e)
			}
		},
		nil
}

func eventsFromAggregate(ctx context.Context, aggregate eventstore.Aggregate) ([]*event, error) {
	events := make([]*event, len(aggregate.Commands()))
	for i, command := range aggregate.Commands() {
		events[i] = eventPool.Get()

		events[i].aggregate = aggregate.ID()
		events[i].action = command.Action()
		events[i].revision = command.Revision()

		if command.Payload() != nil {
			payload, err := json.Marshal(command.Payload())
			if err != nil {
				logger.ErrorContext(ctx, "marshal payload failed", "cause", err, "action", events[i].action.Join("."))
				return nil, err
			}
			if len(payload) > 0 {
				events[i].payload = payload
			}
		}
	}

	return events, nil
}
