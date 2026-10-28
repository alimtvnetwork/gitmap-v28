import type { ClassifiedCell } from "./types";
import DatabaseCellRenderer from "./DatabaseCellRenderer";

interface DatabaseGridProps {
  rows: Array<Record<string, unknown>>;
  onInspect: (col: string, classified: ClassifiedCell) => void;
  classifyCell: (v: unknown) => ClassifiedCell;
}

export const DatabaseGrid = ({ rows, onInspect, classifyCell }: DatabaseGridProps) => {
  if (rows.length === 0) {
    return (
      <div className="p-12 text-center text-sm text-muted-foreground border border-dashed border-border rounded-lg bg-card">
        No records found matching criteria.
      </div>
    );
  }

  const columns = Object.keys(rows[0]);

  return (
    <div className="overflow-x-auto border border-border rounded-lg bg-card shadow-sm">
      <table className="w-full border-collapse text-left text-xs">
        <thead className="bg-muted/60 border-b border-border text-muted-foreground">
          <tr>
            {columns.map((c) => (
              <th key={c} className="p-3 font-semibold">
                {c}
              </th>
            ))}
          </tr>
        </thead>
        <tbody className="divide-y divide-border">
          {rows.map((row, idx) => (
            <tr key={idx} className="hover:bg-muted/30 transition-colors">
              {columns.map((c) => (
                <td key={c} className="p-3 font-mono">
                  <DatabaseCellRenderer
                    columnName={c}
                    value={row[c]}
                    onInspect={onInspect}
                    classifyCell={classifyCell}
                  />
                </td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default DatabaseGrid;
