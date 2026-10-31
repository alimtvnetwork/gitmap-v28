import sys

def main():
    file_path = "src/routes/projects/$projectId/rulesets/$rulesetId.tsx"
    with open(file_path, "r", encoding="utf-8") as f:
        content = f.read()

    # Replaces for Industrial Look
    content = content.replace(
        'className="flex min-w-0 flex-1 flex-col overflow-auto p-hmi-4"',
        'className="flex min-w-0 flex-1 flex-col overflow-auto bg-[#0b0c10] text-ca-ink font-mono antialiased p-4"'
    )
    content = content.replace(
        'className="mb-hmi-3 flex flex-wrap items-center gap-hmi-1 rounded-md border border-ca-border/60 bg-ca-panel/50 p-hmi-1"',
        'className="mb-4 flex flex-wrap items-center gap-2 border-b border-[#333] bg-[#1a1c23] p-2"'
    )
    content = content.replace(
        'className="grid grid-cols-1 gap-hmi-4 lg:grid-cols-[minmax(0,1fr)_380px] xl:grid-cols-[minmax(0,1fr)_420px]"',
        'className="grid grid-cols-1 gap-4 lg:grid-cols-[minmax(0,1fr)_380px] xl:grid-cols-[minmax(0,1fr)_420px]"'
    )
    content = content.replace(
        'className="h-[74vh] min-h-[500px] w-full p-1.5 flex flex-col overflow-hidden bg-ca-panel"',
        'className="h-[74vh] min-h-[500px] w-full p-0 flex flex-col overflow-hidden border border-[#333] bg-[#000]"'
    )
    content = content.replace(
        'className="flex h-[74vh] min-h-[500px] flex-col overflow-hidden rounded-lg border border-ca-border"',
        'className="flex h-[74vh] min-h-[500px] flex-col overflow-hidden border border-[#333] bg-[#1a1c23]"'
    )
    
    # Buttons and tabs
    content = content.replace(
        'className="inline-flex items-center gap-1.5 rounded-sm border border-ca-border bg-ca-panel px-2.5 py-1 text-hmi-caption font-semibold text-ca-ink transition hover:border-ca-select hover:bg-ca-panel-2 focus-visible:outline focus-visible:outline-2 focus-visible:outline-ca-focus"',
        'className="inline-flex items-center gap-2 border border-[#444] bg-[#0b0c10] px-3 py-1.5 text-[10px] font-bold uppercase tracking-wider text-ca-ink transition hover:border-ca-primary hover:text-ca-primary focus-visible:outline-none"'
    )
    
    content = content.replace(
        'className="inline-flex items-center gap-hmi-2 rounded-sm border border-ca-border bg-ca-panel px-hmi-2 py-hmi-1 text-hmi-caption font-semibold text-ca-ink transition hover:border-ca-select hover:bg-ca-panel-2 focus-visible:outline focus-visible:outline-2 focus-visible:outline-ca-focus"',
        'className="inline-flex items-center gap-2 border border-[#444] bg-[#0b0c10] px-3 py-1.5 text-[10px] font-bold uppercase tracking-wider text-ca-ink transition hover:border-ca-primary hover:text-ca-primary focus-visible:outline-none"'
    )
    
    content = content.replace(
        'className="inline-flex items-center gap-hmi-2 rounded-sm bg-ca-select px-hmi-3 py-hmi-1 text-hmi-caption font-semibold text-ca-bg transition hover:brightness-110 focus-visible:outline focus-visible:outline-2 focus-visible:outline-ca-focus disabled:opacity-50"',
        'className="inline-flex items-center gap-2 bg-ca-primary px-4 py-1.5 text-[10px] font-bold uppercase tracking-wider text-black transition hover:brightness-110 focus-visible:outline-none disabled:opacity-30"'
    )

    content = content.replace(
        'className="inline-flex items-center gap-hmi-2 rounded-sm border border-cyan-500/60 bg-cyan-950/30 px-hmi-2 py-hmi-1 text-hmi-caption font-semibold text-cyan-200 transition hover:bg-cyan-900/40 hover:border-cyan-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-ca-focus"',
        'className="inline-flex items-center gap-2 border border-cyan-500 bg-cyan-950/50 px-3 py-1.5 text-[10px] font-bold uppercase tracking-wider text-cyan-400 transition hover:bg-cyan-900 focus-visible:outline-none"'
    )

    with open(file_path, "w", encoding="utf-8") as f:
        f.write(content)

if __name__ == "__main__":
    main()
