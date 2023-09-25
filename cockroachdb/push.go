package cockroachdb

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"text/template"

	// crdbpgx "github.com/cockroachdb/cockroach-go/v2/crdb/crdbpgxv5"
	"github.com/jackc/pgx/v5"

	"github.com/adlerhurst/eventstore/v0"
)

var (
	//go:embed push.sql
	pushStmt string
	//go:embed push_current_sequences.sql
	pushtCurrentSequencesStmt string
)

// Push implements [eventstore.Eventstore]
func (crdb *CockroachDB) Push(ctx context.Context, aggregates ...eventstore.Aggregate) (result []eventstore.Event, err error) {
	indexes := prepareIndexes(aggregates)

	events, err := eventsFromAggregates(aggregates)
	if err != nil {
		return nil, err
	}

	conn, err := crdb.client.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
			return
		}
		err = tx.Commit(ctx)
		if err != nil {
			result = nil
		}
	}()

	if err = currentSequences(ctx, tx, indexes); err != nil {
		return nil, err
	}

	return push(ctx, tx, indexes, events)
}

func currentSequences(ctx context.Context, tx pgx.Tx, indexes *aggregateIndexes) (err error) {
	tmpl := template.
		Must(template.New("push").
			Funcs(template.FuncMap{
				"currentSequencesClauses": indexes.currentSequencesClauses,
			}).
			Parse(pushtCurrentSequencesStmt))

	buf := bytes.NewBuffer(nil)
	if err := tmpl.Execute(buf, nil); err != nil {
		panic(err)
	}

	rows, err := tx.Query(ctx, buf.String(), indexes.toAggregateArgs()...)
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

		index := indexes.byAggregate(aggregate)
		if index.shouldCheckSequence && index.sequence != sequence {
			return eventstore.ErrSequenceNotMatched
		}
		index.index = sequence
	}

	return nil
}

func push(ctx context.Context, tx pgx.Tx, indexes *aggregateIndexes, events []*Event) ([]eventstore.Event, error) {
	pushTmpl := template.
		Must(template.New("push").
			Funcs(template.FuncMap{
				"insertValues": indexes.toValues,
			}).
			Parse(pushStmt))

	buf := bytes.NewBuffer(nil)
	if err := pushTmpl.Execute(buf, nil); err != nil {
		panic(err)
	}

	args := make([]interface{}, 0, len(events)*6)
	result := make([]eventstore.Event, len(events))

	for i, event := range events {
		event.sequence = indexes.increment(event.Aggregate())
		args = append(args,
			event.aggregate,
			event.action,
			event.revision,
			event.payload,
			event.sequence,
			i,
		)
		result[i] = event
	}

	rows, err := tx.Query(ctx, buf.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for i := 0; rows.Next(); i++ {
		if err = rows.Scan(&events[i].creationDate, &events[i].position); err != nil {
			return nil, fmt.Errorf("push failed: %w", err)
		}
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return result, nil
}

func prepareIndexes(aggregates []eventstore.Aggregate) *aggregateIndexes {
	indexes := &aggregateIndexes{
		aggregates: make([]*aggregateIndex, 0, len(aggregates)),
	}

	for _, aggregate := range aggregates {
		indexes.commandCount += len(aggregate.Commands())
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
			index.sequence = sequenceChecker.CurrentSequence()
		}
		indexes.aggregates = append(indexes.aggregates, index)
	}

	return indexes
}

type aggregateIndexes struct {
	aggregates   []*aggregateIndex
	commandCount int
}

type aggregateIndex struct {
	aggregate           eventstore.TextSubjects
	index               uint32
	shouldCheckSequence bool
	sequence            uint32
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

func (indexes *aggregateIndexes) toAggregateArgs() []interface{} {
	args := make([]interface{}, len(indexes.aggregates))

	for i, index := range indexes.aggregates {
		args[i] = index.aggregate
	}

	return args
}

func (indexes *aggregateIndexes) currentSequencesClauses() string {
	clauses := make([]string, len(indexes.aggregates))

	for i := range indexes.aggregates {
		clauses[i] = `"aggregate" = $` + strconv.Itoa(i+1)
	}

	return strings.Join(clauses, " OR ")
}

func (indexes *aggregateIndexes) toValues() string {
	values := make([]string, indexes.commandCount)
	var index = 0
	for i := 0; i < indexes.commandCount; i++ {
		values[i] = fmt.Sprintf("($%d::TEXT[], $%d::TEXT[], $%d::INT2, $%d::JSONB, $%d::INT4, $%d::INT4)", index+1, index+2, index+3, index+4, index+5, index+6)
		index += 6
	}

	return strings.Join(values, ", ")
}
