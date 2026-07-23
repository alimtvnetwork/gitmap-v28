import DocsLayout from "@/components/docs/DocsLayout";
import CodeBlock from "@/components/docs/CodeBlock";
import { Terminal, Keyboard, Layers, Settings2 } from "lucide-react";

const KEY_BINDINGS = [
  { key: "q / Esc", context: "Global", action: "Quit TUI" },
  { key: "Tab", context: "Global", action: "Switch between views" },
  { key: "/", context: "Repos", action: "Focus search input" },
  { key: "j / ↓", context: "Repos / Groups", action: "Move cursor down" },
  { key: "k / ↑", context: "Repos / Groups", action: "Move cursor up" },
  { key: "Space", context: "Repos", action: "Toggle repo selection" },
  { key: "a", context: "Repos", action: "Select all repos" },
  { key: "Enter", context: "Repos", action: "Show repo detail" },
  { key: "p", context: "Actions", action: "Pull selected repos" },
  { key: "x", context: "Actions", action: "Execute command on selected" },
  { key: "s", context: "Actions", action: "Show status for selected" },
  { key: "g", context: "Actions", action: "Add selected to group" },
  { key: "c", context: "Groups", action: "Create new group" },
  { key: "d", context: "Groups", action: "Delete group (with confirm)" },
  { key: "r", context: "Status", action: "Force refresh dashboard" },
  { key: "Enter", context: "Releases", action: "Toggle release detail view" },
  { key: "r", context: "Releases", action: "Refresh release list" },
  { key: "/", context: "Logs", action: "Filter by command, alias, args, or exit code" },
  { key: "Enter", context: "Logs", action: "Toggle log detail view" },
  { key: "r", context: "Logs", action: "Refresh log list" },
];

const VIEWS = [
  {
    name: "Repos",
    icon: "📂",
    description: "Browse all tracked repositories with fuzzy search. Type to filter, select multiple repos for batch actions.",
    features: ["Fuzzy search filtering", "Multi-select with Space", "Detail panel on Enter", "Column display: Slug, Branch, Path, Type"],
  },
  {
    name: "Actions",
    icon: "⚡",
    description: "Perform batch git operations on selected repositories from the Repos view.",
    features: ["Pull selected repos", "Execute arbitrary git commands", "Show status summary", "Add selection to a group"],
  },
  {
    name: "Groups",
    icon: "📁",
    description: "Manage repository groups with inline creation and deletion.",
    features: ["List groups with member counts", "Create groups inline", "Delete with confirmation", "Navigate members with Enter"],
  },
  {
    name: "Status",
    icon: "📊",
    description: "Live-refreshing dashboard showing git status for all tracked repositories.",
    features: ["Dirty/clean indicators", "Branch ahead/behind counts", "Stash count per repo", "Configurable auto-refresh interval"],
  },
  {
    name: "Releases",
    icon: "🏷️",
    description: "Browse release history stored in the database with detail view for each version.",
    features: ["Version, tag, branch, date columns", "Detail view with changelog and notes", "Draft and pre-release indicators", "Refresh from database"],
  },
];

const InteractiveTUI = () => {
  return (
    <DocsLayout>
      <div className="max-w-4xl">
        <div className="flex items-center gap-3 mb-2">
          <Terminal className="h-8 w-8 text-primary" />
          <h1 className="text-3xl font-heading font-bold text-foreground docs-h1">Interactive TUI</h1>
        </div>
        <p className="text-muted-foreground mb-8 text-lg">
          A full-screen terminal interface for browsing, searching, and managing repositories — built with Bubble Tea.
        </p>

        {/* Usage */}
        <section className="mb-10">
          <h2 className="text-xl font-heading font-semibold text-foreground mb-3 flex items-center gap-2">
            <Settings2 className="h-5 w-5 text-primary" />
            Usage
          </h2>
          <CodeBlock code={`gitmap interactive [--refresh <seconds>]\ngitmap i [--refresh <seconds>]`} language="bash" title="Command" />
        </section>

        {/* Flags */}
        <section className="mb-10">
          <h2 className="text-xl font-heading font-semibold text-foreground mb-3">Flags</h2>
          <div className="rounded-lg border border-border overflow-hidden">
            <table className="w-full text-sm">
              <thead>
                <tr className="bg-muted/30 border-b border-border">
                  <th className="text-left px-4 py-2 font-mono text-muted-foreground">Flag</th>
                  <th className="text-left px-4 py-2 font-mono text-muted-foreground">Default</th>
                  <th className="text-left px-4 py-2 font-mono text-muted-foreground">Description</th>
                </tr>
              </thead>
              <tbody>
                <tr className="border-b border-border">
                  <td className="px-4 py-2 font-mono text-primary">--refresh</td>
                  <td className="px-4 py-2 font-mono text-muted-foreground">config or 30</td>
                  <td className="px-4 py-2 text-foreground">Dashboard auto-refresh interval in seconds</td>
                </tr>
              </tbody>
            </table>
          </div>
          <p className="text-sm text-muted-foreground mt-2">
            The refresh interval can also be set via <code className="font-mono text-primary">dashboardRefresh</code> in <code className="font-mono text-primary">config.json</code>. The CLI flag takes priority.
          </p>
        </section>

        {/* Views */}
        <section className="mb-10">
          <h2 className="text-xl font-heading font-semibold text-foreground mb-4 flex items-center gap-2">
            <Layers className="h-5 w-5 text-primary" />
            Views
          </h2>
          <p className="text-muted-foreground mb-4">
            The TUI has four views, accessible via <kbd className="px-1.5 py-0.5 rounded bg-muted text-xs font-mono border border-border">Tab</kbd>:
          </p>
          <div className="grid gap-4">
            {VIEWS.map((view) => (
              <div key={view.name} className="rounded-lg border border-border p-4 bg-card">
                <h3 className="font-mono font-semibold text-foreground mb-1 flex items-center gap-2">
                  <span>{view.icon}</span>
                  {view.name}
                </h3>
                <p className="text-sm text-muted-foreground mb-3">{view.description}</p>
                <ul className="grid grid-cols-1 sm:grid-cols-2 gap-1">
                  {view.features.map((f) => (
                    <li key={f} className="text-sm text-foreground/80 flex items-center gap-1.5">
                      <span className="w-1.5 h-1.5 rounded-full bg-primary shrink-0" />
                      {f}
                    </li>
                  ))}
                </ul>
              </div>
            ))}
          </div>
        </section>

        {/* Key Bindings */}
        <section className="mb-10">
          <h2 className="text-xl font-heading font-semibold text-foreground mb-4 flex items-center gap-2">
            <Keyboard className="h-5 w-5 text-primary" />
            Key Bindings
          </h2>
          <div className="rounded-lg border border-border overflow-hidden">
            <table className="w-full text-sm">
              <thead>
                <tr className="bg-muted/30 border-b border-border">
                  <th className="text-left px-4 py-2 font-mono text-muted-foreground">Key</th>
                  <th className="text-left px-4 py-2 font-mono text-muted-foreground">Context</th>
                  <th className="text-left px-4 py-2 font-mono text-muted-foreground">Action</th>
                </tr>
              </thead>
              <tbody>
                {KEY_BINDINGS.map((kb, i) => (
                  <tr key={i} className="border-b border-border last:border-0 hover:bg-muted/20 transition-colors">
                    <td className="px-4 py-2">
                      <kbd className="px-1.5 py-0.5 rounded bg-muted text-xs font-mono border border-border">{kb.key}</kbd>
                    </td>
                    <td className="px-4 py-2 text-muted-foreground">{kb.context}</td>
                    <td className="px-4 py-2 text-foreground">{kb.action}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </section>

        {/* Configuration */}
        <section className="mb-10">
          <h2 className="text-xl font-heading font-semibold text-foreground mb-3">Configuration</h2>
          <p className="text-muted-foreground mb-3">
            Set the default dashboard refresh interval in your <code className="font-mono text-primary">config.json</code>:
          </p>
          <CodeBlock
            code={`{\n  "dashboardRefresh": 30\n}`}
            language="json"
            title="config.json"
          />
          <p className="text-sm text-muted-foreground">
            Values ≤ 0 fall back to the default of 30 seconds.
          </p>
        </section>

        {/* Requirements */}
        <section className="mb-10">
          <h2 className="text-xl font-heading font-semibold text-foreground mb-3">Requirements</h2>
          <ul className="space-y-2 text-sm text-muted-foreground">
            <li className="flex items-start gap-2">
              <span className="w-1.5 h-1.5 rounded-full bg-primary mt-1.5 shrink-0" />
              Requires a terminal with alternate screen support
            </li>
            <li className="flex items-start gap-2">
              <span className="w-1.5 h-1.5 rounded-full bg-primary mt-1.5 shrink-0" />
              Falls back to an error message if the terminal doesn't support TUI mode
            </li>
            <li className="flex items-start gap-2">
              <span className="w-1.5 h-1.5 rounded-full bg-primary mt-1.5 shrink-0" />
              Run <code className="font-mono text-primary">gitmap scan</code> first to populate the repository database
            </li>
          </ul>
        </section>
      </div>
    </DocsLayout>
  );
};

export default InteractiveTUI;
