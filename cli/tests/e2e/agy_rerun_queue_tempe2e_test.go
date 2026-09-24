//go:build tempe2e

package e2e

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdagy"
)

func TestAgyRerun_OneRunningFiveQueued_And_Images_TempE2E(t *testing.T) {
	if os.Getenv("RUN_TEMP_E2E") != "1" {
		t.Skip("skipping temporary e2e test; run on-demand with RUN_TEMP_E2E=1 and -tags=tempe2e")
	}

	seedOneRunningAndFiveQueuedPrompts(t)

	promptWithImages := "Fix VS Code startup JSON issue and dark blue contrast. Refer to file:///C:/Users/Administrator/.gemini/antigravity/brain/5612a36d/screenshot_8BnieUEKidCr.png and C:/Users/Administrator/.gemini/antigravity/brain/5612a36d/.user_uploaded/media_1790240392825.png"
	paths := cmdagy.ExtractAllImagePathsFromPrompt(promptWithImages, nil)
	if len(paths) != 2 {
		t.Fatalf("expected 2 extracted image paths, got %d: %v", len(paths), paths)
	}

	requeued, err := cmdagy.RequeuePendingWithCheckPrefix("")
	if err != nil {
		t.Fatalf("expected RequeuePendingWithCheckPrefix to succeed, got: %v", err)
	}
	if len(requeued) != 5 {
		t.Fatalf("expected 5 queued prompts, got %d", len(requeued))
	}

	for i, q := range requeued {
		if !strings.HasPrefix(strings.ToLower(q.Prompt), "is it completed properly and released??") {
			t.Fatalf("queued item #%d missing check prefix: %s", i+1, q.Prompt)
		}
		fmt.Printf("  ✓ [Queue #%d] %s\n    Prefixed Payload: %s\n\n", i+1, q.Title, strings.ReplaceAll(q.Prompt, "\n\n", " -> "))
	}

	_ = cmdagy.RestartAndRerunProject("1", false, true, "")
}

func seedOneRunningAndFiveQueuedPrompts(t *testing.T) {
	q := cmdagy.AgyPromptQueueFile{
		Active: &cmdagy.AgyPromptQueueEntry{
			ID:     1,
			Type:   "running_prompt",
			Title:  "Active Task #1: Fix VS Code JSON & High-Contrast Install UI",
			Prompt: "Fix VS Code projects.json and contrast with screenshot_8BnieUEKidCr.png",
			Status: "running",
		},
		Queued: []cmdagy.AgyPromptQueueEntry{
			{ID: 2, Type: "queued", Title: "Queued #1: SSH Common Batch Join (sjc)", Prompt: "Implement gitmap ssh-join-common administrator 192.168.1.3(w1),7(w2),12(w3)", Status: "queued"},
			{ID: 3, Type: "queued", Title: "Queued #2: Remote OS Detection & SQLite Save", Prompt: "Probe OSVersion via os_detect and save to SSHConnection table", Status: "queued"},
			{ID: 4, Type: "queued", Title: "Queued #3: AGY Reread & Optimize Project (rop)", Prompt: "Implement gitmap agy rop 5 with split-DB backup in data/AGY/<slug>.db", Status: "queued"},
			{ID: 5, Type: "queued", Title: "Queued #4: AUM Search vs Python Grep Benchmark", Prompt: "Benchmark GitMap AUM search against Python grep and document in readme.md", Status: "queued"},
			{ID: 6, Type: "queued", Title: "Queued #5: Chrome Profile Auth Export/Import E2E", Prompt: "Verify OAuth refresh token and cookies roundtrip in local VM", Status: "queued"},
		},
	}
	if err := cmdagy.SavePromptQueue(q); err != nil {
		t.Fatalf("failed to seed queue: %v", err)
	}
	t.Cleanup(func() {
		_ = cmdagy.ClearPromptQueue()
	})
}
