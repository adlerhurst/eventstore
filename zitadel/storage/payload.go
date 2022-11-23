package storage

import "database/sql/driver"

// Payload represents a byte array that may be null.
// Payload implements the sql.Scanner interface
type Payload []byte

// Scan implements the Scanner interface.
func (p *Payload) Scan(value interface{}) error {
	if value == nil {
		*p = nil
		return nil
	}
	*p = Payload(value.([]byte))
	return nil
}

// Value implements the driver Valuer interface.
func (p Payload) Value() (driver.Value, error) {
	if len(p) == 0 {
		return nil, nil
	}
	return []byte(p), nil
}
