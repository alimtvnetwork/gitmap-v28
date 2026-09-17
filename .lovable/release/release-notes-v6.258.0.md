## Quick Install v6.258.0

### Windows (PowerShell 5.1+)
```powershell
irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.258.0/install.ps1 | iex
```

### Linux / macOS (Bash)
```bash
curl -fsSL https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.258.0/install.sh | bash
```

## Changelog v6.258.0

- feat(safety): implement IsProtectedProcess in cli/lockcheck/ to protect Antigravity IDE, VS Code, Node, Electron, WebView2, and parent shells from taskkill and process termination
- feat(safety): guard TerminateProcesses in cli/cmdagy/ to prevent terminating running Antigravity IDE instances during cache cleanup
- feat(ci): add memory-aware CPU freeness worker scaling and cap linter concurrency to prevent Windows commit limit exhaustion (errno 1455)
- feat(test): isolate 37 heavy E2E subprocess tests in cli/tests/heavy_test from routine fast unit test runs
- feat(pipeline,agy): add gitmap pipeline fix errors agy, pipeline-fix, and aef commands with 4-part RCA embedding and duplicate error detection (--force)
- feat(pipeline,agy): dual prompt architecture generating primary fix prompt and staging verification check in prompt queue
- feat(cmd): root CLI shortcuts for fix agy, fix-agy, gitmap aef, and compound phrase token normalization
- feat(cluster,ssh): SSH config sanitizer automatically pruning unsupported client options (authorizedkeysfile)
- feat(clone): add --force / -f destination re-cloning and directory overwrite flag
- feat(cluster,sc): complete command and help parity across cluster and servers-clients (sc) subcommands with interactive examples
- guidelines: author coding guideline 24, master prompt 24, and cg-isolate-os-tests skill safeguarding OS state during test execution
- fix(cli): handle ReplaceModeTypeUnknown in cmd/replace.go exhaustive switch and skip non-supported OS calls on Darwin in osuser_test.go
