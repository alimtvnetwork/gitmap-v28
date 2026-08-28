# Subtask 3: Implement `gitmap open`

1. Create `gitmap/cmd/open.go`.
2. Implement `runOpen(args []string) error`.
3. Check `runtime.GOOS` to use the correct opener:
   - Windows: `exec.Command("rundll32", "url.dll,FileProtocolHandler", target)` or `cmd /c start`
   - Darwin: `exec.Command("open", target)`
   - Linux: `exec.Command("xdg-open", target)`
4. Register `CmdOpen = "open"` in `constants_cli.go` and `cmd_constants_test.go`.
5. Map it in `roottooling.go`.
6. Create `helptext/open.md`.
