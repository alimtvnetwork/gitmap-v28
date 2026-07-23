package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// runInstallCtxMac generates one Automator Quick Action .workflow
// bundle per flat menu entry under ~/Library/Services. macOS shows
// these in Finder's right-click "Quick Actions" / "Services" submenu.
func runInstallCtxMac() {
	fmt.Print(constants.MsgCtxMacInstallStart)

	dir, err := macServicesDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.MsgCtxFsWriteFail, "$HOME", err)

		return
	}

	exe := resolveCtxExe()
	flat := flattenCtxMenu()
	ok := 0
	for _, e := range flat {
		if writeMacWorkflow(dir, e, exe) {
			ok++
		}
	}
	fmt.Printf(constants.MsgCtxMacInstallDone, ok, len(flat))
}

// runUninstallCtxMac removes every gitmap-prefixed .workflow bundle
// previously written by runInstallCtxMac. Other Services are left alone.
func runUninstallCtxMac() {
	fmt.Print(constants.MsgCtxMacUninstallStart)

	dir, err := macServicesDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.MsgCtxFsWriteFail, "$HOME", err)

		return
	}

	flat := flattenCtxMenu()
	ok := 0
	for _, e := range flat {
		path := filepath.Join(dir, e.Slug+".workflow")
		if err := os.RemoveAll(path); err != nil {
			fmt.Fprintf(os.Stderr, constants.MsgCtxFsRmFail, path, err)

			continue
		}
		ok++
	}
	fmt.Printf(constants.MsgCtxMacUninstallDone, ok, len(flat))
}

// macServicesDir returns ~/Library/Services, creating it if missing.
func macServicesDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, constants.CtxMacServicesRel)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}

	return dir, nil
}

// writeMacWorkflow creates one <slug>.workflow bundle with the minimal
// Info.plist + document.wflow needed for Finder to surface the entry.
// Returns true on success. Diagnostics go to os.Stderr per zero-swallow
// policy.
func writeMacWorkflow(dir string, e flatCtxEntry, exe string) bool {
	bundle := filepath.Join(dir, e.Slug+".workflow")
	contents := filepath.Join(bundle, "Contents")
	if err := os.MkdirAll(contents, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, constants.MsgCtxFsWriteFail, contents, err)

		return false
	}

	if !writeFileCtx(filepath.Join(contents, "Info.plist"), macInfoPlist(e.Label)) {
		return false
	}
	shell := macShellFor(e, exe)
	if !writeFileCtx(filepath.Join(contents, "document.wflow"), macDocumentWflow(shell)) {
		return false
	}

	return true
}

// writeFileCtx is a small wrapper that logs and returns false on error.
func writeFileCtx(path, body string) bool {
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, constants.MsgCtxFsWriteFail, path, err)

		return false
	}

	return true
}

// macShellFor returns the shell command Automator runs ($1 = clicked
// folder path, passed as argument because inputMethod=1).
func macShellFor(e flatCtxEntry, exe string) string {
	args := strings.Join(e.Args, " ")
	target := exe
	if e.Exe != "" {
		target = e.Exe
	}
	switch e.Mode {
	case constants.CtxModePrefill:
		return `osascript -e 'tell application "Terminal" to do script "cd \"'"$1"'\" && printf \"gitmap \""' -e 'tell application "Terminal" to activate'`
	case constants.CtxModeSilent:
		announce := ctxExplainAnnounce(target, e.Args)

		return fmt.Sprintf(`cd "$1" && OUT=$(printf %%s '%s'; '%s' %s 2>&1); osascript -e "display notification \"$(echo \"$OUT\" | head -c 200)\" with title \"%s\""`, announce, target, args, e.Label)
	default:
		// Embed the explain echo as the first command in the Terminal
		// session so the user sees the resolved invocation above the
		// real output. echo runs inside the new shell, not the Automator
		// host, so quoting stays AppleScript-safe.
		echoPrefix := ""
		if ctxExplainEnabled {
			echoPrefix = fmt.Sprintf(`echo \"> %s %s\" && `, target, args)
		}
		open := fmt.Sprintf(`osascript -e 'tell application "Terminal" to do script "cd \"'"$1"'\" && %s'"'"'%s'"'"' %s"' -e 'tell application "Terminal" to activate'`, echoPrefix, target, args)
		if e.Extended {
			// Power-user fan-out: confirm before running. macOS lacks
			// Shift-filter for Quick Actions, so we gate at the script.
			guard := fmt.Sprintf(`osascript -e 'display dialog "Run %s on every tracked repo? This is a power-user batch action." buttons {"Cancel","Run"} default button "Cancel"' >/dev/null 2>&1 || exit 0`, e.Label)

			return guard + " && " + open
		}

		return open
	}
}

// macInfoPlist returns the minimal Services Info.plist that registers
// one menu entry under Finder's Quick Actions for folders.
func macInfoPlist(label string) string {
	return `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0"><dict>
<key>NSServices</key><array><dict>
<key>NSMenuItem</key><dict><key>default</key><string>` + label + `</string></dict>
<key>NSMessage</key><string>runWorkflowAsService</string>
<key>NSSendFileTypes</key><array><string>public.folder</string></array>
</dict></array></dict></plist>
`
}

// macDocumentWflow returns the Automator document.wflow body that runs
// the given shell command with the clicked folder path as $1.
func macDocumentWflow(shell string) string {
	esc := strings.ReplaceAll(shell, "&", "&amp;")
	esc = strings.ReplaceAll(esc, "<", "&lt;")

	return `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0"><dict>
<key>AMApplicationVersion</key><string>2.10</string>
<key>AMDocumentVersion</key><integer>2</integer>
<key>actions</key><array><dict><key>action</key><dict>
<key>ActionBundlePath</key><string>/System/Library/Automator/Run Shell Script.action</string>
<key>ActionName</key><string>Run Shell Script</string>
<key>ActionParameters</key><dict>
<key>COMMAND_STRING</key><string>` + esc + `</string>
<key>CheckedForUserDefaultShell</key><true/>
<key>inputMethod</key><integer>1</integer>
<key>shell</key><string>/bin/bash</string>
<key>source</key><string></string>
</dict></dict></dict></array></dict></plist>
`
}
