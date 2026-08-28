import os

def fix_amend():
    p = 'gitmap/cmd/amend.go'
    with open(p, 'r', encoding='utf8') as f: c = f.read()
    c = c.replace('func parseAmendOpts(args []string) (amendOptions, []string) {', 'func parseAmendOpts(args []string) (amendOptions, []string, *apperror.AppError) {')
    c = c.replace('return apperror.Wrap(err, "parse flags", nil)', 'return amendOptions{}, nil, apperror.Wrap(err, "parse flags", nil)')
    # wait, amend.go has 'return apperror.New("fatal error", "E9000", nil)' inside parseAmendOpts?
    c = c.replace('return apperror.New("fatal error", "E9000", nil)\n\t}', 'return amendOptions{}, nil, apperror.New("fatal error", "E9000", nil)\n\t}')
    c = c.replace('opts, args := parseAmendOpts(args)', 'opts, args, err := parseAmendOpts(args)\n\tif err != nil { return err }')
    c = c.replace('func loadRecentCommits(limit int) []model.CommitEntry {', 'func loadRecentCommits(limit int) ([]model.CommitEntry, *apperror.AppError) {')
    c = c.replace('return apperror.New("fatal error", "E9000", nil)\n\t}\n\treturn commits', 'return nil, apperror.New("fatal error", "E9000", nil)\n\t}\n\treturn commits, nil')
    c = c.replace('commits := loadRecentCommits(opts.limit)', 'commits, err := loadRecentCommits(opts.limit)\n\tif err != nil { return err }')
    with open(p, 'w', encoding='utf8') as f: f.write(c)

def fix_amendexec():
    p = 'gitmap/cmd/amendexec.go'
    with open(p, 'r', encoding='utf8') as f: c = f.read()
    c = c.replace('func amendCommit(branch string, entry model.CommitEntry, body string, dryRun bool) {', 'func amendCommit(branch string, entry model.CommitEntry, body string, dryRun bool) *apperror.AppError {')
    c = c.replace('return apperror.Wrap(branch, "constants.ErrGitCommitFailed", nil)', 'return apperror.Wrap(err, constants.ErrGitCommitFailed, nil)')
    c = c.replace('amendCommit(opts.branch, chosen, message, opts.dryRun)\n\treturn nil', 'return amendCommit(opts.branch, chosen, message, opts.dryRun)')
    with open(p, 'w', encoding='utf8') as f: f.write(c)

fix_amend()
fix_amendexec()
