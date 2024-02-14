package api

import (
	eventstorev1alpha "github.com/adlerhurst/eventstore/service/api/adlerhurst/eventstore/v1alpha"
	"github.com/adlerhurst/eventstore/v2"
)

func toTextSubjects(action []string) eventstore.TextSubjects {
	subjects := make(eventstore.TextSubjects, len(action))

	for i, subject := range action {
		subjects[i] = eventstore.TextSubject(subject)
	}

	return subjects
}

func protoToSubjects(subjects []*eventstorev1alpha.Subject) []eventstore.Subject {
	list := make([]eventstore.Subject, len(subjects))

	for i, subject := range subjects {
		list[i] = protoToSubject(subject)
	}

	return list
}

func protoToSubject(subject *eventstorev1alpha.Subject) eventstore.Subject {
	switch s := subject.GetSubject().(type) {
	case *eventstorev1alpha.Subject_Text:
		return eventstore.TextSubject(s.Text)
	case *eventstorev1alpha.Subject_Wildcard_:
		switch s.Wildcard {
		case eventstorev1alpha.Subject_WILDCARD_SINGLE_TOKEN:
			return eventstore.SingleToken
		case eventstorev1alpha.Subject_WILDCARD_MULTI_TOKEN:
			return eventstore.MultiToken
		}
	}
	panic("mapping of subject failed")
}
