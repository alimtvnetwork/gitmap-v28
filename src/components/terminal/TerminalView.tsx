import React, { useEffect, useRef } from "react";
import { Terminal } from "@xterm/xterm";
import { FitAddon } from "@xterm/addon-fit";
import "@xterm/xterm/css/xterm.css";
import { ResizablePanelGroup, ResizablePanel, ResizableHandle } from "@/components/ui/resizable";
import { Tabs, TabsList, TabsTrigger, TabsContent } from "@/components/ui/tabs";
import { cn } from "@/lib/utils";

export type OsTheme = "win" | "ubuntu" | "linux";

export interface TerminalTab {
  id: string;
  title: string;
  content: string;
}

export interface TerminalViewProps {
  theme?: OsTheme;
  tabs: TerminalTab[];
  hasSplits?: boolean;
}

function getXTermTheme(os: OsTheme) {
  const isWin = os === "win";
  const isLinux = os === "linux";

  if (isWin) return { background: "#012456", foreground: "#ffffff" };

  if (isLinux) return { background: "#000000", foreground: "#00ff00" };

  return { background: "#300a24", foreground: "#ffffff" };
}

function getThemeBgClass(os: OsTheme) {
  const isWin = os === "win";
  const isLinux = os === "linux";

  if (isWin) return "bg-[#012456]";

  if (isLinux) return "bg-black";

  return "bg-[#300a24]";
}

export function XTermRenderer({ content, theme }: { content: string; theme: OsTheme }) {
  const ref = useRef<HTMLDivElement>(null);
  useEffect(() => {
    if (ref.current) {
      const term = new Terminal({ theme: getXTermTheme(theme), fontFamily: "'Ubuntu Mono', monospace" });
      const fit = new FitAddon();
      term.loadAddon(fit);
      term.open(ref.current);
      fit.fit();
      term.write(content.replace(/\n/g, "\r\n"));

      return () => term.dispose();
    }
  }, [content, theme]);

  return <div ref={ref} className="h-full w-full pl-2 pt-2" />;
}

export function TerminalTabs({ tabs, theme }: { tabs: TerminalTab[]; theme: OsTheme }) {
  return (
    <Tabs defaultValue={tabs[0]?.id} className={cn("h-full flex flex-col", getThemeBgClass(theme))}>
      <TabsList className="justify-start rounded-none border-b bg-black/20 px-2 h-10">
        {tabs.map((t) => (
          <TabsTrigger key={t.id} value={t.id} className="text-white/70 data-[state=active]:bg-white/10 data-[state=active]:text-white">{t.title}</TabsTrigger>
        ))}
      </TabsList>
      {tabs.map((t) => (
        <TabsContent key={t.id} value={t.id} className="flex-1 mt-0 overflow-hidden pb-20">
          <XTermRenderer content={t.content} theme={theme} />
        </TabsContent>
      ))}
    </Tabs>
  );
}

export function TerminalView({ theme = "ubuntu", tabs, hasSplits = false }: TerminalViewProps) {
  const isSplitEnabled = hasSplits && tabs.length > 1;

  if (isSplitEnabled) {
    return (
      <ResizablePanelGroup direction="horizontal" className="h-[500px] border">
        <ResizablePanel defaultSize={50}><TerminalTabs tabs={[tabs[0]]} theme={theme} /></ResizablePanel>
        <ResizableHandle />
        <ResizablePanel defaultSize={50}><TerminalTabs tabs={tabs.slice(1)} theme={theme} /></ResizablePanel>
      </ResizablePanelGroup>
    );
  }

  return <div className="h-[500px] border overflow-hidden"><TerminalTabs tabs={tabs} theme={theme} /></div>;
}
