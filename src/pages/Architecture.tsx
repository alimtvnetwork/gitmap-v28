import DocsLayout from "@/components/docs/DocsLayout";
import CodeBlock from "@/components/docs/CodeBlock";

const ArchitecturePage = () => {
  return (
    <DocsLayout>
      <h1 className="text-3xl font-heading font-bold mb-2 docs-h1">Architecture</h1>
      <p className="text-muted-foreground mb-8">
        High-level overview of gitmap-v28's design and data flow.
      </p>

      <section className="space-y-8">
        <div>
          <h2 className="text-xl font-heading font-semibold mb-3 docs-h2">Project Structure</h2>
          <CodeBlock
            title="Directory Layout"
            code={`gitmap-v28/
├── cmd/              # CLI command handlers (one file per command)
├── cloner/           # Clone and safe-pull logic
├── completion/       # Shell completion (bash, zsh, powershell)
├── config/           # Config file loader + flag merging
├── constants/        # All string constants (15+ files)
├── data/             # Default config and template files
├── desktop/          # GitHub Desktop integration
├── detector/         # Project type detection (Go, Node, React, C++, C#)
├── formatter/        # CSV, JSON, terminal, tree output formatters
├── gitutil/          # Git command wrappers
├── helptext/         # Embedded Markdown help files (41 files)
├── mapper/           # Directory tree scanner
├── model/            # Data structures (ScanRecord, Group, Release, Config)
├── release/          # Release workflow, semver, assets, cross-compile
│   ├── assets.go     # Go cross-compilation orchestration
│   ├── assetstargets.go  # Target matrix + config resolution
│   ├── assetsupload.go   # GitHub API upload with retry
│   ├── compress.go   # .zip/.tar.gz archive creation
│   └── checksums.go  # SHA256 checksum generation
├── scanner/          # Recursive repo discovery
├── setup/            # Git config applicator
├── store/            # SQLite database layer
├── verbose/          # Debug logging
└── main.go           # Entry point`}
          />
        </div>

        <hr className="docs-hr" />

        <div>
          <h2 className="text-xl font-heading font-semibold mb-3 docs-h2">Data Flow</h2>
          <div className="bg-card border border-border rounded-lg p-6">
            <div className="space-y-4 font-mono text-sm">
              <div className="flex items-center gap-3">
                <span className="bg-primary/10 text-foreground border border-primary/20 dark:bg-primary/15 dark:text-primary px-3 py-1 rounded transition-colors duration-300 hover:border-primary/40 hover:shadow-sm hover:shadow-primary/10">scanner</span>
                <span className="text-muted-foreground">→ discovers .git directories recursively</span>
              </div>
              <div className="flex items-center gap-3">
                <span className="bg-primary/10 text-foreground border border-primary/20 dark:bg-primary/15 dark:text-primary px-3 py-1 rounded transition-colors duration-300 hover:border-primary/40 hover:shadow-sm hover:shadow-primary/10">mapper</span>
                <span className="text-muted-foreground">→ extracts metadata (URLs, branches, paths)</span>
              </div>
              <div className="flex items-center gap-3">
                <span className="bg-primary/10 text-foreground border border-primary/20 dark:bg-primary/15 dark:text-primary px-3 py-1 rounded transition-colors duration-300 hover:border-primary/40 hover:shadow-sm hover:shadow-primary/10">formatter</span>
                <span className="text-muted-foreground">→ outputs CSV, JSON, terminal, tree, scripts</span>
              </div>
              <div className="flex items-center gap-3">
                <span className="bg-primary/10 text-foreground border border-primary/20 dark:bg-primary/15 dark:text-primary px-3 py-1 rounded transition-colors duration-300 hover:border-primary/40 hover:shadow-sm hover:shadow-primary/10">store</span>
                <span className="text-muted-foreground">→ upserts records into SQLite database</span>
              </div>
              <div className="flex items-center gap-3">
                <span className="bg-primary/10 text-foreground border border-primary/20 dark:bg-primary/15 dark:text-primary px-3 py-1 rounded transition-colors duration-300 hover:border-primary/40 hover:shadow-sm hover:shadow-primary/10">cloner</span>
                <span className="text-muted-foreground">→ re-clones from structured files with progress</span>
              </div>
            </div>
          </div>
        </div>

        <hr className="docs-hr" />

        <div>
          <h2 className="text-xl font-heading font-semibold mb-3 docs-h2">Database Schema</h2>
          <p className="text-muted-foreground mb-3">
            SQLite database at <code className="docs-inline-code">.gitmap/output/data/gitmap.db</code> with
            PascalCase naming convention:
          </p>
          <CodeBlock
            title="Tables"
            language="sql"
            code={`Repos          — Id, Slug, RepoName, HttpsUrl, SshUrl, Branch,
                 RelativePath, AbsolutePath, CloneInstruction, Notes

Groups         — Id, Name, Description, Color, CreatedAt

GroupRepos     — GroupId, RepoId (join table)

Releases       — Id, Version, Tag, Branch, SourceBranch,
                 CommitSha, Changelog, Draft, IsLatest, Source`}
          />
        </div>

        <hr className="docs-hr" />

        <div>
          <h2 className="text-xl font-heading font-semibold mb-3 docs-h2">Output Artifacts</h2>
          <p className="text-muted-foreground mb-3">
            A scan generates the following in <code className="docs-inline-code">.gitmap/output/</code>:
          </p>
          <div className="bg-card border border-border rounded-lg overflow-hidden">
            <table className="w-full text-sm docs-table">
              <thead>
                <tr className="border-b border-border">
                  <th className="text-left px-4 py-2">File</th>
                  <th className="text-left px-4 py-2">Purpose</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-border">
                {[
                  ["gitmap.csv", "Flat CSV with all repo metadata"],
                  ["gitmap.json", "Structured JSON array of ScanRecords"],
                  ["gitmap.txt", "Plain text clone commands"],
                  ["folder-structure.md", "Markdown directory tree"],
                  ["clone.ps1", "PowerShell clone script with error handling"],
                  ["direct-clone.ps1", "HTTPS clone commands"],
                  ["direct-clone-ssh.ps1", "SSH clone commands"],
                  ["desktop.ps1", "GitHub Desktop registration script"],
                  ["last-scan.json", "Cached scan flags for rescan"],
                ].map(([file, desc]) => (
                  <tr key={file}>
                    <td className="px-4 py-2 font-mono text-primary">{file}</td>
                    <td className="px-4 py-2 text-muted-foreground">{desc}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>

        <hr className="docs-hr" />

        <div>
          <h2 className="text-xl font-heading font-semibold mb-3 docs-h2">Build System</h2>
          <p className="text-muted-foreground mb-3">
            Cross-platform build scripts with feature parity:
          </p>
          <div className="bg-card border border-border rounded-lg overflow-hidden">
            <table className="w-full text-sm docs-table">
              <thead>
                <tr className="border-b border-border">
                  <th className="text-left px-4 py-2">Script</th>
                  <th className="text-left px-4 py-2">Platform</th>
                  <th className="text-left px-4 py-2">Purpose</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-border">
                {[
                  ["run.ps1", "Windows", "PowerShell build pipeline"],
                  ["run.sh", "Linux / macOS", "Bash build pipeline (reads powershell.json)"],
                  ["Makefile", "Any (bash)", "Standard interface: build, run, test, update"],
                ].map(([file, platform, desc]) => (
                  <tr key={file}>
                    <td className="px-4 py-2 font-mono text-primary">{file}</td>
                    <td className="px-4 py-2 text-muted-foreground">{platform}</td>
                    <td className="px-4 py-2 text-muted-foreground">{desc}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
          <p className="text-muted-foreground mt-3 text-sm">
            CI/CD via GitHub Actions: test on push, cross-compile 6 targets, auto-release on tags.
          </p>
        </div>

        <hr className="docs-hr" />

        <div>
          <h2 className="text-xl font-heading font-semibold mb-3 docs-h2">Code Style</h2>
          <ul className="space-y-2 text-sm text-muted-foreground">
            <li className="flex gap-2"><span className="text-primary">•</span> All files under 200 lines</li>
            <li className="flex gap-2"><span className="text-primary">•</span> Functions 8–15 lines (focused, single-purpose)</li>
            <li className="flex gap-2"><span className="text-primary">•</span> No magic strings — all constants in dedicated files</li>
            <li className="flex gap-2"><span className="text-primary">•</span> PascalCase for all DB table and column names</li>
            <li className="flex gap-2"><span className="text-primary">•</span> Positive-logic conditionals (no negation in if-statements)</li>
            <li className="flex gap-2"><span className="text-primary">•</span> SemVer versioning — bump on every behavior change</li>
          </ul>
        </div>
      </section>
    </DocsLayout>
  );
};

export default ArchitecturePage;
