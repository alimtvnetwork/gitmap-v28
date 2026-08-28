import re

p = 'gitmap/cmd/bookmarksave.go'
with open(p, 'r', encoding='utf8') as f: c = f.read()
c = re.sub(r'return apperror.New\(constants.ErrBookmarkExists, "E9000", nil\)\n\}', 'return apperror.New(constants.ErrBookmarkExists, "E9000", nil)\n\t}\n\treturn nil\n}', c)
with open(p, 'w', encoding='utf8') as f: f.write(c)

p = 'gitmap/cmd/cdops.go'
with open(p, 'r', encoding='utf8') as f: c = f.read()
c = c.replace('return repos, nil\n}', 'return repos\n}')
c = c.replace('func loadRecentRepos() ([]model.ScanRecord, *apperror.AppError) {\n\trepos, err := store.LoadRecentRepos()\n\tif err != nil {\n\t\treturn nil, apperror.Wrap(err, "constants.ErrListDBFailed", nil)\n\t}\n\n\treturn repos\n}', 'func loadRecentRepos() ([]model.ScanRecord, *apperror.AppError) {\n\trepos, err := store.LoadRecentRepos()\n\tif err != nil {\n\t\treturn nil, apperror.Wrap(err, "constants.ErrListDBFailed", nil)\n\t}\n\n\treturn repos, nil\n}')
with open(p, 'w', encoding='utf8') as f: f.write(c)

p = 'gitmap/cmd/changelog.go'
with open(p, 'r', encoding='utf8') as f: c = f.read()
c = c.replace('return apperror.Wrap(constants.ChangelogFile, "constants.ErrChangelogRead", nil)', 'return apperror.Wrap(err, constants.ErrChangelogRead, nil)')
c = c.replace('return apperror.Wrap(constants.ChangelogFile, "constants.ErrChangelogOpen", nil)', 'return apperror.Wrap(err, constants.ErrChangelogOpen, nil)')
c = c.replace('return nil\n\t}\n\treturn nil\n}', 'return nil\n\t}\n\treturn nil\n}')
with open(p, 'w', encoding='utf8') as f: f.write(c)

p = 'gitmap/cmd/changeloggen.go'
with open(p, 'r', encoding='utf8') as f: c = f.read()
c = c.replace('func handleGenerate(dir string, opts changelogOptions) {', 'func handleGenerate(dir string, opts changelogOptions) *apperror.AppError {')
c = c.replace('return apperror.Wrap(constants.ChangelogFile, "constants.ErrChangelogOpen", nil)', 'return apperror.Wrap(err, constants.ErrChangelogOpen, nil)')
c = c.replace('return apperror.Wrap(constants.ChangelogFile, "constants.ErrChangelogWrite", nil)', 'return apperror.Wrap(err, constants.ErrChangelogWrite, nil)')
c = c.replace('handleGenerate(dir, opts)\n\treturn nil', 'return handleGenerate(dir, opts)')
c = c.replace('fmt.Printf(constants.MsgChangelogWrittenFmt, constants.ChangelogFile)\n}', 'fmt.Printf(constants.MsgChangelogWrittenFmt, constants.ChangelogFile)\n\treturn nil\n}')
with open(p, 'w', encoding='utf8') as f: f.write(c)
