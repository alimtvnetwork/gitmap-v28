import { commitInFlags as flags } from "@/pages/commitInData";

export function CommitInFlagsTable() {
  return (
    <section>
      <h2 className="text-xl font-semibold mb-3">Flags</h2>
      <div className="overflow-x-auto">
        <table className="w-full text-sm border border-border rounded-lg">
          <thead>
            <tr className="bg-muted/50">
              <th className="text-left px-4 py-2 font-medium">Flag</th>
              <th className="text-left px-4 py-2 font-medium">Default</th>
              <th className="text-left px-4 py-2 font-medium">Description</th>
            </tr>
          </thead>
          <tbody>
            {flags.map((f) => (
              <tr key={f.flag} className="border-t border-border">
                <td className="px-4 py-2 font-mono text-primary">{f.flag}</td>
                <td className="px-4 py-2 font-mono text-muted-foreground">{f.def}</td>
                <td className="px-4 py-2 text-muted-foreground">{f.desc}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </section>
  );
}
