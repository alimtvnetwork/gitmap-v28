import { Link } from "react-router-dom";
import { ChevronRight } from "lucide-react";
import type { SpecEntry } from "./specData";

const SpecEntryRow = ({ entry }: { entry: SpecEntry }) => (
  <div className="flex items-center gap-3 px-5 py-2.5 hover:bg-muted/20 transition-colors group">
    <span className="text-xs font-mono text-muted-foreground w-10 shrink-0">{entry.id}</span>
    <span className="text-sm text-foreground">{entry.title}</span>
    {entry.link && (
      <Link
        to={entry.link}
        className="ml-auto flex items-center gap-1 text-xs font-mono text-primary opacity-0 group-hover:opacity-100 transition-opacity"
      >
        docs <ChevronRight className="h-3 w-3" />
      </Link>
    )}
  </div>
);

export default SpecEntryRow;
