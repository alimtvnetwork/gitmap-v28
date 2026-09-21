package cmd

import "testing"

func TestNormalizeCreateArgs_RepoCreateSuite(t *testing.T) {
	tests := []struct {
		input []string
		want  []string
	}{
		{[]string{"repo", "My Project"}, []string{"My Project"}},
		{[]string{"repository", "My Project"}, []string{"My Project"}},
		{[]string{"My Project"}, []string{"My Project"}},
		{[]string{"My Project", "--local"}, []string{"My Project", "--local"}},
	}

	for _, tc := range tests {
		got := normalizeCreateArgs(tc.input)
		if len(got) != len(tc.want) {
			t.Errorf("normalizeCreateArgs(%v) len = %d, want %d", tc.input, len(got), len(tc.want))
			continue
		}
		for i := range got {
			if got[i] != tc.want[i] {
				t.Errorf("normalizeCreateArgs(%v)[%d] = %q, want %q", tc.input, i, got[i], tc.want[i])
			}
		}
	}
}
