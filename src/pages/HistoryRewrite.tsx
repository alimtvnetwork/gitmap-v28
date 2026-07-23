import DocsLayout from "@/components/docs/DocsLayout";
import CodeBlock from "@/components/docs/CodeBlock";
import { History, ShieldCheck, Pin, Trash2 } from "lucide-react";

const flags = [
  { flag: "-y, --yes", def: "false", desc: "Skip the push confirmation prompt; force-push immediately on success" },
  { flag: "--no-push", def: "false", desc: "Stop after verification; print the manual git push command" },
  { flag: "--dry-run", def: "false", desc: "Run rewrite + verification in the sandbox, then exit without pushing" },
  { flag: "--message <s>", def: '""', desc: "Rewrite the message of every touched commit to this string" },
  { flag: "--keep-sandbox", def: "false", desc: "Don't delete the temp mirror-clone on exit" },
  { flag: "-q, --quiet", def: "false", desc: "Suppress per-phase progress lines" },
];

const exitCodes = [
  { code: "0", meaning: "Success" },
  { code: "2", meaning: "Not in a git repo, or origin remote missing" },
  { code: "3", meaning: "git filter-repo not installed" },
  { code: "4", meaning: "Bad args (zero paths, conflicting flags)" },
  { code: "5", meaning: "filter-repo returned non-zero" },
  { code: "6", meaning: "Verification disagreed with the requested operation" },
  { code: "7", meaning: "Push failed" },
];

const HistoryRewritePage = () => (
  <DocsLayout>
    <div className="max-w-4xl space-y-10">
      <div>
        <div className="flex items-center gap-3 mb-2">
          <History className="h-8 w-8 text-primary" />
          <h1 className="text-3xl font-bold tracking-tight">history-purge &amp; history-pin</h1>
        </div>
        <p className="text-lg text-muted-foreground">
          Two commands that wrap <code>git filter-repo</code> in a mirror-clone sandbox so your
          working repository is <strong>never</strong> rewritten in place. Use{" "}
          <code>history-purge</code> to remove file(s) from all history; use <code>history-pin</code>
          to make every past commit appear to contain the present-day bytes of a file.
        </p>
        <p className="text-xs text-muted-foreground mt-2">
          Spec: <code>spec/04-generic-cli/16-history-rewrite.md</code> · Research:{" "}
          <code>spec/15-research/git-history-rewrite-remove-and-pin-file.md</code>
        </p>
      </div>

      <section>
        <h2 className="text-xl font-semibold mb-3">Two commands, one pipeline</h2>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          {[
            {
              icon: Trash2,
              title: "history-purge (hp)",
              desc: "Remove file(s) and folder(s) from every commit on every branch. Useful for leaked secrets and bloat.",
            },
            {
              icon: Pin,
              title: "history-pin (hpin)",
              desc: "Pin file(s) to current content across all history. Every past revision shows the present-day bytes.",
            },
            {
              icon: ShieldCheck,
              title: "Sandbox-safe",
              desc: "Both run in a temp mirror-clone. Verification gates the push prompt; --dry-run never pushes.",
            },
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
        <CodeBlock code={`# Remove leaked credentials from all history
gitmap history-purge .env secrets/api.key
gitmap hp            .env secrets/api.key            # short alias

# Multi-path: separate args, comma, or comma-space all work
gitmap hp "secret.env, build/cache.bin"
gitmap hp secret.env,build/cache.bin

# Pin a doc to its current content across every past commit
gitmap history-pin docs/README.md
gitmap hpin        docs/README.md

# Dry run (no push, sandbox kept on disk for inspection)
gitmap hp .env --dry-run --keep-sandbox

# Hide what was scrubbed by rewriting touched commit messages
gitmap hp .env --message "history cleanup" --yes`} />
      </section>

      <section>
        <h2 className="text-xl font-semibold mb-3">Pipeline (5 phases)</h2>
        <ol className="list-decimal list-inside space-y-2 text-sm text-muted-foreground">
          <li><strong className="text-foreground">Identify origin</strong> — read <code>git remote get-url origin</code> in cwd. Exit 2 if missing.</li>
          <li><strong className="text-foreground">Mirror-clone</strong> — <code>git clone --mirror</code> into <code>os.MkdirTemp</code> sandbox. Cleaned up on every exit path.</li>
          <li><strong className="text-foreground">filter-repo</strong> — purge: <code>--invert-paths --path P</code>; pin: <code>--blob-callback</code> with current bytes substituted for every historical blob SHA.</li>
          <li><strong className="text-foreground">Verify</strong> — purge: <code>git log --all -- P</code> must be empty. Pin: every <code>git show &lt;sha&gt;:P</code> must hash to the same SHA-256.</li>
          <li><strong className="text-foreground">Push</strong> — print a verification-passed banner (mode, path count, sandbox, remote, warning) and prompt <code>Type 'yes' to force-push…</code>. The user must type the literal token <code>yes</code> — any other input aborts. <code>--yes</code> skips the prompt; <code>--no-push</code> short-circuits and prints the manual command.</li>
        </ol>
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
        <h2 className="text-xl font-semibold mb-3">Dependency: git filter-repo</h2>
        <p className="text-sm text-muted-foreground mb-3">
          <code>filter-repo</code> is not bundled with Git. Install once per machine:
        </p>
        <CodeBlock code={`# Linux / Windows (any OS with Python)
pip install --user git-filter-repo

# macOS
brew install git-filter-repo

# Windows alternative
scoop install git-filter-repo`} />
        <p className="text-xs text-muted-foreground mt-2">
          Missing? gitmap exits <code>3</code> with the install hint.
        </p>
      </section>

      <section>
        <h2 className="text-xl font-semibold mb-3">Safety guarantees</h2>
        <ul className="list-disc list-inside space-y-1 text-sm text-muted-foreground">
          <li>Working repo is <strong>never</strong> rewritten — all mutation happens in a temp sandbox.</li>
          <li>Verification runs <strong>before</strong> the push prompt. A failed verification cannot push.</li>
          <li><code>--dry-run</code> short-circuits before push regardless of <code>--yes</code>.</li>
          <li><code>--no-push</code> and <code>--yes</code> are mutually exclusive (exit 4).</li>
          <li>Sandbox path is printed on every error path; <code>--keep-sandbox</code> preserves it permanently.</li>
        </ul>
      </section>

      <section>
        <h2 className="text-xl font-semibold mb-3">See Also</h2>
        <ul className="list-disc list-inside space-y-1 text-sm">
          <li><a href="/commands" className="text-primary hover:underline">All commands</a></li>
          <li><a href="/flags" className="text-primary hover:underline">Flag reference</a></li>
        </ul>
      </section>
    </div>
  </DocsLayout>
);

export default HistoryRewritePage;