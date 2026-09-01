// Package cmd — chrome_reset.go: clean, wipe, or reset Chrome profile caches,
// cookies, history, extensions, or entire profile state.
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

type chromeResetOptions struct {
	Profile      string
	IsCache      bool
	IsCookies    bool
	IsHistory    bool
	IsExtensions bool
	IsAll        bool
}

func runChromeReset(args []string) error {
	checkHelp(constants.SubCmdChromeReset, args)
	opts := parseChromeResetArgs(args)
	profName := opts.Profile
	if profName == "" {
		profName = constants.ChromeDefaultProfileDir
	}
	srcPath, hasDir := resolveChromeProfileDir(profName)
	if !hasDir {
		return apperror.NewSimple(fmt.Sprintf("profile %s not found", profName), "E4401")
	}
	return executeProfileReset(srcPath, profName, opts)
}

func parseChromeResetArgs(args []string) chromeResetOptions {
	var opts chromeResetOptions
	for _, a := range args {
		if a == "--cache" || a == "-cache" {
			opts.IsCache = true
		} else if a == "--cookies" || a == "-cookies" {
			opts.IsCookies = true
		} else if a == "--history" || a == "-history" {
			opts.IsHistory = true
		} else if a == "--extensions" || a == "-extensions" {
			opts.IsExtensions = true
		} else if a == "--all" || a == "-a" {
			opts.IsAll = true
		} else if strings.HasPrefix(a, "--profile=") {
			opts.Profile = strings.TrimPrefix(a, "--profile=")
		} else if !strings.HasPrefix(a, "-") && opts.Profile == "" {
			opts.Profile = a
		}
	}
	if !opts.IsCache && !opts.IsCookies && !opts.IsHistory && !opts.IsExtensions && !opts.IsAll {
		opts.IsCache = true
	}
	return opts
}

func executeProfileReset(srcPath, profName string, opts chromeResetOptions) error {
	fmt.Printf("\n\033[1;96m▸ chrome reset\033[0m  profile \033[1m%s\033[0m (%s)\n", profName, srcPath)
	if opts.IsAll {
		return performFullProfileReset(srcPath, profName)
	}
	if opts.IsCache {
		wipeCacheDirs(srcPath)
	}
	if opts.IsCookies {
		wipeCookieFiles(srcPath)
	}
	if opts.IsHistory {
		wipeHistoryFiles(srcPath)
	}
	if opts.IsExtensions {
		wipeExtensionDir(srcPath)
	}
	fmt.Printf("\n\033[1;92m✓ reset complete\033[0m  profile %q cleaned\n", profName)
	return nil
}

func wipeCacheDirs(srcPath string) {
	cacheFolders := []string{"Cache", "Code Cache", "GPUCache", "DawnCache", "ShaderCache", "GrShaderCache"}
	cleaned := 0
	for _, folder := range cacheFolders {
		target := filepath.Join(srcPath, folder)
		if _, err := os.Stat(target); err == nil {
			_ = os.RemoveAll(target)
			cleaned++
		}
	}
	fmt.Printf("  \033[1;92m✓\033[0m %d cache folder(s) purged\n", cleaned)
}

func wipeCookieFiles(srcPath string) {
	files := []string{"Cookies", "Cookies-journal", filepath.Join("Network", "Cookies"), filepath.Join("Network", "Cookies-journal")}
	cleaned := 0
	for _, file := range files {
		target := filepath.Join(srcPath, file)
		if _, err := os.Stat(target); err == nil {
			_ = os.Remove(target)
			cleaned++
		}
	}
	fmt.Printf("  \033[1;92m✓\033[0m %d cookie store(s) wiped\n", cleaned)
}

func wipeHistoryFiles(srcPath string) {
	files := []string{"History", "History-journal", "Visited Links", "Top Sites", "Shortcuts"}
	cleaned := 0
	for _, file := range files {
		target := filepath.Join(srcPath, file)
		if _, err := os.Stat(target); err == nil {
			_ = os.Remove(target)
			cleaned++
		}
	}
	fmt.Printf("  \033[1;92m✓\033[0m %d history database(s) cleared\n", cleaned)
}

func wipeExtensionDir(srcPath string) {
	extDir := filepath.Join(srcPath, "Extensions")
	if _, err := os.Stat(extDir); err == nil {
		_ = os.RemoveAll(extDir)
		fmt.Println("  \033[1;92m✓\033[0m extensions directory purged")
	}
}

func performFullProfileReset(srcPath, profName string) error {
	wipeCacheDirs(srcPath)
	wipeCookieFiles(srcPath)
	wipeHistoryFiles(srcPath)
	wipeExtensionDir(srcPath)
	prefPath := filepath.Join(srcPath, "Preferences")
	if _, err := os.Stat(prefPath); err == nil {
		_ = os.WriteFile(prefPath+".bak", readFileBytes(prefPath), constants.FilePermission)
		_ = os.WriteFile(prefPath, []byte("{}"), constants.FilePermission)
		fmt.Println("  \033[1;92m✓\033[0m preferences reset to default (backup saved to .bak)")
	}
	fmt.Printf("\n\033[1;92m✓ full reset complete\033[0m  profile %q completely refreshed\n", profName)
	return nil
}

func readFileBytes(path string) []byte {
	b, _ := os.ReadFile(path)
	return b
}
