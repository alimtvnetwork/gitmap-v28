import { Zap } from "lucide-react";

interface PullCacheCardProps {
  pullCacheWindow: string;
  onPullCacheWindowChange: (value: string) => void;
}

export function PullCacheCard({
  pullCacheWindow,
  onPullCacheWindowChange,
}: PullCacheCardProps) {
  return (
    <div className="rounded-xl border border-border bg-card p-6 shadow-sm space-y-4">
      <div className="flex items-center gap-2 text-foreground font-semibold text-lg">
        <Zap className="h-5 w-5 text-primary" />
        <h2>Git Pull Telemetry Cache</h2>
      </div>
      <p className="text-sm text-muted-foreground">
        Configure the TTL window in seconds for cached repository status during{" "}
        <code className="font-mono text-xs bg-muted px-1.5 py-0.5 rounded text-foreground">
          gitmap pull
        </code>
        . If a pull was executed within this window, up-to-date status and uncommitted files are served instantly from cache.
      </p>
      <div className="space-y-2">
        <label
          htmlFor="pull-cache-window-input"
          className="text-xs font-mono uppercase text-muted-foreground tracking-wider"
        >
          Pull Cache Window (Seconds)
        </label>
        <input
          id="pull-cache-window-input"
          type="number"
          value={pullCacheWindow}
          onChange={(e) => onPullCacheWindowChange(e.target.value)}
          className="w-full px-3 py-2 rounded-md border border-input bg-background text-foreground text-sm font-mono focus:outline-none focus:ring-2 focus:ring-primary"
          min="0"
          max="60"
        />
        <p className="text-xs text-muted-foreground">
          Default is 3 seconds (min: 0 to disable cache, max: 60s).
        </p>
      </div>
    </div>
  );
}

export default PullCacheCard;
