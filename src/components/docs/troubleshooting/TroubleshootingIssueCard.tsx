import CodeBlock from "@/components/docs/CodeBlock";
import { CopyFixButton } from "@/components/docs/troubleshooting/CopyFixButton";
import { CopyLinkButton } from "@/components/docs/troubleshooting/CopyLinkButton";
import { Issue, CATEGORY_META } from "@/data/troubleshootingIssues";

interface IssueProps {
  issue: Issue;
}

interface IssueRelatedLinksProps {
  related?: { label: string; href: string }[];
}

const IssueHeader = ({ issue }: IssueProps) => {
  const meta = CATEGORY_META[issue.category];
  const Icon = meta.icon;

  return (
    <header className="px-5 py-3 border-b border-border bg-muted/30 flex items-start gap-3">
      <Icon className="h-5 w-5 text-primary shrink-0 mt-0.5" />
      <div className="flex-1 min-w-0">
        <h2 className="text-base font-heading font-semibold docs-h3">{issue.title}</h2>
        <p className="text-xs font-mono text-muted-foreground mt-0.5">{meta.label}</p>
      </div>
      <CopyLinkButton issueId={issue.id} />
      {issue.fixCommand && <CopyFixButton command={issue.fixCommand} altCommand={issue.altCommand} />}
    </header>
  );
};

const IssueRelatedLinks = ({ related }: IssueRelatedLinksProps) => {
  const hasRelated = Boolean(related && related.length > 0);

  if (hasRelated === false) return null;

  return (
    <div className="pt-2 border-t border-border flex flex-wrap items-center gap-3 text-xs">
      <span className="text-muted-foreground font-mono uppercase tracking-wider">Related</span>
      {related!.map((linkItem) => (
        <a key={linkItem.href} href={linkItem.href} className="text-primary hover:underline font-mono">
          {linkItem.label}
        </a>
      ))}
    </div>
  );
};

const IssueFixSection = ({ issue }: IssueProps) => {
  return (
    <div>
      <h3 className="text-xs font-mono uppercase tracking-wider text-muted-foreground mb-1">Fix</h3>
      <p className="text-sm text-foreground/90 mb-2">{issue.fix}</p>
      {issue.fixCommand && <CodeBlock language={issue.fixLanguage ?? "bash"} code={issue.fixCommand} title="Run" />}
      {issue.altCommand && (
        <div className="mt-2">
          <p className="text-xs font-mono text-muted-foreground mb-1">{issue.altLabel ?? "Alternative"}</p>
          <CodeBlock language="bash" code={issue.altCommand} title="Alt" />
        </div>
      )}
    </div>
  );
};

export const TroubleshootingIssueCard = ({ issue }: IssueProps) => {
  return (
    <article id={issue.id} className="rounded-lg border border-border bg-card overflow-hidden">
      <IssueHeader issue={issue} />
      <div className="p-5 space-y-4">
        <div>
          <h3 className="text-xs font-mono uppercase tracking-wider text-muted-foreground mb-1">Symptom</h3>
          <pre className="text-sm font-mono bg-muted/40 border border-border rounded p-3 overflow-x-auto">{issue.symptom}</pre>
        </div>
        <div>
          <h3 className="text-xs font-mono uppercase tracking-wider text-muted-foreground mb-1">Cause</h3>
          <p className="text-sm text-foreground/90">{issue.cause}</p>
        </div>
        <IssueFixSection issue={issue} />
        <IssueRelatedLinks related={issue.related} />
      </div>
    </article>
  );
};

export default TroubleshootingIssueCard;
