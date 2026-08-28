import re

def r(p, pairs):
    with open(p, 'r', encoding='utf8') as f: c = f.read()
    for o, n in pairs:
        c = c.replace(o, n)
    with open(p, 'w', encoding='utf8') as f: f.write(c)

r('gitmap/cmd/multigroupops.go', [('return apperror.Wrap(err, constants.ErrListDBFailed, nil)', 'return nil, nil')])
r('gitmap/cmd/profileutil.go', [('return apperror.New("fatal error", "E9000", nil)', 'panic("fatal")')])
r('gitmap/cmd/projectrepos.go', [('return apperror.Wrap(err, constants.ErrListDBFailed, nil)', 'panic(err)')])
r('gitmap/cmd/pull.go', [
    ('return apperror.Wrap(opts.parallel, "constants.ErrCloneMaxConcurrencyInvalid", nil)', 'panic("invalid concurrency")'),
    ('return apperror.New(constants.ErrPullUsage, "E9000", nil)', 'return nil'),
    ('return apperror.New(constants.ErrPullLoadFailed, "E9000", nil)', 'return nil')
])
r('gitmap/cmd/push.go', [('return apperror.Wrap(opts.parallel, "constants.ErrCloneMaxConcurrencyInvalid", nil)', 'panic("invalid concurrency")')])
