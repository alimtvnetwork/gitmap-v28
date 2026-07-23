package cmd

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// ctxExplainEnabled is set by runInstallCtx and read by the per-platform
// templating helpers. When true every generated leaf wraps its real
// command with a print of `> gitmap <args>` so users see the exact
// invocation before it runs. Process-local; never persisted to disk —
// only the rendered registry/script body is.
var ctxExplainEnabled bool

// runInstallCtx dispatches the right-click context-menu install to
// the platform-specific implementation. Spec: spec/04-generic-cli/30-install-ctx.md.
// When explain=true, generated entries print their resolved invocation
// before executing.
func runInstallCtx(explain bool) {
	ctxExplainEnabled = explain
	switch runtime.GOOS {
	case "windows":
		runInstallCtxWindows()
	case "darwin":
		runInstallCtxMac()
	case "linux":
		runInstallCtxLinux()
	default:
		fmt.Fprintf(os.Stderr, constants.MsgCtxOSUnsupported, runtime.GOOS)
	}
}

// runUninstallCtx dispatches removal to the platform implementation.
func runUninstallCtx() {
	switch runtime.GOOS {
	case "windows":
		runUninstallCtxWindows()
	case "darwin":
		runUninstallCtxMac()
	case "linux":
		runUninstallCtxLinux()
	default:
		fmt.Fprintf(os.Stderr, constants.MsgCtxOSUnsupported, runtime.GOOS)
	}
}

// runInstallCtxWindows writes the full nested HKCU registry tree.
func runInstallCtxWindows() {
	fmt.Print(constants.MsgCtxInstallStart)

	exe := resolveCtxExe()
	cmds := buildCtxInstallCommands(exe)

	successes := runRegistryCommandsCtx(cmds)
	fmt.Printf(constants.MsgCtxInstallDone, successes, len(cmds))
}

// runUninstallCtxWindows removes the gitmap subtree from both HKCU roots.
func runUninstallCtxWindows() {
	fmt.Print(constants.MsgCtxUninstallStart)

	cmds := [][]string{
		{"reg", "delete", constants.CtxRootKeyBackground, "/f"},
		{"reg", "delete", constants.CtxRootKeyDirectory, "/f"},
	}

	successes := runRegistryCommandsCtx(cmds)
	fmt.Printf(constants.MsgCtxUninstallDone, successes, len(cmds))
}

// resolveCtxExe returns the absolute path to the running gitmap binary.
func resolveCtxExe() string {
	exe, err := os.Executable()
	if err != nil || exe == "" {
		return "gitmap"
	}

	return exe
}

// buildCtxInstallCommands generates the full registry-write command set
// for both root keys (Background and Directory).
func buildCtxInstallCommands(exe string) [][]string {
	var out [][]string

	for _, root := range []string{constants.CtxRootKeyBackground, constants.CtxRootKeyDirectory} {
		out = append(out, rootCascadeCommands(root, exe)...)
		for _, e := range ctxMenu() {
			out = append(out, entryCommands(root, e, exe)...)
		}
	}

	return out
}

// rootCascadeCommands writes the top-level "gitmap ▸" cascade key.
func rootCascadeCommands(root, exe string) [][]string {
	return [][]string{
		{"reg", "add", root, "/ve", "/d", "", "/f"},
		{"reg", "add", root, "/v", "MUIVerb", "/d", constants.CtxRootMUIVerb, "/f"},
		{"reg", "add", root, "/v", "SubCommands", "/d", "", "/f"},
		{"reg", "add", root, "/v", "Icon", "/d", exe + ",0", "/f"},
	}
}

// entryCommands writes one menu entry — either a category cascade
// (with Children) or a leaf \command key.
func entryCommands(root string, e ctxEntry, exe string) [][]string {
	key := root + `\shell\` + e.KeyName

	if len(e.Children) > 0 {
		return categoryCommands(key, e, exe)
	}

	return leafCommands(key, e, exe)
}

// categoryCommands wires a sub-cascade key plus all its child leaves.
func categoryCommands(key string, e ctxEntry, exe string) [][]string {
	out := [][]string{
		{"reg", "add", key, "/ve", "/d", "", "/f"},
		{"reg", "add", key, "/v", "MUIVerb", "/d", e.MUIVerb, "/f"},
		{"reg", "add", key, "/v", "SubCommands", "/d", "", "/f"},
	}
	if icon := resolveCtxIcon(e.Icon, exe); icon != "" {
		out = append(out, []string{"reg", "add", key, "/v", "Icon", "/d", icon, "/f"})
	}
	for _, child := range e.Children {
		out = append(out, leafCommands(key+`\shell\`+child.KeyName, child, exe)...)
	}

	return out
}

// leafCommands wires one terminal/silent/prefill action key. When
// e.Extended is true an empty "Extended" REG_SZ value is written; the
// Windows shell hides those entries unless Shift is held during the
// right-click (the standard mechanism for power-user actions). When
// e.Icon is non-empty an Icon REG_SZ value is written (the menu
// renders that icon next to the entry); the constants.CtxIconExeToken
// placeholder is substituted with the resolved gitmap binary path.
func leafCommands(key string, e ctxEntry, exe string) [][]string {
	out := [][]string{
		{"reg", "add", key, "/ve", "/d", e.MUIVerb, "/f"},
		{"reg", "add", key + `\command`, "/ve", "/d", commandTemplate(e, exe), "/f"},
	}
	if icon := resolveCtxIcon(e.Icon, exe); icon != "" {
		out = append(out, []string{"reg", "add", key, "/v", "Icon", "/d", icon, "/f"})
	}
	if e.Extended {
		out = append(out, []string{"reg", "add", key, "/v", "Extended", "/d", "", "/f"})
	}

	return out
}

// resolveCtxIcon expands the constants.CtxIconExeToken placeholder
// inside an entry's declared Icon value. Returns "" untouched so
// callers can skip writing the Icon registry value entirely.
func resolveCtxIcon(icon, exe string) string {
	if icon == "" {
		return ""
	}

	return strings.ReplaceAll(icon, constants.CtxIconExeToken, exe)
}

// commandTemplate builds the pwsh invocation string baked into a
// \command key's (Default) value. %V is Explorer's clicked-folder token.
func commandTemplate(e ctxEntry, exe string) string {
	if e.Mode == constants.CtxModePrefill {
		return `pwsh -NoExit -NoProfile -Command "Set-Location '%V'; Write-Host -NoNewline 'gitmap '"`
	}

	target := exe
	if e.Exe != "" {
		target = e.Exe // resolved from PATH at runtime (e.g. "git")
	}
	args := strings.Join(e.Args, " ")
	prefix := ctxExplainPrefixPwsh(target, e.Args)
	if e.Mode == constants.CtxModeSilent {
		return fmt.Sprintf(`pwsh -NoProfile -WindowStyle Hidden -Command "Set-Location '%%V'; %s& '%s' %s 2>&1 | Out-String | %% { msg.exe * $_ }"`, prefix, target, args)
	}

	return fmt.Sprintf(`pwsh -NoExit -NoProfile -Command "Set-Location '%%V'; %s& '%s' %s"`, prefix, target, args)
}
