## Quick Install v6.257.0

### Windows (PowerShell 5.1+)
```powershell
irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.257.0/install.ps1 | iex
```

### Linux / macOS (Bash)
```bash
curl -fsSL https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.257.0/install.sh | bash
```

## Changelog v6.257.0

- feat(schedule): flexible shutdown and restart duration parsing (1:45hr, 1:45h, 2h, 120m, 1day, 1d, now) with native OS elevation
- feat(schedule): power schedule state persistence, countdown inspection (schedule shutdown status), and cancellation abort (schedule shutdown cancel)
- feat(cluster): remote GitMap installation (cluster install gitmap, sj install gitmap) with official curl and PowerShell one-liners
- feat(cluster): automatic preflight GitMap detection and on-the-fly bootstrapping during remote cluster and servers-clients command execution
- feat(os): gitmap os ai-clean and scripts/os-ai-clean.py scanning Antigravity brain caches, system task logs, and temp AI dumps with preflight table confirmation
- docs(help): comprehensive cluster triad architecture documentation unifying ssh-join, cluster, and servers-clients with rich copy-pasteable examples
- docs(help): leaf help topics registered for cluster-install, schedule-shutdown, schedule-restart, and os-ai-clean
- ui(table): fixed-width ASCII column formatting and divider alignment for sj ls and cluster nodes
- guidelines: prompt 24 (isolate destructive OS and heavy unit tests) authored, indexed, and synchronized across gitmap and coding-guidelines
- test(qa): 100% green verification on 38 linters, compile gates, 118 E2E smoke tests, and zero nested ifs across 3059 files
