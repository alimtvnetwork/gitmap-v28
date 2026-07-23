import DocsLayout from "@/components/docs/DocsLayout";
import TerminalDemo from "@/components/docs/TerminalDemo";

const demos = [
  {
    title: "gitmap scan",
    lines: [
      { text: "gitmap scan ~/projects", type: "input" as const, delay: 800 },
      { text: "", delay: 200 },
      { text: "gitmap v2.28.0", type: "header" as const },
      { text: "", delay: 100 },
      { text: "Scanning: /home/user/projects", delay: 150 },
      { text: "Found 12 repositories", type: "accent" as const, delay: 300 },
      { text: "", delay: 100 },
      { text: "  myapp           main    https://github.com/user/myapp.git", delay: 80 },
      { text: "  api-server      main    https://github.com/user/api-server.git", delay: 80 },
      { text: "  shared-lib      develop https://github.com/user/shared-lib.git", delay: 80 },
      { text: "  docs-site       main    https://github.com/user/docs-site.git", delay: 80 },
      { text: "  cli-tools       main    https://github.com/user/cli-tools.git", delay: 80 },
      { text: "", delay: 100 },
      { text: "Output: .gitmap/output/", type: "accent" as const, delay: 200 },
      { text: "  gitmap.csv      (12 repos)", delay: 80 },
      { text: "  gitmap.json     (12 repos)", delay: 80 },
      { text: "  clone.ps1       (ready)", delay: 80 },
    ],
  },
  {
    title: "gitmap watch",
    lines: [
      { text: "gitmap watch --interval 10", type: "input" as const, delay: 800 },
      { text: "", delay: 300 },
      { text: "gitmap watch — refreshing every 10s", type: "header" as const, delay: 200 },
      { text: "", delay: 100 },
      { text: "REPO              BRANCH    STATUS    AHEAD  BEHIND  STASH", delay: 150 },
      { text: "────────────────  ────────  ────────  ─────  ──────  ─────", delay: 100 },
      { text: "myapp             main      clean       0      0      0", delay: 80 },
      { text: "api-server        main      dirty       2      0      1", type: "accent" as const, delay: 80 },
      { text: "shared-lib        develop   clean       0      3      0", delay: 80 },
      { text: "docs-site         main      clean       0      0      0", delay: 80 },
      { text: "cli-tools         main      dirty       1      0      0", type: "accent" as const, delay: 80 },
      { text: "", delay: 200 },
      { text: "5 repos  |  2 dirty  |  3 ahead  |  3 behind  |  1 stash", type: "accent" as const },
    ],
  },
  {
    title: "gitmap latest-branch",
    lines: [
      { text: "gitmap lb --top 3 --format terminal", type: "input" as const, delay: 800 },
      { text: "", delay: 300 },
      { text: "Latest branches (top 3):", type: "header" as const, delay: 200 },
      { text: "", delay: 100 },
      { text: "  1. feature/add-auth     a3f2c1d  14-Mar-2026 02:30 PM", delay: 120 },
      { text: "     → Add OAuth2 login flow", delay: 80 },
      { text: "", delay: 60 },
      { text: "  2. fix/memory-leak      b7e4a9f  13-Mar-2026 11:15 AM", delay: 120 },
      { text: "     → Fix goroutine leak in watcher", delay: 80 },
      { text: "", delay: 60 },
      { text: "  3. release/v2.28.0      c1d8f3e  12-Mar-2026 09:45 AM", delay: 120 },
      { text: "     → Release v2.28.0", delay: 80 },
    ],
  },
  {
    title: "gitmap clone",
    lines: [
      { text: "gitmap clone json --target ~/new-machine", type: "input" as const, delay: 800 },
      { text: "", delay: 300 },
      { text: "Cloning from: .gitmap/output/gitmap.json", type: "header" as const, delay: 200 },
      { text: "Target: /home/user/new-machine", delay: 150 },
      { text: "", delay: 200 },
      { text: "  [1/5]  myapp           ████████████████████  done", type: "accent" as const, delay: 400 },
      { text: "  [2/5]  api-server      ████████████████████  done", type: "accent" as const, delay: 350 },
      { text: "  [3/5]  shared-lib      ████████████████████  done", type: "accent" as const, delay: 300 },
      { text: "  [4/5]  docs-site       ████████████████████  done", type: "accent" as const, delay: 350 },
      { text: "  [5/5]  cli-tools       ████████████████████  done", type: "accent" as const, delay: 300 },
      { text: "", delay: 100 },
      { text: "✓ 5/5 repos cloned successfully", type: "accent" as const },
    ],
  },
  {
    title: "gitmap clone-next",
    lines: [
      { text: "D:\\wp-work\\riseup-asia\\macro-ahk-v11>", type: "header" as const, delay: 200 },
      { text: "gitmap cn v++", type: "input" as const, delay: 800 },
      { text: "", delay: 300 },
      { text: "Cloning macro-ahk-v12 into D:\\wp-work\\riseup-asia...", delay: 400 },
      { text: "✓ Cloned macro-ahk-v12", type: "accent" as const, delay: 350 },
      { text: "✓ Registered macro-ahk-v12 with GitHub Desktop", type: "accent" as const, delay: 200 },
      { text: "Remove current folder macro-ahk-v11? [y/N] y", delay: 300 },
      { text: "✓ Removed macro-ahk-v11", type: "accent" as const, delay: 200 },
      { text: "→ Now in macro-ahk-v12", type: "accent" as const, delay: 200 },
      { text: "", delay: 200 },
      { text: "D:\\wp-work\\riseup-asia\\macro-ahk-v12>", type: "header" as const, delay: 300 },
      { text: "gitmap cn v15 --delete", type: "input" as const, delay: 800 },
      { text: "", delay: 300 },
      { text: "Cloning macro-ahk-v15 into D:\\wp-work\\riseup-asia...", delay: 400 },
      { text: "✓ Cloned macro-ahk-v15", type: "accent" as const, delay: 350 },
      { text: "✓ Registered macro-ahk-v15 with GitHub Desktop", type: "accent" as const, delay: 200 },
      { text: "✓ Removed macro-ahk-v12", type: "accent" as const, delay: 200 },
      { text: "→ Now in macro-ahk-v15", type: "accent" as const, delay: 200 },
    ],
  },
];

const InteractiveExamplesPage = () => {
  return (
    <DocsLayout>
      <h1 className="text-3xl font-heading font-bold mb-2 docs-h1">Interactive Examples</h1>
      <p className="text-muted-foreground mb-8">
        Live terminal demos for key gitmap commands. Click <strong>▶</strong> to play.
      </p>

      <div className="space-y-8">
        {demos.map((demo) => (
          <div key={demo.title}>
            <h2 className="text-lg font-mono font-semibold mb-3">{demo.title}</h2>
            <TerminalDemo title={demo.title} lines={demo.lines} />
          </div>
        ))}
      </div>
    </DocsLayout>
  );
};

export default InteractiveExamplesPage;
