package api

import (
	"strconv"

	eventstorev1alpha "github.com/adlerhurst/eventstore/service/api/adlerhurst/eventstore/v1alpha"
	"github.com/adlerhurst/eventstore/v2"
	"google.golang.org/protobuf/types/known/structpb"
)

var _ eventstore.Reducer = (*StreamReducer)(nil)

type StreamReducer struct {
	events chan *eventstorev1alpha.Event
}

// Reduce implements eventstore.Reducer.
func (r *StreamReducer) Reduce(events ...eventstore.Event) error {
	for _, event := range events {
		e, err := eventToProto(event)
		if err != nil {
			return err
		}
		r.events <- e
	}
	return nil
}

func eventToProto(event eventstore.Event) (*eventstorev1alpha.Event, error) {
	e := &eventstorev1alpha.Event{
		Id: event.Aggregate().Join(".") + "." + strconv.Itoa(int(event.Sequence())),
		Action: &eventstorev1alpha.Action{
			Action:   actionToProto(event.Action()),
			Revision: uint32(event.Revision()),
			Payload:  new(structpb.Struct),
		},
	}

	if err := event.UnmarshalPayload(e.Action.Payload); err != nil {
		return nil, err
	}

	return e, nil
}

func actionToProto(subjects eventstore.TextSubjects) []string {
	action := make([]string, len(subjects))
	for i, subject := range subjects {
		action[i] = string(subject)
	}
	return action
}
