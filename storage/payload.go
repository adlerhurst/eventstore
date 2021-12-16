package storage

import (
	"encoding/json"
	"reflect"
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
	typ := reflect.TypeOf(payload)

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	if typ.Kind() == reflect.Struct {
		bytes, err := json.Marshal(payload)
		if err != nil {
			return nil, ErrInvalidPayload
		}
		return bytes, nil
	}

	return nil, ErrInvalidPayload
}
