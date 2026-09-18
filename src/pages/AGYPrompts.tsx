import DocsLayout from "@/components/docs/DocsLayout";
import CodeBlock from "@/components/docs/CodeBlock";
import { Terminal, Bot, RefreshCw, FileText, ExternalLink, CheckCircle2, Shield, Network, Copy, Layers } from "lucide-react";

const RerunPreview = () => (
  <div className="rounded-lg border border-border overflow-hidden my-6">
    <div className="bg-terminal px-4 py-2 flex items-center gap-2 border-b border-border">
      <div className="flex gap-1.5">
        <span className="w-3 h-3 rounded-full bg-red-500/80" />
        <span className="w-3 h-3 rounded-full bg-yellow-500/80" />
        <span className="w-3 h-3 rounded-full bg-green-500/80" />
      </div>
      <span className="text-xs font-mono text-muted-foreground ml-2">gitmap agy rerun last 1 -p is-done</span>
    </div>
    <div className="bg-terminal p-4 font-mono text-sm leading-relaxed overflow-x-auto text-xs">
      <div className="text-green-400">  ✔ Loaded prefix template: [is-done]</div>
      <div className="text-muted-foreground mt-1">  ────────────────────────────────────────────────────────────────────────────</div>
      <div className="text-cyan-400 font-semibold mt-1">  Constructed Replay Prompt:</div>
      <div className="text-terminal-foreground mt-2 bg-black/40 p-3 rounded border border-border/40">
        Is it done properly? Can we check properly the missing items from the task that is mentioned below? Please check it carefully. Do not make any mistakes.
        <br /><br />
        &lt;USER_REQUEST&gt;<br />
        gitmap agy rerun last N<br />
        gitmap agy list-prompts 10 --projects 10<br />
        &lt;/USER_REQUEST&gt;
      </div>
      <div className="text-green-400 mt-2">  ✔ Copied constructed prompt to OS clipboard (ready to paste)</div>
    </div>
  </div>
);

const ListPromptsPreview = () => (
  <div className="rounded-lg border border-border overflow-hidden my-6">
    <div className="bg-terminal px-4 py-2 flex items-center gap-2 border-b border-border">
      <div className="flex gap-1.5">
        <span className="w-3 h-3 rounded-full bg-red-500/80" />
        <span className="w-3 h-3 rounded-full bg-yellow-500/80" />
        <span className="w-3 h-3 rounded-full bg-green-500/80" />
      </div>
      <span className="text-xs font-mono text-muted-foreground ml-2">gitmap agy list-prompts 10 --projects 10</span>
    </div>
    <div className="bg-terminal p-4 font-mono text-sm leading-relaxed overflow-x-auto text-xs">
      <div className="text-cyan-400 font-bold">  ● Scanning Last 10 Projects with Commits...</div>
      <div className="text-muted-foreground mt-1">  Aggregated 10 projects: [gitmap-v28, macro-ahk, wp-onboarding, lara-licensing-v3, ...]</div>
      <div className="text-green-400 mt-2">  ✔ Generated prompt diff report: C:\Users\ADMINI~1\AppData\Local\Temp\agy-prompts-10.md</div>
      <div className="text-blue-400 mt-1">  ✔ Launched non-admin VS Code workspace: <span className="underline">code -n agy-prompts-10.md</span></div>
      <div className="text-muted-foreground mt-1">  ℹ Active window preserved; launched isolated in new editor window without elevation.</div>
    </div>
  </div>
);

export default function AGYPromptsPage() {
  return (
    <DocsLayout>
      <div className="space-y-8 max-w-4xl">
        <div>
          <div className="flex items-center gap-2 text-sm text-primary font-mono mb-2">
            <Bot className="w-4 h-4" />
            <span>AI Orchestration</span>
          </div>
          <h1 className="text-3xl font-bold tracking-tight text-foreground">AGY Prompts, Templates & Rerun</h1>
          <p className="text-muted-foreground mt-2 text-base leading-relaxed">
            Replay historical agent prompts with verification prefixes, inspect prompt diffs across multiple commit projects via non-admin VS Code, manage reusable prompt templates, and delegate execution remotely across SSH, Cluster, and SC nodes.
          </p>
        </div>

        <section className="space-y-4">
          <h2 className="text-xl font-semibold text-foreground flex items-center gap-2">
            <RefreshCw className="w-5 h-5 text-primary" />
            Prompt Replay with Prefix Templates
          </h2>
          <p className="text-sm text-muted-foreground leading-relaxed">
            Quickly replay recent prompts sent to Google Antigravity. GitMap automatically prefixes your prompt with verification criteria, like the built-in <code className="text-xs bg-muted px-1.5 py-0.5 rounded text-primary">is-done</code> template, ensuring thorough quality checks.
          </p>
          <RerunPreview />
          <CodeBlock
            code={`# Replay last prompt with default is-done prefix\ngitmap agy rerun last 1\n\n# Replay last 5 prompts prefixed with custom template\ngitmap agy rerun last 5 -p code-review\n\n# Preview prompt payload without copying to clipboard\ngitmap agy rerun last 3 --dry-run`}
            language="bash"
          />
        </section>

        <section className="space-y-4">
          <h2 className="text-xl font-semibold text-foreground flex items-center gap-2">
            <FileText className="w-5 h-5 text-primary" />
            Cross-Project Prompt Inspection & VS Code Opening
          </h2>
          <p className="text-sm text-muted-foreground leading-relaxed">
            When reviewing prompt evolutions across projects, <code className="text-xs bg-muted px-1.5 py-0.5 rounded text-primary">--projects &lt;N&gt;</code> collects prompt histories for the last N commit projects and launches VS Code in a non-admin, non-intrusive new window.
          </p>
          <ListPromptsPreview />
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4 my-4">
            <div className="p-4 rounded-lg border border-border bg-card">
              <div className="flex items-center gap-2 font-semibold text-sm mb-1 text-foreground">
                <ExternalLink className="w-4 h-4 text-blue-400" />
                New Window Isolation (`code -n`)
              </div>
              <p className="text-xs text-muted-foreground leading-relaxed">
                VS Code opens via <code className="text-primary font-mono">-n</code>, ensuring your existing active editor workspace and tabs remain completely undisturbed.
              </p>
            </div>
            <div className="p-4 rounded-lg border border-border bg-card">
              <div className="flex items-center gap-2 font-semibold text-sm mb-1 text-foreground">
                <Shield className="w-4 h-4 text-green-400" />
                Non-Admin Security
              </div>
              <p className="text-xs text-muted-foreground leading-relaxed">
                VS Code is strictly launched in non-elevated user mode to protect workspace settings and prevent accidental file permission changes.
              </p>
            </div>
          </div>
        </section>

        <section className="space-y-4">
          <h2 className="text-xl font-semibold text-foreground flex items-center gap-2">
            <Layers className="w-5 h-5 text-primary" />
            Prompts Template Registry
          </h2>
          <p className="text-sm text-muted-foreground leading-relaxed">
            Create, edit, and import/export reusable prompt templates formatted as JSON. The built-in <code className="text-xs bg-muted px-1.5 py-0.5 rounded text-primary">is-done</code> template is installed by default.
          </p>
          <CodeBlock
            code={`# List all prompt templates\ngitmap prompts-template ls\n\n# Add a new custom prompt template\ngitmap prompts-template add test-verification "Ensure all unit tests and edge cases are verified:"\n\n# Export template to JSON file\ngitmap prompts-template export is-done ./is-done.json\n\n# Bulk export all templates\ngitmap prompts-template export-all ./templates.json`}
            language="bash"
          />
        </section>

        <section className="space-y-4">
          <h2 className="text-xl font-semibold text-foreground flex items-center gap-2">
            <Network className="w-5 h-5 text-primary" />
            Remote Triad Delegation
          </h2>
          <p className="text-sm text-muted-foreground leading-relaxed">
            Execute Antigravity commands and schedules across remote machines with full terminal parity:
          </p>
          <CodeBlock
            code={`# Run agy rerun on remote SSH node\ngitmap ssh exec agy rerun last 1 -p is-done\n\n# Run agy scan across entire cluster\ngitmap cluster exec agy scan\n\n# Run agy commands via SC delegation\ngitmap sc exec agy scan`}
            language="bash"
          />
        </section>

        <section className="space-y-4">
          <h2 className="text-xl font-semibold text-foreground flex items-center gap-2">
            <Terminal className="w-5 h-5 text-primary" />
            Direct Antigravity Injection & Multi-Project Batching
          </h2>
          <p className="text-sm text-muted-foreground leading-relaxed">
            Feed failing CI/CD pipeline error logs and 4-part RCA prompts directly into Antigravity IDE sessions, stage follow-up verification prompts in the queue ledger, and batch-fix across multiple tracked repositories:
          </p>
          <CodeBlock
            code={`# Extract errors and inject fix task directly into Antigravity IDE\ngitmap pipeline errors agy fix\n\n# Scan all tracked repositories and batch-fix failing pipelines (default limit: 3)\ngitmap pipeline errors agy fix --all\n\n# Advance to next batch of failing projects on subsequent run\ngitmap pipeline errors agy fix --all\n\n# Reset multi-project batch cursor\ngitmap pipeline errors agy fix --all --reset-batch`}
            language="bash"
          />
        </section>
      </div>
    </DocsLayout>
  );
}
