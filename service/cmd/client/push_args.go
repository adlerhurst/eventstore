package client

import (
	"encoding/json"
	"os"
	"strings"

	eventstorev1alpha "github.com/adlerhurst/eventstore/service/internal/api/eventstore/v1alpha"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var pushRequest *eventstorev1alpha.PushRequest

func parsePushRequest(cmd *cobra.Command, args []string) {
	pushRequest = new(eventstorev1alpha.PushRequest)

	aggregatesIndexes := listArgIndexes("aggregates", args...)

	for i := 1; i < len(aggregatesIndexes); i++ {
		aggregate, err := parseAggregate(args[aggregatesIndexes[i-1]+1 : aggregatesIndexes[i]])
		if err != nil {
			config.Logger.Error("unable to parse aggregate", "index", i, "cause", err)
			os.Exit(1)
		}
		pushRequest.Aggregates = append(pushRequest.Aggregates, aggregate)
	}

	aggregate, err := parseAggregate(args[aggregatesIndexes[len(aggregatesIndexes)-1]+1:])
	if err != nil {
		config.Logger.Error("unable to parse last aggregate", "cause", err)
		os.Exit(1)
	}
	pushRequest.Aggregates = append(pushRequest.Aggregates, aggregate)
}

func parseAggregate(args []string) (aggregate *eventstorev1alpha.Aggregate, err error) {
	aggregate = new(eventstorev1alpha.Aggregate)

	aggregateFlags := pflag.NewFlagSet("aggregate", pflag.ContinueOnError)
	id := aggregateFlags.StringSlice("id", nil, "id of the aggregate")
	currentSequence := aggregateFlags.Uint32("currentsequence", 0, "current sequence of the aggregate")

	commandsIndexes := listArgIndexes("commands", args...)
	if len(commandsIndexes) == 0 {
		config.Logger.Info("skipping aggregate because no commands defined")
		return nil, nil
	}

	if err = aggregateFlags.Parse(args[:commandsIndexes[0]]); err != nil {
		return nil, err
	}

	aggregate.Id = *id
	aggregate.CurrentSequence = currentSequence

	for i := 1; i < len(commandsIndexes); i++ {
		command, err := parseCommand(args[commandsIndexes[i-1]+1 : commandsIndexes[i]])
		if err != nil {
			return nil, err
		}
		aggregate.Commands = append(aggregate.Commands, command)
	}

	command, err := parseCommand(args[commandsIndexes[len(commandsIndexes)-1]+1:])
	if err != nil {
		return nil, err
	}
	aggregate.Commands = append(aggregate.Commands, command)

	return aggregate, nil
}

// func parseList[T any](name string, args []string, flagSet *pflag.FlagSet, parse func(args []string) (*T, error)) (objects []*T, err error) {
// 	commandsIndexes := listArgIndexes(name, args...)
// 	if len(commandsIndexes) == 0 {
// 		config.Logger.Info("skipping because flag not defined in args", "flag", name)
// 		return nil, nil
// 	}

// 	if err = flagSet.Parse(args[:commandsIndexes[0]]); err != nil {
// 		return nil, err
// 	}

// 	// aggregate.Id = *id
// 	// aggregate.CurrentSequence = currentSequence

// 	for i := 1; i < len(commandsIndexes); i++ {
// 		command, err := parseCommand(args[commandsIndexes[i-1]+1 : commandsIndexes[i]])
// 		if err != nil {
// 			return nil, err
// 		}
// 		aggregate.Commands = append(aggregate.Commands, command)
// 	}

// 	command, err := parseCommand(args[commandsIndexes[len(commandsIndexes)-1]+1:])
// 	if err != nil {
// 		return nil, err
// 	}
// 	return nil, nil
// }

func parseCommand(args []string) (command *eventstorev1alpha.Command, err error) {
	command = new(eventstorev1alpha.Command)

	commandFlags := pflag.NewFlagSet("command", pflag.ContinueOnError)
	action := commandFlags.StringSlice("action", nil, "action of the command")
	revision := commandFlags.Uint32("revision", 0, "revision of the command")
	payload := commandFlags.BytesBase64("payload", nil, "payload base64 (RFC 4648) encoded")

	if err = commandFlags.Parse(args); err != nil {
		return nil, err
	}

	command.Action = &eventstorev1alpha.Action{
		Action:   *action,
		Revision: *revision,
	}

	if payload != nil && len(*payload) > 0 {
		if err = json.Unmarshal(*payload, &command.Action.Payload); err != nil {
			return nil, err
		}
	}

	return command, nil
}

func listArgIndexes(name string, args ...string) (indexes []int) {
	for i, arg := range args {
		switch strings.ToLower(arg) {
		case "--" + name, "-" + name:
			indexes = append(indexes, i)
		}
	}

	return indexes
}
