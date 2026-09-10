package cmd

import (
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

func TestToolAntigravityAndAgManagerInProbeMap(t *testing.T) {
	agyCfg, hasAgy := toolProbeMap[constants.ToolAntigravity]
	if !hasAgy || len(agyCfg.bins) == 0 {
		t.Fatalf("expected ToolAntigravity in toolProbeMap, got %v", agyCfg)
	}
	agmCfg, hasAgm := toolProbeMap[constants.ToolAgManager]
	if !hasAgm || len(agmCfg.bins) == 0 {
		t.Fatalf("expected ToolAgManager in toolProbeMap, got %v", agmCfg)
	}
}

func TestToolAntigravityAndAgManagerInBinaryMap(t *testing.T) {
	if bin := toolBinaryName(constants.ToolAntigravity); bin != "antigravity" {
		t.Fatalf("expected ToolAntigravity binary name 'antigravity', got %q", bin)
	}
	if bin := toolBinaryName(constants.ToolAgy); bin != "agy" {
		t.Fatalf("expected ToolAgy binary name 'agy', got %q", bin)
	}
	if bin := toolBinaryName(constants.ToolAgManager); bin != "ag-manager" {
		t.Fatalf("expected ToolAgManager binary name 'ag-manager', got %q", bin)
	}
}

func TestBuildFallbackCandidatesContainsExpectedPaths(t *testing.T) {
	cands := buildFallbackCandidates("pnpm")
	if len(cands) == 0 {
		t.Fatal("expected non-empty fallback candidates for pnpm")
	}
	hasLocalShare := false
	for _, c := range cands {
		if strings.Contains(c, ".local") && strings.Contains(c, "pnpm") {
			hasLocalShare = true
			break
		}
	}
	if !hasLocalShare {
		t.Fatalf("expected candidates to contain .local/share/pnpm path, got %v", cands)
	}
}

func TestIsAppImageFileRecognition(t *testing.T) {
	if !isAppImageFile("/tmp/Antigravity.Tools_4.6.9_amd64.AppImage") {
		t.Fatal("expected .AppImage to be recognized")
	}
	if !isAppImageFile("app.appimage") {
		t.Fatal("expected .appimage to be recognized")
	}
	if isAppImageFile("app.deb") {
		t.Fatal("expected .deb NOT to be recognized as AppImage")
	}
}

func TestIsHtmlContentDetection(t *testing.T) {
	htmlData := []byte("<!DOCTYPE html><html><body>Error</body></html>")
	if !isHtmlContent(htmlData) {
		t.Fatal("expected <!DOCTYPE html> to be recognized as HTML")
	}
	shData := []byte("#!/usr/bin/env bash\necho 'hello'")
	if isHtmlContent(shData) {
		t.Fatal("expected bash script NOT to be recognized as HTML")
	}
}

func TestBuildVerifyAppErrorHasStackTrace(t *testing.T) {
	paths := []string{"/usr/local/bin", "/home/user/.local/bin"}
	appErr := buildVerifyAppError("pnpm", "pnpm", "/tmp/log.txt", paths)
	if appErr == nil {
		t.Fatal("expected non-nil AppError")
	}
	if appErr.Stack == "" {
		t.Fatal("expected non-empty Stack trace in AppError")
	}
	if appErr.Ctx["tool"] != "pnpm" {
		t.Fatalf("expected tool 'pnpm' in context, got %v", appErr.Ctx["tool"])
	}
}
