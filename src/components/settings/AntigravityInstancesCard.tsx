import { Server } from "lucide-react";

export interface AntigravityInstance {
  instanceId: string;
  instanceName: string;
  isRunning: boolean;
  processId: number;
  languageServer?: string;
}

interface AntigravityInstancesCardProps {
  instances: AntigravityInstance[];
  selectedInstanceId: string;
  onSelectInstanceId: (id: string) => void;
  isAutoDiscoverActive: boolean;
  onToggleAutoDiscover: (active: boolean) => void;
  isLoadingInstances: boolean;
}

export const AntigravityInstancesCard = ({
  instances,
  selectedInstanceId,
  onSelectInstanceId,
  isAutoDiscoverActive,
  onToggleAutoDiscover,
  isLoadingInstances,
}: AntigravityInstancesCardProps) => {
  return (
    <div className="rounded-xl bg-card border border-border/50 p-6 shadow-sm space-y-5">
      <div className="flex items-center justify-between border-b border-border/40 pb-4">
        <div className="flex items-center gap-2.5">
          <Server className="h-5 w-5 text-primary" />
          <h2 className="text-lg font-heading font-semibold text-foreground">
            Antigravity Multi-Instance Engine
          </h2>
        </div>
        <span className="text-xs px-2.5 py-1 rounded-full bg-muted text-muted-foreground border border-border/60 font-medium">
          Task-06 REST &amp; CLI API
        </span>
      </div>

      <p className="text-sm text-muted-foreground">
        Select default Antigravity IDE instance, monitor language server host bridges, and control multi-process prompt queue aggregation.
      </p>

      <div className="grid sm:grid-cols-2 gap-4">
        <div className="space-y-2">
          <label className="text-xs font-medium text-foreground uppercase tracking-wider">
            Target Active Instance
          </label>
          <select
            value={selectedInstanceId}
            onChange={(e) => onSelectInstanceId(e.target.value)}
            className="w-full bg-muted/20 border border-border/60 text-foreground px-3.5 py-2 rounded-lg text-sm focus:outline-none focus:ring-1 focus:ring-primary"
          >
            <option value="primary">primary (Default Workspace: ~/.gemini/antigravity)</option>
            {instances.map((inst) => (
              <option key={inst.instanceId} value={inst.instanceId}>
                {inst.instanceId} &mdash; {inst.instanceName} {inst.isRunning ? "(Running)" : "(Idle)"}
              </option>
            ))}
          </select>
        </div>

        <div className="space-y-2">
          <label className="text-xs font-medium text-foreground uppercase tracking-wider">
            Process Discovery Engine
          </label>
          <div className="flex items-center justify-between p-2.5 rounded-lg bg-muted/20 border border-border/40">
            <span className="text-sm text-foreground">Auto-scan language server ports</span>
            <input
              type="checkbox"
              checked={isAutoDiscoverActive}
              onChange={(e) => onToggleAutoDiscover(e.target.checked)}
              className="h-4 w-4 rounded accent-primary cursor-pointer"
            />
          </div>
        </div>
      </div>

      <div className="space-y-3 pt-2">
        <div className="flex items-center justify-between">
          <span className="text-xs font-semibold text-muted-foreground uppercase tracking-wider">
            Detected Antigravity Processes ({instances.length})
          </span>
          {isLoadingInstances && (
            <span className="text-xs text-primary animate-pulse">Scanning processes...</span>
          )}
        </div>
        <div className="grid sm:grid-cols-3 gap-3">
          {instances.map((inst) => (
            <div
              key={inst.instanceId}
              className={`p-3.5 rounded-lg border text-left transition-all ${
                selectedInstanceId === inst.instanceId
                  ? "bg-primary/5 border-primary shadow-sm"
                  : "bg-muted/10 border-border/40 hover:border-border/80"
              }`}
            >
              <div className="flex items-center justify-between mb-1.5">
                <span className="text-xs font-bold text-foreground truncate">{inst.instanceName}</span>
                <span
                  className={`text-[10px] px-2 py-0.5 rounded-full font-medium ${
                    inst.isRunning
                      ? "bg-emerald-500/10 text-emerald-700 dark:text-emerald-400 border border-emerald-500/20"
                      : "bg-muted text-muted-foreground border border-border/40"
                  }`}
                >
                  {inst.isRunning ? "RUNNING" : "IDLE"}
                </span>
              </div>
              <div className="text-[11px] text-muted-foreground space-y-0.5 font-mono">
                <div>ID: {inst.instanceId}</div>
                <div>PID: {inst.processId > 0 ? inst.processId : "N/A"}</div>
                {inst.languageServer && <div>Bridge: {inst.languageServer}</div>}
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};

export default AntigravityInstancesCard;
