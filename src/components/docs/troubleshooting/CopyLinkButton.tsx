import { useState, useCallback } from "react";
import { Check, Link2 } from "lucide-react";
import { copyToClipboard } from "@/lib/clipboard";

interface CopyLinkButtonProps {
  issueId: string;
}

export const CopyLinkButton = ({ issueId }: CopyLinkButtonProps) => {
  const [isCopied, setIsCopied] = useState(false);

  const handleCopy = useCallback(async () => {
    const currentUrl = new URL(window.location.href);
    currentUrl.searchParams.set("id", issueId);
    await copyToClipboard(currentUrl.toString());
    setIsCopied(true);
    window.setTimeout(() => setIsCopied(false), 2000);
  }, [issueId]);

  const buttonClasses = isCopied
    ? "border-primary bg-primary/15 text-foreground dark:bg-primary/20 dark:text-primary"
    : "border-border bg-background text-muted-foreground hover:text-foreground hover:border-foreground/40";

  return (
    <button
      type="button"
      onClick={handleCopy}
      aria-label={isCopied ? "Link copied" : "Copy link to this issue"}
      title={isCopied ? "Link copied!" : "Copy link to this issue"}
      className={`shrink-0 inline-flex items-center gap-1.5 px-2.5 py-1.5 rounded-md text-xs font-mono border transition-colors ${buttonClasses}`}
    >
      {isCopied ? <Check className="h-3.5 w-3.5" /> : <Link2 className="h-3.5 w-3.5" />}
      {isCopied ? "Linked" : "Link"}
    </button>
  );
};

export default CopyLinkButton;
