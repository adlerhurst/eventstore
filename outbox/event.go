package outbox

import (
	"encoding/json"
	"log"
	"time"

	"github.com/adlerhurst/eventstore/v0"
)

var _ eventstore.Event = (*Event)(nil)

type Event struct {
	action       []eventstore.TextSubject
	aggregate    []eventstore.TextSubject
	metadata     []byte
	revision     uint16
	creationDate time.Time
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

// Metadata implements [eventstore.Event]
func (e *Event) Metadata() map[string]interface{} {
	if len(e.metadata) == 0 {
		return nil
	}
	metadata := make(map[string]interface{})
	if err := json.Unmarshal(e.metadata, &metadata); err != nil {
		log.Fatalf("unable to unmarshal metadata: %v", err)
	}
	return metadata
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
func (e *Event) UnmarshalPayload(object interface{}) error {
	if len(e.payload) == 0 {
		return nil
	}
	return json.Unmarshal(e.payload, object)
}

func eventsFromCommands(commands []eventstore.Command) (events []*Event, err error) {
	events = make([]*Event, len(commands))
	for i, command := range commands {
		events[i], err = eventFromCommand(command)
		if err != nil {
			return nil, err
		}
	}

	return events, nil
}

func eventFromCommand(command eventstore.Command) (event *Event, err error) {
	event = &Event{
		action:    command.Action(),
		aggregate: command.Aggregate(),
		revision:  command.Revision(),
	}

	if len(command.Metadata()) > 0 {
		if event.metadata, err = json.Marshal(command.Metadata()); err != nil {
			return nil, err
		}
	}
	if command.Payload() != nil {
		payload, err := json.Marshal(command.Payload())
		if err != nil {
			return nil, err
		}
		if len(payload) > 0 {
			event.payload = payload
		}
	}

	return event, nil
}
