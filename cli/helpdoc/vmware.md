# gitmap vmware

Manage VMware guest integration, shared folders mounting, desktop symlinks, and tool installation across Linux and Windows.

## Usage

```bash
gitmap vmware [subcommand] [flags]
gitmap vm [subcommand]
```

## Subcommands

| Subcommand | Alias | Description |
|------------|-------|-------------|
| install | in | Install `open-vm-tools` via apt (Linux) or verify `VMTools` service & `vmrun` utility (Windows) |
| shared enable | mount | Mount `/mnt/hgfs` & persist in crontab (Linux) or validate UNC share & create Desktop shortcut (Windows) |
| shared status | | Check `/mnt/hgfs` mount & crontab persistence (Linux) or UNC share access & Desktop shortcut (Windows) |
| status | | Display comprehensive VMware hypervisor, guest OS, and tools status |

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| -y, --yes | false | Auto-confirm package installation without prompting |
| -n, --dry-run | false | Simulate VMware actions without modifying the system |
| -h, --help | false | Show command help |

## Cross-Platform Features

### Linux Guests

- Installs `open-vm-tools` and `open-vm-tools-desktop` via `apt-get`.
- Mounts VMware shared folders at `/mnt/hgfs` using `vmhgfs-fuse`.
- Symlinks `/mnt/hgfs` to `~/Desktop/SharedDirectories`.
- Persists automount across reboots via `@reboot` crontab entry.

### Windows Guests

- Detects VMware guest environment via Windows Registry (`SOFTWARE\VMware, Inc.\VMware Tools`) and Service Control Manager (`sc query VMTools`).
- Automatically searches for `vmrun.exe` in `PATH` and standard installation directories (`Program Files`, `Program Files (x86)`).
- Validates UNC shared folder availability at `\\vmware-host\Shared Folders`.
- Generates desktop shortcut to `\\vmware-host\Shared Folders` on the user's Desktop.

## Database Tracking & Audit

All `gitmap vmware` operations automatically record execution telemetry into the isolated system SQLite database (`installation.db`):
- Action name (`install`, `shared`, `shared-enable`, `status`).
- Duration in milliseconds.
- Success status and exit codes.
- Error messages and stack traces if an operation fails.

## Examples

### Install / Verify VMware Tools

```bash
gitmap vmware install -y
```

Output (Linux):

```text
▶ gitmap vmware install
  ✓ Successfully installed open-vm-tools and open-vm-tools-desktop
  ✓ Service open-vm-tools enabled and started
  Next: run 'gitmap vmware shared enable' to mount shared folders
```

Output (Windows):

```text
▶ gitmap vmware install (Windows)
  VMware Guest: true
  VMware Tools Service: RUNNING (running=true)
  vmrun utility: C:\Program Files (x86)\VMware\VMware Workstation\vmrun.exe
  ✓ VMware Tools is installed and active on Windows guest
```

### Enable VMware Shared Folders

```bash
gitmap vmware shared enable
```

Output (Linux):

```text
▶ gitmap vmware shared enable
  ✓ Verified mount point /mnt/hgfs
  ✓ Mounted .host:/ at /mnt/hgfs
  ✓ Created Desktop/SharedDirectories symlink
  ✓ Registered @reboot crontab persistence
```

Output (Windows):

```text
▶ gitmap vmware shared enable (Windows)
  ✓ UNC path \\vmware-host\Shared Folders is accessible
  ✓ Created Desktop shortcut to VMware Shared Folders
```

### Check Shared Folders Status

```bash
gitmap vmware shared status
```

Output (Linux):

```text
▶ gitmap vmware shared status
  Mount (/mnt/hgfs): active=true
  Desktop Symlink: present=true (target=/mnt/hgfs)
  Crontab Persistence: registered=true
```

Output (Windows):

```text
▶ gitmap vmware shared status (Windows)
  UNC Shared Folders (\\vmware-host\Shared Folders): accessible=true
  Desktop Shortcut: present=true
```

### Inspect Hypervisor & Guest Status

```bash
gitmap vmware status
```

Output:

```text
▶ gitmap vmware status
  OS Platform: Windows
  VMware Guest: true
  VMware Tools Service: RUNNING (running=true)
  vmrun Installed: true (path=C:\Program Files (x86)\VMware\VMware Workstation\vmrun.exe)
  Shared Folders Accessible: true
```

## Troubleshooting

### Linux: Error -107 cannot open connection!

If `gitmap vmware shared enable` reports:

```text
mount failed: Error -107 cannot open connection! (exit status 149)
```

This indicates the Linux guest kernel transport endpoint cannot communicate with the VMware host hypervisor.

1. In VMware Workstation / Player / Fusion:
   - Go to **Virtual Machine Settings** -> **Options** tab -> **Shared Folders**.
   - Select **Always enabled**.
   - Under **Folders**, click **Add...** and choose at least one host directory.
   - Click **OK** to save VM settings.
2. Restart guest service: `sudo systemctl restart open-vm-tools`.
3. Re-run: `gitmap vmware shared enable`.

### Windows: UNC Path Not Accessible

If `\\vmware-host\Shared Folders` is not accessible on Windows:
1. Open VM Settings -> Options -> Shared Folders.
2. Ensure **Always enabled** is selected and at least one host directory is shared and checked.
3. Ensure the `VMTools` service is running via `sc query VMTools` or Services Manager (`services.msc`).
