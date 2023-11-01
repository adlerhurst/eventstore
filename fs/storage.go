package fs

// import (
// 	"bytes"
// 	"context"
// 	"encoding/json"
// 	"errors"
// 	"io/fs"
// 	"log"
// 	"os"
// 	"time"

// 	"github.com/adlerhurst/eventstore/v2"
// )

// var _ eventstore.Eventstore = (*FS)(nil)

// func NewFS(path string, fs fs.FS) *FS {
// 	return &FS{
// 		fs:   fs,
// 		path: path,
// 	}
// }

// type FS struct {
// 	// root of the events
// 	fs fs.FS
// 	// path to fs
// 	path string
// }

// // Filter implements eventstore.Eventstore.
// func (store *FS) Filter(context.Context, *eventstore.Filter) ([]eventstore.Event, error) {
// 	return nil, errors.New("unimplemented")
// }

// // Push implements eventstore.Eventstore.
// func (store *FS) Push(ctx context.Context, aggregates ...eventstore.Aggregate) ([]eventstore.Event, error) {
// 	events := []eventstore.Event{}
// 	for _, aggregate := range aggregates {
// 		file, err := store.open(aggregate.ID())
// 		if err != nil {
// 			return nil, err
// 		}
// 		stat, err := file.Stat()
// 		if err != nil {
// 			return nil, err
// 		}
// 		defer func() {
// 			if err != nil {
// 				truncateErr := file.Truncate(stat.Size())
// 				if truncateErr != nil {
// 					log.Println("unable to truncate file: ", truncateErr)
// 				}
// 			}
// 			closeErr := file.Close()
// 			if err == nil {
// 				err = closeErr
// 			}
// 			if closeErr != nil {
// 				log.Println("unable to close file: ", closeErr)
// 			}
// 		}()
// 		aggregateEvents := make([]eventstore.Event, len(aggregate.Commands()))
// 		for i, command := range aggregate.Commands() {
// 			var buf bytes.Buffer
// 			if command.Payload() != nil {
// 				err := json.NewEncoder(&buf).Encode(command.Payload())
// 				// payload, err := json.Marshal(command.Payload())
// 				if err != nil {
// 					return nil, err
// 				}
// 			}
// 			aggregateEvents[i] = &event{
// 				aggregate: aggregate.ID(),
// 				Act:       command.Action(),
// 				CreatedAt: time.Now(), //TODO: maybe use file.modtime

// 				Rev: command.Revision(),
// 				// Data: payload,
// 				Data: buf.Bytes(),
// 			}
// 			if err := json.NewEncoder(file).Encode(aggregateEvents[i]); err != nil {
// 				return nil, err
// 			}
// 		}
// 		events = append(events, aggregateEvents...)
// 		if err = file.Sync(); err != nil {
// 			return nil, err
// 		}
// 	}
// 	return events, nil
// }

// // Ready implements eventstore.Eventstore.
// func (store *FS) Ready(context.Context) error {
// 	file, err := store.fs.Open(".")
// 	if err != nil {
// 		return err
// 	}
// 	return file.Close()
// }

// func (store *FS) open(aggregate eventstore.TextSubjects) (*os.File, error) {
// 	path := store.path + string(os.PathSeparator) + aggregate.Join(string(os.PathSeparator))
// 	if err := os.MkdirAll(path, 0777); err != nil {
// 		return nil, err
// 	}
// 	return os.OpenFile(
// 		path+string(os.PathSeparator)+"events.json",
// 		os.O_WRONLY|os.O_APPEND|os.O_CREATE,
// 		0777,
// 	)
// }

// var _ eventstore.Event = (*event)(nil)

// type event struct {
// 	Act       eventstore.TextSubjects `json:"action"`
// 	CreatedAt time.Time               `json:"createdAt"`
// 	Rev       uint16                  `json:"revision"`
// 	Data      []byte                  `json:"payload,omitempty"`
// 	// calculated from path
// 	aggregate eventstore.TextSubjects
// 	// line number in file
// 	sequence uint64
// }

// // Action implements eventstore.Event.
// func (e *event) Action() eventstore.TextSubjects {
// 	return e.Act
// }

// // Aggregate implements eventstore.Event.
// func (e *event) Aggregate() eventstore.TextSubjects {
// 	return e.aggregate
// }

// // CreationDate implements eventstore.Event.
// func (e *event) CreationDate() time.Time {
// 	return e.CreatedAt
// }

// // Revision implements eventstore.Event.
// func (e *event) Revision() uint16 {
// 	return e.Rev
// }

// // Sequence implements eventstore.Event.
// func (e *event) Sequence() uint64 {
// 	return e.sequence
// }

// // UnmarshalPayload implements eventstore.Event.
// func (e *event) UnmarshalPayload(object any) error {
// 	return json.Unmarshal(e.Data, object)
// }
