package cockroachdb

import (
	"context"
	_ "embed"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/adlerhurst/eventstore/v0"
)

// Push implements [eventstore.Eventstore]
func (crdb *CockroachDB) Push(ctx context.Context, aggregates ...eventstore.Aggregate) (err error) {
	indexes := prepareIndexes(aggregates)

	commands, close, err := commandsFromAggregates(ctx, aggregates)
	if err != nil {
		return err
	}
	defer close()

	conn, err := crdb.client.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
			return
		}
		err = tx.Commit(ctx)
	}()

	if err = currentSequences(ctx, tx, indexes); err != nil {
		return err
	}

	return push(ctx, tx, indexes, commands)
}

var (
	currentSequencesPrefix = []byte(`SELECT "sequence", "aggregate" FROM eventstore.events WHERE ("sequence", "aggregate") IN (SELECT (max("sequence"), "aggregate") FROM eventstore.events WHERE `)
	currentSequencesSuffix = []byte(` GROUP BY "aggregate") FOR UPDATE`)
)

func currentSequences(ctx context.Context, tx pgx.Tx, indexes *aggregateIndexes) (err error) {
	var builder strings.Builder
	builder.Write(currentSequencesPrefix)
	indexes.currentSequencesClauses(&builder)
	builder.Write(currentSequencesSuffix)

	rows, err := tx.Query(ctx, builder.String(), indexes.toAggregateArgs()...)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			aggregate eventstore.TextSubjects
			sequence  uint32
		)

		if err = rows.Scan(&sequence, &aggregate); err != nil {
			return err
		}

		aggIdx := indexes.byAggregate(aggregate)
		aggIdx.index = sequence
	}

	// check is not made during scan to verify that non existing aggregates are also checked
	for _, aggregate := range indexes.aggregates {
		if aggregate.shouldCheckSequence && aggregate.index != aggregate.expectedSequence {
			return eventstore.ErrSequenceNotMatched
		}
	}

	return nil
}

var (
	pushEventsPrefix = []byte(`WITH computed AS (SELECT hlc_to_timestamp(cluster_logical_timestamp()) created_at, cluster_logical_timestamp() "position"), input ("aggregate", "action", revision, payload, "sequence", in_tx_order) AS (VALUES `)
	pushEventsSuffix = []byte(`) INSERT INTO eventstore.events (created_at, "position", "aggregate", "action", revision, payload, "sequence", in_tx_order) SELECT c.created_at, c."position", i."aggregate", i."action", i.revision, i.payload, i."sequence", i.in_tx_order FROM input i, computed c RETURNING id, created_at`)

	pushActionsPrefix = []byte(`INSERT INTO eventstore.actions ("event", "action", depth) VALUES `)
)

func push(ctx context.Context, tx pgx.Tx, indexes *aggregateIndexes, commands []*command) (err error) {
	var pushBuilder strings.Builder
	pushBuilder.Write(pushEventsPrefix)
	eventsArgs := indexes.eventValues(commands, &pushBuilder)
	pushBuilder.Write(pushEventsSuffix)

	rows, err := tx.Query(ctx, pushBuilder.String(), eventsArgs...)
	if err != nil {
		return err
	}
	defer rows.Close()

	for i := 0; rows.Next(); i++ {
		var creationDate time.Time

		if err = rows.Scan(&commands[i].id, &creationDate); err != nil {
			return fmt.Errorf("push failed: %w", err)
		}
		commands[i].SetCreationDate(creationDate)
		commands[i].SetSequence(commands[i].sequence)
	}

	var actionBuilder strings.Builder
	actionBuilder.Write(pushActionsPrefix)
	actionsArgs := actionValues(commands, &actionBuilder)

	_, err = tx.Exec(ctx, actionBuilder.String(), actionsArgs...)
	if err != nil {
		return err
	}

	if rows.Err() != nil {
		return rows.Err()
	}
	return nil
}

func prepareIndexes(aggregates []eventstore.Aggregate) *aggregateIndexes {
	indexes := &aggregateIndexes{
		aggregates: make([]*aggregateIndex, 0, len(aggregates)),
	}

	for _, aggregate := range aggregates {
		index := indexes.byAggregate(aggregate.ID())
		if index != nil {
			continue
		}
		index = &aggregateIndex{
			aggregate: aggregate.ID(),
		}
		sequenceChecker, ok := aggregate.(eventstore.AggregatePredefinedSequence)
		if ok {
			index.shouldCheckSequence = true
			index.expectedSequence = sequenceChecker.CurrentSequence()
		}
		indexes.aggregates = append(indexes.aggregates, index)
	}

	return indexes
}

type aggregateIndexes struct {
	aggregates []*aggregateIndex
}

type aggregateIndex struct {
	aggregate           eventstore.TextSubjects
	index               uint32
	shouldCheckSequence bool
	expectedSequence    uint32
}

func (indexes *aggregateIndexes) byAggregate(aggregate eventstore.TextSubjects) *aggregateIndex {
	for _, index := range indexes.aggregates {
		if !reflect.DeepEqual(index.aggregate, aggregate) {
			continue
		}
		return index
	}
	return nil
}

func (indexes *aggregateIndexes) increment(aggregate eventstore.TextSubjects) uint32 {
	index := indexes.byAggregate(aggregate)
	if index == nil {
		panic(fmt.Sprintf("aggregate not prepared in indexes: %v", aggregate))
	}
	index.index++
	return index.index
}

func (indexes *aggregateIndexes) toAggregateArgs() []any {
	args := make([]any, len(indexes.aggregates))

	for i, index := range indexes.aggregates {
		args[i] = index.aggregate
	}

	return args
}

var (
	or = []byte(" OR ")
)

func (indexes *aggregateIndexes) currentSequencesClauses(builder *strings.Builder) {
	for i := range indexes.aggregates {
		builder.Write([]byte(`"aggregate" = $` + strconv.Itoa(i+1)))
		if i+1 < len(indexes.aggregates) {
			builder.Write(or)
		}
	}
}

var (
	uuidCast      = []byte("::UUID")
	textCast      = []byte("::TEXT")
	textArrayCast = []byte("::TEXT[]")
	smallIntCast  = []byte("::INT2")
	intCast       = []byte("::INT4")
	jsonbCast     = []byte("::JSONB")
)

func (indexes *aggregateIndexes) eventValues(commands []*command, builder *strings.Builder) []any {
	var (
		index = 0
		args  = make([]any, 0, len(commands)*6)
	)

	for i := 0; i < len(commands); i++ {
		builder.WriteRune('(')

		builder.WriteRune('$')
		builder.Write([]byte(strconv.Itoa(index + 1)))
		builder.Write(textArrayCast)
		builder.WriteRune(',')

		builder.WriteRune('$')
		builder.Write([]byte(strconv.Itoa(index + 2)))
		builder.Write(textArrayCast)
		builder.WriteRune(',')

		builder.WriteRune('$')
		builder.Write([]byte(strconv.Itoa(index + 3)))
		builder.Write(smallIntCast)
		builder.WriteRune(',')

		builder.WriteRune('$')
		builder.Write([]byte(strconv.Itoa(index + 4)))
		builder.Write(jsonbCast)
		builder.WriteRune(',')

		builder.WriteRune('$')
		builder.Write([]byte(strconv.Itoa(index + 5)))
		builder.Write(intCast)
		builder.WriteRune(',')

		builder.WriteRune('$')
		builder.Write([]byte(strconv.Itoa(index + 6)))
		builder.Write(intCast)

		builder.WriteRune(')')

		if i+1 < len(commands) {
			builder.WriteRune(',')
		}
		index += 6

		commands[i].sequence = indexes.increment(commands[i].aggregate)
		args = append(args,
			commands[i].aggregate,
			commands[i].Action(),
			commands[i].Revision(),
			commands[i].payload,
			commands[i].sequence,
			i,
		)
	}

	return args
}

func actionValues(commands []*command, builder *strings.Builder) []any {
	var (
		index = 0
		args  = make([]any, 0, len(commands)*3)
	)

	for cmdCount, cmd := range commands {
		for depth, a := range cmd.Action() {
			builder.WriteRune('(')

			builder.WriteRune('$')
			builder.Write([]byte(strconv.Itoa(index + 1)))
			builder.Write(uuidCast)
			builder.WriteRune(',')

			builder.WriteRune('$')
			builder.Write([]byte(strconv.Itoa(index + 2)))
			builder.Write(textCast)
			builder.WriteRune(',')

			builder.WriteRune('$')
			builder.Write([]byte(strconv.Itoa(index + 3)))
			builder.Write(smallIntCast)

			builder.WriteRune(')')

			if depth+1 < len(cmd.Action()) || cmdCount+1 < len(commands) {
				builder.WriteRune(',')
			}

			index += 3

			args = append(args,
				cmd.id,
				a,
				depth,
			)
		}
	}

	return args
}
