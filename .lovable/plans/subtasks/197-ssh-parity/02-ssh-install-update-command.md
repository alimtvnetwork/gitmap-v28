# Subtask 02: SSH Remote GitMap Install and Update Commands

## Status: COMPLETED

## Summary of Accomplishments
1. Implemented `gitmap ssh install [target]` and `gitmap ssh install gitmap [target]` in `cli/cmdssh/ssh_install_remote.go` and `cli/cmdssh/ssh_install_actions.go`.
2. Added target connection loader in `cli/cmdssh/ssh_target_loader.go` to support single aliases/IPs or `all`.
3. Implemented remote GitMap detection (`isGitmapPresent` running `gitmap --version` over SSH):
   - If missing: executes official platform installation one-liner (`curl` bash or PowerShell).
   - If present: triggers `gitmap update` to synchronize to the latest release.
4. Implemented `gitmap ssh update [target]` and `gitmap ssh update gitmap [target]` in `cli/cmdssh/ssh_update_remote.go`.
5. Routed `install` and `update` subcommands through `dispatchPrimarySSH` in `cli/cmdssh/ssh.go`.
6. Enforced canonical <= 100 line limit across all install/update files.
