import DocsLayout from "@/components/docs/DocsLayout";
import { Link } from "react-router-dom";
import {
  GitCompare,
  Move,
  GitMerge,
  GitCommit,
  ArrowRight,
  CheckCircle2,
  Clock,
} from "lucide-react";

type Status = "shipped" | "planned";

interface HelpEntry {
  to: string;
  cmd: string;
  alias: string;
  desc: string;
  status: Status;
}

interface HelpSection {
  id: string;
  title: string;
  blurb: string;
  icon: typeof GitCompare;
  entries: HelpEntry[];
}

const sections: HelpSection[] = [
  {
    id: "diff",
    title: "Read-only inspection",
    blurb:
      "Inspect file-state differences between two repository endpoints without writing anything.",
    icon: GitCompare,
    entries: [
      {
        to: "/diff",
        cmd: "diff",
        alias: "df",
        desc: "Side-by-side preview of file differences between LEFT and RIGHT.",
        status: "shipped",
      },
    ],
  },
  {
    id: "move",
    title: "Move",
    blurb:
      "Single-direction file transfer. Overwrites the destination tree with the source tree.",
    icon: Move,
    entries: [
      {
        to: "/mv",
        cmd: "mv",
        alias: "move",
        desc: "Replace RIGHT's working tree with LEFT's content.",
        status: "shipped",
      },
    ],
  },
  {
    id: "merge",
    title: "Merge family — file-state transfer",
    blurb:
      "Three-way merge of file content with conflict prompts. Transfers files, not commit history.",
    icon: GitMerge,
    entries: [
      {
        to: "/merge-both",
        cmd: "merge-both",
        alias: "mb",
        desc: "Bidirectional merge — both sides converge on the union of changes.",
        status: "shipped",
      },
      {
        to: "/merge-left",
        cmd: "merge-left",
        alias: "ml",
        desc: "Pull RIGHT's changes onto LEFT only.",
        status: "shipped",
      },
      {
        to: "/merge-right",
        cmd: "merge-right",
        alias: "mr",
        desc: "Push LEFT's changes onto RIGHT only.",
        status: "shipped",
      },
    ],
  },
  {
    id: "commit",
    title: "Commit family — history transfer (planned)",
    blurb:
      "Replay one side's commits onto the other as fresh commits — preserves the human-readable evolution, not just the final tree.",
    icon: GitCommit,
    entries: [
      {
        to: "/commit-left",
        cmd: "commit-left",
        alias: "cl",
        desc: "Replay RIGHT's commits onto LEFT (writes to LEFT).",
        status: "planned",
      },
      {
        to: "/commit-right",
        cmd: "commit-right",
        alias: "cr",
        desc: "Replay LEFT's commits onto RIGHT (Phase-1 target).",
        status: "planned",
      },
      {
        to: "/commit-both",
        cmd: "commit-both",
        alias: "cb",
        desc: "Interleave both sides' commits by author date onto each side.",
        status: "planned",
      },
    ],
  },
];

const StatusPill = ({ status }: { status: Status }) =>
  status === "shipped" ? (
    <span className="inline-flex items-center gap-1 text-[10px] font-mono px-1.5 py-0.5 rounded border bg-primary/10 text-foreground border-primary/20 transition-colors duration-300 hover:border-primary/40 hover:shadow-sm hover:shadow-primary/10 dark:bg-primary/20 dark:text-primary dark:border-primary/40">
      <CheckCircle2 className="h-2.5 w-2.5" />
      shipped
    </span>
  ) : (
    <span className="inline-flex items-center gap-1 text-[10px] font-mono px-1.5 py-0.5 rounded border bg-destructive/10 text-foreground border-destructive/30 transition-colors duration-300 hover:border-destructive/50 hover:shadow-sm hover:shadow-destructive/10 dark:bg-destructive/20 dark:text-destructive-foreground dark:border-destructive/50">
      <Clock className="h-2.5 w-2.5" />
      planned
    </span>
  );

const HelpIndexPage = () => {
  const totalShipped = sections.reduce(
    (n, s) => n + s.entries.filter((e) => e.status === "shipped").length,
    0,
  );
  const totalPlanned = sections.reduce(
    (n, s) => n + s.entries.filter((e) => e.status === "planned").length,
    0,
  );

  return (
    <DocsLayout>
      <div className="mb-8">
        <h1 className="text-3xl font-heading font-bold docs-h1">Help Index</h1>
        <p className="text-muted-foreground text-sm mt-2 max-w-3xl leading-relaxed">
          Curated entry point for the diff / move / merge / commit command
          families. Every page below is bundled into the docs-site archive
          and served by{" "}
          <code className="docs-inline-code">gitmap help-dashboard</code> so
          you can browse them offline from any installed binary.
        </p>
        <div className="flex flex-wrap gap-2 mt-4">
          <span className="inline-flex items-center gap-1 text-xs font-mono px-2 py-0.5 rounded border bg-primary/10 text-foreground border-primary/20 dark:bg-primary/20 dark:text-primary dark:border-primary/40">
            <CheckCircle2 className="h-3 w-3" />
            {totalShipped} shipped
          </span>
          <span className="inline-flex items-center gap-1 text-xs font-mono px-2 py-0.5 rounded border bg-destructive/10 text-foreground border-destructive/30 dark:bg-destructive/20 dark:text-destructive-foreground dark:border-destructive/50">
            <Clock className="h-3 w-3" />
            {totalPlanned} planned
          </span>
        </div>
      </div>

      <div className="space-y-10">
        {sections.map((section) => {
          const Icon = section.icon;
          return (
            <section key={section.id} id={section.id}>
              <div className="flex items-center gap-3 mb-2">
                <div className="h-9 w-9 rounded-md bg-primary/10 flex items-center justify-center dark:bg-primary/15">
                  <Icon className="h-5 w-5 text-primary" />
                </div>
                <h2 className="text-xl font-heading font-semibold">
                  {section.title}
                </h2>
              </div>
              <p className="text-sm text-muted-foreground mb-4 max-w-3xl leading-relaxed">
                {section.blurb}
              </p>
              <div className="grid gap-3 md:grid-cols-2">
                {section.entries.map((entry) => (
                  <Link
                    key={entry.to}
                    to={entry.to}
                    className="group block rounded-lg border border-border bg-card p-4 hover:border-primary/40 hover:shadow-lg hover:shadow-primary/5 active:translate-y-px focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 focus-visible:ring-offset-background transition-colors duration-300"
                  >
                    <div className="flex items-center justify-between gap-3 mb-2">
                      <div className="flex items-center gap-2 min-w-0">
                        <code className="font-mono font-semibold text-sm text-foreground group-hover:text-primary transition-colors duration-300">
                          {entry.cmd}
                        </code>
                        <span className="text-[10px] font-mono text-muted-foreground">
                          alias: {entry.alias}
                        </span>
                      </div>
                      <StatusPill status={entry.status} />
                    </div>
                    <p className="text-xs text-muted-foreground leading-relaxed">
                      {entry.desc}
                    </p>
                    <div className="flex items-center gap-1 text-[11px] font-mono text-primary/80 group-hover:text-primary mt-3 transition-colors duration-300">
                      Open page
                      <ArrowRight className="h-3 w-3 transition-transform duration-300 group-hover:translate-x-0.5" />
                    </div>
                  </Link>
                ))}
              </div>
            </section>
          );
        })}
      </div>

      <hr className="docs-hr my-10" />

      <section>
        <h2 className="text-xl font-heading font-semibold mb-3">
          Serving this index from <code className="docs-inline-code">help-dashboard</code>
        </h2>
        <p className="text-sm text-muted-foreground leading-relaxed mb-3 max-w-3xl">
          Every route on this page is part of the React docs site that ships in{" "}
          <code className="docs-inline-code">docs-site.zip</code> alongside the
          gitmap binary. Run the dashboard to browse them offline:
        </p>
        <pre className="rounded-md bg-code-bg border border-border p-3 text-sm font-mono overflow-x-auto">
          <code>{`gitmap help-dashboard         # serves dist/ on http://localhost:5173
gitmap hd --port 8080         # custom port
# then open http://localhost:5173/help-index in your browser`}</code>
        </pre>
        <p className="text-xs text-muted-foreground mt-3">
          See <Link to="/help-dashboard" className="text-primary hover:underline">help-dashboard</Link> for full
          flag and shutdown documentation.
        </p>
      </section>
    </DocsLayout>
  );
};

export default HelpIndexPage;
