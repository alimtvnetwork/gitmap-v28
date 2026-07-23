import DocsLayout from "@/components/docs/DocsLayout";
import CodeBlock from "@/components/docs/CodeBlock";
import { Upload, Shield, Database, GitMerge, AlertTriangle } from "lucide-react";

const TerminalPreview = () => (
  <div className="rounded-lg border border-border overflow-hidden my-6">
    <div className="bg-terminal px-4 py-2 flex items-center gap-2 border-b border-border">
      <div className="flex gap-1.5">
        <span className="w-3 h-3 rounded-full bg-red-500/80" />
        <span className="w-3 h-3 rounded-full bg-yellow-500/80" />
        <span className="w-3 h-3 rounded-full bg-green-500/80" />
      </div>
      <span className="text-xs font-mono text-muted-foreground ml-2">gitmap import --confirm</span>
    </div>
    <div className="bg-terminal p-4 font-mono text-sm leading-relaxed">
      <div className="text-muted-foreground">Importing from gitmap-export.json...</div>
      <div className="text-yellow-400 mt-1">[1/42] my-api... added</div>
      <div className="text-yellow-400">[2/42] ui-lib... added</div>
      <div className="text-muted-foreground">...</div>
      <div className="text-green-400 mt-1">✓ 42 repos imported (3 skipped)</div>
      <div className="text-muted-foreground mt-1 text-xs">
        (42 repos, 3 groups, 8 releases, 156 history, 5 bookmarks)
      </div>
    </div>
  </div>
);

const features = [
  { icon: GitMerge, title: "Merge Semantics", desc: "Upserts and insert-or-ignore — never deletes existing data." },
  { icon: Shield, title: "--confirm Required", desc: "Safety flag prevents accidental data overwrites." },
  { icon: Database, title: "Full Restore", desc: "Repos, groups, releases, history, and bookmarks in one pass." },
  { icon: AlertTriangle, title: "Slug Resolution", desc: "Group members linked by slug; missing repos silently skipped." },
];

const strategies = [
  ["Repos", "Upsert by ID", "Updates if exists"],
  ["Groups", "INSERT OR IGNORE", "Skips if name exists"],
  ["GroupRepos", "Slug resolution", "Links repos to groups by slug"],
  ["Releases", "Upsert by Tag", "Updates if exists"],
  ["CommandHistory", "INSERT OR IGNORE by ID", "Skips duplicates"],
  ["Bookmarks", "INSERT OR IGNORE by ID", "Skips duplicates"],
];

const ImportPage = () => (
  <DocsLayout>
    <div className="max-w-4xl">
      <div className="flex items-center gap-3 mb-2">
        <Upload className="h-8 w-8 text-primary" />
        <div>
          <h1 className="text-3xl font-heading font-bold text-foreground docs-h1">Import</h1>
          <p className="text-muted-foreground font-mono text-sm">gitmap import (im)</p>
        </div>
      </div>
      <p className="text-muted-foreground mb-8 text-lg">
        Restore the database from a portable JSON export file using merge/upsert semantics.
      </p>

      <TerminalPreview />

      {/* Features */}
      <section className="mb-10">
        <h2 className="text-xl font-heading font-bold text-foreground mb-4 docs-h2">Features</h2>
        <div className="grid sm:grid-cols-2 gap-4">
          {features.map((f) => (
            <div key={f.title} className="border border-border rounded-lg p-4 bg-card">
              <div className="flex items-center gap-2 mb-2">
                <f.icon className="h-5 w-5 text-primary" />
                <span className="font-mono font-semibold text-foreground">{f.title}</span>
              </div>
              <p className="text-sm text-muted-foreground">{f.desc}</p>
            </div>
          ))}
        </div>
      </section>

      {/* Usage */}
      <section className="mb-10">
        <h2 className="text-xl font-heading font-bold text-foreground mb-4 docs-h2">Usage</h2>
        <CodeBlock code="gitmap import [file] --confirm" />
        <div className="mt-4 space-y-3">
          <div>
            <h3 className="font-mono font-semibold text-foreground mb-1">Arguments</h3>
            <div className="overflow-x-auto">
              <table className="w-full text-sm font-mono">
                <thead>
                  <tr className="border-b border-border text-muted-foreground">
                    <th className="text-left py-2 pr-4">Argument</th>
                    <th className="text-left py-2 pr-4">Default</th>
                    <th className="text-left py-2">Description</th>
                  </tr>
                </thead>
                <tbody>
                  <tr className="border-b border-border">
                    <td className="py-2 pr-4 text-primary">file</td>
                    <td className="py-2 pr-4 text-muted-foreground">gitmap-export.json</td>
                    <td className="py-2 text-foreground">Input file path</td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
          <div>
            <h3 className="font-mono font-semibold text-foreground mb-1">Flags</h3>
            <div className="overflow-x-auto">
              <table className="w-full text-sm font-mono">
                <thead>
                  <tr className="border-b border-border text-muted-foreground">
                    <th className="text-left py-2 pr-4">Flag</th>
                    <th className="text-left py-2 pr-4">Required</th>
                    <th className="text-left py-2">Description</th>
                  </tr>
                </thead>
                <tbody>
                  <tr className="border-b border-border">
                    <td className="py-2 pr-4 text-primary">--confirm</td>
                    <td className="py-2 pr-4 text-yellow-500">Yes</td>
                    <td className="py-2 text-foreground">Confirm the import (prevents accidents)</td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </section>

      {/* Examples */}
      <section className="mb-10">
        <h2 className="text-xl font-heading font-bold text-foreground mb-4 docs-h2">Examples</h2>
        <div className="space-y-4">
          <div>
            <p className="text-sm text-muted-foreground mb-1">Import from default file:</p>
            <CodeBlock code={`gitmap import --confirm\ngitmap im --confirm`} />
          </div>
          <div>
            <p className="text-sm text-muted-foreground mb-1">Import from custom path:</p>
            <CodeBlock code="gitmap import backup-2026-03.json --confirm" />
          </div>
          <div>
            <p className="text-sm text-muted-foreground mb-1">Without --confirm (error):</p>
            <CodeBlock code={`$ gitmap import\nimport requires --confirm flag (existing data will be merged)`} />
          </div>
        </div>
      </section>

      {/* Import Strategies */}
      <section className="mb-10">
        <h2 className="text-xl font-heading font-bold text-foreground mb-4 docs-h2">Import Behavior</h2>
        <p className="text-sm text-muted-foreground mb-3">
          Each table uses a specific merge strategy to safely restore data without data loss:
        </p>
        <div className="overflow-x-auto">
          <table className="w-full text-sm font-mono">
            <thead>
              <tr className="border-b border-border text-muted-foreground">
                <th className="text-left py-2 pr-4">Table</th>
                <th className="text-left py-2 pr-4">Strategy</th>
                <th className="text-left py-2">Notes</th>
              </tr>
            </thead>
            <tbody>
              {strategies.map(([table, strategy, notes]) => (
                <tr key={table} className="border-b border-border">
                  <td className="py-2 pr-4 text-primary">{table}</td>
                  <td className="py-2 pr-4 text-accent-foreground">{strategy}</td>
                  <td className="py-2 text-foreground">{notes}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </section>

      {/* Group Resolution */}
      <section className="mb-10">
        <h2 className="text-xl font-heading font-bold text-foreground mb-4 docs-h2">Group Slug Resolution</h2>
        <div className="border border-border rounded-lg p-4 bg-card">
          <div className="flex items-center gap-3 mb-3">
            <div className="flex items-center gap-2">
              <span className="px-2 py-1 bg-primary/10 text-foreground border border-primary/20 dark:bg-primary/15 dark:text-primary rounded text-xs font-mono transition-colors duration-300 hover:border-primary/40 hover:shadow-sm hover:shadow-primary/10">1. Import group</span>
              <span className="text-muted-foreground">→</span>
              <span className="px-2 py-1 bg-primary/10 text-foreground border border-primary/20 dark:bg-primary/15 dark:text-primary rounded text-xs font-mono transition-colors duration-300 hover:border-primary/40 hover:shadow-sm hover:shadow-primary/10">2. Resolve slugs</span>
              <span className="text-muted-foreground">→</span>
              <span className="px-2 py-1 bg-primary/10 text-foreground border border-primary/20 dark:bg-primary/15 dark:text-primary rounded text-xs font-mono transition-colors duration-300 hover:border-primary/40 hover:shadow-sm hover:shadow-primary/10">3. Link repos</span>
            </div>
          </div>
          <p className="text-sm text-muted-foreground">
            Groups are created first via INSERT OR IGNORE. Then each <code className="text-primary">repoSlug</code> is
            resolved against the Repos table. If a slug doesn't exist (repo not imported), the link is silently skipped.
          </p>
        </div>
      </section>

      {/* File Layout */}
      <section className="mb-10">
        <h2 className="text-xl font-heading font-bold text-foreground mb-4 docs-h2">File Layout</h2>
        <div className="overflow-x-auto">
          <table className="w-full text-sm font-mono">
            <thead>
              <tr className="border-b border-border text-muted-foreground">
                <th className="text-left py-2 pr-4">File</th>
                <th className="text-left py-2">Purpose</th>
              </tr>
            </thead>
            <tbody>
              {[
                ["constants/constants_import.go", "Command names, messages"],
                ["store/import.go", "ImportAll with per-table restore methods"],
                ["cmd/importcmd.go", "Import command handler"],
                ["helptext/import.md", "CLI help text"],
              ].map(([file, purpose]) => (
                <tr key={file} className="border-b border-border">
                  <td className="py-2 pr-4 text-primary">{file}</td>
                  <td className="py-2 text-foreground">{purpose}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </section>

      {/* See Also */}
      <section className="mb-10">
        <h2 className="text-xl font-heading font-bold text-foreground mb-4 docs-h2">See Also</h2>
        <ul className="space-y-1 text-sm font-mono">
          <li><a href="/export" className="text-primary hover:underline">export</a> — Export the database to a file</li>
          <li><a href="/commands" className="text-primary hover:underline">scan</a> — Scan directories as an alternative to import</li>
          <li><a href="/commands" className="text-primary hover:underline">profile</a> — Manage database profiles</li>
        </ul>
      </section>
    </div>
  </DocsLayout>
);

export default ImportPage;
