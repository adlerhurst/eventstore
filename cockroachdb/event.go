package cockroachdb

import (
	"encoding/json"
	"time"

	"github.com/adlerhurst/eventstore/v0"
)

var _ eventstore.Event = (*event)(nil)

type event struct {
	action       eventstore.TextSubjects
	aggregate    eventstore.TextSubjects
	revision     uint16
	creationDate time.Time
	position     float64
	sequence     uint32
	payload      []byte
}

// Action implements [eventstore.Event]
func (e *event) Action() eventstore.TextSubjects {
	return e.action
}

// Aggregate implements [eventstore.Event]
func (e *event) Aggregate() eventstore.TextSubjects {
	return e.aggregate
}

// Revision implements [eventstore.Event]
func (e *event) Revision() uint16 {
	return e.revision
}

// CreationDate implements [eventstore.Event]
func (e *event) CreationDate() time.Time {
	return e.creationDate
}

// Sequence implements [eventstore.Event]
func (e *event) Sequence() uint64 {
	return uint64(e.sequence)
}

// UnmarshalPayload implements [eventstore.Event]
func (e *event) UnmarshalPayload(object any) error {
	if len(e.payload) == 0 {
		return nil
	}
	return json.Unmarshal(e.payload, object)
}
