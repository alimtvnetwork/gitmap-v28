import { useMemo, useState } from "react";
import DocsLayout from "@/components/docs/DocsLayout";
import CodeBlock from "@/components/docs/CodeBlock";
import { Chrome, Download, Upload, Database, Trash2, Search } from "lucide-react";
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

const commands = [
  { name: "chrome-profile-copy", alias: "cpc", desc: "Copy a Chrome profile (bookmarks, extensions, prefs, flags). Always emits a JSON + CSV pair and upserts ChromeProfile." },
  { name: "chrome-profile-export", alias: "cpe", desc: "Export a profile to JSON + CSV. Both paths print under a single Artifacts: block." },
  { name: "chrome-profile-import", alias: "cpi", desc: "Restore a profile from a .json snapshot (full fidelity) or a .csv snapshot (extension IDs + known preferences)." },
  { name: "chrome-profile-list", alias: "cpl", desc: "List profiles on disk + every profile tracked in SQLite (with export counts and last-seen)." },
  { name: "chrome-profile-delete", alias: "cpd", desc: "Remove a profile and its stored JSON/CSV artifacts from the gitmap DB. Requires --yes." },
];

// Sample rows mirror the ChromeProfileExport list returned by
// ListChromeProfilesDB() — the UI table is a faithful preview of the
// real CLI output until a backend bridge is added.
const sampleRows = [
  { name: "Default", lastSeen: "2026-06-19T08:14:02Z", exports: 4 },
  { name: "Profile 1", lastSeen: "2026-06-18T22:10:11Z", exports: 1 },
  { name: "Work", lastSeen: "2026-06-17T11:02:48Z", exports: 6 },
  { name: "Personal", lastSeen: "2026-05-30T07:55:00Z", exports: 2 },
];

type SortKey = "name" | "lastSeen" | "exports";

const ChromeProfileSpec = () => {
  const [query, setQuery] = useState("");
  const [sortKey, setSortKey] = useState<SortKey>("lastSeen");
  const [minExports, setMinExports] = useState("0");

  const filtered = useMemo(() => {
    const q = query.trim().toLowerCase();
    const min = Number.parseInt(minExports, 10) || 0;
    const rows = sampleRows.filter(
      (r) => r.exports >= min && (q === "" || r.name.toLowerCase().includes(q)),
    );
    return [...rows].sort((a, b) => {
      if (sortKey === "name") return a.name.localeCompare(b.name);
      if (sortKey === "exports") return b.exports - a.exports;
      return b.lastSeen.localeCompare(a.lastSeen);
    });
  }, [query, sortKey, minExports]);

  return (
    <DocsLayout>
      <div className="max-w-4xl space-y-10">
        <div>
          <div className="flex items-center gap-3 mb-2">
            <Chrome className="h-8 w-8 text-primary" />
            <h1 className="text-3xl font-bold tracking-tight">Chrome profile management</h1>
          </div>
          <p className="text-lg text-muted-foreground">
            Offline copy, export, import, audit, and delete of Google Chrome profiles. Every export
            writes both JSON and CSV, both paths are printed under a consistent <code>Artifacts:</code>
            block, and the SQLite DB tracks each snapshot so listing and cleanup stay accurate.
          </p>
        </div>

        <section>
          <h2 className="text-xl font-semibold mb-3">Commands</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            {commands.map((c) => (
              <div key={c.name} className="rounded-lg border border-border p-4 bg-card">
                <div className="flex items-center gap-2 mb-1">
                  <code className="font-mono text-sm text-primary">{c.name}</code>
                  <span className="font-mono text-xs px-2 py-0.5 rounded bg-primary/10 text-foreground border border-primary/20">
                    {c.alias}
                  </span>
                </div>
                <p className="text-xs text-muted-foreground">{c.desc}</p>
              </div>
            ))}
          </div>
        </section>

        <section>
          <h2 className="text-xl font-semibold mb-3 flex items-center gap-2">
            <Download className="h-5 w-5 text-primary" /> Export — consistent Artifacts block
          </h2>
          <CodeBlock code={`$ gitmap cpe Default
Artifacts:
  json:  .gitmap/chrome/Default.json
  csv:   .gitmap/chrome/Default.csv
chrome-profile: db synced (Default)`} />
        </section>

        <section>
          <h2 className="text-xl font-semibold mb-3 flex items-center gap-2">
            <Upload className="h-5 w-5 text-primary" /> Import — JSON or CSV
          </h2>
          <CodeBlock code={`gitmap cpi ./snapshots/work.json              # full fidelity
gitmap cpi ./snapshots/work.csv "Profile 2"   # CSV restore (lossy: no bookmarks)`} />
        </section>

        <section>
          <h2 className="text-xl font-semibold mb-3 flex items-center gap-2">
            <Trash2 className="h-5 w-5 text-primary" /> Delete
          </h2>
          <CodeBlock code={`gitmap cpd "Profile 2" --yes
# rm .gitmap/chrome/Profile 2.json
# rm .gitmap/chrome/Profile 2.csv
# chrome-profile-delete: removed profile "Profile 2" (2 artifacts)`} />
        </section>

        <section>
          <h2 className="text-xl font-semibold mb-3 flex items-center gap-2">
            <Database className="h-5 w-5 text-primary" /> Tracked profiles
          </h2>
          <div className="flex flex-col md:flex-row gap-3 mb-4">
            <div className="relative flex-1">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
              <Input
                placeholder="Search by name…"
                value={query}
                onChange={(e) => setQuery(e.target.value)}
                className="pl-9"
              />
            </div>
            <Select value={sortKey} onValueChange={(v) => setSortKey(v as SortKey)}>
              <SelectTrigger className="w-full md:w-48">
                <SelectValue placeholder="Sort by" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="lastSeen">Last seen (newest)</SelectItem>
                <SelectItem value="name">Name (A→Z)</SelectItem>
                <SelectItem value="exports">Artifact count</SelectItem>
              </SelectContent>
            </Select>
            <Select value={minExports} onValueChange={setMinExports}>
              <SelectTrigger className="w-full md:w-44">
                <SelectValue placeholder="Min artifacts" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="0">Any artifact count</SelectItem>
                <SelectItem value="1">≥ 1 artifact</SelectItem>
                <SelectItem value="3">≥ 3 artifacts</SelectItem>
                <SelectItem value="5">≥ 5 artifacts</SelectItem>
              </SelectContent>
            </Select>
          </div>
          <div className="overflow-x-auto rounded-lg border border-border">
            <table className="w-full text-sm">
              <thead className="bg-muted/50">
                <tr>
                  <th className="text-left px-4 py-2 font-medium">Name</th>
                  <th className="text-left px-4 py-2 font-medium">Last seen (UTC)</th>
                  <th className="text-left px-4 py-2 font-medium">Artifacts</th>
                </tr>
              </thead>
              <tbody>
                {filtered.map((r) => (
                  <tr key={r.name} className="border-t border-border">
                    <td className="px-4 py-2 font-mono text-primary">{r.name}</td>
                    <td className="px-4 py-2 text-muted-foreground">{r.lastSeen}</td>
                    <td className="px-4 py-2 text-muted-foreground">{r.exports}</td>
                  </tr>
                ))}
                {filtered.length === 0 && (
                  <tr>
                    <td colSpan={3} className="px-4 py-6 text-center text-muted-foreground">
                      No profiles match the current filters.
                    </td>
                  </tr>
                )}
              </tbody>
            </table>
          </div>
        </section>
      </div>
    </DocsLayout>
  );
};

export default ChromeProfileSpec;
