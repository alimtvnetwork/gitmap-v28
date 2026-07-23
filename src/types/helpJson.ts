/**
 * TypeScript types for the `gitmap help --json` payload.
 *
 * Generated from spec/08-json-schemas/help-json.schema.json (draft-07).
 * Keep in sync with the schema — the Go contract test
 * `helpjson_jsonschema_contract_test.go` validates runtime output against
 * the schema on every build.
 */

export interface HelpJsonGroup {
  /** Intent-based super-group label (ANSI stripped). */
  group: string;
  /** One entry per command row in this group (ANSI stripped, trimmed). */
  lines: string[];
}

export interface HelpJsonPayload {
  /** Value of constants.Version at emission time (e.g. "5.43.1"). */
  version: string;
  /** Total number of help rows (sum of groups[].lines.length). */
  count: number;
  /** Present only when `--filter <q>` was supplied. Verbatim user query. */
  query?: string;
  /** Intent-grouped help rows. */
  groups: HelpJsonGroup[];
}

/** Narrow runtime guard — safe to use on untyped fetch() responses. */
export const isHelpJsonPayload = (value: unknown): value is HelpJsonPayload => {
  if (!value || typeof value !== "object") return false;
  const v = value as Record<string, unknown>;
  if (typeof v.version !== "string") return false;
  if (typeof v.count !== "number") return false;
  if (!Array.isArray(v.groups)) return false;
  return v.groups.every(
    (g) =>
      g &&
      typeof g === "object" &&
      typeof (g as HelpJsonGroup).group === "string" &&
      Array.isArray((g as HelpJsonGroup).lines) &&
      (g as HelpJsonGroup).lines.every((l) => typeof l === "string"),
  );
};
