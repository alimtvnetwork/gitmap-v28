// Single source of truth for tooltip presentation across the docs
// header and toolbars. Importing these tokens — instead of writing
// `side="bottom"` and `sideOffset={6}` ad-hoc — guarantees every
// docs tooltip opens in the same direction, with the same gap, and
// (via DOCS_TOOLTIP_DELAY_MS) the same hover latency.
//
// The delay is enforced at the provider level in App.tsx; the
// constant lives here so tests/visual audits can reference one
// number rather than chasing a literal across files.

export const DOCS_TOOLTIP_SIDE = "bottom" as const;
export const DOCS_TOOLTIP_SIDE_OFFSET = 6;
export const DOCS_TOOLTIP_ALIGN = "center" as const;

// Hover latency before a tooltip opens (ms). Re-tooltipping a
// neighbour within DOCS_TOOLTIP_SKIP_DELAY_MS opens instantly so
// scrubbing across an icon row feels snappy.
export const DOCS_TOOLTIP_DELAY_MS = 150;
export const DOCS_TOOLTIP_SKIP_DELAY_MS = 300;
