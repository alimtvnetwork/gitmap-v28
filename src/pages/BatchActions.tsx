import DocsLayout from "@/components/docs/DocsLayout";
import CodeBlock from "@/components/docs/CodeBlock";
import TerminalDemo from "@/components/docs/TerminalDemo";
import { Zap, Terminal, FolderGit2, Layers } from "lucide-react";

const ACTIONS = [
  {
    key: "p",
    name: "Pull",
    description: "Fetch and merge latest changes for every selected repository. Equivalent to running `git pull` in each repo directory.",
    examples: [
      { title: "Pull all selected repos", code: "# In TUI: select repos with Space, switch to Actions tab, press p\n# Equivalent CLI:\ngitmap pull --all" },
      { title: "Pull a specific group", code: "# In TUI: select group members, then press p\n# Equivalent CLI:\ngitmap pull --group backend" },
    ],
  },
  {
    key: "x",
    name: "Execute",
    description: "Run an arbitrary git command across all selected repositories. A prompt appears for the command string. Output is captured and displayed per-repo with a success/failure summary.",
    examples: [
      { title: "Check branch across repos", code: "# Press x, then type: branch --show-current\n# Runs `git branch --show-current` in each selected repo" },
      { title: "Stash changes everywhere", code: "# Press x, then type: stash\n# Runs `git stash` in each selected repo" },
      { title: "Fetch with prune", code: "# Press x, then type: fetch --prune\n# Runs `git fetch --prune` in each selected repo" },
    ],
  },
  {
    key: "s",
    name: "Status",
    description: "Show a compact status summary for each selected repository — dirty/clean indicator, current branch, ahead/behind counts, and stash count.",
    examples: [
      { title: "Quick status check", code: "# Select repos with Space, press s\n# Output:\n#   ✓ api-server        main    clean\n#   ✗ web-client        develop dirty  +2 -1\n#   ✓ shared-lib        main    clean  stash:1" },
    ],
  },
  {
    key: "g",
    name: "Add to Group",
    description: "Add the currently selected repositories to an existing group or create a new one. Groups persist in the database and can be used for scoped operations across all gitmap commands.",
    examples: [
      { title: "Add repos to an existing group", code: "# Select repos, press g, choose 'backend'\n# Equivalent CLI:\ngitmap group add backend repo-a repo-b" },
      { title: "Create a new group from selection", code: "# Select repos, press g, type new group name\n# Equivalent CLI:\ngitmap group create my-group\ngitmap group add my-group repo-a repo-b" },
    ],
  },
];

const WORKFLOW_STEPS = [
  { step: "1", label: "Select", description: "Use the Repos view to browse and select repositories with Space. Press a to select all, or / to filter first." },
  { step: "2", label: "Switch", description: "Press Tab to move to the Actions view. Your selection is preserved across views." },
  { step: "3", label: "Act", description: "Press the action key (p, x, s, or g). Results appear inline with a success/failure summary." },
  { step: "4", label: "Review", description: "Output is displayed in the Actions view. Press Tab to return to Repos for another round." },
];

const BatchActions = () => {
  return (
    <DocsLayout>
      <div className="max-w-4xl">
        <div className="flex items-center gap-3 mb-2">
          <Zap className="h-8 w-8 text-primary" />
          <h1 className="text-3xl font-heading font-bold text-foreground docs-h1">Batch Actions</h1>
        </div>
        <p className="text-muted-foreground mb-8 text-lg">
          Perform git operations across multiple repositories at once from the Interactive TUI's Actions view.
        </p>

        {/* Terminal Demo */}
        <section className="mb-10 ">
          <TerminalDemo
            title="gitmap interactive — batch pull"
            lines={[
              { text: "gitmap interactive", type: "input" as const, delay: 800 },
              { text: "", delay: 200 },
              { text: "gitmap TUI v2.17.0", type: "header" as const, delay: 150 },
              { text: "", delay: 100 },
              { text: "[ Repos ]  Actions  Status", type: "accent" as const, delay: 200 },
              { text: "", delay: 100 },
              { text: "  ○ myapp             main      clean", delay: 80 },
              { text: "  ○ api-server        main      dirty  +2", delay: 80 },
              { text: "  ○ shared-lib        develop   clean", delay: 80 },
              { text: "  ○ docs-site         main      clean", delay: 80 },
              { text: "  ○ cli-tools         main      dirty  +1", delay: 80 },
              { text: "", delay: 300 },
              { text: "# Select repos with Space", type: "accent" as const, delay: 400 },
              { text: "  ● myapp             main      clean", type: "accent" as const, delay: 250 },
              { text: "  ● api-server        main      dirty  +2", type: "accent" as const, delay: 250 },
              { text: "  ○ shared-lib        develop   clean", delay: 80 },
              { text: "  ○ docs-site         main      clean", delay: 80 },
              { text: "  ● cli-tools         main      dirty  +1", type: "accent" as const, delay: 250 },
              { text: "", delay: 300 },
              { text: "# Tab → Actions, press p to pull", type: "accent" as const, delay: 400 },
              { text: "", delay: 200 },
              { text: "  Repos  [ Actions ]  Status", type: "accent" as const, delay: 200 },
              { text: "", delay: 150 },
              { text: "Pulling 3 selected repos...", type: "header" as const, delay: 400 },
              { text: "  ✓ myapp             Already up to date", delay: 300 },
              { text: "  ✓ api-server        Fast-forward  2 commits", type: "accent" as const, delay: 350 },
              { text: "  ✓ cli-tools         Fast-forward  1 commit", type: "accent" as const, delay: 300 },
              { text: "", delay: 100 },
              { text: "3/3 repos pulled successfully", type: "accent" as const },
            ]}
            autoPlay
          />
        </section>

        {/* Terminal Demo 2 — exec */}
        <section className="mb-10 ">
          <TerminalDemo
            title="gitmap interactive — batch exec"
            lines={[
              { text: "# Select repos, Tab → Actions, press x", type: "accent" as const, delay: 600 },
              { text: "", delay: 200 },
              { text: "  Repos  [ Actions ]  Status", type: "accent" as const, delay: 200 },
              { text: "", delay: 150 },
              { text: "Enter git command: fetch --prune", type: "header" as const, delay: 800 },
              { text: "", delay: 200 },
              { text: "Executing `git fetch --prune` across 3 repos...", type: "header" as const, delay: 400 },
              { text: "", delay: 150 },
              { text: "  ✓ myapp             pruned 2 stale refs", type: "accent" as const, delay: 350 },
              { text: "  ✓ api-server        up to date", delay: 300 },
              { text: "  ✓ cli-tools         pruned 1 stale ref", type: "accent" as const, delay: 350 },
              { text: "", delay: 100 },
              { text: "3/3 repos executed successfully", type: "accent" as const },
            ]}
            autoPlay
          />
        </section>

        {/* Terminal Demo 3 — status */}
        <section className="mb-10 ">
          <TerminalDemo
            title="gitmap interactive — batch status"
            lines={[
              { text: "# Select repos, Tab → Actions, press s", type: "accent" as const, delay: 600 },
              { text: "", delay: 200 },
              { text: "  Repos  [ Actions ]  Status", type: "accent" as const, delay: 200 },
              { text: "", delay: 150 },
              { text: "Status for 5 selected repos:", type: "header" as const, delay: 400 },
              { text: "", delay: 150 },
              { text: "  ✓ myapp             main      clean", delay: 300 },
              { text: "  ✗ api-server        main      dirty   +2 -1", type: "accent" as const, delay: 350 },
              { text: "  ✓ shared-lib        develop   clean   stash:1", delay: 300 },
              { text: "  ✗ web-client        feature   dirty   +5", type: "accent" as const, delay: 350 },
              { text: "  ✓ cli-tools         main      clean", delay: 300 },
              { text: "", delay: 100 },
              { text: "3 clean  ·  2 dirty  ·  1 stash", type: "accent" as const },
            ]}
            autoPlay
          />
        </section>

        {/* Workflow */}
        <section className="mb-10">
          <h2 className="text-xl font-heading font-semibold text-foreground mb-4 flex items-center gap-2">
            <Layers className="h-5 w-5 text-primary" />
            Workflow
          </h2>
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-3">
            {WORKFLOW_STEPS.map((ws) => (
              <div key={ws.step} className="rounded-lg border border-border p-4 bg-card">
                <div className="flex items-center gap-2 mb-2">
                  <span className="w-6 h-6 rounded-full bg-primary text-primary-foreground text-xs font-mono font-bold flex items-center justify-center">
                    {ws.step}
                  </span>
                  <span className="font-mono font-semibold text-foreground">{ws.label}</span>
                </div>
                <p className="text-sm text-muted-foreground">{ws.description}</p>
              </div>
            ))}
          </div>
        </section>

        {/* Actions Reference */}
        <section className="mb-10">
          <h2 className="text-xl font-heading font-semibold text-foreground mb-4 flex items-center gap-2">
            <Terminal className="h-5 w-5 text-primary" />
            Actions Reference
          </h2>
          <div className="space-y-6">
            {ACTIONS.map((action) => (
              <div key={action.key} className="rounded-lg border border-border bg-card overflow-hidden">
                <div className="px-4 py-3 border-b border-border bg-muted/20 flex items-center gap-3">
                  <kbd className="px-2 py-1 rounded bg-muted text-sm font-mono font-bold border border-border">
                    {action.key}
                  </kbd>
                  <h3 className="font-mono font-semibold text-foreground">{action.name}</h3>
                </div>
                <div className="p-4">
                  <p className="text-sm text-muted-foreground mb-4">{action.description}</p>
                  <div className="space-y-3">
                    {action.examples.map((ex, i) => (
                      <div key={i}>
                        <p className="text-xs font-mono text-primary mb-1">{ex.title}</p>
                        <CodeBlock code={ex.code} language="bash" />
                      </div>
                    ))}
                  </div>
                </div>
              </div>
            ))}
          </div>
        </section>

        {/* Multi-Group Integration */}
        <section className="mb-10">
          <h2 className="text-xl font-heading font-semibold text-foreground mb-4 flex items-center gap-2">
            <FolderGit2 className="h-5 w-5 text-primary" />
            Multi-Group Integration
          </h2>
          <p className="text-muted-foreground mb-3">
            Batch actions work seamlessly with the <code className="font-mono text-primary">multi-group</code> command.
            Set active groups from the CLI, then use the TUI for interactive operations:
          </p>
          <CodeBlock
            code={`# Set active multi-group from CLI\ngitmap multi-group backend,frontend\n\n# Launch TUI — repos from both groups are pre-loaded\ngitmap interactive\n\n# Or use multi-group commands directly\ngitmap multi-group pull\ngitmap multi-group status\ngitmap multi-group exec fetch --prune`}
            language="bash"
            title="Multi-Group + TUI"
          />
        </section>

        {/* CLI Equivalents */}
        <section className="mb-10">
          <h2 className="text-xl font-heading font-semibold text-foreground mb-3">CLI Equivalents</h2>
          <p className="text-muted-foreground mb-3">
            Every TUI batch action maps to a CLI command for scripting and automation:
          </p>
          <div className="rounded-lg border border-border overflow-hidden">
            <table className="w-full text-sm">
              <thead>
                <tr className="bg-muted/30 border-b border-border">
                  <th className="text-left px-4 py-2 font-mono text-muted-foreground">TUI Key</th>
                  <th className="text-left px-4 py-2 font-mono text-muted-foreground">CLI Command</th>
                  <th className="text-left px-4 py-2 font-mono text-muted-foreground">Description</th>
                </tr>
              </thead>
              <tbody>
                <tr className="border-b border-border">
                  <td className="px-4 py-2"><kbd className="px-1.5 py-0.5 rounded bg-muted text-xs font-mono border border-border">p</kbd></td>
                  <td className="px-4 py-2 font-mono text-primary">gitmap pull --group &lt;name&gt;</td>
                  <td className="px-4 py-2 text-foreground">Pull repos in a group</td>
                </tr>
                <tr className="border-b border-border">
                  <td className="px-4 py-2"><kbd className="px-1.5 py-0.5 rounded bg-muted text-xs font-mono border border-border">x</kbd></td>
                  <td className="px-4 py-2 font-mono text-primary">gitmap exec --group &lt;name&gt; &lt;cmd&gt;</td>
                  <td className="px-4 py-2 text-foreground">Execute command across group</td>
                </tr>
                <tr className="border-b border-border">
                  <td className="px-4 py-2"><kbd className="px-1.5 py-0.5 rounded bg-muted text-xs font-mono border border-border">s</kbd></td>
                  <td className="px-4 py-2 font-mono text-primary">gitmap status --group &lt;name&gt;</td>
                  <td className="px-4 py-2 text-foreground">Show status for group</td>
                </tr>
                <tr className="border-b border-border last:border-0">
                  <td className="px-4 py-2"><kbd className="px-1.5 py-0.5 rounded bg-muted text-xs font-mono border border-border">g</kbd></td>
                  <td className="px-4 py-2 font-mono text-primary">gitmap group add &lt;name&gt; &lt;slugs&gt;</td>
                  <td className="px-4 py-2 text-foreground">Add repos to group</td>
                </tr>
              </tbody>
            </table>
          </div>
        </section>

        {/* Progress Tracking */}
        <section className="mb-10">
          <h2 className="text-xl font-heading font-semibold text-foreground mb-4 flex items-center gap-2">
            <Zap className="h-5 w-5 text-primary" />
            Progress Tracking
          </h2>
          <p className="text-muted-foreground mb-4">
            All batch operations — pull, exec, and status — display real-time progress counters with per-repo
            success/failure reporting and a final summary. Progress output goes to stderr so it can be separated
            from command output.
          </p>
          <TerminalDemo
            title="gitmap pull — progress tracking"
            lines={[
              { text: "gitmap pull --group backend", type: "input" as const, delay: 800 },
              { text: "", delay: 200 },
              { text: "  [1/4] api-server...", type: "output" as const, delay: 300 },
              { text: "  ✓ done (1.2s)", type: "accent" as const, delay: 400 },
              { text: "  [2/4] auth-service...", type: "output" as const, delay: 300 },
              { text: "  ✓ done (0.8s)", type: "accent" as const, delay: 350 },
              { text: "  [3/4] data-layer...", type: "output" as const, delay: 300 },
              { text: "  ✗ failed", type: "output" as const, delay: 200 },
              { text: "  [4/4] queue-worker...", type: "output" as const, delay: 300 },
              { text: "  — skipped", type: "output" as const, delay: 200 },
              { text: "", delay: 150 },
              { text: "  pull complete: 4/4 repos · 2.3s elapsed", type: "header" as const },
              { text: "  2 succeeded · 1 failed · 1 skipped", type: "accent" as const },
            ]}
            autoPlay
          />
          <div className="mt-4 rounded-lg border border-border p-4 bg-card">
            <h3 className="font-mono font-semibold text-sm mb-2">Progress Modes</h3>
            <ul className="space-y-2 text-sm text-muted-foreground">
              <li className="flex items-start gap-2">
                <span className="w-1.5 h-1.5 rounded-full bg-primary mt-1.5 shrink-0" />
                <strong className="text-foreground">Normal</strong> — real-time [current/total] counters with per-item status (pull, exec)
              </li>
              <li className="flex items-start gap-2">
                <span className="w-1.5 h-1.5 rounded-full bg-primary mt-1.5 shrink-0" />
                <strong className="text-foreground">Quiet</strong> — suppresses per-item output to preserve table layout (status command)
              </li>
            </ul>
          </div>
        </section>

        {/* Tips */}
        <section className="mb-10">
          <h2 className="text-xl font-heading font-semibold text-foreground mb-3">Tips</h2>
          <ul className="space-y-2 text-sm text-muted-foreground">
            <li className="flex items-start gap-2">
              <span className="w-1.5 h-1.5 rounded-full bg-primary mt-1.5 shrink-0" />
              Use <kbd className="px-1 py-0.5 rounded bg-muted text-xs font-mono border border-border">/</kbd> to filter repos before selecting — batch actions only apply to selected repos
            </li>
            <li className="flex items-start gap-2">
              <span className="w-1.5 h-1.5 rounded-full bg-primary mt-1.5 shrink-0" />
              Press <kbd className="px-1 py-0.5 rounded bg-muted text-xs font-mono border border-border">a</kbd> to select all visible (filtered) repos for a scoped batch
            </li>
            <li className="flex items-start gap-2">
              <span className="w-1.5 h-1.5 rounded-full bg-primary mt-1.5 shrink-0" />
              Selection persists across tab switches — select in Repos, act in Actions, review in Status
            </li>
            <li className="flex items-start gap-2">
              <span className="w-1.5 h-1.5 rounded-full bg-primary mt-1.5 shrink-0" />
              Progress tracking shows elapsed time per item and a final success/fail/skip summary
            </li>
          </ul>
        </section>
      </div>
    </DocsLayout>
  );
};

export default BatchActions;
