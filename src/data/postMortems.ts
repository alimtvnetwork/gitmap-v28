export interface PostMortemEntry {
  id: string;
  title: string;
  summary: string;
  version?: string;
  category: "update" | "database" | "release" | "security" | "migration" | "general";
}

export const postMortems: PostMortemEntry[] = [
  {
    id: "01",
    title: "Update File Lock (Windows)",
    summary: "Windows blocks overwriting a running binary. Resolved via rename-first strategy — rename active binary to .old before copying the new one.",
    version: "v2.3.9",
    category: "update",
  },
  {
    id: "02",
    title: "Update Flow Spec Alignment",
    summary: "Update command behavior diverged from spec. Re-aligned the handoff mechanism and cleanup steps to match documented flow.",
    category: "update",
  },
  {
    id: "03",
    title: "Update Sync Lock Loop",
    summary: "Update process entered an infinite retry loop when the temp binary was also locked. Added PID-based detection and single-retry limit.",
    version: "v2.3.9–v2.3.11",
    category: "update",
  },
  {
    id: "04",
    title: "Database Written to Wrong Directory",
    summary: "SQLite database was created in the scan output directory instead of the binary's location. Fixed by anchoring data/ path to the executable.",
    category: "database",
  },
  {
    id: "05",
    title: "`gitmap ls` Returns Empty After Scan",
    summary: "Database path resolution double-nested the data/ folder. Corrected store package path logic.",
    category: "database",
  },
  {
    id: "06",
    title: "Release: Orphaned Metadata Recovery",
    summary: "Release metadata files left behind after failed releases. Added cleanup and recovery documentation.",
    category: "release",
  },
  {
    id: "07",
    title: "Zip Group Silent Failure",
    summary: "Zip groups were not created or uploaded during release with no error shown. Added explicit error reporting and diagnostic feedback.",
    category: "release",
  },
  {
    id: "08",
    title: "Auto-commit Push Rejection",
    summary: "Push failed when remote branch advanced during release. Implemented git pull --rebase recovery with single retry.",
    version: "v2.33.0",
    category: "release",
  },
  {
    id: "09",
    title: "list-releases Reads from DB Instead of Repo",
    summary: "list-releases used stale DB data instead of local .gitmap/release/ files. Prioritized local metadata over database.",
    version: "v2.34.0",
    category: "release",
  },
  {
    id: "10",
    title: "Legacy UUID Data Detection",
    summary: "DB queries failed with raw SQL errors when legacy UUID string IDs were present. Added detection and recovery prompts.",
    version: "v2.35.1",
    category: "migration",
  },
  {
    id: "11",
    title: "Automatic Legacy Directory Migration",
    summary: "Legacy directories (.release/, gitmap-output/) persisted across branch checkouts. Added startup migration with merge-and-remove strategy.",
    version: "v2.36.0",
    category: "migration",
  },
  {
    id: "12",
    title: "Legacy UUID to Integer ID Migration",
    summary: "Migrated database primary keys from UUID TEXT to INTEGER AUTOINCREMENT. Rebuilt Repos table with automatic migration.",
    version: "v2.36.1",
    category: "migration",
  },
  {
    id: "13",
    title: "Release Pipeline `dist` Directory Error",
    summary: "CI failed with 'cd: dist: No such file or directory' — compress step ran in wrong directory. Fixed with explicit working-directory directive.",
    version: "v2.54.0",
    category: "release",
  },
  {
    id: "14",
    title: "Security Hardening — G305, G110, Format Verbs",
    summary: "Fixed zip path traversal (G305), decompression bomb (G110), fmt.Fprintf argument mismatch, and standardized all error messages with mandatory path context.",
    version: "v2.54.1–v2.54.3",
    category: "security",
  },
  {
    id: "code-red",
    title: "Code Red: File Path Error Management",
    summary: "Every file/path error must include exact resolved path, operation, and failure reason. Generic 'file not found' messages are prohibited.",
    version: "v2.54.0–v2.54.1",
    category: "general",
  },
  {
    id: "15",
    title: "Installer Crashes — Progress Bar & Binary Detection",
    summary: "PowerShell progress bar crashed terminal during irm | iex. Versioned binary names (e.g., gitmap-v27.54.6-windows-amd64.exe) were not detected. Fixed with $ProgressPreference, regex matching, and top-level try/catch.",
    version: "v2.55.0",
    category: "general",
  },
  {
    id: "16",
    title: "CI Passthrough Gate Pattern",
    summary: "Job-level `if` skipping caused cached SHA runs to show grey 'Skipped' in GitHub UI instead of green Success. Replaced with step-level conditionals so every job always runs an 'Already validated' echo step and reports ✅.",
    version: "v2.55.0",
    category: "general",
  },
  {
    id: "17",
    title: "Go Flag Ordering — Silent Flag Drop",
    summary: "Go's `flag` package stops parsing at the first positional argument, silently dropping flags like `-y` after `v2.55.0`. Fixed with `reorderFlagsBeforeArgs()` to move all flags before positional args before parsing.",
    version: "v2.58.0",
    category: "release",
  },
  {
    id: "18",
    title: "CI Release Branch Cancellation Protection",
    summary: "Unconditional `cancel-in-progress: true` cancelled release branch CI runs on rapid pushes, risking incomplete artifacts. Fixed with a conditional expression that protects `release/**` branches from cancellation.",
    version: "v2.62.0",
    category: "general",
  },
];
