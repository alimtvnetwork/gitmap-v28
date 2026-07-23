package cmd

import "testing"

func TestConvertURLToSSH(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
		ok   bool
	}{
		{"https no .git", "https://github.com/owner/repo", "git@github.com:owner/repo.git", true},
		{"https with .git", "https://github.com/owner/repo.git", "git@github.com:owner/repo.git", true},
		{"https trailing slash", "https://github.com/owner/repo/", "git@github.com:owner/repo.git", true},
		{"http upgrade", "http://gitlab.example/owner/repo", "git@gitlab.example:owner/repo.git", true},
		{"shorthand pass-through", "git@github.com:owner/repo.git", "git@github.com:owner/repo.git", true},
		{"shorthand adds .git", "git@github.com:owner/repo", "git@github.com:owner/repo.git", true},
		{"ssh scheme with port", "ssh://git@github.com:22/owner/repo.git", "git@github.com:owner/repo.git", true},
		{"ssh scheme no user", "ssh://github.com/owner/repo", "git@github.com:owner/repo.git", true},
		{"not a url", "json", "json", false},
		{"empty", "", "", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := ConvertURLToSSH(tc.in)
			if got != tc.want || ok != tc.ok {
				t.Fatalf("ConvertURLToSSH(%q) = (%q, %v), want (%q, %v)", tc.in, got, ok, tc.want, tc.ok)
			}
		})
	}
}

func TestConvertURLToHTTPS(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
		ok   bool
	}{
		{"shorthand to https", "git@github.com:owner/repo.git", "https://github.com/owner/repo.git", true},
		{"shorthand no .git", "git@github.com:owner/repo", "https://github.com/owner/repo.git", true},
		{"https pass-through", "https://github.com/owner/repo.git", "https://github.com/owner/repo.git", true},
		{"https adds .git", "https://github.com/owner/repo", "https://github.com/owner/repo.git", true},
		{"http upgraded", "http://github.com/owner/repo", "https://github.com/owner/repo.git", true},
		{"ssh scheme to https", "ssh://git@github.com:22/owner/repo.git", "https://github.com/owner/repo.git", true},
		{"not a url", "csv", "csv", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := ConvertURLToHTTPS(tc.in)
			if got != tc.want || ok != tc.ok {
				t.Fatalf("ConvertURLToHTTPS(%q) = (%q, %v), want (%q, %v)", tc.in, got, ok, tc.want, tc.ok)
			}
		})
	}
}
