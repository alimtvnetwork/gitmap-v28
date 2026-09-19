package cmdautomation

import (
	"os"
	"path/filepath"
	"testing"
)

func TestHelpAudit_MissingShortAndParity(t *testing.T) {
	tempDir := setupHelpAuditTestDir(t)
	monad := RunHelpAudit(HelpAuditOptions{Dir: tempDir, IsStrict: false})
	isFail := monad.IsFailure()
	if isFail {
		t.Fatalf("unexpected error: %v", monad.Err)
	}
	res := monad.Value
	hasVios := len(res.Violations) > 0
	if !hasVios {
		t.Errorf("expected violations for missing Short and doc, got 0")
	}
}

func setupHelpAuditTestDir(t *testing.T) string {
	tempDir := t.TempDir()
	cliDir := filepath.Join(tempDir, "cli")
	_ = os.MkdirAll(filepath.Join(cliDir, "helptext"), 0o755)
	badCmd := "package cli\nvar testCmd = &cobra.Command{Use: \"sample\"}\n"
	_ = os.WriteFile(filepath.Join(cliDir, "sample_cmd.go"), []byte(badCmd), 0o644)
	return tempDir
}

func TestPlanConsolidate_ClustersAndIndex(t *testing.T) {
	plansDir := setupPlanConsolidateTestDir(t)
	opts := PlanConsolidateOptions{Dir: plansDir, Threshold: 2, IsDryRun: false}
	monad := RunPlanConsolidate(opts)
	isFail := monad.IsFailure()
	if isFail {
		t.Fatalf("unexpected error: %v", monad.Err)
	}
	res := monad.Value
	hasClusters := res.ClustersCount > 0
	if !hasClusters {
		t.Errorf("expected at least 1 cluster, got %d", res.ClustersCount)
	}
}

func setupPlanConsolidateTestDir(t *testing.T) string {
	tempDir := t.TempDir()
	plansDir := filepath.Join(tempDir, ".ai-memory", "plans")
	completedDir := filepath.Join(plansDir, "completed")
	_ = os.MkdirAll(completedDir, 0o755)
	_ = os.WriteFile(filepath.Join(completedDir, "01-sample.md"), []byte("# Plan 1\n"), 0o644)
	_ = os.WriteFile(filepath.Join(completedDir, "02-sample.md"), []byte("# Plan 2\n"), 0o644)
	return plansDir
}

func TestDocLinks_DetectsBrokenAndFixes(t *testing.T) {
	tempDir := setupDocLinksTestDir(t)
	monad := RunDocLinks(DocLinksOptions{Dir: tempDir, IsFix: true})
	isFail := monad.IsFailure()
	if isFail {
		t.Fatalf("unexpected error: %v", monad.Err)
	}
	res := monad.Value
	verifyDocLinksOutcome(t, res)
}

func setupDocLinksTestDir(t *testing.T) string {
	tempDir := t.TempDir()
	docFile := filepath.Join(tempDir, "sample.md")
	content := "# Sample\n[Missing](nonexistent.md)\n[Old](.ai-memory/coding-guidelines/coding-guidelines.md)\n"
	_ = os.WriteFile(docFile, []byte(content), 0o644)
	return tempDir
}

func verifyDocLinksOutcome(t *testing.T, res DocLinksResult) {
	hasBroken := res.BrokenLinks > 0
	if !hasBroken {
		t.Errorf("expected broken link detection")
	}
	hasFixed := res.FixedLinks > 0
	if !hasFixed {
		t.Errorf("expected fixed link count > 0")
	}
}

func TestSpecMigrate_ResequencesAndUpdates(t *testing.T) {
	tempDir := setupSpecMigrateTestDir(t)
	opts := SpecMigrateOptions{Dir: tempDir, FromNum: 127, ToNum: 129, IsDryRun: false}
	monad := RunSpecMigrate(opts)
	isFail := monad.IsFailure()
	if isFail {
		t.Fatalf("unexpected error: %v", monad.Err)
	}
	res := monad.Value
	hasMigrated := len(res.Migrated) == 1
	if !hasMigrated {
		t.Errorf("expected 1 migrated file, got %d", len(res.Migrated))
	}
}

func setupSpecMigrateTestDir(t *testing.T) string {
	tempDir := t.TempDir()
	specDir := filepath.Join(tempDir, "02-spec", "21-app")
	_ = os.MkdirAll(specDir, 0o755)
	specFile := filepath.Join(specDir, "127-old-spec.md")
	_ = os.WriteFile(specFile, []byte("# Old Spec\n"), 0o644)
	docFile := filepath.Join(tempDir, "02-spec", "readme.md")
	_ = os.WriteFile(docFile, []byte("Refer to 127-old-spec.md for details.\n"), 0o644)
	return tempDir
}
