import { useState, useMemo } from "react";
import DocsLayout from "@/components/docs/DocsLayout";
import CodeBlock from "@/components/docs/CodeBlock";
import {
  Terminal,
  Bot,
  RefreshCw,
  FileText,
  ExternalLink,
  CheckCircle2,
  Shield,
  Network,
  Copy,
  Layers,
  Search,
  Sparkles,
  X,
  Play,
  Wrench,
  Check,
} from "lucide-react";

export type PromptCategory = "All" | "Rerun" | "Listing" | "Templates" | "Pipeline Fix" | "AI Instructions";

interface AGYPromptExample {
  id: string;
  category: "Rerun" | "Listing" | "Templates" | "Pipeline Fix" | "AI Instructions";
  title: string;
  command: string;
  description: string;
  tag?: string;
}

const CATEGORY_CHIPS: PromptCategory[] = [
  "All",
  "Rerun",
  "Listing",
  "Templates",
  "Pipeline Fix",
  "AI Instructions",
];

const PROMPT_EXAMPLES: AGYPromptExample[] = [
  // Rerun
  {
    id: "rerun-1",
    category: "Rerun",
    title: "Replay last prompt with is-done template",
    command: "gitmap agy rerun last 1 -p is-done",
    description: "Replay the most recent prompt with built-in quality gate prefix and copy to OS clipboard.",
    tag: "Default Prefix",
  },
  {
    id: "rerun-5-custom",
    category: "Rerun",
    title: "Replay last 5 prompts with custom review prefix",
    command: "gitmap agy rerun last 5 -p code-review",
    description: "Prepend code-review template across recent five Antigravity prompt submissions.",
    tag: "Multi-Prompt",
  },
  {
    id: "rerun-dry-run",
    category: "Rerun",
    title: "Dry-run prompt payload preview",
    command: "gitmap agy rerun last 3 --dry-run",
    description: "Preview prompt payload directly in terminal without copying to OS clipboard or altering state.",
  },
  {
    id: "rerun-ssh-remote",
    category: "Rerun",
    title: "Remote SSH node prompt replay",
    command: "gitmap ssh exec agy rerun last 1 -p is-done",
    description: "Forward prompt replay execution to designated remote SSH node with full terminal parity.",
    tag: "SSH Triad",
  },

  // Listing
  {
    id: "list-prompts-10",
    category: "Listing",
    title: "Scan last 10 commit projects & launch VS Code",
    command: "gitmap agy list-prompts 10 --projects 10",
    description: "Collect prompt histories across 10 commit projects and open an isolated non-admin VS Code workspace.",
    tag: "VS Code (code -n)",
  },
  {
    id: "list-prompts-diff",
    category: "Listing",
    title: "Inspect prompt diff report as JSON",
    command: "gitmap agy list-prompts --summary --json",
    description: "Export structured JSON history of Antigravity agent interactions and token usage statistics.",
    tag: "JSON Export",
  },
  {
    id: "list-prompts-cwd",
    category: "Listing",
    title: "List recent prompts for current repository",
    command: "gitmap agy list-prompts 5",
    description: "List the last 5 prompts executed within the current working directory repository.",
  },

  // Templates
  {
    id: "template-ls",
    category: "Templates",
    title: "List registered prompt templates",
    command: "gitmap prompts-template ls",
    description: "Display all installed reusable prefix templates in the local registry with description and keys.",
    tag: "Registry",
  },
  {
    id: "template-add",
    category: "Templates",
    title: "Register custom verification template",
    command: 'gitmap prompts-template add test-verification "Ensure all unit tests and edge cases are verified:"',
    description: "Add a new reusable prefix template for automated test verification and quality audits.",
  },
  {
    id: "template-export-single",
    category: "Templates",
    title: "Export template to JSON file",
    command: "gitmap prompts-template export is-done ./is-done.json",
    description: "Export the built-in is-done template to an external JSON specification file.",
  },
  {
    id: "template-export-bulk",
    category: "Templates",
    title: "Bulk export all prompt templates",
    command: "gitmap prompts-template export-all ./templates.json",
    description: "Export the entire library of registered prompt templates to a single JSON archive.",
    tag: "Backup",
  },

  // Pipeline Fix
  {
    id: "pipeline-errors-fix",
    category: "Pipeline Fix",
    title: "Inject pipeline errors & 4-part RCA directly into Antigravity",
    command: "gitmap pipeline errors agy fix",
    description: "Extract failing CI/CD pipeline error logs and 4-part RCA prompts directly into active Antigravity IDE sessions.",
    tag: "RCA Injection",
  },
  {
    id: "pipeline-batch-fix-all",
    category: "Pipeline Fix",
    title: "Batch fix failing pipelines across all repositories",
    command: "gitmap pipeline errors agy fix --all",
    description: "Scan all tracked repositories and autonomously batch-fix failing pipelines across projects.",
    tag: "Batch Fix",
  },
  {
    id: "pipeline-reset-batch-cursor",
    category: "Pipeline Fix",
    title: "Reset multi-project batch cursor",
    command: "gitmap pipeline errors agy fix --all --reset-batch",
    description: "Reset the multi-repository batch cursor back to the initial failing repository.",
  },

  // AI Instructions
  {
    id: "ai-prompts-show-pretty",
    category: "AI Instructions",
    title: "Show AI instruction template with rich ANSI formatting",
    command: "gitmap agy prompts show is-done --pretty",
    description: "Render markdown-formatted AI instruction guidelines with rich ANSI terminal styling.",
    tag: "ANSI Pretty",
  },
  {
    id: "ai-agent-coordinator",
    category: "AI Instructions",
    title: "Inspect AI multi-agent coordinator leases",
    command: "gitmap agent ai status --verbose",
    description: "View Two-Tier Split-DB coordination, active worker claiming, lease fencing, and pre-touch action logging.",
    tag: "Coordinator",
  },
  {
    id: "ai-heatmap-contention",
    category: "AI Instructions",
    title: "Live agent file contention heatmap",
    command: "gitmap heatmap --watch",
    description: "Display real-time contention ranking, agent concurrency tracking, and file touch visualizer.",
    tag: "Telemetry",
  },
  {
    id: "ai-cluster-delegation",
    category: "AI Instructions",
    title: "Broadcast Antigravity scan across cluster fleet",
    command: "gitmap cluster exec agy scan",
    description: "Delegate AI scan tasks across all connected cluster fleet nodes simultaneously.",
    tag: "Cluster",
  },
];

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
  const [searchQuery, setSearchQuery] = useState("");
  const [selectedCategory, setSelectedCategory] = useState<PromptCategory>("All");
  const [copiedId, setCopiedId] = useState<string | null>(null);

  const filteredPrompts = useMemo(() => {
    return PROMPT_EXAMPLES.filter((item) => {
      const matchesCategory = selectedCategory === "All" || item.category === selectedCategory;
      if (!matchesCategory) return false;

      if (!searchQuery.trim()) return true;
      const q = searchQuery.toLowerCase().trim();

      return (
        item.title.toLowerCase().includes(q) ||
        item.command.toLowerCase().includes(q) ||
        item.description.toLowerCase().includes(q) ||
        item.category.toLowerCase().includes(q) ||
        Boolean(item.tag?.toLowerCase().includes(q))
      );
    });
  }, [searchQuery, selectedCategory]);

  const handleCopy = (id: string, text: string) => {
    navigator.clipboard.writeText(text);
    setCopiedId(id);
    setTimeout(() => {
      setCopiedId((current) => (current === id ? null : current));
    }, 2000);
  };

  const handleRunInTerminal = (cmd: string) => {
    window.dispatchEvent(new CustomEvent("gitmap:terminal-run", { detail: { command: cmd } }));
    window.dispatchEvent(new CustomEvent("gitmap:open-terminal", { detail: { command: cmd } }));
  };

  return (
    <DocsLayout>
      <div className="space-y-8 max-w-4xl">
        <div>
          <div className="flex items-center gap-2 text-sm text-primary font-mono mb-2">
            <Bot className="w-4 h-4" />
            <span>AI Orchestration</span>
          </div>
          <h1 className="text-3xl font-heading font-bold tracking-[-0.02em] text-foreground">AGY Prompts, Templates & Rerun</h1>
          <p className="text-muted-foreground mt-2 text-base leading-relaxed">
            Replay historical agent prompts with verification prefixes, inspect prompt diffs across multiple commit projects via non-admin VS Code, manage reusable prompt templates, and delegate execution remotely across SSH, Cluster, and SC nodes.
          </p>
        </div>

        {/* ── Interactive Prompt Explorer & Search Filter ── */}
        <section className="rounded-xl border border-border bg-card p-5 shadow-sm space-y-4">
          <div className="flex items-center justify-between gap-2 flex-wrap">
            <div className="flex items-center gap-2">
              <Sparkles className="h-5 w-5 text-primary" />
              <h2 className="text-lg font-heading font-semibold text-foreground">Interactive Prompt & Command Explorer</h2>
            </div>
            <span className="text-xs font-mono text-muted-foreground">
              Showing {filteredPrompts.length} of {PROMPT_EXAMPLES.length} commands
            </span>
          </div>

          {/* Search Input Filter */}
          <div className="relative">
            <Search className="absolute left-3.5 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
            <input
              type="text"
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              placeholder="Search prompts, flags, commands, or workflows (e.g. is-done, rerun, rca)..."
              className="w-full pl-10 pr-9 py-2.5 rounded-lg border border-border bg-background text-foreground text-sm font-mono placeholder:text-muted-foreground/60 focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent transition-all"
            />
            {searchQuery && (
              <button
                type="button"
                onClick={() => setSearchQuery("")}
                aria-label="Clear search"
                className="absolute right-3 top-1/2 -translate-y-1/2 p-0.5 rounded text-muted-foreground hover:text-foreground"
              >
                <X className="h-4 w-4" />
              </button>
            )}
          </div>

          {/* Category Filter Chips */}
          <div className="flex items-center gap-2 overflow-x-auto pb-1 text-xs font-mono">
            {CATEGORY_CHIPS.map((cat) => {
              const count = cat === "All"
                ? PROMPT_EXAMPLES.length
                : PROMPT_EXAMPLES.filter((p) => p.category === cat).length;
              const isActive = selectedCategory === cat;

              return (
                <button
                  key={cat}
                  type="button"
                  onClick={() => setSelectedCategory(cat)}
                  className={`px-3 py-1.5 rounded-md border transition-all shrink-0 flex items-center gap-1.5 ${
                    isActive
                      ? "border-primary bg-primary text-primary-foreground font-semibold shadow-sm"
                      : "border-border bg-muted/40 hover:bg-muted text-muted-foreground hover:text-foreground"
                  }`}
                >
                  <span>{cat}</span>
                  <span className={`text-[10px] px-1 py-0.2 rounded-full ${
                    isActive ? "bg-primary-foreground/20 text-primary-foreground" : "bg-background text-muted-foreground"
                  }`}>
                    {count}
                  </span>
                </button>
              );
            })}
          </div>

          {/* Prompts Cards Grid */}
          <div className="space-y-3 pt-2">
            {filteredPrompts.length === 0 ? (
              <div className="p-8 text-center rounded-lg border border-dashed border-border bg-muted/20">
                <Search className="h-8 w-8 text-muted-foreground mx-auto mb-2" />
                <p className="text-sm font-medium text-foreground">No prompt commands found</p>
                <p className="text-xs text-muted-foreground mt-1">Try clearing your search query or choosing another category.</p>
                <button
                  type="button"
                  onClick={() => { setSearchQuery(""); setSelectedCategory("All"); }}
                  className="mt-3 px-3 py-1.5 rounded text-xs font-mono bg-primary text-primary-foreground font-medium hover:bg-primary/90 transition-colors active:scale-[0.98]"
                >
                  Reset Filters
                </button>
              </div>
            ) : (
              filteredPrompts.map((item) => (
                <div
                  key={item.id}
                  className="rounded-lg border border-border bg-card p-4 transition-all hover:border-primary/40 hover:shadow-md space-y-3"
                >
                  <div className="flex items-start justify-between gap-3">
                    <div className="min-w-0">
                      <div className="flex items-center gap-2 flex-wrap mb-1">
                        <span className="text-[11px] font-mono px-2 py-0.5 rounded bg-primary/10 text-primary border border-primary/20">
                          {item.category}
                        </span>
                        {item.tag && (
                          <span className="text-[10px] font-mono px-1.5 py-0.5 rounded bg-muted text-muted-foreground border border-border">
                            {item.tag}
                          </span>
                        )}
                      </div>
                      <h3 className="text-sm font-semibold text-foreground">{item.title}</h3>
                      <p className="text-xs text-muted-foreground mt-0.5 leading-relaxed">{item.description}</p>
                    </div>
                  </div>

                  <div className="flex items-center justify-between gap-2 p-2.5 rounded bg-muted/60 border border-border/60 font-mono text-xs overflow-x-auto">
                    <code className="text-primary font-medium select-all">{item.command}</code>
                    <div className="flex items-center gap-1.5 shrink-0 ml-2">
                      <button
                        type="button"
                        onClick={() => handleCopy(item.id, item.command)}
                        className="inline-flex min-h-[44px] items-center gap-1 px-3 py-1.5 rounded text-xs font-mono border border-border bg-background hover:bg-muted text-muted-foreground hover:text-foreground transition-all duration-200 active:scale-[0.98]"
                        title="Copy command to clipboard"
                      >
                        {copiedId === item.id ? (
                          <>
                            <CheckCircle2 className="w-3.5 h-3.5 text-emerald-400" />
                            <span className="text-emerald-400 font-semibold">Copied!</span>
                          </>
                        ) : (
                          <>
                            <Copy className="w-3.5 h-3.5" />
                            <span>Copy</span>
                          </>
                        )}
                      </button>

                      <button
                        type="button"
                        onClick={() => handleRunInTerminal(item.command)}
                        className="inline-flex min-h-[44px] items-center gap-1 px-3 py-1.5 rounded text-xs font-mono font-medium bg-primary text-primary-foreground hover:bg-primary/90 transition-all duration-200 shadow-sm active:scale-[0.98]"
                        title="Execute command in live web terminal drawer"
                      >
                        <Terminal className="w-3.5 h-3.5" />
                        <span>Run in Terminal</span>
                      </button>
                    </div>
                  </div>
                </div>
              ))
            )}
          </div>
        </section>

        {/* ── Documentation Sections ── */}
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
