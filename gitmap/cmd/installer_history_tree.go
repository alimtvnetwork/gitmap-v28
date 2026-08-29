// Package cmd — installer_history_tree.go renders installer history as a tree hierarchy.
package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

// getInstallerTimestamp extracts the most recent timestamp from an installer record.
func getInstallerTimestamp(scriptRecord model.InstallerScript) string {
	if scriptRecord.UpdatedAt != "" {
		return scriptRecord.UpdatedAt
	}
	return scriptRecord.CreatedAt
}

// isScriptNewer compares two installer script records by timestamp and ID.
func isScriptNewer(candidateRecord, currentRecord model.InstallerScript) bool {
	candidateTimestamp := getInstallerTimestamp(candidateRecord)
	currentTimestamp := getInstallerTimestamp(currentRecord)
	if candidateTimestamp != currentTimestamp {
		return candidateTimestamp > currentTimestamp
	}
	return candidateRecord.ID > currentRecord.ID
}

// sortGroupedInstallers sorts grouped installer records descending by timestamp.
func sortGroupedInstallers(latestBySlug map[string]model.InstallerScript, slugKeys []string) []model.InstallerScript {
	sortedScripts := make([]model.InstallerScript, 0, len(slugKeys))
	for _, slugKey := range slugKeys {
		sortedScripts = append(sortedScripts, latestBySlug[slugKey])
	}
	sort.Slice(sortedScripts, func(firstIndex, secondIndex int) bool {
		return isScriptNewer(sortedScripts[firstIndex], sortedScripts[secondIndex])
	})
	return sortedScripts
}

// updateLatestSlugRecord updates the latest record map if candidate is newer.
func updateLatestSlugRecord(
	latestMap map[string]model.InstallerScript,
	slugKeys *[]string,
	scriptRecord model.InstallerScript,
) {
	slugKey := strings.ToLower(strings.TrimSpace(scriptRecord.Slug))
	existingRecord, hasExisting := latestMap[slugKey]
	if !hasExisting {
		*slugKeys = append(*slugKeys, slugKey)
		latestMap[slugKey] = scriptRecord
		return
	}
	if isScriptNewer(scriptRecord, existingRecord) {
		latestMap[slugKey] = scriptRecord
	}
}

// groupLatestInstallers groups records by slug, preserving only the latest entry.
func groupLatestInstallers(scriptList []model.InstallerScript) []model.InstallerScript {
	latestBySlug := make(map[string]model.InstallerScript)
	var slugKeys []string
	for _, scriptRecord := range scriptList {
		updateLatestSlugRecord(latestBySlug, &slugKeys, scriptRecord)
	}
	return sortGroupedInstallers(latestBySlug, slugKeys)
}

// formatHistoryMetadata constructs the metadata string for a history record.
func formatHistoryMetadata(scriptRecord model.InstallerScript) string {
	timestampText := getInstallerTimestamp(scriptRecord)
	if timestampText != "" {
		return fmt.Sprintf("(%s, %s)", scriptRecord.TargetOS, timestampText)
	}
	return fmt.Sprintf("(%s)", scriptRecord.TargetOS)
}

// getHistoryDescription retrieves the script description or a default fallback.
func getHistoryDescription(scriptRecord model.InstallerScript) string {
	if scriptRecord.Description != "" {
		return scriptRecord.Description
	}
	return "Installed component"
}

// buildSingleHistoryNode constructs a single child node for an installer script.
func buildSingleHistoryNode(scriptRecord model.InstallerScript) InstallerTreeNode {
	return InstallerTreeNode{
		Title:       scriptRecord.Slug,
		Description: formatHistoryMetadata(scriptRecord),
		Children: []InstallerTreeNode{
			{Title: getHistoryDescription(scriptRecord), Description: scriptRecord.Version},
		},
	}
}

// renderSingleHistoryTree formats and outputs a non-profile installer tree node.
func renderSingleHistoryTree(scriptRecord model.InstallerScript) {
	rootNode := buildSingleHistoryNode(scriptRecord)
	printInstallerTree(rootNode, "", true)
}

// renderHistoryEntry routes rendering to profile tree or single history node.
func renderHistoryEntry(scriptRecord model.InstallerScript) {
	profileRecord, hasProfile := resolveProfileTree(scriptRecord.Slug)
	if hasProfile {
		printProfileTree(profileRecord)
		return
	}
	renderSingleHistoryTree(scriptRecord)
}

// printHistoryDivider outputs the DIM divider line between history entries.
func printHistoryDivider() {
	fmt.Printf("  %s------%s\n", constants.ColorDim, constants.ColorReset)
}

// renderGroupedHistoryList prints each history entry separated by dividers.
func renderGroupedHistoryList(groupedScripts []model.InstallerScript) {
	for scriptIndex, scriptRecord := range groupedScripts {
		if scriptIndex > 0 {
			printHistoryDivider()
		}
		renderHistoryEntry(scriptRecord)
	}
}

// printInstallerHistoryTree queries and renders installer ledger history as a tree.
func printInstallerHistoryTree(dbInstance *store.DB) {
	if dbInstance == nil {
		return
	}
	scriptList, errList := dbInstance.ListInstallHistory()
	if errList != nil || len(scriptList) == 0 {
		fmt.Printf("  %sNo installer history records found.%s\n", constants.ColorDim, constants.ColorReset)
		return
	}
	renderGroupedHistoryList(groupLatestInstallers(scriptList))
}
