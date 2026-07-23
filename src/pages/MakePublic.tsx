import DocsLayout from "@/components/docs/DocsLayout";
import CodeBlock from "@/components/docs/CodeBlock";
import { Globe, ShieldCheck, Eye } from "lucide-react";

const flags = [
  { flag: "-y, --yes", def: "false", desc: "Skip the private->public confirmation prompt" },
  { flag: "--dry-run", def: "false", desc: "Print the provider command that would run; do not invoke it" },
  { flag: "--verbose", def: "false", desc: "Echo every shell command to stderr before running it" },
];

const exitCodes = [
  { code: "0", meaning: "Success (or already public)" },
  { code: "2", meaning: "Not inside a git repository" },
  { code: "3", meaning: "No origin remote configured" },
  { code: "4", meaning: "Unsupported provider host or unparseable owner/repo" },
  { code: "5", meaning: "Provider CLI missing, not authenticated, or apply failed" },
  { code: "6", meaning: "Bad flag" },
  { code: "7", meaning: "Confirmation required (re-run with --yes)" },
  { code: "8", meaning: "Verification failed (visibility did not change)" },
];

const MakePublicPage = () => (
  <DocsLayout>
    <div className="max-w-4xl space-y-10">
      <div>
        <div className="flex items-center gap-3 mb-2">
          <Globe className="h-8 w-8 text-primary" />
          <h1 className="text-3xl font-bold tracking-tight">make-public</h1>
        </div>
        <p className="text-lg text-muted-foreground">
          Make the current repository <strong>public</strong> on GitHub or GitLab. Detects the
          provider from <code>origin</code>, requires the matching CLI (<code>gh</code> /{" "}
          <code>glab</code>) already authenticated, prompts before flipping visibility, and
          re-reads visibility to verify the change took effect.
        </p>
      </div>

      <section>
        <h2 className="text-xl font-semibold mb-3">Overview</h2>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          {[
            { icon: Globe, title: "Provider auto-detect", desc: "Parses GitHub or GitLab from origin; uses the matching official CLI." },
            { icon: ShieldCheck, title: "No tokens stored", desc: "Reuses the existing gh / glab auth — gitmap never touches credentials." },
            { icon: Eye, title: "Verifies after apply", desc: "Re-reads visibility post-edit and exits non-zero if the flip didn't stick." },
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
        <CodeBlock code={`gitmap make-public [--yes] [--dry-run] [--verbose]`} />
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
        <CodeBlock code={`# Interactive (will prompt for confirmation)
gitmap make-public

# Non-interactive (CI / scripts)
gitmap make-public --yes

# Preview without touching the provider API
gitmap make-public --dry-run

# Debug auth or argv issues
gitmap make-public --yes --verbose`} />
      </section>

      <section>
        <h2 className="text-xl font-semibold mb-3">Requirements</h2>
        <ul className="list-disc list-inside space-y-2 text-muted-foreground text-sm">
          <li><code>gh</code> (GitHub) or <code>glab</code> (GitLab) installed and on <code>PATH</code>.</li>
          <li>Already authenticated: <code>gh auth login</code> or <code>glab auth login</code>.</li>
          <li>The current directory is inside a git repo with a recognizable <code>origin</code>.</li>
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
          <li><a href="/clone-fix-repo" className="text-primary hover:underline">clone-fix-repo-pub (cfrp)</a> — Clone + fix-repo --all + make-public --yes in one go</li>
        </ul>
      </section>
    </div>
  </DocsLayout>
);

export default MakePublicPage;
