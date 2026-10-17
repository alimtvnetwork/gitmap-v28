package llm

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

func TestIsTrainCommand(t *testing.T) {
	if !isTrainCommand("train") {
		t.Errorf("expected isTrainCommand(train) to be true")
	}
	if !isTrainCommand("chain") {
		t.Errorf("expected isTrainCommand(chain) to be true")
	}
	if isTrainCommand("status") {
		t.Errorf("expected isTrainCommand(status) to be false")
	}
}

func TestRunTrainTextOnly(t *testing.T) {
	err := RunTrain([]string{"--text-only"})
	if err != nil {
		t.Fatalf("expected nil error for RunTrain --text-only, got: %v", err)
	}
}

func TestGenerateSkillFile(t *testing.T) {
	tempDir := t.TempDir()
	targetSkill := filepath.Join(tempDir, "skills", "gitmap", "SKILL.md")

	err := GenerateSkillFile(targetSkill)
	if err != nil {
		t.Fatalf("expected nil error for GenerateSkillFile, got: %v", err)
	}

	content, readErr := os.ReadFile(targetSkill)
	if readErr != nil {
		t.Fatalf("failed to read generated skill file: %v", readErr)
	}

	str := string(content)
	if !strings.Contains(str, "name: gitmap") {
		t.Errorf("expected YAML frontmatter 'name: gitmap'")
	}
	if !strings.Contains(str, AuthorName) {
		t.Errorf("expected author attribution in skill file")
	}
	if !strings.Contains(str, SponsorName) {
		t.Errorf("expected sponsor attribution in skill file")
	}
	if !strings.Contains(str, "gitmap aum search") {
		t.Errorf("expected 'gitmap aum search' in skill file")
	}
}

func TestRunTrainLoop(t *testing.T) {
	err := RunTrain([]string{"--loop"})
	if err != nil {
		t.Fatalf("expected nil error for RunTrain --loop, got: %v", err)
	}
}

func TestRunTrainSelfLoopCount(t *testing.T) {
	err := RunTrain([]string{"--self-loop", "2"})
	if err != nil {
		t.Fatalf("expected nil error for RunTrain --self-loop 2, got: %v", err)
	}
}

func TestRunTrainURL(t *testing.T) {
	err := RunTrain([]string{"--url"})
	if err != nil {
		t.Fatalf("expected nil error for RunTrain --url, got: %v", err)
	}
}

func TestRunTrainJSON(t *testing.T) {
	err := RunTrain([]string{"--json"})
	if err != nil {
		t.Fatalf("expected nil error for RunTrain --json, got: %v", err)
	}
}

func TestRunTrainHelp(t *testing.T) {
	err := RunTrain([]string{"--help"})
	if err != nil {
		t.Fatalf("expected nil error for RunTrain --help, got: %v", err)
	}
}

func TestRunLlmHelp(t *testing.T) {
	err := Run([]string{"--help"})
	if err != nil {
		t.Fatalf("expected nil error for Run --help, got: %v", err)
	}
}

func TestRunTrainURLs(t *testing.T) {
	err := RunTrain([]string{"--urls"})
	if err != nil {
		t.Fatalf("expected nil error for RunTrain --urls, got: %v", err)
	}
}

func TestPublicDocLinksContent(t *testing.T) {
	links := GetPublicDocLinks()
	if len(links) < 5 {
		t.Errorf("expected at least 5 public doc links, got %d", len(links))
	}
	for _, link := range links {
		if !strings.HasPrefix(link.URL, "https://") {
			t.Errorf("expected link URL to start with https://, got %s", link.URL)
		}
	}
}

func TestRecursiveInstructions(t *testing.T) {
	instructions := RenderRecursiveInstructions()
	if !strings.Contains(instructions, "STAGE 4: RECURSIVE GIT NETWORK LEARNING DIRECTIVES") {
		t.Errorf("expected STAGE 4 in recursive instructions")
	}
	if !strings.Contains(instructions, "gitmap pe history-ai") {
		t.Errorf("expected gitmap pe history-ai mention in recursive instructions")
	}

	directive := RenderSkillCreationDirective()
	if !strings.Contains(directive, "AUTHORITATIVE SKILL INGESTION") {
		t.Errorf("expected skill ingestion directive header")
	}
}

func TestTrainPhaseRegistryOrder(t *testing.T) {
	want := []string{
		"1. Discovery",
		"2. Refactoring",
		"3. Verification",
		"4. Semantic Commit",
		"5. Telemetry",
		"6. Heal & Fix",
	}
	if len(trainPhases) != len(want) {
		t.Fatalf("expected %d train phases, got %d", len(want), len(trainPhases))
	}
	for i, phase := range trainPhases {
		if phase.Name() != want[i] {
			t.Errorf("phase %d: expected name %q, got %q", i+1, want[i], phase.Name())
		}
	}
}

func TestParseTrainFlagsHealApply(t *testing.T) {
	opts, err := parseTrainFlags([]string{"--heal-apply", "stash", "--text-only"})
	if err != nil {
		t.Fatalf("expected nil error for --heal-apply stash, got: %v", err)
	}
	if opts.HealApply != "stash" {
		t.Errorf("expected HealApply stash, got %q", opts.HealApply)
	}
	if !opts.IsTextOnly {
		t.Errorf("expected IsTextOnly true")
	}
}

func TestParseTrainFlagsHealApplyDefault(t *testing.T) {
	opts, err := parseTrainFlags([]string{})
	if err != nil {
		t.Fatalf("expected nil error for empty args, got: %v", err)
	}
	if opts.HealApply != "" {
		t.Errorf("expected empty HealApply default (report-only), got %q", opts.HealApply)
	}
}

func TestBucketHealReason(t *testing.T) {
	cases := map[string]string{
		"dirty: +1 modified, +2 untracked": "untracked",
		"diverged (+2/-1)":                 "diverged",
		"behind (3)":                       "diverged",
		"ahead (1)":                        "diverged",
		"uncommitted changes":              "dirty-worktree",
		"merge conflict":                   "dirty-worktree",
		"something unexpected":             "other",
	}
	for reason, want := range cases {
		if got := bucketHealReason(reason); got != want {
			t.Errorf("bucketHealReason(%q): expected %q, got %q", reason, want, got)
		}
	}
}

func TestStripTrainTargetMuse(t *testing.T) {
	target, rest := stripTrainTarget([]string{"muse", "--days", "3"})
	if target != "muse" {
		t.Errorf("expected target muse, got %q", target)
	}
	if len(rest) != 2 || rest[0] != "--days" || rest[1] != "3" {
		t.Errorf("expected rest [--days 3], got %v", rest)
	}
}

func TestStripTrainTargetOnlyFirst(t *testing.T) {
	target, rest := stripTrainTarget([]string{"--text-only", "muse"})
	if target != "muse" {
		t.Errorf("expected target muse, got %q", target)
	}
	if len(rest) != 1 || rest[0] != "--text-only" {
		t.Errorf("expected rest [--text-only], got %v", rest)
	}
}

func TestStripTrainTargetPreservesFlagValue(t *testing.T) {
	target, rest := stripTrainTarget([]string{"--skill-path", "muse"})
	if target != "" {
		t.Errorf("expected empty target for flag value, got %q", target)
	}
	if len(rest) != 2 {
		t.Errorf("expected both args preserved, got %v", rest)
	}
}

func TestParseTrainFlagsDaysDefault(t *testing.T) {
	opts, err := parseTrainFlags([]string{})
	if err != nil {
		t.Fatalf("expected nil error for empty args, got: %v", err)
	}
	if opts.DaysN != DefaultTrainMuseDays {
		t.Errorf("expected DaysN %d, got %d", DefaultTrainMuseDays, opts.DaysN)
	}
	if opts.DaysN != 7 {
		t.Errorf("expected DaysN 7, got %d", opts.DaysN)
	}
}

func TestParseTrainFlagsDaysValue(t *testing.T) {
	opts, err := parseTrainFlags([]string{"--days", "14"})
	if err != nil {
		t.Fatalf("expected nil error for --days 14, got: %v", err)
	}
	if opts.DaysN != 14 {
		t.Errorf("expected DaysN 14, got %d", opts.DaysN)
	}
}

func TestParseTrainFlagsDaysClamp(t *testing.T) {
	opts, err := parseTrainFlags([]string{"--days", "0"})
	if err != nil {
		t.Fatalf("expected nil error for --days 0, got: %v", err)
	}
	if opts.DaysN != DefaultTrainMuseDays {
		t.Errorf("expected DaysN clamped to %d, got %d", DefaultTrainMuseDays, opts.DaysN)
	}
}

func initMuseFixtureRepo(t *testing.T) string {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git binary required for muse train fixture")
	}
	dir := t.TempDir()
	runFixtureGit(t, dir, "init")
	runFixtureGit(t, dir, "-c", "user.email=muse@t", "-c", "user.name=muse", "commit", "--allow-empty", "-m", "seed")
	return dir
}

func runFixtureGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v failed: %v\n%s", args, err, out)
	}
}

func captureTrainOutput(t *testing.T, args []string) string {
	t.Helper()
	old := os.Stdout
	r, w, pipeErr := os.Pipe()
	if pipeErr != nil {
		t.Fatalf("os.Pipe failed: %v", pipeErr)
	}
	os.Stdout = w
	runErr := RunTrain(args)
	_ = w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	if runErr != nil {
		t.Fatalf("expected nil error, got: %v", runErr)
	}
	return buf.String()
}

func TestRunTrainMuseMissingDigest(t *testing.T) {
	fixture := initMuseFixtureRepo(t)
	old := repoSecretsPathResolver
	repoSecretsPathResolver = func() (string, bool, *apperror.AppError) {
		return fixture, false, nil
	}
	defer func() { repoSecretsPathResolver = old }()

	out := captureTrainOutput(t, []string{"muse", "--days", "3"})
	if !strings.Contains(out, "No Muse/ digest directory") {
		t.Errorf("expected missing-digest guidance, got:\n%s", out)
	}
	if !strings.Contains(out, "gitmap llm train muse --days 7") {
		t.Errorf("expected re-run hint, got:\n%s", out)
	}
}
