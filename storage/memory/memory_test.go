package memory

import (
	"testing"

	"github.com/adlerhurst/eventstore"
)

func Test_matchSubject(t *testing.T) {
	type args struct {
		event  *eventstore.Event
		filter []eventstore.Subject
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "filter less exact than event",
			args: args{
				event: &eventstore.Event{
					Subjects: []eventstore.TextSubject{"user", "123"},
				},
				filter: []eventstore.Subject{eventstore.TextSubject("user")},
			},
			want: false,
		},
		{
			name: "filter more exact than event",
			args: args{
				event: &eventstore.Event{
					Subjects: []eventstore.TextSubject{"user"},
				},
				filter: []eventstore.Subject{
					eventstore.TextSubject("user"),
					eventstore.TextSubject("123"),
				},
			},
			want: false,
		},
		{
			name: "exact match",
			args: args{
				event: &eventstore.Event{
					Subjects: []eventstore.TextSubject{"user", "123"},
				},
				filter: []eventstore.Subject{
					eventstore.TextSubject("user"),
					eventstore.TextSubject("123"),
				},
			},
			want: true,
		},
		{
			name: "wrong text",
			args: args{
				event: &eventstore.Event{
					Subjects: []eventstore.TextSubject{"user", "321"},
				},
				filter: []eventstore.Subject{
					eventstore.TextSubject("user"),
					eventstore.TextSubject("123"),
				},
			},
			want: false,
		},
		{
			name: "single token match end",
			args: args{
				event: &eventstore.Event{
					Subjects: []eventstore.TextSubject{"user", "123"},
				},
				filter: []eventstore.Subject{
					eventstore.TextSubject("user"),
					eventstore.SingleToken,
				},
			},
			want: true,
		},
		{
			name: "single token match middle",
			args: args{
				event: &eventstore.Event{
					Subjects: []eventstore.TextSubject{"user", "added", "123"},
				},
				filter: []eventstore.Subject{
					eventstore.TextSubject("user"),
					eventstore.SingleToken,
					eventstore.TextSubject("123"),
				},
			},
			want: true,
		},
		{
			name: "single token match beginning",
			args: args{
				event: &eventstore.Event{
					Subjects: []eventstore.TextSubject{"user", "123"},
				},
				filter: []eventstore.Subject{
					eventstore.SingleToken,
					eventstore.TextSubject("123"),
				},
			},
			want: true,
		},
		{
			name: "single token to many in filter",
			args: args{
				event: &eventstore.Event{
					Subjects: []eventstore.TextSubject{"user", "123"},
				},
				filter: []eventstore.Subject{
					eventstore.TextSubject("user"),
					eventstore.TextSubject("123"),
					eventstore.SingleToken,
				},
			},
			want: false,
		},
		{
			name: "single token too few in filter",
			args: args{
				event: &eventstore.Event{
					Subjects: []eventstore.TextSubject{"user", "123", "added"},
				},
				filter: []eventstore.Subject{
					eventstore.TextSubject("user"),
					eventstore.SingleToken,
				},
			},
			want: false,
		},
		{
			name: "multi token subject",
			args: args{
				event: &eventstore.Event{
					Subjects: []eventstore.TextSubject{"user", "123"},
				},
				filter: []eventstore.Subject{
					eventstore.MultiToken,
				},
			},
			want: true,
		},
		{
			name: "multi token subject to few event sujects",
			args: args{
				event: &eventstore.Event{
					Subjects: []eventstore.TextSubject{"user", "123"},
				},
				filter: []eventstore.Subject{
					eventstore.TextSubject("user"),
					eventstore.TextSubject("123"),
					eventstore.MultiToken,
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := matchSubjects(tt.args.event, tt.args.filter); got != tt.want {
				t.Errorf("matchSubject() = %v, want %v", got, tt.want)
			}
		})
	}
}
