import { useState, useMemo } from "react";
import DocsLayout from "@/components/docs/DocsLayout";
import SearchBar from "@/components/docs/SearchBar";
import { commands } from "@/data/commands";

interface FlagRow {
  flag: string;
  description: string;
  command: string;
  alias?: string;
}

const FlagReferencePage = () => {
  const [search, setSearch] = useState("");
  const [sortCol, setSortCol] = useState<"flag" | "command">("flag");
  const [sortAsc, setSortAsc] = useState(true);

  const allFlags = useMemo<FlagRow[]>(() => {
    const rows: FlagRow[] = [];
    for (const cmd of commands) {
      if (!cmd.flags) continue;
      for (const f of cmd.flags) {
        rows.push({ flag: f.flag, description: f.description, command: cmd.name, alias: cmd.alias });
      }
    }
    return rows;
  }, []);

  const filtered = useMemo(() => {
    let rows = allFlags;
    if (search) {
      const q = search.toLowerCase();
      rows = rows.filter(
        (r) => r.flag.toLowerCase().includes(q) || r.description.toLowerCase().includes(q) || r.command.includes(q)
      );
    }
    rows.sort((a, b) => {
      const va = sortCol === "flag" ? a.flag : a.command;
      const vb = sortCol === "flag" ? b.flag : b.command;
      return sortAsc ? va.localeCompare(vb) : vb.localeCompare(va);
    });
    return rows;
  }, [allFlags, search, sortCol, sortAsc]);

  const handleSort = (col: "flag" | "command") => {
    if (sortCol === col) {
      setSortAsc(!sortAsc);
    } else {
      setSortCol(col);
      setSortAsc(true);
    }
  };

  const sortIndicator = (col: "flag" | "command") =>
    sortCol === col ? (sortAsc ? " ↑" : " ↓") : "";

  return (
    <DocsLayout>
      <h1 className="text-3xl font-heading font-bold mb-2 docs-h1">Flag Reference</h1>
      <p className="text-muted-foreground mb-6">
        {allFlags.length} flags across {commands.filter((c) => c.flags?.length).length} commands.
      </p>

      <SearchBar value={search} onChange={setSearch} placeholder="Search flags..." />

      <div className="mt-6 rounded-lg border border-border overflow-hidden">
        <div className="overflow-x-auto">
          <table className="w-full text-sm">
            <thead>
              <tr className="bg-muted/30 border-b border-border">
                <th
                  onClick={() => handleSort("flag")}
                  className="text-left px-4 py-2.5 font-mono font-semibold text-foreground cursor-pointer hover:text-primary transition-colors select-none"
                >
                  Flag{sortIndicator("flag")}
                </th>
                <th className="text-left px-4 py-2.5 font-mono font-semibold text-foreground">
                  Description
                </th>
                <th
                  onClick={() => handleSort("command")}
                  className="text-left px-4 py-2.5 font-mono font-semibold text-foreground cursor-pointer hover:text-primary transition-colors select-none"
                >
                  Command{sortIndicator("command")}
                </th>
              </tr>
            </thead>
            <tbody>
              {filtered.map((row, i) => (
                <tr
                  key={`${row.command}-${row.flag}-${i}`}
                  className="border-b border-border last:border-0 hover:bg-muted/20 transition-colors"
                >
                  <td className="px-4 py-2 font-mono text-primary whitespace-nowrap">
                    {row.flag}
                  </td>
                  <td className="px-4 py-2 text-muted-foreground">{row.description}</td>
                  <td className="px-4 py-2 font-mono whitespace-nowrap">
                    {row.command}
                    {row.alias && (
                      <span className="text-muted-foreground ml-1">({row.alias})</span>
                    )}
                  </td>
                </tr>
              ))}
              {filtered.length === 0 && (
                <tr>
                  <td colSpan={3} className="px-4 py-8 text-center text-muted-foreground font-mono text-sm">
                    No flags matching "{search}"
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      </div>
    </DocsLayout>
  );
};

export default FlagReferencePage;
