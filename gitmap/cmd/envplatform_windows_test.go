//go:build windows

package cmd

import "testing"

func TestFilterPathPartsRemovesMatch(t *testing.T) {
	parts := []string{`C:\bin`, `C:\tools`, ` C:\bin `, `D:\go\bin`}
	got := filterPathParts(parts, `C:\bin`)
	want := []string{`C:\tools`, `D:\go\bin`}
	if len(got) != len(want) {
		t.Fatalf("filterPathParts len = %d, want %d (%v)", len(got), len(want), got)
	}
	for i := range got {
		if got[i] != want[i] {
			t.Errorf("filterPathParts[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestFilterPathPartsCaseInsensitive(t *testing.T) {
	parts := []string{`C:\Bin`, `C:\tools`}
	got := filterPathParts(parts, `c:\bin`)
	if len(got) != 1 || got[0] != `C:\tools` {
		t.Errorf("expected case-insensitive removal; got %v", got)
	}
}
