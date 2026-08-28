import re

def r(p, pairs):
    with open(p, 'r', encoding='utf8') as f: c = f.read()
    for o, n in pairs:
        c = c.replace(o, n)
    with open(p, 'w', encoding='utf8') as f: f.write(c)

r('gitmap/cmd/multigroupops.go', [('return apperror.Wrap(err, constants.ErrMGReadList, nil)', 'panic(err)')])
r('gitmap/cmd/pendingclear.go', [
    ('return apperror.New("fatal error", "E9000", nil)', 'panic("fatal")'),
    ('return apperror.New(constants.WarnPendingDBOpen, "E9000", nil)', 'return nil'),
    ('return apperror.New(constants.ErrPendingTaskQuery, "E9000", nil)', 'return nil'),
    ('return apperror.Wrap(err, constants.ErrPendingBatchClear, nil)', 'panic(err)')
])
r('gitmap/cmd/probe.go', [('return apperror.New("fatal error", "E9000", nil)', 'return nil')])
r('gitmap/cmd/profileutil.go', [
    ('return apperror.Wrap(err, constants.ErrListDBFailed, nil)', 'panic(err)'),
    ('return apperror.New(constants.ErrProfileNotFound, "E9000", nil)', 'panic("not found")')
])
r('gitmap/cmd/projectrepos.go', [('return apperror.New("fatal error", "E9000", nil)', 'panic("fatal")')])
