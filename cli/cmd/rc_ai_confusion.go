package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// RunRCAiConfusion dispatches subcommands for 'gitmap rc ai-confusion'.
func RunRCAiConfusion(args []string) error {
	hasEmptyArgs := len(args) == 0
	if hasEmptyArgs {
		fmt.Print(RenderAiConfusionHelp())

		return nil
	}

	sub := strings.ToLower(strings.TrimSpace(args[0]))

	return dispatchConfusionSubcmd(sub, args[1:])
}

func dispatchConfusionSubcmd(sub string, rest []string) error {
	if isConfusionHelpSubcmd(sub) {
		fmt.Print(RenderAiConfusionHelp())

		return nil
	}

	return executeConfusionSubcmd(sub, rest)
}

func isConfusionHelpSubcmd(sub string) bool {
	return sub == "help" || sub == "-h" || sub == "--help"
}

func executeConfusionSubcmd(sub string, rest []string) error {
	switch sub {
	case "add":
		return runAiConfusionAdd(rest)
	case "view", "ls", "list", "show":
		return runAiConfusionView(rest)
	case "complete", "resolve", "done":
		return runAiConfusionComplete(rest)
	case "sync":
		return runAiConfusionSync(rest)
	default:
		return apperror.NewValidationError("unknown ai-confusion subcommand: " + sub)
	}
}

type addCmdOptions struct {
	jsonFilePath string
	record       store.AiConfusionRecord
}

func runAiConfusionAdd(args []string) error {
	rec, err := parseAddOptions(args)
	if err != nil {
		return err
	}

	valErr := validateAddRecord(rec)
	if valErr != nil {
		return valErr
	}

	return persistAndRenderConfusion(rec)
}

func persistAndRenderConfusion(rec store.AiConfusionRecord) error {
	st, openErr := store.OpenAiConfusionStore("")
	if openErr != nil {
		return openErr
	}
	defer st.Close()

	if insErr := st.InsertAiConfusion(&rec); insErr != nil {
		return insErr
	}

	renderAddConfirmation(rec, st.DBPath())

	return nil
}

func parseAddOptions(args []string) (store.AiConfusionRecord, error) {
	opts := parseAddFlags(args)
	rec, fileErr := maybeLoadJsonFile(opts.jsonFilePath)
	if fileErr != nil {
		return rec, fileErr
	}

	mergeAddFlags(&rec, opts.record)

	return rec, nil
}

func parseAddFlags(args []string) addCmdOptions {
	var opts addCmdOptions
	for i := 0; i < len(args); i++ {
		i = consumeAddToken(&opts, args, i)
	}

	return opts
}

func consumeAddToken(opts *addCmdOptions, args []string, idx int) int {
	tok := args[idx]
	hasNext := idx+1 < len(args)
	if hasNext && applyAddParamToken(opts, tok, args[idx+1]) {
		return idx + 1
	}

	return idx
}

func applyAddParamToken(opts *addCmdOptions, flag, val string) bool {
	if applyFileOrAgentToken(opts, flag, val) {
		return true
	}

	if applyPathOrSpecToken(opts, flag, val) {
		return true
	}

	if applyTaskToken(opts, flag, val) {
		return true
	}

	if applyContextToken(opts, flag, val) {
		return true
	}

	return applyNarrativeToken(opts, flag, val)
}

func applyFileOrAgentToken(opts *addCmdOptions, flag, val string) bool {
	switch flag {
	case "--file", "-f":
		opts.jsonFilePath = val
		return true
	case "--agent", "-a":
		opts.record.AgentSlug = val
		return true
	default:
		return false
	}
}

func applyPathOrSpecToken(opts *addCmdOptions, flag, val string) bool {
	switch flag {
	case "--files", "--file-path", "--path":
		opts.record.FileNames = val
		return true
	case "--md-files", "--md":
		opts.record.MdFiles = val
		return true
	case "--spec", "--spec-path":
		opts.record.Spec = val
		return true
	default:
		return false
	}
}

func applyTaskToken(opts *addCmdOptions, flag, val string) bool {
	switch flag {
	case "--task", "-t":
		opts.record.Task = val
		return true
	case "--subtask", "-s":
		opts.record.Subtask = val
		return true
	default:
		return false
	}
}

func applyContextToken(opts *addCmdOptions, flag, val string) bool {
	switch flag {
	case "--heatmap", "--heatmap-info":
		opts.record.HeatmapInfo = val
		return true
	case "--lines", "--line-numbers":
		opts.record.LineNumbers = val
		return true
	case "--code", "--snippet":
		opts.record.CodeSnippet = val
		return true
	default:
		return false
	}
}

func applyNarrativeToken(opts *addCmdOptions, flag, val string) bool {
	switch flag {
	case "--reason", "--notes":
		opts.record.Reason = val
		return true
	case "--intention", "--int":
		opts.record.InitialIntention = val
		return true
	case "--thinking", "--thought":
		opts.record.ReasoningThinking = val
		return true
	case "--question", "-q":
		opts.record.ConfusionQuestion = val
		return true
	default:
		return false
	}
}

func maybeLoadJsonFile(filePath string) (store.AiConfusionRecord, error) {
	hasEmptyPath := len(strings.TrimSpace(filePath)) == 0
	if hasEmptyPath {
		return store.AiConfusionRecord{}, nil
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return store.AiConfusionRecord{}, apperror.WrapSimple(err, "read confusion json file")
	}

	var jsonIn store.AiConfusionJSONInput
	if jsonErr := json.Unmarshal(data, &jsonIn); jsonErr != nil {
		return store.AiConfusionRecord{}, apperror.WrapSimple(jsonErr, "parse confusion json file")
	}

	return jsonIn.ToRecord(), nil
}

func mergeAddFlags(dest *store.AiConfusionRecord, flagRec store.AiConfusionRecord) {
	mergeIdentityFlags(dest, flagRec)
	mergeContextFlags(dest, flagRec)
	mergeNarrativeFlags(dest, flagRec)
}

func mergeIdentityFlags(dest *store.AiConfusionRecord, flagRec store.AiConfusionRecord) {
	dest.AgentSlug = resolveNonEmpty(flagRec.AgentSlug, dest.AgentSlug)
	dest.FileNames = resolveNonEmpty(flagRec.FileNames, dest.FileNames)
	dest.MdFiles = resolveNonEmpty(flagRec.MdFiles, dest.MdFiles)
	dest.Spec = resolveNonEmpty(flagRec.Spec, dest.Spec)
	dest.Task = resolveNonEmpty(flagRec.Task, dest.Task)
	dest.Subtask = resolveNonEmpty(flagRec.Subtask, dest.Subtask)
}

func mergeContextFlags(dest *store.AiConfusionRecord, flagRec store.AiConfusionRecord) {
	dest.HeatmapInfo = resolveNonEmpty(flagRec.HeatmapInfo, dest.HeatmapInfo)
	dest.LineNumbers = resolveNonEmpty(flagRec.LineNumbers, dest.LineNumbers)
	dest.CodeSnippet = resolveNonEmpty(flagRec.CodeSnippet, dest.CodeSnippet)
	dest.Reason = resolveNonEmpty(flagRec.Reason, dest.Reason)
}

func mergeNarrativeFlags(dest *store.AiConfusionRecord, flagRec store.AiConfusionRecord) {
	dest.InitialIntention = resolveNonEmpty(flagRec.InitialIntention, dest.InitialIntention)
	dest.ReasoningThinking = resolveNonEmpty(flagRec.ReasoningThinking, dest.ReasoningThinking)
	dest.ConfusionQuestion = resolveNonEmpty(flagRec.ConfusionQuestion, dest.ConfusionQuestion)
}

func resolveNonEmpty(override, fallback string) string {
	hasOverride := len(strings.TrimSpace(override)) > 0
	if hasOverride {
		return override
	}

	return fallback
}

func validateAddRecord(rec store.AiConfusionRecord) error {
	hasPayload := hasAnyConfusionPayload(rec)
	if hasPayload {
		return nil
	}

	return apperror.NewValidationError("confusion record requires at least one of --question, --intention, --thinking, --code, or --files")
}

func hasAnyConfusionPayload(rec store.AiConfusionRecord) bool {
	hasQ := len(strings.TrimSpace(rec.ConfusionQuestion)) > 0
	hasInt := len(strings.TrimSpace(rec.InitialIntention)) > 0
	hasTh := len(strings.TrimSpace(rec.ReasoningThinking)) > 0
	hasCode := len(strings.TrimSpace(rec.CodeSnippet)) > 0
	hasFile := len(strings.TrimSpace(rec.FileNames)) > 0

	return hasQ || hasInt || hasTh || hasCode || hasFile
}

func renderAddConfirmation(rec store.AiConfusionRecord, dbPath string) {
	fmt.Printf("\n┌── AI Confusion Episode Recorded ────────────────────────────────────┐\n")
	fmt.Printf("│ ID          : %d (%s)\n", rec.AiConfusionId, rec.ConfusionUuid)
	fmt.Printf("│ Agent       : %s\n", rec.AgentSlug)
	fmt.Printf("│ Target File : %s (Lines: %s)\n", rec.FileNames, rec.LineNumbers)
	fmt.Printf("│ Task        : %s / %s\n", rec.Task, rec.Subtask)
	fmt.Printf("│ Intention   : %s\n", rec.InitialIntention)
	fmt.Printf("│ Question    : %s\n", rec.ConfusionQuestion)
	fmt.Printf("│ Status      : PENDING\n")
	fmt.Printf("└──────────────────────────────────────────────────────────────────────┘\n")
	fmt.Printf("  ✔ Successfully persisted to %s\n\n", dbPath)
}

type viewCmdOptions struct {
	id        string
	task      string
	agent     string
	isAll     bool
	isJSON    bool
	isPending bool
}

func runAiConfusionView(args []string) error {
	opts := parseViewOptions(args)
	st, err := store.OpenAiConfusionStore("")
	if err != nil {
		return err
	}
	defer st.Close()

	filter := buildStoreFilter(opts)
	records, listErr := st.ListAiConfusions(filter)
	if listErr != nil {
		return listErr
	}

	return renderViewOutput(records, opts.isJSON)
}

func buildStoreFilter(opts viewCmdOptions) store.AiConfusionFilter {
	filter := store.AiConfusionFilter{
		Id:        opts.id,
		Task:      opts.task,
		AgentSlug: opts.agent,
	}
	hasCompletedFilter := opts.isPending || (opts.isAll == false && len(opts.id) == 0)
	if hasCompletedFilter {
		isComp := false
		filter.IsCompleted = &isComp
	}

	return filter
}

func parseViewOptions(args []string) viewCmdOptions {
	var opts viewCmdOptions
	for i := 0; i < len(args); i++ {
		i = consumeViewToken(&opts, args, i)
	}

	return opts
}

func consumeViewToken(opts *viewCmdOptions, args []string, idx int) int {
	tok := args[idx]
	if tok == "--all" || tok == "-a" {
		opts.isAll = true
		return idx
	}

	if tok == "--json" || tok == "-j" {
		opts.isJSON = true
		return idx
	}

	if tok == "--pending" || tok == "-p" {
		opts.isPending = true
		return idx
	}

	return consumeViewParamOrPositional(opts, args, idx)
}

func consumeViewParamOrPositional(opts *viewCmdOptions, args []string, idx int) int {
	tok := args[idx]
	hasNext := idx+1 < len(args)
	if hasNext && applyViewParamToken(opts, tok, args[idx+1]) {
		return idx + 1
	}

	isOptionFlag := strings.HasPrefix(tok, "-")
	hasEmptyId := len(opts.id) == 0
	if isOptionFlag == false && hasEmptyId {
		opts.id = tok
	}

	return idx
}

func applyViewParamToken(opts *viewCmdOptions, flag, val string) bool {
	switch flag {
	case "--id":
		opts.id = val
		return true
	case "--task", "-t":
		opts.task = val
		return true
	case "--agent":
		opts.agent = val
		return true
	default:
		return false
	}
}

func renderViewOutput(records []store.AiConfusionRecord, isJSON bool) error {
	if isJSON {
		return renderRecordsJSON(records)
	}

	hasNoRecords := len(records) == 0
	if hasNoRecords {
		fmt.Println("  No AI confusion episodes found matching filter.")

		return nil
	}

	for _, rec := range records {
		renderSingleConfusionCard(rec)
	}

	return nil
}

func renderRecordsJSON(records []store.AiConfusionRecord) error {
	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return apperror.WrapSimple(err, "marshal confusion records to json")
	}

	fmt.Println(string(data))

	return nil
}

func renderSingleConfusionCard(rec store.AiConfusionRecord) {
	statusStr := "PENDING"
	if rec.IsCompleted {
		statusStr = "COMPLETED"
	}

	fmt.Printf("\n┌── AI Confusion Episode [ID: %d | UUID: %s] ──\n", rec.AiConfusionId, rec.ConfusionUuid)
	fmt.Printf("│ Agent       : %s\n", rec.AgentSlug)
	fmt.Printf("│ Task/Subtask: %s / %s\n", rec.Task, rec.Subtask)
	fmt.Printf("│ Spec        : %s\n", rec.Spec)
	fmt.Printf("│ Files/Lines : %s (Lines: %s)\n", rec.FileNames, rec.LineNumbers)
	fmt.Printf("│ Status      : %s\n", statusStr)
	renderCardNarrativeSection(rec)
	renderCardCodeSection(rec)
	fmt.Printf("└────────────────────────────────────────────────────────────────────\n")
}

func renderCardNarrativeSection(rec store.AiConfusionRecord) {
	fmt.Printf("├────────────────────────────────────────────────────────────────────\n")
	renderCardIntentionAndThinking(rec)
	renderCardQuestionAndReason(rec)
}

func renderCardIntentionAndThinking(rec store.AiConfusionRecord) {
	hasIntention := len(rec.InitialIntention) > 0
	if hasIntention {
		fmt.Printf("│ Intention   : %s\n", rec.InitialIntention)
	}

	hasThinking := len(rec.ReasoningThinking) > 0
	if hasThinking {
		fmt.Printf("│ Thinking    : %s\n", rec.ReasoningThinking)
	}
}

func renderCardQuestionAndReason(rec store.AiConfusionRecord) {
	hasQuestion := len(rec.ConfusionQuestion) > 0
	if hasQuestion {
		fmt.Printf("│ Question    : %s\n", rec.ConfusionQuestion)
	}

	hasReason := len(rec.Reason) > 0
	if hasReason {
		fmt.Printf("│ Reason/Notes: %s\n", rec.Reason)
	}
}

func renderCardCodeSection(rec store.AiConfusionRecord) {
	hasCode := len(strings.TrimSpace(rec.CodeSnippet)) > 0
	if hasCode {
		fmt.Printf("│ Code Snippet:\n│   %s\n", strings.ReplaceAll(rec.CodeSnippet, "\n", "\n│   "))
	}
}

func runAiConfusionComplete(args []string) error {
	id, notes := parseCompleteArgs(args)
	hasEmptyId := len(id) == 0
	if hasEmptyId {
		return apperror.NewValidationError("id or uuid is required: gitmap rc ai-confusion complete <id|uuid> [--reason <notes>]")
	}

	st, err := store.OpenAiConfusionStore("")
	if err != nil {
		return err
	}
	defer st.Close()

	if compErr := st.MarkConfusionCompleted(id, notes); compErr != nil {
		return compErr
	}

	fmt.Printf("  ✔ Marked AI confusion episode %s as COMPLETED.\n", id)

	return nil
}

func parseCompleteArgs(args []string) (string, string) {
	var id, notes string
	for i := 0; i < len(args); i++ {
		i = consumeCompleteToken(args, i, &id, &notes)
	}

	return id, notes
}

func consumeCompleteToken(args []string, idx int, id, notes *string) int {
	tok := args[idx]
	hasNext := idx+1 < len(args)
	if hasNext && isCompleteNotesFlag(tok) {
		*notes = args[idx+1]

		return idx + 1
	}

	if hasNext && tok == "--id" {
		*id = args[idx+1]

		return idx + 1
	}

	isPositional := strings.HasPrefix(tok, "-") == false && len(*id) == 0
	if isPositional {
		*id = tok
	}

	return idx
}

func isCompleteNotesFlag(tok string) bool {
	return tok == "--reason" || tok == "--notes" || tok == "-r"
}

func runAiConfusionSync(args []string) error {
	fmt.Println("  Initiating safe 6-stage AI confusion synchronization...")
	res, err := store.SyncAiConfusion("")
	if err != nil {
		return err
	}

	renderSyncResults(res)

	return nil
}

func renderSyncResults(res *store.AiConfusionSyncResult) {
	fmt.Printf("  ✔ AI confusion database synchronized successfully.\n")
	fmt.Printf("    • Exported : %d record(s)\n", res.TotalExported)
	fmt.Printf("    • Replayed : %d record(s)\n", res.TotalReplayed)
	fmt.Printf("    • Committed: %v\n", res.IsCommitted)
	fmt.Printf("    • Pushed   : %v\n", res.IsPushed)
	fmt.Printf("    • Database : %s\n\n", res.DbPath)
}

// RenderAiConfusionHelp returns the help documentation card for AI confusion commands.
func RenderAiConfusionHelp() string {
	return `╭── GitMap AI Confusion Engine (gitmap rc ai-confusion) ─────────────────────╮
│ Record     : gitmap rc ai-confusion add [--file <path>] [--agent <a>]      │
│              [--task <t>] [--subtask <s>] [--spec <p>] [--files <f>]       │
│              [--lines <l>] [--code <c>] [--intention <i>] [--thinking <t>]│
│              [--question <q>] [--reason <r>]                               │
│ Inspect    : gitmap rc ai-confusion view [--id <id>] [--task <t>]         │
│              [--agent <a>] [--all] [--json]                                │
│ Resolve    : gitmap rc ai-confusion complete <id|uuid> [--reason <notes>] │
│ Sync       : gitmap rc ai-confusion sync                                  │
╰────────────────────────────────────────────────────────────────────────────╯
`
}
