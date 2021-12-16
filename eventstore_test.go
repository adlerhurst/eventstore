package eventstore_test

import (
	"context"
	"fmt"
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

func (u testUser) toAdded() *testUserAdded {
	return &testUserAdded{
		id:        u.id,
		FirstName: u.firstName,
		LastName:  u.lastName,
		Username:  u.username,
	}
}

type testUserAdded struct {
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
	id string
}

func (e *testUserRemoved) EditorService() string { return "svc" }

func (e *testUserRemoved) EditorUser() string { return "usr" }

func (e *testUserRemoved) Subjects() []eventstore.TextSubject {
	return []eventstore.TextSubject{"user", eventstore.TextSubject(e.id), "removed"}
}

func (e *testUserRemoved) ResourceOwner() string { return "ro" }

func (e *testUserRemoved) Payload() interface{} { return e }

func TestPush(t *testing.T) {
	es := eventstore.New(memory.New())
	user := testUser{
		id:        "1",
		firstName: "silvan",
		lastName:  "reusser",
		username:  "adlerhurst",
	}
	es.Push(context.Background(),
		user.toAdded(),
		user.toRemoved(),
	)
	user.id = "2"
	es.Push(context.Background(),
		user.toAdded(),
	)
	events, err := es.Filter(context.Background(), eventstore.TextSubject("user"), eventstore.TextSubject("2"), eventstore.MultiToken)
	fmt.Println(events, err)
}
