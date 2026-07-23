import DocsLayout from "@/components/docs/DocsLayout";
import CodeBlock from "@/components/docs/CodeBlock";
import { GitBranch, Trash2, Terminal, AlertTriangle } from "lucide-react";

const DryRunPreview = () => (
  <div className="rounded-lg border border-border overflow-hidden my-6">
    <div className="bg-terminal px-4 py-2 flex items-center gap-2 border-b border-border">
      <div className="flex gap-1.5">
        <span className="w-3 h-3 rounded-full bg-red-500/80" />
        <span className="w-3 h-3 rounded-full bg-yellow-500/80" />
        <span className="w-3 h-3 rounded-full bg-green-500/80" />
      </div>
      <span className="text-xs font-mono text-muted-foreground ml-2">gitmap prune --dry-run</span>
    </div>
    <div className="bg-terminal p-4 font-mono text-sm leading-relaxed overflow-x-auto text-xs">
      <div className="text-primary font-bold mb-1">{"  "}Stale release branches (3):</div>
      <div className="text-terminal-foreground">{"    "}release/v2.20.0  →  tag v2.20.0 exists</div>
      <div className="text-terminal-foreground">{"    "}release/v2.21.0  →  tag v2.21.0 exists</div>
      <div className="text-terminal-foreground">{"    "}release/v2.22.0  →  tag v2.22.0 exists</div>
      <div className="text-blue-400 mt-2">{"  "}Use --confirm to delete, or run without --dry-run for interactive mode.</div>
    </div>
  </div>
);

const DeletePreview = () => (
  <div className="rounded-lg border border-border overflow-hidden my-6">
    <div className="bg-terminal px-4 py-2 flex items-center gap-2 border-b border-border">
      <div className="flex gap-1.5">
        <span className="w-3 h-3 rounded-full bg-red-500/80" />
        <span className="w-3 h-3 rounded-full bg-yellow-500/80" />
        <span className="w-3 h-3 rounded-full bg-green-500/80" />
      </div>
      <span className="text-xs font-mono text-muted-foreground ml-2">gitmap prune --confirm</span>
    </div>
    <div className="bg-terminal p-4 font-mono text-sm leading-relaxed overflow-x-auto text-xs">
      <div className="text-muted-foreground mb-1">{"  "}Pruning stale release branches...</div>
      <div className="text-green-400">{"    "}✓ Deleted release/v2.20.0</div>
      <div className="text-green-400">{"    "}✓ Deleted release/v2.21.0</div>
      <div className="text-green-400">{"    "}✓ Deleted release/v2.22.0</div>
      <div className="text-terminal-foreground mt-2">{"  "}Summary: 3 deleted, 2 kept.</div>
    </div>
  </div>
);

const PrunePage = () => (
  <DocsLayout>
    <div className="max-w-4xl space-y-10">
      {/* Header */}
      <div>
        <div className="flex items-center gap-3 mb-2">
          <GitBranch className="h-8 w-8 text-primary" />
          <h1 className="text-3xl font-bold tracking-tight">Prune</h1>
        </div>
        <p className="text-lg text-muted-foreground">
          Delete stale release branches that have already been tagged.
        </p>
      </div>

      {/* Overview */}
      <section>
        <h2 className="text-xl font-semibold mb-3 flex items-center gap-2">
          <Trash2 className="h-5 w-5 text-primary" /> Overview
        </h2>
        <p className="text-muted-foreground mb-4">
          Over time, release branches accumulate in your local repository after each release.
          The <code className="text-primary">prune</code> command identifies branches matching{" "}
          <code>release/*</code> whose corresponding tags already exist and removes them.
        </p>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          {[
            { icon: GitBranch, title: "Auto Detection", desc: "Finds all release/* branches with existing tags" },
            { icon: AlertTriangle, title: "Safe by Default", desc: "Interactive confirmation unless --confirm is passed" },
            { icon: Terminal, title: "Remote Support", desc: "Optionally clean up remote branches with --remote" },
          ].map((f) => (
            <div key={f.title} className="rounded-lg border border-border p-4 bg-card">
              <f.icon className="h-5 w-5 text-primary mb-2" />
              <h3 className="font-semibold text-sm mb-1">{f.title}</h3>
              <p className="text-xs text-muted-foreground">{f.desc}</p>
            </div>
          ))}
        </div>
      </section>

      {/* Dry Run Preview */}
      <section>
        <h2 className="text-xl font-semibold mb-3">Dry Run Preview</h2>
        <DryRunPreview />
      </section>

      {/* Delete Preview */}
      <section>
        <h2 className="text-xl font-semibold mb-3">Deletion Output</h2>
        <DeletePreview />
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
                { flag: "--dry-run", def: "false", desc: "List stale branches without deleting" },
                { flag: "--confirm", def: "false", desc: "Skip interactive confirmation prompt" },
                { flag: "--remote", def: "false", desc: "Also delete remote release branches" },
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

      {/* Examples */}
      <section>
        <h2 className="text-xl font-semibold mb-3">Examples</h2>
        <CodeBlock code={`# Preview stale branches
gitmap prune --dry-run

# Delete without prompting
gitmap prune --confirm

# Delete locally and remotely
gitmap prune --confirm --remote`} />
      </section>

      {/* How It Works */}
      <section>
        <h2 className="text-xl font-semibold mb-3">How It Works</h2>
        <ol className="list-decimal list-inside space-y-2 text-muted-foreground">
          <li>Lists all local branches matching <code className="text-primary">release/*</code></li>
          <li>For each branch, extracts the version (e.g., <code>release/v2.20.0</code> → <code>v2.20.0</code>)</li>
          <li>Checks if the corresponding tag exists locally</li>
          <li>Branches with existing tags are marked as <strong>stale</strong></li>
          <li>Stale branches are deleted with <code>git branch -D</code></li>
          <li>With <code>--remote</code>, also runs <code>git push origin --delete</code></li>
        </ol>
      </section>
    </div>
  </DocsLayout>
);

export default PrunePage;
