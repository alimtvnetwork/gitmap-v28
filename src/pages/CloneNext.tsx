import DocsLayout from "@/components/docs/DocsLayout";
import CodeBlock from "@/components/docs/CodeBlock";
import { GitBranch, FolderGit2, Trash2, Monitor, ArrowRight } from "lucide-react";

const TerminalPreview = ({ title, lines }: { title: string; lines: string[] }) => (
  <div className="rounded-lg border border-border overflow-hidden my-6">
    <div className="bg-terminal px-4 py-2 flex items-center gap-2 border-b border-border">
      <div className="flex gap-1.5">
        <span className="w-3 h-3 rounded-full bg-red-500/80" />
        <span className="w-3 h-3 rounded-full bg-yellow-500/80" />
        <span className="w-3 h-3 rounded-full bg-green-500/80" />
      </div>
      <span className="text-xs font-mono text-muted-foreground ml-2">{title}</span>
    </div>
    <div className="bg-terminal p-4 font-mono text-xs leading-relaxed overflow-x-auto">
      {lines.map((line, i) => (
        <div
          key={i}
          className={
            line.startsWith("✓")
              ? "text-green-400"
              : line.startsWith("→")
              ? "text-blue-400"
              : line.startsWith("Error")
              ? "text-red-400"
              : line.startsWith("Remove")
              ? "text-yellow-400"
              : "text-terminal-foreground"
          }
        >
          {line || "\u00A0"}
        </div>
      ))}
    </div>
  </div>
);

const VersionLogicTable = () => (
  <div className="overflow-x-auto my-6">
    <table className="w-full text-sm border border-border rounded-lg overflow-hidden">
      <thead>
        <tr className="bg-muted/50">
          <th className="text-left px-4 py-2 font-mono text-xs text-muted-foreground">Current Folder</th>
          <th className="text-left px-4 py-2 font-mono text-xs text-muted-foreground">Argument</th>
          <th className="text-left px-4 py-2 font-mono text-xs text-muted-foreground">Target Repo</th>
        </tr>
      </thead>
      <tbody className="divide-y divide-border">
        {[
          ["macro-ahk-v11", "v++", "macro-ahk-v12"],
          ["macro-ahk-v11", "v15", "macro-ahk-v15"],
          ["macro-ahk-v1", "v++", "macro-ahk-v2"],
          ["macro-ahk", "v++", "macro-ahk-v2"],
          ["macro-ahk", "v5", "macro-ahk-v5"],
          ["my-app-v100", "v++", "my-app-v101"],
        ].map(([folder, arg, target], i) => (
          <tr key={i} className="hover:bg-muted/30 transition-colors">
            <td className="px-4 py-2 font-mono text-xs text-foreground">{folder}</td>
            <td className="px-4 py-2 font-mono text-xs text-primary font-semibold">{arg}</td>
            <td className="px-4 py-2 font-mono text-xs text-green-500 font-semibold">{target}</td>
          </tr>
        ))}
      </tbody>
    </table>
  </div>
);

const FlagRow = ({ flag, description }: { flag: string; description: string }) => (
  <tr className="hover:bg-muted/30 transition-colors">
    <td className="px-4 py-2 font-mono text-xs text-primary">{flag}</td>
    <td className="px-4 py-2 text-sm text-muted-foreground">{description}</td>
  </tr>
);

const CloneNextPage = () => {
  return (
    <DocsLayout>
      <div className="max-w-4xl">
        {/* Header */}
        <div className="flex items-center gap-3 mb-2">
          <GitBranch className="w-8 h-8 text-primary" />
          <h1 className="text-3xl font-heading font-bold docs-h1">clone-next</h1>
          <span className="text-xs font-mono bg-primary/10 text-foreground border border-primary/20 px-2 py-0.5 rounded dark:bg-primary/15 dark:text-primary transition-colors duration-300 hover:border-primary/40 hover:shadow-sm hover:shadow-primary/10">cn</span>
        </div>
        <p className="text-muted-foreground mb-8 text-lg">
          Clone the next or a specific versioned iteration of the current repository. Automatically flattens
          into the base name folder (no version suffix) and tracks version history in the database.
        </p>

        {/* Usage */}
        <section className="mb-10">
          <h2 className="text-xl font-heading font-semibold mb-3 flex items-center gap-2">
            Usage
          </h2>
          <CodeBlock code="gitmap clone-next <v++|vN> [--delete] [--keep] [--no-desktop] [--create-remote] [--verbose]" />
        </section>

        {/* How it works */}
        <section className="mb-10">
          <h2 className="text-xl font-heading font-semibold mb-3 docs-h2">How It Works</h2>
          <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
            {[
              { icon: FolderGit2, title: "Detect", desc: "Reads remote origin URL and folder name from the current repo" },
              { icon: ArrowRight, title: "Compute", desc: "Parses version suffix and applies v++ or vN to build target" },
              { icon: GitBranch, title: "Clone", desc: "Removes existing base folder and clones target into flattened path" },
              { icon: Monitor, title: "Register", desc: "Adds the new clone to GitHub Desktop automatically" },
            ].map(({ icon: Icon, title, desc }) => (
              <div key={title} className="border border-border rounded-lg p-4 bg-card">
                <Icon className="w-5 h-5 text-primary mb-2" />
                <h3 className="font-mono font-semibold text-sm mb-1">{title}</h3>
                <p className="text-xs text-muted-foreground">{desc}</p>
              </div>
            ))}
          </div>
        </section>

        {/* Version logic diagram */}
        <section className="mb-10">
          <h2 className="text-xl font-heading font-semibold mb-3 docs-h2">Version Logic</h2>
          <p className="text-muted-foreground text-sm mb-4">
            The version argument determines how the target repo name is computed from the current folder.
          </p>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
            <div className="border border-border rounded-lg p-4 bg-card">
              <h3 className="font-mono font-semibold text-sm mb-2 text-primary">v++ (Increment)</h3>
              <p className="text-xs text-muted-foreground mb-3">
                Extracts the current version number from the <code className="text-primary">-vN</code> suffix and adds 1.
                If no suffix exists, the repo is treated as v1 and the target becomes v2.
              </p>
              <div className="font-mono text-xs space-y-1">
                <div className="text-muted-foreground">macro-ahk-v11 <span className="text-primary">→</span> macro-ahk-v12</div>
                <div className="text-muted-foreground">macro-ahk <span className="text-primary">→</span> macro-ahk-v2</div>
              </div>
            </div>
            <div className="border border-border rounded-lg p-4 bg-card">
              <h3 className="font-mono font-semibold text-sm mb-2 text-primary">vN (Jump)</h3>
              <p className="text-xs text-muted-foreground mb-3">
                Replaces the version suffix with the exact number specified.
                The current version is ignored — jumps directly to the target.
              </p>
              <div className="font-mono text-xs space-y-1">
                <div className="text-muted-foreground">macro-ahk-v12 + v15 <span className="text-primary">→</span> macro-ahk-v15</div>
                <div className="text-muted-foreground">macro-ahk + v5 <span className="text-primary">→</span> macro-ahk-v5</div>
              </div>
            </div>
          </div>

          <h3 className="font-mono font-semibold text-sm mb-2">Resolution Table</h3>
          <VersionLogicTable />
        </section>

        {/* Flags */}
        <section className="mb-10">
          <h2 className="text-xl font-heading font-semibold mb-3 docs-h2">Flags</h2>
          <div className="overflow-x-auto">
            <table className="w-full text-sm border border-border rounded-lg overflow-hidden">
              <thead>
                <tr className="bg-muted/50">
                  <th className="text-left px-4 py-2 font-mono text-xs text-muted-foreground">Flag</th>
                  <th className="text-left px-4 py-2 font-mono text-xs text-muted-foreground">Description</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-border">
                <FlagRow flag="--delete" description="Auto-remove current folder after clone (skip prompt)" />
                <FlagRow flag="--keep" description="Keep current folder without prompting for removal" />
                <FlagRow flag="--no-desktop" description="Skip GitHub Desktop registration" />
                <FlagRow flag="--create-remote" description="Create target GitHub repo if missing (requires GITHUB_TOKEN)" />
                <FlagRow flag="--ssh-key <name>" description="Use a named SSH key for the clone" />
                <FlagRow flag="--verbose" description="Write detailed debug log to a timestamped file" />
              </tbody>
            </table>
          </div>
        </section>

        {/* Examples */}
        <section className="mb-10">
          <h2 className="text-xl font-heading font-semibold mb-3 docs-h2">Examples</h2>

          <h3 className="font-mono text-sm font-semibold mb-2 text-muted-foreground">Increment version by one</h3>
          <TerminalPreview
            title="gitmap cn v++"
            lines={[
              "Cloning macro-ahk-v12 into D:\\wp-work\\riseup-asia...",
              "✓ Cloned macro-ahk-v12",
              "✓ Registered macro-ahk-v12 with GitHub Desktop",
              "Remove current folder macro-ahk-v11? [y/N] n",
            ]}
          />

          <h3 className="font-mono text-sm font-semibold mb-2 text-muted-foreground">Increment version (flattened)</h3>
          <TerminalPreview
            title="gitmap cn v++"
            lines={[
              "Removing existing macro-ahk for fresh clone...",
              "Cloning macro-ahk-v12 into macro-ahk (flattened)...",
              "✓ Cloned macro-ahk-v12 into macro-ahk",
              "✓ Recorded version transition v11 -> v12",
              "✓ Registered macro-ahk-v12 with GitHub Desktop",
            ]}
          />

          <h3 className="font-mono text-sm font-semibold mb-2 text-muted-foreground">Jump to specific version</h3>
          <TerminalPreview
            title="gitmap cn v15 --delete"
            lines={[
              "Cloning macro-ahk-v15 into macro-ahk (flattened)...",
              "✓ Cloned macro-ahk-v15 into macro-ahk",
              "✓ Recorded version transition v12 -> v15",
              "✓ Registered macro-ahk-v15 with GitHub Desktop",
              "✓ Removed macro-ahk-v12",
            ]}
          />

          <h3 className="font-mono text-sm font-semibold mb-2 text-muted-foreground">Repo without version suffix</h3>
          <TerminalPreview
            title="gitmap clone-next v++"
            lines={[
              "Cloning macro-ahk-v2 into macro-ahk (flattened)...",
              "✓ Cloned macro-ahk-v2 into macro-ahk",
              "✓ Recorded version transition v1 -> v2",
              "✓ Registered macro-ahk-v2 with GitHub Desktop",
            ]}
          />

          <h3 className="font-mono text-sm font-semibold mb-2 text-muted-foreground">Create remote repo before clone</h3>
          <TerminalPreview
            title="gitmap cn v15 --create-remote --delete"
            lines={[
              "Creating GitHub repo macro-ahk-v15...",
              "✓ Created GitHub repo macro-ahk-v15",
              "Cloning macro-ahk-v15 into macro-ahk (flattened)...",
              "✓ Cloned macro-ahk-v15 into macro-ahk",
              "✓ Recorded version transition v12 -> v15",
              "✓ Registered macro-ahk-v15 with GitHub Desktop",
              "✓ Removed macro-ahk-v12",
            ]}
          />

          <h3 className="font-mono text-sm font-semibold mb-2 text-muted-foreground">Lock detection when folder is in use</h3>
          <TerminalPreview
            title="gitmap cn v++ --delete"
            lines={[
              "Removing existing macro-ahk for fresh clone...",
              "Cloning macro-ahk-v12 into macro-ahk (flattened)...",
              "✓ Cloned macro-ahk-v12 into macro-ahk",
              "✓ Recorded version transition v11 -> v12",
              "✓ Registered macro-ahk-v12 with GitHub Desktop",
              "",
              "Error: remove macro-ahk-v11: access denied",
              "Scanning for processes locking macro-ahk-v11...",
              "",
              "  PID     Process",
              "  ────    ───────",
              "  14320   Code.exe",
              "  8412    explorer.exe",
              "",
              "Terminate these processes and retry? [y/N] y",
              "✓ Terminated Code.exe (PID 14320)",
              "✓ Terminated explorer.exe (PID 8412)",
              "✓ Removed macro-ahk-v11",
            ]}
          />
        </section>

        {/* URL Preservation */}
        <section className="mb-10">
          <h2 className="text-xl font-heading font-semibold mb-3 docs-h2">URL Preservation</h2>
          <p className="text-muted-foreground text-sm mb-4">
            The remote URL scheme (HTTPS or SSH) is automatically preserved from the current repo:
          </p>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="border border-border rounded-lg p-4 bg-card">
              <h3 className="font-mono font-semibold text-xs mb-2 text-muted-foreground">HTTPS</h3>
              <div className="font-mono text-xs space-y-1">
                <div className="text-muted-foreground">Current: https://github.com/user/repo-v11.git</div>
                <div className="text-green-400">Target:  https://github.com/user/repo-v12.git</div>
              </div>
            </div>
            <div className="border border-border rounded-lg p-4 bg-card">
              <h3 className="font-mono font-semibold text-xs mb-2 text-muted-foreground">SSH</h3>
              <div className="font-mono text-xs space-y-1">
                <div className="text-muted-foreground">Current: git@github.com:user/repo-v11.git</div>
                <div className="text-green-400">Target:  git@github.com:user/repo-v12.git</div>
              </div>
            </div>
          </div>
        </section>

        {/* Error handling */}
        <section className="mb-10">
          <h2 className="text-xl font-heading font-semibold mb-3 flex items-center gap-2">
            <Trash2 className="w-5 h-5" />
            Error Handling
          </h2>
          <div className="overflow-x-auto">
            <table className="w-full text-sm border border-border rounded-lg overflow-hidden">
              <thead>
                <tr className="bg-muted/50">
                  <th className="text-left px-4 py-2 font-mono text-xs text-muted-foreground">Condition</th>
                  <th className="text-left px-4 py-2 font-mono text-xs text-muted-foreground">Behavior</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-border">
                {[
                  ["Not inside a git repo", "Print error, exit 1"],
                  ["Cannot parse remote URL", "Print error, exit 1"],
                  ["Target directory already exists", "Print error with suggestion, exit 1"],
                  ["Repo creation fails (--create-remote)", "Print error, stop before clone, exit 1"],
                  ["Clone fails (network/auth)", "Print error, skip deletion, exit 1"],
                  ["Deletion fails", "Scan for locking processes via lockcheck"],
                  ["Locking processes found", "Prompt to terminate, retry deletion"],
                  ["Lock scan fails", "Print warning, exit 0 (clone succeeded)"],
                  ["Process termination fails", "Print warning, exit 0 (clone succeeded)"],
                ].map(([cond, behavior], i) => (
                  <tr key={i} className="hover:bg-muted/30 transition-colors">
                    <td className="px-4 py-2 text-xs text-foreground">{cond}</td>
                    <td className="px-4 py-2 text-xs text-muted-foreground">{behavior}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </section>

        {/* File Layout */}
        <section className="mb-10">
          <h2 className="text-xl font-heading font-semibold mb-3 docs-h2">File Layout</h2>
          <div className="overflow-x-auto">
            <table className="w-full text-sm border border-border rounded-lg overflow-hidden">
              <thead>
                <tr className="bg-muted/50">
                  <th className="text-left px-4 py-2 font-mono text-xs text-muted-foreground">File</th>
                  <th className="text-left px-4 py-2 font-mono text-xs text-muted-foreground">Purpose</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-border">
                {[
                  ["cmd/clonenext.go", "Flag parsing, orchestration, deletion + cd"],
                  ["constants/constants_clonenext.go", "Command names, messages, error strings"],
                  ["lockcheck/lockcheck.go", "LockingProcess struct, FindLockingProcesses interface"],
                  ["lockcheck/lockcheck_windows.go", "handle.exe + WMI fallback, KillProcess via taskkill"],
                  ["lockcheck/lockcheck_unix.go", "lsof-based lock detection, KillProcess via kill(2)"],
                  ["helptext/clone-next.md", "Embedded help text for --help flag"],
                ].map(([file, purpose], i) => (
                  <tr key={i} className="hover:bg-muted/30 transition-colors">
                    <td className="px-4 py-2 font-mono text-xs text-primary">{file}</td>
                    <td className="px-4 py-2 text-xs text-muted-foreground">{purpose}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </section>

        {/* See also */}
        <section className="mb-10">
          <h2 className="text-xl font-heading font-semibold mb-3 docs-h2">See Also</h2>
          <ul className="space-y-1 text-sm">
            <li><a href="/commands" className="text-primary hover:underline font-mono">clone</a> — Clone repos from structured file</li>
            <li><a href="/ssh" className="text-primary hover:underline font-mono">ssh</a> — Manage named SSH keys</li>
            <li><a href="/commands" className="text-primary hover:underline font-mono">desktop-sync</a> — Sync repos to GitHub Desktop</li>
          </ul>
        </section>
      </div>
    </DocsLayout>
  );
};

export default CloneNextPage;
