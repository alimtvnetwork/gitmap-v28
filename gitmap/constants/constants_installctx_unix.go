package constants

// Unix (macOS + Linux) install paths and labels for `gitmap install ctx`.
// Spec: spec/04-generic-cli/30-install-ctx.md §7.

// Path fragments under $HOME — joined at runtime so unit tests stay
// hermetic (no hard-coded absolute paths).
const (
	// macOS
	CtxMacServicesRel = "Library/Services"

	// Linux file managers
	CtxLinuxNautilusRel = ".local/share/nautilus/scripts/gitmap"
	CtxLinuxDolphinRel  = ".local/share/kio/servicemenus"
	CtxLinuxThunarRel   = ".config/Thunar/uca.xml"

	CtxLinuxDolphinFile = "gitmap-ctx.desktop"
)

// Flat-menu prefix used on macOS/Linux where nested submenus are not
// available. Final label = "<prefix>: <Category> — <Child>".
const (
	CtxFlatPrefix       = "gitmap"
	CtxFlatSeparator    = ": "
	CtxFlatChildJoiner  = " — "
	CtxThunarMarkBegin  = "<!-- gitmap-ctx-begin -->"
	CtxThunarMarkEnd    = "<!-- gitmap-ctx-end -->"
	CtxThunarUcaRootTag = "actions"
)

// User-facing messages — Unix path.
const (
	MsgCtxMacInstallStart   = "  Adding gitmap to macOS Services (Quick Actions)...\n"
	MsgCtxMacInstallDone    = "  ✓ gitmap Quick Actions installed (%d/%d workflows). Restart Finder to refresh: pkill -KILL -u $USER cfprefsd\n"
	MsgCtxMacUninstallStart = "  Removing gitmap Quick Actions from ~/Library/Services...\n"
	MsgCtxMacUninstallDone  = "  ✓ gitmap Quick Actions removed (%d/%d workflows).\n"

	MsgCtxLinuxInstallStart   = "  Adding gitmap to Linux file-manager context menus (Nautilus / Dolphin / Thunar)...\n"
	MsgCtxLinuxInstallDone    = "  ✓ gitmap Linux context menu installed (%d entries across %d managers).\n"
	MsgCtxLinuxUninstallStart = "  Removing gitmap from Linux file-manager context menus...\n"
	MsgCtxLinuxUninstallDone  = "  ✓ gitmap Linux context menu removed (%d entries).\n"

	MsgCtxFsWriteFail = "  ! write %s: %v\n"
	MsgCtxFsRmFail    = "  ! remove %s: %v\n"
)
