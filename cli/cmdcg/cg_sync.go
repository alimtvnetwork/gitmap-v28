// Package cmdcg provides unified synchronization for coding guidelines, prompts, and skills.
package cmdcg

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdsync"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdsync/enums"
	"github.com/charmbracelet/lipgloss"
)

var (
	cgSyncHeaderStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#bd93f9"))
	cgSyncSuccessStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#50fa7b"))
)

// RunCGSync executes a full 3-way synchronization across specs, prompts, skills, and scripts.
func RunCGSync(repos []string, isDryRun bool) error {
	args := buildSyncArgs(repos, isDryRun)

	return cmdsync.RunSyncWithScope(enums.SyncScopeAll, args)
}

// RunCGPromptSync executes 3-way synchronization restricted strictly to 01-prompts/.
func RunCGPromptSync(repos []string, isDryRun bool) error {
	args := buildSyncArgs(repos, isDryRun)

	return cmdsync.RunSyncWithScope(enums.SyncScopePrompts, args)
}

// RunCGSkillSync executes 3-way synchronization restricted to .agents/skills/ and .cursor/skills/.
func RunCGSkillSync(repos []string, isDryRun bool) error {
	args := buildSyncArgs(repos, isDryRun)

	return cmdsync.RunSyncWithScope(enums.SyncScopeSkills, args)
}

// RunFastGitmapSkillSync updates strictly the .agents/skills/gitmap/ skill bundle.
func RunFastGitmapSkillSync(repos []string, isDryRun bool) error {
	args := buildSyncArgs(repos, isDryRun)

	return cmdsync.RunSyncWithScope(enums.SyncScopeGitmapSkill, args)
}

func buildSyncArgs(repos []string, isDryRun bool) []string {
	var args []string
	if isDryRun {
		args = append(args, "--dry-run")
	}
	if len(repos) > 0 {
		args = append(args, "--repo", strings.Join(repos, ","))
	}

	return args
}
