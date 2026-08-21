import { useState, useCallback } from "react";
import { Copy, Check, Stethoscope, Terminal, FileText, Wrench, ListChecks, LucideIcon } from "lucide-react";
import { copyToClipboard } from "@/lib/clipboard";

interface DiagnosticStep {
  icon: LucideIcon;
  title: string;
  body: string;
  command: string;
}

interface ChecklistCommandProps {
  command: string;
}

interface DiagnosticStepItemProps {
  step: DiagnosticStep;
  index: number;
}

const DIAGNOSTIC_STEPS: DiagnosticStep[] = [
  {
    icon: Stethoscope,
    title: "Run a full health check",
    body: "Surfaces PATH, config, lock, network, and version-mismatch issues in one pass. Most problems are diagnosed (and named) here.",
    command: "gitmap doctor",
  },
  {
    icon: Wrench,
    title: "Apply the auto-fixer if doctor flagged PATH",
    body: "If you see 'PATH binary version mismatch' or 'gitmap not found on PATH', let doctor sync the active binary from the deployed copy.",
    command: "gitmap doctor --fix-path",
  },
  {
    icon: Terminal,
    title: "Re-run the failing command with --verbose",
    body: "Verbose mode prints every git invocation, resolved paths, and timing — and writes a timestamped debug log next to the output.",
    command: "gitmap <your-failing-command> --verbose",
  },
  {
    icon: FileText,
    title: "Inspect the verbose log",
    body: "Logs land in .gitmap/logs/. Open the newest one and scan for the first 'Error:' or non-zero exit — that is almost always the root cause.",
    command: "# Windows\nGet-ChildItem .gitmap\\logs | Sort-Object LastWriteTime | Select-Object -Last 1\n# Unix\nls -t .gitmap/logs | head -1",
  },
];

const ChecklistCommand = ({ command }: ChecklistCommandProps) => {
  const [isCopied, setIsCopied] = useState(false);

  const handleCopy = useCallback(async () => {
    await copyToClipboard(command);
    setIsCopied(true);
    window.setTimeout(() => setIsCopied(false), 2000);
  }, [command]);

  const buttonClasses = isCopied
    ? "border-primary bg-primary/15 text-foreground dark:bg-primary/20 dark:text-primary"
    : "border-border bg-background text-muted-foreground hover:text-foreground hover:border-foreground/40";

  return (
    <div className="relative group">
      <pre className="text-xs font-mono bg-background border border-border rounded p-2.5 pr-10 overflow-x-auto whitespace-pre">
        {command}
      </pre>
      <button
        type="button"
        onClick={handleCopy}
        aria-label={isCopied ? "Copied" : "Copy command"}
        title={isCopied ? "Copied!" : "Copy command"}
        className={`absolute top-1.5 right-1.5 inline-flex items-center justify-center w-7 h-7 rounded border transition-colors ${buttonClasses}`}
      >
        {isCopied ? <Check className="h-3.5 w-3.5" /> : <Copy className="h-3.5 w-3.5" />}
      </button>
    </div>
  );
};

const DiagnosticStepItem = ({ step, index }: DiagnosticStepItemProps) => {
  const Icon = step.icon;

  return (
    <li className="flex gap-3">
      <span
        aria-hidden="true"
        className="shrink-0 inline-flex items-center justify-center w-7 h-7 rounded-full bg-primary/15 text-foreground border border-primary/25 dark:bg-primary/20 dark:text-primary font-mono text-sm font-semibold"
      >
        {index + 1}
      </span>
      <div className="flex-1 min-w-0">
        <div className="flex items-center gap-2 mb-1">
          <Icon className="h-4 w-4 text-muted-foreground shrink-0" />
          <h3 className="text-sm font-semibold text-foreground">{step.title}</h3>
        </div>
        <p className="text-sm text-foreground/80 mb-2">{step.body}</p>
        <ChecklistCommand command={step.command} />
      </div>
    </li>
  );
};

export const DiagnosticChecklist = () => {
  return (
    <section aria-labelledby="diagnostic-checklist-heading" className="mb-8 rounded-lg border border-border bg-muted/20 p-5">
      <header className="flex items-center gap-2 mb-1">
        <ListChecks className="h-5 w-5 text-primary" />
        <h2 id="diagnostic-checklist-heading" className="text-base font-heading font-semibold docs-h3">
          Diagnostic checklist
        </h2>
      </header>
      <p className="text-sm text-muted-foreground mb-4">
        Run these four steps in order before searching the catalog below. Most issues are identified and resolved by step 2.
      </p>
      <ol className="space-y-4">
        {DIAGNOSTIC_STEPS.map((step, stepIndex) => (
          <DiagnosticStepItem key={step.title} step={step} index={stepIndex} />
        ))}
      </ol>
    </section>
  );
};

export default DiagnosticChecklist;
