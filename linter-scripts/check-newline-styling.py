#!/usr/bin/env python3
import os
import sys

def check_file(filepath):
    violations = []
    
    with open(filepath, 'r', encoding='utf-8') as f:
        try:
            content = f.read()
        except UnicodeDecodeError as exc:
            import sys as _sys; print(f"Error: {exc}", file=_sys.stderr)
            return violations # skip unreadable files
            
    lines = content.split('\n')
    
    empty_streak = 0
    for i, line in enumerate(lines):
        stripped = line.strip()
        
        # 1. No double empty lines (\n\n\n)
        if stripped == '':
            empty_streak += 1
            if empty_streak == 2:
                violations.append((i + 1, "No double empty lines (\\n\\n\\n) allowed"))
        else:
            empty_streak = 0
            
        # 2. No empty line at the start of a function/block
        if stripped.endswith('{'):
            if i + 1 < len(lines) and lines[i+1].strip() == '':
                violations.append((i + 2, "No empty line at the start of a function or block"))
                
        # 3. Blank line before return (for multi-line functions/blocks)
        if stripped.startswith('return ') or stripped == 'return':
            if i > 0:
                prev_line = lines[i-1].strip()
                # Skip if prev_line is empty, '{', '}', ':', or a comment
                if (prev_line != '' and 
                    not prev_line.endswith('{') and 
                    not prev_line.endswith('}') and 
                    not prev_line.endswith(':') and 
                    not prev_line.startswith('//') and 
                    not prev_line.startswith('/*') and 
                    not prev_line.startswith('*')):
                    violations.append((i + 1, "Blank line required before return"))
                    
        # 4. Blank line after } if followed by more code
        if stripped == '}':
            if i + 1 < len(lines):
                next_line = lines[i+1].strip()
                # If next line is not empty and not part of the same construct (else, catch, etc.)
                if (next_line != '' and 
                    not next_line.startswith(('}', 'else', 'catch', 'finally', ')', ']', ',', ';', '//', '/*', '</'))):
                    violations.append((i + 1, "Blank line required after '}' if followed by more code"))
                    
        # 5. Check for \n in Go files
        if filepath.endswith('.go'):
            if '"\\n"' in line:
                violations.append((i + 1, 'Use constants.NewLineUnix instead of "\\n"'))

    return violations

def get_target_dir():
    # Check if src/ exists in current directory
    src_dir = os.path.join(os.getcwd(), 'src')
    if os.path.exists(src_dir):
        return src_dir
        
    return os.getcwd()

def main():
    target_dir = get_target_dir()
    extensions = ('.go', '.ts', '.tsx', '.js')
    ignore_dirs = {'.git', 'node_modules', '.github', '.ci-out', 'dist', 'build', 'tmp'}
    
    total_violations = 0
    
    for root, dirs, files in os.walk(target_dir):
        dirs[:] = [d for d in dirs if d not in ignore_dirs]
        for file in files:
            if file.endswith(extensions):
                filepath = os.path.join(root, file)
                violations = check_file(filepath)
                if violations:
                    for line_num, msg in violations:
                        rel_path = os.path.relpath(filepath, os.getcwd())
                        print(f"{rel_path}:{line_num}: {msg}")
                        total_violations += 1
                        
    if total_violations > 0:
        sys.exit(1)
    else:
        sys.exit(0)

if __name__ == "__main__":
    main()
