import { useState, useMemo } from "react";
import DocsLayout from "@/components/docs/DocsLayout";
import { copyToClipboard } from "@/lib/clipboard";
import { useToast } from "@/hooks/use-toast";
import type { ClassifiedCell, InspectTarget } from "@/components/database/types";
import { classifyCell, generateHexDump } from "@/components/database/cellUtils";
import DatabaseHeader from "@/components/database/DatabaseHeader";
import DatabaseToolbar from "@/components/database/DatabaseToolbar";
import DatabaseGrid from "@/components/database/DatabaseGrid";
import CellInspectionModal from "@/components/database/CellInspectionModal";

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
        Config: { autocommit: true, splitDb: false, port: 8080 },
        MetaBlob: "\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDR\x00\x00\x01\x00",
        CommitLog: "feat: interactive presenter HUD\nstyle: 16:9 typography ramp\nfix: aspect ratio clamp",
      },
    ],
    Branches: [
      { Id: 101, RepoId: 1, Name: "main", IsProtected: true, HeadCommit: "a1b2c3d" },
      { Id: 102, RepoId: 1, Name: "feat/cell-inspect", IsProtected: false, HeadCommit: "f8e7d6c" },
      { Id: 103, RepoId: 2, Name: "stable", IsProtected: true, HeadCommit: "9988776" },
      { Id: 104, RepoId: 3, Name: "dev", IsProtected: false, HeadCommit: "1234567" },
    ],
  },
  pipeline: {
    Runs: [
      {
        RunId: 501,
        Stage: "lint",
        Status: "passed",
        DurationMs: 412,
        Output: "check-relative-paths: PASS\ncheck-nested-ifs: PASS\ncheck-function-lengths: PASS",
      },
      {
        RunId: 502,
        Stage: "test",
        Status: "running",
        DurationMs: 1250,
        Output: "PASS: TestSplitDbMetadataCache (0.02s)\nPASS: TestCellClassifier (0.01s)",
      },
      {
        RunId: 503,
        Stage: "build",
        Status: "pending",
        DurationMs: 0,
        Output: null,
      },
    ],
    Artifacts: [
      { Id: 9001, Name: "gitmap.exe", SizeBytes: 15420112, Hash: "SHA256:4a3b2c1..." },
      { Id: 9002, Name: "coverage.html", SizeBytes: 84200, Hash: "SHA256:9f8e7d6..." },
    ],
  },
  tasks: {
    PendingQueue: [
      { TaskId: "T-01", Prompt: "Extract Database sub-components", Priority: "HIGH", RetryCount: 0 },
      { TaskId: "T-02", Prompt: "Audit 4-plane depth tokens", Priority: "MEDIUM", RetryCount: 1 },
      { TaskId: "T-03", Prompt: "Verify WCAG 2.2 contrast compliance", Priority: "HIGH", RetryCount: 0 },
    ],
  },
};

export default function DatabasePage() {
  const { toast } = useToast();
  const [activeDb, setActiveDb] = useState("master");
  const [activeTable, setActiveTable] = useState("Repositories");
  const [filterQuery, setFilterQuery] = useState("");
  const [inspectTarget, setInspectTarget] = useState<InspectTarget | null>(null);
  const [hasCopied, setHasCopied] = useState(false);

  const availableTables = useMemo(() => {
    return Object.keys(SAMPLE_DATABASES[activeDb] || {});
  }, [activeDb]);

  const rawRows = useMemo(() => {
    return SAMPLE_DATABASES[activeDb]?.[activeTable] || [];
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
        <DatabaseGrid
          rows={filteredRows}
          onInspect={handleOpenInspect}
          classifyCell={classifyCell}
        />
        {inspectTarget && (
          <CellInspectionModal
            target={inspectTarget}
            hasCopied={hasCopied}
            onClose={handleCloseInspect}
            onCopyRaw={handleCopyRaw}
            generateHexDump={generateHexDump}
          />
        )}
      </div>
    </DocsLayout>
  );
}
