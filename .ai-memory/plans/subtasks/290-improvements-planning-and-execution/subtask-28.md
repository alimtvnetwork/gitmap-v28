# Subtask Plan 28: Auth Remediation Error Sovereignty & Positive Booleans (`cli/cmdfixauth/fixauth.go`)

> **Subtask ID:** `subtask-28`  
> **Parent Plan:** `.ai-memory/plans/pending/290-improvements-planning-and-execution.md`  
> **Parent Spec:** `02-spec/21-app/290-improvements-planning-and-execution/03-migration-and-verification-spec.md`  
> **Target Path:** `cli/cmdfixauth/fixauth.go`  
> **Status:** Completed  

---

## 1. Objectives & Scope
1. Fix negative boolean naming in `cli/cmdfixauth/fixauth.go`: replace `isNonGitRepoCWD` with positive `isGitRepo := cmdpull.IsGitRepoCWD()` and guard `if !isGitRepo`.
2. Rename non-standard boolean `keyExists` to positive boolean `hasKey`.
3. Enforce **Rule 1 (Error Return Sovereignty)** across leaf functions: refactor `fixAuthGenerate`, `runSSHKeygenEd25519`, `fixAuthBind`, and `fixAuthAnnounce` to return `error` directly instead of calling `cliexit.HandleError` or `os.Exit`.
4. Ensure all refactored routines adhere to the <= 15 line limit (Rule 2) and vertical whitespace styling (Rule 6).

---

## 2. Targeted Problem & Architectural Rationale
In `cli/cmdfixauth/fixauth.go`, multiple coding guideline violations exist:
1. **Negative Boolean:**
   ```go
   isNonGitRepoCWD := !cmdpull.IsGitRepoCWD()
   if isNonGitRepoCWD { ... }
   ```
   Violates Rule 4, which forbids negative booleans and requires positive polarity (`isGitRepo`).
2. **Non-Conforming Boolean Prefix:**
   `keyExists` does not use `is` or `has`. Must be renamed to `hasKey`.
3. **Leaf Error Sovereignty Violation:**
   `fixAuthGenerate`, `runSSHKeygenEd25519`, `fixAuthBind`, and `fixAuthAnnounce` invoke `cliexit.HandleError(...)`. Leaf functions must NEVER call exit handlers or terminate the process. Errors must bubble up to `RunFixAuth`, which returns `error` to the root dispatcher so the database audit logging can record the outcome.

---

## 3. Concrete Implementation Details

### 3.1 Refactor `RunFixAuth` in `cli/cmdfixauth/fixauth.go`
```go
// RunFixAuth parses flags and orchestrates the auth-fix pipeline.
func RunFixAuth(args []string) error {
	user, email, assumeYes, force := parseFixAuthFlags(args)
	if user == "" {
		fmt.Fprintln(os.Stderr, "✗ --user <github-username> is required")
		fmt.Fprintln(os.Stderr, "  example: gitmap fix-auth --user aukgit --email me@example.com")

		return apperror.NewSimple("--user <github-username> is required", "E_FIXAUTH_USER")
	}

	isGitRepo := cmdpull.IsGitRepoCWD()
	if !isGitRepo {
		fmt.Fprintln(os.Stderr, "✗ not a git repository (run inside the repo you want to fix)")

		return apperror.NewSimple("not a git repository (run inside the repo you want to fix)", "E_NOT_GIT_REPO")
	}

	return executeFixAuthPipeline(user, email, assumeYes, force)
}

func executeFixAuthPipeline(user, email string, assumeYes, force bool) error {
	keyPath := fixAuthKeyPath(user)
	if err := cmdssh.EnsureSSHDir(filepath.Dir(keyPath)); err != nil {
		return apperror.WrapSimple(err, "mkdir ~/.ssh")
	}

	if err := fixAuthGenerate(keyPath, resolveFixAuthEmail(email), assumeYes, force); err != nil {
		return err
	}

	if err := fixAuthBind(keyPath); err != nil {
		return err
	}

	if err := fixAuthAnnounce(user, keyPath); err != nil {
		return err
	}

	_, _ = cmdssh.SanitizeSSHConfigFile("")

	return nil
}
```

### 3.2 Refactor `fixAuthGenerate` and `validateKeyGeneration`
```go
func fixAuthGenerate(keyPath, email string, assumeYes, force bool) error {
	hasKey := false
	if _, err := os.Stat(keyPath); err == nil {
		hasKey = true
	}

	if !validateKeyGeneration(hasKey, keyPath, assumeYes, force) {
		return nil
	}

	if err := cmdssh.ValidateSSHKeygen(); err != nil {
		return apperror.WrapSimple(err, constants.ErrSSHKeygenMissing)
	}

	return runSSHKeygenEd25519(keyPath, email)
}

func validateKeyGeneration(hasKey bool, keyPath string, assumeYes, force bool) bool {
	if !hasKey {
		return true
	}

	if !force {
		fmt.Printf("• key already exists, reusing: %s\n", keyPath)

		return false
	}

	isConfirmed := assumeYes || confirmOverwrite(keyPath)
	if !isConfirmed {
		fmt.Println("• aborted; existing key kept")

		return false
	}

	_ = os.Remove(keyPath)
	_ = os.Remove(keyPath + ".pub")

	return true
}
```

### 3.3 Refactor `runSSHKeygenEd25519`, `fixAuthBind`, `fixAuthAnnounce`
```go
func runSSHKeygenEd25519(keyPath, email string) error {
	cmd := exec.Command(constants.SSHKeygenBin,
		"-t", "ed25519", "-C", email, "-f", keyPath, "-N", "")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return apperror.WrapSimple(err, "ssh-keygen")
	}

	return nil
}

func fixAuthBind(keyPath string) error {
	ref := sshBindRefForKey(keyPath)
	cmdStr := fmt.Sprintf("ssh -i %s -F /dev/null -o IdentitiesOnly=yes", ref)
	out, err := exec.Command("git", "config", "core.sshCommand", cmdStr).CombinedOutput()
	if err != nil {
		return apperror.WrapSimple(err, fmt.Sprintf("git config core.sshCommand: %s", out))
	}

	fmt.Printf("✓ pinned repo → %s\n", cmdStr)

	return nil
}

func fixAuthAnnounce(user, keyPath string) error {
	pub, err := os.ReadFile(keyPath + ".pub")
	if err != nil {
		return apperror.WrapSimple(err, "read public key")
	}

	trimmed := strings.TrimSpace(string(pub))
	fmt.Println("\n=== PUBLIC KEY (add this to GitHub) ===")
	fmt.Println(trimmed)
	fmt.Println("=======================================")
	cmdssh.CopyPubKeyAndAnnounce(trimmed)
	fmt.Printf("\nNext steps for %s:\n", user)
	fmt.Println("  1. Open https://github.com/settings/ssh/new")
	fmt.Println("  2. Paste the key above, save.")
	fmt.Println("  3. Run:  git push")

	return nil
}
```

---

## 4. Strict Coding Guidelines Compliance
- **Rule 1 (Error Return Sovereignty):** All leaf methods return `error` or `*apperror.AppError`. Zero calls to `cliexit.HandleError` or `cliexit.HandleSuccess`.
- **Rule 2 (Function Limit):** All refactored routines are <= 15 lines.
- **Rule 3 (Zero Nested Ifs):** Guards flatten conditional branches.
- **Rule 4 (Positive Booleans):** `isGitRepo` and `hasKey` used with positive polarity.
- **Rule 6 (Vertical Whitespace):** Exactly one blank line before `return` and after `}`.

---

## 5. Verification Criteria
- [x] No occurrences of `isNonGitRepoCWD` exist; replaced with `isGitRepo`.
- [x] No occurrences of `keyExists` exist; replaced with `hasKey`.
- [x] No calls to `cliexit.HandleError` in `fixauth.go`.
- [x] All functions <= 15 lines.
- [x] All file paths are strictly relative lowercase Git paths.
