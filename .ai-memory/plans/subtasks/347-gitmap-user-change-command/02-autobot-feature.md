# Subtask 02: Autobot Identity Feature and Profile Registration

> **Sequence ID:** 347  
> **Parent Plan:** [.ai-memory/plans/pending/347-gitmap-user-change-command.md](../../pending/347-gitmap-user-change-command.md)  
> **Specification:** [02-spec/21-app/347-gitmap-user-change-command/02-component-spec.md](../../../../02-spec/21-app/347-gitmap-user-change-command/02-component-spec.md)  
> **Related Specs:**
> - [02-spec/21-app/347-gitmap-user-change-command/01-architecture-spec.md](../../../../02-spec/21-app/347-gitmap-user-change-command/01-architecture-spec.md)
> - [02-spec/21-app/347-gitmap-user-change-command/03-migration-and-verification-spec.md](../../../../02-spec/21-app/347-gitmap-user-change-command/03-migration-and-verification-spec.md)
> **Status:** Pending  
> **Category:** Feature & CLI Ergonomics  

---

## 1. Scope & Objectives

This subtask governs the implementation of the autonomous bot identity feature, automatic profile seeding, Git configuration unsetting, and command-line execution for Plan 347 (`gitmap-user-change-command`).

### Concrete Deliverables
1. **Autobot Profile Model & Store Seeding:**
   - Define canonical autobot identity attributes (`ID: "autobot"`, `Name: "AI automation bot"`, `Email: "local@local"`, `Type: "automation agent"`).
   - Ensure idempotent auto-seeding of the autobot profile in `cli/store/git_profile_store.go` during `LoadGitProfiles()`.
2. **Git Identity Empty Email Unsetting:**
   - Introduce `UnsetGitConfigKey` in `cli/usercontext/git_identity.go`.
   - Update identity application logic so providing an empty email or `--empty-email` explicitly unsets `user.email` in Git config rather than performing a no-op.
3. **Autobot Command Runner:**
   - Implement `RunUserAutobot` in `cli/cmduser/user_autobot_cmd.go`.
   - Support optional email override (e.g. `gitmap autobot "local@local"`), default fallback, and explicit email unsetting (`--empty-email`).
4. **Dispatcher Integration:**
   - Wire top-level routing in `cli/cmduser/user_cmd.go` and `cli/cmd/root.go` for `gitmap autobot` and `gitmap user autobot`.
5. **Unit Test Suite & Verification:**
   - Verify unsetting behavior, profile persistence, and argument parsing via unit tests in `cli/usercontext/git_identity_test.go` and `cli/cmduser/user_test.go`.

---

## 2. Step-by-Step Execution Plan

### Step 1: Model & Store Profile Seeding
- File: `cli/store/git_profile_store.go`
  - Author helper function `ensureAutobotProfile(cfg *model.GitProfileConfig) bool`.
  - Scan `cfg.Profiles` for existing entry with `ID == "autobot"`.
  - If missing, append the canonical autobot profile:
    ```go
    model.GitProfile{
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
    ```
  - Integrate `ensureAutobotProfile` into `LoadGitProfiles()` and `handleDefaultProfilesLoad()`.
  - If modified, persist changes via `SaveGitProfiles(*cfg)`.
  - Ensure all helper functions adhere strictly to the <= 15 lines rule.

### Step 2: Empty Email Unsetting in User Context
- File: `cli/usercontext/git_identity.go`
  - Implement `UnsetGitConfigKey(ctx context.Context, key string, isGlobal bool) error`:
    - Construct arguments via `buildConfigUnsetArgs(key, isGlobal)`.
    - Execute via `execGitCommand`.
    - Treat exit status for non-existent key as success (idempotent removal).
    - Wrap genuine errors using `apperror.WrapSimple`.
  - Implement `ApplyGitIdentityWithOptions(opts ChangeUserOptions) error`:
    - Set `user.name` if `len(opts.Name) > 0`.
    - If `opts.HasEmptyEmail || opts.IsUnsetEmail || len(opts.Email) == 0`: invoke `UnsetGitConfigKey(ctx, "user.email", opts.IsGlobal)`.
    - Else: invoke `writeConfigKey(ctx, "user.email", opts.Email, opts.IsGlobal)`.
  - Enforce leaf error sovereignty (no `cliexit.HandleError` or `os.Exit`).

### Step 3: Implement Autobot Command Runner
- File: `cli/cmduser/user_autobot_cmd.go`
  - Author `RunUserAutobot(args []string) error`:
    - Parse arguments via `parseAutobotOptions(args)`.
    - Apply Git identity via `usercontext.ApplyGitIdentityWithOptions(opts)`.
    - Render formatted confirmation output.
    - Return `error` directly.
  - Implement `parseAutobotOptions(args []string) ChangeUserOptions`:
    - Resolve scope flags: `isGlobal := hasUserFlag(args, "--global") || hasUserFlag(args, "-g")`, `isProject := hasUserFlag(args, "--project") || hasUserFlag(args, "-p")`.
    - Resolve email flags: `hasEmpty := hasUserFlag(args, "--empty-email") || hasUserFlag(args, "--no-email")`.
    - Extract positional email argument or default to `"local@local"`.
    - Return initialized `ChangeUserOptions`.
  - Maintain function length <= 15 lines per function.

### Step 4: Dispatcher Integration & Help Registration
- File: `cli/cmduser/user_cmd.go`
  - Update `DispatchUser(command string)`:
    - Check if command is autobot: `isAutobotCommand(command)`.
    - If true, invoke `RunUserAutobot(os.Args[2:])`.
  - Update `dispatchUserSubcommand(sub string, args []string)`:
    - Add `case "autobot", "bot":` routing to `RunUserAutobot(args)`.
- File: `cli/cmduser/user_help_menu.go`
  - Register `autobot [email] [--global] [--empty-email]` in the help menu.

### Step 5: Verification & Test Suite Execution
- File: `cli/usercontext/git_identity_test.go`
  - Author unit test `TestUnsetGitConfigKey_Success` verifying `--unset` argument dispatch.
  - Author unit test `TestApplyGitIdentityWithOptions_EmptyEmailUnsets` verifying that passing `HasEmptyEmail: true` executes `git config --unset user.email`.
  - Author unit test `TestApplyGitIdentityWithOptions_ConfiguresBoth` verifying standard name and email setting.
- File: `cli/cmduser/user_test.go`
  - Author unit test `TestDispatchUser_AutobotCommand` confirming CLI dispatch to autobot runner.
  - Author unit test `TestParseAutobotOptions_DefaultsAndFlags` verifying flag resolution.

---

## 3. Invariant & Coding Guideline Checklist

| Guideline | Standard | Compliance Check |
|---|---|---|
| **Relative Git Paths** | Strictly relative paths throughout docs and code | All file references use `cli/...`, `02-spec/...`, `.ai-memory/...` |
| **Positive Booleans** | Only `is*` or `has*` boolean prefixes | `isGlobal`, `isProject`, `hasEmptyEmail`, `isUnsetEmail`, `hasModified` |
| **Zero Nested Ifs** | Conditional flattening with guard clauses | Flattened early returns across all functions |
| **15-Line Function Limit** | Max 15 lines per function (<= 8 lines preferred) | Verified in design specs for all runner and parsing helpers |
| **Leaf Error Sovereignty** | Return errors directly to outer caller | No calls to `cliexit.HandleError` or `os.Exit` in leaf functions |
| **Return Blank Lines** | Blank line after `}` and blank line before `return` | Strictly adhered to across all code blocks |

---

## 4. Acceptance Criteria & Definition of Done

- [ ] `model.GitProfile` contains the `"autobot"` profile entry when `store.LoadGitProfiles()` executes.
- [ ] Running `gitmap autobot` configures local Git identity as `AI automation bot <local@local>`.
- [ ] Running `gitmap autobot --global` configures global Git identity as `AI automation bot <local@local>`.
- [ ] Running `gitmap autobot --empty-email` or `gitmap autobot "" --global` unsets `user.email` in Git configuration.
- [ ] `usercontext.UnsetGitConfigKey` executes idempotently without crashing if the key is already absent.
- [ ] All unit tests pass in `cli/usercontext/` and `cli/cmduser/`.
- [ ] No compilation errors or linter warnings under `gitmap py 03-ai-scripts/50-fastgate.py`.
