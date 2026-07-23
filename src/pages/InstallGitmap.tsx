import DocsLayout from "@/components/docs/DocsLayout";
import CodeBlock from "@/components/docs/CodeBlock";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Card } from "@/components/ui/card";
import { Download, RefreshCw, Trash2, Stethoscope, Apple, Terminal, MonitorSmartphone } from "lucide-react";

// Single source of truth for these one-liners is
// spec/01-app/108-cross-platform-install-update.md.
// Keep these snippets byte-identical with the spec table so users
// landing on README, web docs, or `--help` see the same commands.

const REPO_RAW = "https://raw.githubusercontent.com/alimtvnetwork/gitmap-v27/main";

const installCmds = {
  windows: `irm ${REPO_RAW}/gitmap-v27/scripts/install.ps1 | iex`,
  unix: `curl -fsSL ${REPO_RAW}/gitmap-v27/scripts/install.sh | sh`,
};

const installPromptCmds = {
  windows: `irm ${REPO_RAW}/install-quick.ps1 | iex`,
  unix: `curl -fsSL ${REPO_RAW}/install-quick.sh | bash`,
};

const installPinnedCmds = {
  windows: `$ver = 'v3.99.0'
$installer = irm ${REPO_RAW}/gitmap-v27/scripts/install.ps1
& ([scriptblock]::Create($installer)) -Version $ver -NoDiscovery`,
  unix: `curl -fsSL ${REPO_RAW}/gitmap-v27/scripts/install.sh \
  | bash -s -- --version v3.99.0 --no-discovery`,
};

const uninstallCmds = {
  windows: `irm ${REPO_RAW}/uninstall-quick.ps1 | iex`,
  unix: `curl -fsSL ${REPO_RAW}/uninstall-quick.sh | bash`,
};

type CmdMap = { windows: string; unix: string };

const PlatformBlock = ({ commands, lang }: { commands: CmdMap; lang?: { windows: string; unix: string } }) => (
  <Tabs defaultValue="windows" className="w-full">
    <TabsList className="grid w-full max-w-md grid-cols-3 mb-3">
      <TabsTrigger value="windows" className="gap-1.5">
        <MonitorSmartphone className="w-3.5 h-3.5" /> Windows
      </TabsTrigger>
      <TabsTrigger value="macos" className="gap-1.5">
        <Apple className="w-3.5 h-3.5" /> macOS
      </TabsTrigger>
      <TabsTrigger value="linux" className="gap-1.5">
        <Terminal className="w-3.5 h-3.5" /> Linux
      </TabsTrigger>
    </TabsList>
    <TabsContent value="windows">
      <CodeBlock code={commands.windows} title="PowerShell 5.1+ / pwsh 7+" language={lang?.windows ?? "powershell"} />
    </TabsContent>
    <TabsContent value="macos">
      <CodeBlock code={commands.unix} title="Terminal (zsh / bash)" language={lang?.unix ?? "bash"} />
    </TabsContent>
    <TabsContent value="linux">
      <CodeBlock code={commands.unix} title="Terminal (bash / zsh / fish)" language={lang?.unix ?? "bash"} />
    </TabsContent>
  </Tabs>
);

const SectionHeader = ({ icon: Icon, title, kicker }: { icon: typeof Download; title: string; kicker?: string }) => (
  <div className="flex items-start gap-3 mb-4">
    <div className="p-2 rounded-lg bg-primary/10 border border-primary/20">
      <Icon className="w-5 h-5 text-primary" />
    </div>
    <div>
      <h2 className="text-xl font-heading font-semibold docs-h2">{title}</h2>
      {kicker ? <p className="text-sm text-muted-foreground mt-0.5">{kicker}</p> : null}
    </div>
  </div>
);

const InstallGitmapPage = () => {
  return (
    <DocsLayout>
      <div className="max-w-4xl">
        <div className="flex items-center gap-3 mb-2">
          <Download className="w-8 h-8 text-primary" />
          <h1 className="text-3xl font-heading font-bold docs-h1">Install &amp; Update gitmap-v27</h1>
        </div>
        <p className="text-muted-foreground mb-8 text-lg">
          Cross-platform install, update, and uninstall reference. Same one-liners as{" "}
          <code className="docs-inline-code">spec/01-app/108-cross-platform-install-update.md</code> and{" "}
          <code className="docs-inline-code">README.md</code>.
        </p>

        {/* Default install */}
        <section className="mb-12">
          <SectionHeader
            icon={Download}
            title="Install — default (recommended)"
            kicker="Canonical install.ps1 / install.sh. No prompts. Sensible install location. This is what 99% of users want."
          />
          <PlatformBlock commands={installCmds} />
        </section>

        {/* Install with folder prompt */}
        <section className="mb-12">
          <SectionHeader
            icon={Download}
            title="Install — quick (pick install drive)"
            kicker="Use only when you want to install on a specific drive (e.g. D:\\). Prompts for drive/folder, then delegates to the canonical installer."
          />
          <PlatformBlock commands={installPromptCmds} />
        </section>

        {/* Pinned install */}
        <section className="mb-12">
          <SectionHeader
            icon={Download}
            title="Install — pinned version"
            kicker="Strict mode: missing tag exits 1. No sibling probe, no `latest`, no HEAD fallback. Use this in CI."
          />
          <PlatformBlock commands={installPinnedCmds} />
          <p className="text-xs text-muted-foreground mt-3">
            Replace <code className="docs-inline-code">v3.99.0</code> with the release tag you need.
            Resolution contract:{" "}
            <code className="docs-inline-code">spec/07-generic-release/09-generic-install-script-behavior.md</code>.
          </p>
        </section>

        {/* Update */}
        <section className="mb-12">
          <SectionHeader
            icon={RefreshCw}
            title="Update"
            kicker="Self-updates from the linked source repo. Falls back to gitmap-updater, then to the manual one-liner."
          />
          <Card className="p-4 mb-4">
            <h3 className="text-sm font-semibold mb-2 docs-h3">In-place update</h3>
            <CodeBlock code={`gitmap update`} title="Any platform" language="bash" />
          </Card>
          <Card className="p-4">
            <h3 className="text-sm font-semibold mb-2 docs-h3">Pin to a specific version</h3>
            <CodeBlock code={`gitmap self-install --version v3.99.0 --yes`} title="Any platform" language="bash" />
            <p className="text-xs text-muted-foreground mt-2">
              <code className="docs-inline-code">--yes</code> skips the install-folder prompt.{" "}
              <code className="docs-inline-code">--shell-mode &lt;mode&gt;</code> controls which profiles get the
              PATH snippet (auto / zsh / bash / pwsh / fish / both / combos like{" "}
              <code className="docs-inline-code">zsh+pwsh</code>).
            </p>
          </Card>
        </section>

        {/* Verify */}
        <section className="mb-12">
          <SectionHeader
            icon={Stethoscope}
            title="Verify your install"
            kicker="Run these in order if anything looks off after install or update."
          />
          <CodeBlock
            code={`gitmap version                   # confirms binary on PATH + build info
gitmap doctor                    # PATH, profile snippets, deploy folder, DB
gitmap setup print-path-snippet  # prints the exact bytes the installer wrote`}
            title="Verification"
            language="bash"
          />
        </section>

        {/* Uninstall */}
        <section className="mb-12">
          <SectionHeader
            icon={Trash2}
            title="Uninstall"
            kicker="First tries `gitmap self-uninstall`; falls back to a manual sweep if gitmap is no longer on PATH."
          />
          <PlatformBlock commands={uninstallCmds} />
          <div className="mt-4 overflow-x-auto">
            <table className="w-full text-sm border border-border rounded-lg overflow-hidden">
              <thead>
                <tr className="bg-muted/50">
                  <th className="text-left px-4 py-2 font-mono text-xs text-muted-foreground">Flag</th>
                  <th className="text-left px-4 py-2 font-mono text-xs text-muted-foreground">Effect</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-border">
                {[
                  ["--yes / -Yes", "Skip the “delete user data?” prompt and assume yes"],
                  ["--keep-data", "Always keep %APPDATA%\\gitmap (Windows) / ~/.config/gitmap (Unix)"],
                  ["--dir <path>", "Override the auto-detected deploy root"],
                ].map(([flag, desc]) => (
                  <tr key={flag} className="hover:bg-muted/30 transition-colors">
                    <td className="px-4 py-2 font-mono text-xs text-primary">{flag}</td>
                    <td className="px-4 py-2 text-xs text-muted-foreground">{desc}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </section>

        {/* PATH activation */}
        <section className="mb-12">
          <SectionHeader
            icon={Terminal}
            title="Shell PATH activation"
            kicker="Profiles the installer writes for each --shell-mode value."
          />
          <div className="overflow-x-auto">
            <table className="w-full text-sm border border-border rounded-lg overflow-hidden">
              <thead>
                <tr className="bg-muted/50">
                  <th className="text-left px-4 py-2 font-mono text-xs text-muted-foreground">Mode</th>
                  <th className="text-left px-4 py-2 font-mono text-xs text-muted-foreground">Profiles touched</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-border">
                {[
                  ["auto (default)", "Current shell + ~/.profile on Unix; PowerShell profile on Windows"],
                  ["both", "zsh + bash + ~/.profile + fish (if present) + pwsh"],
                  ["zsh / bash / pwsh / fish", "Only that shell family"],
                  ["zsh+pwsh, bash+fish, …", "Strict union of listed families (no auto-detect, no ~/.profile)"],
                ].map(([mode, desc]) => (
                  <tr key={mode} className="hover:bg-muted/30 transition-colors">
                    <td className="px-4 py-2 font-mono text-xs text-primary whitespace-nowrap">{mode}</td>
                    <td className="px-4 py-2 text-xs text-muted-foreground">{desc}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
          <p className="text-xs text-muted-foreground mt-3">
            Snippet templates live in{" "}
            <code className="docs-inline-code">gitmap-v27/constants/constants_pathsnippet.go</code> so install.sh,
            install.ps1, and{" "}
            <code className="docs-inline-code">gitmap setup print-path-snippet</code> emit byte-identical output.
          </p>
        </section>
      </div>
    </DocsLayout>
  );
};

export default InstallGitmapPage;
