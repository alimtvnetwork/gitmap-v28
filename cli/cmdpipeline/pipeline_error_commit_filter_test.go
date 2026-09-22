package cmdpipeline

import (
	"os"
	"strings"
	"testing"
	"time"
)

func TestParsePipelineErrorFlagsCommitSha(t *testing.T) {
	flags := ParsePipelineErrorFlags([]string{"ee4a694", "--json"})
	if flags.CommitTarget != "ee4a694" {
		t.Fatalf("expected CommitTarget 'ee4a694', got '%s'", flags.CommitTarget)
	}

	if !flags.IsJSON {
		t.Fatalf("expected IsJSON to be true")
	}
}

func TestParsePipelineErrorFlagsNegativeOffsets(t *testing.T) {
	cases := []struct {
		arg            string
		expectedOffset int
	}{
		{"-1", -1},
		{"-2", -2},
		{"-1n", -1},
		{"-2N", -2},
		{"HEAD~1", -1},
		{"~3", -3},
	}

	for _, tc := range cases {
		f := ParsePipelineErrorFlags([]string{tc.arg})
		offset, ok := ParseNegativeIndex(f.CommitTarget)
		if !ok || offset != tc.expectedOffset {
			t.Errorf("arg %s: expected offset=%d ok=true, got offset=%d ok=%v",
				tc.arg, tc.expectedOffset, offset, ok)
		}
	}
}

func TestFilterOutSummaryLine(t *testing.T) {
	lines := []string{
		"error: CVT1100: duplicate resource. type:MANIFEST, name:1, language:0x0409",
		"fatal error LNK1123: failure during conversion to COFF: file invalid or corrupt",
		"note: some additional compiler note",
	}

	filtered := filterOutSummaryLine(lines, "fatal error LNK1123: failure during conversion to COFF: file invalid or corrupt")
	if len(filtered) != 2 {
		t.Fatalf("expected 2 lines after filtering summary, got %d", len(filtered))
	}

	for _, l := range filtered {
		if strings.Contains(l, "LNK1123") {
			t.Fatalf("expected LNK1123 line to be filtered out, found %s", l)
		}
	}
}

func TestCompactLinkerCommandLine(t *testing.T) {
	hugeCmd := `link.exe "/NOLOGO" "/NXCOMPAT" "/LARGEADDRESSAWARE" "/OPT:REF" "/OPT:ICF" "/INCREMENTAL:NO" ` +
		`foo.lib bar.lib baz.rlib qux.rlib extra1.lib extra2.lib extra3.rlib ` +
		`"/LIBPATH:C:\rust\lib" "/LIBPATH:C:\msvc\lib" "/LIBPATH:C:\windows\kits\10\lib" ` +
		`kernel32.lib advapi32.lib shell32.lib ole32.lib oleaut32.lib uuid.lib userenv.lib ` +
		`"/OUT:target\release\deps\my_binary.exe"`

	compacted := compactLinkerCommandLine(hugeCmd)
	if len(compacted) >= len(hugeCmd) {
		t.Fatalf("expected compacted line to be shorter than raw line")
	}

	if !strings.Contains(compacted, "omitted") {
		t.Fatalf("expected compacted line to contain 'omitted', got: %s", compacted)
	}
}

func TestPipelineEtaCache(t *testing.T) {
	tmpDir := t.TempDir()
	originalWd, _ := os.Getwd()
	_ = os.Chdir(tmpDir)
	defer func() {
		_ = os.Chdir(originalWd)
	}()

	now := time.Now().UTC()
	runs := []ghRunItem{
		{
			DatabaseId: 101,
			Name:       "CI",
			Status:     "completed",
			Conclusion: "success",
			CreatedAt:  now.Add(-600 * time.Second).Format(time.RFC3339),
			UpdatedAt:  now.Add(-100 * time.Second).Format(time.RFC3339), // 500s
		},
		{
			DatabaseId: 102,
			Name:       "CI",
			Status:     "completed",
			Conclusion: "success",
			CreatedAt:  now.Add(-1300 * time.Second).Format(time.RFC3339),
			UpdatedAt:  now.Add(-700 * time.Second).Format(time.RFC3339), // 600s
		},
	}

	UpdatePipelineEtaCache("testorg/testrepo", runs)
	cachedEta := GetCachedWorkflowETA("testorg/testrepo", "CI")
	if cachedEta < 500 || cachedEta > 600 {
		t.Fatalf("expected cached ETA between 500 and 600s, got %d", cachedEta)
	}

	cacheFile := ResolveEtaCacheFilePath("testorg/testrepo")
	if _, err := os.Stat(cacheFile); os.IsNotExist(err) {
		t.Fatalf("expected eta cache file at %s, not found", cacheFile)
	}
}

func TestParseFailedLogLinesWithWarnings(t *testing.T) {
	rawLogs := `2026-09-22T23:00:00Z build-windows (pull_request) :: warning: unused variable: foo
2026-09-22T23:00:01Z build-windows (pull_request) :: warning: field is never read: bar
2026-09-22T23:00:02Z build-windows (pull_request) :: error: CVT1100: duplicate resource. type:MANIFEST
2026-09-22T23:00:03Z build-windows (pull_request) :: Process completed with exit code 1.
`
	jobs := ParseFailedLogLines(rawLogs)
	if len(jobs) == 0 {
		t.Fatalf("expected failed jobs to be parsed")
	}

	if len(jobs[0].Warnings) == 0 {
		t.Fatalf("expected warnings to be captured, got 0")
	}
}
