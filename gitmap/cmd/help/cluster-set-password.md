# gitmap cluster set-password

Store a client node password for privileged lifecycle commands (restart, shutdown, logoff).

## Usage

```
gitmap cluster set-password --id <id>
```

## Description

Prompts for a password securely (without echoing to the terminal). The password is encrypted as a bcrypt hash (cost ≥ 12) and stored in the local SQLite `ClusterNode` table. It is used transparently when dispatching lifecycle commands to that node.

**Security Notes:**
- Passwords are only stored as bcrypt hashes.
- The raw text is never written to disk or exported.
- Lifecycle commands using this password require TLS communication.

## Examples

```
gitmap cluster set-password --id 5
```
