package zitadel

import "context"

type Eventstore struct {
	storage Storage
}

func NewEventstore(stor Storage) *Eventstore {
	return &Eventstore{
		storage: stor,
	}
}

func (es *Eventstore) Push(ctx context.Context, cmds []Command) ([]Event, error) {
	return nil, nil
}

type Storage interface {
	//Health checks if the storage is available
	Ready(context.Context) error
	//Push stores the command's and returns the resulting Event's
	// the command's should be stored in a single transaction.
	Push(context.Context, []Command) ([]*Event, error)
}
