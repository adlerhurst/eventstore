package api

import (
	"time"

	eventstorev1alpha "github.com/adlerhurst/eventstore/service/api/adlerhurst/eventstore/v1alpha"
	"github.com/adlerhurst/eventstore/v2"
)

func filterRequestToFilter(req *eventstorev1alpha.FilterRequest) *eventstore.Filter {
	return &eventstore.Filter{
		Queries: protoToQueries(req.Queries),
		Limit:   req.Limit,
	}
}

func protoToQueries(queries []*eventstorev1alpha.Query) []*eventstore.FilterQuery {
	filterQueries := make([]*eventstore.FilterQuery, len(queries))

	for i, query := range queries {
		filterQueries[i] = protoToQuery(query)
	}

	return filterQueries
}

func protoToQuery(query *eventstorev1alpha.Query) *eventstore.FilterQuery {
	var (
		createdAtFrom,
		createdAtTo time.Time
	)
	if query.GetCreatedAt().GetFrom() != nil {
		createdAtFrom = query.GetCreatedAt().GetFrom().AsTime()
	}
	if query.GetCreatedAt().GetTo() != nil {
		createdAtTo = query.GetCreatedAt().GetTo().AsTime()
	}
	return &eventstore.FilterQuery{
		Sequence: eventstore.SequenceFilter{
			From: query.GetSequence().GetFrom(),
			To:   query.GetSequence().GetTo(),
		},
		CreatedAt: eventstore.CreatedAtFilter{
			From: createdAtFrom,
			To:   createdAtTo,
		},
		Subjects: protoToSubjects(query.GetSubjects()),
	}
}
