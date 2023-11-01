package fs

// import (
// 	"context"
// 	_ "embed"
// 	"testing"

// 	"github.com/adlerhurst/eventstore"
// )

// // func Benchmark_Push_ParallelSameAggregate(b *testing.B) {
// // 	b.Run("Benchmark_Push_ParallelSameAggregate", func(b *testing.B) {
// // 		eventstore.PushParallelOnSameAggregate(context.Background(), b, store)
// // 	})
// // }

// func Benchmark_Push_ParallelDifferentAggregate(b *testing.B) {
// 	b.Run("Benchmark_Push_ParallelDifferentAggregate", func(b *testing.B) {
// 		eventstore.PushParallelOnDifferentAggregates(context.Background(), b, store)
// 	})
// }

// func Test_Push_Compliance(t *testing.T) {
// 	eventstore.PushComplianceTests(context.Background(), t, store)
// }
