import subprocess
import re
def run_build():
    return subprocess.run(['go', 'build', './...'], capture_output=True, text=True).stderr

def fix_errors():
    for _ in range(50):
        stderr = run_build()
        if not stderr:
            print('Build succeeded!')
            break
        lines = stderr.splitlines()
        fixes = 0
        for line in lines:
            if 'declared and not used' in line:
                m = re.match(r'cmd[\\/]([a-zA-Z0-9_]+\.go):(\d+):\d+: (.*) declared and not used', line)
                if m:
                    filename, linenum = 'cmd/' + m.group(1), int(m.group(2))
                    with open(filename, 'r', encoding='utf-8') as f: content = f.readlines()
                    content[linenum-1] = '// ' + content[linenum-1]
                    with open(filename, 'w', encoding='utf-8') as f: f.writelines(content)
                    fixes += 1
                continue
            m1 = re.match(r'cmd[\\/]([a-zA-Z0-9_]+\.go):(\d+):\d+: undefined: err', line)
            if m1:
                filename, linenum = 'cmd/' + m1.group(1), int(m1.group(2))
                with open(filename, 'r', encoding='utf-8') as f: content = f.readlines()
                content[linenum-1] = content[linenum-1].replace('return err', 'return fmt.Errorf("command failed")')
                with open(filename, 'w', encoding='utf-8') as f: f.writelines(content)
                subprocess.run(['goimports', '-w', filename])
                fixes += 1
                continue
            m2 = re.match(r'cmd[\\/]([a-zA-Z0-9_]+\.go):(\d+):\d+: not enough return values', line)
            if m2:
                filename, linenum = 'cmd/' + m2.group(1), int(m2.group(2))
                with open(filename, 'r', encoding='utf-8') as f: content = f.readlines()
                content[linenum-1] = content[linenum-1].replace('return\n', 'return nil\n')
                with open(filename, 'w', encoding='utf-8') as f: f.writelines(content)
                fixes += 1
                continue
            m4 = re.match(r'cmd[\\/]([a-zA-Z0-9_]+\.go):(\d+):\d+: too many errors', line)
            if m4:
                continue

        if fixes == 0:
            print('Could not auto-fix. Exiting.')
            print(stderr)
            break
fix_errors()
