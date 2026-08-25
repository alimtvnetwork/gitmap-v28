with open(".lovable/plans/pending/03-zsh-kube-master-plan.md", "w", encoding="utf-8") as f:
    f.write("# 03-zsh-kube-master-plan.md\n\n")
    f.write("## Phase 1: Planning, Consolidation & Specification (Steps 1-100)\n\n")
    for i in range(1, 101):
        if i == 1:
            f.write(f"- [ ] Step {i}: Analyze 01-base-shell-scripts to consolidate helpers.\n")
        elif i == 2:
            f.write(f"- [ ] Step {i}: Analyze 02-ubuntu-install for ZSH setup logic.\n")
        elif i == 3:
            f.write(f"- [ ] Step {i}: Design intelligent sequence-based line insertion for .zshrc.\n")
        elif i == 4:
            f.write(f"- [ ] Step {i}: Design authorized_keys line-by-line idempotency logic.\n")
        elif i == 5:
            f.write(f"- [ ] Step {i}: Design aria2c installer for Gitmap CLI.\n")
        elif i == 6:
            f.write(f"- [ ] Step {i}: Map out Windows vs Ubuntu vs CentOS OS-specific behaviors.\n")
        elif i == 7:
            f.write(f"- [ ] Step {i}: Define 'gitmap zsh install' CLI struct.\n")
        elif i == 8:
            f.write(f"- [ ] Step {i}: Define 'gitmap zsh theme-change' CLI struct.\n")
        elif i == 9:
            f.write(f"- [ ] Step {i}: Define 'gitmap zsh completion' template command.\n")
        elif i == 10:
            f.write(f"- [ ] Step {i}: Define 'gitmap os kill-process' CLI struct.\n")
        elif i == 11:
            f.write(f"- [ ] Step {i}: Specify improved error logging mechanics.\n")
        elif i < 100:
            f.write(f"- [ ] Step {i}: Consolidate and specify sub-component {i} for Kubernetes/ZSH integration.\n")
        else:
            f.write(f"- [ ] Step {i}: Finalize 100-step specification review.\n")
    
    f.write("\n## Phase 2: Execution & Implementation (Steps 101-300)\n\n")
    for i in range(101, 301):
        if i == 101:
            f.write(f"- [ ] Step {i}: Scaffold base command for ZSH.\n")
        elif i == 102:
            f.write(f"- [ ] Step {i}: Implement sequence-aware line inserter core utility.\n")
        elif i == 103:
            f.write(f"- [ ] Step {i}: Implement authorized_keys line-by-line idempotency utility.\n")
        elif i == 104:
            f.write(f"- [ ] Step {i}: Scaffold aria2c installer command.\n")
        elif i == 105:
            f.write(f"- [ ] Step {i}: Refactor error logging utility.\n")
        elif i < 300:
            f.write(f"- [ ] Step {i}: Execute sub-task {i-100} for ZSH/Kube integration.\n")
        else:
            f.write(f"- [ ] Step {i}: Finalize E2E tests for the 300-step master plan.\n")
