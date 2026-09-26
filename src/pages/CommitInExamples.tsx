import CodeBlock from "@/components/docs/CodeBlock";

// CommitInExamples renders the worked walkthroughs for the
// commit-in docs page, including declarative JSON configs, state
// templates, variables, and $files.2.names title synthesis.
const CommitInExamples = () => (
  <section>
    <h2 className="text-xl font-semibold mb-3">Examples</h2>

    <h3 className="font-semibold text-sm mt-4 mb-2 text-foreground">
      1 · Declarative 1-Line Migration with <code>--config</code>, Pre-Compiled Variables &amp; <code>$files.2.names</code>
    </h3>
    <p className="text-sm text-muted-foreground mb-2">
      Feed a single declarative JSON file to <code>gitmap commit-pull --config</code>. It
      auto-creates the target repository, imports state templates into{" "}
      <code>gitmap-templates.db</code> (skipping if the SHA-256 <code>exportId</code> is
      unchanged), pre-compiles all <code>$VAR</code> variables before the commit loop,
      strips matching lines (<code>starts_with</code>, <code>ends_with</code>,{" "}
      <code>contains</code>, <code>regex</code>), replaces generic <code>Changes</code>{" "}
      titles with <code>$files.2.names: $seo.title</code>, and enters <code>--cd</code>:
    </p>
    <CodeBlock
      language="json"
      code={`{
  "target": "D:\\\\test-gitmap\\\\test-gitmap",
  "inputs": [
    "https://github.com/alimtvnetwork/git-repo-navigator",
    "https://github.com/alimtvnetwork/gitmap-v{2..28}"
  ],
  "prMode": "merges",
  "tree": true,
  "finalSync": true,
  "cd": true,
  "imports": [".ai-memory/temp/seo-templates.json"],
  "lineSkippers": [
    { "mode": "starts_with", "pattern": "Co-authored-by:" },
    { "mode": "contains", "pattern": "X-Lovable-Edit-ID" }
  ],
  "titleReplacements": [
    { "matchMode": "equals", "match": "Changes", "replacement": "$files.2.names: $seo.title" }
  ],
  "suffixTemplates": ["seo"]
}`}
    />
    <CodeBlock
      language="bash"
      code={`# Run the entire migration in 1 line:
gitmap commit-pull --config .ai-memory/temp/commit-pull-config.json

# Manage templates & variables via CLI or Browser UI:
gitmap templates import .ai-memory/temp/seo-templates.json
gitmap templates export .ai-memory/temp/exported-templates.json --category seo
gitmap templates ui`}
    />

    <h3 className="font-semibold text-sm mt-6 mb-2 text-foreground">
      2 · Convert a plain folder of files into a git repo + replay history
    </h3>
    <p className="text-sm text-muted-foreground mb-2">
      You have <code>./my-project/</code> with code but no <code>.git/</code> yet.
      Point <code>commit-in</code> at it and pull history from a URL — the folder is
      auto-<code>git init</code>ed in place, your files stay where they are.
    </p>
    <CodeBlock
      language="bash"
      code={`gitmap commit-in ./my-project https://github.com/me/my-project-archive.git`}
    />

    <h3 className="font-semibold text-sm mt-6 mb-2 text-foreground">
      3 · Mix a local folder + a remote URL as INPUTS into one canonical timeline
    </h3>
    <CodeBlock
      language="bash"
      code={`gitmap cin ./canonical \\
    ./old-local-checkout,https://github.com/me/old-fork.git,git@github.com:me/new-fork.git`}
    />

    <h3 className="font-semibold text-sm mt-6 mb-2 text-foreground">
      4 · Replay every versioned sibling automatically (<code>all</code> / <code>-N</code>)
    </h3>
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
    --title-suffix " — via gitmap-v28"`}
    />

    <h3 className="font-semibold text-sm mt-6 mb-2 text-foreground">
      6 · Mirror tags + auto-create release branches
    </h3>
    <CodeBlock
      language="bash"
      code={`gitmap cin ./canonical https://github.com/me/legacy.git
gitmap cin ./canonical ./old --tags All
gitmap cin ./canonical ./old --no-release-branch`}
    />
  </section>
);

export default CommitInExamples;
