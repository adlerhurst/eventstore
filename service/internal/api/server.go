package api

import (
	"context"

	"github.com/adlerhurst/eventstore"
	"github.com/adlerhurst/eventstore/service/internal/api/eventstore/v1alpha"
)

var _ eventstorev1alpha.EventStoreServiceServer = (*Server)(nil)

type Server struct {
	eventstorev1alpha.UnimplementedEventStoreServiceServer
	store eventstore.Eventstore
}

func (s *Server) Push(context.Context, *eventstorev1alpha.PushRequest) (*eventstorev1alpha.PushResponse, error) {
	return nil, nil
}

func (s *Server) Filter(*eventstorev1alpha.FilterRequest, eventstorev1alpha.EventStoreService_FilterServer) error {
	return nil
}
