# gitmap cluster set-password

Securely configure or update the client node credential for privileged lifecycle operations.

## Usage

```bash
gitmap cluster set-password --id <node-id>
```

## Description

Prompts for a password, computes a bcrypt hash with cost 12, and stores the hash in the node record. Passwords are never stored in plaintext or exported.

## Flags

- `--id <id>`: Target node Display ID.

## Examples

```bash
gitmap cluster set-password --id 2
```

See also: `gitmap clients`, `gitmap cluster`
