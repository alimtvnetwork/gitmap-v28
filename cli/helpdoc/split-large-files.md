# gitmap split-large-files

Scan the current Git repository for tracked files exceeding a specified size threshold (default: 10MB) and decompose them into multipart ZIP archives (`.zip`, `.z01`, `.z02`, ...). The raw monolithic blobs are removed from Git tracking and replaced with staged multipart archive pieces, keeping Git repository size lightweight and performant.

## Aliases

- `gitmap split large files`
- `gitmap split-large-files`

## Usage

```bash
gitmap split large files [flags]
gitmap split-large-files [flags]
```

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--size`, `-s` | `10MB` | Size threshold for splitting files (supports `B`, `KB`, `MB`, `GB`) |
| `--dry-run` | `false` | Scan repository and preview projected chunks and compression modes without modifying files or Git index |
| `-y`, `--yes` | `false` | Auto-confirm execution prompt; displays the irreversible operation warning |
| `--force` | `false` | Force execution without interactive prompt |

## Compression Heuristic

Each large file is evaluated by file extension to determine the optimal compression method:

- **Pre-compressed Media:** Files with extensions `.png`, `.jpg`, `.jpeg`, `.webp`, `.mp4`, `.mov`, `.mkv`, `.mp3`, `.wav`, `.aac`, `.m4a` use `zip.Store` (Method 0). Recompressing already compressed audio, video, and image codecs wastes CPU cycles and offers no space savings. `zip.Store` performs pure fast streaming chunking at full disk I/O throughput.
- **Documents, Code, and Data:** Files with text, SQL dumps, database files, JSON, CSV, binary dumps, or uncompressed formats use `zip.Deflate` (Method 8) to achieve maximum compression ratio before chunking.

## Multipart ZIP Archive Structure

The generated multipart archives follow standard multi-part ZIP naming:

- **First chunk:** `<filename>.zip`
- **Subsequent chunks:** `<filename>.z01`, `<filename>.z02`, `<filename>.z03`, ...

All chunks except the final chunk are capped precisely at the configured `--size` threshold.

### Extraction & Reassembly

The resulting multipart archives can be extracted directly using standard archival utilities or concatenated into a single ZIP archive:

#### 1. Direct Extraction (7-Zip / WinRAR)
Open `<filename>.zip` with 7-Zip or WinRAR. The tool automatically detects accompanying `.z01`, `.z02`, ... parts in the same directory and extracts the original file.

#### 2. Linux / macOS / Bash
Concatenate the parts in sequential order and unzip:

```bash
cat filename.iso.zip filename.iso.z* > filename.iso.zip.full
unzip filename.iso.zip.full
```

#### 3. Windows PowerShell
Combine parts byte-by-byte into a single archive:

```powershell
Get-Content filename.iso.zip, filename.iso.z* -AsByteStream | Set-Content filename.iso.zip.full -AsByteStream
Expand-Archive -Path filename.iso.zip.full -DestinationPath .
```

#### 4. Python
Extract using Python standard library:

```python
import glob
import zipfile

base_name = "filename.iso"
parts = [f"{base_name}.zip"] + sorted(glob.glob(f"{base_name}.z*"))

with open(f"{base_name}.full.zip", "wb") as out_f:
    for part in parts:
        with open(part, "rb") as in_f:
            out_f.write(in_f.read())

with zipfile.ZipFile(f"{base_name}.full.zip", "r") as zf:
    zf.extractall(".")
```

## Warning & Irreversibility

When run in active mode without `-y`, the command prompts for confirmation:

```
WARNING: This operation rewrites Git history and splits files into multipart zip archives. It is irreversible.
Do you want to proceed? [y/N]:
```

When `-y` or `--yes` is specified, the prompt is bypassed and the bold warning is displayed prior to execution.

> [!WARNING]
> Splitting removes the original monolithic files from the Git staging index (`git rm --cached`) and deletes the local uncompressed file. Ensure any uncommitted local work is backed up prior to running.

## Examples

Preview files exceeding default 10MB threshold without changing anything:
```bash
gitmap split large files --dry-run
```

Preview files exceeding 25MB:
```bash
gitmap split-large-files --size 25MB --dry-run
```

Execute split on files exceeding 10MB with auto-confirmation:
```bash
gitmap split large files -y
```

Force split with custom 50MB chunk size:
```bash
gitmap split-large-files --size 50MB --force
```

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter split-large-files
```

The JSON schema is published at `02-spec/08-json-schemas/help-json.schema.json`.
