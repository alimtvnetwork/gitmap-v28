import { useState, useCallback } from "react";
import { Copy, Check } from "lucide-react";
import { copyToClipboard } from "@/lib/clipboard";

interface CopyFixButtonProps {
  command: string;
  altCommand?: string;
}

export const CopyFixButton = ({ command, altCommand }: CopyFixButtonProps) => {
  const [isCopied, setIsCopied] = useState(false);

  const handleCopy = useCallback(async () => {
    const payload = altCommand ? `${command}\n\n# Alternative\n${altCommand}` : command;
    await copyToClipboard(payload);
    setIsCopied(true);
    window.setTimeout(() => setIsCopied(false), 2000);
  }, [command, altCommand]);

  const buttonClasses = isCopied
    ? "border-primary bg-primary/15 text-foreground dark:bg-primary/20 dark:text-primary"
    : "border-border bg-background text-muted-foreground hover:text-foreground hover:border-foreground/40";

  return (
    <button
      type="button"
      onClick={handleCopy}
      aria-label={isCopied ? "Fix command copied" : "Copy fix command"}
      title={isCopied ? "Copied!" : "Copy fix command"}
      className={`shrink-0 inline-flex items-center gap-1.5 px-2.5 py-1.5 rounded-md text-xs font-mono border transition-colors ${buttonClasses}`}
    >
      {isCopied ? <Check className="h-3.5 w-3.5" /> : <Copy className="h-3.5 w-3.5" />}
      {isCopied ? "Copied" : "Copy fix"}
    </button>
  );
};

export default CopyFixButton;
