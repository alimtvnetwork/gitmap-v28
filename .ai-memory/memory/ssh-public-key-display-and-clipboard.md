# Memory: Bare `gitmap ssh` Public Key Display & Clipboard Copying

## Principle & Strict Avoidance
When `gitmap ssh` is executed without arguments, it functions as an instant public key retrieval and generation command.
It MUST:
1. Print the full public key string to stdout (`  ssh-...`).
2. Copy the public key to the OS clipboard (`[clip] Public key copied to clipboard...`).
3. Display the available SSH subcommands list.

## Absolute Prohibitions
- NEVER remove or suppress public key display when running bare `gitmap ssh`.
- NEVER remove the automatic clipboard copying action.
- NEVER replace bare `gitmap ssh` with only a help message or table that omits the key.
- NEVER truncate or summarize the public key.

## Call Path
- `gitmap ssh` -> `runSSH(args)` in `cli/cmdssh/ssh.go`.
- If `len(args) == 0`:
  - `runSSHGenerate(args)` in `cli/cmdssh/sshgen.go`.
  - Checks if key exists on disk via `keyExistsOnDisk(keyPath)`.
  - Calls `printExistingKeyOnDisk(db, name, keyPath, host)` in `cli/cmdssh/sshexisting.go`.
  - Prints public key string via `printExistingKeyPublic`.
  - Copies to clipboard via `copyPubKeyAndAnnounce(pub)` in `cli/cmdssh/sshcopy.go`.
  - Displays `constants.MsgSSHAvailableCommands`.
