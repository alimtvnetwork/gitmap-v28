# Subtask Plan 02: Eliminate Swallowed Database Errors in `dispatchAlias` (`cli/cmd/alias.go`)

- **Subtask ID:** `subtask-02`
- **Subtask Name:** `02-alias-dispatch-error-propagation`
- **Parent Spec:** `02-spec/21-app/290-improvements-planning-and-execution/01-architecture-spec.md`
- **Assigned Agent Role:** Spec Author 01 / Implementation Worker 01
- **Target Files:**
  - `cli/cmd/alias.go`
- **Status:** Completed

---

## 1. Targeted Problem & Architectural Rationale

In `cli/cmd/alias.go`, `dispatchAlias(sub string, args []string) *apperror.AppError` handles routing for repository alias commands (`set`, `rm`, `list`, `show`, `suggest`).

Currently, `dispatchAlias` swallows errors across 4 subcommands:
```go
	if sub == constants.SubCmdAliasRm {
		runAliasRemove(args)
		return nil
	}
	if sub == constants.SubCmdAliasList {
		runAliasList()
		return nil
	}
	if sub == constants.SubCmdAliasShow {
		runAliasShow(args)
		return nil
	}
	if sub == constants.SubCmdAliasSug {
		runAliasSuggest(args)
		return nil
	}
```
All four handlers (`runAliasRemove`, `runAliasList`, `runAliasShow`, `runAliasSuggest`) return `error`. However, `dispatchAlias` discards the returned `error` and unconditionally returns `nil`.

As a consequence:
1. If the database connection fails or a table write error occurs, the command falsely reports success to the terminal.
2. The audit log (`finishCommandAudit`) receives a `nil` error and marks the command run as successful.
3. This directly violates **Rule 1 (Error Return Sovereignty)**.

---

## 2. Component Boundaries & Code Signatures

### 2.1 Refactor `dispatchAlias`
- **File:** `cli/cmd/alias.go`
- **Signature:** `func dispatchAlias(sub string, args []string) *apperror.AppError`
- **Requirements:**
  - Delegate each subcommand to a dedicated dispatch helper that checks errors.
  - Keep `dispatchAlias` function length <= 15 lines.
  - Return `*apperror.AppError` for unknown subcommands.

```go
func dispatchAlias(sub string, args []string) *apperror.AppError {
	if sub == constants.SubCmdAliasSet {
		return dispatchAliasSet(args)
	}
	if sub == constants.SubCmdAliasRm {
		return dispatchAliasRemove(args)
	}
	if sub == constants.SubCmdAliasList {
		return dispatchAliasList()
	}
	if sub == constants.SubCmdAliasShow {
		return dispatchAliasShow(args)
	}
	if sub == constants.SubCmdAliasSug {
		return dispatchAliasSuggest(args)
	}

	return apperror.New(fmt.Sprintf(constants.ErrUnknownCommand, sub), "E9000", nil)
}
```

### 2.2 Implement Specialized Dispatch Helpers
- **File:** `cli/cmd/alias.go`
- **Requirements:**
  - Add `dispatchAliasRemove(args []string) *apperror.AppError`
  - Add `dispatchAliasList() *apperror.AppError`
  - Add `dispatchAliasShow(args []string) *apperror.AppError`
  - Add `dispatchAliasSuggest(args []string) *apperror.AppError`
  - Each helper must be <= 8 lines, check the returned `error`, wrap it using `apperror.WrapSimple(err, "alias ...")`, and return `nil` on success.

```go
func dispatchAliasRemove(args []string) *apperror.AppError {
	if err := runAliasRemove(args); err != nil {
		return apperror.WrapSimple(err, "alias remove")
	}

	return nil
}

func dispatchAliasList() *apperror.AppError {
	if err := runAliasList(); err != nil {
		return apperror.WrapSimple(err, "alias list")
	}

	return nil
}

func dispatchAliasShow(args []string) *apperror.AppError {
	if err := runAliasShow(args); err != nil {
		return apperror.WrapSimple(err, "alias show")
	}

	return nil
}

func dispatchAliasSuggest(args []string) *apperror.AppError {
	if err := runAliasSuggest(args); err != nil {
		return apperror.WrapSimple(err, "alias suggest")
	}

	return nil
}
```

---

## 3. Implementation Verification Checklist

- [x] `runAliasRemove` error captured and wrapped into `*apperror.AppError`.
- [x] `runAliasList` error captured and wrapped into `*apperror.AppError`.
- [x] `runAliasShow` error captured and wrapped into `*apperror.AppError`.
- [x] `runAliasSuggest` error captured and wrapped into `*apperror.AppError`.
- [x] All dispatch functions <= 8 lines.
- [x] Main `dispatchAlias` <= 15 lines.
- [x] Zero nested `if` statements.
- [x] Blank line before `return` and after `}` preserved.
