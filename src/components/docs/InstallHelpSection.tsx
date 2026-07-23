import CodeBlock from "@/components/docs/CodeBlock";
import { BadgeCheck, ClipboardList, HelpCircle, ListChecks, Route, ShieldCheck } from "lucide-react";

const HELP_COMMAND = `gitmap install --help
gitmap install --list
gitmap install node --dry-run
gitmap install node --manager choco --version 22.5.0`;

const HELP_BLOCKS = [
  [BadgeCheck, "Recommended path", "Use `gitmap install <tool>` first. It detects the platform, chooses the manager, verifies the binary, then records the install.", "recommended"],
  [ListChecks, "Discovery commands", "`--list`, `--status`, and `--check` answer what is supported, what is installed, and whether a target is already present.", "power"],
  [ShieldCheck, "Safe execution", "`--dry-run`, `--manager`, and `--version` make the command predictable before it touches the machine.", "pinned"],
] as const;

const FLAG_GROUPS = [
  ["Inspect", "--help", "--list", "--status", "--check"],
  ["Control", "--manager", "--version", "--upgrade"],
  ["Safety", "--dry-run", "--verbose"],
] as const;

const COLOR_CLASS = {
  recommended: "border-help-recommended/30 bg-help-recommended/10 text-help-recommended",
  power: "border-help-power/30 bg-help-power/10 text-help-power",
  pinned: "border-help-pinned/30 bg-help-pinned/10 text-help-pinned",
} as const;

function HelpBlock({ block }: { block: (typeof HELP_BLOCKS)[number] }) {
  const [Icon, title, body, color] = block;
  return (
    <div className={`border-l-4 ${COLOR_CLASS[color]} bg-card p-4`}>
      <Icon className="mb-3 h-5 w-5" />
      <h3 className="mb-1 font-mono text-sm font-semibold text-foreground">{title}</h3>
      <p className="text-xs leading-5 text-muted-foreground">{body}</p>
    </div>
  );
}

function FlagRail() {
  return (
    <div className="grid gap-3 md:grid-cols-3">
      {FLAG_GROUPS.map(([title, ...flags]) => (
        <div key={title} className="border border-border bg-muted/40 p-3">
          <div className="mb-2 flex items-center gap-2 text-xs font-semibold text-foreground">
            <Route className="h-3.5 w-3.5 text-primary" />{title}
          </div>
          <div className="flex flex-wrap gap-2">{flags.map((flag) => <code key={flag} className="docs-inline-code">{flag}</code>)}</div>
        </div>
      ))}
    </div>
  );
}

export default function InstallHelpSection() {
  return (
    <section className="mb-10 border border-border bg-card/70 p-5 shadow-sm">
      <div className="mb-4 flex items-start gap-3">
        <div className="border border-primary/30 bg-primary/10 p-2 text-primary"><HelpCircle className="h-5 w-5" /></div>
        <div>
          <h2 className="text-xl font-heading font-semibold docs-h2">CLI help</h2>
          <p className="text-sm text-muted-foreground">The install help path should make the safe default, inspection commands, and controlled installs obvious at a glance.</p>
        </div>
      </div>

      <div className="grid gap-4 md:grid-cols-3">{HELP_BLOCKS.map((block) => <HelpBlock key={block[1]} block={block} />)}</div>
      <div className="my-5"><FlagRail /></div>
      <div className="flex items-center gap-2 text-sm font-semibold text-foreground"><ClipboardList className="h-4 w-4 text-primary" />Quick help commands</div>
      <CodeBlock code={HELP_COMMAND} language="bash" title="install help workflow" />
    </section>
  );
}