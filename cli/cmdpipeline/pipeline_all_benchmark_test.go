package cmdpipeline

import "testing"

// BenchmarkCollectAllPipelineSummary measures the `pipeline errors all`
// aggregation hot path against the real catalog DB. Baseline for
// program 265 (WS6); no optimization in this program.
func BenchmarkCollectAllPipelineSummary(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = collectAllPipelineSummary()
	}
}
