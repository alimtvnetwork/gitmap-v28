package helptext

import (
	"embed"
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/render"
)

//go:embed *.md
var files embed.FS

// Print reads and prints the help file for the given command using the
// default PrettyAuto mode (TTY auto-detect + GITMAP_NO_PRETTY opt-out).
// Kept as a thin wrapper for callers that don't parse a --pretty flag.
func Print(command string) {
	PrintWithMode(command, render.PrettyAuto)
}

// PrintWithMode reads and prints the help file for `command`, routing the
// markdown through render.RenderANSI when render.Decide says so for the
// caller-supplied PrettyMode. This is the preferred entry point for
// command surfaces that parse --pretty / --no-pretty so user intent
// flows all the way to the renderer.
//
// The decision is delegated to render.Decide so help, templates show,
// and changelog all answer "should I emit ANSI?" identically.
func PrintWithMode(command string, mode render.PrettyModeType) {
	data, err := ReadRaw(command)
	if err != nil {
		appErr := apperror.WrapWithDetails(
			err,
			"helptext.PrintWithMode",
			"E1071",
			fmt.Sprintf("No help available for '%s'", command),
			"helptext",
			apperror.ErrorTypeValidation,
			apperror.SeverityError,
			map[string]any{"command": command},
		)
		cliexit.HandleError(appErr, 1)

		return
	}

	if render.Decide(mode, render.StdoutIsTerminal(), true) {
		fmt.Print(render.RenderANSI(string(data)))

		return
	}

	fmt.Print(string(data))
}

// PrintRaw bypasses the pretty renderer and prints the embedded
// markdown verbatim. Useful for callers that pipe help into a pager
// or other tooling that handles its own formatting. Equivalent to
// PrintWithMode(command, render.PrettyOff) but spelled out for clarity
// at call sites that always want raw output regardless of TTY state.
func PrintRaw(command string) {
	data, err := ReadRaw(command)
	if err != nil {
		appErr := apperror.WrapWithDetails(
			err,
			"helptext.PrintRaw",
			"E1072",
			fmt.Sprintf("No help available for '%s'", command),
			"helptext",
			apperror.ErrorTypeValidation,
			apperror.SeverityError,
			map[string]any{"command": command},
		)
		cliexit.HandleError(appErr, 1)

		return
	}

	fmt.Print(string(data))
}

var helpAliases = map[string]string{
	"cpi-all":                  "import-all",
	"all-profile-import":       "import-all",
	"import-all-profiles":      "import-all",
	"cpe-all":                  "export-all",
	"all-profile-export":       "export-all",
	"export-all-profiles":      "export-all",
	"cpc-all":                  "copy-all",
	"all-profile-copy":         "copy-all",
	"copy-all-profiles":        "copy-all",
	"inspect":                  "import-check",
	"import-ls":                "import-check",
	"check":                    "import-check",
	"vm":                       "vmware",
	"ngx":                      "nginx",
	"sj":                       "ssh-join",
	"servers-client":           "servers-clients",
	"client":                   "clients",
	"cluster-ls":               "cluster-nodes",
	"cluster-node-add":         "cluster-add",
	"cluster-node-rm":          "cluster-remove",
	"cluster-node-remove":      "cluster-remove",
	"cluster-rm":               "cluster-remove",
	"cluster-ping":             "cluster-status",
	"cluster-run":              "cluster-exec",
	"cluster-script":           "cluster-run-script",
	"cluster-bs":               "cluster-bootstrap",
	"cluster-kube":             "cluster-k8s",
	"cluster-kubernetes":       "cluster-k8s",
	"cluster-in":               "cluster-install",
	"sj-install":               "cluster-install",
	"schedule-shutdown-status": "schedule-shutdown",
	"schedule-shutdown-cancel": "schedule-shutdown",
	"schedule-restart-status":  "schedule-restart",
	"schedule-restart-cancel":  "schedule-restart",
	"os-clean":                 "os-ai-clean",
	"ai-clean":                 "os-ai-clean",
	"fix-pipeline":             "agy-fix-pipeline",
	"agy-fix":                  "agy-fix-pipeline",
	"agy-fp":                   "agy-fix-pipeline",
	"tasks":                    "task",
	"tk":                       "task",
	"com":                      "clone-only-missing",
	"prompt_templates":         "prompts-template",
	"prompt-templates":         "prompts-template",
	"prompt-template":          "prompts-template",
	"prompts-templates":        "prompts-template",
	"pt":                       "prompts-template",
	"pipe":                     "pipeline",
	"pl":                       "pipeline",
	"disk":                     "storage",
	"df":                       "storage",
	"ag":                       "agy",
	"antigravity":              "agy",
	"scripts":                  "ai",
	"ai-create":                "ai",
	"ai-new":                   "ai",
	"ai-scaffold":              "ai",
	"ai-list":                  "ai",
	"ai-run":                   "ai",
	"ai-fix":                   "ai",
	"scripts-create":           "ai",
	"scripts-new":              "ai",
	"scripts-list":             "ai",
	"scripts-run":              "ai",
	"scripts-fix":              "ai",
	"auto":                     "automation",
	"aum":                      "automation",
	"py-auto":                  "automation",
	"automation-search":        "automation",
	"automation-newlines":      "automation",
	"automation-cache":          "automation",
	"automation-benchmark":      "automation",
	"automation-relative-paths": "automation",
	"automation-naming":         "automation",
	"automation-result-wrapper": "automation",
	"automation-params":         "automation",
	"automation-enums":          "automation",
	"relative-paths":            "automation",
	"rel-paths":                 "automation",
	"naming":                    "automation",
	"result-wrapper":            "automation",
	"params":                    "automation",
	"enums":                       "automation",
	"automation-topology":         "automation",
	"topology":                    "automation",
	"topo":                        "automation",
	"codebase-topology":           "automation",
	"automation-db-generate":      "automation",
	"db-generate":                 "automation",
	"db-gen":                      "automation",
	"gen-db":                      "automation",
	"automation-db-migrate":       "automation",
	"db-migrate":                  "automation",
	"db-up":                       "automation",
	"automation-schema-audit":     "automation",
	"schema-audit":                "automation",
	"db-audit":                    "automation",
	"audit-schema":                "automation",
	"automation-preflight":        "automation",
	"preflight":                   "automation",
	"check-all":                   "automation",
	"ci-local":                    "automation",
	"automation-test-inventory":   "automation",
	"test-inventory":              "automation",
	"tests-inv":                   "automation",
	"automation-purge-actions":    "automation",
	"purge-actions":               "automation",
	"purge-artifacts":             "automation",
	"clean-actions":               "automation",
	"automation-smoke-test":       "automation",
	"smoke-test":                  "automation",
	"smoke":                       "automation",
	"installer-smoke":             "automation",
	"automation-version-sync":     "automation",
	"version-sync":                "automation",
	"sync-version":                "automation",
	"ver-sync":                    "automation",
	"automation-release-bump":     "automation",
	"release-bump":                "automation",
	"semver-bump":                 "automation",
	"automation-milestones":       "automation",
	"milestones":                  "automation",
	"milestone":                   "automation",
	"automation-help-audit":       "automation",
	"help-audit":                  "automation",
	"help-check":                  "automation",
	"doc-audit":                   "automation",
	"automation-plan-consolidate": "automation",
	"plan-consolidate":            "automation",
	"consolidate-plans":           "automation",
	"automation-doc-links":        "automation",
	"doc-links":                   "automation",
	"check-links":                 "automation",
	"doc-paths":                   "automation",
	"automation-spec-migrate":     "automation",
	"spec-migrate":                "automation",
	"resequence-spec":             "automation",
	"migrate-spec":                "automation",
	"automation-clean-artifacts":  "automation",
	"clean-artifacts":             "automation",
	"clean-build":                 "automation",
	"rm-artifacts":                "automation",
	"automation-changed-files":    "automation",
	"changed-files":               "automation",
	"git-changes":                 "automation",
	"diff-files":                  "automation",
	"automation-purge-history":    "automation",
	"purge-history":               "automation",
	"trace-history":               "automation",
	"clean-history":               "automation",
	"automation-format-go":        "automation",
	"format-go":                   "automation",
	"gofmt-ast":                   "automation",
	"fmt-go":                      "automation",
}

func resolveHelpAlias(cmd string) string {
	if alias, hasAlias := helpAliases[strings.ToLower(cmd)]; hasAlias {
		return alias
	}

	return cmd
}

// ReadRaw returns the embedded help markdown for the given command
// without exiting the process on miss. Test-friendly counterpart to
// PrintRaw — callers that want to assert on help contents (alias
// coverage, link integrity, etc.) should use this so the test binary
// is not torn down by os.Exit when a lookup fails.
func ReadRaw(command string) ([]byte, error) {
	command = resolveHelpAlias(command)
	data, err := files.ReadFile(command + ".md")
	if err == nil {
		return data, nil
	}

	low := strings.ToLower(command)
	data, err = files.ReadFile(low + ".md")
	if err == nil {
		return data, nil
	}

	dashed := strings.ReplaceAll(low, " ", "-")

	return files.ReadFile(dashed + ".md")
}
