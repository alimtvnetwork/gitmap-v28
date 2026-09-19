package cmdautomation

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

func TestCleanArtifactsClassification(t *testing.T) {
	opts := CleanArtifactsOptions{IsAll: true}
	assertCategory(t, classifyArtifact("gitmap.exe", ".exe", false, opts), "binary")
	assertCategory(t, classifyArtifact("__pycache__", "", true, opts), "pycache")
	assertCategory(t, classifyArtifact("temp.log", ".log", false, opts), "temp")
}

func assertCategory(t *testing.T, actual string, expected string) {
	if actual != expected {
		t.Errorf("expected %s category, got %s", expected, actual)
	}
}

func TestCleanArtifactsDryRun(t *testing.T) {
	tempDir := t.TempDir()
	sampleExe := filepath.Join(tempDir, "test.exe")
	_ = os.WriteFile(sampleExe, []byte("dummy binary"), 0o644)
	opts := CleanArtifactsOptions{
		Dir:             tempDir,
		IsCleanBinaries: true,
		IsDryRun:        true,
	}
	monad := RunCleanArtifacts(opts)
	if monad.IsFailure() {
		t.Fatalf("unexpected failure: %v", monad.Err)
	}
	if monad.Value.TotalFound == 0 {
		t.Errorf("expected at least 1 artifact found in temp dir")
	}
}

func TestChangedFilesParsing(t *testing.T) {
	line := "M   cli/cmdautomation/clean_artifacts.go"
	item := parseStatusLine(line, false)
	if item.Status != "modified" {
		t.Errorf("expected modified status, got %s", item.Status)
	}
	if item.Extension != ".go" {
		t.Errorf("expected .go extension, got %s", item.Extension)
	}
	addedLine := "A   new_file.ts"
	addedItem := parseStatusLine(addedLine, true)
	if addedItem.Status != "added" || addedItem.IsStaged == false {
		t.Errorf("expected staged added item, got %+v", addedItem)
	}
}

func TestDeduplicateChangedFiles(t *testing.T) {
	items := []ChangedFileItem{
		{Path: "cli/main.go", Status: "modified"},
		{Path: "cli/main.go", Status: "modified"},
		{Path: "cli/root.go", Status: "added"},
	}
	deduped := deduplicateChangedFiles(items, false, ".")
	if len(deduped) != 2 {
		t.Errorf("expected 2 deduplicated items, got %d", len(deduped))
	}
}

func TestPurgeBlobLineParsing(t *testing.T) {
	line := "e69de29bb2d1d6434b8b29ae775ad8c2e48c5391 blob 2097152 large_binary.dat"
	item := parseBlobLine(line, 1024*1024, "")
	if item.SizeBytes != 2097152 {
		t.Errorf("expected size 2097152, got %d", item.SizeBytes)
	}
	if item.Path != "large_binary.dat" {
		t.Errorf("expected path large_binary.dat, got %s", item.Path)
	}
	belowThreshold := parseBlobLine(line, 5*1024*1024, "")
	if belowThreshold.SizeBytes != 0 {
		t.Errorf("expected empty item for size below threshold")
	}
}

func TestFormatGoNormalizer(t *testing.T) {
	rawWithBomAndCrlf := []byte("\xef\xbb\xbfpackage main\r\n\r\nfunc main() {}\r\n")
	norm := normalizeFileBytes(rawWithBomAndCrlf)
	if norm.hasBom == false {
		t.Errorf("expected hasBom to be true")
	}
	if norm.hasCrlf == false {
		t.Errorf("expected hasCrlf to be true")
	}
	if string(norm.data[:12]) != "package main" {
		t.Errorf("expected BOM and CRLF stripped, got: %q", string(norm.data))
	}
}

func TestFormatGoSourceAst(t *testing.T) {
	unformatted := []byte("package main\n\nimport (\n\"os\"\n\"fmt\"\n)\nfunc main(){\nfmt.Println(\"test\")\n_ = os.Args\n}\n")
	formatted := formatGoSource(unformatted)
	if len(formatted) == 0 {
		t.Fatalf("expected non-empty formatted Go source")
	}
}

func TestCleanArtifactsMonad(t *testing.T) {
	cleanRes := CleanArtifactsResult{TotalFound: 3, TotalDeleted: 3, IsSuccess: true}
	cleanMonad := result.Ok(cleanRes)
	if cleanMonad.IsFailure() || cleanMonad.Value.IsSuccess == false {
		t.Errorf("expected successful CleanArtifactsResultMonad")
	}
}

func TestChangedFilesMonad(t *testing.T) {
	changedRes := ChangedFilesResult{TotalFiles: 2, IsSuccess: true, Duration: time.Second}
	changedMonad := result.Ok(changedRes)
	if changedMonad.IsFailure() || changedMonad.Value.TotalFiles != 2 {
		t.Errorf("expected successful ChangedFilesResultMonad")
	}
}

func TestPurgeHistoryMonad(t *testing.T) {
	purgeRes := PurgeHistoryResult{TotalFound: 1, IsSuccess: true}
	purgeMonad := result.Ok(purgeRes)
	if purgeMonad.IsFailure() || purgeMonad.Value.TotalFound != 1 {
		t.Errorf("expected successful PurgeHistoryResultMonad")
	}
}

func TestFormatGoMonad(t *testing.T) {
	formatRes := FormatGoResult{TotalFiles: 4, IsClean: true}
	formatMonad := result.Ok(formatRes)
	if formatMonad.IsFailure() || formatMonad.Value.IsClean == false {
		t.Errorf("expected successful FormatGoResultMonad")
	}
}
