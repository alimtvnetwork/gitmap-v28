# whoami

Diagnose *which* Git identity and *which* SSH key the current repo
will use for the next `git push`. Answers the question:

> Why does GitHub say "Permission denied to <wrong-user>" when my
> `git config user.email` is correct?

Because SSH auth ignores `user.email` entirely — that field only
stamps commit authorship. Auth is decided by whichever private key
the SSH agent offers first. `whoami` prints both so the mismatch
becomes obvious, and lists the SSH keys under `~/.ssh/` you can pin
with `gitmap ssh-bind` or `gitmap fix-auth`.

## Usage

    gitmap whoami

Alias: `who`

## Output

- Local vs global `user.name` / `user.email`.
- Origin URL + detected transport (HTTPS / SSH).
- **HTTPS:** the cached credential username (via `git credential fill`).
- **SSH:** the GitHub account seen by `ssh -T git@github.com` and the
  private-key filename the agent offered.
- List of private keys in `~/.ssh/` (excluding `.pub`, `known_hosts`,
  `config`) — candidates to pass to `gitmap ssh-bind` / `fix-auth`.
- Copy-paste fix hints for the detected transport.

## Examples

```
$ gitmap whoami
identity (local):   Ali Karim <me@example.com>
identity (global):  Ali Karim <me@example.com>
origin:             git@github.com:aukgit/alim.karim.profile.git  (ssh)
ssh principal:      karim-mum-v1        ← WRONG account
ssh key offered:    ~/.ssh/id_ed25519

available keys in ~/.ssh:
  id_ed25519
  id_ed25519_aukgit
  id_rsa_work

fix:  gitmap fix-auth --user aukgit
  or: gitmap ssh-bind id_ed25519_aukgit
```

## See Also

- `gitmap fix-auth` / `fa` — one-shot generate + pin + clipboard.
- `gitmap ssh-bind <key>` / `sb` — pin an existing key to this repo.
- `gitmap ssh list` — list keys tracked in the gitmap database.
