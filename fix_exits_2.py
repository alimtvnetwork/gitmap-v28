import os
import re

def process_file(filepath):
    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()

    original = content

    # Replace cliexit.Fail
    def repl_cliexit(m):
        cmd, op, sub, err, code = m.groups()
        cmd = cmd.strip().strip('"')
        op = op.strip().strip('"')
        sub = sub.strip().strip('"')
        msg = f"{op} {sub}".strip()
        err = err.strip()
        if err == "nil":
            return f'return apperror.New("{msg}", "E9000", nil)'
        # If it already wraps apperror, strip it
        if "apperror.Wrap" in err:
            return f'return {err}'
        return f'return apperror.Wrap({err}, "{msg}", nil)'

    content = re.sub(r'cliexit\.Fail\(([^,]+),\s*([^,]+),\s*([^,]+),\s*([^,]+),\s*([^)]+)\)', repl_cliexit, content)

    # pterm.Error.Println / fmt.Fprintln
    def repl_pterm_println(m):
        arg = m.group(1).strip()
        # If it's a string literal, keep it as is, just pass it to New
        return f'return apperror.New({arg}, "E9000", nil)'
    
    content = re.sub(r'pterm\.Error\.Println\(([^)]+)\)\s*os\.Exit\(1\)', repl_pterm_println, content)
    content = re.sub(r'fmt\.Fprintln\(os\.Stderr,\s*([^)]+)\)\s*os\.Exit\(1\)', repl_pterm_println, content)
    
    # pterm.Error.Printf / fmt.Fprintf
    def repl_pterm_printf(m):
        fmt_str = m.group(1).strip()
        args = m.group(2).strip()
        # Usually it's something like err
        first_arg = args.split(',')[0].strip()
        
        # We need a clean op string. We will just use the fmt_str
        # but remove \n and %v
        clean_msg = fmt_str.replace('\\n', '').replace('%v', '').replace('%s', '').replace('"', '').strip()
        
        return f'return apperror.Wrap({first_arg}, "{clean_msg}", nil)'
        
    content = re.sub(r'pterm\.Error\.Printf\(([^,]+),\s*([^)]+)\)\s*os\.Exit\(1\)', repl_pterm_printf, content)
    content = re.sub(r'fmt\.Fprintf\(os\.Stderr,\s*([^,]+),\s*([^)]+)\)\s*os\.Exit\(1\)', repl_pterm_printf, content)

    content = re.sub(r'(?<!return )os\.Exit\(1\)', 'return apperror.New("fatal error", "E9000", nil)', content)

    if original != content:
        with open(filepath, 'w', encoding='utf-8') as f:
            f.write(content)
        print(f"Modified {filepath}")

for root, _, files in os.walk('gitmap/cmd'):
    for f in files:
        if f.endswith('.go'):
            process_file(os.path.join(root, f))
