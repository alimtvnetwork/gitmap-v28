package cmdchromeprofile

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func countBookmarks(raw json.RawMessage) int {
	if len(raw) == 0 {
		return 0
	}

	var root bookmarkRoot
	if err := json.Unmarshal(raw, &root); err != nil {
		return 0
	}

	count := 0
	for _, node := range root.Roots {
		count += countNodeBookmarks(node)
	}

	return count
}

func countNodeBookmarks(node bookmarkNode) int {
	count := 0
	if node.Type == "url" {
		count++
	}

	for _, child := range node.Children {
		count += countNodeBookmarks(child)
	}

	return count
}

func hasProfileBookmarks(profileDir string) bool {
	bmPath := filepath.Join(profileDir, "Bookmarks")
	info, err := os.Stat(bmPath)
	if err != nil {
		return false
	}

	return info.Size() > 50
}

func findNextAvailableProfileDir() string {
	stateSet := chromeLocalStateProfileDirs()
	root := chromeUserDataDir()
	for n := 1; n < 1000; n++ {
		candidate := fmt.Sprintf("Profile %d", n)
		if stateSet[candidate] {
			continue
		}

		if chromeProfilePathExists(filepath.Join(root, candidate)) {
			continue
		}

		return candidate
	}

	return fmt.Sprintf("Profile %d", time.Now().Unix())
}
