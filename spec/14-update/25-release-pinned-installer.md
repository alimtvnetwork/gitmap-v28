I want to have a `release-install.ps1` file, and that file will only be used for the release tab or release pages. The install script we were using before should not be used directly. The idea is that there will be a `release-install.ps1` and a similar shell file. When a user reaches the Release page from a URL, the script will know from that URL which version was requested. From the requested version, it will redirect its request to the same repo and ensure that it does not walk forward to the next `v1`, `v2`, `v3`, etc. There should be a parameter in the install script that fixates the version — something like `no-update` or `no-latest` — so that it installs exactly that version. Since we want to keep it simple, it should be a separate script that knows from this exact version which one to install. The first job is to write the release script spec file. Update the memory and the plan, and then we will do the implementation.

## Important

- Do not act. Spec only. Implementation begins after explicit approval.
- Pinned version is mandatory. Never resolve `latest`, `main`, `master`, or `HEAD`.
- Generic `install.ps1` / `install.sh` remain unchanged as the "always latest" installers.
- New scripts are additive: `release-install.ps1` and `release-install.sh`.

## Scope

1. New Scripts
   a. `release-install.ps1` (Windows / PowerShell)
   b. `release-install.sh` (macOS / Linux / bash)
   c. Located at repository root, alongside existing installers.

2. Untouched
   a. `install.ps1` and `install.sh` keep their current "latest" behavior.
   b. Distributed via the existing `curl | iex` / `curl | bash` one-liner.

## Behavioral Contract

1. Inputs
   a. Explicit version
   i. PowerShell: `-Version <tag>`
   ii. Bash: `--version <tag>`
   b. Pin enforcement
   i. PowerShell: `-NoUpdate` (default ON, cannot be disabled)
   ii. Bash: `--no-update` (default ON, cannot be disabled)
   c. Help
   i. PowerShell: `-Help` / `-?`
   ii. Bash: `--help` / `-h`

2. Version Resolution Order
   a. Explicit `-Version` / `--version` argument wins.
   b. `$INSTALLER_VERSION` / `$env:INSTALLER_VERSION` environment variable (ratified extension; useful for CI pipelines and Dockerfiles where flags are awkward).
   c. Otherwise use baked-in `VERSION_PLACEHOLDER` replaced at release-build time.
   d. If none of the above is available, hard fail with non-zero exit.
   e. Never fall back to `latest`, `main`, or `HEAD`.
   f. When two sources disagree, higher precedence wins and a warning is emitted to stderr.

3. Pin Enforcement Guarantees
   a. Download only from `https://github.com/<owner>/<repo>/releases/download/<PINNED_TAG>/...`.
   b. Pass `--pinned-by-release-install <tag>` into any chained inner installer.
   c. Inner installer must skip all update / self-upgrade logic when this flag is present.
   d. Print resolved pinned version to stdout before any download.
   e. Never query GitHub API for `latest` releases.
   f. Never clone or fetch the source git repo.

4. Validation
   a. Tag must match `^v?\d+\.\d+\.\d+(-[A-Za-z0-9.]+)?$`.
   b. Reject anything else with a clear error to prevent URL path injection.
   c. HEAD-check pinned asset URL before download. On 404, fail clearly.

5. Checksum and Telemetry
   a. Mirror the existing generic installer's verification approach (e.g., SHA256).
   b. Verify checksums from the pinned tag's assets, not from `latest`.
   c. No new telemetry introduced by the release script.

## Resolution Algorithm

1. Parse arguments.
2. If `-Version` / `--version` provided
   a. `resolved_version = arg`
   b. If baked-in placeholder is set and differs, warn to stderr.
3. Else if `VERSION_PLACEHOLDER` has been replaced with a real tag
   a. `resolved_version = baked-in value`
4. Else
   a. Error: "release-install requires a pinned version. Pass -Version <tag> or use the script from a Release page."
   b. Exit 1.
5. Validate `resolved_version` against semver regex; reject otherwise.
6. Print: "Installing pinned version: <resolved_version>".
7. Build asset URL: `https://github.com/<owner>/<repo>/releases/download/<resolved_version>/<asset>`.
8. HEAD-check URL. On 404, error "release <tag> not found or asset missing"; exit 1.
9. Download, verify checksum, extract.
10. Invoke inner installer with `--pinned-by-release-install <resolved_version>`.
11. Inner installer must skip all update/upgrade logic when this flag is present.

## Failure Modes

1. Exit 1: No version resolvable (no arg, no baked-in tag).
2. Exit 2: Version string fails semver regex.
3. Exit 3: Pinned release or asset not found (404).
4. Exit 4: Checksum mismatch.
5. Exit 5: Inner installer rejected `--pinned-by-release-install` (version skew).

## Forbidden Behaviors

1. Resolving `latest`, `main`, `master`, or `HEAD`.
2. Walking forward to a newer tag.
3. Cloning or fetching the source git repo.
4. Silently disabling `-NoUpdate` / `--no-update`.
5. Continuing past a checksum failure.

## Release-Time Build Step

1. On tag push, the release workflow (e.g., GitHub Action) runs.
2. Steps
   a. Take canonical `release-install.ps1` / `release-install.sh` from repo root.
   b. Replace `VERSION_PLACEHOLDER` with the concrete tag (e.g., `v1.4.2`).
   c. Upload baked copies as release assets on that release.
   d. Ensure the Release page advertises the one-liner pointing to those baked assets.
3. Unbaked copies in repo root remain useful for
   a. Standalone power-user invocation with `-Version`.
   b. Being the canonical source the workflow bakes from.

## Release Page One-Liners

1. PowerShell

   ```powershell
   irm https://github.com/<owner>/<repo>/releases/download/v1.4.2/release-install.ps1 | iex
   ```

2. Bash

   ```bash
   curl -fsSL https://github.com/<owner>/<repo>/releases/download/v1.4.2/release-install.sh | bash
   ```

## Ambiguities and Resolutions

1. AMB-1: Behavior when an explicit version is passed to the generic `install.sh` / `install.ps1`
   a. Resolution: Honor strict pinning for that invocation. "Latest" applies only when no version is specified.

2. AMB-2: What "look at the git repo" forbids when pinned
   a. No `git clone`.
   b. No fetching `main` / `master` / `HEAD` tarballs.
   c. No GitHub API calls to `/releases/latest` or `/tags` beyond verifying the pinned tag exists.
   d. Allowed: HEAD request to confirm the pinned asset URL returns 200.

3. AMB-3: How the Release page advertises the URL
   a. Out of scope for the script. Handled by the Release page generator (release notes template or GitHub Action).

4. AMB-4: When `-Version` disagrees with baked-in tag
   a. `-Version` wins. Print a warning to stderr. Proceed with user-requested version.

5. AMB-5: Asset naming convention
   a. Unchanged. Only the base URL changes from `/releases/latest/download/` to `/releases/download/<PINNED_TAG>/`.

6. AMB-6: Checksum / signature verification
   a. Inherit existing generic installer's verification approach against the pinned tag's assets.

7. AMB-7: Telemetry
   a. Inherit existing installer behavior. No additions.

8. AMB-8: Tag format validation
   a. Accept `^v?\d+\.\d+\.\d+(-[A-Za-z0-9.]+)?$`. Reject anything else.

## Acceptance Criteria

1. Running `release-install.ps1` / `release-install.sh` with no version and no baked-in tag exits non-zero with the documented error.
2. Running with `-Version v1.4.2` downloads only from `/releases/download/v1.4.2/` and never from `/releases/latest/`.
3. Running the baked release-asset copy installs exactly that release's pinned version.
4. Passing `-Version` that disagrees with the baked-in tag prints a stderr warning and uses the argument value.
5. Invalid tag formats are rejected with exit code 2 before any network call.
6. 404 on pinned asset URL exits with code 3 and a clear message.
7. Checksum mismatch exits with code 4 and aborts installation.
8. Inner installer receives `--pinned-by-release-install <tag>` and skips all update logic.
9. Generic `install.ps1` / `install.sh` behavior is unchanged when invoked without a version.
10. Generic installers, when invoked with `--version <tag>`, also enforce strict pinning (AMB-1).

## Open Items for User Confirmation

1. Repository target
   a. Confirm the actual repo (owner/name) where these scripts live. The current Lovable project is a React template and is not the install-script repo.

2. AMB-1 confirmation
   a. Confirm that the generic `install.sh` / `install.ps1` should also enforce strict pinning when a version is passed.

3. Asset naming (AMB-5)
   a. Confirm existing asset naming scheme is unchanged.

4. Checksum (AMB-6)
   a. Confirm the generic installer's exact verification approach so the release variant mirrors it.

## File System References

1. Spec file: `prompts/01-release-install.md`

## Job Sequence

1. Job 1: Write this spec into the spec system. Status: Done on approval.
2. Job 2: Implementation. Blocked on Open Items confirmation.
   a. Create `release-install.ps1` and `release-install.sh` per Behavioral Contract.
   b. Wire `VERSION_PLACEHOLDER` baking into the release workflow.
   c. Add `--pinned-by-release-install` handling to existing inner installers.
   d. Update README to document the two installer modes.
