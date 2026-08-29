# Task 3: OS Startup Macro Hook & OS Management

- Implement OS-level startup triggers in `gitmap/osutil/startup.go`.
- Windows: Add to Registry `Run` keys. Linux: Add to `~/.config/autostart/`.
- Add commands: `gitmap schedule restart` and `gitmap schedule shutdown` executing native OS commands.
- Follow error handling guidelines (AppError).
