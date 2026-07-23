import DocsLayout from "@/components/docs/DocsLayout";
import CodeBlock from "@/components/docs/CodeBlock";
import { GitCompareArrows, Database, FileJson, Eye, AlertTriangle } from "lucide-react";

const TerminalPreview = () => (
  <div className="rounded-lg border border-border overflow-hidden my-6">
    <div className="bg-terminal px-4 py-2 flex items-center gap-2 border-b border-border">
      <div className="flex gap-1.5">
        <span className="w-3 h-3 rounded-full bg-red-500/80" />
        <span className="w-3 h-3 rounded-full bg-yellow-500/80" />
        <span className="w-3 h-3 rounded-full bg-green-500/80" />
      </div>
      <span className="text-xs font-mono text-muted-foreground ml-2">gitmap diff-profiles default work</span>
    </div>
    <div className="bg-terminal p-4 font-mono text-sm leading-relaxed">
      <div className="text-muted-foreground mb-2">Comparing profiles: default ↔ work</div>
      <div className="text-yellow-400 font-bold text-xs mt-2">ONLY IN default:</div>
      <div className="text-foreground ml-2">docs-site{"          "}C:\repos\docs-site</div>
      <div className="text-yellow-400 font-bold text-xs mt-2">ONLY IN work:</div>
      <div className="text-foreground ml-2">client-portal{"      "}D:\work\client-portal</div>
      <div className="text-yellow-400 font-bold text-xs mt-2">DIFFERENT:</div>
      <div className="text-foreground ml-2">shared-lib</div>
      <div className="text-muted-foreground ml-4">default: C:\repos\github\shared-lib{"  "}(https)</div>
      <div className="text-muted-foreground ml-4">work:{"    "}D:\work\shared-lib{"           "}(ssh)</div>
      <div className="text-muted-foreground mt-2">SAME: 12 repos (use --all to show)</div>
      <div className="text-green-400 mt-2 text-xs">Summary: 1 only-left | 1 only-right | 1 different | 12 same</div>
    </div>
  </div>
);

const features = [
  { icon: GitCompareArrows, title: "Side-by-Side Comparison", desc: "Compares repos between any two profiles by name matching." },
  { icon: Eye, title: "Difference Categories", desc: "Only-in-A, only-in-B, different config, and same — all clearly labeled." },
  { icon: FileJson, title: "JSON Output", desc: "Structured JSON output via --json for scripting and automation." },
  { icon: AlertTriangle, title: "Validation", desc: "Exits with error if either profile doesn't exist." },
];

const categories = [
  { label: "Only in A", color: "text-yellow-400", desc: "Repository exists in the first profile but not the second." },
  { label: "Only in B", color: "text-yellow-400", desc: "Repository exists in the second profile but not the first." },
  { label: "Different", color: "text-primary", desc: "Exists in both but path or URL mode differs." },
  { label: "Same", color: "text-muted-foreground", desc: "Identical in both profiles. Hidden by default; use --all to show." },
];

const DiffProfilesPage = () => (
  <DocsLayout>
    <div className="max-w-4xl">
      <div className="flex items-center gap-3 mb-2">
        <GitCompareArrows className="h-8 w-8 text-primary" />
        <div>
          <h1 className="text-3xl font-heading font-bold text-foreground docs-h1">Diff Profiles</h1>
          <p className="text-muted-foreground font-mono text-sm">gitmap diff-profiles (dp)</p>
        </div>
      </div>
      <p className="text-muted-foreground mb-8 text-lg">
        Compare tracked repositories between two database profiles to highlight additions, removals, and configuration differences.
      </p>

      <TerminalPreview />

      {/* Features */}
      <section className="mb-10">
        <h2 className="text-xl font-heading font-bold text-foreground mb-4 docs-h2">Features</h2>
        <div className="grid sm:grid-cols-2 gap-4">
          {features.map((f) => (
            <div key={f.title} className="border border-border rounded-lg p-4 bg-card">
              <div className="flex items-center gap-2 mb-2">
                <f.icon className="h-5 w-5 text-primary" />
                <span className="font-mono font-semibold text-foreground">{f.title}</span>
              </div>
              <p className="text-sm text-muted-foreground">{f.desc}</p>
            </div>
          ))}
        </div>
      </section>

      {/* Resolution Logic */}
      <section className="mb-10">
        <h2 className="text-xl font-heading font-bold text-foreground mb-4 docs-h2">Profile Resolution Logic</h2>
        <div className="border border-border rounded-lg p-4 bg-card">
          <div className="flex items-center gap-3 mb-3 flex-wrap">
            <span className="px-2 py-1 bg-primary/10 text-foreground border border-primary/20 dark:bg-primary/15 dark:text-primary rounded text-xs font-mono transition-colors duration-300 hover:border-primary/40 hover:shadow-sm hover:shadow-primary/10">1. Open Profile A DB</span>
            <span className="text-muted-foreground">→</span>
            <span className="px-2 py-1 bg-primary/10 text-foreground border border-primary/20 dark:bg-primary/15 dark:text-primary rounded text-xs font-mono transition-colors duration-300 hover:border-primary/40 hover:shadow-sm hover:shadow-primary/10">2. Open Profile B DB</span>
            <span className="text-muted-foreground">→</span>
            <span className="px-2 py-1 bg-primary/10 text-foreground border border-primary/20 dark:bg-primary/15 dark:text-primary rounded text-xs font-mono transition-colors duration-300 hover:border-primary/40 hover:shadow-sm hover:shadow-primary/10">3. Load repo lists</span>
            <span className="text-muted-foreground">→</span>
            <span className="px-2 py-1 bg-primary/10 text-foreground border border-primary/20 dark:bg-primary/15 dark:text-primary rounded text-xs font-mono transition-colors duration-300 hover:border-primary/40 hover:shadow-sm hover:shadow-primary/10">4. Match by RepoName</span>
          </div>
          <p className="text-sm text-muted-foreground mb-3">
            Each profile name resolves to its database file via <code className="text-primary">store.OpenProfile()</code>.
            The <code className="text-primary">default</code> profile maps to <code className="text-primary">gitmap.db</code>;
            all others map to <code className="text-primary">gitmap-&lt;name&gt;.db</code>. Both databases are opened
            read-only, repo lists are loaded, and repos are matched by their unique <code className="text-primary">RepoName</code> identifier.
          </p>
          <CodeBlock code={`# Profile A: "default" → gitmap.db\n# Profile B: "work"    → gitmap-work.db\n# Match key: RepoName (unique slug)`} />
        </div>
      </section>

      {/* Comparison Categories */}
      <section className="mb-10">
        <h2 className="text-xl font-heading font-bold text-foreground mb-4 docs-h2">Comparison Categories</h2>
        <div className="space-y-3">
          {categories.map((cat) => (
            <div key={cat.label} className="border border-border rounded-lg p-4 bg-card flex items-start gap-3">
              <span className={`font-mono font-bold text-sm whitespace-nowrap ${cat.color}`}>{cat.label}</span>
              <p className="text-sm text-muted-foreground">{cat.desc}</p>
            </div>
          ))}
        </div>
      </section>

      {/* Usage */}
      <section className="mb-10">
        <h2 className="text-xl font-heading font-bold text-foreground mb-4 docs-h2">Usage</h2>
        <CodeBlock code="gitmap diff-profiles <profileA> <profileB> [--all] [--json]" />
        <div className="mt-4">
          <h3 className="font-mono font-semibold text-foreground mb-2">Flags</h3>
          <div className="overflow-x-auto">
            <table className="w-full text-sm font-mono">
              <thead>
                <tr className="border-b border-border text-muted-foreground">
                  <th className="text-left py-2 pr-4">Flag</th>
                  <th className="text-left py-2 pr-4">Default</th>
                  <th className="text-left py-2">Description</th>
                </tr>
              </thead>
              <tbody>
                <tr className="border-b border-border">
                  <td className="py-2 pr-4 text-primary">--all</td>
                  <td className="py-2 pr-4 text-muted-foreground">false</td>
                  <td className="py-2 text-foreground">Show all repos, not just differences</td>
                </tr>
                <tr className="border-b border-border">
                  <td className="py-2 pr-4 text-primary">--json</td>
                  <td className="py-2 pr-4 text-muted-foreground">false</td>
                  <td className="py-2 text-foreground">Output as structured JSON</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </section>

      {/* JSON Output */}
      <section className="mb-10">
        <h2 className="text-xl font-heading font-bold text-foreground mb-4 docs-h2">JSON Output</h2>
        <p className="text-sm text-muted-foreground mb-3">
          Use <code className="text-primary">--json</code> for structured output suitable for scripting and CI pipelines:
        </p>
        <CodeBlock code={`{\n  "profileA": "default",\n  "profileB": "work",\n  "onlyInA": [{"name": "docs-site", "path": "C:\\\\repos\\\\docs-site"}],\n  "onlyInB": [{"name": "client-portal", "path": "D:\\\\work\\\\client-portal"}],\n  "different": [\n    {\n      "name": "shared-lib",\n      "a": {"path": "C:\\\\repos\\\\github\\\\shared-lib", "mode": "https"},\n      "b": {"path": "D:\\\\work\\\\shared-lib", "mode": "ssh"}\n    }\n  ],\n  "same": 12\n}`} />
      </section>

      {/* Examples */}
      <section className="mb-10">
        <h2 className="text-xl font-heading font-bold text-foreground mb-4 docs-h2">Examples</h2>
        <div className="space-y-4">
          <div>
            <p className="text-sm text-muted-foreground mb-1">Compare two profiles:</p>
            <CodeBlock code="gitmap diff-profiles default work" />
          </div>
          <div>
            <p className="text-sm text-muted-foreground mb-1">Include identical repos:</p>
            <CodeBlock code="gitmap dp default work --all" />
          </div>
          <div>
            <p className="text-sm text-muted-foreground mb-1">JSON output for automation:</p>
            <CodeBlock code="gitmap dp work personal --json" />
          </div>
        </div>
      </section>

      {/* File Layout */}
      <section className="mb-10">
        <h2 className="text-xl font-heading font-bold text-foreground mb-4 docs-h2">File Layout</h2>
        <div className="overflow-x-auto">
          <table className="w-full text-sm font-mono">
            <thead>
              <tr className="border-b border-border text-muted-foreground">
                <th className="text-left py-2 pr-4">File</th>
                <th className="text-left py-2">Purpose</th>
              </tr>
            </thead>
            <tbody>
              {[
                ["constants/constants_diffprofile.go", "Command names, headers, messages"],
                ["cmd/diffprofiles.go", "Command entry and flag parsing"],
                ["cmd/diffprofilesops.go", "Comparison logic and output formatting"],
              ].map(([file, purpose]) => (
                <tr key={file} className="border-b border-border">
                  <td className="py-2 pr-4 text-primary">{file}</td>
                  <td className="py-2 text-foreground">{purpose}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </section>

      {/* See Also */}
      <section className="mb-10">
        <h2 className="text-xl font-heading font-bold text-foreground mb-4 docs-h2">See Also</h2>
        <ul className="space-y-1 text-sm font-mono">
          <li><a href="/profile" className="text-primary hover:underline">profile</a> — Create and manage database profiles</li>
          <li><a href="/export" className="text-primary hover:underline">export</a> — Export current profile data</li>
          <li><a href="/commands" className="text-primary hover:underline">list</a> — View repos in the current profile</li>
        </ul>
      </section>
    </div>
  </DocsLayout>
);

export default DiffProfilesPage;
