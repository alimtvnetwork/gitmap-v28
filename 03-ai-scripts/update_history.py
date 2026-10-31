import sys

def main():
    file_path = "src/routes/projects/$projectId/ai-testing-history.tsx"
    with open(file_path, "r", encoding="utf-8") as f:
        content = f.read()

    # Replaces for Industrial Look
    content = content.replace(
        'className="flex min-w-0 flex-1 flex-col overflow-auto"',
        'className="flex min-w-0 flex-1 flex-col overflow-auto bg-[#0b0c10] text-ca-ink font-mono antialiased"'
    )
    
    content = content.replace(
        'className="p-hmi-6"',
        'className="p-6"'
    )
    
    content = content.replace(
        'className="mb-hmi-5 flex flex-wrap items-end justify-between gap-hmi-3"',
        'className="mb-6 flex flex-wrap items-end justify-between gap-3 border-b border-[#22252a] pb-4"'
    )
    
    content = content.replace(
        'className="text-hmi-caption uppercase tracking-wide text-ca-ink-muted"',
        'className="text-[10px] font-bold uppercase tracking-widest text-ca-ink-muted"'
    )
    
    content = content.replace(
        'className="mt-hmi-1 font-display text-hmi-title font-extrabold uppercase tracking-wide text-ca-ink"',
        'className="mt-1 text-sm font-bold uppercase tracking-widest text-[#f5a623]"'
    )
    
    content = content.replace(
        'className="mt-hmi-1 text-hmi-body text-ca-ink-muted"',
        'className="mt-1 text-xs text-ca-ink-muted"'
    )
    
    content = content.replace(
        'className="rounded-sm border border-ca-border bg-ca-panel-2 px-hmi-3 py-hmi-2 text-hmi-body text-ca-ink hover:border-ca-select"',
        'className="inline-flex items-center gap-2 border border-[#444] bg-[#111318] px-4 py-2 text-[10px] font-bold uppercase tracking-wider text-ca-ink transition hover:border-[#f5a623] hover:text-[#f5a623]"'
    )
    
    content = content.replace(
        'className="mb-hmi-4 flex flex-wrap items-center gap-hmi-3 rounded-lg border border-ca-border bg-ca-panel p-hmi-3"',
        'className="mb-4 flex flex-wrap items-center gap-3 border border-[#22252a] bg-[#111318] p-3"'
    )
    
    content = content.replace(
        'className="flex items-center gap-hmi-2 text-hmi-body text-ca-ink"',
        'className="flex items-center gap-2 text-[10px] font-bold uppercase tracking-widest text-ca-ink"'
    )
    
    content = content.replace(
        'className="rounded-sm border border-ca-border bg-ca-panel-2 px-hmi-3 py-hmi-1 text-hmi-body text-ca-ink focus:border-ca-select focus:outline-none"',
        'className="border border-[#444] bg-[#0b0c10] px-3 py-1.5 text-xs text-ca-ink focus:border-[#f5a623] focus:outline-none"'
    )
    
    content = content.replace(
        'className="text-hmi-caption text-ca-ink-muted"',
        'className="text-[10px] uppercase text-ca-ink-muted"'
    )
    
    content = content.replace(
        'className="ml-auto rounded-sm border border-ca-border bg-ca-panel-2 px-hmi-3 py-hmi-1 text-hmi-caption text-ca-ink-muted hover:border-ca-ng hover:text-ca-ng"',
        'className="ml-auto inline-flex items-center gap-2 border border-[#444] bg-[#0b0c10] px-3 py-1.5 text-[10px] font-bold uppercase tracking-wider text-ca-ink-muted hover:border-red-500 hover:text-red-500"'
    )
    
    content = content.replace(
        'className="rounded-lg border border-dashed border-ca-border bg-ca-panel p-hmi-6 text-center text-hmi-body text-ca-ink-muted"',
        'className="border border-dashed border-[#444] bg-[#111318] p-8 text-center text-xs text-ca-ink-muted font-mono"'
    )
    
    content = content.replace(
        'className="grid grid-cols-1 gap-hmi-4 xl:grid-cols-[minmax(0,1fr)_minmax(0,1fr)]"',
        'className="grid grid-cols-1 gap-6 xl:grid-cols-[minmax(0,1fr)_minmax(0,1fr)]"'
    )
    
    content = content.replace(
        'className="overflow-hidden rounded-lg border border-ca-border bg-ca-panel"',
        'className="flex flex-col border border-[#22252a] bg-[#111318]"'
    )
    
    content = content.replace(
        'className="border-b border-ca-border px-hmi-3 py-hmi-2 text-hmi-caption uppercase tracking-wide text-ca-ink-muted"',
        'className="border-b border-[#22252a] bg-[#15171e] px-4 py-2 text-[10px] font-bold uppercase tracking-widest text-ca-ink-muted"'
    )
    
    content = content.replace(
        'className="divide-y divide-ca-border"',
        'className="divide-y divide-[#22252a]"'
    )
    
    content = content.replace(
        'className="rounded-lg border border-ca-border bg-ca-panel p-hmi-3"',
        'className="border border-[#22252a] bg-[#111318] p-4"'
    )
    
    content = content.replace(
        'className="font-display text-hmi-header font-extrabold uppercase tracking-wide text-ca-ink"',
        'className="mb-2 text-xs font-bold uppercase tracking-widest text-ca-ink"'
    )
    
    content = content.replace(
        'className="mt-hmi-2 text-hmi-body text-ca-ink-muted"',
        'className="text-xs text-ca-ink-muted"'
    )
    
    content = content.replace(
        'className="flex flex-wrap items-center gap-hmi-3 p-hmi-3"',
        'className="flex flex-wrap items-center gap-4 p-4 hover:bg-[#1a1c23] transition-colors"'
    )
    
    content = content.replace(
        'className="truncate text-hmi-body font-semibold text-ca-ink"',
        'className="truncate text-xs font-bold text-ca-ink"'
    )
    
    content = content.replace(
        'className="rounded-sm border border-ca-border bg-ca-panel-2 px-hmi-3 py-hmi-1 text-hmi-caption text-ca-ink hover:border-ca-select"',
        'className="inline-flex items-center gap-2 border border-[#444] bg-[#0b0c10] px-3 py-1.5 text-[10px] font-bold uppercase tracking-wider text-ca-ink transition hover:border-[#f5a623] hover:text-[#f5a623]"'
    )
    
    content = content.replace(
        'className="mt-hmi-3 flex flex-col gap-hmi-3"',
        'className="mt-4 flex flex-col gap-4"'
    )
    
    content = content.replace(
        'className="grid grid-cols-2 gap-hmi-3"',
        'className="grid grid-cols-2 gap-4"'
    )
    
    content = content.replace(
        'className="rounded-sm border border-ca-border bg-ca-panel-2 p-hmi-3"',
        'className="border border-[#22252a] bg-[#0b0c10] p-4"'
    )

    with open(file_path, "w", encoding="utf-8") as f:
        f.write(content)

if __name__ == "__main__":
    main()
