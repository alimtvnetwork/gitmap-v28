# Subtask 77.4: `gitmap ssh create` Overwrite Guard, Backup Undo/Redo & Scripts Fixer Sync

- **Parent Plan:** [77-gitmap-db-reset-ssh-key-management-firewall-and-project-detection-rca.md](../../pending/77-gitmap-db-reset-ssh-key-management-firewall-and-project-detection-rca.md)
- **Spec Reference:** [02-spec/21-app/207-gitmap-db-reset-ssh-key-management-firewall-and-project-detection-rca/01-architecture-spec.md](../../../../02-spec/21-app/207-gitmap-db-reset-ssh-key-management-firewall-and-project-detection-rca/01-architecture-spec.md)
- **Status:** Pending
- **Target Area:** `cli/cmdssh/sshgen.go`, `cli/cmdssh/ssh_undo_cmd.go`, `cli/system/`, `scripts-fixer`

## Objective

Safeguard SSH key generation in `gitmap ssh create` by detecting pre-existing keys, prompting the user for confirmation (`[y/N]`) unless explicitly confirmed (`-y`/`--confirm`), automatically creating timestamped backups of existing keys, registering journal undo/redo actions, and keeping automated key regeneration in `scripts-fixer` non-blocking.

## Requirements & Implementation Details

1. **Pre-Existing Key Detection & Interactive Prompt:**
   - In `cli/cmdssh/sshgen.go` (`runSSHCreate`), inspect `keyExistsOnDisk(keyPath)`.
   - If the target private or public key exists and neither `-y`, `--yes`, `--confirm`, nor `-f`/`--force` is supplied:
     - Render warning: `⚠ Warning: SSH key already exists at <keyPath>`
     - Prompt: `Overwrite and backup existing key? [y/N]: `
     - If the user enters anything other than `y` or `yes`, abort cleanly with: `SSH key creation canceled.`
2. **Timestamped Key Backup:**
   - When confirmed, copy the existing private and public key files to `<keyPath>.bak.<YYYYMMDDHHmmss>`.
   - Announce the backup location: `✔ Backed up existing SSH key to: <backupPath>`.
3. **Undo/Redo Journaling Integration:**
   - Call `RecordSSHKeyBackupTask(ctx, name, keyPath, backupPath)` to store the backup metadata in the GitMap task journal.
   - Ensure `gitmap ssh undo` can restore the backed-up key pair if the user reverts the operation.
4. **Key Generation & Storage Reseed:**
   - Proceed with generating the new RSA/ED25519 key pair via `ssh-keygen`.
   - Store the newly generated public key and fingerprint in the SQLite `SSHKey` table (`upsertExistingKeyToDB`).
   - Print the generated key details card.
5. **Automation & `scripts-fixer` Synchronization:**
   - Ensure automated scripts and `scripts-fixer` pass `-y` / `--confirm` during non-interactive fleet provisioning to avoid hanging on standard input.

## Verification & Acceptance Criteria

- Running `gitmap ssh create` against an existing key without flags halts and prompts `[y/N]`. Entering `n` preserves the original key untouched.
- Entering `y` creates a backup file matching the format `<keyPath>.bak.<timestamp>` before overwriting.
- Running `gitmap ssh create -y` executes immediately without prompting, creates a backup, and generates a new key pair.
- The backup operation is recorded in the task journal for undo/redo capability.
