# fix-auth

End-to-end fix for the classic SSH push failure:

    ERROR: Permission to owner/repo.git denied to <wrong-user>.
    fatal: Could not read from remote repository.

Happens when the system SSH agent silently offers the default private
key (e.g. `id_ed25519`) that belongs to a different GitHub account
than the one with write access to the repo. `fix-auth` generates a
per-account key, pins it to the **current repo** via
`core.sshCommand`, and copies the public key to your clipboard so you
can paste it into GitHub.

## Usage

    gitmap fix-auth --user <github-username> [--email <addr>] [-y] [-f]

Alias: `fa`

## Flags

| Flag        | Short | Description                                        | Default              |
|-------------|-------|----------------------------------------------------|----------------------|
| --user      | -u    | GitHub username (required)                         | —                    |
| --email     | -e    | Email comment on the SSH key                       | git config user.email |
| --yes       | -y    | Skip overwrite confirmation                        | false                |
| --force     | -f    | Regenerate even if the key already exists          | false                |

The first positional token is treated as `--user` when the flag is
omitted (`gitmap fix-auth aukgit`).

## What it does

1. Ensures `~/.ssh` exists with mode `0700`.
2. Generates `~/.ssh/id_ed25519_<user>` via `ssh-keygen -t ed25519`
   (skipped if the key already exists and `--force` is not passed).
3. Runs `git config core.sshCommand "ssh -i ~/.ssh/id_ed25519_<user>
   -F /dev/null -o IdentitiesOnly=yes"` in the current repo. The
   `IdentitiesOnly=yes` bit is the fix — it stops the ssh-agent from
   offering the wrong-account default key first.
4. Prints the public key, copies it to the OS clipboard (`clip` on
   Windows, `pbcopy` on macOS, `wl-copy`/`xclip`/`xsel` on Linux),
   and shows the `https://github.com/settings/ssh/new` link.

## Prerequisites

- Must be run **inside the repo** you want to fix (needs `.git/`).
- `ssh-keygen` on PATH (ships with OpenSSH).
- `git` on PATH.

## Examples

### Full fix for the wrong-account push failure

```
$ cd path/to/broken-repo
$ gitmap fix-auth --user aukgit --email me@example.com
Generating public/private ed25519 key pair.
✓ pinned repo → ssh -i ~/.ssh/id_ed25519_aukgit -F /dev/null -o IdentitiesOnly=yes

=== PUBLIC KEY (add this to GitHub) ===
ssh-ed25519 AAAA... me@example.com
=======================================
  ✓ Public key copied to clipboard (95 bytes)

Next steps for aukgit:
  1. Open https://github.com/settings/ssh/new
  2. Paste the key above, save.
  3. Run:  git push
```

### Reuse an existing per-account key (just re-pin the repo)

```
$ gitmap fix-auth -u aukgit
• key already exists, reusing: /home/me/.ssh/id_ed25519_aukgit
✓ pinned repo → ssh -i ~/.ssh/id_ed25519_aukgit ...
```

### Force-regenerate (rotate the key)

```
$ gitmap fix-auth --user aukgit --force --yes
```

### Positional form

```
$ gitmap fa aukgit
```

## See Also

- `gitmap whoami` / `who` — diagnose which identity Git & SSH will use.
- `gitmap ssh-bind <key>` / `sb` — pin an existing key to a repo without generating a new one.
- `gitmap ssh` — general SSH key generation and management.
