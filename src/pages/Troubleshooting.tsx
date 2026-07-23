import { useState, useMemo, useCallback, useEffect, useRef } from "react";
import { useSearchParams } from "react-router-dom";
import DocsLayout from "@/components/docs/DocsLayout";
import CodeBlock from "@/components/docs/CodeBlock";
import SearchBar from "@/components/docs/SearchBar";
import { AlertTriangle, FolderX, FileWarning, KeyRound, Network, Lock, GitBranch, Wrench, Copy, Check, Link2, Stethoscope, Terminal, FileText, ListChecks } from "lucide-react";
import { copyToClipboard } from "@/lib/clipboard";

type Category = "paths" | "config" | "auth" | "network" | "locks" | "git" | "build";

interface Issue {
  id: string;
  category: Category;
  title: string;
  symptom: string;
  cause: string;
  fix: string;
  fixCommand?: string;
  fixLanguage?: string;
  altCommand?: string;
  altLabel?: string;
  related?: { label: string; href: string }[];
}

const categoryMeta: Record<Category, { label: string; icon: typeof AlertTriangle }> = {
  paths: { label: "Invalid paths", icon: FolderX },
  config: { label: "Missing config", icon: FileWarning },
  auth: { label: "Auth & SSH", icon: KeyRound },
  network: { label: "Network", icon: Network },
  locks: { label: "File locks", icon: Lock },
  git: { label: "Git state", icon: GitBranch },
  build: { label: "Build & PATH", icon: Wrench },
};

const issues: Issue[] = [
  // ── paths ─────────────────────────────────────
  {
    id: "scan-bad-dir",
    category: "paths",
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
    category: "paths",
    title: "scan: --out-file conflicts with --output-path",
    symptom: "Output file written somewhere other than expected.",
    cause: "--out-file is an exact path and overrides --output-path. Passing both leads to confusion.",
    fix: "Pick one: --out-file <exact-path> for a specific file, or --output-path <dir> for a directory.",
    fixCommand: `gitmap scan --out-file ./.gitmap/output/custom.json --output json`,
  },
  {
    id: "clone-target-missing",
    category: "paths",
    title: "clone: --target-dir does not exist",
    symptom: "Error: target directory does not exist: D:\\new-projects",
    cause: "Clone refuses to create deep parent paths to avoid clobbering typos.",
    fix: "Create the directory first, then re-run.",
    fixCommand: `mkdir D:\\new-projects\ngitmap clone json --target-dir D:\\new-projects`,
    fixLanguage: "powershell",
  },
  {
    id: "windows-long-paths",
    category: "paths",
    title: "Windows: 'unable to create file' / path too long",
    symptom: "git clone fails with 'Filename too long' or 'unable to create file'.",
    cause: "Windows enforces a 260-char path limit unless long-path support is enabled in Git.",
    fix: "Enable long paths globally in Git, then retry.",
    fixCommand: `git config --global core.longpaths true`,
  },

  // ── config ────────────────────────────────────
  {
    id: "config-missing",
    category: "config",
    title: "config.json not found",
    symptom: "Doctor warns: config.json not found (using defaults)",
    cause: "Either no config has been created, or you are running from a directory without a ./data/ folder.",
    fix: "Pass --config explicitly, or create the file in the default location.",
    fixCommand: `gitmap scan --config "C:\\Users\\you\\.gitmap\\config.json"`,
    related: [{ label: "Configuration", href: "/config" }],
  },
  {
    id: "config-invalid-json",
    category: "config",
    title: "config.json is not valid JSON",
    symptom: "Doctor: config.json is not valid JSON",
    cause: "Trailing commas, unquoted keys, or a corrupted edit.",
    fix: "Validate the file, then re-save. Run doctor to confirm.",
    fixCommand: `gitmap doctor`,
  },
  {
    id: "setup-config-missing",
    category: "config",
    title: "git-setup.json not found (setup fails)",
    symptom: "Doctor: git-setup.json not found (setup will fail without --config)",
    cause: "First-time setup needs git-setup.json next to the binary or supplied explicitly.",
    fix: "Point setup at the file using --config.",
    fixCommand: `gitmap setup --config ./data/git-setup.json`,
    related: [{ label: "Setup", href: "/setup" }],
  },
  {
    id: "repopath-missing",
    category: "config",
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
    category: "auth",
    title: "clone-next: named SSH key not found",
    symptom: "Error: ssh key 'work' not registered",
    cause: "The named key passed via --ssh-key / -K has not been registered with gitmap.",
    fix: "List, then register the key.",
    fixCommand: `gitmap ssh list\ngitmap ssh add work ~/.ssh/id_ed25519_work`,
    related: [{ label: "SSH keys", href: "/ssh" }],
  },
  {
    id: "permission-denied-publickey",
    category: "auth",
    title: "git: Permission denied (publickey)",
    symptom: "fatal: Could not read from remote repository. Permission denied (publickey).",
    cause: "Either ssh-agent has no key loaded, or the wrong key is selected for this remote.",
    fix: "Switch the scan/clone to a named key, or add the key to ssh-agent.",
    fixCommand: `gitmap clone-next v++ --ssh-key work\n# or\nssh-add ~/.ssh/id_ed25519`,
  },
  {
    id: "https-mode-needed",
    category: "auth",
    title: "Repos cloned with SSH URLs but org only allows HTTPS",
    symptom: "ssh: connect to host github.com port 22: Connection refused",
    cause: "Output was generated with --mode ssh on a network that blocks port 22.",
    fix: "Re-scan with --mode https, then re-clone.",
    fixCommand: `gitmap scan ~/projects --mode https --output json\ngitmap clone json`,
  },
  {
    id: "create-remote-no-token",
    category: "auth",
    title: "clone-next --create-remote: GITHUB_TOKEN missing",
    symptom: "Error: cannot create GitHub repo macro-ahk-v22: GITHUB_TOKEN not set",
    cause: "--create-remote calls the GitHub API and requires a token with 'repo' scope.",
    fix: "Export a token with repo scope, then re-run.",
    fixCommand: `# bash / zsh\nexport GITHUB_TOKEN=ghp_xxx\n# PowerShell\n$env:GITHUB_TOKEN="ghp_xxx"\n\ngitmap cn v++ --create-remote`,
  },

  // ── network ───────────────────────────────────
  {
    id: "github-unreachable",
    category: "network",
    title: "github.com unreachable (offline mode)",
    symptom: "Doctor: Network: github.com unreachable (offline mode)",
    cause: "DNS, proxy, or VPN issue blocking the host.",
    fix: "Confirm reachability, then re-run with verbose logging if it persists.",
    fixCommand: `curl -I https://github.com\ngitmap clone json --verbose`,
  },

  // ── locks ─────────────────────────────────────
  {
    id: "cn-cwd-locked",
    category: "locks",
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
    category: "locks",
    title: "clone-next -f: another process holds the folder",
    symptom: "Error: --force could not remove macro-ahk: unlinkat: access denied",
    cause: "An editor, file explorer, or watcher (not your shell) has an open handle on the folder.",
    fix: "Close the holder (VS Code, Explorer preview pane), then retry. The lock-detector can name PIDs.",
    fixCommand: `gitmap cn v+1 --delete --verbose`,
  },
  {
    id: "stale-lockfile",
    category: "locks",
    title: "Stale gitmap lock file",
    symptom: "Doctor: Lock file exists — another gitmap may be running (or stale)",
    cause: "A previous gitmap process exited without releasing its advisory lock.",
    fix: "If no other gitmap is running, run doctor — it will surface the path; remove it manually only if needed.",
    fixCommand: `gitmap doctor`,
  },

  // ── git ───────────────────────────────────────
  {
    id: "cn-no-remote",
    category: "git",
    title: "clone-next: not a git repo or no remote origin",
    symptom: "Error: not a git repo or no remote origin",
    cause: "clone-next reads the remote of the cwd to know what to clone next; the cwd has no origin.",
    fix: "cd into a real cloned repo first, or set an origin.",
    fixCommand: `git remote add origin https://github.com/you/repo-v1.git\ngitmap cn v++`,
  },
  {
    id: "cn-bad-version-arg",
    category: "git",
    title: "clone-next: invalid version argument",
    symptom: "Error: invalid version argument: foo (expected v++, v+1, or vN)",
    cause: "Only v++, v+1, or vN (positive integer) are accepted.",
    fix: "Use a valid form.",
    fixCommand: `gitmap cn v++\ngitmap cn v15`,
  },
  {
    id: "clone-source-missing",
    category: "git",
    title: "clone: source file not found",
    symptom: "Error: source file not found: .gitmap/output/gitmap.json",
    cause: "You haven't scanned yet, or you're running from a different directory than where output was written.",
    fix: "Run scan first, or pass an explicit path.",
    fixCommand: `gitmap scan ~/projects --output json\ngitmap clone json`,
  },

  // ── build ─────────────────────────────────────
  {
    id: "not-on-path",
    category: "build",
    title: "gitmap not found on PATH",
    symptom: "bash: gitmap-v27: command not found",
    cause: "Deploy directory is not on PATH, or the deployed binary differs from the active one.",
    fix: "Auto-sync the PATH binary from the deployed copy.",
    fixCommand: `gitmap doctor --fix-path`,
    related: [{ label: "Doctor", href: "/doctor" }],
  },
  {
    id: "version-mismatch",
    category: "build",
    title: "PATH binary version mismatch",
    symptom: "Doctor: PATH binary version mismatch (PATH: 3.50.0, Source: 3.52.0)",
    cause: "You rebuilt the source but the PATH binary was not refreshed.",
    fix: "Run the auto-fixer, or rebuild from source.",
    fixCommand: `gitmap doctor --fix-path\n# or full rebuild\n.\\run.ps1`,
    fixLanguage: "powershell",
  },
  {
    id: "wrapper-not-loaded",
    category: "build",
    title: "gitmap cd prints path but doesn't change directory",
    symptom: "Doctor: Shell wrapper not loaded — gitmap cd prints path but cannot change directory",
    cause: "The shell function wrapper that intercepts 'gitmap cd' has not been sourced.",
    fix: "Run setup, then reload your shell profile.",
    fixCommand: `gitmap setup\n# bash\nsource ~/.bashrc\n# zsh\nsource ~/.zshrc\n# PowerShell\n. $PROFILE`,
  },
];

const isValidCategoryKey = (v: string): v is Category =>
  (Object.keys(categoryMeta) as string[]).includes(v);

const Troubleshooting = () => {
  const [searchParams, setSearchParams] = useSearchParams();
  const initialSearch = searchParams.get("search") ?? searchParams.get("q") ?? "";
  const initialCategoryRaw = searchParams.get("category") ?? "all";
  const initialCategory: Category | "all" =
    initialCategoryRaw === "all" || isValidCategoryKey(initialCategoryRaw)
      ? (initialCategoryRaw as Category | "all")
      : "all";

  const [search, setSearch] = useState(initialSearch);
  const [activeCategory, setActiveCategory] = useState<Category | "all">(initialCategory);
  const scrolledIdRef = useRef<string | null>(null);

  // Sync state -> URL (replace, no history entry per keystroke).
  useEffect(() => {
    const next = new URLSearchParams(searchParams);
    if (search) next.set("search", search);
    else next.delete("search");
    next.delete("q");
    if (activeCategory !== "all") next.set("category", activeCategory);
    else next.delete("category");
    if (next.toString() !== searchParams.toString()) {
      setSearchParams(next, { replace: true });
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [search, activeCategory]);

  const filtered = useMemo(() => {
    let rows = issues;
    if (activeCategory !== "all") rows = rows.filter((i) => i.category === activeCategory);
    if (search) {
      const q = search.toLowerCase();
      rows = rows.filter(
        (i) =>
          i.title.toLowerCase().includes(q) ||
          i.symptom.toLowerCase().includes(q) ||
          i.cause.toLowerCase().includes(q) ||
          i.fix.toLowerCase().includes(q) ||
          (i.fixCommand?.toLowerCase().includes(q) ?? false)
      );
    }
    return rows;
  }, [search, activeCategory]);

  const categoryCounts = useMemo(() => {
    const counts: Record<string, number> = { all: issues.length };
    for (const i of issues) counts[i.category] = (counts[i.category] ?? 0) + 1;
    return counts;
  }, []);

  // Deep-link: scroll to ?id=<issue-id> and relax filters that would hide it.
  const targetId = searchParams.get("id");
  useEffect(() => {
    if (!targetId) return;
    const issue = issues.find((i) => i.id === targetId);
    if (!issue) return;
    if (activeCategory !== "all" && activeCategory !== issue.category) {
      setActiveCategory("all");
      return;
    }
    if (filtered.findIndex((i) => i.id === targetId) === -1) {
      if (search) setSearch("");
      return;
    }
    if (scrolledIdRef.current === targetId) return;
    const el = document.getElementById(targetId);
    if (el) {
      scrolledIdRef.current = targetId;
      el.scrollIntoView({ behavior: "smooth", block: "start" });
      el.classList.add("ring-2", "ring-primary");
      window.setTimeout(() => el.classList.remove("ring-2", "ring-primary"), 2400);
    }
  }, [targetId, filtered, activeCategory, search]);

  return (
    <DocsLayout>
      <h1 className="text-3xl font-heading font-bold mb-2 docs-h1">Troubleshooting</h1>
      <p className="text-muted-foreground mb-6">
        Common gitmap errors grouped by category, each with the symptom, root cause, and the exact
        flag or command to fix it. When in doubt, start with{" "}
        <code className="docs-inline-code">gitmap doctor</code>.
      </p>

      <DiagnosticChecklist />

      <SearchBar value={search} onChange={setSearch} placeholder="Search by error, symptom, or fix..." />

      <div className="flex flex-wrap gap-2 mt-4 mb-8">
        <button
          onClick={() => setActiveCategory("all")}
          className={`px-3 py-1.5 rounded-md text-sm font-mono border transition-colors ${
            activeCategory === "all"
              ? "bg-primary text-primary-foreground border-primary"
              : "bg-background text-muted-foreground border-border hover:text-foreground hover:border-foreground/30"
          }`}
        >
          all ({categoryCounts.all})
        </button>
        {(Object.keys(categoryMeta) as Category[]).map((key) => {
          const Icon = categoryMeta[key].icon;
          return (
            <button
              key={key}
              onClick={() => setActiveCategory(key)}
              className={`px-3 py-1.5 rounded-md text-sm font-mono border transition-colors flex items-center gap-1.5 ${
                activeCategory === key
                  ? "bg-primary text-primary-foreground border-primary"
                  : "bg-background text-muted-foreground border-border hover:text-foreground hover:border-foreground/30"
              }`}
            >
              <Icon className="h-3.5 w-3.5" />
              {categoryMeta[key].label} ({categoryCounts[key] ?? 0})
            </button>
          );
        })}
      </div>

      <section className="space-y-6">
        {filtered.map((issue) => {
          const Icon = categoryMeta[issue.category].icon;
          return (
            <article
              key={issue.id}
              id={issue.id}
              className="rounded-lg border border-border bg-card overflow-hidden"
            >
              <header className="px-5 py-3 border-b border-border bg-muted/30 flex items-start gap-3">
                <Icon className="h-5 w-5 text-primary shrink-0 mt-0.5" />
                <div className="flex-1 min-w-0">
                  <h2 className="text-base font-heading font-semibold docs-h3">{issue.title}</h2>
                  <p className="text-xs font-mono text-muted-foreground mt-0.5">
                    {categoryMeta[issue.category].label}
                  </p>
                </div>
                <CopyLinkButton issueId={issue.id} />
                {issue.fixCommand && (
                  <CopyFixButton command={issue.fixCommand} altCommand={issue.altCommand} />
                )}
              </header>

              <div className="p-5 space-y-4">
                <div>
                  <h3 className="text-xs font-mono uppercase tracking-wider text-muted-foreground mb-1">
                    Symptom
                  </h3>
                  <pre className="text-sm font-mono bg-muted/40 border border-border rounded p-3 overflow-x-auto">
                    {issue.symptom}
                  </pre>
                </div>

                <div>
                  <h3 className="text-xs font-mono uppercase tracking-wider text-muted-foreground mb-1">
                    Cause
                  </h3>
                  <p className="text-sm text-foreground/90">{issue.cause}</p>
                </div>

                <div>
                  <h3 className="text-xs font-mono uppercase tracking-wider text-muted-foreground mb-1">
                    Fix
                  </h3>
                  <p className="text-sm text-foreground/90 mb-2">{issue.fix}</p>
                  {issue.fixCommand && (
                    <CodeBlock
                      language={issue.fixLanguage ?? "bash"}
                      code={issue.fixCommand}
                      title="Run"
                    />
                  )}
                  {issue.altCommand && (
                    <div className="mt-2">
                      <p className="text-xs font-mono text-muted-foreground mb-1">
                        {issue.altLabel ?? "Alternative"}
                      </p>
                      <CodeBlock language="bash" code={issue.altCommand} title="Alt" />
                    </div>
                  )}
                </div>

                {issue.related && issue.related.length > 0 && (
                  <div className="pt-2 border-t border-border flex flex-wrap items-center gap-3 text-xs">
                    <span className="text-muted-foreground font-mono uppercase tracking-wider">
                      Related
                    </span>
                    {issue.related.map((r) => (
                      <a
                        key={r.href}
                        href={r.href}
                        className="text-primary hover:underline font-mono"
                      >
                        {r.label}
                      </a>
                    ))}
                  </div>
                )}
              </div>
            </article>
          );
        })}

        {filtered.length === 0 && (
          <div className="rounded-lg border border-border p-8 text-center">
            <AlertTriangle className="h-8 w-8 text-muted-foreground mx-auto mb-3" />
            <p className="font-mono text-sm text-muted-foreground">
              No issues match "{search}". Try a different keyword or clear the filter.
            </p>
          </div>
        )}
      </section>

      <aside className="mt-10 rounded-lg border border-border bg-muted/20 p-5">
        <h2 className="font-heading font-semibold text-foreground mb-2">Still stuck?</h2>
        <ul className="list-disc pl-5 space-y-1 text-sm text-muted-foreground">
          <li>
            Run <code className="docs-inline-code">gitmap doctor</code> for a full health snapshot.
          </li>
          <li>
            Re-run the failing command with <code className="docs-inline-code">--verbose</code> to
            generate a timestamped debug log.
          </li>
          <li>
            Check the <a href="/post-mortems" className="text-primary hover:underline">post-mortems</a>{" "}
            for past incidents and their resolutions.
          </li>
        </ul>
      </aside>
    </DocsLayout>
  );
};

interface CopyFixButtonProps {
  command: string;
  altCommand?: string;
}

// CopyFixButton — one-click copy of the primary fix command (and optional
// alternative) shown in the issue header. Uses the same clipboard API path as
// CodeBlock so behavior is consistent across the page.
const CopyFixButton = ({ command, altCommand }: CopyFixButtonProps) => {
  const [copied, setCopied] = useState(false);

  const handleCopy = useCallback(async () => {
    const payload = altCommand
      ? `${command}\n\n# Alternative\n${altCommand}`
      : command;
    await copyToClipboard(payload);
    setCopied(true);
    window.setTimeout(() => setCopied(false), 2000);
  }, [command, altCommand]);

  return (
    <button
      type="button"
      onClick={handleCopy}
      aria-label={copied ? "Fix command copied" : "Copy fix command"}
      title={copied ? "Copied!" : "Copy fix command"}
      className={`shrink-0 inline-flex items-center gap-1.5 px-2.5 py-1.5 rounded-md text-xs font-mono border transition-colors ${
        copied
          ? "border-primary bg-primary/15 text-foreground dark:bg-primary/20 dark:text-primary"
          : "border-border bg-background text-muted-foreground hover:text-foreground hover:border-foreground/40"
      }`}
    >
      {copied ? (
        <>
          <Check className="h-3.5 w-3.5" />
          Copied
        </>
      ) : (
        <>
          <Copy className="h-3.5 w-3.5" />
          Copy fix
        </>
      )}
    </button>
  );
};

// CopyLinkButton — copies a deep-link to this specific issue card so it can
// be shared and re-opened directly via the ?id= query parameter.
const CopyLinkButton = ({ issueId }: { issueId: string }) => {
  const [copied, setCopied] = useState(false);

  const handleCopy = useCallback(async () => {
    const url = new URL(window.location.href);
    url.searchParams.set("id", issueId);
    await copyToClipboard(url.toString());
    setCopied(true);
    window.setTimeout(() => setCopied(false), 2000);
  }, [issueId]);

  return (
    <button
      type="button"
      onClick={handleCopy}
      aria-label={copied ? "Link copied" : "Copy link to this issue"}
      title={copied ? "Link copied!" : "Copy link to this issue"}
      className={`shrink-0 inline-flex items-center gap-1.5 px-2.5 py-1.5 rounded-md text-xs font-mono border transition-colors ${
        copied
          ? "border-primary bg-primary/15 text-foreground dark:bg-primary/20 dark:text-primary"
          : "border-border bg-background text-muted-foreground hover:text-foreground hover:border-foreground/40"
      }`}
    >
      {copied ? (
        <>
          <Check className="h-3.5 w-3.5" />
          Linked
        </>
      ) : (
        <>
          <Link2 className="h-3.5 w-3.5" />
          Link
        </>
      )}
    </button>
  );
};

// DiagnosticChecklist — short, ordered walkthrough that surfaces the
// canonical "what to run when something breaks" sequence: doctor first,
// auto-fix second, verbose re-run third, then read the log.
const DiagnosticChecklist = () => {
  const steps: {
    icon: typeof Stethoscope;
    title: string;
    body: string;
    command: string;
  }[] = [
    {
      icon: Stethoscope,
      title: "Run a full health check",
      body: "Surfaces PATH, config, lock, network, and version-mismatch issues in one pass. Most problems are diagnosed (and named) here.",
      command: "gitmap doctor",
    },
    {
      icon: Wrench,
      title: "Apply the auto-fixer if doctor flagged PATH",
      body: "If you see 'PATH binary version mismatch' or 'gitmap not found on PATH', let doctor sync the active binary from the deployed copy.",
      command: "gitmap doctor --fix-path",
    },
    {
      icon: Terminal,
      title: "Re-run the failing command with --verbose",
      body: "Verbose mode prints every git invocation, resolved paths, and timing — and writes a timestamped debug log next to the output.",
      command: "gitmap <your-failing-command> --verbose",
    },
    {
      icon: FileText,
      title: "Inspect the verbose log",
      body: "Logs land in .gitmap/logs/. Open the newest one and scan for the first 'Error:' or non-zero exit — that is almost always the root cause.",
      command: "# Windows\nGet-ChildItem .gitmap\\logs | Sort-Object LastWriteTime | Select-Object -Last 1\n# Unix\nls -t .gitmap/logs | head -1",
    },
  ];

  return (
    <section
      aria-labelledby="diagnostic-checklist-heading"
      className="mb-8 rounded-lg border border-border bg-muted/20 p-5"
    >
      <header className="flex items-center gap-2 mb-1">
        <ListChecks className="h-5 w-5 text-primary" />
        <h2
          id="diagnostic-checklist-heading"
          className="text-base font-heading font-semibold docs-h3"
        >
          Diagnostic checklist
        </h2>
      </header>
      <p className="text-sm text-muted-foreground mb-4">
        Run these four steps in order before searching the catalog below. Most issues are
        identified and resolved by step 2.
      </p>
      <ol className="space-y-4">
        {steps.map((step, idx) => {
          const Icon = step.icon;
          return (
            <li key={step.title} className="flex gap-3">
              <span
                aria-hidden="true"
                className="shrink-0 inline-flex items-center justify-center w-7 h-7 rounded-full bg-primary/15 text-foreground border border-primary/25 dark:bg-primary/20 dark:text-primary font-mono text-sm font-semibold"
              >
                {idx + 1}
              </span>
              <div className="flex-1 min-w-0">
                <div className="flex items-center gap-2 mb-1">
                  <Icon className="h-4 w-4 text-muted-foreground shrink-0" />
                  <h3 className="text-sm font-semibold text-foreground">{step.title}</h3>
                </div>
                <p className="text-sm text-foreground/80 mb-2">{step.body}</p>
                <ChecklistCommand command={step.command} />
              </div>
            </li>
          );
        })}
      </ol>
    </section>
  );
};

// ChecklistCommand — compact command block with copy button used inside the
// DiagnosticChecklist. Lighter than CodeBlock so the checklist stays scannable.
const ChecklistCommand = ({ command }: { command: string }) => {
  const [copied, setCopied] = useState(false);

  const handleCopy = useCallback(async () => {
    await copyToClipboard(command);
    setCopied(true);
    window.setTimeout(() => setCopied(false), 2000);
  }, [command]);

  return (
    <div className="relative group">
      <pre className="text-xs font-mono bg-background border border-border rounded p-2.5 pr-10 overflow-x-auto whitespace-pre">
        {command}
      </pre>
      <button
        type="button"
        onClick={handleCopy}
        aria-label={copied ? "Copied" : "Copy command"}
        title={copied ? "Copied!" : "Copy command"}
        className={`absolute top-1.5 right-1.5 inline-flex items-center justify-center w-7 h-7 rounded border transition-colors ${
          copied
            ? "border-primary bg-primary/15 text-foreground dark:bg-primary/20 dark:text-primary"
            : "border-border bg-background text-muted-foreground hover:text-foreground hover:border-foreground/40"
        }`}
      >
        {copied ? <Check className="h-3.5 w-3.5" /> : <Copy className="h-3.5 w-3.5" />}
      </button>
    </div>
  );
};

export default Troubleshooting;
