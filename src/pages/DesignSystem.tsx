import DocsLayout from "@/components/docs/DocsLayout";
import CodeBlock from "@/components/docs/CodeBlock";
import { useState } from "react";
import { Copy, Check } from "lucide-react";
import { copyToClipboard } from "@/lib/clipboard";

const colorTokens = [
  { token: "--background", light: "220 20% 97%", dark: "220 25% 6%", usage: "Page background" },
  { token: "--foreground", light: "220 25% 10%", dark: "220 10% 90%", usage: "Default text" },
  { token: "--card", light: "0 0% 100%", dark: "220 25% 9%", usage: "Card surfaces" },
  { token: "--card-foreground", light: "220 25% 10%", dark: "220 10% 90%", usage: "Card text" },
  { token: "--primary", light: "142 71% 45%", dark: "142 71% 45%", usage: "Brand green" },
  { token: "--primary-foreground", light: "220 25% 5%", dark: "220 25% 5%", usage: "Text on primary" },
  { token: "--secondary", light: "220 14% 92%", dark: "220 20% 14%", usage: "Secondary surfaces" },
  { token: "--secondary-foreground", light: "220 25% 10%", dark: "220 10% 90%", usage: "Secondary text" },
  { token: "--muted", light: "220 14% 92%", dark: "220 20% 14%", usage: "Muted backgrounds" },
  { token: "--muted-foreground", light: "220 10% 46%", dark: "220 10% 55%", usage: "Hint text" },
  { token: "--accent", light: "142 71% 45%", dark: "142 71% 45%", usage: "Accent color" },
  { token: "--destructive", light: "0 84% 60%", dark: "0 62% 30%", usage: "Error/danger" },
  { token: "--border", light: "220 13% 87%", dark: "220 20% 16%", usage: "Borders" },
  { token: "--ring", light: "142 71% 45%", dark: "142 71% 45%", usage: "Focus ring" },
  { token: "--terminal", light: "220 25% 8%", dark: "220 25% 5%", usage: "Terminal bg" },
  { token: "--terminal-foreground", light: "142 71% 55%", dark: "142 71% 55%", usage: "Terminal text" },
  { token: "--code-bg", light: "220 20% 94%", dark: "220 25% 10%", usage: "Inline code bg" },
];

const sidebarTokens = [
  { token: "--sidebar-background", light: "220 20% 95%", dark: "220 25% 8%" },
  { token: "--sidebar-foreground", light: "220 10% 30%", dark: "220 10% 75%" },
  { token: "--sidebar-primary", light: "142 71% 45%", dark: "142 71% 45%" },
  { token: "--sidebar-accent", light: "220 14% 90%", dark: "220 20% 12%" },
  { token: "--sidebar-border", light: "220 13% 87%", dark: "220 20% 14%" },
];

const syntaxTokens = [
  { cls: ".hljs-keyword", color: "280 70% 70%", label: "Keywords" },
  { cls: ".hljs-string", color: "142 60% 55%", label: "Strings" },
  { cls: ".hljs-number", color: "30 90% 65%", label: "Numbers" },
  { cls: ".hljs-comment", color: "220 10% 50%", label: "Comments" },
  { cls: ".hljs-function", color: "200 80% 65%", label: "Functions" },
  { cls: ".hljs-title", color: "50 85% 65%", label: "Titles" },
  { cls: ".hljs-built_in", color: "170 60% 55%", label: "Built-ins" },
  { cls: ".hljs-type", color: "30 80% 65%", label: "Types" },
  { cls: ".hljs-variable", color: "10 80% 65%", label: "Variables" },
  { cls: ".hljs-attr", color: "170 60% 55%", label: "Attributes" },
];

function ColorSwatch({ hsl, label, token }: { hsl: string; label?: string; token?: string }) {
  const [copied, setCopied] = useState(false);
  const cssValue = `hsl(${hsl})`;

  const handleCopy = async () => {
    await copyToClipboard(token ? `var(${token})` : hsl);
    setCopied(true);
    setTimeout(() => setCopied(false), 1200);
  };

  return (
    <button
      onClick={handleCopy}
      className="group flex flex-col items-center gap-1.5 cursor-pointer"
      title={`Click to copy ${token || hsl}`}
    >
      <div
        className="w-12 h-12 rounded-lg border border-border shadow-sm group-hover:scale-110 transition-transform relative"
        style={{ backgroundColor: cssValue }}
      >
        {copied && (
          <div className="absolute inset-0 flex items-center justify-center bg-background/80 rounded-lg">
            <Check className="h-4 w-4 text-primary" />
          </div>
        )}
      </div>
      {label && <span className="text-[10px] font-mono text-muted-foreground leading-tight text-center max-w-[72px] truncate">{label}</span>}
    </button>
  );
}

const goExample = `package main

import "fmt"

func main() {
    repos := scanAll("~/projects")
    for _, r := range repos {
        fmt.Printf("%-30s %s\\n", r.Name, r.Branch)
    }
}`;

const bashExample = `# Install gitmap-v28
curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/gitmap-v28/scripts/install.sh | bash

# Scan all repos
gitmap scan ~/projects --format table

# Watch for changes
gitmap watch ~/projects --interval 5s`;

const jsonExample = `{
  "scan_paths": ["~/projects", "~/work"],
  "default_format": "table",
  "watch": {
    "interval": "5s",
    "notify": true
  }
}`;

const DesignSystemPage = () => {
  const isDark = document.documentElement.classList.contains("dark");

  return (
    <DocsLayout>
      <div className="space-y-10">
        {/* Header */}
        <div>
          <h1 className="docs-h1">Design System</h1>
          <p className="text-muted-foreground text-lg mt-2">
            Interactive reference for colors, typography, and component patterns.
            All tokens are defined as CSS variables in <code className="docs-inline-code">src/index.css</code> —
            change once, update everywhere.
          </p>
        </div>

        {/* Color Tokens */}
        <section className="space-y-4">
          <h2 className="docs-h2">Color Tokens</h2>
          <p className="text-sm text-muted-foreground">
            Click any swatch to copy its CSS variable. All values are HSL without the <code className="docs-inline-code">hsl()</code> wrapper.
          </p>

          <div className="border border-border rounded-lg overflow-hidden">
            <div className="bg-muted/30 px-4 py-2 border-b border-border">
              <span className="text-xs font-mono text-muted-foreground">Core Tokens — showing {isDark ? "dark" : "light"} mode values</span>
            </div>
            <div className="overflow-x-auto">
              <table className="w-full text-sm">
                <thead>
                  <tr className="border-b border-border bg-muted/20">
                    <th className="text-left px-4 py-2 font-mono text-xs text-muted-foreground">Token</th>
                    <th className="text-left px-4 py-2 font-mono text-xs text-muted-foreground">Swatch</th>
                    <th className="text-left px-4 py-2 font-mono text-xs text-muted-foreground">HSL</th>
                    <th className="text-left px-4 py-2 font-mono text-xs text-muted-foreground">Usage</th>
                  </tr>
                </thead>
                <tbody>
                  {colorTokens.map((t) => (
                    <tr key={t.token} className="border-b border-border/50 hover:bg-muted/20 transition-colors">
                      <td className="px-4 py-2 font-mono text-xs text-foreground">{t.token}</td>
                      <td className="px-4 py-2">
                        <div className="flex gap-2">
                          <ColorSwatch hsl={t.light} label="Light" token={t.token} />
                          <ColorSwatch hsl={t.dark} label="Dark" token={t.token} />
                        </div>
                      </td>
                      <td className="px-4 py-2 font-mono text-xs text-muted-foreground">
                        <div>{t.light}</div>
                        <div className="opacity-60">{t.dark}</div>
                      </td>
                      <td className="px-4 py-2 text-xs text-muted-foreground">{t.usage}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        </section>

        {/* Sidebar Tokens */}
        <section className="space-y-4">
          <h3 className="docs-h3">Sidebar Tokens</h3>
          <div className="flex flex-wrap gap-4">
            {sidebarTokens.map((t) => (
              <div key={t.token} className="flex items-center gap-3 border border-border rounded-lg px-3 py-2">
                <ColorSwatch hsl={isDark ? t.dark : t.light} token={t.token} />
                <div>
                  <div className="font-mono text-xs text-foreground">{t.token}</div>
                  <div className="font-mono text-[10px] text-muted-foreground">{isDark ? t.dark : t.light}</div>
                </div>
              </div>
            ))}
          </div>
        </section>

        {/* Syntax Highlighting */}
        <section className="space-y-4">
          <h2 className="docs-h2">Syntax Token Colors</h2>
          <p className="text-sm text-muted-foreground">
            Highlight.js classes mapped to HSL values on the terminal background.
          </p>
          <div className="rounded-lg p-4" style={{ backgroundColor: "hsl(220 25% 8%)" }}>
            <div className="flex flex-wrap gap-3">
              {syntaxTokens.map((t) => (
                <div key={t.cls} className="flex items-center gap-2">
                  <div
                    className="w-4 h-4 rounded-sm"
                    style={{ backgroundColor: `hsl(${t.color})` }}
                  />
                  <span className="font-mono text-xs" style={{ color: `hsl(${t.color})` }}>
                    {t.label}
                  </span>
                  <span className="font-mono text-[10px] opacity-50" style={{ color: "hsl(220 10% 60%)" }}>
                    {t.cls}
                  </span>
                </div>
              ))}
            </div>
          </div>
        </section>

        {/* Typography */}
        <section className="space-y-4">
          <h2 className="docs-h2">Typography</h2>
          <div className="space-y-3 border border-border rounded-lg p-5">
            <div>
              <span className="text-xs font-mono text-muted-foreground block mb-1">H1 — Ubuntu · gradient</span>
              <h1 className="docs-h1">Heading One</h1>
            </div>
            <div>
              <span className="text-xs font-mono text-muted-foreground block mb-1">H2 — Ubuntu · gradient</span>
              <h2 className="docs-h2">Heading Two</h2>
            </div>
            <div>
              <span className="text-xs font-mono text-muted-foreground block mb-1">H3 — left border accent</span>
              <h3 className="docs-h3">Heading Three</h3>
            </div>
            <div>
              <span className="text-xs font-mono text-muted-foreground block mb-1">Body — Poppins 400</span>
              <p className="text-foreground">The quick brown fox jumps over the lazy dog. Regular body text uses Poppins at normal weight.</p>
            </div>
            <div>
              <span className="text-xs font-mono text-muted-foreground block mb-1">Mono — Ubuntu Mono</span>
              <p className="font-mono text-foreground">git clone repo && cd repo && make build</p>
            </div>
            <div>
              <span className="text-xs font-mono text-muted-foreground block mb-1">Inline code</span>
              <p className="text-foreground">Use <code className="docs-inline-code">gitmap scan</code> to discover repositories.</p>
            </div>
          </div>

          {/* Font Usage Guidelines — when to reach for Poppins vs Ubuntu Mono */}
          <div className="border border-border rounded-lg p-5 space-y-4 bg-card/30">
            <div>
              <h3 className="docs-h3 !mt-0">Font Usage Guidelines</h3>
              <p className="text-sm text-muted-foreground mt-1">
                One rule: <span className="text-foreground font-medium">Poppins for humans, Ubuntu Mono for machines.</span>
                Mono is reserved for things the user would literally type into a shell or copy into a config — never for decorative chrome.
              </p>
            </div>

            <div className="grid md:grid-cols-2 gap-4">
              {/* Poppins column */}
              <div className="border border-border rounded-md p-4 space-y-3">
                <div className="flex items-baseline justify-between">
                  <span className="text-foreground font-medium">Poppins (font-sans)</span>
                  <span className="text-xs text-muted-foreground">Default — no class needed</span>
                </div>
                <ul className="text-sm text-foreground space-y-1.5 list-disc pl-5">
                  <li>Headings, body copy, paragraphs</li>
                  <li>Button labels &amp; nav items</li>
                  <li>Command names &amp; aliases (<code className="docs-inline-code">scan</code>, <code className="docs-inline-code">s</code>) on the docs page</li>
                  <li>Category titles (e.g. "Scanning &amp; Discovery")</li>
                  <li>Flag descriptions, table headers, tooltips</li>
                  <li>Anything a designer would call a "label"</li>
                </ul>
                <div className="text-xs text-muted-foreground pt-2 border-t border-border">
                  <span className="text-primary">✓ Do:</span> <code className="docs-inline-code">&lt;span className="font-sans"&gt;scan&lt;/span&gt;</code>
                </div>
              </div>

              {/* Ubuntu Mono column */}
              <div className="border border-border rounded-md p-4 space-y-3">
                <div className="flex items-baseline justify-between">
                  <span className="text-foreground font-medium font-mono">Ubuntu Mono (font-mono)</span>
                  <span className="text-xs text-muted-foreground">Opt-in only</span>
                </div>
                <ul className="text-sm text-foreground space-y-1.5 list-disc pl-5">
                  <li>Shell commands inside <code className="docs-inline-code">CodeBlock</code> / <code className="docs-inline-code">TerminalDemo</code></li>
                  <li>Inline code (<code className="docs-inline-code">.docs-inline-code</code>)</li>
                  <li>Literal CLI output, JSON, YAML, file paths</li>
                  <li>Diff hunks, error messages copied from a terminal</li>
                  <li>Keyboard shortcuts inside <code className="docs-inline-code">&lt;kbd&gt;</code></li>
                </ul>
                <div className="text-xs text-muted-foreground pt-2 border-t border-border">
                  <span className="text-destructive">✗ Don't:</span> use <code className="docs-inline-code">font-mono</code> for prose, badges, chips, or section titles for "techy vibes".
                </div>
              </div>
            </div>

            <div className="border border-border rounded-md p-4 bg-background/40 space-y-2">
              <p className="text-sm text-foreground font-medium">Quick decision test</p>
              <p className="text-sm text-muted-foreground">
                Ask: <span className="text-foreground italic">"Would the user type this into a terminal exactly as shown?"</span>
                If yes → <code className="docs-inline-code">font-mono</code>. If no → leave it as default Poppins.
              </p>
            </div>
          </div>
        </section>

        {/* Code Blocks */}
        <section className="space-y-4">
          <h2 className="docs-h2">Code Blocks</h2>
          <p className="text-sm text-muted-foreground">
            Terminal-themed with language accent bars, line numbers, hover glow, and click-to-pin.
          </p>

          <div className="space-y-4">
            <div>
              <span className="text-xs font-mono text-muted-foreground block mb-2">Go — accent hsl(195 60% 50%)</span>
              <CodeBlock code={goExample} language="go" title="main.go" />
            </div>
            <div>
              <span className="text-xs font-mono text-muted-foreground block mb-2">Bash — accent hsl(120 40% 55%)</span>
              <CodeBlock code={bashExample} language="bash" title="install.sh" />
            </div>
            <div>
              <span className="text-xs font-mono text-muted-foreground block mb-2">JSON — accent hsl(40 80% 55%)</span>
              <CodeBlock code={jsonExample} language="json" title="config.json" />
            </div>
          </div>
        </section>

        {/* Component Patterns */}
        <section className="space-y-4">
          <h2 className="docs-h2">Component Patterns</h2>

          {/* Cards */}
          <h3 className="docs-h3">Cards</h3>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="border border-border rounded-lg p-4 bg-card">
              <h4 className="font-semibold text-foreground mb-1">Default Card</h4>
              <p className="text-sm text-muted-foreground">Uses <code className="docs-inline-code">bg-card</code> and <code className="docs-inline-code">border-border</code>.</p>
            </div>
            <div className="border border-primary/30 rounded-lg p-4 bg-card">
              <h4 className="font-semibold text-primary mb-1">Accent Card</h4>
              <p className="text-sm text-muted-foreground">Uses <code className="docs-inline-code">border-primary/30</code> for emphasis.</p>
            </div>
          </div>

          {/* Table */}
          <h3 className="docs-h3">Tables</h3>
          <div className="docs-table">
            <table>
              <thead>
                <tr>
                  <th>Flag</th>
                  <th>Default</th>
                  <th>Description</th>
                </tr>
              </thead>
              <tbody>
                <tr><td><code className="docs-inline-code">--format</code></td><td>table</td><td>Output format</td></tr>
                <tr><td><code className="docs-inline-code">--depth</code></td><td>3</td><td>Scan depth</td></tr>
                <tr><td><code className="docs-inline-code">--watch</code></td><td>false</td><td>Enable watch mode</td></tr>
              </tbody>
            </table>
          </div>

          {/* Blockquote */}
          <h3 className="docs-h3">Blockquotes</h3>
          <blockquote className="docs-blockquote">
            All color tokens reference CSS variables. Change the root value and every component updates automatically.
          </blockquote>
        </section>

        {/* Usage Guide */}
        <section className="space-y-4">
          <h2 className="docs-h2">How to Change the Theme</h2>
          <CodeBlock
            code={`/* src/index.css — Change brand color from green to blue */

/* Before */
--primary: 142 71% 45%;

/* After */
--primary: 217 91% 60%;

/* Every button, link, heading gradient, focus ring,
   and accent updates automatically. */`}
            language="css"
            title="Theme change example"
          />
          <div className="border border-border rounded-lg p-4 bg-muted/20">
            <h4 className="font-semibold text-foreground mb-2">Rules</h4>
            <ul className="text-sm text-muted-foreground space-y-1 list-disc list-inside">
              <li>Never use raw colors (<code className="docs-inline-code">text-white</code>, <code className="docs-inline-code">bg-black</code>)</li>
              <li>Always use semantic tokens (<code className="docs-inline-code">text-foreground</code>, <code className="docs-inline-code">bg-primary</code>)</li>
              <li>HSL values omit the <code className="docs-inline-code">hsl()</code> wrapper in CSS variables</li>
              <li><code className="docs-inline-code">--primary</code> is identical in light and dark mode by design</li>
            </ul>
          </div>
        </section>
      </div>
    </DocsLayout>
  );
};

export default DesignSystemPage;
