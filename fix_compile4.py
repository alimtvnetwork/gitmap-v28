import os, re

def r(p, o, n):
    with open(p, 'r', encoding='utf8') as f: c = f.read()
    c = c.replace(o, n)
    with open(p, 'w', encoding='utf8') as f: f.write(c)

r('gitmap/cmd/as.go', 'return apperror.New("fatal error", "E9000", nil)\n\t}', 'return model.ScanRecord{}, apperror.New("fatal error", "E9000", nil)\n\t}')
r('gitmap/cmd/asops.go', 'return apperror.New("fatal error", "E9000", nil)\n\t}\n}', 'return apperror.New("fatal error", "E9000", nil)\n\t}\n\treturn nil\n}')
r('gitmap/cmd/bookmarkrun.go', 'func runBookmarkRun(args []string) {', 'func runBookmarkRun(args []string) *apperror.AppError {')
r('gitmap/cmd/bookmarkrun.go', 'func executeBookmark(name string, dryRun bool) {', 'func executeBookmark(name string, dryRun bool) *apperror.AppError {')
r('gitmap/cmd/bookmarkrun.go', 'return apperror.Wrap(findErr, constants.ErrBookmarkNotFound, nil)', 'return apperror.Wrap(err, constants.ErrBookmarkNotFound, nil)')
r('gitmap/cmd/bookmarksave.go', 'func runBookmarkSave(args []string) {', 'func runBookmarkSave(args []string) *apperror.AppError {')
r('gitmap/cmd/bookmarksave.go', 'func upsertBookmark(bookmark model.BookmarkRow) {', 'func upsertBookmark(bookmark model.BookmarkRow) *apperror.AppError {')
r('gitmap/cmd/bookmarksave.go', 'return apperror.Wrap(name, "constants.ErrBookmarkNotFound", nil)', 'return apperror.Wrap(err, constants.ErrBookmarkNotFound, nil)')
r('gitmap/cmd/bookmarksave.go', 'return apperror.Wrap(findErr, constants.ErrBookmarkNotFound, nil)', 'return apperror.Wrap(err, constants.ErrBookmarkNotFound, nil)')

r('gitmap/cmd/branch.go', 'func runBranch(args []string) error {', 'func runBranch(args []string) *apperror.AppError {')
r('gitmap/cmd/branch.go', 'return apperror.Wrap(sub, "constants.ErrUnknownCommand", nil)', 'return apperror.New(constants.ErrUnknownCommand, "E9000", nil)')

