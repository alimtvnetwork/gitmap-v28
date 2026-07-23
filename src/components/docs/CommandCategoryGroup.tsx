import { useState, useEffect } from "react";
import { ChevronDown, ChevronRight } from "lucide-react";
import { motion, AnimatePresence } from "framer-motion";
import CommandCard from "./CommandCard";
import type { CommandDef } from "@/data/commands";

interface Props {
  label: string;
  description: string;
  icon?: string;
  commands: CommandDef[];
  defaultOpen?: boolean;
  forceOpen?: boolean;
  onNavigate?: (commandName: string) => void;
  commandRefs?: React.MutableRefObject<Record<string, HTMLDivElement | null>>;
}

const CommandCategoryGroup = ({ label, description, icon, commands, defaultOpen = true, forceOpen, onNavigate, commandRefs }: Props) => {
  const [open, setOpen] = useState(defaultOpen);

  useEffect(() => {
    if (forceOpen) setOpen(true);
  }, [forceOpen]);

  return (
    <div className="rounded-lg border border-border overflow-hidden">
      <button
        onClick={() => setOpen(!open)}
        className="w-full flex items-center gap-3 px-4 py-3 bg-muted/30 hover:bg-muted/50 transition-colors text-left"
      >
        {open ? (
          <ChevronDown className="h-4 w-4 text-primary shrink-0" />
        ) : (
          <ChevronRight className="h-4 w-4 text-muted-foreground shrink-0" />
        )}
        <div className="flex-1 min-w-0 flex items-center gap-2">
          {icon && <span className="text-base">{icon}</span>}
          <span className="font-sans font-semibold text-sm text-foreground">{label}</span>
          <span className="text-xs text-muted-foreground">({commands.length})</span>
        </div>
        <span className="text-xs text-muted-foreground truncate hidden sm:inline">{description}</span>
      </button>

      <AnimatePresence initial={false}>
        {open && (
          <motion.div
            initial={{ height: 0, opacity: 0 }}
            animate={{ height: "auto", opacity: 1 }}
            exit={{ height: 0, opacity: 0 }}
            transition={{ duration: 0.2, ease: "easeInOut" }}
            className="overflow-hidden"
          >
            <div className="p-2 space-y-1.5">
              {commands.map((cmd) => (
                <div key={cmd.name} ref={(el) => { if (commandRefs) commandRefs.current[cmd.name] = el; }}>
                  <CommandCard {...cmd} onNavigate={onNavigate} />
                </div>
              ))}
            </div>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
};

export default CommandCategoryGroup;
