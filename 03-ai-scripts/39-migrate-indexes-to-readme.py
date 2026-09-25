#!/usr/bin/env python3
"""
39-migrate-indexes-to-readme.py

Autonomous Repository-Wide Index Migration Script:
1. Renames all 'readme.md' and 'readme.md' files to lowercase 'readme.md' across the repository using `git mv` (or filesystem rename if untracked).
2. Rewrites all occurrences of 'readme.md' and 'readme.md' to 'readme.md' in specs, prompts, scripts, links, and configs.
3. Automatically updates `scripts/sync-spec-tree.mjs` so 'readme.md' is recognized as an overview node in the specTree.
4. Regenerates `src/data/specTree.json` via `node scripts/sync-spec-tree.mjs`.
5. Verifies zero index files remain and reports stats.
"""

import os
import sys
import subprocess
from pathlib import Path

# Ensure UTF-8 output
sys.stdout.reconfigure(encoding="utf-8")

ROOT_DIR = Path(".").resolve()

EXCLUDED_DIRS = {
    ".git",
    "node_modules",
    ".gemini",
    "dist",
    "build",
    ".next",
    ".cache",
    ".turbo"
}

BINARY_EXTENSIONS = {
    ".png", ".jpg", ".jpeg", ".gif", ".ico", ".svg", ".webp",
    ".woff", ".woff2", ".ttf", ".eot",
    ".exe", ".dll", ".so", ".dylib", ".bin",
    ".zip", ".tar", ".gz", ".7z", ".rar",
    ".pyc", ".db", ".sqlite", ".sqlite3"
}

def is_excluded_path(p: Path) -> bool:
    for part in p.parts:
        if part in EXCLUDED_DIRS:
            return True
    return False

def find_index_files():
    found = []
    for root, dirs, files in os.walk(ROOT_DIR):
        root_path = Path(root)
        if is_excluded_path(root_path):
            continue
        for f in files:
            if f in ("readme.md", "readme.md"):
                found.append(root_path / f)
    return sorted(found)

def rename_index_files(index_files):
    renamed_count = 0
    for file_path in index_files:
        target_path = file_path.parent / "readme.md"
        rel_src = file_path.relative_to(ROOT_DIR).as_posix()
        rel_dst = target_path.relative_to(ROOT_DIR).as_posix()
        
        # Try git mv first
        res = subprocess.run(["git", "mv", rel_src, rel_dst], capture_output=True, text=True)
        if res.returncode == 0:
            print(f"[git mv] {rel_src} -> {rel_dst}")
            renamed_count += 1
        else:
            # Fallback to standard os.rename
            try:
                os.rename(file_path, target_path)
                print(f"[os.rename] {rel_src} -> {rel_dst}")
                renamed_count += 1
            except Exception as e:
                print(f"[ERROR] Failed to rename {rel_src}: {e}", file=sys.stderr)
    return renamed_count

def update_file_references():
    updated_files_count = 0
    total_replacements = 0

    for root, dirs, files in os.walk(ROOT_DIR):
        root_path = Path(root)
        if is_excluded_path(root_path):
            continue
        for f in files:
            file_path = root_path / f
            if file_path.suffix.lower() in BINARY_EXTENSIONS:
                continue
            
            try:
                content = file_path.read_text(encoding="utf-8")
            except (UnicodeDecodeError, PermissionError):
                continue
            
            new_content = content
            rep_count = 0
            
            if "readme.md" in new_content:
                rep_count += new_content.count("readme.md")
                new_content = new_content.replace("readme.md", "readme.md")
                
            if "readme.md" in new_content:
                rep_count += new_content.count("readme.md")
                new_content = new_content.replace("readme.md", "readme.md")
                
            if rep_count > 0:
                written = False
                for attempt in range(5):
                    try:
                        file_path.write_text(new_content, encoding="utf-8")
                        written = True
                        break
                    except OSError as e:
                        import time
                        time.sleep(0.1)
                if written:
                    updated_files_count += 1
                    total_replacements += rep_count
                    rel_p = file_path.relative_to(ROOT_DIR).as_posix()
                    if updated_files_count <= 20 or updated_files_count % 50 == 0:
                        print(f"[updated] {rel_p} ({rep_count} replacements)")
                else:
                    print(f"[error writing] {file_path}", file=sys.stderr)

    return updated_files_count, total_replacements

def update_spec_tree_script():
    script_path = ROOT_DIR / "scripts" / "sync-spec-tree.mjs"
    if not script_path.exists():
        return
    text = script_path.read_text(encoding="utf-8")
    old_target = """function deriveFileName(filename) {
  // "readme.md" -> "Overview"
  return deriveName(filename.replace(/\\.md$/, ""));
}"""
    new_replacement = """function deriveFileName(filename) {
  const base = filename.replace(/\\.md$/, "");
  if (base.toLowerCase() === "readme" || base === "01-index" || base === "00-index") {
    return "Overview";
  }
  return deriveName(base);
}"""
    if old_target in text:
        text = text.replace(old_target, new_replacement)
        script_path.write_text(text, encoding="utf-8")
        print("[updated] scripts/sync-spec-tree.mjs deriveFileName handler")

def regenerate_spec_tree():
    script_path = ROOT_DIR / "scripts" / "sync-spec-tree.mjs"
    if script_path.exists():
        res = subprocess.run(["node", "scripts/sync-spec-tree.mjs"], capture_output=True, text=True)
        if res.returncode == 0:
            print("[synced] src/data/specTree.json regenerated successfully")
        else:
            print(f"[warning] node scripts/sync-spec-tree.mjs failed: {res.stderr}", file=sys.stderr)

def main():
    print("=== Step 1: Finding readme.md and readme.md files ===")
    index_files = find_index_files()
    print(f"Found {len(index_files)} index files to rename.")

    print("\n=== Step 2: Renaming to readme.md ===")
    renamed = rename_index_files(index_files)
    print(f"Renamed {renamed} files.")

    print("\n=== Step 3: Updating scripts/sync-spec-tree.mjs ===")
    update_spec_tree_script()

    print("\n=== Step 4: Updating references in all repository files ===")
    updated_files, total_replacements = update_file_references()
    print(f"Updated {total_replacements} references across {updated_files} files.")

    print("\n=== Step 5: Regenerating specTree.json ===")
    regenerate_spec_tree()

    print("\n=== Step 6: Verifying remaining index files ===")
    remaining = find_index_files()
    if remaining:
        print(f"[WARNING] {len(remaining)} index files still found:", file=sys.stderr)
        for r in remaining:
            print(f"  {r.relative_to(ROOT_DIR).as_posix()}", file=sys.stderr)
    else:
        print("[SUCCESS] All readme.md and readme.md files successfully migrated to readme.md!")

if __name__ == "__main__":
    main()
