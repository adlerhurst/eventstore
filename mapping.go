package eventstore

import "log"

func (es *Eventstore) RegisterEvent(subs []Subject, mapping mapEvent) *Eventstore {
	for _, t := range es.types {
		if t.check(subs[0]) {
			t.register(mapping, subs[1:])
			return es
		}
	}

	t := newTyp(subs[0])
	es.types = append(es.types, t)
	t.register(mapping, subs[1:])

	return es
}

func (es *Eventstore) MapEvent(e *Event) *Event {
	for _, t := range es.types {
		if event := t.mapEvent(e, 0); event != nil {
			return event
		}
	}
	return e
}

func newTyp(s Subject) *typ {
	return &typ{
		sub: s,
	}
}

type mapEvent func(*Event) *Event

func baseMapping(e *Event) *Event {
	return e
}

type typ struct {
	sub     Subject
	mapping mapEvent

	nodes []*typ
}

func (t *typ) check(s Subject) bool {
	return mapping(t.sub)(s)
}

func (t *typ) mapEvent(e *Event, idx int) *Event {
	if !t.check(e.Subjects[idx]) {
		return nil
	}
	if idx == len(e.Subjects)-1 {
		if t.mapping != nil {
			return t.mapping(e)
		}
		return baseMapping(e)
	}
	for _, n := range t.nodes {
		event := n.mapEvent(e, idx+1)
		if event != nil {
			return event
		}
	}
	return nil
}

func (t *typ) register(eventMapping mapEvent, subs []Subject) {
	if len(subs) == 0 {
		t.mapping = eventMapping
		return
	}
	for _, n := range t.nodes {
		if n.sub == subs[0] {
			n.register(eventMapping, subs[1:])
			return
		}
	}
	n := newTyp(subs[0])
	t.nodes = append(t.nodes, n)
	n.register(eventMapping, subs[1:])
}

func mapping(sub Subject) func(Subject) bool {
	switch s := sub.(type) {
	case TextSubject:
		return func(sub Subject) bool {
			return s == sub.(TextSubject)
		}
	case singleToken:
		return func(sub Subject) bool {
			return sub != nil
		}
	case multiToken:
		return func(sub Subject) bool {
			return sub != nil
		}
	}
	log.Fatal("sub not implemented")
	return nil
}

type Mapping struct {
	s     Subject
	f     []func(...Event)
	nodes []*Mapping
}

func Register(m *Mapping, subs []Subject, f func(...Event)) {
	if len(subs) == 0 {
		m.f = append(m.f, f)
	}

	for _, node := range m.nodes {
		if node.Is(subs[0]) {
			Register(node, subs[1:], f)
			return
		}
	}
	node := &Mapping{
		s: subs[0],
	}
	m.nodes = append(m.nodes, node)
	Register(node, subs[1:], f)
}

func (m *Mapping) Is(sub Subject) bool {
	return m.s == sub
}

func (m *Mapping) postPush(event Event, subIdx int) {
	if event.Subjects[subIdx] != m.s {
		return
	}
	if len(event.Subjects) == subIdx+1 {
		for _, f := range m.f {
			f(event)
		}
	}
}

func (m *Mapping) shouldExec(subs []Subject) bool {

	return false
}

// orgs.01.users.20.grants.300.roles.admin.added
// orgs.>
//
// user.id.added
// user.>
// user.*.added
