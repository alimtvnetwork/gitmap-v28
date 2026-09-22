// Package cmdagy — agy_target_resolver.go resolves project targets by sequence, ID, or slug.
package cmdagy

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// ResolveAgyProjectTargets resolves raw args (tokens/comma-separated) to matching AgyProject list.
func ResolveAgyProjectTargets(args []string, projects []AgyProject) ([]AgyProject, error) {
	tokens := expandCommaSeparatedTokens(args)
	if len(tokens) == 0 {
		return nil, apperror.NewSimple("requires target (sequence, id, or slug)", "E9000")
	}
	return resolveTokensToProjects(tokens, projects)
}

func resolveTokensToProjects(tokens []string, projects []AgyProject) ([]AgyProject, error) {
	var results []AgyProject
	seen := make(map[string]bool)
	for _, token := range tokens {
		matched, err := matchSingleTarget(token, projects)
		if err != nil {
			return nil, err
		}
		results = appendUniqueProjects(results, matched, seen)
	}
	return results, nil
}

func appendUniqueProjects(list, matched []AgyProject, seen map[string]bool) []AgyProject {
	for _, p := range matched {
		if !seen[p.ID] {
			seen[p.ID] = true
			list = append(list, p)
		}
	}
	return list
}

func expandCommaSeparatedTokens(args []string) []string {
	var tokens []string
	for _, a := range args {
		for _, part := range strings.Split(a, ",") {
			t := strings.TrimSpace(part)
			if t != "" {
				tokens = append(tokens, t)
			}
		}
	}

	return tokens
}

func matchSingleTarget(token string, projects []AgyProject) ([]AgyProject, error) {
	if seqMatch, isSeq := matchBySequence(token, projects); isSeq {
		return []AgyProject{seqMatch}, nil
	}

	if idMatches := matchByID(token, projects); len(idMatches) > 0 {
		return idMatches, nil
	}

	if slugMatches := matchBySlug(token, projects); len(slugMatches) > 0 {
		return slugMatches, nil
	}

	return nil, fmt.Errorf("target %q did not match any project sequence (e.g. 001), ID, or slug prefix", token)
}

func matchBySequence(token string, projects []AgyProject) (AgyProject, bool) {
	trimmed := strings.TrimLeft(token, "0")
	if trimmed == "" {
		return AgyProject{}, false
	}
	seq, err := strconv.Atoi(trimmed)
	if err != nil {
		return AgyProject{}, false
	}

	if seq >= 1 && seq <= len(projects) {
		return projects[seq-1], true
	}

	return AgyProject{}, false
}
