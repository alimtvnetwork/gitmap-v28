package cmdos

import "fmt"

const osUsageText = `Usage: gitmap os [subcommand] [flags]

Commands:
  ip                  Inspect, set, change, switch, or revert network IP configuration
  fix                 Register, edit, run, export, and import system repair scripts
  clean (clear)       Clean temporary and ephemeral system cache directories
  dev-clean (dev clean) Clean compiler, package manager, and build tool caches
  ai-clean (aiclean)  Scan and purge Antigravity brain, task, and temp AI cache dumps
  zsh                 Install, theme, switch, profile, and clean ZSH & Oh-My-Zsh
  user                Add, edit, export, import, or remove operating system users
  group (user-group)  List, create, edit, export, import, and remove user groups
  vmware              Discover and mount VMware shared folders (/mnt/hgfs)
  cron                Inspect, append, and remove crontab scheduled jobs
  storage (disk)      Inspect disk drive capacities, partitions, and storage metrics
  autologin           Configure OS auto-login credentials (Windows & Ubuntu)
  tweak               Manage Windows desktop tweaks, power schemes, and hibernation
  dm                  Inspect and configure Linux Display Managers & Wayland
  dns                 Inspect, set, benchmark, and revert network DNS servers
  theme               Switch desktop visual appearance between dark and light
  update              Check and update package repository metadata
  upgrade             Upgrade all packages and toolchains across detected package managers
  display (disp)      Inspect and configure OS display settings, resolution & timeouts
  fix-link (fixlink)  Inspect and repair broken symlinks and shared directories
  status (st)         Display operating system environment and link diagnostics
  help                Show this help message

Flags:
  --target <path>     Explicit target for symlink repair
  --force (-f)        Recreate symlinks even if target missing
  --recursive (-r)    Recursively inspect directories
  --dry-run (-n)      Inspect without modifying disk
  --json              Output as JSON`

func printOSUsage() {
	fmt.Println(osUsageText)
}
