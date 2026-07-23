import DocsLayout from "@/components/docs/DocsLayout";
import CodeBlock from "@/components/docs/CodeBlock";
import { Rocket, Tag, GitPullRequest, Database } from "lucide-react";

const aliasFlags = [
  { flag: "--force, -f", def: "false", desc: "Overwrite an existing alias that points to a different repo" },
];

const releaseAliasFlags = [
  { flag: "--pull", def: "false", desc: "Run git pull --ff-only inside the target repo before releasing" },
  { flag: "--no-stash", def: "false", desc: "Abort if the working tree is dirty (skip the auto-stash safety net)" },
  { flag: "--dry-run", def: "false", desc: "Forwarded to gitmap release — preview the pipeline without tagging" },
];

const ReleaseAliasPage = () => (
  <DocsLayout>
    <div className="max-w-4xl space-y-10">
      <div>
        <div className="flex items-center gap-3 mb-2">
          <Rocket className="h-8 w-8 text-primary" />
          <h1 className="text-3xl font-bold tracking-tight">as · release-alias · release-alias-pull</h1>
        </div>
        <p className="text-lg text-muted-foreground">
          Register a Git repo under a short name with <code>gitmap as</code>, then release it from
          anywhere on disk with <code>release-alias</code> (alias <code>ra</code>) or the
          pull-then-release shortcut <code>release-alias-pull</code> (alias <code>rap</code>).
          Auto-stash protects in-flight work; the stash label is{" "}
          <code>alias-version-unixts</code> so concurrent runs never pop each other&apos;s stash.
        </p>
        <p className="text-xs text-muted-foreground mt-2">
          Available since v3.0.0 · Spec: <code>spec/04-generic-cli/15-release-alias.md</code>
        </p>
      </div>

      <section>
        <h2 className="text-xl font-semibold mb-3">Three commands, one workflow</h2>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          {[
            {
              icon: Tag,
              title: "as (s-alias)",
              desc: "Register the current Git repo under a short name. Mirrors to VS Code Project Manager projects.json when present.",
            },
            {
              icon: Rocket,
              title: "release-alias (ra)",
              desc: "From any directory: chdir into the aliased repo, auto-stash, run the standard release pipeline, then pop the stash.",
            },
            {
              icon: GitPullRequest,
              title: "release-alias-pull (rap)",
              desc: "Same as ra but with --pull always implied. Pulls --ff-only first so you never tag on a divergent tree.",
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
        <h2 className="text-xl font-semibold mb-3">Typical workflow</h2>
        <CodeBlock code={`# 1. Inside the repo, register a short name (folder basename if omitted)
cd ~/work/my-api
gitmap as my-api

# 2. From anywhere on disk, release by alias
gitmap release-alias my-api v1.4.0
gitmap ra            my-api v1.4.0          # short alias

# 3. Or pull --ff-only first, then release in one shot
gitmap release-alias-pull my-api v1.4.0
gitmap rap                my-api v1.4.0     # short alias

# Preview without tagging
gitmap ra  my-api v1.4.0 --dry-run
gitmap rap my-api v1.4.0 --dry-run

# Strict mode: refuse to auto-stash a dirty tree
gitmap ra my-api v1.4.0 --no-stash`} />
      </section>

      <section>
        <h2 className="text-xl font-semibold mb-3">gitmap as — register an alias</h2>
        <CodeBlock code={`gitmap as [alias-name] [--force]

gitmap as                  # alias = folder basename
gitmap as backend          # explicit alias
gitmap as backend -f       # overwrite an existing 'backend' alias`} />
        <FlagTable flags={aliasFlags} />
      </section>

      <section>
        <h2 className="text-xl font-semibold mb-3">gitmap release-alias — flags</h2>
        <p className="text-sm text-muted-foreground mb-3">
          <code>release-alias-pull</code> accepts the same flags except <code>--pull</code>, which
          is always on.
        </p>
        <FlagTable flags={releaseAliasFlags} />
      </section>

      <section>
        <h2 className="text-xl font-semibold mb-3">Auto-stash safety net</h2>
        <ul className="list-disc list-inside space-y-1 text-sm text-muted-foreground">
          <li>If the working tree is dirty, gitmap stashes with the label <code>alias-version-unixts</code> before chdir + release.</li>
          <li>After the release pipeline returns, only the stash whose label matches is popped — concurrent <code>ra</code> calls in the same repo never clobber each other.</li>
          <li><code>--no-stash</code> opts out: a dirty tree aborts the release with a clear error instead of stashing.</li>
          <li>If the release pipeline fails, the labelled stash is left in place so you can recover with <code>git stash list</code> + <code>git stash apply</code>.</li>
        </ul>
      </section>

      <section>
        <h2 className="text-xl font-semibold mb-3">
          <Database className="h-5 w-5 inline mr-1 text-primary" />
          Where aliases live
        </h2>
        <p className="text-sm text-muted-foreground">
          Aliases are stored in the gitmap SQLite database next to the binary, in the{" "}
          <code>Alias</code> table. When VS Code Project Manager is detected, <code>gitmap as</code>{" "}
          also mirrors the alias into <code>projects.json</code> so the Project Manager sidebar and
          gitmap stay in sync. Run <code>gitmap alias list</code> to see every registered alias.
        </p>
      </section>

      <section>
        <h2 className="text-xl font-semibold mb-3">See also</h2>
        <ul className="list-disc list-inside space-y-1 text-sm">
          <li><a href="/alias" className="text-primary hover:underline">alias — list / inspect / remove registered aliases</a></li>
          <li><a href="/release" className="text-primary hover:underline">release — the underlying release pipeline</a></li>
          <li><a href="/cd" className="text-primary hover:underline">cd — jump to an aliased repo in your shell</a></li>
          <li><a href="/commands" className="text-primary hover:underline">All commands</a></li>
        </ul>
      </section>
    </div>
  </DocsLayout>
);

type FlagRow = { flag: string; def: string; desc: string };

const FlagTable = ({ flags }: { flags: FlagRow[] }) => (
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
);

export default ReleaseAliasPage;