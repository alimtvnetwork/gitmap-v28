# CI/CD RCA: Windows Binary Archive Packaging vs Executable PE Target

- Job: Local CLI Installation & Windows PowerShell Invocation
- Type: NativeCommandFailed / ResourceUnavailable
- Tool: Windows PowerShell 5.1 & PowerShell Core (`Microsoft.PowerShell_profile.ps1`)
- Status: ✅ Resolved

---

## 1. Symptom

Running `gitmap update` or `gitmap` in Windows PowerShell (`powershell.exe`) failed with:
```text
Program 'gitmap.exe' failed to run: The specified executable is not a valid application for this OS platform.At
C:\Users\Administrator\Documents\WindowsPowerShell\Microsoft.PowerShell_profile.ps1:123 char:5
+     & $real @args
+     ~~~~~~~~~~~~~.
At C:\Users\Administrator\Documents\WindowsPowerShell\Microsoft.PowerShell_profile.ps1:123 char:5
+     & $real @args
+     ~~~~~~~~~~~~~
    + CategoryInfo          : ResourceUnavailable: (:) [], ParentContainsErrorRecordException
    + FullyQualifiedErrorId : NativeCommandFailed
```

---

## 2. Root Cause

1. **Building Archive Instead of Executable PE**:
   When building locally with `go build -o bin/gitmap.exe ./cmd` in directory `cli/`, the target path `./cmd` is a library package (`package cmd`), not `package main` (which resides in `cli/main.go`). Go built a static archive (`.a`) with magic bytes `!<arch>\n` (`[33, 60, 97, 114, 99, 104, 62]`) instead of a portable executable (`MZ`).
2. **Copying Non-Executable File to Local Bin**:
   The resulting `.a` file was copied to `C:\Users\Administrator\AppData\Local\gitmap-cli\gitmap.exe`. When Windows PowerShell resolved `$real` via `Get-GitmapCommand` and attempted process creation (`& $real`), the Windows OS loader rejected the file as not a valid executable for this platform (`ERROR_BAD_EXE_FORMAT` / 193).

---

## 3. Resolution

1. **Rebuild with Main Package**:
   Rebuilt the Windows binary targeting the main package root:
   ```powershell
   cd cli
   $env:GOOS = "windows"
   $env:GOARCH = "amd64"
   go build -o bin/gitmap.exe .
   ```
2. **Verify PE Executable Header**:
   Verified magic bytes are `MZ` (`b'MZ'`), confirming valid Windows PE structure.
3. **Deploy & End-to-End Test**:
   Installed to `C:\Users\Administrator\AppData\Local\gitmap-cli\gitmap.exe`. Ran `gitmap update` in Windows PowerShell (`powershell.exe`), which downloaded the official `v6.317.0` release package via `aria2c`, verified SHA256 checksums, executed `gitmap setup`, and confirmed:
   ```text
   ✔ Successfully updated from v6.317.0 to v6.317.0
   PASS Version: gitmap v6.317.0
   ```

---

## 4. Prevention & Learnings

- Always compile the CLI executable using `go build -o bin/gitmap.exe .` in `cli/` (or `go build -o bin/gitmap.exe cli/main.go`), never `./cmd` which is a non-main package.
- Include a quick header sanity check (`b'MZ'`) whenever deploying local Windows binaries.
