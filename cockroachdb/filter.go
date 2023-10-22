package cockroachdb

import (
	"bytes"
	"context"
	_ "embed"
	"strconv"
	"strings"
	"text/template"

	"github.com/adlerhurst/eventstore/v0"
	"github.com/jackc/pgx/v5"
)

var (
	//go:embed filter.sql
	filterStmt string
	filterTmpl *template.Template
)

// Filter implements [eventstore.Eventstore]
func (crdb *CockroachDB) Filter(ctx context.Context, filter *eventstore.Filter, reducer eventstore.Reducer) (err error) {
	query, args := prepareStatement(filter)

	tx, err := crdb.client.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted, AccessMode: pgx.ReadOnly})
	if err != nil {
		return err
	}
	defer func() {
		// errors are not handled because it's a read-only transaction
		if err != nil {
			_ = tx.Rollback(ctx)
			return
		}
		_ = tx.Commit(ctx)
	}()

	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		event := eventPool.Get()
		err = rows.Scan(
			&event.aggregate,
			&event.action,
			&event.revision,
			&event.payload,
			&event.sequence,
			&event.creationDate,
			&event.position,
		)
		if err != nil {
			event.payload = nil
			eventPool.Put(event)
			return err
		}

		if err = reducer.Reduce(event); err != nil {
			event.payload = nil
			eventPool.Put(event)
			return err
		}
		event.payload = nil
		eventPool.Put(event)
	}

	return nil
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
		clauses = append(clauses, "e.created_at >= ?")
		args = append(args, filter.CreatedAt.From)
	}
	if !filter.CreatedAt.To.IsZero() {
		clauses = append(clauses, "e.created_at <= ?")
		args = append(args, filter.CreatedAt.To)
	}
	if filter.Sequence.From > 0 {
		clauses = append(clauses, "e.sequence >= ?")
		args = append(args, filter.Sequence.From)
	}
	if filter.Sequence.To > 0 {
		clauses = append(clauses, "e.sequence <= ?")
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
			args = append(args, len(action)-1)
			return []string{"(" + strings.Join(clauses, " OR ") + ")", `a."cardinality" >= ?`}, args
		default:
			clauses = append(clauses, `(a."action" = ? AND a."index" = ?)`)
			args = append(args, a, i)
		}
	}

	args = append(args, len(action))

	return []string{"(" + strings.Join(clauses, " OR ") + ")", `a."cardinality" >= ?`}, args
}
