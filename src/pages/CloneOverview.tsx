import DocsLayout from "@/components/docs/DocsLayout";
import CodeBlock from "@/components/docs/CodeBlock";

const fileFormats = [
  { fmt: "JSON", path: ".gitmap/output/gitmap.json", shorthand: "json", note: "Default — produced by `gitmap scan`. Hierarchical groups + per-repo metadata." },
  { fmt: "CSV",  path: ".gitmap/output/gitmap.csv",  shorthand: "csv",  note: "Spreadsheet-friendly. One row per repo, headers in row 1." },
  { fmt: "Text", path: ".gitmap/output/gitmap.txt",  shorthand: "text", note: "One URL per line. Easiest to hand-edit and pipe in from other tools." },
];

const urlExamples = [
  { url: "https://github.com/org/wp-onboarding-v13.git", folder: "wp-onboarding", note: "Versioned URL → auto-flattened to base name folder, version recorded in RepoVersionHistory." },
  { url: "https://github.com/org/wp-alim.git",            folder: "wp-alim",      note: "Non-versioned URL → cloned into the repo's natural folder name." },
  { url: "git@github.com:org/private-svc.git",            folder: "private-svc",  note: "SSH URL → uses your existing key (or `--ssh-key <name>` for a named pair)." },
];

const CloneOverviewPage = () => {
  return (
    <DocsLayout>
      <h1 className="text-3xl font-heading font-bold mb-2 docs-h1">Clone Overview</h1>
      <p className="text-muted-foreground mb-2">
        <code className="docs-inline-code">gitmap clone</code> serves two distinct workflows: cloning many repos from a
        scan-output file, or cloning a single repo from a direct Git URL. Both modes share the same flags, manifest
        normalization, and Windows path canonicalization.
      </p>
      <p className="text-sm text-muted-foreground mb-8">
        Alias: <code className="docs-inline-code">c</code> — see also{" "}
        <a href="/clone-command" className="text-primary hover:underline">flag reference</a>,{" "}
        <a href="/clone-next" className="text-primary hover:underline">clone-next</a>.
      </p>

      <section className="space-y-10">
        {/* ───────────────────────── File-based ───────────────────────── */}
        <div>
          <h2 className="text-xl font-heading font-semibold mb-3 docs-h2">1. File-based clone (bulk)</h2>
          <p className="text-sm text-muted-foreground mb-4">
            Run <code className="docs-inline-code">gitmap scan</code> first to produce a manifest, then feed it to
            <code className="docs-inline-code">clone</code>. Useful for onboarding a new machine, mirroring an org,
            or rebuilding a workspace after a re-image.
          </p>

          <div className="overflow-x-auto rounded-lg border border-border mb-4">
            <table className="w-full text-sm">
              <thead className="bg-muted/40">
                <tr>
                  <th className="text-left px-4 py-2 font-mono">Format</th>
                  <th className="text-left px-4 py-2 font-mono">Default path</th>
                  <th className="text-left px-4 py-2 font-mono">Shorthand</th>
                  <th className="text-left px-4 py-2 font-mono">When to use</th>
                </tr>
              </thead>
              <tbody>
                {fileFormats.map((f) => (
                  <tr key={f.fmt} className="border-t border-border">
                    <td className="px-4 py-2"><code className="docs-inline-code">{f.fmt}</code></td>
                    <td className="px-4 py-2 text-muted-foreground"><code className="docs-inline-code">{f.path}</code></td>
                    <td className="px-4 py-2"><code className="docs-inline-code">{f.shorthand}</code></td>
                    <td className="px-4 py-2 text-muted-foreground">{f.note}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          <CodeBlock code={`gitmap scan
gitmap clone json --safe-pull`} title="Typical bulk workflow" />
          <p className="text-sm text-muted-foreground mt-3">
            <code className="docs-inline-code">--safe-pull</code> turns existing folders into <em>pull</em> operations
            instead of skipping them, with retry + lockfile diagnostics. Combine with
            <code className="docs-inline-code">--target-dir</code> to relocate the destination root.
          </p>
        </div>

        <hr className="docs-hr" />

        {/* ───────────────────────── URL-based ───────────────────────── */}
        <div>
          <h2 className="text-xl font-heading font-semibold mb-3 docs-h2">2. Direct URL clone (single)</h2>
          <p className="text-sm text-muted-foreground mb-4">
            Pass any HTTPS or SSH Git URL — optionally followed by a custom folder name. Versioned URLs (those whose
            repo name ends in <code className="docs-inline-code">-vN</code>) are auto-flattened into the base name
            and the transition is recorded in the <code className="docs-inline-code">RepoVersionHistory</code> table.
          </p>

          <div className="overflow-x-auto rounded-lg border border-border mb-4">
            <table className="w-full text-sm">
              <thead className="bg-muted/40">
                <tr>
                  <th className="text-left px-4 py-2 font-mono">URL</th>
                  <th className="text-left px-4 py-2 font-mono">Resolved folder</th>
                  <th className="text-left px-4 py-2 font-mono">Behavior</th>
                </tr>
              </thead>
              <tbody>
                {urlExamples.map((u) => (
                  <tr key={u.url} className="border-t border-border">
                    <td className="px-4 py-2"><code className="docs-inline-code">{u.url}</code></td>
                    <td className="px-4 py-2"><code className="docs-inline-code">{u.folder}</code></td>
                    <td className="px-4 py-2 text-muted-foreground">{u.note}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          <CodeBlock code={`gitmap clone https://github.com/org/wp-onboarding-v13.git
gitmap clone https://github.com/org/wp-alim.git "my-project"
gitmap clone git@github.com:org/private-svc.git --ssh-key work`} title="URL workflow" />
        </div>

        <hr className="docs-hr" />

        {/* ───────────────────────── Path normalization ───────────────────────── */}
        <div>
          <h2 className="text-xl font-heading font-semibold mb-3 docs-h2">Windows path canonicalization</h2>
          <p className="text-sm text-muted-foreground mb-3">
            Both modes route every destination through{" "}
            <code className="docs-inline-code">canonicalizePMPath</code> before writing the VS Code Project Manager
            <code className="docs-inline-code"> projects.json</code>:
          </p>
          <ol className="list-decimal list-inside text-sm text-muted-foreground space-y-1 mb-4">
            <li><code className="docs-inline-code">filepath.Clean</code> normalizes separators and removes redundant segments.</li>
            <li><code className="docs-inline-code">filepath.EvalSymlinks</code> resolves symlinks and Windows 8.3 short names (<code className="docs-inline-code">C:\PROGRA~1</code> → <code className="docs-inline-code">C:\Program Files</code>).</li>
            <li>On Windows, paths are lowercased to collapse case-only duplicates (<code className="docs-inline-code">C:\Foo</code> vs <code className="docs-inline-code">c:\foo</code>).</li>
          </ol>
          <p className="text-sm text-muted-foreground mb-3">
            If <code className="docs-inline-code">EvalSymlinks</code> fails (missing path, permissions), the system
            <strong> soft-fails</strong> back to the cleaned absolute path so the clone never blocks on sync resolution.
            Pass <code className="docs-inline-code">--debug-paths</code> (or set
            <code className="docs-inline-code"> GITMAP_DEBUG_PATHS=1</code>) to print before/after path transforms on
            stderr.
          </p>
          <CodeBlock code={`gitmap clone https://github.com/org/wp-alim.git --debug-paths`} title="Diagnose a sync mismatch" />
        </div>

        <hr className="docs-hr" />

        {/* ───────────────────────── See also ───────────────────────── */}
        <div>
          <h2 className="text-xl font-heading font-semibold mb-3 docs-h2">See also</h2>
          <ul className="list-disc list-inside text-muted-foreground space-y-1">
            <li><a href="/clone-command" className="text-primary hover:underline">clone (flag reference)</a> — full flag table</li>
            <li><a href="/scan-command" className="text-primary hover:underline">scan</a> — produce the manifest files</li>
            <li><a href="/clone-next" className="text-primary hover:underline">clone-next</a> — clone the next versioned iteration of the current repo</li>
            <li><a href="/version-history" className="text-primary hover:underline">version-history</a> — inspect transitions recorded by flatten</li>
            <li><a href="/desktop-sync" className="text-primary hover:underline">desktop-sync</a> — register clones with GitHub Desktop</li>
            <li><a href="/ssh" className="text-primary hover:underline">ssh</a> — manage named SSH keys for private repos</li>
          </ul>
        </div>
      </section>
    </DocsLayout>
  );
};

export default CloneOverviewPage;
