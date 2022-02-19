package eventstore

import "testing"

func Test_isSubjectsValid(t *testing.T) {
	type args struct {
		subs []Subject
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "subjects only text",
			args: args{
				subs: []Subject{TextSubject("only"), TextSubject("text")},
			},
			want: true,
		},
		{
			name: "subjects contain single token",
			args: args{
				subs: []Subject{TextSubject("subjects"), TextSubject("contains"), SingleToken, TextSubject("token")},
			},
			want: true,
		},
		{
			name: "subjects contains single and multi token at end",
			args: args{
				subs: []Subject{TextSubject("subjects"), SingleToken, MultiToken},
			},
			want: true,
		},
		{
			name: "full wildcard at end",
			args: args{
				subs: []Subject{TextSubject("subjects"), MultiToken},
			},
			want: true,
		},
		{
			name: "full wildcard not at end",
			args: args{
				subs: []Subject{TextSubject("subjects"), MultiToken, TextSubject("end")},
			},
			want: false,
		},
		{
			name: "no subject",
			args: args{
				subs: []Subject{},
			},
			want: false,
		},
		{
			name: "text contains *",
			args: args{
				subs: []Subject{TextSubject("us*er")},
			},
			want: false,
		},
		{
			name: "text contains >",
			args: args{
				subs: []Subject{TextSubject("user"), TextSubject(">")},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isSubjectsValid(tt.args.subs); got != tt.want {
				t.Errorf("isSubjectsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}
