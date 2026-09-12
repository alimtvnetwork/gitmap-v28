import { RefreshCw } from "lucide-react";

interface PipelineMonitorCardProps {
  pollInterval: string;
  onPollIntervalChange: (value: string) => void;
}

export function PipelineMonitorCard({
  pollInterval,
  onPollIntervalChange,
}: PipelineMonitorCardProps) {
  return (
    <div className="rounded-xl border border-border bg-card p-6 shadow-sm space-y-4">
      <div className="flex items-center gap-2 text-foreground font-semibold text-lg">
        <RefreshCw className="h-5 w-5 text-primary" />
        <h2>CI/CD & Pipeline Monitor</h2>
      </div>
      <p className="text-sm text-muted-foreground">
        Configure automatic background polling frequency for{" "}
        <code className="font-mono text-xs bg-muted px-1.5 py-0.5 rounded text-foreground">
          gitmap pipeline status
        </code>
        .
      </p>
      <div className="space-y-2">
        <label
          htmlFor="poll-interval-input"
          className="text-xs font-mono uppercase text-muted-foreground tracking-wider"
        >
          Status Refresh Interval (Seconds)
        </label>
        <input
          id="poll-interval-input"
          type="number"
          value={pollInterval}
          onChange={(e) => onPollIntervalChange(e.target.value)}
          className="w-full px-3 py-2 rounded-md border border-input bg-background text-foreground text-sm font-mono focus:outline-none focus:ring-2 focus:ring-primary"
          min="5"
          max="120"
        />
      </div>
    </div>
  );
}

export default PipelineMonitorCard;
