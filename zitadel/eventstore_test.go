package zitadel_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/adlerhurst/eventstore/zitadel"
	"github.com/adlerhurst/eventstore/zitadel/storage"

	//sql import
	_ "github.com/jackc/pgx/v4/stdlib"
)

func TestEventstore_Push(t *testing.T) {
	eventstores := createEventstores(t, testCRDBClient)

	type args struct {
		cmds []zitadel.Command
	}
	tests := []struct {
		name    string
		args    args
		want    []*zitadel.Event
		wantErr bool
	}{
		{
			name: "sinlge event",
			args: args{
				cmds: []zitadel.Command{
					&testCommand{},
				},
			},
			want:    []*zitadel.Event{},
			wantErr: false,
		},
		{
			name: "2 events",
			args: args{
				cmds: []zitadel.Command{
					&testCommand{},
					&testCommand{},
				},
			},
			want:    []*zitadel.Event{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		for esKey, es := range eventstores {
			t.Run(fmt.Sprintf("%s %s", esKey, tt.name), func(t *testing.T) {
				_, err := es.Push(context.Background(), tt.args.cmds)
				if (err != nil) != tt.wantErr {
					t.Errorf("Eventstore.Push() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				// if !reflect.DeepEqual(got, tt.want) {
				// 	t.Errorf("Eventstore.Push() = %v, want %v", got, tt.want)
				// }
			})
		}
	}
}

func BenchmarkEventstorePush(b *testing.B) {
	db, err := sql.Open("pgx", "postgresql://root@localhost:26257/defaultdb?sslmode=disable")
	if err != nil {
		b.Fatal("unable to connect to db")
	}
	defer db.Close()
	db.SetMaxIdleConns(20)
	db.SetMaxOpenConns(20)

	eventstores := createEventstores(b, db)
	payloads := createPayloads(b)

	cmd := new(testCommand)

	cmds8000 := make([]zitadel.Command, 8000)
	for i := 0; i < 8000; i++ {
		cmds8000[i] = cmd
	}

	tests := []struct {
		name string
		cmds []zitadel.Command
	}{
		{
			name: "1 event",
			cmds: []zitadel.Command{
				cmd,
			},
		},
		{
			name: "2 events",
			cmds: []zitadel.Command{
				cmd,
				cmd,
			},
		},
		{
			name: "3 events",
			cmds: []zitadel.Command{
				cmd,
				cmd,
				cmd,
			},
		},
		{
			name: "4 events",
			cmds: []zitadel.Command{
				cmd,
				cmd,
				cmd,
				cmd,
			},
		},
		{
			name: "5 event",
			cmds: []zitadel.Command{
				cmd,
				cmd,
				cmd,
				cmd,
				cmd,
			},
		},
		{
			name: "10 event",
			cmds: []zitadel.Command{
				cmd,
				cmd,
				cmd,
				cmd,
				cmd,
				cmd,
				cmd,
				cmd,
				cmd,
				cmd,
			},
		},
		{
			name: "8000 event",
			cmds: cmds8000,
		},
	}
	for esKey, es := range eventstores {
		for payloadKey, payload := range payloads {
			cmd.payload = payload
			for _, tt := range tests {
				b.Run(fmt.Sprintf("%s %s %s", esKey, payloadKey, tt.name), func(b *testing.B) {
					for n := 0; n < b.N; n++ {
						_, err := es.Push(
							context.Background(),
							tt.cmds,
						)
						if err != nil {
							b.Error(err)
						}
					}
				})
			}
		}
	}
}

func BenchmarkEventstorePushParallel(b *testing.B) {
	db, err := sql.Open("pgx", "postgresql://root@localhost:26257/defaultdb?sslmode=disable")
	if err != nil {
		b.Fatal("unable to connect to db")
	}
	defer db.Close()
	db.SetMaxIdleConns(20)
	db.SetConnMaxIdleTime(30 * time.Second)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(10 * time.Minute)

	eventstores := createEventstores(b, db)
	payloads := createPayloads(b)

	cmd := new(testCommand)

	cmds8000 := make([]zitadel.Command, 8000)
	for i := 0; i < 8000; i++ {
		cmds8000[i] = cmd
	}

	tests := []struct {
		name string
		cmds []zitadel.Command
	}{
		{
			name: "1 event",
			cmds: []zitadel.Command{
				cmd,
			},
		},
		{
			name: "2 events",
			cmds: []zitadel.Command{
				cmd,
				cmd,
			},
		},
		{
			name: "3 events",
			cmds: []zitadel.Command{
				cmd,
				cmd,
				cmd,
			},
		},
		{
			name: "4 events",
			cmds: []zitadel.Command{
				cmd,
				cmd,
				cmd,
				cmd,
			},
		},
		{
			name: "5 event",
			cmds: []zitadel.Command{
				cmd,
				cmd,
				cmd,
				cmd,
				cmd,
			},
		},
		{
			name: "10 event",
			cmds: []zitadel.Command{
				cmd,
				cmd,
				cmd,
				cmd,
				cmd,
				cmd,
				cmd,
				cmd,
				cmd,
				cmd,
			},
		},
		// {
		// 	name: "8000 event",
		// 	cmds: cmds8000,
		// },
	}

	for esKey, es := range eventstores {
		cmds := make(chan []zitadel.Command, 10)
		// workers
		for i := 0; i < 10; i++ {
			go pushWorker(b, i, es, cmds)
		}

		for payloadKey, payload := range payloads {
			cmd.payload = payload
			for _, tt := range tests {
				b.Run(fmt.Sprintf("%s %s %s", esKey, payloadKey, tt.name), func(b *testing.B) {
					for n := 0; n < b.N; n++ {
						cmds <- tt.cmds
					}
				})
			}
		}

		close(cmds)
	}
}

func pushWorker(b *testing.B, id int, es *zitadel.Eventstore, cmds <-chan []zitadel.Command) {
	aggID := strconv.Itoa(id)
	for cmd := range cmds {
		cpy := make([]zitadel.Command, len(cmd))
		copy(cpy, cmd)
		for _, c := range cpy {
			testC := c.(*testCommand)
			testC.aggID = aggID
		}
		_, err := es.Push(
			context.Background(),
			cpy,
		)
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkEventstorePushFilterParallel(b *testing.B) {
	db, err := sql.Open("pgx", "postgresql://root@localhost:26257/defaultdb?sslmode=disable")
	if err != nil {
		b.Fatal("unable to connect to db")
	}
	defer db.Close()
	db.SetMaxIdleConns(20)
	db.SetConnMaxIdleTime(30 * time.Second)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(10 * time.Minute)

	eventstores := createEventstores(b, db)
	payloads := createPayloads(b)

	cmd := new(testCommand)

	filter := &zitadel.Filter{
		Aggregates: []*zitadel.AggregateFilter{
			{
				ID: "1",
			},
		},
	}

	tests := []struct {
		name   string
		cmds   []zitadel.Command
		filter *zitadel.Filter
	}{
		{
			name: "1 event",
			cmds: []zitadel.Command{
				cmd,
			},
			filter: filter,
		},
		{
			name: "2 events",
			cmds: []zitadel.Command{
				cmd,
				cmd,
			},
			filter: filter,
		},
		{
			name: "3 events",
			cmds: []zitadel.Command{
				cmd,
				cmd,
				cmd,
			},
			filter: filter,
		},
		{
			name: "4 events",
			cmds: []zitadel.Command{
				cmd,
				cmd,
				cmd,
				cmd,
			},
			filter: filter,
		},
		{
			name: "5 event",
			cmds: []zitadel.Command{
				cmd,
				cmd,
				cmd,
				cmd,
				cmd,
			},
			filter: filter,
		},
		{
			name: "10 event",
			cmds: []zitadel.Command{
				cmd,
				cmd,
				cmd,
				cmd,
				cmd,
				cmd,
				cmd,
				cmd,
				cmd,
				cmd,
			},
			filter: filter,
		},
	}

	for esKey, es := range eventstores {
		cmdCount := 2
		filterCount := 18
		cmds := make(chan []zitadel.Command, cmdCount)
		filters := make(chan *zitadel.Filter, filterCount)
		// workers
		for i := 0; i < 2; i++ {
			go pushWorker(b, i, es, cmds)
		}
		for i := 0; i < 18; i++ {
			go filterWorker(b, i, es, filters)
		}

		for payloadKey, payload := range payloads {
			cmd.payload = payload
			for _, tt := range tests {
				b.Run(fmt.Sprintf("%s %s %s", esKey, payloadKey, tt.name), func(b *testing.B) {
					for n := 0; n < b.N; n++ {
						wg := sync.WaitGroup{}
						wg.Add(20)
						go func() {
							for i := 0; i < filterCount; i++ {
								filters <- tt.filter
								wg.Done()
							}
						}()
						go func() {
							for i := 0; i < cmdCount; i++ {
								cmds <- tt.cmds
								wg.Done()
							}
						}()
						wg.Wait()
					}
				})
			}
		}

		close(cmds)
	}
}

func filterWorker(b *testing.B, id int, es *zitadel.Eventstore, filters <-chan *zitadel.Filter) {
	aggID := strconv.Itoa(id % 2)
	for cmd := range filters {
		cpy := *cmd
		cpy.Aggregates[0].ID = aggID
		cpy.CreationDateGreaterEqual = time.Now().Add(-10 * time.Millisecond)
		_, err := es.Filter(context.Background(), &cpy)
		if err != nil {
			b.Error(err)
		}
	}
}

func TestEventstore_Filter(t *testing.T) {
	db, err := sql.Open("pgx", "postgresql://root@localhost:26257/defaultdb?sslmode=disable")
	if err != nil {
		t.Fatal("unable to connect to db")
	}
	defer db.Close()

	eventstores := createEventstores(t, db)

	defaults := []zitadel.Command{
		&testCommand{},
		&testCommand{},
		&testCommand{},
	}

	type args struct {
		filter *zitadel.Filter
	}
	tests := []struct {
		name    string
		args    args
		want    []*zitadel.Event
		wantErr bool
	}{
		{
			name: "no events",
			args: args{
				filter: &zitadel.Filter{
					InstanceID: "no instance",
				},
			},
			want:    []*zitadel.Event{},
			wantErr: false,
		},
		{
			name: "simple filter",
			args: args{
				filter: &zitadel.Filter{
					InstanceID: "instance",
					Limit:      3,
				},
			},
			want: []*zitadel.Event{
				nil,
				nil,
				nil,
			},
			wantErr: false,
		},
		{
			name: "complex filter",
			args: args{
				filter: &zitadel.Filter{
					InstanceID: "instance",
					Limit:      3,
					Desc:       true,
					Aggregates: []*zitadel.AggregateFilter{
						{
							Type: "testAgg",
							ID:   "1",
						},
					},
				},
			},
			want: []*zitadel.Event{
				nil,
				nil,
				nil,
			},
			wantErr: false,
		},
		{
			name: "complex filter with lists",
			args: args{
				filter: &zitadel.Filter{
					InstanceID: "instance",
					Limit:      3,
					Desc:       true,
					OrgIDs:     []string{"ro"},
					Aggregates: []*zitadel.AggregateFilter{
						{
							Type: "testAgg",
							ID:   "1",
							Events: []*zitadel.EventFilter{
								{
									Types: []string{"event.type"},
								},
							},
						},
					},
				},
			},
			want: []*zitadel.Event{
				nil,
				nil,
				nil,
			},
			wantErr: false,
		},
	}

	for esKey, es := range eventstores {
		if _, err := es.Push(context.Background(), defaults); err != nil {
			t.Fatalf("unable to push default events: %v", err)
		}
		for _, tt := range tests {
			t.Run(fmt.Sprintf("%s %s", esKey, tt.name), func(t *testing.T) {
				events, err := es.Filter(context.Background(), tt.args.filter)
				if errors.Is(err, storage.ErrUnimplemented) {
					return
				}
				if (err != nil) != tt.wantErr {
					t.Errorf("Eventstore.Filter() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if len(tt.want) != len(events) {
					t.Errorf("Eventstore.Filter() expected event count %d, got %d", len(tt.want), len(events))
				}
			})
		}
	}
}

type fataler interface {
	Fatalf(string, ...any)
}

func createEventstores(f fataler, db *sql.DB) map[string]*zitadel.Eventstore {
	// crdb1, err := storage.NewCRDB1(db)
	// if err != nil {
	// 	f.Fatalf("unable to mock database: %v", err)
	// }
	// crdb2, err := storage.NewCRDB1(db)
	// if err != nil {
	// 	f.Fatalf("unable to mock database: %v", err)
	// }
	crdb2_2, err := storage.NewCRDB2_2(db)
	if err != nil {
		f.Fatalf("unable to mock database: %v", err)
	}
	// crdb3, err := storage.NewCRDB3(db)
	// if err != nil {
	// 	f.Fatalf("unable to mock database: %v", err)
	// }

	return map[string]*zitadel.Eventstore{
		// "crdb1":   zitadel.NewEventstore(crdb1),
		// "crdb2":   zitadel.NewEventstore(crdb2),
		"crdb2_2": zitadel.NewEventstore(crdb2_2),
		// "crdb3":   zitadel.NewEventstore(crdb3),
	}
}

type testPayload struct {
	Firstname   string
	Lastname    string
	DisplayName string
	Gender      int8
	LoginNames  []string
}

func createPayloads(f fataler) map[string]interface{} {
	jsonPayload, err := json.Marshal(testPayload{
		Firstname:   "adler",
		Lastname:    "hurst",
		DisplayName: "adlerhurst",
		Gender:      2,
		LoginNames:  []string{"silvan@zitadel.com", "adlerhurst@zitadel.com", "adlerhurst", "adlerhurst@my-comp.zitadel.cloud"},
	})
	if err != nil {
		f.Fatalf("unable to create payload: %v", err)
	}

	return map[string]interface{}{
		"no payload":   nil,
		"json payload": jsonPayload,
		"struct payload": testPayload{
			Firstname:   "adler",
			Lastname:    "hurst",
			DisplayName: "adlerhurst",
			Gender:      2,
			LoginNames:  []string{"silvan@zitadel.com", "adlerhurst@zitadel.com", "adlerhurst", "adlerhurst@my-comp.zitadel.cloud"},
		},
		"pointer payload": &testPayload{
			Firstname:   "adler",
			Lastname:    "hurst",
			DisplayName: "adlerhurst",
			Gender:      2,
			LoginNames:  []string{"silvan@zitadel.com", "adlerhurst@zitadel.com", "adlerhurst", "adlerhurst@my-comp.zitadel.cloud"},
		},
	}
}

type testCommand struct {
	payload interface{}
	aggID   string
}

func (cmd *testCommand) Aggregate() zitadel.Aggregate {
	aggID := cmd.aggID
	if aggID == "" {
		aggID = "1"
	}
	return zitadel.Aggregate{
		ID:            aggID,
		Type:          "testAgg",
		ResourceOwner: "ro",
		InstanceID:    "instance",
	}
}

func (cmd *testCommand) EditorUser() string {
	return "usr"
}

func (cmd *testCommand) Type() string {
	return "event.type"
}

func (cmd *testCommand) Payload() interface{} {
	return cmd.payload
}

func (cmd *testCommand) Version() string {
	return "v1"
}
