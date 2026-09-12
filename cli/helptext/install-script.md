# install-script

Shows and copies the standard one-line `gitmap` installation script to your clipboard. 
It automatically detects your current operating system (Windows vs Linux/macOS) and copies the correct command.

## Usage

`gitmap install-script`

## Examples

```bash
gitmap install-script

# Output: Install script copied to clipboard!

# Windows clipboard: irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/install.ps1 | iex

# Unix clipboard: curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/install.sh | sh

```
