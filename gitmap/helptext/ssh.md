# SSH Key Management

Manages SSH key pairs for Git authentication.

## Usage

    gitmap ssh [subcommand] [flags]

## Subcommands

| Subcommand | Alias | Description                               |
|------------|-------|-------------------------------------------|
| *(none)*   |       | Generate a new SSH key pair               |
| cat        |       | Display the public key for a named key    |
| list       | ls    | List all stored SSH keys                  |
| delete     | rm    | Delete a key record (optionally files)    |
| config     |       | Regenerate ~/.ssh/config managed entries  |

## Flags (generate)

| Flag      | Short | Description                            | Default        |
|-----------|-------|----------------------------------------|----------------|
| --name    | -n    | Label for the key in the database      | default        |
| --path    | -p    | File path for the private key          | ~/.ssh/id_rsa  |
| --email   | -e    | Email comment for the key              | git global     |
| --force   | -f    | Skip prompt if key already exists      | false          |
| --host    | -H    | Git provider hostname                  | github.com     |
| --confirm |       | Require explicit yes before generating | false          |

## Flags (list)

| Flag      | Short | Description                            |
|-----------|-------|----------------------------------------|
| --json    |       | Output keys as JSON for scripting      |

## Flags (delete)

| Flag      | Short | Description                            |
|-----------|-------|----------------------------------------|
| --name    | -n    | Name of the key to delete              |
| --files   |       | Also delete key files from disk        |

## Flags (clone integration)

| Flag      | Short | Description                            |
|-----------|-------|----------------------------------------|
| --ssh-key | -K    | SSH key name to use for cloning        |

## Prerequisites

`ssh-keygen` must be available on PATH (included with OpenSSH).

## Examples

### Generate a default SSH key

    $ gitmap ssh
      Generating public/private rsa key pair.
      ✓ SSH key "default" generated
        Path:        ~/.ssh/id_rsa
        Fingerprint: SHA256:abc123...
        Public key:

      ssh-rsa AAAA... user@example.com

      ℹ  Copy the public key above and add it to your Git provider.

### Generate a named key for work

    $ gitmap ssh --name work --path ~/.ssh/id_rsa_work
      ✓ SSH key "work" generated
        Path:        ~/.ssh/id_rsa_work
        Fingerprint: SHA256:def456...

### Display the public key (`view` / `v` / `cat`)

    $ gitmap ssh view --name work
    ssh-rsa AAAA... user@example.com

### Copy the public key to the clipboard (`copy` / `cp`)

    $ gitmap ssh copy --name work
    ssh-rsa AAAA... user@example.com
      ✓ Public key copied to clipboard (412 bytes)

Uses `clip` on Windows, `pbcopy` on macOS, and `wl-copy` / `xclip` / `xsel`
on Linux (first one found). If none are available, the key is still
printed and a one-line warning is emitted — never fails.

### Create a key (`create` is an explicit alias for the default)

    $ gitmap ssh create --name work --path ~/.ssh/id_rsa_work

### Clone using a specific SSH key

    $ gitmap clone repos.json --ssh-key work
      → Cloning with SSH key "work" (~/.ssh/id_rsa_work)
      ✓ Cloned 5 repos, 0 failed.

### List all stored keys

    $ gitmap ssh list

      SSH Keys (2):

      Name            Path                            Fingerprint                Created
      default         ~/.ssh/id_rsa                   SHA256:abc123...           2026-03-22
      work            ~/.ssh/id_rsa_work              SHA256:def456...           2026-03-22

## See Also

- `gitmap clone` - Clone repositories from structured files
- `gitmap setup` - Configure Git global settings

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter ssh
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
