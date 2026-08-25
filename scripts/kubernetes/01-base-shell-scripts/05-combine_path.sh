#!/bin/bash

# Function to combine two paths ensuring there's only one slash between them
# Parameters:
#   path1: The first path.
#   path2: The second path.
# Returns: Combined path with a single slash.
# Example:
#   path1="/path/to/source/"
#   path2="/directory/inside/"
#   combined_path=$(combine_path "$path1" "$path2")
#   echo $combined_path  # Outputs: /path/to/source/directory/inside/
combine_path() {
    local path1="$1"
    local path2="$2"
    
    # Remove trailing slash from path1 if it exists
    if [[ "$path1" == */ ]]; then
        path1="${path1%/}"
    fi

    # Remove leading slash from path2 if it exists
    if [[ "$path2" == /* ]]; then
        path2="${path2#/}"
    fi

    # Combine the two paths with a single slash in between
    echo "$path1/$path2"
}

# Function to combine the source path with the base directory of the dest path
# Parameters:
#   source_path: The source path.
#   dest_path: The destination path whose base directory will be used.
# Returns: Combined path with the source path and base directory of the destination path.
# Example:
#   source_path="/path/to/source"
#   dest_path="/path/to/machine/destination"
#   final_path=$(combine_with_base_path "$source_path" "$dest_path")
#   echo $final_path  # Outputs: /path/to/source/destination
combine_with_base_path() {
    local source_path="$1"
    local dest_path="$2"
    local base_dir
    base_dir=$(basename "$dest_path")
    combine_path "$source_path" "$base_dir"
}