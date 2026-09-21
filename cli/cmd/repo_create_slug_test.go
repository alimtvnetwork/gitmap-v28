package cmd

import "testing"

func TestSlugifyRepoName(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"My Cool Repo", "my-cool-repo"},
		{"my_cool_repo", "my-cool-repo"},
		{"repo-name-123", "repo-name-123"},
		{"  Spaces  And   Tabs\t", "spaces-and-tabs"},
		{"Special!@#$%^&*Chars", "specialchars"},
		{"Leading and Trailing---", "leading-and-trailing"},
		{"Multiple...Dots_Under/Slashes", "multiple-dots-under-slashes"},
		{"already-a-slug", "already-a-slug"},
	}

	for _, tc := range tests {
		got := SlugifyRepoName(tc.input)
		if got != tc.want {
			t.Errorf("SlugifyRepoName(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}
