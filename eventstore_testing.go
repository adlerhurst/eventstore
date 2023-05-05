package eventstore

import (
	// "context"
	"context"
	"encoding/json"
	"strconv"

	// "fmt"
	"reflect"
	// "strconv"
	"testing"
	"time"
	// "github.com/cockroachdb/cockroach-go/v2/testserver"
	// "github.com/jackc/pgx/v5/pgxpool"
	// . "github.com/adlerhurst/eventstore/v0"
	// "github.com/adlerhurst/eventstore/v0/cockroachdb"
	// "github.com/adlerhurst/eventstore/v0/memory"
)

type testUser struct {
	id        string
	firstName string
	lastName  string
	username  string
}

var defaultTestUser = &testUser{
	id:        "id",
	firstName: "firstName",
	lastName:  "lastName",
	username:  "adlerhurst",
}

func (u testUser) toAdded() *testUserAdded {
	return &testUserAdded{
		id:        u.id,
		FirstName: u.firstName,
		LastName:  u.lastName,
		Username:  u.username,
	}
}

var (
	_ Command = (*testUserAdded)(nil)
	_ Event   = (*testUserAdded)(nil)
)

type testUserAdded struct {
	id        string
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Username  string `json:"username,omitempty"`
}

// Action implements [eventstore.Action]
func (e *testUserAdded) Action() TextSubjects {
	return []TextSubject{"user", TextSubject(e.id), "added"}
}

// Aggregate implements [eventstore.Action]
func (e *testUserAdded) Aggregate() TextSubjects {
	return []TextSubject{"user", TextSubject(e.id)}
}

// Metadata implements [eventstore.Action]
func (*testUserAdded) Metadata() map[string]interface{} {
	return map[string]interface{}{
		"editorService": "svc",
		"editorUser":    "usr",
		"resourceOwner": "ro",
	}
}

// Revision implements [eventstore.Action]
func (*testUserAdded) Revision() uint16 { return 1 }

// Payload implements [eventstore.Command]
func (e *testUserAdded) Payload() interface{} { return e }

// Options implements [eventstore.Command]
func (e *testUserAdded) Options() []func(Command) error { return nil }

// Sequence implements [eventstore.Event]
func (e *testUserAdded) Sequence() uint64 { return 0 }

// CreationDate implements [eventstore.Event]
func (e *testUserAdded) CreationDate() time.Time { return time.Time{} }

// UnmarshalPayload implements [eventstore.Event]
func (e *testUserAdded) UnmarshalPayload(object interface{}) error {
	data, err := json.Marshal(e)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, object)
}

func (u testUser) toFirstNameChanged() *testUserFirstNameChanged {
	return &testUserFirstNameChanged{
		id:        u.id,
		FirstName: u.firstName,
	}
}

var _ Command = (*testUserFirstNameChanged)(nil)

type testUserFirstNameChanged struct {
	id        string
	FirstName string `json:"firstName,omitempty"`
}

// Action implements [eventstore.Action]
func (e *testUserFirstNameChanged) Action() TextSubjects {
	return []TextSubject{"user", TextSubject(e.id), "changed", "firstName"}
}

// Aggregate implements [eventstore.Action]
func (e *testUserFirstNameChanged) Aggregate() TextSubjects {
	return []TextSubject{"user", TextSubject(e.id)}
}

// Metadata implements [eventstore.Action]
func (*testUserFirstNameChanged) Metadata() map[string]interface{} {
	return map[string]interface{}{
		"editorService": "svc",
		"editorUser":    "usr",
		"resourceOwner": "ro",
	}
}

// Revision implements [eventstore.Action]
func (*testUserFirstNameChanged) Revision() uint16 { return 1 }

// Payload implements [eventstore.Command]
func (e *testUserFirstNameChanged) Payload() interface{} { return e }

// Options implements [eventstore.Command]
func (e *testUserFirstNameChanged) Options() []func(Command) error { return nil }

// Sequence implements [eventstore.Event]
func (e *testUserFirstNameChanged) Sequence() uint64 { return 0 }

// CreationDate implements [eventstore.Event]
func (e *testUserFirstNameChanged) CreationDate() time.Time { return time.Time{} }

// UnmarshalPayload implements [eventstore.Event]
func (e *testUserFirstNameChanged) UnmarshalPayload(object interface{}) error {
	data, err := json.Marshal(e)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, object)
}

func (u testUser) toLastNameChanged() *testUserLastNameChanged {
	return &testUserLastNameChanged{
		id:       u.id,
		LastName: u.lastName,
	}
}

var _ Command = (*testUserLastNameChanged)(nil)

type testUserLastNameChanged struct {
	id       string
	LastName string `json:"lastName,omitempty"`
}

// Action implements [eventstore.Action]
func (e *testUserLastNameChanged) Action() TextSubjects {
	return []TextSubject{"user", TextSubject(e.id), "changed", "lastName"}
}

// Aggregate implements [eventstore.Action]
func (e *testUserLastNameChanged) Aggregate() TextSubjects {
	return []TextSubject{"user", TextSubject(e.id)}
}

// Metadata implements [eventstore.Action]
func (*testUserLastNameChanged) Metadata() map[string]interface{} {
	return map[string]interface{}{
		"editorService": "svc",
		"editorUser":    "usr",
		"resourceOwner": "ro",
	}
}

// Revision implements [eventstore.Action]
func (*testUserLastNameChanged) Revision() uint16 { return 1 }

// Payload implements [eventstore.Command]
func (e *testUserLastNameChanged) Payload() interface{} { return e }

// Options implements [eventstore.Command]
func (e *testUserLastNameChanged) Options() []func(Command) error { return nil }

// Sequence implements [eventstore.Event]
func (e *testUserLastNameChanged) Sequence() uint64 { return 0 }

// CreationDate implements [eventstore.Event]
func (e *testUserLastNameChanged) CreationDate() time.Time { return time.Time{} }

// UnmarshalPayload implements [eventstore.Event]
func (e *testUserLastNameChanged) UnmarshalPayload(object interface{}) error {
	data, err := json.Marshal(e)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, object)
}

func (u testUser) toUsernameChanged() *testUsernameChanged {
	return &testUsernameChanged{
		id:       u.id,
		Username: u.username,
	}
}

var _ Command = (*testUsernameChanged)(nil)

type testUsernameChanged struct {
	id       string
	Username string `json:"username,omitempty"`
}

// Action implements [eventstore.Action]
func (e *testUsernameChanged) Action() TextSubjects {
	return []TextSubject{"user", TextSubject(e.id), "changed", "username"}
}

// Aggregate implements [eventstore.Action]
func (e *testUsernameChanged) Aggregate() TextSubjects {
	return []TextSubject{"user", TextSubject(e.id)}
}

// Metadata implements [eventstore.Action]
func (*testUsernameChanged) Metadata() map[string]interface{} {
	return map[string]interface{}{
		"editorService": "svc",
		"editorUser":    "usr",
		"resourceOwner": "ro",
	}
}

// Revision implements [eventstore.Action]
func (*testUsernameChanged) Revision() uint16 { return 1 }

// Payload implements [eventstore.Command]
func (e *testUsernameChanged) Payload() interface{} { return e }

// Options implements [eventstore.Command]
func (e *testUsernameChanged) Options() []func(Command) error { return nil }

// Sequence implements [eventstore.Event]
func (e *testUsernameChanged) Sequence() uint64 { return 0 }

// CreationDate implements [eventstore.Event]
func (e *testUsernameChanged) CreationDate() time.Time { return time.Time{} }

// UnmarshalPayload implements [eventstore.Event]
func (e *testUsernameChanged) UnmarshalPayload(object interface{}) error {
	data, err := json.Marshal(e)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, object)
}

func (u testUser) toRemoved() *testUserRemoved {
	return &testUserRemoved{
		id: u.id,
	}
}

var _ Command = (*testUserRemoved)(nil)

type testUserRemoved struct {
	id string
}

// Action implements [eventstore.Action]
func (e *testUserRemoved) Action() TextSubjects {
	return []TextSubject{"user", TextSubject(e.id), "removed"}
}

// Aggregate implements [eventstore.Action]
func (e *testUserRemoved) Aggregate() TextSubjects {
	return []TextSubject{"user", TextSubject(e.id)}
}

// Metadata implements [eventstore.Action]
func (*testUserRemoved) Metadata() map[string]interface{} {
	return map[string]interface{}{
		"editorService": "svc",
		"editorUser":    "usr",
		"resourceOwner": "ro",
	}
}

// Revision implements [eventstore.Action]
func (*testUserRemoved) Revision() uint16 { return 1 }

// Payload implements [eventstore.Command]
func (e *testUserRemoved) Payload() interface{} { return nil }

// Options implements [eventstore.Command]
func (e *testUserRemoved) Options() []func(Command) error { return nil }

// Sequence implements [eventstore.Event]
func (e *testUserRemoved) Sequence() uint64 { return 0 }

// CreationDate implements [eventstore.Event]
func (e *testUserRemoved) CreationDate() time.Time { return time.Time{} }

// UnmarshalPayload implements [eventstore.Event]
func (e *testUserRemoved) UnmarshalPayload(object interface{}) error {
	return nil
}

type TestEventstore interface {
	Eventstore
	Before(ctx context.Context, t testing.TB) error
	After(ctx context.Context, t testing.TB) error
}

// func TestEventstore_Push(t *testing.T) {
// 	crdb := startCRDB(t)

// 	type fields struct {
// 		storage Eventstore
// 	}
// 	type args struct {
// 		commands []Command
// 	}
// 	user := new(testUser)
// 	*user = *defaultTestUser
// 	user.id = "2"
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		want    []Event
// 		wantErr bool
// 	}{
// 		{
// 			name: "cockroachdb",
// 			fields: fields{
// 				storage: crdb,
// 			},
// 			args: args{
// 				commands: []Command{
// 					defaultTestUser.toAdded(),
// 					defaultTestUser.toRemoved(),
// 					user.toAdded(),
// 				},
// 			},
// 			want: []Event{
// 				defaultTestUser.toAdded(),
// 				defaultTestUser.toRemoved(),
// 				user.toAdded(),
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "multiple events",
// 			fields: fields{
// 				storage: memory.New(),
// 			},
// 			args: args{
// 				commands: []Command{
// 					defaultTestUser.toAdded(),
// 					defaultTestUser.toRemoved(),
// 				},
// 			},
// 			want: []Event{
// 				defaultTestUser.toAdded(),
// 				defaultTestUser.toRemoved(),
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "multiple aggregates",
// 			fields: fields{
// 				storage: memory.New(),
// 			},
// 			args: args{
// 				commands: []Command{
// 					defaultTestUser.toAdded(),
// 					user.toAdded(),
// 					defaultTestUser.toUsernameChanged(),
// 					user.toFirstNameChanged(),
// 					user.toLastNameChanged(),
// 					defaultTestUser.toRemoved(),
// 				},
// 			},
// 			want: []Event{
// 				defaultTestUser.toAdded(),
// 				user.toAdded(),
// 				defaultTestUser.toUsernameChanged(),
// 				user.toFirstNameChanged(),
// 				user.toLastNameChanged(),
// 				defaultTestUser.toRemoved(),
// 			},
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, err := tt.fields.storage.Push(context.Background(), tt.args.commands...)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("Eventstore.Push() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			assertEvents(t, tt.want, got)
// 		})
// 	}
// }

func PushComplianceTests(ctx context.Context, t *testing.T, store TestEventstore) {
	type args struct {
		commands []Command
	}
	user := new(testUser)
	*user = *defaultTestUser
	user.id = "2"
	tests := []struct {
		name    string
		args    args
		want    []Event
		wantErr bool
	}{
		{
			name: "multiple events",
			args: args{
				commands: []Command{
					defaultTestUser.toAdded(),
					defaultTestUser.toRemoved(),
				},
			},
			want: []Event{
				defaultTestUser.toAdded(),
				defaultTestUser.toRemoved(),
			},
			wantErr: false,
		},
		{
			name: "multiple aggregates",
			args: args{
				commands: []Command{
					defaultTestUser.toAdded(),
					user.toAdded(),
					defaultTestUser.toUsernameChanged(),
					user.toFirstNameChanged(),
					user.toLastNameChanged(),
					defaultTestUser.toRemoved(),
				},
			},
			want: []Event{
				defaultTestUser.toAdded(),
				user.toAdded(),
				defaultTestUser.toUsernameChanged(),
				user.toFirstNameChanged(),
				user.toLastNameChanged(),
				defaultTestUser.toRemoved(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		if err := store.Before(ctx, t); err != nil {
			t.Error("unable to execute store.Before: ", err)
		}
		t.Run(tt.name, func(t *testing.T) {
			got, err := store.Push(ctx, tt.args.commands...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Eventstore.Push() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assertEvents(t, tt.want, got)
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
			user := new(testUser)
			*user = *defaultTestUser
			user.id = b.Name()

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
	b.RunParallel(func(p *testing.PB) {
		for n := 0; p.Next(); n++ {
			user := new(testUser)
			*user = *defaultTestUser
			user.id = b.Name() + strconv.Itoa(n)

			pushDefaultCommands(ctx, b, store, user)
		}
	})
	if err := store.After(ctx, b); err != nil {
		b.Error("unable to execute store.After: ", err)
	}
}

func pushDefaultCommands(ctx context.Context, t testing.TB, store TestEventstore, user *testUser) {
	t.Helper()

	cmds := []Command{
		user.toAdded(),
		user.toRemoved(),
	}

	_, err := store.Push(ctx, cmds...)
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
		want    []Event
		wantErr bool
	}{
		{
			name: "multi token",
			args: args{
				filter: &Filter{
					Action: []Subject{TextSubject("user"), TextSubject("id"), MultiToken},
				},
			},
			want: []Event{
				defaultTestUser.toAdded(),
				defaultTestUser.toRemoved(),
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
			want: []Event{
				defaultTestUser.toAdded(),
				defaultTestUser.toRemoved(),
			},
			wantErr: false,
		},
		{
			name: "all",
			args: args{
				filter: &Filter{
					Action: []Subject{
						TextSubject("user"),
						TextSubject("id"),
						TextSubject("added"),
					},
				},
			},
			want: []Event{
				defaultTestUser.toAdded(),
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
			want: []Event{
				defaultTestUser.toAdded(),
				defaultTestUser.toRemoved(),
			},
			wantErr: false,
		},
	}
	cmds := []Command{
		defaultTestUser.toAdded(),
		defaultTestUser.toRemoved(),
	}
	for _, tt := range tests {
		if err := store.Before(ctx, t); err != nil {
			t.Error("unable to execute store.Before: ", err)
		}
		t.Run(tt.name, func(t *testing.T) {
			_, err := store.Push(ctx, cmds...)
			if err != nil {
				t.Fatalf("unable to push events: %v", err)
			}
			got, err := store.Filter(ctx, tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("Eventstore.Filter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assertEvents(t, tt.want, got)
		})
		if err := store.After(ctx, t); err != nil {
			t.Error("unable to execute store.After: ", err)
		}
	}
}

// func BenchmarkEventstoreFilter(b *testing.B) {
// 	tests := []struct {
// 		name    string
// 		storage Eventstore
// 	}{
// 		{
// 			name:    "memory",
// 			storage: memory.New(),
// 		},
// 	}
// 	for _, tt := range tests {
// 		b.Run(tt.name, func(b *testing.B) {
// 			user := new(testUser)
// 			*user = *defaultTestUser
// 			user.id = "2"
// 			cmds := []Command{
// 				user.toAdded(),
// 				defaultTestUser.toAdded(),
// 				defaultTestUser.toRemoved(),
// 				user.toRemoved(),
// 			}
// 			_, err := tt.storage.Push(context.Background(), cmds...)
// 			if err != nil {
// 				b.Error(err)
// 				b.FailNow()
// 			}
// 			for n := 0; n < b.N; n++ {
// 				events, err := tt.storage.Filter(context.Background(), &Filter{
// 					Limit:  2,
// 					Action: []Subject{TextSubject("user"), SingleToken, TextSubject("added")},
// 				})
// 				if err != nil {
// 					b.Error(err)
// 				}
// 				if len(events) != 2 {
// 					b.Errorf("2 events should be returned got %d", len(events))
// 				}
// 			}
// 		})
// 	}
// }

func assertEvents(t *testing.T, want, got []Event) (failed bool) {
	t.Helper()

	if len(want) != len(got) {
		t.Errorf("unexpected amount of events. want %d, got %d", len(want), len(got))
		return true
	}

	for i := 0; i < len(want); i++ {
		failed = failed || assertEvent(t, want[i], got[i])
	}

	return failed
}

func assertEvent(t *testing.T, want, got Event) (failed bool) {
	t.Helper()

	failed = assertAction(t, want, got)
	failed = failed || assertPayload(t, want, got)

	if want.Sequence() > 0 && want.Sequence() != got.Sequence() {
		failed = true
		t.Errorf("expected sequence %d got: %d", want.Sequence(), got.Sequence())
	}
	if !want.CreationDate().IsZero() && !want.CreationDate().Equal(got.CreationDate()) {
		failed = true
		t.Errorf("expected creation date %v got: %v", want.CreationDate(), got.CreationDate())
	}

	return failed
}

func assertAction(t *testing.T, want, got Action) (failed bool) {
	t.Helper()

	if !reflect.DeepEqual(want.Action(), got.Action()) {
		t.Errorf("expected action %q got: %q", want.Action().Join("."), got.Action().Join("."))
		failed = true
	}
	if !reflect.DeepEqual(want.Aggregate(), got.Aggregate()) {
		t.Errorf("expected aggregate %q got: %q", want.Aggregate().Join("."), got.Aggregate().Join("."))
		failed = true
	}
	if want.Revision() > 0 && want.Revision() != got.Revision() {
		t.Errorf("expected revision %d got: %d", want.Revision(), got.Revision())
		failed = true
	}
	if !reflect.DeepEqual(want.Metadata(), got.Metadata()) {
		t.Errorf("expected metadata %v got: %v", want.Metadata(), got.Metadata())
		failed = true
	}

	return failed
}

func assertPayload(t *testing.T, want, got Event) (failed bool) {
	var (
		gotPayload, wantPayload interface{}
	)
	if err := want.UnmarshalPayload(&wantPayload); err != nil {
		t.Errorf("unable to unmarshal want payload: %v", err)
		failed = true
	}
	if err := got.UnmarshalPayload(&gotPayload); err != nil {
		t.Errorf("unable to unmarshal gotten payload: %v", err)
		failed = true
	}
	if !reflect.DeepEqual(gotPayload, wantPayload) {
		t.Errorf("payload not equal want: %#v got: %#v", wantPayload, gotPayload)
		failed = true
	}

	return failed
}

// func assertCommands(t *testing.T, want, got []Command) (failed bool) {
// 	t.Helper()

// 	if len(want) != len(got) {
// 		t.Errorf("unexpected amount of commands. want %d, got %d", len(want), len(got))
// 		return true
// 	}

// 	for i := 0; i < len(want); i++ {
// 		failed = failed || assertCommand(t, want[i], got[i])
// 	}

// 	return failed
// }

// func assertCommand(t *testing.T, want, got Command) (failed bool) {
// 	t.Helper()

// 	failed = assertAction(t, want, got)
// 	failed = failed || assertCommandOption(t, want.Options(), got.Options())

// 	if !reflect.DeepEqual(want.Payload(), got.Payload()) {
// 		failed = true
// 		t.Errorf("expected payload %#v got: %#v", want.Payload(), got.Payload())
// 	}

// 	return failed
// }

// func assertCommandOption(t *testing.T, want, got []func(Command) error) (failed bool) {
// 	t.Helper()

// 	if len(want) != len(got) {
// 		t.Errorf("unequal length of options: want %d, got %d", len(want), len(got))
// 		return true
// 	}
// 	var gotCmd, wantCmd Command
// 	for i := 0; i < len(want); i++ {
// 		if err := want[i](wantCmd); err != nil {
// 			t.Errorf("wanted option %d failed: %v", i, err)
// 			failed = true
// 		}
// 		if err := got[i](gotCmd); err != nil {
// 			t.Errorf("gotten option %d failed: %v", i, err)
// 			failed = true
// 		}
// 	}

// 	if !reflect.DeepEqual(gotCmd, wantCmd) {
// 		t.Errorf("commands unequal after options: want %#v, got: %#v", wantCmd, gotCmd)
// 		failed = true
// 	}

// 	return failed
// }

// func startCRDB(t testing.TB) *cockroachdb.CockroachDB {
// 	t.Helper()

// 	var ts *testserver.TestServer
// 	_ = ts
// 	// ts, err := testserver.NewTestServer()
// 	// if err != nil {
// 	// 	t.Fatal(err)
// 	// }
// 	// dbpool, err := pgxpool.New(context.Background(), ts.PGURL().String())
// 	dbpool, err := pgxpool.New(context.Background(), "postgresql://root@localhost:26257/weekend?sslmode=disable")

// 	if err != nil {
// 		t.Errorf("unable to create database pool: %v", err)
// 		t.FailNow()
// 	}

// 	crdb := cockroachdb.New(&cockroachdb.Config{
// 		Pool: dbpool,
// 	})
// 	if err := crdb.Setup(context.Background()); err != nil {
// 		t.Fatalf("unable to setup cockroach: %v", err)
// 	}

// 	return crdb
// }
