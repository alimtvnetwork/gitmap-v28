# SSH Key Management

## Overview

The `ssh` command lets users generate, store, view, regenerate, and manage
SSH keys directly from gitmap. Keys are stored by name in the database,
auto-managed in `~/.ssh/config`, and integrated with `gitmap-v28 clone` via
the `--ssh-key` flag.

---

## Goals

1. **One-command key generation** — `gitmap-v28 ssh` creates an RSA-4096 key
   with no passphrase at the default SSH location.
2. **Named keys** — each key gets a user-defined label (e.g., `work`,
   `personal`) stored in the database alongside the public key content.
3. **Quick access** — `gitmap-v28 ssh cat` prints the public key for copying.
4. **Regeneration safety** — if a key with the same name/path exists,
   prompt the user before overwriting.
5. **Multi-key support** — generate keys at custom paths, auto-manage
   `~/.ssh/config` Host entries so Git uses the correct key.
6. **Clone integration** — `gitmap-v28 clone --ssh-key <name>` uses the
   named key's SSH config for the clone operation.

---

## Command: `gitmap-v28 ssh`

Alias: `ssh`

Manages SSH key pairs for Git authentication.

### Subcommands

| Subcommand   | Alias | Description                                    |
|--------------|-------|------------------------------------------------|
| *(default)*  |       | Generate a new SSH key pair                    |
| `cat`        |       | Display the public key for a named key         |
| `list`       | `ls`  | List all stored SSH keys                       |
| `delete`     | `rm`  | Delete a stored key record (optionally files)  |
| `config`     |       | Show or regenerate `~/.ssh/config` entries     |

---

## Subcommand: `gitmap-v28 ssh` (Generate)

### Flags

| Flag             | Short | Type   | Default              | Description                              |
|------------------|-------|--------|----------------------|------------------------------------------|
| `--name`         | `-n`  | string | `default`            | Label for the key in the database        |
| `--path`         | `-p`  | string | `~/.ssh/id_rsa`      | File path for the private key            |
| `--email`        | `-e`  | string | Git global email     | Email comment for the key                |
| `--force`        | `-f`  | bool   | `false`              | Skip regeneration prompt if key exists   |

### Behavior

1. Resolve `--email`: if omitted, read from `git config --global user.email`.
2. Check if a key with `--name` already exists in the database.
   - If yes and `--force` is not set:
     - Print the existing key path and fingerprint.
     - Prompt: `Key "<name>" already exists. [R]egenerate / [N]ew path / [C]ancel`
       - **R**: delete old files, generate at same path.
       - **N**: prompt for a new path, generate there, keep both records.
       - **C**: abort.
   - If yes and `--force` is set: overwrite silently.
3. Run: `ssh-keygen -t rsa -b 4096 -C "<email>" -f "<path>" -N ""`
4. Read the generated `.pub` file content.
5. Insert/update the database record: name, private key path, public key
   content, fingerprint, creation timestamp.
6. Update `~/.ssh/config` (see SSH Config Management below).
7. Print success message with the public key.

### Output

```
  ✓ SSH key "work" generated
    Path:        ~/.ssh/id_rsa_work
    Fingerprint: SHA256:abc123...
    Public key:

  ssh-rsa AAAA... user@example.com

  ℹ  Copy the public key above and add it to your Git provider.
```

---

## Subcommand: `gitmap-v28 ssh cat`

### Flags

| Flag     | Short | Type   | Default   | Description                |
|----------|-------|--------|-----------|----------------------------|
| `--name` | `-n`  | string | `default` | Name of the key to display |

### Behavior

1. Look up the key by `--name` in the database.
2. Print the stored public key content to stdout.
3. If not found, print error and list available key names.

### Output

```
ssh-rsa AAAA... user@example.com
```

---

## Subcommand: `gitmap-v28 ssh list`

Lists all SSH keys stored in the database.

### Output

```
  Name       Path                    Fingerprint            Created
  ─────────  ──────────────────────  ─────────────────────  ──────────────
  default    ~/.ssh/id_rsa           SHA256:abc123...       2026-03-22
  work       ~/.ssh/id_rsa_work      SHA256:def456...       2026-03-22
```

---

## Subcommand: `gitmap-v28 ssh delete`

### Flags

| Flag       | Short | Type   | Default | Description                          |
|------------|-------|--------|---------|--------------------------------------|
| `--name`   | `-n`  | string | *(req)* | Name of the key to delete            |
| `--files`  |       | bool   | `false` | Also delete the key files from disk  |

### Behavior

1. Look up by name, confirm deletion.
2. Remove database record.
3. If `--files`: delete private key and `.pub` file.
4. Remove corresponding `~/.ssh/config` Host entry.

---

## Subcommand: `gitmap-v28 ssh config`

### Behavior

1. Regenerate all `~/.ssh/config` Host entries from database records.
2. Print the current managed config block.

Useful after manual edits or when entries get out of sync.

---

## SSH Config Management

When multiple keys exist, gitmap-v28 auto-manages a clearly marked block
in `~/.ssh/config`:

```
# --- gitmap-v28 managed (do not edit) ---
Host github.com-default
    HostName github.com
    User git
    IdentityFile ~/.ssh/id_rsa
    IdentitiesOnly yes

Host github.com-work
    HostName github.com
    User git
    IdentityFile ~/.ssh/id_rsa_work
    IdentitiesOnly yes
# --- end gitmap-v28 managed ---
```

### Rules

- Entries are bounded by start/end markers.
- Only the managed block is modified; user entries are preserved.
- When only one key exists (named `default`), the Host is plain
  `github.com` (no suffix) so standard Git URLs work unchanged.
- When multiple keys exist, Host entries use the `<hostname>-<name>`
  pattern and guidance is printed.

---

## Clone Integration

### Flag

| Flag         | Short | Type   | Description                           |
|--------------|-------|--------|---------------------------------------|
| `--ssh-key`  | `-K`  | string | Name of SSH key to use for cloning    |

### Behavior

1. Look up the key name in the database.
2. Set `GIT_SSH_COMMAND` to `ssh -i <private-key-path> -o IdentitiesOnly=yes`
   for the clone subprocess.
3. Proceed with normal clone logic.

### Example

```bash
gitmap-v28 clone repos.json --ssh-key work
```

---

## Data Model

### Table: `SSHKeys`

| Column        | Type    | Constraints              | Description                       |
|---------------|---------|--------------------------|-----------------------------------|
| `ID`          | TEXT    | PK, UUID                 | Unique identifier                 |
| `Name`        | TEXT    | UNIQUE, NOT NULL         | User-defined label                |
| `PrivatePath` | TEXT    | NOT NULL                 | Absolute path to private key file |
| `PublicKey`    | TEXT    | NOT NULL                 | Full public key content           |
| `Fingerprint` | TEXT    | NOT NULL                 | SHA256 fingerprint                |
| `Email`       | TEXT    |                          | Email comment used in generation  |
| `CreatedAt`   | TEXT    | NOT NULL                 | ISO 8601 timestamp                |

---

## Constants

New file: `constants/constants_ssh.go`

| Constant               | Value / Format                               |
|-------------------------|----------------------------------------------|
| `CmdSSH`               | `"ssh"`                                      |
| `CmdSSHCat`            | `"cat"`                                      |
| `CmdSSHList`           | `"list"` / `"ls"`                            |
| `CmdSSHDelete`         | `"delete"` / `"rm"`                          |
| `CmdSSHConfig`         | `"config"`                                   |
| `FlagSSHName`          | `"--name"` / `"-n"`                          |
| `FlagSSHPath`          | `"--path"` / `"-p"`                          |
| `FlagSSHEmail`         | `"--email"` / `"-e"`                         |
| `FlagSSHForce`         | `"--force"` / `"-f"`                         |
| `FlagSSHFiles`         | `"--files"`                                  |
| `FlagSSHKey`           | `"--ssh-key"` / `"-K"`                       |
| `SSHKeyType`           | `"rsa"`                                      |
| `SSHKeyBits`           | `"4096"`                                     |
| `SSHConfigMarkerStart` | `"# --- gitmap-v28 managed (do not edit) ---"`   |
| `SSHConfigMarkerEnd`   | `"# --- end gitmap-v28 managed ---"`             |
| `DefaultSSHKeyName`    | `"default"`                                  |
| `MsgSSHGenerated`      | `"  ✓ SSH key \"%s\" generated\n"`           |
| `MsgSSHExists`         | `"  Key \"%s\" already exists at %s\n"`      |
| `MsgSSHCatNotFound`    | `"  Key \"%s\" not found. Available: %s\n"`  |
| `MsgSSHDeleted`        | `"  ✓ SSH key \"%s\" deleted\n"`             |
| `ErrSSHKeygen`         | `"Error generating SSH key: %v\n"`           |
| `ErrSSHReadPub`        | `"Error reading public key: %v\n"`           |

---

## Help File

File: `helptext/ssh.md` (≤120 lines)

### Sections

- **Alias**: *(none — `ssh` is already short)*
- **Usage**: `gitmap-v28 ssh [subcommand] [flags]`
- **Subcommands**: table of cat, list, delete, config
- **Flags**: table of --name, --path, --email, --force
- **Prerequisites**: `ssh-keygen` must be available on PATH
- **Examples**: 3 examples covering generate, cat, multi-key clone
- **See Also**: `gitmap-v28 clone`, `gitmap-v28 setup`

---

## Shell Completion

Add SSH-aware completions:

- After `gitmap-v28 ssh`: suggest subcommands (cat, list, delete, config).
- After `gitmap-v28 ssh cat --name`: suggest key names from DB via
  `gitmap-v28 completion --list-ssh-keys`.
- After `gitmap-v28 clone --ssh-key`: suggest key names.
- New completion flag: `--list-ssh-keys` returns one key name per line.

---

## Package Structure

| File                          | Responsibility                          |
|-------------------------------|-----------------------------------------|
| `cmd/ssh.go`                  | Subcommand dispatch and flag parsing    |
| `cmd/sshgen.go`               | Key generation logic                    |
| `cmd/sshcat.go`               | Public key display                      |
| `cmd/sshlist.go`              | List keys                               |
| `cmd/sshdelete.go`            | Delete key records and files            |
| `cmd/sshconfig.go`            | SSH config file management              |
| `model/sshkey.go`             | SSHKey struct                           |
| `store/sshkey.go`             | CRUD operations for SSHKeys table       |
| `constants/constants_ssh.go`  | All SSH-related constants               |
| `helptext/ssh.md`             | Help documentation                      |

---

## Acceptance Criteria

1. `gitmap-v28 ssh` generates an RSA-4096 key with no passphrase and stores
   the public key in the database.
2. `gitmap-v28 ssh --name work --path ~/.ssh/id_rsa_work` creates a named
   key at a custom path.
3. `gitmap-v28 ssh cat` prints the default key's public key to stdout.
4. `gitmap-v28 ssh cat --name work` prints the named key's public key.
5. Running `gitmap-v28 ssh` when a key exists prompts for
   Regenerate/New path/Cancel (unless `--force`).
6. `gitmap-v28 ssh list` shows all stored keys with paths and fingerprints.
7. `gitmap-v28 ssh delete --name work` removes the DB record.
8. `gitmap-v28 ssh delete --name work --files` also removes key files.
9. `~/.ssh/config` is auto-updated with managed Host entries when
   keys are added or deleted.
10. `gitmap-v28 clone repos.json --ssh-key work` clones using the named
    key via `GIT_SSH_COMMAND`.
11. Shell completion suggests subcommands after `ssh` and key names
    after `--name` and `--ssh-key`.
12. All string literals live in `constants/constants_ssh.go`.

---

## Contributors

- [**Md. Alim Ul Karim**](https://www.linkedin.com/in/alimkarim) — Creator & Lead Architect.
  - [Google Profile](https://www.google.com/search?q=Alim+Ul+Karim)
- [Riseup Asia LLC](https://riseup-asia.com) (2026)
