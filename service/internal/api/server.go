package api

import (
	"context"
	"log/slog"
	"sync"

	"github.com/adlerhurst/eventstore"
	"github.com/adlerhurst/eventstore/service/internal/api/eventstore/v1alpha"
)

var (
	_      eventstorev1alpha.EventStoreServiceServer = (*Server)(nil)
	logger                                           = slog.Default()
)

type serverOpt func(s *Server)

func WithLogger(l *slog.Logger) serverOpt {
	return func(s *Server) {
		logger = l
	}
}

func NewServer(ctx context.Context, store eventstore.Eventstore) *Server {
	return &Server{
		store: store,
	}
}

type Server struct {
	eventstorev1alpha.UnimplementedEventStoreServiceServer
	store eventstore.Eventstore
}

func (s *Server) Push(ctx context.Context, req *eventstorev1alpha.PushRequest) (*eventstorev1alpha.PushResponse, error) {
	aggregates := pushRequestToAggregates(req)

	if err := s.store.Push(ctx, aggregates...); err != nil {
		return nil, err
	}

	return &eventstorev1alpha.PushResponse{}, nil
}

func (s *Server) Filter(req *eventstorev1alpha.FilterRequest, server eventstorev1alpha.EventStoreService_FilterServer) error {
	ctx, cancel := context.WithCancel(server.Context())
	defer cancel()

	reducer := &StreamReducer{events: make(chan *eventstorev1alpha.Event, 10)}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for event := range reducer.events {
			err := server.Send(&eventstorev1alpha.FilterResponse{Events: []*eventstorev1alpha.Event{event}})
			if err != nil {
				cancel()
				break
			}
		}
		wg.Done()
	}()

	err := s.store.Filter(ctx, filterRequestToFilter(req), reducer)
	close(reducer.events)
	if err != nil {
		logger.WarnContext(ctx, "filter failed", "cause", err)
		return err
	}

	// wait until all events are sent
	wg.Wait()
	return nil
}
