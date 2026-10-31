import sys

def main():
    file_path = "src/routes/projects/$projectId/trial-run.$runId.tsx"
    with open(file_path, "r", encoding="utf-8") as f:
        content = f.read()

    # Replaces for Industrial Look
    content = content.replace(
        'className="flex min-w-0 flex-1 flex-col overflow-auto p-hmi-6"',
        'className="flex min-w-0 flex-1 flex-col overflow-auto bg-[#0b0c10] text-ca-ink font-mono antialiased p-6"'
    )
    
    content = content.replace(
        'className="mt-hmi-2 font-display text-hmi-title font-extrabold uppercase tracking-wide text-ca-ink"',
        'className="mt-2 text-sm font-bold uppercase tracking-widest text-ca-ink"'
    )
    
    content = content.replace(
        'className="mt-hmi-1 text-hmi-body text-ca-ink-muted"',
        'className="mt-1 text-xs text-ca-ink-muted"'
    )
    
    content = content.replace(
        'className="mt-hmi-1 font-mono text-hmi-caption text-ca-ink-muted"',
        'className="mt-1 text-[10px] text-[#555]"'
    )
    
    content = content.replace(
        'aria-label="Run summary"\n          className={`mb-hmi-5 flex items-center gap-hmi-4 rounded-lg border p-hmi-4 ${',
        'aria-label="Run summary"\n          className={`mb-6 flex items-center gap-4 border p-4 ${'
    )
    
    content = content.replace(
        'className={`font-display text-hmi-header font-extrabold uppercase tracking-wide ${',
        'className={`text-2xl font-bold uppercase tracking-wider ${'
    )
    
    content = content.replace(
        'className="grid grid-cols-1 gap-hmi-5 lg:grid-cols-[minmax(0,1fr)_360px]"',
        'className="grid grid-cols-1 gap-6 lg:grid-cols-[minmax(0,1fr)_400px]"'
    )
    
    content = content.replace(
        'className="rounded-lg border border-ca-border bg-ca-panel"',
        'className="border border-[#22252a] bg-[#111318]"'
    )
    
    content = content.replace(
        'className="border-b border-ca-border px-hmi-3 py-hmi-2 text-hmi-caption uppercase tracking-wide text-ca-ink-muted"',
        'className="border-b border-[#22252a] bg-[#15171e] px-4 py-2 text-[10px] font-bold uppercase tracking-widest text-ca-ink-muted"'
    )
    
    content = content.replace(
        'className="border-b border-ca-border text-hmi-caption uppercase tracking-wide text-ca-ink-muted"',
        'className="border-b border-[#22252a] text-[10px] font-bold uppercase tracking-widest text-ca-ink-muted"'
    )
    
    content = content.replace(
        'className="px-hmi-3 py-hmi-2 text-left"',
        'className="px-4 py-2 text-left"'
    )
    content = content.replace(
        'className="px-hmi-3 py-hmi-2 text-right"',
        'className="px-4 py-2 text-right"'
    )
    
    content = content.replace(
        'className="border-b border-ca-border/60 last:border-0"',
        'className="border-b border-[#22252a] last:border-0"'
    )
    
    content = content.replace(
        'className="px-hmi-3 py-hmi-2 text-ca-ink"',
        'className="px-4 py-3 text-xs text-ca-ink font-bold"'
    )
    
    content = content.replace(
        'className="px-hmi-3 py-hmi-2 text-ca-ink-muted"',
        'className="px-4 py-3 text-xs text-ca-ink-muted"'
    )
    
    content = content.replace(
        'className="px-hmi-3 py-hmi-2 text-right font-mono text-ca-ink"',
        'className="px-4 py-3 text-xs text-right font-mono text-ca-ink"'
    )
    
    content = content.replace(
        'className="flex min-h-0 flex-col rounded-lg border border-ca-border bg-ca-panel"',
        'className="flex min-h-0 flex-col border border-[#22252a] bg-[#000]"'
    )
    
    content = content.replace(
        'className="flex items-center justify-between border-b border-ca-border px-hmi-3 py-hmi-2 text-hmi-caption uppercase tracking-wide text-ca-ink-muted"',
        'className="flex items-center justify-between border-b border-[#22252a] bg-[#111318] px-4 py-2 text-[10px] font-bold uppercase tracking-widest text-ca-ink-muted"'
    )
    
    content = content.replace(
        'className="max-h-[60vh] overflow-auto whitespace-pre-wrap p-hmi-3 font-mono text-hmi-caption text-ca-ink"',
        'className="max-h-[60vh] overflow-auto whitespace-pre-wrap p-4 font-mono text-xs text-[#0f0]"'
    )

    with open(file_path, "w", encoding="utf-8") as f:
        f.write(content)

if __name__ == "__main__":
    main()
