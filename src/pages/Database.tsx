import { useState, useMemo } from "react";
import DocsLayout from "@/components/docs/DocsLayout";
import {
  Database,
  Table as TableIcon,
  Search,
  Copy,
  Check,
  X,
  Code2,
  Binary,
  FileText,
  RefreshCw,
  ChevronLeft,
  ChevronRight,
} from "lucide-react";
import { copyToClipboard } from "@/lib/clipboard";
import { useToast } from "@/hooks/use-toast";

interface ClassifiedCell {
  cellType: "null" | "bool" | "blob" | "json" | "multiline" | "text" | "plain";
  badgeText: string;
  rawValue: unknown;
  sizeBytes: number;
}

interface InspectTarget {
  columnName: string;
  classified: ClassifiedCell;
}

const SAMPLE_DATABASES: Record<string, Record<string, Array<Record<string, unknown>>>> = {
  master: {
    Repositories: [
      {
        Id: 1,
        Name: "gitmap-v28",
        Path: "d:\\work\\gitmap",
        Branch: "main",
        IsClean: true,
        Config: { autocommit: true, splitDb: true, maxWorkers: 4 },
        MetaBlob: "\x00\x01\x02\x03\x04\x05GITMAP_BINARY_HDR\x00\x00\x01\xAA\xBB\xCC\xDD",
        CommitLog: "feat: add thread-safe table metadata cache\nfix: table resolution without keyword\nperf: optimize cell inspect",
      },
      {
        Id: 2,
        Name: "kernel-tools",
        Path: "d:\\work\\kernel-tools",
        Branch: "stable",
        IsClean: false,
        Config: { autocommit: false, splitDb: false },
        MetaBlob: "\x1F\x8B\x08\x00\x00\x00\x00\x00\x02\xFF\xCB\x48\xCD\xC9\xC9\x57\x28\xCF",
        CommitLog: "chore: update build script\nfix: memory leak in runner",
      },
      {
        Id: 3,
        Name: "slides-app",
        Path: "d:\\work\\slides-app",
        Branch: "dev",
        IsClean: true,
        Config: { theme: "dark", autoSave: true, fps: 60, canvasWidth: 1920 },
        MetaBlob: "PNG\x0D\x0A\x1A\x0A\x00\x00\x00\x0DIHDR\x00\x00\x04\x00\x00\x00\x02",
        CommitLog: "feat: add 16:9 canvas layout mode\nstyle: update typography ramp",
      },
    ],
    DatabaseAlias: [
      { DatabaseAliasId: 1, AliasName: "master", DatabasePath: ".gitmap/gitmap.db", DatabaseType: "system", IsActive: true, FileSizeBytes: 1048576 },
      { DatabaseAliasId: 2, AliasName: "pipeline", DatabasePath: ".gitmap/pipeline.db", DatabaseType: "system", IsActive: true, FileSizeBytes: 524288 },
      { DatabaseAliasId: 3, AliasName: "tasks", DatabasePath: ".gitmap/tasks.db", DatabaseType: "system", IsActive: true, FileSizeBytes: 262144 },
      { DatabaseAliasId: 4, AliasName: "heatmap", DatabasePath: ".gitmap/heatmap.db", DatabaseType: "system", IsActive: true, FileSizeBytes: 131072 },
    ],
  },
  pipeline: {
    PipelineRun: [
      {
        RunId: "run_8841",
        WorkflowName: "CI Quality Gates",
        Status: "SUCCESS",
        IsActive: false,
        Payload: { event: "push", branch: "main", commit: "a7f3b89", checks: ["lint", "tests", "security"] },
        Diagnostics: "\x08\x96\x01\x12\x07SUCCESS\x1A\x12diagnostic_summary\x20\x01\x2A\x04PASS",
        Logs: "Phase 1: Environment bootstrap\nPhase 2: Linter pass verified\nPhase 3: E2E matrix completed in 1.4s\nExit code: 0",
      },
      {
        RunId: "run_8842",
        WorkflowName: "Split-DB Sync",
        Status: "IN_PROGRESS",
        IsActive: true,
        Payload: { targetDb: "master", strategy: "wal_checkpoint", timeoutSec: 10 },
        Diagnostics: "\x08\x97\x01\x12\x04BUSY\x1A\x0Echeckpoint_log\x20\x02",
        Logs: "Checkpointing WAL journal...\nFlushing pages to disk...\nAwaiting lock release...",
      },
    ],
  },
  tasks: {
    Subtask: [
      {
        SubtaskId: 101,
        Title: "Implement TableMetadataCache",
        Status: "DONE",
        Priority: "HIGH",
        IsCompleted: true,
        Checkpoints: { ttlSeconds: 8, threadSafe: true, methods: ["Get", "Set", "Invalidate", "InvalidateAll"] },
        Notes: "Thread-safe RWMutex in-memory store.\nCache TTL set to 8 seconds.\nZero nested conditionals.",
      },
      {
        SubtaskId: 102,
        Title: "Modernize Database Visualizer",
        Status: "DONE",
        Priority: "HIGH",
        IsCompleted: true,
        Checkpoints: { hexDump: true, formattedJson: true, decodedString: true, badges: ["BLOB", "JSON", "↵"] },
        Notes: "Replaces raw binary dumps with interactive classification badges.\nAdds 3-tab inspector drawer.",
      },
    ],
  },
};

export default function DatabasePage() {
  const { toast } = useToast();
  const [activeDb, setActiveDb] = useState<string>("master");
  const [activeTable, setActiveTable] = useState<string>("Repositories");
  const [filterQuery, setFilterQuery] = useState<string>("");
  const [inspectTarget, setInspectTarget] = useState<InspectTarget | null>(null);
  const [hasCopied, setHasCopied] = useState<boolean>(false);

  const availableTables = useMemo(() => {
    return Object.keys(SAMPLE_DATABASES[activeDb] || {});
  }, [activeDb]);

  const rawRows = useMemo(() => {
    const dbTables = SAMPLE_DATABASES[activeDb] || {};

    return dbTables[activeTable] || [];
  }, [activeDb, activeTable]);

  const filteredRows = useMemo(() => {
    if (!filterQuery.trim()) {
      return rawRows;
    }

    const q = filterQuery.toLowerCase();

    return rawRows.filter((row) =>
      Object.values(row).some((val) => String(val).toLowerCase().includes(q))
    );
  }, [rawRows, filterQuery]);

  const handleSelectDb = (db: string) => {
    setActiveDb(db);
    const tables = Object.keys(SAMPLE_DATABASES[db] || {});
    if (tables.length > 0) {
      setActiveTable(tables[0]);
    }
  };

  const handleOpenInspect = (columnName: string, classified: ClassifiedCell) => {
    setInspectTarget({ columnName, classified });
    setHasCopied(false);
  };

  const handleCloseInspect = () => {
    setInspectTarget(null);
  };

  const handleCopyRaw = async (text: string) => {
    await copyToClipboard(text);
    setHasCopied(true);
    toast({ title: "Copied!", description: "Content copied to clipboard." });
    setTimeout(() => setHasCopied(false), 2000);
  };

  return (
    <DocsLayout>
      <div className="max-w-7xl mx-auto space-y-6 pb-16">
        <DatabaseHeader />
        <DatabaseToolbar
          activeDb={activeDb}
          activeTable={activeTable}
          availableTables={availableTables}
          filterQuery={filterQuery}
          onSelectDb={handleSelectDb}
          onSelectTable={setActiveTable}
          onChangeFilter={setFilterQuery}
        />
        <DatabaseGrid rows={filteredRows} onInspect={handleOpenInspect} />
        {inspectTarget && (
          <CellInspectionModal
            target={inspectTarget}
            hasCopied={hasCopied}
            onClose={handleCloseInspect}
            onCopyRaw={handleCopyRaw}
          />
        )}
      </div>
    </DocsLayout>
  );
}

function DatabaseHeader() {
  return (
    <div className="flex flex-wrap items-center justify-between gap-4 border-b border-border/40 pb-5">
      <div className="flex items-center gap-3">
        <div className="h-10 w-10 rounded-lg bg-primary/10 flex items-center justify-center">
          <Database className="h-5 w-5 text-primary" />
        </div>
        <div>
          <h1 className="text-2xl font-heading font-bold tracking-[-0.02em] text-foreground">Database Visualizer</h1>
          <p className="text-sm text-muted-foreground">
            Interactive table browser, schema inspector, and binary/JSON payload decoder.
          </p>
        </div>
      </div>
      <div className="flex items-center gap-2">
        <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-semibold bg-emerald-500/10 text-emerald-700 dark:text-emerald-400 border border-emerald-500/20">
          Connected
        </span>
      </div>
    </div>
  );
}

interface ToolbarProps {
  activeDb: string;
  activeTable: string;
  availableTables: string[];
  filterQuery: string;
  onSelectDb: (db: string) => void;
  onSelectTable: (table: string) => void;
  onChangeFilter: (val: string) => void;
}

function DatabaseToolbar({
  activeDb,
  activeTable,
  availableTables,
  filterQuery,
  onSelectDb,
  onSelectTable,
  onChangeFilter,
}: ToolbarProps) {
  const dbs = ["master", "pipeline", "tasks"];

  return (
    <div className="flex flex-col sm:flex-row gap-4 justify-between items-start sm:items-center bg-card p-4 rounded-lg border border-border">
      <div className="flex flex-wrap items-center gap-3">
        <div className="flex items-center gap-1.5 bg-muted/50 p-1 rounded-md border border-border">
          {dbs.map((db) => (
            <button
              key={db}
              type="button"
              onClick={() => onSelectDb(db)}
              className={`px-3 py-1 text-xs font-medium rounded transition-colors ${
                activeDb === db
                  ? "bg-primary text-primary-foreground shadow-sm"
                  : "text-muted-foreground hover:text-foreground"
              }`}
            >
              {db}
            </button>
          ))}
        </div>
        <div className="flex items-center gap-1.5 bg-muted/50 p-1 rounded-md border border-border">
          {availableTables.map((t) => (
            <button
              key={t}
              type="button"
              onClick={() => onSelectTable(t)}
              className={`px-3 py-1 text-xs font-medium rounded transition-colors ${
                activeTable === t
                  ? "bg-background text-foreground shadow-sm border border-border"
                  : "text-muted-foreground hover:text-foreground"
              }`}
            >
              {t}
            </button>
          ))}
        </div>
      </div>
      <div className="relative w-full sm:w-64">
        <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
        <input
          type="text"
          value={filterQuery}
          onChange={(e) => onChangeFilter(e.target.value)}
          placeholder="Filter records..."
          className="w-full bg-background border border-border rounded-md pl-9 pr-3 py-1.5 text-xs text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-1 focus:ring-primary"
        />
      </div>
    </div>
  );
}

interface GridProps {
  rows: Array<Record<string, unknown>>;
  onInspect: (col: string, classified: ClassifiedCell) => void;
}

function DatabaseGrid({ rows, onInspect }: GridProps) {
  if (rows.length === 0) {
    return (
      <div className="p-12 text-center text-sm text-muted-foreground border border-dashed border-border rounded-lg">
        No records found matching criteria.
      </div>
    );
  }

  const columns = Object.keys(rows[0]);

  return (
    <div className="overflow-x-auto border border-border rounded-lg bg-card">
      <table className="w-full border-collapse text-left text-xs">
        <thead className="bg-muted/60 border-b border-border text-muted-foreground">
          <tr>
            {columns.map((c) => (
              <th key={c} className="p-3 font-semibold">
                {c}
              </th>
            ))}
          </tr>
        </thead>
        <tbody className="divide-y divide-border">
          {rows.map((row, idx) => (
            <tr key={idx} className="hover:bg-muted/30 transition-colors">
              {columns.map((c) => (
                <td key={c} className="p-3 font-mono">
                  <DatabaseCellRenderer
                    columnName={c}
                    value={row[c]}
                    onInspect={onInspect}
                  />
                </td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

interface CellRendererProps {
  columnName: string;
  value: unknown;
  onInspect: (col: string, classified: ClassifiedCell) => void;
}

function DatabaseCellRenderer({ columnName, value, onInspect }: CellRendererProps) {
  const classified = useMemo(() => classifyCell(value), [value]);

  if (classified.cellType === "null") {
    return <span className="text-muted-foreground font-medium italic font-sans text-xs">NULL</span>;
  }

  if (classified.cellType === "bool") {
    const isTrue = value === true;

    return (
      <span
        className={`inline-flex items-center px-1.5 py-0.5 rounded text-[10px] font-bold ${
          isTrue
            ? "bg-emerald-500/15 text-emerald-700 dark:text-emerald-400 border border-emerald-500/30"
            : "bg-rose-500/15 text-rose-700 dark:text-rose-400 border border-rose-500/30"
        }`}
      >
        {classified.badgeText}
      </span>
    );
  }

  if (classified.cellType === "blob") {
    return (
      <button
        type="button"
        onClick={() => onInspect(columnName, classified)}
        className="inline-flex items-center gap-1 px-2 py-0.5 rounded text-[11px] font-semibold bg-purple-500/15 text-purple-700 dark:text-purple-300 border border-purple-500/30 hover:bg-purple-500/25 transition-colors cursor-pointer"
      >
        <Binary className="h-3 w-3" />
        {classified.badgeText}
      </button>
    );
  }

  if (classified.cellType === "json") {
    return (
      <button
        type="button"
        onClick={() => onInspect(columnName, classified)}
        className="inline-flex items-center gap-1 px-2 py-0.5 rounded text-[11px] font-semibold bg-sky-500/15 text-sky-700 dark:text-sky-300 border border-sky-500/30 hover:bg-sky-500/25 transition-colors cursor-pointer"
      >
        <Code2 className="h-3 w-3" />
        {classified.badgeText}
      </button>
    );
  }

  if (classified.cellType === "multiline") {
    return (
      <button
        type="button"
        onClick={() => onInspect(columnName, classified)}
        className="inline-flex items-center gap-1 px-2 py-0.5 rounded text-[11px] font-semibold bg-amber-500/15 text-amber-800 dark:text-amber-300 border border-amber-500/30 hover:bg-amber-500/25 transition-colors cursor-pointer max-w-[280px] truncate"
      >
        <FileText className="h-3 w-3 flex-shrink-0" />
        <span className="truncate">{classified.badgeText}</span>
      </button>
    );
  }

  if (classified.cellType === "text") {
    return (
      <button
        type="button"
        onClick={() => onInspect(columnName, classified)}
        className="text-muted-foreground hover:text-sky-400 underline decoration-dotted underline-offset-2 transition-colors cursor-pointer text-left"
      >
        {classified.badgeText}
      </button>
    );
  }

  return <span className="text-foreground">{classified.badgeText}</span>;
}

interface ModalProps {
  target: InspectTarget;
  hasCopied: boolean;
  onClose: () => void;
  onCopyRaw: (text: string) => void;
}

function CellInspectionModal({ target, hasCopied, onClose, onCopyRaw }: ModalProps) {
  const [activeTab, setActiveTab] = useState<"hex" | "json" | "text">(() => {
    if (target.classified.cellType === "json") return "json";
    if (target.classified.cellType === "multiline" || target.classified.cellType === "text") return "text";

    return "hex";
  });

  const rawString = useMemo(() => {
    const raw = target.classified.rawValue;
    if (typeof raw === "string") return raw;

    return JSON.stringify(raw, null, 2);
  }, [target]);

  const hexLines = useMemo(() => {
    return generateHexDump(rawString);
  }, [rawString]);

  const jsonFormatted = useMemo(() => {
    const raw = target.classified.rawValue;
    if (typeof raw === "object" && raw !== null) {
      return JSON.stringify(raw, null, 2);
    }

    if (typeof raw === "string") {
      try {
        return JSON.stringify(JSON.parse(raw), null, 2);
      } catch {
        return raw;
      }
    }

    return String(raw);
  }, [target]);

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/70 backdrop-blur-sm p-4">
      <div className="w-full max-w-3xl max-h-[85vh] bg-card border border-border rounded-xl shadow-2xl flex flex-col overflow-hidden">
        <ModalHeader target={target} onClose={onClose} />
        <ModalNav activeTab={activeTab} onChangeTab={setActiveTab} />
        <div className="flex-1 overflow-auto p-4 bg-muted/20 font-mono text-xs">
          {activeTab === "hex" && <HexDumpView content={hexLines} />}
          {activeTab === "json" && <JsonDumpView content={jsonFormatted} />}
          {activeTab === "text" && <TextDumpView content={rawString} />}
        </div>
        <ModalFooter
          hasCopied={hasCopied}
          activeTab={activeTab}
          onClose={onClose}
          onCopy={() => {
            const textToCopy = activeTab === "hex" ? hexLines : activeTab === "json" ? jsonFormatted : rawString;
            onCopyRaw(textToCopy);
          }}
        />
      </div>
    </div>
  );
}

function ModalHeader({ target, onClose }: { target: InspectTarget; onClose: () => void }) {
  return (
    <div className="flex items-center justify-between px-5 py-3.5 border-b border-border bg-card">
      <div className="flex items-center gap-2.5">
        <h3 className="text-sm font-bold text-foreground">Cell Inspector</h3>
        <span className="px-2 py-0.5 rounded text-[11px] font-mono bg-primary/10 text-primary border border-primary/20">
          {target.columnName}
        </span>
        <span className="px-2 py-0.5 rounded text-[11px] font-mono bg-muted text-muted-foreground border border-border">
          {target.classified.badgeText}
        </span>
      </div>
      <button
        type="button"
        onClick={onClose}
        className="h-7 w-7 rounded flex items-center justify-center text-muted-foreground hover:text-foreground hover:bg-muted transition-colors"
      >
        <X className="h-4 w-4" />
      </button>
    </div>
  );
}

function ModalNav({
  activeTab,
  onChangeTab,
}: {
  activeTab: "hex" | "json" | "text";
  onChangeTab: (t: "hex" | "json" | "text") => void;
}) {
  const tabs = [
    { id: "hex", label: "Hex Dump" },
    { id: "json", label: "Formatted JSON" },
    { id: "text", label: "Decoded String" },
  ] as const;

  return (
    <div className="flex gap-1 px-5 pt-3 border-b border-border bg-card">
      {tabs.map((tab) => (
        <button
          key={tab.id}
          type="button"
          onClick={() => onChangeTab(tab.id)}
          className={`px-3.5 py-1.5 text-xs font-semibold rounded-t border-b-2 transition-colors ${
            activeTab === tab.id
              ? "text-primary border-primary bg-primary/5"
              : "text-muted-foreground border-transparent hover:text-foreground"
          }`}
        >
          {tab.label}
        </button>
      ))}
    </div>
  );
}

function HexDumpView({ content }: { content: string }) {
  return (
    <div className="bg-muted/40 p-3.5 rounded-lg border border-border overflow-x-auto text-[11px] leading-relaxed text-sky-700 dark:text-sky-300 whitespace-pre shadow-inner">
      {content}
    </div>
  );
}

function JsonDumpView({ content }: { content: string }) {
  return (
    <pre className="bg-muted/40 p-3.5 rounded-lg border border-border overflow-x-auto text-[11px] leading-relaxed text-emerald-700 dark:text-emerald-300 whitespace-pre-wrap word-break shadow-inner">
      {content}
    </pre>
  );
}

function TextDumpView({ content }: { content: string }) {
  return (
    <pre className="bg-muted/40 p-3.5 rounded-lg border border-border overflow-x-auto text-[11px] leading-relaxed text-amber-800 dark:text-amber-300 whitespace-pre-wrap word-break shadow-inner">
      {content}
    </pre>
  );
}

interface ModalFooterProps {
  hasCopied: boolean;
  activeTab: string;
  onClose: () => void;
  onCopy: () => void;
}

function ModalFooter({ hasCopied, activeTab, onClose, onCopy }: ModalFooterProps) {
  return (
    <div className="flex items-center justify-between px-5 py-3 border-t border-border bg-card">
      <span className="text-xs text-muted-foreground">
        Active View: <strong className="text-foreground uppercase">{activeTab}</strong>
      </span>
      <div className="flex items-center gap-2">
        <button
          type="button"
          onClick={onCopy}
          className="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-md text-xs font-medium bg-muted text-foreground hover:bg-muted/80 border border-border transition-colors"
        >
          {hasCopied ? <Check className="h-3.5 w-3.5 text-emerald-400" /> : <Copy className="h-3.5 w-3.5" />}
          {hasCopied ? "Copied!" : "Copy Content"}
        </button>
        <button
          type="button"
          onClick={onClose}
          className="px-3.5 py-1.5 rounded-md text-xs font-semibold bg-primary text-primary-foreground hover:bg-primary/90 transition-colors"
        >
          Done
        </button>
      </div>
    </div>
  );
}

function formatByteSize(bytes: number): string {
  if (bytes <= 0) {
    return "0 B";
  }

  if (bytes < 1024) {
    return `${bytes} B`;
  }

  if (bytes < 1024 * 1024) {
    return `${(bytes / 1024).toFixed(1)} KB`;
  }

  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
}

function isBinaryString(str: string): boolean {
  if (!str || str.length === 0) {
    return false;
  }

  if (str.includes("\uFFFD") || str.includes("\x00")) {
    return true;
  }

  const checkLen = Math.min(str.length, 256);
  for (let i = 0; i < checkLen; i++) {
    const code = str.charCodeAt(i);
    if (code < 32 && code !== 9 && code !== 10 && code !== 13) {
      return true;
    }
  }

  return false;
}

function classifyCell(val: unknown): ClassifiedCell {
  if (val === null || val === undefined) {
    return { cellType: "null", badgeText: "NULL", rawValue: val, sizeBytes: 0 };
  }

  if (typeof val === "boolean") {
    return { cellType: "bool", badgeText: val ? "[BOOL true]" : "[BOOL false]", rawValue: val, sizeBytes: 1 };
  }

  if (typeof val === "object") {
    return classifyObjectCell(val as object);
  }

  if (typeof val === "string") {
    return classifyStringCell(val);
  }

  return { cellType: "plain", badgeText: String(val), rawValue: val, sizeBytes: String(val).length };
}

function classifyObjectCell(val: object): ClassifiedCell {
  const isArr = Array.isArray(val);
  const count = isArr ? `${(val as unknown[]).length} items` : `${Object.keys(val).length} keys`;
  const jsonStr = JSON.stringify(val);

  return { cellType: "json", badgeText: `[JSON ${count}]`, rawValue: val, sizeBytes: jsonStr.length };
}

function classifyStringCell(val: string): ClassifiedCell {
  if (isBinaryString(val)) {
    return { cellType: "blob", badgeText: `[BLOB ${formatByteSize(val.length)}]`, rawValue: val, sizeBytes: val.length };
  }

  const jsonResult = tryParseJsonString(val);
  if (jsonResult) {
    return jsonResult;
  }

  if (val.includes("\n") || val.includes("\r")) {
    return classifyMultilineString(val);
  }

  if (val.length > 70) {
    return { cellType: "text", badgeText: `${val.slice(0, 68)}...`, rawValue: val, sizeBytes: val.length };
  }

  return { cellType: "plain", badgeText: val, rawValue: val, sizeBytes: val.length };
}

function tryParseJsonString(val: string): ClassifiedCell | null {
  const trimmed = val.trim();
  const isCandidate = (trimmed.startsWith("{") && trimmed.endsWith("}")) || (trimmed.startsWith("[") && trimmed.endsWith("]"));
  if (!isCandidate) {
    return null;
  }

  try {
    const parsed = JSON.parse(trimmed);
    if (typeof parsed === "object" && parsed !== null) {
      const isArr = Array.isArray(parsed);
      const count = isArr ? `${parsed.length} items` : `${Object.keys(parsed).length} keys`;

      return { cellType: "json", badgeText: `[JSON ${count}]`, rawValue: parsed, sizeBytes: val.length };
    }
  } catch {
    return null;
  }

  return null;
}

function classifyMultilineString(val: string): ClassifiedCell {
  const lines = val.split(/\r?\n/);
  const firstLine = lines[0].slice(0, 32);

  return { cellType: "multiline", badgeText: `${firstLine} ↵ [${lines.length} lines]`, rawValue: val, sizeBytes: val.length };
}

function generateHexDump(raw: unknown): string {
  const str = typeof raw === "string" ? raw : JSON.stringify(raw);
  if (!str) {
    return "(empty)";
  }

  const bytes = extractStringBytes(str);

  return buildHexDumpLines(bytes).join("\n");
}

function extractStringBytes(str: string): number[] {
  const bytes: number[] = [];
  for (let i = 0; i < str.length; i++) {
    bytes.push(str.charCodeAt(i) & 0xff);
  }

  return bytes;
}

function buildHexDumpLines(bytes: number[]): string[] {
  const lines: string[] = [];
  for (let offset = 0; offset < bytes.length; offset += 16) {
    const chunk = bytes.slice(offset, offset + 16);
    lines.push(formatHexLine(offset, chunk));
  }

  return lines;
}

function formatHexLine(offset: number, chunk: number[]): string {
  const hexParts = chunk.map((b) => b.toString(16).padStart(2, "0").toUpperCase());
  while (hexParts.length < 16) {
    hexParts.push("  ");
  }

  const ascii = chunk.map((b) => (b >= 32 && b <= 126 ? String.fromCharCode(b) : ".")).join("");
  const left = hexParts.slice(0, 8).join(" ");
  const right = hexParts.slice(8).join(" ");
  const offStr = offset.toString(16).padStart(8, "0").toUpperCase();

  return `${offStr}  ${left}  ${right}  |${ascii}|`;
}
