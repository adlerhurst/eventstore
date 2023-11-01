package api

import (
	"github.com/adlerhurst/eventstore/service/internal/api/eventstore/v1alpha"
	"github.com/adlerhurst/eventstore/v2"
)

var (
	_ eventstore.Aggregate = (*Aggregate)(nil)
)

// Aggregate implements eventstore.Aggregate
type Aggregate struct {
	id       eventstore.TextSubjects
	commands []eventstore.Command
	sequence *uint32
}

// CurrentSequence implements eventstore.Aggregate.
func (a *Aggregate) CurrentSequence() *uint32 {
	return a.sequence
}

// Commands implements eventstore.Aggregate.
func (a *Aggregate) Commands() []eventstore.Command {
	return a.commands
}

// ID implements eventstore.Aggregate.
func (a *Aggregate) ID() eventstore.TextSubjects {
	return a.id
}

func protoToAggregate(aggregate *eventstorev1alpha.Aggregate) eventstore.Aggregate {
	return &Aggregate{
		id:       toTextSubjects(aggregate.Id),
		commands: protoToCommands(aggregate.Commands),
	}
}

func pushRequestToAggregates(req *eventstorev1alpha.PushRequest) []eventstore.Aggregate {
	aggregates := make([]eventstore.Aggregate, len(req.Aggregates))

	for i, aggregate := range req.Aggregates {
		aggregates[i] = protoToAggregate(aggregate)
	}

	return aggregates
}
