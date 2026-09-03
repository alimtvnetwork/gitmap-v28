# Data Portability, Profiles & Bookmarks

Isolate configurations, bookmark command flows, and relocate repository folders safely.

## Commands

### `gitmap profile <subcommand>`
* **Alias:** `pf`
* Subcommands:
  * `create <name>`: Creates a new isolated database profile.
  * `list` (aliases: `ls`, `status`): Lists available profiles.
  * `switch <name>`: Switches the active database profile.
  * `delete <name>`: Deletes a database profile.
  * `show`: Displays details of the current active profile.

### `gitmap bookmark <subcommand>`
* **Alias:** `bk`
* Subcommands:
  * `save <name> [command...]`: Saves a command and flag configuration.
  * `list`: Lists all saved bookmarks.
  * `run <name>`: Executes a saved bookmark.
  * `delete <name>`: Deletes a saved bookmark.

### Portability & Maintenance
* `gitmap export [file]` (alias: `ex`): Exports all database tables to portable JSON.
* `gitmap import <file>` (alias: `im`): Restores database from a JSON export.
* `gitmap mv <src> <dst>` (alias: `move`): Moves a repository folder and updates VS Code and Desktop links.
* `gitmap rm <name>` (aliases: `remove`, `del`): Untracks repository from GitMap database.
