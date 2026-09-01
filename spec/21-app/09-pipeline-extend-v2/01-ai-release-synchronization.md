# AI Release Synchronization

When an AI Agent handles a release for this repository, it MUST synchronously update all canonical version strings across the codebase *before* pushing any Git tag. Failure to do so will cause the `Installer Smoke Test` to fail.

## Standard Operating Procedure

### 1. Update Central `version.json`

The absolute source of truth for the version is located at the root of the repository in `version.json`. 

You must only bump this file:
```json
{
  "version": "6.91.0"
}
```

The AI MUST NOT attempt to use PowerShell to mutate the `gitmap` codebase (`constants.go`) or `.gitmap/release/latest.json`. The CI pipeline parses `version.json` via `jq` and actively injects it into Golang via `-ldflags` during compilation.

### 2. Update Web/Changelog Metadata

Make sure you update `package.json` and inject the `## [vX.Y.Z] YYYY-MM-DD` block into `changelog.md`.

### 3. Push to Branch

```bash
git add .
git commit -m "chore(release): bump version to vX.Y.Z"
git push origin main
```
You MUST wait for GitHub to register the commit before proceeding to step 4.

### 4. Tag Execution

Only push the tag after the source branch has been synchronized and pushed.
```bash
git tag vX.Y.Z
git push origin vX.Y.Z
```
