package cockroachdb

import (
	"context"
	_ "embed"
	"sync"

	"github.com/adlerhurst/eventstore/v0"
)

var (
	//go:embed subscribe.sql
	subscribeStmt string

	//go:embed unsubscribe.sql
	unsubscribeStmt string
)

func (crdb *CockroachDB) Subscribe(ctx context.Context, subjects []eventstore.Subject) (id string, err error) {
	row := crdb.client.QueryRow(ctx, subscribeStmt, subjects)
	if err := row.Scan(&id); err != nil {
		return "", err
	}
	return id, nil
}

func (crdb *CockroachDB) Unubscribe(ctx context.Context, id string) error {
	_, err := crdb.client.Exec(ctx, unsubscribeStmt, id)
	return err
}

/*
user.1.added

(pattern[0] in ('*', 'user'))
(pattern[1] in ('*', '1'))
(pattern[1] in ('*', 'added'))

user.1.*
user.*.*
#
user.#
user.1.#
*.*.*
*.*.added
*.1.added
user.*.added


user

user
*
#

user.added

#
*.*
*.added
user.*
user.#
user.added


*/

func (crdb *CockroachDB) subscriptions(ctx context.Context, action eventstore.TextSubjects) ([]string, error) {
	var patters []string
	for _, item := range action {

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
