import re
import sys

def main():
    file_path = "src/routes/projects.index.tsx"
    with open(file_path, "r", encoding="utf-8") as f:
        content = f.read()
    
    # We will replace `const mainContent = (` ... to the end of the file.
    
    start_str = "const mainContent = ("
    start_idx = content.find(start_str)
    if start_idx == -1:
        print("Could not find mainContent")
        sys.exit(1)
    
    prefix = content[:start_idx]
    
    new_main_content = """const mainContent = (
    <div className="flex min-w-0 flex-1 flex-col overflow-auto bg-[#0b0c10] text-ca-ink font-mono antialiased p-6">
      <div className="mx-auto w-full max-w-7xl">
        <header className="mb-6 flex flex-wrap items-center justify-between gap-4 border-b border-[#333] pb-4">
          <div className="min-w-0 flex items-baseline gap-4">
            <h1 className="text-xl font-bold uppercase tracking-wider text-ca-ink flex items-center gap-2">
              <Cpu className="text-ca-primary" size={20} />
              Project Database
            </h1>
            <span className="text-xs tracking-widest tabular-nums text-ca-ink-muted">
              {hydrated ? `TOTAL: ${list.length}` : "LOADING..."}
            </span>
          </div>
          <div className="flex flex-wrap items-center gap-3">
            <label
              className="inline-flex cursor-pointer items-center gap-2 border border-[#333] bg-[#1a1c23] px-4 py-2 text-xs font-bold uppercase tracking-wider text-ca-ink hover:border-ca-primary transition-colors focus-within:border-ca-primary"
              aria-label="Import project from JSON or YAML"
            >
              <Upload aria-hidden size={14} />
              Import
              <input
                type="file"
                accept=".json,.yaml,.yml,application/json,application/yaml"
                className="sr-only"
                onChange={(e) => {
                  const f = e.currentTarget.files?.[0];
                  if (f) void handleImportFile(f);
                  e.currentTarget.value = "";
                }}
              />
            </label>
            <button
              type="button"
              onClick={openDialog}
              className="inline-flex items-center gap-2 bg-ca-primary px-4 py-2 text-xs font-bold uppercase tracking-wider text-black hover:bg-ca-primary-hover transition-colors focus-visible:outline-none"
            >
              <Plus aria-hidden size={16} />
              New Project
            </button>
          </div>
        </header>

        {importErr ? (
          <div className="mb-4 rounded-sm border border-red-500/50 bg-red-500/10 p-3 text-xs text-red-500">
            ERR: {importErr}
          </div>
        ) : null}

        <div className="mb-6 flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between bg-[#1a1c23] p-3 border border-[#333]">
          <div className="relative flex-1 sm:max-w-md flex items-center">
            <Search
              aria-hidden
              size={14}
              className="absolute left-3 text-ca-ink-muted pointer-events-none"
            />
            <input
              type="search"
              value={prefs.query}
              onChange={(e) => setPrefs((p) => ({ ...p, query: e.target.value }))}
              placeholder="SEARCH PROJECTS..."
              aria-label="Filter projects by name"
              disabled={!hydrated || totalCount === 0}
              className="w-full bg-[#0b0c10] border border-[#333] py-2 pl-9 pr-3 text-xs text-ca-ink placeholder:text-ca-ink-muted focus:border-ca-primary focus:outline-none disabled:opacity-50 uppercase"
            />
          </div>
          <div className="flex items-center gap-3">
            <span className="text-xs font-bold text-ca-ink-muted uppercase">Sort:</span>
            <select
              value={prefs.sort}
              onChange={(e) => setPrefs((p) => ({ ...p, sort: e.target.value as SortKey }))}
              aria-label="Sort projects"
              disabled={!hydrated || totalCount === 0}
              className="bg-[#0b0c10] border border-[#333] px-3 py-2 text-xs text-ca-ink focus:border-ca-primary focus:outline-none disabled:opacity-50 uppercase appearance-none cursor-pointer"
            >
              <option value="createdDesc">DESC (NEWEST)</option>
              <option value="createdAsc">ASC (OLDEST)</option>
              <option value="nameAsc">ALPHA (A-Z)</option>
              <option value="nameDesc">ALPHA (Z-A)</option>
            </select>
          </div>
        </div>

        {!hydrated ? (
          <div className="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
            {Array.from({ length: 6 }).map((_, i) => (
              <div key={i} className="border border-[#333] bg-[#1a1c23] p-5 h-32 animate-pulse flex flex-col justify-between">
                 <div className="h-4 bg-[#333] w-2/3"></div>
                 <div className="h-3 bg-[#333] w-1/3"></div>
              </div>
            ))}
          </div>
        ) : showEmpty ? (
          <div className="flex flex-col items-center justify-center p-12 border border-dashed border-[#444] bg-[#1a1c23]">
            <FolderOpen size={48} className="text-ca-ink-muted mb-4 opacity-50" />
            <h2 className="text-lg font-bold text-ca-ink uppercase tracking-widest mb-2">DB EMPTY</h2>
            <p className="text-sm text-ca-ink-muted text-center max-w-sm mb-6">
              No inspection projects exist in the current configuration matrix.
            </p>
            <button
              onClick={openDialog}
              className="border border-[#444] bg-transparent px-6 py-2 text-sm font-bold uppercase tracking-wider text-ca-ink hover:border-ca-primary transition-colors"
            >
              INITIALIZE PROJECT
            </button>
          </div>
        ) : showNoMatch ? (
           <div className="flex flex-col items-center justify-center p-12 border border-[#333] bg-[#1a1c23]">
             <span className="text-sm font-bold text-ca-ink-muted uppercase">QUERY_RESULT: 0 MATCHES</span>
           </div>
        ) : (
          <>
            {rowError ? (
              <div className="mb-4 rounded-sm border border-red-500/50 bg-red-500/10 p-3 text-xs text-red-500">
                {rowError}
              </div>
            ) : null}
            <div className="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
              {list.map((p) => {
                const busy = rowBusyId === p.id;
                const canRun = p.rulesetIds.length > 0 && !busy;

                return (
                  <div
                    key={p.id}
                    className="group relative flex flex-col border border-[#333] bg-[#1a1c23] transition-colors hover:border-ca-primary"
                  >
                    <div className="absolute left-0 top-0 bottom-0 w-[3px] bg-[#333] transition-colors group-hover:bg-ca-primary" />
                    
                    <div className="p-4 flex-1">
                      <div className="flex items-start justify-between mb-2">
                        <Link
                          to="/projects/$projectId"
                          params={{ projectId: toIntParam(IntAliasNamespaceType.Project, p.id) }}
                          className="flex-1 min-w-0 pr-4 hover:opacity-80"
                        >
                          <h2 className="truncate text-base font-bold uppercase tracking-wider text-ca-ink">
                            {p.name}
                          </h2>
                          <div className="mt-1 flex items-center gap-2">
                            <span className="text-[10px] text-ca-ink-muted">ID: {p.id.slice(0, 8)}</span>
                          </div>
                        </Link>
                        {p.deviceId && (
                           <div className="shrink-0 flex items-center justify-center border border-ca-primary/50 bg-ca-primary/10 px-2 py-0.5 text-[10px] text-ca-primary tracking-wider">
                              CAM_LINK
                           </div>
                        )}
                      </div>

                      <div className="mt-4 flex flex-wrap gap-2 text-[10px] uppercase font-bold tracking-wider">
                        <div className="border border-[#444] bg-[#0b0c10] px-2 py-1 text-ca-ink-muted">
                          {p.rulesetIds.length} {p.rulesetIds.length === 1 ? "RULESET" : "RULESETS"}
                        </div>
                      </div>
                    </div>

                    <div className="flex items-center justify-between border-t border-[#333] bg-[#0b0c10] px-4 py-2">
                      <div className="flex gap-2">
                         <button
                           title="Rename"
                           onClick={() => { setRowError(null); setRenameFor({ id: p.id, name: p.name }); }}
                           disabled={busy}
                           className="text-ca-ink-muted hover:text-ca-ink transition-colors disabled:opacity-50"
                         >
                           <Pencil size={14} />
                         </button>
                         <button
                           title="Duplicate"
                           onClick={() => void handleDuplicate(p.id)}
                           disabled={busy}
                           className="text-ca-ink-muted hover:text-ca-ink transition-colors disabled:opacity-50"
                         >
                           <Copy size={14} />
                         </button>
                         <button
                           title="Delete"
                           onClick={() => { setRowError(null); setDeleteFor({ id: p.id, name: p.name }); }}
                           disabled={busy}
                           className="text-ca-ink-muted hover:text-red-500 transition-colors disabled:opacity-50"
                         >
                           <Trash2 size={14} />
                         </button>
                      </div>
                      <button
                        onClick={() => void handleRun(p.id)}
                        disabled={!canRun}
                        className="flex items-center gap-2 border border-ca-primary bg-ca-primary/10 px-3 py-1 text-[10px] font-bold text-ca-primary hover:bg-ca-primary hover:text-black transition-colors disabled:opacity-30 disabled:border-[#444] disabled:text-[#444] disabled:bg-transparent"
                      >
                        <Play size={12} className={canRun ? "fill-current" : ""} />
                        EXECUTE
                      </button>
                    </div>
                  </div>
                );
              })}
            </div>
          </>
        )}
      </div>

      <AddDeviceModal
        isOpen={isAddDeviceOpen}
        onClose={() => setIsAddDeviceOpen(false)}
        onDeviceCreated={(created) => setDeviceId(created.id)}
      />

      {renameFor ? (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/80 p-4 backdrop-blur-sm">
          <form
            onSubmit={(e) => { e.preventDefault(); commitRename(); }}
            className="w-full max-w-sm border border-[#444] bg-[#0b0c10] shadow-2xl"
          >
            <div className="border-b border-[#333] bg-[#1a1c23] px-4 py-3">
              <h2 className="text-sm font-bold uppercase tracking-widest text-ca-ink">RENAME PROJECT</h2>
            </div>
            <div className="p-4">
              <label className="text-[10px] font-bold uppercase tracking-wider text-ca-ink-muted mb-2 block">
                NEW NAME
              </label>
              <input
                autoFocus
                value={renameFor.name}
                onChange={(e) => setRenameFor((prev) => (prev ? { ...prev, name: e.target.value } : prev))}
                className="w-full border border-[#444] bg-[#1a1c23] px-3 py-2 text-sm text-ca-ink focus:border-ca-primary focus:outline-none"
              />
            </div>
            <div className="flex gap-2 border-t border-[#333] bg-[#1a1c23] p-4">
              <button
                type="button"
                onClick={() => setRenameFor(null)}
                className="flex-1 border border-[#444] px-4 py-2 text-xs font-bold uppercase tracking-wider text-ca-ink hover:bg-[#333]"
              >
                CANCEL
              </button>
              <button
                type="submit"
                className="flex-1 bg-ca-primary px-4 py-2 text-xs font-bold uppercase tracking-wider text-black hover:bg-ca-primary-hover"
              >
                APPLY
              </button>
            </div>
          </form>
        </div>
      ) : null}

      {deleteFor ? (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/80 p-4 backdrop-blur-sm">
          <div className="w-full max-w-sm border border-red-500/50 bg-[#0b0c10] shadow-2xl">
            <div className="border-b border-red-500/20 bg-red-500/10 px-4 py-3">
              <h2 className="text-sm font-bold uppercase tracking-widest text-red-500">DELETE PROJECT?</h2>
            </div>
            <div className="p-4">
              <p className="text-xs text-ca-ink leading-relaxed">
                <span className="font-bold text-red-400">WARNING:</span> This removes "{deleteFor.name}" and every ruleset it owns from this browser. This cannot be undone.
              </p>
            </div>
            <div className="flex gap-2 border-t border-[#333] bg-[#1a1c23] p-4">
              <button
                type="button"
                onClick={() => setDeleteFor(null)}
                className="flex-1 border border-[#444] px-4 py-2 text-xs font-bold uppercase tracking-wider text-ca-ink hover:bg-[#333]"
              >
                CANCEL
              </button>
              <button
                type="button"
                onClick={commitDelete}
                className="flex-1 bg-red-500 px-4 py-2 text-xs font-bold uppercase tracking-wider text-white hover:bg-red-600"
              >
                PURGE
              </button>
            </div>
          </div>
        </div>
      ) : null}
    </div>
  );

  if (mode === UiModeType.Standard) {
    return (
      <StandardAppShell activeNav="projects" title="PROJECTS_DB">
        {mainContent}
      </StandardAppShell>
    );
  }

  return (
    <HmiShell title="Projects">
      <SectionTopBar section={SectionIdType.Home} active="projects" />
      {mainContent}
    </HmiShell>
  );
}
"""
    
    # We replace from start_str to the very end of the file.
    
    with open(file_path, "w", encoding="utf-8") as f:
        f.write(prefix + new_main_content)

if __name__ == "__main__":
    main()
