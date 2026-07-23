/** @deprecated Use gitmap/constants/constants.go as the single source of truth. */
export const VERSION = "v6.79.0" as const;

/** Repo status values used in watch/status displays. */
export const RepoStatus = {
  Clean: "clean",
  Dirty: "dirty",
} as const;

export type RepoStatus = (typeof RepoStatus)[keyof typeof RepoStatus];

/** Terminal line type for animated demos. */
export const TerminalLineType = {
  Input: "input",
  Output: "output",
  Header: "header",
  Accent: "accent",
} as const;

export type TerminalLineType = (typeof TerminalLineType)[keyof typeof TerminalLineType];

/** Project filter including the "all" option. */
export const FILTER_ALL = "all" as const;

/** Animation timing defaults (ms). */
export const TERMINAL_INPUT_DELAY = 600;
export const TERMINAL_OUTPUT_DELAY = 120;

/** Watch dashboard refresh interval (seconds). */
export const WATCH_REFRESH_INTERVAL = 30;

/** Status indicator symbols. */
export const STATUS_ICON_DIRTY = "●";
export const STATUS_ICON_CLEAN = "✔";

/** Root-relative path placeholder. */
export const ROOT_RELATIVE_PATH = ".";
export const ROOT_RELATIVE_LABEL = "(root)";
