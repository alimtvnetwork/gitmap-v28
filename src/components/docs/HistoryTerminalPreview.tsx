export interface HistoryItem {
  command: string;
  alias: string;
  flags: string;
  status: string;
  duration: string;
  time: string;
}

interface HistoryTerminalPreviewProps {
  items: HistoryItem[];
}

export const HistoryTerminalPreview = ({ items }: HistoryTerminalPreviewProps) => {
  return (
    <div className="rounded-lg border border-border overflow-hidden my-6 shadow-sm">
      <div className="bg-terminal px-4 py-2 flex items-center gap-2 border-b border-border">
        <div className="flex gap-1.5">
          <span className="w-3 h-3 rounded-full bg-red-500/80" />
          <span className="w-3 h-3 rounded-full bg-yellow-500/80" />
          <span className="w-3 h-3 rounded-full bg-green-500/80" />
        </div>
        <span className="text-xs font-mono text-muted-foreground ml-2">gitmap history</span>
      </div>
      <div className="bg-terminal p-4 font-mono text-sm leading-relaxed overflow-x-auto">
        <div className="text-primary font-bold">
          {"  "}COMMAND{"     "}ALIAS{"  "}FLAGS{"                "}STATUS{"  "}DURATION{"   "}TIME
        </div>
        {items.map((r) => (
          <div key={r.time} className="text-terminal-foreground hover:bg-white/5 transition-colors">
            {"  "}
            <span className="inline-block w-[88px] text-foreground font-semibold">{r.command}</span>
            <span className="inline-block w-[48px] text-muted-foreground">{r.alias}</span>
            <span className="inline-block w-[160px] text-muted-foreground">{r.flags}</span>
            <span
              className={`inline-block w-[56px] font-semibold ${
                r.status === "FAIL" ? "text-rose-600 dark:text-rose-400" : "text-emerald-600 dark:text-emerald-400"
              }`}
            >
              {r.status === "FAIL" ? "✗" : "✔"} {r.status}
            </span>
            <span className="inline-block w-[80px] text-muted-foreground">{r.duration}</span>
            <span className="text-muted-foreground">{r.time}</span>
          </div>
        ))}
        <div className="mt-3 text-muted-foreground text-xs border-t border-border/50 pt-2">
          Showing {items.length} entries | 1 failed
        </div>
      </div>
    </div>
  );
};

export default HistoryTerminalPreview;
