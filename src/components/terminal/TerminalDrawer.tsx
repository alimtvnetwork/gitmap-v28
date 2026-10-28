import { useState, useRef, useEffect, useCallback } from "react";
import {
  Terminal as TerminalIcon,
  ChevronUp,
  ChevronDown,
  Maximize2,
  Minimize2,
  X,
  Play,
  RefreshCw,
  Trash2,
  CheckCircle2,
  AlertCircle,
  Clock,
  Palette,
} from "lucide-react";
import { useTheme } from "@/hooks/useTheme";

interface TerminalItem {
  id: string;
  type: "system" | "cmd" | "stdout" | "stderr" | "meta";
  content: string;
  exitCode?: number;
  duration?: string;
}

const INITIAL_ITEMS: TerminalItem[] = [
  { id: "init-1", type: "system", content: "GitMap Integrated Terminal v6.153.0" },
  { id: "init-2", type: "system", content: "Type a gitmap command or click quick action buttons below." },
  { id: "init-3", type: "cmd", content: "gitmap pipeline status" },
  {
    id: "init-4",
    type: "stdout",
    content: "  ● Repo:             alimtvnetwork/gitmap-v28\n  ● Status:           completed (conclusion: success)\n  ● Last Tag Release: v6.153.0\n  ● Pending PRs:       1",
  },
  { id: "init-5", type: "meta", content: "exit 0", exitCode: 0, duration: "32ms" },
];

const THEME_OPTIONS = [
  { id: "dark", label: "Dark+" },
  { id: "ubuntu", label: "Ubuntu" },
  { id: "win", label: "PowerShell" },
] as const;

interface TerminalThemeConfig {
  containerThemeClass: string;
  headerThemeClass: string;
  quickBarThemeClass: string;
  outputThemeClass: string;
  inputBarThemeClass: string;
  promptColorClass: string;
  textColorClass: string;
}

const UBUNTU_THEME: TerminalThemeConfig = {
  containerThemeClass: "border-t border-terminal-ubuntu-border shadow-2xl bg-terminal-ubuntu-bg",
  headerThemeClass: "bg-terminal-ubuntu-header border-b border-terminal-ubuntu-border text-slate-100",
  quickBarThemeClass: "bg-terminal-ubuntu-quickbar border-b border-terminal-ubuntu-border text-slate-300",
  outputThemeClass: "bg-terminal-ubuntu-bg text-white",
  inputBarThemeClass: "bg-terminal-ubuntu-header border-t border-terminal-ubuntu-border text-white",
  promptColorClass: "text-terminal-ubuntu-prompt",
  textColorClass: "text-white",
};

const WIN_THEME: TerminalThemeConfig = {
  containerThemeClass: "border-t border-terminal-win-border shadow-2xl bg-terminal-win-bg",
  headerThemeClass: "bg-terminal-win-header border-b border-terminal-win-border text-slate-100",
  quickBarThemeClass: "bg-terminal-win-quickbar border-b border-terminal-win-border text-slate-300",
  outputThemeClass: "bg-terminal-win-bg text-white",
  inputBarThemeClass: "bg-terminal-win-header border-t border-terminal-win-border text-white",
  promptColorClass: "text-terminal-win-prompt",
  textColorClass: "text-white",
};

const LIGHT_THEME: TerminalThemeConfig = {
  containerThemeClass: "border-t-2 border-slate-700/90 shadow-[0_-12px_30px_rgba(0,0,0,0.35)] ring-1 ring-slate-800 bg-slate-900",
  headerThemeClass: "bg-slate-900 border-b border-slate-800 text-slate-100",
  quickBarThemeClass: "bg-slate-900/95 border-b border-slate-800 text-slate-300",
  outputThemeClass: "bg-slate-950 text-slate-100",
  inputBarThemeClass: "bg-slate-900 border-t border-slate-800 text-slate-100",
  promptColorClass: "text-emerald-400",
  textColorClass: "text-slate-100",
};

const DEFAULT_THEME: TerminalThemeConfig = {
  containerThemeClass: "border-t border-border bg-card shadow-2xl",
  headerThemeClass: "bg-muted/60 border-b border-border text-foreground",
  quickBarThemeClass: "bg-background/50 border-b border-border text-muted-foreground",
  outputThemeClass: "bg-terminal text-terminal-foreground",
  inputBarThemeClass: "bg-background border-t border-border text-foreground",
  promptColorClass: "text-primary",
  textColorClass: "text-slate-200",
};

function getThemeClasses(theme: string, isDark: boolean): TerminalThemeConfig {
  if (theme === "ubuntu") {
    return UBUNTU_THEME;
  }

  if (theme === "win") {
    return WIN_THEME;
  }

  if (!isDark) {
    return LIGHT_THEME;
  }

  return DEFAULT_THEME;
}

export function TerminalDrawer() {
  const { isDark } = useTheme();
  const [isOpen, setIsOpen] = useState(false);
  const [isMaximized, setIsMaximized] = useState(false);
  const [inputVal, setInputVal] = useState("");
  const [isRunning, setIsRunning] = useState(false);
  const [items, setItems] = useState<TerminalItem[]>(INITIAL_ITEMS);

  // Command history navigation
  const [cmdHistory, setCmdHistory] = useState<string[]>(["gitmap pipeline status"]);
  const [historyIndex, setHistoryIndex] = useState(-1);

  // Theme styling (dark, ubuntu, win)
  const [termTheme, setTermTheme] = useState<string>(() => {
    return localStorage.getItem("gitmap_terminal_theme") || "dark";
  });

  const bottomRef = useRef<HTMLDivElement>(null);
  const inputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    const handleStorageChange = () => {
      const storedTheme = localStorage.getItem("gitmap_terminal_theme");
      if (storedTheme) setTermTheme(storedTheme);
    };

    window.addEventListener("storage", handleStorageChange);

    return () => window.removeEventListener("storage", handleStorageChange);
  }, []);

  const handleSetTheme = (themeId: string) => {
    setTermTheme(themeId);
    localStorage.setItem("gitmap_terminal_theme", themeId);
  };

  useEffect(() => {
    if (isOpen) {
      bottomRef.current?.scrollIntoView({ behavior: "smooth" });
      inputRef.current?.focus();
    }
  }, [items, isOpen]);

  const runCommand = useCallback(async (cmdStr: string) => {
    const cleanCmd = cmdStr.trim();
    if (!cleanCmd) return;

    if (cleanCmd.toLowerCase() === "clear") {
      setItems([{ id: Date.now().toString(), type: "system", content: "Terminal cleared." }]);
      setInputVal("");
      setHistoryIndex(-1);
      return;
    }

    // Add command to history navigation
    setCmdHistory((prev) => {
      const filtered = prev.filter((c) => c !== cleanCmd);

      return [...filtered, cleanCmd];
    });
    setHistoryIndex(-1);

    const cmdEntryId = Date.now().toString() + "-cmd";
    setItems((prev) => [
      ...prev,
      { id: cmdEntryId, type: "cmd", content: cleanCmd },
    ]);
    setInputVal("");
    setIsRunning(true);

    const startTime = performance.now();
    let data: any = null;
    let durationMs = 0;

    try {
      let resp = await fetch("/api/command/exec", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ command: cleanCmd }),
      });

      if (resp.status === 404) {
        resp = await fetch("/api/terminal/exec", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ command: cleanCmd, nodeAlias: "local" }),
        });
      }

      durationMs = Math.round(performance.now() - startTime);

      if (resp.ok) {
        data = await resp.json();
      } else {
        const errText = await resp.text();
        try {
          data = JSON.parse(errText);
        } catch {
          data = {
            success: false,
            error: errText || `HTTP ${resp.status}: Command execution failed`,
            exitCode: 1,
          };
        }
      }
    } catch (networkErr: any) {
      durationMs = Math.round(performance.now() - startTime);
      data = {
        success: false,
        error: networkErr?.message || "Failed to contact execution server endpoint",
        exitCode: 1,
      };
    } finally {
      setIsRunning(false);
    }

    const stdout = data?.stdout || data?.output || "";
    const stderr = data?.stderr || data?.error || (Array.isArray(data?.errors) ? data.errors.join("\n") : "");
    const exitCode = typeof data?.exitCode === "number" ? data.exitCode : (data?.success ? 0 : 1);
    const durationStr = typeof data?.durationMs === "number" ? `${data.durationMs}ms` : `${durationMs}ms`;

    const newEntries: TerminalItem[] = [];

    if (stdout && stdout.trim()) {
      newEntries.push({
        id: Date.now().toString() + "-out",
        type: "stdout",
        content: stdout.trimEnd(),
      });
    }

    if (stderr && stderr.trim() && stderr.trim() !== stdout.trim()) {
      newEntries.push({
        id: Date.now().toString() + "-err",
        type: "stderr",
        content: stderr.trimEnd(),
      });
    }

    newEntries.push({
      id: Date.now().toString() + "-meta",
      type: "meta",
      content: `exit ${exitCode}`,
      exitCode,
      duration: durationStr,
    });

    setItems((prev) => [...prev, ...newEntries]);
  }, []);

  // Listen for custom events to trigger terminal commands from other pages (e.g. AGYPrompts)
  const runCommandRef = useRef(runCommand);
  runCommandRef.current = runCommand;

  useEffect(() => {
    const handleCustomTerminalEvent = (e: Event) => {
      const customEvt = e as CustomEvent<{ command?: string } | string>;
      const cmd = typeof customEvt.detail === "string" ? customEvt.detail : customEvt.detail?.command;
      setIsOpen(true);
      if (cmd && typeof cmd === "string") {
        runCommandRef.current(cmd);
      }
    };

    window.addEventListener("gitmap:terminal-run", handleCustomTerminalEvent);
    window.addEventListener("gitmap:open-terminal", handleCustomTerminalEvent);

    return () => {
      window.removeEventListener("gitmap:terminal-run", handleCustomTerminalEvent);
      window.removeEventListener("gitmap:open-terminal", handleCustomTerminalEvent);
    };
  }, []);

  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === "Enter" && !isRunning) {
      e.preventDefault();
      runCommand(inputVal);
      return;
    }

    if (e.key === "ArrowUp") {
      e.preventDefault();
      if (cmdHistory.length === 0) return;
      const nextIdx = historyIndex === -1 ? cmdHistory.length - 1 : Math.max(0, historyIndex - 1);
      setHistoryIndex(nextIdx);
      setInputVal(cmdHistory[nextIdx]);
      return;
    }

    if (e.key === "ArrowDown" && historyIndex !== -1) {
      e.preventDefault();
      const nextIdx = historyIndex + 1;
      const isOutOfBounds = nextIdx >= cmdHistory.length;
      if (isOutOfBounds) {
        setHistoryIndex(-1);
        setInputVal("");
        return;
      }

      setHistoryIndex(nextIdx);
      setInputVal(cmdHistory[nextIdx]);
      return;
    }
  };

  const {
    containerThemeClass,
    headerThemeClass,
    quickBarThemeClass,
    outputThemeClass,
    inputBarThemeClass,
    promptColorClass,
    textColorClass,
  } = getThemeClasses(termTheme, isDark);

  return (
    <>
      {/* Floating Toggle Button */}
      <button
        type="button"
        onClick={() => setIsOpen(!isOpen)}
        aria-label="Toggle Interactive Terminal"
        className="fixed bottom-4 right-4 z-40 flex items-center gap-2 px-3 py-2 rounded-lg bg-primary text-primary-foreground font-mono text-xs font-semibold shadow-lg hover:bg-primary/90 transition-all border border-primary/40 focus:outline-none focus:ring-2 focus:ring-primary"
      >
        <TerminalIcon className="h-4 w-4" />
        <span>Terminal</span>
        {isOpen ? <ChevronDown className="h-3.5 w-3.5" /> : <ChevronUp className="h-3.5 w-3.5" />}
      </button>

      {/* Slide-out Terminal Drawer */}
      {isOpen && (
        <div
          className={`fixed bottom-0 left-0 right-0 z-50 transition-all flex flex-col ${containerThemeClass} ${
            isMaximized ? "h-[80vh]" : "h-80"
          }`}
        >
          {/* Header Bar */}
          <div className={`flex items-center justify-between px-4 py-2 ${headerThemeClass}`}>
            <div className="flex items-center gap-3 font-mono text-xs font-semibold">
              <div className="flex items-center gap-2">
                <TerminalIcon className={`h-4 w-4 ${promptColorClass}`} />
                <span>GitMap Interactive Live Terminal</span>
              </div>

              {/* Execution Status Indicator */}
              {isRunning ? (
                <div className="flex items-center gap-1.5 px-2 py-0.5 rounded-full bg-primary/20 text-primary border border-primary/40 font-mono text-[10px] animate-pulse">
                  <RefreshCw className="h-3 w-3 animate-spin" />
                  <span>Running...</span>
                </div>
              ) : (
                <div className="hidden sm:flex items-center gap-1 px-2 py-0.5 rounded-full bg-emerald-500/10 text-emerald-400 border border-emerald-500/20 font-mono text-[10px]">
                  <span className="w-1.5 h-1.5 rounded-full bg-emerald-400" />
                  <span>Ready</span>
                </div>
              )}
            </div>

            <div className="flex items-center gap-2">
              {/* Terminal Theme Selector Chips */}
              <div className="hidden md:flex items-center gap-1 mr-2 px-1.5 py-0.5 rounded border border-white/10 bg-black/20 text-[10px] font-mono">
                <Palette className="h-3 w-3 opacity-70" />
                {THEME_OPTIONS.map((opt) => (
                  <button
                    key={opt.id}
                    type="button"
                    onClick={() => handleSetTheme(opt.id)}
                    className={`px-1.5 py-0.5 rounded transition-colors ${
                      termTheme === opt.id
                        ? "bg-primary text-primary-foreground font-semibold"
                        : "opacity-70 hover:opacity-100"
                    }`}
                  >
                    {opt.label}
                  </button>
                ))}
              </div>

              <button
                type="button"
                onClick={() => setItems([{ id: Date.now().toString(), type: "system", content: "Terminal cleared." }])}
                title="Clear Terminal"
                className="p-1 rounded opacity-75 hover:opacity-100 hover:bg-white/10"
              >
                <Trash2 className="h-3.5 w-3.5" />
              </button>
              <button
                type="button"
                onClick={() => setIsMaximized(!isMaximized)}
                title={isMaximized ? "Restore Size" : "Maximize"}
                className="p-1 rounded opacity-75 hover:opacity-100 hover:bg-white/10"
              >
                {isMaximized ? <Minimize2 className="h-3.5 w-3.5" /> : <Maximize2 className="h-3.5 w-3.5" />}
              </button>
              <button
                type="button"
                onClick={() => setIsOpen(false)}
                title="Close"
                className="p-1 rounded opacity-75 hover:opacity-100 hover:bg-white/10"
              >
                <X className="h-3.5 w-3.5" />
              </button>
            </div>
          </div>

          {/* Quick Action Badges */}
          <div className={`flex items-center gap-2 px-4 py-1.5 overflow-x-auto text-[11px] font-mono ${quickBarThemeClass}`}>
            <span className="opacity-70 shrink-0">Quick Commands:</span>
            {[
              "gitmap pipeline status",
              "gitmap pipeline waittime",
              "gitmap pipeline error-logs",
              "gitmap help",
              "gitmap reclone --help",
              "gitmap conflicts",
            ].map((qCmd) => (
              <button
                key={qCmd}
                type="button"
                disabled={isRunning}
                onClick={() => runCommand(qCmd)}
                className="px-2 py-0.5 rounded bg-black/20 hover:bg-primary/20 hover:text-primary transition-colors shrink-0 disabled:opacity-50"
              >
                {qCmd}
              </button>
            ))}
          </div>

          {/* Terminal Output Area */}
          <div className={`flex-1 overflow-y-auto p-4 font-mono text-xs space-y-1.5 select-text ${outputThemeClass}`}>
            {items.map((item) => {
              if (item.type === "system") {
                return (
                  <div key={item.id} className="opacity-60 italic leading-relaxed">
                    {item.content}
                  </div>
                );
              }

              if (item.type === "cmd") {
                return (
                  <div key={item.id} className="flex items-center gap-1.5 pt-1.5 font-bold">
                    <span className={promptColorClass}>$</span>
                    <span className={textColorClass}>{item.content}</span>
                  </div>
                );
              }

              if (item.type === "stdout") {
                return (
                  <pre key={item.id} className={`whitespace-pre-wrap font-mono leading-relaxed pl-3 select-text ${textColorClass}`}>
                    {item.content}
                  </pre>
                );
              }

              if (item.type === "stderr") {
                return (
                  <pre key={item.id} className="whitespace-pre-wrap font-mono leading-relaxed pl-3 select-text text-red-400">
                    {item.content}
                  </pre>
                );
              }

              if (item.type === "meta") {
                return (
                  <div key={item.id} className="flex items-center gap-3 text-[11px] font-mono pl-3 pb-1 opacity-80">
                    {item.exitCode === 0 ? (
                      <span className="inline-flex items-center gap-1 text-emerald-400 font-semibold">
                        <CheckCircle2 className="h-3 w-3" />
                        <span>exit 0</span>
                      </span>
                    ) : (
                      <span className="inline-flex items-center gap-1 text-red-400 font-semibold">
                        <AlertCircle className="h-3 w-3" />
                        <span>exit {item.exitCode ?? 1}</span>
                      </span>
                    )}
                    {item.duration && (
                      <span className="inline-flex items-center gap-1 opacity-70 text-[10px]">
                        <Clock className="h-2.5 w-2.5" />
                        <span>{item.duration}</span>
                      </span>
                    )}
                  </div>
                );
              }

              return null;
            })}
            <div ref={bottomRef} />
          </div>

          {/* Command Input Row */}
          <div className={`flex items-center gap-2 px-3 py-2 ${inputBarThemeClass}`}>
            <span className={`font-mono text-xs font-bold ${promptColorClass}`}>$</span>
            <input
              ref={inputRef}
              type="text"
              value={inputVal}
              onChange={(e) => setInputVal(e.target.value)}
              onKeyDown={handleKeyDown}
              disabled={isRunning}
              placeholder="Type command (e.g. gitmap pipeline status) and press Enter... (↑/↓ for history)"
              className="flex-1 bg-transparent text-inherit font-mono text-xs focus:outline-none placeholder:opacity-50"
            />
            <button
              type="button"
              disabled={isRunning || !inputVal.trim()}
              onClick={() => runCommand(inputVal)}
              className="flex items-center gap-1.5 px-3 py-1 rounded bg-primary text-primary-foreground font-mono text-xs font-semibold hover:bg-primary/90 transition-colors disabled:opacity-50"
            >
              {isRunning ? (
                <>
                  <RefreshCw className="h-3 w-3 animate-spin" />
                  <span>Running</span>
                </>
              ) : (
                <>
                  <Play className="h-3 w-3" />
                  <span>Run</span>
                </>
              )}
            </button>
          </div>
        </div>
      )}
    </>
  );
}

export default TerminalDrawer;
