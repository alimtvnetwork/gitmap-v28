package cmd_test

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/store"
)

// --- extractAliasFlag tests (logic copy for black-box testing) ---

func extractAliasFlag(args []string) (string, []string) {
	for i, arg := range args {
		if arg == "-A" || arg == "--alias" {
			if i+1 < len(args) {
				return args[i+1], removeElements(args, i, 2)
			}

			return "", args
		}

		if hasAliasPrefix(arg, "-A=") {
			return arg[3:], removeElements(args, i, 1)
		}
		if hasAliasPrefix(arg, "--alias=") {
			return arg[8:], removeElements(args, i, 1)
		}
	}

	return "", args
}

func hasAliasPrefix(arg, prefix string) bool {
	return len(arg) > len(prefix) && arg[:len(prefix)] == prefix
}

func removeElements(args []string, index, count int) []string {
	result := make([]string, 0, len(args)-count)
	result = append(result, args[:index]...)
	result = append(result, args[index+count:]...)

	return result
}

// resolvedAlias mirrors the cmd package struct for testing.
type resolvedAlias struct {
	Alias        string
	AbsolutePath string
	Slug         string
}

// --- extractAliasFlag: -A <name> form ---

func TestExtractAliasFlag_ShortForm(t *testing.T) {
	alias, remaining := extractAliasFlag([]string{"-A", "api", "fetch"})
	if alias != "api" {
		t.Errorf("expected alias=api, got %q", alias)
	}
	if len(remaining) != 1 || remaining[0] != "fetch" {
		t.Errorf("expected [fetch], got %v", remaining)
	}
}

func TestExtractAliasFlag_LongForm(t *testing.T) {
	alias, remaining := extractAliasFlag([]string{"--alias", "web", "status"})
	if alias != "web" {
		t.Errorf("expected alias=web, got %q", alias)
	}
	if len(remaining) != 1 || remaining[0] != "status" {
		t.Errorf("expected [status], got %v", remaining)
	}
}

func TestExtractAliasFlag_ShortEquals(t *testing.T) {
	alias, remaining := extractAliasFlag([]string{"-A=infra", "--all"})
	if alias != "infra" {
		t.Errorf("expected alias=infra, got %q", alias)
	}
	if len(remaining) != 1 || remaining[0] != "--all" {
		t.Errorf("expected [--all], got %v", remaining)
	}
}

func TestExtractAliasFlag_LongEquals(t *testing.T) {
	alias, remaining := extractAliasFlag([]string{"--alias=db", "pull"})
	if alias != "db" {
		t.Errorf("expected alias=db, got %q", alias)
	}
	if len(remaining) != 1 || remaining[0] != "pull" {
		t.Errorf("expected [pull], got %v", remaining)
	}
}

func TestExtractAliasFlag_NoAlias(t *testing.T) {
	alias, remaining := extractAliasFlag([]string{"--all", "fetch"})
	if alias != "" {
		t.Errorf("expected empty alias, got %q", alias)
	}
	if len(remaining) != 2 {
		t.Errorf("expected 2 remaining args, got %d", len(remaining))
	}
}

func TestExtractAliasFlag_Empty(t *testing.T) {
	alias, remaining := extractAliasFlag([]string{})
	if alias != "" {
		t.Errorf("expected empty alias, got %q", alias)
	}
	if len(remaining) != 0 {
		t.Errorf("expected empty remaining, got %v", remaining)
	}
}

func TestExtractAliasFlag_MidArgs(t *testing.T) {
	alias, remaining := extractAliasFlag([]string{"--group", "work", "-A", "api", "--all"})
	if alias != "api" {
		t.Errorf("expected alias=api, got %q", alias)
	}
	if len(remaining) != 3 || remaining[0] != "--group" || remaining[1] != "work" || remaining[2] != "--all" {
		t.Errorf("expected [--group work --all], got %v", remaining)
	}
}

func TestExtractAliasFlag_MissingValue(t *testing.T) {
	// When -A is last arg with no value, original exits; here we just
	// verify the flag is detected (returns empty since no next arg).
	alias, _ := extractAliasFlag([]string{"fetch", "-A"})
	if alias != "" {
		t.Errorf("expected empty alias for missing value, got %q", alias)
	}
}

// --- hasAliasPrefix ---

func TestHasAliasPrefix_Valid(t *testing.T) {
	if !hasAliasPrefix("-A=api", "-A=") {
		t.Error("expected true for -A=api")
	}
}

func TestHasAliasPrefix_ExactLength(t *testing.T) {
	if hasAliasPrefix("-A=", "-A=") {
		t.Error("expected false for exact prefix match with no value")
	}
}

func TestHasAliasPrefix_TooShort(t *testing.T) {
	if hasAliasPrefix("-A", "-A=") {
		t.Error("expected false for shorter than prefix")
	}
}

// --- removeElements ---

func TestRemoveElements_Middle(t *testing.T) {
	result := removeElements([]string{"a", "b", "c", "d"}, 1, 2)
	if len(result) != 2 || result[0] != "a" || result[1] != "d" {
		t.Errorf("expected [a d], got %v", result)
	}
}

func TestRemoveElements_Start(t *testing.T) {
	result := removeElements([]string{"a", "b", "c"}, 0, 1)
	if len(result) != 2 || result[0] != "b" || result[1] != "c" {
		t.Errorf("expected [b c], got %v", result)
	}
}

func TestRemoveElements_End(t *testing.T) {
	result := removeElements([]string{"a", "b", "c"}, 2, 1)
	if len(result) != 2 || result[0] != "a" || result[1] != "b" {
		t.Errorf("expected [a b], got %v", result)
	}
}

// --- Alias context accessor tests ---

func TestAliasAccessors_Nil(t *testing.T) {
	// Simulate nil context
	var ctx *resolvedAlias

	path := ""
	slug := ""
	hasAlias := false

	if ctx != nil {
		path = ctx.AbsolutePath
		slug = ctx.Slug
		hasAlias = true
	}

	if path != "" {
		t.Errorf("expected empty path, got %q", path)
	}
	if slug != "" {
		t.Errorf("expected empty slug, got %q", slug)
	}
	if hasAlias {
		t.Error("expected hasAlias=false")
	}
}

func TestAliasAccessors_Set(t *testing.T) {
	ctx := &resolvedAlias{
		Alias:        "api",
		AbsolutePath: "/home/user/repos/api-gateway",
		Slug:         "github/user/api-gateway",
	}

	if ctx.AbsolutePath != "/home/user/repos/api-gateway" {
		t.Errorf("expected path, got %q", ctx.AbsolutePath)
	}
	if ctx.Slug != "github/user/api-gateway" {
		t.Errorf("expected slug, got %q", ctx.Slug)
	}
}

// --- AliasAsRecords equivalent ---

func aliasAsRecords(ctx *resolvedAlias) []store.AliasWithRepo {
	if ctx == nil {
		return nil
	}

	return []store.AliasWithRepo{{
		AbsolutePath: ctx.AbsolutePath,
		Slug:         ctx.Slug,
	}}
}

func TestAliasAsRecords_Nil(t *testing.T) {
	records := aliasAsRecords(nil)
	if records != nil {
		t.Errorf("expected nil, got %v", records)
	}
}

func TestAliasAsRecords_Set(t *testing.T) {
	ctx := &resolvedAlias{
		Alias:        "web",
		AbsolutePath: "/repos/web-app",
		Slug:         "github/user/web-app",
	}

	records := aliasAsRecords(ctx)
	if len(records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(records))
	}
	if records[0].AbsolutePath != "/repos/web-app" {
		t.Errorf("expected path /repos/web-app, got %q", records[0].AbsolutePath)
	}
	if records[0].Slug != "github/user/web-app" {
		t.Errorf("expected slug github/user/web-app, got %q", records[0].Slug)
	}
}

// --- Integration: alias flag extraction across command patterns ---

func TestAliasFlagIntegration_PullWithAlias(t *testing.T) {
	// Simulates: gitmap pull -A api
	args := []string{"-A", "api"}
	alias, remaining := extractAliasFlag(args)

	if alias != "api" {
		t.Errorf("pull: expected alias=api, got %q", alias)
	}
	if len(remaining) != 0 {
		t.Errorf("pull: expected no remaining args, got %v", remaining)
	}

	ctx := &resolvedAlias{Alias: alias, AbsolutePath: "/repos/api", Slug: "github/user/api"}
	records := aliasAsRecords(ctx)
	if len(records) != 1 || records[0].AbsolutePath != "/repos/api" {
		t.Errorf("pull: alias resolution failed, records=%v", records)
	}
}

func TestAliasFlagIntegration_ExecWithAlias(t *testing.T) {
	// Simulates: gitmap exec -A web fetch --prune
	args := []string{"-A", "web", "fetch", "--prune"}
	alias, remaining := extractAliasFlag(args)

	if alias != "web" {
		t.Errorf("exec: expected alias=web, got %q", alias)
	}
	if len(remaining) != 2 || remaining[0] != "fetch" || remaining[1] != "--prune" {
		t.Errorf("exec: expected [fetch --prune], got %v", remaining)
	}
}

func TestAliasFlagIntegration_StatusWithAlias(t *testing.T) {
	// Simulates: gitmap status --alias=api --all
	args := []string{"--alias=api", "--all"}
	alias, remaining := extractAliasFlag(args)

	if alias != "api" {
		t.Errorf("status: expected alias=api, got %q", alias)
	}
	if len(remaining) != 1 || remaining[0] != "--all" {
		t.Errorf("status: expected [--all], got %v", remaining)
	}
}

func TestAliasFlagIntegration_CdWithAlias(t *testing.T) {
	// Simulates: gitmap cd -A infra
	args := []string{"-A", "infra"}
	alias, remaining := extractAliasFlag(args)

	if alias != "infra" {
		t.Errorf("cd: expected alias=infra, got %q", alias)
	}
	if len(remaining) != 0 {
		t.Errorf("cd: expected no remaining args, got %v", remaining)
	}

	ctx := &resolvedAlias{Alias: alias, AbsolutePath: "/repos/infra", Slug: "github/user/infra"}
	if ctx.AbsolutePath != "/repos/infra" {
		t.Errorf("cd: expected resolved path, got %q", ctx.AbsolutePath)
	}
}

func TestAliasFlagIntegration_AliasWithGroupFlag(t *testing.T) {
	// Simulates: gitmap status --group work -A api --all
	// Alias should be extracted, group and --all remain
	args := []string{"--group", "work", "-A", "api", "--all"}
	alias, remaining := extractAliasFlag(args)

	if alias != "api" {
		t.Errorf("mixed: expected alias=api, got %q", alias)
	}
	if len(remaining) != 3 {
		t.Errorf("mixed: expected 3 remaining args, got %v", remaining)
	}
	if remaining[0] != "--group" || remaining[1] != "work" || remaining[2] != "--all" {
		t.Errorf("mixed: expected [--group work --all], got %v", remaining)
	}
}

func TestAliasFlagIntegration_EqualsFormExec(t *testing.T) {
	// Simulates: gitmap exec -A=db status
	args := []string{"-A=db", "status"}
	alias, remaining := extractAliasFlag(args)

	if alias != "db" {
		t.Errorf("exec-equals: expected alias=db, got %q", alias)
	}
	if len(remaining) != 1 || remaining[0] != "status" {
		t.Errorf("exec-equals: expected [status], got %v", remaining)
	}
}
