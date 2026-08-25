#!/bin/bash

# Directory containing shell files to be sourced
SHELL_DIR="$PWD/../01-base-shell-scripts"
CURRENT_FILENAME="00-import-all.sh"
EXCLUDE_FILES=("01-logger.sh" "$CURRENT_FILENAME")
# Source the logger file first
source "$SHELL_DIR/01-logger.sh"

# log_message "$CURRENT_FILENAME"

# Function to check if a filename is missing from the exclude list
is_filename_missing() {
    local file="$1"
    for excluded in "${EXCLUDE_FILES[@]}"; do
        if [ "$file" == "$excluded" ]; then
            return 1
        fi
    done
    return 0
}

# Function to source all shell files except excluded ones
source_shell_files() {
    for file in "$SHELL_DIR"/*.sh; do
        filename=$(basename "$file")
        if is_filename_missing "$filename"; then
            # log_message "Sourcing - $file"
            source "$file"
        fi
    done
}

# Call the function to source the shell files
source_shell_files
