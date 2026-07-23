import DocsLayout from "@/components/docs/DocsLayout";
import CodeBlock from "@/components/docs/CodeBlock";

const ConfigPage = () => {
  return (
    <DocsLayout>
      <h1 className="text-3xl font-heading font-bold mb-2 docs-h1">Configuration</h1>
      <p className="text-muted-foreground mb-8">
        Customize gitmap behavior through JSON config files, profiles, and the three-layer merge pattern.
      </p>

      <section className="space-y-8">
        <div>
          <h2 className="text-xl font-heading font-semibold mb-3 docs-h2">Three-Layer Config</h2>
          <div className="bg-card border border-border rounded-lg p-4 mb-4">
            <div className="space-y-2 font-mono text-sm">
              <div className="flex items-center gap-3">
                <span className="bg-muted text-muted-foreground px-3 py-1 rounded">1. Defaults</span>
                <span className="text-muted-foreground">→ Constants in code (lowest priority)</span>
              </div>
              <div className="flex items-center gap-3">
                <span className="bg-primary/10 text-foreground border border-primary/20 dark:bg-primary/15 dark:text-primary px-3 py-1 rounded transition-colors duration-300 hover:border-primary/40 hover:shadow-sm hover:shadow-primary/10">2. Config file</span>
                <span className="text-muted-foreground">→ ./data/config.json</span>
              </div>
              <div className="flex items-center gap-3">
                <span className="bg-primary/20 text-foreground border border-primary/30 dark:bg-primary/20 dark:text-primary px-3 py-1 rounded font-semibold">3. CLI flags</span>
                <span className="text-muted-foreground">→ Always wins (highest priority)</span>
              </div>
            </div>
          </div>
          <div className="docs-blockquote mb-4">
            Missing config file → use defaults silently. Flags always override config file values.
          </div>
        </div>

        <hr className="docs-hr" />

        <div>
          <h2 className="text-xl font-heading font-semibold mb-3 docs-h2">config.json</h2>
          <p className="text-muted-foreground mb-3">
            The main config file controls scan defaults and release settings. Located at <code className="docs-inline-code">./data/config.json</code>:
          </p>
          <CodeBlock
            title="data/config.json"
            language="json"
            code={`{
  "defaultMode": "https",
  "defaultOutput": "terminal",
  "outputDir": ".gitmap/output",
  "excludeDirs": ["node_modules", ".git", "vendor", ".venv"],
  "notes": "",
  "release": {
    "targets": [
      {"goos": "windows", "goarch": "amd64"},
      {"goos": "linux", "goarch": "amd64"},
      {"goos": "darwin", "goarch": "arm64"}
    ],
    "checksums": true,
    "compress": false
  }
}`}
          />
        </div>

        <hr className="docs-hr" />

        <div>
          <h2 className="text-xl font-heading font-semibold mb-3 docs-h2">Config Fields</h2>
          <div className="rounded-lg border border-border overflow-hidden">
            <table className="w-full text-sm docs-table">
              <thead>
                <tr className="border-b border-border">
                  <th className="text-left px-4 py-2">Field</th>
                  <th className="text-left px-4 py-2">Type</th>
                  <th className="text-left px-4 py-2">Default</th>
                  <th className="text-left px-4 py-2">Description</th>
                </tr>
              </thead>
              <tbody>
                {[
                  ["defaultMode", "string", '"https"', 'Clone URL style: "https" or "ssh"'],
                  ["defaultOutput", "string", '"terminal"', 'Output format: "terminal", "csv", or "json"'],
                  ["outputDir", "string", '".gitmap/output"', "Directory for all generated output files"],
                  ["excludeDirs", "[]string", "[]", "Directory names to skip during recursive scan"],
                  ["notes", "string", '""', "Default note for all records"],
                  ["release", "object", "{}", "Release-specific settings (see below)"],
                ].map(([field, type, def, desc]) => (
                  <tr key={field} className="border-t border-border">
                    <td className="px-4 py-2 font-mono text-primary">{field}</td>
                    <td className="px-4 py-2 font-mono text-muted-foreground">{type}</td>
                    <td className="px-4 py-2 font-mono text-muted-foreground">{def}</td>
                    <td className="px-4 py-2 text-muted-foreground">{desc}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>

        <hr className="docs-hr" />

        <div>
          <h2 className="text-xl font-heading font-semibold mb-3 docs-h2">Release Config</h2>
          <p className="text-muted-foreground mb-3">
            The <code className="docs-inline-code">release</code> section configures cross-compilation defaults.
            CLI flags (<code className="docs-inline-code">--targets</code>, <code className="docs-inline-code">--compress</code>,{" "}
            <code className="docs-inline-code">--checksums</code>) always override these values.
          </p>
          <div className="rounded-lg border border-border overflow-hidden">
            <table className="w-full text-sm docs-table">
              <thead>
                <tr className="border-b border-border">
                  <th className="text-left px-4 py-2">Field</th>
                  <th className="text-left px-4 py-2">Type</th>
                  <th className="text-left px-4 py-2">Default</th>
                  <th className="text-left px-4 py-2">Description</th>
                </tr>
              </thead>
              <tbody>
                {[
                  ["targets", "[]object", "[] (all 6)", "Override cross-compile OS/arch matrix"],
                  ["checksums", "bool", "false", "Generate SHA256 checksums.txt for assets"],
                  ["compress", "bool", "false", "Wrap assets in .zip (Windows) or .tar.gz"],
                ].map(([field, type, def, desc]) => (
                  <tr key={field} className="border-t border-border">
                    <td className="px-4 py-2 font-mono text-primary">release.{field}</td>
                    <td className="px-4 py-2 font-mono text-muted-foreground">{type}</td>
                    <td className="px-4 py-2 font-mono text-muted-foreground">{def}</td>
                    <td className="px-4 py-2 text-muted-foreground">{desc}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
          <p className="text-sm text-muted-foreground mt-3">
            Each target object has <code className="docs-inline-code">goos</code> and{" "}
            <code className="docs-inline-code">goarch</code> string fields. Use{" "}
            <code className="docs-inline-code">gitmap release --list-targets</code> to verify
            the resolved matrix.
          </p>
        </div>

        <hr className="docs-hr" />

        <div>
          <h2 className="text-xl font-heading font-semibold mb-3 docs-h2">git-setup.json</h2>
          <p className="text-muted-foreground mb-3">
            Configure global Git settings applied by <code className="docs-inline-code">gitmap setup</code>:
          </p>
          <CodeBlock
            title="data/git-setup.json"
            language="json"
            code={`{
  "settings": [
    { "key": "core.autocrlf", "value": "true" },
    { "key": "diff.tool", "value": "vscode" },
    { "key": "merge.tool", "value": "vscode" }
  ]
}`}
          />
        </div>

        <hr className="docs-hr" />

        <div>
          <h2 className="text-xl font-heading font-semibold mb-3 docs-h2">Profiles</h2>
          <p className="text-muted-foreground mb-3">
            Maintain separate database environments (work, personal, client) using profiles:
          </p>
          <CodeBlock
            code={`# Create a new profile\ngitmap profile create work\n\n# Switch to it\ngitmap profile switch work\n\n# List all profiles\ngitmap profile list\n\n# Compare repos across profiles\ngitmap diff-profiles default work`}
            title="Terminal"
          />
          <p className="text-sm text-muted-foreground mt-2">
            Each profile has its own SQLite database file. The <code className="docs-inline-code">default</code> profile
            always exists and cannot be deleted. Profile config is stored in{" "}
            <code className="docs-inline-code">.gitmap/output/data/profiles.json</code>.
          </p>
        </div>

        <hr className="docs-hr" />

        <div>
          <h2 className="text-xl font-heading font-semibold mb-3 docs-h2">CD Defaults</h2>
          <p className="text-muted-foreground mb-3">
            Set default navigation paths for repos cloned to multiple locations:
          </p>
          <CodeBlock
            code={`gitmap cd set-default myrepo C:\\repos\\github\\myrepo\ngitmap cd clear-default myrepo`}
            title="Terminal"
          />
        </div>
      </section>
    </DocsLayout>
  );
};

export default ConfigPage;
