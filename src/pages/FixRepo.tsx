import DocsLayout from "@/components/docs/DocsLayout";
import CodeBlock from "@/components/docs/CodeBlock";
import { Wrench, GitBranch, ShieldCheck, FileCode } from "lucide-react";

const flags = [
  { flag: "-2 / -3 / -5", def: "-2", desc: "How many prior versions to rewrite (window width)" },
  { flag: "--all", def: "false", desc: "Rewrite every prior version (v1..vK-1 -> vK)" },
  { flag: "--dry-run", def: "false", desc: "Print modifications without writing files" },
  { flag: "--verbose", def: "false", desc: "Per-file diff output" },
  { flag: "--strict", def: "false", desc: "After rewrite + gofmt, run go test on touched Go packages" },
  { flag: "--config <path>", def: "(auto)", desc: "Override fix-repo.config.json location" },
];

const exitCodes = [
  { code: "0", meaning: "Success" },
  { code: "2", meaning: "Not a git repository" },
  { code: "3", meaning: "No origin remote configured" },
  { code: "4", meaning: "Repo name has no -vN suffix" },
  { code: "5", meaning: "Bad / unparseable version suffix" },
  { code: "6", meaning: "Bad flag" },
  { code: "7", meaning: "Write failed" },
  { code: "8", meaning: "Bad config file" },
  { code: "9", meaning: "Tests failed (--strict only)" },
];

const FixRepoPage = () => (
  <DocsLayout>
    <div className="max-w-4xl space-y-10">
      <div>
        <div className="flex items-center gap-3 mb-2">
          <Wrench className="h-8 w-8 text-primary" />
          <h1 className="text-3xl font-bold tracking-tight">fix-repo</h1>
          <span className="font-mono text-xs px-2 py-1 rounded bg-primary/10 text-foreground border border-primary/20 dark:bg-primary/15 dark:text-primary dark:border-primary/40">
            alias: fr
          </span>
        </div>
        <p className="text-lg text-muted-foreground">
          Rewrite prior <code>{`{base}-vN`}</code> versioned-repo-name tokens in every tracked
          text file to the current version. Go-native re-implementation of{" "}
          <code>fix-repo.ps1</code> / <code>fix-repo.sh</code> with byte-identical exit codes
          and config schema.
        </p>
        <p className="text-xs text-muted-foreground mt-2">
          Spec: <code>spec/04-generic-cli/27-fix-repo-command.md</code>
        </p>
      </div>

      <section>
        <h2 className="text-xl font-semibold mb-3">Overview</h2>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          {[
            { icon: GitBranch, title: "Reads identity from git", desc: "Repo name must end with -vN. The current N is parsed from origin." },
            { icon: FileCode, title: "Negative-lookahead guard", desc: "-v1 will never match inside -v10 / -v18 — width-crossing safe." },
            { icon: ShieldCheck, title: "Strict mode", desc: "--strict runs go test on touched packages to catch semantic desyncs." },
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
        <CodeBlock code={`gitmap fix-repo [-2 | -3 | -5 | --all] [--dry-run] [--verbose] [--strict] [--config <path>]
gitmap fr                                                                      # short alias

# PowerShell-style flags also accepted: -DryRun -Verbose -Strict -All -Config <p>`} />
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
        <CodeBlock code={`# Preview the default 2-version window
gitmap fix-repo --dry-run --verbose

# Rewrite every prior version, then run tests on touched Go packages
gitmap fix-repo --all --strict

# Wider window (last 5 prior versions)
gitmap fr -5`} />
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
          <li><a href="/clone-fix-repo" className="text-primary hover:underline">clone-fix-repo</a> — Clone + fix-repo --all in one shot</li>
          <li><a href="/replace" className="text-primary hover:underline">replace</a> — Generic literal find/replace</li>
        </ul>
      </section>
    </div>
  </DocsLayout>
);

export default FixRepoPage;
