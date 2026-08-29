import { useState, useMemo, useEffect, useRef } from "react";
import { useSearchParams, SetURLSearchParams } from "react-router-dom";
import { AlertTriangle } from "lucide-react";
import DocsLayout from "@/components/docs/DocsLayout";
import SearchBar from "@/components/docs/SearchBar";
import { DiagnosticChecklist } from "@/components/docs/troubleshooting/DiagnosticChecklist";
import { TroubleshootingIssueCard } from "@/components/docs/troubleshooting/TroubleshootingIssueCard";
import { TroubleshootingFilterButtons } from "@/components/docs/troubleshooting/TroubleshootingFilterButtons";
import { TroubleshootingHelpSection } from "@/components/docs/troubleshooting/TroubleshootingHelpSection";
import {
  CategoryType,
  Category,
  CategoryFilterType,
  Issue,
  categoryMeta,
  issues,
  isValidCategoryKey,
} from "@/data/troubleshootingIssues";

export { CategoryType, CategoryFilterType, categoryMeta, issues, isValidCategoryKey };
export type { Category, Issue };

interface TroubleshootingIssueListProps {
  items: Issue[];
  search: string;
}

function getInitialCategory(rawCategory: string): Category | CategoryFilterType {
  const isAllOrValid = rawCategory === CategoryFilterType.All || isValidCategoryKey(rawCategory);
  if (isAllOrValid) return rawCategory as Category | CategoryFilterType;

  return CategoryFilterType.All;
}

function matchesIssueSearch(issueItem: Issue, query: string): boolean {
  const isTitleMatch = issueItem.title.toLowerCase().includes(query);
  const isSymptomMatch = issueItem.symptom.toLowerCase().includes(query);
  const isCauseMatch = issueItem.cause.toLowerCase().includes(query);
  const isFixMatch = issueItem.fix.toLowerCase().includes(query);
  const isCommandMatch = Boolean(issueItem.fixCommand?.toLowerCase().includes(query));

  return isTitleMatch || isSymptomMatch || isCauseMatch || isFixMatch || isCommandMatch;
}

function filterIssues(
  allIssues: Issue[],
  activeCategory: Category | CategoryFilterType,
  search: string
): Issue[] {
  let rows = allIssues;
  if (activeCategory !== CategoryFilterType.All) {
    rows = rows.filter((item) => item.category === activeCategory);
  }

  if (search) {
    const query = search.toLowerCase();
    rows = rows.filter((item) => matchesIssueSearch(item, query));
  }

  return rows;
}

function calculateCategoryCounts(allIssues: Issue[]): Record<string, number> {
  const counts: Record<string, number> = { [CategoryFilterType.All]: allIssues.length };
  for (const item of allIssues) counts[item.category] = (counts[item.category] ?? 0) + 1;

  return counts;
}

function highlightIssueElement(element: HTMLElement): void {
  element.scrollIntoView({ behavior: "smooth", block: "start" });
  element.classList.add("ring-2", "ring-primary");
  window.setTimeout(() => element.classList.remove("ring-2", "ring-primary"), 2400);
}

function useTroubleshootingUrlSync(
  searchParams: URLSearchParams,
  setSearchParams: SetURLSearchParams,
  search: string,
  activeCategory: Category | CategoryFilterType
): void {
  useEffect(() => {
    const nextParams = new URLSearchParams(searchParams);
    if (search) nextParams.set("search", search);
    else nextParams.delete("search");
    nextParams.delete("q");
    if (activeCategory !== CategoryFilterType.All) nextParams.set("category", activeCategory);
    else nextParams.delete("category");
    if (nextParams.toString() !== searchParams.toString()) {
      setSearchParams(nextParams, { replace: true });
    }
  }, [search, activeCategory, searchParams, setSearchParams]);
}

function useTroubleshootingDeepLink(
  targetId: string | null,
  filteredIssues: Issue[],
  activeCategory: Category | CategoryFilterType,
  setActiveCategory: (category: Category | CategoryFilterType) => void,
  setSearch: (search: string) => void
): void {
  const scrolledIdRef = useRef<string | null>(null);

  useEffect(() => {
    if (Boolean(targetId) === false) return;
    const targetIssue = issues.find((issueItem) => issueItem.id === targetId);
    if (Boolean(targetIssue) === false) return;

    if (activeCategory !== CategoryFilterType.All && activeCategory !== targetIssue!.category) {
      setActiveCategory(CategoryFilterType.All);
      return;
    }

    if (filteredIssues.findIndex((issueItem) => issueItem.id === targetId) === -1) {
      setSearch("");
      return;
    }

    if (scrolledIdRef.current === targetId) return;
    const element = document.getElementById(targetId!);
    if (element) {
      scrolledIdRef.current = targetId;
      highlightIssueElement(element);
    }
  }, [targetId, filteredIssues, activeCategory, setActiveCategory, setSearch]);
}

const TroubleshootingIssueList = ({ items, search }: TroubleshootingIssueListProps) => {
  const hasItems = items.length > 0;
  if (hasItems === false) {
    return (
      <div className="rounded-lg border border-border p-8 text-center">
        <AlertTriangle className="h-8 w-8 text-muted-foreground mx-auto mb-3" />
        <p className="font-mono text-sm text-muted-foreground">
          No issues match "{search}". Try a different keyword or clear the filter.
        </p>
      </div>
    );
  }

  return (
    <section className="space-y-6">
      {items.map((issueItem) => (
        <TroubleshootingIssueCard key={issueItem.id} issue={issueItem} />
      ))}
    </section>
  );
};

const Troubleshooting = () => {
  const [searchParams, setSearchParams] = useSearchParams();
  const initialSearch = searchParams.get("search") ?? searchParams.get("q") ?? "";
  const initialCategory = getInitialCategory(searchParams.get("category") ?? CategoryFilterType.All);

  const [search, setSearch] = useState(initialSearch);
  const [activeCategory, setActiveCategory] = useState<Category | CategoryFilterType>(initialCategory);

  useTroubleshootingUrlSync(searchParams, setSearchParams, search, activeCategory);
  const filtered = useMemo(() => filterIssues(issues, activeCategory, search), [activeCategory, search]);
  const categoryCounts = useMemo(() => calculateCategoryCounts(issues), []);
  useTroubleshootingDeepLink(searchParams.get("id"), filtered, activeCategory, setActiveCategory, setSearch);

  return (
    <DocsLayout>
      <h1 className="text-3xl font-heading font-bold mb-2 docs-h1">Troubleshooting</h1>
      <p className="text-muted-foreground mb-6">
        Common gitmap errors grouped by category, each with the symptom, root cause, and the exact flag or command to fix it. When in doubt, start with <code className="docs-inline-code">gitmap doctor</code>.
      </p>
      <DiagnosticChecklist />
      <SearchBar value={search} onChange={setSearch} placeholder="Search by error, symptom, or fix..." />
      <TroubleshootingFilterButtons activeCategory={activeCategory} categoryCounts={categoryCounts} onSelectCategory={setActiveCategory} />
      <TroubleshootingIssueList items={filtered} search={search} />
      <TroubleshootingHelpSection />
    </DocsLayout>
  );
};

export default Troubleshooting;
