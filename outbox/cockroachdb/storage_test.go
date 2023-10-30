package cockroachdb

import (
	"context"
	_ "embed"
	"log"
	"os"
	"testing"

	"github.com/adlerhurst/eventstore/v0"
	"github.com/cockroachdb/cockroach-go/v2/testserver"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ eventstore.TestEventstore = (*testStorage)(nil)

type testStorage struct {
	*CockroachDB
}

// After implements eventstore.TestEventstore
func (*testStorage) After(ctx context.Context, t testing.TB) error {
	return nil
}

// Before implements eventstore.TestEventstore
func (s *testStorage) Before(ctx context.Context, t testing.TB) (err error) {
	_, err = s.client.Exec(ctx, "TRUNCATE outbox.events CASCADE")
	return err
}

var store *testStorage

func TestMain(m *testing.M) {
	store = startCRDB()
	os.Exit(m.Run())
}

func startCRDB() *testStorage {
	var ts *testserver.TestServer
	_ = ts
	// ts, err := testserver.NewTestServer()
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// dbpool, err := pgxpool.New(context.Background(), ts.PGURL().String())
	config, err := pgxpool.ParseConfig("postgresql://root@localhost:26257/eventstore?sslmode=disable&application_name=bench4")
	if err != nil {
		log.Fatalf("unable to parse conn string: %v", err)
	}

	dbpool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("unable to create database pool: %v", err)
	}

	crdb := New(&Config{
		Pool: dbpool,
	})

	if err := crdb.Setup(context.Background()); err != nil {
		log.Fatalf("unable to setup cockroach: %v", err)
	}

	return &testStorage{CockroachDB: crdb}
}
