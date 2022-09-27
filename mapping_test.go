package eventstore_test

// import (
// 	"testing"

// 	"github.com/adlerhurst/eventstore"
// 	"github.com/adlerhurst/eventstore/storage/memory"
// )

// func TestEventstore_Register(t *testing.T) {
// 	type fields struct {
// 		storage eventstore.Storage
// 	}
// 	type register struct {
// 		subs    []eventstore.Subject
// 		mapping func(eventstore.EventBase) eventstore.Event
// 	}
// 	type arg struct {
// 		event     eventstore.EventBase
// 		typeCheck func(*testing.T, eventstore.Event)
// 	}
// 	type test struct {
// 		name   string
// 		fields fields
// 		arg    arg
// 	}
// 	tests := struct {
// 		registers []register
// 		tests     []test
// 	}{
// 		registers: []register{
// 			{
// 				subs: []eventstore.Subject{
// 					eventstore.TextSubject("users"),
// 					eventstore.SingleToken,
// 					eventstore.TextSubject("added"),
// 				},
// 				mapping: func(e eventstore.Event) eventstore.Event {
// 					return &testUserAdded{
// 						id: string(e.Subjects[1]),
// 					}
// 				},
// 			},
// 			{
// 				subs: []eventstore.Subject{
// 					eventstore.TextSubject("users"),
// 					eventstore.SingleToken,
// 					eventstore.TextSubject("removed"),
// 				},
// 				mapping: func(e eventstore.EventBase) eventstore.Event {
// 					return &testUserRemoved{
// 						id: string(e.Subjects[1]),
// 					}
// 				},
// 			},
// 			{
// 				subs: []eventstore.Subject{
// 					eventstore.TextSubject("users"),
// 					eventstore.SingleToken,
// 					eventstore.TextSubject("username"),
// 					eventstore.TextSubject("changed"),
// 				},
// 				mapping: func(e eventstore.EventBase) eventstore.Event {
// 					return &testUsernameChanged{
// 						id: string(e.Subjects[1]),
// 					}
// 				},
// 			},
// 		},
// 		tests: []test{
// 			{
// 				name:   "users.added => EventBase",
// 				fields: fields{storage: memory.New()},
// 				arg: arg{
// 					event: eventstore.EventBase{
// 						Subjects: []eventstore.TextSubject{"users", "added"},
// 					},
// 					typeCheck: func(t *testing.T, e eventstore.Event) {
// 						if _, ok := e.(eventstore.EventBase); !ok {
// 							t.Error("returned type should be eventstore.EventBase")
// 						}
// 					},
// 				},
// 			},
// 			{
// 				name:   "users.123.added => testUserAdded",
// 				fields: fields{storage: memory.New()},
// 				arg: arg{
// 					event: eventstore.EventBase{
// 						Subjects: []eventstore.TextSubject{"users", "123", "added"},
// 					},
// 					typeCheck: func(t *testing.T, e eventstore.Event) {
// 						if _, ok := e.(*testUserAdded); !ok {
// 							t.Error("returned type should be *testUserAdded")
// 						}
// 					},
// 				},
// 			},
// 			{
// 				name:   "users => eventstore.BaseEvent",
// 				fields: fields{storage: memory.New()},
// 				arg: arg{
// 					event: eventstore.EventBase{
// 						Subjects: []eventstore.TextSubject{"users"},
// 					},
// 					typeCheck: func(t *testing.T, e eventstore.Event) {
// 						if _, ok := e.(eventstore.EventBase); !ok {
// 							t.Error("returned type should be eventstore.EventBase")
// 						}
// 					},
// 				},
// 			},
// 			{
// 				name:   "orgs.345.added => eventstore.BaseEvent",
// 				fields: fields{storage: memory.New()},
// 				arg: arg{
// 					event: eventstore.EventBase{
// 						Subjects: []eventstore.TextSubject{"orgs", "345", "added"},
// 					},
// 					typeCheck: func(t *testing.T, e eventstore.Event) {
// 						if _, ok := e.(eventstore.EventBase); !ok {
// 							t.Error("returned type should be eventstore.EventBase")
// 						}
// 					},
// 				},
// 			},
// 		},
// 	}
// 	for _, tt := range tests.tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			es := eventstore.New(tt.fields.storage)
// 			for _, r := range tests.registers {
// 				es.RegisterEvent(r.subs, r.mapping)
// 			}
// 			e := es.MapEvent(tt.arg.event)
// 			tt.arg.typeCheck(t, e)
// 		})
// 	}
// }
