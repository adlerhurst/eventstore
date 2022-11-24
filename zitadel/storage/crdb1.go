package storage

import (
	"context"
	"database/sql"
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/adlerhurst/eventstore/zitadel"
)

var _ zitadel.Storage = (*CRDB1)(nil)

type CRDB1 struct {
	client *sql.DB
}

var (
	//go:embed 1_push.sql
	pushStmt1 string
	//go:embed 1_create.sql
	createStmt1 string
)

// NewCRDB1 creates a new client and checks if all requirements are fulfilled.
func NewCRDB1(client *sql.DB) (*CRDB1, error) {
	if _, err := client.Exec(createStmt1); err != nil {
		return nil, err
	}

	return &CRDB1{client}, nil
}

func (crdb *CRDB1) Ready(ctx context.Context) error {
	return crdb.client.PingContext(ctx)
}

func (crdb *CRDB1) Push(ctx context.Context, cmds []zitadel.Command) (_ []*zitadel.Event, err error) {
	events := make([]*zitadel.Event, len(cmds))
	tx, err := crdb.client.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	for i, cmd := range cmds {
		payload, err := payloadToJSON(cmds[i].Payload())
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		row := tx.QueryRow(pushStmt1,
			cmd.Type(),
			cmd.Aggregate().Type,
			cmd.Aggregate().ID,
			cmd.Aggregate().Version,
			payload,
			cmd.EditorUser(),
			cmd.EditorService(),
			cmd.Aggregate().ResourceOwner,
			cmd.Aggregate().InstanceID,
		)
		events[i] = zitadel.EventFromCommand(cmds[i])
		if cmds[i].Payload() != nil {
			events[i].Payload, err = json.Marshal(cmds[i].Payload())
			if err != nil {
				// this error must never occure because
				// it should happen before push
				tx.Rollback()
				panic(fmt.Sprintf("error occured in marshal after push: %v", err))
			}
		}
		if err = row.Scan(&events[i].CreationDate); err != nil {
			// if this error occures we are fucked
			tx.Rollback()
			panic(fmt.Sprintf("error occured in scan after push: %v", err))
		}
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return events, nil
}
