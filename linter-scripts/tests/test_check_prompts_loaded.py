"""Tests for ``linter-scripts/check-prompts-loaded.py``.

Covers the four contractual outcomes:

  * PASS  — index lists every prompt on disk.
  * FAIL  — orphan prompt (file exists, no reference).
  * FAIL  — dangling reference (referenced filename, missing file).
  * ERROR — index file or prompts directory missing.

Each test invokes the script as a subprocess against a temporary repo
layout to keep the assertions black-box (exit code + stderr/stdout).
"""
from __future__ import annotations

import subprocess
import sys
import tempfile
import unittest
from pathlib import Path

REPO_ROOT = Path(__file__).resolve().parents[2]
SCRIPT = REPO_ROOT / "linter-scripts" / "check-prompts-loaded.py"


def write_file(path: Path, body: str) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(body, encoding="utf-8")


def run_check(index: Path, prompts_dir: Path) -> subprocess.CompletedProcess:
    cmd = [sys.executable, str(SCRIPT),
           "--index", str(index),
           "--prompts-dir", str(prompts_dir)]
    return subprocess.run(cmd, capture_output=True, text=True, check=False)


def make_layout(tmp: Path, index_body: str | None,
                prompt_names: list[str]) -> tuple[Path, Path]:
    index = tmp / ".lovable" / "prompts.md"
    prompts_dir = tmp / ".lovable" / "prompts"
    prompts_dir.mkdir(parents=True, exist_ok=True)
    for name in prompt_names:
        write_file(prompts_dir / name, "# stub\n")
    if index_body is not None:
        write_file(index, index_body)
    return index, prompts_dir


class CheckPromptsLoadedTests(unittest.TestCase):

    def test_pass_when_index_lists_every_prompt(self) -> None:
        with tempfile.TemporaryDirectory() as raw:
            tmp = Path(raw)
            body = "# Index\n- prompts/01-a.md\n- prompts/02-b.md\n"
            index, prompts_dir = make_layout(tmp, body,
                                             ["01-a.md", "02-b.md"])
            result = run_check(index, prompts_dir)
            self.assertEqual(result.returncode, 0, msg=result.stdout)

    def test_fail_on_orphan_prompt(self) -> None:
        with tempfile.TemporaryDirectory() as raw:
            tmp = Path(raw)
            body = "# Index\n- prompts/01-a.md\n"
            index, prompts_dir = make_layout(tmp, body,
                                             ["01-a.md", "02-b.md"])
            result = run_check(index, prompts_dir)
            self.assertEqual(result.returncode, 1)
            self.assertIn("Orphan prompts", result.stdout)

    def test_fail_on_dangling_reference(self) -> None:
        with tempfile.TemporaryDirectory() as raw:
            tmp = Path(raw)
            body = "# Index\n- prompts/01-a.md\n- prompts/99-missing.md\n"
            index, prompts_dir = make_layout(tmp, body, ["01-a.md"])
            result = run_check(index, prompts_dir)
            self.assertEqual(result.returncode, 1)
            self.assertIn("Dangling references", result.stdout)

    def test_error_when_index_missing(self) -> None:
        with tempfile.TemporaryDirectory() as raw:
            tmp = Path(raw)
            _, prompts_dir = make_layout(tmp, None, ["01-a.md"])
            missing_index = tmp / ".lovable" / "prompts.md"
            result = run_check(missing_index, prompts_dir)
            self.assertEqual(result.returncode, 2)
            self.assertIn("prompt index not found", result.stderr)

    def test_error_when_prompts_dir_empty(self) -> None:
        with tempfile.TemporaryDirectory() as raw:
            tmp = Path(raw)
            index, prompts_dir = make_layout(tmp, "# Index\n", [])
            result = run_check(index, prompts_dir)
            self.assertEqual(result.returncode, 2)
            self.assertIn("no prompt files found", result.stderr)


if __name__ == "__main__":
    unittest.main()
