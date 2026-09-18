# Macro Extension Installation Spec

## Overview

"Macro extension" is a GitHub-hosted external tool, fundamentally distinct from the terminal-based interactive \gitmap macro record\ feature.

## Installation Behavior

- The macro extension installs directly into the directory where the installation script is executed.
- It cannot be checked globally as "installed". Every time it is requested, it reinstalls into the target directory.
- This command is distinct from \gitmap macro\ (which handles in-terminal macros). The extension command should be isolated, e.g., \gitmap macro-ext install\.
