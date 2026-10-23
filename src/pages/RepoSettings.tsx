import { useState, useEffect, useMemo } from "react";
import DocsLayout from "@/components/docs/DocsLayout";
import {
  FolderGit2,
  ShieldAlert,
  ShieldCheck,
  Search,
  RefreshCw,
  CheckCircle2,
  Sliders,
  ToggleLeft,
  ToggleRight,
  Filter,
} from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Switch } from "@/components/ui/switch";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { useToast } from "@/hooks/use-toast";

export interface RepoScanItem {
  id: number;
  repoName: string;
  relativePath: string;
  branch: string;
  isExcluded: boolean;
  lastScanAt: string;
}

interface ScanReposResponse {
  isSuccess: boolean;
  totalCount: number;
  repos: RepoScanItem[];
}

interface ToggleScanReposResponse {
  isSuccess: boolean;
  modifiedCount: number;
  affectedRepos: string[];
  currentTotal: number;
}

const SAMPLE_REPOS: RepoScanItem[] = [
  {
    id: 1,
    repoName: "gitmap-v28",
    relativePath: "./gitmap-v28",
    branch: "main",
    isExcluded: false,
    lastScanAt: "2026-10-23T11:00:00Z",
  },
  {
    id: 2,
    repoName: "docs",
    relativePath: "./docs",
    branch: "main",
    isExcluded: false,
    lastScanAt: "2026-10-23T11:00:00Z",
  },
  {
    id: 3,
    repoName: "auth-service",
    relativePath: "./auth-service",
    branch: "master",
    isExcluded: false,
    lastScanAt: "2026-10-23T11:00:00Z",
  },
  {
    id: 4,
    repoName: "web-portal",
    relativePath: "./web-portal",
    branch: "main",
    isExcluded: false,
    lastScanAt: "2026-10-23T11:00:00Z",
  },
  {
    id: 5,
    repoName: "payment-worker",
    relativePath: "./payment-worker",
    branch: "main",
    isExcluded: true,
    lastScanAt: "2026-10-23T11:00:00Z",
  },
];

export default function RepoSettingsPage() {
  const { toast } = useToast();
  const [repos, setRepos] = useState<RepoScanItem[]>(SAMPLE_REPOS);
  const [searchQuery, setSearchQuery] = useState("");
  const [rangeInput, setRangeInput] = useState("");
  const [filterMode, setFilterMode] = useState<"all" | "active" | "excluded">("all");
  const [isLoading, setIsLoading] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);

  const fetchRepos = async () => {
    setIsLoading(true);
    try {
      const response = await fetch("/api/scan/repos");
      const hasOk = response.ok;
      if (hasOk) {
        const data: ScanReposResponse = await response.json();
        const isDataSuccess = data.isSuccess && Array.isArray(data.repos);
        if (isDataSuccess) {
          setRepos(data.repos);
        }
      }
    } catch {
      // Fallback to existing state if backend is offline during preview
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    void fetchRepos();
  }, []);

  const totalCount = repos.length;
  const excludedCount = useMemo(
    () => repos.filter((r) => r.isExcluded).length,
    [repos]
  );
  const activeCount = totalCount - excludedCount;

  const filteredRepos = useMemo(() => {
    const query = searchQuery.trim().toLowerCase();
    const hasQuery = query.length > 0;

    return repos.filter((item) => {
      const matchesSearch =
        !hasQuery ||
        item.repoName.toLowerCase().includes(query) ||
        item.relativePath.toLowerCase().includes(query) ||
        item.branch.toLowerCase().includes(query);

      if (!matchesSearch) return false;

      if (filterMode === "active") {
        return !item.isExcluded;
      }
      if (filterMode === "excluded") {
        return item.isExcluded;
      }

      return true;
    });
  }, [repos, searchQuery, filterMode]);

  const toggleRepoStatus = async (item: RepoScanItem) => {
    setIsSubmitting(true);
    const targetEndpoint = item.isExcluded
      ? "/api/scan/enable"
      : "/api/scan/disable";

    try {
      const response = await fetch(targetEndpoint, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ repos: [item.relativePath] }),
      });

      const hasOk = response.ok;
      if (hasOk) {
        const result: ToggleScanReposResponse = await response.json();
        const isActionSuccess = Boolean(result.isSuccess);
        if (isActionSuccess) {
          setRepos((prev) =>
            prev.map((r) =>
              r.id === item.id ? { ...r, isExcluded: !r.isExcluded } : r
            )
          );

          toast({
            title: item.isExcluded ? "Repository Re-Enabled" : "Repository Disabled",
            description: `Scan status updated for ${item.repoName}.`,
          });
          return;
        }
      }

      toast({
        title: "Update Failed",
        description: `Could not update scan status for ${item.repoName}.`,
        variant: "destructive",
      });
    } catch {
      // Optimistic update fallback for standalone preview
      setRepos((prev) =>
        prev.map((r) =>
          r.id === item.id ? { ...r, isExcluded: !r.isExcluded } : r
        )
      );
      toast({
        title: "Mock Mode Updated",
        description: `Status toggled for ${item.repoName} (offline preview).`,
      });
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleRangeAction = async (isDisable: boolean) => {
    const rawTokens = rangeInput.trim();
    const hasTokens = rawTokens.length > 0;
    if (!hasTokens) {
      toast({
        title: "Input Required",
        description: "Please enter indices, ranges (e.g. 1-3, 5), or repo slugs.",
        variant: "destructive",
      });
      return;
    }

    setIsSubmitting(true);
    const targetEndpoint = isDisable ? "/api/scan/disable" : "/api/scan/enable";

    try {
      const response = await fetch(targetEndpoint, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ repos: [rawTokens] }),
      });

      const hasOk = response.ok;
      if (hasOk) {
        const result: ToggleScanReposResponse = await response.json();
        const isActionSuccess = Boolean(result.isSuccess);
        if (isActionSuccess) {
          toast({
            title: isDisable ? "Repositories Disabled" : "Repositories Enabled",
            description: `Modified ${result.modifiedCount} repositories. Total excluded: ${result.currentTotal}.`,
          });
          setRangeInput("");
          void fetchRepos();
          return;
        }
      }

      toast({
        title: "Range Action Failed",
        description: "Could not apply range selection to repositories.",
        variant: "destructive",
      });
    } catch {
      toast({
        title: "Action Recorded",
        description: `Range "${rawTokens}" processed in preview mode.`,
      });
      setRangeInput("");
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <DocsLayout>
      <div className="space-y-6 max-w-7xl mx-auto py-4">
        {/* Header Section */}
        <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 border-b border-border/40 pb-5">
          <div>
            <div className="flex items-center gap-2">
              <FolderGit2 className="h-6 w-6 text-primary" />
              <h1 className="text-2xl font-bold tracking-tight text-foreground">
                Repository Scan Settings & Exclusions
              </h1>
            </div>
            <p className="text-sm text-muted-foreground mt-1">
              Configure which repositories are included or excluded during
              automated directory scanning and fleet pull synchronization.
            </p>
          </div>
          <div className="flex items-center gap-2">
            <Button
              variant="outline"
              size="sm"
              onClick={() => void fetchRepos()}
              disabled={isLoading || isSubmitting}
            >
              <RefreshCw className={`h-4 w-4 mr-2 ${isLoading ? "animate-spin" : ""}`} />
              Refresh
            </Button>
          </div>
        </div>

        {/* Status KPI Cards */}
        <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
          <Card className="bg-card/50 border-border/60">
            <CardHeader className="pb-2">
              <CardDescription className="text-xs font-semibold uppercase tracking-wider">
                Tracked Repositories
              </CardDescription>
              <CardTitle className="text-2xl font-bold flex items-center justify-between">
                <span>{totalCount}</span>
                <Sliders className="h-5 w-5 text-muted-foreground" />
              </CardTitle>
            </CardHeader>
            <CardContent className="text-xs text-muted-foreground">
              Total Git repositories discovered in Split-DB
            </CardContent>
          </Card>

          <Card className="bg-card/50 border-border/60">
            <CardHeader className="pb-2">
              <CardDescription className="text-xs font-semibold uppercase tracking-wider">
                Active Scanning
              </CardDescription>
              <CardTitle className="text-2xl font-bold text-emerald-500 flex items-center justify-between">
                <span>{activeCount}</span>
                <ShieldCheck className="h-5 w-5" />
              </CardTitle>
            </CardHeader>
            <CardContent className="text-xs text-muted-foreground">
              Repositories monitored on every scan & pull
            </CardContent>
          </Card>

          <Card className="bg-card/50 border-border/60">
            <CardHeader className="pb-2">
              <CardDescription className="text-xs font-semibold uppercase tracking-wider">
                Excluded / Disabled
              </CardDescription>
              <CardTitle className="text-2xl font-bold text-rose-500 flex items-center justify-between">
                <span>{excludedCount}</span>
                <ShieldAlert className="h-5 w-5" />
              </CardTitle>
            </CardHeader>
            <CardContent className="text-xs text-muted-foreground">
              Ignored paths in custom_repo_exclusions
            </CardContent>
          </Card>
        </div>

        {/* Range Selector & Batch Controls */}
        <Card className="border-border/60 bg-card/30">
          <CardHeader className="pb-3">
            <CardTitle className="text-base font-semibold flex items-center gap-2">
              <Sliders className="h-4 w-4 text-primary" />
              Range Selection & Batch Toggle
            </CardTitle>
            <CardDescription className="text-xs">
              Enter 1-based index numbers, ranges (e.g.{" "}
              <code className="text-primary font-mono font-medium">1-5, 8, 12</code>
              ), or repo slug names to toggle exclusions in bulk.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="flex flex-col sm:flex-row gap-3">
              <div className="relative flex-1">
                <Input
                  value={rangeInput}
                  onChange={(e) => setRangeInput(e.target.value)}
                  placeholder="e.g. 1-3, 5 or alimtvnetwork/docs, auth-service"
                  className="font-mono text-sm"
                  disabled={isSubmitting}
                />
              </div>
              <div className="flex flex-wrap gap-2">
                <Button
                  size="sm"
                  variant="destructive"
                  onClick={() => void handleRangeAction(true)}
                  disabled={isSubmitting || rangeInput.trim().length === 0}
                >
                  <ToggleLeft className="h-4 w-4 mr-1.5" />
                  Disable Range
                </Button>
                <Button
                  size="sm"
                  variant="secondary"
                  onClick={() => void handleRangeAction(false)}
                  disabled={isSubmitting || rangeInput.trim().length === 0}
                >
                  <ToggleRight className="h-4 w-4 mr-1.5" />
                  Enable Range
                </Button>
                <Button
                  size="sm"
                  variant="outline"
                  onClick={() => {
                    setRangeInput("all");
                  }}
                  disabled={isSubmitting}
                >
                  Select All
                </Button>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Repository Table & Search Filter */}
        <div className="space-y-4">
          <div className="flex flex-col sm:flex-row items-stretch sm:items-center justify-between gap-3">
            <div className="relative flex-1 max-w-sm">
              <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
              <Input
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                placeholder="Filter by repo name, path, branch..."
                className="pl-9 h-9 text-sm"
              />
            </div>
            <div className="flex items-center gap-1.5 bg-muted/50 p-1 rounded-lg border border-border/40 text-xs font-medium">
              <Filter className="h-3.5 w-3.5 text-muted-foreground ml-1 mr-0.5" />
              <button
                type="button"
                onClick={() => setFilterMode("all")}
                className={`px-2.5 py-1 rounded-md transition-colors ${
                  filterMode === "all"
                    ? "bg-background text-foreground shadow-sm"
                    : "text-muted-foreground hover:text-foreground"
                }`}
              >
                All ({totalCount})
              </button>
              <button
                type="button"
                onClick={() => setFilterMode("active")}
                className={`px-2.5 py-1 rounded-md transition-colors ${
                  filterMode === "active"
                    ? "bg-background text-foreground shadow-sm"
                    : "text-muted-foreground hover:text-foreground"
                }`}
              >
                Active ({activeCount})
              </button>
              <button
                type="button"
                onClick={() => setFilterMode("excluded")}
                className={`px-2.5 py-1 rounded-md transition-colors ${
                  filterMode === "excluded"
                    ? "bg-background text-foreground shadow-sm"
                    : "text-muted-foreground hover:text-foreground"
                }`}
              >
                Disabled ({excludedCount})
              </button>
            </div>
          </div>

          <div className="border border-border/60 rounded-xl overflow-hidden bg-card/20">
            <Table>
              <TableHeader className="bg-muted/30">
                <TableRow>
                  <TableHead className="w-16 font-mono text-xs">#</TableHead>
                  <TableHead className="w-32">Status</TableHead>
                  <TableHead>Repository Name</TableHead>
                  <TableHead>Relative Path</TableHead>
                  <TableHead className="w-32">Branch</TableHead>
                  <TableHead className="w-28 text-right">Scan State</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {filteredRepos.length === 0 ? (
                  <TableRow>
                    <TableCell
                      colSpan={6}
                      className="text-center py-8 text-muted-foreground text-sm"
                    >
                      No matching repositories found.
                    </TableCell>
                  </TableRow>
                ) : (
                  filteredRepos.map((item, idx) => (
                    <TableRow
                      key={item.id}
                      className={item.isExcluded ? "bg-rose-500/5" : ""}
                    >
                      <TableCell className="font-mono text-xs text-muted-foreground">
                        {idx + 1}
                      </TableCell>
                      <TableCell>
                        {item.isExcluded ? (
                          <Badge variant="destructive" className="text-xs uppercase">
                            [DISABLED]
                          </Badge>
                        ) : (
                          <Badge
                            variant="secondary"
                            className="text-xs uppercase bg-emerald-500/10 text-emerald-500 border-emerald-500/20"
                          >
                            [ENABLED]
                          </Badge>
                        )}
                      </TableCell>
                      <TableCell className="font-medium text-foreground">
                        {item.repoName}
                      </TableCell>
                      <TableCell className="font-mono text-xs text-muted-foreground">
                        {item.relativePath}
                      </TableCell>
                      <TableCell className="font-mono text-xs">
                        <span className="inline-flex items-center gap-1 text-muted-foreground">
                          {item.branch || "main"}
                        </span>
                      </TableCell>
                      <TableCell className="text-right">
                        <div className="flex items-center justify-end gap-2">
                          <Switch
                            checked={!item.isExcluded}
                            onCheckedChange={() => void toggleRepoStatus(item)}
                            disabled={isSubmitting}
                          />
                        </div>
                      </TableCell>
                    </TableRow>
                  ))
                )}
              </TableBody>
            </Table>
          </div>
        </div>
      </div>
    </DocsLayout>
  );
}
