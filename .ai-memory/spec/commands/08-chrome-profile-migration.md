# 08 - Chrome Profile Migration (Cross-OS)

## Overview

The user requested an enhancement to the existing Chrome Profile export/import commands (gitmap cpe, gitmap cpi). The current implementation only exports Bookmarks, Preferences, and extension IDs to a JSON or CSV file.

The enhancement adds support for bundling Chrome's raw SQLite databases (e.g., History, Web Data) alongside the JSON snapshot into a single deployable .zip archive, or allowing selective export of just the JSON or just the SQLite files. The goal is to easily move a profile between machines of different operating systems (Windows to Linux, etc.).

## Key Requirements (from user prompt)

1. **Cross-OS Compatibility**: Must support migrating profiles regardless of the OS.
2. **Export Formats**: zip, json, sqlite.
3. **Integration**: Enhance existing gitmap CLI commands.
4. **Integrated Tests**: Must write comprehensive tests.
