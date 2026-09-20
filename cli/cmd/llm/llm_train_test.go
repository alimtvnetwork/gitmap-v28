package llm

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
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
