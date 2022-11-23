package zitadel_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/adlerhurst/eventstore/zitadel"
	"github.com/adlerhurst/eventstore/zitadel/storage"
)

func TestEventstore_Push(t *testing.T) {
	crdb, err := storage.NewCRDB(testCRDBClient)
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
		{
			name: "push sinlge event",
			es:   zitadel.NewEventstore(crdb),
			args: args{
				cmds: []zitadel.Command{
					&testCommand{},
				},
			},
			want:    []*zitadel.Event{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.es.Push(context.Background(), tt.args.cmds)
			if (err != nil) != tt.wantErr {
				t.Errorf("Eventstore.Push() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Eventstore.Push() = %v, want %v", got, tt.want)
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
