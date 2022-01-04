package memory

import (
	"github.com/adlerhurst/eventstore"
	"github.com/adlerhurst/eventstore/storage"
)

type events []*event

// event implements a linked list
// the aim is to iterate throug events based on the sequence
type event struct {
	eventstore.Event
}

func NewEvent(cmd eventstore.Command, seq uint64) (*event, error) {
	payload, err := storage.PayloadToBytes(cmd.Payload())
	if err != nil {
		return nil, err
	}
	return &event{
		Event: eventstore.Event{
			EditorService: cmd.EditorService(),
			EditorUser:    cmd.EditorUser(),
			Subjects:      cmd.Subjects(),
			Payload:       payload,
			Sequence:      seq,
			ResourceOwner: cmd.ResourceOwner(),
		},
	}, nil
}

func (e *event) toEventstore() eventstore.Event {
	return eventstore.Event{
		EditorService: e.EditorService,
		EditorUser:    e.EditorUser,
		ResourceOwner: e.ResourceOwner,
		Payload:       e.Payload,
		Subjects:      e.Subjects,
		Sequence:      e.Sequence,
	}
}

func (a events) Len() int { return len(a) }

func (a events) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func (a events) Less(i, j int) bool { return a[i].Sequence < a[j].Sequence }
