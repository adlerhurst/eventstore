package storage

import (
	_ "embed"
	"encoding/json"
)

func payloadToJSON(payload interface{}) (Payload, error) {
	if payload == nil {
		return nil, nil
	}
	if p, ok := payload.([]byte); ok && json.Valid(p) {
		return p, nil
	}
	return json.Marshal(payload)
}
