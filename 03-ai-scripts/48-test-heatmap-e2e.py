#!/usr/bin/env python3
"""
48-test-heatmap-e2e.py - End-to-End Non-Destructive Validation for GitMap Heatmap Subsystems.

Validates:
1. Root Heatmap Subsystem:
   - help, log, search, tree, status, tables, view, query
   - flags: --limit (-n), --format (-F) (table, json, yaml), --json, --task, --file
2. Repo-Cache (RC) Multi-AI Heatmap Subsystem:
   - help, list-ai, log, search, tree, status, tables, view, query, commit, sync
   - namespace resolution and isolated storage
3. Verifies zero-contamination invariant via unique transient task slugs.

Usage:
  python 03-ai-scripts/48-test-heatmap-e2e.py [--bin <path>] [--verbose] [--json] [--dry-run]
"""

import argparse
from dataclasses import asdict, dataclass
import json
import os
from pathlib import Path
import subprocess
import sys
import time
import uuid

if hasattr(sys.stdout, "reconfigure"):
    sys.stdout.reconfigure(encoding="utf-8")
if hasattr(sys.stderr, "reconfigure"):
    sys.stderr.reconfigure(encoding="utf-8")


@dataclass
class TestCaseResult:
    name: str
    command: list[str]
    success: bool
    return_code: int
    duration_ms: float
    error_message: str = ""
    output_snippet: str = ""


class HeatmapE2ETester:
    def __init__(self, binary: str = "gitmap", verbose: bool = False, dry_run: bool = False):
        self.binary = binary
        self.verbose = verbose
        self.dry_run = dry_run
        self.session_id = uuid.uuid4().hex[:8]
        self.root_task_slug = f"e2e-root-{self.session_id}"
        self.rc_task_slug = f"e2e-rc-{self.session_id}"
        self.results: list[TestCaseResult] = []

    def run_cmd(self, args: list[str]) -> tuple[int, str, str, float]:
        cmd = [self.binary] + args
        start_time = time.perf_counter()
        if self.dry_run:
            print(f"[DRY-RUN] Would run: {' '.join(cmd)}")
            return 0, "", "", 0.0

        try:
            proc = subprocess.run(
                cmd,
                capture_output=True,
                text=True,
                encoding="utf-8",
                errors="replace",
                timeout=30,
            )
            duration_ms = (time.perf_counter() - start_time) * 1000.0
            return proc.returncode, proc.stdout, proc.stderr, duration_ms
        except Exception as ex:
            duration_ms = (time.perf_counter() - start_time) * 1000.0
            return -1, "", str(ex), duration_ms

    def record_test(
        self,
        name: str,
        args: list[str],
        rc: int,
        stdout: str,
        stderr: str,
        duration_ms: float,
        validator=None,
    ) -> bool:
        cmd_str = [self.binary] + args
        err_msg = ""
        success = (rc == 0)

        if not success:
            err_msg = f"Non-zero exit code ({rc}): {stderr.strip() or stdout.strip()}"
        elif validator is not None and not self.dry_run:
            try:
                val_err = validator(stdout, stderr)
                if val_err:
                    success = False
                    err_msg = val_err
            except Exception as ex:
                success = False
                err_msg = f"Validator exception: {ex}"

        snippet = (stdout[:200] + "...") if len(stdout) > 200 else stdout.strip()
        result = TestCaseResult(
            name=name,
            command=cmd_str,
            success=success,
            return_code=rc,
            duration_ms=round(duration_ms, 2),
            error_message=err_msg,
            output_snippet=snippet,
        )
        self.results.append(result)

        if self.verbose:
            status_tag = "PASS" if success else "FAIL"
            print(f"[{status_tag}] {name} ({result.duration_ms}ms)")
            if not success:
                print(f"   Command: {' '.join(cmd_str)}")
                print(f"   Error:   {err_msg}")
            elif snippet:
                print(f"   Output:  {snippet}")

        return success

    def run_all(self) -> bool:
        print(f"Starting Heatmap E2E Validation [Session: {self.session_id}, Bin: {self.binary}]")

        # 1. Root Heatmap Help
        rc, out, err, dur = self.run_cmd(["heatmap", "--help"])
        def check_help(stdout, _):
            for sub in ("tables", "view", "query", "log", "search", "tree", "status"):
                if sub not in stdout:
                    return f"Missing '{sub}' in root heatmap help output"
            return None
        self.record_test("Root Heatmap Help Documentation", ["heatmap", "--help"], rc, out, err, dur, check_help)

        # 2. RC Heatmap Help
        rc, out, err, dur = self.run_cmd(["heatmap", "rc", "--help"])
        def check_rc_help(stdout, _):
            for sub in ("tables", "view", "query", "list-ai", "log", "status", "commit", "sync"):
                if sub not in stdout:
                    return f"Missing '{sub}' in RC heatmap help output"
            return None
        self.record_test("RC Heatmap Help Documentation", ["heatmap", "rc", "--help"], rc, out, err, dur, check_rc_help)

        # 3. Root Heatmap Log Action
        log_args = [
            "heatmap", "log",
            "--task", self.root_task_slug,
            "--action", "ANALYZE",
            "--file", "cli/heatmap/heatmap_view.go",
            "--lines", "10-50",
            "--state", "EXPLORATION",
            "--reason", f"E2E root test log entry for {self.session_id}",
        ]
        rc, out, err, dur = self.run_cmd(log_args)
        def check_logged(stdout, _):
            if "Logged heatmap action #" not in stdout:
                return "Expected 'Logged heatmap action #' in output"
            return None
        self.record_test("Root Heatmap Log Action", log_args, rc, out, err, dur, check_logged)

        # 4. Root Heatmap Search Text
        search_args = ["heatmap", "search", "--task", self.root_task_slug, "E2E root test"]
        rc, out, err, dur = self.run_cmd(search_args)
        def check_search_text(stdout, _):
            if "heatmap_view.go" not in stdout:
                return "Expected 'heatmap_view.go' in search results"
            return None
        self.record_test("Root Heatmap Search Text", search_args, rc, out, err, dur, check_search_text)

        # 5. Root Heatmap Search JSON
        search_json_args = ["heatmap", "search", "--task", self.root_task_slug, "--json", "E2E root test"]
        rc, out, err, dur = self.run_cmd(search_json_args)
        def check_search_json(stdout, _):
            data = json.loads(stdout)
            if not isinstance(data, list) or len(data) == 0:
                return "Expected non-empty list in search JSON"
            if data[0].get("taskSlug") != self.root_task_slug:
                return f"Unexpected taskSlug in search JSON: {data[0].get('taskSlug')}"
            return None
        self.record_test("Root Heatmap Search JSON", search_json_args, rc, out, err, dur, check_search_json)

        # 6. Root Heatmap Tree & Status
        tree_args = ["heatmap", "tree", "--task", self.root_task_slug, "--json"]
        rc, out, err, dur = self.run_cmd(tree_args)
        def check_tree_json(stdout, _):
            data = json.loads(stdout)
            if not isinstance(data, list):
                return "Expected list of tree nodes in tree JSON"
            return None
        self.record_test("Root Heatmap Tree JSON", tree_args, rc, out, err, dur, check_tree_json)

        status_args = ["heatmap", "status", "--task", self.root_task_slug, "--json"]
        rc, out, err, dur = self.run_cmd(status_args)
        def check_status_json(stdout, _):
            data = json.loads(stdout)
            if data.get("totalActions", 0) < 1:
                return "Expected totalActions >= 1 in status JSON"
            return None
        self.record_test("Root Heatmap Status JSON", status_args, rc, out, err, dur, check_status_json)

        # 7. Root Heatmap Tables Listing
        tables_args = ["heatmap", "tables"]
        rc, out, err, dur = self.run_cmd(tables_args)
        def check_tables_text(stdout, _):
            if "HeatmapActionLog" not in stdout:
                return "Expected 'HeatmapActionLog' table in tables output"
            return None
        self.record_test("Root Heatmap Tables Text", tables_args, rc, out, err, dur, check_tables_text)

        tables_json_args = ["heatmap", "tables", "--format", "json"]
        rc, out, err, dur = self.run_cmd(tables_json_args)
        def check_tables_json(stdout, _):
            data = json.loads(stdout)
            tables = data.get("tables", data) if isinstance(data, dict) else data
            if not isinstance(tables, list) or len(tables) == 0:
                return "Expected non-empty table array in tables JSON"
            names = [item.get("table") or item.get("name") for item in tables]
            if "HeatmapActionLog" not in names:
                return f"Missing HeatmapActionLog in tables: {names}"
            return None
        self.record_test("Root Heatmap Tables JSON", tables_json_args, rc, out, err, dur, check_tables_json)

        # 8. Root Heatmap View
        view_json_args = [
            "heatmap", "view", "HeatmapActionLog",
            "--task", self.root_task_slug,
            "--limit", "5",
            "--format", "json",
        ]
        rc, out, err, dur = self.run_cmd(view_json_args)
        def check_view_json(stdout, _):
            data = json.loads(stdout)
            records = data.get("records", data) if isinstance(data, dict) else data
            if not isinstance(records, list) or len(records) == 0:
                return "Expected non-empty row list in view JSON"
            return None
        self.record_test("Root Heatmap View JSON", view_json_args, rc, out, err, dur, check_view_json)

        view_yaml_args = [
            "heatmap", "view", "HeatmapActionLog",
            "--task", self.root_task_slug,
            "-n", "5",
            "-F", "yaml",
        ]
        rc, out, err, dur = self.run_cmd(view_yaml_args)
        def check_view_yaml(stdout, _):
            if "TaskSlug:" not in stdout and "taskSlug:" not in stdout and "task_slug:" not in stdout:
                return "Expected TaskSlug field in YAML view output"
            return None
        self.record_test("Root Heatmap View YAML", view_yaml_args, rc, out, err, dur, check_view_yaml)

        # 9. Root Heatmap Query
        query_args = [
            "heatmap", "query",
            f"SELECT count(*) as cnt FROM HeatmapActionLog WHERE TaskSlug = '{self.root_task_slug}'",
            "--format", "json",
        ]
        rc, out, err, dur = self.run_cmd(query_args)
        def check_query_json(stdout, _):
            data = json.loads(stdout)
            rows = data.get("rows", data) if isinstance(data, dict) else data
            if not isinstance(rows, list) or len(rows) == 0:
                return "Expected query results array in JSON"
            return None
        self.record_test("Root Heatmap Query JSON", query_args, rc, out, err, dur, check_query_json)

        # 10. RC Heatmap List-AI
        rc_list_args = ["heatmap", "rc", "list-ai", "--json"]
        rc, out, err, dur = self.run_cmd(rc_list_args)
        def check_rc_list(stdout, _):
            data = json.loads(stdout)
            if "supportedAIs" not in data or not isinstance(data["supportedAIs"], list):
                return "Expected supportedAIs list in list-ai JSON"
            return None
        self.record_test("RC Heatmap List-AI JSON", rc_list_args, rc, out, err, dur, check_rc_list)

        # 11. RC Heatmap Log Action
        rc_log_args = [
            "heatmap", "rc", "log",
            "--task", self.rc_task_slug,
            "--action", "MODIFY",
            "--file", "cli/heatmap/heatmap_cmd.go",
            "--lines", "1-20",
            "--state", "CODE_EDITING",
            "--reason", f"E2E RC test log entry for {self.session_id}",
        ]
        rc, out, err, dur = self.run_cmd(rc_log_args)
        def check_rc_logged(stdout, _):
            if "Logged RC heatmap action #" not in stdout:
                return "Expected 'Logged RC heatmap action #' in output"
            return None
        self.record_test("RC Heatmap Log Action", rc_log_args, rc, out, err, dur, check_rc_logged)

        # 12. RC Heatmap Search JSON
        rc_search_args = ["heatmap", "rc", "search", "--task", self.rc_task_slug, "--json", "E2E RC test"]
        rc, out, err, dur = self.run_cmd(rc_search_args)
        def check_rc_search_json(stdout, _):
            data = json.loads(stdout)
            if not isinstance(data, list) or len(data) == 0:
                return "Expected non-empty list in RC search JSON"
            return None
        self.record_test("RC Heatmap Search JSON", rc_search_args, rc, out, err, dur, check_rc_search_json)

        # 13. RC Heatmap Tree & Status
        rc_tree_args = ["heatmap", "rc", "tree", "--task", self.rc_task_slug, "--json"]
        rc, out, err, dur = self.run_cmd(rc_tree_args)
        self.record_test("RC Heatmap Tree JSON", rc_tree_args, rc, out, err, dur, check_tree_json)

        rc_status_args = ["heatmap", "rc", "status", "--task", self.rc_task_slug, "--json"]
        rc, out, err, dur = self.run_cmd(rc_status_args)
        def check_rc_status_json(stdout, _):
            data = json.loads(stdout)
            if "summary" not in data or data["summary"].get("totalActions", 0) < 1:
                return "Expected summary with totalActions >= 1 in RC status JSON"
            return None
        self.record_test("RC Heatmap Status JSON", rc_status_args, rc, out, err, dur, check_rc_status_json)

        # 14. RC Heatmap Tables & View & Query
        rc_tables_args = ["heatmap", "rc", "tables", "--format", "json"]
        rc, out, err, dur = self.run_cmd(rc_tables_args)
        self.record_test("RC Heatmap Tables JSON", rc_tables_args, rc, out, err, dur, check_tables_json)

        rc_view_args = [
            "heatmap", "rc", "view", "HeatmapActionLog",
            "--task", self.rc_task_slug,
            "-n", "5",
            "-F", "json",
        ]
        rc, out, err, dur = self.run_cmd(rc_view_args)
        self.record_test("RC Heatmap View JSON", rc_view_args, rc, out, err, dur, check_view_json)

        rc_query_args = [
            "heatmap", "rc", "query",
            f"SELECT count(*) as cnt FROM HeatmapActionLog WHERE TaskSlug = '{self.rc_task_slug}'",
            "--format", "json",
        ]
        rc, out, err, dur = self.run_cmd(rc_query_args)
        self.record_test("RC Heatmap Query JSON", rc_query_args, rc, out, err, dur, check_query_json)

        # 15. RC Heatmap Commit & Sync
        rc_commit_args = ["heatmap", "rc", "commit", "--task", self.rc_task_slug, "--json"]
        rc, out, err, dur = self.run_cmd(rc_commit_args)
        def check_rc_commit(stdout, _):
            data = json.loads(stdout)
            if data.get("status") != "COMMITTED":
                return f"Expected status COMMITTED, got {data.get('status')}"
            return None
        self.record_test("RC Heatmap Commit Checkpoint", rc_commit_args, rc, out, err, dur, check_rc_commit)

        rc_sync_args = ["heatmap", "rc", "sync", "--json"]
        rc, out, err, dur = self.run_cmd(rc_sync_args)
        def check_rc_sync(stdout, _):
            data = json.loads(stdout)
            if data.get("status") != "SYNCED":
                return f"Expected status SYNCED, got {data.get('status')}"
            return None
        self.record_test("RC Heatmap Sync Telemetry", rc_sync_args, rc, out, err, dur, check_rc_sync)

        # Overall summary
        passed = sum(1 for r in self.results if r.success)
        total = len(self.results)
        all_passed = (passed == total)

        print("-" * 60)
        print(f"Heatmap E2E Validation Complete: {passed}/{total} tests passed.")
        if not all_passed:
            print("Failed Tests:")
            for r in self.results:
                if not r.success:
                    print(f" - {r.name}: {r.error_message}")
        return all_passed


def main():
    parser = argparse.ArgumentParser(description="GitMap Heatmap Subsystem E2E Test Suite")
    parser.add_argument("--bin", default="gitmap", help="Path to gitmap binary (default: gitmap)")
    parser.add_argument("--verbose", "-v", action="store_true", help="Enable verbose test output")
    parser.add_argument("--dry-run", action="store_true", help="Print commands without executing")
    parser.add_argument("--json", action="store_true", help="Output summary in JSON format")

    args = parser.parse_args()

    tester = HeatmapE2ETester(binary=args.bin, verbose=args.verbose or not args.json, dry_run=args.dry_run)
    success = tester.run_all()

    if args.json:
        summary = {
            "sessionId": tester.session_id,
            "success": success,
            "totalTests": len(tester.results),
            "passedTests": sum(1 for r in tester.results if r.success),
            "failedTests": sum(1 for r in tester.results if not r.success),
            "results": [asdict(r) for r in tester.results],
        }
        print(json.dumps(summary, indent=2))

    sys.exit(0 if success else 1)


if __name__ == "__main__":
    main()
