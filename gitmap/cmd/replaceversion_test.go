package cmd

import "testing"

// TestVersionTargetsClamping pins the spec rule that `-N` and `all`
// both clamp the lower bound at 1 and produce an ascending range that
// stops one short of the current version K.
func TestVersionTargetsClamping(t *testing.T) {
	cases := []struct {
		name string
		k    int
		n    int
		want []int
	}{
		{"minus three on v7", 7, 3, []int{4, 5, 6}},
		{"minus one on v7", 7, 1, []int{6}},
		{"all on v7 (n=0)", 7, 0, []int{1, 2, 3, 4, 5, 6}},
		{"clamp minus ten on v3", 3, 10, []int{1, 2}},
		{"v1 has nothing to bump", 1, 0, nil},
		{"v1 minus one is empty", 1, 1, nil},
		{"all on v2", 2, 0, []int{1}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := versionTargets(tc.k, tc.n)
			if !equalIntSlice(got, tc.want) {
				t.Fatalf("versionTargets(%d,%d) = %v, want %v",
					tc.k, tc.n, got, tc.want)
			}
		})
	}
}

// equalIntSlice keeps the test free of reflect.DeepEqual noise.
func equalIntSlice(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
