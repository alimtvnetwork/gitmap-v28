# Learned: Gitmap Installer Management System

## Architectural Patterns

1. **Forward-Slash Portability**:
   - All relative file paths across Windows and Linux must use forward slashes (`/`).
   - Relative path resolution must anchor to the Gitmap workspace root (`fsutil.MakeRelativeToRoot`).

2. **SQLite Persistence Isolation**:
   - Generic `gitmap reset` does not purge specialized plugin configurations (like installers) unless explicitly scoped (`gitmap reset installer`).
   - SQLite migration triggers and tables are isolated into `installer_scripts` and `installer_versions`.

3. **Cobra Subcommand Argument Partitioning**:
   - Use `DisableFlagParsing: true` on Cobra wrappers when passing through arguments to custom `flag.FlagSet` parsers.
   - Separate non-flag arguments before invoking `fs.Parse()` to ensure flags after positional arguments are parsed without truncation.
