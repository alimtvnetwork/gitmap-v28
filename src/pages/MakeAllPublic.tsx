import DocsLayout from "@/components/docs/DocsLayout";
import CodeBlock from "@/components/docs/CodeBlock";
import { Globe, ListChecks, Undo2 } from "lucide-react";

const flags = [
  { flag: "-Y, --yes", def: "false", desc: "Skip the interactive confirmation prompt" },
  { flag: "--verbose", def: "false", desc: "Echo every shell command to stderr before running it" },
];

const exitCodes = [
  { code: "0", meaning: "All target repos already public OR all flips succeeded" },
  { code: "4", meaning: "Owner resolution / provider detection failed" },
  { code: "5", meaning: "Provider CLI missing OR not authenticated OR every repo failed" },
  { code: "6", meaning: "Bad flag / missing positional argument" },
  { code: "7", meaning: "User aborted at the confirmation prompt" },
  { code: "9", meaning: "Partial — at least one repo flipped AND at least one failed" },
];

const MakeAllPublicPage = () => (
  <DocsLayout>
    <div className="max-w-4xl space-y-10">
      <div>
        <div className="flex items-center gap-3 mb-2">
          <Globe className="h-8 w-8 text-primary" />
          <h1 className="text-3xl font-bold tracking-tight">make-all-public (MAPUB)</h1>
        </div>
        <p className="text-lg text-muted-foreground">
          Bulk-flip every matching repo on an owner to <strong>public</strong> in one call.
          Use the uppercase shorthand <code>MAPUB</code> for the exact same behavior. Every
          run is recorded in the audit DB so it can be reversed with{" "}
          <code>visibility-undo</code>.
        </p>
      </div>

      <section>
        <h2 className="text-xl font-semibold mb-3">Overview</h2>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          {[
            { icon: Globe, title: "Wildcard matching", desc: "Glob patterns + !negation against every repo under the owner." },
            { icon: ListChecks, title: "Interactive preview", desc: "Numbered match table with y/n/exclude (e.g. 1,3-5) before any flip." },
            { icon: Undo2, title: "Full audit trail", desc: "Every run + per-repo result persists, replayable via vu / vr / vish." },
          ].map((f) => (
            <div key={f.title} className="rounded-lg border border-border p-4 bg-card">
              <f.icon className="h-5 w-5 text-primary mb-2" />
              <h3 className="font-semibold text-sm mb-1">{f.title}</h3>
              <p className="text-xs text-muted-foreground">{f.desc}</p>
            </div>
          ))}
        </div>
      </section>

      <section>
        <h2 className="text-xl font-semibold mb-3">Usage</h2>
        <CodeBlock code={`gitmap make-all-public <owner-or-url> <patterns> [-Y|--yes] [--verbose]
gitmap MAPUB           <owner-or-url> <patterns> [-Y|--yes] [--verbose]`} />
      </section>

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

      <section>
        <h2 className="text-xl font-semibold mb-3">Examples</h2>
        <CodeBlock code={`# Interactive — preview matches, then confirm
gitmap make-all-public alice "demo-*"

# Multiple patterns + a negation
gitmap make-all-public alice "demo-*,proto-*,!proto-secret"

# Non-interactive (CI / scripts)
gitmap make-all-public https://github.com/alice "*" -Y --verbose

# Uppercase shorthand
gitmap MAPUB alice "demo-*"
gitmap MAPUB https://github.com/alice "demo-*,!demo-secret" -Y`} />
      </section>

      <section>
        <h2 className="text-xl font-semibold mb-3">Pattern syntax</h2>
        <ul className="list-disc list-inside space-y-2 text-muted-foreground text-sm">
          <li><code>exact</code> — match a single repo name.</li>
          <li><code>prefix*</code> — match anything starting with <code>prefix</code>.</li>
          <li><code>*contains*</code> — match anything containing the substring.</li>
          <li><code>prefix*suffix</code> — match anything starting AND ending with the parts.</li>
          <li><code>!pattern</code> — exclude matches (applied after include patterns).</li>
        </ul>
      </section>

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

      <section>
        <h2 className="text-xl font-semibold mb-3">See Also</h2>
        <ul className="list-disc list-inside space-y-1 text-sm">
          <li><a href="/make-all-private" className="text-primary hover:underline">make-all-private (MAPRI)</a> — opposite direction, same machinery.</li>
          <li><a href="/make-public" className="text-primary hover:underline">make-public</a> — single current-repo flip.</li>
        </ul>
      </section>
    </div>
  </DocsLayout>
);

export default MakeAllPublicPage;
