import os
import codecs

def fix_trailing_newline(filepath):
    try:
        with open(filepath, 'rb') as f:
            raw = f.read()
            
        if not raw:
            return False
            
        has_bom = raw.startswith(codecs.BOM_UTF8)
        if has_bom:
            raw = raw[len(codecs.BOM_UTF8):]
            
        try:
            text = raw.decode('utf-8')
        except UnicodeDecodeError:
            text = raw.decode('latin-1')
            
        # Ensure CRLF is stripped
        text = text.replace('\r\n', '\n').replace('\r', '\n')
        
        modified = False
        # Strip all trailing whitespace and newlines, then add exactly one
        original_text = text
        text = text.rstrip(' \t\n\r') + '\n'
        
        if original_text != text or has_bom:
            with open(filepath, 'w', encoding='utf-8', newline='\n') as f:
                f.write(text)
            return True
        return False
    except Exception as e:
        print(f"Error on {filepath}: {e}")
        return False

count = 0
extensions = ('.md', '.txt', '.go', '.ts', '.js', '.mjs', '.cjs', '.jsx', '.cs', '.vb', '.rs', '.json', '.yml', '.yaml', '.sh', '.ps1')

for root, dirs, files in os.walk('.'):
    # Exclude node_modules, .git, etc
    dirs[:] = [d for d in dirs if d not in ['.git', 'node_modules', 'dist', 'build']]
    for file in files:
        if file.endswith(extensions):
            filepath = os.path.join(root, file)
            if fix_trailing_newline(filepath):
                count += 1

print(f"Fixed trailing newlines on {count} files.")
