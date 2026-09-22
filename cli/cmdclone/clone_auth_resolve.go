// Package cmdclone — clone_auth_resolve.go resolves authentication for repositories.
package cmdclone

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/ghtoken"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

// ResolveRepoAuth performs SSH-first probe and falls back to terminal token/browser prompt.
func ResolveRepoAuth(rec model.ScanRecord) (model.ScanRecord, error) {
	rawURL := resolveRecordRemoteURL(rec)
	if rawURL == "" {
		return rec, nil
	}

	if sshURL, isAccessible := ProbeSSHAccess(rawURL); isAccessible {
		rec.Transport = "ssh"
		rec.SSHUrl = sshURL
		fmt.Printf("  %s SSH access verified for %s; bypassing token.\n", constants.ColorGreen+"✓"+constants.ColorReset, getRecordDisplayName(rec))

		return rec, nil
	}

	return resolveHTTPSAuth(rec, rawURL)
}

func resolveHTTPSAuth(rec model.ScanRecord, rawURL string) (model.ScanRecord, error) {
	if tok, isFound := GetGlobalAccessToken(); isFound {
		rec.HTTPSUrl = InjectTokenIntoHTTPS(rawURL, tok)

		return rec, nil
	}

	if tok, _, err := ghtoken.Resolve(); err == nil && tok != "" {
		rec.HTTPSUrl = InjectTokenIntoHTTPS(rawURL, tok)

		return rec, nil
	}

	promptRes, err := PromptTerminalAuth(getRecordDisplayName(rec), rawURL)
	if err != nil {
		return rec, err
	}

	rec.HTTPSUrl = InjectTokenIntoHTTPS(rawURL, promptRes.Token)

	return rec, nil
}

func getRecordDisplayName(rec model.ScanRecord) string {
	if len(rec.RepoName) > 0 {
		return rec.RepoName
	}
	if len(rec.RelativePath) > 0 {
		return rec.RelativePath
	}

	return "repository"
}
