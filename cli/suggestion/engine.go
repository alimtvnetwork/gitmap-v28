package suggestion

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

type commandMeta struct {
	name        string
	description string
}

type flagMeta struct {
	flag string
	desc string
}

var canonicalCommands = []commandMeta{
	{name: "status", description: "Check working tree status across all repositories"},
	{name: "scan", description: "Fast repository scanner with parallel discovery"},
	{name: "clone", description: "Clone git repositories with high speed"},
	{name: "pull", description: "Pull upstream changes for repositories"},
	{name: "commit", description: "Commit staged changes with conventional format"},
	{name: "release", description: "Orchestrate version release and changelog bump"},
	{name: "changelog", description: "Generate or view repository changelog"},
	{name: "list-versions", description: "List published release versions"},
	{name: "reconcile", description: "Reconcile divergent git branches and conflicts"},
	{name: "stash", description: "Stash uncommitted changes across repositories"},
	{name: "wip", description: "Save temporary work-in-progress checkpoint"},
	{name: "group", description: "Organize repositories into workspaces"},
	{name: "cd", description: "Change active repository directory"},
	{name: "ssh", description: "Manage remote SSH keys and connections"},
	{name: "nodes", description: "List and manage distributed cluster nodes"},
	{name: "cluster", description: "Manage containerized and remote cluster fleet"},
	{name: "deploy", description: "Deploy services and configurations to nodes"},
	{name: "vhost", description: "Configure virtual hosts and reverse proxies"},
	{name: "nginx", description: "Manage nginx web server configurations"},
	{name: "agy", description: "Orchestrate Google Antigravity AI agents"},
	{name: "agm", description: "Manage Antigravity agent instances and logs"},
	{name: "aum", description: "Automate polyglot package and tooling updates"},
	{name: "ai", description: "AI analysis, session tracking, and model training"},
	{name: "macro", description: "Record and replay CLI macros and workflows"},
	{name: "pipeline", description: "Run or manage CI/CD pipeline quality gates"},
	{name: "pull-error", description: "Inspect pipeline errors and trace failing steps"},
	{name: "rerun", description: "Rerun failed agent steps and pipeline jobs"},
	{name: "sug", description: "Shutdown execution until quality gates are green"},
	{name: "os", description: "Cross-platform OS desktop and taskbar integrations"},
	{name: "apps", description: "Manage installed developer tools and applications"},
	{name: "install", description: "Install CLI companion tools and agents"},
	{name: "uninstall", description: "Uninstall CLI companion tools and agents"},
	{name: "storage", description: "Inspect disk usage and cache allocations"},
	{name: "clean", description: "Clean working tree and purge build artifacts"},
	{name: "vscode", description: "Synchronize VS Code workspaces and settings"},
	{name: "sync", description: "Synchronize prompts and repositories across fleet"},
	{name: "doctor", description: "Run diagnostic health checks on workspace"},
	{name: "help", description: "Display help information for commands"},
	{name: "version", description: "Print current gitmap version information"},
	{name: "task", description: "Manage SQLite agent tasks and execution status"},
	{name: "ai-analysis", description: "AI session reasoning and decision telemetry"},
}

var canonicalAliases = map[string]string{
	"s":    "status",
	"st":   "status",
	"pe":   "pull-error",
	"pl":   "pipeline",
	"cl":   "clone",
	"co":   "commit",
	"cpf":  "commit",
	"cpb":  "commit",
	"cpr":  "commit",
	"cpar": "commit",
	"sc":   "cluster",
	"sj":   "ssh",
	"sjc":  "ssh",
	"rr":   "rerun",
	"rra":  "rerun",
	"rrq":  "rerun",
	"sug":  "sug",
	"rel":  "release",
	"lp":   "list-prompts",
	"lv":   "list-versions",
}

var semanticSynonyms = map[string]string{
	"docker":  "cluster",
	"k8s":     "cluster",
	"pod":     "cluster",
	"grep":    "aum search",
	"rg":      "aum search",
	"find":    "scan",
	"search":  "aum search",
	"rm":      "clean",
	"remove":  "clean",
	"del":     "clean",
	"delete":  "clean",
	"history": "log",
	"test":    "pipeline",
	"check":   "status",
	"publish": "release",
	"build":   "pipeline",
	"revert":  "reconcile",
	"merge":   "reconcile",
	"branch":  "status",
	"purge":   "clean",
}

var standardFlags = []flagMeta{
	{flag: "--help", desc: "Show command line help and usage details"},
	{flag: "--verbose", desc: "Enable detailed logging output"},
	{flag: "--all", desc: "Process all target repositories"},
	{flag: "--quiet", desc: "Suppress non-essential terminal output"},
	{flag: "--json", desc: "Output machine-readable JSON format"},
	{flag: "--force", desc: "Bypass non-fatal safety confirmation checks"},
	{flag: "--dry-run", desc: "Simulate actions without modifying files"},
	{flag: "--theme", desc: "Select terminal color palette theme"},
	{flag: "--glyphs", desc: "Configure UTF-8 emoji and glyph mode"},
}

type defaultEngine struct{}

// NewEngine creates a new 4-tier suggestion resolution engine.
func NewEngine() Engine {
	return &defaultEngine{}
}

func (e *defaultEngine) ResolveCommand(token string) SuggestionGroup {
	cleanToken := strings.ToLower(strings.TrimSpace(token))
	if len(cleanToken) == 0 {
		return SuggestionGroup{Title: "Empty Token", Reason: "No command token provided"}
	}
	candidates := e.gatherCommandCandidates(cleanToken)
	uniqueList := deduplicateSuggestions(candidates)
	sortSuggestionsByConfidence(uniqueList)
	limited := limitSuggestions(uniqueList, 5)
	return buildCommandSuggestionGroup(cleanToken, limited)
}

func (e *defaultEngine) gatherCommandCandidates(token string) []Suggestion {
	var out []Suggestion
	out = append(out, matchTier1Aliases(token)...)
	out = append(out, matchTier2Levenshtein(token, 3)...)
	out = append(out, matchTier3Prefix(token)...)
	out = append(out, matchTier4Synonyms(token)...)
	return out
}

func buildCommandSuggestionGroup(token string, items []Suggestion) SuggestionGroup {
	reason := fmt.Sprintf("Command '%s' is not recognized. Closest matches:", token)
	return SuggestionGroup{
		Title:       "Did you mean?",
		Reason:      reason,
		InputToken:  token,
		Suggestions: items,
	}
}

func matchTier1Aliases(token string) []Suggestion {
	targetCmd, hasAlias := canonicalAliases[token]
	if !hasAlias {
		return nil
	}
	desc := findCommandDescription(targetCmd)
	return []Suggestion{{
		Command:     "gitmap " + targetCmd,
		Description: desc,
		Category:    CategoryIntent,
		Confidence:  1.0,
		ActionType:  ActionAutoRun,
	}}
}

func findCommandDescription(cmdName string) string {
	for _, item := range canonicalCommands {
		if item.name == cmdName {
			return item.description
		}
	}
	return "Execute gitmap " + cmdName
}

func matchTier2Levenshtein(token string, maxDist int) []Suggestion {
	var results []Suggestion
	for _, meta := range canonicalCommands {
		dist := computeLevenshtein(token, meta.name)
		if dist > 0 && dist <= maxDist {
			results = append(results, buildLevenshteinSuggestion(token, meta, dist))
		}
	}
	return results
}

func buildLevenshteinSuggestion(token string, meta commandMeta, dist int) Suggestion {
	maxLen := maxInt(len(token), len(meta.name))
	conf := calculateNormalizedConfidence(dist, maxLen)
	action := ActionManual
	if conf >= 0.75 {
		action = ActionAutoRun
	}
	return Suggestion{
		Command:     "gitmap " + meta.name,
		Description: meta.description,
		Category:    CategoryTypo,
		Confidence:  conf,
		ActionType:  action,
	}
}

func calculateNormalizedConfidence(dist, maxLen int) float64 {
	if maxLen == 0 {
		return 0.0
	}
	raw := 1.0 - (float64(dist) / float64(maxLen))
	return math.Round(raw*100) / 100
}

func computeLevenshtein(s1, s2 string) int {
	if s1 == s2 {
		return 0
	}
	if len(s1) == 0 {
		return len(s2)
	}
	if len(s2) == 0 {
		return len(s1)
	}
	return calculateRowLevenshtein(s1, s2)
}

func calculateRowLevenshtein(s1, s2 string) int {
	r1, r2 := []rune(s1), []rune(s2)
	prev := make([]int, len(r2)+1)
	for j := 0; j <= len(r2); j++ {
		prev[j] = j
	}
	for i := 1; i <= len(r1); i++ {
		prev = stepLevenshteinRow(r1[i-1], r2, prev)
	}
	return prev[len(r2)]
}

func stepLevenshteinRow(c1 rune, r2 []rune, prev []int) []int {
	curr := make([]int, len(r2)+1)
	curr[0] = prev[0] + 1
	for j := 1; j <= len(r2); j++ {
		cost := 0
		if c1 != r2[j-1] {
			cost = 1
		}
		curr[j] = min3(curr[j-1]+1, prev[j]+1, prev[j-1]+cost)
	}
	return curr
}

func min3(a, b, c int) int {
	m := a
	if b < m {
		m = b
	}
	if c < m {
		m = c
	}
	return m
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func matchTier3Prefix(token string) []Suggestion {
	if len(token) < 2 {
		return nil
	}
	var results []Suggestion
	for _, meta := range canonicalCommands {
		if isSubsequenceOrPrefix(token, meta.name) {
			results = append(results, buildPrefixSuggestion(token, meta))
		}
	}
	return results
}

func isSubsequenceOrPrefix(needle, haystack string) bool {
	if strings.HasPrefix(haystack, needle) {
		return true
	}
	return isSubsequence(needle, haystack)
}

func isSubsequence(needle, haystack string) bool {
	if len(needle) > len(haystack) {
		return false
	}
	nIdx := 0
	rNeedle := []rune(needle)
	for _, r := range haystack {
		if nIdx < len(rNeedle) && r == rNeedle[nIdx] {
			nIdx++
		}
	}
	return nIdx == len(rNeedle)
}

func buildPrefixSuggestion(token string, meta commandMeta) Suggestion {
	conf := float64(len(token)) / float64(len(meta.name)) * 0.85
	roundedConf := math.Round(conf*100) / 100
	return Suggestion{
		Command:     "gitmap " + meta.name,
		Description: meta.description,
		Category:    CategorySubcommand,
		Confidence:  roundedConf,
		ActionType:  ActionManual,
	}
}

func matchTier4Synonyms(token string) []Suggestion {
	targetCmd, hasSynonym := semanticSynonyms[token]
	if !hasSynonym {
		return nil
	}
	desc := findCommandDescription(targetCmd)
	return []Suggestion{{
		Command:     "gitmap " + targetCmd,
		Description: desc,
		Category:    CategoryIntent,
		Confidence:  0.80,
		ActionType:  ActionManual,
	}}
}

func deduplicateSuggestions(items []Suggestion) []Suggestion {
	seen := make(map[string]bool)
	var out []Suggestion
	for _, s := range items {
		if !seen[s.Command] {
			seen[s.Command] = true
			out = append(out, s)
		}
	}
	return out
}

func sortSuggestionsByConfidence(items []Suggestion) {
	sort.Slice(items, func(i, j int) bool {
		return items[i].Confidence > items[j].Confidence
	})
}

func limitSuggestions(items []Suggestion, limit int) []Suggestion {
	if len(items) <= limit {
		return items
	}
	return items[:limit]
}

func (e *defaultEngine) ResolveFlag(cmdName, flagToken string) SuggestionGroup {
	cleanFlag := strings.ToLower(strings.TrimSpace(flagToken))
	if len(cleanFlag) == 0 {
		return SuggestionGroup{Title: "Empty Flag", Reason: "No flag token provided"}
	}
	results := matchFlagCandidates(cleanFlag)
	sortSuggestionsByConfidence(results)
	limited := limitSuggestions(results, 3)
	return SuggestionGroup{
		Title:       "Flag Suggestion",
		Reason:      fmt.Sprintf("Flag '%s' is not recognized. Did you mean:", flagToken),
		InputToken:  flagToken,
		Suggestions: limited,
	}
}

func matchFlagCandidates(flagToken string) []Suggestion {
	var out []Suggestion
	for _, f := range standardFlags {
		if isFlagMatch(flagToken, f.flag) {
			out = append(out, buildFlagSuggestion(flagToken, f))
		}
	}
	return out
}

func isFlagMatch(token, flag string) bool {
	dist := computeLevenshtein(token, flag)
	return dist <= 3 || strings.HasPrefix(flag, token)
}

func buildFlagSuggestion(token string, f flagMeta) Suggestion {
	dist := computeLevenshtein(token, f.flag)
	conf := calculateNormalizedConfidence(dist, maxInt(len(token), len(f.flag)))
	return Suggestion{
		Command:     f.flag,
		Description: f.desc,
		Category:    CategoryFlag,
		Confidence:  conf,
		ActionType:  ActionManual,
	}
}

func (e *defaultEngine) ResolveRemediation(err error, ctx map[string]any) SuggestionGroup {
	if err == nil {
		return SuggestionGroup{Title: "No Error", Reason: "No error occurred"}
	}
	errMsg := strings.ToLower(err.Error())
	items := buildRemediationSuggestions(errMsg)
	return SuggestionGroup{
		Title:       "Recommended Action",
		Reason:      err.Error(),
		Suggestions: items,
	}
}

func buildRemediationSuggestions(msg string) []Suggestion {
	if strings.Contains(msg, "dirty") || strings.Contains(msg, "working tree") {
		return makeDirtyTreeSuggestions()
	}
	if strings.Contains(msg, "conflict") || strings.Contains(msg, "diverg") {
		return makeConflictSuggestions()
	}
	if strings.Contains(msg, "pipeline") || strings.Contains(msg, "gate") {
		return makePipelineSuggestions()
	}
	return makeDefaultRemediationSuggestions()
}

func makeDirtyTreeSuggestions() []Suggestion {
	return []Suggestion{
		{Command: "gitmap status", Description: "Inspect uncommitted files across workspaces", Category: CategoryRemediation, Confidence: 0.95, ActionType: ActionAutoRun},
		{Command: "gitmap stash", Description: "Stash uncommitted changes safely", Category: CategoryRemediation, Confidence: 0.85, ActionType: ActionManual},
	}
}

func makeConflictSuggestions() []Suggestion {
	return []Suggestion{
		{Command: "gitmap reconcile", Description: "Reconcile divergent branches and auto-resolve known conflict patterns", Category: CategoryRemediation, Confidence: 0.95, ActionType: ActionAutoRun},
	}
}

func makePipelineSuggestions() []Suggestion {
	return []Suggestion{
		{Command: "gitmap pe -t", Description: "Inspect pipeline error logs and dynamic runner telemetry", Category: CategoryRemediation, Confidence: 0.95, ActionType: ActionAutoRun},
		{Command: "gitmap rerun", Description: "Rerun failed pipeline steps", Category: CategoryRemediation, Confidence: 0.90, ActionType: ActionAutoRun},
	}
}

func makeDefaultRemediationSuggestions() []Suggestion {
	return []Suggestion{
		{Command: "gitmap doctor", Description: "Run system diagnostics to inspect environment health", Category: CategoryRemediation, Confidence: 0.70, ActionType: ActionAutoRun},
	}
}
