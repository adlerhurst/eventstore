package zitadel_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"testing"

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

	eventstores := createEventstores(b, db)
	payloads := createPayloads(b)

	cmd := new(testCommand)

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

func TestEventstore_Filter(t *testing.T) {
	eventstores := createEventstores(t, testCRDBClient)

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
				// if !reflect.DeepEqual(got, tt.want) {
				// 	t.Errorf("Eventstore.Push() = %v, want %v", got, tt.want)
				// }
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
	crdb3, err := storage.NewCRDB3(db)
	if err != nil {
		f.Fatalf("unable to mock database: %v", err)
	}

	return map[string]*zitadel.Eventstore{
		// "crdb1":   zitadel.NewEventstore(crdb1),
		// "crdb2":   zitadel.NewEventstore(crdb2),
		"crdb2_2": zitadel.NewEventstore(crdb2_2),
		"crdb3":   zitadel.NewEventstore(crdb3),
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
}

func (cmd *testCommand) Aggregate() zitadel.Aggregate {
	return zitadel.Aggregate{
		ID:            "1",
		Type:          "testAgg",
		ResourceOwner: "ro",
		InstanceID:    "instance",
		Version:       "v1",
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
