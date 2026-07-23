import { motion } from "framer-motion";
import DocsLayout from "@/components/docs/DocsLayout";
import CodeBlock from "@/components/docs/CodeBlock";
import {
  Search, FolderGit2, Code2, FileCode, Braces, Cpu, Database,
  GitBranch, ArrowRight, Layers, ShieldCheck, FileJson, Terminal,
} from "lucide-react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";

const projectTypes = [
  { icon: Code2, type: "Go", key: "go", indicator: "go.mod", color: "text-cyan-400", bg: "bg-cyan-400/10", commands: "go-repos (gr)", desc: "Detects Go modules via go.mod, extracts module name, Go version, and runnable entry points from cmd/ subdirectories." },
  { icon: Braces, type: "Node.js", key: "node", indicator: "package.json", color: "text-green-400", bg: "bg-green-400/10", commands: "node-repos (nr)", desc: "Identifies Node.js projects by package.json presence. Reclassified as React if react dependency found." },
  { icon: FileCode, type: "React", key: "react", indicator: "package.json + react", color: "text-blue-400", bg: "bg-blue-400/10", commands: "react-repos (rr)", desc: "Promoted from Node.js when react, @types/react, react-scripts, next, gatsby, or remix is a dependency." },
  { icon: Cpu, type: "C++", key: "cpp", indicator: "CMakeLists.txt", color: "text-orange-400", bg: "bg-orange-400/10", commands: "cpp-repos (cr)", desc: "Matches CMakeLists.txt, *.vcxproj, or meson.build. Extracts project name from CMake project() directive." },
  { icon: FolderGit2, type: "C#", key: "csharp", indicator: "*.csproj / *.sln", color: "text-purple-400", bg: "bg-purple-400/10", commands: "csharp-repos (csr)", desc: "Parses .csproj XML for TargetFramework/OutputType/Sdk. Solution files (.sln) take precedence over standalone projects." },
];

const detectionRules = [
  { type: "Go", primary: "go.mod exists", secondary: "go.sum, *.go files", falsePositive: "Ignores vendor/, testdata/" },
  { type: "Node.js", primary: "package.json exists", secondary: "package-lock.json, yarn.lock, pnpm-lock.yaml, bun.lock", falsePositive: "Ignores node_modules/, vendor/" },
  { type: "React", primary: "package.json with react dependency", secondary: "@types/react, react-scripts, next, gatsby, remix", falsePositive: "Reclassified from Node.js — exclusive" },
  { type: "C++", primary: "CMakeLists.txt, *.vcxproj, meson.build", secondary: "Makefile + C++ sources, conanfile, vcpkg.json", falsePositive: "Ignores build/, cmake-build-*/, out/" },
  { type: "C#", primary: "*.csproj, *.sln", secondary: "*.fsproj, global.json, *.cs files", falsePositive: "Ignores bin/, obj/, packages/. .sln takes precedence" },
];

const excludeDirs = [
  "node_modules", "vendor", ".git", "dist", "build", "target",
  "bin", "obj", "out", "cmake-build-*", "testdata", "packages", ".venv", ".cache",
];

const dbTables = [
  { name: "ProjectTypes", cols: "Id, Key, Label", desc: "Reference table seeded with go, node, react, cpp, csharp" },
  { name: "DetectedProjects", cols: "Id, RepoId, TypeId, Name, AbsPath, RelPath, Indicator, DetectedAt", desc: "One row per project discovered during scan" },
  { name: "GoProjectMetadata", cols: "ProjectId, ModuleName, GoVersion, GoModPath, GoSumPath", desc: "Go-specific module metadata from go.mod" },
  { name: "GoRunnableFiles", cols: "Id, ProjectId, Name, RelativePath", desc: "Entry points found in cmd/*/main.go or root main.go" },
  { name: "CSharpProjectMetadata", cols: "ProjectId, SlnName, SdkVersion", desc: "C# solution-level metadata from .sln and global.json" },
  { name: "CSharpProjectFiles", cols: "Id, ProjectId, FileName, TargetFramework, OutputType, Sdk", desc: "Individual .csproj file properties parsed from XML" },
  { name: "CSharpKeyFiles", cols: "Id, ProjectId, FileName, Category", desc: "Significant files like .editorconfig, Directory.Build.props" },
];

const pipelineSteps = [
  { step: "1", title: "Tree Walk", desc: "Recursively traverse repo directory, skipping excluded dirs", icon: Search },
  { step: "2", title: "Indicator Match", desc: "Check each directory for primary indicators (go.mod, package.json, etc.)", icon: FileCode },
  { step: "3", title: "Type Classification", desc: "Assign project type; React promoted from Node.js if react dep found", icon: Layers },
  { step: "4", title: "Metadata Extraction", desc: "Parse language-specific files for deep insights (modules, frameworks)", icon: FileJson },
  { step: "5", title: "DB Upsert", desc: "Persist projects + metadata to SQLite with stale record cleanup", icon: Database },
  { step: "6", title: "JSON Export", desc: "Write per-type JSON files (go-projects.json, etc.) to output dir", icon: Terminal },
];

const fileLayout = [
  ["detector/detector.go", "Detection orchestration — tree walk + collect DetectedProject[]"],
  ["detector/rules.go", "Per-type indicator matching, exclusion checks, React dep list"],
  ["detector/goparser.go", "Parse go.mod (module name, Go version), locate go.sum, scan cmd/ for runnables"],
  ["detector/csharpparser.go", "Parse .csproj XML, find .sln, parse global.json SDK version, collect key files"],
  ["detector/parser.go", "Parse package.json name, check React deps, parse CMakeLists.txt project()"],
  ["model/project.go", "DetectedProject struct"],
  ["model/projecttype.go", "ProjectType struct"],
  ["model/gometadata.go", "GoProjectMetadata + GoRunnableFile structs"],
  ["model/csharpmetadata.go", "CSharpProjectMetadata + CSharpProjectFile + CSharpKeyFile structs"],
  ["store/project.go", "Upsert, query by type, count, delete stale projects"],
  ["store/projecttype.go", "Seed ProjectTypes, query by key"],
  ["store/gometadata.go", "Go metadata + runnable CRUD operations"],
  ["store/csharpmetadata.go", "C# metadata + project files + key files CRUD"],
  ["cmd/projectrepos.go", "Query command handler — dispatch by type, --json/--count flags"],
  ["cmd/projectreposoutput.go", "Terminal + JSON formatting for project query results"],
  ["constants/constants_project.go", "Detection IDs, keys, table names, indicators, exclusions"],
  ["constants/constants_project_sql.go", "All SQL: create, seed, upsert, query, stale cleanup, drop"],
];

const fade = { initial: { opacity: 0, y: 10 }, animate: { opacity: 1, y: 0 } };

const ProjectDetectionPage = () => (
  <DocsLayout>
    {/* Hero */}
    <motion.div {...fade} transition={{ duration: 0.3 }}>
      <h1 className="text-2xl sm:text-3xl font-mono font-bold text-foreground flex items-center gap-3 mb-2">
        <Search className="h-7 w-7 sm:h-8 sm:w-8 text-primary" />
        Project Detection
      </h1>
      <p className="text-muted-foreground text-sm sm:text-base mb-8">
        Automatic technology stack detection during <code className="text-primary font-mono">scan</code> and <code className="text-primary font-mono">rescan</code>,
        with deep metadata extraction and dedicated query commands per project type.
      </p>
    </motion.div>

    {/* Detection Pipeline */}
    <motion.div {...fade} transition={{ duration: 0.3, delay: 0.05 }}>
      <h2 className="text-xl font-heading font-semibold mt-6 mb-4 flex items-center gap-2">
        <ArrowRight className="h-5 w-5 text-primary" />
        Detection Pipeline
      </h2>
      <div className="grid grid-cols-2 md:grid-cols-3 gap-3 mb-8">
        {pipelineSteps.map((s) => (
          <Card key={s.step} className="border-border bg-card">
            <CardContent className="p-4">
              <div className="flex items-center gap-2 mb-2">
                <span className="font-mono text-xs font-bold text-foreground bg-primary/10 border border-primary/20 dark:bg-primary/15 dark:text-primary rounded-full w-6 h-6 flex items-center justify-center">{s.step}</span>
                <s.icon className="h-4 w-4 text-muted-foreground" />
              </div>
              <h3 className="font-mono font-semibold text-sm mb-1">{s.title}</h3>
              <p className="text-xs text-muted-foreground leading-relaxed">{s.desc}</p>
            </CardContent>
          </Card>
        ))}
      </div>
    </motion.div>

    {/* Supported Types — Tabbed */}
    <motion.div {...fade} transition={{ duration: 0.3, delay: 0.1 }}>
      <h2 className="text-xl font-heading font-semibold mt-10 mb-4 flex items-center gap-2">
        <Layers className="h-5 w-5 text-primary" />
        Supported Project Types
      </h2>
      <Tabs defaultValue="go" className="mb-8">
        <TabsList className="flex flex-wrap h-auto gap-1 bg-muted/50 p-1">
          {projectTypes.map((p) => (
            <TabsTrigger key={p.key} value={p.key} className="font-mono text-xs gap-1.5 data-[state=active]:bg-background">
              <p.icon className={`h-3.5 w-3.5 ${p.color}`} />
              {p.type}
            </TabsTrigger>
          ))}
        </TabsList>
        {projectTypes.map((p) => (
          <TabsContent key={p.key} value={p.key}>
            <Card className="border-border">
              <CardHeader className="pb-3">
                <CardTitle className="flex items-center gap-3 font-mono text-lg">
                  <span className={`p-2 rounded-lg ${p.bg}`}>
                    <p.icon className={`h-5 w-5 ${p.color}`} />
                  </span>
                  {p.type} Detection
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <p className="text-sm text-muted-foreground">{p.desc}</p>
                <div className="grid sm:grid-cols-2 gap-3">
                  <div className="rounded-lg border border-border bg-muted/30 p-3">
                    <p className="text-xs font-mono font-semibold text-foreground mb-1">Primary Indicator</p>
                    <code className="text-sm text-primary">{p.indicator}</code>
                  </div>
                  <div className="rounded-lg border border-border bg-muted/30 p-3">
                    <p className="text-xs font-mono font-semibold text-foreground mb-1">Query Command</p>
                    <code className="text-sm text-primary">{p.commands}</code>
                  </div>
                </div>
              </CardContent>
            </Card>
          </TabsContent>
        ))}
      </Tabs>
    </motion.div>

    {/* Detection Rules Table */}
    <motion.div {...fade} transition={{ duration: 0.3, delay: 0.15 }}>
      <h2 className="text-xl font-heading font-semibold mt-10 mb-3 flex items-center gap-2">
        <ShieldCheck className="h-5 w-5 text-primary" />
        Detection Rules
      </h2>
      <div className="rounded-lg border border-border overflow-hidden overflow-x-auto mb-8">
        <table className="w-full text-sm">
          <thead>
            <tr className="bg-muted/50">
              <th className="text-left font-mono font-semibold px-4 py-2.5">Type</th>
              <th className="text-left font-mono font-semibold px-4 py-2.5">Primary Indicator</th>
              <th className="text-left font-mono font-semibold px-4 py-2.5 hidden sm:table-cell">Secondary</th>
              <th className="text-left font-mono font-semibold px-4 py-2.5">False Positive Prevention</th>
            </tr>
          </thead>
          <tbody>
            {detectionRules.map((r) => (
              <tr key={r.type} className="border-t border-border">
                <td className="px-4 py-2.5 font-mono text-primary font-semibold whitespace-nowrap">{r.type}</td>
                <td className="px-4 py-2.5 text-foreground">{r.primary}</td>
                <td className="px-4 py-2.5 text-muted-foreground hidden sm:table-cell">{r.secondary}</td>
                <td className="px-4 py-2.5 text-muted-foreground">{r.falsePositive}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </motion.div>

    {/* Metadata Extraction Deep Dive */}
    <motion.div {...fade} transition={{ duration: 0.3, delay: 0.2 }}>
      <h2 className="text-xl font-heading font-semibold mt-10 mb-4 flex items-center gap-2">
        <FileJson className="h-5 w-5 text-primary" />
        Metadata Extraction
      </h2>
      <div className="space-y-4 mb-8">
        <Card className="border-border">
          <CardHeader className="pb-2">
            <CardTitle className="font-mono text-base flex items-center gap-2">
              <Code2 className="h-4 w-4 text-cyan-400" /> Go Metadata
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-3 text-sm">
            <p className="text-muted-foreground">
              Parses <code className="text-primary">go.mod</code> to extract module name and Go version.
              Locates <code className="text-primary">go.sum</code> for dependency verification.
              Scans <code className="text-primary">cmd/*/main.go</code> subdirectories and root <code className="text-primary">main.go</code> to
              identify runnable entry points.
            </p>
            <div className="rounded-lg bg-muted/30 border border-border p-3 font-mono text-xs space-y-1">
              <p className="text-muted-foreground">// Example extracted metadata:</p>
              <p><span className="text-cyan-400">Module:</span> github.com/user/my-api</p>
              <p><span className="text-cyan-400">Go Version:</span> 1.22</p>
              <p><span className="text-cyan-400">Runnables:</span> cmd/server/main.go, cmd/worker/main.go</p>
            </div>
          </CardContent>
        </Card>

        <Card className="border-border">
          <CardHeader className="pb-2">
            <CardTitle className="font-mono text-base flex items-center gap-2">
              <FolderGit2 className="h-4 w-4 text-purple-400" /> C# Metadata
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-3 text-sm">
            <p className="text-muted-foreground">
              Parses <code className="text-primary">.csproj</code> XML to extract <code className="text-primary">TargetFramework</code>,
              <code className="text-primary"> OutputType</code> (Exe/Library), and <code className="text-primary">Sdk</code> attribute.
              Finds <code className="text-primary">.sln</code> files and parses <code className="text-primary">global.json</code> for SDK version.
              Collects key files like <code className="text-primary">Directory.Build.props</code> and <code className="text-primary">.editorconfig</code>.
            </p>
            <div className="rounded-lg bg-muted/30 border border-border p-3 font-mono text-xs space-y-1">
              <p className="text-muted-foreground">// Example extracted metadata:</p>
              <p><span className="text-purple-400">Solution:</span> EnterpriseApp.sln</p>
              <p><span className="text-purple-400">SDK:</span> 8.0.100</p>
              <p><span className="text-purple-400">Projects:</span> Api.csproj (net8.0, Exe), Core.csproj (net8.0, Library)</p>
            </div>
          </CardContent>
        </Card>

        <Card className="border-border">
          <CardHeader className="pb-2">
            <CardTitle className="font-mono text-base flex items-center gap-2">
              <Braces className="h-4 w-4 text-green-400" /> Node.js / React Metadata
            </CardTitle>
          </CardHeader>
          <CardContent className="text-sm">
            <p className="text-muted-foreground">
              Reads <code className="text-primary">package.json</code> <code className="text-primary">name</code> field for the project name.
              Checks <code className="text-primary">dependencies</code> and <code className="text-primary">devDependencies</code> for
              React indicators: <code className="text-primary">react</code>, <code className="text-primary">@types/react</code>,
              <code className="text-primary">react-scripts</code>, <code className="text-primary">next</code>,
              <code className="text-primary">gatsby</code>, <code className="text-primary">remix</code>.
              If any match, the project is promoted from Node.js to React (exclusive classification).
            </p>
          </CardContent>
        </Card>

        <Card className="border-border">
          <CardHeader className="pb-2">
            <CardTitle className="font-mono text-base flex items-center gap-2">
              <Cpu className="h-4 w-4 text-orange-400" /> C++ Metadata
            </CardTitle>
          </CardHeader>
          <CardContent className="text-sm">
            <p className="text-muted-foreground">
              Parses <code className="text-primary">CMakeLists.txt</code> <code className="text-primary">project()</code> directive
              to extract the project name. Supports <code className="text-primary">cmake-build-*</code> prefix matching
              for exclusion of generated build directories.
            </p>
          </CardContent>
        </Card>
      </div>
    </motion.div>

    {/* Database Schema */}
    <motion.div {...fade} transition={{ duration: 0.3, delay: 0.25 }}>
      <h2 className="text-xl font-heading font-semibold mt-10 mb-3 flex items-center gap-2">
        <Database className="h-5 w-5 text-primary" />
        Database Schema
      </h2>
      <p className="text-sm text-muted-foreground mb-4">
        Detection results are persisted to SQLite with a normalized schema. Stale records are cleaned during incremental rescans.
      </p>
      <div className="rounded-lg border border-border overflow-hidden overflow-x-auto mb-8">
        <table className="w-full text-sm">
          <thead>
            <tr className="bg-muted/50">
              <th className="text-left font-mono font-semibold px-4 py-2.5">Table</th>
              <th className="text-left font-mono font-semibold px-4 py-2.5 hidden md:table-cell">Key Columns</th>
              <th className="text-left font-mono font-semibold px-4 py-2.5">Purpose</th>
            </tr>
          </thead>
          <tbody>
            {dbTables.map((t) => (
              <tr key={t.name} className="border-t border-border">
                <td className="px-4 py-2.5 font-mono text-primary font-semibold whitespace-nowrap">{t.name}</td>
                <td className="px-4 py-2.5 font-mono text-xs text-muted-foreground hidden md:table-cell">{t.cols}</td>
                <td className="px-4 py-2.5 text-muted-foreground">{t.desc}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* ER Diagram (text) */}
      <div className="rounded-lg border border-border bg-muted/20 p-4 mb-8 font-mono text-xs leading-relaxed">
        <p className="text-muted-foreground mb-2">// Entity Relationships</p>
        <p className="text-foreground">ProjectTypes <span className="text-primary">1:N</span> DetectedProjects</p>
        <p className="text-foreground">DetectedProjects <span className="text-primary">1:1</span> GoProjectMetadata</p>
        <p className="text-foreground">DetectedProjects <span className="text-primary">1:N</span> GoRunnableFiles</p>
        <p className="text-foreground">DetectedProjects <span className="text-primary">1:1</span> CSharpProjectMetadata</p>
        <p className="text-foreground">DetectedProjects <span className="text-primary">1:N</span> CSharpProjectFiles</p>
        <p className="text-foreground">DetectedProjects <span className="text-primary">1:N</span> CSharpKeyFiles</p>
      </div>
    </motion.div>

    {/* Monorepo & C# Precedence */}
    <motion.div {...fade} transition={{ duration: 0.3, delay: 0.3 }}>
      <h2 className="text-xl font-heading font-semibold mt-10 mb-3 flex items-center gap-2">
        <GitBranch className="h-5 w-5 text-primary" />
        Monorepo & Nesting
      </h2>
      <div className="grid md:grid-cols-2 gap-4 mb-8">
        <Card className="border-border">
          <CardContent className="p-4">
            <Search className="h-5 w-5 text-primary mb-2" />
            <h3 className="font-mono font-semibold text-sm mb-1">Monorepo Support</h3>
            <p className="text-xs text-muted-foreground">
              A single repo with <code className="text-primary">backend/</code> (Go) and <code className="text-primary">frontend/</code> (React)
              produces two separate detection records linked to the same repo via <code className="text-primary">RepoId</code>.
            </p>
          </CardContent>
        </Card>
        <Card className="border-border">
          <CardContent className="p-4">
            <FolderGit2 className="h-5 w-5 text-primary mb-2" />
            <h3 className="font-mono font-semibold text-sm mb-1">Nested Projects</h3>
            <p className="text-xs text-muted-foreground">
              A Node.js project at root with React at <code className="text-primary">web/</code> records both.
              The more specific classification wins at each path level.
            </p>
          </CardContent>
        </Card>
        <Card className="border-border md:col-span-2">
          <CardContent className="p-4">
            <Layers className="h-5 w-5 text-primary mb-2" />
            <h3 className="font-mono font-semibold text-sm mb-1">C# Solution Precedence</h3>
            <p className="text-xs text-muted-foreground">
              When a <code className="text-primary">.sln</code> file is found, it defines a single project entry.
              Individual <code className="text-primary">.csproj</code> files beneath it are stored as child records in
              <code className="text-primary"> CSharpProjectFiles</code>, not separate projects.
              Standalone <code className="text-primary">.csproj</code> files (no parent .sln) become their own project entries.
            </p>
          </CardContent>
        </Card>
      </div>
    </motion.div>

    {/* Query Commands */}
    <motion.div {...fade} transition={{ duration: 0.3, delay: 0.35 }}>
      <h2 className="text-xl font-heading font-semibold mt-10 mb-3 flex items-center gap-2">
        <Terminal className="h-5 w-5 text-primary" />
        Query Commands
      </h2>
      <p className="text-sm text-muted-foreground mb-4">
        Each project type has a dedicated command for instant filtered access from the SQLite index:
      </p>
      <div className="grid md:grid-cols-2 gap-3 mb-6">
        {projectTypes.map((p) => (
          <div key={p.key} className="flex items-center gap-3 rounded-lg border border-border bg-card px-4 py-3">
            <p.icon className={`h-4 w-4 ${p.color} shrink-0`} />
            <code className="font-mono text-sm text-primary">{p.commands}</code>
            <span className="text-xs text-muted-foreground">List {p.type} projects</span>
          </div>
        ))}
      </div>

      <h3 className="text-base font-heading font-semibold mt-6 mb-3">Usage Examples</h3>
      <div className="space-y-2 mb-8">
        <CodeBlock code="gitmap go-repos" title="List all Go projects" />
        <CodeBlock code="gitmap go-repos --json" title="Go projects as JSON" />
        <CodeBlock code="gitmap go-repos --count" title="Count Go projects" />
        <CodeBlock code="gitmap react-repos" title="List React projects" />
        <CodeBlock code="gitmap csharp-repos --json" title="C# projects as JSON" />
      </div>
    </motion.div>

    {/* JSON Output */}
    <motion.div {...fade} transition={{ duration: 0.3, delay: 0.4 }}>
      <h2 className="text-xl font-heading font-semibold mt-10 mb-3 flex items-center gap-2">
        <FileJson className="h-5 w-5 text-primary" />
        JSON Output Files
      </h2>
      <p className="text-sm text-muted-foreground mb-4">
        During scan, per-type JSON files are written to the output directory. Empty types are skipped.
        Records are sorted by <code className="text-primary">repoName</code> then <code className="text-primary">relativePath</code>.
      </p>
      <div className="rounded-lg border border-border bg-muted/20 p-4 mb-8 font-mono text-xs leading-relaxed">
        <p className="text-muted-foreground mb-2">.gitmap/output/</p>
        <p className="text-foreground ml-4">├── go-projects.json</p>
        <p className="text-foreground ml-4">├── node-projects.json</p>
        <p className="text-foreground ml-4">├── react-projects.json</p>
        <p className="text-foreground ml-4">├── cpp-projects.json</p>
        <p className="text-foreground ml-4">└── csharp-projects.json</p>
      </div>
    </motion.div>

    {/* Excluded dirs */}
    <motion.div {...fade} transition={{ duration: 0.3, delay: 0.45 }}>
      <h2 className="text-xl font-heading font-semibold mt-10 mb-3 flex items-center gap-2">
        <ShieldCheck className="h-5 w-5 text-primary" />
        Excluded Directories
      </h2>
      <p className="text-sm text-muted-foreground mb-3">
        These directories are skipped during tree traversal to avoid false positives and reduce scan time:
      </p>
      <div className="flex flex-wrap gap-2 mb-8">
        {excludeDirs.map((dir) => (
          <span key={dir} className="font-mono text-xs bg-muted text-muted-foreground px-2.5 py-1 rounded border border-border">
            {dir}
          </span>
        ))}
      </div>
    </motion.div>

    {/* File layout */}
    <motion.div {...fade} transition={{ duration: 0.3, delay: 0.5 }}>
      <h2 className="text-xl font-heading font-semibold mt-10 mb-3 flex items-center gap-2">
        <Layers className="h-5 w-5 text-primary" />
        Package Layout
      </h2>
      <div className="rounded-lg border border-border overflow-hidden overflow-x-auto">
        <table className="w-full text-sm">
          <thead>
            <tr className="bg-muted/50">
              <th className="text-left font-mono font-semibold px-4 py-2.5">File</th>
              <th className="text-left font-mono font-semibold px-4 py-2.5">Purpose</th>
            </tr>
          </thead>
          <tbody>
            {fileLayout.map(([file, purpose]) => (
              <tr key={file} className="border-t border-border">
                <td className="px-4 py-2 font-mono text-primary text-xs whitespace-nowrap">{file}</td>
                <td className="px-4 py-2 text-muted-foreground text-xs">{purpose}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </motion.div>
  </DocsLayout>
);

export default ProjectDetectionPage;
