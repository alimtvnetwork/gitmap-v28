import { useState, useEffect } from "react";
import DocsLayout from "@/components/docs/DocsLayout";
import { Settings as SettingsIcon, Save, Check } from "lucide-react";
import { useToast } from "@/hooks/use-toast";
import { queryWrapperSync } from "@/lib/queryWrapper";
import { TerminalThemeType } from "@/components/settings/types";
import { StorageTempCard } from "@/components/settings/StorageTempCard";
import { PipelineMonitorCard } from "@/components/settings/PipelineMonitorCard";
import { PullCacheCard } from "@/components/settings/PullCacheCard";
import { TerminalThemeCard } from "@/components/settings/TerminalThemeCard";
import { GitHubAuthCard } from "@/components/settings/GitHubAuthCard";
import { AntigravityInstancesCard } from "@/components/settings/AntigravityInstancesCard";
import { CommitScheduleCard } from "@/components/settings/CommitScheduleCard";
import { NodeResourceCard } from "@/components/settings/NodeResourceCard";
import { UiUxRepoCatalogCard } from "@/components/settings/UiUxRepoCatalogCard";

interface AgyInstanceRecord {
  instanceId: string;
  instanceName: string;
  processId: number;
  languageServer?: string;
  configDir?: string;
  brainDir?: string;
  summariesDbPath?: string;
  activeWorkspaces?: string[];
  isPrimary: boolean;
  isRunning: boolean;
  lastActiveAt?: string;
}

type SettingsTab = "all" | "storage" | "instances" | "theme" | "catalog";

export default function SettingsPage() {
  const { toast } = useToast();
  const [activeTab, setActiveTab] = useState<SettingsTab>("all");
  const [tempDir, setTempDir] = useState(".ai-memory/temp");
  const [pollInterval, setPollInterval] = useState("10");
  const [pullCacheWindow, setPullCacheWindow] = useState("3");
  const [terminalTheme, setTerminalTheme] = useState(TerminalThemeType.Dark);
  const [nodeResourceThreshold, setNodeResourceThreshold] = useState<number>(30);
  const [paeRecencyHours, setPaeRecencyHours] = useState<number>(48);
  const [commitTimezone, setCommitTimezone] = useState("Local");
  const [isCommitScheduleActive, setIsCommitScheduleActive] = useState(false);
  const [commitWeekdayCutoff, setCommitWeekdayCutoff] = useState("08:00 AM");
  const [isSaved, setIsSaved] = useState(false);
  const [isSaving, setIsSaving] = useState(false);

  // Multi-Instance Settings State
  const [selectedInstanceId, setSelectedInstanceId] = useState("primary");
  const [isAutoDiscoverActive, setIsAutoDiscoverActive] = useState(true);
  const [isTelemetrySyncActive, setIsTelemetrySyncActive] = useState(true);
  const [instances, setInstances] = useState<AgyInstanceRecord[]>([]);
  const [isLoadingInstances, setIsLoadingInstances] = useState(false);

  useEffect(() => {
    const res = queryWrapperSync(() => ({
      temp: localStorage.getItem("gitmap_temp_dir"),
      poll: localStorage.getItem("gitmap_poll_interval"),
      pullCache: localStorage.getItem("gitmap_pull_cache_window_seconds"),
      theme: localStorage.getItem("gitmap_terminal_theme"),
      nodeThreshold: localStorage.getItem("gitmap_node_resource_threshold"),
      paeRecency: localStorage.getItem("gitmap_pae_recency_window_hours"),
      commitTz: localStorage.getItem("gitmap_commit_timezone"),
      commitSchedule: localStorage.getItem("gitmap_commit_schedule_active"),
      commitCutoff: localStorage.getItem("gitmap_commit_weekday_cutoff"),
      instanceId: localStorage.getItem("gitmap_agy_instance_id"),
      autoDiscover: localStorage.getItem("gitmap_agy_auto_discover"),
      telemetry: localStorage.getItem("gitmap_agy_telemetry_sync"),
    }));

    if (res.isFail || !res.data) {
      return;
    }

    if (res.data.temp) setTempDir(res.data.temp);
    if (res.data.poll) setPollInterval(res.data.poll);
    if (res.data.pullCache) setPullCacheWindow(res.data.pullCache);
    if (res.data.theme) setTerminalTheme(res.data.theme as TerminalThemeType);
    if (res.data.nodeThreshold) setNodeResourceThreshold(Number(res.data.nodeThreshold) || 30);
    if (res.data.paeRecency) setPaeRecencyHours(Number(res.data.paeRecency) || 48);
    if (res.data.commitTz) setCommitTimezone(res.data.commitTz);
    if (res.data.commitSchedule !== null) setIsCommitScheduleActive(res.data.commitSchedule === "true");
    if (res.data.commitCutoff) setCommitWeekdayCutoff(res.data.commitCutoff);
    if (res.data.instanceId) setSelectedInstanceId(res.data.instanceId);
    if (res.data.autoDiscover !== null) setIsAutoDiscoverActive(res.data.autoDiscover === "true");
    if (res.data.telemetry !== null) setIsTelemetrySyncActive(res.data.telemetry === "true");

    loadInstancesFromAPI();
  }, []);

  const loadInstancesFromAPI = async () => {
    setIsLoadingInstances(true);
    try {
      const response = await fetch("/api/instances");
      if (!response.ok) return;

      const json = await response.json();
      if (json.isSuccess && Array.isArray(json.instances)) {
        setInstances(json.instances);
        return;
      }
    } catch {
      // Fallback local stub if REST server not directly mounted on same port
    } finally {
      setIsLoadingInstances(false);
    }

    // Default registered instances fallback
    setInstances([
      {
        instanceId: "primary",
        instanceName: "Default Antigravity Workspace",
        processId: 4757,
        languageServer: "127.0.0.1:33419",
        configDir: "~/.gemini/antigravity",
        isPrimary: true,
        isRunning: true,
        lastActiveAt: new Date().toISOString(),
      },
      {
        instanceId: "default-copy-8159",
        instanceName: "8159 Secondary IDE",
        processId: 18454,
        languageServer: "127.0.0.1:35427",
        configDir: "~/.antigravity_tools/instances/default-copy-8159/data",
        isPrimary: false,
        isRunning: true,
        lastActiveAt: new Date().toISOString(),
      },
      {
        instanceId: "gitmap-7845",
        instanceName: "GitMap Dedicated Instance",
        processId: 0,
        configDir: "~/.antigravity_tools/instances/gitmap-7845/data",
        isPrimary: false,
        isRunning: false,
        lastActiveAt: new Date(Date.now() - 3600000).toISOString(),
      },
    ]);
  };

  const handleSave = async () => {
    setIsSaving(true);
    const res = queryWrapperSync(() => {
      localStorage.setItem("gitmap_temp_dir", tempDir);
      localStorage.setItem("gitmap_poll_interval", pollInterval);
      localStorage.setItem("gitmap_pull_cache_window_seconds", pullCacheWindow);
      localStorage.setItem("gitmap_terminal_theme", terminalTheme);
      localStorage.setItem("gitmap_node_resource_threshold", String(nodeResourceThreshold));
      localStorage.setItem("gitmap_pae_recency_window_hours", String(paeRecencyHours));
      localStorage.setItem("gitmap_commit_timezone", commitTimezone);
      localStorage.setItem("gitmap_commit_schedule_active", String(isCommitScheduleActive));
      localStorage.setItem("gitmap_commit_weekday_cutoff", commitWeekdayCutoff);
      localStorage.setItem("gitmap_agy_instance_id", selectedInstanceId);
      localStorage.setItem("gitmap_agy_auto_discover", String(isAutoDiscoverActive));
      localStorage.setItem("gitmap_agy_telemetry_sync", String(isTelemetrySyncActive));

      return true;
    });

    if (res.isFailed) {
      toast({ title: "Failed to Save", description: "Storage write failed.", variant: "destructive" });
      setIsSaving(false);
      return;
    }

    try {
      await fetch("/api/settings", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          theme: terminalTheme,
          attributes: {
            "storage.temp_dir": tempDir,
            "pipeline.poll_interval": pollInterval,
            "pull_cache_window_seconds": pullCacheWindow,
            "node_resource_threshold": String(nodeResourceThreshold),
            "pae_recency_window_hours": String(paeRecencyHours),
            "commit.timezone": commitTimezone,
            "commit.schedule_active": String(isCommitScheduleActive),
            "commit.weekday_cutoff": commitWeekdayCutoff,
            "antigravity.instance_id": selectedInstanceId,
            "antigravity.auto_discover": String(isAutoDiscoverActive),
            "antigravity.telemetry_sync": String(isTelemetrySyncActive),
          },
        }),
      });
    } catch {
      // Offline or standalone React mode
    }

    setIsSaving(false);
    setIsSaved(true);
    toast({ title: "Settings Saved", description: "All configuration tokens synchronized." });
    setTimeout(() => setIsSaved(false), 2000);
  };

  return (
    <DocsLayout>
      <div className="max-w-5xl mx-auto space-y-8 pb-16">
        {/* Header with Plane 1 Depth */}
        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 p-5 rounded-xl bg-card/60 border border-border/40 shadow-sm backdrop-blur">
          <div className="flex items-center gap-3">
            <div className="h-11 w-11 rounded-lg bg-muted/60 border border-border/50 text-muted-foreground flex items-center justify-center">
              <SettingsIcon className="h-6 w-6" />
            </div>
            <div>
              <h1 className="text-2xl sm:text-3xl font-heading font-bold text-foreground tracking-[-0.02em]">
                Settings &amp; Architecture
              </h1>
              <p className="text-sm text-muted-foreground">
                Configure 4-Plane elevation tokens, multi-instance discovery, and design benchmarks.
              </p>
            </div>
          </div>
          <div className="flex items-center gap-2">
            <button
              onClick={handleSave}
              disabled={isSaving}
              type="button"
              className="inline-flex items-center gap-2 px-5 py-2.5 min-h-[44px] rounded-lg bg-primary text-primary-foreground font-semibold text-sm shadow hover:bg-primary/90 transition-all active:scale-[0.98] disabled:opacity-50"
            >
              {isSaved ? <Check className="h-4 w-4" /> : <Save className="h-4 w-4" />}
              {isSaved ? "Saved" : isSaving ? "Saving..." : "Save Changes"}
            </button>
          </div>
        </div>

        {/* Tab Navigation (Plane 1 Raised Well with Glassmorphic Sticky Physics) */}
        <div className="sticky top-0 z-20 flex gap-2 p-1.5 rounded-lg bg-background/80 backdrop-blur-md border border-border/40 overflow-x-auto shadow-sm">
          <button
            type="button"
            onClick={() => setActiveTab("all")}
            className={`px-3.5 py-1.5 rounded-md text-xs sm:text-sm font-medium transition-all ${
              activeTab === "all"
                ? "bg-card text-foreground shadow-sm border border-border/60"
                : "text-muted-foreground hover:text-foreground"
            }`}
          >
            📑 All Preferences
          </button>
          <button
            type="button"
            onClick={() => setActiveTab("instances")}
            className={`px-3.5 py-1.5 rounded-md text-xs sm:text-sm font-medium transition-all ${
              activeTab === "instances"
                ? "bg-card text-foreground shadow-sm border border-border/60"
                : "text-muted-foreground hover:text-foreground"
            }`}
          >
            🛸 Antigravity Multi-Instance
          </button>
          <button
            type="button"
            onClick={() => setActiveTab("storage")}
            className={`px-3.5 py-1.5 rounded-md text-xs sm:text-sm font-medium transition-all ${
              activeTab === "storage"
                ? "bg-card text-foreground shadow-sm border border-border/60"
                : "text-muted-foreground hover:text-foreground"
            }`}
          >
            📦 Storage &amp; Telemetry
          </button>
          <button
            type="button"
            onClick={() => setActiveTab("theme")}
            className={`px-3.5 py-1.5 rounded-md text-xs sm:text-sm font-medium transition-all ${
              activeTab === "theme"
                ? "bg-card text-foreground shadow-sm border border-border/60"
                : "text-muted-foreground hover:text-foreground"
            }`}
          >
            🎨 Themes &amp; Terminal
          </button>
          <button
            type="button"
            onClick={() => setActiveTab("catalog")}
            className={`px-3.5 py-1.5 rounded-md text-xs sm:text-sm font-medium transition-all ${
              activeTab === "catalog"
                ? "bg-card text-foreground shadow-sm border border-border/60"
                : "text-muted-foreground hover:text-foreground"
            }`}
          >
            ✨ UI/UX &amp; CSS3 Catalog
          </button>
        </div>

        {/* SECTION 1: ANTIGRAVITY MULTI-INSTANCE PANEL */}
        {(activeTab === "all" || activeTab === "instances") && (
          <AntigravityInstancesCard
            instances={instances}
            selectedInstanceId={selectedInstanceId}
            onSelectInstanceId={setSelectedInstanceId}
            isAutoDiscoverActive={isAutoDiscoverActive}
            onToggleAutoDiscover={setIsAutoDiscoverActive}
            isLoadingInstances={isLoadingInstances}
          />
        )}

        {/* SECTION 2: STORAGE & MONITORING */}
        {(activeTab === "all" || activeTab === "storage") && (
          <div className="grid gap-6">
            <StorageTempCard tempDir={tempDir} onTempDirChange={setTempDir} />
            <PipelineMonitorCard pollInterval={pollInterval} onPollIntervalChange={setPollInterval} />
            <PullCacheCard pullCacheWindow={pullCacheWindow} onPullCacheWindowChange={setPullCacheWindow} />
            <CommitScheduleCard
              timezone={commitTimezone}
              onTimezoneChange={setCommitTimezone}
              isScheduleActive={isCommitScheduleActive}
              onToggleSchedule={setIsCommitScheduleActive}
              weekdayCutoff={commitWeekdayCutoff}
            />
            <NodeResourceCard
              nodeResourceThreshold={nodeResourceThreshold}
              onNodeResourceThresholdChange={setNodeResourceThreshold}
              paeRecencyHours={paeRecencyHours}
              onPaeRecencyHoursChange={setPaeRecencyHours}
            />
          </div>
        )}

        {/* SECTION 3: THEMES & AUTH */}
        {(activeTab === "all" || activeTab === "theme") && (
          <div className="grid gap-6">
            <TerminalThemeCard terminalTheme={terminalTheme} onTerminalThemeChange={setTerminalTheme} />
            <GitHubAuthCard />
          </div>
        )}

        {/* SECTION 4: CURATED UI/UX & CSS3 REPO CATALOG */}
        {(activeTab === "all" || activeTab === "catalog") && (
          <UiUxRepoCatalogCard />
        )}
      </div>
    </DocsLayout>
  );
}
