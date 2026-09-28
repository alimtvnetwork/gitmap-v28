package completion

// generateBash returns the Bash completion script.
func generateBash() string {
	return `_gitmap_completions() {
    local cur="${COMP_WORDS[COMP_CWORD]}"
    local words=("${COMP_WORDS[@]}")
    
    local results
    results=$(gitmap __complete "${words[@]:1}" 2>/dev/null)
    if [[ $? -ne 0 || -z "$results" ]]; then
        return
    fi

    local IFS=$'\n'
    local filtered=()
    for line in $results; do
        if [[ "$line" =~ ^:[0-9]+$ ]]; then
            continue
        fi
        local val="${line%%	*}"
        if [[ -n "$val" ]]; then
            filtered+=("$val")
        fi
    done

    COMPREPLY=($(compgen -W "${filtered[*]}" -- "$cur"))
}
complete -o default -F _gitmap_completions gitmap
`
}
