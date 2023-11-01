package cockroachdb

import (
	"context"
	_ "embed"
	"reflect"
	"strings"
	"testing"

	"github.com/adlerhurst/eventstore/v2"
)

func Benchmark_Push_ParallelSameAggregate(b *testing.B) {
	b.Run("Benchmark_Push_ParallelSameAggregate", func(b *testing.B) {
		eventstore.PushParallelOnSameAggregate(context.Background(), b, store)
	})
}

func Benchmark_Push_ParallelDifferentAggregate(b *testing.B) {
	b.Run("Benchmark_Push_ParallelDifferentAggregate", func(b *testing.B) {
		eventstore.PushParallelOnDifferentAggregates(context.Background(), b, store)
	})
}

func Test_Push_Compliance(t *testing.T) {
	eventstore.PushComplianceTests(context.Background(), t, store)
}

func Test_indexes_eventValues(t *testing.T) {
	type args struct {
		aggregates []eventstore.TextSubjects
		commands   []*command
	}
	type want struct {
		values string
		args   []any
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "1 command",
			args: args{
				aggregates: []eventstore.TextSubjects{{"user", "1"}},
				commands: []*command{
					{
						aggregate: eventstore.TextSubjects{"user", "1"},
						payload:   nil,
						Command: &testCommand{
							testAction: &testAction{
								action:   eventstore.TextSubjects{"user", "1", "added"},
								revision: 1,
							},
						},
					},
				},
			},
			want: want{
				values: "($1::TEXT[],$2::TEXT[],$3::INT2,$4::JSONB,$5::INT4,$6::INT4)",
				args: []any{
					eventstore.TextSubjects{"user", "1"},
					eventstore.TextSubjects{"user", "1", "added"},
					uint16(1),
					[]byte(nil),
					uint32(1),
					0,
				},
			},
		},
		{
			name: "2 commands same aggregate",
			args: args{
				aggregates: []eventstore.TextSubjects{{"user", "1"}},
				commands: []*command{
					{
						aggregate: eventstore.TextSubjects{"user", "1"},
						payload:   nil,
						Command: &testCommand{
							testAction: &testAction{
								action:   eventstore.TextSubjects{"user", "1", "added"},
								revision: 1,
							},
						},
					},
					{
						aggregate: eventstore.TextSubjects{"user", "1"},
						payload:   nil,
						Command: &testCommand{
							testAction: &testAction{
								action:   eventstore.TextSubjects{"user", "1", "changed"},
								revision: 1,
							},
						},
					},
				},
			},
			want: want{
				values: "($1::TEXT[],$2::TEXT[],$3::INT2,$4::JSONB,$5::INT4,$6::INT4),($7::TEXT[],$8::TEXT[],$9::INT2,$10::JSONB,$11::INT4,$12::INT4)",
				args: []any{
					eventstore.TextSubjects{"user", "1"},
					eventstore.TextSubjects{"user", "1", "added"},
					uint16(1),
					[]byte(nil),
					uint32(1),
					0,
					eventstore.TextSubjects{"user", "1"},
					eventstore.TextSubjects{"user", "1", "changed"},
					uint16(1),
					[]byte(nil),
					uint32(2),
					1,
				},
			},
		},
		{
			name: "3 commands 2 aggregates",
			args: args{
				aggregates: []eventstore.TextSubjects{
					{"user", "1"},
					{"user", "2"},
				},
				commands: []*command{
					{
						aggregate: eventstore.TextSubjects{"user", "1"},
						payload:   nil,
						Command: &testCommand{
							testAction: &testAction{
								action:   eventstore.TextSubjects{"user", "1", "added"},
								revision: 1,
							},
						},
					},
					{
						aggregate: eventstore.TextSubjects{"user", "2"},
						payload:   nil,
						Command: &testCommand{
							testAction: &testAction{
								action:   eventstore.TextSubjects{"user", "2", "added"},
								revision: 1,
							},
						},
					},
					{
						aggregate: eventstore.TextSubjects{"user", "2"},
						payload:   nil,
						Command: &testCommand{
							testAction: &testAction{
								action:   eventstore.TextSubjects{"user", "2", "changed"},
								revision: 1,
							},
						},
					},
				},
			},
			want: want{
				values: "($1::TEXT[],$2::TEXT[],$3::INT2,$4::JSONB,$5::INT4,$6::INT4),($7::TEXT[],$8::TEXT[],$9::INT2,$10::JSONB,$11::INT4,$12::INT4),($13::TEXT[],$14::TEXT[],$15::INT2,$16::JSONB,$17::INT4,$18::INT4)",
				args: []any{
					eventstore.TextSubjects{"user", "1"},
					eventstore.TextSubjects{"user", "1", "added"},
					uint16(1),
					[]byte(nil),
					uint32(1),
					0,
					eventstore.TextSubjects{"user", "2"},
					eventstore.TextSubjects{"user", "2", "added"},
					uint16(1),
					[]byte(nil),
					uint32(1),
					1,
					eventstore.TextSubjects{"user", "2"},
					eventstore.TextSubjects{"user", "2", "changed"},
					uint16(1),
					[]byte(nil),
					uint32(2),
					2,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var builder strings.Builder
			indexes := &aggregateIndexes{
				aggregates: make([]*aggregateIndex, len(tt.args.aggregates)),
			}
			for idx, aggregate := range tt.args.aggregates {
				indexes.aggregates[idx] = &aggregateIndex{
					aggregate: aggregate,
				}
			}
			if got := indexes.eventValues(tt.args.commands, &builder); !reflect.DeepEqual(got, tt.want.args) {
				t.Errorf("eventToValue()=\n%q,\nwant\n%q\n", got, tt.want.args)
			}

			if got := builder.String(); got != tt.want.values {
				t.Errorf("unexpected stmt:\n%q,\nwant\n%q\n", got, tt.want.values)
			}
		})
	}
}

func Benchmark_indexes_eventValues(b *testing.B) {
	type args struct {
		aggregates []eventstore.TextSubjects
		commands   []*command
	}
	type want struct {
		values string
		args   []any
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "1 command",
			args: args{
				aggregates: []eventstore.TextSubjects{{"user", "1"}},
				commands: []*command{
					{
						aggregate: eventstore.TextSubjects{"user", "1"},
						payload:   nil,
						Command: &testCommand{
							testAction: &testAction{
								action:   eventstore.TextSubjects{"user", "1", "added"},
								revision: 1,
							},
						},
					},
				},
			},
			want: want{
				values: "($1::TEXT[],$2::TEXT[],$3::INT2,$4::JSONB,$5::INT4,$6::INT4)",
				args: []any{
					eventstore.TextSubjects{"user", "1"},
					eventstore.TextSubjects{"user", "1", "added"},
					uint16(1),
					[]byte(nil),
					uint32(1),
					0,
				},
			},
		},
		{
			name: "2 commands same aggregate",
			args: args{
				aggregates: []eventstore.TextSubjects{{"user", "1"}},
				commands: []*command{
					{
						aggregate: eventstore.TextSubjects{"user", "1"},
						payload:   nil,
						Command: &testCommand{
							testAction: &testAction{
								action:   eventstore.TextSubjects{"user", "1", "added"},
								revision: 1,
							},
						},
					},
					{
						aggregate: eventstore.TextSubjects{"user", "1"},
						payload:   nil,
						Command: &testCommand{
							testAction: &testAction{
								action:   eventstore.TextSubjects{"user", "1", "changed"},
								revision: 1,
							},
						},
					},
				},
			},
			want: want{
				values: "($1::TEXT[],$2::TEXT[],$3::INT2,$4::JSONB,$5::INT4,$6::INT4),($7::TEXT[],$8::TEXT[],$9::INT2,$10::JSONB,$11::INT4,$12::INT4)",
				args: []any{
					eventstore.TextSubjects{"user", "1"},
					eventstore.TextSubjects{"user", "1", "added"},
					uint16(1),
					[]byte(nil),
					uint32(1),
					0,
					eventstore.TextSubjects{"user", "1"},
					eventstore.TextSubjects{"user", "1", "changed"},
					uint16(1),
					[]byte(nil),
					uint32(2),
					1,
				},
			},
		},
		{
			name: "3 commands 2 aggregates",
			args: args{
				aggregates: []eventstore.TextSubjects{
					{"user", "1"},
					{"user", "2"},
				},
				commands: []*command{
					{
						aggregate: eventstore.TextSubjects{"user", "1"},
						payload:   nil,
						Command: &testCommand{
							testAction: &testAction{
								action:   eventstore.TextSubjects{"user", "1", "added"},
								revision: 1,
							},
						},
					},
					{
						aggregate: eventstore.TextSubjects{"user", "2"},
						payload:   nil,
						Command: &testCommand{
							testAction: &testAction{
								action:   eventstore.TextSubjects{"user", "2", "added"},
								revision: 1,
							},
						},
					},
					{
						aggregate: eventstore.TextSubjects{"user", "2"},
						payload:   nil,
						Command: &testCommand{
							testAction: &testAction{
								action:   eventstore.TextSubjects{"user", "2", "changed"},
								revision: 1,
							},
						},
					},
				},
			},
			want: want{
				values: "($1::TEXT[],$2::TEXT[],$3::INT2,$4::JSONB,$5::INT4,$6::INT4),($7::TEXT[],$8::TEXT[],$9::INT2,$10::JSONB,$11::INT4,$12::INT4),($13::TEXT[],$14::TEXT[],$15::INT2,$16::JSONB,$17::INT4,$18::INT4)",
				args: []any{
					eventstore.TextSubjects{"user", "1"},
					eventstore.TextSubjects{"user", "1", "added"},
					uint16(1),
					[]byte(nil),
					uint32(1),
					0,
					eventstore.TextSubjects{"user", "2"},
					eventstore.TextSubjects{"user", "2", "added"},
					uint16(1),
					[]byte(nil),
					uint32(1),
					1,
					eventstore.TextSubjects{"user", "2"},
					eventstore.TextSubjects{"user", "2", "changed"},
					uint16(1),
					[]byte(nil),
					uint32(2),
					2,
				},
			},
		},
	}
	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				indexes := &aggregateIndexes{
					aggregates: make([]*aggregateIndex, len(tt.args.aggregates)),
				}
				for idx, aggregate := range tt.args.aggregates {
					indexes.aggregates[idx] = &aggregateIndex{
						aggregate: aggregate,
					}
				}
				var builder strings.Builder
				if got := indexes.eventValues(tt.args.commands, &builder); !reflect.DeepEqual(got, tt.want.args) {
					b.Errorf("eventToValue()=\n%q,\nwant\n%q\n", got, tt.want.args)
				}

				if got := builder.String(); got != tt.want.values {
					b.Errorf("unexpected stmt:\n%q,\nwant\n%q\n", got, tt.want.values)
				}
			}
		})
	}
}

func Test_actionValues(t *testing.T) {
	type field struct {
		id     string
		action eventstore.TextSubjects
	}
	type args struct {
		fields []*field
	}
	type want struct {
		values string
		args   []any
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "1 command",
			args: args{
				fields: []*field{
					{
						id:     "1",
						action: eventstore.TextSubjects{"user", "1", "added"},
					},
				},
			},
			want: want{
				values: "($1::UUID,$2::TEXT,$3::INT2),($4::UUID,$5::TEXT,$6::INT2),($7::UUID,$8::TEXT,$9::INT2)",
				args: []any{
					"1", eventstore.TextSubject("user"), 0,
					"1", eventstore.TextSubject("1"), 1,
					"1", eventstore.TextSubject("added"), 2,
				},
			},
		},
		{
			name: "2 commands",
			args: args{
				fields: []*field{
					{
						id:     "1",
						action: eventstore.TextSubjects{"user", "1", "added"},
					},
					{
						id:     "2",
						action: eventstore.TextSubjects{"event", "created"},
					},
				},
			},
			want: want{
				values: "($1::UUID,$2::TEXT,$3::INT2),($4::UUID,$5::TEXT,$6::INT2),($7::UUID,$8::TEXT,$9::INT2),($10::UUID,$11::TEXT,$12::INT2),($13::UUID,$14::TEXT,$15::INT2)",
				args: []any{
					"1", eventstore.TextSubject("user"), 0,
					"1", eventstore.TextSubject("1"), 1,
					"1", eventstore.TextSubject("added"), 2,
					"2", eventstore.TextSubject("event"), 0,
					"2", eventstore.TextSubject("created"), 1,
				},
			},
		},
	}
	for _, tt := range tests {
		commands := make([]*command, len(tt.args.fields))
		for i, field := range tt.args.fields {
			commands[i] = &command{
				id: field.id,
				Command: &testCommand{
					testAction: &testAction{
						action: field.action,
					},
				},
			}
		}
		var builder strings.Builder
		t.Run(tt.name, func(t *testing.T) {
			if got := actionValues(commands, &builder); !reflect.DeepEqual(got, tt.want.args) {
				t.Errorf("actionValues() = \n%q\nwant\n%q\n", got, tt.want.args)
			}

			if got := builder.String(); got != tt.want.values {
				t.Errorf("unexpected stmt:\n%q\nwant\n%q\n", got, tt.want.values)
			}
		})
	}
}

func Benchmark_actionValues(b *testing.B) {
	type field struct {
		id     string
		action eventstore.TextSubjects
	}
	type args struct {
		fields []*field
	}
	type want struct {
		values string
		args   []any
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "1 command",
			args: args{
				fields: []*field{
					{
						id:     "1",
						action: eventstore.TextSubjects{"user", "1", "added"},
					},
				},
			},
			want: want{
				values: "($1::UUID,$2::TEXT,$3::INT2),($4::UUID,$5::TEXT,$6::INT2),($7::UUID,$8::TEXT,$9::INT2)",
				args: []any{
					"1", eventstore.TextSubject("user"), 0,
					"1", eventstore.TextSubject("1"), 1,
					"1", eventstore.TextSubject("added"), 2,
				},
			},
		},
		{
			name: "2 commands",
			args: args{
				fields: []*field{
					{
						id:     "1",
						action: eventstore.TextSubjects{"user", "1", "added"},
					},
					{
						id:     "2",
						action: eventstore.TextSubjects{"event", "created"},
					},
				},
			},
			want: want{
				values: "($1::UUID,$2::TEXT,$3::INT2),($4::UUID,$5::TEXT,$6::INT2),($7::UUID,$8::TEXT,$9::INT2),($10::UUID,$11::TEXT,$12::INT2),($13::UUID,$14::TEXT,$15::INT2)",
				args: []any{
					"1", eventstore.TextSubject("user"), 0,
					"1", eventstore.TextSubject("1"), 1,
					"1", eventstore.TextSubject("added"), 2,
					"2", eventstore.TextSubject("event"), 0,
					"2", eventstore.TextSubject("created"), 1,
				},
			},
		},
	}
	for _, tt := range tests {
		commands := make([]*command, len(tt.args.fields))
		for i, field := range tt.args.fields {
			commands[i] = &command{
				id: field.id,
				Command: &testCommand{
					testAction: &testAction{
						action: field.action,
					},
				},
			}
		}
		b.Run(tt.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				var builder strings.Builder
				if got := actionValues(commands, &builder); !reflect.DeepEqual(got, tt.want.args) {
					b.Errorf("actionValues() = %v, want %v\n", got, tt.want.args)
				}

				if got := builder.String(); got != tt.want.values {
					b.Errorf("unexpected stmt:\n%q\nwant\n%q\n", got, tt.want.values)
				}
			}
		})
	}
}
