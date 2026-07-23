import { useState, useEffect, useCallback } from "react";
import DocsLayout from "@/components/docs/DocsLayout";
import CodeBlock from "@/components/docs/CodeBlock";
import { Monitor, RefreshCw, Clock, Zap } from "lucide-react";
import { RepoStatus, STATUS_ICON_DIRTY, STATUS_ICON_CLEAN } from "@/constants";

interface MockRepo {
  name: string;
  status: RepoStatus;
  branch: string;
  ahead: number;
  behind: number;
  stash: number;
}

const MOCK_REPOS: MockRepo[] = [
  { name: "api-gateway", status: RepoStatus.Clean, branch: "main", ahead: 0, behind: 0, stash: 0 },
  { name: "frontend-app", status: RepoStatus.Dirty, branch: "feat/nav", ahead: 2, behind: 0, stash: 1 },
  { name: "shared-lib", status: RepoStatus.Clean, branch: "main", ahead: 0, behind: 3, stash: 0 },
  { name: "docs-site", status: RepoStatus.Dirty, branch: "update-faq", ahead: 1, behind: 1, stash: 0 },
  { name: "infra-config", status: RepoStatus.Clean, branch: "main", ahead: 0, behind: 0, stash: 2 },
  { name: "mobile-app", status: RepoStatus.Clean, branch: "develop", ahead: 0, behind: 5, stash: 0 },
];

const isDirty = (status: RepoStatus): boolean => status === RepoStatus.Dirty;

const statusColor = (status: RepoStatus) =>
  isDirty(status) ? "text-yellow-400" : "text-primary";

const hasCount = (count: number): boolean => count > 0;

const countColor = (count: number) =>
  hasCount(count) ? "text-yellow-400" : "text-muted-foreground";

const TerminalPreview = () => {
  const [tick, setTick] = useState(0);
  const [paused, setPaused] = useState(false);

  useEffect(() => {
    if (paused) return;
    const id = setInterval(() => setTick((t) => t + 1), 4000);
    return () => clearInterval(id);
  }, [paused]);

  const now = new Date();
  const timeStr = now.toLocaleTimeString();

  const dirty = MOCK_REPOS.filter((repo) => isDirty(repo.status)).length;
  const behind = MOCK_REPOS.filter((repo) => hasCount(repo.behind)).length;
  const stash = MOCK_REPOS.reduce((acc, repo) => acc + repo.stash, 0);

  return (
    <div className="rounded-lg border border-border overflow-hidden my-6">
      <div className="bg-terminal px-4 py-2 flex items-center justify-between border-b border-border">
        <div className="flex items-center gap-2">
          <div className="flex gap-1.5">
            <span className="w-3 h-3 rounded-full bg-red-500/80" />
            <span className="w-3 h-3 rounded-full bg-yellow-500/80" />
            <span className="w-3 h-3 rounded-full bg-green-500/80" />
          </div>
          <span className="text-xs font-mono text-muted-foreground ml-2">
            gitmap watch
          </span>
        </div>
        <button
          onClick={() => setPaused(!paused)}
          className="text-xs font-mono text-muted-foreground hover:text-foreground transition-colors flex items-center gap-1"
        >
          {paused ? (
            <>
              <RefreshCw className="h-3 w-3" /> Resume
            </>
          ) : (
            <>
              <Clock className="h-3 w-3" /> Pause
            </>
          )}
        </button>
      </div>

      <div className="bg-terminal p-4 font-mono text-sm leading-relaxed overflow-x-auto">
        {/* Banner */}
        <div className="text-primary">
          <div>╔══════════════════════════════════════╗</div>
          <div>║          gitmap watch{"                "}║</div>
          <div>╚══════════════════════════════════════╝</div>
        </div>

        <div className="text-muted-foreground mt-2 text-xs">
          gitmap watch — refreshing every 30s (Ctrl+C to stop)
        </div>
        <div className="text-muted-foreground text-xs mb-3">
          Last updated: {timeStr}
          {tick > 0 && (
            <span className="text-primary ml-2 animate-pulse">●</span>
          )}
        </div>

        {/* Header */}
        <div className="text-primary font-bold">
          {"  "}REPO{"                   "}STATUS{"     "}BRANCH{"           "}AHEAD{"  "}BEHIND{"  "}STASH
        </div>

        {/* Rows */}
        {MOCK_REPOS.map((repo) => (
          <div key={repo.name} className="text-terminal-foreground">
            {"  "}
            <span className="inline-block w-[170px]">{repo.name}</span>
            <span className={`inline-block w-[72px] ${statusColor(repo.status)}`}>
              {isDirty(repo.status) ? STATUS_ICON_DIRTY : STATUS_ICON_CLEAN} {repo.status}
            </span>
            <span className="inline-block w-[130px] text-muted-foreground">
              {repo.branch}
            </span>
            <span className={`inline-block w-[48px] ${countColor(repo.ahead)}`}>
              {repo.ahead}
            </span>
            <span className={`inline-block w-[56px] ${countColor(repo.behind)}`}>
              {repo.behind}
            </span>
            <span className={`inline-block w-[40px] ${countColor(repo.stash)}`}>
              {repo.stash}
            </span>
          </div>
        ))}

        {/* Summary */}
        <div className="mt-3 text-muted-foreground text-xs border-t border-border/50 pt-2">
          Repos: {MOCK_REPOS.length} | Dirty: {dirty} | Behind: {behind} | Stash: {stash}
        </div>
      </div>
    </div>
  );
};

const features = [
  {
    icon: Monitor,
    title: "Live Dashboard",
    desc: "Full-screen table refreshes automatically, showing branch, dirty status, ahead/behind counts, and stash.",
  },
  {
    icon: RefreshCw,
    title: "Configurable Interval",
    desc: "Set refresh rate from 5 seconds up. Runs git fetch before each cycle unless --no-fetch is passed.",
  },
  {
    icon: Zap,
    title: "JSON Snapshots",
    desc: "Use --json to capture a single point-in-time snapshot for CI pipelines or custom tooling.",
  },
];

const WatchPage = () => {
  return (
    <DocsLayout>
      <h1 className="text-3xl font-heading font-bold mb-2 docs-h1">Watch Command</h1>
      <p className="text-muted-foreground mb-6">
        Real-time monitoring dashboard for all your tracked repositories.
      </p>

      {/* Interactive preview */}
      <h2 className="text-xl font-heading font-semibold mt-8 mb-2">
        Live Preview
      </h2>
      <p className="text-sm text-muted-foreground mb-2">
        Interactive simulation of the watch dashboard. Click Pause/Resume to
        control the refresh indicator.
      </p>
      <TerminalPreview />

      {/* Features */}
      <h2 className="text-xl font-heading font-semibold mt-10 mb-4">Features</h2>
      <div className="grid md:grid-cols-3 gap-4 mb-8">
        {features.map((f) => (
          <div
            key={f.title}
            className="rounded-lg border border-border bg-card p-4"
          >
            <f.icon className="h-5 w-5 text-primary mb-2" />
            <h3 className="font-mono font-semibold text-sm mb-1">{f.title}</h3>
            <p className="text-xs text-muted-foreground">{f.desc}</p>
          </div>
        ))}
      </div>

      {/* Usage */}
      <h2 className="text-xl font-heading font-semibold mt-10 mb-3">Usage</h2>

      <CodeBlock code="gitmap watch" title="Basic — refresh every 30s" />
      <CodeBlock
        code="gitmap watch --interval 10"
        title="Custom refresh interval"
      />
      <CodeBlock
        code="gitmap watch --group work"
        title="Filter by group"
      />
      <CodeBlock
        code="gitmap watch --no-fetch"
        title="Skip git fetch (local refs only)"
      />
      <CodeBlock
        code="gitmap watch --json"
        title="Single JSON snapshot and exit"
      />

      {/* Flags */}
      <h2 className="text-xl font-heading font-semibold mt-10 mb-3">Flags</h2>
      <div className="rounded-lg border border-border overflow-hidden">
        <table className="w-full text-sm">
          <thead>
            <tr className="bg-muted/50">
              <th className="text-left font-mono font-semibold px-4 py-2">
                Flag
              </th>
              <th className="text-left font-mono font-semibold px-4 py-2">
                Description
              </th>
            </tr>
          </thead>
          <tbody>
            {[
              ["--interval <seconds>", "Refresh interval in seconds (min 5, default 30)"],
              ["--group <name>", "Monitor only repos in the specified group"],
              ["--no-fetch", "Skip git fetch; use local refs only"],
              ["--json", "Output a single JSON snapshot and exit"],
            ].map(([flag, desc]) => (
              <tr key={flag} className="border-t border-border">
                <td className="px-4 py-2 font-mono text-primary">{flag}</td>
                <td className="px-4 py-2 text-muted-foreground">{desc}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* JSON output */}
      <h2 className="text-xl font-heading font-semibold mt-10 mb-3">
        JSON Output
      </h2>
      <p className="text-sm text-muted-foreground mb-2">
        When <code className="text-primary font-mono">--json</code> is passed,
        watch prints a single snapshot and exits:
      </p>
      <CodeBlock
        title="gitmap watch --json"
        code={`[
  {
    "repo": "api-gateway",
    "status": "clean",
    "branch": "main",
    "ahead": 0,
    "behind": 0,
    "stash": 0
  },
  {
    "repo": "frontend-app",
    "status": "dirty",
    "branch": "feat/nav",
    "ahead": 2,
    "behind": 0,
    "stash": 1
  }
]`}
      />

      {/* File layout */}
      <h2 className="text-xl font-heading font-semibold mt-10 mb-3">
        File Layout
      </h2>
      <div className="rounded-lg border border-border overflow-hidden">
        <table className="w-full text-sm">
          <thead>
            <tr className="bg-muted/50">
              <th className="text-left font-mono font-semibold px-4 py-2">File</th>
              <th className="text-left font-mono font-semibold px-4 py-2">Purpose</th>
            </tr>
          </thead>
          <tbody>
            {[
              ["constants/constants_watch.go", "Command names, ANSI codes, format strings"],
              ["cmd/watch.go", "Entry point, flag parsing, refresh loop"],
              ["cmd/watchops.go", "Repo loading and snapshot collection"],
              ["cmd/watchformat.go", "Terminal table rendering and summary"],
              ["gitutil/watchstatus.go", "Per-repo git status collection"],
            ].map(([file, purpose]) => (
              <tr key={file} className="border-t border-border">
                <td className="px-4 py-2 font-mono text-primary">{file}</td>
                <td className="px-4 py-2 text-muted-foreground">{purpose}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </DocsLayout>
  );
};

export default WatchPage;
