# ssh-bind

Pin a specific SSH private key to the **current repo** by writing
`git config core.sshCommand`. Fixes the wrong-account push failure
without touching `~/.ssh/config` or any global setting.

## Usage

    gitmap ssh-bind <key-filename-or-path>

Alias: `sb`

The argument may be:

- a bare filename (`id_ed25519_aukgit`) — resolved under `~/.ssh/`,
- a `~/`-relative path (`~/.ssh/id_rsa_work`),
- or an absolute/relative filesystem path.

## What it writes

    core.sshCommand = ssh -i <key> -F /dev/null -o IdentitiesOnly=yes

`IdentitiesOnly=yes` stops the ssh-agent from offering unrelated
keys first, which is the actual cause of the "Permission denied to
<other-user>" error even when the correct key is on disk.

## Prerequisites

- Run **inside the target repo** (needs `.git/`).
- The key file must exist. Run `gitmap whoami` to list candidates
  under `~/.ssh/`.

## Examples

```
$ gitmap ssh-bind id_ed25519_aukgit
✓ pinned SSH key for this repo: /home/me/.ssh/id_ed25519_aukgit
  core.sshCommand = ssh -i ~/.ssh/id_ed25519_aukgit -F /dev/null -o IdentitiesOnly=yes
  test with: git push

$ gitmap sb ~/.ssh/id_rsa_work

$ gitmap sb /etc/ssh/deploy_key_aukgit
```

## See Also

- `gitmap fix-auth` / `fa` — generate a per-account key AND bind it in one step.
- `gitmap whoami` / `who` — diagnose which key/account is currently in use.
- `gitmap ssh` — generate and manage named keys.
