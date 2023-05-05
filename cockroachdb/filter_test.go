package cockroachdb

import (
	"context"
	_ "embed"
	"reflect"
	"testing"
	"time"

	"github.com/adlerhurst/eventstore/v0"
)

func Test_Filter_Compliance(t *testing.T) {
	eventstore.FilterComplianceTests(context.Background(), t, store)
}

func Test_filterToWhere(t *testing.T) {
	timeArg := time.Now()
	type args struct {
		filter *eventstore.Filter
	}
	tests := []struct {
		name      string
		args      args
		wantWhere string
		wantArgs  []interface{}
	}{
		{
			name: "no filter",
			args: args{
				filter: nil,
			},
			wantWhere: "",
			wantArgs:  nil,
		},
		{
			name: "empty filter",
			args: args{
				filter: &eventstore.Filter{},
			},
			wantWhere: "",
			wantArgs:  []interface{}{},
		},
		{
			name: "filter to",
			args: args{
				filter: &eventstore.Filter{
					CreatedAt: eventstore.CreatedAtFilter{
						To: timeArg,
					},
				},
			},
			wantWhere: "WHERE created_at <= $1",
			wantArgs:  []interface{}{timeArg},
		},
		{
			name: "filter from",
			args: args{
				filter: &eventstore.Filter{
					CreatedAt: eventstore.CreatedAtFilter{
						From: timeArg,
					},
				},
			},
			wantWhere: "WHERE created_at >= $1",
			wantArgs:  []interface{}{timeArg},
		},
		{
			name: "filter limit",
			args: args{
				filter: &eventstore.Filter{Limit: 10},
			},
			wantWhere: "",
			wantArgs:  []interface{}{},
		},
		{
			name: "filter action exact",
			args: args{
				filter: &eventstore.Filter{
					Action: []eventstore.Subject{eventstore.TextSubject("test")},
				},
			},
			wantWhere: "WHERE \"action\"[1] = $1 AND cardinality(\"action\") = $2",
			wantArgs:  []interface{}{eventstore.TextSubject("test"), 1},
		},
		{
			name: "filter action with single at end",
			args: args{
				filter: &eventstore.Filter{
					Action: []eventstore.Subject{
						eventstore.TextSubject("test"),
						eventstore.SingleToken,
					},
				},
			},
			wantWhere: "WHERE \"action\"[1] = $1 AND cardinality(\"action\") = $2",
			wantArgs:  []interface{}{eventstore.TextSubject("test"), 2},
		},
		{
			name: "filter action with single at beginning",
			args: args{
				filter: &eventstore.Filter{
					Action: []eventstore.Subject{
						eventstore.SingleToken,
						eventstore.TextSubject("test"),
					},
				},
			},
			wantWhere: "WHERE \"action\"[2] = $1 AND cardinality(\"action\") = $2",
			wantArgs:  []interface{}{eventstore.TextSubject("test"), 2},
		},
		{
			name: "filter action with single between",
			args: args{
				filter: &eventstore.Filter{
					Action: []eventstore.Subject{
						eventstore.TextSubject("test"),
						eventstore.SingleToken,
						eventstore.TextSubject("added"),
					},
				},
			},
			wantWhere: "WHERE \"action\"[1] = $1 AND \"action\"[3] = $2 AND cardinality(\"action\") = $3",
			wantArgs:  []interface{}{eventstore.TextSubject("test"), eventstore.TextSubject("added"), 3},
		},
		{
			name: "filter action only multi",
			args: args{
				filter: &eventstore.Filter{
					Action: []eventstore.Subject{
						eventstore.MultiToken,
					},
				},
			},
			wantWhere: "WHERE cardinality(\"action\") >= $1",
			wantArgs:  []interface{}{0},
		},
		{
			name: "filter action multi at end",
			args: args{
				filter: &eventstore.Filter{
					Action: []eventstore.Subject{
						eventstore.TextSubject("test"),
						eventstore.MultiToken,
					},
				},
			},
			wantWhere: "WHERE \"action\"[1] = $1 AND cardinality(\"action\") >= $2",
			wantArgs:  []interface{}{eventstore.TextSubject("test"), 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotWhere, gotArgs := filterToClauses(tt.args.filter)
			if reflect.DeepEqual(gotWhere, tt.wantWhere) {
				t.Errorf("filterToWhere() gotWhere = %v, want %v", gotWhere, tt.wantWhere)
			}
			if !reflect.DeepEqual(gotArgs, tt.wantArgs) {
				t.Errorf("filterToWhere() gotArgs = %v, want %v", gotArgs, tt.wantArgs)
			}
		})
	}
}
