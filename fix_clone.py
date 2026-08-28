import os

def r(p, o, n):
    with open(p, 'r', encoding='utf8') as f: c = f.read()
    c = c.replace(o, n)
    with open(p, 'w', encoding='utf8') as f: f.write(c)

r('gitmap/cmd/chrome_backup.go', 'func doDryRun(src, dst string) {', 'func doDryRun(src, dst string) *apperror.AppError {')
r('gitmap/cmd/chrome_backup.go', 'fmt.Printf("\\033[1;93m\\u2713 chrome restore (dry-run)\\033[0m  %d file(s) would land in \\033[1;96m%s\\033[0m\\n", n, dst)\n}', 'fmt.Printf("\\033[1;93m\\u2713 chrome restore (dry-run)\\033[0m  %d file(s) would land in \\033[1;96m%s\\033[0m\\n", n, dst)\n\treturn nil\n}')
r('gitmap/cmd/chrome_backup.go', 'doDryRun(src, dst)\n\t\treturn nil', 'return doDryRun(src, dst)')

r('gitmap/cmd/clone.go', 'func applySSHKey(name string) {', 'func applySSHKey(name string) *apperror.AppError {')
r('gitmap/cmd/clone.go', 'if len(name) == 0 {\n\t\treturn\n\t}', 'if len(name) == 0 {\n\t\treturn nil\n\t}')
r('gitmap/cmd/clone.go', 'fmt.Fprintf(os.Stdout, constants.MsgSSHCloneUsing, name, key.PrivatePath)\n}', 'fmt.Fprintf(os.Stdout, constants.MsgSSHCloneUsing, name, key.PrivatePath)\n\treturn nil\n}')
r('gitmap/cmd/clone.go', 'applySSHKey(cf.SSHKey)', 'if err := applySSHKey(cf.SSHKey); err != nil { return err }')

r('gitmap/cmd/clone.go', 'func executeDirectClone(url, folderName string, ghDesktopFlag, noReplace bool, output string, noVSCodeSync bool) {', 'func executeDirectClone(url, folderName string, ghDesktopFlag, noReplace bool, output string, noVSCodeSync bool) *apperror.AppError {')
r('gitmap/cmd/clone.go', 'fmt.Fprintf(os.Stderr, "  \\033[1;92m\\u2714\\033[0m  Task %s spawned and running.\\n", taskID)\n\t}\n}', 'fmt.Fprintf(os.Stderr, "  \\033[1;92m\\u2714\\033[0m  Task %s spawned and running.\\n", taskID)\n\t}\n\treturn nil\n}')
r('gitmap/cmd/clone.go', 'executeDirectClone(url, folderName, cf.GithubDesktop, cf.NoReplace, cf.Output, cf.NoVSCodeSync)\n\t\treturn nil', 'return executeDirectClone(url, folderName, cf.GithubDesktop, cf.NoReplace, cf.Output, cf.NoVSCodeSync)')
r('gitmap/cmd/clone.go', 'executeDirectClone(rootName, rootName, cf.GithubDesktop, cf.NoReplace, cf.Output, cf.NoVSCodeSync)\n\treturn nil', 'return executeDirectClone(rootName, rootName, cf.GithubDesktop, cf.NoReplace, cf.Output, cf.NoVSCodeSync)')

r('gitmap/cmd/clone.go', 'func validateShorthandPath(resolved string) string {', 'func validateShorthandPath(resolved string) (string, *apperror.AppError) {')
r('gitmap/cmd/clone.go', 'return resolved\n\t}', 'return resolved, nil\n\t}')
r('gitmap/cmd/clone.go', 'return apperror.New(constants.ErrShorthandNotFound, "E9000", nil)\n\n\treturn ""\n}', 'return "", apperror.New(constants.ErrShorthandNotFound, "E9000", nil)\n}')
r('gitmap/cmd/clone.go', 'return validateShorthandPath(resolved)', 'path, _ := validateShorthandPath(resolved)\n\t\treturn path')

r('gitmap/cmd/clone.go', 'func executeClone(source, targetDir string, safePull, ghDesktop bool, maxConcurrency int, defaultBranch string, noVSCodeSync bool, clean bool, missingOnly bool) {', 'func executeClone(source, targetDir string, safePull, ghDesktop bool, maxConcurrency int, defaultBranch string, noVSCodeSync bool, clean bool, missingOnly bool) *apperror.AppError {')
r('gitmap/cmd/clone.go', 'fmt.Fprintf(os.Stderr, "  \\033[1;92m\\u2714\\033[0m  Task %s spawned and running.\\n", taskID)\n}', 'fmt.Fprintf(os.Stderr, "  \\033[1;92m\\u2714\\033[0m  Task %s spawned and running.\\n", taskID)\n\treturn nil\n}')
r('gitmap/cmd/clone.go', 'executeClone(source, target, cf.SafePull, cf.GithubDesktop, cf.MaxConcurrency, cf.DefaultBranch, cf.NoVSCodeSync, cf.Clean, cf.MissingOnly)\n\treturn nil', 'return executeClone(source, target, cf.SafePull, cf.GithubDesktop, cf.MaxConcurrency, cf.DefaultBranch, cf.NoVSCodeSync, cf.Clean, cf.MissingOnly)')
