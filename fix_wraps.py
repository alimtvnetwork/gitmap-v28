import os, re

def process_file(filepath):
    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()

    original = content

    # If it's a Wrap wrapping a string instead of err, make it a New
    def repl_wrap(m):
        var, const_name = m.groups()
        if var in ['err', 'findErr', 'readErr', 'mergeErr', 'guardErr', 'pathErr']:
            return f'apperror.Wrap({var}, constants.{const_name}, nil)'
        # It's wrapping a string
        return f'apperror.New(constants.{const_name}, "E9000", nil)'

    content = re.sub(r'apperror\.Wrap\(([a-zA-Z0-9_]+),\s*"constants\.([a-zA-Z0-9_]+)",\s*nil\)', repl_wrap, content)

    # Any remaining apperror.Wrap with "constants.X" where err is err
    content = re.sub(r'apperror\.Wrap\(err,\s*"constants\.([a-zA-Z0-9_]+)",\s*nil\)', r'apperror.Wrap(err, constants.\1, nil)', content)
    
    if original != content:
        with open(filepath, 'w', encoding='utf-8') as f:
            f.write(content)

for root, _, files in os.walk('gitmap/cmd'):
    for f in files:
        if f.endswith('.go'):
            process_file(os.path.join(root, f))
