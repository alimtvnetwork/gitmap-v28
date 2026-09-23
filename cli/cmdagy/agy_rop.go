package cmdagy

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"text/tabwriter"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

type projectSortItem struct {
	project AgyProject
	modTime time.Time
}

func parseProjectTimestamp(raw string, fallback time.Time) time.Time {
	if raw == "" {
		return fallback
	}
	t, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return fallback
	}
	return t
}

func collectAndSortProjects(dirPath string) ([]AgyProject, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, apperror.WrapSimple(err, "read projects dir")
	}

	var items []projectSortItem
	for _, e := range entries {
		if filepath.Ext(e.Name()) != ".json" {
			continue
		}
		info, infoErr := e.Info()
		fileMod := time.Time{}
		if infoErr == nil {
			fileMod = info.ModTime()
		}
		proj, readErr := readAgyProject(filepath.Join(dirPath, e.Name()))
		if readErr == nil {
			modTime := parseProjectTimestamp(proj.UpdatedAt, fileMod)
			items = append(items, projectSortItem{project: proj, modTime: modTime})
		}
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].modTime.After(items[j].modTime)
	})

	projects := make([]AgyProject, 0, len(items))
	for _, it := range items {
		projects = append(projects, it.project)
	}
	return projects, nil
}

// FilterActiveProjects filters projects active in the last 24h, up to maxN.
func FilterActiveProjects(dirPath string, maxN int) ([]AgyProject, error) {
	all, err := collectAndSortProjects(dirPath)
	if err != nil {
		return nil, err
	}

	oneDayAgo := time.Now().Add(-24 * time.Hour)
	var active []AgyProject
	for _, p := range all {
		t := parseProjectTimestamp(p.UpdatedAt, time.Time{})
		if t.After(oneDayAgo) {
			active = append(active, p)
		}
	}

	if len(active) == 0 {
		active = all
	}

	if maxN > 0 && len(active) > maxN {
		active = active[:maxN]
	}
	return active, nil
}

func backupAndRemoveOldConvs(slug string, p AgyProject, isDryRun bool) (int, error) {
	matching, err := FindMatchingConversations(p.GetPath())
	if err != nil || len(matching) == 0 {
		return 0, nil
	}

	if isDryRun {
		return len(matching), nil
	}

	db, err := OpenAGYBackupDB(slug)
	if err != nil {
		return 0, err
	}
	defer db.Close()

	backedUp := 0
	home, _ := os.UserHomeDir()
	convDir := filepath.Join(home, ".gemini", "antigravity", "conversations")

	for _, c := range matching {
		transcript := ReadTranscriptForConv(c.ID)
		backup := AGYConversationBackup{
			ConversationID: c.ID,
			ProjectSlug:    slug,
			ProjectPath:    p.GetPath(),
			Title:          p.Name,
			StepCount:      c.StepCount,
			TranscriptJSON: transcript,
			CreatedAt:      time.Now().UTC(),
		}
		if saveErr := BackupConversation(db, backup); saveErr == nil {
			backedUp++
			_ = os.Remove(filepath.Join(convDir, c.ID+".db"))
		}
	}
	return backedUp, nil
}

func triggerNewProjectConversation(p AgyProject, isDryRun bool) string {
	if isDryRun {
		return "preview-conv-id"
	}
	prompt := fmt.Sprintf("Please inspect the codebase structure and key modules for %s (%s) to refresh your active working memory.", p.Name, p.GetPath())
	res := AgentAPINewConversation(p.Name, prompt)
	if res.IsSuccess() {
		return res.Value
	}
	return ""
}

// ExecuteROP orchestrates the reread-optimize-project workflow.
func ExecuteROP(ctx context.Context, out io.Writer, maxN int, isDryRun bool) (*AGYROPResult, *apperror.AppError) {
	if maxN <= 0 {
		maxN = 5
	}
	dirPath, err := getProjectsDirPath()
	if err != nil {
		return nil, apperror.WrapSimple(err, "get projects dir")
	}

	active, err := FilterActiveProjects(dirPath, maxN)
	if err != nil {
		return nil, apperror.WrapSimple(err, "filter active projects")
	}

	cleanOpts := buildCleanCacheOpts(10, isDryRun, false)
	_ = ExecuteCleanCache(cleanOpts)

	var projectResults []AGYROPProjectResult
	totalBacked := 0

	for _, p := range active {
		slug := ResolveRepoSlug(p.Name, p.GetPath())
		backed, bErr := backupAndRemoveOldConvs(slug, p, isDryRun)
		totalBacked += backed

		newID := triggerNewProjectConversation(p, isDryRun)
		status := "OPTIMIZED"
		errMsg := ""
		if bErr != nil {
			status = "WARNING"
			errMsg = bErr.Error()
		}

		projectResults = append(projectResults, AGYROPProjectResult{
			ProjectName:       p.Name,
			ProjectSlug:       slug,
			ProjectPath:       p.GetPath(),
			ConvsBackedUp:     backed,
			NewConversationID: newID,
			Status:            status,
			ErrorMsg:          errMsg,
		})
	}

	res := &AGYROPResult{
		RequestedN:    maxN,
		ActiveFound:   len(active),
		CacheCleared:  true,
		TotalBackedUp: totalBacked,
		Projects:      projectResults,
	}

	renderROPTable(out, res)
	return res, nil
}

func renderROPTable(out io.Writer, res *AGYROPResult) {
	w := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "STATUS\tPROJECT\tSLUG\tBACKED_UP\tNEW_CONV\tPATH")
	for _, pr := range res.Projects {
		statusStr := constants.ColorGreen + pr.Status + constants.ColorReset
		if pr.Status != "OPTIMIZED" {
			statusStr = constants.ColorYellow + pr.Status + constants.ColorReset
		}
		newConvDisplay := pr.NewConversationID
		if newConvDisplay == "" {
			newConvDisplay = "-"
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%s\t%s\n",
			statusStr, pr.ProjectName, pr.ProjectSlug, pr.ConvsBackedUp, newConvDisplay, pr.ProjectPath)
	}
	_ = w.Flush()
	fmt.Fprintf(out, "\n  ROP summary: Optimized %d projects (Total conversations backed up: %d, retention: --keep 10)\n\n",
		len(res.Projects), res.TotalBackedUp)
}
