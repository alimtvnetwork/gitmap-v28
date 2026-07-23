import DocsLayout from "@/components/docs/DocsLayout";
import CodeBlock from "@/components/docs/CodeBlock";
import { Replace as ReplaceIcon, Search, FileEdit, ShieldCheck } from "lucide-react";

const flags = [
  { flag: "-y, --yes", def: "false", desc: "Skip the y/N confirmation prompt" },
  { flag: "--dry-run", def: "false", desc: "Print summary only, never write" },
  { flag: "-q, --quiet", def: "false", desc: "Suppress per-file diff lines" },
  { flag: "--audit", def: "false", desc: "Report-only scan with line numbers; never writes" },
  { flag: "--ext", def: "(all)", desc: "Comma-separated extension allow-list (e.g. .go,.md). Leading dot optional." },
  { flag: "--ext-case", def: "insensitive", desc: "sensitive | insensitive — controls --ext matching" },
];

const exitCodes = [
  { code: "0", meaning: "Success — replacements applied (or dry-run/audit complete)" },
  { code: "1", meaning: "Write or scan failed" },
  { code: "2", meaning: "Bad arguments / unparseable version suffix" },
];

const ReplacePage = () => (
  <DocsLayout>
    <div className="max-w-4xl space-y-10">
      <div>
        <div className="flex items-center gap-3 mb-2">
          <ReplaceIcon className="h-8 w-8 text-primary" />
          <h1 className="text-3xl font-bold tracking-tight">replace</h1>
          <span className="font-mono text-xs px-2 py-1 rounded bg-primary/10 text-foreground border border-primary/20 dark:bg-primary/15 dark:text-primary dark:border-primary/40">
            alias: rpl
          </span>
        </div>
        <p className="text-lg text-muted-foreground">
          Repo-wide find / replace across every text file. Two modes: literal text swap, or
          version-suffix bump driven by the git remote URL. Atomic temp+rename writer, binary
          files auto-skipped, <code>.git</code> / <code>node_modules</code> / release assets
          excluded.
        </p>
        <p className="text-xs text-muted-foreground mt-2">
          Spec: <code>spec/04-generic-cli/15-replace-command.md</code>
        </p>
      </div>

      <section>
        <h2 className="text-xl font-semibold mb-3">Overview</h2>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          {[
            { icon: FileEdit, title: "Literal swap", desc: "Two positional args do a deterministic find/replace across tracked text files." },
            { icon: Search, title: "Version bump", desc: "-N or all reads the -vK suffix from origin and rewrites prior versions to the current one." },
            { icon: ShieldCheck, title: "Audit mode", desc: "--audit reports every match with line numbers and never writes — perfect for CI gates." },
          ].map((f) => (
            <div key={f.title} className="rounded-lg border border-border p-4 bg-card">
              <f.icon className="h-5 w-5 text-primary mb-2" />
              <h3 className="font-semibold text-sm mb-1">{f.title}</h3>
              <p className="text-xs text-muted-foreground">{f.desc}</p>
            </div>
          ))}
        </div>
      </section>

      <section>
        <h2 className="text-xl font-semibold mb-3">Usage</h2>
        <CodeBlock code={`gitmap replace "<old>" "<new>"     # literal text replace
gitmap replace -N                   # bump v(K-N)..v(K-1) -> vK
gitmap replace --audit              # report-only scan, no writes
gitmap replace all                  # bump v1..v(K-1) -> vK
gitmap rpl "<old>" "<new>"          # short alias`} />
      </section>

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
              {flags.map((f) => (
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

      <section>
        <h2 className="text-xl font-semibold mb-3">Examples</h2>
        <CodeBlock code={`# Literal swap, preview only
gitmap replace "old-name" "new-name" --dry-run

# Bump the previous 3 versions to the current one (reads -vK from origin)
gitmap replace -3 -y

# Bump every prior version
gitmap replace all -y

# CI audit gate: fail if any legacy URL still appears
gitmap replace --audit

# Restrict to Go + Markdown only
gitmap replace "github.com/old" "github.com/new" --ext .go,.md -y`} />
      </section>

      <section>
        <h2 className="text-xl font-semibold mb-3">Excluded paths</h2>
        <p className="text-sm text-muted-foreground">
          <code>.git</code>, <code>.gitmap</code>, <code>.release</code>,{" "}
          <code>node_modules</code>, <code>vendor</code>, <code>.gitmap/release</code>,{" "}
          <code>.gitmap/release-assets</code>, and any file whose first 8 KiB contain a
          NUL byte (treated as binary).
        </p>
      </section>

      <section>
        <h2 className="text-xl font-semibold mb-3">Exit Codes</h2>
        <div className="overflow-x-auto">
          <table className="w-full text-sm border border-border rounded-lg">
            <thead>
              <tr className="bg-muted/50">
                <th className="text-left px-4 py-2 font-medium">Code</th>
                <th className="text-left px-4 py-2 font-medium">Meaning</th>
              </tr>
            </thead>
            <tbody>
              {exitCodes.map((e) => (
                <tr key={e.code} className="border-t border-border">
                  <td className="px-4 py-2 font-mono text-primary">{e.code}</td>
                  <td className="px-4 py-2 text-muted-foreground">{e.meaning}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </section>

      <section>
        <h2 className="text-xl font-semibold mb-3">See Also</h2>
        <ul className="list-disc list-inside space-y-1 text-sm">
          <li><a href="/fix-repo" className="text-primary hover:underline">fix-repo</a> — Rewrite <code>{`{base}-vN`}</code> tokens specifically</li>
          <li><a href="/clone-fix-repo" className="text-primary hover:underline">clone-fix-repo</a> — Clone + fix-repo --all in one shot</li>
          <li><a href="/release-self" className="text-primary hover:underline">release-self</a> — Bump gitmap-v28's own version</li>
        </ul>
      </section>
    </div>
  </DocsLayout>
);

export default ReplacePage;
