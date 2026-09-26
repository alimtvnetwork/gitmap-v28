package commitin

import (
	"encoding/json"
	"os"
	"regexp"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmd/commitin/profile"
)

var configVarRegex = regexp.MustCompile(`\$\{([a-zA-Z0-9_.-]+)\}|\$([a-zA-Z0-9_.-]+)`)

// SyncStateTemplatesHook allows the orchestrator or cmd layer to sync
// imported JSON files into gitmap-templates.db with hash deduplication.
var SyncStateTemplatesHook func(importPaths []string) error

// PrecompileStateCategoryHook loads pre-compiled templates from gitmap-templates.db.
var PrecompileStateCategoryHook func(categoryOrSlug string, extraVars map[string]string) []string

// ConfigLineSkipper defines a line-stripping rule in commit-in config JSON.
type ConfigLineSkipper struct {
	Mode    string `json:"mode"`
	Pattern string `json:"pattern"`
}

// ImportedTemplateItem matches the JSON structure of templates in export/import files.
type ImportedTemplateItem struct {
	ID          string         `json:"id"`
	Category    string         `json:"category"`
	SubCategory string         `json:"subCategory,omitempty"`
	Slug        string         `json:"slug"`
	Title       string         `json:"title"`
	Text        string         `json:"text"`
	Additional  map[string]any `json:"additional,omitempty"`
}

// ImportedTemplateCategory represents a category entry with optional nested items in export/import files.
type ImportedTemplateCategory struct {
	Slug  string                 `json:"slug"`
	Name  string                 `json:"name,omitempty"`
	Items []ImportedTemplateItem `json:"items,omitempty"`
}

// ImportedTemplateFile represents an exported/imported templates + variables JSON payload.
type ImportedTemplateFile struct {
	ExportID   string                     `json:"exportId"`
	Version    string                     `json:"version"`
	Variables  map[string]string          `json:"variables"`
	Categories []ImportedTemplateCategory `json:"categories,omitempty"`
	Templates  []ImportedTemplateItem     `json:"templates"`
}

// CommitInConfigJSON defines the declarative migration configuration file.
type CommitInConfigJSON struct {
	Target            string                         `json:"target"`
	Inputs            []string                       `json:"inputs"`
	PRMode            string                         `json:"prMode,omitempty"`
	ConflictMode      string                         `json:"conflictMode,omitempty"`
	IsTree            bool                           `json:"tree,omitempty"`
	IsFinalSync       bool                           `json:"finalSync,omitempty"`
	IsCD              bool                           `json:"cd,omitempty"`
	IsRecreate        bool                           `json:"recreate,omitempty"`
	IsDryRun          bool                           `json:"dryRun,omitempty"`
	IsNewlineGap      bool                           `json:"newlineGap,omitempty"`
	AuthorName        string                         `json:"authorName,omitempty"`
	AuthorEmail       string                         `json:"authorEmail,omitempty"`
	Exclude           []string                       `json:"exclude,omitempty"`
	Imports           []string                       `json:"imports,omitempty"`
	Variables         map[string]string              `json:"variables,omitempty"`
	LineSkippers      []ConfigLineSkipper            `json:"lineSkippers,omitempty"`
	TitleReplacements []profile.TitleReplacementRule `json:"titleReplacements,omitempty"`
	PrefixTemplates   []string                       `json:"prefixTemplates,omitempty"`
	SuffixTemplates   []string                       `json:"suffixTemplates,omitempty"`
	SuffixSeparator   string                         `json:"suffixSeparator,omitempty"`
	PrefixSeparator   string                         `json:"prefixSeparator,omitempty"`
	PushImmediate     bool                           `json:"pushImmediate,omitempty"`
	SummaryDir        string                         `json:"summaryDir,omitempty"`
	TitlePrefix       string                         `json:"titlePrefix,omitempty"`
	TitleSuffix       string                         `json:"titleSuffix,omitempty"`
}

func applyConfigFileIfPresent(raw *RawArgs) *ParseError {
	if raw.ConfigPath == "" {
		return nil
	}
	data, err := os.ReadFile(raw.ConfigPath)
	if err != nil {
		return newBadArgs("failed to read config %q: %v", raw.ConfigPath, err)
	}
	var cfg CommitInConfigJSON
	if err := json.Unmarshal(data, &cfg); err != nil {
		return newBadArgs("invalid config JSON %q: %v", raw.ConfigPath, err)
	}
	cfg = expandConfigVariables(cfg)
	mergeConfigIntoRaw(raw, cfg)

	return loadAndPrecompileConfigTemplates(raw, cfg)
}

func mergeConfigIntoRaw(raw *RawArgs, cfg CommitInConfigJSON) {
	if raw.Source == "" && cfg.Target != "" {
		raw.Source = cfg.Target
	}
	if len(raw.Inputs) == 0 && len(cfg.Inputs) > 0 {
		raw.Inputs = splitInputs(cfg.Inputs)
	}
	mergeConfigFlags(raw, cfg)
	mergeConfigAuthor(raw, cfg)
	mergeConfigRules(raw, cfg)
}

func mergeConfigFlags(raw *RawArgs, cfg CommitInConfigJSON) {
	if raw.PRMode == "" && cfg.PRMode != "" {
		raw.PRMode = cfg.PRMode
	}
	if raw.ConflictMode == "" && cfg.ConflictMode != "" {
		raw.ConflictMode = cfg.ConflictMode
	}
	raw.IsTree = raw.IsTree || cfg.IsTree
	raw.IsFinalSync = raw.IsFinalSync || cfg.IsFinalSync
	raw.IsCD = raw.IsCD || cfg.IsCD
	raw.IsRecreate = raw.IsRecreate || cfg.IsRecreate
	raw.IsDryRun = raw.IsDryRun || cfg.IsDryRun
	raw.IsPushImmediate = raw.IsPushImmediate || cfg.PushImmediate
	if raw.SuffixSeparator == "" && cfg.SuffixSeparator != "" {
		raw.SuffixSeparator = cfg.SuffixSeparator
	}
	if raw.PrefixSeparator == "" && cfg.PrefixSeparator != "" {
		raw.PrefixSeparator = cfg.PrefixSeparator
	}
	if raw.SummaryDir == "" && cfg.SummaryDir != "" {
		raw.SummaryDir = cfg.SummaryDir
	}
}

func mergeConfigAuthor(raw *RawArgs, cfg CommitInConfigJSON) {
	if raw.AuthorName == "" && cfg.AuthorName != "" {
		raw.AuthorName = cfg.AuthorName
	}
	if raw.AuthorEmail == "" && cfg.AuthorEmail != "" {
		raw.AuthorEmail = cfg.AuthorEmail
	}
}

func mergeConfigRules(raw *RawArgs, cfg CommitInConfigJSON) {
	if len(cfg.Exclude) > 0 {
		raw.Exclude = append(raw.Exclude, cfg.Exclude...)
	}
	for _, ls := range cfg.LineSkippers {
		if ls.Mode != "" && ls.Pattern != "" {
			raw.MessageRules = append(raw.MessageRules, MessageRuleArg{Kind: ls.Mode, Value: ls.Pattern})
		}
	}
	if len(cfg.TitleReplacements) > 0 {
		raw.TitleReplacements = append(raw.TitleReplacements, cfg.TitleReplacements...)
	}
	if raw.TitlePrefix == "" && cfg.TitlePrefix != "" {
		raw.TitlePrefix = cfg.TitlePrefix
	}
	if raw.TitleSuffix == "" && cfg.TitleSuffix != "" {
		raw.TitleSuffix = cfg.TitleSuffix
	}
}

func loadAndPrecompileConfigTemplates(raw *RawArgs, cfg CommitInConfigJSON) *ParseError {
	if SyncStateTemplatesHook != nil && len(cfg.Imports) > 0 {
		_ = SyncStateTemplatesHook(cfg.Imports)
	}
	vars, items := readImportedTemplateFiles(cfg.Imports, cfg.Variables)
	prefixes := resolveTemplateRefs(cfg.PrefixTemplates, items, vars)
	suffixes := resolveTemplateRefs(cfg.SuffixTemplates, items, vars)
	if len(suffixes) == 0 && len(prefixes) == 0 && len(items) > 0 {
		suffixes = precompileAllImportedItems(items, vars)
	}
	raw.MessagePrefix = append(raw.MessagePrefix, prefixes...)
	raw.MessageSuffix = append(raw.MessageSuffix, suffixes...)

	return nil
}

func readImportedTemplateFiles(paths []string, extraVars map[string]string) (map[string]string, []ImportedTemplateItem) {
	vars := make(map[string]string)
	var items []ImportedTemplateItem
	for _, p := range paths {
		items = appendSingleTemplateFile(p, vars, items)
	}
	for k, v := range extraVars {
		vars[strings.TrimPrefix(k, "$")] = v
	}

	return vars, items
}

func appendSingleTemplateFile(path string, vars map[string]string, items []ImportedTemplateItem) []ImportedTemplateItem {
	data, err := os.ReadFile(path)
	if err != nil {
		return items
	}
	var tf ImportedTemplateFile
	if json.Unmarshal(data, &tf) != nil {
		return items
	}
	for k, v := range tf.Variables {
		vars[strings.TrimPrefix(k, "$")] = v
	}

	return appendCategoryAndFlatItems(items, tf)
}

func appendCategoryAndFlatItems(items []ImportedTemplateItem, tf ImportedTemplateFile) []ImportedTemplateItem {
	for _, cat := range tf.Categories {
		for _, it := range cat.Items {
			if it.Category == "" {
				it.Category = cat.Slug
			}
			items = append(items, it)
		}
	}

	return append(items, tf.Templates...)
}

func resolveTemplateRefs(refs []string, imported []ImportedTemplateItem, vars map[string]string) []string {
	var out []string
	for _, ref := range refs {
		out = append(out, resolveSingleTemplateRef(ref, imported, vars)...)
	}

	return out
}

func resolveSingleTemplateRef(ref string, imported []ImportedTemplateItem, vars map[string]string) []string {
	if matched := matchImportedByRef(ref, imported, vars); len(matched) > 0 {
		return matched
	}
	if PrecompileStateCategoryHook != nil {
		if dbMatched := PrecompileStateCategoryHook(ref, vars); len(dbMatched) > 0 {
			return dbMatched
		}
	}

	return []string{expandStaticVars(ref, vars)}
}

func matchImportedByRef(ref string, imported []ImportedTemplateItem, vars map[string]string) []string {
	var out []string
	lowerRef := strings.ToLower(strings.TrimSpace(ref))
	for _, item := range imported {
		isMatch := strings.ToLower(item.Category) == lowerRef ||
			strings.ToLower(item.SubCategory) == lowerRef ||
			strings.ToLower(item.Slug) == lowerRef ||
			strings.ToLower(item.ID) == lowerRef ||
			lowerRef == "all"
		if isMatch {
			out = append(out, formatCompiledItem(item, vars))
		}
	}

	return out
}

func precompileAllImportedItems(imported []ImportedTemplateItem, vars map[string]string) []string {
	out := make([]string, 0, len(imported))
	for _, item := range imported {
		out = append(out, formatCompiledItem(item, vars))
	}

	return out
}

func formatCompiledItem(item ImportedTemplateItem, vars map[string]string) string {
	title := strings.TrimSpace(expandStaticVars(item.Title, vars))
	text := strings.TrimSpace(expandStaticVars(item.Text, vars))
	if strings.HasPrefix(text, "#") || title == "" {
		return text
	}
	cleanTitle := strings.TrimSpace(strings.TrimLeft(title, "#"))

	return "# " + cleanTitle + "\n" + text
}

func expandStaticVars(input string, vars map[string]string) string {
	if !strings.Contains(input, "$") || len(vars) == 0 {
		return input
	}
	return configVarRegex.ReplaceAllStringFunc(input, func(m string) string {
		key := strings.TrimPrefix(m, "$")
		key = strings.TrimPrefix(key, "{")
		key = strings.TrimSuffix(key, "}")
		if val, ok := vars[key]; ok {
			return val
		}
		return m
	})
}

func expandConfigVariables(cfg CommitInConfigJSON) CommitInConfigJSON {
	if len(cfg.Variables) == 0 {
		return cfg
	}
	vars := resolveNestedVariables(cfg.Variables)
	cfg.Target = expandStaticVars(cfg.Target, vars)
	cfg.SummaryDir = expandStaticVars(cfg.SummaryDir, vars)
	for i, in := range cfg.Inputs {
		cfg.Inputs[i] = expandStaticVars(in, vars)
	}
	for i, imp := range cfg.Imports {
		cfg.Imports[i] = expandStaticVars(imp, vars)
	}
	for i, tr := range cfg.TitleReplacements {
		cfg.TitleReplacements[i].Replacement = expandStaticVars(tr.Replacement, vars)
	}

	return cfg
}

func resolveNestedVariables(in map[string]string) map[string]string {
	out := make(map[string]string, len(in))
	for k, v := range in {
		out[k] = v
	}
	for pass := 0; pass < 3; pass++ {
		for k, v := range out {
			out[k] = expandStaticVars(v, out)
		}
	}

	return out
}
