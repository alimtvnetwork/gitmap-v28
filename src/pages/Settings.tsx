import { useState, useEffect } from "react";
import DocsLayout from "@/components/docs/DocsLayout";
import {
  Settings as SettingsIcon,
  Save,
  Check,
  Server,
  Layers,
  Palette,
  ExternalLink,
  Cpu,
  Sparkles,
  ShieldCheck,
  BookmarkCheck,
  FolderGit2,
} from "lucide-react";
import { useToast } from "@/hooks/use-toast";
import { queryWrapperSync } from "@/lib/queryWrapper";
import { TerminalThemeType } from "@/components/settings/types";
import { StorageTempCard } from "@/components/settings/StorageTempCard";
import { PipelineMonitorCard } from "@/components/settings/PipelineMonitorCard";
import { PullCacheCard } from "@/components/settings/PullCacheCard";
import { TerminalThemeCard } from "@/components/settings/TerminalThemeCard";
import { GitHubAuthCard } from "@/components/settings/GitHubAuthCard";
import { AntigravityInstancesCard } from "@/components/settings/AntigravityInstancesCard";

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

interface ReferenceRepoItem {
  name: string;
  repo: string;
  url: string;
  category: string;
  description: string;
  badge: string;
}

const referenceRepositories: ReferenceRepoItem[] = [
  {
    name: "shadcn/ui",
    repo: "shadcn-ui/ui",
    url: "https://github.com/shadcn-ui/ui",
    category: "Design System & Accessible Primitives",
    description: "Copy-paste accessible React components built on Radix UI and Tailwind CSS with first-class HSL token architecture.",
    badge: "Radix + Tailwind",
  },
  {
    name: "Radix Primitives",
    repo: "radix-ui/primitives",
    url: "https://github.com/radix-ui/primitives",
    category: "Headless Component Architecture",
    description: "Unstyled, fully accessible UI components (dialogs, tooltips, dropdowns, popovers) with WAI-ARIA compliance.",
    badge: "WAI-ARIA Standard",
  },
  {
    name: "Aceternity UI",
    repo: "mannupaaji/aceternity-ui",
    url: "https://github.com/mannupaaji/aceternity-ui",
    category: "Visual Craft & Canvas Animations",
    description: "High-craft modern interactive components, border beams, glowing cards, and canvas visual effects.",
    badge: "Framer Motion",
  },
  {
    name: "Magic UI",
    repo: "magicuidesign/magicui",
    url: "https://github.com/magicuidesign/magicui",
    category: "Micro-Interactions & Landing UI",
    description: "20+ copy-paste components for high-conversion landing pages: animated beams, retro grids, and particle UI.",
    badge: "Creative Engineering",
  },
  {
    name: "Tremor",
    repo: "tremorlabs/tremor",
    url: "https://github.com/tremorlabs/tremor",
    category: "Data Telemetry & Dashboards",
    description: "Modular dashboard and telemetry components with clean dark mode neutral depth metrics cards and charts.",
    badge: "Analytics UI",
  },
  {
    name: "Lucide Icons",
    repo: "lucide-icons/lucide",
    url: "https://github.com/lucide-icons/lucide",
    category: "Vector Iconography",
    description: "Consistent, clean vector iconography system with uniform 24px viewports and zero runtime overhead.",
    badge: "Vector Icons",
  },
  {
    name: "Modern CSS Solutions",
    repo: "moderncss.dev",
    url: "https://moderncss.dev",
    category: "CSS3 Architecture & Layout Patterns",
    description: "Deep-dive CSS3 architectural recipes for modern grid, flexbox, fluid typography, and accessible form styling.",
    badge: "CSS3 Standard",
  },
  {
    name: "CodyHouse Framework",
    repo: "CodyHouse/codyhouse-framework",
    url: "https://github.com/CodyHouse/codyhouse-framework",
    category: "Token Architecture & Typography Scales",
    description: "Production-ready design system tokens, typography scales, spacing units, and modular CSS utility classes.",
    badge: "Tokens & Scales",
  },
  {
    name: "CSS-Tricks Guides",
    repo: "css-tricks.com",
    url: "https://css-tricks.com",
    category: "Subgrid & Custom Properties",
    description: "Authoritative guides on fluid typography, CSS custom properties, responsive subgrid layouts, and micro-animations.",
    badge: "Layout Guide",
  },
];

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
      instanceId: localStorage.getItem("gitmap_agy_instance_id"),
      autoDiscover: localStorage.getItem("gitmap_agy_auto_discover"),
      telemetry: localStorage.getItem("gitmap_agy_telemetry_sync"),
    }));

    if (res.isFail || !res.data) {
      return;
    }

    if (res.data.temp) {
      setTempDir(res.data.temp);
    }

    if (res.data.poll) {
      setPollInterval(res.data.poll);
    }

    if (res.data.pullCache) {
      setPullCacheWindow(res.data.pullCache);
    }

    if (res.data.theme) {
      setTerminalTheme(res.data.theme as TerminalThemeType);
    }

    if (res.data.nodeThreshold) {
      setNodeResourceThreshold(Number(res.data.nodeThreshold) || 30);
    }

    if (res.data.paeRecency) {
      setPaeRecencyHours(Number(res.data.paeRecency) || 48);
    }

    if (res.data.instanceId) {
      setSelectedInstanceId(res.data.instanceId);
    }

    if (res.data.autoDiscover !== null) {
      setIsAutoDiscoverActive(res.data.autoDiscover === "true");
    }

    if (res.data.telemetry !== null) {
      setIsTelemetrySyncActive(res.data.telemetry === "true");
    }

    loadInstancesFromAPI();
  }, []);

  const loadInstancesFromAPI = async () => {
    setIsLoadingInstances(true);
    try {
      const response = await fetch("/api/instances");
      if (!response.ok) {
        return;
      }

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

    // Bidirectional sync with backend /api/settings
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
            <div className="h-11 w-11 rounded-lg bg-primary/10 border border-primary/20 flex items-center justify-center dark:bg-primary/20">
              <SettingsIcon className="h-6 w-6 text-primary" />
            </div>
            <div>
              <h1 className="text-2xl sm:text-3xl font-heading font-bold text-foreground tracking-tight">
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
              className="inline-flex items-center gap-2 px-5 py-2.5 rounded-lg bg-primary text-primary-foreground font-semibold text-sm shadow hover:bg-primary/90 transition-all active:scale-95 disabled:opacity-50"
            >
              {isSaved ? <Check className="h-4 w-4" /> : <Save className="h-4 w-4" />}
              {isSaved ? "Saved" : isSaving ? "Saving..." : "Save Changes"}
            </button>
          </div>
        </div>

        {/* Tab Navigation (Plane 1 Raised Well) */}
        <div className="flex gap-2 p-1.5 rounded-lg bg-muted/30 border border-border/40 overflow-x-auto">
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

            {/* Node Resource Benchmark Threshold & PAE Commit Recency Card */}
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
                {/* Node Resource Benchmark Threshold */}
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
                    onChange={(e) => setNodeResourceThreshold(Number(e.target.value))}
                    className="w-full accent-primary cursor-pointer"
                  />
                  <p className="text-xs text-muted-foreground">
                    Idle utilization threshold (default: 30%). When fleet nodes report CPU and RAM below this level, adaptive parallel pull concurrency is boosted to maximum capacity.
                  </p>
                </div>

                {/* PAE Commit Recency Window */}
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
                    onChange={(e) => setPaeRecencyHours(Math.max(1, Number(e.target.value)))}
                    className="w-full bg-muted/20 border border-border/60 text-foreground px-3.5 py-2 rounded-lg text-sm focus:outline-none focus:ring-1 focus:ring-primary font-mono"
                  />
                  <p className="text-xs text-muted-foreground">
                    Lookback hours for <code className="text-primary font-mono font-semibold">gitmap pae</code> (default: 48h). Repositories with commits inside this window are prioritized for fast synchronization.
                  </p>
                </div>
              </div>
            </div>
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
          <div className="rounded-xl bg-card border border-border/50 p-6 shadow-sm space-y-6">
            <div className="flex items-center justify-between border-b border-border/40 pb-4">
              <div className="flex items-center gap-2.5">
                <Sparkles className="h-5 w-5 text-primary" />
                <div>
                  <h2 className="text-lg font-heading font-semibold text-foreground">
                    Authoritative UI/UX &amp; CSS3 Repository Index
                  </h2>
                  <p className="text-xs text-muted-foreground">
                    Curated benchmarks for layout, accessibility primitives, tokens, and micro-interactions.
                  </p>
                </div>
              </div>
              <span className="text-xs px-2.5 py-1 rounded-full bg-muted text-muted-foreground border border-border/60 font-medium">
                Design Reference Spec
              </span>
            </div>

            <div className="grid sm:grid-cols-2 lg:grid-cols-3 gap-4">
              {referenceRepositories.map((repo) => (
                <div
                  key={repo.name}
                  className="p-4 rounded-lg bg-muted/15 border border-border/40 hover:border-primary/50 transition-all flex flex-col justify-between group"
                >
                  <div className="space-y-2">
                    <div className="flex items-start justify-between gap-2">
                      <h3 className="font-heading font-bold text-sm text-foreground group-hover:text-primary transition-colors">
                        {repo.name}
                      </h3>
                      <span className="text-[10px] px-2 py-0.5 rounded-md bg-muted text-muted-foreground border border-border/40 font-mono">
                        {repo.badge}
                      </span>
                    </div>
                    <div className="text-[11px] text-primary font-medium tracking-tight">
                      {repo.category}
                    </div>
                    <p className="text-xs text-muted-foreground leading-relaxed">
                      {repo.description}
                    </p>
                  </div>
                  <div className="pt-3 mt-3 border-t border-border/30 flex items-center justify-between text-xs">
                    <span className="text-[11px] text-muted-foreground font-mono truncate max-w-[170px]">
                      {repo.repo}
                    </span>
                    <a
                      href={repo.url}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="inline-flex items-center gap-1 text-primary hover:underline font-semibold text-xs"
                    >
                      <span>Explore</span>
                      <ExternalLink className="h-3 w-3" />
                    </a>
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}
      </div>
    </DocsLayout>
  );
}
