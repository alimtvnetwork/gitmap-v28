import { commitInExitCodes as exitCodes } from "@/pages/commitInData";

export function CommitInExitCodesTable() {
  return (
    <section>
      <h2 className="text-xl font-semibold mb-3">Exit Codes</h2>
      <div className="overflow-x-auto">
        <table className="w-full text-sm border border-border rounded-lg">
          <thead>
            <tr className="bg-muted/50">
              <th className="text-left px-4 py-2 font-medium">Code</th>
              <th className="text-left px-4 py-2 font-medium">Meaning</th>
            </tr>
          </thead>
          <tbody>
            {exitCodes.map((e) => (
              <tr key={e.code} className="border-t border-border">
                <td className="px-4 py-2 font-mono text-primary">{e.code}</td>
                <td className="px-4 py-2 text-muted-foreground">{e.meaning}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </section>
  );
}
