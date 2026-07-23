import { useState } from "react";
import DocsLayout from "@/components/docs/DocsLayout";
import CodeBlock from "@/components/docs/CodeBlock";
import { ChevronDown, ChevronRight } from "lucide-react";

interface SectionProps {
  id: string;
  title: string;
  children: React.ReactNode;
  defaultOpen?: boolean;
}

const Section = ({ id, title, children, defaultOpen = false }: SectionProps) => {
  const [open, setOpen] = useState(defaultOpen);
  return (
    <div id={id} className="border border-border rounded-lg overflow-hidden">
      <button
        onClick={() => setOpen(!open)}
        className="w-full flex items-center gap-3 px-4 py-3 bg-muted/30 hover:bg-muted/50 transition-colors text-left"
      >
        {open ? <ChevronDown className="h-4 w-4 text-primary shrink-0" /> : <ChevronRight className="h-4 w-4 text-muted-foreground shrink-0" />}
        <span className="font-mono font-semibold text-sm text-foreground">{title}</span>
      </button>
      {open && <div className="px-4 pb-5 pt-3 space-y-4">{children}</div>}
    </div>
  );
};

const P = ({ children }: { children: React.ReactNode }) => (
  <p className="text-sm text-muted-foreground leading-relaxed">{children}</p>
);

const H3 = ({ children }: { children: React.ReactNode }) => (
  <h3 className="text-base font-heading font-semibold text-foreground">{children}</h3>
);

const Table = ({ headers, rows }: { headers: string[]; rows: string[][] }) => (
  <div className="bg-card border border-border rounded-lg overflow-hidden overflow-x-auto">
    <table className="w-full text-sm">
      <thead>
        <tr className="border-b border-border bg-muted/30">
          {headers.map((h) => (
            <th key={h} className="text-left px-4 py-2 font-mono font-semibold text-foreground whitespace-nowrap">{h}</th>
          ))}
        </tr>
      </thead>
      <tbody className="divide-y divide-border">
        {rows.map((row, i) => (
          <tr key={i}>
            {row.map((cell, j) => (
              <td key={j} className={`px-4 py-2 whitespace-nowrap ${j === 0 ? "font-mono text-primary" : "text-muted-foreground"}`}>{cell}</td>
            ))}
          </tr>
        ))}
      </tbody>
    </table>
  </div>
);

const BulletList = ({ items }: { items: string[] }) => (
  <ul className="space-y-1.5 text-sm text-muted-foreground">
    {items.map((item, i) => (
      <li key={i} className="flex gap-2"><span className="text-primary">•</span>{item}</li>
    ))}
  </ul>
);

const sections = [
  { id: "overview", label: "Overview & Philosophy" },
  { id: "structure", label: "Project Structure" },
  { id: "subcommands", label: "Subcommand Architecture" },
  { id: "flags", label: "Flag Parsing" },
  { id: "config", label: "Configuration" },
  { id: "output", label: "Output Formatting" },
  { id: "errors", label: "Error Handling" },
  { id: "style", label: "Code Style" },
  { id: "help", label: "Help System" },
  { id: "database", label: "Database" },
  { id: "build", label: "Build & Deploy" },
  { id: "testing", label: "Testing" },
  { id: "checklist", label: "Implementation Checklist" },
  { id: "dates", label: "Date Formatting" },
  { id: "constants", label: "Constants Reference" },
  { id: "verbose", label: "Verbose Logging" },
  { id: "progress", label: "Progress Tracking" },
  { id: "batch", label: "Batch Execution" },
  { id: "completion", label: "Shell Completion" },
];

const GenericCLIPage = () => {
  const scrollTo = (id: string) => {
    document.getElementById(id)?.scrollIntoView({ behavior: "smooth", block: "start" });
  };

  return (
    <DocsLayout>
      <h1 className="text-3xl font-heading font-bold mb-2 docs-h1">Generic CLI Creation Guidelines</h1>
      <P>A complete, self-contained blueprint for building production-quality CLI tools. Hand this spec to any AI or developer to implement a well-structured CLI from scratch.</P>

      {/* Section jump nav */}
      <div className="flex flex-wrap gap-1.5 my-6">
        {sections.map((s) => (
          <button
            key={s.id}
            onClick={() => scrollTo(s.id)}
            className="text-xs font-mono px-2.5 py-1.5 rounded-md border border-border bg-card text-muted-foreground hover:text-foreground hover:border-primary/40 transition-colors"
          >
            {s.label}
          </button>
        ))}
      </div>

      <div className="space-y-3">
        {/* 01 - Overview */}
        <Section id="overview" title="01 — Overview & Philosophy" defaultOpen>
          <Table
            headers={["Principle", "Detail"]}
            rows={[
              ["Consistency over cleverness", "Predictable patterns across all commands"],
              ["Convention over configuration", "Sensible defaults; config is optional"],
              ["Fail fast, fail clearly", "Bad input → immediate error with actionable message"],
              ["One responsibility per unit", "Each file, function, and package does one thing"],
              ["No magic strings", "Every literal in a constants package"],
              ["Self-documenting", "Help text, version, and examples built into the binary"],
            ]}
          />
        </Section>

        {/* 02 - Project Structure */}
        <Section id="structure" title="02 — Project Structure">
          <CodeBlock code={`toolname/
├── main.go              Entry point — calls cmd.Run()
├── cmd/                 CLI routing, flag parsing, subcommand handlers
│   ├── root.go          Run() + dispatch()
│   ├── rootflags.go     Flag registration helpers
│   ├── rootusage.go     Help/usage printers
│   ├── helpcheck.go     --help interception
│   └── scan.go          One file per subcommand
├── config/              Config file loading + flag merging
├── constants/           All shared string literals and defaults
├── model/               Shared data structures
├── store/               Database init, CRUD operations
├── scanner/             Domain logic — directory walking, detection
├── mapper/              Data transformation (raw → output records)
├── formatter/           Output rendering (terminal, CSV, JSON, scripts)
│   └── templates/       Embedded script templates
├── helptext/            Embedded Markdown help files
│   └── print.go         go:embed + Print function
└── data/                Default config files`} />
          <Table
            headers={["Rule", "Detail"]}
            rows={[
              ["One responsibility per package", "cmd routes, scanner scans, formatter renders"],
              ["No circular imports", "cmd calls others; others never import cmd"],
              ["Leaf packages", "model and constants import nothing project-specific"],
              ["File length", "100–200 lines max per file"],
              ["File naming", "Lowercase, single word or hyphenated"],
            ]}
          />
        </Section>

        {/* 03 - Subcommand Architecture */}
        <Section id="subcommands" title="03 — Subcommand Architecture">
          <H3>Dispatch Pattern</H3>
          <CodeBlock code={`func dispatch(command string) {
    switch command {
    case constants.CmdScan, constants.AliasScan:
        runScan(os.Args[2:])
    case constants.CmdVersion:
        fmt.Println(constants.Version)
    case constants.CmdHelp:
        printUsage()
    default:
        fmt.Fprintf(os.Stderr, "Unknown command: %s\\n", command)
        os.Exit(1)
    }
}`} />
          <H3>Handler Pattern</H3>
          <CodeBlock code={`func runScan(args []string) {
    checkHelp("scan", args)       // 1. Intercept --help
    dir, cfg := parseScanFlags(args) // 2. Parse flags
    records := scanner.Scan(dir, cfg) // 3. Execute logic
    formatter.WriteTerminal(os.Stdout, records) // 4. Output
}`} />
          <Table
            headers={["Rule", "Rationale"]}
            rows={[
              ["One file per subcommand", "Single responsibility"],
              ["Handlers are unexported", "Only Run() is the public API"],
              ["Unknown commands → stderr + exit 1", "Fail fast, fail clearly"],
              ["Each handler ≤ 15 lines", "Extract helpers for complex flows"],
            ]}
          />
        </Section>

        {/* 04 - Flag Parsing */}
        <Section id="flags" title="04 — Flag Parsing">
          <CodeBlock code={`func parseScanFlags(args []string) (dir string, mode string) {
    fs := flag.NewFlagSet("scan", flag.ExitOnError)
    fs.StringVar(&mode, "mode", constants.ModeHTTPS, "Clone URL style")
    fs.Parse(args)
    if fs.NArg() > 0 {
        dir = fs.Arg(0)
    }
    return
}`} />
          <Table
            headers={["Pattern", "Example", "Why"]}
            rows={[
              ["Lowercase with hyphens", "--target-dir", "Readable, standard"],
              ["Boolean flags as switches", "--dry-run", "No value needed"],
              ["Positional args for primary input", "tool scan <dir>", "Natural CLI UX"],
              ["Defaults in constants", "constants.ModeHTTPS", "No inline magic strings"],
            ]}
          />
        </Section>

        {/* 05 - Configuration */}
        <Section id="config" title="05 — Configuration">
          <H3>Three-Layer Config</H3>
          <div className="bg-card border border-border rounded-lg p-4">
            <div className="space-y-2 font-mono text-sm">
              <div className="flex items-center gap-3">
                <span className="bg-muted text-muted-foreground px-3 py-1 rounded">1. Defaults</span>
                <span className="text-muted-foreground">→ Constants in code (lowest priority)</span>
              </div>
              <div className="flex items-center gap-3">
                <span className="bg-primary/10 text-foreground border border-primary/20 dark:bg-primary/15 dark:text-primary px-3 py-1 rounded transition-colors duration-300 hover:border-primary/40 hover:shadow-sm hover:shadow-primary/10">2. Config file</span>
                <span className="text-muted-foreground">→ ./data/config.json</span>
              </div>
              <div className="flex items-center gap-3">
                <span className="bg-primary/20 text-foreground border border-primary/30 dark:bg-primary/20 dark:text-primary px-3 py-1 rounded font-semibold">3. CLI flags</span>
                <span className="text-muted-foreground">→ Always wins (highest priority)</span>
              </div>
            </div>
          </div>
          <BulletList items={[
            "Missing config file → use defaults silently (no error)",
            "Flags always override config file values",
            "Config paths relative to binary unless absolute",
            "JSON schema: flat, camelCase, no nulls (use empty strings/arrays)",
          ]} />
        </Section>

        {/* 06 - Output Formatting */}
        <Section id="output" title="06 — Output Formatting">
          <P>Generate all output formats in one pass — don't make the user choose.</P>
          <Table
            headers={["Format", "Destination", "Purpose"]}
            rows={[
              ["Terminal (colored)", "stdout", "Immediate human feedback"],
              ["CSV", "file", "Spreadsheet / data import"],
              ["JSON", "file", "Machine-readable, re-import"],
              ["Markdown", "file", "Documentation / review"],
              ["Scripts", "file", "Automation / re-execution"],
            ]}
          />
          <H3>Terminal Color Codes</H3>
          <Table
            headers={["Element", "Color", "Purpose"]}
            rows={[
              ["Banner/headers", "Cyan", "Visual identity"],
              ["Success markers (✓)", "Green", "Confirmed items"],
              ["Warnings (⚠)", "Yellow", "Non-fatal issues"],
              ["Data values", "White", "Primary content"],
              ["Metadata", "Dim/Gray", "Secondary info"],
            ]}
          />
        </Section>

        {/* 07 - Error Handling */}
        <Section id="errors" title="07 — Error Handling">
          <Table
            headers={["Exit Code", "Meaning"]}
            rows={[
              ["0", "Success"],
              ["1", "User error (bad args, missing file, invalid input)"],
              ["Non-zero", "Propagated from child processes"],
            ]}
          />
          <BulletList items={[
            "All error format strings in constants package",
            "Errors print to stderr, never stdout",
            "Exit immediately after error — don't continue with bad state",
            "Messages must be actionable — tell the user what to do",
            "Batch operations: log per-item failures, continue, print summary",
          ]} />
        </Section>

        {/* 08 - Code Style */}
        <Section id="style" title="08 — Code Style Rules">
          <Table
            headers={["Constraint", "Rule"]}
            rows={[
              ["if conditions", "Always positive — no !, no !="],
              ["Function length", "8–15 lines (excluding blanks and comments)"],
              ["File length", "100–200 lines max"],
              ["Package granularity", "One responsibility per package"],
              ["Newline before return", "Always, unless return is the only line in an if"],
              ["No magic strings", "All literals in constants package"],
            ]}
          />
          <H3>Naming Conventions</H3>
          <Table
            headers={["Element", "Convention", "Example"]}
            rows={[
              ["Package names", "Lowercase, single word", "scanner, formatter"],
              ["Exported functions", "PascalCase, verb-led", "BuildRecords, WriteCSV"],
              ["Unexported functions", "camelCase, verb-led", "parseFlags, resolveDir"],
              ["Constants", "PascalCase", "DefaultBranch, ModeHTTPS"],
              ["Files", "Lowercase", "terminal.go, csv.go"],
            ]}
          />
        </Section>

        {/* 09 - Help System */}
        <Section id="help" title="09 — Help System">
          <P>Every command supports --help / -h with embedded Markdown files compiled into the binary via go:embed.</P>
          <H3>Help File Template</H3>
          <CodeBlock code={`# toolname <command>

<One-line description>

## Alias
<alias>

## Usage
    toolname <command> [args] [flags]

## Flags
| Flag | Default | Description |
|------|---------|-------------|

## Examples
### Example 1: <title>
    toolname <command> <args>
**Output:**
    <max 3 lines>

## See Also
- [related-command](related-command.md) — description`} />
          <H3>Interception Pattern</H3>
          <CodeBlock code={`func checkHelp(command string, args []string) {
    for _, a := range args {
        if a == "--help" || a == "-h" {
            helptext.Print(command)
            os.Exit(0)
        }
    }
}`} />
        </Section>

        {/* 10 - Database */}
        <Section id="database" title="10 — Database">
          <BulletList items={[
            "CGo-free SQLite driver (e.g., modernc.org/sqlite)",
            "Auto-create database on first data-producing command",
            "Table and column names in PascalCase",
            "Timestamps as TEXT DEFAULT CURRENT_TIMESTAMP",
            "Booleans as INTEGER DEFAULT 0",
            "Upsert strategy: match by unique field, update if exists",
          ]} />
          <H3>Store Package</H3>
          <CodeBlock code={`store/
├── store.go     DB init, open, close, migration, reset
├── repo.go      Item CRUD (upsert, list, find by slug)
├── group.go     Group CRUD
└── history.go   History insert + query`} />
        </Section>

        {/* 11 - Build & Deploy */}
        <Section id="build" title="11 — Build & Deploy">
          <Table
            headers={["Step", "Action", "Skippable"]}
            rows={[
              ["1/4", "Git pull latest source", "-NoPull"],
              ["2/4", "Resolve dependencies", "No"],
              ["3/4", "Compile binary", "No"],
              ["4/4", "Deploy to target directory", "-NoDeploy"],
            ]}
          />
          <H3>Self-Update (Windows-Safe)</H3>
          <BulletList items={[
            "Parent copies itself to temp location",
            "Parent launches temp copy with worker command (blocking)",
            "Worker pulls, builds, deploys new binary",
            "Worker uses rename-first strategy for locked binaries",
            "Always keep .old backup until cleanup runs",
          ]} />
        </Section>

        {/* 12 - Testing */}
        <Section id="testing" title="12 — Testing">
          <Table
            headers={["Layer", "What to Test"]}
            rows={[
              ["mapper", "Data transformation correctness"],
              ["config", "Merge priority (defaults → file → flags)"],
              ["formatter", "Output matches expected format (io.Writer)"],
              ["store", "CRUD operations with in-memory SQLite"],
              ["cmd", "Flag parsing returns correct values"],
              ["scanner", "Detection rules match expected patterns"],
            ]}
          />
          <BulletList items={[
            "Unit tests: same package, same directory as source",
            "Integration tests: under tests/ in separate packages",
            "Table-driven tests for functions with multiple cases",
            "All formatters accept io.Writer for testability",
          ]} />
        </Section>

        {/* 13 - Checklist */}
        <Section id="checklist" title="13 — Implementation Checklist">
          <P>Execute phases in order. Start here when building a new CLI.</P>
          {[
            { phase: "Phase 1: Scaffold", items: ["go mod init", "main.go → cmd.Run()", "constants package (version, CLI names, colors, messages)", "cmd/root.go + rootflags.go + rootusage.go", "version and help commands"] },
            { phase: "Phase 2: Configuration", items: ["model/ package with core structs", "config/config.go with three-layer merge", "data/config.json with defaults"] },
            { phase: "Phase 3: Core Command", items: ["scanner/ or domain logic package", "mapper/ for data transformation", "First real command with flag parsing"] },
            { phase: "Phase 4: Output Formatting", items: ["formatter/ — terminal, CSV, JSON, Markdown, templates", "Output directory structure", "Date formatting utility"] },
            { phase: "Phase 5: Database", items: ["store/ — init, migration, CRUD", "constants_store.go with SQL statements", "db-reset command"] },
            { phase: "Phase 6: Additional Commands", items: ["One file per command", "Flag parsing per command", "Wire into dispatch"] },
            { phase: "Phase 7: Help System", items: ["helptext/print.go with go:embed", "One .md per command", "cmd/helpcheck.go — checkHelp() in every handler"] },
            { phase: "Phase 8: Build & Deploy", items: ["Build script with -ldflags", "Deploy with retry logic", "Self-update command"] },
            { phase: "Phase 9: Testing", items: ["Unit tests for mapper, config, formatter, store", "Integration tests under tests/", "go test ./... passes clean"] },
            { phase: "Phase 10: Polish", items: ["README.md with grouped command reference", "All files ≤ 200 lines, functions ≤ 15 lines", "No magic strings, positive conditionals, blank line before return"] },
          ].map((p) => (
            <div key={p.phase}>
              <h4 className="text-sm font-mono font-semibold text-foreground mb-1">{p.phase}</h4>
              <ul className="space-y-0.5 text-sm text-muted-foreground mb-3">
                {p.items.map((item, i) => (
                  <li key={i} className="flex gap-2"><span className="text-primary/60">☐</span>{item}</li>
                ))}
              </ul>
            </div>
          ))}
        </Section>

        {/* 14 - Date Formatting */}
        <Section id="dates" title="14 — Date Formatting">
          <P>All date/time output passes through a single centralized function. No command formats dates inline.</P>
          <CodeBlock code={`Layout: 02-Jan-2006 03:04 PM

func FormatDisplayDate(t time.Time) string {
    utc := t.UTC()
    local := utc.Local()
    return local.Format(constants.DateDisplayLayout)
}`} />
          <BulletList items={[
            "No time.Format calls in command handlers",
            "Layout constant lives in constants package",
            "UTC → Local conversion inside the function, not at call site",
          ]} />
        </Section>

        {/* 15 - Constants Reference */}
        <Section id="constants" title="15 — Constants Reference">
          <P>The constants package is the single source of truth. Split into focused files ≤ 200 lines each.</P>
          <H3>Naming Prefixes</H3>
          <Table
            headers={["Prefix", "Category", "Example"]}
            rows={[
              ["Cmd", "CLI command names", 'CmdScan = "scan"'],
              ["Mode", "Operation modes", 'ModeHTTPS = "https"'],
              ["Output", "Output formats", 'OutputJSON = "json"'],
              ["Ext", "File extensions", 'ExtCSV = ".csv"'],
              ["Default", "Default values", 'DefaultBranch = "main"'],
              ["Color", "ANSI codes", 'ColorGreen = "\\033[32m"'],
              ["Err", "Error messages", 'ErrSourceRequired = "Error: ..."'],
              ["Msg", "User messages", 'MsgScanComplete = "✓ Scan complete"'],
              ["Git", "Git commands/flags", 'GitClone = "clone"'],
              ["SQL", "SQL statements", "SQLCreateRepos = `CREATE TABLE...`"],
              ["Table", "Table names", 'TableRepos = "Repos"'],
              ["DB", "Database paths", 'DBFile = "toolname.db"'],
              ["Flag", "Flag names", 'FlagVerbose = "verbose"'],
              ["Status", "UI indicators", 'StatusIconClean = "✓ clean"'],
              ["Date", "Date formatting", 'DateDisplayLayout = "02-Jan-..."'],
              ["Tree", "Tree-drawing chars", 'TreeBranch = "├──"'],
            ]}
          />
          <H3>What Does NOT Belong</H3>
          <BulletList items={[
            "Struct definitions → model package",
            "Business logic → domain packages",
            "Template content → go:embed in formatter/templates/",
            "Test data strings → local in test files",
          ]} />
        </Section>

        {/* 16 - Verbose Logging */}
        <Section id="verbose" title="16 — Verbose Logging">
          <P>A shared --verbose flag enables detailed debug logging to a timestamped file. Normal runs produce clean user-facing output only. Verbose runs capture full diagnostics without polluting stdout.</P>
          <Table
            headers={["Rule", "Detail"]}
            rows={[
              ["Off by default", "No log file created unless --verbose is passed"],
              ["File + stderr", "Every verbose entry writes to both the log file and stderr"],
              ["Timestamped entries", "Each line prefixed with [HH:MM:SS.mmm]"],
              ["Timestamped filenames", "Log file named toolname-verbose-YYYY-MM-DD_HH-mm-ss.log"],
              ["Output directory", "Logs written to the tool's default output folder"],
              ["Dim on stderr", "Verbose stderr output uses dim/gray ANSI color"],
              ["No stdout pollution", "Verbose output never mixes with normal command output"],
              ["Global singleton", "One logger instance shared across all packages"],
            ]}
          />
          <H3>Package Structure</H3>
          <CodeBlock code={`verbose/
└── verbose.go     Logger type, Init, Close, Log, IsEnabled, Get`} />
          <H3>Logger API</H3>
          <CodeBlock code={`// Init creates the log file and enables verbose logging.
func Init() (*Logger, error)

// Close flushes and closes the log file.
func (l *Logger) Close()

// Log writes a formatted message to the log file and stderr.
func (l *Logger) Log(format string, args ...interface{})

// IsEnabled returns true if verbose mode is active.
func IsEnabled() bool

// Get returns the global logger (may be nil).
func Get() *Logger`} />
          <H3>Logger Type</H3>
          <CodeBlock code={`type Logger struct {
    file    *os.File
    enabled bool
}

var global *Logger`} />
          <H3>Init Flow</H3>
          <CodeBlock code={`func Init() (*Logger, error) {
    logDir := constants.DefaultOutputFolder
    _ = os.MkdirAll(logDir, constants.DirPermission)

    timestamp := time.Now().Format("2006-01-02_15-04-05")
    logPath := filepath.Join(logDir, fmt.Sprintf(constants.VerboseLogFileFmt, timestamp))

    file, err := os.Create(logPath)
    if err != nil {
        return nil, err
    }

    l := &Logger{file: file, enabled: true}
    global = l
    fmt.Printf(constants.MsgVerboseLogFile, logPath)

    return l, nil
}`} />
          <H3>Log Entry Format</H3>
          <CodeBlock code={`func writeLogEntry(l *Logger, format string, args ...interface{}) {
    line := fmt.Sprintf(format, args...)
    ts := time.Now().Format("15:04:05.000")
    entry := fmt.Sprintf("[%s] %s\\n", ts, line)
    l.file.WriteString(entry)
    fmt.Fprint(os.Stderr, constants.ColorDim+entry+constants.ColorReset)
}

// Example output:
// [14:32:07.123] git clone https://github.com/user/repo.git
// [14:32:09.456] clone completed in 2.3s
// [14:32:09.460] retry attempt 1/4 for locked file`} />
          <H3>Command Handler Pattern</H3>
          <CodeBlock code={`func runPull(args []string) {
    checkHelp("pull", args)
    slug, group, all, verboseFlag := parsePullFlags(args)

    if verboseFlag {
        initVerboseLog()       // Init + defer Close
    }
    // ... command logic ...
}

func initVerboseLog() {
    logger, err := verbose.Init()
    if err != nil {
        fmt.Fprintf(os.Stderr, constants.ErrVerboseInit, err)
        return                 // Non-fatal — continue without logging
    }
    defer logger.Close()
}`} />
          <BulletList items={[
            "Verbose init failure is non-fatal — warn and continue",
            "defer Close() in the same function that calls Init()",
            "Never pass the logger as a parameter — use verbose.Get() or verbose.IsEnabled()",
          ]} />
          <H3>What to Log</H3>
          <Table
            headers={["Category", "Examples"]}
            rows={[
              ["Git operations", "Clone/pull commands, remote URLs, branch names"],
              ["Retry attempts", "Attempt number, delay, reason for retry"],
              ["File I/O", "Paths read/written, file sizes, permissions"],
              ["External processes", "Command lines, exit codes, stdout/stderr"],
              ["Timing", "Operation durations, elapsed time"],
              ["Environment", "OS, paths, config values loaded"],
              ["Errors (detailed)", "Full error chains, stack context"],
            ]}
          />
          <H3>What NOT to Log</H3>
          <BulletList items={[
            "Secrets, tokens, or credentials",
            "Routine success paths that add no diagnostic value",
            "Data that duplicates normal stdout output",
          ]} />
          <H3>Constants</H3>
          <CodeBlock code={`// constants/constants.go
const VerboseLogFileFmt = "toolname-verbose-%s.log"

// constants/constants_cli.go
const FlagVerbose    = "verbose"
const FlagDescVerbose = "Enable verbose debug logging to file"

// constants/constants_messages.go
const MsgVerboseLogFile = "Verbose log: %s\\n"
const ErrVerboseInit    = "Warning: could not initialize verbose log: %v\\n"`} />
          <H3>Conditional Logging in Libraries</H3>
          <CodeBlock code={`func safePullOne(repo model.Record) error {
    logger := verbose.Get()

    if logger != nil {
        logger.Log("pulling %s at %s", repo.Name, repo.Path)
    }

    // ... pull logic ...

    if logger != nil {
        logger.Log("pull complete for %s (%.1fs)", repo.Name, elapsed.Seconds())
    }

    return nil
}`} />
          <BulletList items={[
            "Always nil-check verbose.Get() — verbose may not be active",
            "Keep log calls outside hot loops to avoid performance overhead",
            "Use fmt.Sprintf-style formatting — no structured logging libraries",
          ]} />
        </Section>

        {/* 17 - Progress Tracking */}
        <Section id="progress" title="17 — Progress Tracking">
          <P>Batch operations (clone, pull, exec) process many items sequentially. A progress tracker provides real-time visual feedback so the user knows what is happening, how far along the operation is, and how it concluded.</P>
          <Table
            headers={["Rule", "Detail"]}
            rows={[
              ["Counter format", "[current/total] prefix on every item line"],
              ["Repo name visible", "Current repository name printed alongside the counter"],
              ["Elapsed time", "Duration shown on completion of each item and in the summary"],
              ["Summary at end", "Final line reports total items, elapsed time, success/failure counts"],
              ["Quiet mode", "--quiet flag suppresses all progress output for programmatic use"],
              ["Stderr only", "All progress output goes to stderr — stdout is reserved for data"],
              ["No progress bars", "Counter + name is sufficient; avoid complex progress bar libraries"],
            ]}
          />
          <H3>Package Placement</H3>
          <CodeBlock code={`cloner/
├── cloner.go        Core clone/pull logic
├── progress.go      Progress tracker type
├── safe_pull.go     Retry-aware pull
└── pulldiag.go      Pull diagnostics`} />
          <P>Progress tracking lives in the domain package that owns the batch operation. If multiple packages need it, extract a shared progress/ package.</P>
          <H3>Progress Type</H3>
          <CodeBlock code={`type Progress struct {
    total   int
    current int
    start   time.Time
    quiet   bool
    cloned  int
    pulled  int
    failed  int
}

var global *Logger`} />
          <Table
            headers={["Field", "Purpose"]}
            rows={[
              ["total", "Total number of items to process"],
              ["current", "Running counter, incremented on each Begin() call"],
              ["start", "Timestamp captured at construction for elapsed time"],
              ["quiet", "When true, all output methods become no-ops"],
              ["cloned", "Success counter for new clones"],
              ["pulled", "Success counter for updates (pulls)"],
              ["failed", "Failure counter"],
            ]}
          />
          <H3>Constructor</H3>
          <CodeBlock code={`func NewProgress(total int, quiet bool) *Progress {
    return &Progress{
        total: total,
        start: time.Now(),
        quiet: quiet,
    }
}`} />
          <BulletList items={[
            "time.Now() captured once at construction — not per item",
            "All counters start at zero (Go zero values)",
            "Caller determines quiet from the --quiet CLI flag",
          ]} />
          <H3>Lifecycle Methods</H3>
          <CodeBlock code={`// Begin — called before processing each item
func (p *Progress) Begin(name string) {
    p.current++
    if p.quiet {
        return
    }
    fmt.Fprintf(os.Stderr, constants.ProgressBeginFmt, p.current, p.total, name)
}

// Done — called after an item succeeds
func (p *Progress) Done(result model.CloneResult, pulled bool) {
    if pulled {
        p.pulled++
    } else {
        p.cloned++
    }
    if p.quiet {
        return
    }
    elapsed := time.Since(p.start)
    fmt.Fprintf(os.Stderr, constants.ProgressDoneFmt, formatDuration(elapsed))
}

// Fail — called after an item fails
func (p *Progress) Fail(result model.CloneResult) {
    p.failed++
    if p.quiet {
        return
    }
    fmt.Fprintf(os.Stderr, constants.ProgressFailFmt)
}

// PrintSummary — called once after all items are processed
func (p *Progress) PrintSummary() {
    if p.quiet {
        return
    }
    elapsed := time.Since(p.start)
    fmt.Fprintf(os.Stderr, constants.ProgressSummaryFmt,
        p.current, p.total, formatDuration(elapsed))
    fmt.Fprintf(os.Stderr, constants.ProgressDetailFmt,
        p.cloned, p.pulled, p.failed)
}`} />
          <H3>Example Output</H3>
          <CodeBlock code={`[1/24] repo-name ✓ (2.3s)
[2/24] another-repo ✓ (1.1s)
[3/24] broken-repo ✗ failed
...
Done: 24/24 (1m 12s)
  Cloned: 18  Pulled: 5  Failed: 1`} />
          <H3>Duration Formatting</H3>
          <CodeBlock code={`func formatDuration(d time.Duration) string {
    if d < time.Minute {
        return fmt.Sprintf("%.1fs", d.Seconds())
    }
    mins := int(d.Minutes())
    secs := int(d.Seconds()) % 60
    return fmt.Sprintf("%dm %ds", mins, secs)
}`} />
          <Table
            headers={["Duration", "Output"]}
            rows={[
              ["2.3 seconds", "2.3s"],
              ["72 seconds", "1m 12s"],
              ["5 minutes 3 seconds", "5m 3s"],
            ]}
          />
          <H3>Usage Pattern</H3>
          <CodeBlock code={`func cloneAll(repos []model.Record, quiet bool) {
    progress := NewProgress(len(repos), quiet)

    for _, repo := range repos {
        progress.Begin(repo.Name)

        result, err := cloneOne(repo)
        if err != nil {
            progress.Fail(result)
            continue
        }

        progress.Done(result, repo.Exists)
    }

    progress.PrintSummary()
}`} />
          <BulletList items={[
            "Create progress tracker before the loop",
            "Call Begin() first in each iteration",
            "Call exactly one of Done() or Fail() per item — never both",
            "Call PrintSummary() after the loop, unconditionally",
            "Use continue after Fail() — do not abort the batch",
          ]} />
          <H3>Quiet Mode</H3>
          <CodeBlock code={`// In cmd/clone.go
fs.BoolVar(&quietFlag, constants.FlagQuiet, false, constants.FlagDescQuiet)

// Pass to progress tracker
progress := cloner.NewProgress(len(repos), quietFlag)`} />
          <BulletList items={[
            "All Begin, Done, Fail, and PrintSummary calls become no-ops",
            "Internal counters still increment (for programmatic access)",
            "Stdout data output (JSON, CSV) is unaffected",
          ]} />
          <H3>Constants</H3>
          <CodeBlock code={`// constants/constants_clone.go
const ProgressBeginFmt   = "[%d/%d] %s "
const ProgressDoneFmt    = "✓ (%s)\\n"
const ProgressFailFmt    = "✗ failed\\n"
const ProgressSummaryFmt = "Done: %d/%d (%s)\\n"
const ProgressDetailFmt  = "  Cloned: %d  Pulled: %d  Failed: %d\\n"

// constants/constants_cli.go
const FlagQuiet     = "quiet"
const FlagDescQuiet = "Suppress progress output"`} />
        </Section>

        {/* 18 - Batch Execution */}
        <Section id="batch" title="18 — Batch Execution">
          <P>The exec command runs arbitrary Git commands across multiple repositories in a single invocation, providing consistent output formatting and a summary of success/failure across the entire batch.</P>
          <Table
            headers={["Rule", "Detail"]}
            rows={[
              ["Arbitrary git commands", "Any valid git subcommand and arguments are accepted"],
              ["Two-tier repo lookup", "Database first, JSON file fallback"],
              ["Scope filtering", "--all for every repo, --group for a named subset"],
              ["Per-repo isolation", "Each command runs in the repo's directory independently"],
              ["Continue on failure", "A failed repo does not abort the batch"],
              ["Captured output", "stdout and stderr from each repo are captured and displayed"],
              ["Summary at end", "Final line reports succeeded, failed, and missing counts"],
              ["Banner header", "Visual header showing the command and repo count"],
            ]}
          />
          <H3>Command Signature</H3>
          <CodeBlock code={`toolname exec [--all | --group <name>] <git-args...>

# Examples:
toolname exec status                        # All repos from JSON
toolname exec --all fetch --prune           # All DB repos
toolname exec -g backend log --oneline -5   # Group subset`} />
          <Table
            headers={["Flag", "Short", "Default", "Description"]}
            rows={[
              ["--all", "—", "false", "Run across all repos in the database"],
              ["--group", "-g", '""', "Run across repos in a named group"],
            ]}
          />
          <H3>Flag Parsing</H3>
          <CodeBlock code={`func parseExecFlags(args []string) (groupName string, all bool, gitArgs []string) {
    fs := flag.NewFlagSet(constants.CmdExec, flag.ExitOnError)
    gFlag := fs.String("group", "", constants.FlagDescGroup)
    fs.StringVar(gFlag, "g", "", constants.FlagDescGroup)
    aFlag := fs.Bool("all", false, constants.FlagDescAll)
    fs.Parse(args)

    return *gFlag, *aFlag, fs.Args()
}`} />
          <P>After fs.Parse(args), fs.Args() returns everything not consumed as a flag — these become the git arguments.</P>
          <H3>Repo Loading Strategy</H3>
          <CodeBlock code={`func loadExecByScope(groupName string, all bool) []model.ScanRecord {
    if len(groupName) > 0 {
        return loadRecordsByGroup(groupName)  // DB: group members
    }
    if all {
        return loadAllRecordsDB()             // DB: all repos
    }

    return loadExecRecordsJSON()              // Fallback: JSON file
}`} />
          <Table
            headers={["Priority", "Source", "When"]}
            rows={[
              ["1", "Database (group)", "--group flag provided"],
              ["2", "Database (all)", "--all flag provided"],
              ["3", "JSON file", "Neither flag — use last scan output"],
            ]}
          />
          <H3>Execution Flow</H3>
          <CodeBlock code={`func runExec(args []string) {
    checkHelp("exec", args)
    groupName, all, gitArgs := parseExecFlags(args)
    if len(gitArgs) == 0 {
        fmt.Fprintln(os.Stderr, constants.ErrExecUsage)
        os.Exit(1)
    }

    records := loadExecByScope(groupName, all)
    printExecBanner(gitArgs, len(records))
    succeeded, failed, missing := execAllRepos(records, gitArgs)
    printExecSummary(succeeded, failed, missing, len(records))
}`} />
          <H3>Per-Repo Execution</H3>
          <CodeBlock code={`func execInRepo(rec model.ScanRecord, gitArgs []string) bool {
    cmd := exec.Command(constants.GitBin, gitArgs...)
    cmd.Dir = rec.AbsolutePath
    cmd.Stdout = nil
    cmd.Stderr = nil

    out, err := cmd.CombinedOutput()
    output := strings.TrimSpace(string(out))
    printExecResult(rec.RepoName, output, err)

    return err == nil
}`} />
          <BulletList items={[
            "cmd.Dir sets the working directory — no cd needed",
            "CombinedOutput() captures both stdout and stderr",
            "Trim whitespace from output before display",
            "Return success/failure as boolean — caller tracks counts",
          ]} />
          <H3>Three Outcomes Per Repo</H3>
          <CodeBlock code={`func execOneRepo(rec model.ScanRecord, gitArgs []string) (int, int, int) {
    _, err := os.Stat(rec.AbsolutePath)
    if err == nil && execInRepo(rec, gitArgs) {
        return 1, 0, 0    // succeeded
    }
    if err == nil {
        return 0, 1, 0    // failed (command error)
    }

    // Directory does not exist
    fmt.Printf(constants.ExecMissingFmt, ...)
    return 0, 0, 1        // missing
}`} />
          <H3>Output Format</H3>
          <CodeBlock code={`  ╭──────────────────────╮
  │   Batch Execution    │
  ╰──────────────────────╯

  Command: status
  Repos:   24
  ─────────────────────────

  ✓ repo-name              # Success (green)
      On branch main        # Indented captured output (dim)
  ✗ broken-repo            # Failure (yellow)
      fatal: not a repo     # Error output (dim)
  ⊘ deleted-repo  missing  # Directory not found (dim + yellow)

  ─────────────────────────
  24 repos · 22 succeeded · 1 failed · 1 missing`} />
          <BulletList items={[
            "Repo names truncated to 22 characters for alignment",
            "Output lines indented with 6 spaces",
            "Empty output → no indented lines printed",
            "Colors: green for success, yellow for failure/missing, dim for output",
            "Only non-zero categories included in the summary line",
          ]} />
          <H3>Constants</H3>
          <CodeBlock code={`// constants/constants_cli.go
const CmdExec  = "exec"
const AliasExec = "x"

// constants/constants_messages.go
const ErrExecUsage      = "Error: exec requires at least one git argument"
const ErrExecLoadFailed  = "Error: could not load repos: %v\\n"

// constants/constants_terminal.go
const ExecBannerTop     = "╭──────────────────────╮"
const ExecBannerTitle   = "│   Batch Execution    │"
const ExecBannerBottom  = "╰──────────────────────╯"
const ExecCommandFmt    = "Command: %s"
const ExecRepoCountFmt  = "Repos:   %d"
const ExecSuccessFmt    = "  %s✓ %-22s%s\\n"
const ExecFailFmt       = "  %s✗ %-22s%s\\n"
const ExecMissingFmt    = "  %s%-22s %smissing%s\\n"
const ExecOutputLineFmt = "      %s%s%s\\n"
const SummaryJoinSep    = " · "`} />
          <H3>File Organization</H3>
          <CodeBlock code={`cmd/
└── exec.go          Handler, flag parsing, execution, output`} />
          <P>All exec logic fits in a single file. If it exceeds 200 lines, extract execformat.go for output functions.</P>
        </Section>

        {/* Section 19: Shell Completion */}
        <Section id="completion" title="19 — Shell Completion">
          <P>CLI tools should provide tab-completion for all major shells. Completions cover subcommand names, dynamic values (repo slugs, group names), and per-command flags.</P>

          <H3>Design Rules</H3>
          <Table headers={["Rule", "Detail"]} rows={[
            ["Dynamic data", "Completions invoke hidden --list-* flags that query the DB and print one value per line"],
            ["No decoration", "List output has no color, no banners, no headers — just raw values"],
            ["Three shells", "PowerShell (Register-ArgumentCompleter), Bash (complete -F), Zsh (#compdef)"],
            ["Idempotent install", "Setup writes source line to profile only if absent"],
            ["Fast execution", "List flags open DB, query, print, exit — no formatting overhead"],
          ]} />

          <H3>Completion Subcommand</H3>
          <CodeBlock code={`toolname completion <powershell|bash|zsh>   # print script to stdout
toolname completion --list-repos             # one slug per line
toolname completion --list-groups            # one group name per line
toolname completion --list-commands          # one command per line`} />
          <P>The completion subcommand serves two purposes: printing the full script for manual sourcing, and providing machine-readable lists for the scripts to call at tab-time.</P>

          <H3>What Gets Completed</H3>
          <Table headers={["Context", "Completed Values"]} rows={[
            ["toolname <tab>", "All subcommand names and aliases"],
            ["toolname cd <tab>", "Repo slugs + cd subcommands (repos, set-default, clear-default)"],
            ["toolname pull <tab>", "Repo slugs from database"],
            ["toolname cd repos --group <tab>", "Group names from database"],
            ["toolname exec --group <tab>", "Group names from database"],
            ["toolname group <tab>", "Subcommands: create, add, remove, list, show, delete"],
          ]} />

          <H3>PowerShell Script Pattern</H3>
          <CodeBlock code={`Register-ArgumentCompleter -CommandName toolname -ScriptBlock {
    param($wordToComplete, $commandAst, $cursorPosition)
    $elems = $commandAst.CommandElements | Select-Object -Skip 1
    $cmd = if ($elems.Count -gt 0) { $elems[0].ToString() } else { "" }

    if ($cmd -eq "cd" -or $cmd -eq "go") {
        $items = @(toolname completion --list-repos) + @("repos", "set-default")
        $items | Where-Object { $_ -like "$wordToComplete*" } |
            ForEach-Object { [System.Management.Automation.CompletionResult]::new($_) }
        return
    }

    toolname completion --list-commands | Where-Object { $_ -like "$wordToComplete*" } |
        ForEach-Object { [System.Management.Automation.CompletionResult]::new($_) }
}`} />

          <H3>Bash Script Pattern</H3>
          <CodeBlock code={`_toolname_completions() {
    local cur prev cmd
    cur="\${COMP_WORDS[COMP_CWORD]}"
    cmd="\${COMP_WORDS[1]}"

    if [[ \${COMP_CWORD} -eq 1 ]]; then
        COMPREPLY=($(compgen -W "$(toolname completion --list-commands)" -- "$cur"))
        return
    fi

    case "$cmd" in
        cd|go)
            COMPREPLY=($(compgen -W "$(toolname completion --list-repos) repos" -- "$cur"))
            ;;
        pull)
            COMPREPLY=($(compgen -W "$(toolname completion --list-repos)" -- "$cur"))
            ;;
    esac
}
complete -F _toolname_completions toolname`} />

          <H3>Setup Integration</H3>
          <Table headers={["Shell", "Script Location", "Profile"]} rows={[
            ["PowerShell", "$env:APPDATA/toolname/completions.ps1", "$PROFILE"],
            ["Bash", "~/.local/share/toolname/completions.bash", "~/.bashrc"],
            ["Zsh", "~/.local/share/toolname/completions.zsh", "~/.zshrc"],
          ]} />
          <P>The setup command detects the current shell, generates the completion script, writes it to the appropriate location, and appends a source line to the user's profile. Re-running setup does not duplicate the source line.</P>

          <H3>Install Logic</H3>
          <CodeBlock code={`func Install(shell string) error {
    script, err := Generate(shell)
    if err != nil {
        return err
    }

    scriptPath, profilePath := resolvePaths(shell)

    return writeAndSource(script, scriptPath, profilePath, shell)
}

func addSourceLine(scriptPath, profilePath, shell string) error {
    sourceLine := buildSourceLine(scriptPath, shell)
    existing, _ := os.ReadFile(profilePath)
    if strings.Contains(string(existing), sourceLine) {
        return nil  // already installed
    }

    f, _ := os.OpenFile(profilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    defer f.Close()
    fmt.Fprintf(f, "\\n# toolname shell completion\\n%s\\n", sourceLine)
    return nil
}`} />

          <H3>Constants</H3>
          <CodeBlock code={`// Shells
const ShellPowerShell = "powershell"
const ShellBash       = "bash"
const ShellZsh        = "zsh"

// List flags
const CompListRepos    = "--list-repos"
const CompListGroups   = "--list-groups"
const CompListCommands = "--list-commands"

// File names
const CompFilePS   = "completions.ps1"
const CompFileBash = "completions.bash"
const CompFileZsh  = "completions.zsh"`} />

          <H3>File Organization</H3>
          <CodeBlock code={`cmd/
  completion.go        Subcommand handler + list printers
completion/
  completion.go        Generate() + AllCommands()
  powershell.go        PowerShell script generator
  bash.go              Bash script generator
  zsh.go               Zsh script generator
  install.go           Profile detection + source-line insertion`} />
          <P>Each shell generator is its own file. The install logic is separate from generation. All files stay under 200 lines.</P>
        </Section>
      </div>
    </DocsLayout>
  );
};

export default GenericCLIPage;
