// CommitInWhatItCreates renders the "What it actually creates" overview
// table for the commit-in docs page. Extracted from CommitIn.tsx so the
// page component stays under the project-wide <200-lines rule.
const CommitInWhatItCreates = () => (
  <section>
    <h2 className="text-xl font-semibold mb-3">What it actually creates</h2>
    <p className="text-sm text-muted-foreground mb-3">
      Think of <code>commit-in</code> as a one-shot way to take whatever you point
      it at — a path that exists, a path that <em>doesn't</em> exist yet, a remote
      URL, or a stack of <code>foo-v1 … foo-vN</code> sibling folders — and turn
      it into <strong>one</strong> git repo whose history is the chronological
      union of every input.
    </p>
    <div className="overflow-x-auto mb-3">
      <table className="w-full text-sm border border-border rounded-lg">
        <thead>
          <tr className="bg-muted/50">
            <th className="text-left px-4 py-2 font-medium">You run…</th>
            <th className="text-left px-4 py-2 font-medium">commit-in produces…</th>
          </tr>
        </thead>
        <tbody>
          <tr className="border-t border-border">
            <td className="px-4 py-2 font-mono text-xs text-primary">gitmap cin ./new-repo https://github.com/me/old.git</td>
            <td className="px-4 py-2 text-xs text-muted-foreground"><code>./new-repo/</code> is created + <code>git init</code>'d, then every commit from <code>old.git</code> is appended in author-date order.</td>
          </tr>
          <tr className="border-t border-border">
            <td className="px-4 py-2 font-mono text-xs text-primary">gitmap cin ./gitmap all</td>
            <td className="px-4 py-2 text-xs text-muted-foreground">Every <code>./gitmap-v27</code>, <code>-v2</code>, …, <code>-vN</code> sibling on disk is auto-discovered and replayed into <code>./gitmap-v27/</code> oldest&nbsp;→&nbsp;newest.</td>
          </tr>
          <tr className="border-t border-border">
            <td className="px-4 py-2 font-mono text-xs text-primary">gitmap cin ./gitmap -3</td>
            <td className="px-4 py-2 text-xs text-muted-foreground">Same as above, but only the latest 3 siblings (e.g. v17, v18, v19).</td>
          </tr>
          <tr className="border-t border-border">
            <td className="px-4 py-2 font-mono text-xs text-primary">gitmap cin ./canonical a,b,git@host:c.git</td>
            <td className="px-4 py-2 text-xs text-muted-foreground">Mixed inputs — local + remote — interleaved by author date into one timeline in <code>./canonical/</code>.</td>
          </tr>
        </tbody>
      </table>
    </div>
    <p className="text-sm text-muted-foreground">
      The TARGET is <strong>always</strong> a real git repo when commit-in is
      done — created on the fly if it didn't exist. You do not pre-init,
      pre-clone, or pre-clean anything.
    </p>
  </section>
);

export default CommitInWhatItCreates;
