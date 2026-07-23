import DocsLayout from "@/components/docs/DocsLayout";
import CodeBlock from "@/components/docs/CodeBlock";
import { GitCommit, FileText, Terminal } from "lucide-react";

const ChangelogGenerate = () => {
  return (
    <DocsLayout>
      <div className="max-w-4xl">
        <div className="flex items-center gap-3 mb-2">
          <GitCommit className="h-8 w-8 text-primary" />
          <h1 className="text-3xl font-heading font-bold text-foreground docs-h1">Changelog Generate</h1>
        </div>
        <p className="text-muted-foreground mb-8 text-lg">
          Auto-generate changelog entries from commit messages between Git tags.
        </p>

        {/* Usage */}
        <section className="mb-10">
          <h2 className="text-xl font-heading font-semibold text-foreground mb-3 flex items-center gap-2">
            <Terminal className="h-5 w-5 text-primary" />
            Usage
          </h2>
          <CodeBlock
            code={`gitmap changelog-generate [--from <tag>] [--to <tag>] [--write]\ngitmap cg [--from <tag>] [--to <tag>] [--write]`}
            language="bash"
            title="Command"
          />
        </section>

        {/* Flags */}
        <section className="mb-10">
          <h2 className="text-xl font-heading font-semibold text-foreground mb-3">Flags</h2>
          <div className="rounded-lg border border-border overflow-hidden">
            <table className="w-full text-sm">
              <thead>
                <tr className="bg-muted/30 border-b border-border">
                  <th className="text-left px-4 py-2 font-mono text-muted-foreground">Flag</th>
                  <th className="text-left px-4 py-2 font-mono text-muted-foreground">Default</th>
                  <th className="text-left px-4 py-2 font-mono text-muted-foreground">Description</th>
                </tr>
              </thead>
              <tbody>
                <tr className="border-b border-border">
                  <td className="px-4 py-2 font-mono text-primary">--from</td>
                  <td className="px-4 py-2 font-mono text-muted-foreground">second-latest tag</td>
                  <td className="px-4 py-2 text-foreground">Start tag (older boundary)</td>
                </tr>
                <tr className="border-b border-border">
                  <td className="px-4 py-2 font-mono text-primary">--to</td>
                  <td className="px-4 py-2 font-mono text-muted-foreground">latest tag</td>
                  <td className="px-4 py-2 text-foreground">End tag or HEAD</td>
                </tr>
                <tr className="border-b border-border last:border-0">
                  <td className="px-4 py-2 font-mono text-primary">--write</td>
                  <td className="px-4 py-2 font-mono text-muted-foreground">false</td>
                  <td className="px-4 py-2 text-foreground">Prepend output to CHANGELOG.md instead of printing</td>
                </tr>
              </tbody>
            </table>
          </div>
        </section>

        {/* How It Works */}
        <section className="mb-10">
          <h2 className="text-xl font-heading font-semibold text-foreground mb-3 flex items-center gap-2">
            <FileText className="h-5 w-5 text-primary" />
            How It Works
          </h2>
          <div className="space-y-3 text-sm text-muted-foreground">
            <div className="rounded-lg border border-border p-4 bg-card">
              <h3 className="font-mono font-semibold text-foreground mb-2">Tag Resolution</h3>
              <ul className="space-y-1.5">
                <li className="flex items-start gap-2">
                  <span className="w-1.5 h-1.5 rounded-full bg-primary mt-1.5 shrink-0" />
                  If no flags given, uses the two most recent version tags
                </li>
                <li className="flex items-start gap-2">
                  <span className="w-1.5 h-1.5 rounded-full bg-primary mt-1.5 shrink-0" />
                  <code className="font-mono text-primary">--from</code> only → generates from that tag to HEAD (labeled "Unreleased")
                </li>
                <li className="flex items-start gap-2">
                  <span className="w-1.5 h-1.5 rounded-full bg-primary mt-1.5 shrink-0" />
                  Both <code className="font-mono text-primary">--from</code> and <code className="font-mono text-primary">--to</code> → explicit range
                </li>
              </ul>
            </div>
            <div className="rounded-lg border border-border p-4 bg-card">
              <h3 className="font-mono font-semibold text-foreground mb-2">Commit Filtering</h3>
              <ul className="space-y-1.5">
                <li className="flex items-start gap-2">
                  <span className="w-1.5 h-1.5 rounded-full bg-primary mt-1.5 shrink-0" />
                  Uses <code className="font-mono text-primary">git log --no-merges</code> to exclude merge commits
                </li>
                <li className="flex items-start gap-2">
                  <span className="w-1.5 h-1.5 rounded-full bg-primary mt-1.5 shrink-0" />
                  Extracts commit subject lines only (first line of each message)
                </li>
              </ul>
            </div>
          </div>
        </section>

        {/* Examples */}
        <section className="mb-10">
          <h2 className="text-xl font-heading font-semibold text-foreground mb-3">Examples</h2>
          <div className="space-y-4">
            <div>
              <h3 className="font-mono text-sm font-semibold text-foreground mb-2">Generate between latest two tags</h3>
              <CodeBlock
                code={`$ gitmap changelog-generate

  Changelog: v2.23.0 → v2.24.0

  Preview (use --write to save):

  ## v2.24.0

  - Add TUI log viewer with detail panel
  - Add release rollback on push failure
  - Fix watch interval validation edge case`}
                language="bash"
                title="Terminal"
              />
            </div>
            <div>
              <h3 className="font-mono text-sm font-semibold text-foreground mb-2">From a specific tag to HEAD</h3>
              <CodeBlock
                code={`$ gitmap cg --from v2.22.0

  Changelog: v2.22.0 → HEAD

  Preview (use --write to save):

  ## Unreleased

  - Add changelog-generate command
  - Add TUI log viewer with detail panel
  - Fix zip-group archive naming`}
                language="bash"
                title="Terminal"
              />
            </div>
            <div>
              <h3 className="font-mono text-sm font-semibold text-foreground mb-2">Write directly to CHANGELOG.md</h3>
              <CodeBlock
                code={`$ gitmap cg --from v2.23.0 --to v2.24.0 --write

  Changelog: v2.23.0 → v2.24.0

  ✓ Prepended changelog to CHANGELOG.md`}
                language="bash"
                title="Terminal"
              />
            </div>
          </div>
        </section>

        {/* Prerequisites */}
        <section className="mb-10">
          <h2 className="text-xl font-heading font-semibold text-foreground mb-3">Prerequisites</h2>
          <ul className="space-y-2 text-sm text-muted-foreground">
            <li className="flex items-start gap-2">
              <span className="w-1.5 h-1.5 rounded-full bg-primary mt-1.5 shrink-0" />
              Must be inside a Git repository
            </li>
            <li className="flex items-start gap-2">
              <span className="w-1.5 h-1.5 rounded-full bg-primary mt-1.5 shrink-0" />
              At least one version tag (<code className="font-mono text-primary">v*</code>) must exist
            </li>
          </ul>
        </section>
      </div>
    </DocsLayout>
  );
};

export default ChangelogGenerate;
