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

var (
	_ zitadel.Storage = (*CRDB4)(nil)

	//go:embed 4_create.sql
	createStmt4 string
	//go:embed 4_push.sql
	pushStmt4Fmt string
	//go:embed 4_filter.sql
	filterStmt4 string
)

type CRDB4 struct {
	client *sql.DB
}

// NewCRDB4 creates a new client and checks if all requirements are fulfilled.
func NewCRDB4(client *sql.DB) (*CRDB4, error) {
	if _, err := client.Exec(createStmt4); err != nil {
		return nil, err
	}

	return &CRDB4{client}, nil
}

func (crdb *CRDB4) Ready(ctx context.Context) error {
	return crdb.client.PingContext(ctx)
}

func (crdb *CRDB4) Push(ctx context.Context, cmds []zitadel.Command) ([]*zitadel.Event, error) {
	rows, err := crdb.execPush(ctx, cmds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return eventsFromRows2(cmds, rows), nil
}

func (crdb *CRDB4) Filter(ctx context.Context, filter *zitadel.Filter) ([]*zitadel.Event, error) {
	query := filterStmt2 + " WHERE "
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

func (crdb *CRDB4) execPush(ctx context.Context, cmds []zitadel.Command) (rows *sql.Rows, err error) {
	args := make([]interface{}, 0, len(cmds)*2)
	placeholders := make([]string, len(cmds))

	for i, cmd := range cmds {
		event, err := cmdToEvent4(cmd)
		if err != nil {
			return nil, err
		}

		args = append(args,
			cmd.Aggregate().ID,
			event,
		)

		placeholders[i] = "(" +
			strings.Join(
				[]string{
					"$" + strconv.Itoa(i*2+1),
					"$" + strconv.Itoa(i*2+2),
					"now() + '" + fmt.Sprintf("%f", time.Duration(time.Microsecond*time.Duration(i)).Seconds()) + "s'",
				},
				", ",
			) +
			")"
	}

	return crdb.client.QueryContext(ctx,
		fmt.Sprintf(pushStmt4Fmt,
			strings.Join(placeholders, ", "),
		),
		args...)
}

func cmdToEvent4(cmd zitadel.Command) ([]byte, error) {
	payload, err := payload4FromAny(cmd.Payload())
	if err != nil {
		return nil, err
	}

	return json.Marshal(&Event4{
		Aggregate: &Aggregate4{
			ID:            cmd.Aggregate().ID,
			Type:          cmd.Aggregate().Type,
			ResourceOwner: cmd.Aggregate().ResourceOwner,
			InstanceID:    cmd.Aggregate().InstanceID,
		},
		EditorUser: cmd.EditorUser(),
		Type:       cmd.Type(),
		Payload:    payload,
		Version:    cmd.Version(),
	})
}

type Event4 struct {
	// Aggregate is the metadata of an aggregate
	Aggregate *Aggregate4 `json:"aggregate"`
	// EditorUser is the user who wants to push the event
	EditorUser string `json:"editorUser"`
	// Type must return an event type which should be unique in the aggregate
	Type string `json:"type"`
	// Payload of the event
	Payload map[string]interface{} `json:"payload,omitempty"`
	// Version is the semver this event represents
	Version string `json:"version"`
}

// Aggregate is the basic implementation of Aggregater
type Aggregate4 struct {
	//ID is the unique identitfier of this aggregate
	ID string `json:"id"`
	//Type is the name of the aggregate.
	Type string `json:"type"`
	//ResourceOwner is the org this aggregates belongs to
	ResourceOwner string `json:"owner"`
	//InstanceID is the instance this aggregate belongs to
	InstanceID string `json:"instance"`
}

// Payload4 represents a generic json object that may be null.
// Payload4 implements the sql.Scanner interface
type Payload4 map[string]interface{}

func payload4FromAny(payload any) (pl Payload4, err error) {
	if payload == nil {
		return nil, nil
	}

	if _, ok := payload.([]byte); !ok {
		if payload, err = json.Marshal(payload); err != nil {
			return nil, err
		}
	}

	err = json.Unmarshal(payload.([]byte), &pl)
	if err != nil {
		return nil, err
	}
	return pl, nil
}

// Scan implements the Scanner interface.
func (p *Payload4) Scan(value interface{}) error {
	if value == nil {
		*p = nil
		return nil
	}
	return json.Unmarshal(value.([]byte), p)
}

// // Value implements the driver Valuer interface.
// func (p Payload4) Value() (driver.Value, error) {
// 	if len(p) == 0 {
// 		return nil, nil
// 	}
// 	return []byte(p), nil
// }

func filterToSQL4(filter *zitadel.Filter) (clause string, args []any) {
	if !filter.CreationDateLess.IsZero() {
		clause += "AS OF SYSTEM TIME '" + filter.CreationDateLess.Add(1*time.Microsecond).
			Format(sqlTimeLayout) + "' "
	}

	var argCounter int
	clauses := make([]string, 0, 5)
	args = make([]any, 0, 5)

	if filter.InstanceID != "" {
		var instanceClause string
		instanceClause, args = instanceFilter(filter.InstanceID, &argCounter, args)
		clauses = append(clauses, instanceClause)
	}

	if !filter.CreationDateGreaterEqual.IsZero() {
		clauses = append(clauses, "creation_date >= "+arg(&argCounter))
		args = append(args, filter.CreationDateGreaterEqual)
	}

	if len(filter.OrgIDs) > 0 {
		var orgClause string
		orgClause, args = ownerFilter(filter.OrgIDs, &argCounter, args)
		clauses = append(clauses, orgClause)
	}

	if len(filter.Aggregates) > 0 {
		var aggregateClause string
		aggregateClause, args = aggregatesFilter(filter.Aggregates, &argCounter, args)
		clauses = append(clauses, aggregateClause)
	}

	clause += strings.Join(clauses, " AND ")

	clause += " ORDER BY creation_date"
	if filter.Desc {
		clause += " DESC"
	}

	if filter.Limit > 0 {
		clause += " LIMIT " + arg(&argCounter)
		args = append(args, filter.Limit)
	}

	return clause, args
}

func ownerFilter(orgs []string, argCounter *int, args []any) (string, []any) {
	clauses := make([]string, len(orgs))
	for i, org := range orgs {
		clauses[i] = eventFilter(argCounter)
		args = append(args, aggregateFilter("owner", org))
	}

	return strings.Join(clauses, " OR "), args
}

func instanceFilter(instance string, argCounter *int, args []any) (clause string, _ []any) {
	return eventFilter(argCounter), append(args, aggregateFilter("instance", instance))
}

func eventFilter(argCounter *int) string {
	return "event @> " + arg(argCounter)
}

func aggregateFilter(key, field string) string {
	return `{"aggregate": {"` + key + `": "` + field + `"}}`
}
