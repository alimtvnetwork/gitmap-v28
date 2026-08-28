import os
import re
import glob

def process_file(filepath):
    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()

    original = content

    # Replace cliexit.Fail(command, op, subject, err, code)
    # e.g., cliexit.Fail("commit-in", "parse flag", "limit", err, 1)
    def repl_cliexit(m):
        cmd, op, sub, err, code = m.groups()
        cmd = cmd.strip('"')
        op = op.strip('"')
        sub = sub.strip('"')
        msg = f"{op} {sub}".strip()
        if err == "nil":
            return f'return apperror.New("{msg}", "E9000", nil)'
        return f'return apperror.Wrap({err}, "{msg}", nil)'

    content = re.sub(r'cliexit\.Fail\(([^,]+),\s*([^,]+),\s*([^,]+),\s*([^,]+),\s*([^)]+)\)', repl_cliexit, content)

    # Replace pterm.Error.Println(...) \n os.Exit(1)
    def repl_pterm_println(m):
        msg = m.group(1).replace('\n', ' ').strip('"')
        return f'return apperror.New("{msg}", "E9000", nil)'
    content = re.sub(r'pterm\.Error\.Println\(([^)]+)\)\s*os\.Exit\(1\)', repl_pterm_println, content)
    content = re.sub(r'fmt\.Fprintln\(os\.Stderr,\s*([^)]+)\)\s*os\.Exit\(1\)', repl_pterm_println, content)
    
    # Replace pterm.Error.Printf("... %v", err) \n os.Exit(1)
    def repl_pterm_printf(m):
        msg, args = m.groups()
        msg = msg.strip('"').replace('%v', '').replace('%s', '').replace('\n', '').strip()
        return f'return apperror.Wrap({args.split(",")[0].strip()}, "{msg}", nil)'
    content = re.sub(r'pterm\.Error\.Printf\(([^,]+),\s*([^)]+)\)\s*os\.Exit\(1\)', repl_pterm_printf, content)
    content = re.sub(r'fmt\.Fprintf\(os\.Stderr,\s*([^,]+),\s*([^)]+)\)\s*os\.Exit\(1\)', repl_pterm_printf, content)

    # Replace remaining os.Exit(1) with return apperror.New("fatal error", "E9000", nil)
    content = re.sub(r'(?<!return )os\.Exit\(1\)', 'return apperror.New("fatal error", "E9000", nil)', content)

    # We need to make sure we don't end up with invalid wrap/new from the previous run
    content = re.sub(r'apperror\.Wrap\(([^,]+),\s*"E9000",\s*("[^"]+")\)', r'apperror.Wrap(\1, \2, nil)', content)
    content = re.sub(r'apperror\.New\("E9000",\s*("[^"]+")\)', r'apperror.New(\1, "E9000", nil)', content)

    if original != content:
        with open(filepath, 'w', encoding='utf-8') as f:
            f.write(content)
        print(f"Modified {filepath}")

for root, _, files in os.walk('gitmap/cmd'):
    for f in files:
        if f.endswith('.go'):
            process_file(os.path.join(root, f))
