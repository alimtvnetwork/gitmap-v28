import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";
import { Link } from "react-router-dom";
import { ArrowLeft, FileText } from "lucide-react";
import DocsLayout from "@/components/docs/DocsLayout";

interface SpecPageProps {
  /** Page title shown as H1 above the rendered markdown. */
  title: string;
  /** Short subtitle / lead paragraph. */
  subtitle?: string;
  /** Repository-relative path of the source markdown (e.g. "spec/01-app/100-scan-all.md"). */
  sourcePath: string;
  /** Raw markdown content (import via `?raw`). */
  markdown: string;
  /** Optional related links shown in a footer card. */
  relatedLinks?: { label: string; to: string; description?: string }[];
}

const SpecPage = ({ title, subtitle, sourcePath, markdown, relatedLinks }: SpecPageProps) => {
  return (
    <DocsLayout>
      <div className="mb-6">
        <Link
          to="/spec"
          className="inline-flex items-center gap-1.5 text-sm font-sans text-muted-foreground hover:text-primary transition-colors"
        >
          <ArrowLeft className="h-3.5 w-3.5" />
          Back to spec index
        </Link>
      </div>

      <h1 className="text-3xl font-heading font-bold docs-h1 mb-2">{title}</h1>
      {subtitle && <p className="text-muted-foreground mb-2">{subtitle}</p>}
      <p className="inline-flex items-center gap-1.5 text-xs font-mono text-muted-foreground mb-6">
        <FileText className="h-3 w-3" />
        Source:&nbsp;<code className="docs-inline-code">{sourcePath}</code>
      </p>

      <article className="docs-prose prose prose-sm max-w-none">
        <ReactMarkdown
          remarkPlugins={[remarkGfm]}
          components={{
            h1: ({ children }) => <h2 className="docs-h2 text-2xl font-heading font-bold mt-8 mb-3">{children}</h2>,
            h2: ({ children }) => <h2 className="docs-h2 text-xl font-heading font-bold mt-6 mb-2">{children}</h2>,
            h3: ({ children }) => <h3 className="docs-h3 text-base font-heading font-semibold mt-5 mb-2">{children}</h3>,
            h4: ({ children }) => <h4 className="text-sm font-heading font-semibold mt-4 mb-2 text-foreground">{children}</h4>,
            p:  ({ children }) => <p className="text-sm font-sans text-foreground leading-relaxed mb-3">{children}</p>,
            ul: ({ children }) => <ul className="list-disc pl-5 space-y-1 mb-3 text-sm font-sans text-foreground">{children}</ul>,
            ol: ({ children }) => <ol className="list-decimal pl-5 space-y-1 mb-3 text-sm font-sans text-foreground">{children}</ol>,
            li: ({ children }) => <li className="leading-relaxed">{children}</li>,
            blockquote: ({ children }) => (
              <blockquote className="border-l-4 border-primary/40 bg-muted/40 pl-4 py-2 my-3 text-sm font-sans text-muted-foreground italic">
                {children}
              </blockquote>
            ),
            code: ({ className, children, ...props }) => {
              const isInline = !className;
              if (isInline) return <code className="docs-inline-code">{children}</code>;
              return (
                <code className={`${className ?? ""} font-mono text-xs`} {...props}>
                  {children}
                </code>
              );
            },
            pre: ({ children }) => (
              <pre className="bg-[hsl(var(--code-bg))] border border-border/50 rounded-lg p-4 overflow-x-auto text-xs font-mono mb-4 leading-relaxed">
                {children}
              </pre>
            ),
            table: ({ children }) => (
              <div className="overflow-x-auto mb-4 border border-border rounded-lg">
                <table className="w-full text-sm font-sans">{children}</table>
              </div>
            ),
            thead: ({ children }) => <thead className="bg-muted/50">{children}</thead>,
            th: ({ children }) => <th className="text-left px-3 py-2 font-heading font-semibold text-xs text-foreground border-b border-border">{children}</th>,
            td: ({ children }) => <td className="px-3 py-2 text-foreground border-b border-border/50 align-top">{children}</td>,
            a: ({ href, children }) => {
              const internal = href?.startsWith("/") ?? false;
              if (internal) return <Link to={href!} className="text-primary hover:underline">{children}</Link>;
              return (
                <a href={href} target="_blank" rel="noreferrer" className="text-primary hover:underline">
                  {children}
                </a>
              );
            },
            hr: () => <hr className="my-6 border-border" />,
            strong: ({ children }) => <strong className="font-heading font-semibold text-foreground">{children}</strong>,
          }}
        >
          {markdown}
        </ReactMarkdown>
      </article>

      {relatedLinks && relatedLinks.length > 0 && (
        <div className="mt-10 border border-border rounded-lg p-4 bg-card">
          <h4 className="text-xs font-mono font-semibold text-muted-foreground uppercase tracking-wider mb-3">
            Related
          </h4>
          <ul className="space-y-2">
            {relatedLinks.map((link) => (
              <li key={link.to}>
                <Link to={link.to} className="text-sm font-sans text-primary hover:underline">
                  {link.label}
                </Link>
                {link.description && (
                  <span className="text-sm font-sans text-muted-foreground"> — {link.description}</span>
                )}
              </li>
            ))}
          </ul>
        </div>
      )}
    </DocsLayout>
  );
};

export default SpecPage;
