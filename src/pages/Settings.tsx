import { useState, useEffect } from "react";
import DocsLayout from "@/components/docs/DocsLayout";
import { Settings as SettingsIcon, Save, Check } from "lucide-react";
import { useToast } from "@/hooks/use-toast";
import { queryWrapperSync } from "@/lib/queryWrapper";
import { TerminalThemeType } from "@/components/settings/types";
import { StorageTempCard } from "@/components/settings/StorageTempCard";
import { PipelineMonitorCard } from "@/components/settings/PipelineMonitorCard";
import { TerminalThemeCard } from "@/components/settings/TerminalThemeCard";
import { GitHubAuthCard } from "@/components/settings/GitHubAuthCard";

export default function SettingsPage() {
  const { toast } = useToast();
  const [tempDir, setTempDir] = useState(".lovable/temp");
  const [pollInterval, setPollInterval] = useState("10");
  const [terminalTheme, setTerminalTheme] = useState(TerminalThemeType.Dark);
  const [isSaved, setIsSaved] = useState(false);

  useEffect(() => {
    const res = queryWrapperSync(() => ({
      temp: localStorage.getItem("gitmap_temp_dir"),
      poll: localStorage.getItem("gitmap_poll_interval"),
      theme: localStorage.getItem("gitmap_terminal_theme"),
    }));

    if (res.isSuccess && res.data) {
      if (res.data.temp) setTempDir(res.data.temp);
      if (res.data.poll) setPollInterval(res.data.poll);
      if (res.data.theme) setTerminalTheme(res.data.theme as TerminalThemeType);
    }
  }, []);

  const handleSave = () => {
    const res = queryWrapperSync(() => {
      localStorage.setItem("gitmap_temp_dir", tempDir);
      localStorage.setItem("gitmap_poll_interval", pollInterval);
      localStorage.setItem("gitmap_terminal_theme", terminalTheme);

      return true;
    });

    if (res.isFailed) {
      toast({ title: "Failed to Save", description: "Storage write failed.", variant: "destructive" });

      return;
    }

    setIsSaved(true);
    toast({ title: "Settings Saved", description: "Preferences updated and synced." });
    setTimeout(() => setIsSaved(false), 2000);
  };

  return (
    <DocsLayout>
      <div className="max-w-4xl mx-auto space-y-8 pb-16">
        <div className="flex items-center gap-3 mb-2">
          <div className="h-10 w-10 rounded-lg bg-primary/10 flex items-center justify-center dark:bg-primary/20">
            <SettingsIcon className="h-5 w-5 text-primary" />
          </div>
          <div>
            <h1 className="text-3xl font-heading font-bold text-foreground">Settings & Preferences</h1>
            <p className="text-sm text-muted-foreground">Manage temp folder locations, pipeline monitoring, and UI preferences.</p>
          </div>
        </div>

        <div className="grid gap-6">
          <StorageTempCard tempDir={tempDir} onTempDirChange={setTempDir} />
          <PipelineMonitorCard pollInterval={pollInterval} onPollIntervalChange={setPollInterval} />
          <TerminalThemeCard terminalTheme={terminalTheme} onTerminalThemeChange={setTerminalTheme} />
          <GitHubAuthCard />
          <div className="flex justify-end pt-2">
            <button
              onClick={handleSave}
              type="button"
              className="inline-flex items-center gap-2 px-6 py-2.5 rounded-lg bg-primary text-primary-foreground font-semibold text-sm shadow hover:bg-primary/90 transition-colors"
            >
              {isSaved ? <Check className="h-4 w-4" /> : <Save className="h-4 w-4" />}
              {isSaved ? "Saved" : "Save Changes"}
            </button>
          </div>
        </div>
      </div>
    </DocsLayout>
  );
}
