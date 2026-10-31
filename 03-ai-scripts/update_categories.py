import re
import sys

def main():
    file_path = "src/routes/projects/$projectId/categories.tsx"
    with open(file_path, "r", encoding="utf-8") as f:
        content = f.read()

    start_str = "function ProjectCategoriesTab() {"
    start_idx = content.find(start_str)
    if start_idx == -1:
        print("Could not find ProjectCategoriesTab")
        sys.exit(1)

    prefix = content[:start_idx]

    new_component = """function ProjectCategoriesTab() {
  const { projectId } = Route.useParams();
  const project = useProjectStore((s) => selectProject(s, projectId));
  const rulesets = useProjectStore((s) => selectRulesetsForProject(s, projectId));
  const addProjectCategory = useProjectStore((s) => s.addProjectCategory);
  const renameProjectCategory = useProjectStore((s) => s.renameProjectCategory);
  const deleteProjectCategory = useProjectStore((s) => s.deleteProjectCategory);
  const updateRulesetCategory = useProjectStore((s) => s.updateRulesetCategory);
  const categories = project?.categoryNames ?? [];
  const [newName, setNewName] = useState("");
  const seeded = useSeededEmptyStateAction("categories.list");

  return (
    <section
      aria-labelledby="categories-tab-heading"
      className="flex min-h-0 flex-1 flex-col overflow-auto bg-[#0b0c10] text-ca-ink font-mono antialiased p-6"
      data-project-id={projectId}
    >
      <div className="mx-auto w-full max-w-7xl">
        <header className="mb-6 flex flex-wrap items-center justify-between gap-4 border-b border-[#333] pb-4">
          <div className="flex items-center gap-3">
            <Tag className="text-ca-primary" size={20} aria-hidden />
            <h2 id="categories-tab-heading" className="text-xl font-bold uppercase tracking-wider text-ca-ink">
              Taxonomy Manager
            </h2>
          </div>
          <div className="flex items-center gap-2">
             <span className="text-xs text-ca-ink-muted uppercase font-bold tracking-widest bg-[#1a1c23] border border-[#333] px-3 py-1">
               {categories.length} CATEGORIES
             </span>
             <span className="text-xs text-ca-ink-muted uppercase font-bold tracking-widest bg-[#1a1c23] border border-[#333] px-3 py-1">
               {rulesets.length} RULESETS
             </span>
          </div>
        </header>

        <div className="mb-6 flex gap-3 border border-[#333] bg-[#1a1c23] p-4">
          <div className="flex-1 relative">
             <span className="absolute left-3 top-1/2 -translate-y-1/2 text-xs font-bold text-ca-ink-muted uppercase tracking-wider">NEW:</span>
             <input
               value={newName}
               onChange={(event) => setNewName(event.currentTarget.value)}
               placeholder="CATEGORY_NAME"
               className="w-full bg-[#0b0c10] border border-[#333] py-2 pl-12 pr-3 text-sm text-ca-ink focus:border-ca-primary focus:outline-none uppercase"
             />
          </div>
          <button
            type="button"
            onClick={() => {
              addProjectCategory(projectId, newName);
              setNewName("");
            }}
            disabled={newName.trim() === ""}
            className="flex items-center gap-2 bg-ca-primary px-6 py-2 text-xs font-bold uppercase tracking-wider text-black hover:bg-ca-primary-hover disabled:opacity-30 transition-colors"
          >
            <FolderPlus size={16} />
            ALLOCATE
          </button>
        </div>

        <div className="grid grid-cols-1 gap-6 lg:grid-cols-2">
          {/* Left Column: CRUD */}
          <div className="flex flex-col border border-[#333] bg-[#1a1c23]">
            <div className="border-b border-[#333] bg-[#0b0c10] px-4 py-3">
              <h3 className="text-xs font-bold uppercase tracking-widest text-ca-ink-muted flex items-center gap-2">
                <span className="w-2 h-2 bg-ca-primary"></span>
                ACTIVE CATEGORIES
              </h3>
            </div>
            <div className="p-4 space-y-3">
              {categories.map((category) => (
                <CategoryEditorRow
                  key={category}
                  category={category}
                  count={rulesets.filter((rule) => rule.categoryName === category).length}
                  onRename={(next) => renameProjectCategory(projectId, category, next)}
                  onDelete={() => deleteProjectCategory(projectId, category)}
                />
              ))}
              {categories.length === 0 ? (
                <div className="flex flex-col items-center justify-center p-8 border border-dashed border-[#444] bg-[#0b0c10]">
                  <FolderPlus size={32} className="text-[#444] mb-4" />
                  <span className="text-xs font-bold text-ca-ink-muted uppercase tracking-widest">NO CATEGORIES DECLARED</span>
                </div>
              ) : null}
            </div>
          </div>

          {/* Right Column: Assignment */}
          <div className="flex flex-col border border-[#333] bg-[#1a1c23]">
            <div className="border-b border-[#333] bg-[#0b0c10] px-4 py-3">
              <h3 className="text-xs font-bold uppercase tracking-widest text-ca-ink-muted flex items-center gap-2">
                <span className="w-2 h-2 bg-[#555]"></span>
                RULESET ALLOCATION
              </h3>
            </div>
            <div className="p-4 space-y-3">
              {rulesets.map((rule) => (
                <label
                  key={rule.id}
                  className="grid grid-cols-[minmax(0,1fr)_auto] items-center gap-4 border border-[#333] bg-[#0b0c10] p-2 hover:border-[#555] transition-colors cursor-pointer"
                >
                  <span className="truncate text-xs font-bold text-ca-ink tracking-wider pl-2 uppercase">{rule.name}</span>
                  <select
                    value={rule.categoryName ?? ""}
                    onChange={(event) =>
                      updateRulesetCategory(rule.id, event.currentTarget.value || undefined)
                    }
                    className="w-48 border border-[#333] bg-[#1a1c23] px-3 py-1.5 text-xs text-ca-primary font-bold uppercase tracking-wider focus:border-ca-primary focus:outline-none appearance-none"
                  >
                    <option value="">[ UNASSIGNED ]</option>
                    {categories.map((category) => (
                      <option key={category} value={category}>
                        {category}
                      </option>
                    ))}
                  </select>
                </label>
              ))}
              {rulesets.length === 0 ? (
                <div className="flex flex-col items-center justify-center p-8 border border-dashed border-[#444] bg-[#0b0c10]">
                  <Layers size={32} className="text-[#444] mb-4" />
                  <span className="text-xs font-bold text-ca-ink-muted uppercase tracking-widest">NO RULESETS AWAITING</span>
                </div>
              ) : null}
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}

function CategoryEditorRow({
  category,
  count,
  onRename,
  onDelete,
}: {
  category: string;
  count: number;
  onRename: (name: string) => void;
  onDelete: () => void;
}) {
  const [name, setName] = useState(category);

  return (
    <div className="flex items-center gap-2 border border-[#333] bg-[#0b0c10] p-2 transition-colors focus-within:border-ca-primary hover:border-[#555]">
      <input
        value={name}
        onChange={(event) => setName(event.currentTarget.value)}
        onBlur={() => onRename(name)}
        aria-label={`Rename ${category}`}
        className="flex-1 bg-transparent px-2 text-xs font-bold uppercase tracking-wider text-ca-ink focus:outline-none"
      />
      <div className="flex shrink-0 items-center gap-3 pr-2">
        <span className="text-[10px] font-bold text-ca-ink-muted tracking-widest uppercase">
          {count} RULES
        </span>
        <button
          type="button"
          onClick={onDelete}
          aria-label={`Delete ${category}`}
          className="text-[#666] hover:text-red-500 transition-colors"
        >
          <Trash2 aria-hidden size={16} />
        </button>
      </div>
    </div>
  );
}
"""

    with open(file_path, "w", encoding="utf-8") as f:
        f.write(prefix + new_component)

if __name__ == "__main__":
    main()
