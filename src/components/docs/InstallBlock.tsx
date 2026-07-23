import { useState, useCallback } from "react";
import { Copy, Check, Terminal } from "lucide-react";
import { copyToClipboard } from "@/lib/clipboard";

interface InstallTab {
  label: string;
  command: string;
}

interface InstallBlockProps {
  command?: string;
  tabs?: InstallTab[];
}

const CopyLine = ({ command }: { command: string }) => {
  const [copied, setCopied] = useState(false);

  const handleCopy = useCallback(async () => {
    await copyToClipboard(command);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  }, [command]);

  return (
    <div
      onClick={handleCopy}
      className="group flex cursor-pointer items-start gap-3 rounded-sm border border-border bg-terminal px-5 py-3 transition-colors hover:border-primary/40"
    >
      <Terminal className="h-4 w-4 text-primary mt-0.5 shrink-0" />
      <code className="font-mono text-sm text-terminal-foreground break-all leading-relaxed flex-1">
        {command}
      </code>
      <span className="text-muted-foreground group-hover:text-foreground transition-colors shrink-0 mt-0.5">
        {copied ? <Check className="h-4 w-4 text-primary" /> : <Copy className="h-4 w-4" />}
      </span>
    </div>
  );
};

const InstallBlock = ({ command, tabs }: InstallBlockProps) => {
  const [active, setActive] = useState(0);

  if (tabs && tabs.length > 0) {
    return (
      <div className="space-y-3">
        <div className="inline-flex justify-center gap-1 rounded-md border border-border bg-secondary/70 p-1 mx-auto">
          {tabs.map((tab, i) => (
            <button
              key={tab.label}
              onClick={() => setActive(i)}
              data-active={i === active}
              className={`btn-slide ${
                i === active ? "" : "btn-slide-ghost"
              } px-4 py-1.5 rounded-md text-sm font-sans font-medium transition-all duration-300 ${
                i === active
                  ? "shadow-sm"
                  : "border border-transparent bg-transparent text-muted-foreground hover:border-border hover:bg-card hover:text-foreground"
              }`}
              style={
                i === active
                  ? {
                      borderWidth: "1px",
                      borderStyle: "solid",
                      borderColor: "hsl(var(--primary) / 0.6)",
                      backgroundColor: "hsl(var(--primary) / 0.12)",
                      color: "hsl(var(--primary))",
                    }
                  : undefined
              }
            >
              {tab.label}
            </button>
          ))}
        </div>
        <CopyLine command={tabs[active].command} />
      </div>
    );
  }

  if (command) {
    return <CopyLine command={command} />;
  }

  return null;
};

export default InstallBlock;
