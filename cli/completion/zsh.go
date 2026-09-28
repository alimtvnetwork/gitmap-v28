package completion

// generateZsh returns the Zsh completion script.
func generateZsh() string {
	return `#compdef gitmap

_gitmap() {
    local -a completions
    local -a args
    args=("${words[@]:1}")
    
    local output
    output=$(gitmap __complete "${args[@]}" 2>/dev/null)
    if [[ -z "$output" ]]; then
        _files
        return
    fi

    local IFS=$'\n'
    for line in ${(f)output}; do
        if [[ "$line" =~ ^:[0-9]+$ ]]; then
            continue
        fi
        local val="${line%%	*}"
        local desc="${line#*	}"
        if [[ -n "$val" ]]; then
            if [[ "$val" == "$desc" || -z "$desc" ]]; then
                completions+=("$val")
            else
                completions+=("$val:$desc")
            fi
        fi
    done

    if (( ${#completions[@]} > 0 )); then
        _describe -t gitmap-commands 'gitmap commands' completions
    else
        _files
    fi
}

_gitmap "$@"
`
}
