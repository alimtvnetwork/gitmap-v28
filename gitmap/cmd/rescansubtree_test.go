package cmd

import (
	"reflect"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// TestSplitRescanSubtreeArgs covers the three flow shapes:
// path-only, path-then-flags, flags-then-path, plus error cases.
func TestSplitRescanSubtreeArgs(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name        string
		args        []string
		wantPath    string
		wantFlags   []string
		wantErrSubs string // substring of the error; "" means no error expected
	}{
		{
			name:      "path only",
			args:      []string{"/abs/path"},
			wantPath:  "/abs/path",
			wantFlags: nil,
		},
		{
			name:      "path then flags",
			args:      []string{"/abs/path", "--quiet", "--output", "json"},
			wantPath:  "/abs/path",
			wantFlags: []string{"--quiet", "--output", "json"},
		},
		{
			name:      "flags then path",
			args:      []string{"--quiet", "--output", "json", "/abs/path"},
			wantPath:  "/abs/path",
			wantFlags: []string{"--quiet", "--output", "json"},
		},
		{
			name:      "inline value flag does not eat next positional",
			args:      []string{"--output=json", "/abs/path"},
			wantPath:  "/abs/path",
			wantFlags: []string{"--output=json"},
		},
		{
			name:        "missing path",
			args:        []string{"--quiet"},
			wantErrSubs: "requires <absolutePath>",
		},
		{
			name:        "two paths",
			args:        []string{"/a", "/b"},
			wantErrSubs: "exactly one <absolutePath>",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			gotPath, gotFlags, err := splitRescanSubtreeArgs(tc.args)
			if tc.wantErrSubs != "" {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tc.wantErrSubs)
				}
				if !contains(err.Error(), tc.wantErrSubs) {
					t.Fatalf("error %q does not contain %q", err.Error(), tc.wantErrSubs)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if gotPath != tc.wantPath {
				t.Errorf("path: got %q want %q", gotPath, tc.wantPath)
			}
			if !reflect.DeepEqual(gotFlags, tc.wantFlags) {
				t.Errorf("flags: got %v want %v", gotFlags, tc.wantFlags)
			}
		})
	}
}

// TestBuildRescanSubtreeArgs verifies that the synthetic --max-depth is
// only injected when the caller did not supply their own, and that the
// directory always lands as the final positional regardless.
func TestBuildRescanSubtreeArgs(t *testing.T) {
	t.Parallel()
	wantDefault := []string{
		"--quiet",
		"--" + constants.FlagScanMaxDepth,
		// constants.RescanSubtreeDefaultMaxDepth → string
		intToString(constants.RescanSubtreeDefaultMaxDepth),
		"/abs/dir",
	}
	got := buildRescanSubtreeArgs("/abs/dir", []string{"--quiet"})
	if !reflect.DeepEqual(got, wantDefault) {
		t.Errorf("default injection: got %v want %v", got, wantDefault)
	}

	// User-supplied --max-depth (space form) suppresses injection.
	gotUser := buildRescanSubtreeArgs("/abs/dir",
		[]string{"--" + constants.FlagScanMaxDepth, "12"})
	wantUser := []string{"--" + constants.FlagScanMaxDepth, "12", "/abs/dir"}
	if !reflect.DeepEqual(gotUser, wantUser) {
		t.Errorf("space-form override: got %v want %v", gotUser, wantUser)
	}

	// User-supplied --max-depth=N (inline form) also suppresses.
	gotInline := buildRescanSubtreeArgs("/abs/dir",
		[]string{"--" + constants.FlagScanMaxDepth + "=-1"})
	wantInline := []string{"--" + constants.FlagScanMaxDepth + "=-1", "/abs/dir"}
	if !reflect.DeepEqual(gotInline, wantInline) {
		t.Errorf("inline-form override: got %v want %v", gotInline, wantInline)
	}
}

// TestExtractMaxDepthForLog covers all three forms the banner reads.
func TestExtractMaxDepthForLog(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		args []string
		want string
	}{
		{
			name: "space form",
			args: []string{"--" + constants.FlagScanMaxDepth, "8", "/dir"},
			want: "8",
		},
		{
			name: "inline form",
			args: []string{"--" + constants.FlagScanMaxDepth + "=-1", "/dir"},
			want: "-1",
		},
		{
			name: "single dash inline",
			args: []string{"-" + constants.FlagScanMaxDepth + "=4", "/dir"},
			want: "4",
		},
		{
			name: "missing flag",
			args: []string{"--quiet", "/dir"},
			want: "auto",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := extractMaxDepthForLog(tc.args); got != tc.want {
				t.Errorf("got %q want %q", got, tc.want)
			}
		})
	}
}

// TestRescanSubtreeDefaultIsDeeperThanScanDefault is a guardrail: the
// whole point of this command is to bump the cap, so a future edit
// that lowers RescanSubtreeDefaultMaxDepth to <= the scan default
// must trip CI.
func TestRescanSubtreeDefaultIsDeeperThanScanDefault(t *testing.T) {
	t.Parallel()
	const scanDefaultResolved = 4 // scanner.DefaultMaxDepth — duplicated to avoid the import cycle
	if constants.RescanSubtreeDefaultMaxDepth <= scanDefaultResolved {
		t.Fatalf("RescanSubtreeDefaultMaxDepth=%d must be > scanner.DefaultMaxDepth=%d so rescan-subtree is actually deeper than a plain scan",
			constants.RescanSubtreeDefaultMaxDepth, scanDefaultResolved)
	}
}

// contains is a tiny strings.Contains replacement so this test file
// stays import-light (matches the style of the production file it
// covers).
func contains(haystack, needle string) bool {
	if len(needle) == 0 {
		return true
	}
	for i := 0; i+len(needle) <= len(haystack); i++ {
		if haystack[i:i+len(needle)] == needle {
			return true
		}
	}
	return false
}

// intToString avoids importing strconv just for one call site in
// the test (the production file already uses strconv).
func intToString(n int) string {
	if n == 0 {
		return "0"
	}
	negative := n < 0
	if negative {
		n = -n
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	if negative {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}
