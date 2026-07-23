import { describe, it, expect } from "vitest";
import { readFileSync, existsSync } from "node:fs";
import { resolve } from "node:path";
import { VERSION } from "@/constants/index";

/**
 * version-sync.test.ts — guards against the v3.0.0 follow-up bug
 * where the web `VERSION` constant silently drifted from the Go
 * binary version (was stuck at v4.14.0 while the binary was at
 * v4.22.0 for eight releases).
 *
 * Reads `gitmap-v27/constants/constants.go`, extracts the `Version`
 * literal, and asserts that `src/constants/index.ts` exports the
 * same value with a `v` prefix.
 *
 * If this test fails: bump `VERSION` in `src/constants/index.ts` to
 * match the Go binary -- that file is the second half of every
 * release (the Go bump is the first half). See
 * `.lovable/memory/project/version-bump-procedure.md`.
 */
describe("VERSION sync between web and Go binary", () => {
  const goConstantsPath = resolve(__dirname, "../../gitmap/constants/constants.go");

  it("can locate the Go constants file", () => {
    expect(existsSync(goConstantsPath)).toBe(true);
  });

  it("matches the Go binary Version constant (with v prefix)", () => {
    const source = readFileSync(goConstantsPath, "utf8");
    // Match: const Version = "4.22.0"
    // Tolerates extra whitespace and either single or double quotes
    // (Go only uses double quotes but defensive parsing is cheap).
    const match = source.match(/const\s+Version\s*=\s*["']([^"']+)["']/);
    expect(match, `Version constant not found in ${goConstantsPath}`).not.toBeNull();
    const goVersion = match![1];
    expect(VERSION).toBe(`v${goVersion}`);
  });
});
