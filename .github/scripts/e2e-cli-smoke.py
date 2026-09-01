#!/usr/bin/env python3
import os, subprocess, sys, tempfile, time

def run_cmd(bin_path: str, args: list, cwd: str = None):
    env = os.environ.copy()
    env['GITMAP_SKIP_DELAY'] = '1'
    res = subprocess.run([bin_path] + args, cwd=cwd, env=env, capture_output=True, text=True, encoding='utf-8', errors='replace')
    return res.returncode, res.stdout, res.stderr

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
        (['retry', '--sleep=10ms', '--max-retries=1', 'gitmap version'], [0], 'gitmap v', 'retry command'),
        (['chrome', '--help'], [0], '', 'chrome --help'),
        (['replace', 'history'], [0], '', 'replace history command'),
        (['go-repos', '--json'], [0], '', 'go-repos --json'),
        (['gr', '--json'], [0], '', 'gr --json alias'),
        (['node-repos', '--json'], [0], '', 'node-repos --json'),
        (['react-repos', '--json'], [0], '', 'react-repos --json'),
        (['cpp-repos', '--json'], [0], '', 'cpp-repos --json'),
        (['csharp-repos', '--json'], [0], '', 'csharp-repos --json'),
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
    ]

    failed = 0
    passed = 0
    for args, valid_codes, exp_sub, desc in tests_List:
        code, stdout, stderr = run_cmd(bin_path, args, cwd=repo_root)
        out = stdout + stderr
        if code not in valid_codes:
            print(f'❌ FAIL: {desc} (gitmap {" ".join(args)}) exited with {code}, expected {valid_codes}. Output:\n{out.strip()}', file=sys.stderr)
            failed += 1
        elif exp_sub and exp_sub not in out:
            print(f'❌ FAIL: {desc} (missing substring {exp_sub}). Output:\n{out.strip()}', file=sys.stderr)
            failed += 1
        else:
            print(f'  ✓ PASS: {desc}')
            passed += 1

    print(f'\n=========================================\nE2E Smoke Summary: {passed} passed, {failed} failed (Total: {len(tests_List)})\n=========================================')
    if failed > 0:
        sys.exit(1)
    sys.exit(0)

if __name__ == '__main__':
    main()
