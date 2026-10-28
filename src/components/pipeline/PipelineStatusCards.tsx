import { Activity, Clock, CheckCircle2, XCircle, RefreshCw, Tag, GitPullRequest } from "lucide-react";

export interface PipelineMetrics {
  isRunning: boolean;
  etaSeconds: number;
  lastTagRelease: string;
  pendingPRs: number;
  lastConclusion?: string;
}

interface PipelineStatusCardsProps {
  data: PipelineMetrics;
}

export const PipelineStatusCards = ({ data }: PipelineStatusCardsProps) => {
  return (
    <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
      <div className="rounded-xl border-2 border-primary/70 bg-card p-5 space-y-2 shadow-sm shadow-primary/10 transform-gpu hover:-translate-y-0.5 transition-transform">
        <div className="flex items-center justify-between text-primary text-xs font-mono uppercase tracking-wider font-semibold">
          <span>Pipeline Status</span>
          <Activity className="h-4 w-4 text-primary" />
        </div>
        <div className="flex items-center gap-2">
          {data.isRunning ? (
            <span className="inline-flex items-center gap-1.5 text-sm font-semibold text-amber-800 dark:text-amber-300 bg-amber-500/10 px-2.5 py-1 rounded-full border border-amber-500/30">
              <RefreshCw className="h-3.5 w-3.5 animate-spin" />
              Running
            </span>
          ) : data.lastConclusion === "failure" ? (
            <span className="inline-flex items-center gap-1.5 text-sm font-semibold text-rose-700 dark:text-rose-400 bg-rose-500/10 px-2.5 py-1 rounded-full border border-rose-500/30">
              <XCircle className="h-3.5 w-3.5" />
              Failed
            </span>
          ) : (
            <span className="inline-flex items-center gap-1.5 text-sm font-semibold text-emerald-700 dark:text-emerald-400 bg-emerald-500/10 px-2.5 py-1 rounded-full border border-emerald-500/30">
              <CheckCircle2 className="h-3.5 w-3.5" />
              Passing
            </span>
          )}
        </div>
      </div>

      <div className="rounded-xl border border-border/60 bg-card p-5 space-y-2 shadow-sm transform-gpu hover:-translate-y-0.5 transition-transform">
        <div className="flex items-center justify-between text-muted-foreground text-xs font-mono uppercase tracking-wider">
          <span>ETA Completion</span>
          <Clock className="h-4 w-4 text-muted-foreground" />
        </div>
        <div className="text-2xl font-bold text-foreground font-mono">
          {data.isRunning ? `${data.etaSeconds}s` : "0s (Idle)"}
        </div>
      </div>

      <div className="rounded-xl border border-border/60 bg-card p-5 space-y-2 shadow-sm transform-gpu hover:-translate-y-0.5 transition-transform">
        <div className="flex items-center justify-between text-muted-foreground text-xs font-mono uppercase tracking-wider">
          <span>Latest Release</span>
          <Tag className="h-4 w-4 text-muted-foreground" />
        </div>
        <div className="text-2xl font-bold text-foreground font-mono">
          {data.lastTagRelease}
        </div>
      </div>

      <div className="rounded-xl border border-border/60 bg-card p-5 space-y-2 shadow-sm transform-gpu hover:-translate-y-0.5 transition-transform">
        <div className="flex items-center justify-between text-muted-foreground text-xs font-mono uppercase tracking-wider">
          <span>Pending PRs</span>
          <GitPullRequest className="h-4 w-4 text-muted-foreground" />
        </div>
        <div className="text-2xl font-bold text-foreground font-mono">
          {data.pendingPRs}
        </div>
      </div>
    </div>
  );
};

export default PipelineStatusCards;
