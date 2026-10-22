# Subtask 01: IDE Tokens, Workdir Constants & Environment Utility Refactoring

## Status: Completed
- **Parent Plan**: `.ai-memory/plans/294-code-refactoring-and-coding-guideline-update.md`
- **Architecture Spec**: `02-spec/21-app/294-code-refactoring-and-coding-guideline-update/01-architecture-spec.md`
- **Owner**: Worker 01
- **Owned Files**:
  - `cli/strutil/strutil.go`
  - `cli/cmdide/ide_add.go`
  - `cli/constants/constants_workdir.go`
  - `cli/envutil/envutil.go`
  - `cli/cmdpe/pe_target.go`

---

## 1. Context & Objectives

This subtask addresses magic string elimination, environment variable access standardization, and CSV token construction across the IDE registration and workspace resolution subsystems.

### Key Objectives:
1. Implement `JoinCSV(tokens []string) string` in `cli/strutil/strutil.go`.
2. Refactor `formatRegisteredTokens` in `cli/cmdide/ide_add.go` to leverage `Target*` constants and `strutil.JoinCSV`.
3. Create `cli/constants/constants_workdir.go` defining `EnvWorkdir`, `DefaultWorkdirWorkspace`, and `DefaultWorkdirRepos`.
4. Create `cli/envutil/envutil.go` providing `GetTrimmed(key string) string` and `GetOrDefault(key, fallback string) string`.
5. Refactor `ResolveWorkdir` and `fallbackWorkdir` in `cli/cmdpe/pe_target.go` to use centralized constants and `envutil`, adding comprehensive doc comments explaining fallback behavior.

---

## 2. Invariants & Guardrails

- **Function Length Limit:** Maximum 15 lines per function (<= 8 lines preferred).
- **Conditionals:** Zero nested `if` statements; use guard clauses and early returns.
- **Boolean Naming:** Use positive prefixes `is*` or `has*` only. No negative booleans.
- **Error Sovereignty:** Leaf functions must return errors directly; zero `os.Exit` or `cliexit.HandleError` in leaf code.
- **Path Invariants:** Strictly relative Git paths only. No absolute paths or URI schemes.
- **Safety Gate:** Zero builds (`go build` prohibited) and zero tests (`go test` prohibited).

---

## 3. Step-by-Step Implementation Plan

### Step 1: Implement `JoinCSV` in `cli/strutil/strutil.go`
- **Location:** `cli/strutil/strutil.go`
- **Signature:**
  ```go
  // JoinCSV joins a slice of strings with a comma separator, skipping empty entries.
  func JoinCSV(tokens []string) string
  ```
- **Behavior:**
  - Iterate through `tokens`.
  - Filter out empty strings or whitespace-only strings.
  - Return `strings.Join(cleaned, ",")`.
  - Ensure total function lines <= 12 lines.

### Step 2: Refactor `formatRegisteredTokens` in `cli/cmdide/ide_add.go`
- **Location:** `cli/cmdide/ide_add.go`
- **Existing Code:**
  Lines 209-225 accumulate strings manually and use standard `strings.Join(tokens, ",")`.
- **Refactoring:**
  - Decompose into a token collector helper and a formatter:
    ```go
    func collectRegisteredTokens(res IDEActionResult) []string {
    	var tokens []string
    	if res.IsVSCodeAffected {
    		tokens = append(tokens, string(TargetVSCode))
    	}
    	if res.IsCursorAffected {
    		tokens = append(tokens, string(TargetCursor))
    	}
    	if res.IsAntigravityAffected {
    		tokens = append(tokens, string(TargetAntigravity))
    	}
    	if res.IsDesktopAffected {
    		tokens = append(tokens, string(TargetGitHubDesktop))
    	}

    	return tokens
    }

    func formatRegisteredTokens(res IDEActionResult) string {
    	tokens := collectRegisteredTokens(res)

    	return strutil.JoinCSV(tokens)
    }
    ```
  - Both functions remain <= 15 lines.
  - Imports: ensure `strutil` is cleanly imported.

### Step 3: Create `cli/constants/constants_workdir.go`
- **Location:** `cli/constants/constants_workdir.go`
- **Content:**
  ```go
  package constants

  // Environment variable keys for workspace resolution.
  const (
  	EnvWorkdir = "GITMAP_WORKDIR"
  )

  // Default folder names for workspace and repository storage.
  const (
  	DefaultWorkdirWorkspace = "workspace"
  	DefaultWorkdirRepos     = "repos"
  )
  ```

### Step 4: Create `cli/envutil/envutil.go`
- **Location:** `cli/envutil/envutil.go`
- **Package:** `envutil`
- **Signatures:**
  ```go
  package envutil

  import (
  	"os"
  	"strings"
  )

  // GetTrimmed retrieves the environment variable for key with surrounding whitespace trimmed.
  func GetTrimmed(key string) string {
  	return strings.TrimSpace(os.Getenv(key))
  }

  // GetOrDefault retrieves the trimmed environment variable for key, or fallback if empty.
  func GetOrDefault(key, fallback string) string {
  	val := GetTrimmed(key)
  	if val != "" {
  		return val
  	}

  	return fallback
  }
  ```
- **Rules:**
  - Functions <= 8 lines.
  - No nested conditionals.
  - Guard clauses and early returns.

### Step 5: Refactor `ResolveWorkdir` & `fallbackWorkdir` in `cli/cmdpe/pe_target.go`
- **Location:** `cli/cmdpe/pe_target.go`
- **Refactoring:**
  ```go
  // ResolveWorkdir determines the target working directory from flags, environment, or default path.
  func ResolveWorkdir(flagVal string) string {
  	clean := strings.TrimSpace(flagVal)
  	if clean != "" {
  		return clean
  	}

  	if envVal := envutil.GetTrimmed(constants.EnvWorkdir); envVal != "" {
  		return envVal
  	}

  	return fallbackWorkdir()
  }

  // fallbackWorkdir resolves the default workspace repository folder (~/workspace/repos).
  // If the user's home directory cannot be determined, it falls back to the current directory (".").
  func fallbackWorkdir() string {
  	home, homeErr := os.UserHomeDir()
  	if homeErr != nil {
  		return "."
  	}

  	return filepath.Join(home, constants.DefaultWorkdirWorkspace, constants.DefaultWorkdirRepos)
  }
  ```
- **Imports:** Add `"github.com/alimtvnetwork/gitmap-v28/cli/envutil"`.

---

## 4. Verification & Quality Checklist

- [x] `strutil.JoinCSV` skips empty strings and returns empty string on empty input.
- [x] `cli/cmdide/ide_add.go` uses `Target*` constants and calls `strutil.JoinCSV`.
- [x] `cli/constants/constants_workdir.go` declares `EnvWorkdir`, `DefaultWorkdirWorkspace`, `DefaultWorkdirRepos`.
- [x] `cli/envutil/envutil.go` has `GetTrimmed` and `GetOrDefault` under 15 lines.
- [x] `cli/cmdpe/pe_target.go` resolves workdir using `envutil.GetTrimmed(constants.EnvWorkdir)` and `fallbackWorkdir` uses constants.
- [x] All functions are <= 15 lines (<= 8 lines preferred).
- [x] Zero nested `if` statements.
- [x] All booleans use `is*` or `has*` positive prefixes.
- [x] Strictly relative Git paths only.
- [x] Zero compiler builds and zero tests executed.
