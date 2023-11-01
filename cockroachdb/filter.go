package cockroachdb

import (
	"context"
	_ "embed"
	"strconv"
	"strings"

	"github.com/adlerhurst/eventstore/v2"
	"github.com/jackc/pgx/v5"
)

// Filter implements [eventstore.Eventstore]
func (store *CockroachDB) Filter(ctx context.Context, filter *eventstore.Filter, reducer eventstore.Reducer) (err error) {
	builder, args := store.prepareStatement(filter)

	conn, err := store.client.Acquire(ctx)
	if err != nil {
		logger.ErrorContext(ctx, "acquire connection failed", "cause", err)
		return err
	}
	defer conn.Release()

	// The application_name is required to differentiate between filter and push queries
	// filter queries are not relevant for the [filterIgnoreOpenPush] clause
	_, err = conn.Exec(ctx, "SET application_name = $1", store.filterAppName)
	if err != nil {
		logger.ErrorContext(ctx, "set application name failed", "cause", err)
		return err
	}

	tx, err := conn.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadOnly,
	})
	if err != nil {
		logger.ErrorContext(ctx, "create transaction failed", "cause", err)
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

	rows, err := tx.Query(ctx, builder.String(), args...)
	if err != nil {
		logger.ErrorContext(ctx, "filter events failed", "cause", err)
		return err
	}
	defer rows.Close()

	for rows.Next() {
		event := eventPool.Get()
		err = rows.Scan(
			&event.aggregate,
			&event.revision,
			&event.payload,
			&event.sequence,
			&event.creationDate,
			&event.action,
		)
		if err != nil {
			logger.ErrorContext(ctx, "scan of events failed", "cause", err)
			event.payload = nil
			eventPool.Put(event)
			return err
		}

		if err = reducer.Reduce(event); err != nil {
			logger.DebugContext(ctx, "reduce failed", "cause", err)
			event.payload = nil
			eventPool.Put(event)
			return err
		}
		event.payload = nil
		eventPool.Put(event)
	}

	return rows.Err()
}

var (
	filterColumnSelector = "SELECT e.aggregate, e.revision, e.payload, e.sequence, e.created_at, e.action FROM eventstore.events e "
	filterLimit          = " LIMIT $"
	filterIgnoreOpenPush = `e.created_at < (SELECT COALESCE(MIN(start), NOW())::TIMESTAMPTZ FROM crdb_internal.cluster_transactions where application_name = $`
)

func (store *CockroachDB) prepareStatement(filter *eventstore.Filter) (builder strings.Builder, args []any) {
	var index int

	builder.WriteString(filterColumnSelector)

	builder.WriteString(" WHERE ")
	if len(filter.Queries) > 0 {
		args = queriesToClause(&builder, &index, filter.Queries)
	}

	if len(args) >= 0 {
		builder.WriteString(" AND ")
	}
	builder.WriteString(filterIgnoreOpenPush)
	index++
	builder.WriteString(strconv.Itoa(index))
	args = append(args, store.pushAppName)
	builder.WriteString(") ORDER BY e.position, e.in_tx_order")

	if filter.Limit > 0 {
		builder.WriteString(filterLimit)
		index++
		builder.Write([]byte(strconv.Itoa(index)))
		args = append(args, filter.Limit)
	}

	return builder, args
}

func queriesToClause(builder *strings.Builder, index *int, queries []*eventstore.FilterQuery) (args []any) {
	for i, query := range queries {

		args = append(args, queryToClause(builder, index, query)...)

		if i < len(queries)-1 {
			builder.WriteString(" OR ")
		}
	}

	return args
}

var (
	filterSequenceGt  = " AND e.sequence > $"
	filterSequenceLt  = " AND e.sequence < $"
	filterCreatedAtGt = " AND e.created_at > $"
	filterCreatedAtLt = " AND e.created_at < $"
)

func queryToClause(builder *strings.Builder, index *int, query *eventstore.FilterQuery) []any {
	builder.WriteRune('(')

	args := subjectsToClause(builder, index, query.Subjects)

	if query.Sequence.From > 0 {
		builder.WriteString(filterSequenceGt)
		*index++
		builder.WriteString(strconv.Itoa(*index))
		args = append(args, query.Sequence.From)
	}

	if query.Sequence.To > 0 {
		builder.WriteString(filterSequenceLt)
		*index++
		builder.WriteString(strconv.Itoa(*index))
		args = append(args, query.Sequence.To)
	}

	if !query.CreatedAt.From.IsZero() {
		builder.WriteString(filterCreatedAtGt)
		*index++
		builder.WriteString(strconv.Itoa(*index))
		args = append(args, query.CreatedAt.From)
	}

	if !query.CreatedAt.To.IsZero() {
		builder.WriteString(filterCreatedAtLt)
		*index++
		builder.WriteString(strconv.Itoa(*index))
		args = append(args, query.CreatedAt.To)
	}

	builder.WriteRune(')')

	return args
}

var filterActionsCondition = "e.id IN (SELECT a.event FROM eventstore.actions a"

func subjectsToClause(builder *strings.Builder, index *int, subjects []eventstore.Subject) []any {
	if len(subjects) == 0 {
		return nil
	}

	args := make([]any, 0, len(subjects)*2+1)

	// the loop is used to check if at least 1 subject is a text subject
	// if so text subject queries are written to builder
	for _, subject := range subjects {
		if _, ok := subject.(eventstore.TextSubject); !ok {
			continue
		}

		builder.WriteString(filterActionsCondition)
		args = append(args, subjectsToJoins(builder, index, subjects[1:])...)

		builder.WriteString(" WHERE ")
		if textSubject, ok := subjects[0].(eventstore.TextSubject); ok {
			textSubjectClause(builder, index, "a", textSubject)
			args = append(args, textSubject, 0)
		}
		builder.WriteRune(')')

		break
	}

	if len(args) > 0 {
		builder.WriteString(" AND ")
	}
	actionDepthQuery(builder, index, subjects[len(subjects)-1])
	args = append(args, len(subjects))

	return args
}

func subjectsToJoins(builder *strings.Builder, index *int, subjects []eventstore.Subject) []any {
	args := make([]any, 0, len(subjects)*2)
	for depth, subject := range subjects {
		textSubject, ok := subject.(eventstore.TextSubject)
		if !ok {
			continue
		}
		tableAlias := "a" + strconv.Itoa(depth)
		builder.WriteString(" JOIN eventstore.actions ")
		builder.WriteString(tableAlias)
		builder.WriteString(" ON a.event = ")
		builder.WriteString(tableAlias)
		builder.WriteString(".event")
		builder.WriteString(" AND ")
		textSubjectClause(builder, index, tableAlias, textSubject)
		// depth+1 because depth 0 is handled outside of this function
		args = append(args, textSubject, depth+1)
	}

	return args
}

func actionDepthQuery(builder *strings.Builder, index *int, lastSubject eventstore.Subject) {
	builder.WriteString("e.action_depth")
	switch lastSubject {
	case eventstore.MultiToken:
		builder.WriteString(" >= $")
	default:
		builder.WriteString(" = $")
	}
	*index++
	builder.WriteString(strconv.Itoa(*index))
}

func textSubjectClause(builder *strings.Builder, index *int, tableAlias string, subject eventstore.TextSubject) {
	builder.WriteString(tableAlias)
	builder.WriteString(".action = $")
	*index++
	builder.WriteString(strconv.Itoa(*index))
	builder.WriteString(" AND ")
	builder.WriteString(tableAlias)
	builder.WriteString(".depth = $")
	*index++
	builder.WriteString(strconv.Itoa(*index))
}
