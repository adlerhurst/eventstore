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

var _ zitadel.Storage = (*CRDB3)(nil)

type CRDB3 struct {
	client *sql.DB
}

var (
	//go:embed 3_push.sql
	pushStmt3Fmt string
	//go:embed 3_create.sql
	createStmt3 string
	//go:embed 3_filter.sql
	filterStmt3 string
)

// NewCRDB3 creates a new client and checks if all requirements are fulfilled.
func NewCRDB3(client *sql.DB) (*CRDB3, error) {
	if _, err := client.Exec(createStmt3); err != nil {
		return nil, err
	}

	return &CRDB3{client}, nil
}

func (crdb *CRDB3) Ready(ctx context.Context) error {
	return crdb.client.PingContext(ctx)
}

func (crdb *CRDB3) Push(ctx context.Context, cmds []zitadel.Command) ([]*zitadel.Event, error) {
	rows, err := crdb.execPush(ctx, cmds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return eventsFromRows3(cmds, rows), nil
}

func (crdb *CRDB3) Filter(ctx context.Context, filter *zitadel.Filter) ([]*zitadel.Event, error) {
	query := filterStmt3 + " WHERE "
	clause, args := filterToSQL(filter)
	query += clause

	rows, err := crdb.client.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := make([]*zitadel.Event, 0, filter.Limit)
	for rows.Next() {
		event := new(zitadel.Event)
		var payload Payload
		if err := rows.Scan(
			&event.Type,
			&event.ID,
			&event.CreationDate,
			&event.Type,
			&event.Aggregate.Type,
			&event.Aggregate.ID,
			&event.Aggregate.Version,
			&payload,
			&event.EditorUser,
			&event.Aggregate.ResourceOwner,
			&event.Aggregate.InstanceID,
		); err != nil {
			return nil, err
		}

		event.Payload = payload
		events = append(events, event)
	}

	return events, nil
}

func (crdb *CRDB3) execPush(ctx context.Context, cmds []zitadel.Command) (_ *sql.Rows, err error) {
	args := make([]interface{}, 0, len(cmds)*9)
	placeholders := make([]string, len(cmds))

	for i, cmd := range cmds {
		payload, err := payloadToJSON(cmds[i].Payload())
		if err != nil {
			return nil, err
		}

		args = append(args,
			cmd.Type(),
			cmd.Aggregate().Type,
			cmd.Aggregate().ID,
			cmd.Aggregate().Version,
			payload,
			cmd.EditorUser(),
			cmd.Aggregate().ResourceOwner,
			cmd.Aggregate().InstanceID,
		)

		placeholders[i] = "(" +
			strings.Join(
				[]string{
					"$" + strconv.Itoa(i*8+1),
					"$" + strconv.Itoa(i*8+2),
					"$" + strconv.Itoa(i*8+3),
					"$" + strconv.Itoa(i*8+4),
					"$" + strconv.Itoa(i*8+5),
					"$" + strconv.Itoa(i*8+6),
					"$" + strconv.Itoa(i*8+7),
					"$" + strconv.Itoa(i*8+8),
				},
				", ") +
			")"
	}

	return crdb.client.QueryContext(ctx,
		fmt.Sprintf(pushStmt3Fmt,
			strings.Join(placeholders, ", "),
		),
		args...)
}

func eventsFromRows3(cmds []zitadel.Command, rows *sql.Rows) []*zitadel.Event {
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
		if err = rows.Scan(&events[i].ID, &events[i].CreationDate); err != nil {
			// if this error occures we are fucked
			panic(fmt.Sprintf("error occured in scan after push: %v", err))
		}
	}
	return events
}
