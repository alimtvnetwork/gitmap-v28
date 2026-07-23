import DocsLayout from "@/components/docs/DocsLayout";
import CodeBlock from "@/components/docs/CodeBlock";
import { History, Database, Terminal, Clock } from "lucide-react";

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
            line.startsWith("FROM")
              ? "text-muted-foreground font-semibold"
              : line.includes("transition")
              ? "text-muted-foreground mt-2"
              : line.startsWith("Version history")
              ? "text-primary"
              : line.startsWith("No version")
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

const VersionHistoryPage = () => {
  return (
    <DocsLayout>
      <div className="max-w-4xl">
        {/* Header */}
        <div className="flex items-center gap-3 mb-2">
          <History className="w-8 h-8 text-primary" />
          <h1 className="text-3xl font-heading font-bold docs-h1">version-history</h1>
          <span className="text-xs font-mono bg-primary/10 text-foreground border border-primary/20 px-2 py-0.5 rounded dark:bg-primary/15 dark:text-primary transition-colors duration-300 hover:border-primary/40 hover:shadow-sm hover:shadow-primary/10">vh</span>
        </div>
        <p className="text-muted-foreground mb-8 text-lg">
          Display all version transitions recorded for the current repository, tracked automatically
          by <code className="text-primary">clone-next</code>.
        </p>

        {/* Usage */}
        <section className="mb-10">
          <h2 className="text-xl font-heading font-semibold mb-3">Usage</h2>
          <CodeBlock code="gitmap version-history [--limit N] [--json]" />
        </section>

        {/* How it works */}
        <section className="mb-10">
          <h2 className="text-xl font-heading font-semibold mb-3 docs-h2">How It Works</h2>
          <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
            {[
              { icon: Terminal, title: "Detect", desc: "Reads remote origin URL to identify the repo and its base name" },
              { icon: Database, title: "Lookup", desc: "Finds the repo in the database by its flattened absolute path" },
              { icon: History, title: "Query", desc: "Fetches all RepoVersionHistory rows for that repo, newest first" },
              { icon: Clock, title: "Display", desc: "Prints a formatted table or JSON with from/to versions and timestamps" },
            ].map(({ icon: Icon, title, desc }) => (
              <div key={title} className="border border-border rounded-lg p-4 bg-card">
                <Icon className="w-5 h-5 text-primary mb-2" />
                <h3 className="font-mono font-semibold text-sm mb-1">{title}</h3>
                <p className="text-xs text-muted-foreground">{desc}</p>
              </div>
            ))}
          </div>
        </section>

        {/* Flags */}
        <section className="mb-10">
          <h2 className="text-xl font-heading font-semibold mb-3 docs-h2">Flags</h2>
          <div className="overflow-x-auto">
            <table className="w-full text-sm border border-border rounded-lg overflow-hidden">
              <thead>
                <tr className="bg-muted/50">
                  <th className="text-left px-4 py-2 font-mono text-xs text-muted-foreground">Flag</th>
                  <th className="text-left px-4 py-2 font-mono text-xs text-muted-foreground">Default</th>
                  <th className="text-left px-4 py-2 font-mono text-xs text-muted-foreground">Description</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-border">
                {[
                  ["--limit N", "0 (all)", "Show only the last N transitions"],
                  ["--json", "false", "Output as JSON instead of table"],
                ].map(([flag, def, desc]) => (
                  <tr key={flag} className="bg-card">
                    <td className="px-4 py-2 font-mono text-xs text-primary">{flag}</td>
                    <td className="px-4 py-2 text-xs text-muted-foreground">{def}</td>
                    <td className="px-4 py-2 text-xs">{desc}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </section>

        {/* Examples */}
        <section className="mb-10">
          <h2 className="text-xl font-heading font-semibold mb-3 docs-h2">Examples</h2>

          <h3 className="font-mono text-sm font-semibold mb-2 text-muted-foreground">Show all version transitions</h3>
          <TerminalPreview
            title="gitmap vh"
            lines={[
              "Version history for D:\\wp-work\\riseup-asia\\macro-ahk:",
              "",
              "FROM        TO          FOLDER                    TIMESTAMP",
              "v11         v12         macro-ahk                 2026-04-16T10:30:00Z",
              "v12         v15         macro-ahk                 2026-04-16T14:22:00Z",
              "v15         v16         macro-ahk                 2026-04-16T16:45:00Z",
              "",
              "3 transition(s) recorded.",
            ]}
          />

          <h3 className="font-mono text-sm font-semibold mb-2 text-muted-foreground">Limit to last 2 transitions</h3>
          <TerminalPreview
            title="gitmap vh --limit 2"
            lines={[
              "Version history for D:\\wp-work\\riseup-asia\\macro-ahk:",
              "",
              "FROM        TO          FOLDER                    TIMESTAMP",
              "v12         v15         macro-ahk                 2026-04-16T14:22:00Z",
              "v15         v16         macro-ahk                 2026-04-16T16:45:00Z",
              "",
              "2 transition(s) recorded.",
            ]}
          />

          <h3 className="font-mono text-sm font-semibold mb-2 text-muted-foreground">JSON output</h3>
          <TerminalPreview
            title="gitmap vh --limit 1 --json"
            lines={[
              "[",
              "  {",
              '    "id": 3,',
              '    "repoId": 42,',
              '    "fromVersionTag": "v15",',
              '    "fromVersionNum": 15,',
              '    "toVersionTag": "v16",',
              '    "toVersionNum": 16,',
              '    "flattenedPath": "macro-ahk",',
              '    "createdAt": "2026-04-16T16:45:00Z"',
              "  }",
              "]",
            ]}
          />

          <h3 className="font-mono text-sm font-semibold mb-2 text-muted-foreground">No history found</h3>
          <TerminalPreview
            title="gitmap vh"
            lines={[
              "No version history found for this repo.",
            ]}
          />
        </section>

        {/* Database schema */}
        <section className="mb-10">
          <h2 className="text-xl font-heading font-semibold mb-3 flex items-center gap-2">
            <Database className="w-5 h-5" />
            Database Schema
          </h2>
          <p className="text-muted-foreground text-sm mb-4">
            Version transitions are stored in the <code className="text-primary">RepoVersionHistory</code> table,
            linked to <code className="text-primary">Repos</code> via foreign key.
          </p>
          <div className="overflow-x-auto">
            <table className="w-full text-sm border border-border rounded-lg overflow-hidden">
              <thead>
                <tr className="bg-muted/50">
                  <th className="text-left px-4 py-2 font-mono text-xs text-muted-foreground">Column</th>
                  <th className="text-left px-4 py-2 font-mono text-xs text-muted-foreground">Type</th>
                  <th className="text-left px-4 py-2 font-mono text-xs text-muted-foreground">Description</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-border">
                {[
                  ["Id", "INTEGER PK", "Auto-incrementing primary key"],
                  ["RepoId", "INTEGER FK", "References Repos(Id), cascade delete"],
                  ["FromVersionTag", "TEXT", "Previous version tag (e.g., v15)"],
                  ["FromVersionNum", "INTEGER", "Previous version number (e.g., 15)"],
                  ["ToVersionTag", "TEXT", "New version tag (e.g., v16)"],
                  ["ToVersionNum", "INTEGER", "New version number (e.g., 16)"],
                  ["FlattenedPath", "TEXT", "Base name of the flattened folder"],
                  ["CreatedAt", "TEXT", "ISO 8601 timestamp of the transition"],
                ].map(([col, type, desc]) => (
                  <tr key={col} className="bg-card">
                    <td className="px-4 py-2 font-mono text-xs text-primary">{col}</td>
                    <td className="px-4 py-2 font-mono text-xs text-muted-foreground">{type}</td>
                    <td className="px-4 py-2 text-xs">{desc}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </section>

        {/* File layout */}
        <section className="mb-10">
          <h2 className="text-xl font-heading font-semibold mb-3 docs-h2">File Layout</h2>
          <div className="overflow-x-auto">
            <table className="w-full text-sm border border-border rounded-lg overflow-hidden">
              <thead>
                <tr className="bg-muted/50">
                  <th className="text-left px-4 py-2 font-mono text-xs text-muted-foreground">File</th>
                  <th className="text-left px-4 py-2 font-mono text-xs text-muted-foreground">Responsibility</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-border">
                {[
                  ["cmd/versionhistory.go", "CLI handler, flag parsing, terminal/JSON output"],
                  ["cmd/clonenexthistory.go", "Records version transitions during clone-next"],
                  ["store/version_history.go", "DB insert, update, query for version history"],
                  ["model/version_history.go", "RepoVersionHistoryRecord struct"],
                  ["constants/constants_version_history.go", "SQL statements, table name, error messages"],
                  ["constants/constants_version_history_cmd.go", "CLI command names, help text, output formats"],
                  ["helptext/version-history.md", "Embedded help text for gitmap help version-history"],
                ].map(([file, desc]) => (
                  <tr key={file} className="bg-card">
                    <td className="px-4 py-2 font-mono text-xs text-primary">{file}</td>
                    <td className="px-4 py-2 text-xs">{desc}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </section>

        {/* See also */}
        <section>
          <h2 className="text-xl font-heading font-semibold mb-3 docs-h2">See Also</h2>
          <div className="flex flex-wrap gap-2">
            {[
              { label: "clone-next", href: "/clone-next" },
              { label: "history", href: "/history" },
              { label: "clone", href: "/clone" },
            ].map(({ label, href }) => (
              <a
                key={label}
                href={href}
                className="px-3 py-1.5 rounded-md border border-border bg-card text-xs font-mono hover:bg-primary/10 hover:text-foreground hover:border-primary/30 dark:hover:bg-primary/25 dark:hover:text-background transition-colors"
              >
                {label}
              </a>
            ))}
          </div>
        </section>
      </div>
    </DocsLayout>
  );
};

export default VersionHistoryPage;
