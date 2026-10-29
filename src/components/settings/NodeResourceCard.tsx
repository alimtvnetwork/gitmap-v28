import { Cpu } from "lucide-react";

interface NodeResourceCardProps {
  nodeResourceThreshold: number;
  onNodeResourceThresholdChange: (val: number) => void;
  paeRecencyHours: number;
  onPaeRecencyHoursChange: (val: number) => void;
}

export function NodeResourceCard({
  nodeResourceThreshold,
  onNodeResourceThresholdChange,
  paeRecencyHours,
  onPaeRecencyHoursChange,
}: NodeResourceCardProps) {
  return (
    <div className="rounded-xl bg-card border border-border/50 p-6 shadow-sm space-y-6">
      <div className="flex items-center justify-between border-b border-border/40 pb-4">
        <div className="flex items-center gap-2.5">
          <Cpu className="h-5 w-5 text-primary" />
          <div>
            <h2 className="text-lg font-heading font-semibold text-foreground">
              Node Resource Benchmark &amp; PAE Commit Recency
            </h2>
            <p className="text-xs text-muted-foreground">
              Configure CPU/RAM idle benchmarks for parallel execution and lookback window for efficient pull operations.
            </p>
          </div>
        </div>
        <span className="text-xs px-2.5 py-1 rounded-full bg-muted text-muted-foreground border border-border/60 font-medium">
          Fleet Engine Tuning
        </span>
      </div>

      <div className="grid sm:grid-cols-2 gap-6">
        <div className="space-y-3">
          <div className="flex items-center justify-between">
            <label className="text-xs font-semibold text-foreground uppercase tracking-wider">
              Node Resource Benchmark Threshold
            </label>
            <span className="text-xs font-mono font-bold px-2 py-0.5 rounded bg-primary/10 text-primary">
              {nodeResourceThreshold}%
            </span>
          </div>
          <input
            type="range"
            min="10"
            max="90"
            step="5"
            value={nodeResourceThreshold}
            onChange={(e) => onNodeResourceThresholdChange(Number(e.target.value))}
            className="w-full accent-primary cursor-pointer"
          />
          <p className="text-xs text-muted-foreground">
            Idle utilization threshold (default: 30%). When fleet nodes report CPU and RAM below this level, adaptive parallel pull concurrency is boosted to maximum capacity.
          </p>
        </div>

        <div className="space-y-3">
          <div className="flex items-center justify-between">
            <label className="text-xs font-semibold text-foreground uppercase tracking-wider">
              PAE Commit Recency Window
            </label>
            <span className="text-xs font-mono font-bold px-2 py-0.5 rounded bg-primary/10 text-primary">
              {paeRecencyHours} hours
            </span>
          </div>
          <input
            type="number"
            min="1"
            max="720"
            value={paeRecencyHours}
            onChange={(e) => onPaeRecencyHoursChange(Math.max(1, Number(e.target.value)))}
            className="w-full bg-muted/20 border border-border/60 text-foreground px-3.5 py-2 rounded-lg text-sm focus:outline-none focus:ring-1 focus:ring-primary font-mono"
          />
          <p className="text-xs text-muted-foreground">
            Lookback hours for <code className="text-primary font-mono font-semibold">gitmap pae</code> (default: 48h). Repositories with commits inside this window are prioritized for fast synchronization.
          </p>
        </div>
      </div>
    </div>
  );
}
