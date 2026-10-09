# open

Opens a file, folder, URL, or email address in the host operating system's default application.
Cross-platform compatible: works gracefully on Windows, macOS, and Linux.

## Usage

`gitmap open <target>`

## Examples

```bash
gitmap open "readme.md"       # Opens in default text editor
gitmap open "."               # Opens current folder in Explorer/Finder
gitmap open "https://github.com"  # Opens in default web browser
gitmap open "mailto:test@example.com" # Opens default mail client
```
