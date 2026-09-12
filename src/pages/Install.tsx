import DocsLayout from "@/components/docs/DocsLayout";
import CodeBlock from "@/components/docs/CodeBlock";
import TerminalDemo from "@/components/docs/TerminalDemo";
import InstallHelpSection from "@/components/docs/InstallHelpSection";
import { Download, Trash2, Database, Wrench, FolderDown, Monitor, Terminal, Shield, FileText, AlertTriangle, Cpu, Bot, Server, TerminalSquare, Layers, Sparkles } from "lucide-react";

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
  ["gitmap", "Gitmap core repository orchestration CLI binary"],
  ["vscode", "Visual Studio Code editor"],
  ["node", "Node.js JavaScript runtime"],
  ["yarn", "Yarn package manager"],
  ["bun", "Bun JavaScript runtime"],
  ["pnpm", "pnpm package manager"],
  ["python", "Python programming language"],
  ["go", "Go programming language"],
  ["git", "Git version control system"],
  ["git-lfs", "Git Large File Storage extension"],
  ["gh", "GitHub CLI"],
  ["github-desktop", "GitHub Desktop GUI application"],
  ["cpp", "C++ compiler toolchain (MinGW / GCC)"],
  ["php", "PHP programming language"],
  ["powershell", "PowerShell cross-platform shell (pwsh)"],
  ["chocolatey", "Chocolatey package manager for Windows"],
  ["winget", "Windows Package Manager CLI"],
  ["chrome", "Google Chrome web browser"],
  ["build-essential", "Ubuntu compiler toolchain & dev libraries (gcc, g++, make)"],
  ["wordpress", "WordPress web publishing platform & WP-CLI"],
  ["antigravity", "Antigravity CLI autonomous coding assistant"],
  ["ag-manager", "Antigravity Manager GUI desktop application"],
  ["ag-ctx", "Add Antigravity to system right-click context menus"],
  ["npp", "Notepad++ text editor with synced settings"],
  ["npp-settings", "Sync Notepad++ configuration settings only"],
  ["install-npp", "Install Notepad++ binary only without settings"],
  ["vscode-settings", "Sync VS Code settings, keybindings, and extensions"],
  ["obs", "OBS Studio screen recorder and streaming suite"],
  ["obs-settings", "Sync OBS Studio profiles and scenes"],
  ["wt-settings", "Sync Windows Terminal configurations"],
  ["dbeaver", "Universal database GUI client"],
  ["sticky-notes", "Sticky Notes desktop utility"],
  ["scripts", "Clone gitmap scripts to local folder"],
];

const dbTools: [string, string][] = [
  ["mysql", "MySQL relational database server (port 3306)"],
  ["mariadb", "MariaDB open-source MySQL-compatible database (port 3306)"],
  ["postgresql", "PostgreSQL object-relational database system (port 5432)"],
  ["sqlite", "SQLite embedded zero-configuration database"],
  ["mongodb", "MongoDB document-oriented NoSQL database (port 27017)"],
  ["couchdb", "Apache CouchDB document database with JSON REST API (port 5984)"],
  ["redis", "Redis in-memory key-value data structure store (port 6379)"],
  ["cassandra", "Apache Cassandra distributed wide-column NoSQL (port 9042)"],
  ["neo4j", "Neo4j native graph database platform (port 7474 / 7687)"],
  ["elasticsearch", "Elasticsearch distributed search and analytics engine (port 9200)"],
  ["duckdb", "DuckDB in-process columnar analytical SQL database"],
  ["litedb", "LiteDB embedded NoSQL document store for .NET"],
];

const langTools: [string, string][] = [
  ["rust", "Rust programming language compiler & Cargo toolchain"],
  ["dotnet", ".NET SDK and developer multi-platform runtime"],
  ["java", "OpenJDK Java Development Kit and runtime"],
  ["flutter", "Flutter SDK for cross-platform applications"],
  ["laravel", "Laravel PHP framework CLI installer"],
  ["composer", "Composer dependency manager for PHP"],
];

const aiTools: [string, string][] = [
  ["ollama", "Ollama local large language model runner & server"],
  ["llama-cpp", "llama.cpp high-performance local LLM inference engine"],
  ["python-libs", "Core AI/ML Python stack (numpy, pandas, torch, transformers)"],
  ["antigravity", "Antigravity CLI autonomous coding assistant"],
  ["ag-manager", "Antigravity Manager GUI desktop application"],
];

const devopsTools: [string, string][] = [
  ["docker", "Docker container platform and runtime engine"],
  ["kubernetes", "Kubernetes container orchestration CLI (kubectl)"],
  ["jenkins", "Jenkins continuous integration and automation server"],
  ["nginx", "Nginx high-performance HTTP server and reverse proxy"],
  ["vmware", "VMware Workstation virtualization platform"],
  ["open-vm-tools", "Open Virtual Machine Tools for VMware guest optimization"],
];

const utilityTools: [string, string][] = [
  ["zsh", "Z shell interactive command interpreter"],
  ["flameshot", "Flameshot screenshot capture and annotation tool"],
  ["conemu", "ConEmu Windows console emulator with tabs and splits"],
  ["vlc", "VLC media player cross-platform multimedia player"],
  ["qbittorrent", "qBittorrent free and open-source BitTorrent client"],
  ["utorrent", "uTorrent lightweight BitTorrent client"],
];

const customScripts: [string, string, string, string][] = [
  ["scripts-fixer", "scripts-fixer-v20", "PowerShell / Bash", "Repository scripts and path fixer / auto-repair suite"],
  ["coding-guidelines", "coding-guidelines-v24", "PowerShell / Bash", "AlimTV Network Coding Guidelines (v24) automated compliance installer"],
  ["macro-ahk", "macro-ahk-v55", "PowerShell / Bash", "AutoHotkey v2 automation and shortcut extension suite"],
];

const profilesData: [string, string, string, string][] = [
  ["minimal", "min, basic", "Editor, Git, Node.js, and Python", "vscode, git, node, python"],
  ["base", "ubuntu-basic, ub", "Essential OS compiler toolchain and shell", "curl, git, build-essential, zsh"],
  ["dev", "developer, dev-stack", "Standard dev workstation + runtimes + AI suite", "vscode, git, python, node, pnpm, go, rust, php, antigravity, ag-manager"],
  ["small-dev", "ub+sdev, ubuntu-small-dev", "Lightweight dev suite with runtimes", "ubuntu-basic + vscode + go + node"],
  ["advance", "ub+dev, ubuntu+dev", "Full developer workstation suite", "small-dev + docker + python3"],
  ["dev-advance", "—", "Advanced developer stack with databases & AI", "dev stack + databases + LLM inference"],
  ["terminal", "—", "Terminal power-user environment", "zsh, conemu, git, curl, wget, nano, vim"],
  ["ubuntu", "ubuntu-dev, linux-dev", "Compiler toolchain, shell, browsers, runtimes", "build-essential, git, zsh, vscode, chrome, node, python, go, antigravity, ag-manager"],
  ["ai", "ai-dev, ml, llm", "Local LLM runners, Python ML libs & Antigravity", "python, ollama, llama-cpp, python-libs, antigravity, ag-manager"],
  ["backend", "back, server", "Minimal stack + databases, Docker & languages", "minimal + docker, mysql, postgresql, redis, go, dotnet, java"],
  ["fullstack", "full, web", "Backend + pnpm, PHP, Composer, MongoDB & CI/CD", "backend + pnpm, php, composer, mongodb, jenkins"],
  ["web-dev", "—", "Modern full-stack web developer environment", "node, pnpm, yarn, bun, php, composer, mysql"],
  ["devops", "—", "Containers, clusters, proxies & virtualization", "docker, kubernetes, jenkins, nginx, vmware, open-vm-tools"],
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

const getCategoryIcon = (category: string) => {
  switch (category) {
    case "Core Tools":
      return <Wrench className="h-4 w-4 text-primary" />;
    case "Databases":
      return <Database className="h-4 w-4 text-primary" />;
    case "Languages & Runtimes":
      return <Cpu className="h-4 w-4 text-primary" />;
    case "Local AI":
      return <Bot className="h-4 w-4 text-primary" />;
    case "DevOps & Containers":
      return <Server className="h-4 w-4 text-primary" />;
    case "Terminal & Utilities":
      return <TerminalSquare className="h-4 w-4 text-primary" />;
    default:
      return <Wrench className="h-4 w-4 text-primary" />;
  }
};

const ToolTable = ({ tools, category }: { tools: [string, string][]; category: string }) => (
  <div className="mb-6">
    <h3 className="font-mono font-semibold text-sm mb-2 flex items-center gap-2">
      {getCategoryIcon(category)}
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
            ["--manager <name>", "Force package manager (choco, winget, apt, brew, snap, dnf, pacman)"],
            ["--version <ver>", "Install a specific version"],
            ["--verbose", "Show full installer output"],
            ["--dry-run", "Show install command without executing"],
            ["--check", "Only check if tool is installed"],
            ["--tree, -t", "Preview full tool hierarchy/tree of a profile before installing"],
            ["--list", "List all supported tools and profiles with installation status"],
            ["--yes, -y", "Auto-confirm installation prompts without interactive questions"],
            ["--explain", "Print exact resolved shell command before executing"],
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

        {/* Installation Profiles */}
        <section className="mb-10">
          <h2 className="text-xl font-heading font-semibold mb-3 docs-h2">
            <span className="flex items-center gap-2"><Layers className="h-5 w-5" /> Installation Profiles</span>
          </h2>
          <p className="text-muted-foreground text-sm mb-4">
            Profiles bundle multiple developer tools and packages into a single command. Missing tools are automatically resolved while already-installed components are skipped.
          </p>

          <div className="overflow-x-auto mb-6">
            <table className="w-full text-sm border border-border rounded-lg overflow-hidden">
              <thead>
                <tr className="bg-muted/50">
                  <th className="text-left px-4 py-2 font-mono text-xs text-muted-foreground">Profile</th>
                  <th className="text-left px-4 py-2 font-mono text-xs text-muted-foreground">Aliases</th>
                  <th className="text-left px-4 py-2 font-mono text-xs text-muted-foreground">Description</th>
                  <th className="text-left px-4 py-2 font-mono text-xs text-muted-foreground">Key Tools Included</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-border">
                {profilesData.map(([name, aliases, desc, tools], i) => (
                  <tr key={i} className="hover:bg-muted/30 transition-colors">
                    <td className="px-4 py-2 font-mono text-xs text-primary font-semibold">{name}</td>
                    <td className="px-4 py-2 font-mono text-xs text-muted-foreground">{aliases}</td>
                    <td className="px-4 py-2 text-xs text-foreground">{desc}</td>
                    <td className="px-4 py-2 font-mono text-xs text-muted-foreground">{tools}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          <h3 className="font-mono font-semibold text-sm mb-2">Profile Usage & Tree Preview</h3>
          <CodeBlock code={`# List all available profiles and installed tool count
$ gitmap install profile
$ gitmap install profile --list

# Preview tool hierarchy before installing (--tree / -t)
$ gitmap install profile dev --tree
$ gitmap in dev -t

# Install profile directly
$ gitmap install profile dev
$ gitmap in dev -y
$ gitmap install ubuntu --dry-run`} />
        </section>

        {/* Supported Tools */}
        <section className="mb-10">
          <h2 className="text-xl font-heading font-semibold mb-3 docs-h2">Supported Tools</h2>
          <p className="text-muted-foreground text-sm mb-4">
            {coreTools.length + dbTools.length + langTools.length + aiTools.length + devopsTools.length + utilityTools.length} tools across 6 distinct categories and {managers.length} package managers.
          </p>
          <ToolTable tools={coreTools} category="Core Tools" />
          <ToolTable tools={dbTools} category="Databases" />
          <ToolTable tools={langTools} category="Languages & Runtimes" />
          <ToolTable tools={aiTools} category="Local AI" />
          <ToolTable tools={devopsTools} category="DevOps & Containers" />
          <ToolTable tools={utilityTools} category="Terminal & Utilities" />
        </section>

        {/* Custom Script Tools */}
        <section className="mb-10">
          <h2 className="text-xl font-heading font-semibold mb-3 docs-h2">
            <span className="flex items-center gap-2"><Sparkles className="h-5 w-5" /> Custom Scripts & Automation Suites</span>
          </h2>
          <p className="text-muted-foreground text-sm mb-4">
            Automated bootstrapping one-liners for specialized developer utilities and governance suites.
          </p>

          <div className="overflow-x-auto mb-6">
            <table className="w-full text-sm border border-border rounded-lg overflow-hidden">
              <thead>
                <tr className="bg-muted/50">
                  <th className="text-left px-4 py-2 font-mono text-xs text-muted-foreground">Command</th>
                  <th className="text-left px-4 py-2 font-mono text-xs text-muted-foreground">Package</th>
                  <th className="text-left px-4 py-2 font-mono text-xs text-muted-foreground">Platform</th>
                  <th className="text-left px-4 py-2 font-mono text-xs text-muted-foreground">Description</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-border">
                {customScripts.map(([cmd, pkg, platform, desc], i) => (
                  <tr key={i} className="hover:bg-muted/30 transition-colors">
                    <td className="px-4 py-2 font-mono text-xs text-primary font-semibold">{cmd}</td>
                    <td className="px-4 py-2 font-mono text-xs text-muted-foreground">{pkg}</td>
                    <td className="px-4 py-2 text-xs text-muted-foreground">{platform}</td>
                    <td className="px-4 py-2 text-xs text-foreground">{desc}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          <CodeBlock code={`# Install Gitmap scripts fixer
$ gitmap install scripts-fixer

# Install coding guidelines compliance checker (v24)
$ gitmap install coding-guidelines
$ gitmap in cg

# Install AutoHotkey macro extension suite
$ gitmap install macro-ahk`} />
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
                  ["install.ps1", "gitmap-v28/scripts/", "PowerShell one-liner installer for Windows"],
                  ["install.sh", "gitmap-v28/scripts/", "Bash one-liner installer for Linux/macOS"],
                  ["uninstall.ps1", "gitmap-v28/scripts/", "PowerShell uninstaller"],
                  ["Get-LastRelease.ps1", "gitmap-v28/scripts/", "Resolve latest release version (3-tier fallback)"],
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
    Binary: C:\\Users\\me\\AppData\\Local\\gitmap-v28\\gitmap.exe
    Install Dir: C:\\Users\\me\\AppData\\Local\\gitmap-v28
    PATH Target: User PATH (added)
    Session PATH: refreshed for current PowerShell session`} />
            </div>
            <div className="border border-border rounded-lg p-4 bg-card">
              <h3 className="font-mono font-semibold text-sm mb-2">Unix / macOS</h3>
              <CodeBlock code={`Install summary
    Version: v2.65.0
    Binary: /Users/me/.local/bin/gitmap-v28
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
