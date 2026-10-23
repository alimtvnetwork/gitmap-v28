package cmdexclude

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// ToggleResult represents the outcome of a repository exclusion toggle operation.
type ToggleResult struct {
	ModifiedCount int      `json:"modifiedCount"`
	AffectedRepos []string `json:"affectedRepos"`
	TotalExcluded int      `json:"totalExcluded"`
	CurrentTotal  int      `json:"currentTotal"`
	IsSuccess     bool     `json:"isSuccess"`
}

// ToggleRepoExclusions adds or removes repository paths from custom exclusions.
func ToggleRepoExclusions(selection string, isDisable bool) (*ToggleResult, error) {
	repos, err := loadTrackedRepos()
	if err != nil {
		return nil, err
	}

	targets, err := resolveTargets(repos, selection)
	if err != nil {
		return nil, err
	}

	existing, _ := LoadCustomExclusions()

	return executeExclusionSave(existing, targets, isDisable)
}

func executeExclusionSave(existing, targets []string, isDisable bool) (*ToggleResult, error) {
	updated, affected := computeExclusionDelta(existing, targets, isDisable)
	saveErr := persistUpdatedExclusions(updated, len(affected) > 0)
	if saveErr != nil {
		return nil, saveErr
	}

	total := len(updated)

	return &ToggleResult{
		ModifiedCount: len(affected),
		AffectedRepos: affected,
		TotalExcluded: total,
		CurrentTotal:  total,
		IsSuccess:     true,
	}, nil
}

func persistUpdatedExclusions(updated []string, hasChanges bool) error {
	if !hasChanges {
		return nil
	}

	err := SaveCustomExclusions(updated)
	if err != nil {
		return apperror.WrapSimple(err, "save_custom_exclusions")
	}

	return nil
}

func computeExclusionDelta(existing, targets []string, isDisable bool) ([]string, []string) {
	if isDisable {
		return applyDisableExclusions(existing, targets)
	}

	return applyEnableExclusions(existing, targets)
}

func loadTrackedRepos() ([]model.ScanRecord, error) {
	db, err := store.OpenDefault()
	if err != nil {
		return nil, apperror.WrapSimple(err, "open_default_store")
	}

	defer db.Close()

	repos, listErr := db.ListRepos()
	if listErr != nil {
		return nil, apperror.WrapSimple(listErr, "list_repos")
	}

	return repos, nil
}

func resolveTargets(repos []model.ScanRecord, selection string) ([]string, error) {
	clean := strings.TrimSpace(selection)
	if clean == "" || strings.EqualFold(clean, "none") {
		return nil, nil
	}
	if strings.EqualFold(clean, "all") {
		return resolveAllTargets(repos), nil
	}

	return parseMixedSelectionTargets(repos, clean)
}

func resolveAllTargets(repos []model.ScanRecord) []string {
	targets := make([]string, 0, len(repos))
	for _, r := range repos {
		targets = append(targets, repoIdentifier(r))
	}

	return targets
}

func repoIdentifier(r model.ScanRecord) string {
	if r.Slug != "" {
		return r.Slug
	}
	if r.RepoName != "" {
		return r.RepoName
	}

	return r.RelativePath
}

func parseMixedSelectionTargets(repos []model.ScanRecord, clean string) ([]string, error) {
	tokens := strings.Split(clean, ",")
	seen := make(map[string]bool)
	var targets []string

	for _, token := range tokens {
		t := strings.TrimSpace(token)
		if t != "" {
			appendTokenTargets(repos, t, seen, &targets)
		}
	}

	return targets, nil
}

func appendTokenTargets(repos []model.ScanRecord, t string, seen map[string]bool, targets *[]string) {
	items := matchTokenToTargets(repos, t)
	for _, item := range items {
		low := strings.ToLower(item)
		if !seen[low] {
			seen[low] = true
			*targets = append(*targets, item)
		}
	}
}

func matchTokenToTargets(repos []model.ScanRecord, t string) []string {
	indices, err := ParseNumberedSelection(t, len(repos))
	hasValidIndices := err == nil && len(indices) > 0
	if hasValidIndices {
		return indicesToTargets(repos, indices)
	}

	return matchSlugTargets(repos, t)
}

func indicesToTargets(repos []model.ScanRecord, indices []int) []string {
	var targets []string
	for _, idx := range indices {
		if idx >= 1 && idx <= len(repos) {
			targets = append(targets, repoIdentifier(repos[idx-1]))
		}
	}

	return targets
}

func matchSlugTargets(repos []model.ScanRecord, token string) []string {
	var matches []string
	for _, r := range repos {
		if isRepoMatch(r, token) {
			matches = append(matches, repoIdentifier(r))
		}
	}
	hasMatches := len(matches) > 0
	if hasMatches {
		return matches
	}

	return []string{token}
}

func isRepoMatch(r model.ScanRecord, token string) bool {
	if strings.EqualFold(r.Slug, token) || strings.EqualFold(r.RepoName, token) {
		return true
	}
	if strings.EqualFold(r.RelativePath, token) || strings.EqualFold(r.PrimaryAlias, token) {
		return true
	}

	return false
}

func applyDisableExclusions(existing []string, targets []string) ([]string, []string) {
	set := make(map[string]bool)
	for _, p := range existing {
		set[strings.ToLower(p)] = true
	}

	var affected []string
	updated := append([]string{}, existing...)

	for _, target := range targets {
		low := strings.ToLower(target)
		if !set[low] {
			set[low] = true
			updated = append(updated, target)
			affected = append(affected, target)
		}
	}

	return updated, affected
}

func applyEnableExclusions(existing []string, targets []string) ([]string, []string) {
	toRemove := make(map[string]bool)
	for _, t := range targets {
		toRemove[strings.ToLower(t)] = true
	}

	var updated []string
	var affected []string

	for _, p := range existing {
		low := strings.ToLower(p)
		if toRemove[low] {
			affected = append(affected, p)
		} else {
			updated = append(updated, p)
		}
	}

	return updated, affected
}
