# 347 - GitMap User Change Command: Component Specification

- **Spec ID:** `02-spec/21-app/347-gitmap-user-change-command/02-component-spec.md`
- **Subsystem:** `cli/cmduser`, `cli/usercontext`, `cli/store`, `cli/model`
- **Status:** APPROVED
- **Target Version:** `v6.628.0`
- **Related Specs:**
  - `02-spec/21-app/347-gitmap-user-change-command/readme.md`
  - `02-spec/21-app/347-gitmap-user-change-command/01-architecture-spec.md`
  - `02-spec/21-app/347-gitmap-user-change-command/03-migration-and-verification-spec.md`
- **Related Plans:**
  - `.ai-memory/plans/pending/347-gitmap-user-change-command.md`
  - `.ai-memory/plans/subtasks/347-gitmap-user-change-command/02-autobot-feature.md`

---

## 1. Component Scope & Functional Architecture

This specification formalizes the component contracts, concrete data structures, and state transitions for user identity switching in GitMap. It specifically governs:
1. The `ChangeUserOptions` contract for parameter ingestion and scope resolution.
2. The Autobot profile structure (`ID: 'autobot'`, `Name: 'AI automation bot'`, `Email: 'local@local'`, `Type: 'automation agent'`).
3. Enhanced Git identity management in `cli/usercontext/git_identity.go`, supporting empty email unsetting via `UnsetGitConfigKey`.
4. Automated autobot profile persistence and discovery within `cli/store/git_profile_store.go` via `LoadGitProfiles()`.

### 1.1 Subsystem Boundary Diagram

```
                 +-------------------------------------------------------------+
                 |                  CLI Command Layer                          |
                 |      `gitmap change-user`, `gitmap autobot`                 |
                 +-------------------------------------------------------------+
                                        |
                                        v
                 +-------------------------------------------------------------+
                 |                 cli/cmduser (Handlers)                      |
                 |  - RunUserChange(args)      - RunUserAutobot(args)          |
                 |  - parseChangeUserOptions   - parseAutobotOptions           |
                 +-------------------------------------------------------------+
                                  /            \
                                 /              \
                                v                v
+-----------------------------------+        +-----------------------------------+
|      cli/store (Persistence)      |        |   cli/usercontext (Git Identity)  |
| - LoadGitProfiles()               |        | - SetGitIdentity(name, email, g)  |
| - SaveGitProfiles(cfg)            |        | - UnsetGitConfigKey(ctx, key, g)  |
| - ensureAutobotProfile(&cfg)      |        | - ApplyGitIdentityWithOptions     |
+-----------------------------------+        +-----------------------------------+
                 |                                             |
                 v                                             v
     ~/.gitmap/data/git_profiles.json             git config [--global] key value
```

---

## 2. Data Models & Type Contracts

All data structures adhere strictly to the repository standards: PascalCase exported identifiers, positive booleans (`is*`, `has*` only; zero negative booleans), and explicit JSON serialization tags.

### 2.1 `ChangeUserOptions`

Located in `cli/cmduser/types.go` (or `cli/usercontext/types.go`), `ChangeUserOptions` encapsulates the validated parameters required to switch or reconfigure a Git identity:

```go
package cmduser

// ChangeUserOptions encapsulates parameters for switching or setting Git identities.
type ChangeUserOptions struct {
	Name          string `json:"name"`
	Email         string `json:"email"`
	ProfileAlias  string `json:"profileAlias,omitempty"`
	IsGlobal      bool   `json:"isGlobal"`
	IsProject     bool   `json:"isProject"`
	HasEmptyEmail bool   `json:"hasEmptyEmail"`
	IsUnsetEmail  bool   `json:"isUnsetEmail"`
}
```

#### Field Specifications:
- `Name` (`string`): The Git commit author name (e.g., `"AI automation bot"` or `"MD ALIM UL KARIM"`).
- `Email` (`string`): The Git commit author email (e.g., `"local@local"`). If empty and `HasEmptyEmail` or `IsUnsetEmail` is true, the email will be unset from the Git configuration.
- `ProfileAlias` (`string`): Associated profile alias (e.g., `"autobot"`, `"work"`, `"personal"`).
- `IsGlobal` (`bool`): When true, modifications apply globally (`git config --global`).
- `IsProject` (`bool`): When true, binds the profile locally to the current repository directory.
- `HasEmptyEmail` (`bool`): Explicit flag indicating the user supplied an empty email string.
- `IsUnsetEmail` (`bool`): When true, instructs the runner to unset `user.email` in Git.

### 2.2 Autobot Profile Structure

The Autobot profile is a first-class entity within `model.GitProfile` (`cli/model/git_profile.go`):

```go
package model

// Autobot Profile Canonical Definition:
// ID:         "autobot"
// Name:       "AI automation bot"
// Email:      "local@local"
// Type:       "automation agent"
// Provider:   "git"
// AuthMethod: "config"
// IsDefault:  false
```

#### Field Contract Invariants:
| Field | Value | Description |
|---|---|---|
| `ID` | `"autobot"` | Unique profile identifier and lookup key |
| `Name` | `"AI automation bot"` | Git author display name |
| `Email` | `"local@local"` | Fallback local non-routable email |
| `Provider` | `"git"` | System provider identifier |
| `Type` | `"automation agent"` | Profile type distinguishing autonomous agents from human users |
| `AuthMethod` | `"config"` | Git config-based identity attachment |
| `IsDefault` | `false` | Defaults to non-default until explicitly assigned |
| `UsageCount` | `0` (initial) | Incremented on active switch |

---

## 3. `cli/usercontext/git_identity.go` Enhancements

### 3.1 Empty Email Handling & The Zero-No-Op Dilemma

In the existing implementation of `cli/usercontext/git_identity.go`, `writeConfigKey` contains the following guard clause:

```go
func writeConfigKey(ctx context.Context, key, value string, isGlobal bool) error {
	if len(value) == 0 {
		return nil
	}
	// ... executes git config ...
}
```

**Architectural Defect:** If a caller invokes identity update with an empty email string (`""`), `writeConfigKey` immediately returns `nil` without modifying Git configuration. Consequently, any pre-existing email remains in place, violating the user's explicit directive to clear or keep the email empty.

**Remediation:** Introduce explicit key unsetting support via `UnsetGitConfigKey` and branch execution based on positive boolean flags (`HasEmptyEmail` / `IsUnsetEmail`).

### 3.2 `UnsetGitConfigKey` Contract

```go
// UnsetGitConfigKey removes a configuration key from local or global Git configuration.
// If the key does not exist, it returns nil to ensure idempotent execution.
func UnsetGitConfigKey(ctx context.Context, key string, isGlobal bool) error {
	args := buildConfigUnsetArgs(key, isGlobal)
	out, err := execGitCommand(ctx, args...)
	if err != nil {
		return handleUnsetConfigError(out, err)
	}

	return nil
}

func buildConfigUnsetArgs(key string, isGlobal bool) []string {
	if isGlobal {
		return []string{"config", "--global", "--unset", key}
	}

	return []string{"config", "--unset", key}
}

func handleUnsetConfigError(out []byte, err error) error {
	if isKeyNotFoundError(out, err) {
		return nil
	}

	return apperror.WrapSimple(err, "git config unset failed: "+strings.TrimSpace(string(out)))
}

func isKeyNotFoundError(out []byte, err error) bool {
	msg := strings.TrimSpace(string(out))
	if strings.Contains(msg, "key does not contain") {
		return true
	}

	return strings.Contains(msg, "not in that section")
}
```

### 3.3 `ApplyGitIdentityWithOptions` Contract

To prevent function signature bloat and enforce parameter structs, introduce `ApplyGitIdentityWithOptions`:

```go
// ApplyGitIdentityWithOptions configures user.name and user.email according to options.
func ApplyGitIdentityWithOptions(opts ChangeUserOptions) error {
	ctx, cancel := context.WithTimeout(context.Background(), gitCommandTimeout)
	defer cancel()

	if err := writeConfigKey(ctx, "user.name", opts.Name, opts.IsGlobal); err != nil {
		return err
	}

	return applyEmailIdentity(ctx, opts)
}

func applyEmailIdentity(ctx context.Context, opts ChangeUserOptions) error {
	hasEmptyEmail := len(opts.Email) == 0 || opts.HasEmptyEmail || opts.IsUnsetEmail
	if hasEmptyEmail {
		return UnsetGitConfigKey(ctx, "user.email", opts.IsGlobal)
	}

	return writeConfigKey(ctx, "user.email", opts.Email, opts.IsGlobal)
}
```

#### Invariants Enforced:
1. **Leaf Error Sovereignty:** `ApplyGitIdentityWithOptions` and `UnsetGitConfigKey` return `error` directly. They never call `cliexit.HandleError` or `os.Exit`.
2. **15-Line Function Limit:** Each function adheres strictly to <= 15 lines (average 4-8 lines).
3. **Positive Booleans:** Boolean checks use positive predicates (`hasEmptyEmail`, `isGlobal`).
4. **Blank Line Hygiene:** A blank line precedes every `return` statement and follows every closing brace `}`.

---

## 4. Profile Persistence in `cli/store/git_profile_store.go`

### 4.1 Autobot Profile Discovery & Automatic Seeding

The Autobot profile must be registered within the profile store (`git_profiles.json`) so it is discoverable via `gitmap user list` and bindable to repositories.

In `cli/store/git_profile_store.go`, `LoadGitProfiles()` is enhanced to verify and seed the Autobot profile idempotently:

```go
// LoadGitProfiles loads existing profiles from disk or seeds defaults including autobot.
func LoadGitProfiles() (model.GitProfileConfig, error) {
	path := GitProfilesPath()
	data, err := os.ReadFile(path)
	if err != nil {
		return handleDefaultProfilesLoad()
	}

	cfg, unmarshalErr := unmarshalGitProfiles(data)
	if unmarshalErr != nil {
		return cfg, apperror.WrapSimple(unmarshalErr, "unmarshal git profiles:")
	}

	return ensureAndPersistAutobot(cfg)
}

func ensureAndPersistAutobot(cfg model.GitProfileConfig) (model.GitProfileConfig, error) {
	hasModified := ensureAutobotProfile(&cfg)
	if !hasModified {
		return cfg, nil
	}

	saveErr := SaveGitProfiles(cfg)

	return cfg, saveErr
}
```

### 4.2 `ensureAutobotProfile` Implementation

```go
func ensureAutobotProfile(cfg *model.GitProfileConfig) bool {
	for idx := range cfg.Profiles {
		if cfg.Profiles[idx].ID == "autobot" {
			return false
		}
	}

	cfg.Profiles = append(cfg.Profiles, newAutobotProfile())

	return true
}

func newAutobotProfile() model.GitProfile {
	return model.GitProfile{
		ID:         "autobot",
		Name:       "AI automation bot",
		Email:      "local@local",
		Provider:   "git",
		Type:       "automation agent",
		AuthMethod: "config",
		IsDefault:  false,
		UsageCount: 0,
		LastUsedAt: time.Now(),
	}
}
```

### 4.3 Persistence Guarantees
- **Atomic File Writing:** Preserves JSON envelope formatting with `jsonx.NewEnvelope`.
- **Non-Destructive:** Existing active profile, defaults, and project bindings are left untouched.
- **Auto-Healing:** If a user deletes the autobot profile accidentally, the next call to `LoadGitProfiles()` automatically re-seeds it.

---

## 5. Command Runners in `cli/cmduser/`

### 5.1 `RunUserAutobot`

Implements the dedicated runner for `gitmap autobot`:

```go
package cmduser

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"github.com/alimtvnetwork/gitmap-v28/cli/usercontext"
)

// RunUserAutobot switches the Git identity to the Autobot profile.
func RunUserAutobot(args []string) error {
	opts := parseAutobotOptions(args)
	if err := usercontext.ApplyGitIdentityWithOptions(opts); err != nil {
		return err
	}

	recordAutobotUsage()
	renderAutobotSuccess(opts)

	return nil
}

func parseAutobotOptions(args []string) ChangeUserOptions {
	email := resolveAutobotEmail(args)
	hasEmpty := hasUserFlag(args, "--empty-email") || hasUserFlag(args, "--no-email")
	isGlobal := hasUserFlag(args, "--global") || hasUserFlag(args, "-g")
	isProject := hasUserFlag(args, "--project") || hasUserFlag(args, "-p")

	return ChangeUserOptions{
		Name:          "AI automation bot",
		Email:         email,
		ProfileAlias:  "autobot",
		IsGlobal:      isGlobal,
		IsProject:     isProject,
		HasEmptyEmail: hasEmpty,
		IsUnsetEmail:  hasEmpty,
	}
}

func resolveAutobotEmail(args []string) string {
	for _, a := range args {
		if !strings.HasPrefix(a, "-") && strings.Contains(a, "@") {
			return a
		}
	}

	return "local@local"
}

func renderAutobotSuccess(opts ChangeUserOptions) {
	scope := "local repository"
	if opts.IsGlobal {
		scope = "global system"
	}

	if opts.HasEmptyEmail {
		fmt.Printf("✔ Switched %s Git identity to %s (email unset)\n", scope, opts.Name)
		return
	}

	fmt.Printf("✔ Switched %s Git identity to %s <%s>\n", scope, opts.Name, opts.Email)
}
```

### 5.2 `RunUserChange`

Implements the general one-liner switch runner for `gitmap change-user` and `gitmap change-login`:

```go
// RunUserChange switches or sets the Git identity from one-liner arguments.
func RunUserChange(args []string) error {
	opts, err := parseChangeUserOptions(args)
	if err != nil {
		return err
	}

	if err := usercontext.ApplyGitIdentityWithOptions(opts); err != nil {
		return err
	}

	renderChangeSuccess(opts)

	return nil
}
```

---

## 6. Coding Guidelines & Invariant Checklist

| Invariant | Standard | Enforcement Verification |
|---|---|---|
| **Relative Paths Only** | Only relative Git paths in docs and errors | `02-spec/...`, `cli/...`, `.ai-memory/...` verified |
| **Positive Booleans** | `is*`, `has*` only; zero negative prefixes | `IsGlobal`, `IsProject`, `HasEmptyEmail`, `IsUnsetEmail` |
| **Leaf Error Sovereignty** | Return error directly, no CLI exit calls | `return err`, no `cliexit.HandleError` in leaf functions |
| **15-Line Function Limit** | Every function <= 15 lines (<= 8 lines preferred) | Verified across all proposed helpers |
| **Return Blank Lines** | Blank line before `return` and after `}` | Strictly applied across all Go definitions |
| **Zero Nested Ifs** | Flatten conditional trees with guard returns | Zero nested `if` statements throughout |

---

## 7. Acceptance Criteria

1. **Autobot Seeding:** Calling `store.LoadGitProfiles()` loads or seeds `autobot` with `Name: "AI automation bot"`, `Email: "local@local"`, and `Type: "automation agent"`.
2. **Empty Email Handling:** When `HasEmptyEmail` or `IsUnsetEmail` is true, `usercontext.UnsetGitConfigKey` executes `git config [--global] --unset user.email` without error.
3. **One-Liner Autobot Invocation:** Running `gitmap autobot` updates the active git configuration to `"AI automation bot <local@local>"`.
4. **Unset Email Invocation:** Running `gitmap autobot --empty-email` sets `user.name` to `"AI automation bot"` and unsets `user.email`.
5. **Unit Test Pass:** All tests in `cli/usercontext/git_identity_test.go` and `cli/cmduser/user_test.go` succeed with zero regressions.
