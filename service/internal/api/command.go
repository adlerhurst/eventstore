package api

import (
	"time"

	"github.com/adlerhurst/eventstore/service/internal/api/eventstore/v1alpha"
	"github.com/adlerhurst/eventstore/v2"
)

var _ eventstore.Command = (*Command)(nil)

type Command struct {
	*eventstorev1alpha.Command

	createdAt time.Time
	sequence  uint32
}

// Action implements eventstore.Command.
func (c *Command) Action() eventstore.TextSubjects {
	return toTextSubjects(c.Command.Action.Action)
}

// Payload implements eventstore.Command.
func (c *Command) Payload() any {
	if payload := c.Command.GetAction().GetPayload(); payload != nil {
		return payload
	}
	return nil
}

// Revision implements eventstore.Command.
func (c *Command) Revision() uint16 {
	return uint16(c.Command.Action.Revision)
}

// SetCreationDate implements eventstore.Command.
func (c *Command) SetCreationDate(creationDate time.Time) {
	c.createdAt = creationDate
}

// SetSequence implements eventstore.Command.
func (c *Command) SetSequence(sequence uint32) {
	c.sequence = sequence
}

func protoToCommands(commands []*eventstorev1alpha.Command) []eventstore.Command {
	cmds := make([]eventstore.Command, len(commands))

	for i, command := range commands {
		cmds[i] = protoToCommand(command)
	}

	return cmds
}

func protoToCommand(command *eventstorev1alpha.Command) *Command {
	return &Command{
		Command: command,
	}
}
