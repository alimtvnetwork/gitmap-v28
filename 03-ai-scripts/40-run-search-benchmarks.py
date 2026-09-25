#!/usr/bin/env python3
"""Search and File Discovery Performance Benchmark Suite.

Benchmarks GitMap Native AUM vs. Ripgrep vs. Python vs. PowerShell across:
1. Universal Wildcard File Search (*test*.md)
2. Complex Regex / Content Search (appfault.AppError)
3. Zero-Disk File Streaming / Reading (readme.md)

Outputs exact measured timings, speedup ratios, and formatted Markdown tables.
"""
from __future__ import annotations

import os
from pathlib import Path
import statistics
import subprocess
import sys
import time

REPO_ROOT = Path(__file__).resolve().parent.parent


def run_benchmark(cmd: list[str], cwd: Path, iterations: int = 3) -> dict[str, any]:
    """Run command multiple times and measure elapsed milliseconds."""
    durations: list[float] = []
    output_sample = ""
    exit_code = 0

    for i in range(iterations):
        start = time.perf_counter()
        res = subprocess.run(
            cmd,
            cwd=str(cwd),
            capture_output=True,
            text=True,
            encoding="utf-8",
            errors="replace",
        )
        elapsed_ms = (time.perf_counter() - start) * 1000.0
        durations.append(elapsed_ms)
        if i == 0:
            output_sample = res.stdout.strip()
            exit_code = res.returncode

    avg_ms = statistics.mean(durations)
    median_ms = statistics.median(durations)
    min_ms = min(durations)

    return {
        "avg_ms": avg_ms,
        "median_ms": median_ms,
        "min_ms": min_ms,
        "exit_code": exit_code,
        "output_len": len(output_sample.splitlines()) if output_sample else 0,
        "raw_durations": durations,
    }


def main() -> None:
    print(f"Running benchmarks on repository: {REPO_ROOT}", flush=True)
    print("Warmup GitMap...", flush=True)
    subprocess.run(["gitmap", "lf"], cwd=str(REPO_ROOT), capture_output=True)

    results = {}

    # --- 1. File Search (*test*.md) ---
    print("\n[1/3] Benchmarking File Search (*test*.md)...", flush=True)
    results["file_search"] = {}
    
    print("  -> GitMap find...", flush=True)
    results["file_search"]["GitMap Native (`gitmap find`)"] = run_benchmark(["gitmap", "find", "*test*", "-ext", "md"], REPO_ROOT, iterations=3)
    
    print("  -> Ripgrep files...", flush=True)
    results["file_search"]["Ripgrep (`rg --files`)"] = run_benchmark(["rg", "--files", "-g", "*test*.md"], REPO_ROOT, iterations=3)
    
    print("  -> Python fast-file-scanner...", flush=True)
    results["file_search"]["Python (`11-fast-file-scanner.py`)"] = run_benchmark([sys.executable, "03-ai-scripts/11-fast-file-scanner.py", "--search", "test", "--limit", "100"], REPO_ROOT, iterations=3)
    
    print("  -> PowerShell Get-ChildItem...", flush=True)
    results["file_search"]["PowerShell (`Get-ChildItem`)"] = run_benchmark(["pwsh", "-NoProfile", "-Command", "$null = (Get-ChildItem -Recurse -File -Filter '*test*.md' -Exclude '.git','node_modules')"], REPO_ROOT, iterations=3)

    # --- 2. Content / Regex Search (appfault.AppError) ---
    print("\n[2/3] Benchmarking Content / Regex Search (appfault.AppError)...", flush=True)
    results["content_search"] = {}
    
    print("  -> GitMap AUM Hot-Cache...", flush=True)
    results["content_search"]["GitMap Hot-Cache (`DH2D` SQLite)"] = run_benchmark(["gitmap", "search", "AppError"], REPO_ROOT, iterations=3)
    
    print("  -> Ripgrep regex...", flush=True)
    results["content_search"]["Ripgrep (`rg`)"] = run_benchmark(["rg", "appfault\\.AppError", "."], REPO_ROOT, iterations=3)
    
    print("  -> Python fast-cached-grep...", flush=True)
    results["content_search"]["Python (`12-fast-cached-grep.py`)"] = run_benchmark([sys.executable, "03-ai-scripts/12-fast-cached-grep.py", "--pattern", "appfault\\.AppError", "--limit", "50"], REPO_ROOT, iterations=1)
    
    print("  -> PowerShell Select-String...", flush=True)
    results["content_search"]["PowerShell (`Select-String`)"] = run_benchmark(["pwsh", "-NoProfile", "-Command", "$null = (Get-ChildItem -Path 01-prompts,04-code -Recurse -File -Filter *.go | Select-String -Pattern 'appfault\\.AppError')"], REPO_ROOT, iterations=2)

    # --- 3. File Streaming / Cat (readme.md) ---
    print("\n[3/3] Benchmarking File Streaming / Cat (readme.md)...", flush=True)
    results["file_cat"] = {}
    
    print("  -> GitMap cat...", flush=True)
    results["file_cat"]["GitMap Streaming (`gitmap cat`)"] = run_benchmark(["gitmap", "cat", "readme.md"], REPO_ROOT, iterations=3)
    
    print("  -> Ripgrep streaming...", flush=True)
    results["file_cat"]["Ripgrep (`rg ^`)"] = run_benchmark(["rg", "^", "readme.md"], REPO_ROOT, iterations=3)
    
    print("  -> Python fast-file-reader...", flush=True)
    results["file_cat"]["Python (`17-fast-file-reader.py`)"] = run_benchmark([sys.executable, "03-ai-scripts/17-fast-file-reader.py", "--file", "readme.md", "--limit", "1000"], REPO_ROOT, iterations=3)
    
    print("  -> PowerShell Get-Content...", flush=True)
    results["file_cat"]["PowerShell (`Get-Content`)"] = run_benchmark(["pwsh", "-NoProfile", "-Command", "$null = (Get-Content readme.md)"], REPO_ROOT, iterations=3)

    print("\n================ BENCHMARK RESULTS SUMMARY ================\n", flush=True)

    for category, tools in results.items():
        print(f"### {category.replace('_', ' ').title()}\n", flush=True)
        print("| Tool / Engine | Command / Method | Measured Latency (Avg) | Min Latency | Speedup vs PS |", flush=True)
        print("| :--- | :--- | :--- | :--- | :--- |", flush=True)
        ps_avg = [v["avg_ms"] for k, v in tools.items() if "PowerShell" in k][0]
        for tool, data in tools.items():
            speedup = (ps_avg / data['avg_ms']) if data['avg_ms'] > 0 else 0
            speedup_str = f"**{speedup:.1f}x faster**" if speedup > 1.05 else ("Baseline (1.0x)" if "PowerShell" in tool else f"{speedup:.2f}x")
            print(f"| {tool} | `{data['raw_durations']}` | **{data['avg_ms']:.2f} ms** | {data['min_ms']:.2f} ms | {speedup_str} |", flush=True)
        print("\n", flush=True)


if __name__ == "__main__":
    main()
