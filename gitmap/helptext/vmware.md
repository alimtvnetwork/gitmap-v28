# gitmap vmware

Manage VMware guest integration, shared folders mounting, desktop symlinks, and tool installation.

## Usage

```bash
gitmap vmware [subcommand] [flags]
gitmap vm [subcommand]
```

## Subcommands

| Subcommand | Alias | Description |
|------------|-------|-------------|
| install | in | Install `open-vm-tools` and `open-vm-tools-desktop` via apt and enable systemd service |
| shared enable | mount | Enable VMware shared folders, mount `/mnt/hgfs`, symlink `~/Desktop/SharedDirectories`, and configure `@reboot` crontab |
| shared status | | Check status of VMware tools, mount point, and shared directories |
| status | | Display VMware guest environment detection status |

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| -y, --yes | false | Auto-confirm package installation without prompting |
| -n, --dry-run | false | Simulate VMware actions without mounting or writing crontab |
| -h, --help | false | Show command help |

## Examples

### Install VMware Tools

```bash
gitmap vmware install -y
```

Output:

```text
▶ gitmap vmware install
  ✓ Successfully installed open-vm-tools and open-vm-tools-desktop
  ✓ Service open-vm-tools enabled and started
  Next: run 'gitmap vmware shared enable' to mount shared folders
```

### Enable VMware Shared Folders

```bash
gitmap vmware shared enable
```

Output:

```text
▶ gitmap vmware shared enable
  ✓ Verified mount point /mnt/hgfs
  ✓ Mounted .host:/ at /mnt/hgfs
  ✓ Created Desktop/SharedDirectories symlink
  ✓ Registered @reboot crontab persistence
```

### Check Shared Folders Status

```bash
gitmap vmware shared status
```

Output:

```text
▶ gitmap vmware shared status
  Mount (/mnt/hgfs): active=true
  Desktop Symlink: present=true (target=/mnt/hgfs)
```

## Troubleshooting: Error -107 cannot open connection!

If `gitmap vmware shared enable` reports:

```text
mount failed: Error -107 cannot open connection! (exit status 149)
```

This indicates the Linux guest kernel transport endpoint cannot communicate with the VMware host hypervisor.

### Resolution Steps

1. In VMware Workstation / Player / Fusion:
   - Go to **Virtual Machine Settings** -> **Options** tab -> **Shared Folders**.
   - Select **Always enabled** (or *Enabled until next power off*).
   - Under **Folders**, click **Add...** and choose at least one host directory (e.g. `D:\work` or `C:\Users`).
   - Ensure the folder checkbox is checked (Enabled).
   - Click **OK** to save VM settings.
2. Ensure guest service is running:
   ```bash
   sudo systemctl restart open-vm-tools
   ```
3. Re-run:
   ```bash
   gitmap vmware shared enable
   ```
