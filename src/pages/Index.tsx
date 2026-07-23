
import { Link } from "react-router-dom";
import { FolderGit2, GitBranch, RefreshCw, Eye } from "lucide-react";
import DocsLayout from "@/components/docs/DocsLayout";
import FeatureCard from "@/components/docs/FeatureCard";
import InstallBlock from "@/components/docs/InstallBlock";
import CommandBubbles from "@/components/docs/CommandBubbles";
import TabOrderMap from "@/components/docs/TabOrderMap";
import { VERSION } from "@/constants/index";

const INSTALL_TABS = [
  {
    label: "Windows",
    command:
      "irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/install-quick.ps1 | iex",
  },
  {
    label: "Linux / macOS",
    command:
      "curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/install-quick.sh | bash",
  },
];

const UNINSTALL_TABS = [
  {
    label: "Windows",
    command:
      "irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/uninstall-quick.ps1 | iex",
  },
  {
    label: "Linux / macOS",
    command:
      "curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/uninstall-quick.sh | bash",
  },
];

const HomePage = () => {
  return (
    <DocsLayout>
      <section className="py-14 text-center">
        <div className="reveal">
          <div className="mb-4 text-center">
            <h1 className="text-4xl md:text-6xl font-heading font-bold docs-h1 text-shimmer tracking-tight">
              gitmap
            </h1>
            <p className="mx-auto mt-2 w-fit bg-transparent p-0 font-mono text-sm font-bold leading-tight text-foreground md:text-base">
              {VERSION}
            </p>
          </div>
          <p className="text-lg text-muted-foreground max-w-2xl mx-auto mb-8 leading-relaxed font-sans">
            Scan a folder tree for Git repos, generate structured clone files, and
            re-clone the exact layout on any machine. Track, group, release, and
            manage repositories from a single CLI.
          </p>

          <div className="mx-auto mb-8 max-w-3xl rounded-xl bg-card/40 px-8 py-7 text-center backdrop-blur-sm">
            <div className="mb-6 flex items-center justify-center gap-2 pb-2">
              <span className="h-2.5 w-2.5 rounded-full bg-destructive/80" />
              <span className="h-2.5 w-2.5 rounded-full bg-primary/80" />
              <span className="h-2.5 w-2.5 rounded-full bg-muted-foreground/50" />
              <p className="ml-2 text-xs font-sans uppercase tracking-[0.18em] text-muted-foreground">
                Terminal quick actions
              </p>
            </div>

            <div className="space-y-6">
              <div>
                <p className="mb-3 text-center text-xs font-heading font-semibold uppercase tracking-[0.18em] text-primary">
                  Install
                </p>
                <InstallBlock tabs={INSTALL_TABS} />
              </div>

              <div>
                <p className="mb-3 text-center text-xs font-heading font-semibold uppercase tracking-[0.18em] text-muted-foreground">
                  Uninstall
                </p>
                <InstallBlock tabs={UNINSTALL_TABS} />
              </div>
            </div>

            <p className="mx-auto mt-6 max-w-2xl text-xs text-muted-foreground font-sans leading-relaxed">
              Uninstall removes the <code className="font-mono text-foreground">gitmap</code> binary and its PATH entries, then prompts before deleting your data folder
              (<code className="font-mono text-foreground">%APPDATA%\gitmap-v28</code> on Windows, <code className="font-mono text-foreground">~/.config/gitmap-v28</code> on Linux/macOS).
              Pass <code className="font-mono text-foreground">--keep-data</code> to always keep it, or <code className="font-mono text-foreground">-y</code>/<code className="font-mono text-foreground">--yes</code> to skip the prompt.
            </p>
          </div>

          <div className="flex gap-4 justify-center">
            <Link
              to="/getting-started"
              className="btn-slide group relative rounded-sm border border-primary bg-primary px-6 py-2.5 font-heading text-sm font-medium text-primary-foreground shadow-sm hover:brightness-110 active:translate-y-px focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 focus-visible:ring-offset-background"
            >
              Get Started
            </Link>
            <Link
              to="/commands"
              className="btn-slide btn-slide-ghost group relative rounded-sm border border-border bg-card px-6 py-2.5 font-heading text-sm font-medium text-foreground hover:border-primary/40 hover:bg-secondary active:translate-y-px focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 focus-visible:ring-offset-background"
            >
              View Commands
            </Link>
          </div>
        </div>
      </section>

      <hr className="docs-hr" />

      <section className="reveal grid md:grid-cols-2 gap-4 py-8">
        <FeatureCard
          icon={FolderGit2}
          title="Scan & Map"
          description="Recursively discover Git repos, extract metadata, and output CSV/JSON/terminal views with clone scripts."
        />
        <FeatureCard
          icon={GitBranch}
          title="Clone & Restore"
          description="Re-clone the exact folder structure on a new machine from JSON, CSV, or text files with safe-pull and progress tracking."
        />
        <FeatureCard
          icon={RefreshCw}
          title="Release & Version"
          description="Create releases with tags, branches, changelogs, and semantic versioning — all from the command line."
        />
        <FeatureCard
          icon={Eye}
          title="Watch & Monitor"
          description="Live-refresh dashboard showing dirty/clean status, ahead/behind counts, and stash entries across all tracked repos."
        />
      </section>

      <CommandBubbles />

      <TabOrderMap />

      <script
        type="application/ld+json"
        dangerouslySetInnerHTML={{
          __html: JSON.stringify({
            "@context": "https://schema.org",
            "@type": "SoftwareApplication",
            name: "gitmap-v28",
            applicationCategory: "DeveloperApplication",
            operatingSystem: "Windows, macOS, Linux",
            description: "CLI tool to scan, map, and re-clone Git repository trees.",
          }),
        }}
      />
    </DocsLayout>
  );
};

export default HomePage;
