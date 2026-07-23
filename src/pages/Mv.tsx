import DocsLayout from "@/components/docs/DocsLayout";
import CodeBlock from "@/components/docs/CodeBlock";
import { Move, FolderInput, Trash2, GitBranch } from "lucide-react";

const flags = [
  { flag: "--no-push", def: "false", desc: "Skip git push on URL endpoints (still commits)" },
  { flag: "--no-commit", def: "false", desc: "Skip both commit and push on URL endpoints" },
  { flag: "--force-folder", def: "false", desc: "Replace folder whose origin doesn't match the URL" },
  { flag: "--pull", def: "false", desc: "Force git pull --ff-only on a folder endpoint" },
  { flag: "--init", def: "false", desc: "When RIGHT is auto-created, also git init it" },
  { flag: "--dry-run", def: "false", desc: "Print every action; perform none" },
  { flag: "--include-vcs", def: "false", desc: "Include .git/ in the copy (default: excluded)" },
  { flag: "--include-node-modules", def: "false", desc: "Include node_modules/ in the copy" },
];

const exitCodes = [
  { code: "0", meaning: "Success" },
  { code: "1", meaning: "Resolution, copy, commit, or push failed (message on stderr)" },
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
      <span className="text-xs font-mono text-muted-foreground ml-2">gitmap mv ./gitmap https://github.com/owner/gitmap-v28</span>
    </div>
    <div className="bg-terminal p-4 font-mono text-xs leading-relaxed overflow-x-auto">
      <div className="text-terminal-foreground">[mv] resolving RIGHT : https://github.com/owner/gitmap-v28</div>
      <div className="text-terminal-foreground">[mv]   -&gt; mapped to working folder: /work/gitmap-v28</div>
      <div className="text-terminal-foreground">[mv]   -&gt; folder does not exist; cloning</div>
      <div className="text-green-400">[mv]   -&gt; clone OK</div>
      <div className="text-terminal-foreground">[mv] copying files LEFT -&gt; RIGHT (excluding .git/) ...</div>
      <div className="text-green-400">[mv]   copied 142 files</div>
      <div className="text-terminal-foreground">[mv] committing in https://github.com/owner/gitmap ...</div>
      <div className="text-terminal-foreground">[mv]   commit a1b2c3d "gitmap mv from ./gitmap-v28"</div>
      <div className="text-terminal-foreground">[mv] pushing https://github.com/owner/gitmap ...</div>
      <div className="text-green-400">[mv]   push OK</div>
      <div className="text-blue-400 mt-1">[mv] done</div>
    </div>
  </div>
);

const MvPage = () => (
  <DocsLayout>
    <div className="max-w-4xl space-y-10">
      <div>
        <div className="flex items-center gap-3 mb-2">
          <Move className="h-8 w-8 text-primary" />
          <h1 className="text-3xl font-bold tracking-tight">mv</h1>
          <span className="font-mono text-xs px-2 py-1 rounded bg-primary/10 text-foreground border border-primary/20 dark:bg-primary/15 dark:text-primary dark:border-primary/40 transition-colors duration-300 hover:border-primary/40 hover:shadow-sm hover:shadow-primary/10">alias: move</span>
        </div>
        <p className="text-lg text-muted-foreground">
          Move every file from LEFT into RIGHT, then delete LEFT entirely. Either endpoint can be a local folder OR a remote git URL with an
          optional <code className="text-primary">:branch</code> suffix. URL endpoints are auto-cloned (or re-pulled) and a commit + push is made on the URL side after the file copy.
        </p>
        <p className="text-xs text-muted-foreground mt-2">
          Spec: <code>spec/01-app/97-move-and-merge.md</code>
        </p>
      </div>

      <section>
        <h2 className="text-xl font-semibold mb-3">Overview</h2>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          {[
            { icon: FolderInput, title: "Folder or URL", desc: "Either side accepts a local path or remote git URL with optional :branch." },
            { icon: GitBranch, title: "Auto clone + push", desc: "URL endpoints are cloned on demand and committed + pushed after the copy." },
            { icon: Trash2, title: "Destructive on LEFT", desc: "After a successful copy, LEFT is deleted. Use --dry-run to preview." },
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
        <CodeBlock code={`gitmap mv   LEFT RIGHT [flags]
gitmap move LEFT RIGHT [flags]`} />
        <p className="text-sm text-muted-foreground mt-2">
          LEFT and RIGHT can each be a local folder path (relative or absolute) or a remote git URL with an optional <code>:branch</code> suffix
          (e.g. <code>https://github.com/owner/repo:develop</code>).
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
        <CodeBlock code={`# Move one local folder into another
gitmap mv ./gitmap ./gitmap-v28

# Move a local folder into a remote repo (clone + push)
gitmap mv ./gitmap https://github.com/owner/gitmap-v28

# Preview without writing anything
gitmap mv ./gitmap ./gitmap --dry-run`} />
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
          <li>The <code>.git/</code> folder is never copied; LEFT's <code>.git/</code> is removed along with the rest of LEFT after the copy.</li>
          <li>LEFT and RIGHT must not resolve to the same folder, and one must not be nested inside the other — checked before any write.</li>
          <li>On a URL endpoint, the commit message is <code>gitmap mv from &lt;LEFT&gt;</code>.</li>
        </ul>
      </section>

      <section>
        <h2 className="text-xl font-semibold mb-3">See Also</h2>
        <ul className="list-disc list-inside space-y-1 text-sm">
          <li><a href="/merge-both" className="text-primary hover:underline">merge-both</a> — Two-way file-level merge</li>
          <li><a href="/merge-left" className="text-primary hover:underline">merge-left</a> — Merge into LEFT only</li>
          <li><a href="/merge-right" className="text-primary hover:underline">merge-right</a> — Merge into RIGHT only</li>
          <li><a href="/diff" className="text-primary hover:underline">diff</a> — Preview tree differences before moving</li>
        </ul>
      </section>
    </div>
  </DocsLayout>
);

export default MvPage;
