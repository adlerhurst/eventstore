package eventstore

import "strings"

type Subject interface {
	subject()
}

var (
	SingleToken = Subject(singleToken{})
	MultiToken  = Subject(multiToken{})
)

type TextSubject string

func (TextSubject) subject() {}

type singleToken struct{}

func (singleToken) subject() {}

type multiToken struct{}

func (multiToken) subject() {}

type TextSubjects []TextSubject

func (ts TextSubjects) Join(sep string) string {
	return join(ts, sep)
}

type text interface{ ~string }

// join is a copy of [strings.Join] which allows type constraints
func join[t text](elems []t, sep string) string {
	switch len(elems) {
	case 0:
		return ""
	case 1:
		return string(elems[0])
	}
	n := len(sep) * (len(elems) - 1)
	for i := 0; i < len(elems); i++ {
		n += len(elems[i])
	}

	var b strings.Builder
	b.Grow(n)
	b.WriteString(string(elems[0]))
	for _, s := range elems[1:] {
		b.WriteString(sep)
		b.WriteString(string(s))
	}
	return b.String()
}
