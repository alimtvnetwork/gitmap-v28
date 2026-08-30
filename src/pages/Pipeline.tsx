import { useState, useEffect } from "react";
import DocsLayout from "@/components/docs/DocsLayout";
import { Activity, Clock, CheckCircle2, XCircle, AlertCircle, RefreshCw, GitPullRequest, Tag, ExternalLink, Download, Copy, Check } from "lucide-react";
import { copyToClipboard } from "@/lib/clipboard";
import { useToast } from "@/hooks/use-toast";

interface PipelineStatusData {
  isRunning: boolean;
  etaSeconds: number;
  lastTagRelease: string;
  pendingPipelines: number;
  pendingTasks: number;
  pendingPRs: number;
  repo: string;
  activeWorkflow?: string;
  lastStatus?: string;
  lastConclusion?: string;
  lastRunUrl?: string;
  updatedAt: string;
}

export default function PipelinePage() {
  const { toast } = useToast();
  const [isCopied, setIsCopied] = useState(false);
  const [data, setData] = useState<PipelineStatusData>({
    isRunning: false,
    etaSeconds: 0,
    lastTagRelease: "v6.153.0",
    pendingPipelines: 0,
    pendingTasks: 0,
    pendingPRs: 1,
    repo: "alimtvnetwork/gitmap-v28",
    activeWorkflow: "CI Quality Gates",
    lastStatus: "completed",
    lastConclusion: "success",
    lastRunUrl: "https://github.com/alimtvnetwork/gitmap-v28/actions",
    updatedAt: new Date().toISOString(),
  });

  const handleCopyJSON = () => {
    copyToClipboard(JSON.stringify(data, null, 2));
    setIsCopied(true);
    toast({ title: "Copied!", description: "Pipeline status payload copied to clipboard." });
    setTimeout(() => setIsCopied(false), 2000);
  };

  return (
    <DocsLayout>
      <div className="max-w-5xl mx-auto space-y-8 pb-16">
        {/* Header */}
        <div className="flex flex-wrap items-center justify-between gap-4">
          <div className="flex items-center gap-3">
            <div className="h-10 w-10 rounded-lg bg-primary/10 flex items-center justify-center dark:bg-primary/20">
              <Activity className="h-5 w-5 text-primary" />
            </div>
            <div>
              <h1 className="text-3xl font-heading font-bold text-foreground">CI/CD Pipeline Monitor</h1>
              <p className="text-sm text-muted-foreground">Real-time status, ETA completion estimation, and failure logs.</p>
            </div>
          </div>
          <div className="flex items-center gap-2">
            <button
              type="button"
              onClick={handleCopyJSON}
              className="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-md border border-border bg-card text-xs font-mono text-foreground hover:bg-muted transition-colors"
            >
              {isCopied ? <Check className="h-3.5 w-3.5 text-primary" /> : <Copy className="h-3.5 w-3.5" />}
              {isCopied ? "Copied" : "Copy Status JSON"}
            </button>
          </div>
        </div>

        {/* Status Metrics Cards */}
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          <div className="rounded-xl border border-border bg-card p-5 space-y-2">
            <div className="flex items-center justify-between text-muted-foreground text-xs font-mono uppercase tracking-wider">
              <span>Pipeline Status</span>
              <Activity className="h-4 w-4 text-primary" />
            </div>
            <div className="flex items-center gap-2">
              {data.isRunning ? (
                <span className="inline-flex items-center gap-1.5 text-sm font-semibold text-yellow-500 bg-yellow-500/10 px-2.5 py-1 rounded-full border border-yellow-500/20">
                  <RefreshCw className="h-3.5 w-3.5 animate-spin" />
                  Running
                </span>
              ) : data.lastConclusion === "failure" ? (
                <span className="inline-flex items-center gap-1.5 text-sm font-semibold text-red-500 bg-red-500/10 px-2.5 py-1 rounded-full border border-red-500/20">
                  <XCircle className="h-3.5 w-3.5" />
                  Failed
                </span>
              ) : (
                <span className="inline-flex items-center gap-1.5 text-sm font-semibold text-emerald-500 bg-emerald-500/10 px-2.5 py-1 rounded-full border border-emerald-500/20">
                  <CheckCircle2 className="h-3.5 w-3.5" />
                  Passing
                </span>
              )}
            </div>
          </div>

          <div className="rounded-xl border border-border bg-card p-5 space-y-2">
            <div className="flex items-center justify-between text-muted-foreground text-xs font-mono uppercase tracking-wider">
              <span>ETA Completion</span>
              <Clock className="h-4 w-4 text-primary" />
            </div>
            <div className="text-2xl font-bold text-foreground font-mono">
              {data.isRunning ? `${data.etaSeconds}s` : "0s (Idle)"}
            </div>
          </div>

          <div className="rounded-xl border border-border bg-card p-5 space-y-2">
            <div className="flex items-center justify-between text-muted-foreground text-xs font-mono uppercase tracking-wider">
              <span>Latest Release</span>
              <Tag className="h-4 w-4 text-primary" />
            </div>
            <div className="text-2xl font-bold text-foreground font-mono">
              {data.lastTagRelease}
            </div>
          </div>

          <div className="rounded-xl border border-border bg-card p-5 space-y-2">
            <div className="flex items-center justify-between text-muted-foreground text-xs font-mono uppercase tracking-wider">
              <span>Pending PRs</span>
              <GitPullRequest className="h-4 w-4 text-primary" />
            </div>
            <div className="text-2xl font-bold text-foreground font-mono">
              {data.pendingPRs}
            </div>
          </div>
        </div>

        {/* CLI Usage Commands */}
        <div className="rounded-xl border border-border bg-card p-6 space-y-4">
          <h2 className="text-lg font-heading font-semibold text-foreground">CLI Command Reference</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-3 font-mono text-xs">
            <div className="p-3 rounded-lg border border-border bg-muted/40 space-y-1">
              <div className="text-primary font-semibold">gitmap pipeline status</div>
              <div className="text-muted-foreground text-[11px]">View active workflow status, ETA, and pending PRs.</div>
            </div>
            <div className="p-3 rounded-lg border border-border bg-muted/40 space-y-1">
              <div className="text-primary font-semibold">gitmap pipeline waittime</div>
              <div className="text-muted-foreground text-[11px]">Output ETA seconds directly to stdout for scripts/agents.</div>
            </div>
            <div className="p-3 rounded-lg border border-border bg-muted/40 space-y-1">
              <div className="text-primary font-semibold">gitmap pipeline error-logs --json</div>
              <div className="text-muted-foreground text-[11px]">Fetch step-by-step failure output as structured JSON.</div>
            </div>
            <div className="p-3 rounded-lg border border-border bg-muted/40 space-y-1">
              <div className="text-primary font-semibold">gitmap pipeline error-logs --tempfile ci-err.json</div>
              <div className="text-muted-foreground text-[11px]">Save error log into .lovable/temp/ for LLM auto-fixes.</div>
            </div>
          </div>
        </div>
      </div>
    </DocsLayout>
  );
}
