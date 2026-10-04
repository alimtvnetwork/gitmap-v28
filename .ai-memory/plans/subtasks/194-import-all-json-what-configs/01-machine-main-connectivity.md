# Subtask 01: Machine `main` Connectivity & SSH Join Disambiguation
- Discovered root cause: positional parsing assumed `[user@host] [alias]` without verifying if positional 0 was alias and positional 1 was `user@host`.
- Password ending in `@` (`rtyrty123@`) was incorrectly treated as a multi-target address by naive `strings.Contains(s, "@")`.
- Refactored `isTargetAddress` to `isUserAtHostAddress` checking host validity after `@`.
- Added symmetric resolution in `resolveJoinPositionalTarget` in `cli/cmdssh/ssh_parser.go` and `cli/cmdssh/sshjoin_add_pass_cmd.go`.
- Tested and verified live connectivity to `main` (`node-main`).
- Status: Completed.
