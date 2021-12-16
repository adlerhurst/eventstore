package eventstore

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
