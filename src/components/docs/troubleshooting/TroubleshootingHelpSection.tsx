export const TroubleshootingHelpSection = () => {
  return (
    <aside className="mt-10 rounded-lg border border-border bg-muted/20 p-5">
      <h2 className="font-heading font-semibold text-foreground mb-2">Still stuck?</h2>
      <ul className="list-disc pl-5 space-y-1 text-sm text-muted-foreground">
        <li>
          Run <code className="docs-inline-code">gitmap doctor</code> for a full health snapshot.
        </li>
        <li>
          Re-run the failing command with <code className="docs-inline-code">--verbose</code> to generate a timestamped debug log.
        </li>
        <li>
          Check the <a href="/post-mortems" className="text-primary hover:underline">post-mortems</a> for past incidents and their resolutions.
        </li>
      </ul>
    </aside>
  );
};

export default TroubleshootingHelpSection;
