package eventstore

import "strings"

type Subject interface {
	subject()
}

var (
	// SingleToken represents asterisk character (*)
	SingleToken = Subject(singleToken{})
	// MultiToken represents the greater than symbol (>), also known as full wildcard
	MultiToken = Subject(multiToken{})
)

type TextSubject string

func (TextSubject) subject() {}

type singleToken struct{}

func (singleToken) subject() {}

type multiToken struct{}

func (multiToken) subject() {}

func isSubjectsValid(subs []Subject) bool {
	if len(subs) == 0 {
		return false
	}

	for i, sub := range subs {
		switch s := sub.(type) {
		case TextSubject:
			if strings.ContainsAny(string(s), "*>") {
				return false
			}
		case multiToken:
			if i+1 != len(subs) {
				return false
			}
		}
	}
	return true
}
