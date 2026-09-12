import { Github, ShieldCheck } from "lucide-react";

export function GitHubAuthCard() {
  return (
    <div className="rounded-xl border border-border bg-card p-6 shadow-sm space-y-4">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2 text-foreground font-semibold text-lg">
          <Github className="h-5 w-5 text-primary" />
          <h2>GitHub CLI Authentication</h2>
        </div>
        <span className="inline-flex items-center gap-1 text-xs font-mono px-2 py-0.5 rounded border bg-primary/10 text-primary border-primary/30">
          <ShieldCheck className="h-3.5 w-3.5" />
          GH Token Auto-Resolved
        </span>
      </div>
      <p className="text-sm text-muted-foreground">
        The pipeline monitor communicates directly with GitHub Actions using
        your system credentials (
        <code className="font-mono text-xs bg-muted px-1.5 py-0.5 rounded text-foreground">
          gh auth token
        </code>{" "}
        or{" "}
        <code className="font-mono text-xs bg-muted px-1.5 py-0.5 rounded text-foreground">
          GITHUB_TOKEN
        </code>
        ).
      </p>
    </div>
  );
}

export default GitHubAuthCard;
