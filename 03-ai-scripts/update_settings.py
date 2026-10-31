import sys
import glob
import os

def process_file(file_path):
    with open(file_path, "r", encoding="utf-8") as f:
        content = f.read()

    # Replaces for Industrial Look
    content = content.replace(
        'className="flex-1 overflow-auto"',
        'className="flex-1 overflow-auto bg-[#0b0c10] text-ca-ink font-mono antialiased"'
    )
    
    content = content.replace(
        'className="mx-auto w-full max-w-6xl p-hmi-6"',
        'className="mx-auto w-full max-w-6xl p-6"'
    )
    
    content = content.replace(
        'className="mb-hmi-4 flex flex-wrap items-center justify-between gap-hmi-3 border-b border-ca-border pb-hmi-3"',
        'className="mb-6 flex flex-wrap items-center justify-between gap-3 border-b border-[#22252a] pb-4"'
    )
    
    content = content.replace(
        'className="font-display text-hmi-header font-extrabold uppercase tracking-wide text-ca-ink"',
        'className="text-sm font-bold uppercase tracking-widest text-[#f5a623]"'
    )
    
    content = content.replace(
        'className="mt-hmi-1 text-hmi-caption text-ca-ink-muted"',
        'className="mt-1 text-xs text-ca-ink-muted"'
    )
    
    content = content.replace(
        'className="w-full min-h-9 rounded-md bg-ca-panel-2 border border-ca-border pl-8 pr-8 py-hmi-2 text-hmi-body text-ca-ink hmi-tabular placeholder:text-ca-ink-muted focus:border-ca-select focus:outline-none"',
        'className="w-full min-h-9 border border-[#333] bg-[#0b0c10] pl-8 pr-8 py-2 text-xs text-ca-ink placeholder:text-[#555] focus:border-[#f5a623] focus:outline-none"'
    )
    
    content = content.replace(
        'className="grid grid-cols-1 gap-hmi-5 lg:grid-cols-[minmax(200px,220px)_minmax(0,1fr)]"',
        'className="grid grid-cols-1 gap-6 lg:grid-cols-[minmax(200px,220px)_minmax(0,1fr)]"'
    )
    
    content = content.replace(
        'className="group flex items-center gap-hmi-3 rounded-md border border-ca-border bg-ca-panel px-hmi-3 py-hmi-3 text-hmi-body text-ca-ink transition hover:-translate-y-px hover:border-ca-select hover:shadow-[0_10px_30px_-14px_color-mix(in_oklab,var(--color-ca-select)_60%,transparent)]"',
        'className="group flex items-center gap-3 border border-[#22252a] bg-[#111318] px-4 py-3 text-xs text-ca-ink transition hover:border-[#f5a623]"'
    )
    
    content = content.replace(
        'className="grid h-9 w-9 shrink-0 place-items-center rounded-md border border-ca-border bg-ca-panel-2 text-ca-select transition group-hover:bg-ca-select/10"',
        'className="grid h-8 w-8 shrink-0 place-items-center bg-[#0b0c10] text-[#f5a623] transition group-hover:bg-[#1a1c23]"'
    )
    
    content = content.replace(
        'className="mt-hmi-3 text-hmi-caption text-ca-ink-muted hmi-tabular"',
        'className="mt-3 text-[10px] text-ca-ink-muted hmi-tabular"'
    )
    
    content = content.replace(
        'className="mt-hmi-1 block w-full min-h-10 rounded-md bg-ca-panel-2 border border-ca-border px-hmi-3 py-hmi-2 text-hmi-body text-ca-ink hmi-tabular focus:border-ca-select focus:outline-none"',
        'className="mt-1 block w-full min-h-10 border border-[#333] bg-[#0b0c10] px-3 py-2 text-xs text-ca-ink focus:border-[#f5a623] focus:outline-none"'
    )

    content = content.replace(
        'className="mt-hmi-6 flex justify-end gap-hmi-3"',
        'className="mt-6 flex justify-end gap-3"'
    )

    with open(file_path, "w", encoding="utf-8") as f:
        f.write(content)

def main():
    for f in glob.glob("src/routes/settings/*.tsx"):
        process_file(f)

if __name__ == "__main__":
    main()
