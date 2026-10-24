package cmdplans

import (
	"bufio"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdplans/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdplans/enums"
)

func isMetaFile(fileName string) bool {
	return enums.IsMetaFileName(fileName)
}

func scanFileLinesForMetadata(file *os.File, fallbackStatus enums.PlanStatusType) (string, enums.PlanStatusType) {
	scanner := bufio.NewScanner(file)
	title := ""
	status := fallbackStatus
	hasExplicitStatus := false

	for count := 0; scanner.Scan() && count < 35; count++ {
		line := scanner.Text()
		title = resolveLineTitle(title, line)
		status, hasExplicitStatus = resolveLineStatus(status, line, hasExplicitStatus)
	}

	return title, status
}

func resolveLineTitle(currentTitle, line string) string {
	if len(currentTitle) > 0 {
		return currentTitle
	}
	if t, ok := extractTitleFromLine(line); ok {
		return t
	}

	return currentTitle
}

func resolveLineStatus(currentStatus enums.PlanStatusType, line string, hasExplicit bool) (enums.PlanStatusType, bool) {
	if hasExplicit {
		return currentStatus, true
	}
	if s, ok := extractStatusFromLine(line); ok {
		return s, true
	}

	return currentStatus, false
}

func parsePlanFile(filePath, relPath string, fallbackStatus enums.PlanStatusType) (*PlanItem, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, apperror.WrapSimple(err, "os.Open plan file")
	}
	defer file.Close()

	id, slug := parseIDAndSlug(filepath.Base(filePath))
	title, status := scanFileLinesForMetadata(file, fallbackStatus)
	if len(title) == 0 {
		title = formatFallbackTitle(slug)
	}

	return buildPlanItemRecord(id, slug, title, status, relPath, file), nil
}

func buildPlanItemRecord(id int, slug, title string, status enums.PlanStatusType, relPath string, file *os.File) *PlanItem {
	info, _ := file.Stat()
	updatedAt := ""
	if info != nil {
		updatedAt = info.ModTime().Format(constants.DefaultDateLayout)
	}

	return &PlanItem{
		ID:        id,
		Slug:      slug,
		Title:     title,
		Status:    status,
		FilePath:  filepath.ToSlash(relPath),
		UpdatedAt: updatedAt,
	}
}

func scanPlanDirectory(dirPath, baseRel string, fallbackStatus enums.PlanStatusType) ([]PlanItem, error) {
	entries, err := os.ReadDir(dirPath)
	if os.IsNotExist(err) {
		return nil, nil
	}

	if err != nil {
		return nil, apperror.WrapSimple(err, "os.ReadDir")
	}

	return collectValidPlanEntries(entries, dirPath, baseRel, fallbackStatus), nil
}

func collectValidPlanEntries(entries []os.DirEntry, dirPath, baseRel string, fallbackStatus enums.PlanStatusType) []PlanItem {
	items := make([]PlanItem, 0, len(entries))
	for _, entry := range entries {
		if item, hasItem := processDirEntry(entry, dirPath, baseRel, fallbackStatus); hasItem {
			items = append(items, *item)
		}
	}

	return items
}

func processDirEntry(entry os.DirEntry, dirPath, baseRel string, fallbackStatus enums.PlanStatusType) (*PlanItem, bool) {
	if entry.IsDir() || isMetaFile(entry.Name()) || !strings.HasSuffix(entry.Name(), constants.ExtMarkdown) {
		return nil, false
	}
	fullPath := filepath.Join(dirPath, entry.Name())
	relPath := filepath.Join(baseRel, entry.Name())
	item, parseErr := parsePlanFile(fullPath, relPath, fallbackStatus)
	if parseErr != nil {
		return nil, false
	}

	return item, true
}

func sortCompletedPlans(items []PlanItem) {
	sort.Slice(items, func(i, j int) bool {
		if items[i].ID != items[j].ID {
			return items[i].ID > items[j].ID
		}

		return items[i].Slug < items[j].Slug
	})
}

func sortPendingPlans(items []PlanItem) {
	sort.Slice(items, func(i, j int) bool {
		if items[i].ID != items[j].ID {
			return items[i].ID < items[j].ID
		}

		return items[i].Slug < items[j].Slug
	})
}

func sortPlanSlices(completed, pending []PlanItem) {
	sortCompletedPlans(completed)
	sortPendingPlans(pending)
}

func partitionDiscoveredItems(items []PlanItem, completed, pending *[]PlanItem) {
	for _, item := range items {
		if item.Status.IsCompleted() {
			*completed = append(*completed, item)
		} else {
			*pending = append(*pending, item)
		}
	}
}

func scanPendingAndArchived(plansDir, rootRel string) ([]PlanItem, []PlanItem, error) {
	pendingDir := filepath.Join(plansDir, constants.DirPending)
	pendingRel := filepath.Join(rootRel, constants.DirPending)
	pending, err := scanPlanDirectory(pendingDir, pendingRel, enums.PlanStatusPending)
	if err != nil {
		return nil, nil, err
	}
	completedDir := filepath.Join(plansDir, constants.DirCompleted)
	completedRel := filepath.Join(rootRel, constants.DirCompleted)
	archived, _ := scanPlanDirectory(completedDir, completedRel, enums.PlanStatusCompleted)

	return pending, archived, nil
}

func loadPlanDirectories(repoPath string) ([]PlanItem, []PlanItem, []PlanItem, error) {
	plansDir := filepath.Join(repoPath, constants.DirAiMemory, constants.DirPlans)
	rootRel := filepath.Join(constants.DirAiMemory, constants.DirPlans)
	root, err := scanPlanDirectory(plansDir, rootRel, enums.PlanStatusCompleted)
	if err != nil {
		return nil, nil, nil, err
	}
	pending, archived, err := scanPendingAndArchived(plansDir, rootRel)

	return root, pending, archived, err
}

// DiscoverPlans scans the plans folders in repoPath and segregates items into completed and pending.
func DiscoverPlans(repoPath string) ([]PlanItem, []PlanItem, error) {
	rootItems, pendingItems, archivedItems, err := loadPlanDirectories(repoPath)
	if err != nil {
		return nil, nil, err
	}

	allCompleted := make([]PlanItem, 0, len(rootItems)+len(archivedItems))
	allPending := make([]PlanItem, 0, len(pendingItems))

	partitionDiscoveredItems(rootItems, &allCompleted, &allPending)
	partitionDiscoveredItems(pendingItems, &allCompleted, &allPending)
	allCompleted = append(allCompleted, archivedItems...)
	sortPlanSlices(allCompleted, allPending)

	return allCompleted, allPending, nil
}
