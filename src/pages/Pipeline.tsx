import { useState, useEffect } from "react";
import DocsLayout from "@/components/docs/DocsLayout";
import { Activity, ExternalLink, Download, Copy, Check } from "lucide-react";
import { copyToClipboard } from "@/lib/clipboard";
import { useToast } from "@/hooks/use-toast";
import PipelineStatusCards from "@/components/pipeline/PipelineStatusCards";

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
              <h1 className="text-3xl font-heading font-bold tracking-[-0.02em] text-foreground">CI/CD Pipeline Monitor</h1>
              <p className="text-sm text-muted-foreground">Real-time status, ETA completion estimation, and failure logs.</p>
            </div>
          </div>
          <div className="flex items-center gap-2">
            <button
              type="button"
              onClick={handleCopyJSON}
              className="inline-flex items-center gap-1.5 px-3 py-1.5 min-h-[36px] rounded-md border border-border bg-card text-xs font-mono text-foreground hover:bg-muted transition-all active:scale-[0.98]"
            >
              {isCopied ? <Check className="h-3.5 w-3.5 text-primary" /> : <Copy className="h-3.5 w-3.5" />}
              {isCopied ? "Copied" : "Copy Status JSON"}
            </button>
          </div>
        </div>

        {/* Status Metrics Cards */}
        <PipelineStatusCards data={data} />

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
              <div className="text-primary font-semibold">gitmap pipeline errors agy fix</div>
              <div className="text-muted-foreground text-[11px]">Directly inject pipeline error logs and RCA prompt into Antigravity IDE.</div>
            </div>
            <div className="p-3 rounded-lg border border-border bg-muted/40 space-y-1">
              <div className="text-primary font-semibold">gitmap pipeline errors agy fix --all</div>
              <div className="text-muted-foreground text-[11px]">Parallel scan all repos and batch-fix failing pipelines (3 projects/run).</div>
            </div>
          </div>
        </div>
      </div>
    </DocsLayout>
  );
}
