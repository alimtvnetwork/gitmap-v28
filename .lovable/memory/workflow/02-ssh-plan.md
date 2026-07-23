# SSH Key Management — Implementation Plan

Spec: `spec/01-app/50-ssh-keys.md`

## Atomic Tasks (ordered)

### Phase 1: Foundation

1. **Create `constants/constants_ssh.go`** — all SSH command names, flag names, messages, errors, config markers, defaults.
2. **Create `model/sshkey.go`** — `SSHKey` struct with ID, Name, PrivatePath, PublicKey, Fingerprint, Email, CreatedAt.
3. **Create `store/sshkey.go`** — `CreateSSHKeysTable`, `InsertSSHKey`, `SelectSSHKeyByName`, `SelectAllSSHKeys`, `DeleteSSHKeyByName`, `UpdateSSHKey`.
4. **Register `SSHKeys` table in `store/store.go` Migrate()** — add `CreateSSHKeysTable` call.

### Phase 2: Core Commands

5. **Create `cmd/ssh.go`** — main dispatch: parse subcommand (cat/list/delete/config/default→generate), register in root dispatch.
6. **Create `cmd/sshgen.go`** — key generation: resolve email, check existence, prompt (R/N/C), run `ssh-keygen`, read .pub, store in DB, update SSH config.
7. **Create `cmd/sshcat.go`** — look up key by `--name`, print public key content, error if not found with available names.
8. **Create `cmd/sshlist.go`** — query all keys, format as aligned table (Name, Path, Fingerprint, Created).
9. **Create `cmd/sshdelete.go`** — confirm, delete DB record, optionally delete files, update SSH config.
10. **Create `cmd/sshconfig.go`** — read all keys from DB, regenerate managed block in `~/.ssh/config`, print result.

### Phase 3: SSH Config Management

11. **Implement SSH config parser in `cmd/sshconfig.go`** — read `~/.ssh/config`, find managed block by markers, replace or append. Preserve user entries outside markers.
12. **Implement Host entry generation** — single key (default) uses plain `github.com`, multiple keys use `github.com-<name>` pattern.

### Phase 4: Clone Integration

13. **Add `--ssh-key` / `-K` flag to `cmd/clone.go`** — look up key name from DB, set `GIT_SSH_COMMAND` env var on clone subprocess.
14. **Print SSH clone guidance** — when using non-default key, show how the SSH config Host maps to remote URLs.

### Phase 5: Help & Completion

15. **Create `helptext/ssh.md`** — alias, usage, subcommands, flags, prerequisites, 3 examples, see-also.
16. **Register `ssh` in root dispatch** — add to `dispatchCore` or `dispatchUtility` in `root.go`.
17. **Add `--list-ssh-keys` flag to completion command** — returns key names, one per line.
18. **Update `completion/powershell.go`** — add SSH subcommand and `--ssh-key` flag completions.
19. **Update `completion/bash.go`** — add SSH completions.
20. **Update `completion/zsh.go`** — add SSH completions.

### Phase 6: Documentation

21. **Update `README.md`** — add SSH section to command table.
22. **Update docs site** — add SSH page (`src/pages/SSH.tsx`), route, sidebar entry.
23. **Update `src/data/commands.ts`** — add SSH command entry with flags and examples.
24. **Create release JSON** — `.gitmap/release/v2.25.0.json` with SSH changelog entries.
