package eventstore

import (
	"context"
	"errors"
	"time"
)

// Eventstore abstracts all functions needed to store events
// and filters the stored events
type Eventstore interface {
	// Ready checks if the storage is available
	Ready(ctx context.Context) error
	// Push stores the command's and returns the resulting Event's
	// the commands should be stored in a single transaction
	// if the current sequence of an [AggregatePredefinedSequence] does not match
	// [ErrSequenceNotMatched] is returned
	Push(ctx context.Context, aggregates ...Aggregate) ([]Event, error)
	// Filter returns the events matching the subject
	Filter(ctx context.Context, filter *Filter, reducer Reducer) error
}

// Aggregate represents the stream the events are written to
// If the aggregate implements [AggregatePredefinedSequence],
// the current sequence of the aggregate is verified from the storage
type Aggregate interface {
	// ID is the unique identifier of the stream
	ID() TextSubjects
	// Commands is the list of write intents
	Commands() []Command
}

// AggregatePredefinedSequence is used in storage to determine if the command requires a specific sequence
// If the order doesn't matter the command must not implement this interface
type AggregatePredefinedSequence interface {
	Command
	// CurrentSequence returns the current sequence of the aggregate
	// If it's the first command return 0
	// If it's the nth command return the specific sequence
	CurrentSequence() uint32
}

// Action describes the base data of [Command]'s and [Event]'s
type Action interface {
	// Action represent the change of an object
	//
	// most likely the [Aggregate()] list will be the first elements of the
	// [Action]
	// e.g. add user A was added: {"users", "A", "added"}
	Action() TextSubjects
	// Revision is an upcounting number which represents the version of the schema of the payload
	// the revision must change as soon as the logic to create the payload or schema of the payload changes
	Revision() uint16
}

// Command represents a change to be made
type Command interface {
	Action
	// Payload returns the payload of the event. It represent the changed fields by the event
	// valid types are:
	// - nil (no payload),
	// - struct which can be marshalled
	// - pointer to struct which can be marshalled
	Payload() any
}

// Event is the abstraction if a user wants to get events mapped by the eventstore
type Event interface {
	Action
	// Aggregate represents the object the command belongs to
	// and is used to generate the `Sequence` of the [Event]
	// e.g. user A: {"users", "A"}
	Aggregate() TextSubjects
	// Sequence represents the position of the event inside a specific subject
	Sequence() uint64
	// CreationDate is the timestamp the event was stored to the eventstore
	CreationDate() time.Time
	// UnmarshalPayload maps the stored payload into the given object
	// object must be of type *struct
	UnmarshalPayload(object any) error
}

// Filter represents a query
type Filter struct {
	// Sequence filters the sequences of all the actions
	Sequence SequenceFilter
	// CreatedAt filters the time and event was created
	CreatedAt CreatedAtFilter
	// Limit represents the maximum events returned
	Limit uint64
	// Action represents the event type
	Action []Subject
}

type SequenceFilter struct {
	From uint64
	To   uint64
}

type CreatedAtFilter struct {
	From time.Time
	To   time.Time
}

// Reducer represents a model
type Reducer interface {
	// Reduce maps events to a model
	Reduce(events ...Event) error
}

var (
	ErrSequenceNotMatched = errors.New("sequence of aggregate did not match")
)
