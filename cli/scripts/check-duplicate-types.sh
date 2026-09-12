#!/usr/bin/env bash
# check-duplicate-types.sh — detects duplicate type declarations across files
# in the same Go package. Exits with code 1 if duplicates are found.

set -euo pipefail

root="${1:-.}"
found=0

# Find all Go packages (directories containing .go files)
while IFS= read -r pkg_dir; do
    # Extract all "type <Name>" declarations with their file
    declarations=$(grep -rn '^type [A-Z][A-Za-z0-9_]* ' "$pkg_dir"/*.go 2>/dev/null \
        | grep -v '_test.go:' \
        | sed 's/:.*type \([A-Z][A-Za-z0-9_]*\) .*/:\1/' \
        || true)

    [ -z "$declarations" ] && continue

    # Group by type name and check for duplicates across different files
    echo "$declarations" | awk -F: '{
        file = $1; name = $NF
        if (name in seen && seen[name] != file) {
            printf "  DUPLICATE: type %s\n    -> %s\n    -> %s\n", name, seen[name], file
            found = 1
        }
        seen[name] = file
    } END { exit found }' || {
        echo "✗ Duplicate type(s) found in package: $pkg_dir"
        found=1
    }
done < <(find "$root" -name '*.go' -not -path '*/vendor/*' -not -name '*_test.go' \
    -exec dirname {} \; | sort -u)

if [ "$found" -eq 0 ]; then
    echo "✓ No duplicate type declarations found."
    exit 0
else
    echo ""
    echo "Fix: rename one declaration or move it to a shared file."
    exit 1
fi
