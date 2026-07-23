import DocsLayout from "@/components/docs/DocsLayout";
import CodeBlock from "@/components/docs/CodeBlock";
import TerminalDemo from "@/components/docs/TerminalDemo";
import { Stethoscope, CheckCircle2, AlertTriangle, Terminal, Wrench } from "lucide-react";

const CHECKS = [
  { name: "Repo Path", desc: "Verifies the configured repository root path exists and is accessible", severity: "error" },
  { name: "Active Binary", desc: "Locates gitmap on PATH and prints the resolved absolute path + version", severity: "info" },
  { name: "Deployed Binary", desc: "Reads powershell.json to verify the deployed binary matches expectations", severity: "error" },
  { name: "Version Mismatch", desc: "Compares the running binary version against the deployed version", severity: "warn" },
  { name: "Git", desc: "Checks that git is installed and prints its version", severity: "error" },
  { name: "Go", desc: "Checks that go is installed and prints its version", severity: "warn" },
  { name: "Changelog", desc: "Warns if CHANGELOG.md is missing from the repository root", severity: "warn" },
  { name: "Config File", desc: "Validates config.json exists and contains valid JSON", severity: "warn" },
  { name: "Database", desc: "Opens the SQLite database and runs migrations to verify integrity", severity: "error" },
  { name: "Lock File", desc: "Reports if a stale process lock file exists in the data directory", severity: "warn" },
  { name: "Network", desc: "Tests basic connectivity with a TCP dial to verify online status", severity: "info" },
];

const terminalLines = [
  { text: "gitmap doctor", type: "input" as const, delay: 800 },
  { text: "", type: "output" as const },
  { text: "  gitmap doctor v2.27.0", type: "header" as const },
  { text: "  ════════════════════════", type: "output" as const },
  { text: "", type: "output" as const },
  { text: "  ✓  Repo path          /home/user/repos", type: "accent" as const, delay: 200 },
  { text: "  ✓  Active binary      /usr/local/bin/gitmap (v2.27.0)", type: "accent" as const, delay: 150 },
  { text: "  ✓  Deployed binary    /opt/gitmap-v28/gitmap-v28", type: "accent" as const, delay: 150 },
  { text: "  ✓  Version match      running = deployed", type: "accent" as const, delay: 150 },
  { text: "  ✓  Git                git version 2.43.0", type: "accent" as const, delay: 150 },
  { text: "  ✓  Go                 go1.22.2", type: "accent" as const, delay: 150 },
  { text: "  ✓  Changelog          CHANGELOG.md found", type: "accent" as const, delay: 150 },
  { text: "  ✓  Config             config.json valid", type: "accent" as const, delay: 150 },
  { text: "  ✓  Database           gitmap.db OK (migrations current)", type: "accent" as const, delay: 150 },
  { text: "  ✓  Lock file          no stale lock", type: "accent" as const, delay: 150 },
  { text: "  ✓  Network            online", type: "accent" as const, delay: 150 },
  { text: "", type: "output" as const },
  { text: "  All checks passed ✓", type: "accent" as const },
];

const terminalFixLines = [
  { text: "gitmap doctor --fix-path", type: "input" as const, delay: 800 },
  { text: "", type: "output" as const },
  { text: "  Updating PATH to include gitmap binary location...", type: "output" as const, delay: 400 },
  { text: "  ✓  PATH updated successfully", type: "accent" as const },
];

const DoctorPage = () => (
  <DocsLayout>
    <div className="max-w-4xl space-y-10">
      {/* Header */}
      <div>
        <div className="flex items-center gap-3 mb-2">
          <Stethoscope className="h-8 w-8 text-primary" />
          <h1 className="text-3xl font-bold tracking-tight">Doctor</h1>
        </div>
        <p className="text-lg text-muted-foreground">
          Diagnose your gitmap installation with 11 automated health checks.
        </p>
      </div>

      {/* Terminal Demo */}
      <section>
        <TerminalDemo title="gitmap doctor — all checks passing" lines={terminalLines} autoPlay />
      </section>

      {/* Overview */}
      <section>
        <h2 className="text-xl font-semibold mb-3 flex items-center gap-2">
          <CheckCircle2 className="h-5 w-5 text-primary" /> Overview
        </h2>
        <p className="text-muted-foreground mb-4">
          The <code className="text-primary">doctor</code> command runs a comprehensive suite of checks
          against your environment, configuration, database, and network connectivity. It reports issues
          with actionable fix suggestions and exits with a summary count.
        </p>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          {[
            { icon: Stethoscope, title: "11 Checks", desc: "Covers binary, Git, Go, config, database, lock, and network" },
            { icon: Wrench, title: "--fix-path", desc: "Auto-fix PATH issues with the --fix-path flag" },
            { icon: AlertTriangle, title: "Actionable Output", desc: "Each issue includes a specific fix suggestion" },
          ].map((f) => (
            <div key={f.title} className="rounded-lg border border-border p-4 bg-card">
              <f.icon className="h-5 w-5 text-primary mb-2" />
              <h3 className="font-semibold text-sm mb-1">{f.title}</h3>
              <p className="text-xs text-muted-foreground">{f.desc}</p>
            </div>
          ))}
        </div>
      </section>

      {/* Checks Reference */}
      <section>
        <h2 className="text-xl font-semibold mb-3 flex items-center gap-2">
          <Terminal className="h-5 w-5 text-primary" /> Health Checks
        </h2>
        <div className="overflow-x-auto">
          <table className="w-full text-sm border border-border rounded-lg">
            <thead>
              <tr className="bg-muted/50">
                <th className="text-left px-4 py-2 font-medium">Check</th>
                <th className="text-left px-4 py-2 font-medium">Severity</th>
                <th className="text-left px-4 py-2 font-medium">Description</th>
              </tr>
            </thead>
            <tbody>
              {CHECKS.map((c) => (
                <tr key={c.name} className="border-t border-border">
                  <td className="px-4 py-2 font-mono text-primary text-xs">{c.name}</td>
                  <td className="px-4 py-2">
                    <span className={`text-xs font-mono px-1.5 py-0.5 rounded border ${
                      c.severity === "error" ? "bg-destructive/10 text-foreground border-destructive/20 dark:bg-destructive/20 dark:text-destructive" :
                      c.severity === "warn" ? "bg-yellow-500/10 text-yellow-700 border-yellow-500/30 dark:bg-yellow-500/20 dark:text-yellow-100 dark:border-yellow-400/40" :
                      "bg-primary/10 text-foreground border-primary/20 dark:bg-primary/15 dark:text-primary"
                    }`}>
                      {c.severity}
                    </span>
                  </td>
                  <td className="px-4 py-2 text-muted-foreground text-xs">{c.desc}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </section>

      {/* Flags */}
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
              <tr className="border-t border-border">
                <td className="px-4 py-2 font-mono text-primary">--fix-path</td>
                <td className="px-4 py-2 font-mono text-muted-foreground">false</td>
                <td className="px-4 py-2 text-muted-foreground">Auto-fix PATH to include the gitmap binary directory</td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>

      {/* Fix Path Demo */}
      <section>
        <h2 className="text-xl font-semibold mb-3">Auto-Fix PATH</h2>
        <TerminalDemo title="gitmap doctor --fix-path" lines={terminalFixLines} autoPlay />
      </section>

      {/* Examples */}
      <section>
        <h2 className="text-xl font-semibold mb-3">Examples</h2>
        <CodeBlock code={`# Run all diagnostic checks
gitmap doctor

# Auto-fix PATH issues
gitmap doctor --fix-path`} />
      </section>

      {/* See Also */}
      <section>
        <h2 className="text-xl font-semibold mb-3">See Also</h2>
        <ul className="space-y-1 text-sm">
          {[
            { label: "gitmap setup", href: "/getting-started" },
            { label: "gitmap update", href: "/commands" },
            { label: "Configuration", href: "/config" },
          ].map((link) => (
            <li key={link.label}>
              <a href={link.href} className="text-primary hover:underline font-mono text-xs">{link.label}</a>
            </li>
          ))}
        </ul>
      </section>
    </div>
  </DocsLayout>
);

export default DoctorPage;
