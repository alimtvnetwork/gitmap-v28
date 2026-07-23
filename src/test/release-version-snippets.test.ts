import { describe, expect, it } from "vitest";
import {
  buildReleaseSnippets,
  isValidReleaseVersion,
  normalizeReleaseVersion,
  RELEASE_DOCS_HOST,
  RELEASE_REPO,
} from "@/pages/releaseVersionSnippets";

const VERSION = "v3.39.0";

describe("releaseVersionSnippets — semver guard", () => {
  it("accepts canonical vMAJOR.MINOR.PATCH", () => {
    expect(isValidReleaseVersion("v1.2.3")).toBe(true);
    expect(isValidReleaseVersion("3.39.0")).toBe(true); // auto-prefix
    expect(isValidReleaseVersion("v4.23.0-rc.1")).toBe(true);
  });

  it("rejects garbage and the word latest", () => {
    expect(isValidReleaseVersion("latest")).toBe(false);
    expect(isValidReleaseVersion("v1")).toBe(false);
    expect(isValidReleaseVersion("v1.2")).toBe(false);
    expect(isValidReleaseVersion("vX.Y.Z")).toBe(false);
    expect(isValidReleaseVersion("")).toBe(false);
  });

  it("normalizes raw versions to a leading v", () => {
    expect(normalizeReleaseVersion("3.39.0")).toBe("v3.39.0");
    expect(normalizeReleaseVersion("v3.39.0")).toBe("v3.39.0");
  });
});

describe("releaseVersionSnippets — Windows", () => {
  const snip = buildReleaseSnippets(VERSION, "windows");

  it("pinned snippet hits the snapshot asset on the GitHub release", () => {
    expect(snip.pinned).toContain(
      `https://github.com/${RELEASE_REPO}/releases/download/${VERSION}/release-version-${VERSION}.ps1`,
    );
    expect(snip.pinned).toContain("$env:TEMP\\rv.ps1");
  });

  it("generic snippet hits the docs host with -Version flag", () => {
    expect(snip.generic).toContain(
      `${RELEASE_DOCS_HOST}/scripts/release-version.ps1`,
    );
    expect(snip.generic).toContain(`-Version ${VERSION}`);
  });

  it("matches snapshot", () => {
    expect(snip).toMatchInlineSnapshot(`
      {
        "generic": "# Generic install — same script, version passed as parameter
      iwr https://gitmap.dev/scripts/release-version.ps1 -OutFile $env:TEMP\\rv.ps1
      & $env:TEMP\\rv.ps1 -Version v3.39.0",
        "pinned": "# Pinned install — locks gitmap-v27 to v3.39.0 (no auto-upgrade)
      iwr https://github.com/alimtvnetwork/gitmap-v27/releases/download/v3.39.0/release-version-v3.39.0.ps1 -OutFile $env:TEMP\\rv.ps1
      & $env:TEMP\\rv.ps1",
      }
    `);
  });
});

describe("releaseVersionSnippets — Unix", () => {
  const snip = buildReleaseSnippets(VERSION, "unix");

  it("pinned snippet curl|bash the snapshot .sh asset", () => {
    expect(snip.pinned).toContain(
      `https://github.com/${RELEASE_REPO}/releases/download/${VERSION}/release-version-${VERSION}.sh`,
    );
    expect(snip.pinned).toContain("| bash");
  });

  it("generic snippet pipes the docs host with --version flag", () => {
    expect(snip.generic).toContain(
      `${RELEASE_DOCS_HOST}/scripts/release-version.sh`,
    );
    expect(snip.generic).toContain(`bash -s -- --version ${VERSION}`);
  });

  it("matches snapshot", () => {
    expect(snip).toMatchInlineSnapshot(`
      {
        "generic": "# Generic install — same script, version passed as parameter
      curl -fsSL https://gitmap.dev/scripts/release-version.sh | bash -s -- --version v3.39.0",
        "pinned": "# Pinned install — locks gitmap-v27 to v3.39.0 (no auto-upgrade)
      curl -fsSL https://github.com/alimtvnetwork/gitmap-v27/releases/download/v3.39.0/release-version-v3.39.0.sh | bash",
      }
    `);
  });
});

describe("releaseVersionSnippets — spec 105 invariants", () => {
  it("never contains the word 'latest' on either platform", () => {
    for (const platform of ["windows", "unix"] as const) {
      const snip = buildReleaseSnippets(VERSION, platform);
      expect(snip.pinned.toLowerCase()).not.toContain("latest");
      expect(snip.generic.toLowerCase()).not.toContain("latest");
    }
  });

  it("never references /releases/latest GitHub endpoint", () => {
    for (const platform of ["windows", "unix"] as const) {
      const snip = buildReleaseSnippets(VERSION, platform);
      expect(snip.pinned).not.toContain("/releases/latest");
      expect(snip.generic).not.toContain("/releases/latest");
    }
  });

  it("bakes the exact requested version into every URL", () => {
    const v = "v9.9.9";
    const snip = buildReleaseSnippets(v, "windows");
    const occurrences = (snip.pinned.match(/v9\.9\.9/g) || []).length;
    expect(occurrences).toBeGreaterThanOrEqual(2); // comment + URL
  });
});
