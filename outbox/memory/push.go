package memory

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"text/template"

	crdbpgx "github.com/cockroachdb/cockroach-go/v2/crdb/crdbpgxv5"
	"github.com/jackc/pgx/v5"

	"github.com/adlerhurst/eventstore/v0"
)

var (
	//go:embed push.sql
	pushStmt string
	//go:embed outbox.sql
	outboxStmt string
	//go:embed push_current_sequences.sql
	pushtCurrentSequencesStmt string
)

// Push implements [eventstore.Eventstore]
func (crdb *CockroachDB) Push(ctx context.Context, commands ...eventstore.Command) (result []eventstore.Event, err error) {
	indexes := prepareIndexes(commands)

	events, err := crdb.eventsFromCommands(commands)
	if err != nil {
		return nil, err
	}

	conn, err := crdb.client.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	err = crdbpgx.ExecuteTx(ctx, conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
		if err = currentSequences(ctx, tx, indexes); err != nil {
			return err
		}

		result, err = push(ctx, tx, indexes, events)
		if err != nil {
			return err
		}
		return crdb.pushToOutbox(ctx, tx, events)
	})

	if err != nil {
		return nil, err
	}
	return result, nil
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

	if rows.Err() != nil {
		return rows.Err()
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

	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return result, nil
}

func (crdb *CockroachDB) pushToOutbox(ctx context.Context, tx pgx.Tx, events []*Event) error {
	params := make([]string, 0, len(events))
	args := make([]any, 0, len(events)*4)

	var paramIndex int
	for _, event := range events {
		for _, r := range event.receivers {
			params = append(params, fmt.Sprintf("($%d, $%d, $%d, $%d)", paramIndex+1, paramIndex+2, paramIndex+3, paramIndex+4))
			paramIndex += 4
			args = append(args, event.aggregate, event.sequence, event.creationDate, r)
		}
	}

	if len(params) == 0 {
		return nil
	}

	pushTmpl := template.
		Must(template.New("outbox").
			Funcs(template.FuncMap{
				"insertValues": func() string {
					return strings.Join(params, ", ")
				},
			}).
			Parse(outboxStmt))

	buf := bytes.NewBuffer(nil)
	if err := pushTmpl.Execute(buf, nil); err != nil {
		panic(err)
	}

	_, err := tx.Exec(ctx, buf.String(), args...)
	return err
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
		clauses[i] = `"aggregate" = $` + strconv.Itoa(i+1)
	}

	return strings.Join(clauses, " OR ")
}

func (indexes *aggregateIndexes) toValues() string {
	values := make([]string, indexes.commandCount)
	var index = 0
	for i := 0; i < indexes.commandCount; i++ {
		values[i] = fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d)", index+1, index+2, index+3, index+4, index+5, index+6)
		index += 6
	}

	return strings.Join(values, ", ")
}
