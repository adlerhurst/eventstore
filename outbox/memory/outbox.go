package memory

import (
	"sync"

	"github.com/adlerhurst/eventstore/v0"
)

func (crdb *CockroachDB) Subscribe(r receiver, subjects []eventstore.Subject) {
	crdb.outbox.mu.Lock()
	defer crdb.outbox.mu.Unlock()

	crdb.outbox.addReceiver(r, subjects)
}

func (crdb *CockroachDB) Unubscribe(r receiver, subjects []eventstore.Subject) {
	crdb.outbox.mu.Lock()
	defer crdb.outbox.mu.Unlock()

	sub := crdb.outbox.lookup(subjects)
	if sub == nil {
		return
	}

	for i, existingReceiver := range sub.receivers {
		if existingReceiver != r {
			continue
		}
		sub.receivers[i] = sub.receivers[len(sub.receivers)-1]
		sub.receivers[len(sub.receivers)-1] = ""
		sub.receivers = sub.receivers[:len(sub.receivers)-1]
		return
	}
}

type outbox struct {
	mu            sync.RWMutex
	subscriptions *subscription
}

func (o *outbox) lookup(subjects []eventstore.Subject) *subscription {
	return o.subscriptions.lookup(subjects)
}

func (o *outbox) lookupReceivers(subjects eventstore.TextSubjects) (receivers []receiver) {
	return o.subscriptions.lookupReceivers(subjects)
}

func (o *outbox) addReceiver(r receiver, subjects []eventstore.Subject) {
	o.subscriptions.addReceiver(r, subjects)
}

type subscription struct {
	subject   eventstore.Subject
	Leafs     []*subscription
	receivers []receiver
}

// lookup has to be an exact match e.g. user.*.added != user.id.added
func (sub *subscription) lookup(subjects []eventstore.Subject) *subscription {
	if sub.subject != subjects[0] {
		return nil
	}

	if len(subjects) == 1 {
		return sub
	}

	for _, leaf := range sub.Leafs {
		if s := leaf.lookup(subjects[1:]); s != nil {
			return s
		}
	}

	return nil
}

// lookupReceivers allows wildcard serach e.g. user.*.added == user.id.added
func (sub *subscription) lookupReceivers(subjects eventstore.TextSubjects) (receivers []receiver) {
	if sub.subject == nil {
		// this is the root subscription
		for _, leaf := range sub.Leafs {
			receivers = append(receivers, leaf.lookupReceivers(subjects)...)
		}
		return receivers
	}

	if sub.subject == eventstore.MultiToken {
		return sub.receivers
	}

	if sub.subject == eventstore.SingleToken {
		if len(subjects) == 1 {
			// last subject
			return sub.receivers
		}
		for _, leaf := range sub.Leafs {
			receivers = append(receivers, leaf.lookupReceivers(subjects[1:])...)
		}
		return receivers
	}

	if sub.subject != subjects[0] {
		return nil
	}

	receivers = sub.receivers
	for _, leaf := range sub.Leafs {
		receivers = append(receivers, leaf.lookupReceivers(subjects[1:])...)
	}

	return receivers
}

func (sub *subscription) addReceiver(r receiver, subjects []eventstore.Subject) {
	if len(subjects) == 0 {
		for _, existingReceiver := range sub.receivers {
			if existingReceiver == r {
				return
			}
		}
		sub.receivers = append(sub.receivers, r)
		return
	}

	for _, leaf := range sub.Leafs {
		if leaf.subject == subjects[0] {
			leaf.addReceiver(r, subjects[1:])
			return
		}
	}

	leaf := &subscription{
		subject:   subjects[0],
		receivers: []receiver{},
		Leafs:     []*subscription{},
	}
	sub.Leafs = append(sub.Leafs, leaf)
	leaf.addReceiver(r, subjects[1:])
}

type receiver string
