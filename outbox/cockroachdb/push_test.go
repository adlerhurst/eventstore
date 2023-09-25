package cockroachdb

import (
	"context"
	_ "embed"
	"testing"

	"github.com/adlerhurst/eventstore/v0"
)

func Benchmark_Push_ParallelSameAggregate(b *testing.B) {
	b.Run("Benchmark_Push_ParallelSameAggregate", func(b *testing.B) {
		eventstore.PushParallelOnSameAggregate(context.Background(), b, store)
	})
}

func Benchmark_Push_ParallelDifferentAggregate(b *testing.B) {
	store.Subscribe("all", []eventstore.Subject{eventstore.MultiToken})
	store.Subscribe("all_user_added", []eventstore.Subject{eventstore.TextSubject("user"), eventstore.SingleToken, eventstore.TextSubject("added")})
	store.Subscribe("never", []eventstore.Subject{eventstore.TextSubject("project")})
	store.Subscribe("never2", []eventstore.Subject{eventstore.SingleToken, eventstore.TextSubject("never")})
	b.Run("Benchmark_Push_ParallelDifferentAggregate", func(b *testing.B) {
		eventstore.PushParallelOnDifferentAggregates(context.Background(), b, store)
	})
}

func Test_Push_Compliance(t *testing.T) {
	store.Subscribe("test", []eventstore.Subject{eventstore.MultiToken})
	store.Subscribe("test2", []eventstore.Subject{eventstore.TextSubject("user"), eventstore.SingleToken, eventstore.TextSubject("added")})
	eventstore.PushComplianceTests(context.Background(), t, store)
}

func Test_indexes_toValues(t *testing.T) {
	type args struct {
		indexes *aggregateIndexes
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1",
			args: args{
				indexes: &aggregateIndexes{
					aggregates:   make([]*aggregateIndex, 1),
					commandCount: 1,
				},
			},
			want: `($2, array_to_string($2, ':'), $3, $4, $5, $6, (SELECT IF (EXISTS (SELECT seq FROM current_sequences WHERE "aggregate" = $2), (SELECT seq FROM current_sequences WHERE "aggregate" = $2), 0)) + $7)`,
		},
		{
			name: "2",
			args: args{
				indexes: &aggregateIndexes{
					aggregates:   make([]*aggregateIndex, 2),
					commandCount: 2,
				},
			},
			want: `($3, array_to_string($3, ':'), $4, $5, $6, $7, (SELECT IF (EXISTS (SELECT seq FROM current_sequences WHERE "aggregate" = $3), (SELECT seq FROM current_sequences WHERE "aggregate" = $3), 0)) + $8), ($9, array_to_string($9, ':'), $10, $11, $12, $13, (SELECT IF (EXISTS (SELECT seq FROM current_sequences WHERE "aggregate" = $9), (SELECT seq FROM current_sequences WHERE "aggregate" = $9), 0)) + $14)`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.indexes.toValues(); got != tt.want {
				t.Errorf("eventToValue() = \n%q\n, want \n%q\n", got, tt.want)
			}
		})
	}
}
