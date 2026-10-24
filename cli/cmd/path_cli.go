package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
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
	if constants.IsJSONFlag(clean) {
		opts.isJSON = true

		return
	}

	if constants.IsHelpFlag(clean) {
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
	special := constants.SpecialPathType(clean)
	if special.IsRepoCache() {
		p, err := resolveSpecialRepoCDPath(string(constants.SpecialPathRepoCacheShort))

		return p, err == nil && len(p) > 0
	}

	if special.IsRepoSecrets() {
		p, err := resolveSpecialRepoCDPath(string(constants.SpecialPathRepoSecretsShort))

		return p, err == nil && len(p) > 0
	}

	return "", false
}

func queryDBForRepoPath(db *store.DB, clean string) (string, bool) {
	aliasRow, err := db.ResolveAlias(clean)
	if err == nil && len(aliasRow.AbsolutePath) > 0 {
		return aliasRow.AbsolutePath, true
	}

	repos, err := db.FindBySlug(clean)
	if err == nil && len(repos) > 0 {
		return repos[0].AbsolutePath, true
	}

	all, listErr := db.ListRepos()
	if listErr != nil {
		return "", false
	}

	for _, r := range all {
		if strings.EqualFold(r.RepoName, clean) || strings.EqualFold(r.Slug, clean) {
			return r.AbsolutePath, true
		}
	}

	return "", false
}

func resolveRepoOrLocalPath(clean string) (string, error) {
	specialPath, isSpecial := resolveSpecialPathTarget(clean)
	if isSpecial {
		return specialPath, nil
	}

	db, err := openDB()
	if err != nil {
		return checkLocalDirFallback(clean)
	}

	defer db.Close()

	p, hasRepo := queryDBForRepoPath(db, clean)
	if hasRepo {
		return p, nil
	}

	return checkLocalDirFallback(clean)
}

func checkLocalDirFallback(clean string) (string, error) {
	info, err := os.Stat(clean)
	if err != nil || !info.IsDir() {
		return "", apperror.NewNotFoundError(fmt.Sprintf("repository or target %q not found", clean))
	}

	abs, absErr := filepath.Abs(clean)
	if absErr != nil {
		return clean, nil
	}

	return abs, nil
}

func outputJSONPathResult(target, resolved string, isValid bool) error {
	res := PathResult{Name: target, Path: resolved, IsValid: isValid}
	data, _ := json.MarshalIndent(res, "", "  ")
	fmt.Println(string(data))
	if !isValid {
		return apperror.NewNotFoundError(fmt.Sprintf("target %q not found", target))
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
