package cmd

import (
	"strings"
	"testing"
)

// TestParsePRCEntries covers the comma + whitespace splitter for `gitmap prc`.
func TestParsePRCEntries(t *testing.T) {
	t.Parallel()

	type tc struct {
		name    string
		args    []string
		want    []prcEntry
		wantErr string
	}

	cases := []tc{
		{
			name: "single entry",
			args: []string{"repo-v1", "v2"},
			want: []prcEntry{{token: "repo-v1", version: "v2"}},
		},
		{
			name: "two entries comma-separated",
			args: []string{"alpha", "v1,", "beta", "v2"},
			want: []prcEntry{
				{token: "alpha", version: "v1"},
				{token: "beta", version: "v2"},
			},
		},
		{
			name:    "missing version",
			args:    []string{"lone"},
			wantErr: "missing version",
		},
		{
			name:    "extra tokens",
			args:    []string{"a", "v1", "junk"},
			wantErr: "extra tokens",
		},
		{
			name:    "no args",
			args:    nil,
			wantErr: "usage:",
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			got, err := parsePRCEntries(c.args)
			if c.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), c.wantErr) {
					t.Fatalf("want err containing %q, got %v", c.wantErr, err)
				}

				return
			}
			if err != nil {
				t.Fatalf("unexpected err: %v", err)
			}
			if len(got) != len(c.want) {
				t.Fatalf("len=%d want %d (%+v)", len(got), len(c.want), got)
			}
			for i := range got {
				if got[i] != c.want[i] {
					t.Errorf("entry[%d]=%+v want %+v", i, got[i], c.want[i])
				}
			}
		})
	}
}

// TestIsPRCURL checks URL detection covers https, ssh, and slug rejection.
func TestIsPRCURL(t *testing.T) {
	t.Parallel()

	cases := map[string]bool{
		"https://github.com/a/b.git": true,
		"http://x/y":                 true,
		"git@github.com:a/b.git":     true,
		"alice/repo":                 false,
		"repo-v1":                    false,
		"":                           false,
	}

	for in, want := range cases {
		if got := isPRCURL(in); got != want {
			t.Errorf("isPRCURL(%q)=%v want %v", in, got, want)
		}
	}
}

// TestSlugFromURL verifies last-segment + .git-strip + lowercasing.
func TestSlugFromURL(t *testing.T) {
	t.Parallel()

	cases := map[string]string{
		"https://github.com/Alice/MyRepo.git": "myrepo",
		"https://gitlab.com/g/sub/proj":       "proj",
		"git@github.com:Org/Thing.git":        "thing",
		"PlainName":                           "plainname",
		"":                                    "",
	}

	for in, want := range cases {
		if got := slugFromURL(in); got != want {
			t.Errorf("slugFromURL(%q)=%q want %q", in, got, want)
		}
	}
}
