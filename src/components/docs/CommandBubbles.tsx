import { Link, useLocation } from "react-router-dom";

interface Bubble {
  name: string;
  alias?: string;
  description: string;
  /** Optional deep-link target on /commands. Falls back to /commands. */
  to?: string;
}

const COMMANDS_PATH = "/commands";
const HISTORY_REWRITE_PATH = "/history-rewrite";

const BUBBLES: Bubble[] = [
  { name: "scan", alias: "s", description: "Discover Git repos on disk" },
  { name: "clone", alias: "c", description: "Re-clone from a scan file" },
  { name: "clone-next", alias: "cn", description: "Clone next versioned iteration" },
  { name: "pull", alias: "p", description: "Pull latest for tracked repos" },
  { name: "watch", alias: "w", description: "Live status dashboard" },
  { name: "exec", alias: "x", description: "Run git across all repos" },
  { name: "release", alias: "r", description: "Branch, tag, push, attach" },
  { name: "as", alias: "s-alias", description: "Alias the current repo" },
  { name: "inject", alias: "inj", description: "Register folder + open VS Code" },
  { name: "cd", alias: "go", description: "Jump shell into a tracked repo" },
  { name: "group", alias: "g", description: "Manage repo groups" },
  { name: "changelog", alias: "cl", description: "View release notes" },
  { name: "history-purge", alias: "hp", description: "Remove file(s) from all history (sandbox)", to: HISTORY_REWRITE_PATH },
  { name: "history-pin", alias: "hpin", description: "Pin file(s) to current content across history", to: HISTORY_REWRITE_PATH },
];

const CommandBubblesHeader = ({ isOnCommands }: { isOnCommands: boolean }) => (
  <div className="mb-5 flex items-baseline justify-between gap-4">
    <h2 id="command-bubbles-heading" className="font-heading text-lg font-semibold text-foreground">
      Explore commands
    </h2>
    <Link
      to={COMMANDS_PATH}
      aria-label="View all gitmap commands on the commands reference page"
      aria-current={isOnCommands ? "page" : undefined}
      className="rounded-sm text-xs font-sans text-muted-foreground hover:text-primary transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 focus-visible:ring-offset-background"
    >
      View all →
    </Link>
  </div>
);

interface BubbleItemProps {
  bubble: Bubble;
  isOnCommands: boolean;
}

const BubbleItem = ({ bubble, isOnCommands }: BubbleItemProps) => {
  const target = bubble.to ?? COMMANDS_PATH;
  const aliasText = bubble.alias ? ` (alias: ${bubble.alias})` : "";
  const ariaLabel = `${bubble.name} command${aliasText}: ${bubble.description}. Opens the commands reference.`;

  return (
    <li>
      <Link
        to={target}
        title={bubble.description}
        aria-label={ariaLabel}
        aria-current={isOnCommands ? "page" : undefined}
        className="btn-slide btn-slide-ghost group inline-flex items-center gap-2 rounded-full border border-border bg-card px-4 py-1.5 text-sm font-sans text-foreground hover:border-primary/50 hover:bg-secondary"
      >
        <code aria-hidden="true" className="font-mono text-sm text-primary">
          {bubble.name}
        </code>
        {bubble.alias && (
          <span aria-hidden="true" className="font-mono text-xs text-muted-foreground">
            {bubble.alias}
          </span>
        )}
      </Link>
    </li>
  );
};

const CommandBubbles = () => {
  const { pathname } = useLocation();
  const isOnCommands = pathname === COMMANDS_PATH;

  return (
    <section className="reveal py-8" aria-labelledby="command-bubbles-heading">
      <CommandBubblesHeader isOnCommands={isOnCommands} />
      <ul className="flex flex-wrap gap-2 list-none p-0 m-0" aria-label="Quick links to common gitmap commands">
        {BUBBLES.map((b) => (
          <BubbleItem key={b.name} bubble={b} isOnCommands={isOnCommands} />
        ))}
      </ul>
    </section>
  );
};

export default CommandBubbles;
