package cockroachdb

import (
	"context"
	_ "embed"

	"github.com/adlerhurst/eventstore/v0"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func New(client *pgxpool.Pool) *CockroachDB {
	return &CockroachDB{
		client: client,
	}
}

var (
	_ eventstore.Eventstore = (*CockroachDB)(nil)
	//go:embed 0_setup.sql
	setup string
	//go:embed push.sql
	push string
)

type CockroachDB struct {
	client *pgxpool.Pool
}

func (crdb *CockroachDB) Setup(ctx context.Context) error {
	_, err := crdb.client.Exec(ctx, setup)
	return err
}

// Filter implements eventstore.Eventstore
func (*CockroachDB) Filter(context.Context, *eventstore.Filter) ([]eventstore.Event, error) {
	panic("unimplemented")
}

// Push implements eventstore.Eventstore
func (crdb *CockroachDB) Push(ctx context.Context, commands ...eventstore.Command) ([]eventstore.Event, error) {
	events, err := eventsFromCommands(commands)
	if err != nil {
		return nil, err
	}

	result := make([]eventstore.Event, len(events))
	tx, err := crdb.client.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	for i, event := range events {
		row := tx.QueryRow(ctx, push,
			event.action,
			event.aggregate,
			event.revision,
			event.metadata,
			event.payload,
		)

		result[i] = event
	}

	return result, nil
}

// Ready implements eventstore.Eventstore
func (crdb *CockroachDB) Ready(ctx context.Context) error {
	return crdb.client.Ping(ctx)
}
