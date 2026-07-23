import { forwardRef } from "react";
import { Search, X } from "lucide-react";

interface SpecSearchBarProps {
  query: string;
  onChange: (value: string) => void;
}

const SpecSearchBar = forwardRef<HTMLInputElement, SpecSearchBarProps>(
  ({ query, onChange }, ref) => (
    <div className="relative mb-6">
      <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
      <input
        ref={ref}
        type="text"
        value={query}
        onChange={(e) => onChange(e.target.value)}
        placeholder="Filter specs… (e.g. release, database, TUI)"
        className="w-full pl-9 pr-16 py-2.5 text-sm font-mono bg-muted/30 border border-border rounded-lg text-foreground placeholder:text-muted-foreground/50 focus:outline-none focus:ring-2 focus:ring-primary/40 focus:border-primary/50 transition-colors"
      />
      {query ? (
        <button
          onClick={() => onChange("")}
          className="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground transition-colors"
        >
          <X className="h-4 w-4" />
        </button>
      ) : (
        <kbd className="absolute right-3 top-1/2 -translate-y-1/2 text-[10px] font-mono text-muted-foreground/50 border border-border rounded px-1.5 py-0.5 pointer-events-none">
          /
        </kbd>
      )}
    </div>
  )
);

SpecSearchBar.displayName = "SpecSearchBar";

export default SpecSearchBar;
