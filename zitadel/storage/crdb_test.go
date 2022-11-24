package storage

import (
	_ "embed"
	"reflect"
	"testing"
	"time"

	"github.com/adlerhurst/eventstore/zitadel"
)

func Test_filterToSQL(t *testing.T) {
	date := time.Now()
	orgs := []string{"org1"}

	type args struct {
		filter *zitadel.Filter
	}
	tests := []struct {
		name       string
		args       args
		wantClause string
		wantArgs   []any
	}{
		{
			name: "minimum filter",
			args: args{
				filter: &zitadel.Filter{
					InstanceID: "instance",
				},
			},
			wantClause: "instance_id = $1 ORDER BY creation_date",
			wantArgs:   []any{"instance"},
		},
		{
			name: "all fields of filter",
			args: args{
				filter: &zitadel.Filter{
					InstanceID:               "instance",
					OrgIDs:                   orgs,
					CreationDateGreaterEqual: date,
					CreationDateLess:         date,
					Limit:                    10,
				},
			},
			wantClause: "AS OF SYSTEM TIME '" + date.Add(1*time.Microsecond).Format(sqlTimeLayout) + "' instance_id = $1 AND creation_date >= $2 AND org_id @> $3 ORDER BY creation_date LIMIT $4",
			wantArgs: []any{
				"instance",
				date,
				orgs,
				uint32(10),
			},
		},
		{
			name: "with aggregate",
			args: args{
				filter: &zitadel.Filter{
					InstanceID: "instance",
					Aggregates: []*zitadel.AggregateFilter{
						{
							Type: "agg1",
						},
					},
				},
			},
			wantClause: "instance_id = $1 AND (aggregate_type = $2) ORDER BY creation_date",
			wantArgs: []any{
				"instance",
				"agg1",
			},
		},
		{
			name: "with aggregates",
			args: args{
				filter: &zitadel.Filter{
					InstanceID: "instance",
					Aggregates: []*zitadel.AggregateFilter{
						{
							Type: "agg1",
						},
						{
							Type: "agg2",
						},
					},
				},
			},
			wantClause: "instance_id = $1 AND ((aggregate_type = $2) OR (aggregate_type = $3)) ORDER BY creation_date",
			wantArgs: []any{
				"instance",
				"agg1",
				"agg2",
			},
		},
		{
			name: "with aggregate and events",
			args: args{
				filter: &zitadel.Filter{
					InstanceID: "instance",
					Aggregates: []*zitadel.AggregateFilter{
						{
							Type: "user",
							Events: []*zitadel.EventFilter{
								{
									Types: []string{"user.created", "user.changed"},
								},
								{
									Types: []string{"user.removed"},
								},
							},
						},
					},
				},
			},
			wantClause: "instance_id = $1 AND (aggregate_type = $2 AND (event_type @> $3 OR event_type @> $4)) ORDER BY creation_date",
			wantArgs: []any{
				"instance",
				"user",
				[]string{"user.created", "user.changed"},
				[]string{"user.removed"},
			},
		},
		{
			name: "with aggregate and event",
			args: args{
				filter: &zitadel.Filter{
					InstanceID: "instance",
					Aggregates: []*zitadel.AggregateFilter{
						{
							Type: "user",
							Events: []*zitadel.EventFilter{
								{
									Types: []string{"user.created", "user.changed"},
								},
							},
						},
					},
				},
			},
			wantClause: "instance_id = $1 AND (aggregate_type = $2 AND event_type @> $3) ORDER BY creation_date",
			wantArgs: []any{
				"instance",
				"user",
				[]string{"user.created", "user.changed"},
			},
		},
		// {
		// 	name: "",
		// 	args: args{
		// 	filter: &zitadel.Filter{},
		// },
		// 	wantClause: "",
		// 	wantArgs: []any{},
		// },
		// {
		// 	name: "",
		// 	args: args{
		// 	filter: &zitadel.Filter{},
		// },
		// 	wantClause: "",
		// 	wantArgs: []any{},
		// },
		// {
		// 	name: "",
		// 	args: args{
		// 	filter: &zitadel.Filter{},
		// },
		// 	wantClause: "",
		// 	wantArgs: []any{},
		// },
		// {
		// 	name: "",
		// 	args: args{
		// 	filter: &zitadel.Filter{},
		// },
		// 	wantClause: "",
		// 	wantArgs: []any{},
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotClause, gotArgs := filterToSQL(tt.args.filter)
			if gotClause != tt.wantClause {
				t.Errorf("filterToSQL() \ngotClause = %q, \nwant        %q", gotClause, tt.wantClause)
			}
			if len(gotArgs) != len(tt.wantArgs) {
				t.Errorf("want len of args: %d, got %d", len(tt.wantArgs), len(gotArgs))
			}
			for i := 0; i < len(tt.wantArgs); i++ {
				if !reflect.DeepEqual(gotArgs[i], tt.wantArgs[i]) {
					t.Errorf("filterToSQL() arg %d \ngot = %v\nwant  %v", i, gotArgs[i], tt.wantArgs[i])
				}
			}
		})
	}
}
