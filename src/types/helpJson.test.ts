import { describe, it, expect } from "vitest";
import { isHelpJsonPayload, type HelpJsonPayload } from "./helpJson";

describe("isHelpJsonPayload", () => {
  const valid: HelpJsonPayload = {
    version: "5.43.1",
    count: 2,
    groups: [{ group: "GET STARTED", lines: ["help  Show help", "version  Print version"] }],
  };

  it("accepts a schema-conformant payload", () => {
    expect(isHelpJsonPayload(valid)).toBe(true);
  });

  it("accepts an optional query field", () => {
    expect(isHelpJsonPayload({ ...valid, query: "clone" })).toBe(true);
  });

  it("rejects missing required fields", () => {
    expect(isHelpJsonPayload({ version: "5.43.1", count: 0 })).toBe(false);
    expect(isHelpJsonPayload(null)).toBe(false);
    expect(isHelpJsonPayload("nope")).toBe(false);
  });

  it("rejects malformed groups", () => {
    expect(
      isHelpJsonPayload({ ...valid, groups: [{ group: 1, lines: [] }] }),
    ).toBe(false);
    expect(
      isHelpJsonPayload({ ...valid, groups: [{ group: "x", lines: [42] }] }),
    ).toBe(false);
  });
});
