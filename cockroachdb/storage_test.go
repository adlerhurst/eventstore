package cockroachdb

import (
	"context"
	_ "embed"
	"log"
	"os"
	"testing"
	"time"

	"github.com/adlerhurst/eventstore/v1"
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
	_, err = s.client.Exec(ctx, "TRUNCATE eventstore.events CASCADE")
	return err
}

var store *testStorage

func TestMain(m *testing.M) {
	store = startCRDB()
	os.Exit(m.Run())
}

func startCRDB() *testStorage {
	store := New(&Config{
		Pool: connectToDB(),
	})

	if err := store.Setup(context.Background()); err != nil {
		log.Fatalf("unable to setup cockroach: %v", err)
	}

	return &testStorage{CockroachDB: store}
}

var _ eventstore.Action = (*testAction)(nil)

type testAction struct {
	action   eventstore.TextSubjects
	revision uint16
}

// Action implements eventstore.Action.
func (a *testAction) Action() eventstore.TextSubjects {
	return a.action
}

// Revision implements eventstore.Action.
func (a *testAction) Revision() uint16 {
	return a.revision
}

var _ eventstore.Command = (*testCommand)(nil)

type testCommand struct {
	*testAction
	payload any

	createdAt time.Time
	sequence  uint32
}

// Payload implements eventstore.Command.
func (c *testCommand) Payload() any {
	return c.payload
}

// SetCreationDate implements eventstore.Command.
func (c *testCommand) SetCreationDate(creationDate time.Time) {
	c.createdAt = creationDate
}

// SetSequence implements eventstore.Command.
func (c *testCommand) SetSequence(sequence uint32) {
	c.sequence = sequence
}

func connectToDB() *pgxpool.Pool {
	if isCI := os.Getenv("CI"); isCI == "true" {
		return connectToTestServer()
	}

	return connectToLocalhost()
}

func connectToTestServer() *pgxpool.Pool {
	ts, err := testserver.NewTestServer()
	if err != nil {
		log.Fatal(err)
	}
	pool, err := pgxpool.New(context.Background(), ts.PGURL().String())
	if err != nil {
		log.Fatalf("unable to create database pool: %v", err)
	}
	return pool
}

func connectToLocalhost() *pgxpool.Pool {
	config, err := pgxpool.ParseConfig("postgresql://root@localhost:26257/eventstore?sslmode=disable")
	if err != nil {
		log.Fatalf("unable to parse conn string: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("unable to create database pool: %v", err)
	}

	return pool
}
