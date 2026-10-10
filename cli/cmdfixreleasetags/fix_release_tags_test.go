package cmdfixreleasetags

import (
	"bytes"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

func TestParseFixReleaseTagsFlags_Defaults(t *testing.T) {
	flags, err := ParseFixReleaseTagsFlags([]string{})
	if err != nil {
		t.Fatalf("unexpected error parsing empty flags: %v", err)
	}

	if flags.IsDryRun {
		t.Errorf("expected IsDryRun=false, got true")
	}
	if flags.IsConfirmed {
		t.Errorf("expected IsConfirmed=false, got true")
	}
	if flags.IsJSON {
		t.Errorf("expected IsJSON=false, got true")
	}
	if flags.IsLocalOnly {
		t.Errorf("expected IsLocalOnly=false, got true")
	}
	if flags.IsRemoteOnly {
		t.Errorf("expected IsRemoteOnly=false, got true")
	}
	if flags.IsVerbose {
		t.Errorf("expected IsVerbose=false, got true")
	}
	if flags.IsHelpRequested {
		t.Errorf("expected IsHelpRequested=false, got true")
	}
	if flags.TargetDirectory != "." {
		t.Errorf("expected TargetDirectory='.', got '%s'", flags.TargetDirectory)
	}
}

func TestParseFixReleaseTagsFlags_DryRunPermutations(t *testing.T) {
	cases := [][]string{
		{"--dry-run"},
		{"-n"},
		{"--dryrun"},
		{"-dry-run"},
	}

	for _, args := range cases {
		flags, err := ParseFixReleaseTagsFlags(args)
		if err != nil {
			t.Fatalf("unexpected error for %v: %v", args, err)
		}
		if !flags.IsDryRun {
			t.Errorf("expected IsDryRun=true for args %v", args)
		}
	}
}

func TestParseFixReleaseTagsFlags_ConfirmPermutations(t *testing.T) {
	cases := [][]string{
		{"-y"},
		{"--yes"},
		{"--confirm"},
		{"-confirm"},
	}

	for _, args := range cases {
		flags, err := ParseFixReleaseTagsFlags(args)
		if err != nil {
			t.Fatalf("unexpected error for %v: %v", args, err)
		}
		if !flags.IsConfirmed {
			t.Errorf("expected IsConfirmed=true for args %v", args)
		}
	}
}

func TestParseFixReleaseTagsFlags_JSONAndVerbose(t *testing.T) {
	flags, err := ParseFixReleaseTagsFlags([]string{"--json", "-v"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !flags.IsJSON {
		t.Errorf("expected IsJSON=true")
	}
	if !flags.IsVerbose {
		t.Errorf("expected IsVerbose=true")
	}
}

func TestParseFixReleaseTagsFlags_MutualExclusion(t *testing.T) {
	_, err := ParseFixReleaseTagsFlags([]string{"--local-only", "--remote-only"})
	if err == nil {
		t.Fatalf("expected error when both --local-only and --remote-only are passed")
	}
	if !strings.Contains(err.Error(), "cannot specify both --local-only and --remote-only") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestParseFixReleaseTagsFlags_TargetDirectory(t *testing.T) {
	flags, err := ParseFixReleaseTagsFlags([]string{"--repo", "."})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if flags.TargetDirectory != "." {
		t.Errorf("expected TargetDirectory='.', got '%s'", flags.TargetDirectory)
	}

	flagsEq, err := ParseFixReleaseTagsFlags([]string{"--repo=."})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if flagsEq.TargetDirectory != "." {
		t.Errorf("expected TargetDirectory='.', got '%s'", flagsEq.TargetDirectory)
	}
}

func TestParseFixReleaseTagsFlags_InvalidDirectory(t *testing.T) {
	_, err := ParseFixReleaseTagsFlags([]string{"--repo", "non_existent_directory_xyz123"})
	if err == nil {
		t.Fatalf("expected error for non-existent directory")
	}
	if !strings.Contains(err.Error(), "does not exist or is inaccessible") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestParseFixReleaseTagsFlags_HelpFlag(t *testing.T) {
	for _, arg := range []string{"-h", "--help", "help"} {
		flags, err := ParseFixReleaseTagsFlags([]string{arg})
		if err != nil {
			t.Fatalf("unexpected error for %s: %v", arg, err)
		}
		if !flags.IsHelpRequested {
			t.Errorf("expected IsHelpRequested=true for %s", arg)
		}
	}
}

func TestBuildDiagnosticTableConfig_ColumnsAndRows(t *testing.T) {
	records := []ReleaseTagAuditRecord{
		{
			Tag:                   "v1.0.0",
			CommitSha:             "abcdef123456",
			HasGitHubRelease:      true,
			IsDraft:               false,
			AssetCount:            3,
			AuditReason:           ReasonHealthy,
			ProtectionStatus:      StatusEligible,
			IsEligibleForDeletion: false,
			WorkflowRuns: []CIWorkflowRunInfo{
				{Conclusion: "success", IsSuccessful: true},
			},
		},
		{
			Tag:                   "v1.0.1",
			CommitSha:             "789012abcdef",
			HasGitHubRelease:      false,
			AuditReason:           ReasonOrphanTag,
			ProtectionStatus:      StatusEligible,
			IsEligibleForDeletion: true,
		},
		{
			Tag:                   "v1.0.2",
			CommitSha:             "345678abcdef",
			HasGitHubRelease:      true,
			IsDraft:               true,
			AssetCount:            0,
			AuditReason:           ReasonDraftRelease,
			ProtectionStatus:      StatusProtectedGrace,
			IsProtected:           true,
			IsEligibleForDeletion: false,
			WorkflowRuns: []CIWorkflowRunInfo{
				{Status: "in_progress", IsInProgress: true},
			},
		},
	}

	cfg := BuildDiagnosticTableConfig(records)

	if len(cfg.Columns) != 6 {
		t.Fatalf("expected 6 columns, got %d", len(cfg.Columns))
	}

	expectedCols := []string{"TAG", "COMMIT", "RELEASE STATUS", "ASSETS", "CI/CD STATUS", "ACTION"}
	for i, name := range expectedCols {
		if cfg.Columns[i].Title != name {
			t.Errorf("column %d title expected '%s', got '%s'", i, name, cfg.Columns[i].Title)
		}
		if cfg.Columns[i].Align != termout.AlignLeft {
			t.Errorf("column %d expected AlignLeft", i)
		}
	}

	if len(cfg.Rows) != 3 {
		t.Fatalf("expected 3 rows, got %d", len(cfg.Rows))
	}

	// Verify row 0: healthy
	r0 := cfg.Rows[0].Cells
	if r0[0] != "v1.0.0" {
		t.Errorf("row 0 tag expected 'v1.0.0', got '%s'", r0[0])
	}
	if r0[1] != "abcdef1" {
		t.Errorf("row 0 commit expected 'abcdef1', got '%s'", r0[1])
	}
	if r0[2] != "Published" {
		t.Errorf("row 0 release status expected 'Published', got '%s'", r0[2])
	}
	if r0[3] != "3 Assets" {
		t.Errorf("row 0 assets expected '3 Assets', got '%s'", r0[3])
	}
	if r0[4] != "Passed" {
		t.Errorf("row 0 CI/CD expected 'Passed', got '%s'", r0[4])
	}
	if r0[5] != "Keep (Healthy Release)" {
		t.Errorf("row 0 action expected 'Keep (Healthy Release)', got '%s'", r0[5])
	}

	// Verify row 1: orphan
	r1 := cfg.Rows[1].Cells
	if r1[2] != "Missing Release" {
		t.Errorf("row 1 release status expected 'Missing Release', got '%s'", r1[2])
	}
	if r1[3] != "0 Assets" {
		t.Errorf("row 1 assets expected '0 Assets', got '%s'", r1[3])
	}
	if r1[4] != "No CI" {
		t.Errorf("row 1 CI/CD expected 'No CI', got '%s'", r1[4])
	}
	if r1[5] != "Delete Release & Tag" {
		t.Errorf("row 1 action expected 'Delete Release & Tag', got '%s'", r1[5])
	}

	// Verify row 2: draft protected by grace period
	r2 := cfg.Rows[2].Cells
	if r2[2] != "Draft" {
		t.Errorf("row 2 release status expected 'Draft', got '%s'", r2[2])
	}
	if r2[4] != "In-Progress" {
		t.Errorf("row 2 CI/CD expected 'In-Progress', got '%s'", r2[4])
	}
	if r2[5] != "Skip (Grace Period)" {
		t.Errorf("row 2 action expected 'Skip (Grace Period)', got '%s'", r2[5])
	}
}

func TestConfirmDeletion_ConfirmedBypass(t *testing.T) {
	confirmed, err := ConfirmDeletion(5, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !confirmed {
		t.Errorf("expected confirmed=true when isConfirmed=true")
	}
}

func TestConfirmDeletion_NonInteractiveFailsafe(t *testing.T) {
	// Ensure override is cleared
	StdinReaderOverride = nil
	forceFalse := false
	ForceInteractiveOverride = &forceFalse
	defer func() {
		ForceInteractiveOverride = nil
	}()

	_, err := ConfirmDeletion(3, false)
	if err == nil {
		t.Fatalf("expected error in non-interactive environment")
	}
	if !strings.Contains(err.Error(), "E1025") {
		t.Errorf("expected error to contain E1025, got: %v", err)
	}
	if !strings.Contains(err.Error(), "non-interactive") {
		t.Errorf("expected error to mention non-interactive, got: %v", err)
	}
}

func TestConfirmDeletion_InteractiveInputYes(t *testing.T) {
	forceTrue := true
	ForceInteractiveOverride = &forceTrue
	defer func() {
		ForceInteractiveOverride = nil
		StdinReaderOverride = nil
	}()

	for _, answer := range []string{"y\n", "yes\n", "Y\n", "YES\n"} {
		StdinReaderOverride = bytes.NewBufferString(answer)
		confirmed, err := ConfirmDeletion(2, false)
		if err != nil {
			t.Fatalf("unexpected error for %q: %v", answer, err)
		}
		if !confirmed {
			t.Errorf("expected confirmed=true for %q", answer)
		}
	}
}

func TestConfirmDeletion_InteractiveInputNo(t *testing.T) {
	forceTrue := true
	ForceInteractiveOverride = &forceTrue
	defer func() {
		ForceInteractiveOverride = nil
		StdinReaderOverride = nil
	}()

	for _, answer := range []string{"n\n", "no\n", "\n", "cancel\n"} {
		StdinReaderOverride = bytes.NewBufferString(answer)
		confirmed, err := ConfirmDeletion(2, false)
		if err != nil {
			t.Fatalf("unexpected error for %q: %v", answer, err)
		}
		if confirmed {
			t.Errorf("expected confirmed=false for %q", answer)
		}
	}
}

func TestBuildFixReleaseTagsHelpMenu(t *testing.T) {
	menu := BuildFixReleaseTagsHelpMenu()

	if menu.Title == "" {
		t.Errorf("expected non-empty help menu title")
	}

	if len(menu.UsageLines) == 0 {
		t.Errorf("expected usage lines")
	}

	foundCmd := false
	for _, ul := range menu.UsageLines {
		if strings.Contains(ul, "gitmap fix release tags") {
			foundCmd = true
			break
		}
	}
	if !foundCmd {
		t.Errorf("expected usage lines to include 'gitmap fix release tags'")
	}

	if len(menu.Sections) != 3 {
		t.Fatalf("expected 3 help sections, got %d", len(menu.Sections))
	}

	if len(menu.FooterFlags) == 0 {
		t.Errorf("expected footer flags")
	}

	if len(menu.Tips) == 0 {
		t.Errorf("expected tips")
	}
}

func TestRenderFixReleaseTagsHelp_NoPanic(t *testing.T) {
	// Simply ensure rendering does not crash or panic
	RenderFixReleaseTagsHelp()
}

func TestRunFixReleaseTags_HelpFlag(t *testing.T) {
	err := RunFixReleaseTags([]string{"--help"})
	if err != nil {
		t.Fatalf("expected nil error on --help, got: %v", err)
	}

	err = RunCLI([]string{"-h"})
	if err != nil {
		t.Fatalf("expected nil error on -h, got: %v", err)
	}
}

func TestRunFixReleaseTags_DryRunExecution(t *testing.T) {
	os.Setenv("GITMAP_MOCK_GH", "1")
	defer os.Unsetenv("GITMAP_MOCK_GH")

	// Running dry-run on current repository should succeed without error
	err := RunFixReleaseTags([]string{"--dry-run", "--repo", "."})
	if err != nil {
		t.Fatalf("expected nil error on dry-run, got: %v", err)
	}
}

func TestActionCellFormatting(t *testing.T) {
	now := time.Now()
	_ = now

	cases := []struct {
		name     string
		rec      ReleaseTagAuditRecord
		expected string
	}{
		{
			name: "eligible",
			rec: ReleaseTagAuditRecord{
				IsEligibleForDeletion: true,
			},
			expected: "Delete Release & Tag",
		},
		{
			name: "protected active",
			rec: ReleaseTagAuditRecord{
				IsProtected:      true,
				ProtectionStatus: StatusProtectedActive,
			},
			expected: "Skip (Active Version)",
		},
		{
			name: "protected latest",
			rec: ReleaseTagAuditRecord{
				IsProtected:      true,
				ProtectionStatus: StatusProtectedLatest,
			},
			expected: "Skip (Latest Healthy)",
		},
		{
			name: "protected grace",
			rec: ReleaseTagAuditRecord{
				IsProtected:      true,
				ProtectionStatus: StatusProtectedGrace,
			},
			expected: "Skip (Grace Period)",
		},
		{
			name: "healthy",
			rec: ReleaseTagAuditRecord{
				AuditReason: ReasonHealthy,
			},
			expected: "Keep (Healthy Release)",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual := formatActionCell(tc.rec)
			if actual != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, actual)
			}
		})
	}
}
