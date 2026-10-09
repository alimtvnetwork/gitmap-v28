# gitmap fix-link

Inspect, validate, and repair broken or dangling symlinks across Ubuntu, Debian, and Linux systems.

## Usage

```bash
gitmap fix-link [path] [flags]
gitmap os fix-link [path] [flags]
gitmap fixlink [path]
gitmap fl [path]
```

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --target <path> | "" | Explicit target destination for broken link repair |
| --force (-f) | false | Recreate symlink even if target does not yet exist |
| --recursive (-r) | false | Recursively inspect subdirectories when path is a folder |
| --dry-run (-n) | false | Report broken links and proposed fixes without touching disk |
| --json | false | Output link inspection results as structured JSON |

## Examples

### Repair Default Workstation Links

```bash
gitmap fix-link
```

Output:

```text
▶ Inspecting and repairing symlinks...
  ✓ Healthy:  /home/developer/Desktop/SharedDirectories -> /mnt/hgfs
  ✓ Repaired: /usr/local/bin/gitmap -> /home/developer/.local/bin/gitmap
  • Summary:  2 checked, 1 healthy, 1 repaired, 0 broken
```

### Inspect Directory Recursively With Dry Run

```bash
gitmap fix-link ~/projects --recursive --dry-run
```

Output:

```text
▶ Inspecting symlinks (dry run)...
  • Broken:   /home/developer/projects/link -> /var/old/target (missing)
  • Summary:  12 checked, 11 healthy, 0 repaired, 1 broken
```

### Fix Broken Symlink to Specific Target

```bash
gitmap fix-link ~/Desktop/SharedDirectories --target /mnt/hgfs --force
```

Output:

```text
▶ Inspecting and repairing symlinks...
  ✓ Repaired: /home/developer/Desktop/SharedDirectories -> /mnt/hgfs
  • Summary:  1 checked, 1 repaired, 0 broken
```
