import { useEffect, useMemo, useState } from "react";
import { useNavigate } from "react-router-dom";
import {
  CommandDialog,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from "@/components/ui/command";
import { commands } from "@/data/commands";

/**
 * Global ⌘K / Ctrl+K command palette.
 *
 * Fuzzy-searches every command in `src/data/commands.ts` by name,
 * alias, description, and any example string — so users can jump
 * straight to a usage without scrolling the Commands page.
 */
const CommandPalette = () => {
  const [open, setOpen] = useState(false);
  const navigate = useNavigate();

  useEffect(() => {
    const handler = (e: KeyboardEvent) => {
      const isToggle = (e.key === "k" || e.key === "K") && (e.metaKey || e.ctrlKey);
      if (isToggle) {
        e.preventDefault();
        setOpen((o) => !o);
      }
    };
    window.addEventListener("keydown", handler);
    return () => window.removeEventListener("keydown", handler);
  }, []);

  const items = useMemo(
    () =>
      commands.map((c) => ({
        id: c.name,
        // Concat haystack so cmdk's built-in fuzzy hits name, alias, and examples.
        haystack: [
          c.name,
          c.alias ?? "",
          c.description,
          ...(c.examples?.map((e) => e.command) ?? []),
        ]
          .join(" ")
          .toLowerCase(),
        name: c.name,
        alias: c.alias,
        description: c.description,
      })),
    [],
  );

  const onSelect = (commandName: string) => {
    setOpen(false);
    navigate(`/commands?cmd=${encodeURIComponent(commandName)}`);
  };

  return (
    <CommandDialog open={open} onOpenChange={setOpen}>
      <CommandInput placeholder="Search commands, aliases, examples…  (⌘K)" />
      <CommandList>
        <CommandEmpty>No commands match.</CommandEmpty>
        <CommandGroup heading="Commands">
          {items.map((it) => (
            <CommandItem
              key={it.id}
              value={it.haystack}
              onSelect={() => onSelect(it.name)}
            >
              <div className="flex flex-col gap-0.5">
                <div className="flex items-baseline gap-2 font-mono">
                  <span className="text-foreground">{it.name}</span>
                  {it.alias && (
                    <span className="text-xs text-muted-foreground">({it.alias})</span>
                  )}
                </div>
                <span className="line-clamp-1 text-xs text-muted-foreground">
                  {it.description}
                </span>
              </div>
            </CommandItem>
          ))}
        </CommandGroup>
      </CommandList>
    </CommandDialog>
  );
};

export default CommandPalette;
