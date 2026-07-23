import DocsLayout from "@/components/docs/DocsLayout";
import CodeBlock from "@/components/docs/CodeBlock";
import { GitCommit, AlertTriangle, Clock, Filter, ListOrdered } from "lucide-react";
import { ReactNode } from "react";

export type Direction = "left" | "right" | "both";

interface CommitTransferPageProps {
  direction: Direction;
}

const meta: Record<Direction, {
  cmd: string;
  alias: string;
  short: string;
  arrow: string;
  writes: string;
  intro: ReactNode;
  example: string;
  output: string[];
}> = {
  left: {
    cmd: "commit-left",
    alias: "cl",
    short: "Replay RIGHT's commits onto LEFT",
    arrow: "RIGHT → LEFT",
    writes: "LEFT",
    intro: (
      <>
        <code className="text-primary">commit-left</code> replays every RIGHT-only commit since the merge-base
        onto LEFT as a sequence of fresh commits. Mirror of <code className="text-primary">commit-right</code>;
        useful when LEFT is your "clean" target and you want to import RIGHT's evolution chronologically.
      </>
    ),
    example: `# Replay RIGHT's commits since divergence onto LEFT
gitmap commit-left ./gitmap ./gitmap-v27-experimental

# Bypass the preview prompt
gitmap cl ./mine ./theirs -y`,
    output: [
      "[commit-left] replaying 4 commits from RIGHT onto LEFT:",
      "  [1/4] a3f2c1d  feat: add OAuth flow",
      "  [2/4] b7e4a9f  fix: handle expired tokens",
      "  [3/4] c2d8e1a  refactor: extract token store",
      "  [4/4] e5fa12b  docs: update README",
      "[commit-left] proceed? [y/N]",
    ],
  },
  right: {
    cmd: "commit-right",
    alias: "cr",
    short: "Replay LEFT's commits onto RIGHT",
    arrow: "LEFT → RIGHT",
    writes: "RIGHT",
    intro: (
      <>
        <code className="text-primary">commit-right</code> replays every LEFT-only commit since the merge-base
        onto RIGHT as a sequence of fresh commits. Phase-1 target of the family — see spec 106 §18 for the
        rollout order.
      </>
    ),
    example: `# Replay LEFT's commits onto a remote target
gitmap commit-right ./local https://github.com/owner/repo

# Re-run after adding a few new LEFT commits — only the genuinely-new commits replay
gitmap cr ./local https://github.com/owner/repo`,
    output: [
      "[commit-right] replaying 7 commits from LEFT onto RIGHT:",
      "  [1/7] a3f2c1d  feat: add OAuth flow",
      "  [2/7] b7e4a9f  fix: handle expired tokens",
      "  ...",
      "[commit-right] proceed? [y/N]",
    ],
  },
  both: {
    cmd: "commit-both",
    alias: "cb",
    short: "Interleave both sides' commits by author date",
    arrow: "bidirectional",
    writes: "LEFT and RIGHT",
    intro: (
      <>
        <code className="text-primary">commit-both</code> computes LEFT-only and RIGHT-only commits independently,
        sorts the union by <strong>author date ascending</strong> (LEFT first on ties), and replays each
        commit onto the opposite side. After the run both sides share the same chronological union.
      </>
    ),
    example: `# Interleave LEFT-only and RIGHT-only commits onto both sides
gitmap commit-both ./repo-a ./repo-b

# Same, with author-date-based interleave already implied
gitmap cb ./fork-a ./fork-b -y`,
    output: [
      "[commit-both] LEFT-only: 3 commits, RIGHT-only: 2 commits",
      "[commit-both] interleave (by author date):",
      "  [1/5] 2026-04-15 11:02  L→R  feat: add OAuth flow",
      "  [2/5] 2026-04-15 14:30  R→L  fix: handle expired tokens",
      "  [3/5] 2026-04-16 09:11  L→R  refactor: extract token store",
      "  [4/5] 2026-04-17 22:45  L→R  docs: update README",
      "  [5/5] 2026-04-18 09:55  R→L  test: add token-store tests",
      "[commit-both] proceed? [y/N]",
    ],
  },
};

const sharedFlags = [
  { flag: "-y, --yes", def: "false", desc: "Skip the up-front replay-set preview prompt" },
  { flag: "--include-merges", def: "false", desc: "Include merge commits in the replay set (default: filtered out)" },
  { flag: "--limit N", def: "0", desc: "Cap the replay set at N commits (oldest-first); 0 = unlimited" },
  { flag: "--since DATE", def: "''", desc: "Only replay commits authored on or after DATE (RFC3339 or git-style)" },
  { flag: "--strip / --no-strip", def: "true", desc: "Apply commit-message strip rules from config (prefix/suffix removal)" },
  { flag: "--drop / --no-drop", def: "true", desc: "Apply whole-commit drop filter (Merge branch, Revert ', fixup!, ...)" },
  { flag: "--prefer-source / --prefer-target", def: "false", desc: "Conflict resolution direction (mirrors merge-* preferences)" },
  { flag: "--no-provenance", def: "false", desc: "Skip the gitmap-replay: footer that protects against double-replay" },
  { flag: "--mirror", def: "false", desc: "Delete target-only files between commits (true tree-mirror; destructive)" },
  { flag: "--dry-run", def: "false", desc: "Print every action; perform none" },
];

const CommitTransferPage = ({ direction }: CommitTransferPageProps) => {
  const m = meta[direction];

  return (
    <DocsLayout>
      <div className="max-w-4xl space-y-10">
        <div>
          <div className="flex items-center gap-3 mb-2 flex-wrap">
            <GitCommit className="h-8 w-8 text-primary" />
            <h1 className="text-3xl font-bold tracking-tight">{m.cmd}</h1>
            <span className="font-mono text-xs px-2 py-1 rounded bg-primary/10 text-foreground border border-primary/20 dark:bg-primary/15 dark:text-primary dark:border-primary/40 transition-colors duration-300 hover:border-primary/40 hover:shadow-sm hover:shadow-primary/10">
              alias: {m.alias}
            </span>
            <span className="font-mono text-xs px-2 py-1 rounded bg-destructive/10 text-foreground border border-destructive/30 transition-colors duration-300 hover:border-destructive/50 hover:shadow-sm hover:shadow-destructive/10 dark:bg-destructive/20 dark:text-destructive-foreground dark:border-destructive/50">
              {m.arrow}
            </span>
          </div>
          <p className="text-lg text-muted-foreground">{m.short}.</p>
          <p className="text-xs text-muted-foreground mt-2">
            Spec: <code>spec/01-app/106-commit-left-right-both.md</code>
          </p>
        </div>

        {/* PLANNED banner — high-contrast in both themes; body text uses
            foreground (not muted) so it stays legible on the tinted bg. */}
        <div className="rounded-lg border-2 border-destructive/40 bg-destructive/5 dark:bg-destructive/10 p-4 flex items-start gap-3">
          <AlertTriangle className="h-5 w-5 text-destructive dark:text-[hsl(0_85%_72%)] flex-shrink-0 mt-0.5" />
          <div className="space-y-1.5">
            <p className="font-semibold text-sm text-foreground">
              Status: PLANNED — implementation deferred
            </p>
            <p className="text-sm text-foreground/90 dark:text-foreground/95 leading-relaxed">
              This command is fully specified but not yet shipped in the CLI. The page documents the
              contract so it can be reviewed before code lands. Track progress in spec 106 §18.
            </p>
            <p className="text-xs text-foreground/80 dark:text-foreground/85 leading-relaxed">
              Phasing: <strong className="text-foreground">Phase 1</strong> ships{" "}
              <code className="text-primary">commit-right</code> end-to-end first;{" "}
              <code className="text-primary">commit-left</code> is a wiring flip on top;{" "}
              <code className="text-primary">commit-both</code> adds the interleave-by-timestamp
              planner last.
            </p>
          </div>
        </div>

        <section>
          <h2 className="text-xl font-semibold mb-3">Purpose</h2>
          <p className="text-sm text-muted-foreground mb-3">
            The <code className="text-primary">merge-*</code> family transfers <strong>file state</strong> between two
            repo endpoints. The <code className="text-primary">commit-*</code> family transfers{" "}
            <strong>the history of how that file state was reached</strong> — it replays the source side's
            commit timeline onto the target side as a sequence of fresh commits, preserving order,
            authorship intent, and (cleaned) commit messages.
          </p>
          <p className="text-sm text-muted-foreground">{m.intro}</p>
        </section>

        <section>
          <h2 className="text-xl font-semibold mb-3">Mechanism (planned)</h2>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            {[
              { icon: ListOrdered, title: "Replay set", desc: "git rev-list --reverse --no-merges <merge-base>..<source-HEAD> — oldest-first." },
              { icon: Clock, title: "Manual reconstruct", desc: "Per source commit: checkout, file-snapshot copy, git add -A && commit on target. No cherry-pick." },
              { icon: Filter, title: "Message pipeline", desc: "Drop filter → strip rules → strip footers → trim → optional Conventional-Commits prefix." },
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
          <h2 className="text-xl font-semibold mb-3">Usage (planned)</h2>
          <CodeBlock code={`gitmap ${m.cmd} LEFT RIGHT [flags]
gitmap ${m.alias}            LEFT RIGHT [flags]`} />
          <p className="text-sm text-muted-foreground mt-2">
            LEFT and RIGHT use the same endpoint syntax as the <code>merge-*</code> family — local folder
            paths or <code>https://</code> / <code>git@</code> URLs with optional <code>:branch</code>{" "}
            suffix. Resolution rules are inherited verbatim from <code>movemerge.Endpoint</code>. Writes commits
            to: <strong>{m.writes}</strong>.
          </p>
        </section>

        <section>
          <h2 className="text-xl font-semibold mb-3">Sample preview prompt</h2>
          <div className="rounded-lg border border-border overflow-hidden">
            <div className="bg-terminal px-4 py-2 flex items-center gap-2 border-b border-border">
              <div className="flex gap-1.5">
                <span className="w-3 h-3 rounded-full bg-red-500/80" />
                <span className="w-3 h-3 rounded-full bg-yellow-500/80" />
                <span className="w-3 h-3 rounded-full bg-green-500/80" />
              </div>
              <span className="text-xs font-mono text-muted-foreground ml-2">gitmap {m.cmd} LEFT RIGHT</span>
            </div>
            <div className="bg-terminal p-4 font-mono text-xs leading-relaxed overflow-x-auto">
              {m.output.map((line, i) => (
                <div key={i} className={i === m.output.length - 1 ? "text-blue-400 mt-1" : "text-terminal-foreground"}>
                  {line}
                </div>
              ))}
            </div>
          </div>
        </section>

        <section>
          <h2 className="text-xl font-semibold mb-3">Flags (planned)</h2>
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
                {sharedFlags.map((f) => (
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
          <h2 className="text-xl font-semibold mb-3">Examples (planned)</h2>
          <CodeBlock code={m.example} />
        </section>

        <section>
          <h2 className="text-xl font-semibold mb-3">Caveats</h2>
          <ul className="list-disc list-inside space-y-2 text-muted-foreground text-sm">
            <li>
              Because replay uses <strong>manual file-state reconstruction</strong> rather than git-native
              cherry-pick, the target tree at each step will not byte-match the source commit's tree when
              the target carries files the source never had. This is intentional — the goal is "the same
              human-readable evolution," not a tree-hash-equivalent mirror. Pass <code>--mirror</code> to
              get strict tree equivalence (destructive on target-only files).
            </li>
            <li>
              Re-runs are protected by a <code>gitmap-replay:</code> footer on every replayed commit so the
              same source SHA is never replayed twice onto the same target.
            </li>
            <li>
              Naming mirror: <code>commit-left</code> writes <strong>to LEFT</strong> (source = RIGHT),
              exactly like <code>merge-left</code> writes files to LEFT. The "-left" suffix always names
              the destination.
            </li>
          </ul>
        </section>

        <section>
          <h2 className="text-xl font-semibold mb-3">See Also</h2>
          <ul className="list-disc list-inside space-y-1 text-sm">
            {direction !== "left" && (
              <li><a href="/commit-left" className="text-primary hover:underline">commit-left</a> — Mirror direction (writes to LEFT)</li>
            )}
            {direction !== "right" && (
              <li><a href="/commit-right" className="text-primary hover:underline">commit-right</a> — Phase-1 target (writes to RIGHT)</li>
            )}
            {direction !== "both" && (
              <li><a href="/commit-both" className="text-primary hover:underline">commit-both</a> — Bidirectional interleave by author date</li>
            )}
            <li><a href="/merge-both" className="text-primary hover:underline">merge-both</a> — File-state companion (transfers files, not history)</li>
            <li><a href="/diff" className="text-primary hover:underline">diff</a> — Read-only preview of file-state differences</li>
          </ul>
        </section>
      </div>
    </DocsLayout>
  );
};

export default CommitTransferPage;
