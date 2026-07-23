import DocsLayout from "@/components/docs/DocsLayout";
import CodeBlock from "@/components/docs/CodeBlock";
import { Download, FileJson, Database, Package, Shield } from "lucide-react";

const MOCK_EXPORT = {
  version: "2.21.0",
  exportedAt: "2026-03-18T10:30:00Z",
  repos: [
    { slug: "my-api", branch: "main", httpsUrl: "https://github.com/user/my-api.git" },
    { slug: "ui-lib", branch: "develop", httpsUrl: "https://github.com/user/ui-lib.git" },
  ],
  groups: [{ name: "frontend", color: "blue", repoSlugs: ["ui-lib"] }],
  releases: [{ tag: "v2.21.0", date: "2026-03-18" }],
  history: ["(3 records)"],
  bookmarks: ["(2 records)"],
};

const TerminalPreview = () => (
  <div className="rounded-lg border border-border overflow-hidden my-6">
    <div className="bg-terminal px-4 py-2 flex items-center gap-2 border-b border-border">
      <div className="flex gap-1.5">
        <span className="w-3 h-3 rounded-full bg-red-500/80" />
        <span className="w-3 h-3 rounded-full bg-yellow-500/80" />
        <span className="w-3 h-3 rounded-full bg-green-500/80" />
      </div>
      <span className="text-xs font-mono text-muted-foreground ml-2">gitmap export</span>
    </div>
    <div className="bg-terminal p-4 font-mono text-sm leading-relaxed">
      <div className="text-muted-foreground">Exporting 42 repos...</div>
      <div className="text-green-400 mt-1">✓ Exported to gitmap-export.json</div>
      <div className="text-muted-foreground mt-1 text-xs">
        (42 repos, 3 groups, 8 releases, 156 history, 5 bookmarks)
      </div>
    </div>
  </div>
);

const features = [
  { icon: Database, title: "Full Database Dump", desc: "Exports repos, groups, releases, history, and bookmarks in one file." },
  { icon: FileJson, title: "Portable JSON", desc: "Human-readable JSON with 2-space indentation, ready for version control." },
  { icon: Package, title: "Self-Contained Groups", desc: "Groups include repoSlugs array — no foreign key resolution needed." },
  { icon: Shield, title: "Read-Only", desc: "Export never modifies the database; safe to run anytime." },
];

const schema = [
  ["version", "string", "App version at export time"],
  ["exportedAt", "string", "RFC 3339 timestamp"],
  ["repos", "ScanRecord[]", "All scanned repositories"],
  ["groups", "GroupExport[]", "Groups with embedded repoSlugs"],
  ["releases", "ReleaseRecord[]", "All release metadata"],
  ["history", "CommandHistoryRecord[]", "Full command history"],
  ["bookmarks", "BookmarkRecord[]", "All saved bookmarks"],
];

const ExportPage = () => (
  <DocsLayout>
    <div className="max-w-4xl">
      <div className="flex items-center gap-3 mb-2">
        <Download className="h-8 w-8 text-primary" />
        <div>
          <h1 className="text-3xl font-heading font-bold text-foreground docs-h1">Export</h1>
          <p className="text-muted-foreground font-mono text-sm">gitmap export (ex)</p>
        </div>
      </div>
      <p className="text-muted-foreground mb-8 text-lg">
        Dump the entire gitmap database into a single portable JSON file for backup, sharing, or migration.
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
        <CodeBlock code="gitmap export [file]" />
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
                    <td className="py-2 text-foreground">Output file path</td>
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
            <p className="text-sm text-muted-foreground mb-1">Export to default file:</p>
            <CodeBlock code={`gitmap export\ngitmap ex`} />
          </div>
          <div>
            <p className="text-sm text-muted-foreground mb-1">Export to custom path:</p>
            <CodeBlock code="gitmap export backup-2026-03.json" />
          </div>
        </div>
      </section>

      {/* Export Schema */}
      <section className="mb-10">
        <h2 className="text-xl font-heading font-bold text-foreground mb-4 docs-h2">Export Schema (DatabaseExport)</h2>
        <div className="overflow-x-auto">
          <table className="w-full text-sm font-mono">
            <thead>
              <tr className="border-b border-border text-muted-foreground">
                <th className="text-left py-2 pr-4">Field</th>
                <th className="text-left py-2 pr-4">Type</th>
                <th className="text-left py-2">Description</th>
              </tr>
            </thead>
            <tbody>
              {schema.map(([field, type, desc]) => (
                <tr key={field} className="border-b border-border">
                  <td className="py-2 pr-4 text-primary">{field}</td>
                  <td className="py-2 pr-4 text-accent-foreground">{type}</td>
                  <td className="py-2 text-foreground">{desc}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </section>

      {/* JSON Sample */}
      <section className="mb-10">
        <h2 className="text-xl font-heading font-bold text-foreground mb-4 docs-h2">Sample Output</h2>
        <CodeBlock code={JSON.stringify(MOCK_EXPORT, null, 2)} />
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
                ["constants/constants_export.go", "Command names, messages, defaults"],
                ["model/export.go", "DatabaseExport and GroupExport structs"],
                ["store/export.go", "ExportAll aggregation method"],
                ["cmd/export.go", "Export command handler"],
                ["helptext/export.md", "CLI help text"],
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
          <li><a href="/import" className="text-primary hover:underline">import</a> — Restore database from backup</li>
          <li><a href="/commands" className="text-primary hover:underline">scan</a> — Scan directories to populate the database</li>
          <li><a href="/commands" className="text-primary hover:underline">profile</a> — Manage database profiles</li>
        </ul>
      </section>
    </div>
  </DocsLayout>
);

export default ExportPage;
