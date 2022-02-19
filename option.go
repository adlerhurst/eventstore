package eventstore

type Option func(*Eventstore)

// func WithPusub(ps Pubsub) Option {
// 	return func(e *Eventstore) {
// 		e.ps = ps
// 	}
// }
