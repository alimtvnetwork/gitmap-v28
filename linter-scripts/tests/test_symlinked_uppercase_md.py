"""Tests for uppercase ``.MD`` detection via symlinks inside the repo root.

Background: :func:`iter_markdown_files` is the full-tree walker the
linter uses when no ``--changed-files`` / ``--diff-base`` intake is
provided. It is built on ``Path.rglob("*")`` and lower-cases the
suffix at filter time so ``UPPER.MD`` is matched on every platform
(see :mod:`test_case_insensitive_extensions`).

A separate failure mode that the case-insensitivity tests do NOT
exercise is symlinked entries:

* A **file symlink** placed under the scanned root that points at a
  real ``.MD`` file elsewhere in the repo. ``rglob`` yields the
  symlink itself and ``Path.is_file()`` follows it, so the file
  *should* be classified as a markdown file regardless of the suffix
  case on either side of the link.
* A **directory symlink** placed under the scanned root that points
  at a real directory containing ``.MD`` files. By design ``rglob``
  does NOT recurse into directory symlinks (this protects the linter
  from cycles and from accidentally walking a vendored
  ``node_modules`` pointer). We pin that contract here so a future
  ``follow_symlinks=True`` refactor can't silently change the walk
  semantics without updating the test (and a corresponding
  constraint memory).

The end-to-end CLI test additionally proves that the symlinked
``.MD`` actually reaches the rule engine — a real P-001 / P-004
violation embedded in the *target* file is reported against the
*symlink path under the scanned root*, which is what humans grep
for in CI logs.
"""

from __future__ import annotations

import importlib.util
import json
import subprocess
import sys
import tempfile
import unittest
from pathlib import Path

LINTER = (Path(__file__).resolve().parent.parent
          / "check-placeholder-comments.py")


def _symlinks_supported(root: Path) -> bool:
    """Best-effort probe for symlink support on the current FS.

    Some CI sandboxes (notably Windows runners without Developer Mode
    and a handful of network filesystems) raise ``OSError`` on
    ``Path.symlink_to``. Skipping the test there beats a spurious
    red CI — the platforms we actually care about (Linux + macOS)
    all support symlinks unconditionally.
    """
    probe_target = root / "_probe_target"
    probe_link = root / "_probe_link"
    probe_target.write_text("x", encoding="utf-8")
    try:
        probe_link.symlink_to(probe_target)
    except (OSError, NotImplementedError) as exc:
        import sys as _sys; print(f"Error: {exc}", file=_sys.stderr)
        return False
    finally:
        try:
            probe_link.unlink()
        except FileNotFoundError as exc:
            import sys as _sys; print(f"Error: {exc}", file=_sys.stderr)
            pass
        probe_target.unlink()
    return True


def _load_module():  # type: ignore[no-untyped-def]
    spec = importlib.util.spec_from_file_location("_pc_symlinks", LINTER)
    assert spec is not None and spec.loader is not None
    mod = importlib.util.module_from_spec(spec)
    sys.modules["_pc_symlinks"] = mod
    spec.loader.exec_module(mod)
    return mod


def _run(*args: str, cwd: Path) -> subprocess.CompletedProcess:
    return subprocess.run([sys.executable, str(LINTER), *args],
                          cwd=cwd, capture_output=True, text=True)


class IterMarkdownFilesFollowsFileSymlinks(unittest.TestCase):
    def setUp(self) -> None:
        self.mod = _load_module()
        self._tmp = tempfile.TemporaryDirectory()
        self.addCleanup(self._tmp.cleanup)
        self.root = Path(self._tmp.name)
        if not _symlinks_supported(self.root):
            self.skipTest("filesystem does not support symlinks")

    def test_uppercase_md_via_file_symlink_is_yielded(self) -> None:
        # Real target lives OUTSIDE the scanned subdir but INSIDE
        # the repo root, mirroring a common monorepo layout where
        # ``spec/`` links into a shared ``vendor/docs/`` tree.
        outside = self.root / "vendor"
        outside.mkdir()
        target = outside / "REAL.MD"
        target.write_text("# real\n", encoding="utf-8")
        scanned = self.root / "spec"
        scanned.mkdir()
        link = scanned / "LINKED.MD"
        link.symlink_to(target)

        files = sorted(p.name for p in
                       self.mod.iter_markdown_files(scanned))
        self.assertEqual(files, ["LINKED.MD"],
            msg="file symlink with uppercase .MD suffix must be "
                "discovered by the full-tree walker")

    def test_mixed_case_link_to_lowercase_target_is_yielded(self) -> None:
        # Suffix on the *link* drives classification — the walker
        # never inspects the target's suffix. This pins that
        # behaviour so a future refactor that lowercases the
        # resolved path instead of the link path can't change
        # semantics.
        target = self.root / "real.md"
        target.write_text("# real\n", encoding="utf-8")
        scanned = self.root / "spec"
        scanned.mkdir()
        (scanned / "Mirror.Md").symlink_to(target)

        files = sorted(p.name for p in
                       self.mod.iter_markdown_files(scanned))
        self.assertEqual(files, ["Mirror.Md"])

    def test_link_with_non_md_suffix_is_filtered(self) -> None:
        # A symlink whose own suffix is ``.txt`` must NOT be picked
        # up even if the target is a ``.MD``. Case-insensitivity does
        # not widen the allowlist (see test_case_insensitive_extensions
        # for the non-symlink analogue).
        target = self.root / "real.MD"
        target.write_text("# real\n", encoding="utf-8")
        scanned = self.root / "spec"
        scanned.mkdir()
        (scanned / "alias.txt").symlink_to(target)

        files = [p.name for p in self.mod.iter_markdown_files(scanned)]
        self.assertEqual(files, [])

    def test_directory_symlink_is_not_recursed_into(self) -> None:
        # rglob's default (``follow_symlinks=False``) intentionally
        # skips directory symlinks. We pin the current contract: an
        # uppercase .MD reachable ONLY through a directory symlink is
        # NOT walked. If we ever flip this, update both the test and
        # add a memory under .lovable/memory/architecture/ so the
        # change is documented.
        outside = self.root / "vendor_dir"
        outside.mkdir()
        (outside / "DEEP.MD").write_text("# deep\n", encoding="utf-8")
        scanned = self.root / "spec"
        scanned.mkdir()
        (scanned / "linkdir").symlink_to(outside,
                                          target_is_directory=True)
        # Sanity: a sibling real .MD next to the dir symlink IS
        # found, so we know the walk itself ran.
        (scanned / "Sibling.MD").write_text("# s\n", encoding="utf-8")

        files = sorted(p.name for p in
                       self.mod.iter_markdown_files(scanned))
        self.assertEqual(files, ["Sibling.MD"],
            msg="rglob must not descend into directory symlinks; "
                "only the sibling real .MD should be found")

    def test_symlink_and_target_both_yielded_when_under_same_root(self) -> None:
        # When a symlink lives in the same scanned tree as its
        # target both entries are walked separately — they have
        # distinct ``Path`` keys (different lexical paths), so the
        # dedup ``set`` in the walker keeps both. Pin this so a
        # future caller that passes ``--changed-files`` for the
        # symlink path still gets an audit row even when the
        # target is also listed.
        target = self.root / "real.MD"
        target.write_text("# real\n", encoding="utf-8")
        (self.root / "alias.MD").symlink_to(target)

        files = sorted(p.name for p in
                       self.mod.iter_markdown_files(self.root))
        self.assertEqual(files, ["alias.MD", "real.MD"])


class CliEndToEndSymlinkedUppercaseMd(unittest.TestCase):
    """Full-process smoke tests that the symlinked .MD actually
    reaches the rule engine, not just the file iterator."""

    def setUp(self) -> None:
        self._tmp = tempfile.TemporaryDirectory()
        self.addCleanup(self._tmp.cleanup)
        self.root = Path(self._tmp.name)
        if not _symlinks_supported(self.root):
            self.skipTest("filesystem does not support symlinks")

    def test_violation_in_symlinked_md_is_reported(self) -> None:
        # Target file carries a real placeholder violation. The link
        # sits under the scanned ``spec/`` root.
        vendor = self.root / "vendor"
        vendor.mkdir()
        target = vendor / "TARGET.MD"
        target.write_text(
            "# Title\n\n"
            "<!-- TODO(P-003): describe target -->\n"
            "- [TODO(P-003): describe target]"
            "(relative/path/to/spec.txt)\n",
            encoding="utf-8",
        )
        spec = self.root / "spec"
        spec.mkdir()
        (spec / "LINKED.MD").symlink_to(target)

        proc = _run("--root", "spec", "--repo-root", str(self.root),
                    "--json", cwd=self.root)
        self.assertEqual(proc.returncode, 1,
            msg=f"stdout={proc.stdout!r} stderr={proc.stderr!r}")
        violations = json.loads(proc.stdout)
        self.assertTrue(violations,
            msg="symlinked uppercase .MD reached the iterator but "
                "produced zero violations — rule engine never saw it")
        # Every reported violation must be attributed to the *link
        # path under the scanned root*, not the resolved target. CI
        # log-grep tooling keys on this exact string.
        for v in violations:
            self.assertEqual(v["file"], "spec/LINKED.MD",
                msg=f"unexpected file path in violation: {v!r}")
        codes = {v["code"] for v in violations}
        # P-001 fires because "(P-003):" is not an imperative verb;
        # P-004 fires because the placeholder block has no bullets
        # the parser accepts. Either is acceptable evidence the
        # rules ran — assert at least one of the expected codes.
        self.assertTrue(codes & {"P-001", "P-003", "P-004"},
            msg=f"none of the expected rule codes fired: {codes}")

    def test_clean_symlinked_md_exits_zero(self) -> None:
        # No-violation control: the symlinked .MD reaches the engine
        # but contains nothing to flag, so exit code is 0 and stdout
        # is an empty JSON array.
        vendor = self.root / "vendor"
        vendor.mkdir()
        (vendor / "CLEAN.MD").write_text(
            "# Clean\n\nNo placeholders.\n", encoding="utf-8")
        spec = self.root / "spec"
        spec.mkdir()
        (spec / "MIRROR.MD").symlink_to(vendor / "CLEAN.MD")

        proc = _run("--root", "spec", "--repo-root", str(self.root),
                    "--json", cwd=self.root)
        self.assertEqual(proc.returncode, 0,
            msg=f"stdout={proc.stdout!r} stderr={proc.stderr!r}")
        self.assertEqual(json.loads(proc.stdout), [])


if __name__ == "__main__":
    unittest.main()
