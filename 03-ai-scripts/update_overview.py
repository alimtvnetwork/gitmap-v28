import os
file_path = "src/routes/projects/$projectId/index.tsx"

with open(file_path, "r", encoding="utf-8") as f:
    content = f.read()

replacements = {
    'className="flex flex-1 flex-col gap-hmi-6 overflow-auto p-hmi-6"': 'className="flex flex-1 flex-col gap-6 overflow-auto p-6 bg-[#0b0c10]"',
    'rounded-lg border border-ca-border bg-ca-panel p-hmi-4 shadow-hmi-panel': 'rounded-sm border border-[#22252a] bg-[#111318] p-4 shadow-2xl',
    'rounded-lg border border-ca-border bg-ca-panel p-hmi-4': 'rounded-sm border border-[#22252a] bg-[#111318] p-4',
    'font-display text-hmi-title font-extrabold uppercase tracking-wide text-ca-ink': 'text-2xl font-bold uppercase tracking-widest text-[#f5a623]',
    'text-hmi-body text-ca-ink-muted': 'text-xs font-mono text-ca-ink-muted',
    'text-hmi-caption text-ca-ink-muted': 'text-[10px] font-bold uppercase tracking-widest text-ca-ink-muted',
    'text-hmi-body text-ca-ink': 'text-xs text-ca-ink font-mono',
    'bg-ca-panel-2': 'bg-[#1a1c23]',
    'bg-ca-panel': 'bg-[#111318]',
    'bg-ca-bg': 'bg-[#0b0c10]',
    'gap-hmi-1': 'gap-1',
    'gap-hmi-2': 'gap-2',
    'gap-hmi-3': 'gap-3',
    'gap-hmi-4': 'gap-4',
    'gap-hmi-6': 'gap-6',
    'p-hmi-1': 'p-1',
    'p-hmi-2': 'p-2',
    'p-hmi-3': 'p-3',
    'p-hmi-4': 'p-4',
    'p-hmi-5': 'p-5',
    'p-hmi-6': 'p-6',
    'px-hmi-2': 'px-2',
    'px-hmi-3': 'px-3',
    'px-hmi-4': 'px-4',
    'py-hmi-1': 'py-1',
    'py-hmi-2': 'py-2',
    'mt-hmi-1': 'mt-1',
    'mt-hmi-2': 'mt-2',
    'mt-hmi-3': 'mt-3',
    'mt-hmi-4': 'mt-4',
    'font-display text-hmi-header font-extrabold uppercase tracking-wide text-ca-ink': 'text-sm font-bold uppercase tracking-widest text-ca-ink',
    'rounded-sm bg-ca-select px-hmi-4 py-hmi-2 text-hmi-body font-semibold text-ca-bg hover:brightness-110 focus-visible:outline focus-visible:outline-2 focus-visible:outline-ca-focus': 'rounded bg-[#f5a623] px-4 py-2 text-xs font-bold uppercase tracking-wider text-black transition-colors hover:brightness-110',
    'border-ca-border': 'border-[#22252a]',
    'rounded-md border border-ca-border bg-ca-panel-2 px-hmi-3 py-hmi-2 text-hmi-body text-ca-ink hover:border-ca-select focus-visible:outline focus-visible:outline-2 focus-visible:outline-ca-focus': 'flex items-center gap-2 rounded border border-[#333] bg-[#1a1c23] px-3 py-2 text-xs font-bold uppercase tracking-wider text-ca-ink hover:bg-[#22252a] hover:text-[#f5a623] transition-colors',
    'text-ca-select': 'text-[#f5a623]',
    'bg-ca-select': 'bg-[#f5a623]',
    'group flex items-start gap-3 rounded-sm border border-[#22252a] bg-[#111318] p-4 shadow-2xl transition hover:border-[#f5a623] focus-visible:outline focus-visible:outline-2 focus-visible:outline-ca-focus': 'group flex items-start gap-3 rounded-sm border border-[#22252a] bg-[#111318] p-4 shadow-2xl transition hover:border-[#f5a623] focus-visible:outline focus-visible:outline-2 focus-visible:outline-ca-focus'
}

for old, new in replacements.items():
    content = content.replace(old, new)

with open(file_path, "w", encoding="utf-8") as f:
    f.write(content)

