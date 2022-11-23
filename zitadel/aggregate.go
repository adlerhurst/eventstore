package zitadel

// Aggregate is the basic implementation of Aggregater
type Aggregate struct {
	//ID is the unique identitfier of this aggregate
	ID string `json:"-"`
	//Type is the name of the aggregate.
	Type string `json:"-"`
	//ResourceOwner is the org this aggregates belongs to
	ResourceOwner string `json:"-"`
	//InstanceID is the instance this aggregate belongs to
	InstanceID string `json:"-"`
	//Version is the semver this aggregate represents
	Version string `json:"-"`
}
