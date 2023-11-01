package cockroachdb

import (
	"context"
	_ "embed"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/adlerhurst/eventstore"
	"github.com/adlerhurst/eventstore/x"
)

type Config struct {
	Pool *pgxpool.Pool
}

var (
	_           eventstore.Eventstore = (*CockroachDB)(nil)
	logger                            = slog.Default()
	eventPool                         = x.NewPool[event]()
	commandPool                       = x.NewPool[command]()
)

type CockroachDB struct {
	client        *pgxpool.Pool
	pushAppName   string
	filterAppName string
}

func New(config *Config, opts ...storageOpt) *CockroachDB {
	store := &CockroachDB{
		client:        config.Pool,
		pushAppName:   "es_push",
		filterAppName: "es_filter",
	}

	for _, opt := range opts {
		opt(store)
	}

	return store
}

type storageOpt func(*CockroachDB)

func WithLogger(l *slog.Logger) storageOpt {
	return func(store *CockroachDB) {
		logger = l
	}
}

func WithPushAppName(name string) storageOpt {
	return func(store *CockroachDB) {
		store.pushAppName = name
	}
}

func WithFilterAppName(name string) storageOpt {
	return func(store *CockroachDB) {
		store.filterAppName = name
	}
}

//go:embed 0_setup.sql
var setupStmt string

func (store *CockroachDB) Setup(ctx context.Context) error {
	_, err := store.client.Exec(ctx, setupStmt)
	logger.ErrorContext(ctx, "setup failed", "cause", err)
	return err
}

// Ready implements [eventstore.Eventstore]
func (store *CockroachDB) Ready(ctx context.Context) error {
	return store.client.Ping(ctx)
}
