# Subtask 213.3: GitMap Installer Constants Update

- **Parent Plan:** [83-antigravity-ubuntu-update-and-macro-automation.md](../../83-antigravity-ubuntu-update-and-macro-automation.md)
- **Spec Reference:** [02-spec/21-app/213-antigravity-ubuntu-update-and-macro-automation/02-gitmap-macro-and-installer-spec.md](../../../../02-spec/21-app/213-antigravity-ubuntu-update-and-macro-automation/02-gitmap-macro-and-installer-spec.md)
- **Status:** Ready
- **Target Subsystem:** `cli/cmdinstall`
- **Target File:** `cli/cmdinstall/installantigravity_types.go`

---

## 1. Objective

Update the hardcoded default Antigravity installation coordinates in GitMap's installer package (`cli/cmdinstall/installantigravity_types.go`) from version **2.13.0** (Build `6362815968182272`) to **2.19.1** (Build `6046815158665216`). Ensure that any subsequent invocations of `gitmap install antigravity` across the fleet resolve to the latest stable release artifact without manual override flags.

---

## 2. Code Modification Specification

### 2.1 Target File: `cli/cmdinstall/installantigravity_types.go`

Locate lines 15–22 containing the default constants:

```go
<<<<
const (
	AntigravityDefaultVersion = "2.13.0"
	AntigravityDefaultBuildID = "6362815968182272"
	AntigravityBaseURL        = "https://storage.googleapis.com/antigravity-public/antigravity-hub"
	ErrUnsupportedPlatform    = "E_UNSUPPORTED_PLATFORM"
	ErrPrerequisiteFailed     = "E_PREREQUISITE_FAILED"
	ErrDownloadFailed         = "E_DOWNLOAD_FAILED"
)
====
const (
	AntigravityDefaultVersion = "2.19.1"
	AntigravityDefaultBuildID = "6046815158665216"
	AntigravityBaseURL        = "https://storage.googleapis.com/antigravity-public/antigravity-hub"
	ErrUnsupportedPlatform    = "E_UNSUPPORTED_PLATFORM"
	ErrPrerequisiteFailed     = "E_PREREQUISITE_FAILED"
	ErrDownloadFailed         = "E_DOWNLOAD_FAILED"
)
>>>>
```

---

## 3. URL Resolution & Artifact Verification

The installer resolves the download URL dynamically per platform:

```go
func ResolveAntigravityDownloadURL(info AntigravityPlatformInfo) string {
    return fmt.Sprintf("%s/%s-%s/%s/%s",
        AntigravityBaseURL,
        AntigravityDefaultVersion,
        AntigravityDefaultBuildID,
        info.PlatformID,
        info.ArtifactName,
    )
}
```

### Verified Target Coordinates:
- **Linux x64:** `https://storage.googleapis.com/antigravity-public/antigravity-hub/2.19.1-6046815158665216/linux-x64/Antigravity.tar.gz`
- **Windows x64:** `https://storage.googleapis.com/antigravity-public/antigravity-hub/2.19.1-6046815158665216/win32-x64/AntigravitySetup.exe`
- **Darwin ARM64:** `https://storage.googleapis.com/antigravity-public/antigravity-hub/2.19.1-6046815158665216/darwin-arm64/Antigravity.dmg`

---

## 4. Verification & Testing Procedures

1. **Syntax & Unit Test Verification:**
   Execute all unit tests in the `cli/cmdinstall` package to verify zero regressions:
   ```bash
   go test -v ./cli/cmdinstall/...
   ```
2. **Build Verification:**
   Verify full project compilation with the updated constants:
   ```bash
   go build -v ./...
   ```
3. **URL Connectivity Smoke Check:**
   Confirm HTTP 200 OK header response from the Google Cloud Storage bucket:
   ```powershell
   curl.exe -I "https://storage.googleapis.com/antigravity-public/antigravity-hub/2.19.1-6046815158665216/linux-x64/Antigravity.tar.gz"
   ```

---

## 5. Acceptance Criteria

- [ ] `AntigravityDefaultVersion` equals `"2.19.1"` in `cli/cmdinstall/installantigravity_types.go`.
- [ ] `AntigravityDefaultBuildID` equals `"6046815158665216"` in `cli/cmdinstall/installantigravity_types.go`.
- [ ] `go test ./cli/cmdinstall/...` passes with zero errors.
- [ ] `HasValidPlatform`, `HasDownloadURL`, `IsLinux`, and `IsDarwin` methods continue to function unmodified.
- [ ] Computed Linux artifact URL returns HTTP 200 status code.
