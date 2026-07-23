/**
 * Lightweight changelog item classifier (#10).
 *
 * Until the structured-changelog source from #18 lands, we derive
 * tags heuristically from the human-written item text. This is
 * intentionally conservative — when in doubt we omit the tag rather
 * than mislabel, so the filter UI never hides items behind a wrong tag.
 */

export type ChangelogTag = "breaking" | "added" | "changed" | "flag" | "fix" | "perf";

const RULES: Array<{ tag: ChangelogTag; pattern: RegExp }> = [
  { tag: "breaking", pattern: /\b(breaking|removed|migration|backwards-?incompat|flips? to|default (?:flips|inverted))/i },
  { tag: "flag",     pattern: /(`--[a-z][\w-]*`|--[a-z][\w-]*\b|flag added|new flag)/i },
  { tag: "added",    pattern: /^\s*(?:added|new(?:\s|:))/i },
  { tag: "changed",  pattern: /^\s*(?:changed|updated|refactor|renamed)/i },
  { tag: "fix",      pattern: /^\s*(?:fix(?:ed)?|bug)/i },
  { tag: "perf",     pattern: /\b(?:perf(?:ormance)?|faster|speedup|optimi[sz]e)\b/i },
];

/** Returns the set of tags inferred from a single changelog item. */
export function classifyChangelogItem(item: string): ChangelogTag[] {
  const tags = new Set<ChangelogTag>();
  for (const { tag, pattern } of RULES) {
    if (pattern.test(item)) tags.add(tag);
  }
  return [...tags];
}

/** Filter labels surfaced in the UI, in display order. */
export const TAG_LABELS: Record<ChangelogTag, string> = {
  breaking: "Breaking",
  added:    "Added",
  changed:  "Changed",
  flag:     "Flags",
  fix:      "Fixes",
  perf:     "Perf",
};

export const TAG_ORDER: ChangelogTag[] = ["breaking", "added", "changed", "flag", "fix", "perf"];
