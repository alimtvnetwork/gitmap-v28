# `gitmap open`

Opens the current repository (or a specified target directory) using your operating system's native file explorer or default handler. 

## Usage

```bash
gitmap open [target]
```
*(Alias: `gitmap o`)*

## Behavior

* **No arguments**: 
  Discovers the root of the Git repository for the current working directory and opens it. If you are not in a Git repository, it falls back to opening the current working directory.
* **With `[target]`**:
  Opens the exact file or directory path provided.

## Operating System Integration

The command integrates deeply with native OS protocols to ensure the target is opened just like you double-clicked it:

* **Windows**: `rundll32 url.dll,FileProtocolHandler`
* **macOS**: `open`
* **Linux**: `xdg-open`
