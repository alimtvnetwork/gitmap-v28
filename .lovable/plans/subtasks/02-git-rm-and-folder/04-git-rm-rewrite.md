# Subtask 4: Git History Rewrite & Backup

1. In gitmap/cmd/gitrm/, implement the backup feature. Resolve the GLOBAL .gitmap installation directory (e.g., via os.UserHomeDir() + "/.gitmap/backups/git-rm/").
2. Copy the files out of the repository at HEAD to the backup location before modifying history.
3. Use a system command (like git filter-branch --force --index-filter 'git rm --cached --ignore-unmatch <files>' --prune-empty --tag-name-filter cat -- --all) to wipe the files from git history.
4. Clean up efs/original/ and run git gc --aggressive --prune=now to permanently delete the objects.
