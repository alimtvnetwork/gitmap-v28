// Package cmdclone — clone_profile_probe.go implements multi-login profile probing and fallback.
package cmdclone

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

type profileFailure struct {
	name string
	err  error
}

var gitLsRemoteProber = probeProfileLsRemote

func probeProfileLsRemote(ctx context.Context, prof model.GitProfile, url string) (bool, error) {
	probeURL := resolveProfileCloneURL(prof, url)
	cmd := exec.CommandContext(ctx, "git", "ls-remote", "--heads", probeURL)
	cmd.Env = buildProfileGitEnv(prof)
	err := cmd.Run()
	if err == nil {
		return true, nil
	}
	return false, err
}

func resolveProfileCloneURL(prof model.GitProfile, url string) string {
	if !hasProfileSSH(prof) {
		return url
	}
	converted, ok := ConvertURLToSSH(url)
	if ok {
		return converted
	}
	return url
}

func hasProfileSSH(prof model.GitProfile) bool {
	hasKeyPath := len(prof.SSHKeyPath) > 0
	hasKeyName := len(prof.SSHKeyName) > 0
	hasSSHMethod := prof.AuthMethod == "ssh"
	return hasKeyPath || hasKeyName || hasSSHMethod
}

func buildProfileGitEnv(prof model.GitProfile) []string {
	env := envWithOverride("GIT_TERMINAL_PROMPT", "0")
	sshCmd := resolveProfileSSHCommand(prof)
	if len(sshCmd) > 0 {
		env = append(env, "GIT_SSH_COMMAND="+sshCmd)
	}
	return env
}

func resolveProfileSSHCommand(prof model.GitProfile) string {
	if len(prof.SSHKeyPath) > 0 {
		return fmt.Sprintf("ssh -i %s -F /dev/null -o IdentitiesOnly=yes", filepath.ToSlash(prof.SSHKeyPath))
	}
	if len(prof.SSHKeyName) > 0 {
		return resolveSSHCommandByName(prof.SSHKeyName)
	}
	return ""
}

func resolveSSHCommandByName(keyName string) string {
	mainDB, err := openDB()
	if err != nil {
		return ""
	}
	defer mainDB.Close()
	key, findErr := mainDB.FindSSHKeyByName(keyName)
	if findErr == nil && len(key.PrivatePath) > 0 {
		return fmt.Sprintf("ssh -i %s -o IdentitiesOnly=yes", filepath.ToSlash(key.PrivatePath))
	}
	return ""
}

func isAuthOrNotFoundError(err error, output string) bool {
	combined := strings.ToLower(output)
	if err != nil {
		combined += " " + strings.ToLower(err.Error())
	}
	return hasAuthFailureKeyword(combined)
}

func hasAuthFailureKeyword(msg string) bool {
	keywords := []string{
		"repository not found", "401", "404", "could not read username",
		"permission denied", "authentication failed", "unauthorized",
		"access denied", "not found",
	}
	for _, kw := range keywords {
		if strings.Contains(msg, kw) {
			return true
		}
	}
	return false
}

func ResolveDefaultProfile(cfg model.GitProfileConfig) *model.GitProfile {
	prof := findProfileByName(cfg.Profiles, cfg.Default)
	if prof != nil {
		return prof
	}
	prof = findDefaultFlagProfile(cfg.Profiles)
	if prof != nil {
		return prof
	}
	if len(cfg.Profiles) > 0 {
		return &cfg.Profiles[0]
	}
	return nil
}

func findProfileByName(profiles []model.GitProfile, name string) *model.GitProfile {
	if len(name) == 0 {
		return nil
	}
	for i := range profiles {
		if strings.EqualFold(profiles[i].Name, name) {
			return &profiles[i]
		}
	}
	return nil
}

func findDefaultFlagProfile(profiles []model.GitProfile) *model.GitProfile {
	for i := range profiles {
		if profiles[i].IsDefault {
			return &profiles[i]
		}
	}
	return nil
}

func CollectSecondaryProfiles(cfg model.GitProfileConfig, defaultProf *model.GitProfile) []model.GitProfile {
	var secondaries []model.GitProfile
	for _, p := range cfg.Profiles {
		if isSecondaryProfileCandidate(p, defaultProf) {
			secondaries = append(secondaries, p)
		}
	}
	return secondaries
}

func isSecondaryProfileCandidate(p model.GitProfile, defaultProf *model.GitProfile) bool {
	if defaultProf == nil {
		return true
	}
	isMatch := strings.EqualFold(p.Name, defaultProf.Name) || (len(p.ID) > 0 && p.ID == defaultProf.ID)
	if isMatch {
		return false
	}
	return true
}

func ProbeSecondaryProfiles(ctx context.Context, profiles []model.GitProfile, url string) (*model.GitProfile, error) {
	if len(profiles) == 0 {
		return nil, apperror.NewNotFoundError("no secondary profiles configured")
	}
	probeCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	winner, failures := launchParallelProbes(probeCtx, profiles, url, cancel)
	if winner != nil {
		return winner, nil
	}
	return nil, buildAggregatedDiagnosticError(failures, url)
}

func launchParallelProbes(
	ctx context.Context,
	profiles []model.GitProfile,
	url string,
	cancel context.CancelFunc,
) (*model.GitProfile, []profileFailure) {
	var once sync.Once
	var winningProfile *model.GitProfile
	var wg sync.WaitGroup
	var mu sync.Mutex
	var failures []profileFailure

	for _, prof := range profiles {
		wg.Add(1)
		go func(p model.GitProfile) {
			defer wg.Done()
			executeSingleProbe(ctx, p, url, &once, &winningProfile, cancel, &mu, &failures)
		}(prof)
	}
	wg.Wait()
	return winningProfile, failures
}

func executeSingleProbe(
	ctx context.Context,
	p model.GitProfile,
	url string,
	once *sync.Once,
	winner **model.GitProfile,
	cancel context.CancelFunc,
	mu *sync.Mutex,
	failures *[]profileFailure,
) {
	isSuccess, err := gitLsRemoteProber(ctx, p, url)
	if isSuccess {
		once.Do(func() {
			found := p
			*winner = &found
			cancel()
		})
		return
	}
	mu.Lock()
	*failures = append(*failures, profileFailure{name: p.Name, err: err})
	mu.Unlock()
}

func buildAggregatedDiagnosticError(failures []profileFailure, url string) error {
	var msgs []string
	for _, f := range failures {
		errMsg := "probing failed"
		if f.err != nil {
			errMsg = f.err.Error()
		}
		msgs = append(msgs, fmt.Sprintf("[%s: %s]", f.name, errMsg))
	}
	diag := fmt.Sprintf("all secondary profiles failed probing %s: %s", url, strings.Join(msgs, ", "))
	return apperror.NewSimple(diag, "E4041")
}

func ApplyProfileCredentials(prof model.GitProfile, url string) (string, func()) {
	resolvedURL := resolveProfileCloneURL(prof, url)
	sshCmd := resolveProfileSSHCommand(prof)
	prevSSHCmd := os.Getenv(constants.EnvGitSSHCommand)
	if len(sshCmd) > 0 {
		_ = os.Setenv(constants.EnvGitSSHCommand, sshCmd)
	}
	restore := func() {
		if len(sshCmd) > 0 {
			_ = os.Setenv(constants.EnvGitSSHCommand, prevSSHCmd)
		}
	}
	return resolvedURL, restore
}
