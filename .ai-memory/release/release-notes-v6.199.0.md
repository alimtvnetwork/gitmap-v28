## Quick Install v6.199.0

### Windows (PowerShell 5.1+)
```powershell
irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.199.0/install.ps1 | iex
```

### Linux / macOS (Bash)
```bash
curl -fsSL https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.199.0/install.sh | bash
```

## Changelog v6.199.0

- Added gitmap os display command family (display, disp, screen) for desktop session, display server detection (Wayland, X11, DWM, Quartz), screen idle timeout, and never-sleep blanking inhibition
- Implemented SQLite site registry in sites.db with gitmap nginx add <domain> (auto-detecting WordPress and Laravel roots), gitmap nginx rm <domain>, and gitmap nginx list
- Added gitmap nginx ini and showcase displaying recommended production directives for WordPress and Laravel with idempotent marker blocks
- Enhanced VMware shared folder mount resilience with open-vm-tools integration, desktop symlink repair, and reboot crontab persistence
- Integrated Nginx, WordPress, and Laravel setup engines with dynamic PHP-FPM socket discovery
- Added Chrome token vault with reversible ciphers and enhanced Windows command remediation execution
