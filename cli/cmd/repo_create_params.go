package cmd

import (
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
	Profile      model.GitProfile
}

func isPathLike(s string) bool {
	hasSlash := filepath.IsAbs(s) || len(s) > 1 && (s[0] == '.' || s[0] == '/' || s[0] == '\\')
	if hasSlash {
		return true
	}

	return filepath.Dir(s) != "."
}

func applyPositionalArgs(pos []string, p *createRepoParams) {
	hasOne := len(pos) == 1
	if hasOne {
		p.Name = pos[0]
		p.Slug = SlugifyRepoName(pos[0])
		p.LocalDir = filepath.Join(".", p.Slug)

		return
	}

	applyMultiPositionalArgs(pos, p)
}

func applyTwoPositionalArgs(pos []string, p *createRepoParams) {
	isPath := isPathLike(pos[1])
	if isPath {
		p.LocalDir = pos[1]
		p.Slug = SlugifyRepoName(pos[0])

		return
	}

	p.Slug = SlugifyRepoName(pos[1])
	p.LocalDir = filepath.Join(".", p.Slug)
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
