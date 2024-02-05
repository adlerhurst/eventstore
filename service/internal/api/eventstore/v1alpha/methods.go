package eventstorev1alpha

import (
	"encoding/csv"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type argParser[T any] struct {
	*primitiveParser[T]
	customParser func(field *T, arg string) error

	changed bool
}

func newArgParser[T any](set *pflag.FlagSet, name string, opts ...primitiveParserOpt[T]) *argParser[T] {
	parser := new(argParser[T])
	parser.primitiveParser = newPrimitiveParser[T](set, name, opts...)

	return parser
}

func (parser *argParser[T]) Changed() bool {
	return parser.changed
}

type primitiveParserOpt[T any] func(*primitiveParser[T])

func WithDefaultValue[T any](value T) primitiveParserOpt[T] {
	return func(parser *primitiveParser[T]) {
		parser.defaultValue = value
	}
}

type primitiveParser[T any] struct {
	Value        *T
	defaultValue T

	set  *pflag.FlagSet
	name string
}

func newPrimitiveParser[T any](set *pflag.FlagSet, name string, opts ...primitiveParserOpt[T]) *primitiveParser[T] {
	parser := &primitiveParser[T]{
		set:  set,
		name: name,
	}

	for _, opt := range opts {
		opt(parser)
	}

	return parser
}

func (parser *primitiveParser[T]) applyOpts(opts []primitiveParserOpt[T]) {
	for _, opt := range opts {
		opt(parser)
	}
}

func (parser *primitiveParser[T]) Changed() bool {
	return parser.set.Changed(parser.name)
}

// Set implements pflag.Value.
func (v *argParser[T]) Set(arg string) error {
	v.changed = true
	if v.customParser != nil {
		return v.customParser(v.Value, arg)
	}

	value, ok := interface{}(v.Value).(protoreflect.ProtoMessage)
	if !ok {
		DefaultConfig.Logger.Error("must implement custom parser", "type", fmt.Sprintf("%T", v.Value))
	}
	return protojson.UnmarshalOptions{
		// AllowPartial: true,
		DiscardUnknown: true,
	}.Unmarshal([]byte(arg), value)
}

// String implements pflag.Value.
func (v *argParser[T]) String() string {
	value, ok := interface{}(v.Value).(protoreflect.ProtoMessage)
	if !ok {
		return fmt.Sprint(v.Value)
	}
	return protojson.Format(value)
}

// Type implements pflag.Value.
func (v *argParser[T]) Type() string {
	value, ok := interface{}(v.Value).(protoreflect.ProtoMessage)
	if !ok {
		return fmt.Sprintf("%T", v.Value)
	}
	return string(value.ProtoReflect().Type().Descriptor().FullName())
}

type StructParser struct {
	*argParser[structpb.Struct]
}

func NewStructFlag(set *pflag.FlagSet, name, usage string) *StructParser {
	parser := newArgParser[structpb.Struct](set, name)
	parser.Value = new(structpb.Struct)
	set.Var(parser, name, usage)
	return &StructParser{argParser: parser}
}

type StructSliceParser struct {
	*argParser[[]*structpb.Struct]
}

func NewStructSliceFlag(set *pflag.FlagSet, name, usage string, opts ...primitiveParserOpt[[]*structpb.Struct]) *StructSliceParser {
	parser := newArgParser[[]*structpb.Struct](set, name)
	parser.applyOpts(opts)
	parser.Value = new([]*structpb.Struct)
	set.Var(parser, name, usage)
	return &StructSliceParser{argParser: parser}
}

type AnyParser struct {
	*argParser[anypb.Any]
}

func NewAnyFlag(set *pflag.FlagSet, name, usage string) *AnyParser {
	parser := newArgParser[anypb.Any](set, name)
	// TODO: change to message
	parser.Value = new(anypb.Any)
	set.Var(parser, name, usage)
	return &AnyParser{argParser: parser}
}

type TimestampParser struct {
	*argParser[timestamppb.Timestamp]
}

func NewTimestampFlag(set *pflag.FlagSet, name, usage string) *TimestampParser {
	parser := newArgParser[timestamppb.Timestamp](set, name)
	parser.Value = new(timestamppb.Timestamp)
	parser.customParser = timestampParser
	set.Var(parser, name, usage)
	return &TimestampParser{argParser: parser}
}

type TimestampSliceParser struct {
	*argParser[[]*timestamppb.Timestamp]
}

func NewTimestampSliceFlag(set *pflag.FlagSet, name, usage string, opts ...primitiveParserOpt[[]*timestamppb.Timestamp]) *TimestampSliceParser {
	parser := newArgParser[[]*timestamppb.Timestamp](set, name)
	parser.applyOpts(opts)
	parser.Value = new([]*timestamppb.Timestamp)
	parser.customParser = slicePtrParser[timestamppb.Timestamp](timestampParser)
	set.Var(parser, name, usage)
	return &TimestampSliceParser{argParser: parser}
}

func timestampParser(field *timestamppb.Timestamp, arg string) error {
	timestamp, err := time.Parse(time.RFC3339, arg)
	if err != nil {
		return err
	}
	*field = *timestamppb.New(timestamp)
	return nil
}

type DurationParser struct {
	*argParser[durationpb.Duration]
}

func NewDurationFlag(set *pflag.FlagSet, name, usage string) *DurationParser {
	parser := newArgParser[durationpb.Duration](set, name)
	parser.Value = new(durationpb.Duration)
	parser.customParser = durationParser
	set.Var(parser, name, usage)
	return &DurationParser{argParser: parser}
}

type DurationSliceParser struct {
	*argParser[[]*durationpb.Duration]
}

func NewDurationSliceFlag(set *pflag.FlagSet, name, usage string, opts ...primitiveParserOpt[[]*durationpb.Duration]) *DurationSliceParser {
	parser := newArgParser[[]*durationpb.Duration](set, name)
	parser.applyOpts(opts)
	parser.Value = new([]*durationpb.Duration)
	parser.customParser = slicePtrParser[durationpb.Duration](durationParser)
	set.Var(parser, name, usage)
	return &DurationSliceParser{argParser: parser}
}

func durationParser(field *durationpb.Duration, arg string) error {
	duration, err := time.ParseDuration(arg)
	if err != nil {
		return err
	}
	*field = *durationpb.New(duration)
	return nil
}

type enum interface {
	~int32
	Descriptor() protoreflect.EnumDescriptor
	String() string
}

type EnumParser[E enum] struct {
	*argParser[E]
}

func NewEnumFlag[E enum](set *pflag.FlagSet, name, usage string, opts ...primitiveParserOpt[E]) *EnumParser[E] {
	parser := newArgParser[E](set, name)
	parser.applyOpts(opts)
	parser.Value = new(E)
	parser.customParser = enumParser[E]
	set.Var(parser, name, usage)
	return &EnumParser[E]{argParser: parser}
}

type EnumSliceParser[E enum] struct {
	*argParser[[]E]
}

func NewEnumSliceFlag[E enum](set *pflag.FlagSet, name, usage string, opts ...primitiveParserOpt[[]E]) *EnumSliceParser[E] {
	parser := newArgParser[[]E](set, name)
	parser.applyOpts(opts)
	parser.Value = new([]E)
	parser.customParser = sliceParser(enumParser[E])
	set.Var(parser, name, usage)
	return &EnumSliceParser[E]{argParser: parser}
}

func sliceParser[T any](parser func(*T, string) error) func(*[]T, string) error {
	return func(field *[]T, arg string) error {
		stringReader := strings.NewReader(arg)
		csvReader := csv.NewReader(stringReader)
		records, err := csvReader.Read()
		if err != nil {
			return err
		}

		values := make([]T, len(records))
		for i, record := range records {
			value := new(T)
			err := parser(value, record)
			if err != nil {
				return err
			}
			values[i] = *value
		}
		*field = append(*field, values...)

		return nil
	}
}

func slicePtrParser[T any](parser func(*T, string) error) func(*[]*T, string) error {
	return func(field *[]*T, arg string) error {
		stringReader := strings.NewReader(arg)
		csvReader := csv.NewReader(stringReader)
		records, err := csvReader.Read()
		if err != nil {
			return err
		}

		values := make([]*T, len(records))
		for i, record := range records {
			value := new(T)
			err := parser(value, record)
			if err != nil {
				return err
			}
			values[i] = value
		}
		*field = append(*field, values...)

		return nil
	}
}

func enumParser[E enum](field *E, arg string) error {
	if desc := (*field).Descriptor().Values().ByName(protoreflect.Name(arg)); desc != nil {
		*field = E(desc.Number())
		return nil
	}
	if number, err := strconv.Atoi(arg); err == nil {
		if desc := (*field).Descriptor().Values().ByNumber(protoreflect.EnumNumber(number)); desc != nil {
			*field = E(desc.Number())
			return nil
		}
	}

	return errors.New("unknown enum variable")
}

type StringParser struct {
	*primitiveParser[string]
}

func NewStringFlag(set *pflag.FlagSet, name, usage string, opts ...primitiveParserOpt[string]) *StringParser {
	parser := newPrimitiveParser[string](set, name, opts...)
	parser.Value = new(string)
	set.StringVar(parser.Value, name, parser.defaultValue, usage)
	return &StringParser{primitiveParser: parser}
}

type StringSliceParser struct {
	*primitiveParser[[]string]
}

func NewStringSliceFlag(set *pflag.FlagSet, name, usage string, opts ...primitiveParserOpt[[]string]) *StringSliceParser {
	parser := newPrimitiveParser[[]string](set, name, opts...)
	parser.Value = new([]string)
	set.StringSliceVar(parser.Value, name, parser.defaultValue, usage)
	return &StringSliceParser{primitiveParser: parser}
}

type BoolParser struct {
	*primitiveParser[bool]
}

func NewBoolFlag(set *pflag.FlagSet, name, usage string, opts ...primitiveParserOpt[bool]) *BoolParser {
	parser := newPrimitiveParser[bool](set, name, opts...)
	parser.Value = new(bool)
	set.BoolVar(parser.Value, name, parser.defaultValue, usage)
	return &BoolParser{primitiveParser: parser}
}

type BoolSliceParser struct {
	*primitiveParser[[]bool]
}

func NewBoolSliceFlag(set *pflag.FlagSet, name, usage string, opts ...primitiveParserOpt[[]bool]) *BoolSliceParser {
	parser := newPrimitiveParser[[]bool](set, name, opts...)
	parser.Value = new([]bool)
	set.BoolSliceVar(parser.Value, name, parser.defaultValue, usage)
	return &BoolSliceParser{primitiveParser: parser}
}

type Int32Parser struct {
	*primitiveParser[int32]
}

func NewInt32Flag(set *pflag.FlagSet, name, usage string, opts ...primitiveParserOpt[int32]) *Int32Parser {
	parser := newPrimitiveParser[int32](set, name, opts...)
	parser.Value = new(int32)
	set.Int32Var(parser.Value, name, parser.defaultValue, usage)
	return &Int32Parser{primitiveParser: parser}
}

type Int32SliceParser struct {
	*primitiveParser[[]int32]
}

func NewInt32SliceFlag(set *pflag.FlagSet, name, usage string, opts ...primitiveParserOpt[[]int32]) *Int32SliceParser {
	parser := newPrimitiveParser[[]int32](set, name, opts...)
	parser.Value = new([]int32)
	set.Int32SliceVar(parser.Value, name, parser.defaultValue, usage)
	return &Int32SliceParser{primitiveParser: parser}
}

type Sint32Parser struct {
	*primitiveParser[int32]
}

func NewSint32Flag(set *pflag.FlagSet, name, usage string, opts ...primitiveParserOpt[int32]) *Sint32Parser {
	return (*Sint32Parser)(NewInt32Flag(set, name, usage, opts...))
}

type Sfixed32Parser struct {
	*primitiveParser[int32]
}

func NewSfixed32Flag(set *pflag.FlagSet, name, usage string, opts ...primitiveParserOpt[int32]) *Sfixed32Parser {
	return (*Sfixed32Parser)(NewInt32Flag(set, name, usage, opts...))
}

type Uint32Parser struct {
	*primitiveParser[uint32]
}

func NewUint32Flag(set *pflag.FlagSet, name, usage string, opts ...primitiveParserOpt[uint32]) *Uint32Parser {
	parser := newPrimitiveParser[uint32](set, name, opts...)
	parser.Value = new(uint32)
	set.Uint32Var(parser.Value, name, parser.defaultValue, usage)
	return &Uint32Parser{primitiveParser: parser}
}

type Fixed32Parser struct {
	*primitiveParser[uint32]
}

func NewFixed32Flag(set *pflag.FlagSet, name, usage string, opts ...primitiveParserOpt[uint32]) *Fixed32Parser {
	return (*Fixed32Parser)(NewUint32Flag(set, name, usage, opts...))
}

type Uint32SliceParser struct {
	*primitiveParser[[]uint]
}

func NewUint32SliceFlag(set *pflag.FlagSet, name, usage string, opts ...primitiveParserOpt[[]uint]) *Uint32SliceParser {
	return &Uint32SliceParser{primitiveParser: newUintSliceFlag(set, name, usage, opts...)}
}

type Int64Parser struct {
	*primitiveParser[int64]
}

func NewInt64Flag(set *pflag.FlagSet, name, usage string, opts ...primitiveParserOpt[int64]) *Int64Parser {
	parser := newPrimitiveParser[int64](set, name, opts...)
	parser.Value = new(int64)
	set.Int64Var(parser.Value, name, parser.defaultValue, usage)
	return &Int64Parser{primitiveParser: parser}
}

type Sint64Parser struct {
	*primitiveParser[int64]
}

func NewSint64Flag(set *pflag.FlagSet, name, usage string, opts ...primitiveParserOpt[int64]) *Sint64Parser {
	return (*Sint64Parser)(NewInt64Flag(set, name, usage, opts...))
}

type Sfixed64Parser struct {
	*primitiveParser[int64]
}

func NewSfixed64Flag(set *pflag.FlagSet, name, usage string, opts ...primitiveParserOpt[int64]) *Sfixed64Parser {
	return (*Sfixed64Parser)(NewInt64Flag(set, name, usage, opts...))
}

type Int64SliceParser struct {
	*primitiveParser[[]int64]
}

func NewInt64SliceFlag(set *pflag.FlagSet, name, usage string, opts ...primitiveParserOpt[[]int64]) *Int64SliceParser {
	parser := newPrimitiveParser[[]int64](set, name, opts...)
	parser.Value = new([]int64)
	set.Int64SliceVar(parser.Value, name, parser.defaultValue, usage)
	return &Int64SliceParser{primitiveParser: parser}
}

type Uint64Parser struct {
	*primitiveParser[uint64]
}

func NewUint64Flag(set *pflag.FlagSet, name, usage string, opts ...primitiveParserOpt[uint64]) *Uint64Parser {
	parser := newPrimitiveParser[uint64](set, name, opts...)
	parser.Value = new(uint64)
	set.Uint64Var(parser.Value, name, parser.defaultValue, usage)
	return &Uint64Parser{primitiveParser: parser}
}

type Fixed64Parser struct {
	*primitiveParser[uint64]
}

func NewFixed64Flag(set *pflag.FlagSet, name, usage string, opts ...primitiveParserOpt[uint64]) *Fixed64Parser {
	return (*Fixed64Parser)(NewUint64Flag(set, name, usage, opts...))
}

type Uint64SliceParser struct {
	*primitiveParser[[]uint]
}

func NewUint64SliceFlag(set *pflag.FlagSet, name, usage string, opts ...primitiveParserOpt[[]uint]) *Uint64SliceParser {
	return &Uint64SliceParser{primitiveParser: newUintSliceFlag(set, name, usage, opts...)}
}

func newUintSliceFlag(set *pflag.FlagSet, name, usage string, opts ...primitiveParserOpt[[]uint]) *primitiveParser[[]uint] {
	parser := newPrimitiveParser[[]uint](set, name, opts...)
	parser.Value = new([]uint)
	set.UintSliceVar(parser.Value, name, parser.defaultValue, usage)
	return parser
}

type FloatParser struct {
	*primitiveParser[float32]
}

func NewFloatFlag(set *pflag.FlagSet, name, usage string, opts ...primitiveParserOpt[float32]) *FloatParser {
	parser := newPrimitiveParser[float32](set, name, opts...)
	parser.Value = new(float32)
	set.Float32Var(parser.Value, name, parser.defaultValue, usage)
	return &FloatParser{primitiveParser: parser}
}

type FloatSliceParser struct {
	*primitiveParser[[]float32]
}

func NewFloatSliceFlag(set *pflag.FlagSet, name, usage string, opts ...primitiveParserOpt[[]float32]) *FloatSliceParser {
	parser := newPrimitiveParser[[]float32](set, name, opts...)
	parser.Value = new([]float32)
	set.Float32SliceVar(parser.Value, name, parser.defaultValue, usage)
	return &FloatSliceParser{primitiveParser: parser}
}

type DoubleParser struct {
	*primitiveParser[float64]
}

func NewDoubleFlag(set *pflag.FlagSet, name, usage string, opts ...primitiveParserOpt[float64]) *DoubleParser {
	parser := newPrimitiveParser[float64](set, name, opts...)
	parser.Value = new(float64)
	set.Float64Var(parser.Value, name, parser.defaultValue, usage)
	return &DoubleParser{primitiveParser: parser}
}

type DoubleSliceParser struct {
	*primitiveParser[[]float64]
}

func NewDoubleSliceFlag(set *pflag.FlagSet, name, usage string, opts ...primitiveParserOpt[[]float64]) *DoubleSliceParser {
	parser := newPrimitiveParser[[]float64](set, name, opts...)
	parser.Value = new([]float64)
	set.Float64SliceVar(parser.Value, name, parser.defaultValue, usage)
	return &DoubleSliceParser{primitiveParser: parser}
}

type BytesParser struct {
	*primitiveParser[[]byte]
}

func NewBytesFlag(set *pflag.FlagSet, name, usage string, opts ...primitiveParserOpt[[]byte]) *BytesParser {
	parser := newPrimitiveParser[[]byte](set, name, opts...)
	parser.Value = new([]byte)
	set.BytesBase64Var(parser.Value, name, parser.defaultValue, usage)
	return &BytesParser{primitiveParser: parser}
}
