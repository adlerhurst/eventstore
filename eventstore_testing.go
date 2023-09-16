package eventstore

import (
	"context"
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
)

var _ Aggregate = (*testUser)(nil)

type testUser struct {
	id              string
	currentSequence uint64
	commands        []Command
}

// Commands implements Aggregate.
func (a *testUser) Commands() []Command {
	return a.commands
}

// ID implements Aggregate.
func (a *testUser) ID() TextSubjects {
	return []TextSubject{"user", TextSubject(a.id)}
}

type testUserOpt func(*testUser) *testUser

func newTestUser(id string, opts ...testUserOpt) *testUser {
	tu := &testUser{id: id}

	for _, opt := range opts {
		tu = opt(tu)
	}
	return tu
}

func withAdded(firstName, lastName, username string) testUserOpt {
	return func(tu *testUser) *testUser {
		tu.currentSequence++
		tu.commands = append(tu.commands, &testUserAdded{
			id:        tu.id,
			FirstName: firstName,
			LastName:  lastName,
			Username:  username,
			aggregate: tu.ID(),
			sequence:  tu.currentSequence,
		})
		return tu
	}
}

func withFirstName(firstName string) testUserOpt {
	return func(tu *testUser) *testUser {
		tu.currentSequence++
		tu.commands = append(tu.commands, &testUserFirstNameChanged{
			id:        tu.id,
			FirstName: firstName,
			aggregate: tu.ID(),
			sequence:  tu.currentSequence,
		})
		return tu
	}
}

func withLastName(lastName string) testUserOpt {
	return func(tu *testUser) *testUser {
		tu.currentSequence++
		tu.commands = append(tu.commands, &testUserLastNameChanged{
			id:        tu.id,
			LastName:  lastName,
			aggregate: tu.ID(),
			sequence:  tu.currentSequence,
		})
		return tu
	}
}

func withUsername(username string) testUserOpt {
	return func(tu *testUser) *testUser {
		tu.currentSequence++
		tu.commands = append(tu.commands, &testUsernameChanged{
			id:        tu.id,
			Username:  username,
			aggregate: tu.ID(),
			sequence:  tu.currentSequence,
		})
		return tu
	}
}

func withRemoved() testUserOpt {
	return func(tu *testUser) *testUser {
		tu.currentSequence++
		tu.commands = append(tu.commands, &testUserRemoved{
			id:        tu.id,
			aggregate: tu.ID(),
			sequence:  tu.currentSequence,
		})
		return tu
	}
}

var _ Command = (*testUserAdded)(nil)

type testUserAdded struct {
	id        string
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Username  string `json:"username,omitempty"`
	// the following fields are used for assertion
	sequence  uint64
	aggregate TextSubjects
}

// Action implements [Action]
func (e *testUserAdded) Action() TextSubjects {
	return []TextSubject{"user", TextSubject(e.id), "added"}
}

// Revision implements [Action]
func (*testUserAdded) Revision() uint16 { return 1 }

// Payload implements [Command]
func (e *testUserAdded) Payload() interface{} { return e }

func (e *testUserAdded) assertEvent(t *testing.T, got Event) (failed bool) {
	t.Helper()

	if !reflect.DeepEqual(e.aggregate, got.Aggregate()) {
		t.Errorf("unexpected aggregate, want: %v, got: %v", e.aggregate, got.Aggregate())
		failed = true
	}
	if !reflect.DeepEqual(e.Action(), got.Action()) {
		t.Errorf("unexpected action, want: %v, got: %v", e.Action(), got.Action())
		failed = true
	}
	if e.sequence != got.Sequence() {
		t.Errorf("unexpected sequence, want: %v, got: %v", e.sequence, got.Sequence())
		failed = true
	}
	if e.Revision() != got.Revision() {
		t.Errorf("unexpected revision, want: %v, got: %v", e.Revision(), got.Revision())
		failed = true
	}

	failed = failed || assertPayload(t, e, got.UnmarshalPayload)

	return failed
}

var _ Command = (*testUserFirstNameChanged)(nil)

type testUserFirstNameChanged struct {
	id        string
	FirstName string `json:"firstName,omitempty"`
	// the following fields are used for assertion
	sequence  uint64
	aggregate TextSubjects
}

// Action implements [Action]
func (e *testUserFirstNameChanged) Action() TextSubjects {
	return []TextSubject{"user", TextSubject(e.id), "firstName", "set"}
}

// Revision implements [Action]
func (*testUserFirstNameChanged) Revision() uint16 { return 1 }

// Payload implements [Command]
func (e *testUserFirstNameChanged) Payload() interface{} { return e }

func (e *testUserFirstNameChanged) assertEvent(t *testing.T, got Event) (failed bool) {
	t.Helper()

	if !reflect.DeepEqual(e.aggregate, got.Aggregate()) {
		t.Errorf("unexpected aggregate, want: %v, got: %v", e.aggregate, got.Aggregate())
		failed = true
	}
	if !reflect.DeepEqual(e.Action(), got.Action()) {
		t.Errorf("unexpected action, want: %v, got: %v", e.Action(), got.Action())
		failed = true
	}
	if e.sequence != got.Sequence() {
		t.Errorf("unexpected sequence, want: %v, got: %v", e.sequence, got.Sequence())
		failed = true
	}
	if e.Revision() != got.Revision() {
		t.Errorf("unexpected revision, want: %v, got: %v", e.Revision(), got.Revision())
		failed = true
	}

	failed = failed || assertPayload(t, e, got.UnmarshalPayload)

	return failed
}

var _ Command = (*testUserLastNameChanged)(nil)

type testUserLastNameChanged struct {
	id       string
	LastName string `json:"lastName,omitempty"`
	// the following fields are used for assertion
	sequence  uint64
	aggregate TextSubjects
}

// Action implements [Action]
func (e *testUserLastNameChanged) Action() TextSubjects {
	return []TextSubject{"user", TextSubject(e.id), "lastName", "set"}
}

// Revision implements [Action]
func (*testUserLastNameChanged) Revision() uint16 { return 1 }

// Payload implements [Command]
func (e *testUserLastNameChanged) Payload() interface{} { return e }

func (e *testUserLastNameChanged) assertEvent(t *testing.T, got Event) (failed bool) {
	t.Helper()

	if !reflect.DeepEqual(e.aggregate, got.Aggregate()) {
		t.Errorf("unexpected aggregate, want: %v, got: %v", e.aggregate, got.Aggregate())
		failed = true
	}
	if !reflect.DeepEqual(e.Action(), got.Action()) {
		t.Errorf("unexpected action, want: %v, got: %v", e.Action(), got.Action())
		failed = true
	}
	if e.sequence != got.Sequence() {
		t.Errorf("unexpected sequence, want: %v, got: %v", e.sequence, got.Sequence())
		failed = true
	}
	if e.Revision() != got.Revision() {
		t.Errorf("unexpected revision, want: %v, got: %v", e.Revision(), got.Revision())
		failed = true
	}

	failed = failed || assertPayload(t, e, got.UnmarshalPayload)

	return failed
}

var _ Command = (*testUsernameChanged)(nil)

type testUsernameChanged struct {
	id       string
	Username string `json:"username,omitempty"`
	// the following fields are used for assertion
	sequence  uint64
	aggregate TextSubjects
}

// Action implements [Action]
func (e *testUsernameChanged) Action() TextSubjects {
	return []TextSubject{"user", TextSubject(e.id), "username", "set"}
}

// Revision implements [Action]
func (*testUsernameChanged) Revision() uint16 { return 1 }

// Payload implements [Command]
func (e *testUsernameChanged) Payload() interface{} { return e }

func (e *testUsernameChanged) assertEvent(t *testing.T, got Event) (failed bool) {
	t.Helper()

	if !reflect.DeepEqual(e.aggregate, got.Aggregate()) {
		t.Errorf("unexpected aggregate, want: %v, got: %v", e.aggregate, got.Aggregate())
		failed = true
	}
	if !reflect.DeepEqual(e.Action(), got.Action()) {
		t.Errorf("unexpected action, want: %v, got: %v", e.Action(), got.Action())
		failed = true
	}
	if e.sequence != got.Sequence() {
		t.Errorf("unexpected sequence, want: %v, got: %v", e.sequence, got.Sequence())
		failed = true
	}
	if e.Revision() != got.Revision() {
		t.Errorf("unexpected revision, want: %v, got: %v", e.Revision(), got.Revision())
		failed = true
	}

	failed = failed || assertPayload(t, e, got.UnmarshalPayload)

	return failed
}

var _ Command = (*testUserRemoved)(nil)

type testUserRemoved struct {
	id string
	// the following fields are used for assertion
	sequence  uint64
	aggregate TextSubjects
}

// Action implements [Action]
func (e *testUserRemoved) Action() TextSubjects {
	return []TextSubject{"user", TextSubject(e.id), "removed"}
}

// Revision implements [Action]
func (*testUserRemoved) Revision() uint16 { return 1 }

// Payload implements [Command]
func (e *testUserRemoved) Payload() interface{} { return nil }

func (e *testUserRemoved) assertEvent(t *testing.T, got Event) (failed bool) {
	t.Helper()

	if !reflect.DeepEqual(e.aggregate, got.Aggregate()) {
		t.Errorf("unexpected aggregate, want: %v, got: %v", e.aggregate, got.Aggregate())
		failed = true
	}
	if !reflect.DeepEqual(e.Action(), got.Action()) {
		t.Errorf("unexpected action, want: %v, got: %v", e.Action(), got.Action())
		failed = true
	}
	if e.sequence != got.Sequence() {
		t.Errorf("unexpected sequence, want: %v, got: %v", e.sequence, got.Sequence())
		failed = true
	}
	if e.Revision() != got.Revision() {
		t.Errorf("unexpected revision, want: %v, got: %v", e.Revision(), got.Revision())
		failed = true
	}

	failed = failed || assertPayload(t, e, got.UnmarshalPayload)

	return failed
}

type TestEventstore interface {
	Eventstore
	Before(ctx context.Context, t testing.TB) error
	After(ctx context.Context, t testing.TB) error
}

func PushComplianceTests(ctx context.Context, t *testing.T, store TestEventstore) {
	tests := []struct {
		name       string
		aggregates []*testUser
		wantErr    bool
	}{
		{
			name: "multiple events",
			aggregates: []*testUser{
				newTestUser("id",
					withAdded("first name", "last name", "username"),
					withRemoved(),
				),
			},
			wantErr: false,
		},
		{
			name: "multiple aggregates",
			aggregates: []*testUser{
				newTestUser("id",
					withAdded("first name", "last name", "user name"),
					withUsername("changed username"),
					withRemoved(),
				),
				newTestUser("2",
					withAdded("first name 2", "last name 2", "user name 2"),
					withFirstName("new first name 2"),
					withLastName("new last name 2"),
				),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		if err := store.Before(ctx, t); err != nil {
			t.Error("unable to execute store.Before: ", err)
		}
		aggregates := make([]Aggregate, len(tt.aggregates))
		for i, aggregate := range tt.aggregates {
			aggregates[i] = aggregate
		}
		t.Run(tt.name, func(t *testing.T) {
			got, err := store.Push(ctx, aggregates...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Push() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assertEvents(t, tt.aggregates, got)
		})
		if err := store.After(ctx, t); err != nil {
			t.Error("unable to execute store.After: ", err)
		}
	}
}

func PushParallelOnSameAggregate(ctx context.Context, b *testing.B, store TestEventstore) {
	if err := store.Before(ctx, b); err != nil {
		b.Error("unable to execute store.Before: ", err)
	}
	b.ResetTimer()
	b.RunParallel(func(p *testing.PB) {
		for n := 0; p.Next(); n++ {
			user := newTestUser(b.Name(),
				withAdded("first name", "last name", "username"),
				withRemoved(),
			)

			pushDefaultCommands(ctx, b, store, user)
		}
	})
	if err := store.After(ctx, b); err != nil {
		b.Error("unable to execute store.After: ", err)
	}
}

func PushParallelOnDifferentAggregates(ctx context.Context, b *testing.B, store TestEventstore) {
	if err := store.Before(ctx, b); err != nil {
		b.Error("unable to execute store.Before: ", err)
	}
	b.ResetTimer()

	var n atomic.Int64

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {

			i := n.Add(1)

			user := newTestUser(strings.ReplaceAll(b.Name(), "/", "-")+strconv.Itoa(int(i)),
				withAdded("first name", "last name", "username"),
				withRemoved(),
			)

			pushDefaultCommands(ctx, b, store, user)
		}
	})
	if err := store.After(ctx, b); err != nil {
		b.Error("unable to execute store.After: ", err)
	}
}

func pushDefaultCommands(ctx context.Context, t testing.TB, store TestEventstore, user *testUser) {
	t.Helper()

	_, err := store.Push(ctx, user)
	if err != nil {
		t.Error(err)
	}
}

func FilterComplianceTests(ctx context.Context, t *testing.T, store TestEventstore) {
	type args struct {
		filter *Filter
	}
	tests := []struct {
		name    string
		args    args
		want    []*testUser
		wantErr bool
	}{
		{
			name: "multi token",
			args: args{
				filter: &Filter{
					Action: []Subject{TextSubject("user"), TextSubject("5555"), MultiToken},
				},
			},
			want: []*testUser{
				newTestUser("5555",
					withAdded("first name", "last name", "username"),
					withRemoved(),
				),
			},
			wantErr: false,
		},
		{
			name: "multiple single tokens",
			args: args{
				filter: &Filter{
					Action: []Subject{TextSubject("user"), SingleToken, SingleToken},
				},
			},
			want: []*testUser{
				newTestUser("5555",
					withAdded("first name", "last name", "username"),
					withRemoved(),
				),
			},
			wantErr: false,
		},
		{
			name: "all",
			args: args{
				filter: &Filter{
					Action: []Subject{
						TextSubject("user"),
						TextSubject("5555"),
						TextSubject("added"),
					},
				},
			},
			want: []*testUser{
				newTestUser("5555",
					withAdded("first name", "last name", "username"),
				),
			},
			wantErr: false,
		},
		{
			name: "crdb",
			args: args{
				filter: &Filter{
					Action: []Subject{MultiToken},
				},
			},
			want: []*testUser{
				newTestUser("5555",
					withAdded("first name", "last name", "username"),
					withRemoved(),
				),
			},
			wantErr: false,
		},
	}
	if err := store.Before(ctx, t); err != nil {
		t.Error("unable to execute store.Before: ", err)
	}
	for i := 0; i < 10_000; i++ {
		_, err := store.Push(ctx,
			newTestUser(strconv.Itoa(i),
				withAdded("first name", "last name", "username"),
				withRemoved(),
			),
		)
		if err != nil {
			t.Fatalf("unable to push events: %v", err)
		}
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := store.Filter(ctx, tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("Filter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assertEvents(t, tt.want, got)
		})
	}
	if err := store.After(ctx, t); err != nil {
		t.Error("unable to execute store.After: ", err)
	}
}

func FilterBenchTests(ctx context.Context, b *testing.B, store TestEventstore) {
	type args struct {
		filter *Filter
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "multi token",
			args: args{
				filter: &Filter{
					Action: []Subject{TextSubject("user"), TextSubject("5555"), MultiToken},
				},
			},
			want:    2,
			wantErr: false,
		},
		{
			name: "multiple single tokens",
			args: args{
				filter: &Filter{
					Action: []Subject{TextSubject("user"), SingleToken, SingleToken},
				},
			},
			want:    20_000,
			wantErr: false,
		},
		{
			name: "all",
			args: args{
				filter: &Filter{
					Action: []Subject{
						TextSubject("user"),
						TextSubject("5555"),
						TextSubject("added"),
					},
				},
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "multi token at beginning",
			args: args{
				filter: &Filter{
					Action: []Subject{MultiToken},
				},
			},
			want:    20_000,
			wantErr: false,
		},
		{
			name: "all added",
			args: args{
				filter: &Filter{
					Action: []Subject{TextSubject("user"), SingleToken, TextSubject("added")},
				},
			},
			want:    10_000,
			wantErr: false,
		},
	}
	if err := store.Before(ctx, b); err != nil {
		b.Error("unable to execute store.Before: ", err)
	}
	for i := 0; i < 10_000; i++ {
		_, err := store.Push(ctx,
			newTestUser(strconv.Itoa(i),
				withAdded("first name", "last name", "username"),
				withRemoved(),
			),
		)
		if err != nil {
			b.Fatalf("unable to push events: %v", err)
		}
	}
	b.ResetTimer()
	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			b.RunParallel(func(p *testing.PB) {
				for n := 0; p.Next(); n++ {
					got, err := store.Filter(ctx, tt.args.filter)
					if (err != nil) != tt.wantErr {
						b.Errorf("Filter() error = %v, wantErr %v", err, tt.wantErr)
						return
					}
					if len(got) != tt.want {
						b.Errorf("unexpected amount of events. want: %d, got %d", tt.want, len(got))
					}
				}
			})
		})
	}
	if err := store.After(ctx, b); err != nil {
		b.Error("unable to execute store.After: ", err)
	}
}

type eventAsserter interface {
	assertEvent(t *testing.T, e Event) bool
}

func assertEvents(t *testing.T, want []*testUser, got []Event) (failed bool) {
	t.Helper()

	var index int
	for _, testUser := range want {
		for _, command := range testUser.commands {
			asserter, ok := command.(eventAsserter)
			if !ok {
				t.Fatalf("test command is not assertable: %v", command.Action())
			}
			failed = failed || asserter.assertEvent(t, got[index])
			index++
		}
	}

	return failed
}

func assertPayload(t *testing.T, want Command, got func(object any) error) (failed bool) {
	unmarshalWant := func(object any) error {
		data, err := json.Marshal(want.Payload())
		if err != nil {
			return err
		}
		return json.Unmarshal(data, object)
	}
	var (
		gotPayload, wantPayload interface{}
	)
	if err := unmarshalWant(&wantPayload); err != nil {
		t.Errorf("unable to unmarshal want payload: %v", err)
		failed = true
	}
	if err := got(&gotPayload); err != nil {
		t.Errorf("unable to unmarshal gotten payload: %v", err)
		failed = true
	}
	if !reflect.DeepEqual(gotPayload, wantPayload) {
		t.Errorf("payload not equal want: %#v got: %#v", wantPayload, gotPayload)
		failed = true
	}

	return failed
}
