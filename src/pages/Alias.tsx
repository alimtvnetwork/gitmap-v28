import DocsLayout from "@/components/docs/DocsLayout";
import CodeBlock from "@/components/docs/CodeBlock";
import TerminalDemo from "@/components/docs/TerminalDemo";
import { Link2, Zap, Shield, Search } from "lucide-react";

const helpText = `alias (a) — Assign short names to repositories

Alias:    a
Usage:    gitmap alias <subcommand> [args]

Subcommands:
  set <alias> <slug>    Create or reassign an alias
  remove <alias>        Delete an alias
  list                  Show all defined aliases
  show <alias>          Display alias details (target repo, created date)
  suggest               Auto-suggest aliases for unaliased repos
  suggest --apply       Accept all suggestions without prompting

Flags:
  -A <alias>            Global flag — resolve target repo by alias

Prerequisites:
  • Run 'gitmap scan' first to populate the database

Examples:
  gitmap alias set api github/user/api-gateway
  gitmap a list
  gitmap alias suggest
  gitmap pull -A api
  gitmap cd -A web

See Also:
  cd         Navigate to repo by alias
  list       Show repos with slugs
  scan       Populate database for alias targets`;

const listDemo = [
  { text: "gitmap alias list", type: "input" as const, delay: 800 },
  { text: "", type: "output" as const },
  { text: "  Aliases (4):", type: "header" as const },
  { text: "  ──────────────────────────────────────────────", type: "output" as const },
  { text: "  api        → github/user/api-gateway", type: "accent" as const },
  { text: "  web        → github/user/web-frontend", type: "accent" as const },
  { text: "  infra      → github/user/infrastructure", type: "accent" as const },
  { text: "  libs       → github/user/shared-libs", type: "accent" as const },
  { text: "", type: "output" as const },
  { text: "  Hint: Use -A <alias> with any repo command", type: "output" as const },
];

const suggestDemo = [
  { text: "gitmap alias suggest", type: "input" as const, delay: 800 },
  { text: "", type: "output" as const },
  { text: "  Suggested aliases (3 unaliased repos):", type: "header" as const },
  { text: "", type: "output" as const },
  { text: "    api-gateway       → api       Accept? (y/N): y", type: "accent" as const, delay: 600 },
  { text: "    web-frontend      → web       Accept? (y/N): y", type: "accent" as const, delay: 600 },
  { text: "    shared-libs       → libs      Accept? (y/N): n", type: "output" as const, delay: 600 },
  { text: "", type: "output" as const },
  { text: "  ✓ Created 2 alias(es).", type: "accent" as const },
];

const usageDemo = [
  { text: "gitmap pull -A api", type: "input" as const, delay: 800 },
  { text: "", type: "output" as const },
  { text: "  Resolving alias: api → github/user/api-gateway", type: "output" as const },
  { text: "  → /home/user/projects/github/user/api-gateway", type: "accent" as const },
  { text: "", type: "output" as const },
  { text: "  Already up to date.", type: "output" as const },
  { text: "", type: "output" as const },
  { text: "gitmap cd -A web", type: "input" as const, delay: 800 },
  { text: "  → /home/user/projects/github/user/web-frontend", type: "accent" as const },
];

const features = [
  { icon: Link2, title: "Short Names", desc: "Replace long slugs with concise aliases like 'api', 'web', or 'infra'." },
  { icon: Zap, title: "Run From Anywhere", desc: "Execute any gitmap command against a repo using -A <alias> without changing directory." },
  { icon: Search, title: "Auto-Suggest", desc: "During scan/rescan, aliases are suggested based on repo name or slug." },
  { icon: Shield, title: "Conflict-Safe", desc: "Warns and prompts when an alias collision occurs. Cannot shadow gitmap commands." },
];

const commandInteraction = [
  ["cd", "Navigate to aliased repo directory"],
  ["pull", "Pull in aliased repo"],
  ["exec", "Execute command in aliased repo directory"],
  ["status", "Show status of aliased repo"],
  ["watch", "Watch aliased repo for changes"],
  ["release", "Run release from aliased repo"],
  ["scan", "No effect (operates on directories)"],
  ["group", "No effect (operates on group names)"],
];

const schema = [
  ["Id", "INTEGER", "PRIMARY KEY", "Auto-increment"],
  ["Alias", "TEXT", "NOT NULL UNIQUE", "Short name"],
  ["RepoId", "INTEGER", "FK → Repos(Id) CASCADE", "Target repository"],
  ["CreatedAt", "TEXT", "DEFAULT CURRENT_TIMESTAMP", ""],
];

const AliasPage = () => (
  <DocsLayout>
    <h1 className="text-3xl font-heading font-bold mb-2 docs-h1">Repo Aliases</h1>
    <p className="text-muted-foreground mb-6">
      Assign short, memorable names to repositories for quick access from anywhere.
    </p>

    <h2 className="text-xl font-heading font-semibold mt-8 mb-3 docs-h2">Help Output</h2>
    <CodeBlock code={helpText} language="bash" title="gitmap alias --help" />

    <h2 className="text-xl font-heading font-semibold mt-10 mb-3 docs-h2">Alias List</h2>
    <p className="text-sm text-muted-foreground mb-3">
      View all defined aliases with their target repos:
    </p>
    <TerminalDemo title="gitmap alias list" lines={listDemo} autoPlay />

    <h2 className="text-xl font-heading font-semibold mt-10 mb-3 docs-h2">Auto-Suggestion</h2>
    <p className="text-sm text-muted-foreground mb-3">
      Interactively suggest aliases for unaliased repos based on name analysis:
    </p>
    <TerminalDemo title="gitmap alias suggest" lines={suggestDemo} />

    <h2 className="text-xl font-heading font-semibold mt-10 mb-3 docs-h2">Using Aliases (-A Flag)</h2>
    <p className="text-sm text-muted-foreground mb-3">
      Any repo-targeting command accepts <code className="docs-inline-code">-A &lt;alias&gt;</code> to resolve the target:
    </p>
    <TerminalDemo title="gitmap — alias in action" lines={usageDemo} />

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

    <h2 className="text-xl font-heading font-semibold mt-10 mb-3 docs-h2">Subcommands</h2>
    <CodeBlock code="gitmap alias set api github/user/api-gateway" title="Create an alias" />
    <CodeBlock code="gitmap a set web github/user/web-frontend" title="Using short alias 'a'" />
    <CodeBlock code="gitmap alias list" title="List all aliases" />
    <CodeBlock code="gitmap alias show api" title="Show alias details" />
    <CodeBlock code="gitmap alias remove api" title="Remove an alias" />
    <CodeBlock code="gitmap alias suggest" title="Auto-suggest aliases" />
    <CodeBlock code="gitmap alias suggest --apply" title="Auto-accept all suggestions" />

    <h2 className="text-xl font-heading font-semibold mt-10 mb-3 docs-h2">Global -A Flag Examples</h2>
    <CodeBlock code="gitmap pull -A api" title="Pull via alias" />
    <CodeBlock code="gitmap exec -A web -- npm test" title="Run command in aliased repo" />
    <CodeBlock code="gitmap cd -A infra" title="Navigate to aliased repo" />
    <CodeBlock code="gitmap status -A api" title="Check status via alias" />

    <h2 className="text-xl font-heading font-semibold mt-10 mb-3 docs-h2">Command Interaction</h2>
    <div className="rounded-lg border border-border overflow-hidden mb-8">
      <table className="w-full text-sm docs-table">
        <thead>
          <tr className="bg-muted/50 border-b border-border">
            <th className="text-left px-4 py-2">Command</th>
            <th className="text-left px-4 py-2">Behavior with -A</th>
          </tr>
        </thead>
        <tbody>
          {commandInteraction.map(([cmd, behavior]) => (
            <tr key={cmd} className="border-t border-border">
              <td className="px-4 py-2 font-mono text-primary">{cmd}</td>
              <td className="px-4 py-2 text-muted-foreground">{behavior}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>

    <h2 className="text-xl font-heading font-semibold mt-10 mb-3 docs-h2">Conflict Handling</h2>
    <div className="grid md:grid-cols-2 gap-4 mb-8">
      <div className="rounded-lg border border-border bg-card p-4">
        <h3 className="font-heading font-semibold text-sm mb-1 text-primary">Manual Set</h3>
        <p className="text-xs text-muted-foreground">
          If alias exists, prompts: <code className="docs-inline-code">"Reassign to new repo? (y/N)"</code>
        </p>
      </div>
      <div className="rounded-lg border border-border bg-card p-4">
        <h3 className="font-heading font-semibold text-sm mb-1 text-primary">Auto-Suggest</h3>
        <p className="text-xs text-muted-foreground">
          Conflicting aliases are skipped with a warning message.
        </p>
      </div>
      <div className="rounded-lg border border-border bg-card p-4">
        <h3 className="font-heading font-semibold text-sm mb-1 text-primary">Command Shadowing</h3>
        <p className="text-xs text-muted-foreground">
          Alias names cannot match gitmap command names (e.g., <code className="docs-inline-code">scan</code>, <code className="docs-inline-code">pull</code>).
        </p>
      </div>
      <div className="rounded-lg border border-border bg-card p-4">
        <h3 className="font-heading font-semibold text-sm mb-1 text-primary">Case Sensitive</h3>
        <p className="text-xs text-muted-foreground">
          Aliases are case-sensitive: <code className="docs-inline-code">API</code> and <code className="docs-inline-code">api</code> are different aliases.
        </p>
      </div>
    </div>

    <h2 className="text-xl font-heading font-semibold mt-10 mb-3 docs-h2">Table Schema</h2>
    <div className="rounded-lg border border-border overflow-hidden mb-8">
      <table className="w-full text-sm docs-table">
        <thead>
          <tr className="bg-muted/50 border-b border-border">
            <th className="text-left px-4 py-2">Column</th>
            <th className="text-left px-4 py-2">Type</th>
            <th className="text-left px-4 py-2">Constraints</th>
            <th className="text-left px-4 py-2">Notes</th>
          </tr>
        </thead>
        <tbody>
          {schema.map(([col, type, constraints, notes]) => (
            <tr key={col} className="border-t border-border">
              <td className="px-4 py-2 font-mono text-primary">{col}</td>
              <td className="px-4 py-2 font-mono text-muted-foreground">{type}</td>
              <td className="px-4 py-2 text-muted-foreground">{constraints}</td>
              <td className="px-4 py-2 text-muted-foreground">{notes}</td>
            </tr>
          ))}
        </tbody>
      </table>
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
            ["cmd/alias.go", "Subcommand dispatch"],
            ["cmd/aliasops.go", "Subcommand implementation (set, remove, list, show)"],
            ["cmd/aliasresolve.go", "-A flag resolution logic"],
            ["store/alias.go", "Database CRUD for Aliases table"],
            ["model/alias.go", "Alias data struct"],
            ["constants/constants_alias.go", "Messages, SQL, flag descriptions"],
            ["helptext/alias.md", "Command help text"],
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
      <li><a href="/cd" className="text-primary hover:underline">cd</a> — Navigate to repo by name or alias</li>
      <li><a href="/commands" className="text-primary hover:underline">list</a> — List all tracked repos with slugs</li>
      <li><a href="/getting-started" className="text-primary hover:underline">getting-started</a> — Full setup guide</li>
    </ul>
  </DocsLayout>
);

export default AliasPage;
