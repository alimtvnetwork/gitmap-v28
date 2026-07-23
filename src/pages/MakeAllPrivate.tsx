import DocsLayout from "@/components/docs/DocsLayout";
import CodeBlock from "@/components/docs/CodeBlock";
import { Lock, ListChecks, Undo2 } from "lucide-react";

const flags = [
  { flag: "-Y, --yes", def: "false", desc: "Skip the interactive confirmation prompt" },
  { flag: "--verbose", def: "false", desc: "Echo every shell command to stderr before running it" },
];

const exitCodes = [
  { code: "0", meaning: "All target repos already private OR all flips succeeded" },
  { code: "4", meaning: "Owner resolution / provider detection failed" },
  { code: "5", meaning: "Provider CLI missing OR not authenticated OR every repo failed" },
  { code: "6", meaning: "Bad flag / missing positional argument" },
  { code: "7", meaning: "User aborted at the confirmation prompt" },
  { code: "9", meaning: "Partial — at least one repo flipped AND at least one failed" },
];

const MakeAllPrivatePage = () => (
  <DocsLayout>
    <div className="max-w-4xl space-y-10">
      <div>
        <div className="flex items-center gap-3 mb-2">
          <Lock className="h-8 w-8 text-primary" />
          <h1 className="text-3xl font-bold tracking-tight">make-all-private (MAPRI)</h1>
        </div>
        <p className="text-lg text-muted-foreground">
          Bulk-flip every matching repo on an owner to <strong>private</strong> in one call.
          Mirror image of <code>make-all-public</code> — same flags, exit codes, audit trail,
          and undo/redo wiring. <code>MAPRI</code> is the uppercase shorthand.
        </p>
      </div>

      <section>
        <h2 className="text-xl font-semibold mb-3">Overview</h2>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          {[
            { icon: Lock, title: "Owner-wide flip", desc: "Lists every repo under <owner-or-url> and matches your patterns." },
            { icon: ListChecks, title: "Confirmation prompt", desc: "Numbered match table; accept with y, abort with n, or exclude 1,3-5." },
            { icon: Undo2, title: "Reversible", desc: "Every run recorded — reverse with visibility-undo (vu)." },
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
        <CodeBlock code={`gitmap make-all-private <owner-or-url> <patterns> [-Y|--yes] [--verbose]
gitmap MAPRI            <owner-or-url> <patterns> [-Y|--yes] [--verbose]`} />
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
        <CodeBlock code={`# Interactive — preview, then confirm
gitmap make-all-private alice "archive-*"

# Multiple patterns + a negation, non-interactive
gitmap make-all-private alice "old-*,!old-keep" -Y

# GitLab owner via URL, verbose for debugging
gitmap make-all-private https://gitlab.com/alice "*" --verbose

# Uppercase shorthand
gitmap MAPRI alice "demo-*"
gitmap MAPRI . "archive-*" --verbose`} />
      </section>

      <section>
        <h2 className="text-xl font-semibold mb-3">Pattern syntax</h2>
        <ul className="list-disc list-inside space-y-2 text-muted-foreground text-sm">
          <li><code>exact</code> — match a single repo name.</li>
          <li><code>prefix*</code> / <code>*contains*</code> / <code>prefix*suffix</code> — glob shapes.</li>
          <li><code>!pattern</code> — exclude (applied after include patterns).</li>
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
          <li><a href="/make-all-public" className="text-primary hover:underline">make-all-public (MAPUB)</a> — opposite direction.</li>
          <li><a href="/make-public" className="text-primary hover:underline">make-public</a> — single current-repo flip.</li>
        </ul>
      </section>
    </div>
  </DocsLayout>
);

export default MakeAllPrivatePage;
