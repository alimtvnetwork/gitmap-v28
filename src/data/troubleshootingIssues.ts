import { AlertTriangle, FolderX, FileWarning, KeyRound, Network, Lock, GitBranch, Wrench, LucideIcon } from "lucide-react";

export enum CategoryType {
  Paths = "paths",
  Config = "config",
  Auth = "auth",
  Network = "network",
  Locks = "locks",
  Git = "git",
  Build = "build",
}

export type Category = CategoryType;

export enum CategoryFilterType {
  All = "all",
}

export interface TroubleshootingRelatedLink {
  label: string;
  href: string;
}

export interface Issue {
  id: string;
  category: CategoryType;
  title: string;
  symptom: string;
  cause: string;
  fix: string;
  fixCommand?: string;
  fixLanguage?: string;
  altCommand?: string;
  altLabel?: string;
  related?: TroubleshootingRelatedLink[];
}

export interface TroubleshootingCategoryMetaItem {
  label: string;
  icon: LucideIcon;
}

export const CATEGORY_META: Record<CategoryType, TroubleshootingCategoryMetaItem> = {
  [CategoryType.Paths]: { label: "Invalid paths", icon: FolderX },
  [CategoryType.Config]: { label: "Missing config", icon: FileWarning },
  [CategoryType.Auth]: { label: "Auth & SSH", icon: KeyRound },
  [CategoryType.Network]: { label: "Network", icon: Network },
  [CategoryType.Locks]: { label: "File locks", icon: Lock },
  [CategoryType.Git]: { label: "Git state", icon: GitBranch },
  [CategoryType.Build]: { label: "Build & PATH", icon: Wrench },
};

export const categoryMeta = CATEGORY_META;

export const issues: Issue[] = [
  // ── paths ─────────────────────────────────────
  {
    id: "scan-bad-dir",
    category: CategoryType.Paths,
    title: "scan: target directory does not exist",
    symptom: "Error: scan path does not exist: D:\\repos",
    cause: "The positional [dir] argument points to a folder that has not been created or is on an unmounted drive.",
    fix: "Pass an existing absolute path, or omit it to scan the current directory.",
    fixCommand: `gitmap scan "C:\\Users\\you\\projects"\n# or just\ngitmap scan`,
    related: [
      { label: "scan reference", href: "/scan-command" },
      { label: "doctor", href: "/doctor" },
    ],
  },
  {
    id: "out-file-conflict",
    category: CategoryType.Paths,
    title: "scan: --out-file conflicts with --output-path",
    symptom: "Output file written somewhere other than expected.",
    cause: "--out-file is an exact path and overrides --output-path. Passing both leads to confusion.",
    fix: "Pick one: --out-file <exact-path> for a specific file, or --output-path <dir> for a directory.",
    fixCommand: `gitmap scan --out-file ./.gitmap/output/custom.json --output json`,
  },
  {
    id: "clone-target-missing",
    category: CategoryType.Paths,
    title: "clone: --target-dir does not exist",
    symptom: "Error: target directory does not exist: D:\\new-projects",
    cause: "Clone refuses to create deep parent paths to avoid clobbering typos.",
    fix: "Create the directory first, then re-run.",
    fixCommand: `mkdir D:\\new-projects\ngitmap clone json --target-dir D:\\new-projects`,
    fixLanguage: "powershell",
  },
  {
    id: "windows-long-paths",
    category: CategoryType.Paths,
    title: "Windows: 'unable to create file' / path too long",
    symptom: "git clone fails with 'Filename too long' or 'unable to create file'.",
    cause: "Windows enforces a 260-char path limit unless long-path support is enabled in Git.",
    fix: "Enable long paths globally in Git, then retry.",
    fixCommand: `git config --global core.longpaths true`,
  },

  // ── config ────────────────────────────────────
  {
    id: "config-missing",
    category: CategoryType.Config,
    title: "config.json not found",
    symptom: "Doctor warns: config.json not found (using defaults)",
    cause: "Either no config has been created, or you are running from a directory without a ./data/ folder.",
    fix: "Pass --config explicitly, or create the file in the default location.",
    fixCommand: `gitmap scan --config "C:\\Users\\you\\.gitmap\\config.json"`,
    related: [{ label: "Configuration", href: "/config" }],
  },
  {
    id: "config-invalid-json",
    category: CategoryType.Config,
    title: "config.json is not valid JSON",
    symptom: "Doctor: config.json is not valid JSON",
    cause: "Trailing commas, unquoted keys, or a corrupted edit.",
    fix: "Validate the file, then re-save. Run doctor to confirm.",
    fixCommand: `gitmap doctor`,
  },
  {
    id: "setup-config-missing",
    category: CategoryType.Config,
    title: "git-setup.json not found (setup fails)",
    symptom: "Doctor: git-setup.json not found (setup will fail without --config)",
    cause: "First-time setup needs git-setup.json next to the binary or supplied explicitly.",
    fix: "Point setup at the file using --config.",
    fixCommand: `gitmap setup --config ./data/git-setup.json`,
    related: [{ label: "Setup", href: "/setup" }],
  },
  {
    id: "repopath-missing",
    category: CategoryType.Config,
    title: "RepoPath not embedded — self-update broken",
    symptom: "Doctor: RepoPath not embedded. Binary was not built with run.ps1.",
    cause: "Binary was built directly via 'go build' instead of the project's run.ps1, which embeds RepoPath at compile time.",
    fix: "Rebuild with run.ps1 from the source checkout.",
    fixCommand: `.\\run.ps1`,
    fixLanguage: "powershell",
  },

  // ── auth ──────────────────────────────────────
  {
    id: "ssh-key-missing",
    category: CategoryType.Auth,
    title: "clone-next: named SSH key not found",
    symptom: "Error: ssh key 'work' not registered",
    cause: "The named key passed via --ssh-key / -K has not been registered with gitmap.",
    fix: "List, then register the key.",
    fixCommand: `gitmap ssh list\ngitmap ssh add work ~/.ssh/id_ed25519_work`,
    related: [{ label: "SSH keys", href: "/ssh" }],
  },
  {
    id: "permission-denied-publickey",
    category: CategoryType.Auth,
    title: "git: Permission denied (publickey)",
    symptom: "fatal: Could not read from remote repository. Permission denied (publickey).",
    cause: "Either ssh-agent has no key loaded, or the wrong key is selected for this remote.",
    fix: "Switch the scan/clone to a named key, or add the key to ssh-agent.",
    fixCommand: `gitmap clone-next v++ --ssh-key work\n# or\nssh-add ~/.ssh/id_ed25519`,
  },
  {
    id: "https-mode-needed",
    category: CategoryType.Auth,
    title: "Repos cloned with SSH URLs but org only allows HTTPS",
    symptom: "ssh: connect to host github.com port 22: Connection refused",
    cause: "Output was generated with --mode ssh on a network that blocks port 22.",
    fix: "Re-scan with --mode https, then re-clone.",
    fixCommand: `gitmap scan ~/projects --mode https --output json\ngitmap clone json`,
  },
  {
    id: "create-remote-no-token",
    category: CategoryType.Auth,
    title: "clone-next --create-remote: GITHUB_TOKEN missing",
    symptom: "Error: cannot create GitHub repo macro-ahk-v22: GITHUB_TOKEN not set",
    cause: "--create-remote calls the GitHub API and requires a token with 'repo' scope.",
    fix: "Export a token with repo scope, then re-run.",
    fixCommand: `# bash / zsh\nexport GITHUB_TOKEN=ghp_xxx\n# PowerShell\n$env:GITHUB_TOKEN="ghp_xxx"\n\ngitmap cn v++ --create-remote`,
  },

  // ── network ───────────────────────────────────
  {
    id: "github-unreachable",
    category: CategoryType.Network,
    title: "github.com unreachable (offline mode)",
    symptom: "Doctor: Network: github.com unreachable (offline mode)",
    cause: "DNS, proxy, or VPN issue blocking the host.",
    fix: "Confirm reachability, then re-run with verbose logging if it persists.",
    fixCommand: `curl -I https://github.com\ngitmap clone json --verbose`,
  },

  // ── locks ─────────────────────────────────────
  {
    id: "cn-cwd-locked",
    category: CategoryType.Locks,
    title: "clone-next falls back to versioned folder instead of flattening",
    symptom: "→ Falling back to versioned folder macro-ahk-v22 (current folder is locked by this shell)",
    cause: "Your shell is cwd'd into the target folder, so Windows holds a file lock and gitmap can't replace it.",
    fix: "Pass -f to force a chdir-to-parent + flatten, or 'cd ..' first.",
    fixCommand: `gitmap cn v+1 -f`,
    altLabel: "Manual alternative",
    altCommand: `cd ..\ngitmap cn v+1`,
    related: [{ label: "clone-next reference", href: "/clone-next-command" }],
  },
  {
    id: "cn-force-cant-remove",
    category: CategoryType.Locks,
    title: "clone-next -f: another process holds the folder",
    symptom: "Error: --force could not remove macro-ahk: unlinkat: access denied",
    cause: "An editor, file explorer, or watcher (not your shell) has an open handle on the folder.",
    fix: "Close the holder (VS Code, Explorer preview pane), then retry. The lock-detector can name PIDs.",
    fixCommand: `gitmap cn v+1 --delete --verbose`,
  },
  {
    id: "stale-lockfile",
    category: CategoryType.Locks,
    title: "Stale gitmap lock file",
    symptom: "Doctor: Lock file exists — another gitmap may be running (or stale)",
    cause: "A previous gitmap process exited without releasing its advisory lock.",
    fix: "If no other gitmap is running, run doctor — it will surface the path; remove it manually only if needed.",
    fixCommand: `gitmap doctor`,
  },

  // ── git ───────────────────────────────────────
  {
    id: "cn-no-remote",
    category: CategoryType.Git,
    title: "clone-next: not a git repo or no remote origin",
    symptom: "Error: not a git repo or no remote origin",
    cause: "clone-next reads the remote of the cwd to know what to clone next; the cwd has no origin.",
    fix: "cd into a real cloned repo first, or set an origin.",
    fixCommand: `git remote add origin https://github.com/you/repo-v1.git\ngitmap cn v++`,
  },
  {
    id: "cn-bad-version-arg",
    category: CategoryType.Git,
    title: "clone-next: invalid version argument",
    symptom: "Error: invalid version argument: foo (expected v++, v+1, or vN)",
    cause: "Only v++, v+1, or vN (positive integer) are accepted.",
    fix: "Use a valid form.",
    fixCommand: `gitmap cn v++\ngitmap cn v15`,
  },
  {
    id: "clone-source-missing",
    category: CategoryType.Git,
    title: "clone: source file not found",
    symptom: "Error: source file not found: .gitmap/output/gitmap.json",
    cause: "You haven't scanned yet, or you're running from a different directory than where output was written.",
    fix: "Run scan first, or pass an explicit path.",
    fixCommand: `gitmap scan ~/projects --output json\ngitmap clone json`,
  },

  // ── build ─────────────────────────────────────
  {
    id: "not-on-path",
    category: CategoryType.Build,
    title: "gitmap not found on PATH",
    symptom: "bash: gitmap-v28: command not found",
    cause: "Deploy directory is not on PATH, or the deployed binary differs from the active one.",
    fix: "Auto-sync the PATH binary from the deployed copy.",
    fixCommand: `gitmap doctor --fix-path`,
    related: [{ label: "Doctor", href: "/doctor" }],
  },
  {
    id: "version-mismatch",
    category: CategoryType.Build,
    title: "PATH binary version mismatch",
    symptom: "Doctor: PATH binary version mismatch (PATH: 3.50.0, Source: 3.52.0)",
    cause: "You rebuilt the source but the PATH binary was not refreshed.",
    fix: "Run the auto-fixer, or rebuild from source.",
    fixCommand: `gitmap doctor --fix-path\n# or full rebuild\n.\\run.ps1`,
    fixLanguage: "powershell",
  },
  {
    id: "wrapper-not-loaded",
    category: CategoryType.Build,
    title: "gitmap cd prints path but doesn't change directory",
    symptom: "Doctor: Shell wrapper not loaded — gitmap cd prints path but cannot change directory",
    cause: "The shell function wrapper that intercepts 'gitmap cd' has not been sourced.",
    fix: "Run setup, then reload your shell profile.",
    fixCommand: `gitmap setup\n# bash\nsource ~/.bashrc\n# zsh\nsource ~/.zshrc\n# PowerShell\n. $PROFILE`,
  },
];

export function isValidCategoryKey(value: string): value is CategoryType {
  const allKeys = Object.keys(CATEGORY_META) as string[];

  return allKeys.includes(value);
}
