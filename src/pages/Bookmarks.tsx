import DocsLayout from "@/components/docs/DocsLayout";
import CodeBlock from "@/components/docs/CodeBlock";
import { Bookmark, Play, Trash2, Save, Database } from "lucide-react";

const MOCK_BOOKMARKS = [
  { name: "ssh-scan", command: "scan", args: "", flags: "--mode ssh", created: "2026-03-15 09:12" },
  { name: "quick-status", command: "status", args: "", flags: "", created: "2026-03-14 14:30" },
  { name: "scan-projects", command: "scan", args: "./projects", flags: "--mode ssh --open", created: "2026-03-12 11:45" },
  { name: "daily-pull", command: "pull", args: "", flags: "--all", created: "2026-03-10 08:00" },
];

const TerminalPreview = () => (
  <div className="rounded-lg border border-border overflow-hidden my-6">
    <div className="bg-terminal px-4 py-2 flex items-center gap-2 border-b border-border">
      <div className="flex gap-1.5">
        <span className="w-3 h-3 rounded-full bg-red-500/80" />
        <span className="w-3 h-3 rounded-full bg-yellow-500/80" />
        <span className="w-3 h-3 rounded-full bg-green-500/80" />
      </div>
      <span className="text-xs font-mono text-muted-foreground ml-2">gitmap bookmark list</span>
    </div>
    <div className="bg-terminal p-4 font-mono text-sm leading-relaxed overflow-x-auto">
      <div className="text-primary font-bold text-xs">
        {"  "}NAME{"              "}COMMAND{"    "}ARGS{"          "}FLAGS{"                "}CREATED
      </div>
      {MOCK_BOOKMARKS.map((b) => (
        <div key={b.name} className="text-terminal-foreground text-xs">
          {"  "}
          <span className="inline-block w-[120px] text-foreground">{b.name}</span>
          <span className="inline-block w-[72px] text-primary">{b.command}</span>
          <span className="inline-block w-[100px] text-muted-foreground">{b.args || "—"}</span>
          <span className="inline-block w-[140px] text-muted-foreground">{b.flags || "—"}</span>
          <span className="text-muted-foreground">{b.created}</span>
        </div>
      ))}
      <div className="mt-3 text-muted-foreground text-xs border-t border-border/50 pt-2">
        {MOCK_BOOKMARKS.length} bookmarks saved
      </div>
    </div>
  </div>
);

const features = [
  { icon: Save, title: "Save Commands", desc: "Store frequently-used command+flag combinations under a memorable name for instant replay." },
  { icon: Play, title: "Replay via Dispatch", desc: "Replayed bookmarks go through the standard dispatch function, ensuring full audit trail coverage." },
  { icon: Database, title: "SQLite Persistence", desc: "Bookmarks are stored in a dedicated Bookmarks table with unique name constraints." },
  { icon: Trash2, title: "Safe Deletion", desc: "Delete bookmarks by name. Names must be unique — save refuses if a name already exists." },
];

const schema = [
  ["Id", "TEXT", "PRIMARY KEY", "Timestamp-based unique ID"],
  ["Name", "TEXT", "NOT NULL UNIQUE", "User-chosen bookmark name"],
  ["Command", "TEXT", "NOT NULL", "Command name (e.g. scan)"],
  ["Args", "TEXT", "DEFAULT ''", "Positional arguments"],
  ["Flags", "TEXT", "DEFAULT ''", "Flags (e.g. --mode ssh)"],
  ["CreatedAt", "TEXT", "DEFAULT CURRENT_TIMESTAMP", ""],
];

const BookmarksPage = () => (
  <DocsLayout>
    <h1 className="text-3xl font-heading font-bold mb-2 docs-h1">Bookmarks</h1>
    <p className="text-muted-foreground mb-6">
      Save and replay frequently-used CLI command+flag combinations by name.
    </p>

    <h2 className="text-xl font-heading font-semibold mt-8 mb-2">Live Preview</h2>
    <p className="text-sm text-muted-foreground mb-2">
      Simulated terminal output of the bookmark list command.
    </p>
    <TerminalPreview />

    <h2 className="text-xl font-heading font-semibold mt-10 mb-4">Features</h2>
    <div className="grid md:grid-cols-2 gap-4 mb-8">
      {features.map((f) => (
        <div key={f.title} className="rounded-lg border border-border bg-card p-4">
          <f.icon className="h-5 w-5 text-primary mb-2" />
          <h3 className="font-mono font-semibold text-sm mb-1">{f.title}</h3>
          <p className="text-xs text-muted-foreground">{f.desc}</p>
        </div>
      ))}
    </div>

    <h2 className="text-xl font-heading font-semibold mt-10 mb-3">Subcommands</h2>
    <CodeBlock code='gitmap bookmark save ssh-scan scan --mode ssh' title="Save a command bookmark" />
    <CodeBlock code='gitmap bk save scan-projects scan ./projects --mode ssh --open' title="Save with args and multiple flags" />
    <CodeBlock code="gitmap bookmark list" title="List all saved bookmarks" />
    <CodeBlock code="gitmap bk list --json" title="List as JSON" />
    <CodeBlock code="gitmap bookmark run ssh-scan" title="Replay a saved bookmark" />
    <CodeBlock code="gitmap bookmark delete ssh-scan" title="Remove a bookmark" />

    <h2 className="text-xl font-heading font-semibold mt-10 mb-3">Replay Behavior</h2>
    <div className="grid md:grid-cols-3 gap-4 mb-8">
      <div className="rounded-lg border border-border bg-card p-4">
        <h3 className="font-mono font-semibold text-sm mb-1 text-primary">1. Load</h3>
        <p className="text-xs text-muted-foreground">Bookmark record is fetched from the database by name.</p>
      </div>
      <div className="rounded-lg border border-border bg-card p-4">
        <h3 className="font-mono font-semibold text-sm mb-1 text-primary">2. Reconstruct</h3>
        <p className="text-xs text-muted-foreground"><code className="text-primary">os.Args</code> is rebuilt from saved command, args, and flags.</p>
      </div>
      <div className="rounded-lg border border-border bg-card p-4">
        <h3 className="font-mono font-semibold text-sm mb-1 text-primary">3. Dispatch</h3>
        <p className="text-xs text-muted-foreground">Standard <code className="text-primary">dispatch()</code> is called — audit hook records the execution.</p>
      </div>
    </div>

    <h2 className="text-xl font-heading font-semibold mt-10 mb-3">Table Schema</h2>
    <div className="rounded-lg border border-border overflow-hidden">
      <table className="w-full text-sm">
        <thead>
          <tr className="bg-muted/50">
            <th className="text-left font-mono font-semibold px-4 py-2">Column</th>
            <th className="text-left font-mono font-semibold px-4 py-2">Type</th>
            <th className="text-left font-mono font-semibold px-4 py-2">Constraints</th>
            <th className="text-left font-mono font-semibold px-4 py-2">Notes</th>
          </tr>
        </thead>
        <tbody>
          {schema.map(([col, type, constraints, notes]) => (
            <tr key={col} className="border-t border-border">
              <td className="px-4 py-2 font-mono text-primary">{col}</td>
              <td className="px-4 py-2 font-mono text-muted-foreground">{type}</td>
              <td className="px-4 py-2 text-muted-foreground">{constraints}</td>
              <td className="px-4 py-2 text-muted-foreground">{notes}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>

    <h2 className="text-xl font-heading font-semibold mt-10 mb-3">Constraints</h2>
    <ul className="list-disc list-inside text-sm text-muted-foreground space-y-1 mb-8">
      <li>Bookmark names must be unique (enforced by UNIQUE constraint)</li>
      <li>Save refuses if name exists — user must delete first</li>
      <li><code className="text-primary font-mono">db-reset --confirm</code> also clears the Bookmarks table</li>
      <li>PascalCase table and column names</li>
      <li>All files under 200 lines, all functions 8–15 lines</li>
    </ul>

    <h2 className="text-xl font-heading font-semibold mt-10 mb-3">File Layout</h2>
    <div className="rounded-lg border border-border overflow-hidden">
      <table className="w-full text-sm">
        <thead>
          <tr className="bg-muted/50">
            <th className="text-left font-mono font-semibold px-4 py-2">File</th>
            <th className="text-left font-mono font-semibold px-4 py-2">Purpose</th>
          </tr>
        </thead>
        <tbody>
          {[
            ["constants/constants_bookmark.go", "SQL, command names, messages"],
            ["model/bookmark.go", "BookmarkRecord struct"],
            ["store/bookmark.go", "Bookmark CRUD operations"],
            ["cmd/bookmark.go", "Bookmark command routing"],
            ["cmd/bookmarksave.go", "Save subcommand"],
            ["cmd/bookmarklist.go", "List and delete subcommands"],
            ["cmd/bookmarkrun.go", "Run (replay) subcommand"],
          ].map(([file, purpose]) => (
            <tr key={file} className="border-t border-border">
              <td className="px-4 py-2 font-mono text-primary">{file}</td>
              <td className="px-4 py-2 text-muted-foreground">{purpose}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  </DocsLayout>
);

export default BookmarksPage;
