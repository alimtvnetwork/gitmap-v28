package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/cmdvmware"
)

func resolveDefaultFixLinkPaths() []string {
	desktopDir := cmdvmware.ResolveUserDesktopDir()
	candidates := []string{
		filepath.Join(desktopDir, "SharedDirectories"),
		"/usr/local/bin/gitmap",
		filepath.Join(expandHome("~"), ".local", "bin", "gitmap"),
	}

	var active []string
	for _, c := range candidates {
		if pathExists(c) || isLikelyCandidate(c) {
			active = append(active, c)
		}
	}

	return active
}

func isLikelyCandidate(path string) bool {
	return filepath.Base(path) == "SharedDirectories" && pathExists("/mnt/hgfs")
}

func inspectStandardLinksStatus(desktopDir string) error {
	sharedDir := filepath.Join(desktopDir, "SharedDirectories")
	fi, err := os.Lstat(sharedDir)
	if err == nil && (fi.Mode()&os.ModeSymlink != 0) {
		target, _ := os.Readlink(sharedDir)
		fmt.Printf("  • Desktop Link:     %s -> %s (active)\n", sharedDir, target)

		return nil
	}

	fmt.Printf("  • Desktop Link:     %s (not configured)\n", sharedDir)

	return nil
}
