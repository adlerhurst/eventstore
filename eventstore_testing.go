package eventstore

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

var _ Aggregate = (*testUser)(nil)

type testUser struct {
	id                 string
	currentSequence    uint32
	predefinedSequence uint32
	commands           []Command
}

var _ AggregatePredefinedSequence = (*testUserWithSequence)(nil)

type testUserWithSequence struct {
	*testUser
}

// CurrentSequence implements AggregatePredefinedSequence.
func (a *testUserWithSequence) CurrentSequence() uint32 {
	return a.testUser.predefinedSequence
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

func newTestUser(id string, opts ...testUserOpt) Aggregate {
	tu := &testUser{id: id}

	for _, opt := range opts {
		tu = opt(tu)
	}
	if tu.predefinedSequence > 0 {
		return &testUserWithSequence{tu}
	}
	return tu
}

func withPredefinedSequence(sequence uint32) testUserOpt {
	return func(tu *testUser) *testUser {
		tu.predefinedSequence = sequence
		return tu
	}
}

func withAdded(firstName, lastName, username string) testUserOpt {
	return func(tu *testUser) *testUser {
		tu.currentSequence++
		tu.commands = append(tu.commands, &testUserAdded{
			id:           tu.id,
			FirstName:    firstName,
			LastName:     lastName,
			Username:     username,
			wantSequence: tu.currentSequence,
		})
		return tu
	}
}

func withFirstName(firstName string) testUserOpt {
	return func(tu *testUser) *testUser {
		tu.currentSequence++
		tu.commands = append(tu.commands, &testUserFirstNameChanged{
			id:           tu.id,
			FirstName:    firstName,
			aggregate:    tu.ID(),
			wantSequence: tu.currentSequence,
		})
		return tu
	}
}

func withLastName(lastName string) testUserOpt {
	return func(tu *testUser) *testUser {
		tu.currentSequence++
		tu.commands = append(tu.commands, &testUserLastNameChanged{
			id:           tu.id,
			LastName:     lastName,
			aggregate:    tu.ID(),
			wantSequence: tu.currentSequence,
		})
		return tu
	}
}

func withUsername(username string) testUserOpt {
	return func(tu *testUser) *testUser {
		tu.currentSequence++
		tu.commands = append(tu.commands, &testUsernameChanged{
			id:           tu.id,
			Username:     username,
			aggregate:    tu.ID(),
			wantSequence: tu.currentSequence,
		})
		return tu
	}
}

func withRemoved() testUserOpt {
	return func(tu *testUser) *testUser {
		tu.currentSequence++
		tu.commands = append(tu.commands, &testUserRemoved{
			id:           tu.id,
			aggregate:    tu.ID(),
			wantSequence: tu.currentSequence,
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
	wantSequence uint32

	sequence  uint32
	createdAt time.Time
}

// SetCreationDate implements [Command].
func (e *testUserAdded) SetCreationDate(creationDate time.Time) {
	e.createdAt = creationDate
}

// SetSequence implements [Command].
func (e *testUserAdded) SetSequence(sequence uint32) {
	e.sequence = sequence
}

// Action implements [Action]
func (e *testUserAdded) Action() TextSubjects {
	return []TextSubject{"user", TextSubject(e.id), "added"}
}

// Revision implements [Action]
func (*testUserAdded) Revision() uint16 { return 1 }

// Payload implements [Command]
func (e *testUserAdded) Payload() interface{} { return e }

func (c *testUserAdded) assert(t *testing.T) (failed bool) {
	t.Helper()

	if c.wantSequence != c.sequence {
		t.Errorf("unexpected sequence, want: %v, got: %v", c.wantSequence, c.sequence)
		failed = true
	}

	return failed
}

var _ Command = (*testUserFirstNameChanged)(nil)

type testUserFirstNameChanged struct {
	id        string
	FirstName string `json:"firstName,omitempty"`
	// the following fields are used for assertion
	wantSequence uint32
	aggregate    TextSubjects

	sequence  uint32
	createdAt time.Time
}

// SetCreationDate implements [Command].
func (e *testUserFirstNameChanged) SetCreationDate(creationDate time.Time) {
	e.createdAt = creationDate
}

// SetSequence implements [Command].
func (e *testUserFirstNameChanged) SetSequence(sequence uint32) {
	e.sequence = sequence
}

// Action implements [Action]
func (e *testUserFirstNameChanged) Action() TextSubjects {
	return []TextSubject{"user", TextSubject(e.id), "firstName", "set"}
}

// Revision implements [Action]
func (*testUserFirstNameChanged) Revision() uint16 { return 1 }

// Payload implements [Command]
func (e *testUserFirstNameChanged) Payload() interface{} { return e }

func (e *testUserFirstNameChanged) assert(t *testing.T) (failed bool) {
	t.Helper()

	if e.wantSequence != e.sequence {
		t.Errorf("unexpected sequence, want: %v, got: %v", e.wantSequence, e.sequence)
		failed = true
	}

	return failed
}

var _ Command = (*testUserLastNameChanged)(nil)

type testUserLastNameChanged struct {
	id       string
	LastName string `json:"lastName,omitempty"`
	// the following fields are used for assertion
	wantSequence uint32
	aggregate    TextSubjects

	sequence  uint32
	createdAt time.Time
}

// SetCreationDate implements [Command].
func (e *testUserLastNameChanged) SetCreationDate(creationDate time.Time) {
	e.createdAt = creationDate
}

// SetSequence implements [Command].
func (e *testUserLastNameChanged) SetSequence(sequence uint32) {
	e.sequence = sequence
}

// Action implements [Action]
func (e *testUserLastNameChanged) Action() TextSubjects {
	return []TextSubject{"user", TextSubject(e.id), "lastName", "set"}
}

// Revision implements [Action]
func (*testUserLastNameChanged) Revision() uint16 { return 1 }

// Payload implements [Command]
func (e *testUserLastNameChanged) Payload() interface{} { return e }

func (e *testUserLastNameChanged) assert(t *testing.T) (failed bool) {
	t.Helper()

	if e.wantSequence != e.sequence {
		t.Errorf("unexpected sequence, want: %v, got: %v", e.wantSequence, e.sequence)
		failed = true
	}

	return failed
}

var _ Command = (*testUsernameChanged)(nil)

type testUsernameChanged struct {
	id       string
	Username string `json:"username,omitempty"`
	// the following fields are used for assertion
	wantSequence uint32
	aggregate    TextSubjects

	sequence  uint32
	createdAt time.Time
}

// SetCreationDate implements [Command].
func (e *testUsernameChanged) SetCreationDate(creationDate time.Time) {
	e.createdAt = creationDate
}

// SetSequence implements [Command].
func (e *testUsernameChanged) SetSequence(sequence uint32) {
	e.sequence = sequence
}

// Action implements [Action]
func (e *testUsernameChanged) Action() TextSubjects {
	return []TextSubject{"user", TextSubject(e.id), "username", "set"}
}

// Revision implements [Action]
func (*testUsernameChanged) Revision() uint16 { return 1 }

// Payload implements [Command]
func (e *testUsernameChanged) Payload() interface{} { return e }

func (e *testUsernameChanged) assert(t *testing.T) (failed bool) {
	t.Helper()

	if e.wantSequence != e.sequence {
		t.Errorf("unexpected sequence, want: %v, got: %v", e.wantSequence, e.sequence)
		failed = true
	}

	return failed
}

var _ Command = (*testUserRemoved)(nil)

type testUserRemoved struct {
	id string
	// the following fields are used for assertion
	wantSequence uint32
	aggregate    TextSubjects

	sequence  uint32
	createdAt time.Time
}

// SetCreationDate implements [Command].
func (e *testUserRemoved) SetCreationDate(creationDate time.Time) {
	e.createdAt = creationDate
}

// SetSequence implements [Command].
func (e *testUserRemoved) SetSequence(sequence uint32) {
	e.sequence = sequence
}

// Action implements [Action]
func (e *testUserRemoved) Action() TextSubjects {
	return []TextSubject{"user", TextSubject(e.id), "removed"}
}

// Revision implements [Action]
func (*testUserRemoved) Revision() uint16 { return 1 }

// Payload implements [Command]
func (e *testUserRemoved) Payload() interface{} { return nil }

func (e *testUserRemoved) assert(t *testing.T) (failed bool) {
	t.Helper()

	if e.wantSequence != e.sequence {
		t.Errorf("unexpected sequence, want: %v, got: %v", e.wantSequence, e.sequence)
		failed = true
	}

	return failed
}

type TestEventstore interface {
	Eventstore
	Before(ctx context.Context, t testing.TB) error
	After(ctx context.Context, t testing.TB) error
}

func PushComplianceTests(ctx context.Context, t *testing.T, store TestEventstore) {
	tests := []struct {
		name        string
		aggregates  []Aggregate
		expectedErr error
	}{
		{
			name: "multiple events",
			aggregates: []Aggregate{
				newTestUser("id",
					withAdded("first name", "last name", "username"),
					withRemoved(),
				),
			},
			expectedErr: nil,
		},
		{
			name: "defined sequence",
			aggregates: []Aggregate{
				newTestUser("id",
					withPredefinedSequence(0),
					withAdded("first name", "last name", "username"),
					withRemoved(),
				),
			},
			expectedErr: nil,
		},
		{
			name: "multiple events defined sequence error",
			aggregates: []Aggregate{
				newTestUser("id",
					withPredefinedSequence(2),
					withAdded("first name", "last name", "username"),
					withRemoved(),
				),
			},
			expectedErr: ErrSequenceNotMatched,
		},
		{
			name: "multiple aggregates",
			aggregates: []Aggregate{
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
			expectedErr: nil,
		},
	}
	for _, tt := range tests {
		if err := store.Before(ctx, t); err != nil {
			t.Error("unable to execute store.Before: ", err)
		}
		t.Run(tt.name, func(t *testing.T) {
			err := store.Push(ctx, tt.aggregates...)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("expected error was %v, got: %v", tt.expectedErr, err)
			}
			if tt.expectedErr == nil {
				assertAggregates(t, tt.aggregates)
			}
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

func pushDefaultCommands(ctx context.Context, t testing.TB, store TestEventstore, aggregate Aggregate) {
	t.Helper()

	err := store.Push(ctx, aggregate)
	if err != nil {
		t.Error(err)
	}
}

// func FilterComplianceTests(ctx context.Context, t *testing.T, store TestEventstore) {
// 	type args struct {
// 		filter *Filter
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    []*testUser
// 		wantErr bool
// 	}{
// 		{
// 			name: "multi token",
// 			args: args{
// 				filter: &Filter{
// 					Action: []Subject{TextSubject("user"), TextSubject("id"), MultiToken},
// 				},
// 			},
// 			want: []*testUser{
// 				newTestUser("id",
// 					withAdded("first name", "last name", "username"),
// 					withRemoved(),
// 				),
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "multiple single tokens",
// 			args: args{
// 				filter: &Filter{
// 					Action: []Subject{TextSubject("user"), SingleToken, SingleToken},
// 				},
// 			},
// 			want: []*testUser{
// 				newTestUser("id",
// 					withAdded("first name", "last name", "username"),
// 					withRemoved(),
// 				),
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "all",
// 			args: args{
// 				filter: &Filter{
// 					Action: []Subject{
// 						TextSubject("user"),
// 						TextSubject("id"),
// 						TextSubject("added"),
// 					},
// 				},
// 			},
// 			want: []*testUser{
// 				newTestUser("id",
// 					withAdded("first name", "last name", "username"),
// 				),
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "crdb",
// 			args: args{
// 				filter: &Filter{
// 					Action: []Subject{MultiToken},
// 				},
// 			},
// 			want: []*testUser{
// 				newTestUser("id",
// 					withAdded("first name", "last name", "username"),
// 					withRemoved(),
// 				),
// 			},
// 			wantErr: false,
// 		},
// 	}
// 	if err := store.Before(ctx, t); err != nil {
// 		t.Error("unable to execute store.Before: ", err)
// 	}
// 	_, err := store.Push(ctx,
// 		newTestUser("id",
// 			withAdded("first name", "last name", "username"),
// 			withRemoved(),
// 		),
// 	)
// 	if err != nil {
// 		t.Fatalf("unable to push events: %v", err)
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, err := store.Filter(ctx, tt.args.filter)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("Filter() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			assertEvents(t, tt.want, got)
// 		})
// 	}
// 	if err := store.After(ctx, t); err != nil {
// 		t.Error("unable to execute store.After: ", err)
// 	}
// }

// func FilterBenchTests(ctx context.Context, b *testing.B, store TestEventstore) {
// 	type args struct {
// 		filter *Filter
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    int
// 		wantErr bool
// 	}{
// 		{
// 			name: "multi token",
// 			args: args{
// 				filter: &Filter{
// 					Action: []Subject{TextSubject("user"), TextSubject("5555"), MultiToken},
// 				},
// 			},
// 			want:    2,
// 			wantErr: false,
// 		},
// 		{
// 			name: "multiple single tokens",
// 			args: args{
// 				filter: &Filter{
// 					Action: []Subject{TextSubject("user"), SingleToken, SingleToken},
// 				},
// 			},
// 			want:    20_000,
// 			wantErr: false,
// 		},
// 		{
// 			name: "all",
// 			args: args{
// 				filter: &Filter{
// 					Action: []Subject{
// 						TextSubject("user"),
// 						TextSubject("5555"),
// 						TextSubject("added"),
// 					},
// 				},
// 			},
// 			want:    1,
// 			wantErr: false,
// 		},
// 		{
// 			name: "multi token at beginning",
// 			args: args{
// 				filter: &Filter{
// 					Action: []Subject{MultiToken},
// 				},
// 			},
// 			want:    20_000,
// 			wantErr: false,
// 		},
// 		{
// 			name: "all added",
// 			args: args{
// 				filter: &Filter{
// 					Action: []Subject{TextSubject("user"), SingleToken, TextSubject("added")},
// 				},
// 			},
// 			want:    10_000,
// 			wantErr: false,
// 		},
// 	}
// 	if err := store.Before(ctx, b); err != nil {
// 		b.Error("unable to execute store.Before: ", err)
// 	}
// 	for i := 0; i < 10_000; i++ {
// 		err := store.Push(ctx,
// 			newTestUser(strconv.Itoa(i),
// 				withAdded("first name", "last name", "username"),
// 				withRemoved(),
// 			),
// 		)
// 		if err != nil {
// 			b.Fatalf("unable to push events: %v", err)
// 		}
// 	}
// 	b.ResetTimer()
// 	for _, tt := range tests {
// 		b.Run(tt.name, func(b *testing.B) {
// 			b.RunParallel(func(p *testing.PB) {
// 				for n := 0; p.Next(); n++ {
// 					got, err := store.Filter(ctx, tt.args.filter)
// 					if (err != nil) != tt.wantErr {
// 						b.Errorf("Filter() error = %v, wantErr %v", err, tt.wantErr)
// 						return
// 					}
// 					if len(got) != tt.want {
// 						b.Errorf("unexpected amount of events. want: %d, got %d", tt.want, len(got))
// 					}
// 				}
// 			})
// 		})
// 	}
// 	if err := store.After(ctx, b); err != nil {
// 		b.Error("unable to execute store.After: ", err)
// 	}
// }

type commandAsserter interface {
	assert(t *testing.T) bool
}

func assertAggregates(t *testing.T, aggregates []Aggregate) (failed bool) {
	t.Helper()

	var index int
	for _, aggregate := range aggregates {
		for _, command := range aggregate.Commands() {
			asserter, ok := command.(commandAsserter)
			if !ok {
				t.Fatalf("test command is not assertable: %v", command.Action())
			}
			failed = failed || asserter.assert(t)
			index++
		}
	}

	return failed
}

// func assertPayload(t *testing.T, want Command, got func(object any) error) (failed bool) {
// 	unmarshalWant := func(object any) error {
// 		data, err := json.Marshal(want.Payload())
// 		if err != nil {
// 			return err
// 		}
// 		return json.Unmarshal(data, object)
// 	}
// 	var (
// 		gotPayload, wantPayload interface{}
// 	)
// 	if err := unmarshalWant(&wantPayload); err != nil {
// 		t.Errorf("unable to unmarshal want payload: %v", err)
// 		failed = true
// 	}
// 	if err := got(&gotPayload); err != nil {
// 		t.Errorf("unable to unmarshal gotten payload: %v", err)
// 		failed = true
// 	}
// 	if !reflect.DeepEqual(gotPayload, wantPayload) {
// 		t.Errorf("payload not equal want: %#v got: %#v", wantPayload, gotPayload)
// 		failed = true
// 	}

// 	return failed
// }
