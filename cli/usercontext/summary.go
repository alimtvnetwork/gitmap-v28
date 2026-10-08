package usercontext

import (
	"context"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// CollectUserSummary gathers full user context across GitHub and Git.
func CollectUserSummary() UserSummary {
	ghAuth := DetectGhAuth()
	localUser, _ := GetGitIdentity(false)
	globalUser, _ := GetGitIdentity(true)
	activeProfile := resolveActiveProfileName()
	boundProfile := resolveBoundProfileName()
	isInside := isInsideGitRepo()

	return UserSummary{
		GhAuth:        ghAuth,
		LocalGitUser:  localUser,
		GlobalGitUser: globalUser,
		ActiveProfile: activeProfile,
		BoundProfile:  boundProfile,
		IsInsideRepo:  isInside,
	}
}

func resolveActiveProfileName() string {
	cfg, err := store.LoadGitProfiles()
	if err != nil {
		return "default"
	}

	if len(cfg.Active) > 0 {
		return cfg.Active
	}

	if len(cfg.Default) > 0 {
		return cfg.Default
	}

	return "default"
}

func resolveBoundProfileName() string {
	prof, isBound := ResolveProjectUser(".")
	if !isBound || prof == nil {
		return "none"
	}

	if len(prof.Name) > 0 {
		return prof.Name
	}

	return prof.ID
}

func isInsideGitRepo() bool {
	ctx, cancel := context.WithTimeout(context.Background(), gitCommandTimeout)
	defer cancel()

	out, err := execGitCommand(ctx, "rev-parse", "--is-inside-work-tree")
	if err != nil {
		return false
	}

	return strings.TrimSpace(string(out)) == "true"
}
