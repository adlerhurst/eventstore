package api

import (
	"github.com/adlerhurst/eventstore"
	"github.com/adlerhurst/eventstore/service/internal/api/eventstore/v1alpha"
)

var (
	_ eventstore.Aggregate                   = (*AggregateWithoutSequence)(nil)
	_ eventstore.AggregatePredefinedSequence = (*AggregateWithSequence)(nil)
)

// AggregateWithoutSequence implements eventstore.Aggregate
type AggregateWithoutSequence struct {
	id       eventstore.TextSubjects
	commands []eventstore.Command
}

// AggregateWithSequence implements eventstore.AggregatePredefinedSequence
type AggregateWithSequence struct {
	AggregateWithoutSequence
	sequence uint32
}

// CurrentSequence implements eventstore.AggregatePredefinedSequence.
func (a *AggregateWithSequence) CurrentSequence() uint32 {
	return a.sequence
}

// Commands implements eventstore.Aggregate.
func (a *AggregateWithoutSequence) Commands() []eventstore.Command {
	return a.commands
}

// ID implements eventstore.Aggregate.
func (a *AggregateWithoutSequence) ID() eventstore.TextSubjects {
	return a.id
}

func protoToAggregate(aggregate *eventstorev1alpha.Aggregate) eventstore.Aggregate {
	if aggregate.CurrentSequence != nil {
		return &AggregateWithSequence{
			AggregateWithoutSequence: AggregateWithoutSequence{
				id:       toTextSubjects(aggregate.Id),
				commands: protoToCommands(aggregate.Commands),
			},
			sequence: *aggregate.CurrentSequence,
		}
	}
	return &AggregateWithoutSequence{
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
