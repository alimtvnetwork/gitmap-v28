## Quick Install v6.297.0

### Windows (PowerShell 5.1+)
```powershell
irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.297.0/install.ps1 | iex
```

### Linux / macOS (Bash)
```bash
curl -fsSL https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.297.0/install.sh | bash
```

## Changelog v6.297.0

- Fixed SSH Host Key verification: implemented auto-pruning of stale known_hosts entries (handling hashed entries via ssh-keygen -R and line deletion) to prevent REMOTE HOST IDENTIFICATION HAS CHANGED
- Added auto-recovery & retry in SpawnSSH: seamlessly auto-prunes stale host keys and re-trusts remote machines upon detecting host key changes
- Preserved bare gitmap ssh public key display and clipboard copying: strictly enforced TOTAL BAN in .ai-memory/strictly-avoid.md
- Guaranteed bounded stack traces on all errors: updated global error handler so no error is ever emitted without an informative stack trace
- Added repository creation commands suite: gitmap repo-create (repoc), gitmap create-repo (crepo), gitmap create-local-repo (clr) with automatic space slugification
- Enhanced commit transfer & PR workflows: automated destination repository directory provision and GitHub creation in commit-in, commit-left, commit-right, and PR counterparts (cin-pr, cml-pr, cmr-pr)
- Resolved CI/CD pipeline issues: eliminated constants collision (CmdCompareAlias to "comp"), AST registry discrepancies, unused functions, and staticcheck context warnings across all 35 gates
