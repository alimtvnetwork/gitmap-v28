import DocsLayout from "@/components/docs/DocsLayout";
import CodeBlock from "@/components/docs/CodeBlock";
import { Clock, RotateCcw, Shield, Database } from "lucide-react";

const detailLevels = [
  { level: "basic", columns: "Command, Timestamp, Status (OK/FAIL)" },
  { level: "standard", columns: "Command, Timestamp, Flags, Status, Duration" },
  { level: "detailed", columns: "Command, Timestamp, Args, Flags, Status, Duration, Repo Count, Summary" },
];

const MOCK_HISTORY = [
  { command: "scan", alias: "s", flags: "--output json", status: "OK", duration: "342ms", time: "14:23:05" },
  { command: "pull", alias: "p", flags: "--all", status: "OK", duration: "1204ms", time: "14:21:18" },
  { command: "release", alias: "r", flags: "--bump patch", status: "FAIL", duration: "89ms", time: "14:19:42" },
  { command: "watch", alias: "w", flags: "--interval 10", status: "OK", duration: "30012ms", time: "14:15:00" },
  { command: "clone", alias: "c", flags: "--safe-pull", status: "OK", duration: "5621ms", time: "14:10:33" },
];

const TerminalPreview = () => (
  <div className="rounded-lg border border-border overflow-hidden my-6">
    <div className="bg-terminal px-4 py-2 flex items-center gap-2 border-b border-border">
      <div className="flex gap-1.5">
        <span className="w-3 h-3 rounded-full bg-red-500/80" />
        <span className="w-3 h-3 rounded-full bg-yellow-500/80" />
        <span className="w-3 h-3 rounded-full bg-green-500/80" />
      </div>
      <span className="text-xs font-mono text-muted-foreground ml-2">gitmap history</span>
    </div>
    <div className="bg-terminal p-4 font-mono text-sm leading-relaxed overflow-x-auto">
      <div className="text-primary font-bold">
        {"  "}COMMAND{"     "}ALIAS{"  "}FLAGS{"                "}STATUS{"  "}DURATION{"   "}TIME
      </div>
      {MOCK_HISTORY.map((r) => (
        <div key={r.time} className="text-terminal-foreground">
          {"  "}
          <span className="inline-block w-[88px]">{r.command}</span>
          <span className="inline-block w-[48px] text-muted-foreground">{r.alias}</span>
          <span className="inline-block w-[160px] text-muted-foreground">{r.flags}</span>
          <span className={`inline-block w-[56px] ${r.status === "FAIL" ? "text-destructive" : "text-primary"}`}>
            {r.status === "FAIL" ? "✗" : "✔"} {r.status}
          </span>
          <span className="inline-block w-[80px] text-muted-foreground">{r.duration}</span>
          <span className="text-muted-foreground">{r.time}</span>
        </div>
      ))}
      <div className="mt-3 text-muted-foreground text-xs border-t border-border/50 pt-2">
        Showing {MOCK_HISTORY.length} entries | 1 failed
      </div>
    </div>
  </div>
);

const features = [
  { icon: Clock, title: "Automatic Audit Trail", desc: "Every CLI execution is logged automatically via a two-phase hook — no manual action required." },
  { icon: Database, title: "Queryable History", desc: "Filter by command name, limit results, and choose detail levels (basic, standard, detailed)." },
  { icon: Shield, title: "Non-Blocking", desc: "Audit failures never prevent command execution. The hook is silently ignored if the DB is unavailable." },
  { icon: RotateCcw, title: "Safe Reset", desc: "Clear history with --confirm flag to prevent accidental deletion." },
];

const HistoryPage = () => (
  <DocsLayout>
    <h1 className="text-3xl font-heading font-bold mb-2 docs-h1">Command History</h1>
    <p className="text-muted-foreground mb-6">
      Automatic audit trail of every gitmap CLI execution with queryable history and usage tracking.
    </p>

    <h2 className="text-xl font-heading font-semibold mt-8 mb-2">Live Preview</h2>
    <p className="text-sm text-muted-foreground mb-2">
      Simulated terminal output of the history command in standard detail mode.
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

    <h2 className="text-xl font-heading font-semibold mt-10 mb-3">Usage</h2>
    <CodeBlock code="gitmap history" title="Show recent history (standard detail)" />
    <CodeBlock code="gitmap history --detail basic" title="Basic view" />
    <CodeBlock code="gitmap history --detail detailed --command scan" title="Detailed view filtered by command" />
    <CodeBlock code="gitmap history --json --limit 10" title="Last 10 entries as JSON" />
    <CodeBlock code="gitmap history-reset --confirm" title="Clear all history" />

    <h2 className="text-xl font-heading font-semibold mt-10 mb-3">Flags — history</h2>
    <div className="rounded-lg border border-border overflow-hidden">
      <table className="w-full text-sm">
        <thead>
          <tr className="bg-muted/50">
            <th className="text-left font-mono font-semibold px-4 py-2">Flag</th>
            <th className="text-left font-mono font-semibold px-4 py-2">Description</th>
            <th className="text-left font-mono font-semibold px-4 py-2">Default</th>
          </tr>
        </thead>
        <tbody>
          {[
            ["--detail <level>", "Output detail: basic, standard, detailed", "standard"],
            ["--command <name>", "Filter by command name", "(all)"],
            ["--limit N", "Show only the last N entries", "0 (all)"],
            ["--json", "Output as JSON", "false"],
          ].map(([flag, desc, def]) => (
            <tr key={flag} className="border-t border-border">
              <td className="px-4 py-2 font-mono text-primary">{flag}</td>
              <td className="px-4 py-2 text-muted-foreground">{desc}</td>
              <td className="px-4 py-2 text-muted-foreground">{def}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>

    <h2 className="text-xl font-heading font-semibold mt-10 mb-3">Detail Levels</h2>
    <div className="rounded-lg border border-border overflow-hidden">
      <table className="w-full text-sm">
        <thead>
          <tr className="bg-muted/50">
            <th className="text-left font-mono font-semibold px-4 py-2">Level</th>
            <th className="text-left font-mono font-semibold px-4 py-2">Columns Shown</th>
          </tr>
        </thead>
        <tbody>
          {detailLevels.map((d) => (
            <tr key={d.level} className="border-t border-border">
              <td className="px-4 py-2 font-mono text-primary">{d.level}</td>
              <td className="px-4 py-2 text-muted-foreground">{d.columns}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>

    <h2 className="text-xl font-heading font-semibold mt-10 mb-3">Audit Hook</h2>
    <p className="text-sm text-muted-foreground mb-4">
      The audit system uses a two-phase approach in <code className="text-primary font-mono">root.go</code>:
    </p>
    <div className="grid md:grid-cols-2 gap-4 mb-8">
      <div className="rounded-lg border border-border bg-card p-4">
        <h3 className="font-mono font-semibold text-sm mb-1 text-primary">1. Start Phase</h3>
        <p className="text-xs text-muted-foreground">Insert record with command name, args, flags, and start timestamp before execution.</p>
      </div>
      <div className="rounded-lg border border-border bg-card p-4">
        <h3 className="font-mono font-semibold text-sm mb-1 text-primary">2. End Phase</h3>
        <p className="text-xs text-muted-foreground">Update record with finish timestamp, duration, exit code, summary, and repo count after completion.</p>
      </div>
    </div>

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
            ["constants/constants_history.go", "SQL, command names, messages"],
            ["model/history.go", "CommandHistoryRecord struct"],
            ["store/history.go", "History CRUD operations"],
            ["cmd/history.go", "History command (display)"],
            ["cmd/historyreset.go", "History-reset command"],
            ["cmd/audit.go", "Audit hook (start/end recording)"],
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

export default HistoryPage;
