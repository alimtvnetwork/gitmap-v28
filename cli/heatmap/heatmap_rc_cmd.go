package heatmap

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// DispatchRcHeatmap routes gitmap heatmap rc subcommands.
func DispatchRcHeatmap(ctx context.Context, args []string, out io.Writer) error {
	if len(args) == 0 || args[0] == "help" || args[0] == "--help" || args[0] == "-h" {
		printRcHeatmapHelp(out)
		return nil
	}

	subcmd := strings.ToLower(args[0])
	flags, pos := parseHeatmapArgs(args[1:])

	switch subcmd {
	case "log":
		return runRcLog(flags, out)
	case "search":
		return runRcSearch(flags, pos, out)
	case "tree":
		return runRcTree(flags, out)
	case "status":
		return runRcStatus(flags, out)
	case "list-ai", "list_ai", "list":
		return runRcListAi(flags, out)
	case "commit":
		return runRcCommit(flags, out)
	case "sync":
		return runRcSync(flags, out)
	default:
		return apperror.NewValidation("heatmap_rc", "E1000", fmt.Sprintf("unknown rc subcommand: %s", subcmd))
	}
}

func runRcLog(flags map[string]string, out io.Writer) error {
	aiName := resolveTargetAi(flags["ai"])
	repoSlug := resolveTargetRepo(flags["repo"])

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

	fmt.Fprintf(out, "Logged RC heatmap action #%d [ai=%s, repo=%s, task=%s, action=%s, file=%s]\n",
		id, aiName, repoSlug, rec.TaskSlug, rec.ActionType, rec.FilePath)
	return nil
}

func runRcSearch(flags map[string]string, pos []string, out io.Writer) error {
	aiName := resolveTargetAi(flags["ai"])
	repoSlug := resolveTargetRepo(flags["repo"])

	st, err := store.OpenRcHeatmapStore(aiName, repoSlug)
	if err != nil {
		return err
	}
	defer st.Close()

	limit := 50
	if l, err := strconv.Atoi(flags["limit"]); err == nil && l > 0 {
		limit = l
	}

	query := strings.Join(pos, " ")
	records, err := st.SearchActions(query, flags["task"], flags["file"], limit)
	if err != nil {
		return err
	}

	if flags["json"] == "true" {
		return json.NewEncoder(out).Encode(records)
	}
	renderSearchText(out, records)
	return nil
}

func runRcTree(flags map[string]string, out io.Writer) error {
	aiName := resolveTargetAi(flags["ai"])
	repoSlug := resolveTargetRepo(flags["repo"])

	st, err := store.OpenRcHeatmapStore(aiName, repoSlug)
	if err != nil {
		return err
	}
	defer st.Close()

	tree, err := st.GetExecutionTree(flags["task"])
	if err != nil {
		return err
	}

	if flags["json"] == "true" {
		return json.NewEncoder(out).Encode(tree)
	}
	renderRcTreeText(out, tree, "")
	return nil
}

func runRcStatus(flags map[string]string, out io.Writer) error {
	aiName := resolveTargetAi(flags["ai"])
	repoSlug := resolveTargetRepo(flags["repo"])

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

	if flags["json"] == "true" {
		resp := map[string]any{
			"aiName":      aiName,
			"repoSlug":    repoSlug,
			"uncommitted": uncommitted,
			"summary":     summary,
		}
		return json.NewEncoder(out).Encode(resp)
	}

	fmt.Fprintf(out, "RC Heatmap Status [AI: %s, Repo: %s]\n", aiName, repoSlug)
	fmt.Fprintf(out, "Uncommitted Actions: %d\n", uncommitted)
	renderStatusText(out, *summary)
	return nil
}

func runRcListAi(flags map[string]string, out io.Writer) error {
	names := store.ListSupportedAINames()
	current := store.DetectCurrentAIName()

	if flags["json"] == "true" {
		resp := map[string]any{
			"currentAI":    current,
			"supportedAIs": names,
		}
		return json.NewEncoder(out).Encode(resp)
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
	aiName := resolveTargetAi(flags["ai"])
	repoSlug := resolveTargetRepo(flags["repo"])
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

	rec.SyncStatus = "COMMITTED"
	rec.IsDirty = false
	rec.UncommittedChangeCount = 0
	_ = st.UpdateLocalTracking(*rec)

	if flags["json"] == "true" {
		return json.NewEncoder(out).Encode(map[string]any{
			"status":   "COMMITTED",
			"aiName":   aiName,
			"repoSlug": repoSlug,
			"message":  msg,
		})
	}

	fmt.Fprintf(out, "Committed RC heatmap telemetry [AI: %s, Repo: %s]\nMessage: %s\n", aiName, repoSlug, msg)
	return nil
}

func runRcSync(flags map[string]string, out io.Writer) error {
	aiName := resolveTargetAi(flags["ai"])
	repoSlug := resolveTargetRepo(flags["repo"])

	st, err := store.OpenRcHeatmapStore(aiName, repoSlug)
	if err != nil {
		return err
	}
	defer st.Close()

	rec, err := st.GetOrCreateLocalTracking(repoSlug, aiName)
	if err != nil {
		return err
	}

	rec.SyncStatus = "SYNCED"
	rec.IsDirty = false
	_ = st.UpdateLocalTracking(*rec)

	if flags["json"] == "true" {
		return json.NewEncoder(out).Encode(map[string]any{
			"status":   "SYNCED",
			"aiName":   aiName,
			"repoSlug": repoSlug,
		})
	}

	fmt.Fprintf(out, "Synchronized RC heatmap telemetry with repo-cache [AI: %s, Repo: %s]\n", aiName, repoSlug)
	return nil
}

func resolveTargetAi(flag string) string {
	if strings.TrimSpace(flag) != "" {
		return store.NormalizeAIName(flag)
	}
	return store.DetectCurrentAIName()
}

func resolveTargetRepo(flag string) string {
	if strings.TrimSpace(flag) != "" {
		return store.SanitizeSlug(flag)
	}
	return store.ResolveAmbientRepoSlug("")
}

func renderRcTreeText(out io.Writer, nodes []store.RcHeatmapExecutionTreeNode, indent string) {
	for i, node := range nodes {
		branch, nextIndent := "├── ", indent+"│   "
		if i == len(nodes)-1 {
			branch, nextIndent = "└── ", indent+"    "
		}
		if indent == "" {
			branch, nextIndent = "• ", "  "
		}
		lineInfo := ""
		if node.FilePath != "" {
			lineInfo = fmt.Sprintf(" [%s:%d-%d]", node.FilePath, node.StartLine, node.EndLine)
		}
		fmt.Fprintf(out, "%s%s#%d (%s) %s%s: %s\n", indent, branch, node.HeatmapActionLogId,
			node.ReasoningState, node.ActionType, lineInfo, node.ReasoningNotes)
		if len(node.Children) > 0 {
			renderRcTreeText(out, node.Children, nextIndent)
		}
	}
}

func printRcHeatmapHelp(out io.Writer) {
	fmt.Fprintln(out, "Usage: gitmap heatmap rc <command> [flags]\n\nSubcommands:")
	fmt.Fprintln(out, "  log      Log an AI cognitive decision to repo-cache split-db")
	fmt.Fprintln(out, "  search   Search logged actions by query, task, or file")
	fmt.Fprintln(out, "  tree     Display hierarchical causal execution tree")
	fmt.Fprintln(out, "  status   Summarize action volume, states, and uncommitted delta")
	fmt.Fprintln(out, "  list-ai  List supported AI namespaces and ambient detected AI")
	fmt.Fprintln(out, "  commit   Record commit checkpoint and mark local changes clean")
	fmt.Fprintln(out, "  sync     Synchronize local telemetry with repo-cache")
	fmt.Fprintln(out, "  help     Show this help message\n\nFlags:")
	fmt.Fprintln(out, "  --ai      Target AI namespace (antigravity, muse, cloud-code, codex, grok, cursor)")
	fmt.Fprintln(out, "  --repo    Target repository slug (defaults to ambient working repo)")
	fmt.Fprintln(out, "  --task    Task slug")
	fmt.Fprintln(out, "  --file    Target file path")
	fmt.Fprintln(out, "  --json    Output results in JSON format")
}
