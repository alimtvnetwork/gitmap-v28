package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// PathResult represents the JSON payload for high-speed path queries.
type PathResult struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	IsValid bool   `json:"isValid"`
}

type pathCLIOpts struct {
	target string
	isJSON bool
	isHelp bool
}

func parsePathCLIOpts(args []string) pathCLIOpts {
	var opts pathCLIOpts
	for _, a := range args {
		consumePathArg(&opts, a)
	}
	return opts
}

func consumePathArg(opts *pathCLIOpts, a string) {
	clean := strings.ToLower(strings.TrimSpace(a))
	if clean == "--json" || clean == "-j" {
		opts.isJSON = true
		return
	}
	if clean == "--help" || clean == "-h" || clean == "help" {
		opts.isHelp = true
		return
	}
	if !strings.HasPrefix(clean, "-") && len(opts.target) == 0 {
		opts.target = strings.TrimSpace(a)
	}
}

func printPathHelp() error {
	fmt.Println("GitMap Native Fast Path Resolver")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  gitmap path <target> [--json]")
	fmt.Println()
	fmt.Println("Targets:")
	fmt.Println("  rc, repo-cache       Path to repo-cache directory")
	fmt.Println("  rs, repo-secrets     Path to repo-secrets directory")
	fmt.Println("  <repo|alias>         Path to repository or alias")
	return nil
}

func resolveSpecialPathTarget(clean string) (string, bool) {
	switch clean {
	case "rc", "repo-cache":
		p, err := resolveSpecialRepoCDPath("rc")
		return p, err == nil && len(p) > 0
	case "rs", "repo-secrets":
		p, err := resolveSpecialRepoCDPath("rs")
		return p, err == nil && len(p) > 0
	default:
		return "", false
	}
}

func queryDBForRepoPath(db *store.DB, clean string) (string, bool) {
	if aliasRow, err := db.ResolveAlias(clean); err == nil && len(aliasRow.AbsolutePath) > 0 {
		return aliasRow.AbsolutePath, true
	}
	if repos, err := db.FindBySlug(clean); err == nil && len(repos) > 0 {
		return repos[0].AbsolutePath, true
	}
	all, listErr := db.ListRepos()
	if listErr == nil {
		for _, r := range all {
			if strings.EqualFold(r.RepoName, clean) || strings.EqualFold(r.Slug, clean) {
				return r.AbsolutePath, true
			}
		}
	}
	return "", false
}

func resolveRepoOrLocalPath(clean string) (string, error) {
	if specialPath, isSpecial := resolveSpecialPathTarget(clean); isSpecial {
		return specialPath, nil
	}
	db, err := openDB()
	if err == nil {
		defer db.Close()
		if p, hasRepo := queryDBForRepoPath(db, clean); hasRepo {
			return p, nil
		}
	}
	return checkLocalDirFallback(clean)
}

func checkLocalDirFallback(clean string) (string, error) {
	if info, err := os.Stat(clean); err == nil && info.IsDir() {
		if abs, absErr := filepath.Abs(clean); absErr == nil {
			return abs, nil
		}
		return clean, nil
	}
	return "", fmt.Errorf("repository or target %q not found", clean)
}

func outputJSONPathResult(target, resolved string, isValid bool) error {
	res := PathResult{Name: target, Path: resolved, IsValid: isValid}
	data, _ := json.MarshalIndent(res, "", "  ")
	fmt.Println(string(data))
	if !isValid {
		return fmt.Errorf("target %q not found", target)
	}
	return nil
}

// RunPathCLI resolves and prints the requested target path.
func RunPathCLI(args []string) error {
	opts := parsePathCLIOpts(args)
	if opts.isHelp || len(opts.target) == 0 {
		return printPathHelp()
	}
	clean := strings.ToLower(opts.target)
	resolved, err := resolveRepoOrLocalPath(clean)
	isValid := err == nil && len(resolved) > 0
	if opts.isJSON {
		return outputJSONPathResult(opts.target, resolved, isValid)
	}
	if !isValid {
		return err
	}
	fmt.Println(resolved)
	return nil
}
