// This implementation is too slow, querying 1732 rows using a single filter takes 88ms
package storage

import (
	"context"
	"database/sql"
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
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

func (crdb *CRDB4) Push(ctx context.Context, cmds []zitadel.Command) (events []*zitadel.Event, err error) {
	args := make([]interface{}, 0, len(cmds)*2)
	placeholders := make([]string, len(cmds))
	events = make([]*zitadel.Event, len(cmds))

	for i, cmd := range cmds {
		sqlEvent, payload, err := cmdToEvent4(cmd)
		if err != nil {
			return nil, err
		}

		args = append(args,
			cmd.Aggregate().ID,
			sqlEvent,
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

		events[i] = zitadel.EventFromCommand(cmd)
		events[i].Payload = payload
	}

	rows, err := crdb.client.QueryContext(ctx,
		fmt.Sprintf(pushStmt4Fmt,
			strings.Join(placeholders, ", "),
		),
		args...)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for i := 0; rows.Next(); i++ {
		err = rows.Scan(&events[i].CreationDate)
		if err != nil {
			return nil, err
		}
	}

	return events, nil
}

func (crdb *CRDB4) Filter(ctx context.Context, filter *zitadel.Filter) ([]*zitadel.Event, error) {
	query := filterStmt4 + " WHERE "
	clause, args := filterToSQL4(filter)
	query += clause

	rows, err := crdb.client.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := make([]*zitadel.Event, 0, filter.Limit)
	for rows.Next() {
		event := new(zitadel.Event)
		var payload []byte
		if err := rows.Scan(
			&event.CreationDate,
			&payload,
		); err != nil {
			return nil, err
		}

		event.Payload = payload
		events = append(events, event)
	}

	return events, nil
}

func cmdToEvent4(cmd zitadel.Command) (event, payloadBytes []byte, err error) {
	payload, payloadBytes, err := payload4FromAny(cmd.Payload())
	if err != nil {
		return nil, nil, err
	}

	event, err = json.Marshal(&Event4{
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
	if err != nil {
		return nil, nil, err
	}

	return event, payloadBytes, nil
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

func payload4FromAny(payload any) (pl Payload4, payloadBytes []byte, err error) {
	if payload == nil {
		return nil, nil, nil
	}

	if p, ok := payload.([]byte); !ok {
		if payloadBytes, err = json.Marshal(payload); err != nil {
			return nil, nil, err
		}
	} else {
		payloadBytes = p
	}

	err = json.Unmarshal(payloadBytes, &pl)
	if err != nil {
		return nil, nil, err
	}
	return pl, payloadBytes, nil
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

	if !filter.CreationDateGreaterEqual.IsZero() {
		clauses = append(clauses, "creation_date >= "+arg(&argCounter))
		args = append(args, filter.CreationDateGreaterEqual)
	}

	if len(filter.Aggregates) > 0 {
		var aggregatesClause string
		aggregatesClause, args = aggregatesFilter(filter.Aggregates, &argCounter, args)
		clauses = append(clauses, aggregatesClause)
	}

	if len(filter.OrgIDs) > 0 {
		var orgClause string
		orgClause, args = ownerFilter(filter.OrgIDs, &argCounter, args)
		clauses = append(clauses, orgClause)
	}

	if filter.InstanceID != "" {
		var instanceClause string
		instanceClause, args = instanceFilter(filter.InstanceID, &argCounter, args)
		clauses = append(clauses, instanceClause)
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
		clauses[i] = eventClause(argCounter)
		args = append(args, aggregateArg("owner", org))
	}

	return strings.Join(clauses, " OR "), args
}

func instanceFilter(instance string, argCounter *int, args []any) (clause string, _ []any) {
	return eventClause(argCounter), append(args, aggregateArg("instance", instance))
}

func aggregatesFilter(aggregates []*zitadel.AggregateFilter, argCounter *int, args []any) (string, []any) {
	clauses := make([]string, len(aggregates))
	for i, filter := range aggregates {
		clauses[i], args = aggregateFilter(filter, argCounter, args)
		clauses[i] = "(" + clauses[i] + ")"
	}

	return strings.Join(clauses, " OR "), args
}

func aggregateFilter(aggregate *zitadel.AggregateFilter, argCounter *int, args []any) (string, []any) {
	clauses := make([]string, 0, 3)

	clauses = append(clauses, eventClause(argCounter))
	args = append(args, aggregateArg("type", aggregate.Type))

	if aggregate.ID != "" {
		clauses = append(clauses, eventClause(argCounter))
		args = append(args, aggregateArg("id", aggregate.ID))
	}

	if len(aggregate.Events) > 0 {
		var eventsClause string
		eventsClause, args = eventsFilter(aggregate.Events, argCounter, args)
		if len(aggregate.Events) > 1 {
			eventsClause = "(" + eventsClause + ")"
		}
		clauses = append(clauses, eventsClause)
	}

	return strings.Join(clauses, " AND "), args
}

func eventsFilter(events []*zitadel.EventFilter, argCounter *int, args []any) (string, []any) {
	clauses := make([]string, len(events))

	for i, event := range events {
		clauses[i], args = eventFilter(event, argCounter, args)
	}

	return strings.Join(clauses, " OR "), args
}

func eventFilter(event *zitadel.EventFilter, argCounter *int, args []any) (string, []any) {
	clauses := make([]string, len(event.Types), len(event.Types)+len(event.Payload))

	for i, typ := range event.Types {
		clauses[i] = eventClause(argCounter)
		args = append(args, eventArg("type", typ))
	}

	if len(event.Payload) > 0 {
		clauses = append(clauses, eventClause(argCounter))
		arg, err := json.Marshal(event.Payload)
		if err != nil {
			log.Fatalf("unable to marshal payload: %v", err)
		}
		args = append(args, arg)
	}

	return strings.Join(clauses, " OR "), args
}

func eventClause(argCounter *int) string {
	return "event @> " + arg(argCounter)
}

func aggregateArg(key, field string) string {
	return `{"aggregate": {"` + key + `": "` + field + `"}}`
}

func eventArg(key, field string) string {
	return `{"` + key + `": "` + field + `"}`
}
