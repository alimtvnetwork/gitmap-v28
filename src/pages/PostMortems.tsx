import DocsLayout from "@/components/docs/DocsLayout";
import { postMortems, type PostMortemEntry } from "@/data/postMortems";
import { AlertTriangle, Shield, Database, GitBranch, FolderSync, Settings, Tag } from "lucide-react";
import { useState } from "react";

const categoryConfig: Record<PostMortemEntry["category"], { label: string; icon: typeof AlertTriangle; color: string }> = {
  update: { label: "Update", icon: Settings, color: "text-yellow-500" },
  database: { label: "Database", icon: Database, color: "text-blue-500" },
  release: { label: "Release", icon: GitBranch, color: "text-purple-500" },
  security: { label: "Security", icon: Shield, color: "text-red-500" },
  migration: { label: "Migration", icon: FolderSync, color: "text-orange-500" },
  general: { label: "General", icon: AlertTriangle, color: "text-muted-foreground" },
};

const allCategories = Object.keys(categoryConfig) as PostMortemEntry["category"][];

const PostMortemsPage = () => {
  const [activeFilter, setActiveFilter] = useState<PostMortemEntry["category"] | "all">("all");

  const filtered = activeFilter === "all"
    ? postMortems
    : postMortems.filter((pm) => pm.category === activeFilter);

  return (
    <DocsLayout>
      <div className="mb-6">
        <h1 className="text-3xl font-heading font-bold docs-h1">Post-Mortems</h1>
        <p className="text-muted-foreground text-sm mt-1">
          {postMortems.length} documented issues &middot; lessons learned &amp; prevention rules
        </p>
      </div>

      {/* Category filters */}
      <div className="flex flex-wrap gap-2 mb-6">
        <button
          onClick={() => setActiveFilter("all")}
          className={`text-xs font-mono px-2.5 py-1 rounded border transition-colors duration-300 hover:shadow-sm hover:shadow-primary/10 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 focus-visible:ring-offset-background ${
            activeFilter === "all"
              ? "border-primary bg-primary/10 text-foreground dark:bg-primary/15 dark:text-primary"
              : "border-border text-muted-foreground hover:text-foreground hover:border-primary/40"
          }`}
        >
          All ({postMortems.length})
        </button>
        {allCategories.map((cat) => {
          const count = postMortems.filter((pm) => pm.category === cat).length;
          if (count === 0) return null;
          const config = categoryConfig[cat];
          return (
            <button
              key={cat}
              onClick={() => setActiveFilter(cat)}
              className={`text-xs font-mono px-2.5 py-1 rounded border transition-colors duration-300 hover:shadow-sm hover:shadow-primary/10 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 focus-visible:ring-offset-background ${
                activeFilter === cat
                  ? "border-primary bg-primary/10 text-foreground dark:bg-primary/15 dark:text-primary"
                  : "border-border text-muted-foreground hover:text-foreground hover:border-primary/40"
              }`}
            >
              {config.label} ({count})
            </button>
          );
        })}
      </div>

      {/* Post-mortem list */}
      <div className="space-y-3">
        {filtered.map((pm) => {
          const config = categoryConfig[pm.category];
          const Icon = config.icon;

          return (
            <div
              key={pm.id}
              className="flex items-start gap-4 px-4 py-3 rounded-lg border border-border bg-card hover:bg-muted/50 transition-colors"
            >
              <div className={`mt-0.5 shrink-0 ${config.color}`}>
                <Icon className="h-4 w-4" />
              </div>
              <div className="flex-1 min-w-0">
                <div className="flex items-center gap-2 flex-wrap">
                  <span className="font-mono font-semibold text-sm">
                    #{pm.id}
                  </span>
                  <span className="text-sm font-medium">{pm.title}</span>
                  {pm.version && (
                    <span className="inline-flex items-center gap-1 text-[10px] font-mono px-1.5 py-0.5 rounded bg-primary/10 text-foreground border border-primary/20 dark:bg-primary/15 dark:text-primary transition-colors duration-300 hover:border-primary/40 hover:shadow-sm hover:shadow-primary/10">
                      <Tag className="h-2.5 w-2.5" />
                      {pm.version}
                    </span>
                  )}
                </div>
                <p className="text-xs text-muted-foreground mt-1 leading-relaxed">
                  {pm.summary}
                </p>
              </div>
              <span className={`text-[10px] font-mono px-1.5 py-0.5 rounded border border-border shrink-0 ${config.color}`}>
                {config.label}
              </span>
            </div>
          );
        })}
      </div>

      {filtered.length === 0 && (
        <div className="text-center py-12 text-muted-foreground text-sm">
          No post-mortems in this category.
        </div>
      )}
    </DocsLayout>
  );
};

export default PostMortemsPage;
