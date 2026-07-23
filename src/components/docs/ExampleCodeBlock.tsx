import { useState } from "react";
import { Check, Copy, Terminal } from "lucide-react";
import { copyToClipboard } from "@/lib/clipboard";

interface ExampleCodeBlockProps {
  command: string;
  description?: string;
  /** Hide the "Run in terminal" hint (e.g. when the snippet is non-runnable output). */
  hideRunHint?: boolean;
}

/**
 * Renders a fenced example with a copy-to-clipboard button and a
 * "Run in terminal" hint. Used by every command card and helptext
 * embed so every example block has uniform affordances.
 */
const ExampleCodeBlock = ({ command, description, hideRunHint }: ExampleCodeBlockProps) => {
  const [copied, setCopied] = useState(false);

  const onCopy = async () => {
    await copyToClipboard(command);
    setCopied(true);
    setTimeout(() => setCopied(false), 1400);
  };

  return (
    <div className="group relative rounded-md border border-border bg-muted/40">
      <div className="flex items-center justify-between gap-2 px-3 py-1.5 border-b border-border/60">
        {!hideRunHint && (
          <div className="flex items-center gap-1.5 text-[10px] font-mono uppercase tracking-wide text-muted-foreground">
            <Terminal className="h-3 w-3" aria-hidden />
            Run in terminal
          </div>
        )}
        <button
          type="button"
          onClick={onCopy}
          className="ml-auto inline-flex items-center gap-1 rounded px-2 py-0.5 text-[10px] font-mono text-muted-foreground hover:bg-muted hover:text-foreground transition-colors"
          aria-label={copied ? "Copied" : "Copy command to clipboard"}
        >
          {copied ? <Check className="h-3 w-3 text-primary" /> : <Copy className="h-3 w-3" />}
          {copied ? "Copied" : "Copy"}
        </button>
      </div>
      <pre className="overflow-x-auto px-3 py-2 text-xs font-mono text-foreground">
        <code>{command}</code>
      </pre>
      {description && (
        <p className="px-3 pb-2 text-xs text-muted-foreground leading-snug">{description}</p>
      )}
    </div>
  );
};

export default ExampleCodeBlock;
