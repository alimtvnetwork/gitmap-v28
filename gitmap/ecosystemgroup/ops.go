// Package ecosystemgroup provides operations on ecosystem groups.
package ecosystemgroup

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// AddGroup adds or updates a group within the specified ecosystem.
func AddGroup(ecosystem string, name string, desc string, targets []string) *apperror.AppError {
	store, loadErr := LoadStore(ecosystem)
	if loadErr != nil {
		return loadErr
	}

	existing, exists := store.Groups[name]
	if !exists {
		store.Groups[name] = NewGroup(name, desc, targets)
		return SaveStore(store)
	}

	existing.Targets = AppendUniqueTargets(existing.Targets, targets)
	existing.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	store.Groups[name] = existing

	return SaveStore(store)
}

// AppendUniqueTargets appends target items if not already present.
func AppendUniqueTargets(existing []string, targets []string) []string {
	result := existing
	for _, t := range targets {
		if !SliceContains(result, t) {
			result = append(result, t)
		}
	}
	return result
}

// RemoveTarget removes a single target from an ecosystem group.
func RemoveTarget(ecosystem string, name string, target string) *apperror.AppError {
	store, loadErr := LoadStore(ecosystem)
	if loadErr != nil {
		return loadErr
	}

	group, exists := store.Groups[name]
	if !exists {
		return apperror.NewSimple("group_not_found", fmt.Sprintf("group %s not found", name))
	}

	group.Targets = filterTargets(group.Targets, target)
	group.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	store.Groups[name] = group

	return SaveStore(store)
}

func filterTargets(targets []string, exclude string) []string {
	filtered := make([]string, 0, len(targets))
	for _, t := range targets {
		if t != exclude {
			filtered = append(filtered, t)
		}
	}
	return filtered
}

// DeleteGroup deletes a group from the specified ecosystem.
func DeleteGroup(ecosystem string, name string) *apperror.AppError {
	store, loadErr := LoadStore(ecosystem)
	if loadErr != nil {
		return loadErr
	}

	delete(store.Groups, name)

	return SaveStore(store)
}

// ExportGroup exports a group as JSON to outPath.
func ExportGroup(ecosystem string, name string, outPath string) *apperror.AppError {
	store, loadErr := LoadStore(ecosystem)
	if loadErr != nil {
		return loadErr
	}

	group, exists := store.Groups[name]
	if !exists {
		return apperror.NewSimple("group_not_found", fmt.Sprintf("group %s not found", name))
	}

	data, marshalErr := json.MarshalIndent(group, "", "  ")
	if marshalErr != nil {
		return apperror.WrapSimple(marshalErr, "marshal group")
	}

	if writeErr := os.WriteFile(outPath, data, constants.FilePermission); writeErr != nil {
		return apperror.WrapSimple(writeErr, "write group export")
	}

	return nil
}

// ImportGroup imports a group definition from JSON inPath.
func ImportGroup(ecosystem string, inPath string) *apperror.AppError {
	data, readErr := os.ReadFile(inPath)
	if readErr != nil {
		return apperror.WrapSimple(readErr, "read group import")
	}

	var group EcosystemGroup
	if unmarshalErr := json.Unmarshal(data, &group); unmarshalErr != nil {
		return apperror.WrapSimple(unmarshalErr, "unmarshal group import")
	}

	return AddGroup(ecosystem, group.Name, group.Description, group.Targets)
}

// GetSortedGroupNames returns sorted keys from groups map.
func GetSortedGroupNames(groups map[string]EcosystemGroup) []string {
	names := make([]string, 0, len(groups))
	for n := range groups {
		names = append(names, n)
	}

	sort.Strings(names)

	return names
}

// PrintGroupList renders the formatted list of groups to stdout.
func PrintGroupList(ecosystem string, label string) error {
	store, err := LoadStore(ecosystem)
	if err != nil {
		return fmt.Errorf("load groups: %w", err.Unwrap())
	}

	if len(store.Groups) == 0 {
		fmt.Printf("%s No %s groups configured.\n", constants.ColorYellow+"ℹ"+constants.ColorReset, label)
		return nil
	}

	names := GetSortedGroupNames(store.Groups)
	fmt.Printf("%s %s Groups (%d):\n", constants.ColorCyan+"▶"+constants.ColorReset, label, len(names))
	for _, name := range names {
		g := store.Groups[name]
		fmt.Printf("  • \033[1m%-16s\033[0m %d target(s) (Updated: %s)\n", g.Name, len(g.Targets), g.UpdatedAt)
	}

	return nil
}
