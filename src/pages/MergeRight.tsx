import DocsLayout from "@/components/docs/DocsLayout";
import CodeBlock from "@/components/docs/CodeBlock";
import { GitMerge, ArrowRight, ShieldCheck } from "lucide-react";

const flags = [
  { flag: "-y, --yes, -a, --accept-all", def: "false", desc: "Bypass prompt; default is --prefer-left" },
  { flag: "--prefer-left", def: "false", desc: "LEFT always wins (overwrite RIGHT)" },
  { flag: "--prefer-right", def: "false", desc: "RIGHT always wins (skip LEFT's version)" },
  { flag: "--prefer-newer", def: "false", desc: "Newer mtime wins" },
  { flag: "--prefer-skip", def: "false", desc: "Skip every conflict" },
  { flag: "--no-push", def: "false", desc: "Skip git push on URL RIGHT" },
  { flag: "--no-commit", def: "false", desc: "Skip commit and push on URL RIGHT" },
  { flag: "--force-folder", def: "false", desc: "Replace folder whose origin doesn't match URL" },
  { flag: "--pull", def: "false", desc: "Force git pull --ff-only on a folder endpoint" },
  { flag: "--dry-run", def: "false", desc: "Print every action; perform none" },
  { flag: "--include-vcs", def: "false", desc: "Include .git/ in copy/diff" },
  { flag: "--include-node-modules", def: "false", desc: "Include node_modules/ in copy/diff" },
];

const exitCodes = [
  { code: "0", meaning: "Success" },
  { code: "1", meaning: "Resolution, copy, commit, push failed, or user pressed Q" },
  { code: "2", meaning: "Wrong number of positional arguments" },
];

const MergeRightPage = () => (
  <DocsLayout>
    <div className="max-w-4xl space-y-10">
      <div>
        <div className="flex items-center gap-3 mb-2">
          <GitMerge className="h-8 w-8 text-primary" />
          <h1 className="text-3xl font-bold tracking-tight">merge-right</h1>
          <span className="font-mono text-xs px-2 py-1 rounded bg-primary/10 text-foreground border border-primary/20 dark:bg-primary/15 dark:text-primary dark:border-primary/40 transition-colors duration-300 hover:border-primary/40 hover:shadow-sm hover:shadow-primary/10">alias: mr</span>
        </div>
        <p className="text-lg text-muted-foreground">
          One-way file-level merge that writes only into RIGHT. Files missing on RIGHT are copied from LEFT;
          conflicts are resolved into RIGHT. <strong>LEFT is never modified.</strong> If RIGHT originated from a URL it is committed + pushed after the merge.
        </p>
        <p className="text-xs text-muted-foreground mt-2">
          Spec: <code>spec/01-app/97-move-and-merge.md</code>
        </p>
      </div>

      <section>
        <h2 className="text-xl font-semibold mb-3">Overview</h2>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          {[
            { icon: ArrowRight, title: "Writes RIGHT only", desc: "LEFT is treated as read-only — no commit or push happens on LEFT." },
            { icon: ShieldCheck, title: "Default: prefer-left", desc: "With -y, LEFT wins on conflict (treat LEFT as the source being published)." },
            { icon: GitMerge, title: "URL-aware commits", desc: "If RIGHT is a URL, the merge is committed and pushed automatically." },
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
        <CodeBlock code={`gitmap merge-right LEFT RIGHT [flags]
gitmap mr          LEFT RIGHT [flags]`} />
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
        <CodeBlock code={`# Push LEFT's changes into a remote repo (LEFT wins by default)
gitmap merge-right ./local https://github.com/owner/repo -y

# Push to a specific branch
gitmap mr ./local https://github.com/owner/repo:develop -y

# Stage changes locally without pushing
gitmap merge-right ./local ./mirror --no-push`} />
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
          <li>LEFT is read-only for <code>merge-right</code>; no commit or push happens on LEFT even when it is a URL endpoint.</li>
          <li>With <code>-y</code>, the per-command default is <code>--prefer-left</code> (treat LEFT as the source of truth being published).</li>
        </ul>
      </section>

      <section>
        <h2 className="text-xl font-semibold mb-3">See Also</h2>
        <ul className="list-disc list-inside space-y-1 text-sm">
          <li><a href="/merge-left" className="text-primary hover:underline">merge-left</a> — Mirror operation: write into LEFT only</li>
          <li><a href="/merge-both" className="text-primary hover:underline">merge-both</a> — Two-way merge</li>
          <li><a href="/mv" className="text-primary hover:underline">mv</a> — Move LEFT into RIGHT and delete LEFT</li>
        </ul>
      </section>
    </div>
  </DocsLayout>
);

export default MergeRightPage;
