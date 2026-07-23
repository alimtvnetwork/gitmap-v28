import { useMemo, useState } from "react";
import { useParams, Link } from "react-router-dom";
import DocsLayout from "@/components/docs/DocsLayout";
import CodeBlock from "@/components/docs/CodeBlock";
import { Pin, Wrench, ShieldCheck, Tag, ExternalLink, AlertTriangle } from "lucide-react";
import { Tabs, TabsList, TabsTrigger, TabsContent } from "@/components/ui/tabs";
import {
  RELEASE_REPO,
  buildReleaseSnippets,
  isValidReleaseVersion,
  normalizeReleaseVersion,
  type ReleasePlatform,
} from "./releaseVersionSnippets";

const REPO = RELEASE_REPO;

const ReleaseVersionPage = () => {
  const { version: rawVersion = "" } = useParams<{ version: string }>();
  const [platform, setPlatform] = useState<ReleasePlatform>("windows");

  const version = normalizeReleaseVersion(rawVersion);
  const isValid = isValidReleaseVersion(rawVersion);

  const snippets = useMemo(
    () => (isValid ? buildReleaseSnippets(version, platform) : null),
    [version, platform, isValid],
  );

  if (!isValid) {
    return (
      <DocsLayout>
        <div className="space-y-6 max-w-2xl">
          <div className="flex items-center gap-3 text-destructive">
            <AlertTriangle className="h-6 w-6" />
            <h1 className="text-2xl font-heading font-bold">Invalid version</h1>
          </div>
          <p className="text-muted-foreground">
            <code className="font-mono text-foreground">{rawVersion}</code> is not a
            valid semantic version tag. Expected format:{" "}
            <code className="font-mono text-primary">vMAJOR.MINOR.PATCH</code>.
          </p>
          <Link
            to="/changelog"
            className="inline-flex items-center gap-2 text-primary hover:underline"
          >
            Browse all releases
            <ExternalLink className="h-4 w-4" />
          </Link>
        </div>
      </DocsLayout>
    );
  }

  const releaseUrl = `https://github.com/${REPO}/releases/tag/${version}`;

  return (
    <DocsLayout>
      <div className="space-y-10 max-w-3xl">
        {/* Header */}
        <header className="space-y-3">
          <div className="inline-flex items-center gap-2 rounded-full border border-primary/30 bg-primary/10 px-3 py-1 text-xs font-mono text-foreground dark:bg-primary/15 dark:text-primary">
            <Tag className="h-3.5 w-3.5" />
            release {version}
          </div>
          <h1 className="text-3xl font-heading font-bold text-foreground docs-h1">
            Install gitmap {version}
          </h1>
          <p className="text-muted-foreground leading-relaxed">
            This page installs <strong className="text-foreground">exactly</strong>{" "}
            this version. The pinned installer never resolves <em>latest</em>, never
            auto-upgrades, and verifies the SHA256 checksum before extracting.
          </p>
          <a
            href={releaseUrl}
            target="_blank"
            rel="noreferrer"
            className="inline-flex items-center gap-2 text-sm text-primary hover:underline"
          >
            View release {version} on GitHub
            <ExternalLink className="h-3.5 w-3.5" />
          </a>
        </header>

        {/* Platform tabs */}
        <Tabs value={platform} onValueChange={(v) => setPlatform(v as ReleasePlatform)}>
          <TabsList className="grid w-full max-w-md grid-cols-2">
            <TabsTrigger value="windows">Windows (PowerShell)</TabsTrigger>
            <TabsTrigger value="unix">macOS / Linux</TabsTrigger>
          </TabsList>

          {(["windows", "unix"] as ReleasePlatform[]).map((p) => (
            <TabsContent key={p} value={p} className="space-y-6 pt-6">
              {/* Pinned card */}
              <section className="rounded-lg border border-primary/40 bg-primary/5 p-5 space-y-3">
                <div className="flex items-start justify-between gap-3">
                  <div className="flex items-center gap-2">
                    <Pin className="h-5 w-5 text-primary" />
                    <h2 className="font-heading font-semibold text-foreground">
                      Pinned install
                    </h2>
                    <span className="rounded-full bg-primary px-2 py-0.5 text-[10px] font-semibold uppercase tracking-wide text-primary-foreground">
                      Recommended
                    </span>
                  </div>
                  <ShieldCheck className="h-5 w-5 text-primary/70" />
                </div>
                <p className="text-sm text-muted-foreground">
                  Downloads the per-version snapshot script attached to release{" "}
                  <code className="font-mono text-foreground">{version}</code>. Drift-
                  proof: the version is baked into the script.
                </p>
                {snippets && (
                  <CodeBlock
                    code={snippets.pinned}
                    title={
                      p === "windows"
                        ? `release-version-${version}.ps1`
                        : `release-version-${version}.sh`
                    }
                  />
                )}
              </section>

              {/* Generic card */}
              <section className="rounded-lg border border-border bg-card p-5 space-y-3">
                <div className="flex items-center gap-2">
                  <Wrench className="h-5 w-5 text-muted-foreground" />
                  <h2 className="font-heading font-semibold text-foreground">
                    Generic install
                  </h2>
                  <span className="rounded-full bg-muted px-2 py-0.5 text-[10px] font-semibold uppercase tracking-wide text-muted-foreground">
                    Advanced
                  </span>
                </div>
                <p className="text-sm text-muted-foreground">
                  Uses the always-current <code className="font-mono">release-version</code>{" "}
                  script with an explicit{" "}
                  <code className="font-mono text-primary">
                    {p === "windows" ? "-Version" : "--version"}
                  </code>{" "}
                  argument. Useful when scripting installs for multiple versions.
                </p>
                {snippets && <CodeBlock code={snippets.generic} />}
              </section>
            </TabsContent>
          ))}
        </Tabs>

        {/* Guarantees */}
        <section className="space-y-3">
          <h2 className="text-xl font-heading font-bold text-foreground docs-h2">
            What the script guarantees
          </h2>
          <ul className="space-y-2 text-sm text-muted-foreground">
            <li>
              <strong className="text-foreground">Version-pinned.</strong> Hard-fails
              if the GitHub release for{" "}
              <code className="font-mono text-primary">{version}</code> doesn&apos;t
              exist — never silently falls back to <em>latest</em>.
            </li>
            <li>
              <strong className="text-foreground">OS / arch matched.</strong> Detects{" "}
              <code className="font-mono">windows|linux|darwin</code> +{" "}
              <code className="font-mono">amd64|arm64</code> and picks the matching
              release asset.
            </li>
            <li>
              <strong className="text-foreground">Checksum verified.</strong> SHA256
              hash is checked against{" "}
              <code className="font-mono">checksums.txt</code> before extraction.
            </li>
            <li>
              <strong className="text-foreground">Self-install chained.</strong> The
              freshly-installed binary runs{" "}
              <code className="font-mono text-primary">gitmap self-install</code> to
              register completions and profiles.
            </li>
            <li>
              <strong className="text-foreground">Interactive on miss.</strong> If
              the version is unavailable, prompts you to pick from the 5 most-recent
              releases — or exits 1 in non-interactive shells.
            </li>
          </ul>
        </section>

        {/* See also */}
        <section className="space-y-2">
          <h2 className="text-xl font-heading font-bold text-foreground docs-h2">
            See also
          </h2>
          <ul className="space-y-1 text-sm">
            <li>
              <Link to="/install" className="text-primary hover:underline">
                /install
              </Link>{" "}
              — generic latest installer (front-page Get Started)
            </li>
            <li>
              <Link to="/changelog" className="text-primary hover:underline">
                /changelog
              </Link>{" "}
              — every released version
            </li>
            <li>
              <Link to="/release" className="text-primary hover:underline">
                /release
              </Link>{" "}
              — how releases are produced
            </li>
          </ul>
        </section>
      </div>
    </DocsLayout>
  );
};

export default ReleaseVersionPage;
