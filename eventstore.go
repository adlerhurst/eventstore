package eventstore

import (
	"context"
	"errors"
)

//Eventstore abstracts all functions needed to store valid events
// and filters the stored events
type Eventstore struct {
	storage Storage
}

func New(s Storage) *Eventstore {
	return &Eventstore{
		storage: s,
	}
}

//Command represents a change to be made
type Command interface {
	// EditorService is the service who wants to push the event
	EditorService() string
	//EditorUser is the user who wants to push the event
	EditorUser() string
	//Type must return an event type which should be unique in the aggregate
	Subjects() []TextSubject
	//ResourceOwner is the owner of the command
	ResourceOwner() string
	//Payload returns the payload of the event. It represent the changed fields by the event
	// valid types are:
	// * nil (no payload),
	// * json byte array
	// * struct which can be marshalled to json
	// * pointer to struct which can be marshalled to json
	Payload() interface{}
}

//Event represents a change
type Event struct {
	//EditorService is the service which pushed the event
	EditorService string
	//EditorUser is the user which pushed the event
	EditorUser string
	//ResourceOwner is the owner of the event
	ResourceOwner string
	//Payload represents the data as json
	Payload []byte

	//Subject represent the object the event belongs to
	//e.g. add user A event: {"users", "A", "added"}
	Subjects []TextSubject

	Sequence uint64
}

var (
	ErrInvalidCommand = errors.New("invalid command")
)

type Storage interface {
	//Health checks if the storage is available
	Ready(context.Context) error
	//Push stores the command's and returns the resulting Event's
	// the command's should be stored in a single transaction
	Push(context.Context, []Command) ([]Event, error)
	//Filter returns the events matching the subject
	Filter(context.Context, Filter) ([]Event, error)
}

//Ready checks if the eventstore can properly work
// It checks if the repository can serve load
func (es *Eventstore) Ready(ctx context.Context) error {
	return es.storage.Ready(ctx)
}

//PushEvents pushes the events in a single transaction
// an event needs at least an aggregate
func (es *Eventstore) Push(ctx context.Context, commands ...Command) ([]Event, error) {
	return es.storage.Push(ctx, commands)
}

type Filter struct {
	//From represents the lowest sequence
	From uint64
	//To represents the highest sequence
	To uint64
	//Limit represents the maximum events returned
	Limit uint64
	//Subjects represent the subjects to look
	Subjects []Subject
}

func (es *Eventstore) Filter(ctx context.Context, f Filter) ([]Event, error) {
	return es.storage.Filter(ctx, f)
}
