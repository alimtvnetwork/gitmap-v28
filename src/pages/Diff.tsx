import DocsLayout from "@/components/docs/DocsLayout";
import CodeBlock from "@/components/docs/CodeBlock";
import { GitCompare, FileSearch, ShieldCheck, Terminal } from "lucide-react";

const flags = [
  { flag: "--json", def: "false", desc: "Emit a JSON object {summary, entries} instead of text" },
  { flag: "--only-conflicts", def: "false", desc: "Show only files that differ on both sides" },
  { flag: "--only-missing", def: "false", desc: "Show only files present on one side" },
  { flag: "--include-identical", def: "false", desc: "Include byte-equal files in the output" },
  { flag: "--include-vcs", def: "false", desc: "Walk .git/ (default: skipped)" },
  { flag: "--include-node-modules", def: "false", desc: "Walk node_modules/ (default: skipped)" },
];

const exitCodes = [
  { code: "0", meaning: "Diff produced (regardless of whether differences were found)" },
  { code: "1", meaning: "One endpoint missing, not a directory, or walk failed" },
  { code: "2", meaning: "Wrong number of positional arguments (need exactly LEFT RIGHT)" },
];

const TerminalPreview = () => (
  <div className="rounded-lg border border-border overflow-hidden my-6">
    <div className="bg-terminal px-4 py-2 flex items-center gap-2 border-b border-border">
      <div className="flex gap-1.5">
        <span className="w-3 h-3 rounded-full bg-red-500/80" />
        <span className="w-3 h-3 rounded-full bg-yellow-500/80" />
        <span className="w-3 h-3 rounded-full bg-green-500/80" />
      </div>
      <span className="text-xs font-mono text-muted-foreground ml-2">gitmap diff ./gitmap ./gitmap-v28</span>
    </div>
    <div className="bg-terminal p-4 font-mono text-xs leading-relaxed overflow-x-auto">
      <div className="text-primary font-bold">  Conflicts (different content on both sides):</div>
      <div className="text-terminal-foreground">    README.md  (L: 4.2 KB @ 2026-04-17 14:02 | R: 4.1 KB @ 2026-04-18 09:11)</div>
      <div className="text-terminal-foreground">    src/app.ts  (L: 2.0 KB @ 2026-04-16 11:00 | R: 2.3 KB @ 2026-04-18 09:55)</div>
      <div className="text-primary font-bold mt-2">  Missing on RIGHT (would be added by merge-right / merge-both):</div>
      <div className="text-terminal-foreground">    docs/changelog.md  (L: 1.1 KB @ 2026-04-15 08:30)</div>
      <div className="text-primary font-bold mt-2">  Missing on LEFT (would be added by merge-left / merge-both):</div>
      <div className="text-terminal-foreground">    scripts/build.sh  (R: 512 B @ 2026-04-17 22:45)</div>
      <div className="text-blue-400 mt-2">[diff] summary: 1 missing-on-left, 1 missing-on-right, 2 conflicts, 137 identical</div>
    </div>
  </div>
);

const DiffPage = () => (
  <DocsLayout>
    <div className="max-w-4xl space-y-10">
      <div>
        <div className="flex items-center gap-3 mb-2">
          <GitCompare className="h-8 w-8 text-primary" />
          <h1 className="text-3xl font-bold tracking-tight">diff</h1>
          <span className="font-mono text-xs px-2 py-1 rounded bg-primary/10 text-foreground border border-primary/20 dark:bg-primary/15 dark:text-primary dark:border-primary/40 transition-colors duration-300 hover:border-primary/40 hover:shadow-sm hover:shadow-primary/10">alias: df</span>
        </div>
        <p className="text-lg text-muted-foreground">
          Read-only preview of what <code className="text-primary">merge-both</code>,{" "}
          <code className="text-primary">merge-left</code>, or <code className="text-primary">merge-right</code> would change between two folders.
          Lists files present on only one side and files whose content differs on both sides. Writes nothing, commits nothing, pushes nothing.
        </p>
        <p className="text-xs text-muted-foreground mt-2">
          Spec: companion to <code>spec/01-app/97-move-and-merge.md</code>
        </p>
      </div>

      <section>
        <h2 className="text-xl font-semibold mb-3 flex items-center gap-2">
          <FileSearch className="h-5 w-5 text-primary" /> Overview
        </h2>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          {[
            { icon: ShieldCheck, title: "Side-effect free", desc: "Never writes, commits, or pushes. URL endpoints are rejected — clone first." },
            { icon: FileSearch, title: "Three categories", desc: "Conflicts, missing-on-LEFT, missing-on-RIGHT — exactly what merge-* will act on." },
            { icon: Terminal, title: "JSON or text", desc: "Default human-readable output, --json for machine consumption." },
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
        <CodeBlock code={`gitmap diff LEFT RIGHT [flags]
gitmap df   LEFT RIGHT [flags]`} />
        <p className="text-sm text-muted-foreground mt-2">
          LEFT and RIGHT must both be local folder paths. URL endpoints are intentionally rejected — clone them first with{" "}
          <code className="text-primary">gitmap clone</code> so <code className="text-primary">diff</code> stays strictly side-effect-free.
        </p>
      </section>

      <section>
        <h2 className="text-xl font-semibold mb-3">Sample Output</h2>
        <TerminalPreview />
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
        <CodeBlock code={`# Plain diff between two local folders
gitmap diff ./gitmap ./gitmap-v28

# Conflicts only (preview before merge-both)
gitmap diff ./gitmap ./gitmap --only-conflicts

# Machine-readable output
gitmap df ./gitmap ./gitmap --json`} />
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
        <h2 className="text-xl font-semibold mb-3">Notes</h2>
        <ul className="list-disc list-inside space-y-2 text-muted-foreground text-sm">
          <li>The recommended dry-run preview before <code className="text-primary">gitmap merge-both</code> — every conflict listed here will trigger the <code>[L]eft / [R]ight / [S]kip / [A]ll-left / [B]all-right / [Q]uit</code> prompt during merge-both.</li>
          <li>The same default ignore list as <code>merge-*</code> applies: <code>.git/</code>, <code>node_modules/</code>, and <code>.gitmap/release-assets/</code> are skipped unless the corresponding <code>--include-*</code> flag is set.</li>
          <li>Identical files are tallied in the summary but not listed by default (use <code>--include-identical</code> to dump them).</li>
        </ul>
      </section>

      <section>
        <h2 className="text-xl font-semibold mb-3">See Also</h2>
        <ul className="list-disc list-inside space-y-1 text-sm">
          <li><a href="/merge-both" className="text-primary hover:underline">merge-both</a> — Apply a two-way merge after previewing</li>
          <li><a href="/merge-left" className="text-primary hover:underline">merge-left</a> — Apply RIGHT's changes into LEFT</li>
          <li><a href="/merge-right" className="text-primary hover:underline">merge-right</a> — Apply LEFT's changes into RIGHT</li>
          <li><a href="/mv" className="text-primary hover:underline">mv</a> — Move LEFT into RIGHT and delete LEFT</li>
        </ul>
      </section>
    </div>
  </DocsLayout>
);

export default DiffPage;
