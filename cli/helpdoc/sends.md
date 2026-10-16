# gitmap sends

Dispatches conventional commits across targeted repositories with standardized semantic prefixes, non-destructive clean repository skipping, and push automation.

## Usage

```bash
gitmap sends <verb> <target> "<message>" [flags]
```

## Semantic Verbs

| Verb | Description |
|------|-------------|
| `cpf` | Injects `"Feature: "` prefix, stages, commits, and pushes |
| `cpb` | Injects `"Bug: "` prefix, stages, commits, and pushes |
| `cpr` | Injects `"Release: "` prefix, stages, commits, and pushes |
| `commit-fix` | Injects `"Fix: "` prefix, stages, commits, and pushes |
| `cp` | Flat commit without prefix, stages, commits, and pushes |
| `cm` | Local commit only without prefix, does not push |

## Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--dry-run` | `-n` | `false` | Preview commit and push actions without mutating Git refs |
| `--no-push` | | `false` | Commit locally but do not push to upstream remote |
| `--json` | `-j` | `false` | Output structured JSON telemetry |
| `--verbose` | `-v` | `false` | Print detailed git subprocess execution output |
| `-h`, `--help` | | `false` | Display command help menu |

## Examples

### Example 1: Stage, commit, and push feature across all dirty repositories

```bash
gitmap sends cpf all "add unified pending commits suite"
```

### Example 2: Simulate a bug fix commit on a specific repository

```bash
gitmap sends cpb gitmap "fix null pointer in table renderer" -n
```

### Example 3: Commit locally without pushing upstream

```bash
gitmap sends cm all "wip: save local edits" --no-push
```

### Example 4: Workspace-Wide WIP Commits (`sends cp all`)

```bash
gitmap sends cp all "wip: save changes"
```

When executing `gitmap sends cp all`:
1. **Workspace-Wide Discovery**: Queries `gitmap.db` across the entire workspace directory to locate every repository with uncommitted or unstaged changes, rather than scoping only to the current working directory.
2. **Pattern Change Categorization**: Automatically groups changed files into recognizable semantic categories:
   - **Styling & Presentation**: CSS, Tailwind, stylesheets, and visual theme tokens
   - **Core Architecture & Logic**: Go, TypeScript, Rust, Python source files
   - **Boolean Conventions & Guards**: Boolean normalization, positive polarity conditions, guard clauses
   - **Branding & Assets**: Icons, SVG graphics, avatar imagery, and visual media
3. **Relative Path Markdown Linking**: Generates clean commit bodies with markdown links to relative paths (`[path/to/file](path/to/file)`), documenting file changes without absolute path leaks.
4. **Parallel Staging & Push**: Concurrently stages modified files, creates commits with the specified WIP title and formatted body, and pushes to remote tracking branches.

### Example 5: Dry-Run Preview of Workspace Commit Bundles

```bash
gitmap sends cp all "wip: inspect before commit" --dry-run
```
