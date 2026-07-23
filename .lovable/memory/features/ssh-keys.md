# Memory: features/ssh-keys
Updated: 2026-03-22

The 'ssh' command manages SSH key pairs for Git authentication. It supports named keys (--name), RSA-4096 generation with no passphrase via ssh-keygen, public key storage in the SSHKeys database table, and auto-management of ~/.ssh/config Host entries using gitmap-owned markers. The 'ssh cat' subcommand prints the public key for quick copying. Clone integration uses '--ssh-key <name>' flag on the existing clone command, setting GIT_SSH_COMMAND for the subprocess. Shell completion suggests SSH subcommands and key names. Spec: spec/01-app/50-ssh-keys.md.
