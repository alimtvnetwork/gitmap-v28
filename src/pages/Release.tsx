import DocsLayout from "@/components/docs/DocsLayout";
import CodeBlock from "@/components/docs/CodeBlock";
import { GitBranch, Tag, Upload, Clock, Shield, Eye, Package, Target, FileCheck, Archive, FileText } from "lucide-react";

const features = [
  { icon: GitBranch, title: "Branch + Tag", desc: "Creates release/vX.Y.Z branch and vX.Y.Z tag in one step." },
  { icon: Tag, title: "Semver Padding", desc: "Partial versions auto-pad: v1 → v1.0.0, v1.2 → v1.2.0." },
  { icon: Upload, title: "Auto Push", desc: "Pushes branch and tag to origin after creation." },
  { icon: Clock, title: "Auto Increment", desc: "Use --bump to increment from the latest released version." },
  { icon: Shield, title: "Duplicate Detection", desc: "Aborts if the version tag or metadata file already exists. Recovers orphaned metadata automatically." },
  { icon: Eye, title: "Dry Run", desc: "Preview all steps without executing with --dry-run." },
  { icon: Package, title: "Go Cross-Compile", desc: "Auto-detect go.mod and build binaries for 6 OS/arch targets." },
  { icon: Archive, title: "Compress & Checksum", desc: "Wrap assets in .zip/.tar.gz and generate SHA256 checksums." },
  { icon: Target, title: "Custom Targets", desc: "Override default matrix via --targets flag or config.json." },
  { icon: FileText, title: "Release Notes", desc: "Use --notes / -N to set a title and annotation for the release." },
];

const releaseFlags = [
  { flag: "--assets <path>", description: "Directory or file to record as release assets" },
  { flag: "--commit <sha>", description: "Create release from a specific commit" },
  { flag: "--branch <name>", description: "Create release from latest commit of a branch" },
  { flag: "--bump major|minor|patch", description: "Auto-increment from the latest released version" },
  { flag: "--notes <text> / -N", description: "Release notes or title for the release (used as tag annotation and GitHub title)" },
  { flag: "--draft", description: "Mark release metadata as draft" },
  { flag: "--dry-run", description: "Preview release steps without executing" },
  { flag: "--compress", description: "Wrap assets in .zip (Windows) or .tar.gz (Linux/macOS)" },
  { flag: "--checksums", description: "Generate SHA256 checksums.txt for assets" },
  { flag: "--no-assets", description: "Skip Go binary cross-compilation" },
  { flag: "--targets <list>", description: "Cross-compile targets (e.g. windows/amd64,linux/arm64)" },
  { flag: "--list-targets", description: "Print resolved target matrix and exit" },
  { flag: "--zip-group <name>", description: "Include a persistent zip group as a release asset (repeatable)" },
  { flag: "-Z <path>", description: "Add ad-hoc file or folder to zip as a release asset (repeatable)" },
  { flag: "--bundle <name.zip>", description: "Bundle all -Z items into a single named archive" },
  { flag: "--verbose", description: "Write detailed debug log" },
];

const branchFlags = [
  { flag: "--assets <path>", description: "Directory or file to record" },
  { flag: "--notes <text> / -N", description: "Release notes or title for the release" },
  { flag: "--draft", description: "Mark release metadata as draft" },
  { flag: "--dry-run", description: "Preview steps without executing" },
  { flag: "--no-commit", description: "Skip post-release auto-commit and push" },
  { flag: "--verbose", description: "Write detailed debug log" },
];

const bumpExamples = [
  { current: "1.2.3", patch: "1.2.4", minor: "1.3.0", major: "2.0.0" },
  { current: "0.9.1", patch: "0.9.2", minor: "0.10.0", major: "1.0.0" },
];

const paddingExamples = [
  { input: "v1", resolved: "v1.0.0", branch: "release/v1.0.0", tag: "v1.0.0" },
  { input: "v1.2", resolved: "v1.2.0", branch: "release/v1.2.0", tag: "v1.2.0" },
  { input: "v1.2.3", resolved: "v1.2.3", branch: "release/v1.2.3", tag: "v1.2.3" },
];

const defaultTargets = [
  { goos: "windows", goarch: "amd64", suffix: "_windows_amd64.exe" },
  { goos: "windows", goarch: "arm64", suffix: "_windows_arm64.exe" },
  { goos: "linux", goarch: "amd64", suffix: "_linux_amd64" },
  { goos: "linux", goarch: "arm64", suffix: "_linux_arm64" },
  { goos: "darwin", goarch: "amd64", suffix: "_darwin_amd64" },
  { goos: "darwin", goarch: "arm64", suffix: "_darwin_arm64" },
];

const errorScenarios = [
  { scenario: "Invalid version string", behavior: "'abc' is not a valid version." },
  { scenario: "--commit SHA not found", behavior: "commit abc123 not found." },
  { scenario: "--branch does not exist", behavior: "branch develop does not exist." },
  { scenario: "Push to remote fails", behavior: "failed to push to remote: <detail>" },
  { scenario: "Version already released", behavior: "Version v1.2.3 is already released." },
  { scenario: "Orphaned metadata (no tag/branch)", behavior: "Prompts to remove stale JSON, then proceeds if confirmed." },
  { scenario: "--bump + version argument", behavior: "--bump cannot be used with an explicit version argument." },
  { scenario: "--commit + --branch", behavior: "--commit and --branch are mutually exclusive." },
  { scenario: "Go build fails for target", behavior: "Logs error, continues with remaining targets." },
  { scenario: "Asset upload fails", behavior: "Retries once, then logs and continues." },
];

const ReleasePage = () => {
  return (
    <DocsLayout>
      <div className="space-y-10">
        {/* Header */}
        <div>
          <h1 className="text-3xl font-heading font-bold text-foreground mb-3 docs-h1">Release Command</h1>
          <p className="text-muted-foreground leading-relaxed max-w-2xl">
            Automate Git release workflows: create branches, tags, push to remote, cross-compile Go binaries,
            upload to GitHub, and track release history. Supports semver, compression, checksums,
            config-driven targets, and dry-run preview.
          </p>
        </div>

        {/* Features */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {features.map((f) => (
            <div key={f.title} className="p-5 rounded-lg border border-border bg-card hover:border-primary/40 transition-colors group">
              <div className="h-9 w-9 rounded-md bg-primary/10 flex items-center justify-center mb-3 group-hover:bg-primary/20 transition-colors">
                <f.icon className="h-4 w-4 text-primary" />
              </div>
              <h3 className="font-mono font-semibold text-foreground text-sm mb-1">{f.title}</h3>
              <p className="text-xs text-muted-foreground leading-relaxed">{f.desc}</p>
            </div>
          ))}
        </div>

        {/* Commands */}
        <div>
          <h2 className="text-xl font-heading font-bold text-foreground mb-4 docs-h2">Commands</h2>
          <div className="space-y-4">
            <div className="p-4 rounded-lg border border-border bg-card">
              <h3 className="font-mono font-semibold text-foreground mb-1">gitmap release [version] <span className="text-muted-foreground font-normal text-sm">(alias: r)</span></h3>
              <p className="text-sm text-muted-foreground">Create a release branch, Git tag, push to remote, and optionally cross-compile + upload Go binaries.</p>
            </div>
            <div className="p-4 rounded-lg border border-border bg-card">
              <h3 className="font-mono font-semibold text-foreground mb-1">gitmap release-branch &lt;branch&gt; <span className="text-muted-foreground font-normal text-sm">(alias: rb)</span></h3>
              <p className="text-sm text-muted-foreground">Complete a release from an existing release/vX.Y.Z branch.</p>
            </div>
            <div className="p-4 rounded-lg border border-border bg-card">
              <h3 className="font-mono font-semibold text-foreground mb-1">gitmap release-pending <span className="text-muted-foreground font-normal text-sm">(alias: rp)</span></h3>
              <p className="text-sm text-muted-foreground">Release all release/v* branches that are missing tags.</p>
            </div>
          </div>
        </div>

        {/* Workflow Diagram */}
        <div>
          <h2 className="text-xl font-heading font-bold text-foreground mb-4 docs-h2">Release Workflow</h2>

          <h3 className="text-base font-heading font-semibold text-foreground mb-3">gitmap release [version]</h3>
          <div className="p-5 rounded-lg border border-border bg-card font-mono text-sm space-y-1 mb-6">
            {[
              "1. Resolve version (CLI → --bump → version.json → error)",
              "2. Pad partial version to full semver",
              "3. Check .gitmap/release/ and git tags for duplicates",
              "4. Resolve source commit (--commit / --branch / HEAD)",
              "5. Create/switch to branch release/vX.Y.Z",
              "6. Create git tag vX.Y.Z (annotated with --notes if provided)",
              "7. Push branch + tag to origin",
              "7a. Cross-compile Go binaries (if go.mod detected)",
              "7a′. Resolve and compress zip groups / ad-hoc -Z items",
              "7b. Compress assets (.zip/.tar.gz) if --compress",
              "7c. Generate checksums.txt if --checksums",
              "7d. Upload assets to GitHub Releases API",
              "8. Return to original branch",
              "9. Write .gitmap/release/vX.Y.Z.json + update latest.json on original branch",
              "10. Auto-commit .gitmap/release/ metadata files",
            ].map((step) => (
              <p key={step} className="text-foreground/80 pl-2">{step}</p>
            ))}
          </div>

          <h3 className="text-base font-heading font-semibold text-foreground mb-3">gitmap release-branch / release-pending</h3>
          <div className="p-5 rounded-lg border border-border bg-card font-mono text-sm space-y-1">
            {[
              "1. Validate branch exists / discover pending branches + metadata",
              "2. Extract version from branch name, pad to semver",
              "3. Check if tag already exists → abort if so",
              "4. Checkout the release branch",
              "5. Create tag (annotated with --notes if provided)",
              "6. Push branch + tag to origin, upload assets",
              "7. Return to original branch",
            ].map((step) => (
              <p key={step} className="text-foreground/80 pl-2">{step}</p>
            ))}
            <div className="mt-3 pt-3 border-t border-border/50">
              <p className="text-muted-foreground text-xs">
                <span className="text-primary font-semibold">Note:</span> Steps 9–10 from the release command (metadata write + auto-commit) are skipped.
                These commands process existing branches/metadata — they only tag, push, and upload.
              </p>
            </div>
          </div>
        </div>

        {/* Go Cross-Compilation */}
        <div>
          <h2 className="text-xl font-heading font-bold text-foreground mb-4 docs-h2">Go Cross-Compilation</h2>
          <p className="text-sm text-muted-foreground mb-4">
            When a <code className="font-mono text-primary">go.mod</code> file is detected, gitmap automatically
            cross-compiles binaries for all OS/arch targets using <code className="font-mono text-primary">CGO_ENABLED=0</code>.
            No external tools required — uses Go's native cross-compilation.
          </p>

          <h3 className="text-base font-heading font-semibold text-foreground mb-3">Default Target Matrix</h3>
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-border">
                  <th className="text-left py-2 px-3 font-mono text-muted-foreground font-medium">GOOS</th>
                  <th className="text-left py-2 px-3 font-mono text-muted-foreground font-medium">GOARCH</th>
                  <th className="text-left py-2 px-3 font-mono text-muted-foreground font-medium">Filename Suffix</th>
                </tr>
              </thead>
              <tbody>
                {defaultTargets.map((t) => (
                  <tr key={`${t.goos}-${t.goarch}`} className="border-b border-border/50">
                    <td className="py-2 px-3 font-mono text-primary">{t.goos}</td>
                    <td className="py-2 px-3 font-mono text-foreground">{t.goarch}</td>
                    <td className="py-2 px-3 font-mono text-muted-foreground">{t.suffix}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          <h3 className="text-base font-heading font-semibold text-foreground mt-6 mb-3">Target Resolution (Three-Layer)</h3>
          <div className="bg-card border border-border rounded-lg p-4 space-y-2">
            <div className="flex items-center gap-3 font-mono text-sm">
              <span className="bg-primary/20 text-foreground border border-primary/30 dark:bg-primary/20 dark:text-primary px-3 py-1 rounded font-semibold">1. --targets flag</span>
              <span className="text-muted-foreground">→ Highest priority (always wins)</span>
            </div>
            <div className="flex items-center gap-3 font-mono text-sm">
              <span className="bg-primary/10 text-foreground border border-primary/20 dark:bg-primary/15 dark:text-primary px-3 py-1 rounded transition-colors duration-300 hover:border-primary/40 hover:shadow-sm hover:shadow-primary/10">2. config.json</span>
              <span className="text-muted-foreground">→ release.targets array</span>
            </div>
            <div className="flex items-center gap-3 font-mono text-sm">
              <span className="bg-muted text-muted-foreground px-3 py-1 rounded">3. Built-in defaults</span>
              <span className="text-muted-foreground">→ All 6 targets</span>
            </div>
          </div>

          <p className="text-sm text-muted-foreground mt-4">
            Use <code className="font-mono text-primary">gitmap release --list-targets</code> to inspect
            the resolved matrix:
          </p>
          <CodeBlock code={`$ gitmap release --list-targets\nRelease targets (6):\nSource: built-in defaults\n\n  windows/amd64\n  windows/arm64\n  linux/amd64\n  linux/arm64\n  darwin/amd64\n  darwin/arm64`} />

          <p className="text-sm text-muted-foreground mt-4">
            With a <code className="font-mono text-primary">--targets</code> override:
          </p>
          <CodeBlock code={`$ gitmap release --list-targets --targets windows/amd64,linux/amd64\nRelease targets (2):\nSource: --targets flag\n\n  windows/amd64\n  linux/amd64`} />
        </div>

        {/* Version Resolution */}
        <div>
          <h2 className="text-xl font-heading font-bold text-foreground mb-4 docs-h2">Version Resolution</h2>
          <p className="text-sm text-muted-foreground mb-3">Version is resolved in priority order:</p>
          <div className="space-y-2 mb-6">
            {[
              { label: "CLI argument", example: "gitmap release v1.2.3" },
              { label: "--bump flag", example: "reads latest, increments" },
              { label: "version.json", example: '{ "version": "1.2.3" }' },
              { label: "Error", example: "no version source found" },
            ].map((item, i) => (
              <div key={item.label} className="flex items-start gap-3 text-sm">
                <span className="font-mono text-primary font-bold">{i + 1}.</span>
                <span className="font-mono font-semibold text-foreground">{item.label}</span>
                <span className="text-muted-foreground">— {item.example}</span>
              </div>
            ))}
          </div>

          <h3 className="text-base font-heading font-semibold text-foreground mb-3">Partial Version Padding</h3>
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-border">
                  <th className="text-left py-2 px-3 font-mono text-muted-foreground font-medium">Input</th>
                  <th className="text-left py-2 px-3 font-mono text-muted-foreground font-medium">Resolved</th>
                  <th className="text-left py-2 px-3 font-mono text-muted-foreground font-medium">Branch</th>
                  <th className="text-left py-2 px-3 font-mono text-muted-foreground font-medium">Tag</th>
                </tr>
              </thead>
              <tbody>
                {paddingExamples.map((row) => (
                  <tr key={row.input} className="border-b border-border/50">
                    <td className="py-2 px-3 font-mono text-primary">{row.input}</td>
                    <td className="py-2 px-3 font-mono text-foreground">{row.resolved}</td>
                    <td className="py-2 px-3 font-mono text-foreground/80">{row.branch}</td>
                    <td className="py-2 px-3 font-mono text-foreground/80">{row.tag}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>

        {/* Release Flags */}
        <div>
          <h2 className="text-xl font-heading font-bold text-foreground mb-4 docs-h2">Release Flags</h2>
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-border">
                  <th className="text-left py-2 px-3 font-mono text-muted-foreground font-medium">Flag</th>
                  <th className="text-left py-2 px-3 font-mono text-muted-foreground font-medium">Description</th>
                </tr>
              </thead>
              <tbody>
                {releaseFlags.map((f) => (
                  <tr key={f.flag} className="border-b border-border/50">
                    <td className="py-2 px-3 font-mono text-primary whitespace-nowrap">{f.flag}</td>
                    <td className="py-2 px-3 text-muted-foreground">{f.description}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          <h3 className="text-base font-heading font-semibold text-foreground mt-6 mb-3">Release-Branch / Release-Pending Flags</h3>
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-border">
                  <th className="text-left py-2 px-3 font-mono text-muted-foreground font-medium">Flag</th>
                  <th className="text-left py-2 px-3 font-mono text-muted-foreground font-medium">Description</th>
                </tr>
              </thead>
              <tbody>
                {branchFlags.map((f) => (
                  <tr key={f.flag} className="border-b border-border/50">
                    <td className="py-2 px-3 font-mono text-primary whitespace-nowrap">{f.flag}</td>
                    <td className="py-2 px-3 text-muted-foreground">{f.description}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>

        {/* Auto-Increment */}
        <div>
          <h2 className="text-xl font-heading font-bold text-foreground mb-4 docs-h2">Auto-Increment (--bump)</h2>
          <p className="text-sm text-muted-foreground mb-3">
            Reads the latest version from <span className="font-mono text-foreground">.gitmap/release/latest.json</span> and increments.
            Falls back to scanning local Git tags.
          </p>
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-border">
                  <th className="text-left py-2 px-3 font-mono text-muted-foreground font-medium">Current</th>
                  <th className="text-left py-2 px-3 font-mono text-muted-foreground font-medium">--bump patch</th>
                  <th className="text-left py-2 px-3 font-mono text-muted-foreground font-medium">--bump minor</th>
                  <th className="text-left py-2 px-3 font-mono text-muted-foreground font-medium">--bump major</th>
                </tr>
              </thead>
              <tbody>
                {bumpExamples.map((row) => (
                  <tr key={row.current} className="border-b border-border/50">
                    <td className="py-2 px-3 font-mono text-foreground">{row.current}</td>
                    <td className="py-2 px-3 font-mono text-primary">{row.patch}</td>
                    <td className="py-2 px-3 font-mono text-primary">{row.minor}</td>
                    <td className="py-2 px-3 font-mono text-primary">{row.major}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>

        {/* Usage Examples */}
        <div>
          <h2 className="text-xl font-heading font-bold text-foreground mb-4 docs-h2">Usage Examples</h2>
          <div className="space-y-4">
            <div>
              <p className="text-sm text-muted-foreground mb-2">Full semver release from HEAD</p>
              <CodeBlock code="gitmap release v1.2.3" />
            </div>
            <div>
              <p className="text-sm text-muted-foreground mb-2">Auto-increment with compression and checksums</p>
              <CodeBlock code="gitmap release --bump patch --compress --checksums" />
            </div>
            <div>
              <p className="text-sm text-muted-foreground mb-2">Custom cross-compile targets</p>
              <CodeBlock code="gitmap release v2.0.0 --targets windows/amd64,linux/amd64" />
            </div>
            <div>
              <p className="text-sm text-muted-foreground mb-2">Inspect resolved target matrix</p>
              <CodeBlock code="gitmap release --list-targets" />
            </div>
            <div>
              <p className="text-sm text-muted-foreground mb-2">Skip Go binary compilation</p>
              <CodeBlock code="gitmap release v1.0.0 --no-assets" />
            </div>
            <div>
              <p className="text-sm text-muted-foreground mb-2">Release with a persistent zip group</p>
              <CodeBlock code="gitmap release v3.0.0 --zip-group docs-bundle" />
            </div>
            <div>
              <p className="text-sm text-muted-foreground mb-2">Ad-hoc zip with bundle</p>
              <CodeBlock code={`gitmap release v3.0.0 -Z ./dist/report.pdf -Z ./dist/manual.pdf --bundle docs.zip`} />
            </div>
            <div>
              <p className="text-sm text-muted-foreground mb-2">Combine persistent groups and ad-hoc items</p>
              <CodeBlock code="gitmap release v3.0.0 --zip-group docs-bundle -Z ./extras/notes.txt" />
            </div>
            <div>
              <p className="text-sm text-muted-foreground mb-2">Release with notes (title and tag annotation)</p>
              <CodeBlock code={`gitmap release --bump patch -N 'Hotfix for auth timeout'`} />
            </div>
            <div>
              <p className="text-sm text-muted-foreground mb-2">Draft and dry-run</p>
              <CodeBlock code={`gitmap release v3.0.0-rc.1 --draft\ngitmap release v1.0.0 --dry-run`} />
            </div>
            <div>
              <p className="text-sm text-muted-foreground mb-2">Release from existing branch / pending</p>
              <CodeBlock code={`gitmap release-branch release/v1.2.0\ngitmap release-pending\ngitmap release-pending --dry-run`} />
            </div>
          </div>
        </div>

        {/* Dry-Run Output */}
        <div>
          <h2 className="text-xl font-heading font-bold text-foreground mb-4 docs-h2">Dry-Run Preview</h2>
          <div className="bg-card border border-border rounded-lg p-5 font-mono text-sm space-y-1">
            <p className="text-muted-foreground">{`$ gitmap release v1.2.3 --dry-run`}</p>
            <p className="text-foreground/80">&nbsp;&nbsp;[dry-run] Create branch release/v1.2.3 from main</p>
            <p className="text-foreground/80">&nbsp;&nbsp;[dry-run] Create tag v1.2.3</p>
            <p className="text-foreground/80">&nbsp;&nbsp;[dry-run] Push branch and tag to origin</p>
            <p className="text-foreground/80">&nbsp;&nbsp;[dry-run] Would cross-compile 6 binaries:</p>
            <p className="text-foreground/60">&nbsp;&nbsp;&nbsp;&nbsp;→ gitmap_v1.2.3_windows_amd64.exe</p>
            <p className="text-foreground/60">&nbsp;&nbsp;&nbsp;&nbsp;→ gitmap_v1.2.3_linux_amd64</p>
            <p className="text-foreground/60">&nbsp;&nbsp;&nbsp;&nbsp;→ gitmap_v1.2.3_darwin_arm64</p>
            <p className="text-foreground/60">&nbsp;&nbsp;&nbsp;&nbsp;...</p>
            <p className="text-foreground/80">&nbsp;&nbsp;[dry-run] Would upload 6 assets</p>
            <p className="text-foreground/80">&nbsp;&nbsp;[dry-run] Write metadata to .gitmap/release/v1.2.3.json</p>
            <p className="text-foreground/80">&nbsp;&nbsp;[dry-run] Mark v1.2.3 as latest</p>
          </div>
        </div>

        {/* Orphaned Metadata Recovery */}
        <div>
          <h2 className="text-xl font-heading font-bold text-foreground mb-4 docs-h2">Orphaned Metadata Recovery</h2>
          <p className="text-sm text-muted-foreground mb-4 max-w-2xl">
            If a <code className="font-mono text-primary">.gitmap/release/vX.Y.Z.json</code> file exists but neither the
            Git tag nor the release branch is found, the release command detects orphaned metadata and prompts
            for recovery instead of failing with a duplicate error.
          </p>

          <div className="bg-card border border-border rounded-lg p-5 font-mono text-sm space-y-1 mb-6">
            <p className="text-muted-foreground">$ gitmap release --bump patch</p>
            <p className="text-foreground/80">&nbsp;&nbsp;→ Bumped v2.3.9 → v2.3.10</p>
            <p className="text-yellow-500">&nbsp;&nbsp;⚠ Release metadata exists for v2.3.10 but no tag or branch was found.</p>
            <p className="text-foreground/80">&nbsp;&nbsp;→ Do you want to remove the release JSON and proceed? (y/N): <span className="text-primary">y</span></p>
            <p className="text-green-500">&nbsp;&nbsp;✓ Removed orphaned release metadata for v2.3.10</p>
            <p className="text-foreground/80">&nbsp;</p>
            <p className="text-foreground/80">&nbsp;&nbsp;Creating release v2.3.10...</p>
            <p className="text-green-500">&nbsp;&nbsp;✓ Created branch release/v2.3.10</p>
            <p className="text-green-500">&nbsp;&nbsp;✓ Created tag v2.3.10</p>
          </div>

          <h3 className="text-base font-heading font-semibold text-foreground mb-3">Detection Logic</h3>
          <div className="space-y-2 mb-4">
            {[
              { step: "1", text: "Release JSON exists for the target version" },
              { step: "2", text: "Git tag does not exist locally or on remote" },
              { step: "3", text: "Release branch does not exist" },
              { step: "→", text: "Prompt user to remove stale JSON and proceed" },
            ].map((item) => (
              <div key={item.step} className="flex items-start gap-3 text-sm">
                <span className="font-mono text-primary font-bold w-5 text-right">{item.step}</span>
                <span className="text-muted-foreground">{item.text}</span>
              </div>
            ))}
          </div>

          <div className="p-4 rounded-lg border border-border bg-card">
            <p className="text-sm text-muted-foreground">
              If the user declines (<code className="font-mono text-primary">n</code> or Enter), the release is aborted.
              If confirmed (<code className="font-mono text-primary">y</code>), the orphaned JSON file is deleted and
              the normal release workflow continues from step 5.
            </p>
          </div>
        </div>

        {/* Zip Groups */}
        <div>
          <h2 className="text-xl font-heading font-bold text-foreground mb-4 docs-h2">Zip Groups</h2>
          <p className="text-sm text-muted-foreground mb-4 max-w-2xl">
            Attach compressed file bundles to releases using persistent groups (stored in SQLite) or
            ad-hoc <code className="font-mono text-primary">-Z</code> items. Archives use maximum compression
            (Deflate level 9) and are placed in the staging directory before upload.
          </p>

          <h3 className="text-base font-heading font-semibold text-foreground mb-3">Flags</h3>
          <div className="overflow-x-auto mb-6">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-border">
                  <th className="text-left py-2 px-3 font-mono text-muted-foreground font-medium">Flag</th>
                  <th className="text-left py-2 px-3 font-mono text-muted-foreground font-medium">Behavior</th>
                </tr>
              </thead>
              <tbody>
                {[
                  { flag: "--zip-group <name>", desc: "Look up a named group from the database, compress its items into a single archive" },
                  { flag: "-Z <path>", desc: "Add an ad-hoc file or folder; each item becomes its own archive unless --bundle is used" },
                  { flag: "--bundle <name.zip>", desc: "Combine all -Z items into a single named archive instead of individual ones" },
                ].map((row) => (
                  <tr key={row.flag} className="border-b border-border/50">
                    <td className="py-2 px-3 font-mono text-primary whitespace-nowrap">{row.flag}</td>
                    <td className="py-2 px-3 text-muted-foreground">{row.desc}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          <h3 className="text-base font-heading font-semibold text-foreground mb-3">Archive Naming</h3>
          <div className="overflow-x-auto mb-6">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-border">
                  <th className="text-left py-2 px-3 font-mono text-muted-foreground font-medium">Scenario</th>
                  <th className="text-left py-2 px-3 font-mono text-muted-foreground font-medium">Output Filename</th>
                </tr>
              </thead>
              <tbody>
                {[
                  { scenario: "Persistent group, no custom name", output: "<group>_<version>.zip" },
                  { scenario: "Persistent group, custom archive name", output: "<archive-name> (as-is)" },
                  { scenario: "Ad-hoc file, no --bundle", output: "<filename>.zip" },
                  { scenario: "Ad-hoc folder, no --bundle", output: "<foldername>.zip" },
                  { scenario: "Ad-hoc with --bundle", output: "<bundle-name>.zip" },
                ].map((row) => (
                  <tr key={row.scenario} className="border-b border-border/50">
                    <td className="py-2 px-3 text-foreground">{row.scenario}</td>
                    <td className="py-2 px-3 font-mono text-primary">{row.output}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          <h3 className="text-base font-heading font-semibold text-foreground mb-3">Metadata (zipGroups)</h3>
          <p className="text-sm text-muted-foreground mb-3">
            Release metadata in <code className="font-mono text-primary">.gitmap/release/vX.Y.Z.json</code> includes a
            <code className="font-mono text-primary"> zipGroups</code> array recording each group and its archive:
          </p>
          <CodeBlock code={`{
  "version": "3.0.0",
  "tag": "v3.0.0",
  "zipGroups": [
    {
      "name": "docs-bundle",
      "archive": "docs-bundle_v3.0.0.zip",
      "items": ["README.md", "CHANGELOG.md", "docs/"]
    }
  ],
  "assets": [
    "gitmap_v3.0.0_windows_amd64.exe.zip",
    "docs-bundle_v3.0.0.zip"
  ]
}`} />

          <h3 className="text-base font-heading font-semibold text-foreground mt-6 mb-3">Dry-Run Output</h3>
          <div className="bg-card border border-border rounded-lg p-5 font-mono text-sm space-y-1">
            <p className="text-muted-foreground">$ gitmap release v3.0.0 --zip-group docs-bundle -Z ./report.pdf --dry-run</p>
            <p className="text-foreground/80">&nbsp;&nbsp;[dry-run] Would create 2 zip archive(s):</p>
            <p className="text-foreground/60">&nbsp;&nbsp;&nbsp;&nbsp;→ docs-bundle_v3.0.0.zip (3 items: README.md, CHANGELOG.md, docs/)</p>
            <p className="text-foreground/60">&nbsp;&nbsp;&nbsp;&nbsp;→ report.pdf.zip (1 item: dist/report.pdf)</p>
            <p className="text-foreground/80">&nbsp;&nbsp;[dry-run] Would upload 8 assets + checksums.txt</p>
          </div>
        </div>

        {/* Error Scenarios */}
        <div>
          <h2 className="text-xl font-heading font-bold text-foreground mb-4 docs-h2">Error Scenarios</h2>
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-border">
                  <th className="text-left py-2 px-3 font-mono text-muted-foreground font-medium">Scenario</th>
                  <th className="text-left py-2 px-3 font-mono text-muted-foreground font-medium">Behavior</th>
                </tr>
              </thead>
              <tbody>
                {errorScenarios.map((row) => (
                  <tr key={row.scenario} className="border-b border-border/50">
                    <td className="py-2 px-3 text-foreground">{row.scenario}</td>
                    <td className="py-2 px-3 font-mono text-destructive/80">{row.behavior}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
      </div>

        {/* Package Layout */}
        <div>
          <h2 className="text-xl font-heading font-bold text-foreground mb-4 docs-h2">Package Layout</h2>
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-border">
                  <th className="text-left py-2 px-3 font-mono text-muted-foreground font-medium">File</th>
                  <th className="text-left py-2 px-3 font-mono text-muted-foreground font-medium">Responsibility</th>
                </tr>
              </thead>
              <tbody>
                {[
                  { file: "release/semver.go", desc: "Version parsing, padding, comparison, bumping" },
                  { file: "release/metadata.go", desc: "Read/write .gitmap/release/*.json, latest.json, version.json" },
                  { file: "release/gitops.go", desc: "Branch, tag, push, checkout Git operations" },
                  { file: "release/github.go", desc: "Asset collection, changelog/readme detection" },
                  { file: "release/workflow.go", desc: "Orchestration: Execute(), metadata-first sequencing" },
                  { file: "release/workflowbranch.go", desc: "release-branch and release-pending workflows" },
                  { file: "release/assets.go", desc: "Go cross-compilation orchestration" },
                  { file: "release/assetstargets.go", desc: "Target matrix, config parsing, resolution" },
                  { file: "release/assetsupload.go", desc: "GitHub API upload with retry" },
                  { file: "release/compress.go", desc: ".zip/.tar.gz archive creation" },
                  { file: "release/checksums.go", desc: "SHA256 checksum generation" },
                  { file: "release/ziparchive.go", desc: "Zip group archive creation (max compression)" },
                ].map((row) => (
                  <tr key={row.file} className="border-b border-border/50">
                    <td className="py-2 px-3 font-mono text-primary">{row.file}</td>
                    <td className="py-2 px-3 text-muted-foreground">{row.desc}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>

        {/* See Also */}
        <section className="mt-10">
          <h2 className="text-xl font-heading font-bold text-foreground mb-4 docs-h2">See Also</h2>
          <ul className="space-y-1 text-sm font-mono">
            <li><a href="/zip-group" className="text-primary hover:underline">zip-group</a> — Manage named file/folder collections for release bundling <span className="text-muted-foreground">↗</span></li>
            <li><a href="/clear-release-json" className="text-primary hover:underline">clear-release-json</a> — Remove stale release metadata files <span className="text-muted-foreground">↗</span></li>
          </ul>
        </section>
      </div>
    </DocsLayout>
  );
};

export default ReleasePage;
