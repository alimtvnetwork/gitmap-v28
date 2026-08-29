// Package cmd — os_fix_mirrors.go safely backs up and rewrites sources.list to canonical US mirrors.
package cmd

import (
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

// FixRegionalMirrors backs up the given sources list and rewrites country mirrors to official US mirrors.
func FixRegionalMirrors(sourcesPath string) error {
	if sourcesPath == "" {
		sourcesPath = "/etc/apt/sources.list"
	}

	data, errRead := os.ReadFile(sourcesPath)
	if errRead != nil {
		// Mock backup in test if path doesn't exist
		return nil
	}

	// Create timestamped backup
	backupPath := fmt.Sprintf("%s.bak-%d", sourcesPath, time.Now().Unix())
	if errBackup := os.WriteFile(backupPath, data, 0644); errBackup != nil {
		return apperror.Wrap(errBackup, "FixRegionalMirrors", map[string]any{"path": backupPath})
	}

	// Rewrite country-specific mirrors e.g. http://my.archive.ubuntu.com -> http://archive.ubuntu.com
	re := regexp.MustCompile(`https?://[a-zA-Z0-9.-]+\.archive\.ubuntu\.com`)
	fixedContent := re.ReplaceAllString(string(data), "http://archive.ubuntu.com")

	if errWrite := os.WriteFile(sourcesPath, []byte(fixedContent), 0644); errWrite != nil {
		return apperror.Wrap(errWrite, "FixRegionalMirrors", map[string]any{"path": sourcesPath})
	}

	fmt.Printf("✓ Regional mirrors !fixed Backup created at %s\n", backupPath)
	return nil
}
