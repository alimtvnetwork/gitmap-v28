import DocsLayout from "@/components/docs/DocsLayout";
import CodeBlock from "@/components/docs/CodeBlock";

const flags = [
  { flag: "--config <path>", default: "./data/config.json", desc: "Config file path" },
  { flag: "--mode ssh|https", default: "https", desc: "Clone URL style" },
  { flag: "--output csv|json|terminal", default: "terminal", desc: "Output format" },
  { flag: "--output-path <dir>", default: "./.gitmap/output", desc: "Output directory" },
  { flag: "--out-file <path>", default: "—", desc: "Exact output file path" },
  { flag: "--github-desktop", default: "false", desc: "Add scanned repos to GitHub Desktop" },
  { flag: "--open", default: "false", desc: "Open output folder after scan" },
  { flag: "--quiet", default: "false", desc: "Suppress clone help section (CI/scripted use)" },
  { flag: "--no-vscode-sync", default: "false", desc: "Skip syncing into VS Code Project Manager" },
  { flag: "--no-auto-tags", default: "false", desc: "Skip auto-derived tags (git/node/go/...)" },
];

const ScanCommandPage = () => {
  return (
    <DocsLayout>
      <h1 className="text-3xl font-heading font-bold mb-2 docs-h1">gitmap scan</h1>
      <p className="text-muted-foreground mb-2">
        Scan a directory tree for Git repositories and record them in the local database.
      </p>
      <p className="text-sm text-muted-foreground mb-8">
        Alias: <code className="docs-inline-code">s</code>
      </p>

      <section className="space-y-8">
        <div>
          <h2 className="text-xl font-heading font-semibold mb-3 docs-h2">Usage</h2>
          <CodeBlock code={`gitmap scan [dir] [flags]`} title="Syntax" />
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

          <h3 className="text-base font-heading font-semibold mb-2 mt-4 docs-h3">Scan a directory</h3>
          <CodeBlock code={`gitmap scan D:\\wp-work`} title="Terminal" />
          <CodeBlock
            language="bash"
            title="Output"
            code={`Scanning D:\\wp-work...
[1/42] github/user/my-api
[2/42] github/user/web-app
[3/42] github/org/billing-svc
...
Found 42 repositories
✓ Output written to ./.gitmap/output/
✓ Database updated (42 repos)`}
          />

          <h3 className="text-base font-heading font-semibold mb-2 mt-6 docs-h3">JSON output with SSH URLs</h3>
          <CodeBlock code={`gitmap scan ~/work --output json --mode ssh`} title="Terminal" />
          <CodeBlock
            language="bash"
            title="Output"
            code={`Scanning ~/work...
Found 18 repositories
✓ .gitmap/output/gitmap.json written
✓ .gitmap/output/gitmap.csv written
✓ Clone URLs use SSH format (git@github.com:...)`}
          />

          <h3 className="text-base font-heading font-semibold mb-2 mt-6 docs-h3">Register with GitHub Desktop</h3>
          <CodeBlock code={`gitmap scan D:\\repos --github-desktop`} title="Terminal" />

          <h3 className="text-base font-heading font-semibold mb-2 mt-6 docs-h3">Quiet CSV scan</h3>
          <CodeBlock code={`gitmap s . --quiet --output csv`} title="Terminal" />
        </div>

        <hr className="docs-hr" />

        <div>
          <h2 className="text-xl font-heading font-semibold mb-3 docs-h2">See also</h2>
          <ul className="list-disc list-inside text-muted-foreground space-y-1">
            <li><a href="/clone-command" className="text-primary hover:underline">clone</a> — Clone repos from scan output</li>
            <li><a href="/commands" className="text-primary hover:underline">commands</a> — Full command reference</li>
            <li><a href="/flags" className="text-primary hover:underline">flag reference</a> — All flags across commands</li>
          </ul>
        </div>
      </section>
    </DocsLayout>
  );
};

export default ScanCommandPage;
