import { Sparkles, ExternalLink } from "lucide-react";
import { referenceRepositories } from "./referenceReposData";

export function UiUxRepoCatalogCard() {
  return (
    <div className="rounded-xl bg-card border border-border/50 p-6 shadow-sm space-y-6">
      <div className="flex items-center justify-between border-b border-border/40 pb-4">
        <div className="flex items-center gap-2.5">
          <Sparkles className="h-5 w-5 text-primary" />
          <div>
            <h2 className="text-lg font-heading font-semibold text-foreground">
              Authoritative UI/UX &amp; CSS3 Repository Index
            </h2>
            <p className="text-xs text-muted-foreground">
              Curated benchmarks for layout, accessibility primitives, tokens, and micro-interactions.
            </p>
          </div>
        </div>
        <span className="text-xs px-2.5 py-1 rounded-full bg-muted text-muted-foreground border border-border/60 font-medium">
          Design Reference Spec
        </span>
      </div>

      <div className="grid sm:grid-cols-2 lg:grid-cols-3 gap-4">
        {referenceRepositories.map((repo) => (
          <div
            key={repo.name}
            className="p-4 rounded-lg bg-muted/15 border border-border/40 hover:border-primary/50 transition-all flex flex-col justify-between group"
          >
            <div className="space-y-2">
              <div className="flex items-start justify-between gap-2">
                <h3 className="font-heading font-bold text-sm text-foreground group-hover:text-primary transition-colors">
                  {repo.name}
                </h3>
                <span className="text-[10px] px-2 py-0.5 rounded-md bg-muted text-muted-foreground border border-border/40 font-mono">
                  {repo.badge}
                </span>
              </div>
              <div className="text-[11px] text-primary font-medium tracking-tight">
                {repo.category}
              </div>
              <p className="text-xs text-muted-foreground leading-relaxed">
                {repo.description}
              </p>
            </div>
            <div className="pt-3 mt-3 border-t border-border/30 flex items-center justify-between text-xs">
              <span className="text-[11px] text-muted-foreground font-mono truncate max-w-[170px]">
                {repo.repo}
              </span>
              <a
                href={repo.url}
                target="_blank"
                rel="noopener noreferrer"
                className="inline-flex items-center gap-1 text-primary hover:underline font-semibold text-xs"
              >
                <span>Explore</span>
                <ExternalLink className="h-3 w-3" />
              </a>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
