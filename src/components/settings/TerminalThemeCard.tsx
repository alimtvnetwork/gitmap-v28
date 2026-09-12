import { Terminal } from "lucide-react";
import { TerminalThemeType } from "./types";

interface TerminalThemeCardProps {
  terminalTheme: string;
  onTerminalThemeChange: (theme: string) => void;
}

const THEME_OPTIONS = [
  { id: TerminalThemeType.Dark, label: "Dark+" },
  { id: TerminalThemeType.Ubuntu, label: "Ubuntu Purple" },
  { id: TerminalThemeType.Win, label: "PowerShell Navy" },
];

export function TerminalThemeCard({
  terminalTheme,
  onTerminalThemeChange,
}: TerminalThemeCardProps) {
  return (
    <div className="rounded-xl border border-border bg-card p-6 shadow-sm space-y-4">
      <div className="flex items-center gap-2 text-foreground font-semibold text-lg">
        <Terminal className="h-5 w-5 text-primary" />
        <h2>Interactive UI Terminal</h2>
      </div>
      <p className="text-sm text-muted-foreground">
        Select the color palette theme for the integrated web terminal drawer.
      </p>
      <div className="grid grid-cols-3 gap-3">
        {THEME_OPTIONS.map((themeOpt) => {
          const isSelected = terminalTheme === themeOpt.id;

          return (
            <button
              key={themeOpt.id}
              type="button"
              onClick={() => onTerminalThemeChange(themeOpt.id)}
              className={`px-4 py-3 rounded-lg border text-sm font-medium transition-colors ${
                isSelected
                  ? "border-primary bg-primary/10 text-foreground font-semibold"
                  : "border-border bg-background text-muted-foreground hover:text-foreground"
              }`}
            >
              {themeOpt.label}
            </button>
          );
        })}
      </div>
    </div>
  );
}

export default TerminalThemeCard;
