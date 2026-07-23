import DocsLayout from "@/components/docs/DocsLayout";
import CodeBlock from "@/components/docs/CodeBlock";
import { GitBranch, Wrench, Globe } from "lucide-react";

const exitCodes = [
  { code: "0", meaning: "Success" },
  { code: "6", meaning: "Bad flag (e.g. missing URL)" },
  { code: "9", meaning: "chdir into cloned folder failed" },
  { code: "10", meaning: "Chained step (clone or fix-repo) failed — its exit code is propagated as-is" },
];

const CloneFixRepoPage = () => (
  <DocsLayout>
    <div className="max-w-4xl space-y-10">
      <div>
        <div className="flex items-center gap-3 mb-2">
          <GitBranch className="h-8 w-8 text-primary" />
          <h1 className="text-3xl font-bold tracking-tight">clone-fix-repo</h1>
          <span className="font-mono text-xs px-2 py-1 rounded bg-primary/10 text-foreground border border-primary/20 dark:bg-primary/15 dark:text-primary dark:border-primary/40">
            alias: cfr
          </span>
        </div>
        <p className="text-lg text-muted-foreground">
          Clone a repository, then immediately run <code>fix-repo --all</code> inside the new
          folder. One-shot replacement for the manual sequence: <code>gitmap clone</code>{" "}
          → <code>cd</code> → <code>gitmap fix-repo --all</code>.
        </p>
      </div>

      <section>
        <h2 className="text-xl font-semibold mb-3">Overview</h2>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          {[
            { icon: GitBranch, title: "Clones first", desc: "Versioned URLs auto-flatten (e.g. myrepo-v13 -> myrepo/). Optional folder name." },
            { icon: Wrench, title: "Then fix-repo --all", desc: "Re-execs the same gitmap binary with fix-repo --all in the cloned folder." },
            { icon: Globe, title: "Public variant", desc: "Use clone-fix-repo-pub (cfrp) to also run make-public --yes at the end." },
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
        <CodeBlock code={`gitmap clone-fix-repo <url> [folder]
gitmap cfr <url> [folder]                        # short alias

# Public variant (clone + fix-repo --all + make-public --yes)
gitmap clone-fix-repo-pub <url> [folder]
gitmap cfrp <url> [folder]`} />
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
              <tr className="border-t border-border">
                <td className="px-4 py-2 font-mono text-primary">--no-vscode-sync</td>
                <td className="px-4 py-2 font-mono text-muted-foreground">false</td>
                <td className="px-4 py-2 text-muted-foreground">
                  Forwarded to <code>clone</code>; skips writing the folder into VS Code Project
                  Manager <code>projects.json</code>. <code>fix-repo --all</code> is unaffected.
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>

      <section>
        <h2 className="text-xl font-semibold mb-3">Examples</h2>
        <CodeBlock code={`# HTTPS clone + fix
gitmap clone-fix-repo https://github.com/acme/myrepo-v13.git

# SSH clone with explicit folder name
gitmap cfr git@github.com:acme/myrepo-v13.git myrepo-fresh

# Public-publish in one go (clone + fix + make-public)
gitmap cfrp https://github.com/acme/myrepo-v13.git`} />
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
          <li><a href="/clone-command" className="text-primary hover:underline">clone</a> — The underlying clone step</li>
          <li><a href="/fix-repo" className="text-primary hover:underline">fix-repo</a> — The underlying rewrite step</li>
          <li><a href="/make-public" className="text-primary hover:underline">make-public</a> — Used by the cfrp variant</li>
        </ul>
      </section>
    </div>
  </DocsLayout>
);

export default CloneFixRepoPage;
