import { useState, useMemo, useRef, useCallback } from "react";
import { Copy, Check, Download } from "lucide-react";
import DocsLayout from "@/components/docs/DocsLayout";
import CommandCard from "@/components/docs/CommandCard";
import CommandCategoryGroup from "@/components/docs/CommandCategoryGroup";
import SearchBar from "@/components/docs/SearchBar";
import { commands, Categories } from "@/data/commands";
import { copyToClipboard } from "@/lib/clipboard";

const CommandsPage = () => {
  const [search, setSearch] = useState("");
  const [forceOpen, setForceOpen] = useState<string | null>(null);
  const [highlightCmd, setHighlightCmd] = useState<string | null>(null);
  const categoryRefs = useRef<Record<string, HTMLDivElement | null>>({});
  const commandRefs = useRef<Record<string, HTMLDivElement | null>>({});

  const scrollToCategory = useCallback((key: string) => {
    setSearch("");
    setForceOpen(key);
    setTimeout(() => {
      categoryRefs.current[key]?.scrollIntoView({ behavior: "smooth", block: "start" });
      setForceOpen(null);
    }, 50);
  }, []);

  const handleNavigate = useCallback((commandName: string) => {
    const target = commands.find((c) => c.name === commandName);
    if (!target) return;

    setSearch("");
    setForceOpen(target.category);
    setHighlightCmd(commandName);

    setTimeout(() => {
      const el = commandRefs.current[commandName];
      if (el) {
        el.scrollIntoView({ behavior: "smooth", block: "center" });
        el.classList.add("ring-2", "ring-primary/50", "rounded-lg");
        setTimeout(() => {
          el.classList.remove("ring-2", "ring-primary/50", "rounded-lg");
          setHighlightCmd(null);
        }, 1500);
      }
      setForceOpen(null);
    }, 100);
  }, []);

  const filtered = useMemo(() => {
    if (!search) return commands;
    const q = search.toLowerCase();
    return commands.filter(
      (c) =>
        c.name.includes(q) ||
        c.alias?.includes(q) ||
        c.description.toLowerCase().includes(q)
    );
  }, [search]);

  const isSearching = search.length > 0;

  const [copied, setCopied] = useState(false);

  const generateMarkdown = useCallback(() => {
    let md = `# gitmap Command Reference\n\n`;
    md += `> ${commands.length} commands organized by category.\n\n`;
    Categories.forEach((cat) => {
      const cmds = commands.filter((c) => c.category === cat.key);
      if (cmds.length === 0) return;
      md += `## ${cat.icon || ""} ${cat.label}\n\n${cat.description}\n\n`;
      cmds.forEach((cmd) => {
        md += `### \`${cmd.name}\`${cmd.alias ? ` (alias: \`${cmd.alias}\`)` : ""}\n\n`;
        md += `${cmd.description}\n\n`;
        if (cmd.usage) md += `**Usage:**\n\`\`\`\n${cmd.usage}\n\`\`\`\n\n`;
        if (cmd.flags?.length) {
          md += `**Flags:**\n`;
          cmd.flags.forEach((f) => { md += `- \`${f.flag}\` — ${f.description}\n`; });
          md += `\n`;
        }
        if (cmd.examples?.length) {
          md += `**Examples:**\n`;
          cmd.examples.forEach((ex) => {
            if (ex.description) md += `${ex.description}:\n`;
            md += `\`\`\`bash\n${ex.command}\n\`\`\`\n\n`;
          });
        }
      });
      md += `---\n\n`;
    });
    return md;
  }, []);

  const handleCopyAll = useCallback(async () => {
    await copyToClipboard(generateMarkdown());
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  }, [generateMarkdown]);

  const handleDownloadMd = useCallback(() => {
    const blob = new Blob([generateMarkdown()], { type: "text/markdown" });
    const url = URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = "gitmap-commands.md";
    a.click();
    URL.revokeObjectURL(url);
  }, [generateMarkdown]);

  return (
    <DocsLayout>
      <div className="flex items-center justify-between mb-2">
        <h1 className="text-3xl font-heading font-bold docs-h1">Command Reference</h1>
        <div className="flex items-center gap-1">
          <button
            onClick={handleCopyAll}
            className="p-2 rounded-lg border border-border text-muted-foreground hover:text-foreground hover:bg-muted/50 transition-colors"
            title="Copy all as Markdown"
          >
            {copied ? <Check className="h-4 w-4 text-primary" /> : <Copy className="h-4 w-4" />}
          </button>
          <button
            onClick={handleDownloadMd}
            className="p-2 rounded-lg border border-border text-muted-foreground hover:text-foreground hover:bg-muted/50 transition-colors"
            title="Download as .md"
          >
            <Download className="h-4 w-4" />
          </button>
        </div>
      </div>
      <p className="text-muted-foreground mb-6">
        All {commands.length} gitmap commands organized by category.
      </p>

      <div className="mb-8 p-6 bg-muted/30 rounded-lg border border-border">
        <h2 className="text-xl font-bold mb-3 docs-h2">Cluster Command Delegation</h2>
        <p className="text-sm text-muted-foreground mb-4">
          Once <code className="text-primary">gitmap serve</code> and <code className="text-primary">gitmap join</code> establish a cluster of machines, operators use a unified CLI surface to broadcast shell commands, Git operations, or package installations across the network.
        </p>
        <ul className="list-disc list-inside text-sm text-muted-foreground space-y-2 mb-4">
          <li><strong>Target Selectors</strong>: <code className="text-primary">servers-clients</code> targets all nodes including the server, while <code className="text-primary">clients</code> excludes the server.</li>
          <li><strong>Filtering</strong>: Use <code className="text-primary">--except &lt;list&gt;</code> to exclude nodes by IP or ID, or <code className="text-primary">--ip</code> and <code className="text-primary">--id</code> to target exactly.</li>
          <li><strong>Chaining</strong>: You can chain multiple sub-commands sequentially on targeted nodes (e.g., <code className="text-primary">gitmap clients cmd "whoami", ps "Get-Date" --id 1,3</code>).</li>
          <li><strong>Audit Trail</strong>: All executions maintain a persistent history verifiable via <code className="text-primary">gitmap cluster history</code>.</li>
        </ul>
        <div className="mt-4 p-4 border border-red-500/20 bg-red-500/10 rounded-md">
          <h3 className="font-semibold text-red-500 mb-2">Security Note</h3>
          <p className="text-sm text-foreground">
            Client node credentials used for lifecycle commands are securely stored as bcrypt hashes. Passwords are never exported (e.g., via <code className="text-primary">cluster export</code>), and all cluster communication is strictly encrypted over TLS.
          </p>
        </div>
      </div>

      {/* Category summary banner */}
      <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-5 gap-2 mb-6">
        {Categories.map((cat) => {
          const count = commands.filter((c) => c.category === cat.key).length;
          return (
            <button
              key={cat.key}
              onClick={() => scrollToCategory(cat.key)}
              className="rounded-lg border border-border bg-card px-3 py-2.5 text-left hover:bg-muted/50 hover:border-primary/40 transition-all duration-200 cursor-pointer group"
            >
              <div className="flex items-center gap-2 mb-1">
                {cat.icon && <span className="text-base">{cat.icon}</span>}
                <span className="text-lg font-mono font-bold text-primary">{count}</span>
              </div>
              <div className="text-[11px] text-muted-foreground font-mono leading-tight truncate group-hover:text-foreground transition-colors">{cat.label}</div>
            </button>
          );
        })}
      </div>

      <SearchBar value={search} onChange={setSearch} />

      <div className="mt-6 space-y-3">
        {isSearching ? (
          <>
            {filtered.map((cmd) => (
              <div key={cmd.name} ref={(el) => { commandRefs.current[cmd.name] = el; }}>
                <CommandCard {...cmd} onNavigate={handleNavigate} />
              </div>
            ))}
            {filtered.length === 0 && (
              <p className="text-center text-muted-foreground py-8 font-mono text-sm">
                No commands matching "{search}"
              </p>
            )}
          </>
        ) : (
          Categories.map((cat) => {
            const cmds = filtered.filter((c) => c.category === cat.key);
            if (cmds.length === 0) return null;
            return (
              <div key={cat.key} ref={(el) => { categoryRefs.current[cat.key] = el; }}>
                <CommandCategoryGroup
                  label={cat.label}
                  description={cat.description}
                  icon={cat.icon}
                  commands={cmds}
                  forceOpen={forceOpen === cat.key}
                  onNavigate={handleNavigate}
                  commandRefs={commandRefs}
                />
              </div>
            );
          })
        )}
      </div>
    </DocsLayout>
  );
};

export default CommandsPage;
