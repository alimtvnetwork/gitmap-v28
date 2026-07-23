// Pure helpers for the /release/:version page. Extracted so they can be
// unit-tested without mounting React. Spec: spec/01-app/105-release-version-script.md.

export const RELEASE_REPO = "alimtvnetwork/gitmap-v27";
export const RELEASE_DOCS_HOST = "https://gitmap.dev";
export const SEMVER_TAG = /^v\d+\.\d+\.\d+(-[A-Za-z0-9.]+)?$/;

export type ReleasePlatform = "windows" | "unix";

export interface InstallSnippet {
  pinned: string;
  generic: string;
}

export const buildReleaseSnippets = (
  version: string,
  platform: ReleasePlatform,
): InstallSnippet => {
  const releaseBase = `https://github.com/${RELEASE_REPO}/releases/download/${version}`;

  if (platform === "windows") {
    return {
      pinned: [
        `# Pinned install — locks gitmap to ${version} (no auto-upgrade)`,
        `iwr ${releaseBase}/release-version-${version}.ps1 -OutFile $env:TEMP\\rv.ps1`,
        `& $env:TEMP\\rv.ps1`,
      ].join("\n"),
      generic: [
        `# Generic install — same script, version passed as parameter`,
        `iwr ${RELEASE_DOCS_HOST}/scripts/release-version.ps1 -OutFile $env:TEMP\\rv.ps1`,
        `& $env:TEMP\\rv.ps1 -Version ${version}`,
      ].join("\n"),
    };
  }

  return {
    pinned: [
      `# Pinned install — locks gitmap to ${version} (no auto-upgrade)`,
      `curl -fsSL ${releaseBase}/release-version-${version}.sh | bash`,
    ].join("\n"),
    generic: [
      `# Generic install — same script, version passed as parameter`,
      `curl -fsSL ${RELEASE_DOCS_HOST}/scripts/release-version.sh | bash -s -- --version ${version}`,
    ].join("\n"),
  };
};

export const normalizeReleaseVersion = (raw: string): string =>
  raw.startsWith("v") ? raw : `v${raw}`;

export const isValidReleaseVersion = (raw: string): boolean =>
  SEMVER_TAG.test(normalizeReleaseVersion(raw));
