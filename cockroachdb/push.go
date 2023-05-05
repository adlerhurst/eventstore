package cockroachdb

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"text/template"

	"github.com/adlerhurst/eventstore/v0"
	"github.com/jackc/pgx/v5"
)

var (
	//go:embed push.sql
	pushStmt string
	//go:embed push_current_sequences.sql
	pushtCurrentSequencesStmt string
)

// Push implements [eventstore.Eventstore]
func (crdb *CockroachDB) Push(ctx context.Context, commands ...eventstore.Command) (result []eventstore.Event, err error) {
	indexes := prepareIndexes(commands)

	events, err := eventsFromCommands(commands)
	if err != nil {
		return nil, err
	}

	for i := 0; i < 10; i++ {
		result, err = crdb.retriablePush(ctx, indexes, events)
		if err == nil {
			break
		}
		log.Printf("push %d failed: %v\n", i, err)
	}

	return result, err
}

func (crdb *CockroachDB) retriablePush(ctx context.Context, indexes *aggregateIndexes, events []*Event) ([]eventstore.Event, error) {
	tx, err := crdb.client.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
			return
		}
		err = tx.Commit(ctx)
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
			sequence  uint64
		)

		if err = rows.Scan(&sequence, &aggregate); err != nil {
			return err
		}

		index := indexes.byAggregate(aggregate)
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
		args = append(args,
			event.aggregate,
			event.action,
			event.revision,
			event.metadata,
			event.payload,
			indexes.increment(event.Aggregate()),
		)
		result[i] = event
	}

	rows, err := tx.Query(ctx, buf.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for i := 0; rows.Next(); i++ {
		if err = rows.Scan(&events[i].sequence, &events[i].creationDate); err != nil {
			return nil, fmt.Errorf("push failed: %w", err)
		}
	}

	return result, nil
}

func prepareIndexes(commands []eventstore.Command) *aggregateIndexes {
	indexes := &aggregateIndexes{
		aggregates:   make([]*aggregateIndex, 0, len(commands)),
		commandCount: len(commands),
	}

	for _, command := range commands {
		index := indexes.byAggregate(command.Aggregate())
		if index != nil {
			continue
		}
		indexes.aggregates = append(indexes.aggregates, &aggregateIndex{
			aggregate: command.Aggregate(),
		})
	}

	return indexes
}

type aggregateIndexes struct {
	aggregates   []*aggregateIndex
	commandCount int
}

type aggregateIndex struct {
	aggregate eventstore.TextSubjects
	index     uint64
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

func (indexes *aggregateIndexes) increment(aggregate eventstore.TextSubjects) uint64 {
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
		clauses[i] = `"aggregate" @> $` + strconv.Itoa(i+1)
	}

	return strings.Join(clauses, " OR ")
}

func (indexes *aggregateIndexes) toValues() string {
	values := make([]string, indexes.commandCount)
	var index = 0
	for i := 0; i < indexes.commandCount; i++ {
		values[i] = fmt.Sprintf("($%[1]d, array_to_string($%[1]d, ':'), $%[2]d, $%[3]d, $%[4]d, $%[5]d, $%[6]d)", index+1, index+2, index+3, index+4, index+5, index+6)
		index += 6
	}

	return strings.Join(values, ", ")
}
