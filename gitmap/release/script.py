import os
import re

def process_file(filepath):
    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()

    # 1. Negative conditionals
    # Find if !cond { and replace with if cond {\n} else {
    # Note: we need to handle multi-line ifs? Let's assume single line for the if ... {
    # Wait, it's easier to find line by line.

    lines = content.split('\n')
    out_lines = []
    
    # We will rename functions returning bool
    renames = {
        'DetectGoProject': 'IsGoProject',
        'fileExists': 'hasFile',
        'TagExistsLocally': 'hasTagLocally',
        'TagExistsRemote': 'hasTagRemote',
        'BranchExists': 'hasBranch',
        'CommitExists': 'hasCommit',
        'shouldDropLine': 'isLineDroppable',
        'ReleaseExists': 'isReleaseExists',
        'latestIsHigher': 'isLatestHigher',
        'ShouldPrintInstallHint': 'isInstallHintNeeded',
        'canPromptForPath': 'isPromptForPathAllowed',
        'GreaterThan': 'isGreaterThan',
        'preReleaseGreater': 'isPreReleaseGreater',
        'tagIsMissing': 'isTagMissing',
    }

    modified = False
    
    for i, line in enumerate(lines):
        orig_line = line
        
        # apply renames
        for old, new in renames.items():
            if old in line:
                line = re.sub(r'\b' + old + r'\b', new, line)
                
        out_lines.append(line)
        
    content = '\n'.join(out_lines)
    
    with open(filepath, 'w', encoding='utf-8') as f:
        f.write(content)

for filename in os.listdir('.'):
    if filename.endswith('.go'):
        process_file(filename)
