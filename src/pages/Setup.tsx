import DocsLayout from "@/components/docs/DocsLayout";
import CodeBlock from "@/components/docs/CodeBlock";
import TerminalDemo from "@/components/docs/TerminalDemo";
import { Settings, Terminal, KeyRound, CheckCircle, Zap } from "lucide-react";

const setupDemo = [
  { text: "gitmap setup", type: "input" as const, delay: 800 },
  { text: "", type: "output" as const },
  { text: "  ✓ Git config: core.longpaths = true", type: "accent" as const },
  { text: "  ✓ Git config: diff.tool = vscode", type: "accent" as const },
  { text: "  ✓ Git config: merge.tool = vscode", type: "accent" as const },
  { text: "  ✓ Git config: credential.helper = manager", type: "accent" as const },
  { text: "  – Git config: user.name (already set, skipped)", type: "output" as const },
  { text: "", type: "output" as const },
  { text: "  ✓ Installed 'gcd' shell function → restart your terminal or source your profile", type: "accent" as const },
  { text: "  ✓ Shell completions installed for PowerShell", type: "accent" as const },
  { text: "", type: "output" as const },
  { text: "  Setup complete — 5 applied, 1 skipped", type: "header" as const },
];

const dryRunDemo = [
  { text: "gitmap setup --dry-run", type: "input" as const, delay: 800 },
  { text: "", type: "output" as const },
  { text: "  [dry-run] Would set core.longpaths = true", type: "output" as const },
  { text: "  [dry-run] Would set diff.tool = vscode", type: "output" as const },
  { text: "  [dry-run] Would install gcd shell function", type: "output" as const },
  { text: "  [dry-run] Would install shell completions", type: "output" as const },
  { text: "", type: "output" as const },
  { text: "  No changes made (dry run)", type: "header" as const },
];

const helpText = `# gitmap setup

Apply global Git configuration, install the gcd shell wrapper,
and set up tab completions.

## Usage

    gitmap setup [flags]

## Flags

| Flag        | Default | Description                          |
|-------------|---------|--------------------------------------|
| --dry-run   | false   | Preview changes without applying     |
| --force     | false   | Re-apply even if already configured  |

## Prerequisites

- Git installed and on PATH
- Shell profile writable (~/.bashrc, ~/.zshrc, or $PROFILE)

## What It Does

1. Reads data/git-setup.json for desired Git config values
2. Compares each key against current global Git config
3. Sets only values that differ (skips already-matching)
4. Installs the gcd() shell function for cd navigation
5. Installs shell tab-completions for the active shell

## Examples

### Example 1: First-time setup

    gitmap setup

    ✓ Git config: core.longpaths = true
    ✓ Git config: diff.tool = vscode
    ✓ Installed 'gcd' shell function
    ✓ Shell completions installed for PowerShell
    Setup complete — 4 applied, 0 skipped

### Example 2: Preview changes

    gitmap setup --dry-run

    [dry-run] Would set core.longpaths = true
    [dry-run] Would install gcd shell function
    No changes made (dry run)

## See Also

- cd (go) — Navigate to repos using gcd
- completion — Manage shell completions separately
- config — View/edit gitmap configuration`;

const bashWrapper = `# Installed by gitmap setup into ~/.bashrc
gcd() {
  local dest
  dest="$(gitmap cd "$@")"
  if [ -n "$dest" ] && [ -d "$dest" ]; then
    cd "$dest" || return
  fi
}`;

const zshWrapper = `# Installed by gitmap setup into ~/.zshrc
gcd() {
  local dest
  dest="$(gitmap cd "$@")"
  if [[ -n "$dest" ]] && [[ -d "$dest" ]]; then
    cd "$dest" || return
  fi
}`;

const powershellWrapper = `# Installed by gitmap setup into $PROFILE
function gcd {
  $dest = gitmap cd @args
  if ($dest -and (Test-Path $dest)) {
    Set-Location $dest
  }
}`;

const gitSetupJson = `{
  "core.longpaths": "true",
  "core.autocrlf": "input",
  "diff.tool": "vscode",
  "difftool.vscode.cmd": "code --wait --diff $LOCAL $REMOTE",
  "merge.tool": "vscode",
  "mergetool.vscode.cmd": "code --wait $MERGED",
  "credential.helper": "manager",
  "init.defaultBranch": "main"
}`;

const features = [
  {
    icon: Settings,
    title: "Git Config",
    desc: "Applies global Git settings from a JSON profile — diff/merge tools, credential helpers, core options",
  },
  {
    icon: Terminal,
    title: "gcd Shell Wrapper",
    desc: "Installs the gcd() function so cd navigation works in Bash, Zsh, and PowerShell",
  },
  {
    icon: KeyRound,
    title: "Tab Completions",
    desc: "Auto-installs shell completions for all gitmap commands, repos, groups, and aliases",
  },
  {
    icon: CheckCircle,
    title: "Idempotent",
    desc: "Skips settings already configured; uses markers to prevent duplicate shell function installs",
  },
];

const Setup = () => {
  return (
    <DocsLayout>
      <div className="max-w-4xl mx-auto">
        <h1 className="docs-h1">Setup</h1>
        <p className="text-lg text-muted-foreground mb-8 font-body">
          One command to configure Git globals, install the <code className="docs-inline-code">gcd</code> navigation wrapper,
          and enable shell tab-completions for every gitmap command.
        </p>

        {/* Terminal Preview */}
        <section className="mb-12">
          <h2 className="docs-h2">First-Time Setup</h2>
          <TerminalDemo lines={setupDemo} title="gitmap setup" />
        </section>

        {/* Features */}
        <section className="mb-12">
          <h2 className="docs-h2">What Gets Installed</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            {features.map((f) => (
              <div key={f.title} className="p-4 rounded-lg border border-border bg-card">
                <div className="flex items-center gap-2 mb-2">
                  <f.icon className="h-5 w-5 text-primary" />
                  <h3 className="font-heading font-semibold text-foreground">{f.title}</h3>
                </div>
                <p className="text-sm text-muted-foreground font-body">{f.desc}</p>
              </div>
            ))}
          </div>
        </section>

        {/* Help Output */}
        <section className="mb-12">
          <h2 className="docs-h2">--help Output</h2>
          <CodeBlock code={helpText} language="markdown" title="gitmap setup --help" />
        </section>

        {/* Git Config Profile */}
        <section className="mb-12">
          <h2 className="docs-h2">Git Config Profile</h2>
          <p className="text-muted-foreground mb-4 font-body">
            Setup reads <code className="docs-inline-code">data/git-setup.json</code> and applies each key
            as a global Git config value. Existing values that already match are skipped.
          </p>
          <CodeBlock code={gitSetupJson} language="json" title="data/git-setup.json" />
        </section>

        {/* Shell Wrappers */}
        <section className="mb-12">
          <h2 className="docs-h2">Shell Wrapper: gcd</h2>
          <p className="text-muted-foreground mb-4 font-body">
            The <code className="docs-inline-code">gcd</code> function captures the path printed by{" "}
            <code className="docs-inline-code">gitmap cd</code> and performs the actual directory change.
            Setup detects the active shell and installs the appropriate version.
          </p>

          <div className="space-y-4">
            <CodeBlock code={bashWrapper} language="bash" title="Bash (~/.bashrc)" />
            <CodeBlock code={zshWrapper} language="bash" title="Zsh (~/.zshrc)" />
            <CodeBlock code={powershellWrapper} language="powershell" title="PowerShell ($PROFILE)" />
          </div>

          <div className="mt-4 p-3 rounded-lg border border-border bg-muted/30">
            <p className="text-sm text-muted-foreground font-body">
              <strong>Marker-based idempotency:</strong> Each wrapper is preceded by a{" "}
              <code className="docs-inline-code"># gitmap cd wrapper</code> comment. Setup checks for this
              marker before writing to avoid duplicate installs.
            </p>
          </div>
        </section>

        {/* Dry Run */}
        <section className="mb-12">
          <h2 className="docs-h2">Dry Run Preview</h2>
          <p className="text-muted-foreground mb-4 font-body">
            Use <code className="docs-inline-code">--dry-run</code> to preview all changes without modifying
            any files or Git config.
          </p>
          <TerminalDemo lines={dryRunDemo} title="gitmap setup --dry-run" />
        </section>

        {/* Flags Table */}
        <section className="mb-12">
          <h2 className="docs-h2">Flags</h2>
          <div className="overflow-x-auto">
            <table className="docs-table">
              <thead>
                <tr>
                  <th>Flag</th>
                  <th>Default</th>
                  <th>Description</th>
                </tr>
              </thead>
              <tbody>
                <tr>
                  <td><code className="docs-inline-code">--dry-run</code></td>
                  <td>false</td>
                  <td>Preview changes without applying them</td>
                </tr>
                <tr>
                  <td><code className="docs-inline-code">--force</code></td>
                  <td>false</td>
                  <td>Re-apply all settings even if already configured</td>
                </tr>
              </tbody>
            </table>
          </div>
        </section>

        {/* See Also */}
        <section className="mb-12">
          <h2 className="docs-h2">See Also</h2>
          <ul className="space-y-2 text-muted-foreground font-body">
            <li>
              <a href="/cd" className="text-primary hover:underline font-medium">cd (go)</a>{" "}
              — Navigate to repos using the gcd wrapper
            </li>
            <li>
              <a href="/alias" className="text-primary hover:underline font-medium">alias</a>{" "}
              — Manage repo aliases for quick access
            </li>
            <li>
              <a href="/config" className="text-primary hover:underline font-medium">config</a>{" "}
              — View and edit gitmap configuration
            </li>
            <li>
              <a href="/commands" className="text-primary hover:underline font-medium">commands</a>{" "}
              — Full command reference
            </li>
          </ul>
        </section>
      </div>
    </DocsLayout>
  );
};

export default Setup;
