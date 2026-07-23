import DocsLayout from "@/components/docs/DocsLayout";
import TerminalDemo from "@/components/docs/TerminalDemo";
import { changelog } from "@/data/changelog";
import { Tag, ChevronDown, ChevronRight } from "lucide-react";
import { useMemo, useState } from "react";
import { motion, AnimatePresence } from "framer-motion";
import {
  classifyChangelogItem,
  TAG_LABELS,
  TAG_ORDER,
  type ChangelogTag,
} from "@/lib/changelogTags";

const terminalLines = [
  { text: "gitmap list-versions", type: "input" as const, delay: 800 },
  { text: "", type: "output" as const },
  { text: "  VERSION    CHANGES  HIGHLIGHTS", type: "header" as const },
  { text: "  ───────    ───────  ──────────", type: "output" as const },
  { text: "  v6.65.0        6   release-undo (ru): undo a tagged release", type: "accent" as const },
  { text: "  v6.63.0        4   cfrp: re-enabled privatize prior 5 versions", type: "output" as const },
  { text: "  v6.60.0        5   commit-in: resumable via state.json checkpoint", type: "output" as const },
  { text: "  v6.57.0        7   backup ls/prune, ssh status health check", type: "output" as const },
  { text: "  v6.55.0        9   cpm merge: settings/bookmarks/extensions", type: "output" as const },
  { text: "  v6.50.0       12   cfr/cfrp pretty runner, spinners, --dry-run", type: "output" as const },
  { text: "", type: "output" as const },
  { text: "  6 versions shown · 43 total changes", type: "accent" as const },
  { text: "", type: "output" as const },
  { text: "gitmap changelog v6.65.0", type: "input" as const, delay: 1000 },
  { text: "", type: "output" as const },
  { text: "  ## v6.65.0", type: "header" as const },
  { text: "  • release-undo (ru): delete local + remote tag + sidecar JSON", type: "output" as const },
  { text: "  • Defaults to latest release when no version is given", type: "output" as const },
  { text: "  • --keep-remote skips remote tag deletion", type: "output" as const },
  { text: "  • Copy-friendly summary line for task-completion reports", type: "output" as const },
];

const ChangelogPage = () => {
  const [expandedVersions, setExpandedVersions] = useState<Set<string>>(
    new Set(changelog.slice(0, 3).map((e) => e.version))
  );

  const [activeTags, setActiveTags] = useState<Set<ChangelogTag>>(new Set());

  const toggle = (version: string) => {
    setExpandedVersions((prev) => {
      const next = new Set(prev);
      next.has(version) ? next.delete(version) : next.add(version);
      return next;
    });
  };

  const toggleTag = (tag: ChangelogTag) =>
    setActiveTags((prev) => {
      const next = new Set(prev);
      next.has(tag) ? next.delete(tag) : next.add(tag);
      return next;
    });

  const expandAll = () => setExpandedVersions(new Set(changelog.map((e) => e.version)));
  const collapseAll = () => setExpandedVersions(new Set());

  const filteredChangelog = useMemo(() => {
    if (activeTags.size === 0) return changelog;
    return changelog
      .map((entry) => ({
        ...entry,
        items: entry.items.filter((item) =>
          classifyChangelogItem(item).some((t) => activeTags.has(t)),
        ),
      }))
      .filter((entry) => entry.items.length > 0);
  }, [activeTags]);

  return (
    <DocsLayout>
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-3xl font-heading font-bold docs-h1">Changelog</h1>
          <p className="text-muted-foreground text-sm mt-1">
            {filteredChangelog.length} of {changelog.length} releases shown
          </p>
        </div>
        <div className="flex gap-2">
          <button onClick={expandAll} className="text-xs font-mono text-muted-foreground hover:text-foreground transition-colors px-2 py-1 rounded border border-border">
            Expand all
          </button>
          <button onClick={collapseAll} className="text-xs font-mono text-muted-foreground hover:text-foreground transition-colors px-2 py-1 rounded border border-border">
            Collapse all
          </button>
        </div>
      </div>

      <div className="mb-6 rounded-lg border border-primary/20 bg-primary/5 px-4 py-3">
        <p className="text-xs font-mono text-foreground/90">
          <span className="text-primary font-semibold">Tip · Help UX (v5.42.0+):</span>{" "}
          discover commands fast with{" "}
          <code className="px-1 py-0.5 rounded bg-muted text-foreground">gitmap help --compact</code>,{" "}
          <code className="px-1 py-0.5 rounded bg-muted text-foreground">--groups</code>,{" "}
          <code className="px-1 py-0.5 rounded bg-muted text-foreground">--filter &lt;q&gt;</code> /{" "}
          <code className="px-1 py-0.5 rounded bg-muted text-foreground">-f</code>, or{" "}
          <code className="px-1 py-0.5 rounded bg-muted text-foreground">--json</code>{" "}
          (v5.43.0+, schema:{" "}
          <a
            href="https://github.com/alimtvnetwork/gitmap-v28/blob/main/spec/08-json-schemas/help-json.schema.json"
            target="_blank"
            rel="noreferrer"
            className="text-primary underline-offset-2 hover:underline"
          >
            help-json.schema.json
          </a>
          ).
        </p>
      </div>

      <div className="mb-6 flex flex-wrap items-center gap-2">
        <span className="text-xs font-mono text-muted-foreground mr-1">Filter:</span>
        {TAG_ORDER.map((tag) => {
          const active = activeTags.has(tag);
          return (
            <button
              key={tag}
              onClick={() => toggleTag(tag)}
              className={`text-xs font-mono px-2 py-0.5 rounded border transition-colors ${
                active
                  ? "border-primary/60 bg-primary/10 text-foreground"
                  : "border-border text-muted-foreground hover:text-foreground hover:border-border/80"
              }`}
              aria-pressed={active}
            >
              {TAG_LABELS[tag]}
            </button>
          );
        })}
        {activeTags.size > 0 && (
          <button
            onClick={() => setActiveTags(new Set())}
            className="text-xs font-mono text-muted-foreground hover:text-foreground underline-offset-2 hover:underline ml-1"
          >
            clear
          </button>
        )}
      </div>

      <div className="mb-8">
        <TerminalDemo title="gitmap — version history" lines={terminalLines} autoPlay />
      </div>

      <div className="relative">
        {/* Timeline line */}
        <div className="absolute left-[15px] top-0 bottom-0 w-px bg-border" />

        <div className="space-y-2">
          {filteredChangelog.map((entry, i) => {
            const isOpen = expandedVersions.has(entry.version);
            const isLatest = i === 0;

            return (
              <div key={entry.version} className="relative pl-10">
                {/* Timeline dot */}
                <div className={`absolute left-[10px] top-3 h-[11px] w-[11px] rounded-full border-2 ${isLatest ? "border-primary bg-primary" : "border-muted-foreground/40 bg-background"}`} />

                <button
                  onClick={() => toggle(entry.version)}
                  className="w-full flex flex-col gap-1 px-4 py-2.5 rounded-lg border border-border bg-card hover:bg-muted/50 transition-colors text-left"
                >
                  <div className="flex items-center gap-3">
                    {isOpen ? (
                      <ChevronDown className="h-4 w-4 text-primary shrink-0" />
                    ) : (
                      <ChevronRight className="h-4 w-4 text-muted-foreground shrink-0" />
                    )}
                    <Tag className="h-3.5 w-3.5 text-primary shrink-0" />
                    <span className="font-mono font-semibold text-sm">{entry.version}</span>
                    {isLatest && (
                      <span className="text-[10px] font-mono px-1.5 py-0.5 rounded bg-primary/10 text-foreground border border-primary/20 dark:bg-primary/15 dark:text-primary transition-colors duration-300 hover:border-primary/40 hover:shadow-sm hover:shadow-primary/10">
                        latest
                      </span>
                    )}
                    {entry.date && (
                      <span className="text-xs text-muted-foreground font-mono">
                        {entry.date}
                      </span>
                    )}
                    <span className="text-xs text-muted-foreground ml-auto">
                      {entry.items.length} change{entry.items.length !== 1 ? "s" : ""}
                    </span>
                  </div>
                  {entry.subtitle && (
                    <p className="text-xs text-muted-foreground/90 pl-[3.25rem] pr-4 leading-snug font-sans">
                      {entry.subtitle}
                    </p>
                  )}
                </button>

                <AnimatePresence initial={false}>
                  {isOpen && (
                    <motion.div
                      initial={{ height: 0, opacity: 0 }}
                      animate={{ height: "auto", opacity: 1 }}
                      exit={{ height: 0, opacity: 0 }}
                      transition={{ duration: 0.2 }}
                      className="overflow-hidden"
                    >
                      <ul className="mt-1 ml-12 mr-2 space-y-1 pb-2 border-l border-border/60 pl-4">
                        {entry.items.map((item, j) => (
                          <li key={j} className="text-sm text-muted-foreground flex gap-2">
                            <span className="text-primary mt-1.5 shrink-0">•</span>
                            <span className="leading-relaxed">{item}</span>
                          </li>
                        ))}
                      </ul>
                    </motion.div>
                  )}
                </AnimatePresence>
              </div>
            );
          })}
        </div>
      </div>
    </DocsLayout>
  );
};

export default ChangelogPage;
