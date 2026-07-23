package cmd

// TestTopLevelCmdConstantsAreUnique parses every `gitmap/constants/constants_*.go`
// file and asserts that no two top-level command identifiers (Cmd* string
// constants in const blocks marked `gitmap:cmd top-level`, minus per-spec
// `gitmap:cmd skip` lines) share the same string value.
//
// Catches in CI:
//   - Go redeclaration (e.g. CmdReleaseAlias defined in two files)
//   - Same string value reused as the canonical command name of two
//     different commands (e.g. "cd" used by both CmdCD and CmdCDCmd)
//
// Mirrors the marker rules used by completion/internal/gencommands so a
// passing test guarantees the dispatch table and shell completion list
// stay internally consistent.

import (
	"path/filepath"
	"testing"
)

func TestTopLevelCmdConstantsAreUnique(t *testing.T) {
	dir := constantsDirForTest(t)

	files, err := filepath.Glob(filepath.Join(dir, "constants_*.go"))
	if err != nil {
		t.Fatalf("glob constants dir: %v", err)
	}
	if len(files) == 0 {
		t.Fatalf("no constants_*.go files found under %s", dir)
	}

	byValue := map[string][]cmdConstantOccurrence{}
	for _, path := range files {
		collectTopLevelCmdConstants(t, path, byValue)
	}

	reportDuplicateValues(t, byValue)
}
