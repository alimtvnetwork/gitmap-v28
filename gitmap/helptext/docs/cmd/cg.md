# gitmap cg (Coding Guidelines, Prompts & Versioning)

Manage and install Coding Guidelines, Prompt Architect, and Version JSON in any repository.

## Usage
```bash
# Canonical version.json (SSOT) Installation
gitmap cg install-version-json [targets...] [--version=<ver>]  # Install version.json + docs + enqueue what-to-read
gitmap init-version [--version=<ver>]                         # Shortcut for current repository

# Coding Guidelines
gitmap cg status                                              # Tabulated status of Coding Guidelines
gitmap cg version                                             # Show installed Coding Guidelines version
gitmap cg install [targets...]                                # Install Coding Guidelines
gitmap cg update [targets...]                                 # Selective update of installed repos

# Prompt Architect v2
gitmap cg install-prompts [targets...]                        # Install Prompt Architect v2
gitmap cg update-prompts [targets...]                         # In-place update Prompt Architect v2
gitmap cg prompts-status [targets...]                         # Check Prompt Architect status
gitmap cg prompts-version [targets...]                        # Show Prompt Architect version

# Aliases
gitmap install-prompts                                        # Current directory shortcut
gitmap ct install-prompts                                     # Backward-compatible alias
```
