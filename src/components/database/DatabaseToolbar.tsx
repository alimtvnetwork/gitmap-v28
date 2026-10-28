import { Search } from "lucide-react";

interface DatabaseToolbarProps {
  activeDb: string;
  activeTable: string;
  availableTables: string[];
  filterQuery: string;
  onSelectDb: (db: string) => void;
  onSelectTable: (table: string) => void;
  onChangeFilter: (val: string) => void;
}

const DBS = ["master", "pipeline", "tasks"];

export const DatabaseToolbar = ({
  activeDb,
  activeTable,
  availableTables,
  filterQuery,
  onSelectDb,
  onSelectTable,
  onChangeFilter,
}: DatabaseToolbarProps) => {
  return (
    <div className="flex flex-col sm:flex-row gap-4 justify-between items-start sm:items-center bg-card p-4 rounded-lg border border-border shadow-sm">
      <div className="flex flex-wrap items-center gap-3">
        <div className="flex items-center gap-1.5 bg-muted/50 p-1 rounded-md border border-border">
          {DBS.map((db) => (
            <button
              key={db}
              type="button"
              onClick={() => onSelectDb(db)}
              className={`px-3 py-1.5 min-h-[36px] text-xs font-medium rounded transition-all active:scale-[0.98] ${
                activeDb === db
                  ? "bg-primary text-primary-foreground shadow-sm"
                  : "text-muted-foreground hover:text-foreground"
              }`}
            >
              {db}
            </button>
          ))}
        </div>
        <div className="flex items-center gap-1.5 bg-muted/50 p-1 rounded-md border border-border">
          {availableTables.map((t) => (
            <button
              key={t}
              type="button"
              onClick={() => onSelectTable(t)}
              className={`px-3 py-1.5 min-h-[36px] text-xs font-medium rounded transition-all active:scale-[0.98] ${
                activeTable === t
                  ? "bg-background text-foreground shadow-sm border border-border"
                  : "text-muted-foreground hover:text-foreground"
              }`}
            >
              {t}
            </button>
          ))}
        </div>
      </div>
      <div className="relative w-full sm:w-64">
        <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
        <input
          type="text"
          value={filterQuery}
          onChange={(e) => onChangeFilter(e.target.value)}
          placeholder="Filter records..."
          className="w-full bg-background border border-border rounded-md pl-9 pr-3 py-1.5 text-xs text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-1 focus:ring-primary"
        />
      </div>
    </div>
  );
};

export default DatabaseToolbar;
