package zitadel

// Command is the intend to store an event into the eventstore
type Command interface {
	//Aggregate is the metadata of an aggregate
	Aggregate() Aggregate
	// EditorService is the service who wants to push the event
	EditorService() string
	//EditorUser is the user who wants to push the event
	EditorUser() string
	//Type must return an event type which should be unique in the aggregate
	Type() string
	//Payload returns the payload of the event. It represent the changed fields by the event
	// valid types are:
	// * nil (no payload),
	// * json byte array
	// * struct which can be marshalled to json
	// * pointer to struct which can be marshalled to json
	Payload() interface{}
	//TODO: UniqueConstraints should be added for unique attributes of an event, if nil constraints will not be checked
	// UniqueConstraints() []*EventUniqueConstraint
}