#!/usr/bin/env python3
"""
03-ai-scripts/46-agent-sqlite-task-manager.py
=============================================
Antigravity Multi-Agent SQLite Task Manager & Crash Forensics Engine.

Provides ACID-compliant, concurrency-safe task coordination and crash forensics
for multi-agent runs (V6 workflow). Replaces fragile markdown file editing with
WAL-mode SQLite micro-transactions, preventing Windows OS file lock collisions.

Capabilities:
  1. Deterministic task slug generation from prompt / task name.
  2. Automatic scan of .ai-memory/temp-agents/ for existing / similar runs to resume.
  3. Run-scoped SQLite database initialization in .ai-memory/temp-agents/<nn>-<slug>/agent-task.db.
  4. Atomic subtask claiming, status transitions, and in-flight action logging.
  5. Crash forensics: pinpointing exactly which agent crashed, what file it was touching,
     and what action it was executing when it failed.

Usage:
  # 1. Initialize task run (or inspect prior matching run):
  python 03-ai-scripts/46-agent-sqlite-task-manager.py init --name "validate regex syntax and prevent nil fallback"

  # 2. Add decomposed subtasks:
  python 03-ai-scripts/46-agent-sqlite-task-manager.py add-subtasks --db <db> --tasks-json '[...]'

  # 3. Worker claims next subtask:
  python 03-ai-scripts/46-agent-sqlite-task-manager.py claim --db <db> --agent "Worker 01"

  # 4. Worker logs in-flight action before touching a file:
  python 03-ai-scripts/46-agent-sqlite-task-manager.py log-action --db <db> --subtask-id 1 --agent "Worker 01" --action "write_to_file" --file "pkg/aum/regex.go" --details "updating parser"

  # 5. Worker completes subtask:
  python 03-ai-scripts/46-agent-sqlite-task-manager.py complete --db <db> --subtask-id 1 --agent "Worker 01" --evidence "PASS exit 0"

  # 6. Diagnose crashes or inspect status:
  python 03-ai-scripts/46-agent-sqlite-task-manager.py diagnose --db <db>
  python 03-ai-scripts/46-agent-sqlite-task-manager.py status --db <db>
"""

import argparse
import datetime
import difflib
import json
import os
import re
import sqlite3
import sys
from pathlib import Path
from typing import Any, Dict, List, Optional, Tuple


def get_repo_root() -> Path:
    """Resolve repository root directory."""
    current = Path(__file__).resolve().parent
    while current != current.parent:
        if (current / ".git").is_dir() or (current / "agents.md").is_file():
            return current
        current = current.parent
    return Path.cwd()


REPO_ROOT = get_repo_root()
TEMP_AGENTS_DIR = REPO_ROOT / ".ai-memory" / "temp-agents"


def slugify(text: str) -> str:
    """Generate a clean, deterministic, lowercase kebab-case slug."""
    text = text.lower().strip()
    text = re.sub(r"[^\w\s-]", "", text)
    slug = re.sub(r"[\s_-]+", "-", text).strip("-")
    words = slug.split("-")
    if len(words) > 6:
        slug = "-".join(words[:6])
    return slug or "task"


def get_db_connection(db_path: Path) -> sqlite3.Connection:
    """Open SQLite connection with WAL mode and busy timeout."""
    db_path.parent.mkdir(parents=True, exist_ok=True)
    conn = sqlite3.connect(str(db_path), timeout=10.0)
    conn.row_factory = sqlite3.Row
    conn.execute("PRAGMA journal_mode = WAL;")
    conn.execute("PRAGMA busy_timeout = 5000;")
    conn.execute("PRAGMA synchronous = NORMAL;")
    return conn


def init_schema(conn: sqlite3.Connection) -> None:
    """Initialize PascalCase SQLite tables with positive boolean columns."""
    with conn:
        conn.executescript("""
        CREATE TABLE IF NOT EXISTS ParentTask (
            ParentTaskId INTEGER PRIMARY KEY AUTOINCREMENT,
            TaskName TEXT NOT NULL,
            TaskSlug TEXT NOT NULL,
            RunDirectory TEXT NOT NULL,
            Status TEXT NOT NULL DEFAULT 'ACTIVE',
            IsActive INTEGER NOT NULL DEFAULT 1,
            HasCompleted INTEGER NOT NULL DEFAULT 0,
            TotalStepsBudget INTEGER NOT NULL DEFAULT 300,
            CurrentStep INTEGER NOT NULL DEFAULT 1,
            Notes TEXT NULL,
            CreatedAt TEXT NOT NULL,
            UpdatedAt TEXT NOT NULL
        );

        CREATE TABLE IF NOT EXISTS Subtask (
            SubtaskId INTEGER PRIMARY KEY AUTOINCREMENT,
            ParentTaskId INTEGER NOT NULL,
            TaskCode TEXT NOT NULL,
            Title TEXT NOT NULL,
            AssignedAgentRole TEXT NULL,
            OwnedFilesJson TEXT NOT NULL DEFAULT '[]',
            Status TEXT NOT NULL DEFAULT 'PENDING',
            IsBlocked INTEGER NOT NULL DEFAULT 0,
            HasCompleted INTEGER NOT NULL DEFAULT 0,
            Evidence TEXT NULL,
            CreatedAt TEXT NOT NULL,
            UpdatedAt TEXT NOT NULL,
            FOREIGN KEY (ParentTaskId) REFERENCES ParentTask(ParentTaskId)
        );

        CREATE TABLE IF NOT EXISTS AgentActionLog (
            ActionLogId INTEGER PRIMARY KEY AUTOINCREMENT,
            SubtaskId INTEGER NOT NULL,
            AgentRole TEXT NOT NULL,
            ActionType TEXT NOT NULL,
            TargetFile TEXT NULL,
            ActionDetails TEXT NOT NULL,
            Status TEXT NOT NULL DEFAULT 'IN_PROGRESS',
            CreatedAt TEXT NOT NULL,
            FOREIGN KEY (SubtaskId) REFERENCES Subtask(SubtaskId)
        );

        CREATE INDEX IF NOT EXISTS idx_subtask_status ON Subtask(Status);
        CREATE INDEX IF NOT EXISTS idx_action_subtask ON AgentActionLog(SubtaskId);
        """)


def find_similar_existing_run(target_slug: str) -> Optional[Tuple[Path, float]]:
    """Scan .ai-memory/temp-agents/ for existing runs with matching or similar slug."""
    if not TEMP_AGENTS_DIR.exists():
        return None

    best_match: Optional[Path] = None
    highest_ratio = 0.0

    for item in TEMP_AGENTS_DIR.iterdir():
        if not item.is_dir():
            continue
        # Strip numeric prefix: '01-aum-validate-regex' -> 'aum-validate-regex'
        m = re.match(r"^\d+-(.+)$", item.name)
        existing_slug = m.group(1) if m else item.name

        if existing_slug == target_slug:
            return (item, 1.0)

        # Check substring match
        if target_slug in existing_slug or existing_slug in target_slug:
            ratio = 0.85
        else:
            ratio = difflib.SequenceMatcher(None, target_slug, existing_slug).ratio()

        if ratio > highest_ratio and ratio >= 0.70:
            highest_ratio = ratio
            best_match = item

    if best_match is not None:
        return (best_match, highest_ratio)
    return None


def get_next_run_dir(slug: str) -> Path:
    """Find next numeric prefix in .ai-memory/temp-agents/."""
    TEMP_AGENTS_DIR.mkdir(parents=True, exist_ok=True)
    existing_nums = []
    for item in TEMP_AGENTS_DIR.iterdir():
        if item.is_dir():
            m = re.match(r"^(\d+)-", item.name)
            if m:
                existing_nums.append(int(m.group(1)))
    next_num = max(existing_nums, default=0) + 1
    return TEMP_AGENTS_DIR / f"{next_num:02d}-{slug}"


def now_iso() -> str:
    """Return current UTC timestamp in ISO 8601."""
    return datetime.datetime.now(datetime.timezone.utc).isoformat()


# ---------------------------------------------------------------------------
# CLI Command Implementations
# ---------------------------------------------------------------------------

def cmd_init(args: argparse.Namespace) -> int:
    """Initialize task run, checking for existing runs to resume."""
    slug = slugify(args.name)
    existing = find_similar_existing_run(slug)

    if existing is not None:
        existing_dir, similarity = existing
        db_file = existing_dir / "agent-task.db"
        if db_file.exists():
            conn = get_db_connection(db_file)
            init_schema(conn)
            cur = conn.cursor()
            cur.execute("SELECT * FROM ParentTask ORDER BY ParentTaskId DESC LIMIT 1;")
            parent = cur.fetchone()
            if parent is not None:
                # Run crash diagnostics on this existing run
                diagnostics = diagnose_db(conn)
                out = {
                    "action": "RESUME_FOUND",
                    "similarity": round(similarity, 3),
                    "runDirectory": str(existing_dir.relative_to(REPO_ROOT)).replace("\\", "/"),
                    "databasePath": str(db_file.relative_to(REPO_ROOT)).replace("\\", "/"),
                    "parentTaskId": parent["ParentTaskId"],
                    "slug": parent["TaskSlug"],
                    "status": parent["Status"],
                    "isActive": bool(parent["IsActive"]),
                    "hasCompleted": bool(parent["HasCompleted"]),
                    "diagnostics": diagnostics,
                }
                print(json.dumps(out, indent=2))
                conn.close()
                return 0
            conn.close()

    # Create new run directory and DB
    run_dir = get_next_run_dir(slug)
    run_dir.mkdir(parents=True, exist_ok=True)
    db_file = run_dir / "agent-task.db"

    conn = get_db_connection(db_file)
    init_schema(conn)

    rel_run_dir = str(run_dir.relative_to(REPO_ROOT)).replace("\\", "/")
    rel_db_path = str(db_file.relative_to(REPO_ROOT)).replace("\\", "/")
    created_at = now_iso()

    with conn:
        cur = conn.cursor()
        cur.execute("""
            INSERT INTO ParentTask (
                TaskName, TaskSlug, RunDirectory, Status, IsActive, HasCompleted,
                TotalStepsBudget, CurrentStep, CreatedAt, UpdatedAt
            ) VALUES (?, ?, ?, 'ACTIVE', 1, 0, ?, 1, ?, ?);
        """, (args.name, slug, rel_run_dir, args.budget, created_at, created_at))
        parent_id = cur.lastrowid

    conn.close()

    out = {
        "action": "INITIALIZED",
        "parentTaskId": parent_id,
        "slug": slug,
        "runDirectory": rel_run_dir,
        "databasePath": rel_db_path,
        "totalStepsBudget": args.budget,
        "status": "ACTIVE",
    }
    print(json.dumps(out, indent=2))
    return 0


def cmd_add_subtasks(args: argparse.Namespace) -> int:
    """Populate subtasks for the task run."""
    db_path = Path(args.db)
    if not db_path.is_absolute():
        db_path = REPO_ROOT / db_path

    if not db_path.exists():
        print(json.dumps({"error": f"Database not found at {args.db}"}), file=sys.stderr)
        return 1

    tasks_raw = args.tasks_json
    if not tasks_raw and args.tasks_file:
        file_path = Path(args.tasks_file)
        if not file_path.is_absolute():
            file_path = REPO_ROOT / file_path
        tasks_raw = file_path.read_text(encoding="utf-8")

    if not tasks_raw:
        print(json.dumps({"error": "No subtasks provided via --tasks-json or --tasks-file"}), file=sys.stderr)
        return 1

    try:
        subtasks_data = json.loads(tasks_raw)
    except json.JSONDecodeError as exc:
        print(json.dumps({"error": f"Invalid JSON for subtasks: {exc}"}), file=sys.stderr)
        return 1

    conn = get_db_connection(db_path)
    init_schema(conn)

    cur = conn.cursor()
    cur.execute("SELECT ParentTaskId FROM ParentTask ORDER BY ParentTaskId DESC LIMIT 1;")
    row = cur.fetchone()
    if not row:
        print(json.dumps({"error": "ParentTask record missing in database"}), file=sys.stderr)
        conn.close()
        return 1

    parent_id = args.parent_id if args.parent_id else row["ParentTaskId"]
    created_at = now_iso()
    inserted_ids = []

    with conn:
        for idx, task in enumerate(subtasks_data, start=1):
            code = task.get("code") or f"Task-{idx:02d}"
            title = task.get("title", f"Subtask {idx}")
            agent = task.get("agent_role") or task.get("assigned_agent")
            owned = json.dumps(task.get("owned_files", []))
            cur.execute("""
                INSERT INTO Subtask (
                    ParentTaskId, TaskCode, Title, AssignedAgentRole,
                    OwnedFilesJson, Status, IsBlocked, HasCompleted,
                    CreatedAt, UpdatedAt
                ) VALUES (?, ?, ?, ?, ?, 'PENDING', 0, 0, ?, ?);
            """, (parent_id, code, title, agent, owned, created_at, created_at))
            inserted_ids.append(cur.lastrowid)

    conn.close()
    print(json.dumps({
        "status": "SUCCESS",
        "subtasksCreated": len(inserted_ids),
        "subtaskIds": inserted_ids,
    }, indent=2))
    return 0


def cmd_claim(args: argparse.Namespace) -> int:
    """Worker atomically claims next available subtask."""
    db_path = Path(args.db)
    if not db_path.is_absolute():
        db_path = REPO_ROOT / db_path

    conn = get_db_connection(db_path)
    init_schema(conn)

    with conn:
        cur = conn.cursor()
        # Find next pending task (preferring assigned or unassigned)
        if args.agent:
            cur.execute("""
                SELECT * FROM Subtask
                WHERE Status = 'PENDING'
                  AND (AssignedAgentRole = ? OR AssignedAgentRole IS NULL OR AssignedAgentRole = '')
                ORDER BY SubtaskId ASC LIMIT 1;
            """, (args.agent,))
        else:
            cur.execute("""
                SELECT * FROM Subtask
                WHERE Status = 'PENDING'
                ORDER BY SubtaskId ASC LIMIT 1;
            """)

        task = cur.fetchone()
        if not task:
            conn.close()
            print(json.dumps({"status": "NO_TASKS_AVAILABLE"}))
            return 0

        subtask_id = task["SubtaskId"]
        assigned_role = args.agent or task["AssignedAgentRole"] or "Worker"
        updated_at = now_iso()

        cur.execute("""
            UPDATE Subtask
            SET Status = 'IN_PROGRESS',
                AssignedAgentRole = ?,
                UpdatedAt = ?
            WHERE SubtaskId = ?;
        """, (assigned_role, updated_at, subtask_id))

        # Log start action
        cur.execute("""
            INSERT INTO AgentActionLog (
                SubtaskId, AgentRole, ActionType, ActionDetails, Status, CreatedAt
            ) VALUES (?, ?, 'START', 'Subtask claimed and started', 'IN_PROGRESS', ?);
        """, (subtask_id, assigned_role, updated_at))

    owned_files = json.loads(task["OwnedFilesJson"]) if task["OwnedFilesJson"] else []
    out = {
        "status": "CLAIMED",
        "subtaskId": subtask_id,
        "taskCode": task["TaskCode"],
        "title": task["Title"],
        "assignedAgentRole": assigned_role,
        "ownedFiles": owned_files,
    }
    conn.close()
    print(json.dumps(out, indent=2))
    return 0


def cmd_log_action(args: argparse.Namespace) -> int:
    """Log in-flight agent action before touching a file or running a check."""
    db_path = Path(args.db)
    if not db_path.is_absolute():
        db_path = REPO_ROOT / db_path

    conn = get_db_connection(db_path)
    init_schema(conn)

    created_at = now_iso()
    details = args.details or f"Executing {args.action} on {args.file or 'workspace'}"

    with conn:
        cur = conn.cursor()
        cur.execute("""
            INSERT INTO AgentActionLog (
                SubtaskId, AgentRole, ActionType, TargetFile, ActionDetails, Status, CreatedAt
            ) VALUES (?, ?, ?, ?, ?, 'IN_PROGRESS', ?);
        """, (args.subtask_id, args.agent, args.action, args.file, details, created_at))
        action_id = cur.lastrowid

        cur.execute("UPDATE Subtask SET UpdatedAt = ? WHERE SubtaskId = ?;", (created_at, args.subtask_id))

    conn.close()
    print(json.dumps({"status": "LOGGED", "actionLogId": action_id}))
    return 0


def cmd_complete(args: argparse.Namespace) -> int:
    """Mark subtask completed with evidence."""
    db_path = Path(args.db)
    if not db_path.is_absolute():
        db_path = REPO_ROOT / db_path

    conn = get_db_connection(db_path)
    init_schema(conn)

    updated_at = now_iso()
    with conn:
        cur = conn.cursor()
        cur.execute("""
            UPDATE Subtask
            SET Status = 'DONE',
                HasCompleted = 1,
                IsBlocked = 0,
                Evidence = ?,
                UpdatedAt = ?
            WHERE SubtaskId = ?;
        """, (args.evidence or "PASS exit 0", updated_at, args.subtask_id))

        cur.execute("""
            INSERT INTO AgentActionLog (
                SubtaskId, AgentRole, ActionType, ActionDetails, Status, CreatedAt
            ) VALUES (?, ?, 'COMPLETE', 'Subtask verified and completed', 'DONE', ?);
        """, (args.subtask_id, args.agent or "Worker", updated_at))

    conn.close()
    print(json.dumps({"status": "COMPLETED", "subtaskId": args.subtask_id}))
    return 0


def cmd_fail(args: argparse.Namespace) -> int:
    """Mark subtask failed or blocked."""
    db_path = Path(args.db)
    if not db_path.is_absolute():
        db_path = REPO_ROOT / db_path

    conn = get_db_connection(db_path)
    init_schema(conn)

    updated_at = now_iso()
    with conn:
        cur = conn.cursor()
        cur.execute("""
            UPDATE Subtask
            SET Status = 'FAILED',
                IsBlocked = 1,
                Evidence = ?,
                UpdatedAt = ?
            WHERE SubtaskId = ?;
        """, (args.reason or "Failed", updated_at, args.subtask_id))

        cur.execute("""
            INSERT INTO AgentActionLog (
                SubtaskId, AgentRole, ActionType, ActionDetails, Status, CreatedAt
            ) VALUES (?, ?, 'FAIL', ?, 'FAILED', ?);
        """, (args.subtask_id, args.agent or "Worker", args.reason or "Task failed", updated_at))

    conn.close()
    print(json.dumps({"status": "FAILED", "subtaskId": args.subtask_id}))
    return 0


def diagnose_db(conn: sqlite3.Connection) -> Dict[str, Any]:
    """Diagnose uncompleted or crashed agent actions."""
    cur = conn.cursor()

    cur.execute("""
        SELECT s.*,
               a.ActionType AS LastActionType,
               a.TargetFile AS LastTargetFile,
               a.ActionDetails AS LastActionDetails,
               a.CreatedAt AS LastActionTimestamp
        FROM Subtask s
        LEFT JOIN AgentActionLog a ON a.ActionLogId = (
            SELECT MAX(ActionLogId) FROM AgentActionLog WHERE SubtaskId = s.SubtaskId
        )
        WHERE s.Status = 'IN_PROGRESS'
        ORDER BY s.SubtaskId ASC;
    """)
    in_progress = cur.fetchall()

    crashes = []
    for row in in_progress:
        crashes.append({
            "subtaskId": row["SubtaskId"],
            "taskCode": row["TaskCode"],
            "title": row["Title"],
            "assignedAgentRole": row["AssignedAgentRole"],
            "lastActionType": row["LastActionType"],
            "lastTargetFile": row["LastTargetFile"],
            "lastActionDetails": row["LastActionDetails"],
            "lastTimestamp": row["LastActionTimestamp"],
            "diagnosis": (
                f"Agent '{row['AssignedAgentRole']}' was in progress on {row['TaskCode']} ('{row['Title']}'). "
                f"Last logged action: '{row['LastActionType']}' targeting '{row['LastTargetFile'] or 'workspace'}'. "
                "Agent terminated before completing. Likely caused by tool execution crash, OS file lock, or turn timeout."
            )
        })

    cur.execute("SELECT COUNT(*) AS c FROM Subtask WHERE Status = 'DONE';")
    done_count = cur.fetchone()["c"]

    cur.execute("SELECT COUNT(*) AS c FROM Subtask WHERE Status = 'PENDING';")
    pending_count = cur.fetchone()["c"]

    cur.execute("SELECT COUNT(*) AS c FROM Subtask WHERE Status = 'FAILED';")
    failed_count = cur.fetchone()["c"]

    return {
        "hasCrashesDetected": len(crashes) > 0,
        "crashedAgents": crashes,
        "counts": {
            "done": done_count,
            "pending": pending_count,
            "inProgress": len(in_progress),
            "failed": failed_count,
        }
    }


def cmd_diagnose(args: argparse.Namespace) -> int:
    """Diagnose crashes and report forensics."""
    db_path = Path(args.db)
    if not db_path.is_absolute():
        db_path = REPO_ROOT / db_path

    if not db_path.exists():
        print(json.dumps({"error": f"Database not found at {args.db}"}), file=sys.stderr)
        return 1

    conn = get_db_connection(db_path)
    init_schema(conn)

    diagnostics = diagnose_db(conn)
    conn.close()

    print(json.dumps(diagnostics, indent=2))
    return 0


def cmd_status(args: argparse.Namespace) -> int:
    """Output high-level progress status."""
    db_path = Path(args.db)
    if not db_path.is_absolute():
        db_path = REPO_ROOT / db_path

    if not db_path.exists():
        print(json.dumps({"error": f"Database not found at {args.db}"}), file=sys.stderr)
        return 1

    conn = get_db_connection(db_path)
    init_schema(conn)

    cur = conn.cursor()
    cur.execute("SELECT * FROM ParentTask ORDER BY ParentTaskId DESC LIMIT 1;")
    parent = cur.fetchone()

    cur.execute("SELECT SubtaskId, TaskCode, Title, AssignedAgentRole, Status, Evidence FROM Subtask ORDER BY SubtaskId ASC;")
    subtasks = [dict(row) for row in cur.fetchall()]

    counts: Dict[str, int] = {}
    for st in subtasks:
        s = st["Status"]
        counts[s] = counts.get(s, 0) + 1

    total = len(subtasks)
    done = counts.get("DONE", 0)
    percent = round((done / total * 100), 1) if total > 0 else 0.0

    conn.close()

    out = {
        "taskSlug": parent["TaskSlug"] if parent else "unknown",
        "taskName": parent["TaskName"] if parent else "unknown",
        "status": parent["Status"] if parent else "unknown",
        "totalSubtasks": total,
        "completed": done,
        "pending": counts.get("PENDING", 0),
        "inProgress": counts.get("IN_PROGRESS", 0),
        "failed": counts.get("FAILED", 0),
        "percentComplete": percent,
        "subtasks": subtasks,
    }
    print(json.dumps(out, indent=2))
    return 0


def main() -> int:
    parser = argparse.ArgumentParser(description="Antigravity Multi-Agent SQLite Task Manager & Crash Forensics Engine")
    subparsers = parser.add_subparsers(dest="command", required=True)

    # init
    p_init = subparsers.add_parser("init", help="Initialize task run or detect resume/crash state")
    p_init.add_argument("--name", required=True, help="Task name or prompt description")
    p_init.add_argument("--budget", type=int, default=300, help="Total step budget (default 300)")

    # add-subtasks
    p_add = subparsers.add_parser("add-subtasks", help="Populate decomposed subtasks")
    p_add.add_argument("--db", required=True, help="Path to SQLite database")
    p_add.add_argument("--parent-id", type=int, help="Optional ParentTaskId")
    p_add.add_argument("--tasks-json", help="JSON string representing list of subtasks")
    p_add.add_argument("--tasks-file", help="Path to JSON file representing list of subtasks")

    # claim
    p_claim = subparsers.add_parser("claim", help="Atomically claim next available subtask")
    p_claim.add_argument("--db", required=True, help="Path to SQLite database")
    p_claim.add_argument("--agent", required=True, help="Agent role name (e.g. Worker 01)")

    # log-action
    p_log = subparsers.add_parser("log-action", help="Log an in-flight agent action before touching a file")
    p_log.add_argument("--db", required=True, help="Path to SQLite database")
    p_log.add_argument("--subtask-id", type=int, required=True, help="SubtaskId")
    p_log.add_argument("--agent", required=True, help="Agent role name")
    p_log.add_argument("--action", required=True, help="Action type (e.g. write_to_file, run_linter)")
    p_log.add_argument("--file", help="Target relative file path being touched")
    p_log.add_argument("--details", help="Human-readable details of action")

    # complete
    p_comp = subparsers.add_parser("complete", help="Mark subtask completed with evidence")
    p_comp.add_argument("--db", required=True, help="Path to SQLite database")
    p_comp.add_argument("--subtask-id", type=int, required=True, help="SubtaskId")
    p_comp.add_argument("--agent", help="Agent role name")
    p_comp.add_argument("--evidence", help="Evidence string (e.g. PASS exit 0)")

    # fail
    p_fail = subparsers.add_parser("fail", help="Mark subtask failed")
    p_fail.add_argument("--db", required=True, help="Path to SQLite database")
    p_fail.add_argument("--subtask-id", type=int, required=True, help="SubtaskId")
    p_fail.add_argument("--agent", help="Agent role name")
    p_fail.add_argument("--reason", help="Failure reason or blocker details")

    # diagnose
    p_diag = subparsers.add_parser("diagnose", help="Diagnose crashed or abandoned agent tasks")
    p_diag.add_argument("--db", required=True, help="Path to SQLite database")

    # status
    p_stat = subparsers.add_parser("status", help="Get task execution status summary")
    p_stat.add_argument("--db", required=True, help="Path to SQLite database")

    args = parser.parse_args()

    dispatch = {
        "init": cmd_init,
        "add-subtasks": cmd_add_subtasks,
        "claim": cmd_claim,
        "log-action": cmd_log_action,
        "complete": cmd_complete,
        "fail": cmd_fail,
        "diagnose": cmd_diagnose,
        "status": cmd_status,
    }

    return dispatch[args.command](args)


if __name__ == "__main__":
    sys.exit(main())
