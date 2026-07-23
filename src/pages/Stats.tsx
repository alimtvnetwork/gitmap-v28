import DocsLayout from "@/components/docs/DocsLayout";
import CodeBlock from "@/components/docs/CodeBlock";
import { BarChart3, TrendingUp, AlertTriangle, Zap } from "lucide-react";

const MOCK_STATS = [
  { command: "scan", runs: 42, success: 40, fail: 2, failPct: "4.8%", avg: "312", min: "189", max: "1204", last: "14:23" },
  { command: "pull", runs: 28, success: 28, fail: 0, failPct: "0.0%", avg: "845", min: "201", max: "3102", last: "14:21" },
  { command: "watch", runs: 15, success: 15, fail: 0, failPct: "0.0%", avg: "30021", min: "5004", max: "60033", last: "14:15" },
  { command: "release", runs: 8, success: 5, fail: 3, failPct: "37.5%", avg: "156", min: "42", max: "892", last: "14:19" },
  { command: "clone", runs: 6, success: 6, fail: 0, failPct: "0.0%", avg: "5230", min: "2101", max: "9842", last: "14:10" },
];

const TerminalPreview = () => (
  <div className="rounded-lg border border-border overflow-hidden my-6">
    <div className="bg-terminal px-4 py-2 flex items-center gap-2 border-b border-border">
      <div className="flex gap-1.5">
        <span className="w-3 h-3 rounded-full bg-red-500/80" />
        <span className="w-3 h-3 rounded-full bg-yellow-500/80" />
        <span className="w-3 h-3 rounded-full bg-green-500/80" />
      </div>
      <span className="text-xs font-mono text-muted-foreground ml-2">gitmap stats</span>
    </div>
    <div className="bg-terminal p-4 font-mono text-sm leading-relaxed overflow-x-auto">
      <div className="text-muted-foreground text-xs mb-2">
        Overall: 99 runs | 6 commands | 94 success | 5 failed | 5.1% fail rate | avg 4112ms
      </div>
      <div className="text-primary font-bold text-xs">
        {"  "}COMMAND{"     "}RUNS{"  "}OK{"    "}FAIL{"  "}FAIL%{"   "}AVG(ms){"  "}MIN{"    "}MAX{"      "}LAST
      </div>
      {MOCK_STATS.map((r) => (
        <div key={r.command} className="text-terminal-foreground text-xs">
          {"  "}
          <span className="inline-block w-[80px]">{r.command}</span>
          <span className="inline-block w-[40px] text-right">{r.runs}</span>
          <span className="inline-block w-[40px] text-right text-primary">{r.success}</span>
          <span className={`inline-block w-[40px] text-right ${r.fail > 0 ? "text-destructive" : "text-muted-foreground"}`}>{r.fail}</span>
          <span className={`inline-block w-[56px] text-right ${parseFloat(r.failPct) > 10 ? "text-destructive" : "text-muted-foreground"}`}>{r.failPct}</span>
          <span className="inline-block w-[64px] text-right text-muted-foreground">{r.avg}</span>
          <span className="inline-block w-[56px] text-right text-muted-foreground">{r.min}</span>
          <span className="inline-block w-[64px] text-right text-muted-foreground">{r.max}</span>
          <span className="inline-block w-[56px] text-right text-muted-foreground">{r.last}</span>
        </div>
      ))}
    </div>
  </div>
);

const features = [
  { icon: BarChart3, title: "Aggregated Metrics", desc: "Total runs, success/fail counts, failure rates — all computed from the CommandHistory table." },
  { icon: TrendingUp, title: "Performance Insights", desc: "Average, min, and max execution times help identify slow or problematic commands." },
  { icon: AlertTriangle, title: "Failure Tracking", desc: "Per-command failure percentage highlights unreliable operations at a glance." },
  { icon: Zap, title: "Zero Setup", desc: "No additional tables or configuration — stats are derived from the existing audit trail." },
];

const StatsPage = () => (
  <DocsLayout>
    <h1 className="text-3xl font-heading font-bold mb-2 docs-h1">Stats Command</h1>
    <p className="text-muted-foreground mb-6">
      Aggregated usage statistics and performance metrics for all gitmap CLI commands.
    </p>

    <h2 className="text-xl font-heading font-semibold mt-8 mb-2 docs-h2">Live Preview</h2>
    <p className="text-sm text-muted-foreground mb-2">
      Simulated terminal output showing per-command statistics.
    </p>
    <TerminalPreview />

    <h2 className="text-xl font-heading font-semibold mt-10 mb-4 docs-h2">Features</h2>
    <div className="grid md:grid-cols-2 gap-4 mb-8">
      {features.map((f) => (
        <div key={f.title} className="rounded-lg border border-border bg-card p-4">
          <f.icon className="h-5 w-5 text-primary mb-2" />
          <h3 className="font-mono font-semibold text-sm mb-1">{f.title}</h3>
          <p className="text-xs text-muted-foreground">{f.desc}</p>
        </div>
      ))}
    </div>

    <h2 className="text-xl font-heading font-semibold mt-10 mb-3 docs-h2">Usage</h2>
    <CodeBlock code="gitmap stats" title="Show all command stats" />
    <CodeBlock code="gitmap stats --command scan" title="Stats for a specific command" />
    <CodeBlock code="gitmap stats --json" title="JSON output for scripting" />

    <h2 className="text-xl font-heading font-semibold mt-10 mb-3 docs-h2">Flags</h2>
    <div className="rounded-lg border border-border overflow-hidden">
      <table className="w-full text-sm docs-table">
        <thead>
          <tr className="bg-muted/50">
            <th className="text-left font-mono font-semibold px-4 py-2">Flag</th>
            <th className="text-left font-mono font-semibold px-4 py-2">Description</th>
            <th className="text-left font-mono font-semibold px-4 py-2">Default</th>
          </tr>
        </thead>
        <tbody>
          {[
            ["--command <name>", "Show stats for a specific command only", "(all)"],
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

    <h2 className="text-xl font-heading font-semibold mt-10 mb-3 docs-h2">Output Fields</h2>
    <div className="rounded-lg border border-border overflow-hidden">
      <table className="w-full text-sm docs-table">
        <thead>
          <tr className="bg-muted/50">
            <th className="text-left font-mono font-semibold px-4 py-2">Field</th>
            <th className="text-left font-mono font-semibold px-4 py-2">Description</th>
          </tr>
        </thead>
        <tbody>
          {[
            ["Command", "Command name"],
            ["Runs", "Total number of executions"],
            ["Success", "Count of successful runs (exit code 0)"],
            ["Fail", "Count of failed runs (exit code ≠ 0)"],
            ["Fail%", "Failure rate as percentage"],
            ["Avg(ms)", "Average execution duration"],
            ["Min(ms)", "Fastest execution"],
            ["Max(ms)", "Slowest execution"],
            ["Last Used", "Timestamp of most recent execution"],
          ].map(([field, desc]) => (
            <tr key={field} className="border-t border-border">
              <td className="px-4 py-2 font-mono text-primary">{field}</td>
              <td className="px-4 py-2 text-muted-foreground">{desc}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>

    <h2 className="text-xl font-heading font-semibold mt-10 mb-3 docs-h2">File Layout</h2>
    <div className="rounded-lg border border-border overflow-hidden">
      <table className="w-full text-sm docs-table">
        <thead>
          <tr className="bg-muted/50">
            <th className="text-left font-mono font-semibold px-4 py-2">File</th>
            <th className="text-left font-mono font-semibold px-4 py-2">Purpose</th>
          </tr>
        </thead>
        <tbody>
          {[
            ["constants/constants_stats.go", "SQL queries, messages, formatting"],
            ["model/stats.go", "CommandStats and OverallStats structs"],
            ["store/stats.go", "Stats query methods"],
            ["cmd/stats.go", "Stats command handler"],
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

export default StatsPage;
