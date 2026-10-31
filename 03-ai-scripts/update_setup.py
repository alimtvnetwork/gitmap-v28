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
        'className="bg-ca-panel"',
        'className="bg-[#111318]"'
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
    for f in glob.glob("src/routes/setup/*.tsx"):
        process_file(f)

if __name__ == "__main__":
    main()
