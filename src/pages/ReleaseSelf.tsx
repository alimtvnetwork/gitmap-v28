import DocsLayout from "@/components/docs/DocsLayout";
import CodeBlock from "@/components/docs/CodeBlock";
import TerminalDemo from "@/components/docs/TerminalDemo";
import { Tag, GitBranch, Terminal, ArrowLeftRight } from "lucide-react";

const features = [
  { icon: Tag, title: "Self-Release", desc: "Release gitmap itself from any directory" },
  { icon: ArrowLeftRight, title: "Auto-Fallback", desc: "release auto-detects non-Git dirs and falls back to self-release" },
  { icon: GitBranch, title: "Full Flag Parity", desc: "All release flags work identically in self-release mode" },
  { icon: Terminal, title: "Directory Return", desc: "Returns to your original working directory after release" },
];

const selfReleaseDemo = {
  title: "gitmap release-self --bump patch",
  lines: [
    { text: "$ gitmap rs --bump patch", type: "input" as const, delay: 0 },
    { text: "→ Self-release: switching to /home/user/go/src/gitmap-v27", type: "accent" as const, delay: 600 },
    { text: "v2.45.0 → v2.45.1", type: "output" as const, delay: 400 },
    { text: "Creating release v2.45.1...", type: "header" as const, delay: 400 },
    { text: "  ✓ Created branch release/v2.45.1", type: "output" as const, delay: 300 },
    { text: "  ✓ Created tag v2.45.1", type: "output" as const, delay: 300 },
    { text: "  ✓ Pushed branch and tag to origin", type: "output" as const, delay: 300 },
    { text: "  Release v2.45.1 complete.", type: "accent" as const, delay: 400 },
    { text: "✓ Returned to /home/user/projects/other-repo", type: "accent" as const, delay: 400 },
  ],
};

const autoFallbackDemo = {
  title: "Auto-fallback when not in a Git repo",
  lines: [
    { text: "$ cd /tmp", type: "input" as const, delay: 0 },
    { text: "$ gitmap release --bump minor", type: "input" as const, delay: 600 },
    { text: "→ Self-release: switching to /home/user/go/src/gitmap-v27", type: "accent" as const, delay: 600 },
    { text: "v2.45.0 → v2.46.0", type: "output" as const, delay: 400 },
    { text: "Creating release v2.46.0...", type: "header" as const, delay: 400 },
    { text: "  ✓ Created branch release/v2.46.0", type: "output" as const, delay: 300 },
    { text: "  ✓ Created tag v2.46.0", type: "output" as const, delay: 300 },
    { text: "  ✓ Pushed branch and tag to origin", type: "output" as const, delay: 300 },
    { text: "  Release v2.46.0 complete.", type: "accent" as const, delay: 400 },
    { text: "✓ Returned to /tmp", type: "accent" as const, delay: 400 },
  ],
};

const errorScenarios = [
  { scenario: "Executable path + DB fallback both fail", behavior: "Exit 1: could not locate gitmap source repository" },
  { scenario: "DB path stale (no .git)", behavior: "Falls through to error" },
  { scenario: "Release fails", behavior: "Standard release error handling; still returns to original dir" },
  { scenario: "Return chdir fails", behavior: "Warning printed; exit 0 (release succeeded)" },
];

const ReleaseSelfPage = () => {
  return (
    <DocsLayout>
      <div className="space-y-8">
        {/* Header */}
        <div>
          <div className="flex items-center gap-3 mb-2">
            <Tag className="h-8 w-8 text-primary" />
            <h1 className="text-3xl font-heading font-bold text-foreground docs-h1">release-self</h1>
            <span className="text-xs font-mono bg-muted px-2 py-0.5 rounded text-muted-foreground">rs</span>
            <span className="text-xs font-mono bg-muted px-2 py-0.5 rounded text-muted-foreground">rself</span>
          </div>
          <p className="text-muted-foreground text-lg">
            Release gitmap itself from any directory — explicitly or via auto-fallback.
          </p>
        </div>

        {/* Features */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {features.map((f) => (
            <div key={f.title} className="flex items-start gap-3 p-4 rounded-lg border border-border bg-card">
              <f.icon className="h-5 w-5 text-primary mt-0.5" />
              <div>
                <h3 className="font-mono font-semibold text-foreground text-sm">{f.title}</h3>
                <p className="text-muted-foreground text-sm">{f.desc}</p>
              </div>
            </div>
          ))}
        </div>

        {/* Usage */}
        <section>
          <h2 className="text-xl font-heading font-semibold text-foreground mb-3">Usage</h2>
          <CodeBlock code="gitmap release-self [version] [flags]" />
          <p className="text-sm text-muted-foreground mt-2">
            Accepts all flags supported by <code className="text-primary font-mono">release</code> —
            <code className="text-primary font-mono"> --assets</code>,{" "}
            <code className="text-primary font-mono">--bump</code>,{" "}
            <code className="text-primary font-mono">--draft</code>,{" "}
            <code className="text-primary font-mono">--dry-run</code>,{" "}
            <code className="text-primary font-mono">--compress</code>,{" "}
            <code className="text-primary font-mono">--checksums</code>, etc.
          </p>
        </section>

        {/* How It Works */}
        <section>
          <h2 className="text-xl font-heading font-semibold text-foreground mb-3">How It Works</h2>
          <div className="space-y-3">
            {[
              { step: "1", title: "Resolve source repo", desc: "Tries os.Executable() + symlink resolution first. If that fails (e.g., binary installed outside source tree), falls back to the source_repo_path stored in the SQLite Settings table." },
              { step: "2", title: "Check current directory", desc: "If already in the source repo, skips the directory switch and proceeds directly." },
              { step: "3", title: "Switch directory", desc: "Records the caller's working directory, then os.Chdir() into the resolved source repo root." },
              { step: "4", title: "Execute release", desc: "Runs the full release workflow (identical to gitmap release) from the source repo." },
              { step: "5", title: "Return to caller", desc: "os.Chdir() back to the original working directory and prints a confirmation message." },
            ].map((s) => (
              <div key={s.step} className="flex items-start gap-3 p-3 rounded-lg border border-border bg-card">
                <span className="flex-shrink-0 h-6 w-6 rounded-full bg-primary text-primary-foreground flex items-center justify-center text-xs font-mono font-bold">
                  {s.step}
                </span>
                <div>
                  <h4 className="font-mono font-semibold text-foreground text-sm">{s.title}</h4>
                  <p className="text-muted-foreground text-sm">{s.desc}</p>
                </div>
              </div>
            ))}
          </div>
        </section>

        {/* Terminal Demos */}
        <section>
          <h2 className="text-xl font-heading font-semibold text-foreground mb-3">Examples</h2>
          <div className="space-y-6">
            <div>
              <h3 className="text-sm font-mono font-semibold text-muted-foreground mb-2">Explicit self-release with bump</h3>
              <TerminalDemo title={selfReleaseDemo.title} lines={selfReleaseDemo.lines} autoPlay />
            </div>
            <div>
              <h3 className="text-sm font-mono font-semibold text-muted-foreground mb-2">Auto-fallback from non-Git directory</h3>
              <TerminalDemo title={autoFallbackDemo.title} lines={autoFallbackDemo.lines} />
            </div>
          </div>
        </section>

        {/* CLI Examples */}
        <section>
          <h2 className="text-xl font-heading font-semibold text-foreground mb-3">More Examples</h2>
          <div className="space-y-4">
            <div>
              <p className="text-sm text-muted-foreground mb-1">Dry-run self-release</p>
              <CodeBlock code="gitmap rs --bump minor --dry-run" />
            </div>
            <div>
              <p className="text-sm text-muted-foreground mb-1">Self-release with assets and compression</p>
              <CodeBlock code="gitmap rs v2.46.0 --assets ./dist --compress --checksums" />
            </div>
            <div>
              <p className="text-sm text-muted-foreground mb-1">Draft self-release</p>
              <CodeBlock code="gitmap rs --bump patch --draft --notes 'pre-release test'" />
            </div>
          </div>
        </section>

        {/* Error Scenarios */}
        <section>
          <h2 className="text-xl font-heading font-semibold text-foreground mb-3">Error Scenarios</h2>
          <div className="overflow-x-auto">
            <table className="w-full text-sm border border-border rounded-lg overflow-hidden">
              <thead>
                <tr className="bg-muted">
                  <th className="text-left p-3 font-mono font-semibold text-foreground">Scenario</th>
                  <th className="text-left p-3 font-mono font-semibold text-foreground">Behavior</th>
                </tr>
              </thead>
              <tbody>
                {errorScenarios.map((e, i) => (
                  <tr key={i} className="border-t border-border">
                    <td className="p-3 text-muted-foreground">{e.scenario}</td>
                    <td className="p-3 font-mono text-xs text-foreground">{e.behavior}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </section>

        {/* See Also */}
        <section>
          <h2 className="text-xl font-heading font-semibold text-foreground mb-3">See Also</h2>
          <div className="flex flex-wrap gap-2">
            {[
              { label: "Release", url: "/release" },
              { label: "Temp Release", url: "/temp-release" },
              { label: "Prune", url: "/prune" },
            ].map((link) => (
              <a
                key={link.label}
                href={link.url}
                className="text-sm font-mono text-primary hover:underline px-3 py-1.5 rounded-md border border-border bg-card hover:bg-muted transition-colors"
              >
                {link.label}
              </a>
            ))}
          </div>
        </section>
      </div>
    </DocsLayout>
  );
};

export default ReleaseSelfPage;
