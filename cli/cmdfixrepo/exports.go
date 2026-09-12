package cmdfixrepo

import (
	"fmt"
	"strings"
)

// FixRepoBackupManifest captures the per-snapshot index file metadata.
type FixRepoBackupManifest struct {
	SchemaVersion int      `json:"schemaVersion"`
	Repo          string   `json:"repo"`
	CurrentV      int      `json:"currentVersion"`
	Timestamp     string   `json:"timestamp"`
	GitmapVersion string   `json:"gitmapVersion"`
	Files         []string `json:"files"`
}

// FixRepoIdentity represents the resolved repo identity.
type FixRepoIdentity struct {
	Root    string
	Host    string
	Owner   string
	Base    string
	Current int
}

// ChunkerDoctorResult captures probe check result.
type ChunkerDoctorResult struct {
	Name   string `json:"name"`
	OK     bool   `json:"ok"`
	Detail string `json:"detail"`
}

// RunFixRepo executes the fix-repo CLI command.
func RunFixRepo(args []string) error {
	return runFixRepo(args)
}

// ResolveFixRepoIdentity resolves the current repository root and version.
func ResolveFixRepoIdentity() FixRepoIdentity {
	id := resolveFixRepoIdentity()
	return FixRepoIdentity{
		Root:    id.root,
		Host:    id.host,
		Owner:   id.owner,
		Base:    id.base,
		Current: id.current,
	}
}

// CopyFileForBackup safely duplicates a file before modification.
func CopyFileForBackup(src, dst string) error {
	return copyFileForBackup(src, dst)
}

// ChunkPathsForGofmt splits files into CLI-safe chunks for gofmt.
func ChunkPathsForGofmt(paths []string, maxCmdLen int) [][]string {
	return chunkPathsForGofmt(paths, maxCmdLen)
}

// BatchCmdLen calculates the total command line length of a batch.
func BatchCmdLen(batch []string) int {
	return batchCmdLen(batch)
}

// GofmtArgvOverhead returns standard argv overhead constant.
const GofmtArgvOverhead = gofmtArgvOverhead

func probeChunkerSelfTest(budget int) ChunkerDoctorResult {
	if got := chunkPathsForGofmt(nil, budget); got != nil {
		return ChunkerDoctorResult{Name: "chunker-selftest", OK: false, Detail: "empty input did not return nil"}
	}

	small := []string{"a.go", "b.go", "c.go"}
	if got := chunkPathsForGofmt(small, budget); len(got) != 1 {
		return ChunkerDoctorResult{
			Name: "chunker-selftest", OK: false,
			Detail: fmt.Sprintf("expected 1 chunk for %d small paths, got %d", len(small), len(got)),
		}
	}

	long := strings.Repeat("x", 200)
	overflow := make([]string, 500)
	for i := range overflow {
		overflow[i] = long
	}

	batches := chunkPathsForGofmt(overflow, budget)
	if len(batches) < 2 {
		return ChunkerDoctorResult{
			Name: "chunker-selftest", OK: false,
			Detail: fmt.Sprintf("overflow input yielded only %d batch(es)", len(batches)),
		}
	}

	for i, b := range batches {
		if len(b) > 1 && batchCmdLen(b)-gofmtArgvOverhead > budget {
			return ChunkerDoctorResult{
				Name: "chunker-selftest", OK: false,
				Detail: fmt.Sprintf("batch %d exceeds budget %d", i+1, budget),
			}
		}
	}

	return ChunkerDoctorResult{
		Name: "chunker-selftest", OK: true,
		Detail: fmt.Sprintf("500 synthetic paths -> %d batches under budget %d", len(batches), budget),
	}
}

// ProbeChunkerSelfTest runs the chunker self-test.
func ProbeChunkerSelfTest(budget int) ChunkerDoctorResult {
	return probeChunkerSelfTest(budget)
}

// RewriteFixRepoFile applies target rewrites to a file.
func RewriteFixRepoFile(fullPath, base string, current int, targets []int, dryRun bool) (int, error) {
	return rewriteFixRepoFile(fullPath, base, current, targets, dryRun)
}
