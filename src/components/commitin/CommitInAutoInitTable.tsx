import { commitInAutoInit as autoInit } from "@/pages/commitInData";

export function CommitInAutoInitTable() {
  return (
    <section>
      <h2 className="text-xl font-semibold mb-3">How &lt;source&gt; auto-init works</h2>
      <p className="text-sm text-muted-foreground mb-3">
        You never have to <code>git init</code> first. <code>commit-in</code> resolves
        <code> &lt;source&gt;</code> through a fixed dispatch table — no prompts, no flags:
      </p>
      <div className="overflow-x-auto">
        <table className="w-full text-sm border border-border rounded-lg">
          <thead>
            <tr className="bg-muted/50">
              <th className="text-left px-4 py-2 font-medium">If &lt;source&gt; is…</th>
              <th className="text-left px-4 py-2 font-medium">commit-in does…</th>
            </tr>
          </thead>
          <tbody>
            {autoInit.map((row) => (
              <tr key={row.when} className="border-t border-border">
                <td className="px-4 py-2 text-muted-foreground">{row.when}</td>
                <td className="px-4 py-2 font-mono text-xs">{row.then}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </section>
  );
}
