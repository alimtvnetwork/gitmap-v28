import DocsLayout from "@/components/docs/DocsLayout";
import CodeBlock from "@/components/docs/CodeBlock";
import TerminalDemo from "@/components/docs/TerminalDemo";
import { Globe, Terminal, Server, MonitorPlay } from "lucide-react";

const features = [
  { icon: Server, title: "Static Serving", desc: "Serves pre-built dist/ files via Go's built-in HTTP server — no Node.js required." },
  { icon: Terminal, title: "Dev Mode Fallback", desc: "Falls back to npm install && npm run dev when no dist/ folder is found." },
  { icon: Globe, title: "Auto-Open Browser", desc: "Automatically opens the docs site in your default browser on launch." },
  { icon: MonitorPlay, title: "Configurable Port", desc: "Serve on any port with --port (default: 5173)." },
];

const flags = [
  ["--port <number>", "Port to serve the dashboard on", "5173"],
];

const fileLayout = [
  ["constants/constants_helpdashboard.go", "CLI names, flag descriptions, messages, defaults"],
  ["cmd/helpdashboard.go", "Command handler, static server, dev fallback, browser open"],
  ["helptext/help-dashboard.md", "Embedded help text for the command"],
];

const demoStaticLines = [
  { text: "gitmap help-dashboard", type: "input" as const, delay: 800 },
  { text: "  Serving docs from /usr/local/bin/docs-site/dist on http://localhost:5173", type: "accent" as const },
  { text: "  Opening http://localhost:5173 in browser...", type: "output" as const },
  { text: "", type: "output" as const, delay: 1200 },
  { text: "  ^C", type: "input" as const, delay: 2000 },
  { text: "", type: "output" as const },
  { text: "  Server stopped.", type: "output" as const },
];

const demoDevLines = [
  { text: "gitmap hd", type: "input" as const, delay: 800 },
  { text: "  No pre-built dist/ found, falling back to npm run dev", type: "output" as const },
  { text: "  Running npm install...", type: "output" as const, delay: 600 },
  { text: "  Starting dev server from /usr/local/bin/docs-site...", type: "accent" as const, delay: 800 },
  { text: "  Opening http://localhost:5173 in browser...", type: "output" as const },
  { text: "", type: "output" as const, delay: 1200 },
  { text: "  ^C", type: "input" as const, delay: 2000 },
  { text: "", type: "output" as const },
  { text: "  Server stopped.", type: "output" as const },
];

const demoPortLines = [
  { text: "gitmap hd --port 8080", type: "input" as const, delay: 800 },
  { text: "  Serving docs from /usr/local/bin/docs-site/dist on http://localhost:8080", type: "accent" as const },
  { text: "  Opening http://localhost:8080 in browser...", type: "output" as const },
];

const HelpDashboardPage = () => (
  <DocsLayout>
    <h1 className="text-3xl font-heading font-bold mb-2 docs-h1">Help Dashboard Command</h1>
    <p className="text-muted-foreground mb-6">
      Serve the interactive documentation site locally in your browser.
      <span className="ml-2 text-xs font-mono bg-muted px-2 py-0.5 rounded">alias: hd</span>
    </p>

    <h2 className="text-xl font-heading font-semibold mt-8 mb-4">Features</h2>
    <div className="grid md:grid-cols-2 gap-4 mb-8">
      {features.map((f) => (
        <div key={f.title} className="rounded-lg border border-border bg-card p-4">
          <f.icon className="h-5 w-5 text-primary mb-2" />
          <h3 className="font-mono font-semibold text-sm mb-1">{f.title}</h3>
          <p className="text-xs text-muted-foreground">{f.desc}</p>
        </div>
      ))}
    </div>

    <h2 className="text-xl font-heading font-semibold mt-10 mb-3">How It Works</h2>
    <div className="rounded-lg border border-border bg-card p-4 mb-8 text-sm text-muted-foreground space-y-2">
      <p><strong className="text-foreground">1.</strong> Locates the <code className="text-primary">docs-site/</code> directory relative to the gitmap binary.</p>
      <p><strong className="text-foreground">2.</strong> If a pre-built <code className="text-primary">dist/</code> folder exists, serves it with Go's built-in HTTP server (no dependencies).</p>
      <p><strong className="text-foreground">3.</strong> If no <code className="text-primary">dist/</code> is found, falls back to <code className="text-primary">npm install && npm run dev</code>.</p>
      <p><strong className="text-foreground">4.</strong> Opens the dashboard in your default browser automatically.</p>
    </div>

    <h2 className="text-xl font-heading font-semibold mt-10 mb-3">Interactive Examples</h2>
    <div className="space-y-6 mb-8">
      <TerminalDemo title="Static mode (pre-built dist/)" lines={demoStaticLines} autoPlay />
      <TerminalDemo title="Dev mode fallback (no dist/)" lines={demoDevLines} />
      <TerminalDemo title="Custom port" lines={demoPortLines} />
    </div>

    <h2 className="text-xl font-heading font-semibold mt-10 mb-3">Usage</h2>
    <CodeBlock code="gitmap help-dashboard [flags]" title="Basic usage" />
    <CodeBlock code="gitmap hd --port 8080" title="Custom port" />

    <h2 className="text-xl font-heading font-semibold mt-10 mb-3">Flags</h2>
    <div className="rounded-lg border border-border overflow-hidden">
      <table className="w-full text-sm">
        <thead>
          <tr className="bg-muted/50">
            <th className="text-left font-mono font-semibold px-4 py-2">Flag</th>
            <th className="text-left font-mono font-semibold px-4 py-2">Description</th>
            <th className="text-left font-mono font-semibold px-4 py-2">Default</th>
          </tr>
        </thead>
        <tbody>
          {flags.map(([flag, desc, def]) => (
            <tr key={flag} className="border-t border-border">
              <td className="px-4 py-2 font-mono text-primary">{flag}</td>
              <td className="px-4 py-2 text-muted-foreground">{desc}</td>
              <td className="px-4 py-2 text-muted-foreground">{def}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>

    <h2 className="text-xl font-heading font-semibold mt-10 mb-3">Prerequisites</h2>
    <div className="rounded-lg border border-border overflow-hidden">
      <table className="w-full text-sm">
        <thead>
          <tr className="bg-muted/50">
            <th className="text-left font-mono font-semibold px-4 py-2">Mode</th>
            <th className="text-left font-mono font-semibold px-4 py-2">Requirements</th>
          </tr>
        </thead>
        <tbody>
          <tr className="border-t border-border">
            <td className="px-4 py-2 font-mono text-primary">Static</td>
            <td className="px-4 py-2 text-muted-foreground">None — serves pre-built files directly</td>
          </tr>
          <tr className="border-t border-border">
            <td className="px-4 py-2 font-mono text-primary">Dev (fallback)</td>
            <td className="px-4 py-2 text-muted-foreground">Node.js and npm on PATH</td>
          </tr>
        </tbody>
      </table>
    </div>

    <h2 className="text-xl font-heading font-semibold mt-10 mb-3">File Layout</h2>
    <div className="rounded-lg border border-border overflow-hidden">
      <table className="w-full text-sm">
        <thead>
          <tr className="bg-muted/50">
            <th className="text-left font-mono font-semibold px-4 py-2">File</th>
            <th className="text-left font-mono font-semibold px-4 py-2">Purpose</th>
          </tr>
        </thead>
        <tbody>
          {fileLayout.map(([file, purpose]) => (
            <tr key={file} className="border-t border-border">
              <td className="px-4 py-2 font-mono text-primary">{file}</td>
              <td className="px-4 py-2 text-muted-foreground">{purpose}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>

    <h2 className="text-xl font-heading font-semibold mt-10 mb-3">See Also</h2>
    <ul className="list-disc list-inside text-sm text-muted-foreground space-y-1">
      <li><a href="/help-index" className="text-primary hover:underline">help-index</a> — Curated index of diff / mv / merge / commit pages served by this dashboard</li>
      <li><code className="text-primary">docs</code> — Open the hosted documentation website</li>
      <li><code className="text-primary">dashboard</code> — Generate an HTML analytics dashboard for a repo</li>
    </ul>
  </DocsLayout>
);

export default HelpDashboardPage;
