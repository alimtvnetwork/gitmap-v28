# 2026-09-30 SSH Join Re-Enrollment Auth Failure, Premature Interactive Prompt & Windows Double Shell Wrapping RCA

- **Spec Document:** [02-spec/22-app-issues/52-ssh-join-re-enrollment-auth-failure-and-double-shell-wrap-rca.md](../../02-spec/22-app-issues/52-ssh-join-re-enrollment-auth-failure-and-double-shell-wrap-rca.md)
- **Status:** Resolved
- **Impact Area:** `cli/cmdssh/sshjoin_enroll.go`, `cli/cmdssh/ssh_auth_key_deploy.go`, `cli/crypto/ssh_client.go`

## Summary

1. `gitmap ssh join` previously called an interactive password prompt prematurely before testing default SSH keys, local vault credentials (`gitmap.db`), or `vmpass.json` / `06-vmpass.json` configuration fallbacks.
2. `wrapCommandForShell` in `cli/crypto/ssh_client.go` unconditionally wrapped PowerShell commands in `powershell -NoProfile -Command "..."`, which caused nested quote syntax errors when the input script (`winAuthKeyDeployScript`) already began with `powershell`. As a result, the public key was never successfully written to `authorized_keys` during the initial join on Windows hosts.
3. Windows OpenSSH configurations vary in whether they check `__PROGRAMDATA__/ssh/administrators_authorized_keys` or `$HOME\.ssh\authorized_keys`. Writing exclusively to one location caused key authentication to fail for non-admin accounts or when administrative key matching was disabled.
4. Resolved by establishing a strict authentication precedence (Keys -> Vault -> Fallback Store -> Interactive Prompt), eliminating redundant outer PowerShell invocations, and writing authorized keys to both user and administrator locations on Windows.
