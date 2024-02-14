package client

import (
	"encoding/base64"
	"reflect"
	"strings"
	"testing"

	eventstorev1alpha "github.com/adlerhurst/eventstore/service/api/adlerhurst/eventstore/v1alpha"
	"google.golang.org/protobuf/types/known/structpb"
)

func Test_listArgIndexes(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name        string
		args        args
		wantIndexes []int
	}{
		{
			name: "--aggregates --id 123",
			args: args{
				name: "aggregates",
			},
			wantIndexes: []int{0},
		},
		{
			name: "--aggregates --id --aggregates --aggregates",
			args: args{
				name: "aggregates",
			},
			wantIndexes: []int{0, 2, 3},
		},
		{
			name: "nope",
			args: args{
				name: "aggregates",
			},
			wantIndexes: nil,
		},
		{
			name: "--aggreGates",
			args: args{
				name: "aggregates",
			},
			wantIndexes: []int{0},
		},
		{
			name: "-aggregates",
			args: args{
				name: "aggregates",
			},
			wantIndexes: []int{0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotIndexes := listArgIndexes(tt.args.name, strings.Split(tt.name, " ")...); !reflect.DeepEqual(gotIndexes, tt.wantIndexes) {
				t.Errorf("listArgIndexes() = %v, want %v", gotIndexes, tt.wantIndexes)
			}
		})
	}
}

func Test_parseCommand(t *testing.T) {
	addedPayload := base64.StdEncoding.EncodeToString([]byte(`{"username": "adlerhurst"}`))
	tests := []struct {
		name        string
		wantCommand *eventstorev1alpha.Command
		wantErr     bool
	}{
		{
			name: "--revision=1",
			wantCommand: &eventstorev1alpha.Command{
				Action: &eventstorev1alpha.Action{
					Action:   nil,
					Revision: 1,
					Payload:  nil,
				},
			},
			wantErr: false,
		},
		{
			name: "--revision=1 --action user,1 --action removed",
			wantCommand: &eventstorev1alpha.Command{
				Action: &eventstorev1alpha.Action{
					Action:   []string{"user", "1", "removed"},
					Revision: 1,
					Payload:  nil,
				},
			},
			wantErr: false,
		},
		{
			name: "--revision=1 --action user,1,added --payload=" + addedPayload,
			wantCommand: &eventstorev1alpha.Command{
				Action: &eventstorev1alpha.Action{
					Action:   []string{"user", "1", "added"},
					Revision: 1,
					Payload: &structpb.Struct{
						Fields: map[string]*structpb.Value{
							"username": structpb.NewStringValue("adlerhurst"),
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "empty",
			wantCommand: &eventstorev1alpha.Command{
				Action: &eventstorev1alpha.Action{
					Action:   nil,
					Revision: 0,
					Payload:  nil,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCommand, err := parseCommand(strings.Split(tt.name, " "))
			if (err != nil) != tt.wantErr {
				t.Errorf("parsePushRequestAggregateCommand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assertCommand(t, gotCommand, tt.wantCommand)
		})
	}
}

func Test_parseAggregate(t *testing.T) {
	sequence := uint32(5)
	addedPayload := base64.StdEncoding.EncodeToString([]byte(`{"username": "adlerhurst"}`))
	tests := []struct {
		name          string
		wantAggregate *eventstorev1alpha.Aggregate
		wantErr       bool
	}{
		{
			name:          "--id",
			wantAggregate: nil,
		},
		{
			name:          "--id user,1 --currentsequence 5",
			wantAggregate: nil,
		},
		{
			name: "--id user,1 --currentsequence 5 " +
				"--commands --action user,1,removed --revision 5",
			wantAggregate: &eventstorev1alpha.Aggregate{
				Id:              []string{"user", "1"},
				CurrentSequence: &sequence,
				Commands: []*eventstorev1alpha.Command{
					{
						Action: &eventstorev1alpha.Action{
							Action:   []string{"user", "1", "removed"},
							Revision: 5,
							Payload: &structpb.Struct{
								Fields: map[string]*structpb.Value{},
							},
						},
					},
				},
			},
		},
		{
			name: "--id user,1 --currentsequence 5" +
				" --commands --action user,1,added --revision 5 --payload " + addedPayload +
				" --commands --action user,1,removed --revision 5",
			wantAggregate: &eventstorev1alpha.Aggregate{
				Id:              []string{"user", "1"},
				CurrentSequence: &sequence,
				Commands: []*eventstorev1alpha.Command{
					{
						Action: &eventstorev1alpha.Action{
							Action:   []string{"user", "1", "added"},
							Revision: 5,
							Payload: &structpb.Struct{
								Fields: map[string]*structpb.Value{
									"username": structpb.NewStringValue("adlerhurst"),
								},
							},
						},
					},
					{
						Action: &eventstorev1alpha.Action{
							Action:   []string{"user", "1", "removed"},
							Revision: 5,
							Payload: &structpb.Struct{
								Fields: map[string]*structpb.Value{},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAggregate, err := parseAggregate(strings.Split(tt.name, " "))
			if (err != nil) != tt.wantErr {
				t.Errorf("parseAggregate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if gotAggregate == nil && tt.wantAggregate != nil {
				t.Errorf("gotAggregate was nil")
				return
			}

			assertAggregate(t, gotAggregate, tt.wantAggregate)
		})
	}
}

func assertAggregate(t *testing.T, got, want *eventstorev1alpha.Aggregate) {
	t.Helper()

	if got == nil && want == nil {
		return
	}

	if !reflect.DeepEqual(got.Id, want.Id) {
		t.Errorf("unexpected aggregate.Id: want %v, got %v", want.Id, got.Id)
	}

	if !reflect.DeepEqual(got.CurrentSequence, want.CurrentSequence) {
		t.Errorf("unexpected aggregate.CurrentSequence: want %v, got %v", *want.CurrentSequence, *got.CurrentSequence)
	}

	if len(got.Commands) != len(want.Commands) {
		t.Errorf("unexpected amount of aggregate.Commands: want %v, got %v", len(want.Commands), len(got.Commands))
		return
	}

	for i, gotCommand := range got.Commands {
		assertCommand(t, gotCommand, want.Commands[i])
	}
}

func assertCommand(t *testing.T, got, want *eventstorev1alpha.Command) {
	t.Helper()
	assertAction(t, got.Action, want.Action)
}

func assertAction(t *testing.T, got, want *eventstorev1alpha.Action) {
	t.Helper()
	if !reflect.DeepEqual(got.Action, want.Action) {
		t.Errorf("unexpected action.Action: want %v, got %v", want.Action, got.Action)
	}
	if got.Revision != want.Revision {
		t.Errorf("unexpected action.Revision: want %d, got %d", want.Revision, got.Revision)
	}
	if !reflect.DeepEqual(got.Payload.AsMap(), want.Payload.AsMap()) {
		t.Errorf("unexpected action.Payload: want %v, got %v", want.Payload.AsMap(), got.Payload.AsMap())
	}

}
