import DocsLayout from "@/components/docs/DocsLayout";
import CodeBlock from "@/components/docs/CodeBlock";
import { GitCommit, Layers, Clock, ShieldCheck } from "lucide-react";
import { commitInProfileJson as profileJson } from "./commitInData";
import CommitInExamples from "./CommitInExamples";
import CommitInWhatItCreates from "./CommitInWhatItCreates";
import { CommitInFlagsTable } from "@/components/commitin/CommitInFlagsTable";
import { CommitInAutoInitTable } from "@/components/commitin/CommitInAutoInitTable";
import { CommitInSponsorTemplates } from "@/components/commitin/CommitInSponsorTemplates";
import { CommitInExitCodesTable } from "@/components/commitin/CommitInExitCodesTable";

const OVERVIEW_CARDS = [
  { icon: Clock, title: "Chronological replay", desc: "Inputs walked oldest -> newest by author date; dates preserved byte-for-byte." },
  { icon: Layers, title: "Multi-source", desc: "Comma-separated inputs, or use all / -N to pull every versioned sibling." },
  { icon: ShieldCheck, title: "Idempotent", desc: "Dedupe via ShaMap means re-running never replays a commit twice across runs." },
];

const CommitInPage = () => (
  <DocsLayout>
    <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 space-y-8">
      <div>
        <div className="flex items-center gap-3 mb-2">
          <GitCommit className="h-8 w-8 text-primary" />
          <h1 className="text-3xl font-bold tracking-[-0.02em]">commit-in</h1>
          <span className="font-mono text-xs px-2 py-1 rounded bg-primary/10 text-foreground border border-primary/20 dark:bg-primary/15 dark:text-primary dark:border-primary/40">
            alias: cin
          </span>
        </div>
        <p className="text-lg text-muted-foreground">
          Walk one or more SOURCE git repos in author-date order and APPEND each commit
          (preserving BOTH <code>AuthorDate</code> and <code>CommitterDate</code>) into a TARGET repo without rewriting history.
        </p>
      </div>

      <section>
        <h2 className="text-xl font-semibold mb-3">Overview</h2>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          {OVERVIEW_CARDS.map((f) => (
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
        <CodeBlock code={"gitmap commit-in <source> <input1,input2,...> [flags]\ngitmap cin       <source> all                    [flags]\ngitmap cin       <source> -5                     [flags]"} />
        <p className="text-sm text-muted-foreground mt-3">
          <code>&lt;source&gt;</code> is the TARGET repo receiving appended commits. Auto-init is fixed: clone, reuse, or init.
        </p>
      </section>

      <CommitInFlagsTable />
      <CommitInAutoInitTable />
      <CommitInExamples />

      <section>
        <h2 className="text-xl font-semibold mb-3">Sample profile JSON</h2>
        <CodeBlock language="json" title=".gitmap/commit-in/profiles/Default.json" code={profileJson} />
      </section>

      <section>
        <h2 className="text-xl font-semibold mb-3">Resume &amp; Idempotency</h2>
        <p className="text-sm text-muted-foreground">
          Every replayed SHA is written to <code>ShaMap</code> and mirrored into <code>.gitmap/commit-in/state.json</code> for seamless recovery.
        </p>
      </section>

      <CommitInExitCodesTable />
      <CommitInSponsorTemplates />

      <section>
        <h2 className="text-xl font-semibold mb-3">See Also</h2>
        <ul className="list-disc list-inside space-y-1 text-sm text-muted-foreground">
          <li><a href="/commit-left" className="text-primary hover:underline">commit-left</a> / <a href="/commit-right" className="text-primary hover:underline">commit-right</a></li>
          <li><a href="/merge-left" className="text-primary hover:underline">merge-left</a> / <a href="/merge-right" className="text-primary hover:underline">merge-right</a></li>
        </ul>
      </section>
    </div>
  </DocsLayout>
);

export default CommitInPage;
