# AI Release Synchronization

When an AI Agent handles a release for this repository, it MUST synchronously update all canonical version strings across the codebase *before* pushing any Git tag. Failure to do so will cause the `Installer Smoke Test` to fail.

## Standard Operating Procedure

### 1. Perform PowerShell Version Bump
You must dynamically sweep and replace the old version string (`OLD.VER.SION`) with the new version string (`X.Y.Z`) across the following canonical files:

```powershell
$newVer = "X.Y.Z"
$files = @(
  ".gitmap/release/latest.json", 
  ".lovable/plan.md", 
  "gitmap/constants/constants.go", 
  "readme.md", 
  "src/constants/index.ts", 
  "package.json"
)

foreach ($f in $files) { 
  $c = Get-Content $f -Raw
  $c = $c -replace 'OLD\.VER\.SION', $newVer
  Set-Content $f -Value $c -NoNewline 
}
```
**CRITICAL**: You must also duplicate the release manifest JSON:
```powershell
Copy-Item ".gitmap/release/vOLD.VER.SION.json" ".gitmap/release/vX.Y.Z.json"
$c = Get-Content ".gitmap/release/vX.Y.Z.json" -Raw
$c = $c -replace 'OLD\.VER\.SION', $newVer
Set-Content ".gitmap/release/vX.Y.Z.json" -Value $c -NoNewline
```

### 2. Push to Branch
```bash
git add .
git commit -m "chore(release): bump version to vX.Y.Z"
git push origin main
```
You MUST wait for GitHub to register the commit before proceeding to step 3.

### 3. Tag Execution
Only push the tag after the source branch has been synchronized and pushed.
```bash
git tag vX.Y.Z
git push origin vX.Y.Z
```
