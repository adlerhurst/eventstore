package memory

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"strconv"
	"strings"
	"text/template"

	"github.com/adlerhurst/eventstore/v0"
)

var (
	//go:embed filter.sql
	filterStmt string
	filterTmpl *template.Template
)

// Filter implements [eventstore.Eventstore]
func (crdb *CockroachDB) Filter(ctx context.Context, filter *eventstore.Filter) (events []eventstore.Event, err error) {
	query, args := prepareStatement(filter)

	rows, err := crdb.client.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		event := new(Event)
		err = rows.Scan(
			&event.aggregate,
			&event.action,
			&event.revision,
			&event.metadata,
			&event.payload,
			&event.sequence,
			&event.creationDate,
		)
		if err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	return events, nil
}

type filterData struct {
	Where string
	Limit string
}

func prepareStatement(filter *eventstore.Filter) (string, []interface{}) {
	var data filterData
	clauses, args := filterToClauses(filter)

	if len(clauses) > 0 {
		data.Where = "WHERE " + strings.Join(clauses, " AND ")
	}

	if filter.Limit > 0 {
		data.Limit = "LIMIT ?"
		args = append(args, filter.Limit)
	}

	buf := bytes.NewBuffer(nil)
	if err := filterTmpl.Execute(buf, data); err != nil {
		panic(err)
	}

	query := buf.String()
	for i := range args {
		query = strings.Replace(query, "?", "$"+strconv.Itoa(i+1), 1)
	}

	return query, args
}

func filterToClauses(filter *eventstore.Filter) (clauses []string, args []interface{}) {
	if filter == nil {
		return nil, nil
	}

	clauses = make([]string, 0, 4)
	args = make([]interface{}, 0, 4)

	if !filter.CreatedAt.From.IsZero() {
		clauses = append(clauses, "created_at >= ?")
		args = append(args, filter.CreatedAt.From)
	}
	if !filter.CreatedAt.To.IsZero() {
		clauses = append(clauses, "created_at <= ?")
		args = append(args, filter.CreatedAt.To)
	}
	if filter.Sequence.From > 0 {
		clauses = append(clauses, "sequence >= ?")
		args = append(args, filter.Sequence.From)
	}
	if filter.Sequence.To > 0 {
		clauses = append(clauses, "sequence <= ?")
		args = append(args, filter.Sequence.To)
	}
	if len(filter.Action) > 0 {
		c, a := actionToClauses(filter.Action)
		clauses = append(clauses, c...)
		args = append(args, a...)
	}

	return clauses, args
}

func actionToClauses(action []eventstore.Subject) (clauses []string, args []interface{}) {
	for i, a := range action {
		switch a {
		case eventstore.SingleToken:
			continue
		case eventstore.MultiToken:
			clauses = append(clauses, "cardinality(\"action\") >= ?")
			args = append(args, len(action)-1)
			return clauses, args
		default:
			clauses = append(clauses, fmt.Sprintf("\"action\"[%d] = ?", i+1))
			args = append(args, a)
		}
	}

	clauses = append(clauses, "cardinality(\"action\") = ?")
	args = append(args, len(action))

	return clauses, args
}
