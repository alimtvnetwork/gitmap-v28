package cmdhealth

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	cmdconstants "github.com/alimtvnetwork/gitmap-v28/cli/cmdhealth/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdhealth/enums"
)

func TestRunHealthHelp(t *testing.T) {
	if err := RunHealth([]string{"--help"}); err != nil {
		t.Errorf("expected nil error on --help, got %v", err)
	}
	if err := RunHealth([]string{"-h"}); err != nil {
		t.Errorf("expected nil error on -h, got %v", err)
	}
}

func TestRunHealthFlagError(t *testing.T) {
	err := RunHealth([]string{"--invalid-flag-that-does-not-exist"})
	if err == nil {
		t.Error("expected non-nil error on unknown flag")
	}
}

func TestResolveHealthWorkdirFlag(t *testing.T) {
	expected := filepath.Join("custom", "dir")
	resolved := resolveHealthWorkdir(expected)
	if resolved != expected {
		t.Errorf("expected %s, got %s", expected, resolved)
	}
}

func TestResolveHealthWorkdirEnv(t *testing.T) {
	targetEnv := filepath.Join("env", "path")
	_ = os.Setenv(cmdconstants.EnvGitmapWorkdir, targetEnv)
	defer os.Unsetenv(cmdconstants.EnvGitmapWorkdir)

	resolved := resolveHealthWorkdir("")
	if resolved != targetEnv {
		t.Errorf("expected %s from env, got %s", targetEnv, resolved)
	}
}

func TestWorstItemPointerPipeline(t *testing.T) {
	facts := healthFacts{
		pipelineRed: 2,
		redNames:    []string{"owner/repo1", "owner/repo2"},
	}
	worst := worstItemPointer(facts)
	if !strings.Contains(worst, "pe triage -> owner/repo1") {
		t.Errorf("expected pipeline triage recommendation, got %s", worst)
	}
}

func TestWorstItemPointerStaleDb(t *testing.T) {
	facts := healthFacts{
		staleDBs:   1,
		staleNames: []string{"master.db"},
	}
	worst := worstItemPointer(facts)
	if worst != cmdconstants.WorstDbDoctorRepair {
		t.Errorf("expected db repair recommendation, got %s", worst)
	}
}

func TestWorstItemPointerDirtyRepos(t *testing.T) {
	facts := healthFacts{
		repoDirty:  3,
		dirtyNames: []string{"repo1", "repo2", "repo3"},
	}
	worst := worstItemPointer(facts)
	if !strings.Contains(worst, "dirty working trees (3)") {
		t.Errorf("expected dirty working trees recommendation, got %s", worst)
	}
}

func TestWorstItemPointerClean(t *testing.T) {
	facts := healthFacts{
		repoTotal: 10,
	}
	worst := worstItemPointer(facts)
	if worst != cmdconstants.WorstNothingUrgent {
		t.Errorf("expected all checks green recommendation, got %s", worst)
	}
}

func TestDetermineWorstSeverityAndCategory(t *testing.T) {
	testSeverityAndCategoryPipeline(t)
	testSeverityAndCategoryDb(t)
	testSeverityAndCategoryGit(t)
	testSeverityAndCategoryClean(t)
}

func testSeverityAndCategoryPipeline(t *testing.T) {
	facts := healthFacts{pipelineRed: 1, redNames: []string{"r1"}}
	if DetermineWorstSeverity(facts) != enums.HealthSeverityError {
		t.Error("expected error severity for red pipeline")
	}
	if DetermineWorstCategory(facts) != enums.HealthCategoryPipeline {
		t.Error("expected pipeline category for red pipeline")
	}
}

func testSeverityAndCategoryDb(t *testing.T) {
	facts := healthFacts{staleDBs: 2}
	if DetermineWorstSeverity(facts) != enums.HealthSeverityWarn {
		t.Error("expected warn severity for stale dbs")
	}
	if DetermineWorstCategory(facts) != enums.HealthCategoryDb {
		t.Error("expected db category for stale dbs")
	}
}

func testSeverityAndCategoryGit(t *testing.T) {
	facts := healthFacts{repoDirty: 4}
	if DetermineWorstSeverity(facts) != enums.HealthSeverityWarn {
		t.Error("expected warn severity for dirty repos")
	}
	if DetermineWorstCategory(facts) != enums.HealthCategoryGit {
		t.Error("expected git category for dirty repos")
	}
}

func testSeverityAndCategoryClean(t *testing.T) {
	facts := healthFacts{repoTotal: 5}
	if DetermineWorstSeverity(facts) != enums.HealthSeverityOk {
		t.Error("expected ok severity for clean state")
	}
	if DetermineWorstCategory(facts) != enums.HealthCategorySystem {
		t.Error("expected system category for clean state")
	}
}

func TestParseHealthSlug(t *testing.T) {
	testValidSlugs(t)
	testInvalidSlugs(t)
}

func testValidSlugs(t *testing.T) {
	slug, hasSlug := parseHealthSlug("https://github.com/alimtvnetwork/gitmap-v28.git")
	if !hasSlug || slug != "alimtvnetwork/gitmap-v28" {
		t.Errorf("unexpected slug result: %s, %v", slug, hasSlug)
	}
	slugSsh, hasSlugSsh := parseHealthSlug("git@github.com:alimtvnetwork/gitmap-v28.git")
	if !hasSlugSsh || slugSsh != "alimtvnetwork/gitmap-v28" {
		t.Errorf("unexpected ssh slug result: %s, %v", slugSsh, hasSlugSsh)
	}
}

func testInvalidSlugs(t *testing.T) {
	if _, hasSlug := parseHealthSlug("https://gitlab.com/other/project.git"); hasSlug {
		t.Error("expected non-github url to not resolve slug")
	}
	if _, hasSlug := parseHealthSlug("not-a-url"); hasSlug {
		t.Error("expected invalid url to not resolve slug")
	}
}
