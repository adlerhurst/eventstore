package storage

import (
	"database/sql"
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/adlerhurst/eventstore/zitadel"
)

var (
	//go:embed create.sql
	createTableStmt string
)

func eventsFromRows(cmds []zitadel.Command, rows *sql.Rows) []*zitadel.Event {
	var err error
	events := make([]*zitadel.Event, len(cmds))

	for i := 0; rows.Next(); i++ {
		events[i] = zitadel.EventFromCommand(cmds[i])
		if cmds[i].Payload() != nil {
			events[i].Payload, err = json.Marshal(cmds[i].Payload())
			if err != nil {
				// this error must never occure because
				// it should happen before push
				panic(fmt.Sprintf("error occured in marshal after push: %v", err))
			}
		}
		if err = rows.Scan(&events[i].CreationDate); err != nil {
			// if this error occures we are fucked
			panic(fmt.Sprintf("error occured in scan after push: %v", err))
		}
	}
	return events
}

func payloadToJSON(payload interface{}) (Payload, error) {
	if payload == nil {
		return nil, nil
	}
	if p, ok := payload.([]byte); ok && json.Valid(p) {
		return p, nil
	}
	return json.Marshal(payload)
}
