# Learned: SSH Login and Host Join

## Architectural Patterns

1. **Host Target Format**:
   - Accepts targets in format `user@ip:port` or `ip@user` or `alias`.
   - Defaults: user `root`, port `22`.

2. **Cross-Platform IP Management**:
   - Windows uses `netsh interface ip set address`.
   - Linux uses `ip addr add` / `ip addr del`.
   - Multi-platform ping validation tests network reachability before confirming IP change, with automated rollback on failure.
