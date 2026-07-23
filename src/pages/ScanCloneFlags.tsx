import { useState, useMemo } from "react";
import DocsLayout from "@/components/docs/DocsLayout";
import SearchBar from "@/components/docs/SearchBar";

interface FlagSpec {
  flag: string;
  default: string;
  values: string;
  description: string;
  command: "scan" | "clone";
}

const flags: FlagSpec[] = [
  // ── scan ─────────────────────────────────────────────
  {
    command: "scan",
    flag: "--config <path>",
    default: "./data/config.json",
    values: "filesystem path",
    description: "Config file path to load scan settings from.",
  },
  {
    command: "scan",
    flag: "--mode <style>",
    default: "https",
    values: "ssh | https",
    description: "Clone URL style emitted in output files.",
  },
  {
    command: "scan",
    flag: "--output <fmt>",
    default: "terminal",
    values: "csv | json | terminal",
    description: "Output format for scan results.",
  },
  {
    command: "scan",
    flag: "--output-path <dir>",
    default: "./.gitmap/output",
    values: "filesystem path",
    description: "Directory where output files are written.",
  },
  {
    command: "scan",
    flag: "--out-file <path>",
    default: "(derived)",
    values: "filesystem path",
    description: "Exact output file path; overrides --output-path.",
  },
  {
    command: "scan",
    flag: "--github-desktop",
    default: "false",
    values: "boolean (flag)",
    description: "Register discovered repos with GitHub Desktop.",
  },
  {
    command: "scan",
    flag: "--open",
    default: "false",
    values: "boolean (flag)",
    description: "Open output folder in file explorer after scan.",
  },
  {
    command: "scan",
    flag: "--quiet",
    default: "false",
    values: "boolean (flag)",
    description: "Suppress the clone-help section (CI / scripted use).",
  },
  {
    command: "scan",
    flag: "--no-vscode-sync",
    default: "false (sync ON)",
    values: "boolean (flag)",
    description: "Skip syncing scanned repos into VS Code Project Manager projects.json.",
  },
  {
    command: "scan",
    flag: "--no-auto-tags",
    default: "false (tags ON)",
    values: "boolean (flag)",
    description: "Skip auto-derived tags (git/node/go/python/rust/docker) when syncing.",
  },

  // ── clone ────────────────────────────────────────────
  {
    command: "clone",
    flag: "--target-dir <dir>",
    default: "current directory",
    values: "filesystem path",
    description: "Base directory where repos are cloned.",
  },
  {
    command: "clone",
    flag: "--safe-pull",
    default: "false (auto-enabled)",
    values: "boolean (flag)",
    description: "Pull existing repos with retry + unlock diagnostics.",
  },
  {
    command: "clone",
    flag: "--github-desktop",
    default: "false",
    values: "boolean (flag)",
    description: "Auto-register cloned repos with GitHub Desktop (no prompt).",
  },
  {
    command: "clone",
    flag: "--verbose",
    default: "false",
    values: "boolean (flag)",
    description: "Write a detailed debug log to a timestamped file.",
  },
];

const ScanCloneFlagsPage = () => {
  const [search, setSearch] = useState("");
  const [filter, setFilter] = useState<"all" | "scan" | "clone">("all");

  const filtered = useMemo(() => {
    let rows = flags;
    if (filter !== "all") rows = rows.filter((r) => r.command === filter);
    if (search) {
      const q = search.toLowerCase();
      rows = rows.filter(
        (r) =>
          r.flag.toLowerCase().includes(q) ||
          r.description.toLowerCase().includes(q) ||
          r.values.toLowerCase().includes(q) ||
          r.default.toLowerCase().includes(q)
      );
    }
    return rows;
  }, [search, filter]);

  const scanCount = flags.filter((f) => f.command === "scan").length;
  const cloneCount = flags.filter((f) => f.command === "clone").length;

  return (
    <DocsLayout>
      <h1 className="text-3xl font-heading font-bold mb-2 docs-h1">
        scan &amp; clone — Flag Reference
      </h1>
      <p className="text-muted-foreground mb-6">
        Every supported flag for <code className="font-mono text-primary">gitmap scan</code> and{" "}
        <code className="font-mono text-primary">gitmap clone</code>, with defaults and valid values.
        {" "}
        <span className="text-xs">
          ({scanCount} scan flags · {cloneCount} clone flags)
        </span>
      </p>

      <div className="flex flex-wrap items-center gap-2 mb-4">
        {(["all", "scan", "clone"] as const).map((key) => (
          <button
            key={key}
            onClick={() => setFilter(key)}
            className={`px-3 py-1.5 rounded-md text-sm font-mono border transition-colors ${
              filter === key
                ? "bg-primary text-primary-foreground border-primary"
                : "bg-background text-muted-foreground border-border hover:text-foreground hover:border-foreground/30"
            }`}
          >
            {key === "all" ? `all (${flags.length})` : key === "scan" ? `scan (${scanCount})` : `clone (${cloneCount})`}
          </button>
        ))}
      </div>

      <SearchBar value={search} onChange={setSearch} placeholder="Search flag, value, or description..." />

      <div className="mt-6 rounded-lg border border-border overflow-hidden">
        <div className="overflow-x-auto">
          <table className="w-full text-sm">
            <thead>
              <tr className="bg-muted/30 border-b border-border">
                <th className="text-left px-4 py-2.5 font-mono font-semibold text-foreground whitespace-nowrap">
                  Command
                </th>
                <th className="text-left px-4 py-2.5 font-mono font-semibold text-foreground">
                  Flag
                </th>
                <th className="text-left px-4 py-2.5 font-mono font-semibold text-foreground">
                  Valid values
                </th>
                <th className="text-left px-4 py-2.5 font-mono font-semibold text-foreground">
                  Default
                </th>
                <th className="text-left px-4 py-2.5 font-mono font-semibold text-foreground">
                  Description
                </th>
              </tr>
            </thead>
            <tbody>
              {filtered.map((row, i) => (
                <tr
                  key={`${row.command}-${row.flag}-${i}`}
                  className="border-b border-border last:border-0 hover:bg-muted/20 transition-colors"
                >
                  <td className="px-4 py-2 font-mono text-xs whitespace-nowrap">
                    <span
                      className={`inline-block px-2 py-0.5 rounded-md border ${
                        row.command === "scan"
                          ? "bg-primary/10 text-foreground border-primary/20 dark:bg-primary/15 dark:text-primary"
                          : "bg-accent/30 text-accent-foreground border-accent/40 dark:bg-accent/40"
                      }`}
                    >
                      {row.command}
                    </span>
                  </td>
                  <td className="px-4 py-2 font-mono text-primary whitespace-nowrap">
                    {row.flag}
                  </td>
                  <td className="px-4 py-2 font-mono text-xs text-muted-foreground whitespace-nowrap">
                    {row.values}
                  </td>
                  <td className="px-4 py-2 font-mono text-xs text-muted-foreground whitespace-nowrap">
                    {row.default}
                  </td>
                  <td className="px-4 py-2 text-foreground/90">{row.description}</td>
                </tr>
              ))}
              {filtered.length === 0 && (
                <tr>
                  <td colSpan={5} className="px-4 py-8 text-center text-muted-foreground font-mono text-sm">
                    No flags matching "{search}"
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      </div>

      <div className="mt-8 rounded-lg border border-border bg-muted/20 p-4 text-sm">
        <h2 className="font-heading font-semibold text-foreground mb-2">Notes</h2>
        <ul className="list-disc pl-5 space-y-1 text-muted-foreground">
          <li>
            Boolean flags toggle on when present; pass nothing to keep the default.
          </li>
          <li>
            <code className="font-mono text-primary">--safe-pull</code> is auto-enabled for re-clones
            of existing folders even when not passed explicitly.
          </li>
          <li>
            For per-command examples, see the dedicated{" "}
            <a href="/scan-command" className="text-primary underline">scan</a> and{" "}
            <a href="/clone-command" className="text-primary underline">clone</a> reference pages.
          </li>
        </ul>
      </div>
    </DocsLayout>
  );
};

export default ScanCloneFlagsPage;
