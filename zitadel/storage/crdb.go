package storage

import (
	"context"
	"database/sql"

	"github.com/adlerhurst/eventstore/zitadel"
)

type CRDB struct {
	client *sql.DB
}

func (crdb *CRDB) Ready(ctx context.Context) error {
	return crdb.client.PingContext(ctx)
}

func (crdb *CRDB) Push(context.Context, []zitadel.Command) ([]*zitadel.Event, error) {

	return nil, nil
}
