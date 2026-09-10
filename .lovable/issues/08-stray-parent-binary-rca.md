# Issue 08: Stray Binary in Workspace Parent Directory (`../gitmap.exe`) Root Cause Analysis

## 1. Problem Description

A compiled binary `gitmap.exe` (29.35 MB, dated September 5, 2026) was found deposited in the workspace parent directory (`../gitmap.exe`), outside the repository boundaries and outside canonical deployment targets (`bin/gitmap.exe` and `%LOCALAPPDATA%\gitmap-cli\gitmap.exe`).

## 2. Root Cause Analysis

1. **Relative Directory Ambiguity in Build Documentation:**
   Multiple documentation guides (`readme.md`, `what-to-read.md`, and `gitmap/cmd/llmdocssections.go`) historically prescribed the build snippet:
   ```bash
   cd gitmap && go build -o ../gitmap .
   ```
   When developers or automation scripts ran this while already inside the repository root (`gitmap`), `../gitmap` resolved to the parent directory (`../gitmap.exe`), unintentionally placing a 29+ MB binary in the parent directory.

2. **Makefile Build Target Outputting to Relative Parent:**
   The `Makefile` target `build` previously contained:
   ```makefile
   build:
       @cd $(MODULE) && CGO_ENABLED=0 $(GO) build -ldflags "$(LDFLAGS)" -o ../$(BINARY) .
   ```
   While this was intended to place the binary in the repository root when invoked from the root directory, it contaminated the top-level repo or parent workspace if directory context shifted.

3. **Workspace Parent Scanning Defaults:**
   `run.ps1` default scanning tasks often target `$parentDir` (`scan $parentDir`), increasing the likelihood that automated runners or developers invoked build scripts with the working directory set to either the parent or the root.

## 3. Resolution & Safeguards

1. **Stray Binary Removal:**
   The obsolete `../gitmap.exe` binary in the parent directory was safely removed.

2. **Canonical Build Target Guard in `Makefile`:**
   Updated `Makefile` build and clean targets:
   ```makefile
   build:
       @mkdir -p bin
       @cd $(MODULE) && CGO_ENABLED=0 $(GO) build -ldflags "$(LDFLAGS)" -o ../bin/$(BINARY) .
       @echo "Built bin/$(BINARY) ($(VERSION))"

   clean:
       @rm -f bin/$(BINARY) bin/$(BINARY).exe $(BINARY) $(BINARY).exe
       @rm -rf $(MODULE)/.gitmap/release-assets
       @echo "Cleaned."
   ```

3. **Total Ban in `.lovable/strictly-avoid.md`:**
   Added a strict rule prohibiting outputting or depositing binaries in parent directories or outside `bin/` and user AppData directories.
