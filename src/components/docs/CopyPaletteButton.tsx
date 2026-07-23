import { useState } from "react";
import { Check, ClipboardCopy } from "lucide-react";
import { useTheme } from "@/hooks/useTheme";
import { DocsTooltip } from "@/components/docs/DocsTooltip";

// CSS custom properties that make up the active VS Code-inspired palette.
// Read live from the document so the snapshot always reflects what's rendered.
const PALETTE_VARS = [
  "--background",
  "--foreground",
  "--card",
  "--card-foreground",
  "--popover",
  "--popover-foreground",
  "--primary",
  "--primary-foreground",
  "--secondary",
  "--secondary-foreground",
  "--muted",
  "--muted-foreground",
  "--accent",
  "--accent-foreground",
  "--destructive",
  "--destructive-foreground",
  "--border",
  "--input",
  "--ring",
  "--sidebar-background",
  "--sidebar-foreground",
  "--sidebar-primary",
  "--sidebar-primary-foreground",
  "--sidebar-accent",
  "--sidebar-accent-foreground",
  "--sidebar-border",
  "--sidebar-ring",
] as const;

const COPY_FEEDBACK_MS = 1500;

function buildPaletteCss(themeName: string, selector: string): string {
  const styles = getComputedStyle(document.documentElement);
  const lines = PALETTE_VARS
    .map((name) => {
      const value = styles.getPropertyValue(name).trim();
      return value ? `  ${name}: ${value};` : null;
    })
    .filter((line): line is string => line !== null);
  return `/* ${themeName} */\n${selector} {\n${lines.join("\n")}\n}\n`;
}

export function CopyPaletteButton() {
  const { isDark } = useTheme();
  const [copied, setCopied] = useState(false);
  const themeLabel = isDark ? "Dark" : "Light";
  const selector = isDark ? ".dark" : ".light";

  const handleCopy = async () => {
    try {
      const css = buildPaletteCss(themeLabel, selector);
      await navigator.clipboard.writeText(css);
      setCopied(true);
      window.setTimeout(() => setCopied(false), COPY_FEEDBACK_MS);
    } catch (error) {
      // Surface failure to the console — clipboard can be blocked by permissions.
      console.error("[copy-palette] failed to write clipboard:", error);
    }
  };

  const tooltipText = copied
    ? "Copied!"
    : `Copy ${themeLabel} theme palette (CSS variables) to clipboard`;

  return (
    <DocsTooltip label={tooltipText}>
      <button
        type="button"
        onClick={handleCopy}
        aria-label={`Copy ${themeLabel} theme palette to clipboard`}
        className="docs-focus-ring inline-flex h-7 w-7 items-center justify-center rounded-sm border border-border bg-card text-muted-foreground shadow-sm transition-colors hover:text-foreground"
      >
        {copied ? (
          <Check className="h-3.5 w-3.5" aria-hidden="true" />
        ) : (
          <ClipboardCopy className="h-3.5 w-3.5" aria-hidden="true" />
        )}
      </button>
    </DocsTooltip>
  );
}

export default CopyPaletteButton;
