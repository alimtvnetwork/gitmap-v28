import sys
import glob
import os

def process_file(file_path):
    with open(file_path, "r", encoding="utf-8") as f:
        content = f.read()

    # Replaces for Industrial Look
    content = content.replace(
        'className="p-hmi-6"',
        'className="p-6 bg-[#0b0c10] text-ca-ink font-mono"'
    )
    
    content = content.replace(
        'className="p-hmi-4"',
        'className="p-4 bg-[#0b0c10] text-ca-ink font-mono"'
    )
    
    content = content.replace(
        'bg-ca-panel"',
        'bg-[#111318]"'
    )
    
    content = content.replace(
        'bg-ca-panel ',
        'bg-[#111318] '
    )
    
    content = content.replace(
        'bg-ca-panel-2',
        'bg-[#1a1c23]'
    )
    
    content = content.replace(
        'border-ca-border',
        'border-[#22252a]'
    )
    
    content = content.replace(
        'rounded-lg',
        'rounded-none'
    )
    
    content = content.replace(
        'rounded-md',
        'rounded-sm'
    )

    content = content.replace(
        'text-hmi-title font-extrabold',
        'text-sm font-bold uppercase tracking-widest text-[#f5a623]'
    )
    
    content = content.replace(
        'text-hmi-header font-extrabold',
        'text-xs font-bold uppercase tracking-widest'
    )
    
    content = content.replace(
        'text-hmi-caption',
        'text-[10px] uppercase tracking-wider'
    )
    
    content = content.replace(
        'text-hmi-body',
        'text-xs'
    )

    with open(file_path, "w", encoding="utf-8") as f:
        f.write(content)

def main():
    files_to_process = [
        "src/routes/diagnostics.tsx",
        "src/routes/ops.tsx",
        "src/routes/run.tsx",
        "src/routes/analyze.tsx",
        "src/routes/results.tsx",
        "src/routes/cli.tsx"
    ]
    for pattern in ["src/routes/observability/*.tsx"]:
        files_to_process.extend(glob.glob(pattern))

    for f in files_to_process:
        if os.path.exists(f):
            process_file(f)

if __name__ == "__main__":
    main()
