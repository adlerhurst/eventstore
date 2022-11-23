package storage

import (
	"context"
	"database/sql"
	_ "embed"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/adlerhurst/eventstore/zitadel"
)

var _ zitadel.Storage = (*CRDB)(nil)

type CRDB struct {
	client *sql.DB
}

var (
	//go:embed create.sql
	createTableStmt string
	//go:embed push.sql
	pushStmtFmt string
)

// NewCRDB creates a new client and checks if all requirements are fulfilled.
func NewCRDB(client *sql.DB) (*CRDB, error) {
	if _, err := client.Exec(createTableStmt); err != nil {
		return nil, err
	}

	return &CRDB{client}, nil
}

func (crdb *CRDB) Ready(ctx context.Context) error {
	return crdb.client.PingContext(ctx)
}

func (crdb *CRDB) Push(ctx context.Context, cmds []zitadel.Command) ([]*zitadel.Event, error) {
	rows, err := crdb.execPush(ctx, cmds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return eventsFromRows(cmds, rows), nil
}

func (crdb *CRDB) execPush(ctx context.Context, cmds []zitadel.Command) (_ *sql.Rows, err error) {
	args := make([]interface{}, 0, len(cmds)*9)
	placeholders := make([]string, len(cmds))
	for i, cmd := range cmds {
		var payload Payload
		if cmds[i].Payload() != nil {
			payload, err = json.Marshal(cmds[i].Payload())
			if err != nil {
				return nil, err
			}
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
		)
		placeholders[i] = "(" +
			strings.Join(
				[]string{
					"$" + strconv.Itoa(1*(i+1)),
					"$" + strconv.Itoa(2*(i+1)),
					"$" + strconv.Itoa(3*(i+1)),
					"$" + strconv.Itoa(4*(i+1)),
					"$" + strconv.Itoa(5*(i+1)),
					"$" + strconv.Itoa(6*(i+1)),
					"$" + strconv.Itoa(7*(i+1)),
					"$" + strconv.Itoa(8*(i+1)),
					"$" + strconv.Itoa(9*(i+1)),
				},
				", ") +
			")"
	}

	return crdb.client.QueryContext(ctx,
		fmt.Sprintf(pushStmtFmt,
			strings.Join(placeholders, ", "),
		),
		args...)
}

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
