const files = [
  ["constants/constants_bookmark.go", "SQL, command names, messages"],
  ["model/bookmark.go", "BookmarkRecord struct"],
  ["store/bookmark.go", "Bookmark CRUD operations"],
  ["cmd/bookmark.go", "Bookmark command routing"],
  ["cmd/bookmarksave.go", "Save subcommand"],
  ["cmd/bookmarklist.go", "List and delete subcommands"],
  ["cmd/bookmarkrun.go", "Run (replay) subcommand"],
];

export function BookmarkFileLayoutTable() {
  return (
    <div className="rounded-lg border border-border overflow-hidden">
      <table className="w-full text-sm">
        <thead>
          <tr className="bg-muted/50">
            <th className="text-left font-mono font-semibold px-4 py-2">File</th>
            <th className="text-left font-mono font-semibold px-4 py-2">Purpose</th>
          </tr>
        </thead>
        <tbody>
          {files.map(([file, purpose]) => (
            <tr key={file} className="border-t border-border">
              <td className="px-4 py-2 font-mono text-primary">{file}</td>
              <td className="px-4 py-2 text-muted-foreground">{purpose}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
