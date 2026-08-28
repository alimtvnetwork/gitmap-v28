import re

def r(p, pairs):
    with open(p, 'r', encoding='utf8') as f: c = f.read()
    for o, n in pairs:
        c = c.replace(o, n)
    with open(p, 'w', encoding='utf8') as f: f.write(c)

r('gitmap/cmd/reinstall.go', [
    ('return apperror.New(constants.ErrNoInstallation, "E9000", nil)', 'return "", False'),
    ('return apperror.Wrap(err, constants.ErrStatFailed, nil)', 'return "", False')
])
r('gitmap/cmd/releasealias.go', [
    ('return apperror.New("fatal error", "E9000", nil)', 'panic("fatal")'),
    ('var empty string\n\tvar empty string\n\treturn empty\n\treturn empty', 'var empty string\n\treturn empty'),
    ('return apperror.Wrap(err, "✗ Set branch error", nil)', 'panic(err)')
])
r('gitmap/cmd/releasepull.go', [
    ('return apperror.New("fatal error", "E9000", nil)', 'return ""'),
    ('return apperror.Wrap(err, constants.ErrRPCwdFailedFmt, nil)', 'return ""'),
    ('return apperror.Wrap(err, constants.ErrRPMergeFailedFmt, nil)', 'panic(err)'),
    ('return apperror.New(constants.ErrRPPushFailedFmt, "E9000", nil)', 'panic("fatal")')
])
r('gitmap/cmd/releaserebase.go', [('return apperror.Wrap(err, "git push failed", nil)', 'panic(err)')])
r('gitmap/cmd/releaseself.go', [('return apperror.New(constants.ErrSelfCwdFailed, "E9000", nil)', 'panic("fatal")')])
