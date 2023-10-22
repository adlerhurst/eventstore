package x_test

import (
	"strings"
	"testing"
)

func BenchmarkStringConcat(b *testing.B) {
	var long = "asdfasfkjljsakldkjfsadjlkasdjfklasjfklajsdkfjaslkjfkas"
	type test struct {
		name string
		exec func(string) string
		res  string
	}
	for _, test := range []*test{
		{
			name: "+",
			exec: func(input string) string {
				input += "a"
				input += "b"
				input += "c"
				input += long
				return input
			},
			res: "+abcasdfasfkjljsakldkjfsadjlkasdjfklasjfklajsdkfjaslkjfkas",
		},
		{
			name: "b",
			exec: func(input string) string {
				var builder strings.Builder
				builder.WriteString(input)
				builder.WriteRune('a')
				builder.WriteRune('b')
				builder.WriteRune('c')
				builder.WriteString(long)
				return builder.String()
			},
			res: "babcasdfasfkjljsakldkjfsadjlkasdjfklasjfklajsdkfjaslkjfkas",
		},
	} {
		b.Run(test.name, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				if r := test.exec(test.name); r != test.res {
					b.Errorf("unexpected result, want %q, got: %q", test.res, r)
				}
			}
		})
	}
}
