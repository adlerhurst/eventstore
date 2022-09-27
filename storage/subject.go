package storage

// import "github.com/adlerhurst/eventstore"

// func CheckSubject(filter []eventstore.Subject, eventType []eventstore.TextSubject) bool {
// 	for i, filterSubject := range filter {
// 		if len(eventType) <= i {
// 			return false
// 		}

// 		if s, ok := filterSubject.(eventstore.TextSubject); ok && s != eventType[i] {
// 			return false
// 		}

// 		if filterSubject == eventstore.SingleToken {
// 			continue
// 		}

// 		if filterSubject == eventstore.MultiToken {
// 			return true
// 		}
// 	}
// 	return len(filter) == len(eventType)
// }
