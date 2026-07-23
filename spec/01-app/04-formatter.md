# Formatter

## Responsibility

Render a list of `ScanRecord` into multiple output formats.

## Output Behavior

Every scan **always produces all outputs simultaneously**:

1. **Terminal** — colored, professional output to stdout.
2. **CSV** — `gitmap.csv` written to the output directory.
3. **JSON** — `gitmap.json` written to the output directory.
4. **Folder Structure** — `folder-structure.md` written to the output directory.
5. **Clone Script** — `clone.ps1` — self-contained PowerShell script that clones all repos.
6. **Direct Clone** — `direct-clone.ps1` — plain HTTPS git clone commands, one per line.
7. **Direct Clone SSH** — `direct-clone-ssh.ps1` — plain SSH git clone commands, one per line.
8. **Desktop Script** — `register-desktop.ps1` — registers cloned repos with GitHub Desktop.

The output directory defaults to `.gitmap/output/` inside the scanned directory.

## Formats

### Terminal

Colored output with ANSI codes showing:
- Banner with repo count
- Each repo: name (📦), path, and clone instruction
- Folder tree with 📁 folders and 📦 repos with branch names
- Clone help: step-by-step instructions for cloning on another machine

### CSV

Write a CSV file with headers:

```
repoName,httpsUrl,sshUrl,branch,relativePath,absolutePath,cloneInstruction,notes
```

### JSON

Write a JSON array of `ScanRecord` objects with 2-space indentation.

### Folder Structure (Markdown)

Write a tree view of discovered repos:

```markdown
# Folder Structure

Git repositories discovered by gitmap.

├── 📦 **my-app** (`main`) — https://github.com/user/my-app.git
├── libs/
│   └── 📦 **utils** (`main`) — https://github.com/user/utils.git
└── 📦 **docs** (`main`) — https://github.com/user/docs.git
```

### Clone Script (`clone.ps1`)

Generated from `formatter/templates/clone.ps1.tmpl` using Go `text/template`
with `go:embed`. The template receives an array of repo entries and produces
a PowerShell script with a **data-driven loop** instead of repeating blocks:

```powershell
$repos = @(
    @{ Name = "my-app"; Branch = "main"; URL = "https://..."; Path = "my-app" }
    @{ Name = "utils";  Branch = "main"; URL = "https://..."; Path = "libs\utils" }
)

foreach ($repo in $repos) {
    # mkdir, clone, report
}
```

- Accepts a `-TargetDir` parameter (defaults to `.`)
- Creates the folder structure under the target directory
- Clones each repo with `git clone -b <branch> <url> <path>`
- Shows progress (`[1/N]`, `[2/N]`, …) with colored output
- Prints a summary of succeeded/failed clones

### Desktop Registration Script (`register-desktop.ps1`)

Generated from `formatter/templates/desktop.ps1.tmpl` using the same
template engine. Uses the same array-loop pattern:

- Accepts a `-BaseDir` parameter (defaults to `.`)
- Checks if GitHub Desktop CLI (`github`) is available
- Registers each cloned repo with GitHub Desktop
- Shows progress with colored output
- Prints a summary of registered/failed repos

### Template Architecture

Templates live in `formatter/templates/` and are embedded into the binary
via `//go:embed templates/*` in `formatter/template.go`. This means:

- **No external file dependencies** — templates are compiled into the binary
- **Editable without Go knowledge** — PowerShell templates are standalone files
- **Testable independently** — templates can be validated outside Go
- Template data structs: `CloneData`, `DesktopData`, `RepoEntry`

## Output Location

- Terminal: stdout.
- All files: `.gitmap/output/` inside the scanned directory,
  or path from `--output-path` flag, or exact path from `--out-file`.

## Output Directory Contents

```
.gitmap/output/
├── gitmap.csv
├── gitmap.json
├── folder-structure.md
├── clone.ps1
├── direct-clone.ps1
├── direct-clone-ssh.ps1
└── register-desktop.ps1
```

## Template Files

```
formatter/
├── templates/
│   ├── clone.ps1.tmpl             # PowerShell clone script template
│   ├── direct-clone.ps1.tmpl      # Plain HTTPS clone commands template
│   ├── direct-clone-ssh.ps1.tmpl  # Plain SSH clone commands template
│   └── desktop.ps1.tmpl           # GitHub Desktop registration template
├── template.go              # go:embed loader + shared types
├── clonescript.go            # WriteCloneScript (uses clone.ps1.tmpl)
├── directclone.go            # WriteDirectCloneScript + WriteDirectCloneSSHScript
├── desktopscript.go          # WriteDesktopScript (uses desktop.ps1.tmpl)
└── terminal.go               # Terminal output with version banner
```
