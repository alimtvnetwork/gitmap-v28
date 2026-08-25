# Coding Guidelines (CG) Command & Installation Spec

## Overview
The \cg\ (Coding Guidelines) commands provide automated installation and updates for the project's coding guidelines.

## Installation Scripts (Non-Negotiable)
When running \cg install\ or \cg update\, the CLI MUST execute the following exact commands based on the OS:

**Windows (PowerShell):**
\\\powershell
irm https://raw.githubusercontent.com/alimtvnetwork/coding-guidelines-v24/main/install.ps1 | iex
\\\

**Unix/Linux/macOS:**
\\\ash
curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/coding-guidelines-v24/main/install.sh | bash
\\\

## Command Routing
- \gitmap cg\ or \gitmap cg help\: Displays all available subcommands, help text, and explicit \--except\ examples.
- \gitmap cg install [repo-name]\: Installs CG into the given repo. If omitted, installs in the current Git repo.
- \gitmap cg install --workspace --except <repo1,repo2>\: Installs CG on all repos in the current workspace, skipping the exceptions.
- \gitmap cg update\: Updates existing guidelines using the same scripts.
