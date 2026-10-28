import { useState, useMemo } from "react";
import { X, Copy, Check } from "lucide-react";
import type { InspectTarget } from "./types";

interface CellInspectionModalProps {
  target: InspectTarget;
  hasCopied: boolean;
  onClose: () => void;
  onCopyRaw: (text: string) => void;
  generateHexDump: (input: string) => string;
}

export const CellInspectionModal = ({
  target,
  hasCopied,
  onClose,
  onCopyRaw,
  generateHexDump,
}: CellInspectionModalProps) => {
  const [activeTab, setActiveTab] = useState<"hex" | "json" | "text">(() => {
    if (target.classified.cellType === "json") return "json";
    if (target.classified.cellType === "multiline" || target.classified.cellType === "text") return "text";

    return "hex";
  });

  const rawString = useMemo(() => {
    const raw = target.classified.rawValue;
    if (typeof raw === "string") return raw;

    return JSON.stringify(raw, null, 2);
  }, [target]);

  const hexLines = useMemo(() => {
    return generateHexDump(rawString);
  }, [rawString, generateHexDump]);

  const jsonFormatted = useMemo(() => {
    const raw = target.classified.rawValue;
    if (typeof raw === "object" && raw !== null) {
      return JSON.stringify(raw, null, 2);
    }

    if (typeof raw === "string") {
      try {
        return JSON.stringify(JSON.parse(raw), null, 2);
      } catch {
        return raw;
      }
    }

    return String(raw);
  }, [target]);

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/70 backdrop-blur-sm p-4">
      <div className="w-full max-w-3xl max-h-[85vh] bg-card border border-border rounded-xl shadow-2xl flex flex-col overflow-hidden transform-gpu">
        {/* Modal Header */}
        <div className="flex items-center justify-between px-5 py-3.5 border-b border-border bg-card">
          <div className="flex items-center gap-2.5">
            <h3 className="text-sm font-bold text-foreground">Cell Inspector</h3>
            <span className="px-2 py-0.5 rounded text-[11px] font-mono bg-primary/10 text-primary border border-primary/20">
              {target.columnName}
            </span>
            <span className="px-2 py-0.5 rounded text-[11px] font-mono bg-muted text-muted-foreground border border-border">
              {target.classified.badgeText}
            </span>
          </div>
          <button
            type="button"
            onClick={onClose}
            className="h-8 w-8 rounded flex items-center justify-center text-muted-foreground hover:text-foreground hover:bg-muted transition-colors active:scale-[0.98]"
          >
            <X className="h-4 w-4" />
          </button>
        </div>

        {/* Modal Nav */}
        <div className="flex gap-1 px-5 pt-3 border-b border-border bg-card">
          {[
            { id: "hex", label: "Hex Dump" },
            { id: "json", label: "Formatted JSON" },
            { id: "text", label: "Decoded String" },
          ].map((tab) => (
            <button
              key={tab.id}
              type="button"
              onClick={() => setActiveTab(tab.id as "hex" | "json" | "text")}
              className={`px-3.5 py-1.5 text-xs font-semibold rounded-t border-b-2 transition-colors ${
                activeTab === tab.id
                  ? "text-primary border-primary bg-primary/5"
                  : "text-muted-foreground border-transparent hover:text-foreground"
              }`}
            >
              {tab.label}
            </button>
          ))}
        </div>

        {/* View Content */}
        <div className="flex-1 overflow-auto p-4 bg-muted/20 font-mono text-xs">
          {activeTab === "hex" && (
            <div className="bg-muted/40 p-3.5 rounded-lg border border-border overflow-x-auto text-[11px] leading-relaxed text-sky-700 dark:text-sky-300 whitespace-pre shadow-inner">
              {hexLines}
            </div>
          )}
          {activeTab === "json" && (
            <pre className="bg-muted/40 p-3.5 rounded-lg border border-border overflow-x-auto text-[11px] leading-relaxed text-emerald-700 dark:text-emerald-300 whitespace-pre-wrap word-break shadow-inner">
              {jsonFormatted}
            </pre>
          )}
          {activeTab === "text" && (
            <pre className="bg-muted/40 p-3.5 rounded-lg border border-border overflow-x-auto text-[11px] leading-relaxed text-amber-800 dark:text-amber-300 whitespace-pre-wrap word-break shadow-inner">
              {rawString}
            </pre>
          )}
        </div>

        {/* Modal Footer */}
        <div className="flex items-center justify-between px-5 py-3 border-t border-border bg-card">
          <span className="text-xs text-muted-foreground">
            Active View: <strong className="text-foreground uppercase">{activeTab}</strong>
          </span>
          <div className="flex items-center gap-2">
            <button
              type="button"
              onClick={() => {
                const textToCopy = activeTab === "hex" ? hexLines : activeTab === "json" ? jsonFormatted : rawString;
                onCopyRaw(textToCopy);
              }}
              className="inline-flex items-center gap-1.5 px-3 py-1.5 min-h-[36px] rounded-md text-xs font-medium bg-muted text-foreground hover:bg-muted/80 border border-border transition-colors active:scale-[0.98]"
            >
              {hasCopied ? <Check className="h-3.5 w-3.5 text-emerald-600 dark:text-emerald-400" /> : <Copy className="h-3.5 w-3.5" />}
              {hasCopied ? "Copied!" : "Copy Content"}
            </button>
            <button
              type="button"
              onClick={onClose}
              className="px-3.5 py-1.5 min-h-[36px] rounded-md text-xs font-semibold bg-primary text-primary-foreground hover:bg-primary/90 transition-colors active:scale-[0.98]"
            >
              Done
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default CellInspectionModal;
