# gitmap lfs-convert

Scan the current Git repository for tracked files exceeding a specified size threshold (default: 5MB), generate appropriate Git LFS tracking rules in `.gitattributes`, and migrate historical and working-tree binary blobs into lightweight Git LFS pointer files.

## Aliases

- `gitmap lfs convert`
- `gitmap lfs-convert`

## Usage

```bash
gitmap lfs convert [flags]
gitmap lfs-convert [flags]
```

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--size`, `-s` | `5MB` | Size threshold for LFS candidate files (supports `B`, `KB`, `MB`, `GB`) |
| `--dry-run` | `false` | Scan repository and preview candidate files and proposed `.gitattributes` rules without modifying files or Git index |
| `-y`, `--yes` | `false` | Auto-confirm execution prompt; displays the irreversible operation warning |
| `--force` | `false` | Force execution without interactive prompt |

## What It Does

1. **Pre-flight Checks:** Verifies that the current directory is inside a Git repository and checks that `git lfs` is installed and accessible on PATH.
2. **Candidate Discovery:** Scans all Git tracked files (`git ls-files`) and identifies any file whose current size exceeds the threshold (default: 5MB).
3. **Pattern Derivation:** Determines the matching pattern for `.gitattributes` based on file extension (e.g., `*.mp3`, `*.mp4`, `*.psd`, `*.iso`) or full file path for extensionless files.
4. **Git LFS Initialization:** Executes `git lfs install --local` to configure the local Git hooks and smudge/clean filters.
5. **Attributes Configuration:** Appends the canonical tracking rules into `.gitattributes` for any new patterns:
   ```
   <pattern> filter=lfs diff=lfs merge=lfs -text
   ```
6. **Blob Migration:** Executes `git lfs migrate import --everything --include="<patterns>"` to rewrite Git commit objects, replace raw binary blobs with tiny text pointer files (~130 bytes), and drastically decrease repository clone times and disk usage.

## Warning & Irreversibility

Converting historical Git blobs to Git LFS pointers alters commit SHAs across rewritten branch history. When running in active mode without `-y`, the command prompts:

```
WARNING: This operation rewrites Git history and converts files to Git LFS pointers. It is irreversible.
Do you want to proceed? [y/N]:
```

When `-y` or `--yes` is specified, the prompt is skipped and the bold yellow warning is output directly to the terminal.

> [!WARNING]
> Because Git history is rewritten, any collaborator branches will diverge. Downstream clones must run `git lfs pull` or re-clone the repository once changes are published to a remote.

## Recovery & Inspection

To inspect files currently tracked by Git LFS:

```bash
git lfs ls-files
```

To verify Git LFS status in the repository:

```bash
git lfs status
```

To pull the full binary payloads for LFS pointers:

```bash
git lfs pull
```

## Examples

Preview files exceeding default 5MB threshold and their corresponding LFS rules:
```bash
gitmap lfs convert --dry-run
```

Preview files exceeding 20MB:
```bash
gitmap lfs-convert --size 20MB --dry-run
```

Convert all files exceeding 5MB to Git LFS with auto-confirmation:
```bash
gitmap lfs convert -y
```

Convert files exceeding 10MB using force mode:
```bash
gitmap lfs-convert --size 10MB --force
```

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter lfs-convert
```

The JSON schema is published at `02-spec/08-json-schemas/help-json.schema.json`.
