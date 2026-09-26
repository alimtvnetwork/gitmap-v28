// Package cmdagy — agy_recreate_ops.go executes the complete 4-phase recreate lifecycle.
package cmdagy

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/workspacesync"
)

// ExecuteAgyRecreate orchestrates recreation across all resolved project targets.
func ExecuteAgyRecreate(targets []AgyProject, opts AgyRecreateOptions) error {
	if len(targets) == 0 {
		return apperror.NewSimple("no projects selected to recreate", "E9001")
	}
	var lastErr error
	successCount := 0

	for _, p := range targets {
		res, err := processSingleRecreate(p, opts)
		if err != nil {
			lastErr = err
			renderRecreateError(p.Name, err)
			continue
		}
		successCount++
		renderRecreateSuccess(res)
	}

	if successCount == 0 && lastErr != nil {
		return apperror.WrapSimple(lastErr, "all project recreate operations failed")
	}
	return nil
}

func processSingleRecreate(p AgyProject, opts AgyRecreateOptions) (*AgyRecreateResult, error) {
	renderRecreateBanner(p)
	if opts.IsDryRun {
		return renderDryRunPreview(p), nil
	}
	purged := purgeCacheAndConversations(p)
	removeProjectRegistration(p.ID)
	newId, regErr := registerFreshProject(p)
	if regErr != nil {
		return nil, regErr
	}
	p.ID = newId
	convId, dispErr := dispatchRecreateConversation(p, opts)
	if dispErr != nil {
		return nil, dispErr
	}
	return buildRecreateResult(p, purged, convId), nil
}

func purgeCacheAndConversations(p AgyProject) int {
	purged, _ := purgeProjectConversations(p)
	cleanStagedPromptFiles(p.GetPath())
	cleanTempPromptFiles(p.ID)
	fmt.Printf("  %s→ Phase 1:%s Purged %d conversation(s) and cleared cache\n",
		constants.ColorDim, constants.ColorReset, purged)
	return purged
}

func cleanStagedPromptFiles(projectPath string) {
	if projectPath == "" {
		return
	}
	_ = os.Remove(filepath.Join(projectPath, activeAgyPromptRelativePath))
	_ = os.Remove(filepath.Join(projectPath, "agy-prompt-queue.json"))
}

func cleanTempPromptFiles(projectId string) {
	if projectId == "" {
		return
	}
	promptFile := filepath.Join(os.TempDir(), fmt.Sprintf("agy-read-prompt-%s.txt", projectId))
	_ = os.Remove(promptFile)
}

func removeProjectRegistration(projectId string) {
	if projectId == "" {
		return
	}
	_ = deleteProjectFile(projectId)
	fmt.Printf("  %s→ Phase 2:%s Removed old project configuration (%s)\n",
		constants.ColorDim, constants.ColorReset, shortProjectId(projectId))
}

func registerFreshProject(p AgyProject) (string, error) {
	projectPath := p.GetPath()
	if !workspacesync.SyncAntigravity(projectPath, p.Name) {
		return "", apperror.NewSimple("failed to register project with Antigravity", "E9003")
	}
	newId := discoverFreshProjectId(projectPath)
	fmt.Printf("  %s→ Phase 3:%s Registered fresh Antigravity project (%s)\n",
		constants.ColorDim, constants.ColorReset, shortProjectId(newId))
	return newId, nil
}

func discoverFreshProjectId(projectPath string) string {
	configDir, err := getProjectsDirPath()
	if err != nil {
		return ""
	}
	fileURI := buildFolderURI(projectPath)
	return workspacesync.FindExistingProjectID(configDir, fileURI)
}

func buildRecreateResult(p AgyProject, purged int, convId string) *AgyRecreateResult {
	fmt.Printf("  %s→ Phase 4:%s Created new conversation session\n",
		constants.ColorDim, constants.ColorReset)
	return &AgyRecreateResult{
		ProjectName:  p.Name,
		ProjectPath:  p.GetPath(),
		NewProjectId: p.ID,
		PurgedConvs:  purged,
		ConvId:       convId,
		IsSuccess:    true,
	}
}

func renderRecreateBanner(p AgyProject) {
	fmt.Printf("\n  %s┌── Antigravity Recreate Project ──────────────────────────────────────┐%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  │ Project:    %-57s │\n", p.Name)
	fmt.Printf("  │ Path:       %-57s │\n", p.GetPath())
	fmt.Printf("  │ Action:     Purge cache/convs → Re-add to AGY → Read Memory conv     │\n")
	fmt.Printf("  %s└──────────────────────────────────────────────────────────────────────┘%s\n", constants.ColorCyan, constants.ColorReset)
}

func renderDryRunPreview(p AgyProject) *AgyRecreateResult {
	fmt.Printf("  %s[dry-run] Would purge cache, re-register '%s', and spawn read-memory conv%s\n",
		constants.ColorYellow, p.Name, constants.ColorReset)
	return &AgyRecreateResult{ProjectName: p.Name, ProjectPath: p.GetPath(), IsSuccess: true}
}

func renderRecreateSuccess(res *AgyRecreateResult) {
	fmt.Printf("  %s Successfully recreated project '%s'!\n",
		constants.ColorGreen+"✓"+constants.ColorReset, res.ProjectName)
	if res.ConvId != "" && res.ConvId != "offline" && res.ConvId != "staged-offline" {
		fmt.Printf("    • New Conversation ID: %s%s%s\n", constants.ColorCyan, res.ConvId, constants.ColorReset)
	}
	fmt.Printf("    • Staged Read Memory:  Thorough inspection of files, memory & specs\n\n")
}

func renderRecreateError(name string, err error) {
	fmt.Printf("  %s Failed recreating project '%s': %v\n\n",
		constants.ColorRed+"✗"+constants.ColorReset, name, err)
}
