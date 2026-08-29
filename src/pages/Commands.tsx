import { useState, useMemo, useRef, useCallback, MutableRefObject } from "react";
import { Copy, Check, Download } from "lucide-react";
import DocsLayout from "@/components/docs/DocsLayout";
import CommandCard from "@/components/docs/CommandCard";
import CommandCategoryGroup from "@/components/docs/CommandCategoryGroup";
import SearchBar from "@/components/docs/SearchBar";
import ClusterCommandDelegation from "@/components/docs/ClusterCommandDelegation";
import { generateCommandsMarkdown, downloadMarkdownFile } from "@/components/docs/commandsMarkdown";
import { commands, Categories, CommandDef, CommandCategory } from "@/data/commands";
import { copyToClipboard } from "@/lib/clipboard";

interface CategorySummaryButtonProps {
  category: CommandCategory;
  count: number;
  onClick: () => void;
}

interface CategorySummaryBannerProps {
  onSelectCategory: (key: string) => void;
}

interface CommandsHeaderProps {
  isCopied: boolean;
  onCopyAll: () => void;
  onDownload: () => void;
}

interface CommandsSearchResultsProps {
  items: CommandDef[];
  search: string;
  onNavigate: (commandName: string) => void;
  commandRefs: MutableRefObject<Record<string, HTMLDivElement | null>>;
}

interface CommandsCategoryGroupsProps {
  filteredCommands: CommandDef[];
  forceOpen: string | null;
  onNavigate: (commandName: string) => void;
  categoryRefs: MutableRefObject<Record<string, HTMLDivElement | null>>;
  commandRefs: MutableRefObject<Record<string, HTMLDivElement | null>>;
}

function highlightElement(element: HTMLElement, onClear: () => void): void {
  element.scrollIntoView({ behavior: "smooth", block: "center" });
  element.classList.add("ring-2", "ring-primary/50", "rounded-lg");
  setTimeout(() => {
    element.classList.remove("ring-2", "ring-primary/50", "rounded-lg");
    onClear();
  }, 1500);
}

function filterCommands(allCommands: CommandDef[], searchQuery: string): CommandDef[] {
  const hasSearch = searchQuery.length > 0;

  if (hasSearch === false) return allCommands;

  const query = searchQuery.toLowerCase();

  return allCommands.filter((cmd) => {
    const isNameMatch = cmd.name.includes(query);
    const isAliasMatch = Boolean(cmd.alias?.includes(query));
    const isDescMatch = cmd.description.toLowerCase().includes(query);

    return isNameMatch || isAliasMatch || isDescMatch;
  });
}

function useCategoryScroll(
  categoryRefs: MutableRefObject<Record<string, HTMLDivElement | null>>,
  setSearch: (value: string) => void,
  setForceOpen: (value: string | null) => void
) {
  return useCallback((key: string) => {
    setSearch("");
    setForceOpen(key);
    setTimeout(() => {
      categoryRefs.current[key]?.scrollIntoView({ behavior: "smooth", block: "start" });
      setForceOpen(null);
    }, 50);
  }, [categoryRefs, setForceOpen, setSearch]);
}

function useCommandNavigate(
  commandRefs: MutableRefObject<Record<string, HTMLDivElement | null>>,
  setSearch: (value: string) => void,
  setForceOpen: (value: string | null) => void,
  setHighlightCmd: (value: string | null) => void
) {
  return useCallback((commandName: string) => {
    const target = commands.find((cmd) => cmd.name === commandName);

    if (Boolean(target) === false) return;

    setSearch("");
    setForceOpen(target!.category);
    setHighlightCmd(commandName);
    setTimeout(() => {
      const element = commandRefs.current[commandName];

      if (element) {
        highlightElement(element, () => setHighlightCmd(null));
      }

      setForceOpen(null);
    }, 100);
  }, [commandRefs, setForceOpen, setHighlightCmd, setSearch]);
}

const CategorySummaryButton = ({ category, count, onClick }: CategorySummaryButtonProps) => {
  return (
    <button
      onClick={onClick}
      className="rounded-lg border border-border bg-card px-3 py-2.5 text-left hover:bg-muted/50 hover:border-primary/40 transition-all duration-200 cursor-pointer group"
    >
      <div className="flex items-center gap-2 mb-1">
        {category.icon && <span className="text-base">{category.icon}</span>}
        <span className="text-lg font-mono font-bold text-primary">{count}</span>
      </div>
      <div className="text-[11px] text-muted-foreground font-mono leading-tight truncate group-hover:text-foreground transition-colors">
        {category.label}
      </div>
    </button>
  );
};

const CategorySummaryBanner = ({ onSelectCategory }: CategorySummaryBannerProps) => {
  return (
    <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-5 gap-2 mb-6">
      {Categories.map((category) => (
        <CategorySummaryButton
          key={category.key}
          category={category}
          count={commands.filter((cmd) => cmd.category === category.key).length}
          onClick={() => onSelectCategory(category.key)}
        />
      ))}
    </div>
  );
};

const CommandsHeader = ({ isCopied, onCopyAll, onDownload }: CommandsHeaderProps) => {
  return (
    <div className="flex items-center justify-between mb-2">
      <h1 className="text-3xl font-heading font-bold docs-h1">Command Reference</h1>
      <div className="flex items-center gap-1">
        <button
          onClick={onCopyAll}
          className="p-2 rounded-lg border border-border text-muted-foreground hover:text-foreground hover:bg-muted/50 transition-colors"
          title="Copy all as Markdown"
        >
          {isCopied ? <Check className="h-4 w-4 text-primary" /> : <Copy className="h-4 w-4" />}
        </button>
        <button
          onClick={onDownload}
          className="p-2 rounded-lg border border-border text-muted-foreground hover:text-foreground hover:bg-muted/50 transition-colors"
          title="Download as .md"
        >
          <Download className="h-4 w-4" />
        </button>
      </div>
    </div>
  );
};

const CommandsSearchResults = ({ items, search, onNavigate, commandRefs }: CommandsSearchResultsProps) => {
  const hasMatches = items.length > 0;

  if (hasMatches === false) {
    return (
      <p className="text-center text-muted-foreground py-8 font-mono text-sm">
        No commands matching "{search}"
      </p>
    );
  }

  return (
    <>
      {items.map((cmd) => (
        <div key={cmd.name} ref={(element) => { commandRefs.current[cmd.name] = element; }}>
          <CommandCard {...cmd} onNavigate={onNavigate} />
        </div>
      ))}
    </>
  );
};

const CommandsCategoryGroups = ({
  filteredCommands,
  forceOpen,
  onNavigate,
  categoryRefs,
  commandRefs,
}: CommandsCategoryGroupsProps) => {
  return (
    <>
      {Categories.map((category) => {
        const categoryCommands = filteredCommands.filter((cmd) => cmd.category === category.key);

        if (categoryCommands.length === 0) return null;

        return (
          <div key={category.key} ref={(element) => { categoryRefs.current[category.key] = element; }}>
            <CommandCategoryGroup
              label={category.label}
              description={category.description}
              icon={category.icon}
              commands={categoryCommands}
              forceOpen={forceOpen === category.key}
              onNavigate={onNavigate}
              commandRefs={commandRefs}
            />
          </div>
        );
      })}
    </>
  );
};

const CommandsPage = () => {
  const [search, setSearch] = useState("");
  const [forceOpen, setForceOpen] = useState<string | null>(null);
  const [, setHighlightCmd] = useState<string | null>(null);
  const [isCopied, setIsCopied] = useState(false);
  const categoryRefs = useRef<Record<string, HTMLDivElement | null>>({});
  const commandRefs = useRef<Record<string, HTMLDivElement | null>>({});

  const scrollToCategory = useCategoryScroll(categoryRefs, setSearch, setForceOpen);
  const handleNavigate = useCommandNavigate(commandRefs, setSearch, setForceOpen, setHighlightCmd);
  const filtered = useMemo(() => filterCommands(commands, search), [search]);
  const isSearching = search.length > 0;

  const handleCopyAll = useCallback(async () => {
    await copyToClipboard(generateCommandsMarkdown(commands, Categories));
    setIsCopied(true);
    setTimeout(() => setIsCopied(false), 2000);
  }, []);

  const handleDownloadMd = useCallback(() => {
    downloadMarkdownFile(generateCommandsMarkdown(commands, Categories), "gitmap-commands.md");
  }, []);

  return (
    <DocsLayout>
      <CommandsHeader isCopied={isCopied} onCopyAll={handleCopyAll} onDownload={handleDownloadMd} />
      <p className="text-muted-foreground mb-6">All {commands.length} gitmap commands organized by category.</p>
      <ClusterCommandDelegation />
      <CategorySummaryBanner onSelectCategory={scrollToCategory} />
      <SearchBar value={search} onChange={setSearch} />
      <div className="mt-6 space-y-3">
        {isSearching ? (
          <CommandsSearchResults items={filtered} search={search} onNavigate={handleNavigate} commandRefs={commandRefs} />
        ) : (
          <CommandsCategoryGroups filteredCommands={filtered} forceOpen={forceOpen} onNavigate={handleNavigate} categoryRefs={categoryRefs} commandRefs={commandRefs} />
        )}
      </div>
    </DocsLayout>
  );
};

export default CommandsPage;
