package storage

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/adlerhurst/eventstore/zitadel"
)

var _ zitadel.Storage = (*CRDB2_2)(nil)

type CRDB2_2 struct {
	client *sql.DB
}

// NewCRDB2_2 creates a new client and checks if all requirements are fulfilled.
func NewCRDB2_2(client *sql.DB) (*CRDB2_2, error) {
	if _, err := client.Exec(createStmt2); err != nil {
		return nil, err
	}

	return &CRDB2_2{client}, nil
}

func (crdb *CRDB2_2) Ready(ctx context.Context) error {
	return crdb.client.PingContext(ctx)
}

func (crdb *CRDB2_2) Push(ctx context.Context, cmds []zitadel.Command) ([]*zitadel.Event, error) {
	rows, err := crdb.execPush(ctx, cmds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return eventsFromRows2(cmds, rows), nil
}

func (crdb *CRDB2_2) Filter(ctx context.Context, filter *zitadel.Filter) ([]*zitadel.Event, error) {
	query := filterStmt2
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
			&event.CreationDate,
			&event.Type,
			&event.Aggregate.Type,
			&event.Aggregate.ID,
			&event.Version,
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

func (crdb *CRDB2_2) execPush(ctx context.Context, cmds []zitadel.Command) (rows *sql.Rows, err error) {
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
			cmd.Version(),
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
					"now() + '" + fmt.Sprintf("%f", time.Duration(time.Microsecond*time.Duration(i)).Seconds()) + "s'",
				},
				", ",
			) +
			")"
	}

	return crdb.client.QueryContext(ctx,
		fmt.Sprintf(pushStmt2Fmt,
			strings.Join(placeholders, ", "),
		),
		args...)
}
