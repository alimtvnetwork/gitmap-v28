#!/usr/bin/env python3
"""
check-tunable-constants.py — Tunable-constants linter (FU-15).

Enforces the contract in `spec/19-main-worker-service/15-tunable-constants.md` §6:

  Rule T1 — Presence: every numeric literal in spec/19/ prose and
            spec/14-update/28-worker-push-instruction.md that is followed
            by a time/count unit (s, sec, seconds, min, minutes, h, hours,
            attempts, retries, times) MUST be either:
              (a) named in §2 of 15-tunable-constants.md (the catalogue
                  table — its key appears in `2.x` rows), OR
              (b) explicitly waivered by `<!-- TUNABLE-WAIVER: ... -->`
                  on the same line.
  Rule T2 — Unique keys: no two §2 rows share the same Key.
  Rule T3 — Seed parity: every §4 `config.seed.json` Settings.<Name>
            `Default` value matches the §2 row for the same logical key.
  Rule T4 — Session-TTL invariant (FU-16):
            MainWorker.Auth.MainSessionAbsoluteMaxSeconds
              >= MainWorker.Auth.MainSessionTtlSeconds
            in BOTH §2 catalogue defaults and §4 seed defaults.

Exit codes:
  0   all rules pass
  1   one or more rules failed
  2   linter setup error

CODE RED compliance:
  - Functions <=15 lines, zero nested ifs, positive guards.
  - Errors are surfaced with file/line context; never swallowed.
"""

from __future__ import annotations
import re
import sys
from pathlib import Path

REPO_ROOT = Path(__file__).resolve().parent.parent
TUNABLES_FILE = REPO_ROOT / "spec/19-main-worker-service/15-tunable-constants.md"
SCAN_GLOBS = (
    "spec/19-main-worker-service/*.md",
    "spec/14-update/28-worker-push-instruction.md",
)
# Files exempt from T1 scanning (catalogue + diagrams + JSON snippets).
EXEMPT_NAMES = {"15-tunable-constants.md"}

# A numeric followed by a unit token. Unit list mirrors §6.
UNIT_RE = re.compile(
    r"\b(\d+(?:\.\d+)?)\s*"
    r"(s|sec|secs|second|seconds|min|mins|minute|minutes|"
    r"h|hr|hrs|hour|hours|attempts?|retries|times)\b",
    re.IGNORECASE,
)
WAIVER_RE = re.compile(r"<!--\s*TUNABLE-WAIVER:")
# Headings/code-fence detection for narrowing prose scan.
CODE_FENCE_RE = re.compile(r"^\s*```")
# §2 row: starts with `| `, has a backticked `Key` and a numeric default.
# §2 row: starts with `| `, has a backticked `Key` in col 1.
ROW_RE = re.compile(r"^\|\s*`([A-Za-z][A-Za-z0-9_.]+)`\s*\|(.*)$")
# Default in col-2: prefer **bold** number; else first backticked token; else first token.
DEFAULT_BOLD_RE = re.compile(r"\*\*([^*|]+?)\*\*")
DEFAULT_TICK_RE = re.compile(r"`([^`|]+?)`")
# §4 Settings.<Name> with Default.
SEED_RE = re.compile(
    r"\"([A-Za-z][A-Za-z0-9_.]+)\"\s*:\s*\{\s*"
    r"\"Type\"\s*:\s*\"[^\"]+\"\s*,\s*\"Default\"\s*:\s*"
    r"(\"[^\"]*\"|true|false|[0-9.]+)"
)


def fail(msg: str) -> None:
    print(f"ERROR: {msg}", file=sys.stderr)


def file_must_exist(path: Path) -> None:
    if path.is_file():
        return
    fail(f"required file missing: {path}")
    sys.exit(2)


def read_lines(path: Path) -> list[str]:
    return path.read_text(encoding="utf-8").splitlines()


def is_prose_line(line: str, in_fence: bool) -> bool:
    return not in_fence and not line.lstrip().startswith("|")


def toggle_fence(line: str, in_fence: bool) -> bool:
    if CODE_FENCE_RE.match(line):
        return not in_fence
    return in_fence


CATALOGUE_START_RE = re.compile(r"^##\s+2\.\s")
CATALOGUE_END_RE = re.compile(r"^##\s+(?!2(?:\.\d+)?\b)\d")


def _is_inside_catalogue(line: str, inside: bool) -> bool:
    if CATALOGUE_START_RE.match(line):
        return True
    if inside and CATALOGUE_END_RE.match(line):
        return False
    return inside


def collect_catalogue_keys() -> set[str]:
    keys: set[str] = set()
    in_fence = False
    for raw in read_lines(TUNABLES_FILE):
        in_fence = toggle_fence(raw, in_fence)
        if in_fence:
            continue
        match = ROW_RE.match(raw)
        if match:
            keys.add(match.group(1))
    return keys


def find_unit_hits(line: str) -> list[tuple[str, str]]:
    return [(m.group(1), m.group(2)) for m in UNIT_RE.finditer(line)]


def line_is_waivered(line: str) -> bool:
    return bool(WAIVER_RE.search(line))


def referenced_keys_in_line(line: str, keys: set[str]) -> bool:
    if any(key in line for key in keys):
        return True
    tails = {key.rsplit(".", 1)[-1] for key in keys if "." in key}
    return any(tail in line for tail in tails)


def scan_file_for_t1(path: Path, keys: set[str]) -> list[str]:
    violations: list[str] = []
    in_fence = False
    for idx, raw in enumerate(read_lines(path), start=1):
        in_fence = toggle_fence(raw, in_fence)
        violation = check_line_t1(path, idx, raw, in_fence, keys)
        if violation:
            violations.append(violation)
    return violations


def check_line_t1(
    path: Path, idx: int, raw: str, in_fence: bool, keys: set[str]
) -> str | None:
    if in_fence:
        return None
    hits = find_unit_hits(raw)
    if not hits:
        return None
    if line_is_waivered(raw):
        return None
    if referenced_keys_in_line(raw, keys):
        return None
    sample = hits[0]
    return f"{path}:{idx}: untracked tunable `{sample[0]} {sample[1]}` — cite 15-tunable-constants.md §2 or add TUNABLE-WAIVER"


def expand_scan_targets() -> list[Path]:
    out: list[Path] = []
    for pattern in SCAN_GLOBS:
        out.extend(sorted(REPO_ROOT.glob(pattern)))
    return [p for p in out if p.name not in EXEMPT_NAMES]


def rule_t1() -> list[str]:
    keys = collect_catalogue_keys()
    if not keys:
        return ["15-tunable-constants.md §2: no catalogue rows parsed"]
    out: list[str] = []
    for path in expand_scan_targets():
        out.extend(scan_file_for_t1(path, keys))
    return out


def rule_t2() -> list[str]:
    # Same §2-only scope as T3: §4.1 alias map's first column reuses the
    # prose key, so we'd get spurious duplicates if we scanned globally.
    seen: dict[str, int] = {}
    in_fence = False
    in_catalogue = False
    for raw in read_lines(TUNABLES_FILE):
        in_fence = toggle_fence(raw, in_fence)
        in_catalogue = _is_inside_catalogue(raw, in_catalogue)
        if in_fence or not in_catalogue:
            continue
        match = ROW_RE.match(raw)
        if not match:
            continue
        seen[match.group(1)] = seen.get(match.group(1), 0) + 1
    return [f"duplicate key `{k}` in §2" for k, n in seen.items() if n > 1]


def collect_seed_defaults() -> dict[str, str]:
    text = TUNABLES_FILE.read_text(encoding="utf-8")
    return {m.group(1): m.group(2).strip('"') for m in SEED_RE.finditer(text)}


# Mapping §4 short-name -> §2 catalogue key. Pinned for explicitness.
SEED_TO_KEY = {
    "RetryMaxAttempts": "MainWorker.Retry.MaxAttempts",
    "RetryJitterPct": "MainWorker.Retry.JitterPct",
    "IdempotencyKeyTtlSeconds": "MainWorker.Idempotency.KeyTtlSeconds",
    "IdempotencyKeyMaxLength": "MainWorker.Idempotency.KeyMaxLength",
    "IdempotencyStoreCleanupSec": "MainWorker.Idempotency.StoreCleanupSeconds",
    "HeartbeatIntervalSeconds": "MainWorker.Heartbeat.IntervalSeconds",
    "HeartbeatMissedThreshold": "MainWorker.Heartbeat.MissedThreshold",
    "HeartbeatQuarantineCooldown": "MainWorker.Heartbeat.QuarantineCooldownSeconds",
    "HeartbeatGraceWindowSeconds": "MainWorker.Heartbeat.GraceWindowSeconds",
    "WorkerJwtTtlSeconds": "MainWorker.Auth.WorkerJwtTtlSeconds",
    "JwtRefreshLeadSeconds": "MainWorker.Auth.JwtRefreshLeadSeconds",
    "MainSessionTtlSeconds": "MainWorker.Auth.MainSessionTtlSeconds",
    "MainSessionAbsoluteMaxSeconds": "MainWorker.Auth.MainSessionAbsoluteMaxSeconds",
    "SessionSlidingExtendOnRead": "MainWorker.Auth.SessionSlidingExtendOnReadOnly",
    "ClockSkewToleranceSeconds": "MainWorker.Auth.ClockSkewToleranceSeconds",
    "RoutingHttpTimeoutSeconds": "MainWorker.Routing.HttpTimeoutSeconds",
    "RoutingHandshakeTimeoutSec": "MainWorker.Routing.HttpHandshakeTimeoutSeconds",
    "RoutingMaxConcurrentPerNode": "MainWorker.Routing.MaxConcurrentPerWorker",
    "RateAuthPerMinutePerIp": "MainWorker.RateLimit.AuthEndpointsPerMinutePerIp",
    "RateWorkerPerMinutePerToken": "MainWorker.RateLimit.WorkerEndpointsPerMinutePerToken",
    "RateOtherPerMinutePerUser": "MainWorker.RateLimit.OtherAuthenticatedPerMinutePerUser",
    "PushUpdateMaxRunSeconds": "WorkerPushUpdate.MaxRunDurationSeconds",
    "PushUpdateHandoffTimeoutSec": "WorkerPushUpdate.HandoffTimeoutSeconds",
    "PushUpdateRetentionDays": "WorkerPushUpdate.InstructionRetentionDays",
    "PushUpdateIssuedSkewSec": "WorkerPushUpdate.IssuedSkewSeconds",
    "SelfUpdateRedirectStaleHours": "MainWorker.SelfUpdate.RedirectStaleHours",
    "BootstrapRetryBackoffSec": "MainWorker.Bootstrap.RetryBackoffSeconds",
    "BootstrapRetryMaxAttempts": "MainWorker.Bootstrap.RetryMaxAttempts",
    "CacheCompanyToWorkerTtlSeconds": "MainWorker.Cache.CompanyToWorkerTtlSeconds",
    "CacheWorkerRegistryTtlSeconds": "MainWorker.Cache.WorkerRegistryTtlSeconds",
    # CacheRecentCompanyPerUserTtlSeconds intentionally omitted from T3
    # parity: §2.16 binds its default to MainSessionTtlSeconds (not a
    # numeric literal), so seed-vs-prose equality is undefined. §4.1
    # alias map still requires the row.
    # ── v2.0.0 backup-tier materialization (lifts the prior
    # MainWorker.Backup.* T3 waiver per §4 line 272 / audit-09 §2.1) ──
    "Backup.Enabled": "MainWorker.Backup.Enabled",
    "Backup.MaxBackupsPerPrimary": "MainWorker.Backup.MaxBackupsPerPrimary",
    "Backup.LagWarningSeconds": "MainWorker.Backup.LagWarningSeconds",
    "Backup.HeartbeatIntervalSeconds": "MainWorker.Backup.HeartbeatIntervalSeconds",
    "Backup.SyncIntervalSeconds": "MainWorker.Backup.SyncIntervalSeconds",
    "Backup.MaxRowsPerEnvelope": "MainWorker.Backup.MaxRowsPerEnvelope",
    "Backup.TombstoneRetentionSeconds": "MainWorker.Backup.TombstoneRetentionSeconds",
    "Backup.LogRetentionSeconds": "MainWorker.Backup.LogRetentionSeconds",
    "Backup.QuarantineCompactionOverrideSeconds": "MainWorker.Backup.QuarantineCompactionOverrideSeconds",
    "Backup.MaxKeyAgeSeconds": "MainWorker.Backup.MaxKeyAgeSeconds",
    "Backup.RotationAckTimeoutSeconds": "MainWorker.Backup.RotationAckTimeoutSeconds",
    "Backup.RotationActivationDelaySeconds": "MainWorker.Backup.RotationActivationDelaySeconds",
    "Backup.RetiredKeyGraceSeconds": "MainWorker.Backup.RetiredKeyGraceSeconds",
    "Backup.RsaKeySizeBits": "MainWorker.Backup.RsaKeySizeBits",
    "Backup.Endpoint.IncrementalDiffTimeoutSeconds": "MainWorker.Backup.Endpoint.IncrementalDiffTimeoutSeconds",
    "Backup.Endpoint.RotateKeysTimeoutSeconds": "MainWorker.Backup.Endpoint.RotateKeysTimeoutSeconds",
    "Backup.Endpoint.RestoreByDateTimeoutSeconds": "MainWorker.Backup.Endpoint.RestoreByDateTimeoutSeconds",
    "Backup.Endpoint.SnapshotsTimeoutSeconds": "MainWorker.Backup.Endpoint.SnapshotsTimeoutSeconds",
    "Backup.Endpoint.HealthTimeoutSeconds": "MainWorker.Backup.Endpoint.HealthTimeoutSeconds",
    "Backup.Apply.MaxRetriesPerEnvelope": "MainWorker.Backup.Apply.MaxRetriesPerEnvelope",
    "Backup.Apply.TransactionTimeoutSeconds": "MainWorker.Backup.Apply.TransactionTimeoutSeconds",
    "Backup.Apply.DeadLetterRetentionDays": "MainWorker.Backup.Apply.DeadLetterRetentionDays",
    "Backup.Apply.IdempotencyRowRetentionDays": "MainWorker.Backup.Apply.IdempotencyRowRetentionDays",
    "Backup.SnapshotRetentionDays": "MainWorker.Backup.SnapshotRetentionDays",
    "Backup.Snapshot.BuildHourUtc": "MainWorker.Backup.Snapshot.BuildHourUtc",
    "Backup.Snapshot.QuiesceTimeoutSeconds": "MainWorker.Backup.Snapshot.QuiesceTimeoutSeconds",
    "Backup.Snapshot.MaxBuildSeconds": "MainWorker.Backup.Snapshot.MaxBuildSeconds",
    "Backup.Restore.PrimaryAckTimeoutSeconds": "MainWorker.Backup.Restore.PrimaryAckTimeoutSeconds",
}




def collect_catalogue_defaults() -> dict[str, str]:
    # T3 catalogue parser: only §2 rows count. Skipping §4.1's alias map
    # (which reuses ROW_RE shape) prevents the alias map's seed-short-name
    # cells from overwriting real defaults. Lifts the silent T3 waiver
    # exposed by the v2.0.0 backup-tier materialization.
    out: dict[str, str] = {}
    in_fence = False
    in_catalogue = False
    for raw in read_lines(TUNABLES_FILE):
        in_fence = toggle_fence(raw, in_fence)
        in_catalogue = _is_inside_catalogue(raw, in_catalogue)
        if in_fence or not in_catalogue:
            continue
        match = ROW_RE.match(raw)
        if not match:
            continue
        out[match.group(1)] = normalize_default(match.group(2))
    return out


def normalize_default(rest_of_row: str) -> str:
    cells = rest_of_row.split("|")
    cell = cells[0] if cells else ""
    bold = DEFAULT_BOLD_RE.search(cell)
    if bold:
        return bold.group(1).strip()
    tick = DEFAULT_TICK_RE.search(cell)
    if tick:
        return tick.group(1).strip()
    return cell.strip().split()[0] if cell.strip() else ""


def rule_t3() -> list[str]:
    seed = collect_seed_defaults()
    cat = collect_catalogue_defaults()
    out: list[str] = []
    for short, key in SEED_TO_KEY.items():
        out.extend(diff_seed_row(short, key, seed, cat))
    return out


def diff_seed_row(
    short: str, key: str, seed: dict[str, str], cat: dict[str, str]
) -> list[str]:
    seed_val = seed.get(short)
    cat_val = cat.get(key)
    if seed_val is None:
        return [f"§4 missing seed entry `{short}` (expected for `{key}`)"]
    if cat_val is None:
        return [f"§2 missing catalogue row `{key}` (referenced by §4 `{short}`)"]
    if values_match(seed_val, cat_val):
        return []
    return [f"default mismatch: §4 `{short}`={seed_val!r} vs §2 `{key}`={cat_val!r}"]


def values_match(seed_val: str, cat_val: str) -> bool:
    if seed_val == cat_val:
        return True
    return _strip_array_chars(seed_val) == _strip_array_chars(cat_val)


def _strip_array_chars(value: str) -> str:
    return value.replace("[", "").replace("]", "").replace(" ", "").replace(",", "")


SESSION_TTL_KEY = "MainWorker.Auth.MainSessionTtlSeconds"
SESSION_ABS_KEY = "MainWorker.Auth.MainSessionAbsoluteMaxSeconds"
SESSION_TTL_SEED = "MainSessionTtlSeconds"
SESSION_ABS_SEED = "MainSessionAbsoluteMaxSeconds"


_DURATION_PATTERN = re.compile(r"(\d+(?:\.\d+)?)\s*([a-zA-Z]*)")
_UNIT_TO_SECONDS = {
    "": 1, "s": 1, "sec": 1, "secs": 1, "second": 1, "seconds": 1,
    "ms": 0.001,
    "m": 60, "min": 60, "mins": 60, "minute": 60, "minutes": 60,
    "h": 3600, "hr": 3600, "hrs": 3600, "hour": 3600, "hours": 3600,
    "d": 86400, "day": 86400, "days": 86400,
}


def _scrub_for_duration(raw: str) -> str:
    # Strip markdown bold/italic and grouping punctuation that hugs digits/units.
    cleaned = raw
    for ch in ("*", "`", "(", ")", "[", "]", ","):
        cleaned = cleaned.replace(ch, " ")
    return cleaned.strip()


# parse_seconds: extract whole seconds from spec values. Unchanged on existing
# strings (bare ints like `28800`, bolded `**28800** (8h)`). Hardened to accept
# `28800s`, `1h`, `12 hours`, `500ms`, `8 h`, `1 day`, code-fenced, slashed.
# First numeric token wins (matches prior head.split()[0]); attached unit is
# honoured; no unit → seconds.
def parse_seconds(raw: str) -> int | None:
    if not raw:
        return None
    cleaned = _scrub_for_duration(raw)
    if not cleaned:
        return None
    match = _DURATION_PATTERN.search(cleaned)
    if match is None:
        return None
    return _coerce_to_whole_seconds(match.group(1), match.group(2).lower())


def _coerce_to_whole_seconds(number_text: str, unit_text: str) -> int | None:
    if unit_text not in _UNIT_TO_SECONDS:
        return None
    multiplier = _UNIT_TO_SECONDS[unit_text]
    seconds_float = float(number_text) * multiplier
    if seconds_float < 0:
        return None
    return int(seconds_float)


def rule_t4_pair(label: str, ttl_raw: str | None, abs_raw: str | None) -> list[str]:
    if ttl_raw is None:
        return [f"{label} missing `{SESSION_TTL_KEY}` value"]
    if abs_raw is None:
        return [f"{label} missing `{SESSION_ABS_KEY}` value"]
    ttl = parse_seconds(ttl_raw)
    cap = parse_seconds(abs_raw)
    if ttl is None or cap is None:
        return [f"{label} non-numeric session TTL: ttl={ttl_raw!r} cap={abs_raw!r}"]
    if cap >= ttl:
        return []
    return [f"{label} invariant violated: AbsoluteMax({cap}) < Ttl({ttl})"]


def rule_t4() -> list[str]:
    cat = collect_catalogue_defaults()
    seed = collect_seed_defaults()
    out: list[str] = []
    out.extend(rule_t4_pair("§2", cat.get(SESSION_TTL_KEY), cat.get(SESSION_ABS_KEY)))
    out.extend(rule_t4_pair("§4", seed.get(SESSION_TTL_SEED), seed.get(SESSION_ABS_SEED)))
    return out


def main() -> int:
    file_must_exist(TUNABLES_FILE)
    failures: list[str] = []
    failures.extend(prefix_each("T1", rule_t1()))
    failures.extend(prefix_each("T2", rule_t2()))
    failures.extend(prefix_each("T3", rule_t3()))
    failures.extend(prefix_each("T4", rule_t4()))
    return report(failures)


def prefix_each(rule: str, items: list[str]) -> list[str]:
    return [f"[{rule}] {item}" for item in items]


def report(failures: list[str]) -> int:
    if not failures:
        print("check-tunable-constants: OK")
        return 0
    for line in failures:
        print(line, file=sys.stderr)
    print(f"check-tunable-constants: {len(failures)} violation(s)", file=sys.stderr)
    return 1


if __name__ == "__main__":
    sys.exit(main())
