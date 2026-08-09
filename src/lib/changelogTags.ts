/**
 * Lightweight changelog item classifier (#10).
 *
 * Until the structured-changelog source from #18 lands, we derive
 * tags heuristically from the human-written item text. This is
 * intentionally conservative — when in doubt we omit the tag rather
 * than mislabel, so the filter UI never hides items behind a wrong tag.
 */

export enum ChangelogTagType {
  Breaking = "breaking",
  Added = "added",
  Changed = "changed",
  Flag = "flag",
  Fix = "fix",
  Perf = "perf",
}
export type ChangelogTag = ChangelogTagType;

const RULES: Array<{ tag: ChangelogTagType; pattern: RegExp }> = [
  { tag: ChangelogTagType.Breaking, pattern: /\b(breaking|removed|migration|backwards-?incompat|flips? to|default (?:flips|inverted))/i },
  { tag: ChangelogTagType.Flag,     pattern: /(`--[a-z][\w-]*`|--[a-z][\w-]*\b|flag added|new flag)/i },
  { tag: ChangelogTagType.Added,    pattern: /^\s*(?:added|new(?:\s|:))/i },
  { tag: ChangelogTagType.Changed,  pattern: /^\s*(?:changed|updated|refactor|renamed)/i },
  { tag: ChangelogTagType.Fix,      pattern: /^\s*(?:fix(?:ed)?|bug)/i },
  { tag: ChangelogTagType.Perf,     pattern: /\b(?:perf(?:ormance)?|faster|speedup|optimi[sz]e)\b/i },
];

/** Returns the set of tags inferred from a single changelog item. */
export function classifyChangelogItem(item: string): ChangelogTagType[] {
  const tags = new Set<ChangelogTagType>();
  for (const { tag, pattern } of RULES) {
    if (pattern.test(item)) tags.add(tag);
  }
  return [...tags];
}

/** Filter labels surfaced in the UI, in display order. */
export const TAG_LABELS: Record<ChangelogTagType, string> = {
  [ChangelogTagType.Breaking]: "Breaking",
  [ChangelogTagType.Added]:    "Added",
  [ChangelogTagType.Changed]:  "Changed",
  [ChangelogTagType.Flag]:     "Flags",
  [ChangelogTagType.Fix]:      "Fixes",
  [ChangelogTagType.Perf]:     "Perf",
};

export const TAG_ORDER: ChangelogTagType[] = [
  ChangelogTagType.Breaking,
  ChangelogTagType.Added,
  ChangelogTagType.Changed,
  ChangelogTagType.Flag,
  ChangelogTagType.Fix,
  ChangelogTagType.Perf,
];
