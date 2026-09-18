import re

def fix_colors(filepath):
    with open(filepath, 'r') as f:
        content = f.read()

    # Find things like `--config` or `--verbose` and wrap them in constants.ColorCyan + "..." + constants.ColorReset
    # Only within strings that look like help lines.
    def repl(m):
        flag = m.group(1)
        # Avoid double-wrapping if already wrapped
        if 'ColorCyan' in m.group(0):
            return m.group(0)
        return f'"{flag}"' # Wait, in Go we need `" + constants.ColorCyan + "`" + flag + "`" + constants.ColorReset + "`"

    # Actually it's easier to just do regex replacement in the file text.
    # We want to replace `--[a-z0-9-]+` inside help string literals.
    lines = content.split('\n')
    new_lines = []
    for line in lines:
        if 'Help' in line and '=' in line and '"' in line and '--' in line:
            # simple replacement for the flag part
            # match flag e.g. --config
            line = re.sub(r'(--[a-zA-Z0-9-]+)', r'" + constants.ColorCyan + "\1" + constants.ColorReset + "', line)
            # fix potential double quotes
            line = line.replace('"" + ', '').replace(' + ""', '')
        new_lines.append(line)

    with open(filepath, 'w') as f:
        f.write('\n'.join(new_lines))

fix_colors(r'gitmap\constants\constants_helpsections.go')
fix_colors(r'gitmap\constants\constants_cli.go')
fix_colors(r'gitmap\constants\constants_fixrepohelp.go')
fix_colors(r'gitmap\constants\constants_clonefixrepo.go')
