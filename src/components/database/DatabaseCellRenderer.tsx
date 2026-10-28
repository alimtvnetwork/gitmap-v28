import { useMemo } from "react";
import { Binary, Code2, FileText } from "lucide-react";
import type { ClassifiedCell } from "./types";

interface CellRendererProps {
  columnName: string;
  value: unknown;
  onInspect: (col: string, classified: ClassifiedCell) => void;
  classifyCell: (v: unknown) => ClassifiedCell;
}

export const DatabaseCellRenderer = ({
  columnName,
  value,
  onInspect,
  classifyCell,
}: CellRendererProps) => {
  const classified = useMemo(() => classifyCell(value), [value, classifyCell]);

  if (classified.cellType === "null") {
    return <span className="text-muted-foreground font-medium italic font-sans text-xs">NULL</span>;
  }

  if (classified.cellType === "bool") {
    const isTrue = value === true;

    return (
      <span
        className={`inline-flex items-center px-1.5 py-0.5 rounded text-[10px] font-bold ${
          isTrue
            ? "bg-emerald-500/15 text-emerald-700 dark:text-emerald-400 border border-emerald-500/30"
            : "bg-rose-500/15 text-rose-700 dark:text-rose-400 border border-rose-500/30"
        }`}
      >
        {classified.badgeText}
      </span>
    );
  }

  if (classified.cellType === "blob") {
    return (
      <button
        type="button"
        onClick={() => onInspect(columnName, classified)}
        className="inline-flex items-center gap-1 px-2 py-0.5 rounded text-[11px] font-semibold bg-purple-500/15 text-purple-700 dark:text-purple-300 border border-purple-500/30 hover:bg-purple-500/25 transition-colors cursor-pointer active:scale-[0.98]"
      >
        <Binary className="h-3 w-3" />
        {classified.badgeText}
      </button>
    );
  }

  if (classified.cellType === "json") {
    return (
      <button
        type="button"
        onClick={() => onInspect(columnName, classified)}
        className="inline-flex items-center gap-1 px-2 py-0.5 rounded text-[11px] font-semibold bg-sky-500/15 text-sky-700 dark:text-sky-300 border border-sky-500/30 hover:bg-sky-500/25 transition-colors cursor-pointer active:scale-[0.98]"
      >
        <Code2 className="h-3 w-3" />
        {classified.badgeText}
      </button>
    );
  }

  if (classified.cellType === "multiline") {
    return (
      <button
        type="button"
        onClick={() => onInspect(columnName, classified)}
        className="inline-flex items-center gap-1 px-2 py-0.5 rounded text-[11px] font-semibold bg-amber-500/15 text-amber-800 dark:text-amber-300 border border-amber-500/30 hover:bg-amber-500/25 transition-colors cursor-pointer max-w-[280px] truncate active:scale-[0.98]"
      >
        <FileText className="h-3 w-3 flex-shrink-0" />
        <span className="truncate">{classified.badgeText}</span>
      </button>
    );
  }

  if (classified.cellType === "text") {
    return (
      <button
        type="button"
        onClick={() => onInspect(columnName, classified)}
        className="text-muted-foreground hover:text-sky-600 dark:hover:text-sky-400 underline decoration-dotted underline-offset-2 transition-colors cursor-pointer text-left"
      >
        {classified.badgeText}
      </button>
    );
  }

  return <span className="text-foreground">{classified.badgeText}</span>;
};

export default DatabaseCellRenderer;
