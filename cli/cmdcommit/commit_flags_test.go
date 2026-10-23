package cmdcommit

import (
	"reflect"
	"testing"
)

func TestExtractCommitFlags_SpaceDelimited(t *testing.T) {
	args := []string{
		"--agent-id", "agy-01",
		"--heatmap-id", "hm-99",
		"--ai", "gemini-pro",
		"--spec", "02-spec/test.md",
		"--task", "my-task",
		"-p", "-n",
		"feat: new feature",
	}

	clean, meta := ExtractCommitFlags(args)
	assertSpaceMetaFields(t, meta)
	assertCleanMessage(t, clean, []string{"feat: new feature"})
}

func assertSpaceMetaFields(t *testing.T, meta CommitMetadata) {
	if meta.AgentID != "agy-01" || meta.HeatmapID != "hm-99" {
		t.Errorf("unexpected agent/heatmap: %+v", meta)
	}
	if meta.AIName != "gemini-pro" || meta.SpecPath != "02-spec/test.md" {
		t.Errorf("unexpected ai/spec: %+v", meta)
	}
	if meta.TaskSlug != "my-task" || !meta.HasPush || !meta.IsDryRun || !meta.IsAIMode {
		t.Errorf("unexpected flags state: %+v", meta)
	}
}

func assertCleanMessage(t *testing.T, got, want []string) {
	if !reflect.DeepEqual(got, want) {
		t.Errorf("clean message mismatch: got %v, want %v", got, want)
	}
}

func TestExtractCommitFlags_EqualsDelimited(t *testing.T) {
	args := []string{
		"--agent=worker-7",
		"--heatmap=hm-88",
		"--ai=claude-3",
		"--spec=02-spec/equals.md",
		"--task=task-abc",
		"fix: bug",
	}

	clean, meta := ExtractCommitFlags(args)
	assertEqualsMetaFields(t, meta)
	assertCleanMessage(t, clean, []string{"fix: bug"})
}

func assertEqualsMetaFields(t *testing.T, meta CommitMetadata) {
	if meta.AgentID != "worker-7" || meta.HeatmapID != "hm-88" {
		t.Errorf("unexpected agent/heatmap: %+v", meta)
	}
	if meta.AIName != "claude-3" || meta.SpecPath != "02-spec/equals.md" {
		t.Errorf("unexpected ai/spec: %+v", meta)
	}
	if meta.TaskSlug != "task-abc" || !meta.IsAIMode {
		t.Errorf("unexpected task/ai state: %+v", meta)
	}
}

func TestExtractCommitFlags_ShortFlags(t *testing.T) {
	args := []string{"-t", "short-task", "-p", "-n", "chore: tidy up"}
	clean, meta := ExtractCommitFlags(args)

	if meta.TaskSlug != "short-task" || !meta.HasPush || !meta.IsDryRun {
		t.Errorf("short flags extraction failed: %+v", meta)
	}
	assertCleanMessage(t, clean, []string{"chore: tidy up"})
}

func TestExtractCommitFlags_AmbientFallbacks(t *testing.T) {
	t.Setenv(EnvGitmapAgentID, "ambient-agent-42")
	t.Setenv(EnvGitmapHeatmapID, "ambient-hm-55")
	t.Setenv(EnvGitmapTask, "ambient-task-99")
	t.Setenv(EnvGitmapSpecPath, "02-spec/ambient.md")
	t.Setenv(EnvGitmapAIName, "ambient-ai-model")

	clean, meta := ExtractCommitFlags([]string{"docs: update"})
	assertAmbientMeta(t, meta)
	assertCleanMessage(t, clean, []string{"docs: update"})
}

func assertAmbientMeta(t *testing.T, meta CommitMetadata) {
	if meta.AgentID != "ambient-agent-42" || meta.HeatmapID != "ambient-hm-55" {
		t.Errorf("unexpected ambient agent/heatmap: %+v", meta)
	}
	if meta.TaskSlug != "ambient-task-99" || meta.SpecPath != "02-spec/ambient.md" {
		t.Errorf("unexpected ambient task/spec: %+v", meta)
	}
	if meta.AIName != "ambient-ai-model" || !meta.IsAIMode {
		t.Errorf("unexpected ambient ai: %+v", meta)
	}
}

func TestExtractCommitFlags_AIModeDetection(t *testing.T) {
	clearAllAIEnvs(t)

	_, metaNoAI := ExtractCommitFlags([]string{"chore: plain human"})
	if metaNoAI.IsAIMode {
		t.Errorf("expected IsAIMode to be false without flags or envs")
	}

	_, metaFlag := ExtractCommitFlags([]string{"--ai", "docs: ai commit"})
	if !metaFlag.IsAIMode {
		t.Errorf("expected IsAIMode to be true with --ai flag")
	}
}

func clearAllAIEnvs(t *testing.T) {
	envs := []string{
		EnvGitmapAgentID, EnvAntigravityAgentID, EnvAgentID,
		EnvGitmapAIName, EnvAIModel, EnvGeminiModel, EnvMuseModel,
		"ANTIGRAVITY_INSTANCE_ID", "GEMINI_CLI_VERSION", "MUSE_RUNNER_ID",
	}
	for _, env := range envs {
		t.Setenv(env, "")
	}
}

func TestExtractCommitFlags_NoFlagLeakage(t *testing.T) {
	args := []string{"word1", "--push", "word2", "-n", "word3", "--task=t1", "word4"}
	clean, meta := ExtractCommitFlags(args)

	assertCleanMessage(t, clean, []string{"word1", "word2", "word3", "word4"})
	if !meta.HasPush || !meta.IsDryRun || meta.TaskSlug != "t1" {
		t.Errorf("flags did not parse cleanly: %+v", meta)
	}
}
