// Package cmdcg provides unified synchronization for coding guidelines, prompts, and skills.
package cmdcg

// NOTE: the cgSyncHeaderStyle/cgSyncSuccessStyle vars were deleted here
// on 31-Oct-26 - they were unused (golangci-lint unused findings on the
// v6.643.0 baseline diff) and keeping them would keep failing that gate.
import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdsync"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdsync/enums"
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
