// Package completion — fish.go generates Fish shell tab-completion scripts.
package completion

// generateFish returns the Fish completion script.
func generateFish() string {
	return `# Fish completion script for gitmap

function __fish_gitmap_needs_command
    set -l cmd (commandline -opc)
    if test (count $cmd) -eq 1
        return 0
    end
    return 1
end

function __fish_gitmap_using_command
    set -l cmd (commandline -opc)
    if test (count $cmd) -gt 1
        if test $cmd[2] = $argv[1]
            return 0
        end
    end
    return 1
end

complete -c gitmap -n '__fish_gitmap_needs_command' -f -a '(gitmap completion --list-commands)' -d 'GitMap command'
complete -c gitmap -n '__fish_gitmap_using_command cd' -f -a '(gitmap completion --list-repos)' -d 'Repository'
complete -c gitmap -n '__fish_gitmap_using_command pull' -f -a '(gitmap completion --list-repos)' -d 'Repository'
complete -c gitmap -n '__fish_gitmap_using_command group' -f -a 'create add remove list show delete pull status exec clear (gitmap completion --list-groups)' -d 'Group subcommand'
complete -c gitmap -n '__fish_gitmap_using_command alias' -f -a 'set remove list show suggest (gitmap completion --list-aliases)' -d 'Alias subcommand'
complete -c gitmap -n '__fish_gitmap_using_command zip-group' -f -a 'create add remove list show delete rename (gitmap completion --list-zip-groups)' -d 'Zip group subcommand'
complete -c gitmap -n '__fish_gitmap_using_command agy' -f -a 'ls add rm update scan reconcile remove-missing pin-projects group find-duplicates optimize-projects clean-cache status stats all-projects-read-memory-prompt fix-pipeline prompt sync export import plugins' -d 'Antigravity command'
complete -c gitmap -n '__fish_gitmap_using_command tasks' -f -a 'list history undo redo clear' -d 'Tasks command'
complete -c gitmap -l help -s h -d 'Show help'
complete -c gitmap -l verbose -s v -d 'Verbose output'
complete -c gitmap -l version -d 'Show version'
`
}
