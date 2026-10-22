package cmdplans

import (
	"bufio"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

func isMetaFile(fileName string) bool {
	lower := strings.ToLower(fileName)
	if lower == "01-index.md" || lower == "ledger.md" || lower == "readme.md" {
		return true
	}

	return lower == "what-to-read.md"
}

func scanFileLinesForMetadata(file *os.File, fallbackStatus PlanStatusType) (string, PlanStatusType) {
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

func resolveLineStatus(currentStatus PlanStatusType, line string, hasExplicit bool) (PlanStatusType, bool) {
	if hasExplicit {
		return currentStatus, true
	}
	if s, ok := extractStatusFromLine(line); ok {
		return s, true
	}

	return currentStatus, false
}

func parsePlanFile(filePath, relPath string, fallbackStatus PlanStatusType) (*PlanItem, error) {
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

func buildPlanItemRecord(id int, slug, title string, status PlanStatusType, relPath string, file *os.File) *PlanItem {
	info, _ := file.Stat()
	updatedAt := ""
	if info != nil {
		updatedAt = info.ModTime().Format("2006-01-02")
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

func scanPlanDirectory(dirPath, baseRel string, fallbackStatus PlanStatusType) ([]PlanItem, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, apperror.WrapSimple(err, "os.ReadDir")
	}

	return collectValidPlanEntries(entries, dirPath, baseRel, fallbackStatus), nil
}

func collectValidPlanEntries(entries []os.DirEntry, dirPath, baseRel string, fallbackStatus PlanStatusType) []PlanItem {
	items := make([]PlanItem, 0, len(entries))
	for _, entry := range entries {
		if item, hasItem := processDirEntry(entry, dirPath, baseRel, fallbackStatus); hasItem {
			items = append(items, *item)
		}
	}

	return items
}

func processDirEntry(entry os.DirEntry, dirPath, baseRel string, fallbackStatus PlanStatusType) (*PlanItem, bool) {
	if entry.IsDir() || isMetaFile(entry.Name()) || !strings.HasSuffix(entry.Name(), ".md") {
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

func sortPlanSlices(completed, pending []PlanItem) {
	sort.Slice(completed, func(i, j int) bool {
		if completed[i].ID != completed[j].ID {
			return completed[i].ID > completed[j].ID
		}

		return completed[i].Slug < completed[j].Slug
	})

	sort.Slice(pending, func(i, j int) bool {
		if pending[i].ID != pending[j].ID {
			return pending[i].ID < pending[j].ID
		}

		return pending[i].Slug < pending[j].Slug
	})
}

func partitionDiscoveredItems(items []PlanItem, completed, pending *[]PlanItem) {
	for _, item := range items {
		if item.Status == StatusCompleted {
			*completed = append(*completed, item)
		} else {
			*pending = append(*pending, item)
		}
	}
}

func loadPlanDirectories(repoPath string) ([]PlanItem, []PlanItem, []PlanItem, error) {
	plansDir := filepath.Join(repoPath, ".ai-memory", "plans")
	pendingDir := filepath.Join(plansDir, "pending")
	completedDir := filepath.Join(plansDir, "completed")

	root, err := scanPlanDirectory(plansDir, filepath.Join(".ai-memory", "plans"), StatusCompleted)
	if err != nil {
		return nil, nil, nil, err
	}

	pending, err := scanPlanDirectory(pendingDir, filepath.Join(".ai-memory", "plans", "pending"), StatusPending)
	if err != nil {
		return nil, nil, nil, err
	}

	archived, _ := scanPlanDirectory(completedDir, filepath.Join(".ai-memory", "plans", "completed"), StatusCompleted)

	return root, pending, archived, nil
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
