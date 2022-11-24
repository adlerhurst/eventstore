package eventstore_test

import (
	"context"
	"encoding/json"
	"reflect"
	"strconv"
	"testing"

	"github.com/adlerhurst/eventstore"
	"github.com/adlerhurst/eventstore/storage/memory"
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

type testUserAdded struct {
	eventstore.Event `json:"-"`

	id        string
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Username  string `json:"username,omitempty"`
}

func (e *testUserAdded) EditorService() string { return "svc" }

func (e *testUserAdded) EditorUser() string { return "usr" }

func (e *testUserAdded) Subjects() []eventstore.TextSubject {
	return []eventstore.TextSubject{"user", eventstore.TextSubject(e.id), "added"}
}

func (e *testUserAdded) ResourceOwner() string { return "ro" }

func (e *testUserAdded) Payload() interface{} { return e }

func (u testUser) toFirstNameChanged() *testUserFirstNameChanged {
	return &testUserFirstNameChanged{
		id:        u.id,
		FirstName: u.firstName,
	}
}

type testUserFirstNameChanged struct {
	eventstore.Event `json:"-"`

	id        string
	FirstName string `json:"firstName,omitempty"`
}

func (e *testUserFirstNameChanged) EditorService() string { return "svc" }

func (e *testUserFirstNameChanged) EditorUser() string { return "usr" }

func (e *testUserFirstNameChanged) Subjects() []eventstore.TextSubject {
	return []eventstore.TextSubject{"user", eventstore.TextSubject(e.id), "changed", "firstName"}
}

func (e *testUserFirstNameChanged) ResourceOwner() string { return "ro" }

func (e *testUserFirstNameChanged) Payload() interface{} { return e }

func (u testUser) toLastNameChanged() *testUserLastNameChanged {
	return &testUserLastNameChanged{
		id:       u.id,
		LastName: u.lastName,
	}
}

type testUserLastNameChanged struct {
	eventstore.Event `json:"-"`

	id       string
	LastName string `json:"lastName,omitempty"`
}

func (e *testUserLastNameChanged) EditorService() string { return "svc" }

func (e *testUserLastNameChanged) EditorUser() string { return "usr" }

func (e *testUserLastNameChanged) Subjects() []eventstore.TextSubject {
	return []eventstore.TextSubject{"user", eventstore.TextSubject(e.id), "changed", "lastName"}
}

func (e *testUserLastNameChanged) ResourceOwner() string { return "ro" }

func (e *testUserLastNameChanged) Payload() interface{} { return e }

func (u testUser) toUsernameChanged() *testUsernameChanged {
	return &testUsernameChanged{
		id:       u.id,
		Username: u.username,
	}
}

type testUsernameChanged struct {
	eventstore.Event `json:"-"`

	id       string
	Username string `json:"username,omitempty"`
}

func (e *testUsernameChanged) EditorService() string { return "svc" }

func (e *testUsernameChanged) EditorUser() string { return "usr" }

func (e *testUsernameChanged) Subjects() []eventstore.TextSubject {
	return []eventstore.TextSubject{"user", eventstore.TextSubject(e.id), "changed", "username"}
}

func (e *testUsernameChanged) ResourceOwner() string { return "ro" }

func (e *testUsernameChanged) Payload() interface{} { return e }

func (u testUser) toRemoved() *testUserRemoved {
	return &testUserRemoved{
		id: u.id,
	}
}

type testUserRemoved struct {
	eventstore.Event `json:"-"`

	id string
}

func (e *testUserRemoved) EditorService() string { return "svc" }

func (e *testUserRemoved) EditorUser() string { return "usr" }

func (e *testUserRemoved) Subjects() []eventstore.TextSubject {
	return []eventstore.TextSubject{"user", eventstore.TextSubject(e.id), "removed"}
}

func (e *testUserRemoved) ResourceOwner() string { return "ro" }

func (e *testUserRemoved) Payload() interface{} { return nil }

func TestEventstore_Push(t *testing.T) {
	type fields struct {
		storage eventstore.Storage
	}
	type args struct {
		commands []eventstore.Command
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []eventstore.Event
		wantErr bool
	}{
		{
			name: "multiple events",
			fields: fields{
				storage: memory.New(),
			},
			args: args{
				commands: []eventstore.Command{
					defaultTestUser.toAdded(),
					defaultTestUser.toRemoved(),
				},
			},
			want: []eventstore.Event{
				{
					EditorUser: "usr",
					// ResourceOwner: "ro",
					Subjects: []eventstore.TextSubject{"user", "id", "added"},
					Sequence: 1,
					Payload:  mustJSON(t, defaultTestUser.toAdded()),
				},
				{
					EditorUser: "usr",
					// ResourceOwner: "ro",
					Subjects: []eventstore.TextSubject{"user", "id", "removed"},
					Sequence: 2,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			es := eventstore.New(tt.fields.storage)
			got, err := es.Push(context.Background(), tt.args.commands...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Eventstore.Push() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for i, wanted := range tt.want {
				if !reflect.DeepEqual(got[i], wanted) {
					t.Errorf("Eventstore.Push() %d = %v, want %v", i, got[i], wanted)
				}
			}
		})
	}
}

func BenchmarkEventstorePush(b *testing.B) {
	tests := []struct {
		name    string
		storage eventstore.Storage
	}{
		{
			name:    "memory",
			storage: memory.New(),
		},
	}
	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			es := eventstore.New(tt.storage)
			for n := 0; n < b.N; n++ {
				user := new(testUser)
				*user = *defaultTestUser
				user.id = strconv.Itoa(n)
				cmds := []eventstore.Command{
					user.toAdded(),
					user.toRemoved(),
				}
				_, err := es.Push(context.Background(), cmds...)
				if err != nil {
					b.Error(err)
				}
			}
		})
	}
}

func TestEventstore_Filter(t *testing.T) {
	type fields struct {
		storage eventstore.Storage
	}
	type args struct {
		filter *eventstore.Filter
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []eventstore.Event
		wantErr bool
	}{
		{
			name: "multiple events",
			fields: fields{
				storage: memory.New(),
			},
			args: args{
				filter: &eventstore.Filter{
					Subjects: []*eventstore.SubjectFilter{
						{
							Subjects: []eventstore.Subject{
								eventstore.TextSubject("user"),
								eventstore.TextSubject("id"),
								eventstore.MultiToken,
							},
						},
					},
				},
			},
			want: []eventstore.Event{
				{
					EditorUser: "usr",
					Subjects:   []eventstore.TextSubject{"user", "id", "added"},
					Sequence:   1,
					Payload:    mustJSON(t, defaultTestUser.toAdded()),
				},
				{
					EditorUser: "usr",
					Subjects:   []eventstore.TextSubject{"user", "id", "removed"},
					Sequence:   2,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		ctx := context.Background()
		cmds := []eventstore.Command{
			defaultTestUser.toAdded(),
			defaultTestUser.toRemoved(),
		}
		t.Run(tt.name, func(t *testing.T) {
			es := eventstore.New(tt.fields.storage)
			_, err := es.Push(ctx, cmds...)
			if err != nil {
				t.Fatalf("unable to push events: %v", err)
			}
			got, err := es.Filter(ctx, tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("Eventstore.Filter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(tt.want) != len(got) {
				t.Errorf("unexpected length of filtered events want %d, got %d", len(tt.want), len(got))
				return
			}
			for i, wanted := range tt.want {
				if !reflect.DeepEqual(got[i], wanted) {
					t.Errorf("Eventstore.Filter() %d = %v, want %v", i, got[i], wanted)
				}
			}
		})
	}
}

func BenchmarkEventstoreFilter(b *testing.B) {
	tests := []struct {
		name    string
		storage eventstore.Storage
	}{
		{
			name:    "memory",
			storage: memory.New(),
		},
	}
	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			es := eventstore.New(tt.storage)
			user := new(testUser)
			*user = *defaultTestUser
			user.id = "2"
			cmds := []eventstore.Command{
				user.toAdded(),
				defaultTestUser.toAdded(),
				defaultTestUser.toRemoved(),
				user.toRemoved(),
			}
			_, err := es.Push(context.Background(), cmds...)
			if err != nil {
				b.Error(err)
				b.FailNow()
			}
			for n := 0; n < b.N; n++ {
				events, err := es.Filter(context.Background(), &eventstore.Filter{
					Limit: 2,
					Subjects: []*eventstore.SubjectFilter{
						{
							From: 1,
							Subjects: []eventstore.Subject{
								eventstore.TextSubject("user"),
								eventstore.SingleToken,
								eventstore.TextSubject("added"),
							},
						},
					},
				})
				if err != nil {
					b.Error(err)
				}
				if len(events) != 2 {
					b.Errorf("%d: 2 events should be returned got %d", n, len(events))
				}
			}
		})
	}
}

func mustJSON(t *testing.T, object interface{}) []byte {
	t.Helper()
	data, err := json.Marshal(object)
	if err != nil {
		t.Fatal(err)
	}
	return data
}
