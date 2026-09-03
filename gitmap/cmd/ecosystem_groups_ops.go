// Package cmd — ecosystem_groups_ops.go provides operations on ecosystem groups.
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

func addEcosystemGroup(ecosystem string, name string, desc string, targets []string) *apperror.AppError {
	store, loadErr := loadEcosystemGroupStore(ecosystem)
	if loadErr != nil {
		return loadErr
	}

	existing, exists := store.Groups[name]
	if !exists {
		store.Groups[name] = newEcosystemGroup(name, desc, targets)
		return saveEcosystemGroupStore(store)
	}

	existing.Targets = appendUniqueTargets(existing.Targets, targets)
	existing.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	store.Groups[name] = existing

	return saveEcosystemGroupStore(store)
}

func appendUniqueTargets(existing []string, targets []string) []string {
	result := existing
	for _, t := range targets {
		if !sliceContains(result, t) {
			result = append(result, t)
		}
	}

	return result
}

func removeEcosystemGroupTarget(ecosystem string, name string, target string) *apperror.AppError {
	store, loadErr := loadEcosystemGroupStore(ecosystem)
	if loadErr != nil {
		return loadErr
	}

	group, exists := store.Groups[name]
	if !exists {
		return apperror.NewSimple("group_not_found", fmt.Sprintf("group %s not found", name))
	}

	newTargets := make([]string, 0, len(group.Targets))
	for _, t := range group.Targets {
		if t != target {
			newTargets = append(newTargets, t)
		}
	}

	group.Targets = newTargets
	group.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	store.Groups[name] = group

	return saveEcosystemGroupStore(store)
}

func deleteEcosystemGroup(ecosystem string, name string) *apperror.AppError {
	store, loadErr := loadEcosystemGroupStore(ecosystem)
	if loadErr != nil {
		return loadErr
	}

	delete(store.Groups, name)

	return saveEcosystemGroupStore(store)
}

func exportEcosystemGroup(ecosystem string, name string, outPath string) *apperror.AppError {
	store, loadErr := loadEcosystemGroupStore(ecosystem)
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

func importEcosystemGroup(ecosystem string, inPath string) *apperror.AppError {
	data, readErr := os.ReadFile(inPath)
	if readErr != nil {
		return apperror.WrapSimple(readErr, "read group import")
	}

	var group EcosystemGroup
	if unmarshalErr := json.Unmarshal(data, &group); unmarshalErr != nil {
		return apperror.WrapSimple(unmarshalErr, "unmarshal group import")
	}

	return addEcosystemGroup(ecosystem, group.Name, group.Description, group.Targets)
}

func getSortedGroupNames(groups map[string]EcosystemGroup) []string {
	names := make([]string, 0, len(groups))
	for n := range groups {
		names = append(names, n)
	}

	sort.Strings(names)

	return names
}
