package cmd

import (
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

type createRepoParams struct {
	Name         string
	Slug         string
	LocalDir     string
	Description  string
	IsPublic     bool
	IsSkipRemote bool
	IsJSON       bool
	IsYAML       bool
	IsCommon     bool
	IsCG         bool
	IsCD         bool
	Profile      model.GitProfile
}

func isPathLike(s string) bool {
	if s == "." || s == ".." {
		return true
	}
	hasSlash := filepath.IsAbs(s) || len(s) > 1 && (s[0] == '.' || s[0] == '/' || s[0] == '\\')
	if hasSlash {
		return true
	}

	return filepath.Dir(s) != "."
}

func isExistingDirectory(path string) bool {
	info, err := os.Stat(path)

	return err == nil && info.IsDir()
}

func applyPositionalArgs(pos []string, p *createRepoParams) {
	if len(pos) == 0 {
		applyCurrentDirPositional(p)
		return
	}

	if len(pos) == 1 {
		applySinglePositionalArg(pos[0], p)
		return
	}

	applyMultiPositionalArgs(pos, p)
}

func applyCurrentDirPositional(p *createRepoParams) {
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "."
	}
	p.LocalDir = cwd
	p.Name = filepath.Base(cwd)
	p.Slug = SlugifyRepoName(p.Name)
}

func applySinglePositionalArg(token string, p *createRepoParams) {
	if token != "." && !isExistingDirectory(token) && !isPathLike(token) {
		p.Name = token
		p.Slug = SlugifyRepoName(token)
		p.LocalDir = filepath.Join(".", p.Slug)
		return
	}

	absDir, err := filepath.Abs(token)
	p.LocalDir = absDir
	p.Name = filepath.Base(absDir)
	if err != nil {
		p.LocalDir = token
		p.Name = filepath.Base(token)
	}
	p.Slug = SlugifyRepoName(p.Name)
}

func applyTwoPositionalArgs(pos []string, p *createRepoParams) {
	if pos[0] == "." || isExistingDirectory(pos[0]) || isPathLike(pos[0]) {
		absDir, _ := filepath.Abs(pos[0])
		p.LocalDir = absDir
		p.Name = pos[1]
		p.Slug = SlugifyRepoName(pos[1])
		return
	}

	if pos[1] == "." || isExistingDirectory(pos[1]) || isPathLike(pos[1]) {
		absDir, _ := filepath.Abs(pos[1])
		p.LocalDir = absDir
		p.Name = pos[0]
		p.Slug = SlugifyRepoName(pos[0])
		return
	}

	p.Name = pos[0]
	p.Slug = SlugifyRepoName(pos[0])
	p.Description = pos[1]
}

func applyMultiPositionalArgs(pos []string, p *createRepoParams) {
	p.Name = pos[0]
	hasTwo := len(pos) == 2
	if hasTwo {
		applyTwoPositionalArgs(pos, p)

		return
	}

	p.LocalDir = pos[1]
	p.Slug = SlugifyRepoName(pos[2])
}

func parseCreateParams(args []string, defaultLocal bool) (createRepoParams, error) {
	p := createRepoParams{
		IsPublic:     hasArgFlag(args, "--public") && !hasArgFlag(args, "--private"),
		IsSkipRemote: defaultLocal || hasArgFlag(args, "--local") || hasArgFlag(args, "--no-remote"),
		IsJSON:       hasArgFlag(args, "--json") || hasArgFlag(args, "-json"),
		IsYAML:       hasArgFlag(args, "--yaml") || hasArgFlag(args, "--yml") || hasArgFlag(args, "-yaml") || hasArgFlag(args, "-y"),
		IsCommon:     hasArgFlag(args, "--common"),
		IsCG:         hasArgFlag(args, "--cg") || hasArgFlag(args, "--coding-guideline"),
		IsCD:         hasArgFlag(args, "--cd"),
		Description:  extractFlagVal(args, "--description"),
	}
	if p.Description == "" {
		p.Description = extractFlagVal(args, "-d")
	}

	applyPositionalArgs(extractNonFlagTokens(args), &p)
	if dir := extractFlagVal(args, "--dir"); dir != "" {
		p.LocalDir = dir
	}
	if slug := extractFlagVal(args, "--slug"); slug != "" {
		p.Slug = SlugifyRepoName(slug)
	}

	prof, profErr := resolveCreationProfile(args)
	p.Profile = prof

	return p, profErr
}
