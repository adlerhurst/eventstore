package cockroachdb

import (
	"context"
	"encoding/json"

	"github.com/adlerhurst/eventstore"
)

type command struct {
	eventstore.Command
	payload   []byte
	aggregate eventstore.TextSubjects

	id       string
	sequence uint32
}

func commandsFromAggregates(ctx context.Context, aggregates []eventstore.Aggregate) (commands []*command, close func(), err error) {
	commands = make([]*command, 0, len(aggregates))
	for _, aggregate := range aggregates {
		aggregateEvents, err := commandsFromAggregate(ctx, aggregate)
		if err != nil {
			return nil, func() {}, err
		}
		commands = append(commands, aggregateEvents...)
	}

	return commands,
		func() {
			for _, cmd := range commands {
				cmd.payload = nil
				commandPool.Put(cmd)
			}
		},
		nil
}

func commandsFromAggregate(ctx context.Context, aggregate eventstore.Aggregate) ([]*command, error) {
	commands := make([]*command, len(aggregate.Commands()))
	for i, command := range aggregate.Commands() {
		commands[i] = commandPool.Get()

		commands[i].Command = command
		commands[i].aggregate = aggregate.ID()

		if payload := command.Payload(); payload != nil {
			var err error
			commands[i].payload, err = json.Marshal(payload)
			if err != nil {
				logger.ErrorContext(ctx, "marshal payload failed", "cause", err, "action", commands[i].Action().Join("."))
				return nil, err
			}
		}
	}

	return commands, nil
}
