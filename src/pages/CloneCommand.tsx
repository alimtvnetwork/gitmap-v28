import DocsLayout from "@/components/docs/DocsLayout";
import CodeBlock from "@/components/docs/CodeBlock";

const flags = [
  { flag: "--target-dir <dir>", default: "current dir", desc: "Base directory for clones" },
  { flag: "--safe-pull", default: "false", desc: "Pull existing repos with retry + unlock diagnostics" },
  { flag: "--github-desktop", default: "false", desc: "Auto-register with GitHub Desktop (no prompt)" },
  { flag: "--verbose", default: "false", desc: "Write detailed debug log to a timestamped file" },
];

const CloneCommandPage = () => {
  return (
    <DocsLayout>
      <h1 className="text-3xl font-heading font-bold mb-2 docs-h1">gitmap clone</h1>
      <p className="text-muted-foreground mb-2">
        Clone repositories from a structured output file (JSON, CSV, or text), or clone a single
        repository directly from a Git URL.
      </p>
      <p className="text-sm text-muted-foreground mb-8">
        Alias: <code className="docs-inline-code">c</code>
      </p>

      <section className="space-y-8">
        <div>
          <h2 className="text-xl font-heading font-semibold mb-3 docs-h2">Usage</h2>
          <CodeBlock code={`gitmap clone <source|json|csv|text|url> [folder] [flags]`} title="Syntax" />
          <p className="text-sm text-muted-foreground mt-2">
            Shorthands <code className="docs-inline-code">json</code>,{" "}
            <code className="docs-inline-code">csv</code>, and{" "}
            <code className="docs-inline-code">text</code> auto-resolve to the default scan output files.
          </p>
        </div>

        <hr className="docs-hr" />

        <div>
          <h2 className="text-xl font-heading font-semibold mb-3 docs-h2">Flags</h2>
          <div className="overflow-x-auto rounded-lg border border-border">
            <table className="w-full text-sm">
              <thead className="bg-muted/40">
                <tr>
                  <th className="text-left px-4 py-2 font-mono">Flag</th>
                  <th className="text-left px-4 py-2 font-mono">Default</th>
                  <th className="text-left px-4 py-2 font-mono">Description</th>
                </tr>
              </thead>
              <tbody>
                {flags.map((f) => (
                  <tr key={f.flag} className="border-t border-border">
                    <td className="px-4 py-2"><code className="docs-inline-code">{f.flag}</code></td>
                    <td className="px-4 py-2 text-muted-foreground">{f.default}</td>
                    <td className="px-4 py-2 text-muted-foreground">{f.desc}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>

        <hr className="docs-hr" />

        <div>
          <h2 className="text-xl font-heading font-semibold mb-3 docs-h2">Examples</h2>

          <h3 className="text-base font-heading font-semibold mb-2 mt-4 docs-h3">Clone from a direct URL (versioned — auto-flattened)</h3>
          <CodeBlock code={`gitmap clone https://github.com/alimtvnetwork/wp-onboarding-v13.git`} title="Terminal" />
          <CodeBlock
            language="bash"
            title="Output"
            code={`Cloning wp-onboarding-v13 into wp-onboarding...
Cloned wp-onboarding-v13 successfully.
  + 1 repo(s) added to GitHub Desktop, 0 failed.
  Opening wp-onboarding in VS Code...
  VS Code opened.`}
          />

          <h3 className="text-base font-heading font-semibold mb-2 mt-6 docs-h3">Clone URL into a custom folder</h3>
          <CodeBlock code={`gitmap clone https://github.com/alimtvnetwork/wp-alim.git "my-project"`} title="Terminal" />

          <h3 className="text-base font-heading font-semibold mb-2 mt-6 docs-h3">Clone from JSON scan output</h3>
          <CodeBlock code={`gitmap clone json --target-dir D:\\projects`} title="Terminal" />
          <CodeBlock
            language="bash"
            title="Output"
            code={`Cloning from .gitmap/output/gitmap.json...
[1/12] Cloning my-api... done
[2/12] Cloning web-app... done
...
Clone complete: 12 succeeded, 0 failed`}
          />

          <h3 className="text-base font-heading font-semibold mb-2 mt-6 docs-h3">Safe-pull existing repos</h3>
          <CodeBlock code={`gitmap c csv --safe-pull`} title="Terminal" />
          <CodeBlock
            language="bash"
            title="Output"
            code={`[1/8] my-api exists -> pulling... Already up to date.
[2/8] web-app exists -> pulling... Updated (3 new commits)
[3/8] Cloning billing-svc... done
...
Clone complete: 8 succeeded, 0 failed`}
          />

          <h3 className="text-base font-heading font-semibold mb-2 mt-6 docs-h3">Verbose text-file clone</h3>
          <CodeBlock code={`gitmap clone text --verbose`} title="Terminal" />
        </div>

        <hr className="docs-hr" />

        <div>
          <h2 className="text-xl font-heading font-semibold mb-3 docs-h2">See also</h2>
          <ul className="list-disc list-inside text-muted-foreground space-y-1">
            <li><a href="/scan-command" className="text-primary hover:underline">scan</a> — Generate output files first</li>
            <li><a href="/clone-next" className="text-primary hover:underline">clone-next</a> — Clone next version of a repo</li>
            <li><a href="/desktop-sync" className="text-primary hover:underline">desktop-sync</a> — Sync repos to GitHub Desktop</li>
          </ul>
        </div>
      </section>
    </DocsLayout>
  );
};

export default CloneCommandPage;
