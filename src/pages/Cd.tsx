import DocsLayout from "@/components/docs/DocsLayout";
import CodeBlock from "@/components/docs/CodeBlock";
import TerminalDemo from "@/components/docs/TerminalDemo";
import { FolderOpen, Search, Terminal, Layers, Zap } from "lucide-react";

const directJumpDemo = [
  { text: "gitmap cd myrepo", type: "input" as const, delay: 800 },
  { text: "", type: "output" as const },
  { text: "  → /home/user/projects/github/user/myrepo", type: "accent" as const },
  { text: "", type: "output" as const },
  { text: "pwd", type: "input" as const, delay: 600 },
  { text: "  /home/user/projects/github/user/myrepo", type: "output" as const },
];

const fuzzyDemo = [
  { text: "gitmap cd apii", type: "input" as const, delay: 800 },
  { text: "", type: "output" as const },
  { text: '  ✗ No repo found matching "apii"', type: "output" as const },
  { text: "", type: "output" as const },
  { text: "  Did you mean?", type: "header" as const },
  { text: "    1. api-gateway         (distance: 1)", type: "accent" as const },
  { text: "    2. api-client          (distance: 2)", type: "accent" as const },
  { text: "    3. app-infra           (distance: 3)", type: "output" as const },
  { text: "", type: "output" as const },
  { text: "  Select [1-3] or press Enter to cancel: 1", type: "input" as const, delay: 1000 },
  { text: "  → /home/user/projects/github/user/api-gateway", type: "accent" as const },
];

const pickerDemo = [
  { text: "gitmap cd repos", type: "input" as const, delay: 800 },
  { text: "", type: "output" as const },
  { text: "  REPOS (42 tracked)", type: "header" as const },
  { text: "  ──────────────────", type: "output" as const },
  { text: "  > api-gateway           github/user/api-gateway", type: "accent" as const },
  { text: "    web-frontend          github/user/web-frontend", type: "output" as const },
  { text: "    shared-libs           github/user/shared-libs", type: "output" as const },
  { text: "    infrastructure        github/user/infrastructure", type: "output" as const },
  { text: "", type: "output" as const },
  { text: "  Type to filter · ↑↓ navigate · Enter select", type: "output" as const },
];

const features = [
  { icon: FolderOpen, title: "Direct Navigation", desc: "Jump to any tracked repo by name or slug — no need to remember full paths." },
  { icon: Search, title: "Fuzzy Suggestions", desc: "When no exact match is found, suggests closest repos ranked by edit distance." },
  { icon: Layers, title: "Interactive Picker", desc: "Use 'repos' keyword for an fzf-style interactive picker with type-to-filter." },
  { icon: Zap, title: "Shell Wrapper (gcd)", desc: "Auto-installed shell function that performs the actual directory change." },
];

const flags = [
  ["--group <name>", "Filter repos by group in interactive picker", "(all)"],
  ["--pick", "Force interactive picker even with exact match", "false"],
  ["--verbose", "Show full path resolution details", "false"],
];

const helpText = `cd (go) — Navigate to a tracked repo directory

Alias:    go
Usage:    gitmap cd <repo-name|repos> [flags]

Arguments:
  <repo-name>    Exact repo name, slug, or alias to navigate to
  repos          Open interactive picker (fzf-style)

Flags:
  --group <name>   Filter picker by group
  --pick           Force interactive picker
  --verbose        Show path resolution details

Prerequisites:
  • Run 'gitmap scan' first to populate the database
  • Run 'gitmap setup' to install the gcd shell wrapper

Examples:
  gitmap cd myrepo              Jump to repo directory
  gitmap cd repos               Interactive repo picker
  gitmap cd repos --group work  Pick from work group only
  gcd myrepo                    Shell wrapper (same behavior)
  gcd repos                     Shell wrapper with picker

See Also:
  list       List all tracked repos with slugs
  scan       Scan directories to populate database
  alias      Assign short names to repos
  group      Manage repo groups`;

const CdPage = () => (
  <DocsLayout>
    <h1 className="text-3xl font-heading font-bold mb-2 docs-h1">CD Command</h1>
    <p className="text-muted-foreground mb-6">
      Navigate your shell to any tracked repository directory by name, with fuzzy matching and interactive picking.
    </p>

    <h2 className="text-xl font-heading font-semibold mt-8 mb-3 docs-h2">Help Output</h2>
    <CodeBlock code={helpText} language="bash" title="gitmap cd --help" />

    <h2 className="text-xl font-heading font-semibold mt-10 mb-3 docs-h2">Direct Navigation</h2>
    <p className="text-sm text-muted-foreground mb-3">
      Jump directly to a repo by its name, slug, or alias:
    </p>
    <TerminalDemo title="gitmap cd — direct jump" lines={directJumpDemo} autoPlay />

    <h2 className="text-xl font-heading font-semibold mt-10 mb-3 docs-h2">Fuzzy Matching</h2>
    <p className="text-sm text-muted-foreground mb-3">
      When no exact match is found, gitmap suggests the closest repos ranked by Levenshtein distance:
    </p>
    <TerminalDemo title="gitmap cd — fuzzy suggestion" lines={fuzzyDemo} />

    <h2 className="text-xl font-heading font-semibold mt-10 mb-3 docs-h2">Interactive Picker</h2>
    <p className="text-sm text-muted-foreground mb-3">
      Use <code className="docs-inline-code">repos</code> as the argument for an fzf-style picker with type-to-filter:
    </p>
    <TerminalDemo title="gitmap cd repos — interactive picker" lines={pickerDemo} />

    <h2 className="text-xl font-heading font-semibold mt-10 mb-4 docs-h2">Features</h2>
    <div className="grid md:grid-cols-2 gap-4 mb-8">
      {features.map((f) => (
        <div key={f.title} className="rounded-lg border border-border bg-card p-4">
          <f.icon className="h-5 w-5 text-primary mb-2" />
          <h3 className="font-heading font-semibold text-sm mb-1">{f.title}</h3>
          <p className="text-xs text-muted-foreground">{f.desc}</p>
        </div>
      ))}
    </div>

    <h2 className="text-xl font-heading font-semibold mt-10 mb-3 docs-h2">Flags</h2>
    <div className="rounded-lg border border-border overflow-hidden mb-8">
      <table className="w-full text-sm docs-table">
        <thead>
          <tr className="bg-muted/50 border-b border-border">
            <th className="text-left px-4 py-2">Flag</th>
            <th className="text-left px-4 py-2">Description</th>
            <th className="text-left px-4 py-2">Default</th>
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

    <h2 className="text-xl font-heading font-semibold mt-10 mb-3 docs-h2">Shell Wrapper — gcd</h2>
    <p className="text-sm text-muted-foreground mb-3">
      The <code className="docs-inline-code">gcd</code> function is auto-installed by <code className="docs-inline-code">gitmap setup</code> into
      your shell profile. It captures the stdout path and performs the actual <code className="docs-inline-code">cd</code>:
    </p>
    <CodeBlock
      language="bash"
      title="~/.bashrc (auto-installed)"
      code={`gcd() {
  local dir
  dir=$(gitmap cd "$@" 2>/dev/null)
  if [ -n "$dir" ] && [ -d "$dir" ]; then
    cd "$dir" || return
  else
    gitmap cd "$@"
  fi
}`}
    />
    <p className="text-sm text-muted-foreground mt-3">
      All interactive prompts and errors go to <code className="docs-inline-code">stderr</code>, keeping <code className="docs-inline-code">stdout</code> clean for path resolution.
    </p>

    <h2 className="text-xl font-heading font-semibold mt-10 mb-3 docs-h2">Examples</h2>
    <CodeBlock code="gitmap cd myrepo" title="Jump to repo by name" />
    <CodeBlock code="gitmap go myrepo" title="Using alias 'go'" />
    <CodeBlock code="gitmap cd repos" title="Interactive picker" />
    <CodeBlock code="gitmap cd repos --group backend" title="Pick from group" />
    <CodeBlock code="gcd myrepo" title="Shell wrapper" />
    <CodeBlock code="gitmap cd -A api" title="Navigate via alias" />

    <h2 className="text-xl font-heading font-semibold mt-10 mb-3 docs-h2">Resolution Order</h2>
    <div className="space-y-2 text-sm text-muted-foreground mb-8">
      <div className="flex gap-2"><span className="text-primary font-mono">1.</span> Exact repo name match</div>
      <div className="flex gap-2"><span className="text-primary font-mono">2.</span> Exact slug match</div>
      <div className="flex gap-2"><span className="text-primary font-mono">3.</span> Alias resolution (if <code className="docs-inline-code">-A</code> flag used)</div>
      <div className="flex gap-2"><span className="text-primary font-mono">4.</span> Fuzzy suggestion (Levenshtein distance ≤ 3)</div>
      <div className="flex gap-2"><span className="text-primary font-mono">5.</span> Error with "not found" message</div>
    </div>

    <h2 className="text-xl font-heading font-semibold mt-10 mb-3 docs-h2">File Layout</h2>
    <div className="rounded-lg border border-border overflow-hidden">
      <table className="w-full text-sm docs-table">
        <thead>
          <tr className="bg-muted/50 border-b border-border">
            <th className="text-left px-4 py-2">File</th>
            <th className="text-left px-4 py-2">Purpose</th>
          </tr>
        </thead>
        <tbody>
          {[
            ["cmd/cd.go", "Command handler and argument dispatch"],
            ["cmd/cdpicker.go", "Interactive repo picker logic"],
            ["cmd/cdfuzzy.go", "Fuzzy matching and suggestion engine"],
            ["store/repo.go", "Database queries for repo lookup"],
            ["constants/constants_cd.go", "Messages, prompts, format strings"],
            ["helptext/cd.md", "Command help text"],
          ].map(([file, purpose]) => (
            <tr key={file} className="border-t border-border">
              <td className="px-4 py-2 font-mono text-primary">{file}</td>
              <td className="px-4 py-2 text-muted-foreground">{purpose}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>

    <h2 className="text-xl font-heading font-semibold mt-10 mb-3 docs-h2">See Also</h2>
    <ul className="space-y-1 text-sm font-mono">
      <li><a href="/alias" className="text-primary hover:underline">alias</a> — Assign short names to repos for quick access</li>
      <li><a href="/commands" className="text-primary hover:underline">list</a> — List all tracked repos with slugs</li>
      <li><a href="/getting-started" className="text-primary hover:underline">getting-started</a> — Setup guide including gcd</li>
    </ul>
  </DocsLayout>
);

export default CdPage;
