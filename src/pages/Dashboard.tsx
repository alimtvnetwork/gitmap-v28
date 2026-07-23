import DocsLayout from "@/components/docs/DocsLayout";
import CodeBlock from "@/components/docs/CodeBlock";
import TerminalDemo from "@/components/docs/TerminalDemo";
import { BarChart3, Globe, Filter, FolderOpen } from "lucide-react";

const features = [
  { icon: BarChart3, title: "Commit Analytics", desc: "Timeline charts, author breakdowns, and contribution heatmaps — all rendered client-side." },
  { icon: Globe, title: "Fully Offline", desc: "Self-contained HTML with no CDN dependencies. Works in air-gapped environments." },
  { icon: Filter, title: "Client-Side Filtering", desc: "Filter by date range, author, tag, and search commits interactively in the browser." },
  { icon: FolderOpen, title: "JSON + HTML Output", desc: "Generates both a structured JSON data file and a standalone HTML dashboard." },
];

const flags = [
  ["--limit <n>", "Maximum number of commits to include", "0 (all)"],
  ["--since <date>", "Only include commits after this date (YYYY-MM-DD)", "(none)"],
  ["--no-merges", "Exclude merge commits from the output", "false"],
  ["--out-dir <path>", "Output directory for dashboard files", ".gitmap/output"],
  ["--open", "Open the generated dashboard in the default browser", "false"],
];

const fileLayout = [
  ["constants/constants_dashboard.go", "Git formats, flag descriptions, messages"],
  ["model/dashboard.go", "DashboardData, CommitEntry, AuthorSummary structs"],
  ["dashboard/gitquery.go", "Raw git log/branch/tag extraction"],
  ["dashboard/aggregate.go", "buildAuthors, buildFrequency, attachTags"],
  ["dashboard/collector.go", "Orchestrates query → aggregate pipeline"],
  ["dashboard/writer.go", "JSON and HTML file generation"],
  ["dashboard/templates/dashboard.html", "Embedded self-contained HTML template"],
  ["cmd/dashboard.go", "CLI flag parsing and command handler"],
];

const demoLines = [
  { text: "gitmap dashboard", type: "input" as const, delay: 800 },
  { text: "Collecting repository data...", type: "output" as const },
  { text: "Wrote .gitmap/output/dashboard.json (482 commits, 7 authors)", type: "accent" as const },
  { text: "Wrote .gitmap/output/dashboard.html", type: "accent" as const },
  { text: "Dashboard generated in .gitmap/output", type: "output" as const },
];

const demoLimitLines = [
  { text: "gitmap db --limit 100 --open", type: "input" as const, delay: 800 },
  { text: "Collecting repository data...", type: "output" as const },
  { text: "Wrote .gitmap/output/dashboard.json (100 commits, 5 authors)", type: "accent" as const },
  { text: "Wrote .gitmap/output/dashboard.html", type: "accent" as const },
  { text: "Dashboard generated in .gitmap/output", type: "output" as const },
  { text: "Opening dashboard in browser...", type: "output" as const },
];

const demoFilterLines = [
  { text: "gitmap dashboard --since 2025-01-01 --no-merges --out-dir ./report", type: "input" as const, delay: 800 },
  { text: "Collecting repository data...", type: "output" as const },
  { text: "Wrote ./report/dashboard.json (63 commits, 4 authors)", type: "accent" as const },
  { text: "Wrote ./report/dashboard.html", type: "accent" as const },
  { text: "Dashboard generated in ./report", type: "output" as const },
];

const DashboardPage = () => (
  <DocsLayout>
    <h1 className="text-3xl font-heading font-bold mb-2 docs-h1">Dashboard Command</h1>
    <p className="text-muted-foreground mb-6">
      Generate an interactive, self-contained HTML dashboard for repository analytics.
      <span className="ml-2 text-xs font-mono bg-muted px-2 py-0.5 rounded">alias: db</span>
    </p>

    <h2 className="text-xl font-heading font-semibold mt-8 mb-4">Features</h2>
    <div className="grid md:grid-cols-2 gap-4 mb-8">
      {features.map((f) => (
        <div key={f.title} className="rounded-lg border border-border bg-card p-4">
          <f.icon className="h-5 w-5 text-primary mb-2" />
          <h3 className="font-mono font-semibold text-sm mb-1">{f.title}</h3>
          <p className="text-xs text-muted-foreground">{f.desc}</p>
        </div>
      ))}
    </div>

    <h2 className="text-xl font-heading font-semibold mt-10 mb-3">Interactive Examples</h2>

    <div className="space-y-6 mb-8">
      <TerminalDemo title="Full dashboard generation" lines={demoLines} autoPlay />
      <TerminalDemo title="Last 100 commits, open in browser" lines={demoLimitLines} />
      <TerminalDemo title="Filtered: since date, no merges, custom output" lines={demoFilterLines} />
    </div>

    <h2 className="text-xl font-heading font-semibold mt-10 mb-3">Usage</h2>
    <CodeBlock code="gitmap dashboard [flags]" title="Basic usage" />
    <CodeBlock code="gitmap db --limit 100 --open" title="Alias with flags" />
    <CodeBlock code="gitmap dashboard --since 2025-01-01 --no-merges --out-dir ./report" title="Filtered output" />

    <h2 className="text-xl font-heading font-semibold mt-10 mb-3">Flags</h2>
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
          {flags.map(([flag, desc, def]) => (
            <tr key={flag} className="border-t border-border">
              <td className="px-4 py-2 font-mono text-primary">{flag}</td>
              <td className="px-4 py-2 text-muted-foreground">{desc}</td>
              <td className="px-4 py-2 text-muted-foreground">{def}</td>
            </tr>
          ))}
        </tbody>
      </table>
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
          {fileLayout.map(([file, purpose]) => (
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

export default DashboardPage;
