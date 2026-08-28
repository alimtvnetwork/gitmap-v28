import re

def r(p, pairs):
    with open(p, 'r', encoding='utf8') as f: c = f.read()
    for o, n in pairs:
        c = c.replace(o, n)
    with open(p, 'w', encoding='utf8') as f: f.write(c)

r('gitmap/cmd/install.go', [('return apperror.New(constants.ErrInstallUnknownTool, "E9000", nil)', 'panic("unknown tool")')])
r('gitmap/cmd/install_custom.go', [
    ('return apperror.Wrap(tool, "Unknown custom tool:", nil)', 'return apperror.New("Unknown custom tool: " + tool, "E9000", nil)'),
    ('return apperror.Wrap(tool, "Failed to install :", nil)', 'return apperror.New("Failed to install: " + tool, "E9000", nil)')
])
r('gitmap/cmd/installcleancode.go', [('return apperror.Wrap(constants.DefaultCleanCodeURL, "constants.MsgCleanCodeNoPwsh", nil)', 'return apperror.New(constants.MsgCleanCodeNoPwsh, "E9000", nil)')])
r('gitmap/cmd/installtools.go', [('return apperror.Wrap(err, "constants.ErrInstallDownloadFailed", nil)', 'panic(err)')])
r('gitmap/cmd/latestbranch.go', [
    ('return apperror.Wrap(err, "constants.ErrLatestBranchFatal", nil)', 'panic(err)'),
    ('return apperror.Wrap(cfg.remote, "constants.ErrLatestBranchNoRefs", nil)', 'return nil'),
    ('return apperror.Wrap(cfg.filter, "constants.ErrLatestBranchNoMatch", nil)', 'return nil')
])
r('gitmap/cmd/latestbranchswitch.go', [
    ('return apperror.Wrap(err, "constants.ErrLatestBranchFatal", nil)', 'panic(err)'),
    ('return apperror.Wrap(opts.Remote, "constants.ErrLatestBranchNoRefs", nil)', 'panic(opts.Remote)')
])
