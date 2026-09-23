//go:build e2e

package cmdautomation

import (
	"testing"
)

func TestE2E_GoVsPythonSearchBenchmark(t *testing.T) {
	err := RunBenchmark("search")
	if err != nil {
		t.Fatalf("search benchmark failed: %v", err)
	}
}

func TestE2E_AllAutomationBenchmarks(t *testing.T) {
	err := RunBenchmark("all")
	if err != nil {
		t.Fatalf("all benchmarks failed: %v", err)
	}
}
