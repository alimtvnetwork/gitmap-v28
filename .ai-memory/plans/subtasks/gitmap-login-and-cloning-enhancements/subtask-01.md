# Subtask 01: SSH Probe Fix, Direct SSH Clone Bypass, and `--ssh` Pipeline Routing

- **Task ID:** subtask-login-clone-01
- **Parent Task:** gitmap-login-and-cloning-enhancements
- **Status:** DONE
- **Owner:** Worker 01

## Objectives
1. In `cli/cmdclone/clone_auth_probe.go`:
   - In `executeSSHProbe`: remove the conflicting `-h` flag from `git ls-remote --exit-code -h <url> HEAD` so it runs `git ls-remote --exit-code <url> HEAD`.
2. In `cli/cmdclone/clone_auth_resolve.go`:
   - In `ResolveRepoAuth`: check if `isSSHCloneURL(rawURL)`. If true, set `rec.Transport = "ssh"`, `rec.SSHUrl = rawURL`, and return `rec, nil` immediately without probing or falling back to HTTPS/PAT prompts!
3. In `cli/cmdclone/clone.go`:
   - In `DirectCloneParams`: add `UseSSH bool`.
   - In `dispatchCloneExecution`: pass `UseSSH: cf.UseSSH` to `DirectCloneParams`.
   - In `executeDirectClone`: if `isSSHCloneURL(url)` or `params.UseSSH`, skip `ResolveRepoAuth` and execute the clone directly via git.
4. In `cli/cmdclone/clonemulti.go`:
   - In `executeDirectCloneOne`: if `isSSHCloneURL(url)`, skip `ResolveRepoAuth` and execute directly.
5. In `cli/cmdclone/clonefixrepo.go`:
   - Ensure `cfr` and `cfrp` preserve SSH transport when supplied with SSH URLs.
6. Verify all modified Go files remain strictly <= 300 lines with nesting depth <= 1 (extracting helper files if needed).
