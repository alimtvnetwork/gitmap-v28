const MOCK_BOOKMARKS = [
  { name: "ssh-scan", command: "scan", args: "", flags: "--mode ssh", created: "2026-03-15 09:12" },
  { name: "quick-status", command: "status", args: "", flags: "", created: "2026-03-14 14:30" },
  { name: "scan-projects", command: "scan", args: "./projects", flags: "--mode ssh --open", created: "2026-03-12 11:45" },
  { name: "daily-pull", command: "pull", args: "", flags: "--all", created: "2026-03-10 08:00" },
];

export function BookmarkTerminalPreview() {
  return (
    <div className="rounded-lg border border-border overflow-hidden my-6">
      <div className="bg-terminal px-4 py-2 flex items-center gap-2 border-b border-border">
        <div className="flex gap-1.5">
          <span className="w-3 h-3 rounded-full bg-red-500/80" />
          <span className="w-3 h-3 rounded-full bg-yellow-500/80" />
          <span className="w-3 h-3 rounded-full bg-green-500/80" />
        </div>
        <span className="text-xs font-mono text-muted-foreground ml-2">gitmap bookmark list</span>
      </div>
      <div className="bg-terminal p-4 font-mono text-sm leading-relaxed overflow-x-auto">
        <div className="text-primary font-bold text-xs">
          {"  "}NAME{"              "}COMMAND{"    "}ARGS{"          "}FLAGS{"                "}CREATED
        </div>
        {MOCK_BOOKMARKS.map((b) => (
          <div key={b.name} className="text-terminal-foreground text-xs">
            {"  "}
            <span className="inline-block w-[120px] text-foreground">{b.name}</span>
            <span className="inline-block w-[72px] text-primary">{b.command}</span>
            <span className="inline-block w-[100px] text-muted-foreground">{b.args || "—"}</span>
            <span className="inline-block w-[140px] text-muted-foreground">{b.flags || "—"}</span>
            <span className="text-muted-foreground">{b.created}</span>
          </div>
        ))}
        <div className="mt-3 text-muted-foreground text-xs border-t border-border/50 pt-2">
          {MOCK_BOOKMARKS.length} bookmarks saved
        </div>
      </div>
    </div>
  );
}
