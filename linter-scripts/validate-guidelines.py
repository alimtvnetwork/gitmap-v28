#!/usr/bin/env python3
"""
Cross-Language Coding Guidelines Validator
==========================================

Version: 1.5.0  (2026-04-19) — Added P2/P3/P5/P7 boolean-principle checks.

Validates Go, PHP, TypeScript, and Rust source files against the coding
guidelines defined in spec/02-coding-guidelines/03-coding-guidelines-spec/.

This script enforces rules that ESLint cannot cover (Go, PHP, Rust)
and provides a unified validation report across all languages.

Usage:
    python3 scripts/validate-guidelines.py [--path <dir>] [--fix] [--json]

Rules Enforced:
    CODE-RED-001  No nested if statements
    CODE-RED-002  Boolean naming (is/has/can/should/was prefix)         [P1]
    CODE-RED-003  No magic strings in comparisons
    CODE-RED-004  Max 15 lines per function
    CODE-RED-005  No fmt.Errorf() in Go (use apperror)
    CODE-RED-006  No (T, error) returns in Go services (use Result[T])
    CODE-RED-007  No string-based Go enums (use byte + iota)
    CODE-RED-008  No raw string error codes — use apperrtype enum
    CODE-RED-009  Generic file-not-found without exact path or reason
    CODE-RED-011  No magic numbers in logic
    CODE-RED-012  Immutable by default (const over let/var, no reassignment)
    CODE-RED-013  File length ≤ 300 lines (hard max 400)
    CODE-RED-014  Max 3 parameters per function
    CODE-RED-015  No any/interface{}/unknown in business logic
    CODE-RED-016  No silent error swallowing (empty catch, ignored err)
    CODE-RED-017  PHP: catch Throwable, not Exception
    CODE-RED-018  Sequential independent async (use Promise.all / errgroup)
    CODE-RED-019  SQL string concatenation (injection risk)
    CODE-RED-020  Go: missing stack trace (raw errors.New instead of apperror)
    CODE-RED-021  Mixed && and || in single expression                  [P4]
    CODE-RED-022  Negative words in boolean identifiers (isNot*, hasNo*) [P2]
    CODE-RED-023  Raw `!` on function/method calls (use semantic inverse) [P3]
    CODE-RED-024  Bare true/false as positional argument                [P5]
    CODE-RED-025  Assignment inside if/while condition                  [P7]
    STYLE-001     Blank line before return (R4)
    STYLE-002     No else after return (R7)
    STYLE-003     Blank line after closing brace (R5)
    STYLE-004     Blank line before if/else if block
"""

import re
import os
import sys
import json
import glob
import argparse
from dataclasses import dataclass, field, asdict
from typing import List, Optional
from collections import defaultdict


@dataclass
class Violation:
    file: str
    line: int
    rule: str
    severity: str  # "CODE-RED" or "STYLE"
    message: str
    code_snippet: str = ""


@dataclass
class ValidationReport:
    total_files: int = 0
    total_violations: int = 0
    code_red_count: int = 0
    style_count: int = 0
    violations: List[Violation] = field(default_factory=list)
    by_rule: dict = field(default_factory=lambda: defaultdict(int))
    by_file: dict = field(default_factory=lambda: defaultdict(int))


# ═══════════════════════════════════════════════════════════════════════
# Language Detection
# ═══════════════════════════════════════════════════════════════════════

def detect_language(filepath: str) -> Optional[str]:
    ext = os.path.splitext(filepath)[1]
    mapping = {
        ".go": "go",
        ".ts": "typescript", ".tsx": "typescript",
        ".js": "javascript", ".jsx": "javascript",
        ".php": "php",
        ".rs": "rust",
    }
    return mapping.get(ext)


# ═══════════════════════════════════════════════════════════════════════
# Rule Checkers
# ═══════════════════════════════════════════════════════════════════════

def check_nested_if(lines: List[str], filepath: str) -> List[Violation]:
    """CODE-RED-001: No nested if statements."""
    violations = []
    indent_stack = []

    for i, line in enumerate(lines):
        stripped = line.strip()

        # Track if-block nesting by indentation
        if stripped.startswith("if ") or stripped.startswith("if("):
            current_indent = len(line) - len(line.lstrip())

            # Check if we're inside another if block
            for parent_indent, parent_line in indent_stack:
                if current_indent > parent_indent:
                    violations.append(Violation(
                        file=filepath,
                        line=i + 1,
                        rule="CODE-RED-001",
                        severity="CODE-RED",
                        message="Nested `if` is forbidden. Flatten with early returns or combined conditions.",
                        code_snippet=stripped[:120],
                    ))
                    break

            indent_stack.append((current_indent, i))

        # Reset tracking on closing braces at lower indent
        if stripped == "}" and indent_stack:
            current_indent = len(line) - len(line.lstrip())
            indent_stack = [(ind, ln) for ind, ln in indent_stack if ind < current_indent]

    return violations


def check_boolean_naming(lines: List[str], filepath: str, lang: str) -> List[Violation]:
    """CODE-RED-002: Boolean variables must use is/has/can/should/was prefix."""
    violations = []
    prefixes = ("is", "has", "can", "should", "was", "will")
    exempt = {"ok", "done", "found", "exists", "err", "error", "true", "false"}

    if lang == "go":
        pattern = re.compile(r"(\w+)\s*:=\s*(true|false)\b")
    elif lang == "php":
        pattern = re.compile(r"\$(\w+)\s*=\s*(true|false)\b", re.IGNORECASE)
    elif lang in ("typescript", "javascript"):
        pattern = re.compile(r"(?:const|let|var)\s+(\w+)\s*(?::\s*boolean)?\s*=\s*(true|false)\b")
    else:
        return violations

    for i, line in enumerate(lines):
        for m in pattern.finditer(line):
            name = m.group(1)

            if name.lower() in exempt:
                continue
            if name.startswith("_"):
                continue
            if any(name.startswith(p) and len(name) > len(p) and name[len(p)].isupper() for p in prefixes):
                continue

            violations.append(Violation(
                file=filepath,
                line=i + 1,
                rule="CODE-RED-002",
                severity="CODE-RED",
                message=f'Boolean variable "${name}" must start with is/has/can/should/was/will.',
                code_snippet=line.strip()[:120],
            ))

    return violations


def check_magic_strings(lines: List[str], filepath: str, lang: str) -> List[Violation]:
    """CODE-RED-003: No raw string literals in comparisons."""
    violations = []
    exempt_values = {"", "GET", "POST", "PUT", "DELETE", "PATCH", "string", "number",
                     "boolean", "undefined", "object", "function", "utf-8", "utf8",
                     # UI/layout prop values — not business logic
                     "horizontal", "vertical", "left", "right", "top", "bottom",
                     "sm", "md", "lg", "xl", "2xl", "none", "auto",
                     "popper", "dot", "line", "dashed", "collapsed", "expanded",
                     "floating", "inset", "sidebar", "header", "footer",
                     "default", "destructive", "outline", "secondary", "ghost", "link",
                     }

    if lang == "go":
        pattern = re.compile(r'(?:==|!=)\s*"([^"]+)"')
    elif lang == "php":
        pattern = re.compile(r"(?:===|!==|==|!=)\s*['\"]([^'\"]+)['\"]")
    elif lang in ("typescript", "javascript"):
        pattern = re.compile(r"(?:===|!==)\s*['\"]([^'\"]+)['\"]")
    else:
        return violations

    for i, line in enumerate(lines):
        stripped = line.strip()

        # Skip comments
        if stripped.startswith("//") or stripped.startswith("#") or stripped.startswith("*"):
            continue

        for m in pattern.finditer(line):
            value = m.group(1)

            if value in exempt_values or len(value) <= 1:
                continue
            if value.startswith("/") or value.startswith("http") or value.startswith("."):
                continue

            violations.append(Violation(
                file=filepath,
                line=i + 1,
                rule="CODE-RED-003",
                severity="CODE-RED",
                message=f'Magic string "{value}" in comparison. Use an enum constant.',
                code_snippet=stripped[:120],
            ))

    return violations


def check_function_length(lines: List[str], filepath: str, lang: str, max_lines: int = 15) -> List[Violation]:
    """CODE-RED-004: Max 15 lines per function body."""
    violations = []

    if lang == "go":
        func_pattern = re.compile(r"^func\s+(?:\(.*?\)\s+)?(\w+)")
    elif lang == "php":
        func_pattern = re.compile(r"(?:public|private|protected|static)?\s*function\s+(\w+)")
    elif lang in ("typescript", "javascript"):
        func_pattern = re.compile(r"(?:function\s+(\w+)|(?:const|let|var)\s+(\w+)\s*=\s*(?:async\s+)?(?:\([^)]*\)|[^=])\s*=>)")
    else:
        return violations

    i = 0

    while i < len(lines):
        m = func_pattern.search(lines[i])

        if m:
            func_name = m.group(1) or (m.group(2) if m.lastindex >= 2 else None) or "anonymous"
            # Find the opening brace
            brace_line = i

            while brace_line < len(lines) and "{" not in lines[brace_line]:
                brace_line += 1

            if brace_line >= len(lines):
                i += 1
                continue

            # Count lines in function body
            depth = 0
            body_start = brace_line
            body_lines = 0

            for j in range(brace_line, len(lines)):
                depth += lines[j].count("{") - lines[j].count("}")

                if j > brace_line:
                    stripped = lines[j].strip()

                    if stripped and not stripped.startswith("//") and not stripped.startswith("*") and stripped != "}":
                        body_lines += 1

                if depth <= 0 and j > brace_line:
                    break

            if body_lines > max_lines:
                violations.append(Violation(
                    file=filepath,
                    line=i + 1,
                    rule="CODE-RED-004",
                    severity="CODE-RED",
                    message=f'Function "{func_name}" has {body_lines} lines (max {max_lines}). Extract into smaller helpers.',
                    code_snippet=lines[i].strip()[:120],
                ))

        i += 1

    return violations


def check_go_specific(lines: List[str], filepath: str) -> List[Violation]:
    """Go-specific CODE RED rules."""
    violations = []

    for i, line in enumerate(lines):
        stripped = line.strip()

        # CODE-RED-005: No fmt.Errorf()
        if "fmt.Errorf(" in stripped and not stripped.startswith("//"):
            violations.append(Violation(
                file=filepath, line=i + 1,
                rule="CODE-RED-005", severity="CODE-RED",
                message="Use apperror.New()/apperror.Wrap() instead of fmt.Errorf().",
                code_snippet=stripped[:120],
            ))

        # CODE-RED-007: No string-based enums
        if re.match(r"type\s+\w+\s+string\s*$", stripped):
            violations.append(Violation(
                file=filepath, line=i + 1,
                rule="CODE-RED-007", severity="CODE-RED",
                message="String-based enums forbidden. Use `type Variant byte` with iota.",
                code_snippet=stripped[:120],
            ))

        # CODE-RED-008: No raw string error codes — use apperrtype enum
        if re.search(r'apperror\.(New|Wrap|Fail\w*)\(\s*"E\d+', stripped) and not stripped.startswith("//"):
            violations.append(Violation(
                file=filepath, line=i + 1,
                rule="CODE-RED-008", severity="CODE-RED",
                message='Raw string error code (e.g. "E2010") found. Use apperrtype enum: apperror.NewType(apperrtype.SiteNotFound).',
                code_snippet=stripped[:120],
            ))
    # CODE-RED-006: No (T, error) returns in service functions
    for i, line in enumerate(lines):
        if re.search(r"func\s.*\)\s*\(\*?\w+,\s*error\)", line):
            if "test" not in filepath.lower() and "_test.go" not in filepath:
                violations.append(Violation(
                    file=filepath, line=i + 1,
                    rule="CODE-RED-006", severity="CODE-RED",
                    message="Service functions must return Result[T], not (T, error).",
                    code_snippet=line.strip()[:120],
                ))

    return violations


def check_magic_numbers(lines: List[str], filepath: str, lang: str) -> List[Violation]:
    """CODE-RED-011: No magic numbers in logic (comparisons, assignments, arithmetic)."""
    violations = []
    exempt_numbers = {"0", "1", "-1", "0.0", "1.0", "100", "2", "404"}
    # Match numeric literals in comparisons and arithmetic operators
    pattern = re.compile(
        r'(?:==|!=|===|!==|>=|<=|>|<|\*|/|%|\+\s*(?!\+)|-\s*(?!-))'
        r'\s*(-?\d+\.?\d*)\b'
    )
    # Tailwind/CSS patterns: bg-black/80, top-1/2, w-3/4, opacity-50, z-50, p-2.5
    tailwind_number_pattern = re.compile(r'[a-z]-\d|/\d|opacity-|z-\d|gap-|p-|m-|w-|h-|text-\d|rounded-')
    # Line is inside a string (className, template literal, JSX attribute)
    string_context_pattern = re.compile(r'className=|class=|["\'`].*[-/]\d')

    # Track whether we're inside a const object/map/Record literal
    in_const_object = False
    const_brace_depth = 0
    # Pattern to detect start of const object/map declarations
    const_object_start = re.compile(
        r'^\s*(?:export\s+)?(?:const|var|let)\s+\w+\s*'
        r'(?::\s*(?:Record|Map|Readonly)?[<\[{(].*?[>\]})]\s*)?=\s*\{',
    )
    # Go map literal: var/const name = map[...]Type{
    go_map_start = re.compile(r'^\s*(?:var|const)?\s*\w+\s*(?::?=)\s*map\[')

    for i, line in enumerate(lines):
        stripped = line.strip()

        # --- Track const object/map scope ---
        if not in_const_object:
            if const_object_start.search(stripped) or go_map_start.search(stripped):
                in_const_object = True
                const_brace_depth = stripped.count("{") - stripped.count("}")
                if const_brace_depth <= 0:
                    in_const_object = False
                continue
        else:
            const_brace_depth += stripped.count("{") - stripped.count("}")
            if const_brace_depth <= 0:
                in_const_object = False
            continue  # skip all lines inside const object literals

        # Skip comments, imports, const declarations
        if stripped.startswith("//") or stripped.startswith("#") or stripped.startswith("*"):
            continue
        if stripped.startswith("import") or stripped.startswith("require"):
            continue
        if stripped.startswith("const ") or stripped.startswith("const("):
            continue
        if "= iota" in stripped:
            continue

        # Skip lines that are primarily Tailwind/CSS class strings
        if tailwind_number_pattern.search(stripped) or string_context_pattern.search(stripped):
            continue

        # Skip numbers that only appear inside quoted strings (HSL colors, config values)
        no_strings = re.sub(r'(["\']).*?\1', '', stripped)
        if not pattern.search(no_strings):
            continue

        # Skip JSX text content like <h1>404</h1>
        if re.search(r'>\s*\d+\s*</', stripped):
            continue

        # Skip object literal value lines (key: value patterns in config maps)
        if re.match(r'^\s*[\w"\']+\s*:\s*["\'\d]', stripped):
            continue

        # Skip array/slice index access like [0], [1]
        if re.search(r'\[\s*\d+\s*\]', stripped):
            cleaned = re.sub(r'\[\s*\d+\s*\]', '', stripped)
            if not pattern.search(cleaned):
                continue

        for m in pattern.finditer(line):
            value = m.group(1)

            if value in exempt_numbers:
                continue
            if "line" in stripped.lower() or "column" in stripped.lower():
                continue

            violations.append(Violation(
                file=filepath,
                line=i + 1,
                rule="CODE-RED-011",
                severity="CODE-RED",
                message=f'Magic number {value} in logic. Use a named constant.',
                code_snippet=stripped[:120],
            ))

    return violations


def check_variable_mutation(lines: List[str], filepath: str, lang: str) -> List[Violation]:
    """CODE-RED-012: Prefer const/immutable. Flag let/var declarations and reassignments."""
    violations = []

    if lang in ("typescript", "javascript"):
        # Flag let/var declarations (suggest const)
        let_var_pattern = re.compile(r'^\s*(?:let|var)\s+(\w+)\s*(?::\s*\w+)?\s*=')
        reassign_pattern = re.compile(r'^\s*(\w+)\s*(?:\+=|-=|\*=|/=|%=|&&=|\|\|=|\?\?=|=(?!=))')

        declared_lets = {}  # name -> line number

        for i, line in enumerate(lines):
            stripped = line.strip()

            if stripped.startswith("//") or stripped.startswith("/*") or stripped.startswith("*"):
                continue

            # Track let/var declarations
            m = let_var_pattern.match(line)
            if m:
                name = m.group(1)
                # Skip common exemptions: loop vars, React state
                if name in ("i", "j", "k", "idx", "index"):
                    continue
                if "useState" in line or "useRef" in line:
                    continue
                # Skip for-loop declarations
                if stripped.startswith("for ") or stripped.startswith("for("):
                    continue

                declared_lets[name] = i + 1
                violations.append(Violation(
                    file=filepath,
                    line=i + 1,
                    rule="CODE-RED-012",
                    severity="CODE-RED",
                    message=f'Use `const` instead of `let`/`var` for "{name}". Assign once.',
                    code_snippet=stripped[:120],
                ))

    elif lang == "go":
        # In Go, flag variables that are reassigned with = (not :=)
        # after initial declaration — limited heuristic
        reassign_pattern = re.compile(r'^\s+(\w+)\s*=\s+')
        decl_pattern = re.compile(r'^\s+(\w+)\s*:=\s+')
        declared = {}

        for i, line in enumerate(lines):
            stripped = line.strip()

            if stripped.startswith("//"):
                continue

            m = decl_pattern.match(line)
            if m:
                declared[m.group(1)] = i + 1

            m2 = reassign_pattern.match(line)
            if m2:
                name = m2.group(1)
                if name in declared and (i + 1) - declared[name] > 1:
                    # Skip common exemptions
                    if name in ("err", "ok", "ctx", "cancel", "mu", "wg"):
                        continue
                    if name.startswith("_"):
                        continue

                    violations.append(Violation(
                        file=filepath,
                        line=i + 1,
                        rule="CODE-RED-012",
                        severity="CODE-RED",
                        message=f'Variable "{name}" reassigned (first declared L{declared[name]}). Avoid mutation.',
                        code_snippet=stripped[:120],
                    ))

    return violations


# ═══════════════════════════════════════════════════════════════════════
# New Rule Checkers (Consolidated Review Guide)
# ═══════════════════════════════════════════════════════════════════════


def check_file_length(lines: List[str], filepath: str) -> List[Violation]:
    """CODE-RED-013: File length ≤ 300 lines (hard max 400)."""
    total = len(lines)

    if total <= 300:
        return []

    severity = "CODE-RED" if total > 400 else "STYLE"
    rule = "CODE-RED-013" if total > 400 else "STYLE-005"
    msg = f"File has {total} lines (limit 300, hard max 400). Split by concern."

    return [Violation(file=filepath, line=1, rule=rule, severity=severity, message=msg)]


def check_parameter_count(lines: List[str], filepath: str, lang: str) -> List[Violation]:
    """CODE-RED-014: Max 3 parameters per function."""
    violations = []

    if lang == "go":
        pattern = re.compile(r"^func\s+(?:\(.*?\)\s+)?(\w+)\(([^)]*)\)")
    elif lang == "php":
        pattern = re.compile(r"function\s+(\w+)\(([^)]*)\)")
    elif lang in ("typescript", "javascript"):
        pattern = re.compile(r"(?:function\s+(\w+)|(?:const|let)\s+(\w+)\s*=\s*(?:async\s*)?\()\s*([^)]*)\)")
        # React components use a single destructured props object: function Foo({ a, b, c, d }: Props)
        # These are NOT multi-param functions — they have 1 param (the props object).
        react_destructured = re.compile(r"function\s+\w+\s*\(\s*\{")
    else:
        return violations

    for i, line in enumerate(lines):
        m = pattern.search(line)

        if not m:
            continue

        # Skip React component destructured props (single object param)
        if lang in ("typescript", "javascript") and react_destructured.search(line):
            continue

        groups = m.groups()
        func_name = next((g for g in groups[:-1] if g), "anonymous")
        params_str = groups[-1].strip()

        if not params_str:
            continue

        param_count = len([p for p in params_str.split(",") if p.strip()])

        if param_count > 3:
            violations.append(Violation(
                file=filepath, line=i + 1,
                rule="CODE-RED-014", severity="CODE-RED",
                message=f'Function "{func_name}" has {param_count} params (max 3). Use options object.',
                code_snippet=line.strip()[:120],
            ))

    return violations


def check_no_any_type(lines: List[str], filepath: str, lang: str) -> List[Violation]:
    """CODE-RED-015: No any/interface{}/unknown in business logic."""
    violations = []

    for i, line in enumerate(lines):
        stripped = line.strip()

        if stripped.startswith("//") or stripped.startswith("*"):
            continue

        if lang in ("typescript", "javascript"):
            if re.search(r":\s*any\b", stripped) or re.search(r"as\s+any\b", stripped):
                violations.append(Violation(
                    file=filepath, line=i + 1,
                    rule="CODE-RED-015", severity="CODE-RED",
                    message="Type `any` is forbidden. Use a concrete type or generic.",
                    code_snippet=stripped[:120],
                ))

        if lang == "go":
            if "interface{}" in stripped and "// allow:" not in stripped:
                violations.append(Violation(
                    file=filepath, line=i + 1,
                    rule="CODE-RED-015", severity="CODE-RED",
                    message="`interface{}` forbidden in business logic. Use concrete type or generic.",
                    code_snippet=stripped[:120],
                ))

    return violations


def check_error_swallowing(lines: List[str], filepath: str, lang: str) -> List[Violation]:
    """CODE-RED-016: No silent error swallowing."""
    violations = []

    for i, line in enumerate(lines):
        stripped = line.strip()

        # TS/JS: empty catch block
        if lang in ("typescript", "javascript"):
            if stripped == "} catch {" or stripped == "} catch (e) {" or stripped == "} catch (error) {":
                if i + 1 < len(lines) and lines[i + 1].strip() == "}":
                    violations.append(Violation(
                        file=filepath, line=i + 1,
                        rule="CODE-RED-016", severity="CODE-RED",
                        message="Empty catch block — never swallow errors. Log or rethrow.",
                        code_snippet=stripped[:120],
                    ))

        # Go: err assigned but never checked
        if lang == "go":
            if re.match(r'\s*\w+,\s*_\s*:?=', line) or re.match(r'\s*_\s*=\s*\w+\.', line):
                if "err" not in stripped and "test" not in filepath.lower():
                    pass  # blank assignment is sometimes OK

            if re.match(r'\s*(\w+)\(.*\)\s*$', stripped) and "err" not in stripped:
                pass  # function call without error capture — heuristic only

    # Go: assigned err but no check on next lines
    if lang == "go":
        for i, line in enumerate(lines):
            if re.search(r'\berr\s*:?=\s*\w+', line):
                has_check = False

                for j in range(i + 1, min(i + 4, len(lines))):
                    if "err" in lines[j] and ("!= nil" in lines[j] or "HasError" in lines[j]):
                        has_check = True
                        break

                if not has_check and i + 1 < len(lines) and "return" not in lines[i + 1]:
                    stripped = line.strip()

                    if not stripped.startswith("//") and "defer" not in stripped:
                        violations.append(Violation(
                            file=filepath, line=i + 1,
                            rule="CODE-RED-016", severity="CODE-RED",
                            message="Error assigned but not checked within 3 lines. Handle or return.",
                            code_snippet=stripped[:120],
                        ))

    return violations


def check_php_throwable(lines: List[str], filepath: str) -> List[Violation]:
    """CODE-RED-017: PHP must catch Throwable, not just Exception."""
    violations = []

    for i, line in enumerate(lines):
        if re.search(r'catch\s*\(\s*\\?Exception\s', line):
            violations.append(Violation(
                file=filepath, line=i + 1,
                rule="CODE-RED-017", severity="CODE-RED",
                message="Catch `Throwable`, not `Exception`. PHP Errors are not Exceptions.",
                code_snippet=line.strip()[:120],
            ))

    return violations


def check_sequential_async(lines: List[str], filepath: str, lang: str) -> List[Violation]:
    """CODE-RED-018: Sequential independent async calls should use Promise.all / errgroup."""
    violations = []

    if lang not in ("typescript", "javascript"):
        return violations

    for i in range(len(lines) - 1):
        curr = lines[i].strip()
        nxt = lines[i + 1].strip()

        if curr.startswith("const ") and "await " in curr and nxt.startswith("const ") and "await " in nxt:
            violations.append(Violation(
                file=filepath, line=i + 2,
                rule="CODE-RED-018", severity="CODE-RED",
                message="Consecutive awaits detected. If independent, use Promise.all().",
                code_snippet=nxt[:120],
            ))

    return violations


def check_sql_injection(lines: List[str], filepath: str, lang: str) -> List[Violation]:
    """CODE-RED-019: No SQL string concatenation."""
    violations = []
    sql_keywords = re.compile(r'\b(SELECT|INSERT|UPDATE|DELETE|FROM|WHERE)\b', re.IGNORECASE)

    for i, line in enumerate(lines):
        stripped = line.strip()

        if stripped.startswith("//") or stripped.startswith("#"):
            continue

        if not sql_keywords.search(stripped):
            continue

        # Skip template literals / string interpolation that don't contain actual SQL verbs
        # (e.g. `Lines ${from}–${to}` triggers FROM but is not SQL)
        actual_sql_verbs = re.compile(r'\b(SELECT|INSERT\s+INTO|UPDATE\s+\w+\s+SET|DELETE\s+FROM)\b', re.IGNORECASE)

        if not actual_sql_verbs.search(stripped):
            continue

        # String concat with variables in SQL
        has_concat = False

        if lang == "go" and ('" +' in stripped or '+ "' in stripped):
            has_concat = True
        if lang in ("typescript", "javascript") and ("${" in stripped or '" +' in stripped or "' +" in stripped):
            if "${" in stripped and "`" in stripped:
                has_concat = True
        if lang == "php" and ('"$' in stripped or ".$" in stripped or '. $' in stripped):
            has_concat = True

        if has_concat:
            violations.append(Violation(
                file=filepath, line=i + 1,
                rule="CODE-RED-019", severity="CODE-RED",
                message="SQL string concatenation detected — use parameterized queries.",
                code_snippet=stripped[:120],
            ))

    return violations


def check_go_raw_errors(lines: List[str], filepath: str) -> List[Violation]:
    """CODE-RED-020: Go must use apperror, not errors.New/fmt.Errorf for stack traces."""
    violations = []

    for i, line in enumerate(lines):
        stripped = line.strip()

        if stripped.startswith("//"):
            continue

        if "errors.New(" in stripped:
            violations.append(Violation(
                file=filepath, line=i + 1,
                rule="CODE-RED-020", severity="CODE-RED",
                message="Use apperror.New() instead of errors.New() — stack trace required.",
                code_snippet=stripped[:120],
            ))

    return violations


def check_mixed_operators(lines: List[str], filepath: str) -> List[Violation]:
    """CODE-RED-021: Never mix && and || in a single expression."""
    violations = []

    for i, line in enumerate(lines):
        stripped = line.strip()

        if stripped.startswith("//") or stripped.startswith("*"):
            continue

        if "&&" in stripped and "||" in stripped:
            if stripped.startswith("if ") or stripped.startswith("if(") or stripped.startswith("} else if"):
                violations.append(Violation(
                    file=filepath, line=i + 1,
                    rule="CODE-RED-021", severity="CODE-RED",
                    message="Mixed && and || in condition. Extract to named booleans.",
                    code_snippet=stripped[:120],
                ))

    return violations


# ═══════════════════════════════════════════════════════════════════════
# Style Checkers
# ═══════════════════════════════════════════════════════════════════════

def check_style_rules(lines: List[str], filepath: str) -> List[Violation]:
    """STYLE-001/002/003/004: Blank-line-before-return, no-else-after-return, blank-line-after-brace, blank-line-before-if."""
    violations = []

    for i, line in enumerate(lines):
        stripped = line.strip()

        # STYLE-001: Blank line before return
        if stripped.startswith("return ") or stripped == "return":
            if i > 0:
                prev = lines[i - 1].strip()

                # Skip: return is the first statement in a block (prev is or ends with '{')
                if prev == "{" or prev.endswith("{") or prev.endswith("=> {") or prev.endswith("=> ("):
                    continue

                # Skip: return is the very first non-blank line of the function/file
                is_first_statement = all(l.strip() == "" or l.strip().startswith("//") or l.strip().startswith("#") or l.strip() == "{" for l in lines[:i])

                if is_first_statement:
                    continue

                if prev and not prev.startswith("//") and not prev.startswith("#"):
                    violations.append(Violation(
                        file=filepath, line=i + 1,
                        rule="STYLE-001", severity="STYLE",
                        message="Add a blank line before `return` statement.",
                        code_snippet=stripped[:120],
                    ))

        # STYLE-002: No else after return
        if stripped in ("} else {",) or stripped.startswith("} else if") or stripped.startswith("} else "):
            for j in range(i - 1, max(i - 6, -1), -1):
                ps = lines[j].strip()
                if ps.startswith("return ") or ps == "return":
                    violations.append(Violation(
                        file=filepath, line=i + 1,
                        rule="STYLE-002", severity="STYLE",
                        message="No `else` after `return`. Use early return pattern.",
                        code_snippet=stripped[:120],
                    ))
                    break
                if ps in ("}", "{"):
                    break

        # STYLE-003: Blank line after closing brace
        if stripped == "}" and i + 1 < len(lines):
            next_stripped = lines[i + 1].strip()
            skip = ("", "}", "//", ")", "} else", "case ", "default:")
            if next_stripped and not any(next_stripped.startswith(s) for s in skip if s):
                if next_stripped not in ("}", ")"):
                    violations.append(Violation(
                        file=filepath, line=i + 1,
                        rule="STYLE-003", severity="STYLE",
                        message="Add a blank line after closing `}`.",
                        code_snippet=stripped[:120],
                    ))

        # STYLE-004: Blank line before if/else if (when preceded by a statement)
        if stripped.startswith("if ") or stripped.startswith("if("):
            if i > 0:
                prev = lines[i - 1].strip()

                # Skip: if is the first statement after opening brace (prev is or ends with '{')
                if prev == "{" or prev.endswith("{") or prev.endswith("=> {"):
                    continue

                # Skip: if after another closing brace (already handled by STYLE-003)
                if prev == "}":
                    continue

                # prev must be a non-empty statement (not blank, not comment)
                if prev and not prev.startswith("//") and not prev.startswith("#") and not prev.startswith("*"):
                    violations.append(Violation(
                        file=filepath, line=i + 1,
                        rule="STYLE-004", severity="STYLE",
                        message="Add a blank line before `if` block.",
                        code_snippet=stripped[:120],
                    ))

    return violations

def check_generic_file_errors(lines: List[str], filepath: str, lang: str) -> List[Violation]:
    """CODE-RED-009: No generic file-not-found messages without exact path or reason."""
    violations = []

    # Patterns that indicate a generic file/path error without context
    generic_patterns = [
        re.compile(r'''(?:console\.(?:error|warn|log)|log\.(?:Error|Warn|Info|error|warn|info))\s*\(\s*["'`](?:file not found|path not found|missing file|file does not exist|no such file|file is missing|could not (?:find|open|read|load|access) (?:the )?file)["'`]\s*\)''', re.IGNORECASE),
        re.compile(r'''(?:throw\s+new\s+(?:Error|AppError|Exception)\s*\(\s*["'`](?:file not found|path not found|missing file|file does not exist|no such file|file is missing)["'`])''', re.IGNORECASE),
        re.compile(r'''fmt\.Errorf\s*\(\s*"(?:file not found|path not found|missing file|file does not exist|no such file)"''', re.IGNORECASE),
    ]

    # Patterns that indicate proper context (path variable or reason included)
    has_context = re.compile(r'''(?:%[sdvw]|%\{|\$\{|\bpath\b|\bfilePath\b|\bfilename\b|\breason\b|\bcause\b)''', re.IGNORECASE)

    for i, line in enumerate(lines):
        stripped = line.strip()

        if stripped.startswith("//") or stripped.startswith("#") or stripped.startswith("*"):
            continue

        for pat in generic_patterns:
            if pat.search(stripped) and not has_context.search(stripped):
                violations.append(Violation(
                    file=filepath,
                    line=i + 1,
                    rule="CODE-RED-009",
                    severity="CODE-RED",
                    message="Generic file/path error without exact path or failure reason. Include the file path and reason per the Code Red file-path logging rule.",
                    code_snippet=stripped[:120],
                ))
                break

    return violations


# ═══════════════════════════════════════════════════════════════════════
# Main Validation
# ═══════════════════════════════════════════════════════════════════════

# ═══════════════════════════════════════════════════════════════════════
# Boolean Principles — P2/P3/P5/P7 (added v1.5.0)
# ═══════════════════════════════════════════════════════════════════════

# P2: identifiers using negative words
NEGATIVE_BOOL_PATTERNS = [
    re.compile(r"\b(is|has|can|should|was|will)(Not|No)[A-Z]\w*"),
    re.compile(r"\b(is|has|can|should|was|will)(Not|No)_\w+"),
]

# P3: raw `!` applied to a function/method call
P3_BANG_CALL = re.compile(r"(?:^|[\s(=,&|!])!\s*[A-Za-z_$][\w$.]*\s*\(")

# P5: bare true/false as a positional argument in a call
P5_BARE_BOOL_ARG = re.compile(r"[A-Za-z_$][\w$.]*\s*\([^)\n]*?(?<![\w.])(true|false)(?![\w])(?:\s*,\s*[^)\n]*?)?\s*\)")

# P7: assignment inside a condition (excludes ==, !=, <=, >=, :=)
P7_ASSIGN_IN_COND = re.compile(r"\b(?:if|while)\s*\(?[^()=!<>]*?[^=!<>]=[^=][^)]*\)?\s*\{")


def check_negative_words(lines, filepath, lang):
    """CODE-RED-022 (P2): No negative words in boolean identifiers."""
    violations = []

    for i, line in enumerate(lines):
        stripped = line.strip()

        if isCommentLine(stripped):
            continue

        for pat in NEGATIVE_BOOL_PATTERNS:
            match = pat.search(line)

            if match is None:
                continue

            violations.append(Violation(
                file=filepath,
                line=i + 1,
                rule="CODE-RED-022",
                severity="CODE-RED",
                message=f'Negative-word boolean identifier "{match.group(0)}". Use a positive synonym (e.g. isPending, isInvalid, lacksAccess).',
                code_snippet=stripped[:120],
            ))

    return violations


def check_bang_on_call(lines, filepath, lang):
    """CODE-RED-023 (P3): Raw `!` on function/method calls is forbidden."""
    violations = []

    for i, line in enumerate(lines):
        stripped = line.strip()

        if isCommentLine(stripped):
            continue

        match = P3_BANG_CALL.search(line)

        if match is None:
            continue

        if isAllowedBangContext(line, match):
            continue

        violations.append(Violation(
            file=filepath,
            line=i + 1,
            rule="CODE-RED-023",
            severity="CODE-RED",
            message="Raw `!` on a function/method call is forbidden. Use a positive guard function or semantic inverse method.",
            code_snippet=stripped[:120],
        ))

    return violations


def check_bare_bool_args(lines, filepath, lang):
    """CODE-RED-024 (P5): Bare true/false as positional argument."""
    violations = []
    exempt_callers = ("expect", "assert", "should", "describe", "it", "test", "console.log", "JSON.stringify")

    for i, line in enumerate(lines):
        stripped = line.strip()

        if isCommentLine(stripped):
            continue

        if isAllowedBoolArgCall(stripped, exempt_callers):
            continue

        match = P5_BARE_BOOL_ARG.search(stripped)

        if match is None:
            continue

        violations.append(Violation(
            file=filepath,
            line=i + 1,
            rule="CODE-RED-024",
            severity="CODE-RED",
            message=f'Bare `{match.group(1)}` as positional argument. Use a named flag, options object, or dedicated method.',
            code_snippet=stripped[:120],
        ))

    return violations


def check_assignment_in_condition(lines, filepath, lang):
    """CODE-RED-025 (P7): No assignment inside if/while conditions."""
    violations = []

    for i, line in enumerate(lines):
        stripped = line.strip()

        if isCommentLine(stripped):
            continue

        if lang == "go" and ":=" in stripped:
            continue  # Go comma-ok / short var decl is exempt

        match = P7_ASSIGN_IN_COND.search(stripped)

        if match is None:
            continue

        violations.append(Violation(
            file=filepath,
            line=i + 1,
            rule="CODE-RED-025",
            severity="CODE-RED",
            message="Assignment inside if/while condition is forbidden. Hoist the assignment to a prior line.",
            code_snippet=stripped[:120],
        ))

    return violations


def isCommentLine(stripped):
    return stripped.startswith("//") or stripped.startswith("#") or stripped.startswith("*")


def isAllowedBangContext(line, match):
    # Skip `!=`, `!==`, `!!`
    after = line[match.end() - 1:]

    if after.startswith("!="):
        return True

    return False


def isAllowedBoolArgCall(stripped, exempt_callers):
    for caller in exempt_callers:
        if caller in stripped:
            return True

    return False


# ═══════════════════════════════════════════════════════════════════════
# Main Validation
# ═══════════════════════════════════════════════════════════════════════

def validate_file(filepath: str) -> List[Violation]:
    lang = detect_language(filepath)

    if not lang:
        return []

    # Skip test files, generated files, vendor
    basename = os.path.basename(filepath)
    normalized = filepath.replace("\\", "/")

    if any(skip in normalized for skip in [
        "vendor/", "node_modules/", "dist/", ".min.", "_test.go",
        "components/ui/",  # Auto-generated shadcn/ui — not business logic
    ]):
        return []

    if basename.endswith(".test.ts") or basename.endswith(".spec.ts"):
        return []

    try:
        with open(filepath, "r", encoding="utf-8", errors="ignore") as f:
            content = f.read()
    except (IOError, OSError) as exc:
        import sys as _sys; print(f"Error: {exc}", file=_sys.stderr)
        return []

    lines = content.split("\n")
    violations = []

    # Universal rules
    violations.extend(check_nested_if(lines, filepath))
    violations.extend(check_boolean_naming(lines, filepath, lang))
    violations.extend(check_magic_strings(lines, filepath, lang))
    violations.extend(check_magic_numbers(lines, filepath, lang))
    violations.extend(check_function_length(lines, filepath, lang))
    violations.extend(check_variable_mutation(lines, filepath, lang))
    violations.extend(check_file_length(lines, filepath))
    violations.extend(check_parameter_count(lines, filepath, lang))
    violations.extend(check_no_any_type(lines, filepath, lang))
    violations.extend(check_error_swallowing(lines, filepath, lang))
    violations.extend(check_sequential_async(lines, filepath, lang))
    violations.extend(check_sql_injection(lines, filepath, lang))
    violations.extend(check_mixed_operators(lines, filepath))
    violations.extend(check_style_rules(lines, filepath))
    violations.extend(check_generic_file_errors(lines, filepath, lang))
    violations.extend(check_negative_words(lines, filepath, lang))
    violations.extend(check_bang_on_call(lines, filepath, lang))
    violations.extend(check_bare_bool_args(lines, filepath, lang))
    violations.extend(check_assignment_in_condition(lines, filepath, lang))

    # Language-specific rules
    if lang == "go":
        violations.extend(check_go_specific(lines, filepath))
        violations.extend(check_go_raw_errors(lines, filepath))

    if lang == "php":
        violations.extend(check_php_throwable(lines, filepath))

    return violations


def main() -> int:
    import sys
    if hasattr(sys.stdout, "reconfigure"):
        sys.stdout.reconfigure(encoding="utf-8")
    parser = argparse.ArgumentParser(description="Cross-Language Coding Guidelines Validator")
    parser.add_argument("--path", default="src", help="Directory to scan (default: src)")
    parser.add_argument("--json", action="store_true", help="Output as JSON")
    parser.add_argument("--max-lines", type=int, default=15, help="Max function lines (default: 15)")
    args = parser.parse_args()

    report = ValidationReport()

    extensions = ("*.go", "*.ts", "*.tsx", "*.js", "*.jsx", "*.php", "*.rs")

    for ext in extensions:
        for filepath in glob.glob(os.path.join(args.path, "**", ext), recursive=True):
            report.total_files += 1
            violations = validate_file(filepath)

            for v in violations:
                report.violations.append(v)
                report.total_violations += 1
                report.by_rule[v.rule] += 1
                report.by_file[v.file] += 1

                if v.severity == "CODE-RED":
                    report.code_red_count += 1
                else:
                    report.style_count += 1

    if args.json:
        output = {
            "totalFiles": report.total_files,
            "totalViolations": report.total_violations,
            "codeRedCount": report.code_red_count,
            "styleCount": report.style_count,
            "byRule": dict(report.by_rule),
            "violations": [asdict(v) for v in report.violations],
        }
        print(json.dumps(output, indent=2))
    else:
        print_report(report)

    # Exit with error code if CODE RED violations found
    sys.exit(1 if report.code_red_count > 0 else 0)


def print_report(report: ValidationReport):
    print("=" * 72)
    print("  CODING GUIDELINES VALIDATION REPORT")
    print("=" * 72)
    print(f"  Files scanned:      {report.total_files}")
    print(f"  Total violations:   {report.total_violations}")
    print(f"  🔴 CODE RED:        {report.code_red_count}")
    print(f"  ⚠️  Style:           {report.style_count}")
    print("=" * 72)

    if not report.violations:
        print("\n  ✅ ALL CLEAR — No violations found.\n")
        return

    # Group by file
    by_file = defaultdict(list)

    for v in report.violations:
        by_file[v.file].append(v)

    for filepath in sorted(by_file.keys()):
        violations = by_file[filepath]
        print(f"\n📄 {filepath} ({len(violations)} violations)")

        for v in sorted(violations, key=lambda x: x.line):
            icon = "🔴" if v.severity == "CODE-RED" else "⚠️ "
            print(f"  {icon} L{v.line:<5d} [{v.rule}] {v.message}")

            if v.code_snippet:
                print(f"           │ {v.line}: {v.code_snippet}")

    print(f"\n{'=' * 72}")
    print("  BY RULE:")

    for rule, count in sorted(report.by_rule.items()):
        print(f"    {rule}: {count}")

    print(f"{'=' * 72}\n")

    if report.code_red_count > 0:
        print(f"  ❌ FAILED — {report.code_red_count} CODE RED violation(s) must be fixed before merge.\n")
    else:
        print(f"  ⚠️  PASSED with {report.style_count} style warning(s).\n")


if __name__ == "__main__":
    main()
