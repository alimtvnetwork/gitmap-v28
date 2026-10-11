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

  # 7. Inspect SQLite database schema:
  python 03-ai-scripts/46-agent-sqlite-task-manager.py schema
  python 03-ai-scripts/46-agent-sqlite-task-manager.py schema --json
  python 03-ai-scripts/46-agent-sqlite-task-manager.py schema --ddl
"""

import argparse
import datetime
import difflib
import json
import os
import re
import socket
import sqlite3
import sys
import time
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


def get_machine_id() -> str:
    """Resolve deterministic machine identity in sanitized lowercase kebab-case."""
    env_id = os.environ.get("GITMAP_MACHINE_ID", "").strip()
    if env_id:
        cleaned = re.sub(r"[^a-z0-9_-]", "-", env_id.lower()).strip("-")
        return cleaned[:32] or "local-machine"

    raw = ""
    try:
        raw = socket.gethostname().strip().lower()
    except Exception:
        pass

    if not raw:
        raw = os.environ.get("COMPUTERNAME", os.environ.get("HOSTNAME", "local-machine")).strip().lower()

    cleaned = re.sub(r"[^a-z0-9_-]", "-", raw).strip("-")
    return cleaned[:32] or "local-machine"


CANONICAL_ACTION_TYPES = {"modify", "analyze", "read", "remove", "refactor", "create"}
LIFECYCLE_ACTION_TYPES = {"START", "COMPLETE", "FAIL"}

ACTION_TYPE_ALIASES: Dict[str, str] = {
    # modify aliases
    "write": "modify",
    "replace_file_content": "modify",
    "edit": "modify",
    "edit_file": "modify",
    "patch": "modify",
    "update": "modify",
    # read aliases
    "view_file": "read",
    "read_file": "read",
    "cat": "read",
    "read_memory": "read",
    "gitmap cat": "read",
    "read_resource": "read",
    # analyze aliases
    "search": "analyze",
    "find": "analyze",
    "grep": "analyze",
    "diagnose": "analyze",
    "check": "analyze",
    "inspect": "analyze",
    "lint": "analyze",
    "run_linter": "analyze",
    "run_test": "analyze",
    "analyze_file": "analyze",
    # remove aliases
    "delete": "remove",
    "delete_file": "remove",
    "drop": "remove",
    "clean": "remove",
    "rm": "remove",
    "unlink": "remove",
    "purge": "remove",
    # refactor aliases
    "restructure": "refactor",
    "optimize": "refactor",
    "reorganize": "refactor",
    "refactor_module": "refactor",
    "extract_type": "refactor",
    # create aliases
    "touch": "create",
    "create_file": "create",
    "mkdir": "create",
    "add": "create",
    "new": "create",
    "write_to_file": "create",
}


def normalize_action_type(raw_action: str, target_file: Optional[str] = None) -> Optional[str]:
    """Normalize raw tool/IDE action to canonical ActionType enum or lifecycle action."""
    action_raw = raw_action.strip()
    action_upper = action_raw.upper()
    if action_upper in LIFECYCLE_ACTION_TYPES:
        return action_upper

    action_lower = action_raw.lower()
    if action_lower in CANONICAL_ACTION_TYPES:
        return action_lower

    if action_lower == "write_to_file" and target_file:
        try:
            target_path = REPO_ROOT / target_file
            if target_path.exists():
                return "modify"
        except Exception:
            pass
        return "create"

    if action_lower in ACTION_TYPE_ALIASES:
        return ACTION_TYPE_ALIASES[action_lower]

    return None


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


SCHEMA_DDL = """
CREATE TABLE IF NOT EXISTS ParentTask (
    ParentTaskId INTEGER PRIMARY KEY AUTOINCREMENT,
    TaskName TEXT NOT NULL,
    TaskSlug TEXT NOT NULL,
    MachineId TEXT NOT NULL DEFAULT 'local-machine',
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
    MachineId TEXT NOT NULL DEFAULT 'local-machine',
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
    MachineId TEXT NOT NULL DEFAULT 'local-machine',
    AgentRole TEXT NOT NULL,
    ActionType TEXT NOT NULL,
    TargetFile TEXT NULL,
    ActionDetails TEXT NOT NULL,
    Status TEXT NOT NULL DEFAULT 'IN_PROGRESS',
    CreatedAt TEXT NOT NULL,
    FOREIGN KEY (SubtaskId) REFERENCES Subtask(SubtaskId)
);

CREATE INDEX IF NOT EXISTS idx_subtask_status ON Subtask(Status);
CREATE INDEX IF NOT EXISTS idx_subtask_machine ON Subtask(MachineId);
CREATE INDEX IF NOT EXISTS idx_action_subtask ON AgentActionLog(SubtaskId);
CREATE INDEX IF NOT EXISTS idx_action_file ON AgentActionLog(TargetFile);
CREATE INDEX IF NOT EXISTS idx_action_created ON AgentActionLog(CreatedAt);
""".strip()

SCHEMA_METADATA: Dict[str, Any] = {
    "version": "1.0",
    "conventions": {
        "naming": "PascalCase for tables and columns",
        "booleans": "Positive booleans only (Is*, Has*) storing integer 0 or 1",
        "concurrency": "WAL mode enabled with busy_timeout=5000ms",
        "paths": "Strict relative git paths only (forward slashes)",
        "machineId": "Deterministic lowercase kebab-case host identifier",
        "actions": "Strict canonical ActionType enum (modify, analyze, read, remove, refactor, create)",
    },
    "actionTypes": {
        "canonical": sorted(list(CANONICAL_ACTION_TYPES)),
        "lifecycle": sorted(list(LIFECYCLE_ACTION_TYPES)),
        "aliases": ACTION_TYPE_ALIASES,
    },
    "tables": {
        "ParentTask": {
            "description": "Represents the top-level parent task run and budget.",
            "primaryKey": "ParentTaskId",
            "columns": [
                {"name": "ParentTaskId", "type": "INTEGER", "primaryKey": True, "nullable": False, "default": None, "description": "Auto-incrementing primary key"},
                {"name": "TaskName", "type": "TEXT", "primaryKey": False, "nullable": False, "default": None, "description": "Original task title or prompt"},
                {"name": "TaskSlug", "type": "TEXT", "primaryKey": False, "nullable": False, "default": None, "description": "Deterministic kebab-case identifier"},
                {"name": "MachineId", "type": "TEXT", "primaryKey": False, "nullable": False, "default": "'local-machine'", "description": "Deterministic host machine identity"},
                {"name": "RunDirectory", "type": "TEXT", "primaryKey": False, "nullable": False, "default": None, "description": "Relative directory in .ai-memory/temp-agents/"},
                {"name": "Status", "type": "TEXT", "primaryKey": False, "nullable": False, "default": "'ACTIVE'", "description": "ACTIVE, COMPLETED, FAILED"},
                {"name": "IsActive", "type": "INTEGER", "primaryKey": False, "nullable": False, "default": "1", "isPositiveBoolean": True, "description": "1 if active, 0 if paused or finished"},
                {"name": "HasCompleted", "type": "INTEGER", "primaryKey": False, "nullable": False, "default": "0", "isPositiveBoolean": True, "description": "1 when task has completed"},
                {"name": "TotalStepsBudget", "type": "INTEGER", "primaryKey": False, "nullable": False, "default": "300", "description": "Budgeted step count (default 300)"},
                {"name": "CurrentStep", "type": "INTEGER", "primaryKey": False, "nullable": False, "default": "1", "description": "Current monotonic step count"},
                {"name": "Notes", "type": "TEXT", "primaryKey": False, "nullable": True, "default": None, "description": "Optional notes or summary"},
                {"name": "CreatedAt", "type": "TEXT", "primaryKey": False, "nullable": False, "default": None, "description": "ISO 8601 UTC timestamp"},
                {"name": "UpdatedAt", "type": "TEXT", "primaryKey": False, "nullable": False, "default": None, "description": "ISO 8601 UTC timestamp"},
            ],
            "indexes": []
        },
        "Subtask": {
            "description": "Decomposed isolated work units assigned to worker agents.",
            "primaryKey": "SubtaskId",
            "columns": [
                {"name": "SubtaskId", "type": "INTEGER", "primaryKey": True, "nullable": False, "default": None, "description": "Auto-incrementing primary key"},
                {"name": "ParentTaskId", "type": "INTEGER", "primaryKey": False, "nullable": False, "default": None, "foreignKey": "ParentTask(ParentTaskId)", "description": "Parent task foreign key"},
                {"name": "TaskCode", "type": "TEXT", "primaryKey": False, "nullable": False, "default": None, "description": "Task code, e.g. Task-01"},
                {"name": "Title", "type": "TEXT", "primaryKey": False, "nullable": False, "default": None, "description": "Subtask title description"},
                {"name": "MachineId", "type": "TEXT", "primaryKey": False, "nullable": False, "default": "'local-machine'", "description": "Deterministic host machine identity"},
                {"name": "AssignedAgentRole", "type": "TEXT", "primaryKey": False, "nullable": True, "default": None, "description": "Assigned agent role (e.g. Worker 01)"},
                {"name": "OwnedFilesJson", "type": "TEXT", "primaryKey": False, "nullable": False, "default": "'[]'", "isJson": True, "description": "JSON array of disjoint relative file paths"},
                {"name": "Status", "type": "TEXT", "primaryKey": False, "nullable": False, "default": "'PENDING'", "description": "PENDING, IN_PROGRESS, DONE, FAILED"},
                {"name": "IsBlocked", "type": "INTEGER", "primaryKey": False, "nullable": False, "default": "0", "isPositiveBoolean": True, "description": "1 if blocked by failure, 0 otherwise"},
                {"name": "HasCompleted", "type": "INTEGER", "primaryKey": False, "nullable": False, "default": "0", "isPositiveBoolean": True, "description": "1 if completed successfully, 0 otherwise"},
                {"name": "Evidence", "type": "TEXT", "primaryKey": False, "nullable": True, "default": None, "description": "Verification evidence or exit code"},
                {"name": "CreatedAt", "type": "TEXT", "primaryKey": False, "nullable": False, "default": None, "description": "ISO 8601 UTC timestamp"},
                {"name": "UpdatedAt", "type": "TEXT", "primaryKey": False, "nullable": False, "default": None, "description": "ISO 8601 UTC timestamp"},
            ],
            "indexes": [
                "idx_subtask_status ON Subtask(Status)",
                "idx_subtask_machine ON Subtask(MachineId)",
            ]
        },
        "AgentActionLog": {
            "description": "Fine-grained audit trail capturing in-flight actions for crash forensics.",
            "primaryKey": "ActionLogId",
            "columns": [
                {"name": "ActionLogId", "type": "INTEGER", "primaryKey": True, "nullable": False, "default": None, "description": "Auto-incrementing primary key"},
                {"name": "SubtaskId", "type": "INTEGER", "primaryKey": False, "nullable": False, "default": None, "foreignKey": "Subtask(SubtaskId)", "description": "Subtask foreign key"},
                {"name": "MachineId", "type": "TEXT", "primaryKey": False, "nullable": False, "default": "'local-machine'", "description": "Deterministic host machine identity"},
                {"name": "AgentRole", "type": "TEXT", "primaryKey": False, "nullable": False, "default": None, "description": "Worker role name executing the action"},
                {"name": "ActionType", "type": "TEXT", "primaryKey": False, "nullable": False, "default": None, "description": "Canonical ActionType (modify, analyze, read, remove, refactor, create, START, COMPLETE, FAIL)"},
                {"name": "TargetFile", "type": "TEXT", "primaryKey": False, "nullable": True, "default": None, "description": "Target relative file path"},
                {"name": "ActionDetails", "type": "TEXT", "primaryKey": False, "nullable": False, "default": None, "description": "Human-readable description of action"},
                {"name": "Status", "type": "TEXT", "primaryKey": False, "nullable": False, "default": "'IN_PROGRESS'", "description": "IN_PROGRESS, DONE, FAILED"},
                {"name": "CreatedAt", "type": "TEXT", "primaryKey": False, "nullable": False, "default": None, "description": "ISO 8601 UTC timestamp"},
            ],
            "indexes": [
                "idx_action_subtask ON AgentActionLog(SubtaskId)",
                "idx_action_file ON AgentActionLog(TargetFile)",
                "idx_action_created ON AgentActionLog(CreatedAt)",
            ]
        }
    },
    "payloadSpecifications": {
        "Subtask.OwnedFilesJson": {
            "description": "Serialized JSON array of strictly relative repository file paths",
            "type": "array of strings",
            "rules": [
                "Must be a valid JSON array",
                "Each entry must be a non-empty string",
                "Must be strictly relative to the repository root",
                "Zero absolute filesystem paths (C:\\ or leading /)",
                "Zero file:// or file:/// URIs"
            ],
            "example": ["03-ai-scripts/46-agent-sqlite-task-manager.py", "01-prompts/spec.md"]
        },
        "add-subtasks": {
            "description": "Input payload for populating subtasks",
            "supportedFormats": [
                "Direct list: [ { ... }, { ... } ]",
                "Object with tasks key: { \"tasks\": [ ... ] }",
                "Object with subtasks key: { \"subtasks\": [ ... ] }"
            ],
            "subtaskFields": {
                "title": {"type": "string", "required": True, "description": "Non-empty subtask title"},
                "code": {"type": "string", "required": False, "description": "Subtask identifier (default: Task-NN)"},
                "agent_role": {"type": "string", "required": False, "description": "Worker role name (e.g. Worker 01)"},
                "owned_files": {"type": "array of strings", "required": False, "description": "Disjoint relative paths owned by this task"}
            }
        }
    }
}


def ensure_schema_migrations(conn: sqlite3.Connection) -> None:
    """Apply non-destructive schema migrations for backward compatibility on existing databases."""
    cur = conn.cursor()

    cur.execute("SELECT name FROM sqlite_master WHERE type='table' AND name='ParentTask';")
    has_parent = bool(cur.fetchone())
    if has_parent:
        cur.execute("PRAGMA table_info(ParentTask);")
        columns = [row["name"] for row in cur.fetchall()]
        if "MachineId" not in columns:
            cur.execute("ALTER TABLE ParentTask ADD COLUMN MachineId TEXT NOT NULL DEFAULT 'local-machine';")

    cur.execute("SELECT name FROM sqlite_master WHERE type='table' AND name='Subtask';")
    has_subtask = bool(cur.fetchone())
    if has_subtask:
        cur.execute("PRAGMA table_info(Subtask);")
        columns = [row["name"] for row in cur.fetchall()]
        if "MachineId" not in columns:
            cur.execute("ALTER TABLE Subtask ADD COLUMN MachineId TEXT NOT NULL DEFAULT 'local-machine';")

    cur.execute("SELECT name FROM sqlite_master WHERE type='table' AND name='AgentActionLog';")
    has_action = bool(cur.fetchone())
    if has_action:
        cur.execute("PRAGMA table_info(AgentActionLog);")
        columns = [row["name"] for row in cur.fetchall()]
        if "MachineId" not in columns:
            cur.execute("ALTER TABLE AgentActionLog ADD COLUMN MachineId TEXT NOT NULL DEFAULT 'local-machine';")


def init_schema(conn: sqlite3.Connection) -> None:
    """Initialize PascalCase SQLite tables with positive boolean columns."""
    with conn:
        ensure_schema_migrations(conn)
        conn.executescript(SCHEMA_DDL)


def validate_relative_path(path_str: str) -> Optional[str]:
    """Validate that path is a clean relative git path without drive letters, leading slashes, or file URIs."""
    raw = path_str.strip()

    if not raw:
        return "File path string cannot be empty"

    if raw.startswith("file://") or raw.startswith("file:///"):
        return f"file:// URI rejected: '{raw}'"

    if raw.startswith("/") or raw.startswith("\\"):
        return f"Leading slash rejected (must be relative path): '{raw}'"

    if re.match(r"^[a-zA-Z]:[\\/]", raw):
        return f"Windows drive letter rejected (must be relative path): '{raw}'"

    try:
        p = Path(raw)

        if p.is_absolute():
            return f"Absolute path rejected (must be relative path): '{raw}'"
    except Exception:
        pass

    return None


def resolve_db_path(db_arg: Optional[str]) -> Optional[Path]:
    """Resolve SQLite database path from CLI argument or discover the latest run."""
    if db_arg:
        p = Path(db_arg)

        if not p.is_absolute():
            return REPO_ROOT / p

        return p

    if not TEMP_AGENTS_DIR.exists():
        return None

    runs = []

    for item in TEMP_AGENTS_DIR.iterdir():
        if item.is_dir():
            db_file = item / "agent-task.db"

            if db_file.exists():
                m = re.match(r"^(\d+)-", item.name)
                num = int(m.group(1)) if m else 0
                runs.append((num, item.stat().st_mtime, db_file))

    if not runs:
        return None

    runs.sort(key=lambda x: (x[0], x[1]), reverse=True)

    return runs[0][2]


def format_schema_markdown() -> str:
    """Format schema definition into human-readable Markdown tables."""
    lines: List[str] = [
        "# Antigravity Multi-Agent SQLite Task Manager Schema",
        "",
        "## Architecture Conventions",
        "- **Table & Column Naming:** PascalCase",
        "- **Boolean Naming:** Positive booleans only (`IsActive`, `HasCompleted`, `IsBlocked`) storing integer 0 or 1",
        "- **Concurrency:** WAL mode enabled with busy_timeout=5000ms",
        "- **Path Hygiene:** Strict relative git paths only; zero absolute paths and zero `file:///` URIs",
        "",
        "## Database Tables",
    ]

    for table_name, table_info in SCHEMA_METADATA["tables"].items():
        lines.append(f"### Table: `{table_name}`")
        lines.append(f"{table_info['description']}\n")
        lines.append("| Column | Type | Nullable | Default | Constraints / Flags | Description |")
        lines.append("| :--- | :--- | :--- | :--- | :--- | :--- |")

        for col in table_info["columns"]:
            nullable_str = "YES" if col["nullable"] else "NO"
            default_str = str(col["default"]) if col["default"] is not None else "-"
            flags = []

            if col.get("primaryKey"):
                flags.append("PRIMARY KEY")

            if col.get("foreignKey"):
                flags.append(f"FK -> {col['foreignKey']}")

            if col.get("isPositiveBoolean"):
                flags.append("Positive Boolean (0/1)")

            if col.get("isJson"):
                flags.append("JSON Payload")

            flags_str = ", ".join(flags) if flags else "-"
            lines.append(f"| `{col['name']}` | `{col['type']}` | {nullable_str} | `{default_str}` | {flags_str} | {col['description']} |")

        if table_info.get("indexes"):
            lines.append("")
            lines.append("**Indexes:**")
            for idx in table_info["indexes"]:
                lines.append(f"- `{idx}`")

        lines.append("")

    lines.append("## JSON Payload Specifications")

    for payload_name, spec in SCHEMA_METADATA["payloadSpecifications"].items():
        lines.append(f"### `{payload_name}`")
        lines.append(f"{spec['description']}\n")

        if "rules" in spec:
            lines.append("**Validation Rules:**")
            for rule in spec["rules"]:
                lines.append(f"- {rule}")
            lines.append("")

        if "example" in spec:
            lines.append(f"**Example:** `{json.dumps(spec['example'])}`\n")

        if "supportedFormats" in spec:
            lines.append("**Accepted Formats:**")
            for fmt in spec["supportedFormats"]:
                lines.append(f"- {fmt}")
            lines.append("")

        if "subtaskFields" in spec:
            lines.append("| Field | Type | Required | Description |")
            lines.append("| :--- | :--- | :--- | :--- |")
            for f_name, f_info in spec["subtaskFields"].items():
                req_str = "YES" if f_info["required"] else "NO"
                lines.append(f"| `{f_name}` | `{f_info['type']}` | {req_str} | {f_info['description']} |")
            lines.append("")

    return "\n".join(lines)


def cmd_schema(args: argparse.Namespace) -> int:
    """Display SQLite schema and JSON data integrity specifications."""
    if args.is_ddl:
        print(SCHEMA_DDL)

        return 0

    if args.is_json:
        print(json.dumps(SCHEMA_METADATA, indent=2))

        return 0

    print(format_schema_markdown())

    return 0


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
    machine_id = get_machine_id()
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
                rel_run_dir = str(existing_dir.relative_to(REPO_ROOT)).replace("\\", "/")
                rel_db_path = str(db_file.relative_to(REPO_ROOT)).replace("\\", "/")
                parent_keys = parent.keys() if hasattr(parent, "keys") else []
                task_name = parent["TaskName"] if "TaskName" in parent_keys else args.name
                parent_machine = parent["MachineId"] if "MachineId" in parent_keys and parent["MachineId"] else machine_id
                budget = parent["TotalStepsBudget"] if "TotalStepsBudget" in parent_keys else args.budget
                out = {
                    "isSuccess": True,
                    "action": "RESUME_FOUND",
                    "parentTaskId": parent["ParentTaskId"],
                    "taskSlug": parent["TaskSlug"],
                    "taskName": task_name,
                    "machineId": parent_machine,
                    "runDirectory": rel_run_dir,
                    "databasePath": rel_db_path,
                    "totalStepsBudget": budget,
                    "status": parent["Status"],
                    "similarity": round(similarity, 3),
                    "slug": parent["TaskSlug"],
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
                TaskName, TaskSlug, MachineId, RunDirectory, Status, IsActive, HasCompleted,
                TotalStepsBudget, CurrentStep, CreatedAt, UpdatedAt
            ) VALUES (?, ?, ?, ?, 'ACTIVE', 1, 0, ?, 1, ?, ?);
        """, (args.name, slug, machine_id, rel_run_dir, args.budget, created_at, created_at))
        parent_id = cur.lastrowid

    conn.close()

    out = {
        "isSuccess": True,
        "action": "INITIALIZED",
        "parentTaskId": parent_id,
        "taskSlug": slug,
        "taskName": args.name,
        "machineId": machine_id,
        "runDirectory": rel_run_dir,
        "databasePath": rel_db_path,
        "totalStepsBudget": args.budget,
        "status": "ACTIVE",
        "slug": slug,
    }
    print(json.dumps(out, indent=2))
    return 0


def cmd_add_subtasks(args: argparse.Namespace) -> int:
    """Populate subtasks for the task run with enhanced validation."""
    db_path = Path(args.db)

    if not db_path.is_absolute():
        db_path = REPO_ROOT / db_path

    if not db_path.exists():
        print(json.dumps({"error": f"Database not found at {args.db}"}), file=sys.stderr)

        return 1

    tasks_raw = args.tasks_json

    if not tasks_raw:
        if args.tasks_file:
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

    if isinstance(subtasks_data, dict):
        if "subtasks" in subtasks_data:
            if isinstance(subtasks_data["subtasks"], list):
                subtasks_data = subtasks_data["subtasks"]
            else:
                print(json.dumps({"error": "'subtasks' field in dict must be a list"}), file=sys.stderr)

                return 1
        elif "tasks" in subtasks_data:
            if isinstance(subtasks_data["tasks"], list):
                subtasks_data = subtasks_data["tasks"]
            else:
                print(json.dumps({"error": "'tasks' field in dict must be a list"}), file=sys.stderr)

                return 1
        else:
            print(json.dumps({"error": "Dict input must contain 'tasks' or 'subtasks' list"}), file=sys.stderr)

            return 1
    elif not isinstance(subtasks_data, list):
        print(json.dumps({"error": "Input JSON must be a list or an object containing 'tasks' or 'subtasks' list"}), file=sys.stderr)

        return 1

    if not subtasks_data:
        print(json.dumps({"error": "Subtask list is empty"}), file=sys.stderr)

        return 1

    conn = get_db_connection(db_path)
    init_schema(conn)

    cur = conn.cursor()
    cur.execute("SELECT ParentTaskId, TaskSlug FROM ParentTask ORDER BY ParentTaskId DESC LIMIT 1;")
    row = cur.fetchone()

    if not row:
        print(json.dumps({"error": "ParentTask record missing in database"}), file=sys.stderr)
        conn.close()

        return 1

    parent_id = args.parent_id if args.parent_id else row["ParentTaskId"]
    task_slug = row["TaskSlug"] if "TaskSlug" in row.keys() else "task"
    machine_id = get_machine_id()
    created_at = now_iso()
    inserted_ids = []

    validated_tasks: List[Dict[str, Any]] = []

    for idx, task in enumerate(subtasks_data, start=1):
        if not isinstance(task, dict):
            print(json.dumps({"error": f"Subtask item at index {idx} must be a JSON object (dict), got {type(task).__name__}"}), file=sys.stderr)
            conn.close()

            return 1

        raw_title = task.get("title")

        if not raw_title:
            print(json.dumps({"error": f"Subtask item at index {idx} must have a non-empty 'title'"}), file=sys.stderr)
            conn.close()

            return 1

        if not isinstance(raw_title, str):
            print(json.dumps({"error": f"Subtask item at index {idx} 'title' must be a string"}), file=sys.stderr)
            conn.close()

            return 1

        title = raw_title.strip()

        if not title:
            print(json.dumps({"error": f"Subtask item at index {idx} 'title' cannot be blank"}), file=sys.stderr)
            conn.close()

            return 1

        owned_files = task.get("owned_files", [])

        if owned_files is None:
            owned_files = []

        if not isinstance(owned_files, list):
            print(json.dumps({"error": f"Subtask item at index {idx} 'owned_files' must be a list of strings"}), file=sys.stderr)
            conn.close()

            return 1

        normalized_files: List[str] = []

        for file_idx, f_entry in enumerate(owned_files, start=1):
            if not isinstance(f_entry, str):
                print(json.dumps({"error": f"Subtask item at index {idx} owned_files entry {file_idx} must be a string"}), file=sys.stderr)
                conn.close()

                return 1

            path_error = validate_relative_path(f_entry)

            if path_error is not None:
                print(json.dumps({"error": f"Subtask item at index {idx} owned_files entry '{f_entry}' invalid: {path_error}"}), file=sys.stderr)
                conn.close()

                return 1

            normalized_files.append(f_entry.strip().replace("\\", "/"))

        code = task.get("code") or f"Task-{idx:02d}"
        agent = task.get("agent_role") or task.get("assigned_agent")

        validated_tasks.append({
            "code": str(code).strip(),
            "title": title,
            "agent": str(agent).strip() if agent else None,
            "owned_json": json.dumps(normalized_files),
        })

    # Check disjoint file ownership across subtasks in the incoming batch
    all_seen_files: Dict[str, str] = {}

    for vtask in validated_tasks:
        for f in json.loads(vtask["owned_json"]):
            if f in all_seen_files:
                prior_code = all_seen_files[f]
                print(json.dumps({"error": f"Disjoint file ownership collision: file '{f}' is claimed by both '{prior_code}' and '{vtask['code']}'"}), file=sys.stderr)
                conn.close()

                return 1

            all_seen_files[f] = vtask["code"]

    with conn:
        for vtask in validated_tasks:
            cur.execute("""
                INSERT INTO Subtask (
                    ParentTaskId, TaskCode, Title, MachineId, AssignedAgentRole,
                    OwnedFilesJson, Status, IsBlocked, HasCompleted,
                    CreatedAt, UpdatedAt
                ) VALUES (?, ?, ?, ?, ?, ?, 'PENDING', 0, 0, ?, ?);
            """, (parent_id, vtask["code"], vtask["title"], machine_id, vtask["agent"], vtask["owned_json"], created_at, created_at))
            inserted_ids.append(cur.lastrowid)

    conn.close()

    subtasks_list: List[Dict[str, Any]] = []
    for sid, vtask in zip(inserted_ids, validated_tasks):
        subtasks_list.append({
            "subtaskId": sid,
            "parentTaskId": parent_id,
            "taskCode": vtask["code"],
            "title": vtask["title"],
            "machineId": machine_id,
            "assignedAgentRole": vtask["agent"],
            "ownedFiles": json.loads(vtask["owned_json"]),
            "status": "PENDING",
        })

    print(json.dumps({
        "status": "SUCCESS",
        "isSuccess": True,
        "parentTaskId": parent_id,
        "taskSlug": task_slug,
        "machineId": machine_id,
        "subtasksCreated": len(inserted_ids),
        "subtaskIds": inserted_ids,
        "subtasks": subtasks_list,
    }, indent=2))

    return 0


def cmd_claim(args: argparse.Namespace) -> int:
    """Worker atomically claims next available subtask."""
    db_path = Path(args.db)
    if not db_path.is_absolute():
        db_path = REPO_ROOT / db_path

    conn = get_db_connection(db_path)
    init_schema(conn)

    machine_id = get_machine_id()

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
            print(json.dumps({"status": "NO_TASKS_AVAILABLE", "isSuccess": False}))
            return 0

        subtask_id = task["SubtaskId"]
        assigned_role = args.agent or task["AssignedAgentRole"] or "Worker"
        updated_at = now_iso()

        cur.execute("""
            UPDATE Subtask
            SET Status = 'IN_PROGRESS',
                AssignedAgentRole = ?,
                MachineId = ?,
                UpdatedAt = ?
            WHERE SubtaskId = ?;
        """, (assigned_role, machine_id, updated_at, subtask_id))

        # Log start action
        cur.execute("""
            INSERT INTO AgentActionLog (
                SubtaskId, MachineId, AgentRole, ActionType, ActionDetails, Status, CreatedAt
            ) VALUES (?, ?, ?, 'START', 'Subtask claimed and started', 'IN_PROGRESS', ?);
        """, (subtask_id, machine_id, assigned_role, updated_at))

    owned_files = json.loads(task["OwnedFilesJson"]) if task["OwnedFilesJson"] else []
    out = {
        "isSuccess": True,
        "status": "CLAIMED",
        "subtaskId": subtask_id,
        "taskCode": task["TaskCode"],
        "title": task["Title"],
        "machineId": machine_id,
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

    if not db_path.exists():
        print(json.dumps({"error": f"Database not found at {args.db}"}), file=sys.stderr)
        return 1

    # Validate relative file path if provided
    if args.file:
        path_err = validate_relative_path(args.file)
        if path_err is not None:
            print(json.dumps({"error": f"Invalid target file: {path_err}"}), file=sys.stderr)
            return 1

    # Validate and normalize action type
    normalized_action = normalize_action_type(args.action, args.file)
    if not normalized_action:
        allowed_list = sorted(list(CANONICAL_ACTION_TYPES))
        err_msg = {
            "error": f"Invalid --action '{args.action}'. Must be one of: {', '.join(allowed_list)}.",
            "allowedActionTypes": allowed_list,
        }
        print(json.dumps(err_msg), file=sys.stderr)
        return 1

    conn = get_db_connection(db_path)
    init_schema(conn)

    machine_id = get_machine_id()
    created_at = now_iso()

    reasoning_text = getattr(args, "reasoning", None) or ""
    lines_text = f"[Lines: {args.lines}] " if getattr(args, "lines", None) else ""
    if reasoning_text:
        details = f"{lines_text}{reasoning_text}"
    else:
        details = args.details or f"Executing {normalized_action} on {args.file or 'workspace'}"

    with conn:
        cur = conn.cursor()
        cur.execute("""
            INSERT INTO AgentActionLog (
                SubtaskId, MachineId, AgentRole, ActionType, TargetFile, ActionDetails, Status, CreatedAt
            ) VALUES (?, ?, ?, ?, ?, ?, 'IN_PROGRESS', ?);
        """, (args.subtask_id, machine_id, args.agent, normalized_action, args.file, details, created_at))
        action_id = cur.lastrowid

        cur.execute("UPDATE Subtask SET UpdatedAt = ? WHERE SubtaskId = ?;", (created_at, args.subtask_id))

    conn.close()
    out = {
        "isSuccess": True,
        "status": "LOGGED",
        "actionLogId": action_id,
        "subtaskId": args.subtask_id,
        "agentRole": args.agent,
        "actionType": normalized_action,
        "targetFile": args.file,
        "machineId": machine_id,
    }
    print(json.dumps(out, indent=2))
    return 0


def cmd_complete(args: argparse.Namespace) -> int:
    """Mark subtask completed with evidence."""
    db_path = Path(args.db)
    if not db_path.is_absolute():
        db_path = REPO_ROOT / db_path

    conn = get_db_connection(db_path)
    init_schema(conn)

    machine_id = get_machine_id()
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
                SubtaskId, MachineId, AgentRole, ActionType, ActionDetails, Status, CreatedAt
            ) VALUES (?, ?, ?, 'COMPLETE', 'Subtask verified and completed', 'DONE', ?);
        """, (args.subtask_id, machine_id, args.agent or "Worker", updated_at))

    conn.close()
    print(json.dumps({"isSuccess": True, "status": "COMPLETED", "subtaskId": args.subtask_id, "machineId": machine_id}))
    return 0


def cmd_fail(args: argparse.Namespace) -> int:
    """Mark subtask failed or blocked."""
    db_path = Path(args.db)
    if not db_path.is_absolute():
        db_path = REPO_ROOT / db_path

    conn = get_db_connection(db_path)
    init_schema(conn)

    machine_id = get_machine_id()
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
                SubtaskId, MachineId, AgentRole, ActionType, ActionDetails, Status, CreatedAt
            ) VALUES (?, ?, ?, 'FAIL', ?, 'FAILED', ?);
        """, (args.subtask_id, machine_id, args.agent or "Worker", args.reason or "Task failed", updated_at))

    conn.close()
    print(json.dumps({"isSuccess": False, "status": "FAILED", "subtaskId": args.subtask_id, "machineId": machine_id}))
    return 0


def diagnose_db(conn: sqlite3.Connection) -> Dict[str, Any]:
    """Diagnose uncompleted or crashed agent actions."""
    cur = conn.cursor()

    cur.execute("""
        SELECT s.*,
               a.ActionType AS LastActionType,
               a.TargetFile AS LastTargetFile,
               a.ActionDetails AS LastActionDetails,
               a.CreatedAt AS LastActionTimestamp,
               a.MachineId AS LastActionMachineId
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
        row_keys = row.keys() if hasattr(row, "keys") else []
        crashes.append({
            "subtaskId": row["SubtaskId"],
            "taskCode": row["TaskCode"],
            "title": row["Title"],
            "machineId": row["MachineId"] if "MachineId" in row_keys else "local-machine",
            "assignedAgentRole": row["AssignedAgentRole"],
            "lastActionType": row["LastActionType"],
            "lastTargetFile": row["LastTargetFile"],
            "lastActionDetails": row["LastActionDetails"],
            "lastTimestamp": row["LastActionTimestamp"],
            "lastActionMachineId": row["LastActionMachineId"] if "LastActionMachineId" in row_keys else "local-machine",
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

    cur.execute("SELECT SubtaskId, TaskCode, Title, MachineId, AssignedAgentRole, Status, Evidence FROM Subtask ORDER BY SubtaskId ASC;")
    subtasks = [dict(row) for row in cur.fetchall()]

    counts: Dict[str, int] = {}
    for st in subtasks:
        s = st["Status"]
        counts[s] = counts.get(s, 0) + 1

    total = len(subtasks)
    done = counts.get("DONE", 0)
    percent = round((done / total * 100), 1) if total > 0 else 0.0

    conn.close()

    machine_id = get_machine_id()
    parent_keys = parent.keys() if parent and hasattr(parent, "keys") else []
    parent_machine = parent["MachineId"] if parent and "MachineId" in parent_keys and parent["MachineId"] else machine_id

    out = {
        "taskSlug": parent["TaskSlug"] if parent else "unknown",
        "taskName": parent["TaskName"] if parent else "unknown",
        "machineId": parent_machine,
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


def cmd_heatmap_discuss(args: argparse.Namespace) -> int:
    """Real-time fast heatmap discussion engine (<5ms) for multi-agent file contention."""
    start_time = time.perf_counter()

    db_path = resolve_db_path(getattr(args, "db", None))
    if db_path is None or not db_path.exists():
        err_msg = f"Database not found at {getattr(args, 'db', None) or '.ai-memory/temp-agents/'}"
        if getattr(args, "is_json", False):
            print(json.dumps({"isSuccess": False, "error": err_msg}))
        else:
            print(f"Error: {err_msg}", file=sys.stderr)
        return 1

    conn = get_db_connection(db_path)
    init_schema(conn)

    machine_id = get_machine_id()
    cur = conn.cursor()

    # 1. Parent Task context
    cur.execute("SELECT TaskSlug, TaskName, MachineId FROM ParentTask ORDER BY ParentTaskId DESC LIMIT 1;")
    parent_row = cur.fetchone()
    task_slug = parent_row["TaskSlug"] if parent_row else "unknown"
    task_name = parent_row["TaskName"] if parent_row else "unknown"
    parent_machine = parent_row["MachineId"] if parent_row else machine_id

    # 2. File claims from Subtask
    cur.execute("SELECT SubtaskId, TaskCode, AssignedAgentRole, OwnedFilesJson, Status FROM Subtask;")
    subtask_rows = cur.fetchall()

    claims_by_file: Dict[str, List[Dict[str, str]]] = {}
    total_active_claims = 0

    for st in subtask_rows:
        code = st["TaskCode"]
        role = st["AssignedAgentRole"] or "Worker"
        raw_owned = st["OwnedFilesJson"]
        if raw_owned:
            try:
                owned_files = json.loads(raw_owned)
                if isinstance(owned_files, list):
                    for f in owned_files:
                        if isinstance(f, str) and f.strip():
                            clean_path = f.strip().replace("\\", "/")
                            if clean_path not in claims_by_file:
                                claims_by_file[clean_path] = []
                            claims_by_file[clean_path].append({
                                "agentRole": role,
                                "subtaskCode": code,
                            })
                            total_active_claims += 1
            except Exception:
                pass

    # 3. Actions from AgentActionLog
    window_sec = getattr(args, "window", 1800) or 1800
    cutoff = (datetime.datetime.now(datetime.timezone.utc) - datetime.timedelta(seconds=window_sec)).isoformat()

    cur.execute("""
        SELECT ActionLogId, SubtaskId, MachineId, AgentRole, ActionType, TargetFile, ActionDetails, Status, CreatedAt
        FROM AgentActionLog
        WHERE TargetFile IS NOT NULL AND TargetFile != ''
        ORDER BY ActionLogId DESC
        LIMIT 2000;
    """)
    all_action_rows = cur.fetchall()
    conn.close()

    # Filter actions in time window, fallback to all recent if window has none
    window_actions = [row for row in all_action_rows if row["CreatedAt"] >= cutoff]
    actions_to_analyze = window_actions if window_actions else all_action_rows

    actions_by_file: Dict[str, List[sqlite3.Row]] = {}
    for act in actions_to_analyze:
        t_file = act["TargetFile"].strip().replace("\\", "/")
        if t_file:
            if t_file not in actions_by_file:
                actions_by_file[t_file] = []
            actions_by_file[t_file].append(act)

    all_candidate_files = sorted(list(set(claims_by_file.keys()) | set(actions_by_file.keys())))

    file_entries: List[Dict[str, Any]] = []
    active_agents: set[str] = set()

    for f_path in all_candidate_files:
        f_actions = actions_by_file.get(f_path, [])
        f_claims = claims_by_file.get(f_path, [])

        agents_set = set()
        for act in f_actions:
            if act["AgentRole"]:
                agents_set.add(act["AgentRole"])
        for clm in f_claims:
            if clm["agentRole"]:
                agents_set.add(clm["agentRole"])

        active_agents.update(agents_set)

        mutations_count = 0
        inspections_count = 0
        in_flight_count = 0

        for act in f_actions:
            a_type = act["ActionType"].lower()
            if a_type in ("modify", "create", "remove", "refactor"):
                mutations_count += 1
            elif a_type in ("read", "analyze"):
                inspections_count += 1
            if act["Status"] == "IN_PROGRESS":
                in_flight_count += 1

        claims_count = len(f_claims)
        touch_count = len(f_actions)
        unique_agents_count = len(agents_set)

        # Contention scoring formula
        agent_penalty = 5.0 if unique_agents_count > 1 else 0.0
        score = (mutations_count * 2.0) + (inspections_count * 0.5) + agent_penalty + (claims_count * 3.0) + (in_flight_count * 2.0)
        contention_score = round(score, 1)

        # Severity / Contention Level
        if contention_score >= 10.0 or (unique_agents_count >= 2 and (mutations_count > 0 or claims_count >= 2)):
            severity = "COLLISION_RISK"
            contention_level = "CRITICAL"
        elif contention_score >= 4.0:
            severity = "WARM"
            contention_level = "MEDIUM"
        elif contention_score >= 0.5:
            severity = "COOL"
            contention_level = "LOW"
        else:
            severity = "IDLE"
            contention_level = "IDLE"

        recent_actions = []
        for act in f_actions[:5]:
            recent_actions.append({
                "agentRole": act["AgentRole"],
                "actionType": act["ActionType"],
                "details": act["ActionDetails"],
                "status": act["Status"],
                "timestamp": act["CreatedAt"],
            })

        last_action_type = recent_actions[0]["actionType"] if recent_actions else None
        last_touched_at = recent_actions[0]["timestamp"] if recent_actions else None

        file_entries.append({
            "filePath": f_path,
            "contentionScore": contention_score,
            "contentionLevel": contention_level,
            "severity": severity,
            "claimsCount": claims_count,
            "inFlightTouches": in_flight_count,
            "touchCount": touch_count,
            "agents": sorted(list(agents_set)),
            "activeClaimants": f_claims,
            "recentActions": recent_actions,
            "lastActionType": last_action_type,
            "lastTouchedAt": last_touched_at,
        })

    # Sort files by contentionScore descending, then touchCount descending
    file_entries.sort(key=lambda x: (x["contentionScore"], x["touchCount"]), reverse=True)

    high_contention_files = [f for f in file_entries if f["severity"] == "COLLISION_RISK"]
    high_contention_count = len(high_contention_files)
    collision_warning = high_contention_count > 0

    if collision_warning:
        c_names = [f["filePath"] for f in high_contention_files[:3]]
        discussion_summary = f"Collision risk detected on {', '.join(c_names)}: overlapping claims or concurrent worker touches."
    else:
        discussion_summary = "No collision risks detected. All active claims are disjoint and safe."

    elapsed_ms = round((time.perf_counter() - start_time) * 1000, 2)
    now_ts = now_iso()

    if getattr(args, "is_json", False):
        payload = {
            "isSuccess": True,
            "timestamp": now_ts,
            "evaluatedAt": now_ts,
            "taskSlug": task_slug,
            "taskName": task_name,
            "machineId": parent_machine,
            "executionTimeMs": elapsed_ms,
            "activeAgentsCount": len(active_agents),
            "totalHeatedFiles": len(file_entries),
            "totalActiveClaims": total_active_claims,
            "highContentionFilesCount": high_contention_count,
            "collisionWarning": collision_warning,
            "hotspotFiles": file_entries,
            "files": file_entries,
            "discussionSummary": discussion_summary,
        }
        print(json.dumps(payload, indent=2))
        return 0

    # Terminal human formatted output
    print(f"[{elapsed_ms:.1f}ms] Heatmap Discussion Summary: {task_slug}")
    print(f"Evaluated: {len(file_entries)} heated files | Active Claims: {total_active_claims} | Critical Contention: {high_contention_count}\n")

    if not file_entries:
        print("No heated files or active file claims found in current window.")
        return 0

    print(f"{'FILE PATH':<48} {'CLAIMS':>6} {'TOUCHES':>8} {'SCORE':>6} {'LEVEL':<10} {'ACTIVE WORKERS'}")
    print("-" * 95)
    for fe in file_entries[:15]:
        workers_str = ", ".join(fe["agents"]) if fe["agents"] else "-"
        print(f"{fe['filePath'][:46]:<48} {fe['claimsCount']:>6} {fe['touchCount']:>8} {fe['contentionScore']:>6} {fe['contentionLevel']:<10} {workers_str}")

    print("-" * 95)
    print("DISCUSSION ADVISORY & CONTENTION MITIGATION:")
    if collision_warning:
        for cf in high_contention_files[:3]:
            claimant_roles = [c["agentRole"] for c in cf["activeClaimants"]]
            print(f"  * Critical contention on '{cf['filePath']}' (Score: {cf['contentionScore']}).")
            if claimant_roles:
                print(f"    Active claimants: {', '.join(claimant_roles)}. Subtasks must hold strictly disjoint files.")
    else:
        print("  * No active file collisions detected. Safe for concurrent worker operations.")

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
    p_log.add_argument("--action", required=True, help="Action type (e.g. modify, replace_file_content)")
    p_log.add_argument("--file", help="Target relative file path being touched")
    p_log.add_argument("--details", help="Human-readable details of action")
    p_log.add_argument("--lines", help="Line range e.g. 10,20")
    p_log.add_argument("--reasoning", help="Technical rationale for the action")

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

    # schema
    p_schema = subparsers.add_parser("schema", help="Display SQLite database schema and data integrity specifications")
    p_schema.add_argument("--json", action="store_true", dest="is_json", help="Output machine-readable JSON schema definition")
    p_schema.add_argument("--ddl", action="store_true", dest="is_ddl", help="Output raw SQL CREATE TABLE and INDEX statements")
    p_schema.add_argument("--db", help="Optional path to SQLite database")

    # heatmap
    p_heat = subparsers.add_parser("heatmap", help="Real-time file churn and contention heatmap")
    p_heat.add_argument("subaction", nargs="?", default=None, help="Optional subaction: 'discuss'")
    p_heat.add_argument("--discuss", action="store_true", help="Enable discussion summary mode")
    p_heat.add_argument("--db", help="Path to SQLite database")
    p_heat.add_argument("--slug", help="Filter by task slug")
    p_heat.add_argument("--window", type=int, default=1800, help="Time window in seconds (default: 1800)")
    p_heat.add_argument("--fast", action="store_true", default=True, help="Fast query mode without disk stat")
    p_heat.add_argument("--json", action="store_true", dest="is_json", help="Emit structured JSON output")

    # heatmap-discuss alias
    p_hd = subparsers.add_parser("heatmap-discuss", help="Alias for heatmap discuss")
    p_hd.add_argument("--db", help="Path to SQLite database")
    p_hd.add_argument("--slug", help="Filter by task slug")
    p_hd.add_argument("--window", type=int, default=1800, help="Time window in seconds (default: 1800)")
    p_hd.add_argument("--fast", action="store_true", default=True, help="Fast query mode without disk stat")
    p_hd.add_argument("--json", action="store_true", dest="is_json", help="Emit structured JSON output")

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
        "schema": cmd_schema,
        "heatmap": cmd_heatmap_discuss,
        "heatmap-discuss": cmd_heatmap_discuss,
    }

    return dispatch[args.command](args)


if __name__ == "__main__":
    sys.exit(main())
