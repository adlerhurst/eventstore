package zitadel_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/adlerhurst/eventstore/zitadel"
	"github.com/adlerhurst/eventstore/zitadel/storage"

	//sql import
	_ "github.com/jackc/pgx/v4/stdlib"
)

func TestEventstore_Push(t *testing.T) {
	db, err := sql.Open("pgx", "postgresql://root@localhost:26257/defaultdb?sslmode=disable")
	if err != nil {
		t.Fatal("unable to connect to db")
	}
	defer db.Close()
	// crdb1, err := storage.NewCRDB1(testCRDBClient)
	// if err != nil {
	// 	t.Fatalf("unable to mock database: %v", err)
	// }
	// crdb2, err := storage.NewCRDB2(testCRDBClient)
	// if err != nil {
	// 	t.Fatalf("unable to mock database: %v", err)
	// }
	crdb2_2, err := storage.NewCRDB2_2(db)
	if err != nil {
		t.Fatalf("unable to mock database: %v", err)
	}

	type args struct {
		cmds []zitadel.Command
	}
	tests := []struct {
		name    string
		es      *zitadel.Eventstore
		args    args
		want    []*zitadel.Event
		wantErr bool
	}{
		// {
		// 	name: "crdb1 sinlge event",
		// 	es:   zitadel.NewEventstore(crdb1),
		// 	args: args{
		// 		cmds: []zitadel.Command{
		// 			&testCommand{},
		// 		},
		// 	},
		// 	want:    []*zitadel.Event{},
		// 	wantErr: false,
		// },
		// {
		// 	name: "crdb1 2 events",
		// 	es:   zitadel.NewEventstore(crdb1),
		// 	args: args{
		// 		cmds: []zitadel.Command{
		// 			&testCommand{},
		// 			&testCommand{},
		// 		},
		// 	},
		// 	want:    []*zitadel.Event{},
		// 	wantErr: false,
		// },
		// {
		// 	name: "crdb2 sinlge event",
		// 	es:   zitadel.NewEventstore(crdb2),
		// 	args: args{
		// 		cmds: []zitadel.Command{
		// 			&testCommand{},
		// 		},
		// 	},
		// 	want:    []*zitadel.Event{},
		// 	wantErr: false,
		// },
		// {
		// 	name: "crdb2 2 events",
		// 	es:   zitadel.NewEventstore(crdb2),
		// 	args: args{
		// 		cmds: []zitadel.Command{
		// 			&testCommand{},
		// 			&testCommand{},
		// 		},
		// 	},
		// 	want:    []*zitadel.Event{},
		// 	wantErr: false,
		// },
		// {
		// 	name: "crdb2_2 sinlge event",
		// 	es:   zitadel.NewEventstore(crdb2_2),
		// 	args: args{
		// 		cmds: []zitadel.Command{
		// 			&testCommand{},
		// 		},
		// 	},
		// 	want:    []*zitadel.Event{},
		// 	wantErr: false,
		// },
		{
			name: "crdb2_2 2 events",
			es:   zitadel.NewEventstore(crdb2_2),
			args: args{
				cmds: []zitadel.Command{
					&testCommand{},
					&testCommand{},
					&testCommand{},
					&testCommand{},
					&testCommand{},
					&testCommand{},
					&testCommand{},
					&testCommand{},
					&testCommand{},
					&testCommand{},
					&testCommand{},
					&testCommand{},
					&testCommand{},
				},
			},
			want:    []*zitadel.Event{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.es.Push(context.Background(), tt.args.cmds)
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

func BenchmarkEventstorePush(b *testing.B) {
	db, err := sql.Open("pgx", "postgresql://root@localhost:26257/defaultdb?sslmode=disable")
	if err != nil {
		b.Fatal("unable to connect to db")
	}
	defer db.Close()

	crdb1, err := storage.NewCRDB1(db)
	if err != nil {
		b.Fatalf("unable to mock database: %v", err)
	}
	crdb2, err := storage.NewCRDB1(db)
	if err != nil {
		b.Fatalf("unable to mock database: %v", err)
	}
	crdb2_2, err := storage.NewCRDB2_2(db)
	if err != nil {
		b.Fatalf("unable to mock database: %v", err)
	}

	tests := []struct {
		name    string
		storage zitadel.Storage
		cmds    []zitadel.Command
	}{
		{
			name:    "crdb1 1 event no payload",
			storage: crdb1,
			cmds: []zitadel.Command{
				&testCommand{},
			},
		},
		{
			name:    "crdb1 2 events no payload",
			storage: crdb1,
			cmds: []zitadel.Command{
				&testCommand{},
				&testCommand{},
			},
		},
		{
			name:    "crdb1 3 events no payload",
			storage: crdb1,
			cmds: []zitadel.Command{
				&testCommand{},
				&testCommand{},
				&testCommand{},
			},
		},
		{
			name:    "crdb1 4 events no payload",
			storage: crdb1,
			cmds: []zitadel.Command{
				&testCommand{},
				&testCommand{},
				&testCommand{},
				&testCommand{},
			},
		},
		{
			name:    "crdb1 5 event no payload",
			storage: crdb1,
			cmds: []zitadel.Command{
				&testCommand{},
				&testCommand{},
				&testCommand{},
				&testCommand{},
				&testCommand{},
			},
		},
		{
			name:    "crdb1 10 event no payload",
			storage: crdb1,
			cmds: []zitadel.Command{
				&testCommand{},
				&testCommand{},
				&testCommand{},
				&testCommand{},
				&testCommand{},
				&testCommand{},
				&testCommand{},
				&testCommand{},
				&testCommand{},
				&testCommand{},
			},
		},
		{
			name:    "crdb2 1 event no payload",
			storage: crdb2,
			cmds: []zitadel.Command{
				&testCommand{},
			},
		},
		{
			name:    "crdb2 2 events no payload",
			storage: crdb2,
			cmds: []zitadel.Command{
				&testCommand{},
				&testCommand{},
			},
		},
		{
			name:    "crdb2 3 events no payload",
			storage: crdb2,
			cmds: []zitadel.Command{
				&testCommand{},
				&testCommand{},
				&testCommand{},
			},
		},
		{
			name:    "crdb2 4 events no payload",
			storage: crdb2,
			cmds: []zitadel.Command{
				&testCommand{},
				&testCommand{},
				&testCommand{},
				&testCommand{},
			},
		},
		{
			name:    "crdb2 5 event no payload",
			storage: crdb2,
			cmds: []zitadel.Command{
				&testCommand{},
				&testCommand{},
				&testCommand{},
				&testCommand{},
				&testCommand{},
			},
		},
		{
			name:    "crdb2 10 event no payload",
			storage: crdb2,
			cmds: []zitadel.Command{
				&testCommand{},
				&testCommand{},
				&testCommand{},
				&testCommand{},
				&testCommand{},
				&testCommand{},
				&testCommand{},
				&testCommand{},
				&testCommand{},
				&testCommand{},
			},
		},
		{
			name:    "crdb2_2 1 event no payload",
			storage: crdb2_2,
			cmds: []zitadel.Command{
				&testCommand{},
			},
		},
		{
			name:    "crdb2_2 2 events no payload",
			storage: crdb2_2,
			cmds: []zitadel.Command{
				&testCommand{},
				&testCommand{},
			},
		},
		{
			name:    "crdb2_2 3 events no payload",
			storage: crdb2_2,
			cmds: []zitadel.Command{
				&testCommand{},
				&testCommand{},
				&testCommand{},
			},
		},
		{
			name:    "crdb2_2 4 events no payload",
			storage: crdb2_2,
			cmds: []zitadel.Command{
				&testCommand{},
				&testCommand{},
				&testCommand{},
				&testCommand{},
			},
		},
		{
			name:    "crdb2_2 5 event no payload",
			storage: crdb2_2,
			cmds: []zitadel.Command{
				&testCommand{},
				&testCommand{},
				&testCommand{},
				&testCommand{},
				&testCommand{},
			},
		},
		{
			name:    "crdb2_2 10 event no payload",
			storage: crdb2_2,
			cmds: []zitadel.Command{
				&testCommand{},
				&testCommand{},
				&testCommand{},
				&testCommand{},
				&testCommand{},
				&testCommand{},
				&testCommand{},
				&testCommand{},
				&testCommand{},
				&testCommand{},
			},
		},
	}
	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			es := zitadel.NewEventstore(tt.storage)
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

func (cmd *testCommand) EditorService() string {
	return "svc"
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
