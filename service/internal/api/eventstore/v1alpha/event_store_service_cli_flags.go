// Code generated by protoc-gen-go-cli. DO NOT EDIT.

package eventstorev1alpha

import (
	pflag "github.com/spf13/pflag"
	os "os"
)

type ActionFlag struct {
	*Action

	changed bool
	set     *pflag.FlagSet

	actionFlag   *StringSliceParser
	revisionFlag *Uint32Parser
	payloadFlag  *StructParser
}

func (x *ActionFlag) AddFlags(parent *pflag.FlagSet) {
	x.set = pflag.NewFlagSet("Action", pflag.ContinueOnError)

	x.actionFlag = NewStringSliceFlag(x.set, "action", "")
	x.revisionFlag = NewUint32Flag(x.set, "revision", "")
	x.payloadFlag = NewStructFlag(x.set, "payload", "")
	parent.AddFlagSet(x.set)
}

func (x *ActionFlag) ParseFlags(parent *pflag.FlagSet, args []string) {
	flagIndexes := fieldIndexes(args)

	if err := x.set.Parse(flagIndexes.primitives().args); err != nil {
		DefaultConfig.Logger.Error("failed to parse flags", "cause", err)
		os.Exit(1)
	}

	if x.actionFlag.Changed() {
		x.changed = true
		x.Action.Action = *x.actionFlag.Value
	}
	if x.revisionFlag.Changed() {
		x.changed = true
		x.Revision = *x.revisionFlag.Value
	}
	if x.payloadFlag.Changed() {
		x.changed = true
		x.Payload = x.payloadFlag.Value
	}
}

func (x *ActionFlag) Changed() bool {
	return x.changed
}

type AggregateFlag struct {
	*Aggregate

	changed bool
	set     *pflag.FlagSet

	idFlag              *StringSliceParser
	commandsFlag        []*CommandFlag
	currentSequenceFlag *Uint32Parser
}

func (x *AggregateFlag) AddFlags(parent *pflag.FlagSet) {
	x.set = pflag.NewFlagSet("Aggregate", pflag.ContinueOnError)

	x.idFlag = NewStringSliceFlag(x.set, "id", "")
	x.commandsFlag = []*CommandFlag{}
	x.currentSequenceFlag = NewUint32Flag(x.set, "current-sequence", "")
	parent.AddFlagSet(x.set)
}

func (x *AggregateFlag) ParseFlags(parent *pflag.FlagSet, args []string) {
	flagIndexes := fieldIndexes(args, "commands")

	if err := x.set.Parse(flagIndexes.primitives().args); err != nil {
		DefaultConfig.Logger.Error("failed to parse flags", "cause", err)
		os.Exit(1)
	}

	for _, flagIdx := range flagIndexes.byName("commands") {
		x.commandsFlag = append(x.commandsFlag, &CommandFlag{Command: new(Command)})
		x.commandsFlag[len(x.commandsFlag)-1].AddFlags(x.set)
		x.commandsFlag[len(x.commandsFlag)-1].ParseFlags(x.set, flagIdx.args)
	}

	if x.idFlag.Changed() {
		x.changed = true
		x.Id = *x.idFlag.Value
	}
	if len(x.commandsFlag) > 0 {
		x.changed = true
		x.Commands = make([]*Command, len(x.commandsFlag))
		for i, value := range x.commandsFlag {
			x.Commands[i] = value.Command
		}
	}

	if x.currentSequenceFlag.Changed() {
		x.changed = true
		x.CurrentSequence = x.currentSequenceFlag.Value
	}
}

func (x *AggregateFlag) Changed() bool {
	return x.changed
}

type CommandFlag struct {
	*Command

	changed bool
	set     *pflag.FlagSet

	actionFlag *ActionFlag
}

func (x *CommandFlag) AddFlags(parent *pflag.FlagSet) {
	x.set = pflag.NewFlagSet("Command", pflag.ContinueOnError)

	x.actionFlag = &ActionFlag{Action: new(Action)}
	x.actionFlag.AddFlags(x.set)
	parent.AddFlagSet(x.set)
}

func (x *CommandFlag) ParseFlags(parent *pflag.FlagSet, args []string) {
	flagIndexes := fieldIndexes(args, "action")

	if err := x.set.Parse(flagIndexes.primitives().args); err != nil {
		DefaultConfig.Logger.Error("failed to parse flags", "cause", err)
		os.Exit(1)
	}

	if flagIdx := flagIndexes.lastByName("action"); flagIdx != nil {
		x.actionFlag.ParseFlags(x.set, flagIdx.args)
	}

	if x.actionFlag.Changed() {
		x.changed = true
		x.Action = x.actionFlag.Action
	}

}

func (x *CommandFlag) Changed() bool {
	return x.changed
}

type FilterRequestFlag struct {
	*FilterRequest

	changed bool
	set     *pflag.FlagSet

	queriesFlag []*QueryFlag
	limitFlag   *Uint64Parser
}

func (x *FilterRequestFlag) AddFlags(parent *pflag.FlagSet) {
	x.set = pflag.NewFlagSet("FilterRequest", pflag.ContinueOnError)

	x.queriesFlag = []*QueryFlag{}
	x.limitFlag = NewUint64Flag(x.set, "limit", "")
	parent.AddFlagSet(x.set)
}

func (x *FilterRequestFlag) ParseFlags(parent *pflag.FlagSet, args []string) {
	flagIndexes := fieldIndexes(args, "queries")

	if err := x.set.Parse(flagIndexes.primitives().args); err != nil {
		DefaultConfig.Logger.Error("failed to parse flags", "cause", err)
		os.Exit(1)
	}

	for _, flagIdx := range flagIndexes.byName("queries") {
		x.queriesFlag = append(x.queriesFlag, &QueryFlag{Query: new(Query)})
		x.queriesFlag[len(x.queriesFlag)-1].AddFlags(x.set)
		x.queriesFlag[len(x.queriesFlag)-1].ParseFlags(x.set, flagIdx.args)
	}

	if len(x.queriesFlag) > 0 {
		x.changed = true
		x.Queries = make([]*Query, len(x.queriesFlag))
		for i, value := range x.queriesFlag {
			x.Queries[i] = value.Query
		}
	}

	if x.limitFlag.Changed() {
		x.changed = true
		x.Limit = *x.limitFlag.Value
	}
}

func (x *FilterRequestFlag) Changed() bool {
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
	flagIndexes := fieldIndexes(args, "aggregates")

	if err := x.set.Parse(flagIndexes.primitives().args); err != nil {
		DefaultConfig.Logger.Error("failed to parse flags", "cause", err)
		os.Exit(1)
	}

	for _, flagIdx := range flagIndexes.byName("aggregates") {
		x.aggregatesFlag = append(x.aggregatesFlag, &AggregateFlag{Aggregate: new(Aggregate)})
		x.aggregatesFlag[len(x.aggregatesFlag)-1].AddFlags(x.set)
		x.aggregatesFlag[len(x.aggregatesFlag)-1].ParseFlags(x.set, flagIdx.args)
	}
	if len(x.aggregatesFlag) > 0 {
		x.changed = true
		x.Aggregates = make([]*Aggregate, len(x.aggregatesFlag))
		for i, value := range x.aggregatesFlag {
			x.Aggregates[i] = value.Aggregate
		}
	}

}

func (x *PushRequestFlag) Changed() bool {
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
	flagIndexes := fieldIndexes(args, "subjects", "sequence", "created-at")

	if err := x.set.Parse(flagIndexes.primitives().args); err != nil {
		DefaultConfig.Logger.Error("failed to parse flags", "cause", err)
		os.Exit(1)
	}

	if flagIdx := flagIndexes.lastByName("sequence"); flagIdx != nil {
		x.sequenceFlag.ParseFlags(x.set, flagIdx.args)
	}

	if flagIdx := flagIndexes.lastByName("created-at"); flagIdx != nil {
		x.createdAtFlag.ParseFlags(x.set, flagIdx.args)
	}

	for _, flagIdx := range flagIndexes.byName("subjects") {
		x.subjectsFlag = append(x.subjectsFlag, &SubjectFlag{Subject: new(Subject)})
		x.subjectsFlag[len(x.subjectsFlag)-1].AddFlags(x.set)
		x.subjectsFlag[len(x.subjectsFlag)-1].ParseFlags(x.set, flagIdx.args)
	}

	if len(x.subjectsFlag) > 0 {
		x.changed = true
		x.Subjects = make([]*Subject, len(x.subjectsFlag))
		for i, value := range x.subjectsFlag {
			x.Subjects[i] = value.Subject
		}
	}

	if x.sequenceFlag.Changed() {
		x.changed = true
		x.Sequence = x.sequenceFlag.Query_Sequence
	}

	if x.createdAtFlag.Changed() {
		x.changed = true
		x.CreatedAt = x.createdAtFlag.Query_CreatedAt
	}

}

func (x *QueryFlag) Changed() bool {
	return x.changed
}

type Query_CreatedAtFlag struct {
	*Query_CreatedAt

	changed bool
	set     *pflag.FlagSet

	fromFlag *TimestampParser
	toFlag   *TimestampParser
}

func (x *Query_CreatedAtFlag) AddFlags(parent *pflag.FlagSet) {
	x.set = pflag.NewFlagSet("Query_CreatedAt", pflag.ContinueOnError)

	x.fromFlag = NewTimestampFlag(x.set, "from", "")
	x.toFlag = NewTimestampFlag(x.set, "to", "")
	parent.AddFlagSet(x.set)
}

func (x *Query_CreatedAtFlag) ParseFlags(parent *pflag.FlagSet, args []string) {
	flagIndexes := fieldIndexes(args)

	if err := x.set.Parse(flagIndexes.primitives().args); err != nil {
		DefaultConfig.Logger.Error("failed to parse flags", "cause", err)
		os.Exit(1)
	}

	if x.fromFlag.Changed() {
		x.changed = true
		x.From = x.fromFlag.Value
	}
	if x.toFlag.Changed() {
		x.changed = true
		x.To = x.toFlag.Value
	}
}

func (x *Query_CreatedAtFlag) Changed() bool {
	return x.changed
}

type Query_SequenceFlag struct {
	*Query_Sequence

	changed bool
	set     *pflag.FlagSet

	fromFlag *Uint32Parser
	toFlag   *Uint32Parser
}

func (x *Query_SequenceFlag) AddFlags(parent *pflag.FlagSet) {
	x.set = pflag.NewFlagSet("Query_Sequence", pflag.ContinueOnError)

	x.fromFlag = NewUint32Flag(x.set, "from", "")
	x.toFlag = NewUint32Flag(x.set, "to", "")
	parent.AddFlagSet(x.set)
}

func (x *Query_SequenceFlag) ParseFlags(parent *pflag.FlagSet, args []string) {
	flagIndexes := fieldIndexes(args)

	if err := x.set.Parse(flagIndexes.primitives().args); err != nil {
		DefaultConfig.Logger.Error("failed to parse flags", "cause", err)
		os.Exit(1)
	}

	if x.fromFlag.Changed() {
		x.changed = true
		x.From = *x.fromFlag.Value
	}
	if x.toFlag.Changed() {
		x.changed = true
		x.To = *x.toFlag.Value
	}
}

func (x *Query_SequenceFlag) Changed() bool {
	return x.changed
}

type SubjectFlag struct {
	*Subject

	changed bool
	set     *pflag.FlagSet

	textFlag     *StringParser
	wildcardFlag *EnumParser[Subject_Wildcard]
}

func (x *SubjectFlag) AddFlags(parent *pflag.FlagSet) {
	x.set = pflag.NewFlagSet("Subject", pflag.ContinueOnError)

	x.textFlag = NewStringFlag(x.set, "text", "")
	x.wildcardFlag = NewEnumFlag[Subject_Wildcard](x.set, "wildcard", "")
	parent.AddFlagSet(x.set)
}

func (x *SubjectFlag) ParseFlags(parent *pflag.FlagSet, args []string) {
	flagIndexes := fieldIndexes(args)

	if err := x.set.Parse(flagIndexes.primitives().args); err != nil {
		DefaultConfig.Logger.Error("failed to parse flags", "cause", err)
		os.Exit(1)
	}

	switch fieldIndexes(args, "text", "wildcard").last().flag {
	case "text":
		if x.textFlag.Changed() {
			x.changed = true
			x.Subject.Subject = &Subject_Text{Text: *x.textFlag.Value}
		}
	case "wildcard":
		if x.wildcardFlag.Changed() {
			x.changed = true
			x.Subject.Subject = &Subject_Wildcard_{Wildcard: *x.wildcardFlag.Value}
		}
	}
}

func (x *SubjectFlag) Changed() bool {
	return x.changed
}
