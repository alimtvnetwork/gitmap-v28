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
