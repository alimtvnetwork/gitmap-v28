package gitutil

import "testing"

// TestParseLsRemoteSymref covers the pure-string parsing of
// `git ls-remote --symref origin HEAD` output without invoking git.
// The fixtures mirror real outputs we've seen across GitHub, GitLab,
// and Gitea — including the leading-blank-line quirk some servers
// emit and the rare multi-symref case where the first match wins.
func TestParseLsRemoteSymref(t *testing.T) {
	cases := []struct {
		name   string
		input  string
		want   string
		wantOK bool
	}{
		{
			name:   "github main",
			input:  "ref: refs/heads/main\tHEAD\nabc123\tHEAD\n",
			want:   "main",
			wantOK: true,
		},
		{
			name:   "gitlab master with leading blank",
			input:  "\nref: refs/heads/master\tHEAD\ndeadbeef\tHEAD\n",
			want:   "master",
			wantOK: true,
		},
		{
			name:   "develop branch",
			input:  "ref: refs/heads/develop\tHEAD\n",
			want:   "develop",
			wantOK: true,
		},
		{
			name:   "no symref present",
			input:  "abc123\tHEAD\n",
			want:   "",
			wantOK: false,
		},
		{
			name:   "empty output",
			input:  "",
			want:   "",
			wantOK: false,
		},
		{
			name:   "malformed ref line",
			input:  "ref: not-a-real-ref\tHEAD\n",
			want:   "",
			wantOK: false,
		},
		{
			name:   "first symref wins",
			input:  "ref: refs/heads/trunk\tHEAD\nref: refs/heads/main\tHEAD\n",
			want:   "trunk",
			wantOK: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := parseLsRemoteSymref(tc.input)
			if ok != tc.wantOK || got != tc.want {
				t.Fatalf("parseLsRemoteSymref(%q) = (%q,%v), want (%q,%v)",
					tc.input, got, ok, tc.want, tc.wantOK)
			}
		})
	}
}
