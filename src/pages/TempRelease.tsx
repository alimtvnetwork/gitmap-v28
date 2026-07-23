import DocsLayout from "@/components/docs/DocsLayout";
import CodeBlock from "@/components/docs/CodeBlock";
import { GitBranch, Trash2, Terminal, Layers, Zap } from "lucide-react";

const CreatePreview = () => (
  <div className="rounded-lg border border-border overflow-hidden my-6">
    <div className="bg-terminal px-4 py-2 flex items-center gap-2 border-b border-border">
      <div className="flex gap-1.5">
        <span className="w-3 h-3 rounded-full bg-red-500/80" />
        <span className="w-3 h-3 rounded-full bg-yellow-500/80" />
        <span className="w-3 h-3 rounded-full bg-green-500/80" />
      </div>
      <span className="text-xs font-mono text-muted-foreground ml-2">gitmap tr 5 v1.$$ -s 10</span>
    </div>
    <div className="bg-terminal p-4 font-mono text-sm leading-relaxed overflow-x-auto text-xs">
      <div className="text-blue-400 mb-1">{"  "}→ Starting sequence: 10</div>
      <div className="text-muted-foreground mb-1">{"  "}Creating 5 temp-release branch(es)...</div>
      <div className="text-green-400">{"  "}✓ Created temp-release/v1.10 from abc1234</div>
      <div className="text-green-400">{"  "}✓ Created temp-release/v1.11 from def5678</div>
      <div className="text-green-400">{"  "}✓ Created temp-release/v1.12 from 789abcd</div>
      <div className="text-green-400">{"  "}✓ Created temp-release/v1.13 from bbb2222</div>
      <div className="text-green-400">{"  "}✓ Created temp-release/v1.14 from ccc3333</div>
      <div className="text-muted-foreground mt-1">{"  "}Pushing 5 branch(es) to origin...</div>
      <div className="text-green-400">{"  "}✓ Pushed 5 branch(es) to origin</div>
      <div className="text-primary mt-1">{"  "}Temp-release complete.</div>
    </div>
  </div>
);

const RemovePreview = () => (
  <div className="rounded-lg border border-border overflow-hidden my-6">
    <div className="bg-terminal px-4 py-2 flex items-center gap-2 border-b border-border">
      <div className="flex gap-1.5">
        <span className="w-3 h-3 rounded-full bg-red-500/80" />
        <span className="w-3 h-3 rounded-full bg-yellow-500/80" />
        <span className="w-3 h-3 rounded-full bg-green-500/80" />
      </div>
      <span className="text-xs font-mono text-muted-foreground ml-2">gitmap tr remove v1.10 to v1.14</span>
    </div>
    <div className="bg-terminal p-4 font-mono text-sm leading-relaxed overflow-x-auto text-xs">
      <div className="text-muted-foreground mb-1">{"  "}Remove 5 temp-release branch(es):</div>
      <div className="text-terminal-foreground">{"    "}temp-release/v1.10</div>
      <div className="text-terminal-foreground">{"    "}temp-release/v1.11</div>
      <div className="text-terminal-foreground">{"    "}temp-release/v1.12</div>
      <div className="text-terminal-foreground">{"    "}temp-release/v1.13</div>
      <div className="text-terminal-foreground">{"    "}temp-release/v1.14</div>
      <div className="text-blue-400 mt-1">{"  "}Proceed? (y/N): <span className="text-terminal-foreground">y</span></div>
      <div className="text-green-400 mt-1">{"  "}✓ Removed 5 temp-release branch(es) (local + remote)</div>
    </div>
  </div>
);

const TempReleasePage = () => (
  <DocsLayout>
    <div className="max-w-4xl space-y-10">
      {/* Header */}
      <div>
        <div className="flex items-center gap-3 mb-2">
          <Layers className="h-8 w-8 text-primary" />
          <h1 className="text-3xl font-bold tracking-tight">Temp Release</h1>
        </div>
        <p className="text-lg text-muted-foreground">
          Create lightweight, temporary release branches from recent commits — no tags, no metadata.
        </p>
      </div>

      {/* Overview */}
      <section>
        <h2 className="text-xl font-semibold mb-3 flex items-center gap-2">
          <Zap className="h-5 w-5 text-primary" /> Overview
        </h2>
        <p className="text-muted-foreground mb-4">
          The <code className="text-primary">temp-release</code> command spins up candidate release
          branches from your recent commits for quick experimentation. Unlike{" "}
          <code className="text-primary">release</code>, it creates <strong>no tags</strong> and
          writes <strong>no metadata</strong>. When you're done, clean up with{" "}
          <code className="text-primary">tr remove</code>.
        </p>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          {[
            { icon: GitBranch, title: "No Checkout", desc: "Branches created via git branch <name> <sha> — stays on current branch" },
            { icon: Layers, title: "Batch Push", desc: "All branches pushed to origin in a single git push command" },
            { icon: Trash2, title: "Easy Cleanup", desc: "Remove single, range, or all branches with confirmation prompts" },
          ].map((f) => (
            <div key={f.title} className="rounded-lg border border-border p-4 bg-card">
              <f.icon className="h-5 w-5 text-primary mb-2" />
              <h3 className="font-semibold text-sm mb-1">{f.title}</h3>
              <p className="text-xs text-muted-foreground">{f.desc}</p>
            </div>
          ))}
        </div>
      </section>

      {/* Version Pattern */}
      <section>
        <h2 className="text-xl font-semibold mb-3">Version Pattern</h2>
        <p className="text-muted-foreground mb-4">
          Use <code className="text-primary">$</code> placeholders for zero-padded sequence numbers.
          The number of <code>$</code> characters determines the digit width:
        </p>
        <div className="overflow-x-auto">
          <table className="w-full text-sm border border-border rounded-lg">
            <thead>
              <tr className="bg-muted/50">
                <th className="text-left px-4 py-2 font-medium">Pattern</th>
                <th className="text-left px-4 py-2 font-medium">Digits</th>
                <th className="text-left px-4 py-2 font-medium">Example Output</th>
              </tr>
            </thead>
            <tbody>
              {[
                { pattern: "v1.$$", digits: "2", example: "v1.05, v1.12, v1.99" },
                { pattern: "v2.$$$", digits: "3", example: "v2.005, v2.012, v2.099" },
                { pattern: "v3.$$$$", digits: "4", example: "v3.0005, v3.0012, v3.0099" },
              ].map((p) => (
                <tr key={p.pattern} className="border-t border-border">
                  <td className="px-4 py-2 font-mono text-primary">{p.pattern}</td>
                  <td className="px-4 py-2 text-muted-foreground">{p.digits}</td>
                  <td className="px-4 py-2 font-mono text-muted-foreground">{p.example}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </section>

      {/* Create Preview */}
      <section>
        <h2 className="text-xl font-semibold mb-3">Branch Creation</h2>
        <CreatePreview />
      </section>

      {/* Remove Preview */}
      <section>
        <h2 className="text-xl font-semibold mb-3">Removing Branches</h2>
        <RemovePreview />
      </section>

      {/* Flags */}
      <section>
        <h2 className="text-xl font-semibold mb-3">Flags</h2>
        <div className="overflow-x-auto">
          <table className="w-full text-sm border border-border rounded-lg">
            <thead>
              <tr className="bg-muted/50">
                <th className="text-left px-4 py-2 font-medium">Flag</th>
                <th className="text-left px-4 py-2 font-medium">Default</th>
                <th className="text-left px-4 py-2 font-medium">Description</th>
              </tr>
            </thead>
            <tbody>
              {[
                { flag: "-s, --start", def: "auto", desc: "Starting sequence number" },
                { flag: "--dry-run", def: "false", desc: "Preview branch names without creating" },
                { flag: "--json", def: "false", desc: "JSON output for list subcommand" },
                { flag: "--verbose", def: "false", desc: "Detailed logging" },
              ].map((f) => (
                <tr key={f.flag} className="border-t border-border">
                  <td className="px-4 py-2 font-mono text-primary">{f.flag}</td>
                  <td className="px-4 py-2 font-mono text-muted-foreground">{f.def}</td>
                  <td className="px-4 py-2 text-muted-foreground">{f.desc}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </section>

      {/* Subcommands */}
      <section>
        <h2 className="text-xl font-semibold mb-3">Subcommands</h2>
        <div className="overflow-x-auto">
          <table className="w-full text-sm border border-border rounded-lg">
            <thead>
              <tr className="bg-muted/50">
                <th className="text-left px-4 py-2 font-medium">Command</th>
                <th className="text-left px-4 py-2 font-medium">Description</th>
              </tr>
            </thead>
            <tbody>
              {[
                { cmd: "tr list [--json]", desc: "List all temp-release branches with SHA, message, and date" },
                { cmd: "tr remove <version>", desc: "Remove a single temp-release branch" },
                { cmd: "tr remove <v1> to <v2>", desc: "Remove a range of branches (inclusive)" },
                { cmd: "tr remove all", desc: "Remove all temp-release branches" },
              ].map((s) => (
                <tr key={s.cmd} className="border-t border-border">
                  <td className="px-4 py-2 font-mono text-primary">{s.cmd}</td>
                  <td className="px-4 py-2 text-muted-foreground">{s.desc}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </section>

      {/* Examples */}
      <section>
        <h2 className="text-xl font-semibold mb-3">Examples</h2>
        <CodeBlock code={`# Create 10 branches from last 10 commits, starting at sequence 5
gitmap tr 10 v1.$$ -s 5

# Create 1 branch, auto-increment from last temp-release
gitmap tr 1 v1.$$

# Preview without creating
gitmap tr 5 v2.$$$ --dry-run

# List all temp-release branches
gitmap tr list

# Remove a single branch
gitmap tr remove v1.05

# Remove a range
gitmap tr remove v1.05 to v1.10

# Remove all
gitmap tr remove all`} />
      </section>

      {/* How It Works */}
      <section>
        <h2 className="text-xl font-semibold mb-3">How It Works</h2>
        <ol className="list-decimal list-inside space-y-2 text-muted-foreground">
          <li>Fetches the last <strong>N</strong> commits from <code className="text-primary">git log</code> (oldest first)</li>
          <li>Resolves the starting sequence from <code className="text-primary">-s</code> flag or auto-detects from the database</li>
          <li>Creates branches using <code className="text-primary">git branch &lt;name&gt; &lt;sha&gt;</code> — <strong>no checkout</strong></li>
          <li>Pushes all branches to origin in a single batch command</li>
          <li>Records each branch in the <code className="text-primary">TempReleases</code> SQLite table</li>
          <li>On removal, deletes both local and remote branches and cleans up DB records</li>
        </ol>
      </section>
    </div>
  </DocsLayout>
);

export default TempReleasePage;
