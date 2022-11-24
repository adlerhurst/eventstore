package zitadel

import (
	"time"
)

type Filter struct {
	// If required in future: Tx *sql.Tx
	// have to add in `storage.filterToSQL`,
	// check if tx is nil for as of system time
	// else creation_date < `CreationDateLessEqual`

	InstanceID               string //required
	Aggregates               []*AggregateFilter
	OrgIDs                   []string //mandatory
	CreationDateGreaterEqual time.Time
	CreationDateLess         time.Time
	Limit                    uint32
	Desc                     bool
	//If required in future: EventIDs []string
}

type AggregateFilter struct {
	Type   string //required
	ID     string //mandatory
	Events []*EventFilter
}

type EventFilter struct {
	Types []string //required
	// FUTURE: Payload map[string]any
}
