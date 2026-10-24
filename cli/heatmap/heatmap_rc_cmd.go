package heatmap

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func isRcOverviewRequested(args []string) bool {
	if len(args) == 0 {
		return true
	}

	first := strings.ToLower(strings.TrimSpace(args[0]))
	if first == "help" || first == "--help" || first == "-h" || first == "overview" {
		return true
	}

	return strings.HasPrefix(first, "-")
}

func dispatchRcAction(subcmd string, flags map[string]string, pos []string, out io.Writer) (bool, error) {
	switch subcmd {
	case "log":
		return true, runRcLog(flags, out)
	case "update", "outcome":
		return true, runRcUpdate(flags, pos, out)
	case "search":
		return true, runRcSearch(flags, pos, out)
	case "tree":
		return true, runRcTree(flags, out)
	case "status":
		return true, runRcStatus(flags, out)
	}

	return false, nil
}

func dispatchRcAdmin(ctx context.Context, subcmd string, flags map[string]string, pos []string, out io.Writer) error {
	switch subcmd {
	case "list-ai", "list_ai", "list":
		return runRcListAi(flags, out)
	case "commit":
		return runRcCommit(flags, out)
	case "sync":
		return runRcSync(flags, out)
	case "tables", "table":
		return RunHeatmapTables(true, flags, out)
	case "view":
		return RunHeatmapView(true, flags, pos, out)
	case "query", "sql":
		return ExecuteHeatmapQuery(true, flags, pos, out)
	case "overview":
		return DispatchHeatmapOverview(ctx, out, true)
	default:
		return apperror.NewValidation("heatmap_rc", "E1000", fmt.Sprintf("unknown rc subcommand: %s", subcmd))
	}
}

// DispatchRcHeatmap routes gitmap heatmap rc subcommands.
func DispatchRcHeatmap(ctx context.Context, args []string, out io.Writer) error {
	if isRcOverviewRequested(args) {
		return DispatchHeatmapOverview(ctx, out, true)
	}

	subcmd := strings.ToLower(args[0])
	flags, pos := parseHeatmapArgs(args[1:])
	isHandled, err := dispatchRcAction(subcmd, flags, pos, out)
	if isHandled {
		return err
	}

	return dispatchRcAdmin(ctx, subcmd, flags, pos, out)
}

func runRcLog(flags map[string]string, out io.Writer) error {
	aiName, repoSlug := resolveTargetAi(flags["ai"]), resolveTargetRepo(flags["repo"])
	st, err := store.OpenRcHeatmapStore(aiName, repoSlug)
	if err != nil {
		return err
	}
	defer st.Close()

	rec := buildActionLogRecord(flags)
	rec.RepoName = repoSlug
	id, err := st.LogAction(rec)
	if err != nil {
		return err
	}
	rec.HeatmapActionLogId = id

	return renderRcHeatmapLogResponse(out, rec, st.RepoDbPath(), aiName, repoSlug, flags)
}

func runRcSearch(flags map[string]string, pos []string, out io.Writer) error {
	aiName, repoSlug := resolveTargetAi(flags["ai"]), resolveTargetRepo(flags["repo"])
	st, err := store.OpenRcHeatmapStore(aiName, repoSlug)
	if err != nil {
		return err
	}
	defer st.Close()

	limit := 50
	if l, err := strconv.Atoi(flags["limit"]); err == nil && l > 0 {
		limit = l
	}

	records, err := st.SearchActions(strings.Join(pos, " "), flags["task"], flags["file"], limit)
	if err != nil {
		return err
	}

	if flags["json"] == "true" || isJsonRequested(flags) {
		return writeIndentedJSON(out, records)
	}
	renderSearchText(out, records)
	return nil
}

func runRcTree(flags map[string]string, out io.Writer) error {
	aiName, repoSlug := resolveTargetAi(flags["ai"]), resolveTargetRepo(flags["repo"])
	st, err := store.OpenRcHeatmapStore(aiName, repoSlug)
	if err != nil {
		return err
	}
	defer st.Close()

	tree, err := st.GetExecutionTree(flags["task"])
	if err != nil {
		return err
	}

	if flags["json"] == "true" || isJsonRequested(flags) {
		return writeIndentedJSON(out, tree)
	}
	renderRcTreeText(out, tree, "")
	return nil
}

func runRcStatus(flags map[string]string, out io.Writer) error {
	aiName, repoSlug := resolveTargetAi(flags["ai"]), resolveTargetRepo(flags["repo"])
	st, err := store.OpenRcHeatmapStore(aiName, repoSlug)
	if err != nil {
		return err
	}
	defer st.Close()

	summary, err := st.GetStatusSummary(flags["task"])
	if err != nil {
		return err
	}

	uncommitted, _ := st.GetUncommittedActionCount(repoSlug, aiName)
	if flags["json"] == "true" || isJsonRequested(flags) {
		return writeIndentedJSON(out, map[string]any{
			"aiName": aiName, "repoSlug": repoSlug, "uncommitted": uncommitted, "summary": summary,
		})
	}

	fmt.Fprintf(out, "RC Heatmap Status [AI: %s, Repo: %s]\nUncommitted Actions: %d\n", aiName, repoSlug, uncommitted)
	renderStatusText(out, *summary)
	return nil
}

func runRcListAi(flags map[string]string, out io.Writer) error {
	names := store.ListSupportedAINames()
	current := store.DetectCurrentAIName()

	if flags["json"] == "true" || isJsonRequested(flags) {
		return writeIndentedJSON(out, map[string]any{
			"currentAI": current, "supportedAIs": names,
		})
	}

	fmt.Fprintln(out, "Supported AI Engine Namespaces in repo-cache (.heatmap/):")
	for _, n := range names {
		prefix := store.AISequencePrefix(n)
		activeMark := " "
		if n == current {
			activeMark = "*"
		}
		fmt.Fprintf(out, " %s %s-%-15s (slug: %s)\n", activeMark, prefix, n, n)
	}
	fmt.Fprintf(out, "\nAmbient Detected AI: %s\n", current)
	return nil
}

func runRcCommit(flags map[string]string, out io.Writer) error {
	aiName, repoSlug := resolveTargetAi(flags["ai"]), resolveTargetRepo(flags["repo"])
	taskSlug := flags["task"]
	if taskSlug == "" {
		taskSlug = "default-task"
	}
	msg := flags["msg"]
	if msg == "" {
		msg = fmt.Sprintf("Heatmap - %s - %s", repoSlug, taskSlug)
	}

	st, err := store.OpenRcHeatmapStore(aiName, repoSlug)
	if err != nil {
		return err
	}
	defer st.Close()

	rec, err := st.GetOrCreateLocalTracking(repoSlug, aiName)
	if err != nil {
		return err
	}

	rec.SyncStatus, rec.IsDirty, rec.UncommittedChangeCount = "COMMITTED", false, 0
	_ = st.UpdateLocalTracking(*rec)

	if flags["json"] == "true" || isJsonRequested(flags) {
		return writeIndentedJSON(out, map[string]any{
			"status": "COMMITTED", "aiName": aiName, "repoSlug": repoSlug, "message": msg,
		})
	}

	fmt.Fprintf(out, "Committed RC heatmap telemetry [AI: %s, Repo: %s]\nMessage: %s\n", aiName, repoSlug, msg)
	return nil
}

func runRcSync(flags map[string]string, out io.Writer) error {
	aiName, repoSlug := resolveTargetAi(flags["ai"]), resolveTargetRepo(flags["repo"])
	st, err := store.OpenRcHeatmapStore(aiName, repoSlug)
	if err != nil {
		return err
	}
	defer st.Close()

	rec, err := st.GetOrCreateLocalTracking(repoSlug, aiName)
	if err != nil {
		return err
	}

	rec.SyncStatus, rec.IsDirty = "SYNCED", false
	_ = st.UpdateLocalTracking(*rec)

	if flags["json"] == "true" || isJsonRequested(flags) {
		return writeIndentedJSON(out, map[string]any{
			"status": "SYNCED", "aiName": aiName, "repoSlug": repoSlug,
		})
	}

	fmt.Fprintf(out, "Synchronized RC heatmap telemetry with repo-cache [AI: %s, Repo: %s]\n", aiName, repoSlug)
	return nil
}

func resolveTargetAi(flag string) string {
	if s := strings.TrimSpace(flag); s != "" {
		return store.NormalizeAIName(s)
	}
	return store.DetectCurrentAIName()
}

func resolveTargetRepo(flag string) string {
	if s := strings.TrimSpace(flag); s != "" {
		return store.SanitizeSlug(s)
	}
	return store.ResolveAmbientRepoSlug("")
}
}
