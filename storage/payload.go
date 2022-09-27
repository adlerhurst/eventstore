package storage

import (
	"encoding/json"
)

func PayloadToBytes(payload interface{}) ([]byte, error) {
	switch p := payload.(type) {
	case nil:
		return nil, nil
	case []byte:
		if !json.Valid(p) {
			return nil, ErrInvalidPayload
		}
		return p, nil
	}

	bytes, err := json.Marshal(payload)
	if err != nil {
		return nil, ErrInvalidPayload
	}
	return bytes, nil
}
