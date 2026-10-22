# Subtask 05: Root Utility Help Error Propagation & Return Formatting

> **Subtask ID:** Subtask-05  
> **Parent Plan:** `.ai-memory/plans/291-improvements-planning-and-execution.md`  
> **Target Path:** `.ai-memory/plans/subtasks/291-improvements-planning-and-execution/subtask-05.md`  
> **Status:** PENDING  

---

## 1. Objectives & Scope

1. Modernize error handling in `cli/cmd/rootutility_help.go` so that help subcommands bubble up errors rather than discarding them.
2. Convert void helper routines (`dispatchHelpTopic`, `dispatchHelpExample`) into functions returning explicit `error` types.
3. Propagate errors from underlying documentation loaders (`helpdoc.ReadRaw`, `helpdoc.PrintWithMode`) to root dispatch.
4. Enforce strict vertical whitespace hygiene:
   - Mandatory blank line before every `return` statement.
   - Mandatory blank line after closing brace `}`.
5. Enforce 15-line limit per function and positive boolean naming (`isHelpExample`, `hasTopic`).

---

## 2. Detailed Technical Plan

### 2.1 Refactor `runHelpDispatch`
- In `cli/cmd/rootutility_help.go`:
  - Locate `runHelpDispatch() error`.
  - When `isHelpExample` is true:
    - Invoke `return dispatchHelpExample(os.Args[3:])`.
  - When `isHelpTopicCmd` is true:
    - Invoke `return cmdhelp.RunHelpTopic(os.Args[3:])` (with proper appfault-to-error conversion if required).
  - When `hasTopic` is true:
    - Invoke `return dispatchHelpTopic(os.Args[2])`.
  - Add required vertical blank lines before every `return` statement to conform with coding guidelines.

### 2.2 Refactor `dispatchHelpTopic` & `dispatchHelpExample` Signatures
- Convert `dispatchHelpTopic(rawTopic string)` to `dispatchHelpTopic(rawTopic string) error`:
  - Capture and inspect errors from `helpdoc.ReadRaw(topic)`:
    - If error indicates topic not found, safely fallback to `printUsageFiltered(topic)` and return `nil`.
    - If error is an IO or filesystem fault, return the error.
  - Return `nil` upon successful rendering.
- Convert `dispatchHelpExample(args []string)` to `dispatchHelpExample(args []string) error`:
  - Inspect errors from example renderers and return `nil` or error explicitly.

### 2.3 Guideline & Formatting Enforcement
- Ensure every function in `rootutility_help.go` is $\le 15$ lines.
- Fix all missing blank lines before returns:
  ```go
  if isHelpExample {
      return dispatchHelpExample(os.Args[3:])
  }

  if isHelpTopicCmd {
      return cmdhelp.RunHelpTopic(os.Args[3:])
  }
  ```
- Eliminate any nested `if` blocks.

---

## 3. Verification Criteria

- [ ] `runHelpDispatch`, `dispatchHelpTopic`, and `dispatchHelpExample` return `error` types.
- [ ] No documentation loading errors are silently swallowed without handling.
- [ ] Every `return` statement is preceded by an empty line.
- [ ] Every closing brace `}` followed by non-brace statements is followed by an empty line.
- [ ] All functions in `cli/cmd/rootutility_help.go` are $\le 15$ lines.
