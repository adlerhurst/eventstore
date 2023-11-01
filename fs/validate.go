package fs

// import (
// 	"context"
// 	"errors"
// 	"os"
// 	"strings"

// 	"github.com/adlerhurst/eventstore"
// )

// var _ eventstore.Eventstore = (*Validator)(nil)

// var (
// 	InvalidAggregateIDErr = errors.New("invalid aggregate id")
// )

// func NewValidator(storage eventstore.Eventstore) *Validator {
// 	return &Validator{
// 		storage: storage,
// 	}
// }

// type Validator struct {
// 	storage eventstore.Eventstore
// }

// // Filter implements eventstore.Eventstore.
// func (v *Validator) Filter(ctx context.Context, f *eventstore.Filter) ([]eventstore.Event, error) {
// 	return v.storage.Filter(ctx, f)
// }

// // Push implements eventstore.Eventstore.
// // It checks if the aggreate.ID() contains file path attributes like ".", "..", "/"
// func (v *Validator) Push(ctx context.Context, aggregates ...eventstore.Aggregate) ([]eventstore.Event, error) {
// 	for _, aggregate := range aggregates {
// 		for _, field := range aggregate.ID() {
// 			if field == "." || field == ".." || strings.Contains(string(field), string(os.PathSeparator)) {
// 				return nil, InvalidAggregateIDErr
// 			}
// 		}
// 	}
// 	return v.storage.Push(ctx, aggregates...)
// }

// // Ready implements eventstore.Eventstore.
// func (v *Validator) Ready(ctx context.Context) error {
// 	return v.storage.Ready(ctx)
// }
