#!/usr/bin/env python3
import concurrent.futures
import os
import subprocess
import sys
import time

def run_cmd(bin_path: str, args: list, cwd: str = None):
    env = os.environ.copy()
    env['GITMAP_SKIP_DELAY'] = '1'
    env['GITMAP_NON_INTERACTIVE'] = '1'
    env['CI'] = '1'
    try:
        res = subprocess.run(
            [bin_path] + args,
            cwd=cwd,
            env=env,
            capture_output=True,
            text=True,
            encoding='utf-8',
            errors='replace',
            timeout=35,
        )
        return res.returncode, res.stdout, res.stderr
    except subprocess.TimeoutExpired:
        return 124, '', 'Error: Command timed out after 35 seconds'

def execute_single_test(bin_path, repo_root, idx, test_item):
    args, valid_codes, exp_sub, desc = test_item
    code, stdout, stderr = run_cmd(bin_path, args, cwd=repo_root)
    out = stdout + stderr
    if code not in valid_codes:
        msg = f'❌ FAIL: {desc} (gitmap {" ".join(args)}) exited with {code}, expected {valid_codes}. Output:\n{out.strip()}'
        return idx, False, msg
    if exp_sub and exp_sub not in out:
        msg = f'❌ FAIL: {desc} (missing substring {exp_sub}). Output:\n{out.strip()}'
        return idx, False, msg
    return idx, True, f'  ✓ PASS: {desc}'

def execute_sequential_chain(bin_path, repo_root, items_with_indices):
    results = []
    for idx, test_item in items_with_indices:
        results.append(execute_single_test(bin_path, repo_root, idx, test_item))
    return results

def main():
    if hasattr(sys.stdout, 'reconfigure'):
        sys.stdout.reconfigure(encoding='utf-8')
    if hasattr(sys.stderr, 'reconfigure'):
        sys.stderr.reconfigure(encoding='utf-8')

    repo_root = os.path.abspath(os.path.join(os.path.dirname(__file__), '..', '..'))
    bin_name = 'gitmap.exe' if os.name == 'nt' else 'gitmap'
    bin_path = os.path.abspath(sys.argv[1]) if len(sys.argv) > 1 else os.path.join(repo_root, bin_name)

    if not os.path.isfile(bin_path):
        print(f'Path {bin_path} does not exist', file=sys.stderr)
        sys.exit(1)

    print(f'Running E2E CLI smoke tests against: {bin_path}')
    start_time = time.time()

    tests_List = [
        (['version'], [0], 'gitmap v', 'Version command'),
        (['v'], [0], 'gitmap v', 'Version alias v'),
        (['help'], [0], 'Usage: gitmap', 'Help command'),
        (['--help'], [0], 'Usage: gitmap', '--help flag'),
        (['-h'], [0], 'Usage: gitmap', '-h flag'),
        (['llm'], [0], 'Gitmap LLM Specification', 'llm command'),
        (['llm-docs', '--stdout'], [0], 'LLM.md', 'llm-docs --stdout'),
        (['find'], [0], 'Usage:', 'find zero-args help'),
        (['find', '*version*.json'], [0], '', 'find with wildcard pattern'),
        (['f'], [0], 'Usage:', 'f zero-args alias help'),
        (['find-files'], [0], 'Usage:', 'find-files zero-args help'),
        (['find-files', 'version.json'], [0], '', 'find-files exact query'),
        (['ff'], [0], 'Usage:', 'ff zero-args alias help'),
        (['find-files-any'], [0], 'Usage:', 'find-files-any zero-args help'),
        (['find-files-any', 'smoke', '-ext', 'py'], [0], '', 'find-files-any with -ext filter'),
        (['ffa'], [0], 'Usage:', 'ffa zero-args alias help'),
        (['find-files-startswith'], [0], 'Usage:', 'find-files-startswith zero-args help'),
        (['find-files-startswith', 'version'], [0], '', 'find-files-startswith prefix query'),
        (['ffs'], [0], 'Usage:', 'ffs zero-args alias help'),
        (['find-files-endswith'], [0], 'Usage:', 'find-files-endswith zero-args help'),
        (['find-files-endswith', '.json'], [0], '', 'find-files-endswith suffix query'),
        (['ffe'], [0], 'Usage:', 'ffe zero-args alias help'),
        (['list-files', '--help'], [0], '', 'list-files --help'),
        (['search'], [0], 'Usage:', 'search zero-args help'),
        (['search', '--help'], [0], '', 'search --help'),
        (['search', 'Version'], [0, 1], '', 'search query'),
        (['repo-search'], [0], 'Usage:', 'repo-search zero-args help'),
        (['repo-search-json', '--help'], [0], '', 'repo-search-json --help'),
        (['find-regex'], [0], 'Usage:', 'find-regex zero-args help'),
        (['find-read'], [0], 'Usage:', 'find-read zero-args help'),
        (['find-regex-read'], [0], 'Usage:', 'find-regex-read zero-args help'),
        (['doctor'], [0, 1], 'gitmap doctor', 'doctor command'),
        (['pipeline', 'status', '--json'], [0], 'isRunning', 'pipeline status --json'),
        (['pipeline', 'eta'], [0], '', 'pipeline eta'),
        (['pipeline-ai', '--help'], [0], 'Usage', 'pipeline-ai --help'),
        (['pipeline-ai', 'status', '--json'], [0], 'isRunning', 'pipeline-ai status --json'),
        (['pipeline-ai', 'status', '-t', '25', '--json'], [0], 'isRunning', 'pipeline-ai status with -t'),
        (['pipeline-ai', 'eta', '--json'], [0], 'isRunning', 'pipeline-ai eta --json'),
        (['pl-ai', 'status', '--json'], [0], 'isRunning', 'pl-ai alias'),
        (['eta'], [0], '', 'eta alias'),
        (['waittime'], [0], '', 'waittime alias'),
        (['pipeline', 'error-logs', '--help'], [0], '', 'pipeline error-logs --help'),
        (['pipeline', 'logs', '--help'], [0], '', 'pipeline logs --help'),
        (['status'], [0], '', 'status command'),
        (['st'], [0], '', 'status alias st'),
        (['antigravity', 'stats'], [0], 'Account: Default', 'antigravity stats'),
        (['agy', 'stats'], [0], 'Account: Default', 'agy stats alias'),
        (['ag', 'stats'], [0], 'Account: Default', 'ag stats alias'),
        (['agy', 'ls'], [0], '', 'agy ls command'),
        (['vscode', 'ls'], [0], '', 'vscode ls command'),
        (['vsc', 'ls'], [0], '', 'vsc ls alias'),
        (['schedule', '--help'], [0], '', 'schedule --help'),
        (['schedule', 'list', '--json'], [0], '', 'schedule list --json'),
        (['schedule', 'add', 'smoke-job', 'echo smoke', '--every', '1h'], [0], '', 'schedule add command'),
        (['schedule', 'status', 'smoke-job', '--json'], [0], 'smoke-job', 'schedule status --json'),
        (['schedule', 'run', 'smoke-job'], [0], '', 'schedule run command'),
        (['schedule', 'log', 'smoke-job', '--json'], [0], 'runnerUser', 'schedule log --json'),
        (['schedule', 'export-all', '--json'], [0], '', 'schedule export-all --json'),
        (['schedule', 'disable', 'smoke-job'], [0], '', 'schedule disable command'),
        (['schedule', 'enable', 'smoke-job'], [0], '', 'schedule enable command'),
        (['schedule', 'reset', 'smoke-job'], [0], '', 'schedule reset command'),
        (['schedule', 'rm', 'smoke-job'], [0], '', 'schedule rm command'),
        (['macro', '--help'], [0], '', 'macro --help'),
        (['macro', 'add', 'smoke-macro', 'echo smoke-macro'], [0], '', 'macro add command'),
        (['macro', 'list', '--json'], [0], '', 'macro list --json'),
        (['macro', 'rm', 'smoke-macro'], [0], '', 'macro rm command'),
        (['retry', '--sleep=10ms', '--max-retries=1', 'echo retry-smoke-test'], [0], 'retry-smoke-test', 'retry command'),
        (['chrome', '--help'], [0], '', 'chrome --help'),
        (['replace', 'history'], [0], '', 'replace history command'),
        (['go-repos', '--json'], [0, 1], '', 'go-repos --json'),
        (['gr', '--json'], [0, 1], '', 'gr --json alias'),
        (['node-repos', '--json'], [0, 1], '', 'node-repos --json'),
        (['react-repos', '--json'], [0, 1], '', 'react-repos --json'),
        (['cpp-repos', '--json'], [0, 1], '', 'cpp-repos --json'),
        (['csharp-repos', '--json'], [0, 1], '', 'csharp-repos --json'),
        (['error', 'ls'], [0], '', 'error ls command'),
        (['error', 'warnings'], [0], '', 'error warnings command'),
        (['completion', 'bash'], [0], '', 'completion bash'),
        (['completion', 'powershell'], [0], '', 'completion powershell'),
        (['completion', 'zsh'], [0], '', 'completion zsh'),
        (['diff-profiles', '--help'], [0], '', 'diff-profiles --help'),
        (['zip-group', '--help'], [0], '', 'zip-group --help'),
        (['env', '--help'], [0], '', 'env --help'),
        (['task', '--help'], [0], '', 'task --help'),
        (['amend', '--help'], [0], '', 'amend --help'),
        (['amend-list', '--help'], [0], '', 'amend-list --help'),
        (['seo-write', '--help'], [0], '', 'seo-write --help'),
        (['fix-repo', '--help'], [0], '', 'fix-repo --help'),
        (['gomod', '--help'], [0], '', 'gomod --help'),
        (['serve', '--help'], [0], '', 'serve --help'),
        (['open', '--help'], [0], '', 'open --help'),
        (['make-public', '--help'], [0], '', 'make-public --help'),
        (['make-private', '--help'], [0], '', 'make-private --help'),
        (['make-all-public', '--help'], [0], '', 'make-all-public --help'),
        (['make-all-private', '--help'], [0], '', 'make-all-private --help'),
        (['visibility-history', '--help'], [0], '', 'visibility-history --help'),
        (['visibility-undo', '--help'], [0], '', 'visibility-undo --help'),
        (['profiles', '--help'], [0], '', 'profiles --help'),
        (['profiles', 'ls', '--json'], [0], '', 'profiles ls --json'),
        (['profiles', 'status'], [0], '', 'profiles status'),
        (['create', '--help'], [0], '', 'create --help'),
        (['backup', '--help'], [0], '', 'backup --help'),
        (['backup', 'status'], [0], '', 'backup status'),
        (['backup', 'ls'], [0], '', 'backup ls'),
    ]

    # Stateful indices that must run in strict sequence to avoid SQLite locking
    schedule_indices = set(range(51, 62))  # schedule --help -> schedule rm
    macro_indices = set(range(62, 66))     # macro --help -> macro rm

    schedule_chain = [(i, tests_List[i]) for i in range(51, 62)]
    macro_chain = [(i, tests_List[i]) for i in range(62, 66)]

    results_map = {}
    futures = []

    # Parallel execution with thread pool
    max_workers = min(12, os.cpu_count() * 2 if os.cpu_count() else 8)
    with concurrent.futures.ThreadPoolExecutor(max_workers=max_workers) as executor:
        # Submit all independent tests in parallel
        for i, item in enumerate(tests_List):
            if i in schedule_indices or i in macro_indices:
                continue
            futures.append(executor.submit(execute_single_test, bin_path, repo_root, i, item))

        for f in concurrent.futures.as_completed(futures):
            res = f.result()
            results_map[res[0]] = (res[1], res[2])

    # Run stateful chains in sequential isolation
    for r in execute_sequential_chain(bin_path, repo_root, schedule_chain):
        results_map[r[0]] = (r[1], r[2])
    for r in execute_sequential_chain(bin_path, repo_root, macro_chain):
        results_map[r[0]] = (r[1], r[2])

    passed = 0
    failed = 0
    for i in range(len(tests_List)):
        ok, msg = results_map[i]
        if ok:
            print(msg, flush=True)
            passed += 1
        else:
            print(msg, file=sys.stderr, flush=True)
            failed += 1

    elapsed = time.time() - start_time
    print(f'\n=========================================\nE2E Smoke Summary: {passed} passed, {failed} failed (Total: {len(tests_List)}) in {elapsed:.2f}s\n=========================================')
    if failed > 0:
        sys.exit(1)
    sys.exit(0)

if __name__ == '__main__':
    main()
