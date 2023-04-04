package memory

import (
	"encoding/json"
	"time"

	"github.com/adlerhurst/eventstore/v0"
)

// Event implements [eventstore.Event]
type Event struct {
	eventstore.Command
	sequence     uint64
	creationDate time.Time
	payload      []byte
}

// Sequence implements [eventstore.Event]
func (e *Event) Sequence() uint64 {
	return e.sequence
}

// CreationDate implements [eventstore.Event]
func (e *Event) CreationDate() time.Time {
	return e.creationDate
}

// UnmarshalPayload implements [eventstore.Event]
func (e *Event) UnmarshalPayload(object interface{}) error {
	if len(e.payload) == 0 {
		return nil
	}
	return json.Unmarshal(e.payload, object)
}
