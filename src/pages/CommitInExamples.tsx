import CodeBlock from "@/components/docs/CodeBlock";

// CommitInExamples renders the seven worked walkthroughs for the
// commit-in docs page. Extracted from CommitIn.tsx so the page
// component stays under the project-wide <200-lines rule.
const CommitInExamples = () => (
  <section>
    <h2 className="text-xl font-semibold mb-3">Examples</h2>

    <h3 className="font-semibold text-sm mt-4 mb-2 text-foreground">
      1 · Convert a plain folder of files into a git repo + replay history
    </h3>
    <p className="text-sm text-muted-foreground mb-2">
      You have <code>./my-project/</code> with code but no <code>.git/</code> yet.
      Point <code>commit-in</code> at it and pull history from a URL — the folder is
      auto-<code>git init</code>ed in place, your files stay where they are.
    </p>
    <CodeBlock
      language="bash"
      code={`# folder exists, no .git/ yet — commit-in will run \`git init\` for you
gitmap commit-in ./my-project https://github.com/me/my-project-archive.git`}
    />

    <h3 className="font-semibold text-sm mt-6 mb-2 text-foreground">
      2 · Mix a local folder + a remote URL as INPUTS into one canonical timeline
    </h3>
    <p className="text-sm text-muted-foreground mb-2">
      The first positional is the TARGET. The second is the comma-separated INPUTS to
      walk in author-date order. You can freely mix a local checkout with one or more
      remote URLs — each URL is shallow-cloned into{" "}
      <code>.gitmap/temp/&lt;runId&gt;/</code> and walked just like the local one.
    </p>
    <CodeBlock
      language="bash"
      code={`# target = ./canonical (auto-init if missing)
# inputs = local folder + 2 remote forks, walked oldest -> newest
gitmap cin ./canonical \\
    ./old-local-checkout,https://github.com/me/old-fork.git,git@github.com:me/new-fork.git`}
    />

    <h3 className="font-semibold text-sm mt-6 mb-2 text-foreground">
      3 · Brand-new target folder from scratch (mkdir + init + replay)
    </h3>
    <p className="text-sm text-muted-foreground mb-2">
      Pass a path that does not exist. <code>commit-in</code> creates the folder, runs
      <code> git init</code>, and starts appending — one command, zero setup.
    </p>
    <CodeBlock
      language="bash"
      code={`gitmap commit-in ./brand-new-canonical \\
    https://github.com/me/legacy-v1.git,https://github.com/me/legacy-v2.git`}
    />

    <h3 className="font-semibold text-sm mt-6 mb-2 text-foreground">
      4 · Replay every versioned sibling automatically (<code>all</code> / <code>-N</code>)
    </h3>
    <p className="text-sm text-muted-foreground mb-2">
      The <code>all</code> keyword expands to every <code>&lt;target&gt;-vN</code>{" "}
      sibling sitting next to your TARGET on disk, walked oldest&nbsp;→&nbsp;newest.
      Use <code>-N</code> to take only the most recent N siblings.
    </p>
    <p className="text-sm text-muted-foreground mb-2">
      <strong>Concrete example.</strong> Say your TARGET is{" "}
      <code>./gitmap-v27</code> and the parent directory looks like this:
    </p>
    <CodeBlock
      language="bash"
      code={`$ ls ./
gitmap         <-- TARGET (the one receiving appended commits)
gitmap-v27
gitmap-v27
gitmap-v27
...
gitmap-v27
gitmap     <-- newest sibling`}
    />
    <p className="text-sm text-muted-foreground mb-2 mt-3">
      Then the keywords expand like this — <em>no manual list needed</em>:
    </p>
    <div className="overflow-x-auto mb-3">
      <table className="w-full text-sm border border-border rounded-lg">
        <thead>
          <tr className="bg-muted/50">
            <th className="text-left px-4 py-2 font-medium">You type…</th>
            <th className="text-left px-4 py-2 font-medium">commit-in walks…</th>
          </tr>
        </thead>
        <tbody>
          <tr className="border-t border-border">
            <td className="px-4 py-2 font-mono text-primary">gitmap cin ./gitmap all</td>
            <td className="px-4 py-2 font-mono text-xs text-muted-foreground">gitmap-v27, gitmap-v27, …, gitmap (every sibling, oldest first)</td>
          </tr>
          <tr className="border-t border-border">
            <td className="px-4 py-2 font-mono text-primary">gitmap cin ./gitmap -5</td>
            <td className="px-4 py-2 font-mono text-xs text-muted-foreground">gitmap-v27, gitmap-v27, gitmap-v27, gitmap-v27, gitmap-v27</td>
          </tr>
          <tr className="border-t border-border">
            <td className="px-4 py-2 font-mono text-primary">gitmap cin ./gitmap -1</td>
            <td className="px-4 py-2 font-mono text-xs text-muted-foreground">gitmap (just the newest)</td>
          </tr>
        </tbody>
      </table>
    </div>
    <p className="text-sm text-muted-foreground mb-2">
      And remember the TARGET (<code>./gitmap-v27</code>) doesn't have to exist
      yet — if it's missing, commit-in does <code>mkdir -p ./gitmap && git init</code>{" "}
      first, then appends every sibling's history into the new repo.
      One command takes you from <em>"I have 19 versioned snapshots"</em> to{" "}
      <em>"I have one git repo with 19 versions of history in author order."</em>
    </p>
    <CodeBlock
      language="bash"
      code={`# Every sibling, save the resolved settings as the default profile
gitmap commit-in ./gitmap all --save-profile Default --set-default

# Just the last 3 siblings, dry-run, with per-language new-function intel
gitmap cin ./gitmap -3 --dry-run --function-intel on --languages Go,TypeScript`}
    />

    <h3 className="font-semibold text-sm mt-6 mb-2 text-foreground">
      5 · Override author + scrub commit messages
    </h3>
    <CodeBlock
      language="bash"
      code={`gitmap cin git@github.com:me/canonical.git \\
    https://github.com/me/old-fork.git,https://github.com/me/new-fork.git \\
    --author-name "Jane Doe" --author-email jane@example.com \\
    --message-exclude "StartsWith:Signed-off-by:,Contains:[skip ci]" \\
    --title-suffix " — via gitmap-v27"`}
    />

    <h3 className="font-semibold text-sm mt-6 mb-2 text-foreground">
      6 · Reuse a saved profile + only rewrite weak titles
    </h3>
    <CodeBlock
      language="bash"
      code={`gitmap cin ./canonical all --default \\
    --override-messages "Refine implementation,Improve module" \\
    --override-only-weak`}
    />

    <h3 className="font-semibold text-sm mt-6 mb-2 text-foreground">
      7 · Headless CI run (fail loudly on any unset value)
    </h3>
    <CodeBlock language="bash" code={`gitmap cin ./canonical all --profile CI --no-prompt`} />

    <h3 className="font-semibold text-sm mt-6 mb-2 text-foreground">
      8 · Mirror tags + auto-create release branches (NEW)
    </h3>
    <p className="text-sm text-muted-foreground mb-2">
      When the source carries an annotated tag like <code>v1.2.3</code>,{" "}
      <code>commit-in</code> re-creates that tag on the destination — pointing at the{" "}
      NEW commit SHA produced by the replay (not the original source SHA, which doesn't
      exist in the destination history). If the tag matches the canonical semver shape,
      it ALSO creates a <code>release/&lt;tag&gt;</code> branch at the same new SHA, so
      it's interchangeable with <code>gitmap release-branch</code> tooling.
    </p>
    <CodeBlock
      language="bash"
      code={`# Default: annotated tags only, auto release branch ON
gitmap cin ./canonical https://github.com/me/legacy.git

# Mirror EVERY tag including lightweight bookmarks
gitmap cin ./canonical ./old --tags All

# Mirror tags but skip the auto release branch (you'll cut releases manually)
gitmap cin ./canonical ./old --no-release-branch

# Custom prefix — branches become 'releases/v1.2.3' instead of 'release/v1.2.3'
gitmap cin ./canonical ./old --release-branch-prefix releases/

# Disable tag mirroring entirely
gitmap cin ./canonical ./old --tags None`}
    />
    <p className="text-sm text-muted-foreground mt-3 mb-2">
      The three-way relationship (old SHA → new SHA → mirrored tag → release branch) is
      persisted in the <code>RewrittenCommit</code> SQLite row. Query it any time:
    </p>
    <CodeBlock
      language="sql"
      code={`SELECT sc.SourceSha       AS OldSha,
       rc.NewSha          AS NewSha,
       rc.MirroredTagName AS Tag,
       rc.MirroredReleaseBranch AS ReleaseBranch
FROM   RewrittenCommit rc
JOIN   SourceCommit    sc ON sc.SourceCommitId = rc.SourceCommitId
WHERE  rc.MirroredTagName IS NOT NULL
ORDER  BY rc.RewrittenCommitId;`}
    />
  </section>
);

export default CommitInExamples;
