package zitadel

import (
	"database/sql/driver"
	"time"

	"github.com/jackc/pgtype"
)

type Filter struct {
	// If required in future: Tx *sql.Tx
	// have to add in `storage.filterToSQL`,
	// check if tx is nil for as of system time
	// else creation_date < `CreationDateLessEqual`

	InstanceID               string //required
	Aggregates               []*AggregateFilter
	OrgIDs                   StringArray //mandatory
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
	Types StringArray //required
	// FUTURE: Payload map[string]any
}

type StringArray []string

// Scan implements the `database/sql.Scanner` interface.
func (s *StringArray) Scan(src any) error {
	array := new(pgtype.TextArray)
	if err := array.Scan(src); err != nil {
		return err
	}
	if err := array.AssignTo(s); err != nil {
		return err
	}
	return nil
}

// Value implements the `database/sql/driver.Valuer` interface.
func (s StringArray) Value() (driver.Value, error) {
	if len(s) == 0 {
		return nil, nil
	}

	array := pgtype.TextArray{}
	if err := array.Set(s); err != nil {
		return nil, err
	}

	return array.Value()
}
