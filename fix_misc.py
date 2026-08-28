import os
def fix_file(f, from_str, to_str):
    with open(f, 'r', encoding='utf-8') as file:
        c = file.read()
    c = c.replace(from_str, to_str)
    with open(f, 'w', encoding='utf-8') as file:
        file.write(c)

fix_file('gitmap/cmd/diff.go', 'apperror.Wrap(guardErr, "E9000", "guard-paths left+" vs "+right")', 'apperror.Wrap(guardErr, "guard-paths", nil)')

fix_file('gitmap/cmd/history.go', '"constants.ErrHistoryQuery+\\n"', 'constants.ErrHistoryQuery')
fix_file('gitmap/cmd/history.go', '"constants.ErrHistoryQuery+\n"', 'constants.ErrHistoryQuery')

# Wait, macro_cmd.go is corrupted. I should git checkout macro_cmd.go and fix manually or let fix_exits do it properly.
