import { useState } from "react";
import { ChevronDown, ChevronRight, ArrowRight, ExternalLink } from "lucide-react";
import { Link } from "react-router-dom";
import CodeBlock from "./CodeBlock";
import type { CommandSeeAlso } from "@/data/commands";

interface CommandCardProps {
  name: string;
  alias?: string;
  description: string;
  usage?: string;
  flags?: { flag: string; description: string }[];
  examples?: { command: string; description?: string }[];
  howToProceed?: { step: number; title: string; action: string }[];
  notes?: string[];
  seeAlso?: CommandSeeAlso[];
  onNavigate?: (commandName: string) => void;
}

const CommandCard = ({ name, alias, description, usage, flags, examples, howToProceed, notes, seeAlso, onNavigate }: CommandCardProps) => {
  const [open, setOpen] = useState(false);

  return (
    <div className="border border-border rounded-lg overflow-hidden transition-colors hover:border-primary/40">
      <button
        onClick={() => setOpen(!open)}
        className="w-full flex items-center gap-3 px-4 py-3 text-left hover:bg-muted/50 transition-colors"
      >
        {open ? (
          <ChevronDown className="h-4 w-4 text-primary shrink-0" />
        ) : (
          <ChevronRight className="h-4 w-4 text-muted-foreground shrink-0" />
        )}
        <div className="flex items-center gap-2 shrink-0">
          <code className="font-sans font-semibold text-sm text-foreground dark:text-white whitespace-nowrap">{name}</code>
          {alias && (
            <span className="text-xs font-sans font-medium text-foreground bg-primary/10 border border-primary/20 px-1.5 py-0.5 rounded whitespace-nowrap dark:bg-primary/20 dark:text-white dark:border-primary/40">{alias}</span>
          )}
        </div>
        <span className="text-sm text-muted-foreground dark:text-slate-300 truncate min-w-0 flex-1">{description}</span>
      </button>

      {open && (
        <div className="px-4 pb-4 border-t border-border pt-3 space-y-3">
          {usage && <CodeBlock code={usage} />}

          {flags && flags.length > 0 && (
            <div>
              <h4 className="text-xs font-mono font-semibold text-muted-foreground dark:text-slate-300 uppercase tracking-wider mb-2">Flags</h4>
              <div className="space-y-1">
                {flags.map((f) => (
                  <div key={f.flag} className="flex gap-4 text-sm">
                    <code className="font-mono text-foreground font-medium dark:text-white whitespace-nowrap">{f.flag}</code>
                    <span className="text-muted-foreground dark:text-slate-300">{f.description}</span>
                  </div>
                ))}
              </div>
            </div>
          )}

          {howToProceed && howToProceed.length > 0 && (
            <div className="rounded-md border border-primary/25 bg-primary/5 p-3.5 space-y-2.5">
              <h4 className="text-xs font-mono font-bold text-primary uppercase tracking-wider flex items-center gap-1.5">
                <span>🚀 How to Proceed</span>
              </h4>
              <div className="space-y-3">
                {howToProceed.map((step) => (
                  <div key={step.step} className="text-sm">
                    <div className="flex items-center gap-2 font-medium text-foreground dark:text-white">
                      <span className="flex h-5 w-5 shrink-0 items-center justify-center rounded-full bg-primary text-primary-foreground text-xs font-bold">
                        {step.step}
                      </span>
                      <span>{step.title}</span>
                    </div>
                    {step.action && (
                      <div className="ml-7 mt-1.5">
                        <CodeBlock code={step.action} />
                      </div>
                    )}
                  </div>
                ))}
              </div>
            </div>
          )}

          {notes && notes.length > 0 && (
            <div className="rounded-md border border-amber-500/30 bg-amber-500/5 p-3.5 space-y-1.5">
              <h4 className="text-xs font-mono font-bold text-amber-600 dark:text-amber-400 uppercase tracking-wider flex items-center gap-1.5">
                <span>⚠️ Important Notes & Safety</span>
              </h4>
              <ul className="list-disc list-inside space-y-1 text-xs text-muted-foreground dark:text-slate-300">
                {notes.map((note, idx) => (
                  <li key={idx}>{note}</li>
                ))}
              </ul>
            </div>
          )}

          {examples && examples.length > 0 && (
            <div>
              <h4 className="text-xs font-mono font-semibold text-muted-foreground uppercase tracking-wider mb-2">Examples</h4>
              {examples.map((ex, i) => (
                <div key={i}>
                  {ex.description && <p className="text-sm text-muted-foreground mb-1">{ex.description}</p>}
                  <CodeBlock code={ex.command} />
                </div>
              ))}
            </div>
          )}

          {seeAlso && seeAlso.length > 0 && (
            <div>
              <h4 className="text-xs font-mono font-semibold text-muted-foreground uppercase tracking-wider mb-2">See Also</h4>
              <div className="flex flex-wrap gap-2">
                {seeAlso.map((ref) =>
                  ref.url ? (
                    <Link
                      key={ref.name}
                      to={ref.url}
                      onClick={(e) => e.stopPropagation()}
                      className="group inline-flex items-center gap-1.5 px-2.5 py-1.5 rounded-md border border-border bg-card text-sm font-mono text-foreground hover:border-primary/60 hover:bg-primary/5 focus:outline-none focus-visible:ring-2 focus-visible:ring-primary/50 focus-visible:border-primary/60 transition-colors"
                      title={ref.description}
                    >
                      <span>{ref.name}</span>
                      <ExternalLink className="h-3 w-3 text-muted-foreground group-hover:text-primary transition-colors" />
                    </Link>
                  ) : (
                    <button
                      key={ref.name}
                      onClick={(e) => {
                        e.stopPropagation();
                        onNavigate?.(ref.name);
                      }}
                      className="group inline-flex items-center gap-1.5 px-2.5 py-1.5 rounded-md border border-border bg-card text-sm font-mono text-foreground hover:border-primary/60 hover:bg-primary/5 focus:outline-none focus-visible:ring-2 focus-visible:ring-primary/50 focus-visible:border-primary/60 transition-colors"
                      title={ref.description}
                    >
                      <span>{ref.name}</span>
                      <ArrowRight className="h-3 w-3 text-muted-foreground group-hover:text-primary transition-colors" />
                    </button>
                  )
                )}
              </div>
            </div>
          )}
        </div>
      )}
    </div>
  );
};

export default CommandCard;
