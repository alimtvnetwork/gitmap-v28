import DocsLayout from "@/components/docs/DocsLayout";
import CodeBlock from "@/components/docs/CodeBlock";
import { GitBranch, FileCode, Shield, Eye } from "lucide-react";

const features = [
  {
    icon: GitBranch,
    title: "Branch Safety",
    desc: "Creates a backup branch before changes and a feature branch for the rename, with auto-merge back to your starting branch.",
  },
  {
    icon: FileCode,
    title: "All-Files-by-Default",
    desc: "Scans every file in the repo for matching module paths. Use --ext to restrict to specific extensions.",
  },
  {
    icon: Shield,
    title: "Dirty-Tree Guard",
    desc: "Refuses to run if your working tree has uncommitted changes, preventing accidental data loss.",
  },
  {
    icon: Eye,
    title: "Dry Run Preview",
    desc: "See exactly which files would change without modifying anything using --dry-run.",
  },
];

const GoModPage = () => {
  return (
    <DocsLayout>
      <h1 className="text-3xl font-heading font-bold mb-2 docs-h1">GoMod Command</h1>
      <p className="text-muted-foreground mb-6">
        Rename a Go module path across an entire repository — updates{" "}
        <code className="text-primary font-mono">go.mod</code> and every file
        that references the old module path. Wraps the operation in a safe
        branch workflow with backup and auto-merge.
      </p>

      {/* Features */}
      <h2 className="text-xl font-heading font-semibold mt-8 mb-4">Features</h2>
      <div className="grid md:grid-cols-2 gap-4 mb-8">
        {features.map((f) => (
          <div
            key={f.title}
            className="rounded-lg border border-border bg-card p-4"
          >
            <f.icon className="h-5 w-5 text-primary mb-2" />
            <h3 className="font-mono font-semibold text-sm mb-1">{f.title}</h3>
            <p className="text-xs text-muted-foreground">{f.desc}</p>
          </div>
        ))}
      </div>

      {/* Usage */}
      <h2 className="text-xl font-heading font-semibold mt-10 mb-3">Usage</h2>

      <CodeBlock
        code="gitmap gomod github.com/new/module"
        title="Rename module path in all files"
      />
      <CodeBlock
        code='gitmap gm github.com/new/module --ext "*.go,*.md"'
        title="Only replace in .go and .md files"
      />
      <CodeBlock
        code="gitmap gomod github.com/new/module --dry-run"
        title="Preview changes without modifying anything"
      />
      <CodeBlock
        code="gitmap gomod github.com/new/module --no-merge"
        title="Commit on feature branch but don't merge back"
      />
      <CodeBlock
        code="gitmap gomod github.com/new/module --no-tidy --verbose"
        title="Skip go mod tidy and print each file"
      />

      {/* Flags */}
      <h2 className="text-xl font-heading font-semibold mt-10 mb-3">Flags</h2>
      <div className="rounded-lg border border-border overflow-hidden">
        <table className="w-full text-sm">
          <thead>
            <tr className="bg-muted/50">
              <th className="text-left font-mono font-semibold px-4 py-2">
                Flag
              </th>
              <th className="text-left font-mono font-semibold px-4 py-2">
                Description
              </th>
            </tr>
          </thead>
          <tbody>
            {[
              ["--dry-run", "Preview changes without modifying files or branches"],
              ["--no-merge", "Commit on feature branch but do not merge back"],
              ["--no-tidy", "Skip go mod tidy after replacement"],
              ["--verbose", "Print each file path as it is modified"],
              ['--ext <exts>', 'Comma-separated file extensions to filter (e.g. "*.go,*.md"); default: all files'],
            ].map(([flag, desc]) => (
              <tr key={flag} className="border-t border-border">
                <td className="px-4 py-2 font-mono text-primary">{flag}</td>
                <td className="px-4 py-2 text-muted-foreground">{desc}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* Terminal preview */}
      <h2 className="text-xl font-heading font-semibold mt-10 mb-3">
        Example Output
      </h2>
      <div className="rounded-lg border border-border overflow-hidden my-4">
        <div className="bg-terminal px-4 py-2 flex items-center gap-2 border-b border-border">
          <div className="flex gap-1.5">
            <span className="w-3 h-3 rounded-full bg-red-500/80" />
            <span className="w-3 h-3 rounded-full bg-yellow-500/80" />
            <span className="w-3 h-3 rounded-full bg-green-500/80" />
          </div>
          <span className="text-xs font-mono text-muted-foreground ml-2">
            gitmap gomod
          </span>
        </div>
        <div className="bg-terminal p-4 font-mono text-sm leading-relaxed">
          <div className="text-primary">✔ Module path renamed</div>
          <div className="text-terminal-foreground">
            {"  "}Old: github.com/oldorg/core
          </div>
          <div className="text-terminal-foreground">
            {"  "}New: github.com/neworg/core
          </div>
          <div className="text-terminal-foreground">
            {"  "}Files updated: 47
          </div>
          <div className="text-terminal-foreground">
            {"  "}Backup branch: backup/before-replace-github-com-neworg-core
          </div>
          <div className="text-terminal-foreground">
            {"  "}Feature branch: feature/replace-github-com-neworg-core
          </div>
          <div className="text-terminal-foreground">
            {"  "}Merged into: main
          </div>
        </div>
      </div>

      {/* Dry run output */}
      <h2 className="text-xl font-heading font-semibold mt-10 mb-3">
        Dry Run Output
      </h2>
      <div className="rounded-lg border border-border overflow-hidden my-4">
        <div className="bg-terminal px-4 py-2 flex items-center gap-2 border-b border-border">
          <div className="flex gap-1.5">
            <span className="w-3 h-3 rounded-full bg-red-500/80" />
            <span className="w-3 h-3 rounded-full bg-yellow-500/80" />
            <span className="w-3 h-3 rounded-full bg-green-500/80" />
          </div>
          <span className="text-xs font-mono text-muted-foreground ml-2">
            gitmap gomod --dry-run
          </span>
        </div>
        <div className="bg-terminal p-4 font-mono text-sm leading-relaxed">
          <div className="text-yellow-400">
            gomod (dry-run): would rename module path
          </div>
          <div className="text-terminal-foreground">
            {"  "}Old: github.com/oldorg/core
          </div>
          <div className="text-terminal-foreground">
            {"  "}New: github.com/neworg/core
          </div>
          <div className="text-terminal-foreground">
            {"  "}Files that would change: 47
          </div>
          <div className="text-muted-foreground mt-1">
            {"  "}cmd/root.go
          </div>
          <div className="text-muted-foreground">
            {"  "}cmd/serve.go
          </div>
          <div className="text-muted-foreground">
            {"  "}pkg/handler/auth.go
          </div>
          <div className="text-muted-foreground">{"  "}...</div>
        </div>
      </div>

      {/* Edge cases */}
      <h2 className="text-xl font-heading font-semibold mt-10 mb-3">
        Edge Cases
      </h2>
      <div className="rounded-lg border border-border overflow-hidden">
        <table className="w-full text-sm">
          <thead>
            <tr className="bg-muted/50">
              <th className="text-left font-mono font-semibold px-4 py-2">
                Scenario
              </th>
              <th className="text-left font-mono font-semibold px-4 py-2">
                Behavior
              </th>
            </tr>
          </thead>
          <tbody>
            {[
              ["No go.mod in directory", "Error: go.mod not found in current directory"],
              ["Not inside a Git repo", "Error: not inside a git repository"],
              ["Old path == new path", "Info: module path is already <path>, nothing to rename"],
              ["Branch already exists", "Error: branch <name> already exists, aborting"],
              ["No files contain old path", "Warning: only go.mod updated"],
              ["Dirty working tree", "Error: commit or stash first"],
              ["Merge conflict", "Error: resolve manually on <branch>"],
              ["go mod tidy fails", "Warning: continuing without tidy"],
            ].map(([scenario, behavior]) => (
              <tr key={scenario} className="border-t border-border">
                <td className="px-4 py-2 font-mono text-primary text-xs">
                  {scenario}
                </td>
                <td className="px-4 py-2 text-muted-foreground text-xs">
                  {behavior}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* File layout */}
      <h2 className="text-xl font-heading font-semibold mt-10 mb-3">
        File Layout
      </h2>
      <div className="rounded-lg border border-border overflow-hidden">
        <table className="w-full text-sm">
          <thead>
            <tr className="bg-muted/50">
              <th className="text-left font-mono font-semibold px-4 py-2">
                File
              </th>
              <th className="text-left font-mono font-semibold px-4 py-2">
                Purpose
              </th>
            </tr>
          </thead>
          <tbody>
            {[
              ["constants/constants_gomod.go", "Command names, flag names, messages, error strings"],
              ["cmd/gomod.go", "Flag parsing, orchestration, summary output"],
              ["cmd/gomodreplace.go", "File walking, path replacement logic"],
              ["cmd/gomodbranch.go", "Branch creation, merge, slug generation"],
            ].map(([file, purpose]) => (
              <tr key={file} className="border-t border-border">
                <td className="px-4 py-2 font-mono text-primary">{file}</td>
                <td className="px-4 py-2 text-muted-foreground">{purpose}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </DocsLayout>
  );
};

export default GoModPage;
