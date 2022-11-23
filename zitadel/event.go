package zitadel

import "time"

type Event struct {
	//CreationDate is the date when the event was written to the store
	CreationDate time.Time
	//Aggregate is the metadata of an aggregate
	Aggregate
	// EditorService is the service who wants to push the event
	EditorService string
	//EditorUser is the user who wants to push the event
	EditorUser string
	//Type must return an event type which should be unique in the aggregate
	Type string
	//Payload returns the payload of the event. It represent the changed fields by the event
	// valid types are:
	// * nil (no payload),
	// * json byte array
	// * struct which can be marshalled to json
	// * pointer to struct which can be marshalled to json
	Payload []byte
}

func EventFromCommand(cmd Command) *Event {
	return &Event{
		Aggregate: Aggregate{
			ID:            cmd.Aggregate().ID,
			Type:          cmd.Aggregate().Type,
			ResourceOwner: cmd.Aggregate().ResourceOwner,
			InstanceID:    cmd.Aggregate().ResourceOwner,
			Version:       cmd.Aggregate().Version,
		},
		EditorService: cmd.EditorService(),
		EditorUser:    cmd.EditorUser(),
		Type:          cmd.Type(),
	}
}