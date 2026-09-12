# explorer

Launch the native graphical desktop file manager regardless of operating system.

## Aliases

`open-explorer`, `folder`, `open-folder`, `browse-folder`

## Usage

    gitmap explorer [path]

## Description

Opens the native file explorer at the target path:
- **Windows**: Launches `explorer.exe <path>` (or `explorer.exe /select,<path>` if target is a file, highlighting it in the folder).
- **macOS**: Launches `open <path>` (or `open -R <path>` if target is a file, revealing it in Finder).
- **Linux**: Launches `xdg-open <path>`.

If no `[path]` is provided, it defaults to the current working directory.

## Examples

```bash
# Open file explorer in the current repository directory
gitmap explorer

# Open a specific folder
gitmap explorer src/components

# Highlight a specific file in the desktop file manager
gitmap explorer .env
```
