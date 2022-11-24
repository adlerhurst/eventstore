package storage

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/adlerhurst/eventstore/zitadel"
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

// 2022-11-24 16:00:23.33129+00
var sqlTimeLayout = "2006-01-02 15:04:05.999999-07"

func filterToSQL(filter *zitadel.Filter) (clause string, args []any) {
	var argCounter int

	if !filter.CreationDateLess.IsZero() {
		clause += "AS OF SYSTEM TIME '" + filter.CreationDateLess.Add(1*time.Microsecond).
			Format(sqlTimeLayout) + "' "
	}

	clauses := make([]string, 0, 5)
	args = make([]any, 0, 5)

	clauses = append(clauses, "instance_id = "+arg(&argCounter))
	args = append(args, filter.InstanceID)

	if !filter.CreationDateGreaterEqual.IsZero() {
		clauses = append(clauses, "creation_date >= "+arg(&argCounter))
		args = append(args, filter.CreationDateGreaterEqual)
	}

	if len(filter.OrgIDs) > 0 {
		clauses = append(clauses, "resource_owner = ANY "+arg(&argCounter))
		args = append(args, filter.OrgIDs)
	}

	if len(filter.Aggregates) > 0 {
		var c string
		c, args = aggregateFiltersToSQL(filter.Aggregates, &argCounter, args)
		if len(filter.Aggregates) > 1 {
			c = "(" + c + ")"
		}
		clauses = append(clauses, c)
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

func aggregateFiltersToSQL(filters []*zitadel.AggregateFilter, argCount *int, args []any) (string, []any) {
	clauses := make([]string, len(filters))
	for i, filter := range filters {
		clauses[i], args = aggregateFilterToSQL(filter, argCount, args)
		clauses[i] = "(" + clauses[i] + ")"
	}

	return strings.Join(clauses, " OR "), args
}

func aggregateFilterToSQL(filter *zitadel.AggregateFilter, argCount *int, args []any) (string, []any) {
	clauses := make([]string, 0, 3)

	clauses = append(clauses, "aggregate_type = "+arg(argCount))
	args = append(args, filter.Type)

	if filter.ID != "" {
		clauses = append(clauses, "aggregate_id = "+arg(argCount))
		args = append(args, filter.ID)
	}

	if len(filter.Events) > 0 {
		var eventsClause string
		eventsClause, args = eventFiltersToSQL(filter.Events, argCount, args)
		if len(filter.Events) > 1 {
			eventsClause = "(" + eventsClause + ")"
		}
		clauses = append(clauses, eventsClause)
	}

	return strings.Join(clauses, " AND "), args
}

func eventFiltersToSQL(filters []*zitadel.EventFilter, argCount *int, args []any) (string, []any) {
	clauses := make([]string, len(filters))
	for i, filter := range filters {
		clauses[i], args = eventFilterToSQL(filter, argCount, args)
	}

	return strings.Join(clauses, " OR "), args
}

func eventFilterToSQL(filter *zitadel.EventFilter, argCount *int, args []any) (string, []any) {
	clause := "event_type = ANY " + arg(argCount)
	args = append(args, filter.Types)

	return clause, args
}

func arg(counter *int) string {
	*counter++
	return "$" + strconv.Itoa(*counter)
}

type unimplementedFilter struct{}

var ErrUnimplemented = errors.New("unimplemented")

func (_ unimplementedFilter) Filter(context.Context, *zitadel.Filter) ([]*zitadel.Event, error) {
	return nil, ErrUnimplemented
}
