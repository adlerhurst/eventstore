package eventstore

import (
	"context"
	"time"
)

// Eventstore abstracts all functions needed to store events
// and filters the stored events
type Eventstore interface {
	// Health checks if the storage is available
	Ready(context.Context) error
	// Push stores the command's and returns the resulting Event's
	// the command's should be stored in a single transaction
	Push(context.Context, ...Command) ([]Event, error)
	// Filter returns the events matching the subject
	Filter(context.Context, *Filter) ([]Event, error)
}

// Action describes the base data of [Command]'s and [Event]'s
type Action interface {
	// Aggregate represents the object the command belongs to
	// and is used to generate the `Sequence` of the [Event]
	// e.g. user A: {"users", "A"}
	Aggregate() TextSubjects
	// Action represent the change of an object
	//
	// most likely the [Aggregate()] list will be the first elements of the
	// [Action]
	// e.g. add user A was added: {"users", "A", "added"}
	Action() TextSubjects
	// Revision is an upcounting number which represents the version of the schema of the payload
	// the revision must change as soon as the logic to create the payload or schema of the payload changes
	Revision() uint16
	// Metadata are additional data relevant for the event
	// e.g. the service which created the event.
	// The value must be a primitive type
	Metadata() map[string]interface{}
}

// Command represents a change to be made
type Command interface {
	Action
	// Payload returns the payload of the event. It represent the changed fields by the event
	// valid types are:
	// - nil (no payload),
	// - struct which can be marshalled
	// - pointer to struct which can be marshalled
	Payload() interface{}

	// Options allow to configure the behaviour of commands during writes.
	// They are defined by the different eventstore layers
	Options() []func(Command) error
}

// Event is the abstraction if a user wants to get events mapped by the eventstore
type Event interface {
	Action
	// Sequence represents the position of the event inside a specific subject
	Sequence() uint64
	// CreationDate is the timestamp the event was stored to the eventstore
	CreationDate() time.Time
	// UnmarshalPayload maps the stored payload into the given object
	// object must be of type *struct
	UnmarshalPayload(object interface{}) error
}

// Filter represents a query
type Filter struct {
	// From represents the lowest sequence
	From uint64
	// To represents the highest sequence
	To uint64
	// Limit represents the maximum events returned
	Limit uint64
	// Action represents the event type
	Action []Subject
}
