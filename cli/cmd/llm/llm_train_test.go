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

