package memory

import (
	"context"
	_ "embed"
	"sync"
	"text/template"

	"github.com/adlerhurst/eventstore/v0"
	"github.com/jackc/pgx/v5/pgxpool"
)

func init() {
	filterTmpl = template.Must(template.New("filter").Parse(filterStmt))
}

type Config struct {
	Pool *pgxpool.Pool
}

var _ eventstore.Eventstore = (*CockroachDB)(nil)

type CockroachDB struct {
	client *pgxpool.Pool
	outbox outbox
}

func New(config *Config) *CockroachDB {
	return &CockroachDB{
		config.Pool,
		outbox{
			mu:            sync.RWMutex{},
			subscriptions: &subscription{},
		},
	}
}

//go:embed 0_setup.sql
var setupStmt string

func (crdb *CockroachDB) Setup(ctx context.Context) error {
	_, err := crdb.client.Exec(ctx, setupStmt)
	return err
}

// Ready implements [eventstore.Eventstore]
func (crdb *CockroachDB) Ready(ctx context.Context) error {
	return crdb.client.Ping(ctx)
}
