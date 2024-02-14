// Code generated by protoc-gen-cli-client. DO NOT EDIT.

package v1alpha

import (
	cli_client "github.com/adlerhurst/cli-client"
	pflag "github.com/spf13/pflag"
	os "os"
)

type FilterRequestFlag struct {
	*FilterRequest

	changed bool
	set     *pflag.FlagSet

	queriesFlag []*QueryFlag
	limitFlag   *cli_client.Uint64Parser
}

func (x *FilterRequestFlag) AddFlags(parent *pflag.FlagSet) {
	x.set = pflag.NewFlagSet("FilterRequest", pflag.ContinueOnError)

	x.queriesFlag = []*QueryFlag{}
	x.limitFlag = cli_client.NewUint64Parser(x.set, "limit", "")
	parent.AddFlagSet(x.set)
}

func (x *FilterRequestFlag) ParseFlags(parent *pflag.FlagSet, args []string) {
	flagIndexes := cli_client.FieldIndexes(args, "queries")

	if err := x.set.Parse(flagIndexes.Primitives().Args); err != nil {
		cli_client.Logger().Error("failed to parse flags", "cause", err)
		os.Exit(1)
	}

	for _, flagIdx := range flagIndexes.ByName("queries") {
		x.queriesFlag = append(x.queriesFlag, &QueryFlag{Query: new(Query)})
		x.queriesFlag[len(x.queriesFlag)-1].AddFlags(x.set)
		x.queriesFlag[len(x.queriesFlag)-1].ParseFlags(x.set, flagIdx.Args)
	}

	if len(x.queriesFlag) > 0 {
		x.changed = true
		x.Queries = make([]*Query, len(x.queriesFlag))
		for i, value := range x.queriesFlag {
			x.FilterRequest.Queries[i] = value.Query
		}
	}

	if x.limitFlag.Changed() {
		x.changed = true
		x.FilterRequest.Limit = *x.limitFlag.Value
	}
}

func (x *FilterRequestFlag) Changed() bool {
	return x.changed
}

type FilterResponseFlag struct {
	*FilterResponse

	changed bool
	set     *pflag.FlagSet

	eventsFlag []*EventFlag
}

func (x *FilterResponseFlag) AddFlags(parent *pflag.FlagSet) {
	x.set = pflag.NewFlagSet("FilterResponse", pflag.ContinueOnError)

	x.eventsFlag = []*EventFlag{}
	parent.AddFlagSet(x.set)
}

func (x *FilterResponseFlag) ParseFlags(parent *pflag.FlagSet, args []string) {
	flagIndexes := cli_client.FieldIndexes(args, "events")

	if err := x.set.Parse(flagIndexes.Primitives().Args); err != nil {
		cli_client.Logger().Error("failed to parse flags", "cause", err)
		os.Exit(1)
	}

	for _, flagIdx := range flagIndexes.ByName("events") {
		x.eventsFlag = append(x.eventsFlag, &EventFlag{Event: new(Event)})
		x.eventsFlag[len(x.eventsFlag)-1].AddFlags(x.set)
		x.eventsFlag[len(x.eventsFlag)-1].ParseFlags(x.set, flagIdx.Args)
	}
	if len(x.eventsFlag) > 0 {
		x.changed = true
		x.Events = make([]*Event, len(x.eventsFlag))
		for i, value := range x.eventsFlag {
			x.FilterResponse.Events[i] = value.Event
		}
	}

}

func (x *FilterResponseFlag) Changed() bool {
	return x.changed
}

type PushRequestFlag struct {
	*PushRequest

	changed bool
	set     *pflag.FlagSet

	aggregatesFlag []*AggregateFlag
}

func (x *PushRequestFlag) AddFlags(parent *pflag.FlagSet) {
	x.set = pflag.NewFlagSet("PushRequest", pflag.ContinueOnError)

	x.aggregatesFlag = []*AggregateFlag{}
	parent.AddFlagSet(x.set)
}

func (x *PushRequestFlag) ParseFlags(parent *pflag.FlagSet, args []string) {
	flagIndexes := cli_client.FieldIndexes(args, "aggregates")

	if err := x.set.Parse(flagIndexes.Primitives().Args); err != nil {
		cli_client.Logger().Error("failed to parse flags", "cause", err)
		os.Exit(1)
	}

	for _, flagIdx := range flagIndexes.ByName("aggregates") {
		x.aggregatesFlag = append(x.aggregatesFlag, &AggregateFlag{Aggregate: new(Aggregate)})
		x.aggregatesFlag[len(x.aggregatesFlag)-1].AddFlags(x.set)
		x.aggregatesFlag[len(x.aggregatesFlag)-1].ParseFlags(x.set, flagIdx.Args)
	}
	if len(x.aggregatesFlag) > 0 {
		x.changed = true
		x.Aggregates = make([]*Aggregate, len(x.aggregatesFlag))
		for i, value := range x.aggregatesFlag {
			x.PushRequest.Aggregates[i] = value.Aggregate
		}
	}

}

func (x *PushRequestFlag) Changed() bool {
	return x.changed
}

type PushResponseFlag struct {
	*PushResponse

	changed bool
	set     *pflag.FlagSet
}

func (x *PushResponseFlag) AddFlags(parent *pflag.FlagSet) {
	x.set = pflag.NewFlagSet("PushResponse", pflag.ContinueOnError)

	parent.AddFlagSet(x.set)
}

func (x *PushResponseFlag) ParseFlags(parent *pflag.FlagSet, args []string) {
	flagIndexes := cli_client.FieldIndexes(args)

	if err := x.set.Parse(flagIndexes.Primitives().Args); err != nil {
		cli_client.Logger().Error("failed to parse flags", "cause", err)
		os.Exit(1)
	}

}

func (x *PushResponseFlag) Changed() bool {
	return x.changed
}

type QueryFlag struct {
	*Query

	changed bool
	set     *pflag.FlagSet

	subjectsFlag  []*SubjectFlag
	sequenceFlag  *Query_SequenceFlag
	createdAtFlag *Query_CreatedAtFlag
}

func (x *QueryFlag) AddFlags(parent *pflag.FlagSet) {
	x.set = pflag.NewFlagSet("Query", pflag.ContinueOnError)

	x.subjectsFlag = []*SubjectFlag{}
	x.sequenceFlag = &Query_SequenceFlag{Query_Sequence: new(Query_Sequence)}
	x.sequenceFlag.AddFlags(x.set)
	x.createdAtFlag = &Query_CreatedAtFlag{Query_CreatedAt: new(Query_CreatedAt)}
	x.createdAtFlag.AddFlags(x.set)
	parent.AddFlagSet(x.set)
}

func (x *QueryFlag) ParseFlags(parent *pflag.FlagSet, args []string) {
	flagIndexes := cli_client.FieldIndexes(args, "subjects", "sequence", "created-at")

	if err := x.set.Parse(flagIndexes.Primitives().Args); err != nil {
		cli_client.Logger().Error("failed to parse flags", "cause", err)
		os.Exit(1)
	}

	if flagIdx := flagIndexes.LastByName("sequence"); flagIdx != nil {
		x.sequenceFlag.ParseFlags(x.set, flagIdx.Args)
	}

	if flagIdx := flagIndexes.LastByName("created-at"); flagIdx != nil {
		x.createdAtFlag.ParseFlags(x.set, flagIdx.Args)
	}

	for _, flagIdx := range flagIndexes.ByName("subjects") {
		x.subjectsFlag = append(x.subjectsFlag, &SubjectFlag{Subject: new(Subject)})
		x.subjectsFlag[len(x.subjectsFlag)-1].AddFlags(x.set)
		x.subjectsFlag[len(x.subjectsFlag)-1].ParseFlags(x.set, flagIdx.Args)
	}

	if len(x.subjectsFlag) > 0 {
		x.changed = true
		x.Subjects = make([]*Subject, len(x.subjectsFlag))
		for i, value := range x.subjectsFlag {
			x.Query.Subjects[i] = value.Subject
		}
	}

	if x.sequenceFlag.Changed() {
		x.changed = true
		x.Query.Sequence = x.sequenceFlag.Query_Sequence
	}

	if x.createdAtFlag.Changed() {
		x.changed = true
		x.Query.CreatedAt = x.createdAtFlag.Query_CreatedAt
	}

}

func (x *QueryFlag) Changed() bool {
	return x.changed
}

type Query_CreatedAtFlag struct {
	*Query_CreatedAt

	changed bool
	set     *pflag.FlagSet

	fromFlag *cli_client.TimestampParser
	toFlag   *cli_client.TimestampParser
}

func (x *Query_CreatedAtFlag) AddFlags(parent *pflag.FlagSet) {
	x.set = pflag.NewFlagSet("Query_CreatedAt", pflag.ContinueOnError)

	x.fromFlag = cli_client.NewTimestampParser(x.set, "from", "")
	x.toFlag = cli_client.NewTimestampParser(x.set, "to", "")
	parent.AddFlagSet(x.set)
}

func (x *Query_CreatedAtFlag) ParseFlags(parent *pflag.FlagSet, args []string) {
	flagIndexes := cli_client.FieldIndexes(args)

	if err := x.set.Parse(flagIndexes.Primitives().Args); err != nil {
		cli_client.Logger().Error("failed to parse flags", "cause", err)
		os.Exit(1)
	}

	if x.fromFlag.Changed() {
		x.changed = true
		x.Query_CreatedAt.From = x.fromFlag.Value
	}
	if x.toFlag.Changed() {
		x.changed = true
		x.Query_CreatedAt.To = x.toFlag.Value
	}
}

func (x *Query_CreatedAtFlag) Changed() bool {
	return x.changed
}

type Query_SequenceFlag struct {
	*Query_Sequence

	changed bool
	set     *pflag.FlagSet

	fromFlag *cli_client.Uint32Parser
	toFlag   *cli_client.Uint32Parser
}

func (x *Query_SequenceFlag) AddFlags(parent *pflag.FlagSet) {
	x.set = pflag.NewFlagSet("Query_Sequence", pflag.ContinueOnError)

	x.fromFlag = cli_client.NewUint32Parser(x.set, "from", "")
	x.toFlag = cli_client.NewUint32Parser(x.set, "to", "")
	parent.AddFlagSet(x.set)
}

func (x *Query_SequenceFlag) ParseFlags(parent *pflag.FlagSet, args []string) {
	flagIndexes := cli_client.FieldIndexes(args)

	if err := x.set.Parse(flagIndexes.Primitives().Args); err != nil {
		cli_client.Logger().Error("failed to parse flags", "cause", err)
		os.Exit(1)
	}

	if x.fromFlag.Changed() {
		x.changed = true
		x.Query_Sequence.From = *x.fromFlag.Value
	}
	if x.toFlag.Changed() {
		x.changed = true
		x.Query_Sequence.To = *x.toFlag.Value
	}
}

func (x *Query_SequenceFlag) Changed() bool {
	return x.changed
}
