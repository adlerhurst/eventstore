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
	// Push stores the commands and sets the resulting metadata on the command
	// the commands should be stored in a single transaction
	// if the current sequence of an [AggregatePredefinedSequence] does not match
	// [ErrSequenceNotMatched] is returned
	Push(ctx context.Context, aggregates ...Aggregate) error
	// Filter applies the events matching the subjects on the reducer
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
	// CurrentSequence returns the current sequence of the aggregate
	// If the aggregate doesn't care about the current sequence it returns nil
	// If it's the first command return 0
	// If it's the nth command return the specific sequence
	CurrentSequence() *uint32
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

	SetSequence(sequence uint32)
	SetCreationDate(creationDate time.Time)
}

// Event is the abstraction if a user wants to get events mapped by the eventstore
type Event interface {
	Action
	// Aggregate represents the object the command belongs to
	// and is used to generate the `Sequence` of the [Event]
	// e.g. user A: {"users", "A"}
	Aggregate() TextSubjects
	// Sequence represents the position of the event inside a specific subject
	Sequence() uint32
	// CreationDate is the timestamp the event was stored to the eventstore
	CreationDate() time.Time
	// UnmarshalPayload maps the stored payload into the given object
	// object must be of type *struct
	UnmarshalPayload(object any) error
}

// Filter represents a query
type Filter struct {
	// Queries are queries on subjects
	Queries []*FilterQuery
	// Limit represents the maximum events returned
	Limit uint64
}

type FilterQuery struct {
	// Sequence limits the sequences for this query
	Sequence SequenceFilter
	// CreatedAt filters the time and event was created
	CreatedAt CreatedAtFilter
	// Action represents the event type
	Subjects []Subject
}

type SequenceFilter struct {
	From uint32
	To   uint32
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
