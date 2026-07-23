import DocsLayout from "@/components/docs/DocsLayout";
import CodeBlock from "@/components/docs/CodeBlock";
import { GitMerge, ArrowLeftRight, Keyboard } from "lucide-react";

const flags = [
  { flag: "-y, --yes, -a, --accept-all", def: "false", desc: "Bypass prompt; default is --prefer-newer" },
  { flag: "--prefer-left", def: "false", desc: "LEFT always wins on conflict" },
  { flag: "--prefer-right", def: "false", desc: "RIGHT always wins on conflict" },
  { flag: "--prefer-newer", def: "false", desc: "Newer mtime wins on conflict" },
  { flag: "--prefer-skip", def: "false", desc: "Skip every conflict; only missing files copied" },
  { flag: "--no-push", def: "false", desc: "Skip git push on URL endpoints" },
  { flag: "--no-commit", def: "false", desc: "Skip commit and push on URL endpoints" },
  { flag: "--force-folder", def: "false", desc: "Replace folder whose origin doesn't match URL" },
  { flag: "--pull", def: "false", desc: "Force git pull --ff-only on a folder endpoint" },
  { flag: "--dry-run", def: "false", desc: "Print every action; perform none" },
  { flag: "--include-vcs", def: "false", desc: "Include .git/ in copy/diff" },
  { flag: "--include-node-modules", def: "false", desc: "Include node_modules/ in copy/diff" },
];

const promptKeys = [
  { key: "L", action: "Take LEFT's version (write into RIGHT)" },
  { key: "R", action: "Take RIGHT's version (write into LEFT)" },
  { key: "S", action: "Skip this file" },
  { key: "A", action: "Sticky: All-Left for the rest of the run" },
  { key: "B", action: "Sticky: All-Right for the rest of the run" },
  { key: "Q", action: "Quit immediately (already-applied changes are kept)" },
];

const exitCodes = [
  { code: "0", meaning: "Success" },
  { code: "1", meaning: "Resolution, copy, commit, push failed, or user pressed Q" },
  { code: "2", meaning: "Wrong number of positional arguments" },
];

const MergeBothPage = () => (
  <DocsLayout>
    <div className="max-w-4xl space-y-10">
      <div>
        <div className="flex items-center gap-3 mb-2">
          <GitMerge className="h-8 w-8 text-primary" />
          <h1 className="text-3xl font-bold tracking-tight">merge-both</h1>
          <span className="font-mono text-xs px-2 py-1 rounded bg-primary/10 text-foreground border border-primary/20 dark:bg-primary/15 dark:text-primary dark:border-primary/40 transition-colors duration-300 hover:border-primary/40 hover:shadow-sm hover:shadow-primary/10">alias: mb</span>
        </div>
        <p className="text-lg text-muted-foreground">
          Two-way file-level merge between LEFT and RIGHT. Files present on only one side are copied to the other; files present on both with
          different content trigger an interactive conflict prompt. Each side that originated from a URL is committed + pushed independently.
        </p>
        <p className="text-xs text-muted-foreground mt-2">
          Spec: <code>spec/01-app/97-move-and-merge.md</code>
        </p>
      </div>

      <section>
        <h2 className="text-xl font-semibold mb-3">Overview</h2>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          {[
            { icon: ArrowLeftRight, title: "Bidirectional", desc: "Both sides gain each other's missing files; conflicts are negotiated." },
            { icon: Keyboard, title: "Interactive prompt", desc: "Per-file [L]/[R]/[S]/[A]/[B]/[Q] choice with sticky modes." },
            { icon: GitMerge, title: "URL-aware commits", desc: "Each URL endpoint is committed and pushed independently after the merge." },
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
        <CodeBlock code={`gitmap merge-both LEFT RIGHT [flags]
gitmap mb         LEFT RIGHT [flags]`} />
        <p className="text-sm text-muted-foreground mt-2">
          LEFT and RIGHT can each be a folder path or a remote git URL (optionally suffixed with <code>:branch</code>).
        </p>
      </section>

      <section>
        <h2 className="text-xl font-semibold mb-3 flex items-center gap-2">
          <Keyboard className="h-5 w-5 text-primary" /> Conflict Prompt
        </h2>
        <div className="rounded-lg border border-border p-4 bg-muted/30 font-mono text-sm mb-4">
          [L]eft  [R]ight  [S]kip  [A]ll-left  [B]all-right  [Q]uit
        </div>
        <div className="overflow-x-auto">
          <table className="w-full text-sm border border-border rounded-lg">
            <thead>
              <tr className="bg-muted/50">
                <th className="text-left px-4 py-2 font-medium">Key</th>
                <th className="text-left px-4 py-2 font-medium">Action</th>
              </tr>
            </thead>
            <tbody>
              {promptKeys.map((p) => (
                <tr key={p.key} className="border-t border-border">
                  <td className="px-4 py-2 font-mono text-primary font-bold">{p.key}</td>
                  <td className="px-4 py-2 text-muted-foreground">{p.action}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
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
        <CodeBlock code={`# Interactive two-way merge between two local folders
gitmap merge-both ./gitmap ./gitmap-v27

# Non-interactive (newer wins by default) — commits + pushes the URL side
gitmap mb ./local https://github.com/owner/repo -y

# Preview a LEFT-wins merge without writing
gitmap merge-both ./a ./b -y --prefer-left --dry-run`} />
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
          <li><a href="/diff" className="text-primary hover:underline">diff</a> — Recommended dry-run preview before merge-both</li>
          <li><a href="/merge-left" className="text-primary hover:underline">merge-left</a> — One-way merge into LEFT only</li>
          <li><a href="/merge-right" className="text-primary hover:underline">merge-right</a> — One-way merge into RIGHT only</li>
          <li><a href="/mv" className="text-primary hover:underline">mv</a> — Move LEFT into RIGHT and delete LEFT</li>
        </ul>
      </section>
    </div>
  </DocsLayout>
);

export default MergeBothPage;
