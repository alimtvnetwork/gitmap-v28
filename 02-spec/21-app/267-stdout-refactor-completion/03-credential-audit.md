# Program 267 / WS5 — Credential Storage & Handling Audit (READ-ONLY)

> Audit scope: every place gitmap **stores** secrets, every path that can
> **expose** them (logs, argv, URLs), and how they are **transported**.
> Read-only — no code changed, no behavior modified.
> Method: code reading via `gitmap aum search`; files referenced by relative
> path from repo root. Audit date: 2026-10-10.

## 1. Secret stores

### S1. GitHub token from `gitmap login` / `gitmap token add`

- **Writer:** `cli/cmdlogin/login_token.go:storeGitHubToken` — runs
  `git config --global github.token <token>`; also
  `cli/cmdtoken/token_cmd.go:runTokenAdd` (same mechanism).
- **Location:** user's global git config — `~/.gitconfig` (or
  `$XDG_CONFIG_HOME/git/config` when XDG set). Key `github.token`.
- **Format:** plaintext, one INI-style line under `[github]`.
- **Permissions:** git writes the file with `0666 & ~umask` → **0644** on a
  standard umask (readable by every local user). gitmap sets no explicit mode
  (no `os.Chmod` anywhere in `cli/` — verified via `aum search`).
- **Keychain / OS store:** NOT used. Flat file only.
- **Resolution precedence** (`cli/secrets/ghtoken.go:Resolve`):
  1. `GITHUB_TOKEN` / `GH_TOKEN` env vars (plaintext, inherited env)
  2. Windows registry user/machine `Environment` values (Windows only,
     `cli/secrets/ghtoken_windows.go` — machine-level key is system-wide)
  3. `git credential fill` for github.com (Git Credential Manager — the
     **only** OS-backed store gitmap reads; gitmap never writes there)
  4. `gh auth token` / `~/.config/gh/hosts.yml` (plaintext YAML `oauth_token:`,
     owned by gh, 0600 by gh itself)
  5. `git config github.token` (the gitmap-written one)
- **Removal:** `gitmap logout` (`cli/cmdlogin/login_status.go:removeStoredTokens`)
  unsets `git config --global github.token` and issues a
  `git credential reject` for github.com.

### S2. SSH fleet vault — node passwords

- **Location:** split-DB table `SSHConnection` (`EncryptedPassword` column),
  in the **installation** split DB (`cli/db/sshconnection.go`); legacy
  `ssh_hosts` table is merged in at read time. DB files are created by the
  SQLite driver with **no explicit chmod** → umask default (typically 0644).
- **Encryption formats** (all stored as prefixed strings):
  - `rsa:` — RSA-OAEP-SHA256 ciphertext. Key: dedicated vault keypair
    `~/.gitmap/keys/vault_rsa` (0600) / `vault_rsa.pub` (0644),
    auto-generated 2048-bit on first use
    (`cli/cmdssh/ssh_vault_rsa.go`); fallback to the user's own
    `~/.ssh/id_rsa` if the vault key is absent
    (`cli/cmdssh/ssh_crypto_pass.go`).
  - `aes:` — AES-GCM with **hardcoded-in-source** keys
    (`gitmap-ssh-secret-key-0123456789` / `gitmap-ssh-fallback-secret-01234`,
    `cli/secrets/ssh_password.go` + `cli/cmdssh/ssh_crypto_pass.go`).
    Anyone who can read the repo or the binary can decrypt these.
  - `salt:` / `caesar:` — rotation ciphers (base64 + character shift), **not
    real encryption**; reversible by anyone. Kept for legacy compatibility.
  - Unprefixed legacy strings pass through as-is (**plaintext**).
- **Verdict:** the `rsa:` path is genuinely encrypted (key material 0600).
  Everything else is obfuscation — `aes:` with a published key and the
  salt/caesar formats are decryptable by any local process.

### S3. Supabase project keys (`gitmap supabase`)

- **Location:** installation split-DB `supabase` table — `AnonKeyEnc`,
  `ServiceKeyEnc`, `DbUrlEnc` (`cli/cmdsupabase/supabase_db.go`).
- **Encryption:** AES-256-GCM via `EncryptSecret`
  (`cli/cmdsupabase/supabase_vault.go:77`) with a master key derived as
  `HMAC-SHA256(salt, hostname)`; salt = 32 random bytes in a `vault.salt`
  file (0600, auto-created) sitting two dirs above the installation DB.
- **Format:** raw base64 GCM ciphertext (no prefix), or legacy `rsa:`/`aes:`
  handled by the SSH vault fallback on decrypt.
- **Display:** `supabase list` masks with `MaskSecret` (first 6 chars + `...****`);
  note the first 4–6 chars of the secret are still shown.
- **Weakness:** key derivation uses only the hostname as the fingerprint —
  a hostname is low-entropy, but the random 32-byte `vault.salt` (0600) is
  the real secret, so this is reasonable defense against cross-machine
  offline DB theft only.

### S4. In-process secrets (transient, not on disk)

- **Session token cache:** `cli/cmdclone/clone_auth_store.go` —
  `globalToken` package var (plaintext, mutex-guarded, process lifetime).
  `--token <PAT>` is stored here and **never written to disk** (per comment in
  `clone.go:prepareCloneTransport`), but it lives in argv/ps while the
  process runs (see §3, T2).
- **AskPass bridge:** `cli/cmdssh/ssh_askpass.go` — password passed as env var
  `GITMAP_SSH_PASS=<plaintext>` to a temp script (script itself is 0700,
  contains no secret, and is deleted after use). Env vars are readable via
  `/proc/<pid>/environ` by the same user while the child runs (transient).
- **SSH connections:** `ssh.InsecureIgnoreHostKey()` in
  `cli/secrets/ssh_client.go` — MITM host verification is disabled, so fleet
  password/token transport is only as safe as the network path.

## 2. Transport (how secrets move)

- **GitHub REST API:** all callers (`cli/clonenext/github.go`,
  `cli/cmdlogin/login_token.go`, `cli/release/*`, `cli/cmdpurge/release_prune.go`,
  `cli/cmddoctor/doctor_extra.go`, `cli/cmd/workflow_open_pr.go`) send the
  token as an **HTTP header** (`Authorization: Bearer|token <tok>`) over HTTPS.
  Clean — no token in URLs at the API layer.
- **Supabase:** headers `apikey` + `Authorization: Bearer` over HTTPS
  (`cli/cmdsupabase/supabase_cmd.go:260`).
- **`git clone` transport — token injected into the URL** (the big one):
  `InjectTokenIntoHTTPS` (`cli/cmdclone/clone_auth_inject.go:11`) rewrites the
  remote to `https://<token>@github.com/org/repo.git` in
  `cli/cmdclone/clone_auth_resolve.go` (from session cache, resolved config,
  or terminal prompt). Consequences:
  1. **`.git/config` of every token-cloned repo persists the token** in
     `url = https://<token>@...` (0644 plaintext) — git does not strip it.
     `git remote -v` then displays it to anyone in the terminal.
  2. Token is on the **process command line** while `git clone` runs →
     visible in `ps aux` to all local users.
  3. Token appears in **CI/CD logs** — git echoes the URL on failure
     (`fatal: unable to access 'https://<token>@github.com/...'`).
- **`gitmap token deploy` to fleet nodes** (`cli/cmdssh/ssh_token_deploy.go:61`):
  the full GitHub token is embedded in a remote shell string
  `git config --global github.token <token>` sent over SSH — plaintext in
  transit on the remote's `ps` (remote `sh -c "git config --global
  github.token <token>"`), and stored **plaintext in remote `~/.gitconfig`**
  (remote umask, 0644) on every fleet node. Reporting prints only the remote
  output, not the command — small mercy.

## 3. Logging / display exposure

- **L1 — Terminal RepoTermBlock prints the token-embedded URL verbatim.**
  `render.RenderRepoTermBlock` prints `to:` (= TargetURL) and `command:` (=
  full `git clone https://<token>@... <dest>` line) with **no masking**
  (`cli/render/repotermblock.go:buildBlockBody`,
  `cli/cmdclone/clonetermstream.go:buildCloneCommand`). Any `--output terminal`
  clone with token-injected auth leaks the full token to stdout (screenshots,
  scrollback, CI logs, piped consumers).
- **L2 — `--print-clone-argv` argv dump** (`cli/cmdclone/cloneprintargv.go`):
  prints the exact argv including the token-embedded URL to stderr when the
  audit flag is set. Opt-in, but it is literally a secret dump flag.
- **L3 — `gitmap login --token <PAT>` and `gitmap token add <token>` take the
  token as a CLI argument** → visible in `ps` (and shell history) for the
  whole process lifetime. (`cli/cmdclone/flags.go:87`, `parseLoginArgs`.)
- **L4 — Interactive terminal token prompt** (`cli/cmdclone/clone_auth_token_input.go`):
  `bufio.NewReader(os.Stdin)` with **no echo suppression** — the pasted/typed
  PAT is visible on screen as typed. (Contrast: `login_browser.go:67`
  uses a secure read for `gitmap login`'s own prompt.)
- **L5 — Masked displays are present but partial:** `gitmap login --status`
  and `gitmap token list` show `first4****last4` of the token. The 8 known
  characters materially reduce brute-force search space for a short token,
  but more importantly they confirm validity to anyone reading the screen.
- **L6 — Env var resolution:** `GITHUB_TOKEN`/`GH_TOKEN` are read from the
  environment, so any child process, crash dump, or `pe`-style diagnostic
  that dumps the environment inherits them. No evidence of gitmap itself
  logging env.
- **What does NOT leak:** API calls use headers; error paths checked
  (`clonenexterrorreport.go`, `reportRemoteExecution`, `cmd_whoami`,
  `doctor_extra.go`) do not print tokens; no external telemetry sink carries
  secrets (machine telemetry is fleet-local over SSH).

## 4. Severity ratings

| ID | Finding | Severity | One-line rationale |
|----|---------|----------|--------------------|
| F1 | `github.token` stored **plaintext in `~/.gitconfig` (0644)** by `gitmap login` / `gitmap token add` | **Critical** | Any local user/process reads the user's GitHub PAT with all granted scopes |
| F2 | RepoTermBlock prints **full token in URL** (`to:` + `command:` lines) on stdout | **Critical** | Every terminal-output clone leaks the PAT to scrollback, CI logs, pipes |
| F3 | Token persists in **cloned repo `.git/config`** (0644) and remote `-v` output | **High** | PAT leaks to every machine/user the repo directory is copied to or backed up on |
| F4 | `--token <PAT>` / `token add <token>` on **command line** | **High** | PAT visible in `ps` + shell history for process lifetime |
| F5 | `token deploy` writes **plaintext PAT to remote `~/.gitconfig`** on fleet nodes and puts it in remote `ps` | **High** | Multiplies F1 across the fleet; remote history/`ps` window per node |
| F6 | Vault `aes:` keys and `salt:`/`caesar:` formats are **hardcoded/reversible** | **High** | Anyone reading source/binary decrypts "encrypted" SSH passwords from the DB |
| F7 | Interactive clone token prompt **echoes PAT to terminal** | **Medium** | Shoulder-surfing / screen-share leak; login's own prompt does this right |
| F8 | `--print-clone-argv` dumps argv with token URL to stderr | **Medium** | Opt-in debug flag, but it's a designed secret-print path |
| F9 | Masked displays show `first4****last4` | **Low** | Partial confirmation of token validity on screen |
| F10 | `GITMAP_SSH_PASS` env var bridge for askpass | **Low** | Same-user `/proc` visibility, transient, deleted script |
| F11 | Vault RSA key 0600 + vault.salt 0600; Supabase AES-256-GCM with 0600 salt | **Low** (positive) | Key-material hygiene is actually good; weak link is key *derivation* entropy (hostname) |

## 5. Recommended follow-ups (fixes OUT of scope — listed only)

1. Stop storing PAT in `git config github.token`; use the OS credential store
   (Git Credential Manager via `git credential approve`, already read by
   `secrets.Resolve`) or `gh`'s hosts.yml; migrate existing plaintext entries.
2. Replace URL-embedded tokens for clone/pull/push with a `credential.helper`
   or `GIT_ASKPASS`-style flow; redact `user:pass@` from all rendered URLs
   (`render.RepoTermBlock`, `buildCloneCommand`, `runCmdPrintArgv`) as an
   immediate stopgap.
3. Stop writing `github.token` to remote nodes' `~/.gitconfig` in
   `token deploy`; prefer per-node credential helper or short-lived tokens.
4. Drop `aes:`-with-hardcoded-key and `salt:`/`caesar:` write paths; keep
   RSA vault as the only SSH-password format (migration: re-encrypt on next
   read).
5. Suppress echo in `readTokenFromTerminal` (reuse the secure reader from
   `login_browser.go`).
6. Set restrictive perms on created DB/key files (`0600` for installation DB,
   already done for `vault_rsa`/`vault.salt`); document that `~/.gitconfig`
   entry is plaintext and world-readable.
7. Consider redacting `first4/last4` masking to a fixed-length mask in
   `login --status` / `token list` (leak-confirmation hardening).

## Source files consulted (all relative to repo root)

- `cli/cmdlogin/login_token.go`, `login_cmd.go`, `login_status.go`, `login_browser.go`
- `cli/cmdtoken/token_cmd.go`
- `cli/secrets/ghtoken.go`, `ghtoken_windows.go`, `ghtoken_other.go`, `secrets_resolver.go`, `encrypt.go`, `salted_cipher.go`, `ssh_password.go`, `ssh_client.go`
- `cli/cmdssh/ssh_token_deploy.go`, `ssh_crypto_pass.go`, `ssh_vault_rsa.go`, `ssh_askpass.go`, `ssh_install_actions.go`, `fleet_parallel.go`
- `cli/cmdclone/clone.go`, `flags.go`, `clone_auth_inject.go`, `clone_auth_resolve.go`, `clone_auth_store.go`, `clone_auth_terminal.go`, `clone_auth_prompt.go`, `clone_auth_token_input.go`, `clone_auth_probe.go`, `clone_auth_browser.go`, `clonetermstream.go`, `clonetermrow.go`, `clonetermurl.go`, `clonetermplan.go`, `cloneprintargv.go`, `clonenexterrorreport.go`, `clone_json.go`, `directclone_record.go`
- `cli/cmdsupabase/supabase_vault.go`, `supabase_cmd.go`, `supabase_db.go`, `supabase_types.go`
- `cli/render/repotermblock.go`, `pretty_emit.go`
- `cli/store/split_db_path.go`, `cli/db/sshconnection.go`, `cli/db/clusternode.go`
- `cli/cmd/cmd_whoami.go`, `cli/cmddoctor/doctor_extra.go`, `cli/cmdfixauth/fixcredential.go`, `cli/clonenext/github.go`, `cli/release/githubapi.go`, `cli/release/assetsupload.go`, `cli/cmdpurge/release_prune.go`
