import os

def r(p, o, n):
    with open(p, 'r', encoding='utf8') as f: c = f.read()
    c = c.replace(o, n)
    with open(p, 'w', encoding='utf8') as f: f.write(c)

r('gitmap/cmd/bookmarkrun.go', 'func runBookmarkRun(args []string) error {', 'func runBookmarkRun(args []string) *apperror.AppError {')
r('gitmap/cmd/bookmarkrun.go', 'func loadAndDispatchBookmark(name string) {', 'func loadAndDispatchBookmark(name string) *apperror.AppError {')
r('gitmap/cmd/bookmarkrun.go', 'fmt.Fprintf(os.Stderr, "  ? Could not query bookmarks: %v\\n", err)\n\t\tos.Exit(1)', 'return apperror.Wrap(err, "constants.ErrBookmarkQuery+", nil)')
r('gitmap/cmd/bookmarkrun.go', 'fmt.Fprintf(os.Stderr, "  ? Could not find bookmark \\"%s\\"\\n", name)\n\t\tos.Exit(1)', 'return apperror.Wrap(err, constants.ErrBookmarkNotFound, nil)')
r('gitmap/cmd/bookmarkrun.go', 'replayBookmark(bk.Command, bk.Args, bk.Flags)\n}', 'replayBookmark(bk.Command, bk.Args, bk.Flags)\n\treturn nil\n}')
r('gitmap/cmd/bookmarkrun.go', 'loadAndDispatchBookmark(name)\n\treturn nil', 'return loadAndDispatchBookmark(name)')

r('gitmap/cmd/bookmarksave.go', 'func runBookmarkSave(args []string) error {', 'func runBookmarkSave(args []string) *apperror.AppError {')
r('gitmap/cmd/bookmarksave.go', 'func saveBookmarkToDB(name, command, args, flags string) {', 'func saveBookmarkToDB(name, command, args, flags string) *apperror.AppError {')
r('gitmap/cmd/bookmarksave.go', 'checkBookmarkNotExists(db, name)', 'if err := checkBookmarkNotExists(db, name); err != nil { return err }')
r('gitmap/cmd/bookmarksave.go', 'fmt.Fprintf(os.Stderr, "  ? Could not query bookmarks: %v\\n", err)\n\t\tos.Exit(1)', 'return apperror.Wrap(err, constants.ErrBookmarkSave, nil)')
r('gitmap/cmd/bookmarksave.go', 'fmt.Fprintf(os.Stderr, "  ? Could not save bookmark: %v\\n", err)\n\t\tos.Exit(1)', 'return apperror.Wrap(err, constants.ErrBookmarkSave, nil)')
r('gitmap/cmd/bookmarksave.go', 'fmt.Printf(constants.MsgBookmarkSaved, name, command, args, flags)\n}', 'fmt.Printf(constants.MsgBookmarkSaved, name, command, args, flags)\n\treturn nil\n}')
r('gitmap/cmd/bookmarksave.go', '}, name string) {', '}, name string) *apperror.AppError {')
r('gitmap/cmd/bookmarksave.go', 'if err != nil {\n\t\treturn\n\t}', 'if err != nil {\n\t\treturn nil\n\t}')
r('gitmap/cmd/bookmarksave.go', 'fmt.Fprintf(os.Stderr, "  ? Bookmark \\"%s\\" already exists. Please choose a different name.\\n", name)\n\t\tos.Exit(1)\n}', 'return apperror.New(constants.ErrBookmarkExists, "E9000", nil)\n}')
r('gitmap/cmd/bookmarksave.go', 'saveBookmarkToDB(name, cmdStr, argStr, flagStr)\n\treturn nil', 'return saveBookmarkToDB(name, cmdStr, argStr, flagStr)')

r('gitmap/cmd/cdops.go', 'func loadRecentRepos() []model.ScanRecord {', 'func loadRecentRepos() ([]model.ScanRecord, *apperror.AppError) {')
r('gitmap/cmd/cdops.go', 'fmt.Fprintf(os.Stderr, "Error: %v\\n", err)\n\t\tos.Exit(1)', 'return nil, apperror.Wrap(err, "constants.ErrListDBFailed", nil)')
r('gitmap/cmd/cdops.go', 'return repos\n}', 'return repos, nil\n}')
r('gitmap/cmd/cdops.go', 'repos := loadRecentRepos()', 'repos, err := loadRecentRepos()\n\tif err != nil { return err }')
r('gitmap/cmd/cdops.go', 'func lookupProjectDir(id string) string {', 'func lookupProjectDir(id string) (string, *apperror.AppError) {')
r('gitmap/cmd/cdops.go', 'fmt.Fprint(os.Stderr, constants.ErrCdBare)\n\t\tos.Exit(1)', 'return "", apperror.New("fatal error", "E9000", nil)')
r('gitmap/cmd/cdops.go', 'fmt.Fprintf(os.Stderr, constants.ErrAsResolveFmt, id, err)\n\t\tos.Exit(1)', 'return "", apperror.Wrap(err, constants.ErrListDBFailed, nil)')
r('gitmap/cmd/cdops.go', 'return p.Path\n}', 'return p.Path, nil\n}')
r('gitmap/cmd/cdops.go', 'dir := lookupProjectDir(term)', 'dir, err := lookupProjectDir(term)\n\tif err != nil { return err }')

r('gitmap/cmd/cg.go', 'func parseCgFlags(args []string) cgOptions {', 'func parseCgFlags(args []string) (cgOptions, *apperror.AppError) {')
r('gitmap/cmd/cg.go', 'fmt.Fprintf(os.Stderr, "  ? Could not parse flags: %v\\n", err)\n\t\tos.Exit(1)', 'return cgOptions{}, apperror.Wrap(err, "parse flags", nil)')
r('gitmap/cmd/cg.go', '\tall:    *all,\n\t}\n}', '\tall:    *all,\n\t}, nil\n}')
r('gitmap/cmd/cg.go', 'opts := parseCgFlags(args)', 'opts, err := parseCgFlags(args)\n\tif err != nil { return err }')
