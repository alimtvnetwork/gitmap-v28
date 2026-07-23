import DocsLayout from "@/components/docs/DocsLayout";
import CodeBlock from "@/components/docs/CodeBlock";

const targets = [
  {
    name: "build",
    desc: "Full pipeline: pull, tidy, build, and deploy.",
    example: "make build",
  },
  {
    name: "run",
    desc: "Build, deploy, and run gitmap with optional arguments.",
    example: 'make run ARGS="list --json"',
  },
  {
    name: "test",
    desc: "Run all Go tests and generate a report.",
    example: "make test",
  },
  {
    name: "update",
    desc: "Pull, build, deploy, and sync the PATH binary.",
    example: "make update",
  },
  {
    name: "no-pull",
    desc: "Build without running git pull first.",
    example: "make no-pull",
  },
  {
    name: "no-deploy",
    desc: "Build without copying the binary to the deploy path.",
    example: "make no-deploy",
  },
  {
    name: "clean",
    desc: "Remove the bin/ build artifacts directory.",
    example: "make clean",
  },
  {
    name: "help",
    desc: "Show all available targets with descriptions.",
    example: "make help",
  },
];

const MakefilePage = () => {
  return (
    <DocsLayout>
      <h1 className="text-3xl font-heading font-bold mb-2 docs-h1">Makefile</h1>
      <p className="text-muted-foreground mb-8">
        A thin wrapper around <code className="font-mono text-primary">run.sh</code> for
        standard <code className="font-mono text-primary">make</code> workflows on Linux and macOS.
      </p>

      <section className="space-y-6">
        <div>
          <h2 className="text-xl font-heading font-semibold mb-3 text-foreground">Quick start</h2>
          <CodeBlock code={`cd gitmap-v28\nmake build          # full pipeline\nmake run ARGS="scan ~/code"   # build + run\nmake test           # tests with report`} title="Terminal" />
        </div>

        <div>
          <h2 className="text-xl font-heading font-semibold mb-4 text-foreground">Targets</h2>
          <div className="space-y-4">
            {targets.map((t) => (
              <div
                key={t.name}
                className="rounded-lg border border-border bg-muted/30 p-4"
              >
                <div className="flex items-baseline gap-3 mb-1">
                  <code className="font-mono font-semibold text-primary text-base">{t.name}</code>
                  <span className="text-sm text-muted-foreground">{t.desc}</span>
                </div>
                <CodeBlock code={t.example} title="Terminal" />
              </div>
            ))}
          </div>
        </div>

        <div>
          <h2 className="text-xl font-heading font-semibold mb-3 text-foreground">Passing arguments</h2>
          <p className="text-muted-foreground mb-3">
            The <code className="font-mono text-primary">run</code> target forwards
            everything in <code className="font-mono text-primary">ARGS</code> to gitmap after building:
          </p>
          <CodeBlock
            code={`make run ARGS="list --json"\nmake run ARGS="scan ~/projects"\nmake run ARGS="latest-branch --top 5"`}
            title="Terminal"
          />
        </div>

        <div>
          <h2 className="text-xl font-heading font-semibold mb-3 text-foreground">How it works</h2>
          <p className="text-muted-foreground">
            Each target delegates to{" "}
            <code className="font-mono text-primary">run.sh</code> with the appropriate flag.
            The Makefile lives at{" "}
            <code className="font-mono text-primary">gitmap-v28/Makefile</code> and invokes{" "}
            <code className="font-mono text-primary">../run.sh</code> relative to the repo root.
            On Windows, use <code className="font-mono text-primary">run.ps1</code> directly
            instead of Make.
          </p>
        </div>
      </section>
    </DocsLayout>
  );
};

export default MakefilePage;
