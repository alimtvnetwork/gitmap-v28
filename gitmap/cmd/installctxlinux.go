package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// runInstallCtxLinux installs the gitmap context menu into Nautilus,
// Dolphin and Thunar (the three file managers covering ~95% of Linux
// desktops). Each manager uses its own format; the flat menu list is
// the single source of truth.
func runInstallCtxLinux() {
	fmt.Print(constants.MsgCtxLinuxInstallStart)

	exe := resolveCtxExe()
	flat := flattenCtxMenu()

	ok := 0
	managers := 0
	if installNautilus(flat, exe) {
		managers++
		ok += len(flat)
	}
	if installDolphin(flat, exe) {
		managers++
		ok += len(flat)
	}
	if installThunar(flat, exe) {
		managers++
		ok += len(flat)
	}
	fmt.Printf(constants.MsgCtxLinuxInstallDone, ok, managers)
}

// runUninstallCtxLinux removes the previously installed entries from
// every supported file manager, leaving foreign entries untouched.
func runUninstallCtxLinux() {
	fmt.Print(constants.MsgCtxLinuxUninstallStart)

	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.MsgCtxFsRmFail, "$HOME", err)

		return
	}
	ok := 0
	ok += rmDirCtx(filepath.Join(home, constants.CtxLinuxNautilusRel))
	ok += rmFileCtx(filepath.Join(home, constants.CtxLinuxDolphinRel, constants.CtxLinuxDolphinFile))
	ok += stripThunarBlock(filepath.Join(home, constants.CtxLinuxThunarRel))
	fmt.Printf(constants.MsgCtxLinuxUninstallDone, ok)
}

// installNautilus drops one executable shell script per menu entry
// under ~/.local/share/nautilus/scripts/gitmap/. Nautilus uses the
// filename as the menu label, so we use the full flat label.
func installNautilus(flat []flatCtxEntry, exe string) bool {
	home, err := os.UserHomeDir()
	if err != nil {
		return false
	}
	dir := filepath.Join(home, constants.CtxLinuxNautilusRel)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, constants.MsgCtxFsWriteFail, dir, err)

		return false
	}
	for _, e := range flat {
		path := filepath.Join(dir, e.Label)
		if err := os.WriteFile(path, []byte(linuxShellScript(e, exe)), 0o755); err != nil {
			fmt.Fprintf(os.Stderr, constants.MsgCtxFsWriteFail, path, err)
		}
	}

	return true
}

// installDolphin writes a single .desktop service-menu file with one
// Action= per flat entry — KDE renders this as a cascading submenu.
func installDolphin(flat []flatCtxEntry, exe string) bool {
	home, err := os.UserHomeDir()
	if err != nil {
		return false
	}
	dir := filepath.Join(home, constants.CtxLinuxDolphinRel)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, constants.MsgCtxFsWriteFail, dir, err)

		return false
	}
	path := filepath.Join(dir, constants.CtxLinuxDolphinFile)
	body := dolphinDesktop(flat, exe)
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, constants.MsgCtxFsWriteFail, path, err)

		return false
	}

	return true
}

// linuxShellScript wraps one entry's command. $NAUTILUS_SCRIPT_CURRENT_URI
// is unset under Dolphin/Thunar — we fall back to the first arg.
func linuxShellScript(e flatCtxEntry, exe string) string {
	args := strings.Join(e.Args, " ")
	target := exe
	if e.Exe != "" {
		target = e.Exe
	}
	cd := `D="${1:-$PWD}"; cd "$D" || exit 1`
	guard := extendedGuard(e)
	echoSh := ctxExplainPrefixSh(target, e.Args)
	announce := ctxExplainAnnounce(target, e.Args)
	switch e.Mode {
	case constants.CtxModePrefill:
		return "#!/bin/sh\n" + cd + "\nx-terminal-emulator -e sh -c 'printf \"gitmap \"; exec $SHELL' &\n"
	case constants.CtxModeSilent:
		return fmt.Sprintf("#!/bin/sh\n%s\n%sOUT=$(printf %%s '%s'; '%s' %s 2>&1)\nnotify-send 'gitmap' \"$(echo \"$OUT\" | head -c 200)\" || echo \"$OUT\"\n", cd, guard, announce, target, args)
	default:
		return fmt.Sprintf("#!/bin/sh\n%s\n%sx-terminal-emulator -e sh -c \"%s'%s' %s; exec $SHELL\" &\n", cd, guard, echoSh, target, args)
	}
}

// extendedGuard returns a shell snippet (or empty string) that prompts
// for confirmation before continuing. zenity/kdialog/xmessage are
// tried in order; if none exist we fall back to a stdin prompt the
// terminal will surface. Used for power-user batch actions like
// pull-all that fan out across the entire catalog.
func extendedGuard(e flatCtxEntry) string {
	if !e.Extended {
		return ""
	}
	msg := fmt.Sprintf("Run %s on every tracked repo? This is a power-user batch action.", e.Label)

	return fmt.Sprintf(`{ zenity --question --text=%q 2>/dev/null || kdialog --yesno %q 2>/dev/null || xmessage -buttons Cancel:1,Run:0 %q 2>/dev/null; } || exit 0
`, msg, msg, msg)
}

// dolphinDesktop returns one .desktop body wiring all flat entries as
// Actions= under a cascading top-level "gitmap" service menu.
func dolphinDesktop(flat []flatCtxEntry, exe string) string {
	var actionIDs []string
	var sections strings.Builder
	for _, e := range flat {
		actionIDs = append(actionIDs, e.Slug)
		fmt.Fprintf(&sections, "[Desktop Action %s]\nName=%s\nExec=sh -c %q\n\n", e.Slug, e.Label, dolphinExec(e, exe))
	}

	return "[Desktop Entry]\nType=Service\nServiceTypes=KonqPopupMenu/Plugin\nMimeType=inode/directory;\nActions=" +
		strings.Join(actionIDs, ";") + ";\nX-KDE-Priority=TopLevel\nX-KDE-Submenu=gitmap\nIcon=utilities-terminal\n\n" +
		sections.String()
}

// dolphinExec returns the Exec= command string. %f expands to the
// clicked folder per the .desktop spec.
func dolphinExec(e flatCtxEntry, exe string) string {
	args := strings.Join(e.Args, " ")
	target := exe
	if e.Exe != "" {
		target = e.Exe
	}
	guard := strings.TrimRight(extendedGuard(e), "\n")
	if guard != "" {
		guard += " && "
	}
	echoSh := ctxExplainPrefixSh(target, e.Args)
	announce := ctxExplainAnnounce(target, e.Args)
	switch e.Mode {
	case constants.CtxModePrefill:
		return `cd "%f" && x-terminal-emulator -e sh -c 'printf "gitmap "; exec $SHELL'`
	case constants.CtxModeSilent:
		return fmt.Sprintf(`cd "%%f" && %sOUT=$(printf %%%%s '%s'; '%s' %s 2>&1) && notify-send 'gitmap' "$OUT"`, guard, announce, target, args)
	default:
		return fmt.Sprintf(`cd "%%f" && %sx-terminal-emulator -e sh -c "%s'%s' %s; exec $SHELL"`, guard, echoSh, target, args)
	}
}
