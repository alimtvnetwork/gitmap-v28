import { useEffect, useMemo, useState, Dispatch, SetStateAction } from "react";
import { useNavigate } from "react-router-dom";
import { Search } from "lucide-react";
import { DocsTooltip } from "@/components/docs/DocsTooltip";
import {
  CommandDialog,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from "@/components/ui/command";
import { commands, CommandDef } from "@/data/commands";

interface CommandPaletteItem {
  id: string;
  haystack: string;
  name: string;
  alias?: string;
  description: string;
}

interface CommandPaletteRowProps {
  item: CommandPaletteItem;
  onSelect: (commandName: string) => void;
}

interface CommandPaletteDialogProps {
  open: boolean;
  onOpenChange: Dispatch<SetStateAction<boolean>>;
  items: CommandPaletteItem[];
  onSelect: (commandName: string) => void;
}

interface CommandPaletteTriggerProps {
  onClick: () => void;
}

function buildHaystack(commandItem: CommandDef): string {
  const exampleCommands = commandItem.examples?.map((exampleItem) => exampleItem.command) ?? [];
  const parts = [commandItem.name, commandItem.alias ?? "", commandItem.description, ...exampleCommands];

  return parts.join(" ").toLowerCase();
}

function buildPaletteItem(commandItem: CommandDef): CommandPaletteItem {
  return {
    id: commandItem.name,
    haystack: buildHaystack(commandItem),
    name: commandItem.name,
    alias: commandItem.alias,
    description: commandItem.description,
  };
}

function useCommandPaletteShortcut(setOpen: Dispatch<SetStateAction<boolean>>): void {
  useEffect(() => {
    const handleKeyDown = (event: KeyboardEvent) => {
      const isKeyK = event.key === "k" || event.key === "K";
      const hasModifier = Boolean(event.metaKey || event.ctrlKey);
      if (isKeyK && hasModifier) {
        event.preventDefault();
        setOpen((isOpen) => isOpen === false);
      }
    };
    window.addEventListener("keydown", handleKeyDown);

    return () => window.removeEventListener("keydown", handleKeyDown);
  }, [setOpen]);
}

const CommandPaletteTrigger = ({ onClick }: CommandPaletteTriggerProps) => {
  const triggerClasses = "docs-focus-ring inline-flex h-7 items-center gap-2 rounded-sm border border-sidebar-border bg-sidebar-accent/60 px-2 text-xs text-muted-foreground hover:bg-sidebar-accent hover:text-sidebar-accent-foreground";

  return (
    <DocsTooltip label="Search commands (⌘K)">
      <button type="button" aria-label="Open command palette (search commands, flags, pages)" onClick={onClick} className={triggerClasses}>
        <Search className="h-3.5 w-3.5" aria-hidden="true" />
        <span className="hidden sm:inline font-mono">Search...</span>
        <kbd className="pointer-events-none hidden h-4 select-none items-center gap-0.5 rounded border border-border bg-muted px-1.5 font-mono text-[10px] font-medium opacity-100 sm:flex">⌘K</kbd>
      </button>
    </DocsTooltip>
  );
};

const CommandPaletteRow = ({ item, onSelect }: CommandPaletteRowProps) => {
  return (
    <CommandItem value={item.haystack} onSelect={() => onSelect(item.name)}>
      <div className="flex flex-col gap-0.5">
        <div className="flex items-baseline gap-2 font-mono">
          <span className="text-foreground">{item.name}</span>
          {item.alias && <span className="text-xs text-muted-foreground">({item.alias})</span>}
        </div>
        <span className="line-clamp-1 text-xs text-muted-foreground">{item.description}</span>
      </div>
    </CommandItem>
  );
};

const CommandPaletteDialog = ({ open, onOpenChange, items, onSelect }: CommandPaletteDialogProps) => {
  return (
    <CommandDialog open={open} onOpenChange={onOpenChange}>
      <CommandInput placeholder="Search commands, aliases, examples…  (⌘K)" />
      <CommandList>
        <CommandEmpty>No commands match.</CommandEmpty>
        <CommandGroup heading="Commands">
          {items.map((paletteItem) => (
            <CommandPaletteRow key={paletteItem.id} item={paletteItem} onSelect={onSelect} />
          ))}
        </CommandGroup>
      </CommandList>
    </CommandDialog>
  );
};

const CommandPalette = () => {
  const [open, setOpen] = useState(false);
  const navigate = useNavigate();
  useCommandPaletteShortcut(setOpen);

  const items = useMemo(() => commands.map(buildPaletteItem), []);

  const handleSelect = (commandName: string) => {
    setOpen(false);
    navigate(`/commands?cmd=${encodeURIComponent(commandName)}`);
  };

  return (
    <>
      <CommandPaletteTrigger onClick={() => setOpen(true)} />
      <CommandPaletteDialog open={open} onOpenChange={setOpen} items={items} onSelect={handleSelect} />
    </>
  );
};

export default CommandPalette;
