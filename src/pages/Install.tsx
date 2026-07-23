import DocsLayout from "@/components/docs/DocsLayout";
import CodeBlock from "@/components/docs/CodeBlock";
import TerminalDemo from "@/components/docs/TerminalDemo";
import InstallHelpSection from "@/components/docs/InstallHelpSection";
import { Download, Trash2, Database, Wrench, FolderDown, Monitor, Terminal, Shield, FileText, AlertTriangle } from "lucide-react";

const terminalLines = [
  { text: "gitmap install --list", type: "input" as const, delay: 800 },
  { text: "", type: "output" as const },
  { text: "  Core Tools:", type: "header" as const },
  { text: "  vscode              Visual Studio Code editor", type: "output" as const },
  { text: "  node                Node.js JavaScript runtime", type: "output" as const },
  { text: "  go                  Go programming language", type: "output" as const },
  { text: "  git                 Git version control", type: "output" as const },
  { text: "  python              Python programming language", type: "output" as const },
  { text: "", type: "output" as const },
  { text: "  Databases:", type: "header" as const },
  { text: "  postgresql          PostgreSQL relational database", type: "output" as const },
  { text: "  redis               Redis in-memory key-value store", type: "output" as const },
  { text: "  mongodb             MongoDB document database", type: "output" as const },
  { text: "", type: "output" as const },
  { text: "gitmap install node", type: "input" as const, delay: 1000 },
  { text: "", type: "output" as const },
  { text: "  Checking if node is installed...", type: "output" as const },
  { text: "  node is not installed.", type: "output" as const },
  { text: "", type: "output" as const },
  { text: "  ┌─ Install Plan ─────────────────────", type: "header" as const },
  { text: "  │ Tool:    node", type: "output" as const },
  { text: "  │ Version: latest", type: "output" as const },
  { text: "  │ Manager: choco", type: "output" as const },
  { text: "  │ Command: choco install nodejs -y --no-progress", type: "output" as const },
  { text: "  └────────────────────────────────────", type: "output" as const },
  { text: "", type: "output" as const },
  { text: "  [1/4] Installing node via choco...", type: "output" as const },
  { text: "  ✓ node install command completed successfully.", type: "accent" as const },
  { text: "  [3/4] Verifying installation...", type: "output" as const },
  { text: "  ✓ node installed successfully.", type: "accent" as const },
  { text: "  → Detected version: v22.5.0", type: "output" as const },
  { text: "", type: "output" as const },
  { text: "  Install summary", type: "header" as const },
  { text: "    Version: v22.5.0", type: "output" as const },
  { text: "    Binary: C:\\Program Files\\nodejs\\node.exe", type: "output" as const },
  { text: "    Install dir: C:\\Program Files\\nodejs", type: "output" as const },
  { text: "    PATH target: User PATH (already present)", type: "output" as const },
  { text: "  [4/4] Recording installation...", type: "accent" as const },
  { text: "  Recorded node v22.5.0 in database.", type: "accent" as const },
];

const FlagTable = ({ flags }: { flags: [string, string][] }) => (
  <div className="overflow-x-auto my-4">
    <table className="w-full text-sm border border-border rounded-lg overflow-hidden">
      <thead>
        <tr className="bg-muted/50">
          <th className="text-left px-4 py-2 font-mono text-xs text-muted-foreground">Flag</th>
          <th className="text-left px-4 py-2 font-mono text-xs text-muted-foreground">Description</th>
        </tr>
      </thead>
      <tbody className="divide-y divide-border">
        {flags.map(([flag, desc], i) => (
          <tr key={i} className="hover:bg-muted/30 transition-colors">
            <td className="px-4 py-2 font-mono text-xs text-primary">{flag}</td>
            <td className="px-4 py-2 text-sm text-muted-foreground">{desc}</td>
          </tr>
        ))}
      </tbody>
    </table>
  </div>
);

const coreTools: [string, string][] = [
  ["vscode", "Visual Studio Code editor"],
  ["node", "Node.js JavaScript runtime"],
  ["yarn", "Yarn package manager"],
  ["bun", "Bun JavaScript runtime"],
  ["pnpm", "pnpm package manager"],
  ["python", "Python programming language"],
  ["go", "Go programming language"],
  ["git", "Git version control"],
  ["git-lfs", "Git Large File Storage"],
  ["gh", "GitHub CLI"],
  ["github-desktop", "GitHub Desktop application"],
  ["cpp", "C++ compiler (MinGW/g++)"],
  ["php", "PHP programming language"],
  ["powershell", "PowerShell shell"],
  ["chocolatey", "Chocolatey package manager"],
  ["winget", "Winget package manager"],
  ["scripts", "Clone gitmap scripts to local folder"],
];

const dbTools: [string, string][] = [
  ["mysql", "MySQL relational database"],
  ["mariadb", "MariaDB (MySQL-compatible fork)"],
  ["postgresql", "PostgreSQL relational database"],
  ["sqlite", "SQLite embedded database"],
  ["mongodb", "MongoDB document database"],
  ["couchdb", "CouchDB document database (REST API)"],
  ["redis", "Redis in-memory key-value store"],
  ["cassandra", "Apache Cassandra wide-column NoSQL"],
  ["neo4j", "Neo4j graph database"],
  ["elasticsearch", "Elasticsearch search and analytics"],
  ["duckdb", "DuckDB analytical columnar database"],
];

const managers: [string, string, string][] = [
  ["choco", "Chocolatey", "Windows"],
  ["winget", "Winget", "Windows"],
  ["apt", "APT", "Debian / Ubuntu"],
  ["brew", "Homebrew", "macOS / Linux"],
  ["snap", "Snap", "Linux"],
  ["dnf", "DNF", "Fedora / RHEL"],
  ["pacman", "Pacman", "Arch Linux"],
];

const ToolTable = ({ tools, category }: { tools: [string, string][]; category: string }) => (
  <div className="mb-6">
    <h3 className="font-mono font-semibold text-sm mb-2 flex items-center gap-2">
      {category === "Core Tools" ? <Wrench className="h-4 w-4 text-primary" /> : <Database className="h-4 w-4 text-primary" />}
      {category}
    </h3>
    <div className="overflow-x-auto">
      <table className="w-full text-sm border border-border rounded-lg overflow-hidden">
        <thead>
          <tr className="bg-muted/50">
            <th className="text-left px-4 py-2 font-mono text-xs text-muted-foreground">Tool Name</th>
            <th className="text-left px-4 py-2 font-mono text-xs text-muted-foreground">Description</th>
          </tr>
        </thead>
        <tbody className="divide-y divide-border">
          {tools.map(([name, desc], i) => (
            <tr key={i} className="hover:bg-muted/30 transition-colors">
              <td className="px-4 py-2 font-mono text-xs text-primary font-semibold">{name}</td>
              <td className="px-4 py-2 text-xs text-muted-foreground">{desc}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  </div>
);

const InstallPage = () => {
  return (
    <DocsLayout>
      <div className="max-w-4xl">
        {/* Header */}
        <div className="flex items-center gap-3 mb-2">
          <Download className="w-8 h-8 text-primary" />
          <h1 className="text-3xl font-heading font-bold docs-h1">install / uninstall</h1>
          <span className="text-xs font-mono bg-primary/10 text-foreground border border-primary/20 px-2 py-0.5 rounded dark:bg-primary/15 dark:text-primary transition-colors duration-300 hover:border-primary/40 hover:shadow-sm hover:shadow-primary/10">in</span>
          <span className="text-xs font-mono bg-primary/10 text-foreground border border-primary/20 px-2 py-0.5 rounded dark:bg-primary/15 dark:text-primary transition-colors duration-300 hover:border-primary/40 hover:shadow-sm hover:shadow-primary/10">un</span>
        </div>
        <p className="text-muted-foreground mb-8 text-lg">
          Install and manage developer tools and databases with automatic version tracking.
        </p>

        {/* Terminal Demo */}
        <div className="mb-10">
          <TerminalDemo title="gitmap — install tools" lines={terminalLines} autoPlay />
        </div>

        {/* Usage */}
        <section className="mb-10">
          <h2 className="text-xl font-heading font-semibold mb-3">Usage</h2>
          <CodeBlock code={`gitmap install <tool> [flags]
gitmap uninstall <tool> [flags]`} />
        </section>

        <InstallHelpSection />

        {/* How it works */}
        <section className="mb-10">
          <h2 className="text-xl font-heading font-semibold mb-3 docs-h2">How It Works</h2>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            {[
              { icon: Download, title: "Detect & Install", desc: "Resolves the platform package manager and installs the tool with a single command" },
              { icon: Database, title: "Track in SQLite", desc: "Records tool name, version (major.minor.patch.build), manager, and timestamps in InstalledTools" },
              { icon: Trash2, title: "Uninstall", desc: "Resolves the original manager from the database and removes the tool cleanly" },
            ].map(({ icon: Icon, title, desc }) => (
              <div key={title} className="border border-border rounded-lg p-4 bg-card">
                <Icon className="w-5 h-5 text-primary mb-2" />
                <h3 className="font-mono font-semibold text-sm mb-1">{title}</h3>
                <p className="text-xs text-muted-foreground">{desc}</p>
              </div>
            ))}
          </div>
        </section>

        {/* v2.65.0 Install UX */}
        <section className="mb-10">
          <h2 className="text-xl font-heading font-semibold mb-3 docs-h2">
            <span className="flex items-center gap-2"><Shield className="h-5 w-5" /> v2.65.0 — Install UX</span>
          </h2>
          <p className="text-muted-foreground text-sm mb-4">
            v2.65.0 overhauled the install experience with structured output, GUI-safe verification, and detailed error logging.
          </p>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
            {[
              {
                icon: FileText,
                title: "Install Plan Box",
                desc: "Every install starts with a structured plan showing tool, version, manager, and the exact command before execution.",
              },
              {
                icon: Download,
                title: "Numbered Steps",
                desc: "Progress is shown as [1/4] Update → [2/4] Install → [3/4] Verify → [4/4] Record for clear tracking.",
              },
              {
                icon: Shield,
                title: "GUI-Safe Verification",
                desc: "GUI tools (Notepad++, GitHub Desktop) skip --version checks that would open a window and block the terminal. Only exe path verification is used.",
              },
              {
                icon: AlertTriangle,
                title: "Error Logging",
                desc: "On failure, a detailed log is written to .gitmap/logs/<tool>-error-<timestamp>.log with version, command, output, and error reason.",
              },
              {
                icon: Terminal,
                title: "Install Summary",
                desc: "Installers now print the installed version, binary path, install directory, and PATH target/status so users know exactly what changed.",
              },
            ].map(({ icon: Icon, title, desc }) => (
              <div key={title} className="border border-border rounded-lg p-4 bg-card">
                <Icon className="w-5 h-5 text-primary mb-2" />
                <h3 className="font-mono font-semibold text-sm mb-1">{title}</h3>
                <p className="text-xs text-muted-foreground">{desc}</p>
              </div>
            ))}
          </div>

          <h3 className="font-mono font-semibold text-sm mb-2">Silent Install Flags</h3>
          <p className="text-muted-foreground text-xs mb-3">
            GUI applications use silent flags to prevent blocking the terminal during installation:
          </p>
          <div className="overflow-x-auto mb-4">
            <table className="w-full text-sm border border-border rounded-lg overflow-hidden">
              <thead>
                <tr className="bg-muted/50">
                  <th className="text-left px-4 py-2 font-mono text-xs text-muted-foreground">Manager</th>
                  <th className="text-left px-4 py-2 font-mono text-xs text-muted-foreground">Silent Flag</th>
                  <th className="text-left px-4 py-2 font-mono text-xs text-muted-foreground">Purpose</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-border">
                {[
                  ["Chocolatey", "--no-progress", "Suppresses progress bar UI during download"],
                  ["Winget", "--silent", "Runs installer without GUI interaction"],
                  ["APT", "-y", "Auto-confirms prompts without blocking"],
                ].map(([mgr, flag, purpose], i) => (
                  <tr key={i} className="hover:bg-muted/30 transition-colors">
                    <td className="px-4 py-2 text-sm text-foreground">{mgr}</td>
                    <td className="px-4 py-2 font-mono text-xs text-primary">{flag}</td>
                    <td className="px-4 py-2 text-xs text-muted-foreground">{purpose}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          <h3 className="font-mono font-semibold text-sm mb-2">Error Log Format</h3>
          <CodeBlock code={`gitmap install error log
========================

Tool:            npp
Version:         latest
Package Manager: choco
Command:         choco install notepadplusplus -y --no-progress
Timestamp:       2026-04-08T14:32:01Z
Error:           exit status 1

--- Installer Output ---

Chocolatey v2.4.0
Installing notepadplusplus...`} />
        </section>

        {/* Install Flags */}
        <section className="mb-10">
          <h2 className="text-xl font-heading font-semibold mb-3 docs-h2">
            <span className="flex items-center gap-2"><Download className="h-5 w-5" /> Install Flags</span>
          </h2>
          <FlagTable flags={[
            ["--manager <name>", "Force package manager (choco, winget, apt, brew, snap)"],
            ["--version <ver>", "Install a specific version"],
            ["--verbose", "Show full installer output"],
            ["--dry-run", "Show install command without executing"],
            ["--check", "Only check if tool is installed"],
            ["--list", "List all supported tools"],
            ["--status", "Show installed tools from database"],
            ["--upgrade", "Upgrade an already-installed tool"],
          ]} />
        </section>

        {/* Uninstall Flags */}
        <section className="mb-10">
          <h2 className="text-xl font-heading font-semibold mb-3 docs-h2">
            <span className="flex items-center gap-2"><Trash2 className="h-5 w-5" /> Uninstall Flags</span>
          </h2>
          <FlagTable flags={[
            ["--dry-run", "Show uninstall command without executing"],
            ["--force", "Skip confirmation prompt"],
            ["--purge", "Remove config files too"],
          ]} />
        </section>

        {/* Supported Tools */}
        <section className="mb-10">
          <h2 className="text-xl font-heading font-semibold mb-3 docs-h2">Supported Tools</h2>
          <p className="text-muted-foreground text-sm mb-4">
            {coreTools.length + dbTools.length} tools across {managers.length} package managers.
          </p>
          <ToolTable tools={coreTools} category="Core Tools" />
          <ToolTable tools={dbTools} category="Databases" />
        </section>

        {/* Scripts */}
        <section className="mb-10">
          <h2 className="text-xl font-heading font-semibold mb-3 docs-h2">
            <span className="flex items-center gap-2"><FolderDown className="h-5 w-5" /> Install Scripts</span>
          </h2>
          <p className="text-muted-foreground text-sm mb-4">
            Clone all gitmap utility scripts to a local folder with <code className="text-primary">gitmap install scripts</code>.
            The scripts are shallow-cloned from the repository and copied to a platform-specific directory.
          </p>

          {/* Platform paths */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
            <div className="border border-border rounded-lg p-4 bg-card">
              <div className="flex items-center gap-2 mb-2">
                <Monitor className="w-4 h-4 text-primary" />
                <h3 className="font-mono font-semibold text-sm">Windows</h3>
              </div>
              <p className="text-xs text-muted-foreground mb-2">
                Reads the deploy drive from <code className="text-primary">powershell.json</code> → <code className="text-primary">deployPath</code>.
                Falls back to <code className="text-primary">D:\gitmap-scripts</code>.
              </p>
              <CodeBlock code={`D:\\gitmap-scripts\\
├── install.ps1
├── uninstall.ps1
├── Get-LastRelease.ps1
└── run.ps1`} />
            </div>
            <div className="border border-border rounded-lg p-4 bg-card">
              <div className="flex items-center gap-2 mb-2">
                <Terminal className="w-4 h-4 text-primary" />
                <h3 className="font-mono font-semibold text-sm">Linux / macOS</h3>
              </div>
              <p className="text-xs text-muted-foreground mb-2">
                Installs to <code className="text-primary">~/Desktop/gitmap-scripts</code>.
              </p>
              <CodeBlock code={`~/Desktop/gitmap-scripts/
├── install.sh
├── install.ps1
├── run.sh
├── run.ps1
├── uninstall.ps1
└── Get-LastRelease.ps1`} />
            </div>
          </div>

          {/* Copied files table */}
          <h3 className="font-mono font-semibold text-sm mb-2">Copied Files</h3>
          <div className="overflow-x-auto">
            <table className="w-full text-sm border border-border rounded-lg overflow-hidden">
              <thead>
                <tr className="bg-muted/50">
                  <th className="text-left px-4 py-2 font-mono text-xs text-muted-foreground">File</th>
                  <th className="text-left px-4 py-2 font-mono text-xs text-muted-foreground">Source</th>
                  <th className="text-left px-4 py-2 font-mono text-xs text-muted-foreground">Description</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-border">
                {[
                  ["install.ps1", "gitmap-v27/scripts/", "PowerShell one-liner installer for Windows"],
                  ["install.sh", "gitmap-v27/scripts/", "Bash one-liner installer for Linux/macOS"],
                  ["uninstall.ps1", "gitmap-v27/scripts/", "PowerShell uninstaller"],
                  ["Get-LastRelease.ps1", "gitmap-v27/scripts/", "Resolve latest release version (3-tier fallback)"],
                  ["run.ps1", "repo root", "PowerShell build, deploy, and self-update script"],
                  ["run.sh", "repo root", "Bash build, deploy, and self-update script"],
                ].map(([file, source, desc], i) => (
                  <tr key={i} className="hover:bg-muted/30 transition-colors">
                    <td className="px-4 py-2 font-mono text-xs text-primary font-semibold">{file}</td>
                    <td className="px-4 py-2 font-mono text-xs text-muted-foreground">{source}</td>
                    <td className="px-4 py-2 text-xs text-muted-foreground">{desc}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          {/* Example */}
          <div className="mt-4">
            <CodeBlock code={`$ gitmap install scripts
  → Scripts target: /home/alim/Desktop/gitmap-scripts
  Cloning gitmap repo for scripts...
  ✓ Copied: install.ps1
  ✓ Copied: install.sh
  ✓ Copied: run.ps1
  ✓ Copied: run.sh
  ✓ Copied: uninstall.ps1
  ✓ Copied: Get-LastRelease.ps1

  ✅ 6 scripts installed to /home/alim/Desktop/gitmap-scripts`} />
          </div>
        </section>

        <section className="mb-10">
          <h2 className="text-xl font-heading font-semibold mb-3 docs-h2">Installer Summary Output</h2>
          <p className="text-muted-foreground text-sm mb-4">
            Every installer now ends with a clear summary so you can see the installed version, exact binary location, and where the PATH change was applied.
          </p>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="border border-border rounded-lg p-4 bg-card">
              <h3 className="font-mono font-semibold text-sm mb-2">PowerShell / Windows</h3>
              <CodeBlock code={`Install summary
    Version: v2.65.0
    Binary: C:\\Users\\me\\AppData\\Local\\gitmap-v27\\gitmap.exe
    Install Dir: C:\\Users\\me\\AppData\\Local\\gitmap-v27
    PATH Target: User PATH (added)
    Session PATH: refreshed for current PowerShell session`} />
            </div>
            <div className="border border-border rounded-lg p-4 bg-card">
              <h3 className="font-mono font-semibold text-sm mb-2">Unix / macOS</h3>
              <CodeBlock code={`Install summary
    Version: v2.65.0
    Binary: /Users/me/.local/bin/gitmap-v27
    Install dir: /Users/me/.local/bin
    Shell: zsh
    PATH target: /Users/me/.zshrc (added)
    Reload: . /Users/me/.zshrc`} />
            </div>
          </div>
        </section>

        {/* Package Managers */}
        <section className="mb-10">
          <h2 className="text-xl font-heading font-semibold mb-3 docs-h2">Package Managers</h2>
          <p className="text-muted-foreground text-sm mb-4">
            gitmap auto-detects the best manager for your platform, or use <code className="text-primary">--manager</code> to override.
          </p>
          <div className="overflow-x-auto">
            <table className="w-full text-sm border border-border rounded-lg overflow-hidden">
              <thead>
                <tr className="bg-muted/50">
                  <th className="text-left px-4 py-2 font-mono text-xs text-muted-foreground">ID</th>
                  <th className="text-left px-4 py-2 font-mono text-xs text-muted-foreground">Manager</th>
                  <th className="text-left px-4 py-2 font-mono text-xs text-muted-foreground">Platform</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-border">
                {managers.map(([id, name, platform], i) => (
                  <tr key={i} className="hover:bg-muted/30 transition-colors">
                    <td className="px-4 py-2 font-mono text-xs text-primary font-semibold">{id}</td>
                    <td className="px-4 py-2 text-sm text-foreground">{name}</td>
                    <td className="px-4 py-2 text-xs text-muted-foreground">{platform}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </section>

        {/* SQLite Tracking */}
        <section className="mb-10">
          <h2 className="text-xl font-heading font-semibold mb-3 docs-h2">SQLite Tracking</h2>
          <p className="text-muted-foreground text-sm mb-4">
            Every install is recorded in the <code className="text-primary">InstalledTools</code> table for version comparison and uninstall resolution.
          </p>
          <CodeBlock code={`CREATE TABLE InstalledTools (
  ID             INTEGER PRIMARY KEY AUTOINCREMENT,
  Tool           TEXT NOT NULL,
  VersionMajor   INTEGER DEFAULT 0,
  VersionMinor   INTEGER DEFAULT 0,
  VersionPatch   INTEGER DEFAULT 0,
  VersionBuild   INTEGER DEFAULT 0,
  VersionString  TEXT DEFAULT '',
  PackageManager TEXT DEFAULT '',
  InstallPath    TEXT DEFAULT '',
  InstalledAt    TEXT DEFAULT (datetime('now')),
  UpdatedAt      TEXT DEFAULT (datetime('now'))
);`} />
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
                  ["constants/constants_install.go", "Tool names, package IDs, flag descriptions, messages"],
                  ["constants/constants_installedtools.go", "InstalledTools table SQL and column constants"],
                  ["cmd/install.go", "Flag parsing, manager resolution, install orchestration"],
                  ["cmd/uninstall.go", "Uninstall flow with confirmation and DB cleanup"],
                  ["store/installedtool.go", "CRUD operations, version parsing, comparison"],
                  ["helptext/install.md", "Embedded help text for --help flag"],
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
            <li><a href="/doctor" className="text-primary hover:underline font-mono">doctor</a> — Diagnose PATH, deploy, and version issues</li>
            <li><a href="/setup" className="text-primary hover:underline font-mono">setup</a> — Configure Git settings and shell completions</li>
            <li><a href="/commands" className="text-primary hover:underline font-mono">env</a> — Check environment variables and tool availability</li>
          </ul>
        </section>
      </div>
    </DocsLayout>
  );
};

export default InstallPage;
