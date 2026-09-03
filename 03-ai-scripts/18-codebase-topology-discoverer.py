#!/usr/bin/env python3
"""
Universal Polyglot Codebase & Topology Discovery Engine
Automatically inspects any codebase (Go, Rust, Python, TypeScript, PHP, C#, SQL),
classifies subsystems (Backend, Database, Frontend, CI/CD, Docs), and maintains
a high-speed TTL-cached topology map in tmp/cache/paths/codebase-topology-cache.json.

All Enums, Cache Keys, Default Directories, Constants, and Functions are imported
directly from 02-shared-engine.py as the single source of truth.

Usage:
  python 03-ai-scripts/18-codebase-topology-discoverer.py [--summary]
  python 03-ai-scripts/18-codebase-topology-discoverer.py --query <subsystem-or-language>
  python 03-ai-scripts/18-codebase-topology-discoverer.py --refresh [--ttl <seconds>]
"""

import argparse
import datetime
from importlib import import_module
import json
import os
from pathlib import Path
import sys
import time
from typing import Any

if hasattr(sys.stdout, "reconfigure"):
    sys.stdout.reconfigure(encoding="utf-8")

sys.path.insert(0, str(Path(__file__).parent))
engine = import_module("02-shared-engine")

# Centralized Enums
LanguageType = engine.LanguageType
SubsystemType = engine.SubsystemType
EncodingType = engine.EncodingType
ExitCodeType = engine.ExitCodeType
CacheKeyType = engine.CacheKeyType

# Centralized Configurations & Manifests
LANGUAGE_MANIFESTS = engine.LANGUAGE_MANIFESTS
SUBSYSTEM_DIR_HINTS = engine.SUBSYSTEM_DIR_HINTS
SUBSYSTEM_ENTRYPOINTS = engine.SUBSYSTEM_ENTRYPOINTS
QUERY_ALIASES = engine.QUERY_ALIASES
LANG_EXT_MAP = engine.LANG_EXT_MAP
PRIMARY_TOPOLOGY_CACHE_FILE = engine.PRIMARY_TOPOLOGY_CACHE_FILE
LEGACY_TOPOLOGY_CACHE_FILE = engine.LEGACY_TOPOLOGY_CACHE_FILE
DEFAULT_TTL_SECONDS = engine.DEFAULT_TTL_SECONDS
DEFAULT_ENCODING = engine.DEFAULT_ENCODING
LINE_SEPARATOR = engine.LINE_SEPARATOR
PATH_SEPARATOR = engine.PATH_SEPARATOR
CURRENT_DIR = engine.CURRENT_DIR
EMPTY_STRING = engine.EMPTY_STRING
CACHE_PATHS_DIR = engine.CACHE_PATHS_DIR

# Centralized Cache Keys
CACHE_KEY_FILES = engine.CACHE_KEY_FILES
CACHE_KEY_TOTAL_FILES = engine.CACHE_KEY_TOTAL_FILES
CACHE_KEY_VERSION = engine.CACHE_KEY_VERSION
CACHE_KEY_GENERATED_AT = engine.CACHE_KEY_GENERATED_AT
CACHE_KEY_EXPIRES_AT = engine.CACHE_KEY_EXPIRES_AT
CACHE_KEY_TTL_SECONDS = engine.CACHE_KEY_TTL_SECONDS
CACHE_KEY_SCAN_DURATION_MS = engine.CACHE_KEY_SCAN_DURATION_MS
CACHE_KEY_ROOT_PATH = engine.CACHE_KEY_ROOT_PATH
CACHE_KEY_MANIFESTS = engine.CACHE_KEY_MANIFESTS
CACHE_KEY_LANGUAGE_DISTRIBUTION = engine.CACHE_KEY_LANGUAGE_DISTRIBUTION
CACHE_KEY_LANGUAGE_ROOTS = engine.CACHE_KEY_LANGUAGE_ROOTS
CACHE_KEY_SUBSYSTEMS = engine.CACHE_KEY_SUBSYSTEMS
CACHE_KEY_ROOTS = engine.CACHE_KEY_ROOTS
CACHE_KEY_ENTRYPOINTS = engine.CACHE_KEY_ENTRYPOINTS
CACHE_KEY_SCHEMA_FILES = engine.CACHE_KEY_SCHEMA_FILES
CACHE_KEY_WORKFLOWS = engine.CACHE_KEY_WORKFLOWS
CACHE_KEY_SPEC_ROOTS = engine.CACHE_KEY_SPEC_ROOTS
CACHE_KEY_TEST_RUNNERS = engine.CACHE_KEY_TEST_RUNNERS

# Centralized Utility Functions
process_repository_files = engine.process_repository_files
normalize_rel_path = engine.normalize_rel_path
is_ignored_directory = engine.is_ignored_directory
is_binary_file = engine.is_binary_file
atomic_cache_lock = engine.atomic_cache_lock

# Build Reverse Lookup for Extensions -> LanguageType from Centralized LANG_EXT_MAP
EXT_TO_LANGUAGE_MAP: dict[str, LanguageType] = {}
for lang_key, exts in LANG_EXT_MAP.items():
    matching_enum = getattr(LanguageType, lang_key.upper(), None)
    if matching_enum:
        for ext in exts:
            EXT_TO_LANGUAGE_MAP[ext.lower()] = matching_enum

# --- Dynamic Topology Discovery Logic ---

def match_language_manifest(filename: str) -> LanguageType | None:
    """Matches a filename against centralized language manifests dynamically."""
    f_lower = filename.lower()
    for lang, manifests in LANGUAGE_MANIFESTS.items():
        for pattern in manifests:
            if pattern.startswith("*."):
                if f_lower.endswith(pattern[1:]):
                    return lang
            elif f_lower == pattern.lower():
                return lang
    return None

def detect_manifests(root_dir: str = CURRENT_DIR) -> dict[str, list[str]]:
    """Detects top-level and submodule package manifests using centralized rules."""
    detected: dict[str, list[str]] = {lang.value: [] for lang in LanguageType}

    for root, dirs, files in os.walk(root_dir):
        dirs[:] = [d for d in dirs if not is_ignored_directory(d)]
        for f in files:
            matched_lang = match_language_manifest(f)
            has_match = bool(matched_lang)
            if has_match:
                norm_rel = normalize_rel_path(os.path.join(root, f))
                detected[matched_lang.value].append(norm_rel)

    return {k: v for k, v in detected.items() if v}

def classify_codebase_subsystems(root_dir: str = CURRENT_DIR) -> dict[str, dict[str, Any]]:
    """Scans and categorizes directory roots into functional subsystems dynamically."""
    subsystems: dict[str, dict[str, Any]] = {
        SubsystemType.BACKEND.value: {CACHE_KEY_ROOTS: set(), CACHE_KEY_ENTRYPOINTS: []},
        SubsystemType.DATABASE.value: {CACHE_KEY_ROOTS: set(), CACHE_KEY_SCHEMA_FILES: []},
        SubsystemType.FRONTEND.value: {CACHE_KEY_ROOTS: set(), CACHE_KEY_ENTRYPOINTS: []},
        SubsystemType.CICD.value: {CACHE_KEY_ROOTS: set(), CACHE_KEY_WORKFLOWS: []},
        SubsystemType.DOCS.value: {CACHE_KEY_ROOTS: set(), CACHE_KEY_SPEC_ROOTS: []},
        SubsystemType.CLI.value: {CACHE_KEY_ROOTS: set(), CACHE_KEY_ENTRYPOINTS: []},
        SubsystemType.TESTS.value: {CACHE_KEY_ROOTS: set(), CACHE_KEY_TEST_RUNNERS: []},
    }

    hint_to_subsystems: dict[str, set[str]] = {}
    for subsys_enum, hints in SUBSYSTEM_DIR_HINTS.items():
        for h in hints:
            h_clean = h.lower()
            if h_clean not in hint_to_subsystems:
                hint_to_subsystems[h_clean] = set()
            hint_to_subsystems[h_clean].add(subsys_enum.value)

    for root, dirs, files in os.walk(root_dir):
        dirs[:] = [d for d in dirs if not is_ignored_directory(d)]
        norm_dir = normalize_rel_path(root)
        dir_parts_lower = {p.lower() for p in Path(root).parts}

        for part in dir_parts_lower:
            matched_subsystems = hint_to_subsystems.get(part)
            has_matches = bool(matched_subsystems)
            if has_matches:
                for subsys_name in matched_subsystems:
                    subsystems[subsys_name][CACHE_KEY_ROOTS].add(norm_dir)

        for f in files:
            norm_file = normalize_rel_path(os.path.join(root, f))
            ext = os.path.splitext(f)[1].lower()
            f_lower = f.lower()

            is_sql_schema = (ext == ".sql" or "schema" in f_lower or "migration" in f_lower)
            if is_sql_schema:
                subsystems[SubsystemType.DATABASE.value][CACHE_KEY_SCHEMA_FILES].append(norm_file)

            is_entrypoint = (f in SUBSYSTEM_ENTRYPOINTS)
            if is_entrypoint:
                subsystems[SubsystemType.BACKEND.value][CACHE_KEY_ENTRYPOINTS].append(norm_file)

            is_cli = ("cli" in f_lower or f.startswith("cmd") or "command" in f_lower)
            if is_cli:
                subsystems[SubsystemType.CLI.value][CACHE_KEY_ENTRYPOINTS].append(norm_file)

            is_test_file = ("test" in f_lower or f.startswith("test_") or f.endswith("_test.go") or f.endswith(".test.ts"))
            if is_test_file:
                subsystems[SubsystemType.TESTS.value][CACHE_KEY_TEST_RUNNERS].append(norm_file)

            is_workflow = (ext in {".yml", ".yaml"} and ".github" in norm_file)
            if is_workflow:
                subsystems[SubsystemType.CICD.value][CACHE_KEY_WORKFLOWS].append(norm_file)

    serialized: dict[str, dict[str, Any]] = {}
    for st_name, data in subsystems.items():
        serialized[st_name] = {
            k: sorted(list(v)) if isinstance(v, set) else sorted(v)
            for k, v in data.items()
        }
    return serialized

def build_topology_map(root_dir: str = CURRENT_DIR, ttl_seconds: int = DEFAULT_TTL_SECONDS) -> dict[str, Any]:
    """Builds complete topology map with timestamps and TTL expiry information."""
    start_time = time.perf_counter()
    now_utc = datetime.datetime.now(datetime.timezone.utc)
    expires_at = now_utc + datetime.timedelta(seconds=ttl_seconds)

    manifests = detect_manifests(root_dir=root_dir)
    subsystems = classify_codebase_subsystems(root_dir=root_dir)

    lang_file_counts: dict[str, int] = {}
    lang_file_roots: dict[str, set[str]] = {}
    total_files = 0

    for root, dirs, files in os.walk(root_dir):
        dirs[:] = [d for d in dirs if not is_ignored_directory(d)]
        norm_dir = normalize_rel_path(root)
        for f in files:
            total_files += 1
            ext = os.path.splitext(f)[1].lower()

            matched_lang = EXT_TO_LANGUAGE_MAP.get(ext)
            has_lang = bool(matched_lang)
            if has_lang:
                lang_str = matched_lang.value
                lang_file_counts[lang_str] = lang_file_counts.get(lang_str, 0) + 1
                if lang_str not in lang_file_roots:
                    lang_file_roots[lang_str] = set()
                lang_file_roots[lang_str].add(norm_dir)

    elapsed_ms = (time.perf_counter() - start_time) * 1000

    return {
        CACHE_KEY_VERSION: "1.0.0",
        CACHE_KEY_GENERATED_AT: now_utc.isoformat(),
        CACHE_KEY_EXPIRES_AT: expires_at.isoformat(),
        CACHE_KEY_TTL_SECONDS: ttl_seconds,
        CACHE_KEY_SCAN_DURATION_MS: round(elapsed_ms, 2),
        CACHE_KEY_TOTAL_FILES: total_files,
        CACHE_KEY_ROOT_PATH: normalize_rel_path(root_dir),
        CACHE_KEY_MANIFESTS: manifests,
        CACHE_KEY_LANGUAGE_DISTRIBUTION: dict(sorted(lang_file_counts.items(), key=lambda x: x[1], reverse=True)),
        CACHE_KEY_LANGUAGE_ROOTS: {k: sorted(list(v)) for k, v in lang_file_roots.items()},
        CACHE_KEY_SUBSYSTEMS: subsystems,
    }

# --- TTL Cache Management ---

def load_cached_topology() -> dict[str, Any] | None:
    """Loads cached topology and checks if TTL is still valid."""
    for cache_p in (PRIMARY_TOPOLOGY_CACHE_FILE, LEGACY_TOPOLOGY_CACHE_FILE):
        if cache_p.exists():
            try:
                with open(cache_p, "r", encoding=DEFAULT_ENCODING) as f:
                    data = json.load(f)
                if isinstance(data, dict):
                    expires_str = data.get(CACHE_KEY_EXPIRES_AT)
                    if expires_str:
                        expires_at = datetime.datetime.fromisoformat(expires_str)
                        now_utc = datetime.datetime.now(datetime.timezone.utc)
                        is_valid = (now_utc < expires_at)
                        if is_valid:
                            return data
            except Exception:
                pass
    return None

def save_topology_cache(topology_data: dict[str, Any]) -> None:
    """Saves topology data to primary and legacy cache locations with atomic locking."""
    CACHE_PATHS_DIR.mkdir(parents=True, exist_ok=True)
    for target in (PRIMARY_TOPOLOGY_CACHE_FILE, LEGACY_TOPOLOGY_CACHE_FILE):
        try:
            temp_file = target.with_suffix(".json.tmp")
            with open(temp_file, "w", encoding=DEFAULT_ENCODING) as f:
                json.dump(topology_data, f, indent=2)
            temp_file.replace(target)
        except Exception:
            pass

def get_or_create_topology(
    root_dir: str = CURRENT_DIR,
    is_force_refresh: bool = False,
    ttl_seconds: int = DEFAULT_TTL_SECONDS
) -> dict[str, Any]:
    """Retrieves cached topology if valid, or regenerates and stores a new cache."""
    if not is_force_refresh:
        cached = load_cached_topology()
        has_valid_cache = bool(cached)
        if has_valid_cache:
            return cached

    topology = build_topology_map(root_dir=root_dir, ttl_seconds=ttl_seconds)
    save_topology_cache(topology)
    return topology

# --- Query & Routing Functions ---

def query_topology(query_term: str, topology: dict[str, Any]) -> None:
    """Searches topology for language or subsystem matches and outputs routing paths."""
    q_clean = query_term.strip().lower()
    subsystems = topology.get(CACHE_KEY_SUBSYSTEMS, {})
    manifests = topology.get(CACHE_KEY_MANIFESTS, {})
    lang_dist = topology.get(CACHE_KEY_LANGUAGE_DISTRIBUTION, {})
    lang_roots = topology.get(CACHE_KEY_LANGUAGE_ROOTS, {})

    print(f"{LINE_SEPARATOR}⚡ Topology Routing Query for: `{query_term}`{LINE_SEPARATOR}")
    found = False

    matched_target = QUERY_ALIASES.get(q_clean)

    # 1. Match Subsystems
    for subsys_name, info in subsystems.items():
        is_subsys_match = (
            q_clean in subsys_name.lower() or
            (matched_target and isinstance(matched_target, SubsystemType) and matched_target.value == subsys_name)
        )
        if is_subsys_match:
            found = True
            print(f"📦 Subsystem: [{subsys_name}]")
            roots = info.get(CACHE_KEY_ROOTS, [])
            has_roots = bool(roots)
            if has_roots:
                print("   📁 Directory Roots:")
                for r in roots[:15]:
                    print(f"      • {r}")
                if len(roots) > 15:
                    print(f"      ... and {len(roots) - 15} more roots.")
            for k, v in info.items():
                if k != CACHE_KEY_ROOTS and isinstance(v, list) and v:
                    print(f"   📄 {k}:")
                    for item in v[:10]:
                        print(f"      • {item}")
                    if len(v) > 10:
                        print(f"      ... and {len(v) - 10} more files.")
            print()

    # 2. Match Languages
    for lang_name, count in lang_dist.items():
        is_lang_match = (
            q_clean in lang_name.lower() or
            (matched_target and isinstance(matched_target, LanguageType) and matched_target.value == lang_name)
        )
        if is_lang_match:
            found = True
            print(f"🔤 Language: [{lang_name}] — {count} tracked files")
            man = manifests.get(lang_name, [])
            has_man = bool(man)
            if has_man:
                print("   📋 Manifests & Configs:")
                for m in man:
                    print(f"      • {m}")
            roots = lang_roots.get(lang_name, [])
            has_roots = bool(roots)
            if has_roots:
                print("   📁 Primary Directories:")
                for r in roots[:12]:
                    print(f"      • {r}")
                if len(roots) > 12:
                    print(f"      ... and {len(roots) - 12} more folders.")
            print()

    if not found:
        print(f"⚠️ No direct subsystem or language matches for '{query_term}'.")
        print("💡 Available Subsystems: BACKEND, DATABASE, FRONTEND, CICD, DOCS, CLI, TESTS")
        print(f"💡 Available Languages: {', '.join(lang_dist.keys())}")

def print_topology_summary(topology: dict[str, Any]) -> None:
    """Prints a clear, high-density terminal summary of the codebase topology."""
    print("================================================================================")
    print(f"🌐 Polyglot Codebase Topology Map (Generated: {topology[CACHE_KEY_GENERATED_AT]})")
    print(f"⏱️ TTL Expiry: {topology[CACHE_KEY_EXPIRES_AT]} | Scan Time: {topology[CACHE_KEY_SCAN_DURATION_MS]}ms")
    print("================================================================================")

    print(f"{LINE_SEPARATOR}📊 Polyglot Language Breakdown:")
    for lang, count in topology.get(CACHE_KEY_LANGUAGE_DISTRIBUTION, {}).items():
        man = topology.get(CACHE_KEY_MANIFESTS, {}).get(lang, [])
        man_str = f" (Manifests: {', '.join(man[:2])})" if man else EMPTY_STRING
        print(f"   • {lang:<14} : {count:>5} files{man_str}")

    print(f"{LINE_SEPARATOR}🏛️ Subsystem Map & Navigation Roots:")
    subsystems = topology.get(CACHE_KEY_SUBSYSTEMS, {})
    for subsys_name, data in subsystems.items():
        roots = data.get(CACHE_KEY_ROOTS, [])
        if roots:
            display_roots = ", ".join(roots[:4])
            if len(roots) > 4:
                display_roots += f" (+{len(roots)-4} more)"
            print(f"   • {subsys_name:<10} : {display_roots}")
        else:
            print(f"   • {subsys_name:<10} : (none detected)")

    print(f"{LINE_SEPARATOR}💡 AI Navigation Tip:")
    print("   Query specific folders instantly using:")
    print("   `python 03-ai-scripts/18-codebase-topology-discoverer.py --query <go|python|db|backend>`")
    print("================================================================================")

# --- CLI Entry Point ---

def main():
    parser = argparse.ArgumentParser(
        description="Universal Polyglot Codebase & Topology Discovery Engine",
        formatter_class=argparse.RawDescriptionHelpFormatter
    )
    parser.add_argument("--query", "-q", help="Search subsystem or language routing (e.g. go, rust, db, backend)")
    parser.add_argument("--summary", "-s", action="store_true", help="Print topology summary breakdown")
    parser.add_argument("--refresh", "-r", action="store_true", help="Force regenerate topology cache")
    parser.add_argument("--ttl", type=int, default=DEFAULT_TTL_SECONDS, help="TTL cache duration in seconds (default: 1800)")
    parser.add_argument("--json", "-j", action="store_true", help="Output raw JSON topology map")
    parser.add_argument("--path", "-p", default=CURRENT_DIR, help="Root directory to discover (default: .)")
    args = parser.parse_args()

    topology = get_or_create_topology(
        root_dir=args.path,
        is_force_refresh=args.refresh,
        ttl_seconds=args.ttl
    )

    has_json = args.json
    if has_json:
        print(json.dumps(topology, indent=2))
        sys.exit(0)

    has_query = bool(args.query)
    if has_query:
        query_topology(args.query, topology)
        sys.exit(0)

    print_topology_summary(topology)
    sys.exit(0)

if __name__ == "__main__":
    main()
