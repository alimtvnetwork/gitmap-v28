const schema = [
  ["Id", "TEXT", "PRIMARY KEY", "Timestamp-based unique ID"],
  ["Name", "TEXT", "NOT NULL UNIQUE", "User-chosen bookmark name"],
  ["Command", "TEXT", "NOT NULL", "Command name (e.g. scan)"],
  ["Args", "TEXT", "DEFAULT ''", "Positional arguments"],
  ["Flags", "TEXT", "DEFAULT ''", "Flags (e.g. --mode ssh)"],
  ["CreatedAt", "TEXT", "DEFAULT CURRENT_TIMESTAMP", ""],
];

export function BookmarkSchemaTable() {
  return (
    <div className="rounded-lg border border-border overflow-hidden">
      <table className="w-full text-sm">
        <thead>
          <tr className="bg-muted/50">
            <th className="text-left font-mono font-semibold px-4 py-2">Column</th>
            <th className="text-left font-mono font-semibold px-4 py-2">Type</th>
            <th className="text-left font-mono font-semibold px-4 py-2">Constraints</th>
            <th className="text-left font-mono font-semibold px-4 py-2">Notes</th>
          </tr>
        </thead>
        <tbody>
          {schema.map(([col, type, constraints, notes]) => (
            <tr key={col} className="border-t border-border">
              <td className="px-4 py-2 font-mono text-primary">{col}</td>
              <td className="px-4 py-2 font-mono text-muted-foreground">{type}</td>
              <td className="px-4 py-2 text-muted-foreground">{constraints}</td>
              <td className="px-4 py-2 text-muted-foreground">{notes}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
