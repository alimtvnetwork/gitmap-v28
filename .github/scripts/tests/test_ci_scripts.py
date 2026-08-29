#!/usr/bin/env python3
"""Cross-platform unit test suite for all CI/CD Python verification scripts.

Runnable on Windows, Linux, and macOS without external dependencies.
"""
import json
import os
import shutil
import subprocess
import sys
import tempfile
import unittest

SCRIPTS_DIR = os.path.abspath(os.path.join(os.path.dirname(__file__), ".."))
REPO_ROOT = os.path.abspath(os.path.join(SCRIPTS_DIR, ".."))


class TestSingleLinterDiff(unittest.TestCase):
    def setUp(self):
        self.temp_dir = tempfile.mkdtemp()
        self.script = os.path.join(SCRIPTS_DIR, "check-single-linter-diff.py")

    def tearDown(self):
        shutil.rmtree(self.temp_dir, ignore_errors=True)

    def test_clean_report_passes(self):
        current_out = os.path.join(self.temp_dir, "cur.json")
        with open(current_out, "w", encoding="utf-8") as f:
            json.dump({"Issues": []}, f)

        res = subprocess.run(
            [sys.executable, self.script, self.temp_dir, "--linter", "gocritic",
             "--current-out", current_out, "--skip-run"],
            capture_output=True, text=True, encoding="utf-8"
        )
        self.assertEqual(res.returncode, 0)
        self.assertIn("OK: no new gocritic findings", res.stdout)

    def test_seeding_mode_warns_and_exits_zero(self):
        current_out = os.path.join(self.temp_dir, "cur.json")
        with open(current_out, "w", encoding="utf-8") as f:
            json.dump({
                "Issues": [{
                    "FromLinter": "gocritic",
                    "Text": "ifElseChain: rewrite if-else",
                    "Pos": {"Filename": "cmd/test.go", "Line": 10, "Column": 2}
                }]
            }, f)

        res = subprocess.run(
            [sys.executable, self.script, self.temp_dir, "--linter", "gocritic",
             "--current-out", current_out, "--skip-run"],
            capture_output=True, text=True, encoding="utf-8"
        )
        self.assertEqual(res.returncode, 0)
        self.assertIn("::warning file=cmd/test.go,line=10,col=2::[gocritic]", res.stdout)
        self.assertIn("Seeding mode", res.stderr)

    def test_diff_with_new_violation_fails(self):
        baseline_file = os.path.join(self.temp_dir, "base.json")
        with open(baseline_file, "w", encoding="utf-8") as f:
            json.dump({
                "Issues": [{
                    "FromLinter": "gocritic",
                    "Text": "preExisting: issue",
                    "Pos": {"Filename": "cmd/old.go", "Line": 5, "Column": 1}
                }]
            }, f)

        current_out = os.path.join(self.temp_dir, "cur.json")
        with open(current_out, "w", encoding="utf-8") as f:
            json.dump({
                "Issues": [
                    {
                        "FromLinter": "gocritic",
                        "Text": "preExisting: issue",
                        "Pos": {"Filename": "cmd/old.go", "Line": 5, "Column": 1}
                    },
                    {
                        "FromLinter": "gocritic",
                        "Text": "appendAssign: append result not assigned",
                        "Pos": {"Filename": "cmd/new.go", "Line": 20, "Column": 4}
                    }
                ]
            }, f)

        res = subprocess.run(
            [sys.executable, self.script, self.temp_dir, "--linter", "gocritic",
             "--baseline", baseline_file, "--current-out", current_out, "--skip-run"],
            capture_output=True, text=True, encoding="utf-8"
        )
        self.assertEqual(res.returncode, 1)
        self.assertIn("::error file=cmd/new.go,line=20,col=4::[gocritic]", res.stdout)
        self.assertIn("FAIL: 1 new gocritic finding(s)", res.stderr)


class TestBareStderrCheck(unittest.TestCase):
    def test_bare_stderr_script_runs_clean_on_repo(self):
        script = os.path.join(SCRIPTS_DIR, "check-bare-stderr-err.py")
        res = subprocess.run([sys.executable, script], capture_output=True, text=True, encoding="utf-8")
        self.assertEqual(res.returncode, 0)
        self.assertIn("OK: no bare 'fmt.Fprintln(os.Stderr, err)' in gitmap/cmd/", res.stdout)


class TestCmdNamingCheck(unittest.TestCase):
    def test_cmd_naming_runs_clean_on_repo(self):
        script = os.path.join(SCRIPTS_DIR, "check-cmd-naming.py")
        res = subprocess.run([sys.executable, script], capture_output=True, text=True, encoding="utf-8")
        self.assertEqual(res.returncode, 0)
        self.assertIn("All cmd/ helper names are domain-qualified.", res.stdout)


class TestConstantsNamingCheck(unittest.TestCase):
    def test_constants_naming_runs_clean_on_repo(self):
        script = os.path.join(SCRIPTS_DIR, "check-constants-naming.py")
        res = subprocess.run([sys.executable, script], capture_output=True, text=True, encoding="utf-8")
        self.assertEqual(res.returncode, 0)
        self.assertIn("constants pass naming checks", res.stdout)


class TestDeployLayoutCheck(unittest.TestCase):
    def test_deploy_layout_runs_clean_on_repo(self):
        script = os.path.join(SCRIPTS_DIR, "check-deploy-layout.py")
        res = subprocess.run([sys.executable, script], capture_output=True, text=True, encoding="utf-8")
        self.assertEqual(res.returncode, 0)
        self.assertIn("Deploy layout check: OK", res.stdout)


class TestLegacyRefsCheck(unittest.TestCase):
    def test_legacy_refs_runs_clean_on_repo(self):
        script = os.path.join(SCRIPTS_DIR, "check-legacy-refs.py")
        res = subprocess.run([sys.executable, script], capture_output=True, text=True, encoding="utf-8")
        self.assertEqual(res.returncode, 0)
        self.assertIn("no forbidden legacy refs found", res.stdout)


class TestNoGoldenAllowLeak(unittest.TestCase):
    def test_golden_leak_runs_clean_on_repo(self):
        script = os.path.join(SCRIPTS_DIR, "check-no-golden-allow-leak.py")
        res = subprocess.run([sys.executable, script], capture_output=True, text=True, encoding="utf-8")
        self.assertEqual(res.returncode, 0)
        self.assertIn("GITMAP_ALLOW_GOLDEN_UPDATE check: OK", res.stdout)


class TestChangelogVersionSync(unittest.TestCase):
    def test_changelog_sync_runs_clean_on_repo(self):
        script = os.path.join(SCRIPTS_DIR, "check-changelog-version-sync.py")
        res = subprocess.run([sys.executable, script], capture_output=True, text=True, encoding="utf-8")
        self.assertEqual(res.returncode, 0)
        self.assertIn("changelog.md has entry for v", res.stdout)


class TestSmokeInstaller(unittest.TestCase):
    def test_smoke_installer_source_mode(self):
        script = os.path.join(SCRIPTS_DIR, "smoke-installer.py")
        res = subprocess.run([sys.executable, script, "source"], capture_output=True, text=True, encoding="utf-8")
        self.assertEqual(res.returncode, 0)
        self.assertIn("Installer smoke test passed", res.stdout)


if __name__ == "__main__":
    unittest.main()
