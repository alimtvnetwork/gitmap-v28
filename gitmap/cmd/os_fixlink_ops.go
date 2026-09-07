package cmd

import (
	"os"
	"path/filepath"
)

// LinkResult records the diagnostic and repair state of one symlink.
type LinkResult struct {
	Path       string `json:"path"`
	Target     string `json:"target"`
	OldTarget  string `json:"old_target,omitempty"`
	IsHealthy  bool   `json:"is_healthy"`
	IsRepaired bool   `json:"is_repaired"`
	IsBroken   bool   `json:"is_broken"`
	Message    string `json:"message"`
}

// FixLinkOptions holds operational flags for fix-link execution.
type FixLinkOptions struct {
	TargetOverride string
	IsForce        bool
	IsDryRun       bool
	IsRecursive    bool
	IsJSON         bool
}

func inspectAndRepairPath(path string, opts FixLinkOptions) ([]LinkResult, error) {
	if isNonSymlinkDirectory(path) {
		return processDirectoryBrokenLinks(path, opts)
	}

	res := repairSingleLink(path, opts)

	return []LinkResult{res}, nil
}

func isNonSymlinkDirectory(path string) bool {
	fi, err := os.Lstat(path)
	if err != nil {
		return false
	}

	return fi.IsDir() && (fi.Mode()&os.ModeSymlink == 0)
}

func repairSingleLink(linkPath string, opts FixLinkOptions) LinkResult {
	fi, err := os.Lstat(linkPath)
	if os.IsNotExist(err) {
		return handleMissingLink(linkPath, opts)
	}
	if err != nil {
		return LinkResult{Path: linkPath, IsBroken: true, Message: err.Error()}
	}
	if fi.Mode()&os.ModeSymlink == 0 {
		return handleRegularFileLink(linkPath, opts)
	}

	return handleExistingSymlink(linkPath, opts)
}

func handleMissingLink(linkPath string, opts FixLinkOptions) LinkResult {
	target := resolveMissingLinkTarget(linkPath, opts.TargetOverride)
	if target == "" {
		return LinkResult{Path: linkPath, IsBroken: true, Message: "path does not exist"}
	}
	if opts.IsDryRun {
		return LinkResult{Path: linkPath, Target: target, IsRepaired: true, Message: "would create symlink"}
	}

	err := makeSymlink(linkPath, target)
	if err != nil {
		return LinkResult{Path: linkPath, Target: target, IsBroken: true, Message: err.Error()}
	}

	return LinkResult{Path: linkPath, Target: target, IsRepaired: true, Message: "created symlink"}
}

func resolveMissingLinkTarget(linkPath, override string) string {
	if override != "" {
		return override
	}
	if filepath.Base(linkPath) == "SharedDirectories" && pathExists("/mnt/hgfs") {
		return "/mnt/hgfs"
	}

	return ""
}

func handleRegularFileLink(linkPath string, opts FixLinkOptions) LinkResult {
	if !opts.IsForce || opts.TargetOverride == "" {
		return LinkResult{Path: linkPath, IsHealthy: true, Message: "file exists (not a symlink)"}
	}
	if opts.IsDryRun {
		return LinkResult{Path: linkPath, Target: opts.TargetOverride, OldTarget: "file", IsRepaired: true, Message: "would replace file with symlink"}
	}

	_ = os.Remove(linkPath)
	err := makeSymlink(linkPath, opts.TargetOverride)
	if err != nil {
		return LinkResult{Path: linkPath, Target: opts.TargetOverride, IsBroken: true, Message: err.Error()}
	}

	return LinkResult{Path: linkPath, Target: opts.TargetOverride, OldTarget: "file", IsRepaired: true, Message: "replaced file with symlink"}
}

func handleExistingSymlink(linkPath string, opts FixLinkOptions) LinkResult {
	currentTarget, err := os.Readlink(linkPath)
	if err != nil {
		return LinkResult{Path: linkPath, IsBroken: true, Message: err.Error()}
	}

	return evaluateSymlinkHealth(linkPath, currentTarget, opts)
}

func evaluateSymlinkHealth(linkPath, currentTarget string, opts FixLinkOptions) LinkResult {
	isValid := checkTargetExists(linkPath, currentTarget)
	isTargetMatch := opts.TargetOverride == "" || opts.TargetOverride == currentTarget

	if isValid && isTargetMatch {
		return LinkResult{Path: linkPath, Target: currentTarget, IsHealthy: true, Message: "valid symlink"}
	}

	newTarget := resolveNewTarget(linkPath, currentTarget, opts.TargetOverride)
	if newTarget == "" {
		return LinkResult{Path: linkPath, Target: currentTarget, IsBroken: true, Message: "dangling symlink (target missing)"}
	}
	if opts.IsDryRun {
		return LinkResult{Path: linkPath, Target: newTarget, OldTarget: currentTarget, IsRepaired: true, Message: "would repair symlink"}
	}

	_ = os.Remove(linkPath)
	err := makeSymlink(linkPath, newTarget)
	if err != nil {
		return LinkResult{Path: linkPath, Target: newTarget, OldTarget: currentTarget, IsBroken: true, Message: err.Error()}
	}

	return LinkResult{Path: linkPath, Target: newTarget, OldTarget: currentTarget, IsRepaired: true, Message: "repaired symlink"}
}

func resolveNewTarget(linkPath, currentTarget, override string) string {
	if override != "" {
		return override
	}

	return findHeuristicTarget(linkPath, currentTarget)
}

func checkTargetExists(linkPath, target string) bool {
	if filepath.IsAbs(target) {
		return pathExists(target)
	}

	cand := filepath.Join(filepath.Dir(linkPath), target)

	return pathExists(cand)
}

func makeSymlink(linkPath, target string) error {
	_ = os.MkdirAll(filepath.Dir(linkPath), 0755)
	_ = os.Remove(linkPath)

	return os.Symlink(target, linkPath)
}

func findHeuristicTarget(linkPath, currentTarget string) string {
	if filepath.Base(linkPath) == "SharedDirectories" && pathExists("/mnt/hgfs") {
		return "/mnt/hgfs"
	}

	candDir := filepath.Join(filepath.Dir(linkPath), currentTarget)
	if pathExists(candDir) {
		return candDir
	}

	candBase := filepath.Join(filepath.Dir(linkPath), filepath.Base(currentTarget))
	if pathExists(candBase) {
		return candBase
	}

	return ""
}
