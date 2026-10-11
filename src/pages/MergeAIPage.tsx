import { useState, useMemo } from "react";
import DocsLayout from "@/components/docs/DocsLayout";
import CodeBlock from "@/components/docs/CodeBlock";
import TerminalDemo, { TerminalLine } from "@/components/docs/TerminalDemo";
import { TerminalLineType } from "@/constants";
import { copyToClipboard } from "@/lib/clipboard";
import {
  GitMerge,
  Plus,
  Trash2,
  Copy,
  Check,
  CheckCircle2,
  Layers,
  FileCode,
  FileText,
  Sliders,
  Sparkles,
  ArrowRight,
  ShieldCheck,
  Clock,
  Terminal,
} from "lucide-react";

type MergeModality = "positional" | "flags" | "comma" | "manifest" | "config";

const flags = [
  { flag: "-d, --destination <dir>", def: "none (required)", desc: "Target directory or repository slug where merged files will be staged" },
  { flag: "-s, --source <url|path>", def: "none (required)", desc: "Source repository URL or local path. Can be repeated multiple times or comma-separated" },
  { flag: "-u, --url <url>", def: "none", desc: "Shorthand alias for --source. Repeatable for multiple remote repositories" },
  { flag: "-c, --config <file.json>", def: "none", desc: "Path to JSON configuration file specifying destination and source repositories" },
  { flag: "--no-ai", def: "false", desc: "Stage files into numbered collision buckets without invoking downstream AI consolidation" },
  { flag: "-h, --help", def: "false", desc: "Print comprehensive usage manual, parameter descriptions, and syntax examples" },
];

const exitCodes = [
  { code: "0", meaning: "Success — All source repositories staged, collisions sequenced (01_/02_), manifest and instruction.md generated" },
  { code: "1", meaning: "Resolution or staging failure — Target directory exists without confirmation, or Git clone operation failed" },
  { code: "2", meaning: "Syntax or argument error — Missing required destination or insufficient source repositories" },
];

const collisionMatrix = [
  {
    scenario: "Overlapping relative path (e.g. src/main.go)",
    baseStaged: "src/01_main.go (earlier commit timestamp)",
    incomingStaged: "src/02_main.go (subsequent commit timestamp)",
    manifestEntry: "fileCollisions[] with variants array",
    aiAction: "Synthesize union into canonical main.go, then delete 01_ and 02_ files",
  },
  {
    scenario: "Unique file in Source 1 only (e.g. src/auth.go)",
    baseStaged: "src/auth.go (direct placement)",
    incomingStaged: "None",
    manifestEntry: "Recorded in repoSequence folderTree",
    aiAction: "Preserve intact; adjust package declarations if needed",
  },
  {
    scenario: "Unique file in Source 2 only (e.g. src/billing.go)",
    baseStaged: "None",
    incomingStaged: "src/billing.go (direct placement)",
    manifestEntry: "Recorded in repoSequence folderTree",
    aiAction: "Preserve intact; wire exports into consolidated module",
  },
  {
    scenario: "Project Documentation (README.md)",
    baseStaged: "01_README.md (Base repository)",
    incomingStaged: "02_README.md (Incoming repository)",
    manifestEntry: "fileCollisions[] collision entry",
    aiAction: "Combine feature lists, prerequisites, and architectures into unified README.md",
  },
  {
    scenario: "Configuration files (e.g. config.json)",
    baseStaged: "01_config.json",
    incomingStaged: "02_config.json",
    manifestEntry: "fileCollisions[] collision entry",
    aiAction: "Perform deep key merge across JSON attributes and environment settings",
  },
];

const terminalDemoLines: TerminalLine[] = [
  {
    text: "gitmap merge-ai ./unified-app https://github.com/org/repo-auth.git https://github.com/org/repo-billing.git",
    type: TerminalLineType.Input,
  },
  {
    text: "● Resolving source repositories & commit histories...",
    type: TerminalLineType.Header,
  },
  {
    text: "  ├── [01] repo-auth    (Last commit: 2026-09-15 14:30:00 UTC) -> Base repository",
    type: TerminalLineType.Output,
  },
  {
    text: "  └── [02] repo-billing (Last commit: 2026-10-01 09:15:00 UTC) -> Incoming repository",
    type: TerminalLineType.Output,
  },
  {
    text: "● Analyzing file paths and sequencing collisions...",
    type: TerminalLineType.Header,
  },
  {
    text: "  ✔ Placed 4 unique files directly into destination tree",
    type: TerminalLineType.Output,
  },
  {
    text: "  ! Collision: 'src/main.go' -> Sequenced to 01_main.go & 02_main.go",
    type: TerminalLineType.Accent,
  },
  {
    text: "  ! Collision: 'README.md'   -> Sequenced to 01_README.md & 02_README.md",
    type: TerminalLineType.Accent,
  },
  {
    text: "● Generating AI handoff artifacts...",
    type: TerminalLineType.Header,
  },
  {
    text: "  ✔ Wrote merge-ai-manifest.json (2 collision pairs recorded)",
    type: TerminalLineType.Output,
  },
  {
    text: "  ✔ Authored instruction.md with 6-step AI consolidation checklist",
    type: TerminalLineType.Output,
  },
  {
    text: "✔ Staging complete! Ready for downstream AI agent consolidation.",
    type: TerminalLineType.Accent,
  },
];

const sampleManifestJson = `{
  "attributes": {
    "generatedAt": "2026-10-11T19:00:00Z",
    "gitMapVersion": "6.534.0",
    "tool": "gitmap merge-ai",
    "collisionStrategy": "prefixed_sequence_01_02"
  },
  "data": {
    "destination": {
      "targetType": "local_directory",
      "slug": "unified-app",
      "resolvedPath": "./unified-app"
    },
    "repoSequence": [
      {
        "sequenceOrder": 1,
        "repoUrl": "https://github.com/org/repo-auth.git",
        "branch": "main",
        "lastCommitDate": "2026-09-15T14:30:00Z"
      },
      {
        "sequenceOrder": 2,
        "repoUrl": "https://github.com/org/repo-billing.git",
        "branch": "main",
        "lastCommitDate": "2026-10-01T09:15:00Z"
      }
    ],
    "fileCollisions": [
      {
        "canonicalPath": "src/main.go",
        "status": "pending_ai_consolidation",
        "variants": [
          { "sequence": "01", "file": "src/01_main.go", "sourceRepo": "repo-auth" },
          { "sequence": "02", "file": "src/02_main.go", "sourceRepo": "repo-billing" }
        ]
      }
    ]
  }
}`;

const sampleInstructionMd = `# AI Merge Consolidation Instruction & Verification Guide

This repository has been prepared by \`gitmap merge-ai\`. Multiple source repositories have been assembled into this staging tree.

## Non-Negotiable Instructions for the Consolidating AI Agent:
1. **Preserve All Business Logic:** Do NOT delete features, endpoints, or data models. Consolidate overlapping files (\`01_*\`, \`02_*\`) into canonical names (e.g. synthesize \`01_main.go\` and \`02_main.go\` into \`main.go\`).
2. **Flatten Sequence Prefixes:** Once consolidated, delete temporary \`01_\` and \`02_\` variant files.
3. **Reference the Manifest:** Open and inspect \`merge-ai-manifest.json\` for repository histories and collision maps.
4. **Zero Code Deletion:** When resolving duplicate utilities, favor additive composition and unified interfaces.`;

const MergeAIPage = () => {
  const [destination, setDestination] = useState("./unified-app");
  const [sources, setSources] = useState<string[]>([
    "https://github.com/org/repo-auth.git",
    "https://github.com/org/repo-billing.git",
  ]);
  const [modality, setModality] = useState<MergeModality>("positional");
  const [manifestFile, setManifestFile] = useState("sources.txt");
  const [configFile, setConfigFile] = useState("merge-config.json");
  const [isCopied, setIsCopied] = useState(false);

  const handleAddSource = () => {
    setSources((prev) => [...prev, ""]);
  };

  const handleUpdateSource = (index: number, value: string) => {
    setSources((prev) => {
      const updated = [...prev];
      updated[index] = value;
      return updated;
    });
  };

  const handleRemoveSource = (index: number) => {
    setSources((prev) => prev.filter((_, i) => i !== index));
  };

  const generatedCommand = useMemo(() => {
    const dest = destination.trim() || "./target";
    const validSources = sources.filter((s) => s.trim().length > 0);

    if (modality === "positional") {
      const srcList = validSources.length > 0 ? validSources.join(" ") : "<source1> <source2>";
      return `gitmap merge-ai ${dest} ${srcList}`;
    }

    if (modality === "flags") {
      const srcFlags = validSources.length > 0 ? validSources.map((s) => `-s ${s}`).join(" ") : "-s <source1> -s <source2>";
      return `gitmap ma -d ${dest} ${srcFlags}`;
    }

    if (modality === "comma") {
      const commaList = validSources.length > 0 ? validSources.join(",") : "<source1>,<source2>";
      return `gitmap ma -d ${dest} --source "${commaList}"`;
    }

    if (modality === "manifest") {
      return `gitmap merge-ai --destination ${dest} ${manifestFile}`;
    }

    return `gitmap merge-ai --config ${configFile}`;
  }, [destination, sources, modality, manifestFile, configFile]);

  const jsonConfigPreview = useMemo(() => {
    return JSON.stringify(
      {
        destination: destination.trim() || "./unified-app",
        sources: sources.filter((s) => s.trim().length > 0),
      },
      null,
      2,
    );
  }, [destination, sources]);

  const handleCopy = () => {
    copyToClipboard(generatedCommand);
    setIsCopied(true);
    setTimeout(() => setIsCopied(false), 2000);
  };

  return (
    <DocsLayout>
      <div className="max-w-4xl space-y-10">
        <div>
          <div className="flex flex-wrap items-center gap-3 mb-2">
            <GitMerge className="h-8 w-8 text-primary" />
            <h1 className="text-3xl font-bold tracking-tight">merge-ai</h1>
            <span className="font-mono text-xs px-2 py-1 rounded bg-primary/10 text-foreground border border-primary/20 dark:bg-primary/15 dark:text-primary dark:border-primary/40 transition-colors">
              alias: ma
            </span>
            <span className="inline-flex items-center gap-1 text-[11px] font-mono px-2 py-0.5 rounded border bg-primary/10 text-foreground border-primary/20 dark:bg-primary/20 dark:text-primary dark:border-primary/40">
              <CheckCircle2 className="h-3 w-3" />
              Shipped (v6.534.0)
            </span>
          </div>
          <p className="text-lg text-muted-foreground leading-relaxed">
            Multi-source Git repository and directory amalgamation engine with non-destructive collision sequencing (
            <code className="docs-inline-code">01_</code> / <code className="docs-inline-code">02_</code>), machine-readable manifest generation, and actionable AI consolidation checklist authoring.
          </p>
          <p className="text-xs text-muted-foreground mt-2">
            Spec: <code>02-spec/21-app/git-prime-machine-os-version-and-ui-help-docs-improvement/02-component-and-cli-spec.md</code>
          </p>
        </div>

        <section>
          <h2 className="text-xl font-semibold mb-3 flex items-center gap-2">
            <Sliders className="h-5 w-5 text-primary" /> Interactive Parameter Builder & Live Generator
          </h2>
          <p className="text-sm text-muted-foreground mb-4">
            Select your preferred input modality and configure destination and source parameters. The exact shell command is generated dynamically in real time.
          </p>

          <div className="rounded-lg border border-border bg-card p-5 space-y-6">
            <div className="flex flex-wrap gap-2 border-b border-border pb-4">
              {[
                { id: "positional", label: "1. Positional CLI" },
                { id: "flags", label: "2. Named Flags" },
                { id: "comma", label: "3. Comma-Separated" },
                { id: "manifest", label: "4. Text Manifest" },
                { id: "config", label: "5. JSON Config" },
              ].map((tab) => (
                <button
                  key={tab.id}
                  onClick={() => setModality(tab.id as MergeModality)}
                  className={`px-3 py-1.5 text-xs font-mono rounded-md transition-colors ${
                    modality === tab.id
                      ? "bg-primary text-primary-foreground font-semibold"
                      : "bg-muted/50 hover:bg-muted text-muted-foreground hover:text-foreground border border-border"
                  }`}
                >
                  {tab.label}
                </button>
              ))}
            </div>

            <div className="space-y-4">
              <div>
                <label className="block text-xs font-mono font-medium text-muted-foreground mb-1">
                  Destination Directory (--destination, -d):
                </label>
                <input
                  type="text"
                  value={destination}
                  onChange={(e) => setDestination(e.target.value)}
                  placeholder="./unified-app"
                  className="w-full rounded-md border border-input bg-background px-3 py-2 text-sm font-mono text-foreground focus:outline-none focus:ring-2 focus:ring-primary"
                />
              </div>

              {modality !== "manifest" && modality !== "config" && (
                <div className="space-y-2">
                  <div className="flex items-center justify-between">
                    <label className="block text-xs font-mono font-medium text-muted-foreground">
                      Source Repositories or Directories (--source, -s):
                    </label>
                    <button
                      onClick={handleAddSource}
                      className="inline-flex items-center gap-1 text-xs font-mono text-primary hover:underline"
                    >
                      <Plus className="h-3 w-3" /> Add Source
                    </button>
                  </div>
                  {sources.map((src, idx) => (
                    <div key={idx} className="flex items-center gap-2">
                      <span className="font-mono text-xs text-muted-foreground w-6">#{idx + 1}</span>
                      <input
                        type="text"
                        value={src}
                        onChange={(e) => handleUpdateSource(idx, e.target.value)}
                        placeholder="https://github.com/org/repo.git or ../local-repo"
                        className="flex-1 rounded-md border border-input bg-background px-3 py-1.5 text-xs font-mono text-foreground focus:outline-none focus:ring-2 focus:ring-primary"
                      />
                      {sources.length > 1 && (
                        <button
                          onClick={() => handleRemoveSource(idx)}
                          className="p-1.5 text-muted-foreground hover:text-destructive rounded hover:bg-muted"
                          title="Remove source"
                        >
                          <Trash2 className="h-3.5 w-3.5" />
                        </button>
                      )}
                    </div>
                  ))}
                </div>
              )}

              {modality === "manifest" && (
                <div className="space-y-2">
                  <label className="block text-xs font-mono font-medium text-muted-foreground">
                    Manifest File Path:
                  </label>
                  <input
                    type="text"
                    value={manifestFile}
                    onChange={(e) => setManifestFile(e.target.value)}
                    placeholder="sources.txt"
                    className="w-full rounded-md border border-input bg-background px-3 py-2 text-sm font-mono text-foreground focus:outline-none focus:ring-2 focus:ring-primary"
                  />
                  <div className="rounded border border-border bg-muted/30 p-3 text-xs font-mono text-muted-foreground">
                    <div className="text-primary font-semibold mb-1"># Example {manifestFile} format:</div>
                    <div>https://github.com/org/service-auth.git</div>
                    <div>https://github.com/org/service-billing.git</div>
                    <div>https://github.com/org/service-notifications.git</div>
                  </div>
                </div>
              )}

              {modality === "config" && (
                <div className="space-y-2">
                  <label className="block text-xs font-mono font-medium text-muted-foreground">
                    Config File Path (--config, -c):
                  </label>
                  <input
                    type="text"
                    value={configFile}
                    onChange={(e) => setConfigFile(e.target.value)}
                    placeholder="merge-config.json"
                    className="w-full rounded-md border border-input bg-background px-3 py-2 text-sm font-mono text-foreground focus:outline-none focus:ring-2 focus:ring-primary"
                  />
                  <div className="rounded border border-border bg-muted/30 p-3">
                    <div className="text-xs font-mono text-primary font-semibold mb-1"># Generated {configFile} preview:</div>
                    <pre className="text-xs font-mono overflow-x-auto text-muted-foreground">
                      <code>{jsonConfigPreview}</code>
                    </pre>
                  </div>
                </div>
              )}
            </div>

            <div className="pt-2">
              <div className="flex items-center justify-between mb-2">
                <span className="text-xs font-mono text-muted-foreground font-semibold flex items-center gap-1.5">
                  <Terminal className="h-3.5 w-3.5 text-primary" /> Generated Shell Command:
                </span>
                <button
                  onClick={handleCopy}
                  className="inline-flex items-center gap-1.5 text-xs font-mono px-2 py-1 rounded border border-border bg-muted/40 hover:bg-muted text-foreground transition-colors"
                >
                  {isCopied ? <Check className="h-3.5 w-3.5 text-primary" /> : <Copy className="h-3.5 w-3.5" />}
                  {isCopied ? "Copied!" : "Copy Command"}
                </button>
              </div>
              <div className="rounded-md border border-border bg-code-bg p-3 font-mono text-sm overflow-x-auto text-primary">
                <code>{generatedCommand}</code>
              </div>
            </div>
          </div>
        </section>

        <section>
          <h2 className="text-xl font-semibold mb-3 flex items-center gap-2">
            <Terminal className="h-5 w-5 text-primary" /> Execution Simulator
          </h2>
          <p className="text-sm text-muted-foreground mb-4">
            Interactive terminal replay demonstrating chronological repository ordering, collision detection, and handoff artifact generation.
          </p>
          <TerminalDemo
            title="gitmap merge-ai execution simulator"
            lines={terminalDemoLines}
            autoPlay={false}
          />
        </section>

        <section>
          <h2 className="text-xl font-semibold mb-3 flex items-center gap-2">
            <Layers className="h-5 w-5 text-primary" /> File Collision Resolution Matrix
          </h2>
          <p className="text-sm text-muted-foreground mb-4">
            When merging multiple repositories, files with identical relative paths never overwrite each other silently. GitMap sequences collisions using deterministic chronological prefixes.
          </p>
          <div className="overflow-x-auto">
            <table className="w-full text-xs border border-border rounded-lg">
              <thead>
                <tr className="bg-muted/50 border-b border-border">
                  <th className="text-left px-3 py-2 font-medium">Collision Scenario</th>
                  <th className="text-left px-3 py-2 font-medium">Base Staged Variant</th>
                  <th className="text-left px-3 py-2 font-medium">Incoming Staged Variant</th>
                  <th className="text-left px-3 py-2 font-medium">Manifest Registry</th>
                  <th className="text-left px-3 py-2 font-medium">Downstream AI Action</th>
                </tr>
              </thead>
              <tbody>
                {collisionMatrix.map((row, idx) => (
                  <tr key={idx} className="border-t border-border hover:bg-muted/20">
                    <td className="px-3 py-2 font-mono font-semibold text-foreground">{row.scenario}</td>
                    <td className="px-3 py-2 font-mono text-primary">{row.baseStaged}</td>
                    <td className="px-3 py-2 font-mono text-muted-foreground">{row.incomingStaged}</td>
                    <td className="px-3 py-2 text-muted-foreground">{row.manifestEntry}</td>
                    <td className="px-3 py-2 text-muted-foreground">{row.aiAction}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </section>

        <section>
          <h2 className="text-xl font-semibold mb-3 flex items-center gap-2">
            <FileCode className="h-5 w-5 text-primary" /> Generated Staging Artifacts
          </h2>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="rounded-lg border border-border bg-card p-4 space-y-2">
              <div className="flex items-center gap-2 text-sm font-semibold text-foreground">
                <FileText className="h-4 w-4 text-primary" />
                <span>merge-ai-manifest.json</span>
              </div>
              <p className="text-xs text-muted-foreground">
                Structured JSON catalog describing repository sequence order, commit hashes, timestamps, and collision variant pairings.
              </p>
              <div className="pt-2">
                <CodeBlock code={sampleManifestJson} language="json" />
              </div>
            </div>

            <div className="rounded-lg border border-border bg-card p-4 space-y-2">
              <div className="flex items-center gap-2 text-sm font-semibold text-foreground">
                <Sparkles className="h-4 w-4 text-primary" />
                <span>instruction.md</span>
              </div>
              <p className="text-xs text-muted-foreground">
                Actionable markdown checklist designed for downstream AI coding agents to consolidate sequenced files and delete temporary prefixes.
              </p>
              <div className="pt-2">
                <CodeBlock code={sampleInstructionMd} language="markdown" />
              </div>
            </div>
          </div>
        </section>

        <section>
          <h2 className="text-xl font-semibold mb-3">Flags Reference</h2>
          <div className="overflow-x-auto">
            <table className="w-full text-sm border border-border rounded-lg">
              <thead>
                <tr className="bg-muted/50 border-b border-border">
                  <th className="text-left px-4 py-2 font-medium">Flag</th>
                  <th className="text-left px-4 py-2 font-medium">Default</th>
                  <th className="text-left px-4 py-2 font-medium">Description</th>
                </tr>
              </thead>
              <tbody>
                {flags.map((f) => (
                  <tr key={f.flag} className="border-t border-border">
                    <td className="px-4 py-2 font-mono text-primary">{f.flag}</td>
                    <td className="px-4 py-2 font-mono text-muted-foreground text-xs">{f.def}</td>
                    <td className="px-4 py-2 text-muted-foreground text-xs">{f.desc}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </section>

        <section>
          <h2 className="text-xl font-semibold mb-3">Exit Codes</h2>
          <div className="overflow-x-auto">
            <table className="w-full text-sm border border-border rounded-lg">
              <thead>
                <tr className="bg-muted/50 border-b border-border">
                  <th className="text-left px-4 py-2 font-medium">Code</th>
                  <th className="text-left px-4 py-2 font-medium">Meaning</th>
                </tr>
              </thead>
              <tbody>
                {exitCodes.map((e) => (
                  <tr key={e.code} className="border-t border-border">
                    <td className="px-4 py-2 font-mono text-primary font-bold">{e.code}</td>
                    <td className="px-4 py-2 text-muted-foreground text-xs">{e.meaning}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </section>

        <section>
          <h2 className="text-xl font-semibold mb-3 flex items-center gap-2">
            <ShieldCheck className="h-5 w-5 text-primary" /> AI Consolidation Protocol & Guidelines
          </h2>
          <div className="rounded-lg border border-border bg-card p-4 space-y-3 text-xs leading-relaxed text-muted-foreground">
            <div className="flex items-start gap-2">
              <CheckCircle2 className="h-4 w-4 text-primary shrink-0 mt-0.5" />
              <div>
                <strong className="text-foreground">Preserve Full Business Logic:</strong> Overlapping services or endpoints must be additively composed into the single canonical target file. Zero code or endpoint deletion.
              </div>
            </div>
            <div className="flex items-start gap-2">
              <CheckCircle2 className="h-4 w-4 text-primary shrink-0 mt-0.5" />
              <div>
                <strong className="text-foreground">Flatten Sequence Prefixes:</strong> After synthesizing the canonical implementation, remove all <code className="docs-inline-code">01_*</code> and <code className="docs-inline-code">02_*</code> files so clean project filenames remain.
              </div>
            </div>
            <div className="flex items-start gap-2">
              <CheckCircle2 className="h-4 w-4 text-primary shrink-0 mt-0.5" />
              <div>
                <strong className="text-foreground">Atomic Commit Ceremony:</strong> Once consolidated, verify builds and run <code className="docs-inline-code">gitmap cpf "feat: consolidate multi-repo merge into unified codebase"</code> to finalize.
              </div>
            </div>
          </div>
        </section>
      </div>
    </DocsLayout>
  );
};

export default MergeAIPage;
