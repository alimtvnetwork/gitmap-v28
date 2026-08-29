# Specification: Gitmap SSH Joiner, Coding Guidelines, and Install Extensions

Provide a brief description of the problem, any background context, and what the change accomplishes.
This specification details the implementation of several new features for the Gitmap CLI:
1. **Extended `gitmap install`**: Add support for installing custom tools (Scripts Fixer, Coding Guidelines, Macro AHK) and listing installed tools (`gitmap install ls`).
2. **Coding Guidelines Manager (`gitmap cg`)**: Install and update coding guidelines on a per-repository basis, with parallel execution across workspaces.
3. **SSH Joiner (`gitmap sj`)**: Manage SSH connections, save credentials securely in SQLite, and support import/export.
4. **SSH Executor (`gitmap se`)**: Execute shell commands concurrently across joined SSH machines.
5. **Help Dashboard (`gitmap hd`)**: Update the React-based UI to include documentation, search, and examples for all new commands.

## User Review Required

> [!IMPORTANT]
> The specification is broken down into 83 actionable steps over 8 phases to ensure nothing is missed and end-to-end testing is conducted at every phase.
> Please review the actionable items below and approve so the implementation can begin.

## Proposed Changes

### Phase 1: Database and Data Model

1. Define `SSHConnection` struct in `gitmap/data/models.go` (ID, Alias, IP, Username, EncryptedPassword, KeyPath, OS).
2. Create SQLite migration file (`gitmap/db/migrations/005_ssh_connections.sql`) to create `ssh_connections` table.
3. Update `gitmap/db/schema.go` to include the new table.
4. Implement `db.InsertSSHConnection(ctx, conn)`.
5. Implement `db.GetSSHConnections(ctx)`.
6. Implement `db.DeleteSSHConnection(ctx, alias)`.
7. Define `RepoCodingGuideline` struct in `models.go` to track version history per repo.
8. Create SQLite migration (`006_repo_cg_versions.sql`) for `repo_cg_versions` table.
9. Implement DB methods for updating and fetching CG versions per repo.
10. Add unit tests for SSH DB operations.
11. Add unit tests for CG version DB operations.

### Phase 2: Crypto & Security Utilities

12. Create `gitmap/crypto/encrypt.go` for symmetric encryption (AES-GCM) of SSH passwords.
13. Implement `crypto.Encrypt(plaintext, key)` using a machine-specific key or user-provided passphrase.
14. Implement `crypto.Decrypt(ciphertext, key)`.
15. Create `gitmap/crypto/ssh_client.go` to wrap `golang.org/x/crypto/ssh`.
16. Implement `ssh_client.ConnectWithPassword(ip, user, password)`.
17. Implement `ssh_client.ConnectWithKey(ip, user, keyPath)`.
18. Implement `ssh_client.RunCommand(session, cmd, shellType)`.
19. Add unit tests for encryption/decryption.
20. Add mock SSH server tests for `ssh_client`.

### Phase 3: `gitmap install` Extensions

21. Locate `gitmap/cmd/install.go` and `gitmap/cmd/installtools.go`.
22. Define a map/registry of custom install scripts (Scripts Fixer, Coding Guidelines, Macro AHK) with their Windows and Unix URLs.
23. Create `gitmap/cmd/install_custom.go`.
24. Implement `installCustomTool(name, osType)` to download and pipe to `iex` or `bash`.
25. Integrate custom tools into the main `runInstall` switch statement.
26. Add `ls` / `list` subcommand to `install`.
27. Implement `checkInstalledStatus(name)` to verify if a tool is already on the machine (check PATH or default directories).
28. Design the `gitmap install ls` UI using `charmbracelet/lipgloss` for tabular, colorful output.
29. Implement the Lipgloss rendering function in `gitmap/tui/install_list.go`.
30. Add end-to-end test for `gitmap install ls`.
31. Add end-to-end test for installing custom tools (mocking the HTTP endpoints).

### Phase 4: `gitmap cg` (Coding Guidelines Manager)

32. Create `gitmap/cmd/cg.go` and register aliases (`coding-guide`, `coding-guidelines`, `cg`).
33. Parse subcommands: `install`, `update`.
34. Parse flags: `--all`, `--exclude`, and variadic positional arguments (repo aliases/paths).
35. Implement `resolveRepos(args, allFlag, excludeCSV)` to fetch target repositories from the database.
36. Create `gitmap/cmd/cg_worker.go` for concurrent execution.
37. Implement a worker pool pattern (e.g., `errgroup.Group`) to run installations in parallel.
38. For each repo, detect OS, construct the install command (`irm ... | iex` or `curl ... | bash`), and set `exec.Command` Dir to the repo path.
39. Capture stdout/stderr from each parallel execution.
40. Upon success, write `{"version": "v24", "installed_at": "..."}` to `<repo>/.gitmap/cg-version.json`.
41. Update the SQLite `repo_cg_versions` table with the new version.
42. Design a colorful progress UI using Lipgloss (e.g., a spinner or progress bar for each repo).
43. Implement `tui.RenderCGProgress(results)`.
44. Add end-to-end tests for `gitmap cg install --all` (using temporary mock repos).

### Phase 5: `gitmap sj` (SSH Joiner)

45. Create `gitmap/cmd/sshjoin.go` and register aliases (`ssh-joiner`, `ssh-join`, `sj`).
46. Parse flags: `--import`, `--export`, and positional args.
47. Implement `runSSHJoinLs(args)` to fetch and display connections from the DB.
48. Design Lipgloss UI for `sj ls` (table showing Alias, IP, User, OS, Status).
49. Implement interactive prompt for adding a new machine using `fmt.Scan` or `charmbracelet/huh` if available.
50. Hash/Encrypt the entered password using the crypto module.
51. Save the new connection to the DB.
52. Implement JSON Export: fetch all records, decrypt (or export encrypted with a master password), serialize to JSON, write to file.
53. Implement JSON Import: read JSON, validate schema, encrypt/store in DB.
54. Add end-to-end tests for `gitmap sj ls`, interactive join (via simulated stdin), import, and export.

### Phase 6: `gitmap se` (SSH Executor)

55. Create `gitmap/cmd/sshexec.go` and register aliases (`ssh-exec`, `ssh-execute`, `se`).
56. Parse shell type (`ps`, `sh`, `cmd`) and the command string.
57. Parse `--exclude` flag.
58. Fetch active SSH connections from the DB, applying exclusions.
59. Initialize an `errgroup` for parallel SSH execution.
60. For each machine, decrypt the password, establish the SSH connection.
61. Wrap the command string based on the shell type (e.g., `pwsh -NoProfile -Command "..."`).
62. Execute the command over the SSH session and capture output.
63. Design a terminal UI to display real-time output from multiple machines, prefixing lines with the machine alias/IP (colorful prefixes via Lipgloss).
64. Handle connection timeouts and execution errors gracefully, logging them in red.
65. Add end-to-end tests for `gitmap se ps "echo test"` using a local mock SSH server.

### Phase 7: UI & Help Dashboard (`gitmap hd`) Updates

66. Navigate to `src/data/commands.ts` (or equivalent) in the React frontend.
67. Add metadata for `gitmap install scripts-fixer`, `coding-guidelines`, `macro-ahk`, and `install ls`.
68. Add metadata for `gitmap cg` (all aliases, flags, examples).
69. Add metadata for `gitmap sj` (joiner, ls, import, export).
70. Add metadata for `gitmap se` (shell types, exclusion, examples).
71. Ensure search indexes these new commands.
72. Update Markdown help files in `gitmap/helptext/`:
    - `install.md` (add custom scripts and `ls`).
    - `coding-guidelines.md` (new file).
    - `ssh-join.md` (new file).
    - `ssh-exec.md` (new file).
73. Ensure the embedded help system (`go:embed`) picks up the new `.md` files.
74. Verify that `gitmap hd` serves the updated UI with the new commands and install script examples.

### Phase 8: E2E Testing, Finalization, and Release

75. Run the full Go test suite (`go test ./... -v -race`).
76. Run Vite frontend tests (`npm run test`).
77. Validate standard `gitmap install` (e.g., node, go) still works without regression.
78. Verify parallel execution limits (throttle `gitmap cg` and `gitmap se` to max 10-20 concurrent goroutines to avoid exhausting file descriptors).
79. Check Lipgloss UI rendering on both Windows (CMD/PowerShell) and Unix (Bash/Zsh) for color compatibility.
80. Perform a manual dry-run release check.
81. Commit all changes with descriptive commit messages following the coding guidelines.
82. Ensure the code passes all linter checks (`golangci-lint run`, `npm run lint`).
83. Write the final release notes in `changelog.md`.

## Verification Plan

### Automated Tests

- `go test ./...`
- `npm run test`

### Manual Verification

- Manually test `gitmap cg`, `gitmap sj`, `gitmap se` on test infrastructure.
