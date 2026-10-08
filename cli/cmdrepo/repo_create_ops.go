package cmdrepo

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/workspacesync"
)

func ExecuteCreateRepo(args []string, defaultLocal bool) error {
	params, parseErr := ParseCreateParams(args, defaultLocal)
	if parseErr != nil {
		return parseErr
	}

	if initErr := InitLocalRepo(params); initErr != nil {
		return initErr
	}

	remoteURL, pushErr := PushRemoteRepo(params)
	if pushErr != nil {
		return pushErr
	}

	absDir, _ := filepath.Abs(params.LocalDir)
	applyWorkspaceSync(params, absDir)
	recordProfileUsage(params.Profile)
	if params.IsCD {
		WriteShellHandoff(absDir)
	}

	return reportCreatedRepo(params, remoteURL)
}

func applyWorkspaceSync(params CreateRepoParams, absDir string) {
	if params.IsNoSync {
		workspacesync.SyncAgyOnly(absDir, params.Name)
		return
	}
	if params.IsNoDesktop {
		workspacesync.SyncWithoutDesktop(absDir, params.Name)
		return
	}
	workspacesync.SyncAll(absDir, params.Name)
}

func WriteShellHandoff(targetPath string) {
	handoffFile := os.Getenv(constants.EnvGitmapHandoffFile)
	if len(handoffFile) == 0 {
		return
	}
	if len(targetPath) == 0 {
		return
	}
	if err := os.WriteFile(handoffFile, []byte(targetPath), 0o600); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrShellHandoffWriteFmt, handoffFile, err)
	}
}
