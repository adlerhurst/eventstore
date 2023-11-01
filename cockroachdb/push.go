package cockroachdb

import (
	"context"
	_ "embed"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	crdb "github.com/cockroachdb/cockroach-go/v2/crdb/crdbpgxv5"
	"github.com/jackc/pgx/v5"

	"github.com/adlerhurst/eventstore"
)

var pushTxOptions = pgx.TxOptions{
	IsoLevel:   pgx.Serializable,
	AccessMode: pgx.ReadWrite,
}

// Push implements [eventstore.Eventstore]
func (store *CockroachDB) Push(ctx context.Context, aggregates ...eventstore.Aggregate) (err error) {
	indexes := prepareIndexes(aggregates)

	commands, close, err := commandsFromAggregates(ctx, aggregates)
	if err != nil {
		return err
	}
	defer close()

	conn, err := store.client.Acquire(ctx)
	if err != nil {
		logger.ErrorContext(ctx, "acquire connection failed", "cause", err)
		return err
	}
	defer conn.Release()

	// The application_name is required to differentiate between filter and push queries
	// filter queries are not relevant for the [filterIgnoreOpenPush] clause
	_, err = conn.Exec(ctx, "SET application_name = $1", store.pushAppName)
	if err != nil {
		logger.ErrorContext(ctx, "set application name failed", "cause", err)
		return err
	}

	return crdb.ExecuteTx(ctx, conn, pushTxOptions, func(tx pgx.Tx) error {
		if err = currentSequences(ctx, tx, indexes); err != nil {
			return err
		}

		return push(ctx, tx, indexes, commands)
	})
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
		logger.ErrorContext(ctx, "query current sequences failed", "cause", err)
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			aggregate eventstore.TextSubjects
			sequence  uint32
		)

		if err = rows.Scan(&sequence, &aggregate); err != nil {
			logger.ErrorContext(ctx, "scan of sequences failed", "cause", err)
			return err
		}

		aggIdx := indexes.byAggregate(aggregate)
		aggIdx.index = sequence
	}

	// check is not made during scan to verify that non existing aggregates are also checked
	for _, aggregate := range indexes.aggregates {
		if aggregate.shouldCheckSequence && aggregate.index != aggregate.expectedSequence {
			logger.DebugContext(ctx, "unexpected sequence", "expected", aggregate.expectedSequence, "got", aggregate.index)
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
		logger.ErrorContext(ctx, "store commands failed", "cause", err)
		return err
	}
	defer rows.Close()

	for i := 0; rows.Next(); i++ {
		var creationDate time.Time

		if err = rows.Scan(&commands[i].id, &creationDate); err != nil {
			logger.ErrorContext(ctx, "scan of returned command metadata failed", "cause", err)
			return fmt.Errorf("push failed: %w", err)
		}
		commands[i].SetCreationDate(creationDate)
		commands[i].SetSequence(commands[i].sequence)
	}

	if rows.Err() != nil {
		logger.ErrorContext(ctx, "push failed", "cause", rows.Err())
		return rows.Err()
	}

	var actionBuilder strings.Builder
	actionBuilder.Write(pushActionsPrefix)
	actionsArgs := actionValues(commands, &actionBuilder)

	_, err = tx.Exec(ctx, actionBuilder.String(), actionsArgs...)
	if err != nil {
		logger.ErrorContext(ctx, "store actions failed", "cause", err)
		return err
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
		if sequence := aggregate.CurrentSequence(); sequence != nil {
			index.shouldCheckSequence = true
			index.expectedSequence = *sequence
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
