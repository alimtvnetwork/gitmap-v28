import DocsLayout from "@/components/docs/DocsLayout";
import CodeBlock from "@/components/docs/CodeBlock";
import { GitCommit, Layers, Clock, ShieldCheck } from "lucide-react";
import {
  commitInFlags as flags,
  commitInExitCodes as exitCodes,
  commitInAutoInit as autoInit,
  commitInProfileJson as profileJson,
} from "./commitInData";
import CommitInExamples from "./CommitInExamples";
import CommitInWhatItCreates from "./CommitInWhatItCreates";

const CommitInPage = () => (
  <DocsLayout>
    <div className="max-w-4xl space-y-10">
      <div>
        <div className="flex items-center gap-3 mb-2">
          <GitCommit className="h-8 w-8 text-primary" />
          <h1 className="text-3xl font-bold tracking-tight">commit-in</h1>
          <span className="font-mono text-xs px-2 py-1 rounded bg-primary/10 text-foreground border border-primary/20 dark:bg-primary/15 dark:text-primary dark:border-primary/40">
            alias: cin
          </span>
        </div>
        <p className="text-lg text-muted-foreground">
          Walk one or more SOURCE git repos in author-date order and APPEND each commit
          (preserving BOTH <code>AuthorDate</code> and <code>CommitterDate</code>) into a
          TARGET repo. Useful for stitching together project history that lives across forks,
          archives, or versioned siblings into a single canonical timeline — without ever
          rewriting an existing commit.
        </p>
        <p className="text-xs text-muted-foreground mt-2">
          Spec: <code>spec/03-commit-in/</code>
        </p>
      </div>

      <section>
        <h2 className="text-xl font-semibold mb-3">Overview</h2>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          {[
            { icon: Clock, title: "Chronological replay", desc: "Inputs walked oldest -> newest by author date; both AuthorDate and CommitterDate preserved byte-for-byte." },
            { icon: Layers, title: "Multi-source", desc: "Comma-separated inputs, or use all / -N to pull every (or the latest N) versioned siblings." },
            { icon: ShieldCheck, title: "Idempotent", desc: "Dedupe via ShaMap means re-running never replays a commit twice across runs." },
          ].map((f) => (
            <div key={f.title} className="rounded-lg border border-border p-4 bg-card">
              <f.icon className="h-5 w-5 text-primary mb-2" />
              <h3 className="font-semibold text-sm mb-1">{f.title}</h3>
              <p className="text-xs text-muted-foreground">{f.desc}</p>
            </div>
          ))}
        </div>
      </section>

      <CommitInWhatItCreates />

      <section>
        <h2 className="text-xl font-semibold mb-3">Usage</h2>
        <CodeBlock code={`gitmap commit-in <source> <input1,input2,...> [flags]
gitmap cin       <source> all                    [flags]
gitmap cin       <source> -5                     [flags]`} />
        <p className="text-sm text-muted-foreground mt-3">
          <code>&lt;source&gt;</code> is the TARGET repo (the one receiving appended commits).
          Auto-init is fixed: URL → <code>git clone</code>; existing repo → reuse; existing
          non-repo folder → <code>git init</code> in place; missing path →{" "}
          <code>mkdir -p && git init</code>. No prompt, no flag.
        </p>
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
        <h2 className="text-xl font-semibold mb-3">How &lt;source&gt; auto-init works</h2>
        <p className="text-sm text-muted-foreground mb-3">
          You never have to <code>git init</code> first. <code>commit-in</code> resolves
          <code> &lt;source&gt;</code> through a fixed dispatch table — no prompts, no flags,
          no surprises:
        </p>
        <div className="overflow-x-auto">
          <table className="w-full text-sm border border-border rounded-lg">
            <thead>
              <tr className="bg-muted/50">
                <th className="text-left px-4 py-2 font-medium">If &lt;source&gt; is…</th>
                <th className="text-left px-4 py-2 font-medium">commit-in does…</th>
              </tr>
            </thead>
            <tbody>
              {autoInit.map((row) => (
                <tr key={row.when} className="border-t border-border">
                  <td className="px-4 py-2 text-muted-foreground">{row.when}</td>
                  <td className="px-4 py-2 font-mono text-xs">{row.then}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </section>

      <CommitInExamples />

      <section>
        <h2 className="text-xl font-semibold mb-3">Sample profile JSON</h2>
        <p className="text-sm text-muted-foreground mb-3">
          Drop this file at{" "}
          <code>.gitmap/commit-in/profiles/Default.json</code> (relative to your workspace
          root — the nearest ancestor containing <code>.gitmap/</code>) and load it with{" "}
          <code>--profile Default</code> or <code>--default</code>. Keys and enum values are
          <strong> PascalCase</strong>; the loader uses <em>strict</em> decoding, so unknown
          keys are an error. Edit anything you like — every field maps 1:1 to a CLI flag
          above.
        </p>
        <CodeBlock
          language="json"
          title=".gitmap/commit-in/profiles/Default.json"
          code={profileJson}
        />
        <p className="text-xs text-muted-foreground mt-3">
          <strong>Tip:</strong> let gitmap write the file for you the first time —
          <code> gitmap cin ./canonical all --save-profile Default --set-default</code> —
          then open the resulting JSON and tweak. Re-saving requires{" "}
          <code>--save-profile-overwrite</code>. Profiles bind by absolute symlink-resolved
          path, NOT by remote URL, so two clones of the same upstream can carry different
          policies.
        </p>
      </section>

      <section>
        <h2 className="text-xl font-semibold mb-3">Resume &amp; idempotency</h2>
        <p className="text-sm text-muted-foreground mb-3">
          <code>commit-in</code> never replays a source commit twice. Every successfully
          replayed source SHA is written to the <code>ShaMap</code> table; on every
          subsequent run the walker looks each source SHA up and SKIPs with reason{" "}
          <code>DuplicateSourceSha</code> when a row already exists. That makes broken /
          partial runs naturally resumable — just re-issue the same command and the
          chain picks up where it left off.
        </p>
        <p className="text-sm text-muted-foreground mb-3">
          For cross-machine continuity (or when the SQLite DB is lost), gitmap mirrors the
          per-input "last completed source SHA" into a sidecar JSON file at{" "}
          <code>.gitmap/commit-in/state.json</code>. Re-seed the DB from it with{" "}
          <code>gitmap cin &lt;source&gt; all</code> — the walker pre-loads every source
          SHA into <code>ShaMap</code> before the first replay, so the resume probe is a
          single SQLite hit per commit.
        </p>
        <CodeBlock
          language="json"
          title=".gitmap/commit-in/state.json"
          code={`{
  "SchemaVersion": 1,
  "Source": "/abs/path/to/canonical",
  "Inputs": {
    "myrepo-v1": {
      "LastCompletedSourceSha": "9f3a…",
      "LastCompletedAt": "2026-06-24T10:12:33Z",
      "CommitsReplayed": 312
    },
    "myrepo-v2": {
      "LastCompletedSourceSha": "ab12…",
      "LastCompletedAt": "2026-06-24T10:14:01Z",
      "CommitsReplayed": 87
    },
    "myrepo-v3": {
      "LastCompletedSourceSha": "",
      "LastCompletedAt": "",
      "CommitsReplayed": 0
    }
  }
}`}
        />
        <p className="text-xs text-muted-foreground mt-3">
          <strong>Recovery flow.</strong> If <code>v2</code> aborted mid-walk, re-running{" "}
          <code>gitmap cin ./canonical myrepo-v1,myrepo-v2,myrepo-v3</code> will: (1) seed
          every source SHA from each input into <code>ShaMap</code>, (2) skip all of{" "}
          <code>v1</code> and the first 87 commits of <code>v2</code> as{" "}
          <code>DuplicateSourceSha</code>, (3) resume replay from the next un-mapped SHA,
          (4) continue into <code>v3</code> without re-asking.
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
          <li><a href="/commit-left" className="text-primary hover:underline">commit-left</a> / <a href="/commit-right" className="text-primary hover:underline">commit-right</a> / <a href="/commit-both" className="text-primary hover:underline">commit-both</a></li>
          <li><a href="/merge-left" className="text-primary hover:underline">merge-left</a> / <a href="/merge-right" className="text-primary hover:underline">merge-right</a> / <a href="/merge-both" className="text-primary hover:underline">merge-both</a></li>
        </ul>
      </section>
    </div>
  </DocsLayout>
);

export default CommitInPage;
