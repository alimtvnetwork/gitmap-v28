import sys
import re

def main():
    file_path = "src/routes/index.tsx"
    with open(file_path, "r", encoding="utf-8") as f:
        content = f.read()

    # We will replace the entire Index component with a new one
    start_str = "function Index() {"
    end_str = "export default Index;\n" # Or something like it. Wait, the end of file might not be export default.
    # Actually, Index is exported in the Route definition: component: Index
    
    # We can just run a python script that completely replaces everything below `function Index() {`
    # Let's find `function Index() {`
    idx = content.find("function Index() {")
    if idx == -1:
        print("Could not find function Index() {")
        sys.exit(1)
        
    prefix = content[:idx]
    
    new_index = """function Index() {
  return (
    <div className="flex h-full w-full flex-col bg-[#0b0c10] text-ca-ink font-mono antialiased overflow-hidden">
      {/* Header */}
      <header className="flex h-14 shrink-0 items-center border-b border-[#22252a] bg-[#111318] px-6">
        <div className="flex flex-col">
          <span className="text-[10px] font-bold text-[#f5a623] uppercase tracking-wider">SYSTEM STATUS</span>
          <h1 className="text-sm font-bold text-ca-ink tracking-wide">HMI MAIN CONSOLE</h1>
        </div>
        <div className="ml-auto flex items-center gap-4">
          <div className="flex items-center gap-2">
            <span className="relative flex h-2 w-2">
              <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-[#00ff00] opacity-75"></span>
              <span className="relative inline-flex rounded-full h-2 w-2 bg-[#00ff00]"></span>
            </span>
            <span className="text-[10px] font-bold uppercase tracking-wider text-ca-ink-muted">VISION ENGINE ACTIVE</span>
          </div>
        </div>
      </header>

      {/* Main Grid */}
      <div className="flex-1 overflow-auto p-6">
        <div className="mx-auto grid max-w-7xl grid-cols-1 gap-6 lg:grid-cols-2">
          
          {/* Workflows */}
          {WORKFLOWS.map((wf) => (
            <div key={wf.id} className="flex flex-col border border-[#22252a] bg-[#111318] transition-colors hover:border-[#444]">
              <div className="flex items-center justify-between border-b border-[#22252a] bg-[#15171e] px-4 py-3">
                <div className="flex items-center gap-3">
                  <wf.icon size={16} className={wf.tone === "amber" ? "text-[#f5a623]" : wf.tone === "cyan" ? "text-cyan-400" : wf.tone === "green" ? "text-[#00ff00]" : "text-purple-400"} />
                  <h2 className="text-xs font-bold uppercase tracking-widest text-ca-ink">{wf.label}</h2>
                </div>
                <Link
                  to={wf.to}
                  className="inline-flex items-center gap-2 border border-[#444] bg-[#0b0c10] px-3 py-1.5 text-[10px] font-bold uppercase tracking-wider text-ca-ink transition hover:border-[#f5a623] hover:text-[#f5a623]"
                >
                  ENTER <ArrowRight size={12} />
                </Link>
              </div>
              <div className="flex-1 p-4">
                <p className="mb-6 text-xs text-ca-ink-muted">{wf.description}</p>
                <div className="grid grid-cols-2 gap-3">
                  {wf.quickActions.map((qa) => (
                    <Link
                      key={qa.label}
                      to={qa.to as any}
                      className="flex items-center justify-between border border-[#333] bg-[#0b0c10] px-3 py-2 transition-colors hover:border-[#666] hover:bg-[#1a1c23]"
                    >
                      <div className="flex items-center gap-2">
                        <qa.icon size={14} className="text-ca-ink-muted" />
                        <span className="text-[10px] font-bold uppercase tracking-wider text-ca-ink">{qa.label}</span>
                      </div>
                      <ArrowRight size={10} className="text-[#333]" />
                    </Link>
                  ))}
                </div>
              </div>
            </div>
          ))}

        </div>
      </div>
    </div>
  );
}
"""
    with open(file_path, "w", encoding="utf-8") as f:
        f.write(prefix + new_index)

if __name__ == "__main__":
    main()
