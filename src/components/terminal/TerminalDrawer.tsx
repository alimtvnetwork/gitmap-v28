import { useState, useRef, useEffect } from "react";
import { Terminal as TerminalIcon, ChevronUp, ChevronDown, Maximize2, Minimize2, X, Play, RefreshCw, Trash2 } from "lucide-react";
import { XTermRenderer, OsTheme } from "@/components/terminal/TerminalView";

export function TerminalDrawer() {
  const [isOpen, setIsOpen] = useState(false);
  const [isMaximized, setIsMaximized] = useState(false);
  const [activeTab, setActiveTab] = useState("bash");
  const [inputVal, setInputVal] = useState("");
  const [history, setHistory] = useState<string[]>([
    "GitMap Integrated Terminal v6.153.0",
    "Type a gitmap command or click quick actions below.",
    "$ gitmap pipeline status",
    "  ● Repo:             alimtvnetwork/gitmap-v28",
    "  ● Status:           completed (conclusion: success)",
    "  ● Last Tag Release: v6.153.0",
    "  ● Pending PRs:       1",
    "",
  ]);

  const bottomRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (isOpen) {
      bottomRef.current?.scrollIntoView({ behavior: "smooth" });
    }
  }, [history, isOpen]);

  const runCommand = (cmdStr: string) => {
    const cleanCmd = cmdStr.trim();
    if (!cleanCmd) return;

    setHistory((prev) => [
      ...prev,
      `$ ${cleanCmd}`,
      `Executing '${cleanCmd}'...`,
      cleanCmd.includes("pipeline") ? "  ● [pipeline] Status: OK. All 23 quality gates green." : "  ● Command executed successfully.",
      "",
    ]);
    setInputVal("");
  };

  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === "Enter") {
      runCommand(inputVal);
    }
  };

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
          className={`fixed bottom-0 left-0 right-0 z-50 border-t border-border bg-card shadow-2xl transition-all flex flex-col ${
            isMaximized ? "h-[80vh]" : "h-72"
          }`}
        >
          {/* Header Bar */}
          <div className="flex items-center justify-between px-4 py-2 border-b border-border bg-muted/60">
            <div className="flex items-center gap-2 font-mono text-xs text-foreground font-semibold">
              <TerminalIcon className="h-3.5 w-3.5 text-primary" />
              <span>GitMap Interactive Web Terminal</span>
            </div>
            <div className="flex items-center gap-2">
              <button
                type="button"
                onClick={() => setHistory(["Terminal cleared."])}
                title="Clear Terminal"
                className="p-1 rounded text-muted-foreground hover:text-foreground hover:bg-muted"
              >
                <Trash2 className="h-3.5 w-3.5" />
              </button>
              <button
                type="button"
                onClick={() => setIsMaximized(!isMaximized)}
                title={isMaximized ? "Restore Size" : "Maximize"}
                className="p-1 rounded text-muted-foreground hover:text-foreground hover:bg-muted"
              >
                {isMaximized ? <Minimize2 className="h-3.5 w-3.5" /> : <Maximize2 className="h-3.5 w-3.5" />}
              </button>
              <button
                type="button"
                onClick={() => setIsOpen(false)}
                title="Close"
                className="p-1 rounded text-muted-foreground hover:text-foreground hover:bg-muted"
              >
                <X className="h-3.5 w-3.5" />
              </button>
            </div>
          </div>

          {/* Quick Action Badges */}
          <div className="flex items-center gap-2 px-4 py-1.5 border-b border-border bg-background/50 overflow-x-auto text-[11px] font-mono">
            <span className="text-muted-foreground shrink-0">Quick Commands:</span>
            {[
              "gitmap pipeline status",
              "gitmap pipeline waittime",
              "gitmap pipeline error-logs",
              "gitmap help",
              "gitmap ls",
            ].map((qCmd) => (
              <button
                key={qCmd}
                type="button"
                onClick={() => runCommand(qCmd)}
                className="px-2 py-0.5 rounded bg-muted hover:bg-primary/20 text-foreground hover:text-primary transition-colors shrink-0"
              >
                {qCmd}
              </button>
            ))}
          </div>

          {/* Terminal Output Area */}
          <div className="flex-1 overflow-y-auto p-4 font-mono text-xs bg-[#0c0c0c] text-[#f1f1f1] space-y-1 select-text">
            {history.map((line, idx) => (
              <div key={idx} className={line.startsWith("$") ? "text-primary font-bold" : "text-slate-200 leading-relaxed"}>
                {line}
              </div>
            ))}
            <div ref={bottomRef} />
          </div>

          {/* Command Input Row */}
          <div className="flex items-center gap-2 px-3 py-2 border-t border-border bg-background">
            <span className="font-mono text-xs text-primary font-bold">$</span>
            <input
              type="text"
              value={inputVal}
              onChange={(e) => setInputVal(e.target.value)}
              onKeyDown={handleKeyDown}
              placeholder="Type command (e.g. gitmap pipeline status) and press Enter..."
              className="flex-1 bg-transparent text-foreground font-mono text-xs focus:outline-none"
            />
            <button
              type="button"
              onClick={() => runCommand(inputVal)}
              className="px-3 py-1 rounded bg-primary text-primary-foreground font-mono text-xs font-semibold hover:bg-primary/90 transition-colors"
            >
              Run
            </button>
          </div>
        </div>
      )}
    </>
  );
}
