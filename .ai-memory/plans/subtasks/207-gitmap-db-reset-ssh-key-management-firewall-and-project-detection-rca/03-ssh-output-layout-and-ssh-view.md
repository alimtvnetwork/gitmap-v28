# Subtask 77.3: SSH Output Layout Reorder & `gitmap ssh view` Command

- **Parent Plan:** [77-gitmap-db-reset-ssh-key-management-firewall-and-project-detection-rca.md](../../pending/77-gitmap-db-reset-ssh-key-management-firewall-and-project-detection-rca.md)
- **Spec Reference:** [02-spec/21-app/207-gitmap-db-reset-ssh-key-management-firewall-and-project-detection-rca/01-architecture-spec.md](../../../../02-spec/21-app/207-gitmap-db-reset-ssh-key-management-firewall-and-project-detection-rca/01-architecture-spec.md)
- **Status:** Pending
- **Target Area:** `cli/cmdssh/ssh.go`, `cli/cmdssh/sshcat.go`, `cli/cmdssh/sshgen.go`

## Objective

Reorganize the terminal presentation of `gitmap ssh` so help categories and commands are displayed first, while the active SSH key details card appears at the bottom as a summary. Implement a dedicated `gitmap ssh view` command (`v`, `show`, `key`, `pubkey`) that shows only the public key card and automatically copies the public key to the OS clipboard without rendering command help.

## Requirements & Implementation Details

1. **Terminal Output Layout Reordering:**
   - In `cli/cmdssh/ssh.go`, adjust `handleEmptySSHArgs` to render the CLI help menu and command options first (`RenderSSHHelp()`).
   - Position the SSH key identity status, fingerprint, and public key card at the very bottom as the focal summary card.
2. **Dedicated Key Viewer Command (`gitmap ssh view`):**
   - Register top-level aliases in SSH dispatch table: `view`, `v`, `show`, `key`, `pubkey`.
   - Dispatch to `runSSHCat(args)` in `cli/cmdssh/sshcat.go`.
   - Support optional key name argument: `gitmap ssh view [key-name]`.
   - Provide `--raw` / `-r` flag to print purely the raw public key string (e.g. for scripting, piping to remote authorized_keys, or redirects).
3. **Auto-Clipboard Synchronization:**
   - On execution of `gitmap ssh view`, automatically copy the active public key to the operating system clipboard (using Windows `clip.exe` / PowerShell clipboard, Linux `xclip` / `wl-copy`, macOS `pbcopy`).
   - Print confirmation: `Status: Copied to clipboard! ✔`.
4. **Clean Presentation:**
   - Render formatted ANSI box with Key Label, Path, Algorithm, Fingerprint, and Public Key preview.
   - Display regeneration guidance: `[tip] To regenerate this key: gitmap ssh create [email] [-y]`.
   - Do NOT display the general command help menu when executing `gitmap ssh view`.

## Verification & Acceptance Criteria

- Running `gitmap ssh` displays command options first, followed by the SSH key card at the bottom.
- Running `gitmap ssh view` outputs only the public key card and copies the key to the clipboard.
- Running `gitmap ssh view --raw` outputs only the raw public key line with no box decorations or help text.
