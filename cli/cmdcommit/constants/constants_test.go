package constants

import (
	"regexp"
	"testing"
)

func TestCommitFlagsConstants(t *testing.T) {
	flags := []string{
		FlagAgentID, FlagHeatmapID, FlagAI, FlagAIAgent,
		FlagSpec, FlagTask, FlagAuthor, FlagPush, FlagDryRun,
	}
	for _, f := range flags {
		if f == "" {
			t.Errorf("expected non-empty flag constant")
		}
	}
}

func TestCommitEnvConstants(t *testing.T) {
	envs := []string{
		EnvGitmapAgentID,
		EnvAntigravityAgentID,
		EnvGitmapHeatmapID,
		EnvAIModel,
	}

	for _, e := range envs {
		if e == "" {
			t.Errorf("expected non-empty env constant")
		}
	}
}

func TestCommitPatternRegex(t *testing.T) {
	if _, err := regexp.Compile(PatternHyphen); err != nil {
		t.Fatalf("PatternHyphen failed to compile: %v", err)
	}

	if _, err := regexp.Compile(PatternConventional); err != nil {
		t.Fatalf("PatternConventional failed to compile: %v", err)
	}
}
