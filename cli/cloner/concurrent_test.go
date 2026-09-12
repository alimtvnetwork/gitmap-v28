package cloner

import (
	"testing"
)

// TestNormalizeWorkers locks in the clamping rules so that the
// --max-concurrency flag never spins up more goroutines than there is
// work to do, and never silently degrades a 0/negative request to
// something other than sequential.
func TestNormalizeWorkers(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name      string
		requested int
		jobs      int
		want      int
	}{
		{"zero clamps to one", 0, 5, 1},
		{"negative clamps to one", -7, 5, 1},
		{"one stays one", 1, 5, 1},
		{"under job count is preserved", 3, 5, 3},
		{"equal to job count is preserved", 5, 5, 5},
		{"over job count clamps down", 99, 5, 5},
		{"empty job list keeps requested", 4, 0, 4},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := normalizeWorkers(tc.requested, tc.jobs)
			if got != tc.want {
				t.Fatalf("normalizeWorkers(%d,%d) = %d, want %d",
					tc.requested, tc.jobs, got, tc.want)
			}
		})
	}
}
