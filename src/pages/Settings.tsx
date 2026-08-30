import { useState, useEffect } from "react";
import DocsLayout from "@/components/docs/DocsLayout";
import { Settings as SettingsIcon, Folder, RefreshCw, Terminal, Github, Save, Check, ShieldCheck } from "lucide-react";
import { useToast } from "@/hooks/use-toast";

export default function SettingsPage() {
  const { toast } = useToast();
  const [tempDir, setTempDir] = useState(".lovable/temp");
  const [pollInterval, setPollInterval] = useState("10");
  const [terminalTheme, setTerminalTheme] = useState("dark");
  const [isSaved, setIsSaved] = useState(false);

  useEffect(() => {
    const savedTemp = localStorage.getItem("gitmap_temp_dir");
    const savedPoll = localStorage.getItem("gitmap_poll_interval");
    const savedTheme = localStorage.getItem("gitmap_terminal_theme");

    if (savedTemp) setTempDir(savedTemp);

    if (savedPoll) setPollInterval(savedPoll);

    if (savedTheme) setTerminalTheme(savedTheme);
  }, []);

  const handleSave = () => {
    localStorage.setItem("gitmap_temp_dir", tempDir);
    localStorage.setItem("gitmap_poll_interval", pollInterval);
    localStorage.setItem("gitmap_terminal_theme", terminalTheme);
    setIsSaved(true);

    toast({
      title: "Settings Saved",
      description: "Preferences updated successfully and synced with local client.",
    });

    setTimeout(() => setIsSaved(false), 2000);
  };

  return (
    <DocsLayout>
      <div className="max-w-4xl mx-auto space-y-8 pb-16">
        <div>
          <div className="flex items-center gap-3 mb-2">
            <div className="h-10 w-10 rounded-lg bg-primary/10 flex items-center justify-center dark:bg-primary/20">
              <SettingsIcon className="h-5 w-5 text-primary" />
            </div>
            <div>
              <h1 className="text-3xl font-heading font-bold text-foreground">Settings & Preferences</h1>
              <p className="text-sm text-muted-foreground">Manage temp folder locations, pipeline monitoring options, and UI preferences.</p>
            </div>
          </div>
        </div>

        <div className="grid gap-6">
          {/* Storage & Temp Directory */}
          <div className="rounded-xl border border-border bg-card p-6 shadow-sm space-y-4">
            <div className="flex items-center gap-2 text-foreground font-semibold text-lg">
              <Folder className="h-5 w-5 text-primary" />
              <h2>Storage & Temp Directory</h2>
            </div>
            <p className="text-sm text-muted-foreground">
              Directory path used for storing error logs with <code className="font-mono text-xs bg-muted px-1.5 py-0.5 rounded text-foreground">--tempfile</code> and automated CI artifacts.
            </p>
            <div className="space-y-2">
              <label htmlFor="temp-dir-input" className="text-xs font-mono uppercase text-muted-foreground tracking-wider">Temp Directory Path</label>
              <input
                id="temp-dir-input"
                type="text"
                value={tempDir}
                onChange={(e) => setTempDir(e.target.value)}
                className="w-full px-3 py-2 rounded-md border border-input bg-background text-foreground text-sm font-mono focus:outline-none focus:ring-2 focus:ring-primary"
                placeholder=".lovable/temp"
              />
            </div>
          </div>

          {/* CI/CD & Pipeline Polling */}
          <div className="rounded-xl border border-border bg-card p-6 shadow-sm space-y-4">
            <div className="flex items-center gap-2 text-foreground font-semibold text-lg">
              <RefreshCw className="h-5 w-5 text-primary" />
              <h2>CI/CD & Pipeline Monitor</h2>
            </div>
            <p className="text-sm text-muted-foreground">
              Configure automatic background polling frequency for <code className="font-mono text-xs bg-muted px-1.5 py-0.5 rounded text-foreground">gitmap pipeline status</code>.
            </p>
            <div className="space-y-2">
              <label htmlFor="poll-interval-input" className="text-xs font-mono uppercase text-muted-foreground tracking-wider">Status Refresh Interval (Seconds)</label>
              <input
                id="poll-interval-input"
                type="number"
                value={pollInterval}
                onChange={(e) => setPollInterval(e.target.value)}
                className="w-full px-3 py-2 rounded-md border border-input bg-background text-foreground text-sm font-mono focus:outline-none focus:ring-2 focus:ring-primary"
                min="5"
                max="120"
              />
            </div>
          </div>

          {/* Terminal Defaults */}
          <div className="rounded-xl border border-border bg-card p-6 shadow-sm space-y-4">
            <div className="flex items-center gap-2 text-foreground font-semibold text-lg">
              <Terminal className="h-5 w-5 text-primary" />
              <h2>Interactive UI Terminal</h2>
            </div>
            <p className="text-sm text-muted-foreground">
              Select the color palette theme for the integrated web terminal drawer.
            </p>
            <div className="grid grid-cols-3 gap-3">
              {[
                { id: "dark", label: "Dark+" },
                { id: "ubuntu", label: "Ubuntu Purple" },
                { id: "win", label: "PowerShell Navy" },
              ].map((themeOpt) => (
                <button
                  key={themeOpt.id}
                  type="button"
                  onClick={() => setTerminalTheme(themeOpt.id)}
                  className={`px-4 py-3 rounded-lg border text-sm font-medium transition-colors ${
                    terminalTheme === themeOpt.id
                      ? "border-primary bg-primary/10 text-foreground font-semibold"
                      : "border-border bg-background text-muted-foreground hover:text-foreground"
                  }`}
                >
                  {themeOpt.label}
                </button>
              ))}
            </div>
          </div>

          {/* GitHub CLI Authentication */}
          <div className="rounded-xl border border-border bg-card p-6 shadow-sm space-y-4">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-2 text-foreground font-semibold text-lg">
                <Github className="h-5 w-5 text-primary" />
                <h2>GitHub CLI Authentication</h2>
              </div>
              <span className="inline-flex items-center gap-1 text-xs font-mono px-2 py-0.5 rounded border bg-primary/10 text-primary border-primary/30">
                <ShieldCheck className="h-3.5 w-3.5" />
                GH Token Auto-Resolved
              </span>
            </div>
            <p className="text-sm text-muted-foreground">
              The pipeline monitor communicates directly with GitHub Actions using your system credentials (<code className="font-mono text-xs bg-muted px-1.5 py-0.5 rounded text-foreground">gh auth token</code> or <code className="font-mono text-xs bg-muted px-1.5 py-0.5 rounded text-foreground">GITHUB_TOKEN</code>).
            </p>
          </div>

          {/* Action Button */}
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
