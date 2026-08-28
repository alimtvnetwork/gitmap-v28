import re

def r(p, o, n):
    with open(p, 'r', encoding='utf8') as f: c = f.read()
    c = c.replace(o, n)
    with open(p, 'w', encoding='utf8') as f: f.write(c)

# clustercommand.go
r('gitmap/cmd/clustercommand.go', 'fmt.Fprintf(os.Stderr, "Preflight error: %v\\n", err)\n\t\tos.Exit(1)', 'return apperror.Wrap(err, "Preflight error", nil)')
r('gitmap/cmd/clustercommand.go', 'fmt.Fprintln(os.Stderr, "Operation aborted.")\n\t\tos.Exit(1)', 'return apperror.New("Operation aborted", "E9000", nil)')
r('gitmap/cmd/clustercommand.go', 'return apperror.New("Operation aborted", "E9000", nil)\n\t}', 'return apperror.New("Operation aborted", "E9000", nil)\n\t}\n\treturn nil')

# code.go
r('gitmap/cmd/code.go', 'fmt.Fprintln(os.Stderr, err.Error())\n\t\tos.Exit(1)', 'return apperror.Wrap(err, "error", nil)')
r('gitmap/cmd/code.go', 'return apperror.Wrap(err, "error", nil)\n\t}', 'return apperror.Wrap(err, "error", nil)\n\t}\n\treturn nil')

# committransfer.go
r('gitmap/cmd/committransfer.go', 'func checkAuthorIdentity(spec *commitTransferSpec) {', 'func checkAuthorIdentity(spec *commitTransferSpec) *apperror.AppError {')
r('gitmap/cmd/committransfer.go', 'fmt.Fprintf(os.Stderr, "[%s] \\033[1;31mError\\033[0m: No author identity detected (missing both local and global git config for name/email).\\n", spec.Name)\n\t\tos.Exit(1)', 'return apperror.New("No author identity detected", "E9000", nil)')
r('gitmap/cmd/committransfer.go', 'fmt.Fprintf(os.Stderr, "         Please configure git user.name and user.email first.\\n")\n\t\tos.Exit(1)', '')
r('gitmap/cmd/committransfer.go', 'checkAuthorIdentity(spec)\n\n\treturn spec', 'if err := checkAuthorIdentity(spec); err != nil { return nil, err }\n\n\treturn spec')
r('gitmap/cmd/committransfer.go', 'return spec, nil\n}', 'return spec, nil\n}')
r('gitmap/cmd/committransfer.go', 'func handleNoCommits(opts commitTransferOptions) {', 'func handleNoCommits(opts commitTransferOptions) *apperror.AppError {')
r('gitmap/cmd/committransfer.go', 'fmt.Fprintf(os.Stderr, "[%s] No commits found on %s\\n", opts.LogPrefix, opts.SourceBranch)\n\t\tos.Exit(0)', 'return nil')
r('gitmap/cmd/committransfer.go', 'handleNoCommits(opts)\n\t\treturn', 'if err := handleNoCommits(opts); err != nil { return err }\n\t\treturn nil')

r('gitmap/cmd/completion.go', 'func requireShellArg(args []string) {', 'func requireShellArg(args []string) *apperror.AppError {')
r('gitmap/cmd/completion.go', 'fmt.Fprintln(os.Stderr, constants.ErrCompletionUsage)\n\t\tos.Exit(1)', 'return apperror.New(constants.ErrCompletionUsage, "E9000", nil)')
r('gitmap/cmd/completion.go', 'return apperror.New(constants.ErrCompletionUsage, "E9000", nil)\n\t}', 'return apperror.New(constants.ErrCompletionUsage, "E9000", nil)\n\t}\n\treturn nil')
r('gitmap/cmd/completion.go', 'requireShellArg(args)', 'if err := requireShellArg(args); err != nil { return err }')

