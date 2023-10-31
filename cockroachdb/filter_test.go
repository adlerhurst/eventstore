package cockroachdb

import (
	"context"
	_ "embed"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/adlerhurst/eventstore/v1"
)

func Benchmark_Filter(b *testing.B) {
	b.Run("Benchmark_Filter", func(b *testing.B) {
		eventstore.FilterBenchTests(context.Background(), b, store)
	})
}

func Test_Filter_Compliance(t *testing.T) {
	eventstore.FilterComplianceTests(context.Background(), t, store)
}

func Test_textSubjectClause(t *testing.T) {
	type args struct {
		index      int
		tableAlias string
		subject    eventstore.TextSubject
	}
	type want struct {
		query string
		index int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "first index",
			args: args{
				index:      0,
				tableAlias: "alias",
				subject:    "user",
			},
			want: want{
				query: "alias.action = $1 AND alias.depth = $2",
				index: 2,
			},
		},
		{
			name: "second index",
			args: args{
				index:      2,
				tableAlias: "alias",
				subject:    "user",
			},
			want: want{
				query: "alias.action = $3 AND alias.depth = $4",
				index: 4,
			},
		},
	}
	for _, tt := range tests {
		var builder strings.Builder
		t.Run(tt.name, func(t *testing.T) {
			textSubjectClause(&builder, &tt.args.index, tt.args.tableAlias, tt.args.subject)

			if got := builder.String(); got != tt.want.query {
				t.Errorf("unexpected query want:\n%q\ngot:\n%q", tt.want.query, got)
			}

			if tt.want.index != tt.args.index {
				t.Errorf("unexpected index: want %d, got: %d", tt.want.index, tt.args.index)
			}
		})
	}
}

func Test_queriesToClause(t *testing.T) {
	createdAtFrom := time.Now()
	createdAtTo := time.Now().Add(10 * time.Second)
	type args struct {
		index   int
		queries []*eventstore.FilterQuery
	}
	type want struct {
		query string
		index int
		args  []any
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "empty query",
			args: args{
				index:   0,
				queries: []*eventstore.FilterQuery{},
			},
			want: want{
				query: "",
				args:  nil,
				index: 0,
			},
		},
		{
			name: "1 subject",
			args: args{
				index: 0,
				queries: []*eventstore.FilterQuery{
					{
						Subjects: []eventstore.Subject{
							eventstore.TextSubject("user"),
						},
					},
				},
			},
			want: want{
				query: "(e.id IN (SELECT a.event FROM eventstore.actions a WHERE a.action = $1 AND a.depth = $2) AND e.action_depth = $3)",
				args: []any{
					eventstore.TextSubject("user"),
					0,
					1,
				},
				index: 3,
			},
		},
		{
			name: "2 subjects",
			args: args{
				index: 0,
				queries: []*eventstore.FilterQuery{
					{
						Subjects: []eventstore.Subject{
							eventstore.TextSubject("user"),
							eventstore.TextSubject("id"),
						},
					},
				},
			},
			want: want{
				query: "(e.id IN (SELECT a.event FROM eventstore.actions a JOIN eventstore.actions a0 ON a.event = a0.event AND a0.action = $1 AND a0.depth = $2 WHERE a.action = $3 AND a.depth = $4) AND e.action_depth = $5)",
				args: []any{
					eventstore.TextSubject("id"),
					1,
					eventstore.TextSubject("user"),
					0,
					2,
				},
				index: 5,
			},
		},
		{
			name: "3 subjects",
			args: args{
				index: 0,
				queries: []*eventstore.FilterQuery{
					{
						Subjects: []eventstore.Subject{
							eventstore.TextSubject("user"),
							eventstore.TextSubject("id"),
							eventstore.TextSubject("added"),
						},
					},
				},
			},
			want: want{
				query: "(e.id IN (SELECT a.event FROM eventstore.actions a JOIN eventstore.actions a0 ON a.event = a0.event AND a0.action = $1 AND a0.depth = $2 JOIN eventstore.actions a1 ON a.event = a1.event AND a1.action = $3 AND a1.depth = $4 WHERE a.action = $5 AND a.depth = $6) AND e.action_depth = $7)",
				args: []any{
					eventstore.TextSubject("id"),
					1,
					eventstore.TextSubject("added"),
					2,
					eventstore.TextSubject("user"),
					0,
					3,
				},
				index: 7,
			},
		},
		{
			name: "2 queries",
			args: args{
				index: 0,
				queries: []*eventstore.FilterQuery{
					{
						Subjects: []eventstore.Subject{
							eventstore.TextSubject("user"),
							eventstore.TextSubject("id"),
						},
					},
					{
						Subjects: []eventstore.Subject{
							eventstore.TextSubject("user"),
							eventstore.TextSubject("id2"),
						},
					},
				},
			},
			want: want{
				query: "(e.id IN (SELECT a.event FROM eventstore.actions a JOIN eventstore.actions a0 ON a.event = a0.event AND a0.action = $1 AND a0.depth = $2 WHERE a.action = $3 AND a.depth = $4) AND e.action_depth = $5) OR (e.id IN (SELECT a.event FROM eventstore.actions a JOIN eventstore.actions a0 ON a.event = a0.event AND a0.action = $6 AND a0.depth = $7 WHERE a.action = $8 AND a.depth = $9) AND e.action_depth = $10)",
				args: []any{
					eventstore.TextSubject("id"),
					1,
					eventstore.TextSubject("user"),
					0,
					2,
					eventstore.TextSubject("id2"),
					1,
					eventstore.TextSubject("user"),
					0,
					2,
				},
				index: 10,
			},
		},
		{
			name: "single token",
			args: args{
				index: 0,
				queries: []*eventstore.FilterQuery{
					{
						Subjects: []eventstore.Subject{
							eventstore.SingleToken,
						},
					},
				},
			},
			want: want{
				query: "(e.action_depth = $1)",
				args: []any{
					1,
				},
				index: 1,
			},
		},
		{
			name: "multi token",
			args: args{
				index: 0,
				queries: []*eventstore.FilterQuery{
					{
						Subjects: []eventstore.Subject{
							eventstore.MultiToken,
						},
					},
				},
			},
			want: want{
				query: "(e.action_depth >= $1)",
				args: []any{
					1,
				},
				index: 1,
			},
		},
		{
			name: "sequence from",
			args: args{
				index: 0,
				queries: []*eventstore.FilterQuery{
					{
						Subjects: []eventstore.Subject{eventstore.MultiToken},
						Sequence: eventstore.SequenceFilter{
							From: 100,
						},
					},
				},
			},
			want: want{
				query: "(e.action_depth >= $1 AND e.sequence > $2)",
				args: []any{
					1,
					uint32(100),
				},
				index: 2,
			},
		},
		{
			name: "sequence to",
			args: args{
				index: 0,
				queries: []*eventstore.FilterQuery{
					{
						Subjects: []eventstore.Subject{eventstore.MultiToken},
						Sequence: eventstore.SequenceFilter{
							To: 100,
						},
					},
				},
			},
			want: want{
				query: "(e.action_depth >= $1 AND e.sequence < $2)",
				args: []any{
					1,
					uint32(100),
				},
				index: 2,
			},
		},
		{
			name: "sequence",
			args: args{
				index: 0,
				queries: []*eventstore.FilterQuery{
					{
						Subjects: []eventstore.Subject{eventstore.MultiToken},
						Sequence: eventstore.SequenceFilter{
							From: 100,
							To:   200,
						},
					},
				},
			},
			want: want{
				query: "(e.action_depth >= $1 AND e.sequence > $2 AND e.sequence < $3)",
				args: []any{
					1,
					uint32(100),
					uint32(200),
				},
				index: 3,
			},
		},
		{
			name: "created_at from",
			args: args{
				index: 0,
				queries: []*eventstore.FilterQuery{
					{
						Subjects: []eventstore.Subject{eventstore.MultiToken},
						CreatedAt: eventstore.CreatedAtFilter{
							From: createdAtFrom,
						},
					},
				},
			},
			want: want{
				query: "(e.action_depth >= $1 AND e.created_at > $2)",
				args: []any{
					1,
					createdAtFrom,
				},
				index: 2,
			},
		},
		{
			name: "created_at to",
			args: args{
				index: 0,
				queries: []*eventstore.FilterQuery{
					{
						Subjects: []eventstore.Subject{eventstore.MultiToken},
						CreatedAt: eventstore.CreatedAtFilter{
							To: createdAtTo,
						},
					},
				},
			},
			want: want{
				query: "(e.action_depth >= $1 AND e.created_at < $2)",
				args: []any{
					1,
					createdAtTo,
				},
				index: 2,
			},
		},
		{
			name: "created_at",
			args: args{
				index: 0,
				queries: []*eventstore.FilterQuery{
					{
						Subjects: []eventstore.Subject{eventstore.MultiToken},
						CreatedAt: eventstore.CreatedAtFilter{
							From: createdAtFrom,
							To:   createdAtTo,
						},
					},
				},
			},
			want: want{
				query: "(e.action_depth >= $1 AND e.created_at > $2 AND e.created_at < $3)",
				args: []any{
					1,
					createdAtFrom,
					createdAtTo,
				},
				index: 3,
			},
		},
	}
	for _, tt := range tests {
		var builder strings.Builder
		t.Run(tt.name, func(t *testing.T) {
			if gotArgs := queriesToClause(&builder, &tt.args.index, tt.args.queries); !reflect.DeepEqual(gotArgs, tt.want.args) {
				t.Errorf("queriesToClause() = %v, want %v", gotArgs, tt.want.args)
			}

			if got := builder.String(); got != tt.want.query {
				t.Errorf("unexpected query want:\n%q\ngot:\n%q", tt.want.query, got)
			}

			if tt.want.index != tt.args.index {
				t.Errorf("unexpected index: want %d, got: %d", tt.want.index, tt.args.index)
			}
		})
	}
}

func Test_subjectsToJoins(t *testing.T) {
	type args struct {
		index    int
		subjects []eventstore.Subject
	}
	type want struct {
		query string
		args  []any
		index int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "no subjects",
			args: args{
				index:    0,
				subjects: []eventstore.Subject{},
			},
			want: want{
				query: "",
				args:  []any{},
				index: 0,
			},
		},
		{
			name: "id",
			args: args{
				index: 0,
				subjects: []eventstore.Subject{
					eventstore.TextSubject("id"),
				},
			},
			want: want{
				query: " JOIN eventstore.actions a0 ON a.event = a0.event AND a0.action = $1 AND a0.depth = $2",
				args: []any{
					eventstore.TextSubject("id"),
					1,
				},
				index: 2,
			},
		},
		{
			name: "id.added",
			args: args{
				index: 0,
				subjects: []eventstore.Subject{
					eventstore.TextSubject("id"),
					eventstore.TextSubject("added"),
				},
			},
			want: want{
				query: " JOIN eventstore.actions a0 ON a.event = a0.event AND a0.action = $1 AND a0.depth = $2 JOIN eventstore.actions a1 ON a.event = a1.event AND a1.action = $3 AND a1.depth = $4",
				args: []any{
					eventstore.TextSubject("id"),
					1,
					eventstore.TextSubject("added"),
					2,
				},
				index: 4,
			},
		},
		{
			name: "*",
			args: args{
				index: 0,
				subjects: []eventstore.Subject{
					eventstore.SingleToken,
				},
			},
			want: want{
				query: "",
				args:  []any{},
				index: 0,
			},
		},
		{
			name: "id.*",
			args: args{
				index: 0,
				subjects: []eventstore.Subject{
					eventstore.TextSubject("id"),
					eventstore.SingleToken,
				},
			},
			want: want{
				query: " JOIN eventstore.actions a0 ON a.event = a0.event AND a0.action = $1 AND a0.depth = $2",
				args: []any{
					eventstore.TextSubject("id"),
					1,
				},
				index: 2,
			},
		},
		{
			name: "*.added",
			args: args{
				index: 0,
				subjects: []eventstore.Subject{
					eventstore.SingleToken,
					eventstore.TextSubject("added"),
				},
			},
			want: want{
				query: " JOIN eventstore.actions a1 ON a.event = a1.event AND a1.action = $1 AND a1.depth = $2",
				args: []any{
					eventstore.TextSubject("added"),
					2,
				},
				index: 2,
			},
		},
		{
			name: "id.*.set",
			args: args{
				index: 0,
				subjects: []eventstore.Subject{
					eventstore.TextSubject("id"),
					eventstore.SingleToken,
					eventstore.TextSubject("set"),
				},
			},
			want: want{
				query: " JOIN eventstore.actions a0 ON a.event = a0.event AND a0.action = $1 AND a0.depth = $2 JOIN eventstore.actions a2 ON a.event = a2.event AND a2.action = $3 AND a2.depth = $4",
				args: []any{
					eventstore.TextSubject("id"),
					1,
					eventstore.TextSubject("set"),
					3,
				},
				index: 4,
			},
		},
		{
			name: "#",
			args: args{
				index: 0,
				subjects: []eventstore.Subject{
					eventstore.MultiToken,
				},
			},
			want: want{
				query: "",
				args:  []any{},
				index: 0,
			},
		},
		{
			name: "id.#",
			args: args{
				index: 0,
				subjects: []eventstore.Subject{
					eventstore.TextSubject("id"),
					eventstore.MultiToken,
				},
			},
			want: want{
				query: " JOIN eventstore.actions a0 ON a.event = a0.event AND a0.action = $1 AND a0.depth = $2",
				args: []any{
					eventstore.TextSubject("id"),
					1,
				},
				index: 2,
			},
		},
	}
	for _, tt := range tests {
		var builder strings.Builder
		t.Run(tt.name, func(t *testing.T) {
			if got := subjectsToJoins(&builder, &tt.args.index, tt.args.subjects); !reflect.DeepEqual(got, tt.want.args) {
				t.Errorf("subjectsJoinQuery() = %v, want %v", got, tt.want.args)
			}

			if got := builder.String(); got != tt.want.query {
				t.Errorf("unexpected query want:\n%q\ngot:\n%q", tt.want.query, got)
			}

			if tt.want.index != tt.args.index {
				t.Errorf("unexpected index: want %d, got: %d", tt.want.index, tt.args.index)
			}
		})
	}
}

func Test_queryToClause(t *testing.T) {
	type args struct {
		index int
		query *eventstore.FilterQuery
	}
	type want struct {
		query string
		args  []any
		index int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "one text subject",
			args: args{
				index: 0,
				query: &eventstore.FilterQuery{
					Subjects: []eventstore.Subject{
						eventstore.TextSubject("user"),
					},
				},
			},
			want: want{
				query: "(e.id IN (SELECT a.event FROM eventstore.actions a WHERE a.action = $1 AND a.depth = $2) AND e.action_depth = $3)",
				args: []any{
					eventstore.TextSubject("user"),
					0,
					1,
				},
				index: 3,
			},
		},
		{
			name: "only single token",
			args: args{
				index: 0,
				query: &eventstore.FilterQuery{
					Subjects: []eventstore.Subject{
						eventstore.SingleToken,
					},
				},
			},
			want: want{
				query: "(e.action_depth = $1)",
				args: []any{
					1,
				},
				index: 1,
			},
		},
	}
	for _, tt := range tests {
		var builder strings.Builder
		t.Run(tt.name, func(t *testing.T) {
			if got := queryToClause(&builder, &tt.args.index, tt.args.query); !reflect.DeepEqual(got, tt.want.args) {
				t.Errorf("queryToClause() = %v, want %v", got, tt.want.args)
			}

			if got := builder.String(); got != tt.want.query {
				t.Errorf("unexpected query want:\n%q\ngot:\n%q", tt.want.query, got)
			}

			if tt.want.index != tt.args.index {
				t.Errorf("unexpected index: want %d, got: %d", tt.want.index, tt.args.index)
			}
		})
	}
}

func Test_actionDepthQuery(t *testing.T) {
	type args struct {
		index       int
		lastSubject eventstore.Subject
	}
	type want struct {
		query string
		index int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "text subject",
			args: args{
				index:       0,
				lastSubject: eventstore.TextSubject("user"),
			},
			want: want{
				query: "e.action_depth = $1",
				index: 1,
			},
		},
		{
			name: "single token",
			args: args{
				index:       0,
				lastSubject: eventstore.SingleToken,
			},
			want: want{
				query: "e.action_depth = $1",
				index: 1,
			},
		},
		{
			name: "multi token",
			args: args{
				index:       0,
				lastSubject: eventstore.MultiToken,
			},
			want: want{
				query: "e.action_depth >= $1",
				index: 1,
			},
		},
	}
	for _, tt := range tests {
		var builder strings.Builder
		t.Run(tt.name, func(t *testing.T) {
			actionDepthQuery(&builder, &tt.args.index, tt.args.lastSubject)

			if got := builder.String(); got != tt.want.query {
				t.Errorf("unexpected query want:\n%q\ngot:\n%q", tt.want.query, got)
			}

			if tt.want.index != tt.args.index {
				t.Errorf("unexpected index: want %d, got: %d", tt.want.index, tt.args.index)
			}
		})
	}
}

func Test_subjectsToClause(t *testing.T) {
	type args struct {
		index    int
		subjects []eventstore.Subject
	}
	type want struct {
		query string
		args  []any
		index int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "no subjects",
			args: args{
				index:    0,
				subjects: []eventstore.Subject{},
			},
			want: want{
				query: "",
				args:  nil,
				index: 0,
			},
		},
		{
			name: "*",
			args: args{
				index: 0,
				subjects: []eventstore.Subject{
					eventstore.SingleToken,
				},
			},
			want: want{
				query: "e.action_depth = $1",
				args: []any{
					1,
				},
				index: 1,
			},
		},
		{
			name: "#",
			args: args{
				index: 0,
				subjects: []eventstore.Subject{
					eventstore.MultiToken,
				},
			},
			want: want{
				query: "e.action_depth >= $1",
				args: []any{
					1,
				},
				index: 1,
			},
		},
		{
			name: "*.#",
			args: args{
				index: 0,
				subjects: []eventstore.Subject{
					eventstore.SingleToken,
					eventstore.MultiToken,
				},
			},
			want: want{
				query: "e.action_depth >= $1",
				args: []any{
					2,
				},
				index: 1,
			},
		},
		{
			name: "user",
			args: args{
				index: 0,
				subjects: []eventstore.Subject{
					eventstore.TextSubject("user"),
				},
			},
			want: want{
				query: "e.id IN (SELECT a.event FROM eventstore.actions a WHERE a.action = $1 AND a.depth = $2) AND e.action_depth = $3",
				args: []any{
					eventstore.TextSubject("user"),
					0,
					1,
				},
				index: 3,
			},
		},
		{
			name: "user.*",
			args: args{
				index: 0,
				subjects: []eventstore.Subject{
					eventstore.TextSubject("user"),
					eventstore.SingleToken,
				},
			},
			want: want{
				query: "e.id IN (SELECT a.event FROM eventstore.actions a WHERE a.action = $1 AND a.depth = $2) AND e.action_depth = $3",
				args: []any{
					eventstore.TextSubject("user"),
					0,
					2,
				},
				index: 3,
			},
		},
		{
			name: "user.#",
			args: args{
				index: 0,
				subjects: []eventstore.Subject{
					eventstore.TextSubject("user"),
					eventstore.MultiToken,
				},
			},
			want: want{
				query: "e.id IN (SELECT a.event FROM eventstore.actions a WHERE a.action = $1 AND a.depth = $2) AND e.action_depth >= $3",
				args: []any{
					eventstore.TextSubject("user"),
					0,
					2,
				},
				index: 3,
			},
		},
		{
			name: "user.id",
			args: args{
				index: 0,
				subjects: []eventstore.Subject{
					eventstore.TextSubject("user"),
					eventstore.TextSubject("id"),
				},
			},
			want: want{
				query: "e.id IN (SELECT a.event FROM eventstore.actions a JOIN eventstore.actions a0 ON a.event = a0.event AND a0.action = $1 AND a0.depth = $2 WHERE a.action = $3 AND a.depth = $4) AND e.action_depth = $5",
				args: []any{
					eventstore.TextSubject("id"),
					1,
					eventstore.TextSubject("user"),
					0,
					2,
				},
				index: 5,
			},
		},
		{
			name: "user.*.added",
			args: args{
				index: 0,
				subjects: []eventstore.Subject{
					eventstore.TextSubject("user"),
					eventstore.SingleToken,
					eventstore.TextSubject("added"),
				},
			},
			want: want{
				query: "e.id IN (SELECT a.event FROM eventstore.actions a JOIN eventstore.actions a1 ON a.event = a1.event AND a1.action = $1 AND a1.depth = $2 WHERE a.action = $3 AND a.depth = $4) AND e.action_depth = $5",
				args: []any{
					eventstore.TextSubject("added"),
					2,
					eventstore.TextSubject("user"),
					0,
					3,
				},
				index: 5,
			},
		},
	}
	for _, tt := range tests {
		var builder strings.Builder
		t.Run(tt.name, func(t *testing.T) {
			if got := subjectsToClause(&builder, &tt.args.index, tt.args.subjects); !reflect.DeepEqual(got, tt.want.args) {
				t.Errorf("subjectsToClause() = %v, want %v", got, tt.want.args)
			}

			if got := builder.String(); got != tt.want.query {
				t.Errorf("unexpected query want:\n%q\ngot:\n%q", tt.want.query, got)
			}

			if tt.want.index != tt.args.index {
				t.Errorf("unexpected index: want %d, got: %d", tt.want.index, tt.args.index)
			}
		})
	}
}
