#!/usr/bin/env python3
import argparse
import asyncio
import os
import sys
import time

if hasattr(sys.stdout, "reconfigure"):
    sys.stdout.reconfigure(encoding="utf-8")
if hasattr(sys.stderr, "reconfigure"):
    sys.stderr.reconfigure(encoding="utf-8")

# ─────────────────────────────────────────────────────────────────────────────
# Test Definitions
# ─────────────────────────────────────────────────────────────────────────────

INDEPENDENT_TESTS_PART1 = [
    (["version"], [0], "gitmap v", "Version command"),
    (["v"], [0], "gitmap v", "Version alias v"),
    (["help"], [0], "Usage: gitmap", "Help command"),
    (["--help"], [0], "Usage: gitmap", "--help flag"),
    (["-h"], [0], "Usage: gitmap", "-h flag"),
    (["llm"], [0], "Gitmap LLM Specification", "llm command"),
    (["llm-docs", "--stdout"], [0], "LLM.md", "llm-docs --stdout"),
    (["find"], [0], "Usage:", "find zero-args help"),
    (["find", "*version*.json"], [0], "", "find with wildcard pattern"),
    (["f"], [0], "Usage:", "f zero-args alias help"),
    (["find-files"], [0], "Usage:", "find-files zero-args help"),
    (["find-files", "version.json"], [0], "", "find-files exact query"),
    (["ff"], [0], "Usage:", "ff zero-args alias help"),
    (["find-files-any"], [0], "Usage:", "find-files-any zero-args help"),
    (["find-files-any", "smoke", "-ext", "py"], [0], "", "find-files-any with -ext filter"),
    (["ffa"], [0], "Usage:", "ffa zero-args alias help"),
    (["find-files-startswith"], [0], "Usage:", "find-files-startswith zero-args help"),
    (["find-files-startswith", "version"], [0], "", "find-files-startswith prefix query"),
    (["ffs"], [0], "Usage:", "ffs zero-args alias help"),
    (["find-files-endswith"], [0], "Usage:", "find-files-endswith zero-args help"),
    (["find-files-endswith", ".json"], [0], "", "find-files-endswith suffix query"),
    (["ffe"], [0], "Usage:", "ffe zero-args alias help"),
    (["list-files", "--help"], [0], "", "list-files --help"),
    (["search"], [0], "Usage:", "search zero-args help"),
    (["search", "--help"], [0], "", "search --help"),
    (["search", "Version"], [0, 1], "", "search query"),
    (["repo-search"], [0], "Usage:", "repo-search zero-args help"),
    (["repo-search-json", "--help"], [0], "", "repo-search-json --help"),
    (["find-regex"], [0], "Usage:", "find-regex zero-args help"),
    (["find-read"], [0], "Usage:", "find-read zero-args help"),
    (["find-regex-read"], [0], "Usage:", "find-regex-read zero-args help"),
    (["doctor"], [0, 1], "gitmap doctor", "doctor command"),
    (["pipeline", "status", "--json"], [0], "isRunning", "pipeline status --json"),
    (["pipeline", "eta"], [0], "", "pipeline eta"),
    (["pipeline-ai", "--help"], [0], "Usage", "pipeline-ai --help"),
    (["pipeline-ai", "status", "--json"], [0], "isRunning", "pipeline-ai status --json"),
    (["pipeline-ai", "status", "-t", "25", "--json"], [0], "isRunning", "pipeline-ai status with -t"),
    (["pipeline-ai", "eta", "--json"], [0], "isRunning", "pipeline-ai eta --json"),
    (["pl-ai", "status", "--json"], [0], "isRunning", "pl-ai alias"),
    (["eta"], [0], "", "eta alias"),
    (["waittime"], [0], "", "waittime alias"),
    (["pipeline", "error-logs", "--help"], [0], "", "pipeline error-logs --help"),
    (["pipeline", "logs", "--help"], [0], "", "pipeline logs --help"),
    (["status"], [0], "", "status command"),
    (["st"], [0], "", "status alias st"),
    (["antigravity", "stats"], [0], "Account: Default", "antigravity stats"),
    (["agy", "stats"], [0], "Account: Default", "agy stats alias"),
    (["ag", "stats"], [0], "Account: Default", "ag stats alias"),
    (["agy", "ls"], [0], "", "agy ls command"),
    (["vscode", "ls"], [0], "", "vscode ls command"),
    (["vsc", "ls"], [0], "", "vsc ls alias"),
]

SCHEDULE_CHAIN = [
    (["schedule", "--help"], [0], "", "schedule --help"),
    (["schedule", "list", "--json"], [0], "", "schedule list --json"),
    (["schedule", "add", "smoke-job", "echo smoke", "--every", "1h"], [0], "", "schedule add command"),
    (["schedule", "status", "smoke-job", "--json"], [0], "smoke-job", "schedule status --json"),
    (["schedule", "run", "smoke-job"], [0], "", "schedule run command"),
    (["schedule", "log", "smoke-job", "--json"], [0], "runnerUser", "schedule log --json"),
    (["schedule", "export-all", "--json"], [0], "", "schedule export-all --json"),
    (["schedule", "disable", "smoke-job"], [0], "", "schedule disable command"),
    (["schedule", "enable", "smoke-job"], [0], "", "schedule enable command"),
    (["schedule", "reset", "smoke-job"], [0], "", "schedule reset command"),
    (["schedule", "rm", "smoke-job"], [0], "", "schedule rm command"),
]

MACRO_CHAIN = [
    (["macro", "--help"], [0], "", "macro --help"),
    (["macro", "add", "smoke-macro", "echo smoke-macro"], [0], "", "macro add command"),
    (["macro", "list", "--json"], [0], "", "macro list --json"),
    (["macro", "rm", "smoke-macro"], [0], "", "macro rm command"),
]

INDEPENDENT_TESTS_PART2 = [
    (["retry", "--sleep=10ms", "--max-retries=1", "echo retry-smoke-test"], [0], "retry-smoke-test", "retry command"),
    (["chrome", "--help"], [0], "", "chrome --help"),
    (["replace", "history"], [0], "", "replace history command"),
    (["go-repos", "--json"], [0, 1], "", "go-repos --json"),
    (["gr", "--json"], [0, 1], "", "gr --json alias"),
    (["node-repos", "--json"], [0, 1], "", "node-repos --json"),
    (["react-repos", "--json"], [0, 1], "", "react-repos --json"),
    (["cpp-repos", "--json"], [0, 1], "", "cpp-repos --json"),
    (["csharp-repos", "--json"], [0, 1], "", "csharp-repos --json"),
    (["error", "ls"], [0], "", "error ls command"),
    (["error", "warnings"], [0], "", "error warnings command"),
    (["completion", "bash"], [0], "", "completion bash"),
    (["completion", "powershell"], [0], "", "completion powershell"),
    (["completion", "zsh"], [0], "", "completion zsh"),
    (["diff-profiles", "--help"], [0], "", "diff-profiles --help"),
    (["zip-group", "--help"], [0], "", "zip-group --help"),
    (["env", "--help"], [0], "", "env --help"),
    (["task", "--help"], [0], "", "task --help"),
    (["amend", "--help"], [0], "", "amend --help"),
    (["amend-list", "--help"], [0], "", "amend-list --help"),
    (["seo-write", "--help"], [0], "", "seo-write --help"),
    (["fix-repo", "--help"], [0], "", "fix-repo --help"),
    (["fix-git", "--help"], [0], "", "fix-git --help"),
    (["fix-git", "--json"], [0], "", "fix-git --json"),
    (["fg", "--help"], [0], "", "fg --help alias"),
    (["gomod", "--help"], [0], "", "gomod --help"),
    (["serve", "--help"], [0], "", "serve --help"),
    (["open", "--help"], [0], "", "open --help"),
    (["make-public", "--help"], [0], "", "make-public --help"),
    (["make-private", "--help"], [0], "", "make-private --help"),
    (["make-all-public", "--help"], [0], "", "make-all-public --help"),
    (["make-all-private", "--help"], [0], "", "make-all-private --help"),
    (["visibility-history", "--help"], [0], "", "visibility-history --help"),
    (["visibility-undo", "--help"], [0], "", "visibility-undo --help"),
    (["profiles", "--help"], [0], "", "profiles --help"),
    (["profiles", "ls", "--json"], [0], "", "profiles ls --json"),
    (["profiles", "status"], [0], "", "profiles status"),
    (["backup", "--help"], [0], "", "backup --help"),
    (["backup", "status"], [0], "", "backup status"),
    (["backup", "ls"], [0], "", "backup ls"),
    (["agy", "group", "--help"], [0], "", "agy group --help"),
    (["agy", "group", "ls"], [0], "", "agy group ls"),
    (["agy", "undo", "--help"], [0], "", "agy undo --help"),
    (["agy", "redo", "--help"], [0], "", "agy redo --help"),
    (["agy", "plugin", "ls"], [0], "", "agy plugin ls"),
    (["agy", "settings", "--help"], [0], "", "agy settings --help"),
    (["prompt", "--help"], [0], "", "prompt --help"),
    (["prompt", "ls"], [0], "", "prompt ls"),
    (["chrome", "undo", "--help"], [0], "", "chrome undo --help"),
    (["vscode", "group", "ls"], [0], "", "vscode group ls"),
    (["gd", "group", "ls"], [0], "", "gd group ls"),
    (["installer", "undo", "--help"], [0], "", "installer undo --help"),
]

# ─────────────────────────────────────────────────────────────────────────────
# Execution Engine
# ─────────────────────────────────────────────────────────────────────────────

async def run_cmd_async(bin_path: str, args: list[str], cwd: str | None = None) -> tuple[int, str, str]:
    env = os.environ.copy()
    env["GITMAP_SKIP_DELAY"] = "1"
    env["GITMAP_NON_INTERACTIVE"] = "1"
    env["CI"] = "1"
    try:
        proc = await asyncio.create_subprocess_exec(
            bin_path,
            *args,
            cwd=cwd,
            env=env,
            stdout=asyncio.subprocess.PIPE,
            stderr=asyncio.subprocess.PIPE,
        )
        stdout_bytes, stderr_bytes = await asyncio.wait_for(proc.communicate(), timeout=35.0)
        stdout = stdout_bytes.decode("utf-8", errors="replace")
        stderr = stderr_bytes.decode("utf-8", errors="replace")
        retcode = proc.returncode
        if retcode is None:
            retcode = 0
        return retcode, stdout, stderr
    except asyncio.TimeoutError:
        try:
            proc.kill()
        except ProcessLookupError:
            pass
        return 124, "", "Error: Command timed out after 35 seconds"
    except Exception as exc:
        return 1, "", f"Error executing command: {exc}"


async def execute_single_test_async(
    bin_path: str, repo_root: str, idx: int, test_item: tuple[list[str], list[int], str, str]
) -> tuple[int, bool, str]:
    args, valid_codes, exp_sub, desc = test_item
    code, stdout, stderr = await run_cmd_async(bin_path, args, cwd=repo_root)
    out = stdout + stderr
    if code not in valid_codes:
        msg = f"❌ FAIL: {desc} (gitmap {' '.join(args)}) exited with {code}, expected {valid_codes}. Output:\n{out.strip()}"
        return idx, False, msg
    if exp_sub and exp_sub not in out:
        msg = f"❌ FAIL: {desc} (missing substring {exp_sub}). Output:\n{out.strip()}"
        return idx, False, msg
    return idx, True, f"  ✓ PASS: {desc}"


async def process_queue_item(
    task_item,
    bin_path: str,
    repo_root: str,
    results_map: dict[int, tuple[bool, str]],
) -> None:
    if isinstance(task_item, list):
        for idx, item in task_item:
            res = await execute_single_test_async(bin_path, repo_root, idx, item)
            results_map[res[0]] = (res[1], res[2])
        return

    idx, item = task_item
    res = await execute_single_test_async(bin_path, repo_root, idx, item)
    results_map[res[0]] = (res[1], res[2])


async def async_worker(
    queue: asyncio.Queue,
    bin_path: str,
    repo_root: str,
    results_map: dict[int, tuple[bool, str]],
) -> None:
    while True:
        task_item = await queue.get()
        if task_item is None:
            queue.task_done()
            break

        await process_queue_item(task_item, bin_path, repo_root, results_map)
        queue.task_done()


def resolve_binary_path(raw_path: str | None, repo_root: str) -> str:
    bin_name = "gitmap.exe" if os.name == "nt" else "gitmap"
    candidates: list[str] = []
    if raw_path:
        candidates.append(os.path.abspath(raw_path))
        candidates.append(os.path.abspath(os.path.join(repo_root, raw_path)))
    candidates.append(os.path.abspath(os.path.join(repo_root, "bin", bin_name)))
    candidates.append(os.path.abspath(os.path.join(repo_root, bin_name)))

    for path in candidates:
        if os.path.isfile(path):
            return path

    searched = "\n  - ".join(candidates)
    print(f"Gitmap binary not found. Searched locations:\n  - {searched}", file=sys.stderr)
    sys.exit(1)


def get_worker_count(requested: int | None) -> int:
    if requested is not None and requested > 0:
        return requested
    cpu = os.cpu_count() or 4
    return min(24, max(8, cpu * 2))


def print_all_results(
    results_map: dict[int, tuple[bool, str]],
    total: int,
    passed: int,
    failed: int,
    elapsed: float,
    worker_count: int,
    bin_path: str,
) -> None:
    print(f"Running E2E CLI smoke tests against: {bin_path} ({worker_count} async workers)\n")
    for i in range(total):
        ok, msg = results_map[i]
        if ok:
            print(msg, flush=True)
        else:
            print(msg, file=sys.stderr, flush=True)
    print(
        f"\n=========================================\n"
        f"E2E Smoke Summary: {passed} passed, {failed} failed (Total: {total}) in {elapsed:.2f}s\n"
        f"========================================="
    )


def print_failed_results(
    results_map: dict[int, tuple[bool, str]],
    total: int,
    passed: int,
    failed: int,
    elapsed: float,
    worker_count: int,
    bin_path: str,
) -> None:
    print(f"Running E2E CLI smoke tests against: {bin_path} ({worker_count} async workers)\n", file=sys.stderr)
    for i in range(total):
        ok, msg = results_map[i]
        if not ok:
            print(msg, file=sys.stderr, flush=True)
    print(
        f"\n=========================================\n"
        f"E2E Smoke Summary: {passed} passed, {failed} failed (Total: {total}) in {elapsed:.2f}s\n"
        f"=========================================",
        file=sys.stderr,
    )


def report_results(
    results_map: dict[int, tuple[bool, str]],
    total: int,
    elapsed: float,
    worker_count: int,
    bin_path: str,
    show_all: bool,
) -> int:
    passed = sum(1 for ok, _ in results_map.values() if ok)
    failed = sum(1 for ok, _ in results_map.values() if not ok)

    if show_all:
        print_all_results(results_map, total, passed, failed, elapsed, worker_count, bin_path)
        return 1 if failed > 0 else 0

    if failed > 0:
        print_failed_results(results_map, total, passed, failed, elapsed, worker_count, bin_path)
        return 1

    print(
        f"✓ All {total} E2E CLI smoke tests passed successfully in {elapsed:.2f}s ({worker_count} async workers)",
        flush=True,
    )
    return 0


def parse_args():
    parser = argparse.ArgumentParser(
        description="Run E2E CLI smoke tests against gitmap binary with async worker group."
    )
    parser.add_argument(
        "bin_path",
        nargs="?",
        default=None,
        help="Path to gitmap binary (defaults to bin/gitmap.exe or gitmap)",
    )
    parser.add_argument(
        "--all",
        "-a",
        action="store_true",
        default=False,
        help="Print all test results (passes and failures). Without this flag, only failures or summary is printed.",
    )
    parser.add_argument(
        "--workers",
        "-w",
        type=int,
        default=None,
        help="Number of concurrent worker tasks in the worker group (default: auto-detected)",
    )
    parser.add_argument(
        "--bin",
        "-b",
        dest="bin_override",
        default=None,
        help="Explicit path to gitmap binary",
    )
    return parser.parse_args()


async def async_main():
    args = parse_args()
    repo_root = os.path.abspath(os.path.join(os.path.dirname(__file__), "..", ".."))

    target_bin = args.bin_override or args.bin_path
    bin_path = resolve_binary_path(target_bin, repo_root)
    worker_count = get_worker_count(args.workers)

    # Assemble full ordered test suite with original deterministic indexing
    all_indexed_items = []
    current_idx = 0

    independent_items_1 = []
    for item in INDEPENDENT_TESTS_PART1:
        all_indexed_items.append((current_idx, item))
        independent_items_1.append((current_idx, item))
        current_idx += 1

    schedule_chain_indexed = []
    for item in SCHEDULE_CHAIN:
        all_indexed_items.append((current_idx, item))
        schedule_chain_indexed.append((current_idx, item))
        current_idx += 1

    macro_chain_indexed = []
    for item in MACRO_CHAIN:
        all_indexed_items.append((current_idx, item))
        macro_chain_indexed.append((current_idx, item))
        current_idx += 1

    independent_items_2 = []
    for item in INDEPENDENT_TESTS_PART2:
        all_indexed_items.append((current_idx, item))
        independent_items_2.append((current_idx, item))
        current_idx += 1

    queue: asyncio.Queue = asyncio.Queue()
    results_map: dict[int, tuple[bool, str]] = {}

    # Enqueue tasks: independent tests individually, stateful chains as grouped sequential tasks
    for item in independent_items_1:
        queue.put_nowait(item)

    queue.put_nowait(schedule_chain_indexed)
    queue.put_nowait(macro_chain_indexed)

    for item in independent_items_2:
        queue.put_nowait(item)

    # Enqueue poison pills for workers
    for _ in range(worker_count):
        queue.put_nowait(None)

    start_time = time.time()

    workers = [
        asyncio.create_task(async_worker(queue, bin_path, repo_root, results_map))
        for _ in range(worker_count)
    ]

    await asyncio.gather(*workers)
    elapsed = time.time() - start_time

    exit_code = report_results(
        results_map=results_map,
        total=len(all_indexed_items),
        elapsed=elapsed,
        worker_count=worker_count,
        bin_path=bin_path,
        show_all=args.all,
    )
    sys.exit(exit_code)


def main():
    asyncio.run(async_main())


if __name__ == "__main__":
    main()
