package cmdai

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

var masterCatalog = []ScriptMetadata{
	newScript("01", "01-index.md", "index", []string{"catalog", "docs"}, CategoryCore, "Master index and script catalog guide", false),
	newScript("02", "02-shared-engine.py", "shared-engine", []string{"engine", "shared"}, CategoryCore, "Shared constants, regex registry, and locks", false),
	newScript("03", "03-file-manipulator.py", "file-manipulator", []string{"manipulator", "rename"}, CategoryFormatting, "Mass lowercasing, sequence fixing, UTF-8 LF normalization", true),
	newScript("04", "04-newline-fixer.py", "newline-fixer", []string{"newlines", "whitespace"}, CategoryFormatting, "Fixes trailing whitespace and missing final newlines", true),
	newScript("05", "05-guideline-autofixer.py", "guideline-autofixer", []string{"guidelines", "autofixer"}, CategoryGuideline, "Composite runner for newlines and boolean conventions", true),
	newScript("06", "06-cicd-local-runner.py", "cicd-local-runner", []string{"cicd", "runner", "ci"}, CategoryCicd, "Runs all 18 CI quality checks in parallel", false),
	newScript("07", "07-relative-path-fixer.py", "relative-path-fixer", []string{"paths", "relative-paths"}, CategoryPaths, "Detects and fixes absolute paths and file URIs", true),
	newScript("08", "08-naming-autofixer.py", "naming-autofixer", []string{"naming", "booleans"}, CategoryGuideline, "Enforces lowercase filenames and boolean naming", true),
	newScript("09", "09-cli-help-auditor.py", "cli-help-auditor", []string{"help-auditor", "cli-help"}, CategoryAudit, "Validates CLI help examples against implementations", false),
	newScript("10", "10-encoding-normalizer.py", "encoding-normalizer", []string{"encoding", "utf8"}, CategoryFormatting, "Normalizes files to strict UTF-8 with UNIX LF endings", true),
	newScript("11", "11-fast-file-scanner.py", "fast-file-scanner", []string{"scanner", "files"}, CategoryDiscovery, "High-speed repo file scanner with cache query", false),
	newScript("12", "12-fast-cached-grep.py", "fast-cached-grep", []string{"grep", "search"}, CategoryDiscovery, "Parallel regex matcher leveraging pre-warmed file cache", false),
	newScript("13", "13-file-size-guard.py", "file-size-guard", []string{"file-size", "blob-guard"}, CategoryAudit, "Audits repository files for oversized binary blobs", false),
	newScript("14", "14-version-sync-checker.py", "version-sync-checker", []string{"version-sync", "version"}, CategoryAudit, "Verifies synchronization of version files and changelog", false),
	newScript("15", "15-sequence-and-title-auditor.py", "sequence-and-title-auditor", []string{"sequence-title", "headers"}, CategoryAudit, "Audits numeric file prefixes and H1 titles", false),
	newScript("16", "16-installer-smoke-tester.py", "installer-smoke-tester", []string{"installer-smoke", "smoke-test"}, CategoryMaintenance, "Installer smoke test validating placeholders and hashes", false),
	newScript("17", "17-fast-file-reader.py", "fast-file-reader", []string{"reader", "fast-read"}, CategoryDiscovery, "Fast file reader and folder explorer via cache", false),
	newScript("18", "18-codebase-topology-discoverer.py", "codebase-topology-discoverer", []string{"topology", "discoverer"}, CategoryDiscovery, "Universal polyglot codebase and topology discovery", false),
	newScript("19", "19-artifact-remover.py", "artifact-remover", []string{"artifacts", "clean-artifacts"}, CategoryMaintenance, "Safe interactive artifact remover with git untracking", true),
	newScript("20", "20-plan-consolidator.py", "plan-consolidator", []string{"plans", "consolidate-plans"}, CategoryPlans, "Fast Lovable plans and subtasks consolidator", true),
	newScript("21", "21-sequence-integrity-linter.py", "sequence-integrity-linter", []string{"sequence-integrity", "sequence-linter"}, CategoryAudit, "Verifies numeric sequences and headers across plans", false),
	newScript("22", "22-doc-path-linter.py", "doc-path-linter", []string{"doc-paths", "path-linter"}, CategoryPaths, "Lints markdown references and verifies doc paths", false),
	newScript("23", "23-coding-guideline-path-consolidator.py", "coding-guideline-path-consolidator", []string{"guideline-paths"}, CategoryGuideline, "Consolidates coding guideline references to specs", true),
	newScript("24", "24-spec-path-migrator.py", "spec-path-migrator", []string{"spec-paths", "spec-migrator"}, CategoryPaths, "Migrates legacy spec references to updated paths", true),
	newScript("25", "25-repo-migrator.py", "repo-migrator", []string{"migrator", "repo-migrate"}, CategoryMaintenance, "Repository-wide asset and structural migration utility", false),
	newScript("26", "26-go-code-formatter.py", "go-code-formatter", []string{"gofmt", "format-go"}, CategoryFormatting, "Cross-platform Go code formatter via gofmt", true),
	newScript("27a", "27-git-changed-files.py", "git-changed-files", []string{"changed-files", "diff-files"}, CategoryDiscovery, "Lists git changed and untracked files rapidly", false),
	newScript("27b", "27-misspell-auditor.py", "misspell-auditor", []string{"misspell", "spelling"}, CategoryAudit, "Audits and auto-fixes British to American English spelling", true),
	newScript("28", "28-go-preflight-ci.py", "go-preflight-ci", []string{"go-preflight", "preflight"}, CategoryCicd, "Runs local Go test and golangci-lint verification", false),
	newScript("29a", "29-release-bumper.py", "release-bumper", []string{"bumper", "bump-version"}, CategoryRelease, "Bumps repository version across all metadata files", true),
	newScript("29b", "29-release-orchestrator.py", "release-orchestrator", []string{"release", "orchestrator"}, CategoryRelease, "Automated release orchestrator, branch and tag manager", false),
	newScript("30a", "30-db-struct-enum-generator.py", "db-struct-enum-generator", []string{"db-enums", "enum-generator"}, CategoryDatabase, "Generates type-safe column name enums from Go structs", true),
	newScript("30b", "30-purge-history.py", "purge-history", []string{"purge-history", "git-purge"}, CategoryMaintenance, "Purges sensitive data and large blobs from git history", false),
	newScript("31a", "31-db-migration-runner.py", "db-migration-runner", []string{"db-migrations", "migrate-db"}, CategoryDatabase, "Standalone SQLite and schema migration runner", false),
	newScript("31b", "31-md-gap-fixer.py", "md-gap-fixer", []string{"md-gap", "gap-fixer"}, CategoryFormatting, "Fixes missing line gaps and headers in markdown files", true),
	newScript("32", "32-deep-consolidator.py", "deep-consolidator", []string{"deep-consolidate"}, CategoryPlans, "Deep plans and subtasks consolidator preserving detail", true),
	newScript("33a", "33-git-history-tracer-and-purger.py", "git-history-tracer-and-purger", []string{"tracer", "history-tracer"}, CategoryMaintenance, "Traces deleted files and workspace restoration", false),
	newScript("33b", "33-test-inventory-generator.py", "test-inventory-generator", []string{"test-inventory", "inventory"}, CategoryDatabase, "Generates test inventory and tracks file changes", true),
	newScript("34a", "34-purge-github-actions-artifacts.py", "purge-github-actions-artifacts", []string{"purge-artifacts", "gh-purge"}, CategoryMaintenance, "Purges GitHub Actions runner artifacts and logs", false),
	newScript("34b", "34-schema-scanner.py", "schema-scanner", []string{"schema-scanner", "schema"}, CategoryDatabase, "Scans SQL table definitions for conventions", false),
	newScript("35", "35-result-wrapper-auditor.py", "result-wrapper-auditor", []string{"result-auditor", "result-wrapper"}, CategoryAudit, "Audits Go functions for Result wrapper compliance", false),
	newScript("36", "36-param-struct-auditor.py", "param-struct-auditor", []string{"param-auditor", "params"}, CategoryGuideline, "Audits Go signatures for param structs and DTOs", false),
	newScript("37", "37-enum-guideline-auditor.py", "enum-guideline-auditor", []string{"enum-auditor", "enums"}, CategoryGuideline, "Audits enums for Type suffix and rune casts", false),
	newScript("38", "38-milestone-consolidator.py", "milestone-consolidator", []string{"milestone", "milestone-consolidate"}, CategoryPlans, "Consolidates completed plans into milestone summaries", true),
}

// MasterScriptCatalog returns all registered master automation scripts.
func MasterScriptCatalog() []ScriptMetadata {
	return masterCatalog
}

// AllScripts returns all registered AI scripts in the master catalog.
func AllScripts() []ScriptMetadata {
	return masterCatalog
}

// FindScriptByToken searches the catalog by numeric prefix, slug, filename, or alias.
func FindScriptByToken(token string) (ScriptMetadata, *apperror.AppError) {
	clean := strings.ToLower(strings.TrimSpace(token))
	for _, s := range masterCatalog {
		isMatch := matchesScriptToken(s, clean)
		if isMatch {
			return s, nil
		}
	}

	ctx := map[string]any{"token": token}

	return ScriptMetadata{}, apperror.New("find_script", "E_SCRIPT_NOT_FOUND", ctx)
}

func matchesScriptToken(s ScriptMetadata, token string) bool {
	isNameMatch := matchesNameOrNumber(s, token)
	if isNameMatch {
		return true
	}

	return matchesAlias(s.Aliases, token)
}

func matchesNameOrNumber(s ScriptMetadata, token string) bool {
	isExactNum := (s.Number == token || strings.TrimRight(s.Number, "ab") == token)
	isExactSlug := (s.Slug == token)
	isExactFile := (s.Filename == token)
	withoutExt := strings.TrimSuffix(s.Filename, ".py")

	return isExactNum || isExactSlug || isExactFile || withoutExt == token
}

func matchesAlias(aliases []string, token string) bool {
	for _, alias := range aliases {
		isMatch := (alias == token)
		if isMatch {
			return true
		}
	}

	return false
}

func defaultArgsFor(hasFix bool) []string {
	if hasFix {
		return []string{"--fix"}
	}

	return nil
}

func newScript(num, file, slug string, aliases []string, cat ScriptCategoryType, desc string, hasFix bool) ScriptMetadata {
	return ScriptMetadata{
		Number:      num,
		Filename:    file,
		Slug:        slug,
		Aliases:     aliases,
		Category:    cat,
		Description: desc,
		HasFixMode:  hasFix,
		DefaultArgs: defaultArgsFor(hasFix),
	}
}
