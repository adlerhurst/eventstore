package storage

import (
	"context"
	"database/sql"
	_ "embed"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/adlerhurst/eventstore/zitadel"
)

var _ zitadel.Storage = (*CRDB2)(nil)

type CRDB2 struct {
	client *sql.DB
}

var (
	//go:embed 2_push.sql
	pushStmt2Fmt string
	//go:embed 2_create.sql
	createStmt2 string
)

// NewCRDB2 creates a new client and checks if all requirements are fulfilled.
func NewCRDB2(client *sql.DB) (*CRDB2, error) {
	if _, err := client.Exec(createStmt2); err != nil {
		return nil, err
	}

	return &CRDB2{client}, nil
}

func (crdb *CRDB2) Ready(ctx context.Context) error {
	return crdb.client.PingContext(ctx)
}

func (crdb *CRDB2) Push(ctx context.Context, cmds []zitadel.Command) ([]*zitadel.Event, error) {
	rows, err := crdb.execPush(ctx, cmds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return eventsFromRows2(cmds, rows), nil
}

func (crdb *CRDB2) execPush(ctx context.Context, cmds []zitadel.Command) (_ *sql.Rows, err error) {
	args := make([]interface{}, 0, len(cmds)*9)
	placeholders := make([]string, len(cmds))

	tx, err := crdb.client.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	for i, cmd := range cmds {
		var creationDate time.Time
		if err := tx.QueryRowContext(ctx, "SELECT statement_timestamp()").Scan(&creationDate); err != nil {
			tx.Rollback()
			return nil, err
		}

		payload, err := payloadToJSON(cmds[i].Payload())
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		args = append(args,
			cmd.Type(),
			cmd.Aggregate().Type,
			cmd.Aggregate().ID,
			cmd.Aggregate().Version,
			payload,
			cmd.EditorUser(),
			cmd.EditorService(),
			cmd.Aggregate().ResourceOwner,
			cmd.Aggregate().InstanceID,
			creationDate,
		)
		placeholders[i] = "(" +
			strings.Join(
				[]string{
					"$" + strconv.Itoa(i*10+1),
					"$" + strconv.Itoa(i*10+2),
					"$" + strconv.Itoa(i*10+3),
					"$" + strconv.Itoa(i*10+4),
					"$" + strconv.Itoa(i*10+5),
					"$" + strconv.Itoa(i*10+6),
					"$" + strconv.Itoa(i*10+7),
					"$" + strconv.Itoa(i*10+8),
					"$" + strconv.Itoa(i*10+9),
					"$" + strconv.Itoa(i*10+10),
				},
				", ") +
			")"
	}

	rows, err := crdb.client.QueryContext(ctx,
		fmt.Sprintf(pushStmt2Fmt,
			strings.Join(placeholders, ", "),
		),
		args...)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return rows, err
}

func eventsFromRows2(cmds []zitadel.Command, rows *sql.Rows) []*zitadel.Event {
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
