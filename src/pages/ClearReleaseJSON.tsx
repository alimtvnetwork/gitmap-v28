import DocsLayout from "@/components/docs/DocsLayout";
import CodeBlock from "@/components/docs/CodeBlock";
import { Trash2, Eye, AlertTriangle, Terminal, FileJson } from "lucide-react";

const features = [
  { icon: Trash2, title: "Targeted Removal", desc: "Delete a specific .gitmap/release/vX.Y.Z.json metadata file." },
  { icon: Eye, title: "Dry Run", desc: "Preview which file would be removed without deleting it." },
  { icon: FileJson, title: "Semver Normalisation", desc: "Partial versions auto-pad: v2 → v2.0.0, v2.1 → v2.1.0." },
  { icon: AlertTriangle, title: "Existence Check", desc: "Validates the file exists before attempting removal." },
];

const flags = [
  { flag: "--dry-run", description: "Preview which file would be removed without deleting it" },
];

const edgeCases = [
  { scenario: "No version argument", behavior: "Print usage message and exit 1" },
  { scenario: "Invalid version string (e.g. abc)", behavior: "Print error with invalid input and exit 1" },
  { scenario: "File does not exist", behavior: "Print 'no release file found' error and exit 1" },
  { scenario: "File is read-only", behavior: "os.Remove fails; print removal error and exit 1" },
  { scenario: "--dry-run with missing file", behavior: "Same missing-file error — dry-run still validates existence" },
  { scenario: "--dry-run with valid file", behavior: "Print preview message and exit 0; file untouched" },
  { scenario: "Partial version v2", behavior: "Normalized to v2.0.0; targets .gitmap/release/v2.0.0.json" },
];

const exitCodes = [
  { code: "0", meaning: "File removed successfully, or dry-run preview printed" },
  { code: "1", meaning: "Missing argument, invalid version, file not found, or removal failed" },
];

const constants = [
  { constant: "MsgClearReleaseDone", format: "✓ Removed .gitmap/release/%s.json" },
  { constant: "MsgClearReleaseDryRun", format: "[dry-run] Would remove %s" },
  { constant: "ErrClearReleaseUsage", format: "Usage: gitmap clear-release-json <version> [--dry-run]" },
  { constant: "ErrClearReleaseNotFound", format: "Error: no release file found for %s" },
  { constant: "ErrClearReleaseFailed", format: "Error: could not remove release file: %v" },
];

const ClearReleaseJSONPage = () => {
  return (
    <DocsLayout>
      <div className="space-y-8">
        <div>
          <h1 className="text-3xl font-heading font-bold mb-2 docs-h1">Clear Release JSON</h1>
          <p className="text-muted-foreground text-lg">
            Remove a specific release metadata JSON file from the <code className="text-primary font-mono text-base">.gitmap/release/</code> directory.
          </p>
        </div>

        {/* Command & Alias */}
        <section>
          <h2 className="text-xl font-heading font-semibold mb-3 docs-h2">Command</h2>
          <CodeBlock code="gitmap clear-release-json <version> [--dry-run]" />
          <p className="text-sm text-muted-foreground mt-2">
            Alias: <code className="font-mono text-primary">crj</code>
          </p>
        </section>

        {/* Features */}
        <section>
          <h2 className="text-xl font-heading font-semibold mb-4">Features</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            {features.map((f) => (
              <div key={f.title} className="rounded-lg border border-border bg-card p-4 flex gap-3">
                <f.icon className="h-5 w-5 text-primary mt-0.5 shrink-0" />
                <div>
                  <div className="font-mono font-medium text-sm">{f.title}</div>
                  <div className="text-sm text-muted-foreground">{f.desc}</div>
                </div>
              </div>
            ))}
          </div>
        </section>

        {/* Flags */}
        <section>
          <h2 className="text-xl font-heading font-semibold mb-3 docs-h2">Flags</h2>
          <div className="rounded-lg border border-border overflow-hidden">
            <table className="w-full text-sm">
              <thead>
                <tr className="bg-muted/50">
                  <th className="text-left px-4 py-2 font-mono font-medium">Flag</th>
                  <th className="text-left px-4 py-2 font-mono font-medium">Description</th>
                </tr>
              </thead>
              <tbody>
                {flags.map((f) => (
                  <tr key={f.flag} className="border-t border-border">
                    <td className="px-4 py-2 font-mono text-primary">{f.flag}</td>
                    <td className="px-4 py-2 text-muted-foreground">{f.description}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </section>

        {/* Version Resolution */}
        <section>
          <h2 className="text-xl font-heading font-semibold mb-3 docs-h2">Version Resolution</h2>
          <p className="text-muted-foreground mb-3">
            The version argument is parsed through <code className="font-mono text-primary">release.Parse</code>, which applies standard semver normalisation:
          </p>
          <ul className="list-disc list-inside space-y-1 text-muted-foreground text-sm">
            <li>A leading <code className="font-mono text-primary">v</code> prefix is optional — <code className="font-mono">2.20.0</code> and <code className="font-mono">v2.20.0</code> are equivalent.</li>
            <li>Partial versions are zero-padded — <code className="font-mono">v2</code> becomes <code className="font-mono">v2.0.0</code>.</li>
            <li>Pre-release suffixes (e.g. <code className="font-mono">v3.0.0-rc.1</code>) are preserved as-is.</li>
          </ul>
        </section>

        {/* Behavior */}
        <section>
          <h2 className="text-xl font-heading font-semibold mb-3 docs-h2">Behavior</h2>

          <h3 className="text-lg font-mono font-medium mb-2 mt-4">Normal Mode</h3>
          <CodeBlock code={`# Remove release metadata for v2.20.0
gitmap clear-release-json v2.20.0
# ✓ Removed .gitmap/release/v2.20.0.json`} title="Normal removal" />

          <h3 className="text-lg font-mono font-medium mb-2 mt-4">Dry-Run Mode</h3>
          <CodeBlock code={`# Preview without deleting
gitmap clear-release-json v2.20.0 --dry-run
# [dry-run] Would remove .gitmap/release/v2.20.0.json`} title="Dry-run preview" />
        </section>

        {/* Edge Cases */}
        <section>
          <h2 className="text-xl font-heading font-semibold mb-3 docs-h2">Edge Cases</h2>
          <div className="rounded-lg border border-border overflow-hidden">
            <table className="w-full text-sm">
              <thead>
                <tr className="bg-muted/50">
                  <th className="text-left px-4 py-2 font-mono font-medium">Scenario</th>
                  <th className="text-left px-4 py-2 font-mono font-medium">Behavior</th>
                </tr>
              </thead>
              <tbody>
                {edgeCases.map((ec) => (
                  <tr key={ec.scenario} className="border-t border-border">
                    <td className="px-4 py-2 text-foreground">{ec.scenario}</td>
                    <td className="px-4 py-2 text-muted-foreground">{ec.behavior}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </section>

        {/* Exit Codes */}
        <section>
          <h2 className="text-xl font-heading font-semibold mb-3 docs-h2">Exit Codes</h2>
          <div className="rounded-lg border border-border overflow-hidden">
            <table className="w-full text-sm">
              <thead>
                <tr className="bg-muted/50">
                  <th className="text-left px-4 py-2 font-mono font-medium">Code</th>
                  <th className="text-left px-4 py-2 font-mono font-medium">Meaning</th>
                </tr>
              </thead>
              <tbody>
                {exitCodes.map((ec) => (
                  <tr key={ec.code} className="border-t border-border">
                    <td className="px-4 py-2 font-mono text-primary">{ec.code}</td>
                    <td className="px-4 py-2 text-muted-foreground">{ec.meaning}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </section>

        {/* Output Constants */}
        <section>
          <h2 className="text-xl font-heading font-semibold mb-3 docs-h2">Output Constants</h2>
          <div className="rounded-lg border border-border overflow-hidden">
            <table className="w-full text-sm">
              <thead>
                <tr className="bg-muted/50">
                  <th className="text-left px-4 py-2 font-mono font-medium">Constant</th>
                  <th className="text-left px-4 py-2 font-mono font-medium">Format String</th>
                </tr>
              </thead>
              <tbody>
                {constants.map((c) => (
                  <tr key={c.constant} className="border-t border-border">
                    <td className="px-4 py-2 font-mono text-primary text-xs">{c.constant}</td>
                    <td className="px-4 py-2 font-mono text-muted-foreground text-xs">{c.format}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </section>

        {/* Implementation */}
        <section>
          <h2 className="text-xl font-heading font-semibold mb-3 docs-h2">Implementation</h2>
          <div className="rounded-lg border border-border overflow-hidden">
            <table className="w-full text-sm">
              <thead>
                <tr className="bg-muted/50">
                  <th className="text-left px-4 py-2 font-mono font-medium">File</th>
                  <th className="text-left px-4 py-2 font-mono font-medium">Responsibility</th>
                </tr>
              </thead>
              <tbody>
                {[
                  { file: "cmd/clearreleasejson.go", resp: "Flag parsing and handler" },
                  { file: "constants/constants_messages.go", resp: "Message and error format strings" },
                  { file: "release/metadata.go", resp: "ReleaseExists, metaFilePath" },
                  { file: "release/semver.go", resp: "Version normalisation and validation" },
                  { file: "helptext/clear-release-json.md", resp: "Embedded help text" },
                ].map((row) => (
                  <tr key={row.file} className="border-t border-border">
                    <td className="px-4 py-2 font-mono text-primary text-xs">{row.file}</td>
                    <td className="px-4 py-2 text-muted-foreground">{row.resp}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </section>

        {/* See Also */}
        <section>
          <h2 className="text-xl font-heading font-semibold mb-3 docs-h2">See Also</h2>
          <ul className="space-y-1 text-sm">
            <li><a href="/release" className="text-primary hover:underline font-mono">release</a> — Create a release</li>
            <li><a href="/commands" className="text-primary hover:underline font-mono">list-releases</a> — Show stored releases</li>
            <li><a href="/commands" className="text-primary hover:underline font-mono">commands</a> — Full command reference</li>
          </ul>
        </section>
      </div>
    </DocsLayout>
  );
};

export default ClearReleaseJSONPage;
