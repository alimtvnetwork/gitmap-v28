# gitmap vmware

Manage VMware guest integration, shared folders mounting, and desktop symlinks.

## Usage

```bash
gitmap vmware [subcommand] [flags]
gitmap vm [subcommand]
```

## Subcommands

| Subcommand | Description |
|------------|-------------|
| shared enable | Enable VMware shared folders, mount /mnt/hgfs, symlink ~/Desktop/SharedDirectories, and configure @reboot crontab |
| shared status | Check status of VMware tools, mount point, and shared directories |

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --dry-run | false | Simulate VMware commands without executing |

## Examples

### Enable VMware Shared Folders

```bash
gitmap vmware shared enable
```

Output:

```text
▶ Configuring VMware Shared Folders...
  ✓ Installed open-vm-tools
  ✓ Mounted .host:/ to /mnt/hgfs
  ✓ Created symlink ~/Desktop/SharedDirectories
  ✓ Persisted mount command in crontab (@reboot)
```

### Check Shared Folders Status

```bash
gitmap vmware shared status
```

Output:

```text
▶ VMware Shared Folders Status:
  • open-vm-tools: installed
  • Mount point (/mnt/hgfs): mounted
  • Desktop symlink: active
```
