import DocsLayout from "@/components/docs/DocsLayout";
import CodeBlock from "@/components/docs/CodeBlock";

const GettingStartedPage = () => {
  return (
    <DocsLayout>
      <h1 className="text-3xl font-heading font-bold mb-2 docs-h1">Getting Started</h1>
      <p className="text-muted-foreground mb-8">
        Get up and running with gitmap in under 5 minutes.
      </p>

      <section className="space-y-8">
        <div>
          <h2 className="text-xl font-heading font-semibold mb-3 docs-h2">1. Install gitmap-v27</h2>
          <p className="text-muted-foreground mb-3">
            Build from source using Go 1.21+:
          </p>
          <CodeBlock
            code={`go install github.com/alimtvnetwork/gitmap-v27/gitmap-v27@latest`}
            title="Terminal"
          />
          <p className="text-sm text-muted-foreground mt-2">
            Or clone the repo and build with the platform-appropriate script:
          </p>
          <CodeBlock code={`# Windows (PowerShell)\ngit clone https://github.com/alimtvnetwork/gitmap-v27/gitmap.git\ncd gitmap-v27\n./run.ps1`} title="PowerShell" language="powershell" />
          <CodeBlock code={`# Linux / macOS (Bash)\ngit clone https://github.com/alimtvnetwork/gitmap-v27/gitmap.git\ncd gitmap-v27\nchmod +x run.sh\n./run.sh`} title="Bash" language="bash" />
          <CodeBlock code={`# Or use Make (requires run.sh)\ncd gitmap-v27\nmake build`} title="Makefile" language="bash" />

          <div className="mt-4 p-4 rounded-lg border border-border bg-muted/30">
            <h3 className="text-sm font-mono font-semibold mb-2 docs-h3">Build script flags</h3>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-x-6 gap-y-1 text-sm text-muted-foreground">
              <div><code className="docs-inline-code">run.ps1 -NoPull</code> / <code className="docs-inline-code">run.sh -n</code> — skip git pull</div>
              <div><code className="docs-inline-code">run.ps1 -NoDeploy</code> / <code className="docs-inline-code">run.sh -d</code> — build only</div>
              <div><code className="docs-inline-code">run.ps1 -Update</code> / <code className="docs-inline-code">run.sh -u</code> — full update pipeline</div>
              <div><code className="docs-inline-code">run.ps1 -R list</code> / <code className="docs-inline-code">run.sh -r list</code> — build &amp; run</div>
              <div><code className="docs-inline-code">run.sh -t</code> / <code className="docs-inline-code">make test</code> — run tests</div>
            </div>
          </div>
        </div>

        <hr className="docs-hr" />

        <div>
          <h2 className="text-xl font-heading font-semibold mb-3 docs-h2">2. Run your first scan</h2>
          <p className="text-muted-foreground mb-3">
            Point gitmap at a directory containing Git repositories:
          </p>
          <CodeBlock code={`gitmap scan ~/projects`} title="Terminal" />
          <p className="text-sm text-muted-foreground mt-2">
            This generates <code className="docs-inline-code">.gitmap/output/</code> containing CSV, JSON,
            folder structure, and clone scripts.
          </p>
        </div>

        <hr className="docs-hr" />

        <div>
          <h2 className="text-xl font-heading font-semibold mb-3 docs-h2">3. Clone on another machine</h2>
          <p className="text-muted-foreground mb-3">
            Copy the output files and restore the exact folder structure:
          </p>
          <CodeBlock
            code={`gitmap clone json --target-dir ./projects`}
            title="Terminal"
          />
          <p className="text-sm text-muted-foreground mt-2">
            Shorthands <code className="docs-inline-code">json</code>,{" "}
            <code className="docs-inline-code">csv</code>, and{" "}
            <code className="docs-inline-code">text</code> auto-resolve to the default output files.
          </p>
        </div>

        <hr className="docs-hr" />

        <div>
          <h2 className="text-xl font-heading font-semibold mb-3 docs-h2">4. Set up shell navigation</h2>
          <p className="text-muted-foreground mb-3">
            Run <code className="docs-inline-code">gitmap setup</code> to auto-install the <code className="docs-inline-code">gcd</code> wrapper function. After restarting your terminal:
          </p>
          <CodeBlock
            code={`gcd myrepo          # jumps to repo directory\ngcd repos           # interactive picker\ngcd repos --group backend`}
            title="Usage"
          />
          <p className="text-sm text-muted-foreground mt-2">
            The <code className="docs-inline-code">gcd</code> function is added to your shell profile automatically. You can also install it manually — see the <a href="/cd" className="text-primary hover:underline">cd docs</a>.
          </p>
        </div>

        <hr className="docs-hr" />

        <div>
          <h2 className="text-xl font-heading font-semibold mb-3 docs-h2">5. Monitor your repos</h2>
          <p className="text-muted-foreground mb-3">
            Start a live dashboard to watch all tracked repos:
          </p>
          <CodeBlock code={`gitmap watch --interval 15`} title="Terminal" />
          <p className="text-sm text-muted-foreground mt-2">
            The dashboard auto-refreshes, showing dirty/clean status, ahead/behind counts, and stash entries.
          </p>
        </div>
      </section>
    </DocsLayout>
  );
};

export default GettingStartedPage;
