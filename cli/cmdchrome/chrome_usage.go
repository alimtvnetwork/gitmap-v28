package cmdchrome

import (
	"fmt"
	"os"
)

func printChromeUsage() {
	fmt.Fprintln(os.Stderr, "usage: gitmap chrome <open|extensions|flags|reset|observe|install|copy|export|import|list|delete|diff|which> [args]")
}

func printChromeRichUsage() {
	fmt.Printf("\n\033[1;96mGitMap Chrome Management\033[0m\n\n")
	fmt.Printf("Usage: gitmap chrome <command> [arguments]\n")
	fmt.Printf("Alias: gitmap cprof <command> [arguments]\n\n")
	printChromeLaunchTable()
	printChromeExtensionTable()
	printChromeFlagsAndResetTable()
	printChromeProfileTable()
	printChromeBatchTable()
	printChromeMaintenanceTable()
	fmt.Println()
}

func printChromeLaunchTable() {
	fmt.Printf("\033[1;94mLaunching & Tabs:\033[0m\n")
	fmt.Printf("  open (launch, tab) [urls]           Open URL(s) in specified or active profile\n")
	fmt.Printf("  observe (tabs, status) [profile]    Inspect running Chrome processes & open tabs (JSON/YAML)\n\n")
}

func printChromeExtensionTable() {
	fmt.Printf("\033[1;94mExtensions & Plugins:\033[0m\n")
	fmt.Printf("  extensions (ext, plugins) [prof]    List installed extensions across profiles\n")
	fmt.Printf("  extension-install (ext-in) <path>   Inject unpacked extension or .crx into profile\n")
	fmt.Printf("  extension-enable (ext-on) <pat>     Enable extension(s) matching pattern or ID\n")
	fmt.Printf("  extension-disable (ext-off) <pat>   Disable extension(s) matching pattern or ID\n")
	fmt.Printf("  extension-disable-all (ext-off-all) Disable ALL extensions in target profile\n\n")
}

func printChromeFlagsAndResetTable() {
	fmt.Printf("\033[1;94mFlags & Profile Reset:\033[0m\n")
	fmt.Printf("  flags (experiments) [ls|enable|..]  Inspect, toggle, or reset Chrome feature flags\n")
	fmt.Printf("  reset (clean, clear) [prof] [flags] Purge caches, cookies, history, or full profile\n\n")
}

func printChromeProfileTable() {
	fmt.Printf("\033[1;94mSingle Profile Operations:\033[0m\n")
	fmt.Printf("  copy (cpc) <src> <dst>              Copy a Chrome profile into an offline profile\n")
	fmt.Printf("  export (cpe) <name> [out]           Export profile snapshot (json, zip, sqlite, yaml)\n")
	fmt.Printf("  import (cpi) [file|dir|glob]        Import profile snapshot(s) safely with email mapping\n")
	fmt.Printf("  import-check (inspect) [dir|ls]     Preview snapshot metadata & actions before importing\n")
	fmt.Printf("  list (ls, cpl)                      List Chrome profiles known to gitmap\n")
	fmt.Printf("  reconcile (repair, sync-state)      Synchronize on-disk profile directories with Local State\n")
	fmt.Printf("  delete (rm, cpd) <name>             Remove a profile & stored artifacts\n")
	fmt.Printf("  merge (cpm) <src> <dst>             Merge selected pieces of one profile into another\n\n")
}

func printChromeBatchTable() {
	fmt.Printf("\033[1;94mBatch Profile Operations:\033[0m\n")
	fmt.Printf("  copy-all (cpc-all) <dst-dir>        Copy ALL discovered Chrome profiles to a directory\n")
	fmt.Printf("  export-all (cpe-all) [out-dir]      Export ALL Chrome profiles in batch\n")
	fmt.Printf("  import-all (cpi-all) <dir>          Import ALL profile snapshots from a directory\n\n")
}

func printChromeMaintenanceTable() {
	fmt.Printf("\033[1;94mBackup, Diff & Maintenance:\033[0m\n")
	fmt.Printf("  install (in)                        Install Chrome via system package manager\n")
	fmt.Printf("  backup [out]                        Snapshot all Chrome profiles into a tar.gz archive\n")
	fmt.Printf("  restore <tarball>                   Restore Chrome profiles from tar.gz archive\n")
	fmt.Printf("  diff <prof1> <prof2>                Compare two Chrome profiles\n")
	fmt.Printf("  bookmarks (export-bookmarks) <name> Export bookmarks to HTML or JSON\n")
	fmt.Printf("  which                               Print Chrome User Data & binary paths\n")
}
