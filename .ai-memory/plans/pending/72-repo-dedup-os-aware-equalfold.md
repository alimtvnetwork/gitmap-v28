# Execution Plan: 72-repo-dedup-os-aware-equalfold

## User Request (Verbatim)

```text
Do you have the code to improve the optimize the redundant repos that has been created? That's it. Do you have the code to do that? Confirm it. And also, I've seen in the code you are trying to compare the code in lowercase. Try not to do that. the specific method strings unfold, which is a lot more faster. Try to use that when you are comparing two strings without the case sensitivity. So try to apply that in terms of check-in, because in Unix, the different paths, different cases actually mean same thing. So you can only check ignoring the path in Windows. So you need to understand which OS you are in. So I think that is a bug we need to fix, and again, make a release and make sure you also optimize in the next pool if there is a redundancy. Okay? You can find the redundancy from your SQLite database, not from file system. Okay, so in future, when you add new files, you make sure that it is unique. Do you understand?
```

## Scope and Subtasks

- **Task-01: OS-Aware Store Indexes & Path Normalization**
  - Files: `cli/store/store.go`, `cli/store/storage_inventory.go`, `cli/vscodepm/path_filter.go`, `cli/workspacesync/path_guard.go`
  - Assigned: Worker 01

- **Task-02: Zero-Allocation `strutil.EqualFoldAny` and CLI Refactoring**
  - Files: `cli/strutil/strutil.go`, `cli/cmdvscode/vscode_cmd.go`, `cli/cmdvmware/vmware.go`
  - Assigned: Worker 02
