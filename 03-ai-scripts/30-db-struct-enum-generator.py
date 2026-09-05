#!/usr/bin/env python3
"""
30-db-struct-enum-generator.py — Inspects Go model structs and auto-generates type-safe column enums.

Usage:
  python 03-ai-scripts/30-db-struct-enum-generator.py --dir gitmap/store
  python 03-ai-scripts/30-db-struct-enum-generator.py --file gitmap/pipelinedb/pipeline_split_db.go
  python 03-ai-scripts/30-db-struct-enum-generator.py --dry-run
"""

from __future__ import annotations

import argparse
from importlib import import_module
from pathlib import Path
import re
import sys

sys.path.insert(0, str(Path(__file__).parent))
engine = import_module("02-shared-engine")

ExitCodeType = engine.ExitCodeType

STRUCT_PATTERN = re.compile(r"type\s+([A-Za-z0-9_]+)\s+struct\s*\{([^}]+)\}", re.MULTILINE)
FIELD_PATTERN = re.compile(r"^\s*([A-Za-z0-9_]+)\s+([A-Za-z0-9_\[\]*]+)(?:\s+`([^`]+)`)?", re.MULTILINE)
PACKAGE_PATTERN = re.compile(r"package\s+([A-Za-z0-9_]+)")


def parse_structs_from_file(file_path: Path) -> tuple[str, list[dict]]:
    content = file_path.read_text(encoding="utf-8", errors="replace")
    pkg_match = PACKAGE_PATTERN.search(content)
    pkg_name = pkg_match.group(1) if pkg_match else "models"

    structs = []
    for match in STRUCT_PATTERN.finditer(content):
        struct_name = match.group(1)
        body = match.group(2)

        fields = []
        for f_match in FIELD_PATTERN.finditer(body):
            field_name = f_match.group(1)
            field_type = f_match.group(2)
            tags = f_match.group(3) or ""

            # Check if public field and not embedded
            if field_name[0].isupper() and not field_name.startswith("XXX_"):
                fields.append({
                    "name": field_name,
                    "type": field_type,
                    "tags": tags,
                })

        if fields:
            structs.append({
                "name": struct_name,
                "fields": fields,
            })

    return pkg_name, structs


def generate_enums_for_struct(struct_info: dict) -> str:
    s_name = struct_info["name"]
    enum_type = f"{s_name}FieldType"
    reg_type = f"{s_name[0].lower() + s_name[1:]}DbRegistry"
    valid_map = f"{s_name[0].lower() + s_name[1:]}ValidMap"
    db_var = f"{s_name}Db"
    field_var = f"{s_name}Field"

    lines = [
        f"// {enum_type} represents column name enums for {s_name}.",
        f"type {enum_type} string",
        "",
        f"// Name returns the identifier name of the field enum.",
        f"func (e {enum_type}) Name() string {{",
        "\treturn string(e)",
        "}",
        "",
        f"// String returns the string representation of the field enum.",
        f"func (e {enum_type}) String() string {{",
        "\treturn string(e)",
        "}",
        "",
        f"// Value returns the raw string value of the field enum.",
        f"func (e {enum_type}) Value() string {{",
        "\treturn string(e)",
        "}",
        "",
        f"// IsCompare checks equality against another field enum object.",
        f"func (e {enum_type}) IsCompare(target {enum_type}) bool {{",
        "\treturn e == target",
        "}",
        "",
        f"// IsEnum checks whether this field enum exists in the valid enum map.",
        f"func (e {enum_type}) IsEnum() bool {{",
        f"\treturn {valid_map}[e]",
        "}",
        "",
        f"// MarshalJSON implements json.Marshaler.",
        f"func (e {enum_type}) MarshalJSON() ([]byte, error) {{",
        "\treturn json.Marshal(string(e))",
        "}",
        "",
        f"// UnmarshalJSON implements json.Unmarshaler with strict map validation.",
        f"func (e *{enum_type}) UnmarshalJSON(data []byte) error {{",
        "\tvar s string",
        "\tif err := json.Unmarshal(data, &s); err != nil {",
        "\t\treturn err",
        "\t}",
        f"\ttarget := {enum_type}(s)",
        f"\tif !{valid_map}[target] {{",
        f'\t\treturn fmt.Errorf("invalid %s enum: %s", "{enum_type}", s)',
        "\t}",
        "\t*e = target",
        "\treturn nil",
        "}",
        "",
        f"// ToJSON converts the field enum to a JSON string representation, returning an AppError on failure.",
        f"func (e {enum_type}) ToJSON() (string, *apperror.AppError) {{",
        "\tb, err := json.Marshal(string(e))",
        "\tif err != nil {",
        '\t\treturn "", apperror.WrapSimple(err, "serialize field to json")',
        "\t}",
        "\treturn string(b), nil",
        "}",
        "",
        f"// FromJSON parses a field enum from a JSON string representation, returning an AppError on failure.",
        f"func (e *{enum_type}) FromJSON(s string) *apperror.AppError {{",
        "\tvar str string",
        "\tif err := json.Unmarshal([]byte(s), &str); err != nil {",
        '\t\treturn apperror.WrapSimple(err, "deserialize field from json")',
        "\t}",
        f"\ttarget := {enum_type}(str)",
        f"\tif !{valid_map}[target] {{",
        f'\t\treturn apperror.WrapSimple(fmt.Errorf("invalid %s enum: %s", "{enum_type}", str), "validate field enum from json")',
        "\t}",
        "\t*e = target",
        "\treturn nil",
        "}",
        "",
    ]

    for f in struct_info["fields"]:
        f_name = f["name"]
        lines.extend([
            f"// Is{f_name} checks whether this field enum instance is {f_name}.",
            f"func (e {enum_type}) Is{f_name}() bool {{",
            f"\treturn e == {db_var}.{f_name}",
            "}",
            "",
        ])

    lines.append(f"type {reg_type} struct {{")
    for f in struct_info["fields"]:
        lines.append(f"\t{f['name']} {enum_type}")
    lines.append("}")
    lines.append("")

    lines.extend([
        f"// All returns a slice of all field enums in {s_name}.",
        f"func (r {reg_type}) All() []{enum_type} {{",
        f"\treturn []{enum_type}{{",
    ])
    for f in struct_info["fields"]:
        lines.append(f"\t\tr.{f['name']},")
    lines.extend([
        "\t}",
        "}",
        "",
        f"// Names returns a slice of string names for all fields in {s_name}.",
        f"func (r {reg_type}) Names() []string {{",
        "\treturn []string{",
    ])
    for f in struct_info["fields"]:
        lines.append(f'\t\t"{f["name"]}",')
    lines.extend([
        "\t}",
        "}",
        "",
        f"// IsEnum checks whether the target object matches any registered field enum in {s_name}.",
        f"func (r {reg_type}) IsEnum(target {enum_type}) bool {{",
        f"\treturn {valid_map}[target]",
        "}",
        "",
    ])

    for f in struct_info["fields"]:
        f_name = f["name"]
        lines.extend([
            f"// Is{f_name} checks whether the target object is {f_name}.",
            f"func (r {reg_type}) Is{f_name}(target {enum_type}) bool {{",
            f"\treturn target == r.{f_name}",
            "}",
            "",
        ])

    lines.extend([
        f"// ToJSON converts the field registry to a JSON string representation, returning an AppError on failure.",
        f"func (r {reg_type}) ToJSON() (string, *apperror.AppError) {{",
        "\tb, err := json.Marshal(r)",
        "\tif err != nil {",
        '\t\treturn "", apperror.WrapSimple(err, "serialize registry to json")',
        "\t}",
        "\treturn string(b), nil",
        "}",
        "",
    ])

    lines.append(f"// {db_var} provides scoped access to field enums: {db_var}.<Field>.")
    lines.append(f"var {db_var} = {reg_type}{{")
    for f in struct_info["fields"]:
        f_name = f["name"]
        lines.append(f'\t{f_name}: "{f_name}",')
    lines.append("}")
    lines.append("")

    lines.append(f"// {valid_map} provides O(1) map validation for field enums.")
    lines.append(f"var {valid_map} = map[{enum_type}]bool{{")
    for f in struct_info["fields"]:
        f_name = f["name"]
        lines.append(f"\t{db_var}.{f_name}: true,")
    lines.append("}")
    lines.append("")

    lines.append(f"// {field_var} is an alias to {db_var}.")
    lines.append(f"var {field_var} = {db_var}")
    lines.append("")
    lines.append(f"// {s_name}Table represents the canonical table name.")
    lines.append(f'const {s_name}Table = "{s_name}"')
    lines.append("")

    return "\n".join(lines)


def process_target_file(target_file: Path, dry_run: bool = False) -> bool:
    pkg_name, structs = parse_structs_from_file(target_file)
    if not structs:
        return False

    out_file = target_file.parent / f"{target_file.stem}_fields_gen.go"

    import_apperror = '\t"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"\n' if pkg_name != "apperror" else ""

    header = f"""// Code generated by gitmap db generate. DO NOT EDIT.

package {pkg_name}

import (
\t"encoding/json"
\t"fmt"
{import_apperror})

"""
    body = "\n".join(generate_enums_for_struct(s) for s in structs)
    full_content = header + body

    if dry_run:
        print(f"[DRY RUN] Would generate: {out_file} ({len(structs)} structs)")
        return True

    out_file.write_text(full_content, encoding="utf-8")
    print(f"  ✔ Generated {out_file.name} ({len(structs)} structs)")
    return True


def main() -> int:
    parser = argparse.ArgumentParser(description="Generate type-safe column name enums from Go model structs")
    parser.add_argument("--dir", help="Directory containing Go model files", default="")
    parser.add_argument("--file", help="Specific Go model file to parse", default="")
    parser.add_argument("--dry-run", action="store_true", help="Preview output without writing files")

    args = parser.parse_args()
    repo_root = Path(__file__).resolve().parent.parent

    if args.file:
        file_path = Path(args.file)
        if not file_path.is_absolute():
            file_path = repo_root / file_path
        if not file_path.is_file():
            print(f"Error: file not found: {file_path}", file=sys.stderr)
            return 1
        process_target_file(file_path, args.dry_run)
        return 0

    target_dir = Path(args.dir) if args.dir else repo_root / "gitmap" / "pipelinedb"
    if not target_dir.is_absolute():
        target_dir = repo_root / target_dir

    if not target_dir.is_dir():
        print(f"Error: directory not found: {target_dir}", file=sys.stderr)
        return 1

    generated_count = 0
    for go_file in target_dir.glob("*.go"):
        if go_file.name.endswith("_test.go") or go_file.name.endswith("_gen.go"):
            continue
        if process_target_file(go_file, args.dry_run):
            generated_count += 1

    print(f"Completed: processed files in {target_dir} ({generated_count} files generated enums).")
    return 0


if __name__ == "__main__":
    sys.exit(main())
