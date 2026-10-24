package gitutil

import (
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/fsutil"
)

// FolderSSHBinding models the persistent link between a directory and an SSH credential.
type FolderSSHBinding struct {
	ID          int64     `json:"id" db:"id"`
	FolderPath  string    `json:"folderPath" db:"folder_path"`
	WorkgroupID int64     `json:"workgroupId,omitempty" db:"workgroup_id"`
	SSHKeyID    int64     `json:"sshKeyId" db:"ssh_key_id"`
	KeyPath     string    `json:"keyPath,omitempty" db:"-"`
	AuthorName  string    `json:"authorName,omitempty" db:"author_name"`
	AuthorEmail string    `json:"authorEmail,omitempty" db:"author_email"`
	Priority    int       `json:"priority" db:"priority"`
	IsActive    bool      `json:"isActive" db:"is_active"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
}

// ResolvedGitContext contains all injected parameters resolved for a target repository.
type ResolvedGitContext struct {
	RepoPath       string   `json:"repoPath"`
	WorkdirAnchor  string   `json:"workdirAnchor"`
	BoundKeyPath   string   `json:"boundKeyPath"`
	AuthorName     string   `json:"authorName"`
	AuthorEmail    string   `json:"authorEmail"`
	GitSSHCommand  string   `json:"gitSSHCommand"`
	HasInjectedSSH bool     `json:"hasInjectedSSH"`
	Environment    []string `json:"-"`
}

// isPathPrefixMatch checks if target matches folder directly or at directory boundary.
func isPathPrefixMatch(targetKey, folderKey string) bool {
	if targetKey == folderKey {
		return true
	}

	if folderKey == "/" && strings.HasPrefix(targetKey, "/") {
		return true
	}

	if !strings.HasPrefix(targetKey, folderKey) {
		return false
	}

	return strings.HasPrefix(targetKey[len(folderKey):], "/")
}

// isBetterMatch compares candidate specificity, priority, and update timestamps.
func isBetterMatch(candidate, currentBest FolderSSHBinding, candidateLen, currentBestLen int) bool {
	if candidateLen != currentBestLen {
		return candidateLen > currentBestLen
	}

	if candidate.Priority != currentBest.Priority {
		return candidate.Priority > currentBest.Priority
	}

	return candidate.UpdatedAt.After(currentBest.UpdatedAt)
}

func updateBestMatch(b FolderSSHBinding, matchLen int, bestMatch *FolderSSHBinding, maxLen *int, hasMatch *bool) {
	if !*hasMatch || isBetterMatch(b, *bestMatch, matchLen, *maxLen) {
		*maxLen = matchLen
		*bestMatch = b
		*hasMatch = true
	}
}

// evaluateCandidate evaluates a single folder binding candidate against targetKey.
func evaluateCandidate(b FolderSSHBinding, targetKey string, bestMatch *FolderSSHBinding, maxLen *int, hasMatch *bool) {
	if !b.IsActive {
		return
	}

	folderKey := fsutil.CanonicalPathKey(b.FolderPath)
	if !isPathPrefixMatch(targetKey, folderKey) {
		return
	}

	updateBestMatch(b, len(folderKey), bestMatch, maxLen, hasMatch)
}

// FindLongestPrefixBinding identifies the most specific folder binding for targetPath.
func FindLongestPrefixBinding(targetPath string, bindings []FolderSSHBinding) (FolderSSHBinding, bool) {
	targetKey := fsutil.CanonicalPathKey(targetPath)
	var bestMatch FolderSSHBinding
	hasMatch := false
	maxLen := -1

	for _, b := range bindings {
		evaluateCandidate(b, targetKey, &bestMatch, &maxLen, &hasMatch)
	}

	return bestMatch, hasMatch
}

// ResolveGitContextFromBinding constructs a ResolvedGitContext from a matched FolderSSHBinding.
func ResolveGitContextFromBinding(repoPath string, binding FolderSSHBinding, hasMatch bool) ResolvedGitContext {
	ctx := ResolvedGitContext{RepoPath: repoPath}
	if !hasMatch || len(binding.KeyPath) == 0 {
		return ctx
	}

	ctx.WorkdirAnchor = binding.FolderPath
	ctx.BoundKeyPath = binding.KeyPath
	ctx.AuthorName = binding.AuthorName
	ctx.AuthorEmail = binding.AuthorEmail
	ctx.GitSSHCommand = BuildDeterministicSSHCommand(binding.KeyPath)
	ctx.HasInjectedSSH = true

	return ctx
}

// ResolveGitContextFromBindings matches bindings and builds the effective ResolvedGitContext.
func ResolveGitContextFromBindings(repoPath string, bindings []FolderSSHBinding) ResolvedGitContext {
	binding, hasMatch := FindLongestPrefixBinding(repoPath, bindings)

	return ResolveGitContextFromBinding(repoPath, binding, hasMatch)
}
