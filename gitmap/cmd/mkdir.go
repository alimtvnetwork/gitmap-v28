package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/macro"
)

type mkdirConfig struct {
	hasParents bool
	hasFile    bool
	isVerbose  bool
	targets    []string
}

func runMkdir(args []string) error {
	cfg := parseMkdirConfig(args)
	if len(cfg.targets) == 0 {
		return apperror.NewSimple("Usage: gitmap mkdir [-p|--parents] [-f|--file] [-v|--verbose] <path>...", "E9000")
	}
	cwd, _ := os.Getwd()
	for _, rawTarget := range cfg.targets {
		target := macro.NormalizeTargetPath(rawTarget, cwd)
		if err := executeMkdirTarget(target, cfg); err != nil {
			return err
		}
	}

	return nil
}

func parseMkdirConfig(args []string) mkdirConfig {
	cfg := mkdirConfig{}
	for _, arg := range args {
		if isMkdirFlag(arg, &cfg) {
			continue
		}
		if len(arg) > 0 {
			cfg.targets = append(cfg.targets, arg)
		}
	}

	return cfg
}

func isMkdirFlag(arg string, cfg *mkdirConfig) bool {
	if arg == "-p" || arg == "--parents" {
		cfg.hasParents = true
		return true
	}
	if arg == "-f" || arg == "--file" {
		cfg.hasFile = true
		cfg.hasParents = true
		return true
	}
	if arg == "-v" || arg == "--verbose" {
		cfg.isVerbose = true
		return true
	}

	return false
}

func executeMkdirTarget(target string, cfg mkdirConfig) error {
	if cfg.hasFile {
		return createTargetFile(target, cfg)
	}

	return createTargetDirectory(target, cfg)
}

func createTargetDirectory(target string, cfg mkdirConfig) error {
	if err := dispatchDirectoryCreation(target, cfg); err != nil {
		return err
	}
	fmt.Printf("✔ Created: %s\n", target)

	return nil
}

func dispatchDirectoryCreation(target string, cfg mkdirConfig) error {
	if cfg.hasParents {
		return createMissingAncestors(target, cfg.isVerbose)
	}

	return createSingleDirectory(target, cfg.isVerbose)
}

func createSingleDirectory(path string, isVerbose bool) error {
	info, err := os.Stat(path)
	if err == nil && info.IsDir() {
		fmt.Printf("  [DIR]  exists:  %s\n", path)
		return nil
	}
	if err := os.Mkdir(path, 0755); err != nil && !os.IsExist(err) {
		return apperror.WrapSimple(err, fmt.Sprintf("failed to create directory %s", path))
	}
	fmt.Printf("  [DIR]  created: %s\n", path)

	return nil
}

func collectMissingAncestors(target string) []string {
	var missing []string
	curr := target
	for len(curr) > 0 && curr != filepath.Dir(curr) {
		info, err := os.Stat(curr)
		if err == nil && info.IsDir() {
			break
		}
		missing = append([]string{curr}, missing...)
		curr = filepath.Dir(curr)
	}

	return missing
}

func createMissingAncestors(target string, isVerbose bool) error {
	missing := collectMissingAncestors(target)
	if len(missing) == 0 {
		fmt.Printf("  [DIR]  exists:  %s\n", target)
		return nil
	}
	for _, dir := range missing {
		if err := os.Mkdir(dir, 0755); err != nil && !os.IsExist(err) {
			return apperror.WrapSimple(err, fmt.Sprintf("failed to create directory %s", dir))
		}
		fmt.Printf("  [DIR]  created: %s\n", dir)
	}

	return nil
}

func createTargetFile(filePath string, cfg mkdirConfig) error {
	parentDir := filepath.Dir(filePath)
	if err := createMissingAncestors(parentDir, cfg.isVerbose); err != nil {
		return err
	}
	if err := touchFile(filePath); err != nil {
		return err
	}
	fmt.Printf("✔ Created: %s\n", filePath)

	return nil
}

func touchFile(filePath string) error {
	info, err := os.Stat(filePath)
	if err == nil && !info.IsDir() {
		fmt.Printf("  [FILE] exists:  %s\n", filePath)
		return nil
	}
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return apperror.WrapSimple(err, fmt.Sprintf("failed to create file %s", filePath))
	}
	defer f.Close()
	fmt.Printf("  [FILE] created: %s\n", filePath)

	return nil
}
