import os
import sys
import argparse
import subprocess
import fnmatch
import re
import platform
from pathlib import Path

DEFAULT_IGNORES = ['.git', 'node_modules']

def normalize_path(path_str):
    p = os.path.abspath(path_str)
    if platform.system() == 'Windows' and not p.startswith('\\\\?\\'):
        return '\\\\?\\' + p
    return p

def run_git_mv(src, dst, cwd):
    # Try git mv
    try:
        res = subprocess.run(
            ['git', 'mv', src, dst],
            cwd=cwd, stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True
        )
        if res.returncode == 0:
            return True
        else:
            return False
    except Exception:
        return False

def rename_file(src_path, dst_path):
    if src_path == dst_path:
        return

    src_norm = normalize_path(src_path)
    dst_norm = normalize_path(dst_path)
    parent_dir = os.path.dirname(src_norm)
    src_name = os.path.basename(src_norm)
    dst_name = os.path.basename(dst_norm)

    if run_git_mv(src_name, dst_name, cwd=parent_dir):
        print(f"✅ git mv: {src_path} -> {dst_path}")
        return
    
    # Fallback to os.rename
    try:
        os.rename(src_norm, dst_norm)
        print(f"✅ os.rename: {src_path} -> {dst_path}")
    except Exception as e:
        print(f"❌ Failed to rename {src_path}: {e}")

def should_ignore(path_name, ignores):
    for pattern in ignores:
        if fnmatch.fnmatch(path_name, pattern) or fnmatch.fnmatch(os.path.basename(path_name), pattern):
            return True
    return False

def cmd_lowercase(args):
    target = args.target_directory
    ignores = DEFAULT_IGNORES.copy()
    if args.except_list:
        ignores.extend([x.strip() for x in args.except_list.split(',')])

    target_norm = normalize_path(target)
    for root, dirs, files in os.walk(target_norm):
        # Filter dirs in place to respect ignores
        dirs[:] = [d for d in dirs if not should_ignore(os.path.join(root, d), ignores)]
        
        for f in files:
            file_path = os.path.join(root, f)
            if should_ignore(file_path, ignores):
                continue
            
            lower_f = f.lower()
            if f != lower_f:
                dst_path = os.path.join(root, lower_f)
                rename_file(file_path, dst_path)

def extract_sequence(filename):
    match = re.match(r'^(\d+)[-_](.*)$', filename)
    if match:
        return int(match.group(1)), match.group(2)
    return None, filename

def format_sequence(seq, rest, max_digits=2):
    return f"{seq:0{max_digits}d}-{rest}"

def parse_pins(pin_str):
    pin_map = {}
    if not pin_str:
        return pin_map
    for pair in pin_str.split(','):
        if '=' not in pair:
            continue
        name, seq = pair.split('=', 1)
        pin_map[name.strip().lower()] = int(seq.strip())
    return pin_map


def cmd_fix_seq(args):
    target = args.target_directory
    target_norm = normalize_path(target)
    pin_map = parse_pins(args.pin)

    if not os.path.isdir(target_norm):
        print(f"❌ Error: {target} is not a directory.")
        return

    files = [f for f in os.listdir(target_norm) if os.path.isfile(os.path.join(target_norm, f))]
    
    # Parse files
    parsed_files = []
    for f in files:
        seq, rest = extract_sequence(f)
        stat = os.stat(os.path.join(target_norm, f))
        parsed_files.append({
            'original': f,
            'seq': seq,
            'rest': rest,
            'rest_lower': rest.lower(),
            'time': stat.st_mtime
        })

    # Sort logic
    if args.order_by_time:
        parsed_files.sort(key=lambda x: x['time'])
    elif args.order_by_az:
        parsed_files.sort(key=lambda x: x['rest_lower'])
    
    # Apply pins
    pinned = []
    unpinned = []
    for pf in parsed_files:
        base_no_ext = os.path.splitext(pf['rest_lower'])[0]
        if base_no_ext in pin_map:
            pf['new_seq'] = pin_map[base_no_ext]
            pinned.append(pf)
        else:
            unpinned.append(pf)

    if args.keep_old_order:
        # Keep existing sequences where possible, sort the rest
        unpinned.sort(key=lambda x: (x['seq'] if x['seq'] is not None else float('inf')))

    used_seqs = {pf['new_seq'] for pf in pinned}
    
    current_seq = 0
    for pf in unpinned:
        while current_seq in used_seqs:
            current_seq += 1
        pf['new_seq'] = current_seq
        used_seqs.add(current_seq)
        current_seq += 1

    # Perform renames
    all_files = pinned + unpinned
    # Figure out max digits for padding
    max_seq = max([pf['new_seq'] for pf in all_files] + [0])
    digits = max(2, len(str(max_seq)))

    for pf in all_files:
        new_name = format_sequence(pf['new_seq'], pf['rest'], digits)
        if new_name != pf['original']:
            src_path = os.path.join(target_norm, pf['original'])
            dst_path = os.path.join(target_norm, new_name)
            rename_file(src_path, dst_path)

def main():
    parser = argparse.ArgumentParser(description="File Manipulator CLI (Lowercase & Fix Sequencing)")
    subparsers = parser.add_subparsers(dest='command', required=True)

    # Lowercase command
    lowercase_p = subparsers.add_parser('lowercase', help="Convert files to lowercase recursively")
    lowercase_p.add_argument('target_directory', help="Directory to process")
    lowercase_p.add_argument('--except', dest='except_list', help="Comma-separated list of patterns to ignore (e.g., 'docs/*, temp.md')")
    
    # Fix sequence command
    fixseq_p = subparsers.add_parser('fix-seq-files', help="Re-sequence files in a directory")
    fixseq_p.add_argument('target_directory', help="Directory to process")
    fixseq_p.add_argument('--order-by-time', action='store_true', help="Order by modification time")
    fixseq_p.add_argument('--order-by-az', action='store_true', help="Order alphabetically")
    fixseq_p.add_argument('--keep-old-order', action='store_true', help="Preserve existing sequences as much as possible")
    fixseq_p.add_argument('--pin', help="Pin specific files to specific sequence numbers (e.g., 'readme=00,intro=01')")

    args = parser.parse_args()

    if args.command == 'lowercase':
        cmd_lowercase(args)
    elif args.command == 'fix-seq-files':
        cmd_fix_seq(args)

if __name__ == '__main__':
    main()
