import re
import sys

def main():
    file_path = "src/routes/projects/$projectId/rulesets/index.tsx"
    with open(file_path, "r", encoding="utf-8") as f:
        content = f.read()

    start_str = "return ("
    start_idx = content.find(start_str, content.find("function RulesetsList"))
    end_idx = content.find("function RulesetsError", start_idx)

    if start_idx == -1 or end_idx == -1:
        print("Could not find RulesetsList return statement or end")
        sys.exit(1)

    prefix = content[:start_idx]
    suffix = content[end_idx:]

    new_content = """return (
    <div className="flex min-w-0 flex-1 flex-col overflow-auto bg-[#0b0c10] font-mono text-ca-ink antialiased p-6">
      <div className="mx-auto w-full max-w-7xl">
        <header className="mb-6 flex flex-wrap items-center justify-between gap-4 border-b border-[#333] pb-4">
          <div className="flex items-center gap-4 min-w-0">
            <h1 className="text-xl font-bold uppercase tracking-wider text-ca-ink flex items-center gap-2">
              <span className="w-2 h-2 bg-ca-primary rounded-full animate-pulse" />
              INSPECTION RULESETS
            </h1>
            <div className="flex items-center gap-2 border border-[#333] bg-[#1a1c23] px-3 py-1">
               <span className="text-xs font-bold text-ca-ink-muted uppercase tracking-widest">{project.name}</span>
            </div>
            <span className="text-xs font-bold tracking-widest tabular-nums text-ca-ink-muted">
              {rulesets.length} {rulesets.length === 1 ? "PROGRAM" : "PROGRAMS"}
            </span>
          </div>
          <Link
            to="/projects/$projectId/rulesets/new"
            params={{ projectId }}
            className="inline-flex items-center gap-2 bg-ca-primary px-4 py-2 text-xs font-bold uppercase tracking-wider text-black hover:bg-ca-primary-hover transition-colors focus-visible:outline-none"
          >
            <Plus aria-hidden size={16} />
            NEW RULESET
          </Link>
        </header>

        {rulesets.length === 0 ? (
          <div className="flex flex-col items-center justify-center p-12 border border-dashed border-[#444] bg-[#1a1c23]">
            <ImageIcon size={48} className="text-[#444] mb-4" />
            <h2 className="text-lg font-bold text-ca-ink uppercase tracking-widest mb-2">NO RULESETS FOUND</h2>
            <p className="text-sm text-ca-ink-muted text-center max-w-sm mb-6">
              Author a ruleset from a reference image. Trial runs and AI testing use the rules you draw here.
            </p>
            <Link
              to="/projects/$projectId/rulesets/new"
              params={{ projectId }}
              className="border border-[#444] bg-transparent px-6 py-2 text-sm font-bold uppercase tracking-wider text-ca-ink hover:border-ca-primary transition-colors"
            >
              INITIALIZE RULESET
            </Link>
          </div>
        ) : (
          <ul className="grid grid-cols-1 gap-4 md:grid-cols-2 xl:grid-cols-3">
            {rulesets.map((r) => (
              <li key={r.id} className="group relative flex flex-col border border-[#333] bg-[#1a1c23] transition-colors hover:border-ca-primary">
                <div className="absolute left-0 top-0 bottom-0 w-[3px] bg-[#333] transition-colors group-hover:bg-ca-primary" />
                
                <Link
                  to="/projects/$projectId/rulesets/$rulesetId"
                  params={{ projectId, rulesetId: toIntParam(IntAliasNamespaceType.Ruleset, r.id) }}
                  className="flex flex-col h-full p-4 focus-visible:outline-none"
                >
                  <div className="flex justify-between items-start mb-2">
                    <h2 className="font-bold text-base uppercase tracking-wider text-ca-ink truncate pr-8">
                      {r.name}
                    </h2>
                  </div>
                  
                  <div className="mt-2 flex flex-wrap gap-2 text-[10px] uppercase font-bold tracking-wider">
                     <span className="border border-[#444] bg-[#0b0c10] px-2 py-1 text-ca-ink-muted">
                       {r.rules.length} {r.rules.length === 1 ? "RULE" : "RULES"}
                     </span>
                     {r.imageRef ? (
                        <span className="border border-ca-primary/50 bg-ca-primary/10 px-2 py-1 text-ca-primary">
                          IMAGE SYNCED
                        </span>
                     ) : (
                        <span className="border border-yellow-500/50 bg-yellow-500/10 px-2 py-1 text-yellow-500">
                          NO IMAGE
                        </span>
                     )}
                  </div>
                  <p className="mt-4 font-mono text-[10px] text-ca-ink-muted">ID: {r.id}</p>
                </Link>
                
                <button
                  type="button"
                  onClick={(e) => {
                    e.preventDefault();
                    e.stopPropagation();
                    const isConfirmed = window.confirm(`Delete ruleset "${r.name}"? This removes its rules from this project.`);
                    if (!isConfirmed) return;
                    deleteRuleset(r.id);
                    notifySuccess(`Ruleset "${r.name}" deleted.`);
                  }}
                  className="absolute right-3 top-3 inline-flex border border-[#444] bg-[#0b0c10] p-1.5 text-[#666] transition-colors hover:border-red-500 hover:text-red-500 focus-visible:outline-none"
                  title="Delete rule set"
                  aria-label={`Delete rule set ${r.name}`}
                >
                  <Trash2 aria-hidden size={14} />
                </button>
                
                <div className="border-t border-[#333] bg-[#0b0c10] px-4 py-2 flex items-center justify-between">
                   <span className="text-[10px] font-bold uppercase tracking-widest text-ca-ink-muted">TRIAL RUN</span>
                   <RulesetTestRunPill ruleset={r} />
                </div>
              </li>
            ))}
          </ul>
        )}
      </div>
    </div>
  );
}

"""

    with open(file_path, "w", encoding="utf-8") as f:
        f.write(prefix + new_content + suffix)

if __name__ == "__main__":
    main()
