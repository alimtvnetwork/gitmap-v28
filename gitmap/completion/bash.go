package completion

// generateBash returns the Bash completion script.
func generateBash() string {
	return `_gitmap_completions() {
    local cur prev cmd sub
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"
    cmd="${COMP_WORDS[1]}"
    sub="${COMP_WORDS[2]}"

    if [[ ${COMP_CWORD} -eq 1 ]]; then
        COMPREPLY=($(compgen -W "$(gitmap completion --list-commands)" -- "$cur"))
        return
    fi

    case "$cmd" in
        cd|go)
            if [[ "$prev" == "--group" || "$prev" == "-g" ]]; then
                COMPREPLY=($(compgen -W "$(gitmap completion --list-groups)" -- "$cur"))
            else
                COMPREPLY=($(compgen -W "$(gitmap completion --list-repos) repos set-default clear-default" -- "$cur"))
            fi
            ;;
        pull)
            COMPREPLY=($(compgen -W "$(gitmap completion --list-repos)" -- "$cur"))
            ;;
        exec)
            if [[ "$prev" == "--group" ]]; then
                COMPREPLY=($(compgen -W "$(gitmap completion --list-groups)" -- "$cur"))
            fi
            ;;
        group|g)
            COMPREPLY=($(compgen -W "create add remove list show delete pull status exec clear $(gitmap completion --list-groups)" -- "$cur"))
            ;;
        list|ls)
            COMPREPLY=($(compgen -W "go node nodejs react cpp csharp groups --group --verbose" -- "$cur"))
            ;;
        multi-group|mg)
            COMPREPLY=($(compgen -W "pull status exec clear $(gitmap completion --list-groups)" -- "$cur"))
            ;;
        release|r)
            if [[ "$prev" == "--zip-group" ]]; then
                COMPREPLY=($(compgen -W "$(gitmap completion --list-zip-groups)" -- "$cur"))
            else
                COMPREPLY=($(compgen -W "--assets --commit --branch --bump --draft --dry-run --compress --checksums --bin --targets --list-targets --verbose --zip-group -Z --bundle" -- "$cur"))
            fi
            ;;
        release-branch|rb)
            COMPREPLY=($(compgen -W "--assets --draft --dry-run --compress --checksums --bin --targets" -- "$cur"))
            ;;
        alias|a)
            COMPREPLY=($(compgen -W "set remove list show suggest $(gitmap completion --list-aliases)" -- "$cur"))
            ;;
        zip-group|z)
            if [[ ${COMP_CWORD} -ge 3 ]] && [[ "$sub" == "add" || "$sub" == "show" || "$sub" == "delete" || "$sub" == "remove" || "$sub" == "rename" ]]; then
                COMPREPLY=($(compgen -W "$(gitmap completion --list-zip-groups)" -- "$cur"))
            else
                COMPREPLY=($(compgen -W "create add remove list show delete rename $(gitmap completion --list-zip-groups)" -- "$cur"))
            fi
            ;;
        dashboard|db)
            COMPREPLY=($(compgen -W "--limit --since --no-merges --out-dir --open" -- "$cur"))
            ;;
        ssh)
            if [[ "$prev" == "--name" || "$prev" == "-n" ]]; then
                COMPREPLY=($(compgen -W "$(gitmap completion --list-ssh-keys)" -- "$cur"))
            else
                COMPREPLY=($(compgen -W "cat list ls delete rm config --name --path --email --force" -- "$cur"))
            fi
            ;;
        clone-next|cn)
            COMPREPLY=($(compgen -W "v++ --delete --keep --no-desktop --ssh-key --verbose" -- "$cur"))
            ;;
        version-history|vh)
            COMPREPLY=($(compgen -W "--limit --json" -- "$cur"))
            ;;
        llm-docs|ld)
            if [[ "$prev" == "--format" ]]; then
                COMPREPLY=($(compgen -W "markdown json" -- "$cur"))
            elif [[ "$prev" == "--sections" ]]; then
                COMPREPLY=($(compgen -W "commands architecture flags conventions structure database installation patterns" -- "$cur"))
            else
                COMPREPLY=($(compgen -W "--stdout --format --sections" -- "$cur"))
            fi
            ;;
        help)
            if [[ "$prev" == "--compact" ]]; then
                COMPREPLY=($(compgen -W "$(gitmap completion --list-help-groups)" -- "$cur"))
            else
                COMPREPLY=($(compgen -W "--compact" -- "$cur"))
            fi
            ;;
        mv|move)
            COMPREPLY=($(compgen -W "--prefer-newer --prefer-left --prefer-right --dry-run --verbose" -- "$cur"))
            ;;
        merge-both|mb|merge-left|ml|merge-right|mr)
            COMPREPLY=($(compgen -W "--prefer-newer --prefer-left --prefer-right --prefer-larger --dry-run --verbose" -- "$cur"))
            ;;
        vscode-pm-sync|vpm)
            if [[ "$prev" == "--mode" ]]; then
                COMPREPLY=($(compgen -W "union replace intersection" -- "$cur"))
            elif [[ "$prev" == "--projects-json" ]]; then
                COMPREPLY=($(compgen -f -- "$cur"))
            else
                COMPREPLY=($(compgen -W "--dry-run --projects-json --tag --mode" -- "$cur"))
            fi
            ;;
        *)
            if [[ "$prev" == "-A" || "$prev" == "--alias" ]]; then
                COMPREPLY=($(compgen -W "$(gitmap completion --list-aliases)" -- "$cur"))
            elif [[ "$prev" == "--zip-group" ]]; then
                COMPREPLY=($(compgen -W "$(gitmap completion --list-zip-groups)" -- "$cur"))
            elif [[ "$prev" == "--ssh-key" || "$prev" == "-K" ]]; then
                COMPREPLY=($(compgen -W "$(gitmap completion --list-ssh-keys)" -- "$cur"))
            fi
            ;;
    esac
}
complete -F _gitmap_completions gitmap
`
}
