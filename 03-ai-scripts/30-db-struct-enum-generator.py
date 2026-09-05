#!/usr/bin/env python3
"""
30-db-struct-enum-generator.py — Inspects Go model structs and auto-generates type-safe column enums.

Usage:
  python 03-ai-scripts/30-db-struct-enum-generator.py --dir gitmap/store
  python 03-ai-scripts/30-db-struct-enum-generator.py --file gitmap/pipelinedb/pipeline_split_db.go
  python 03-ai-scripts/30-db-struct-enum-generator.py --file gitmap/pipelinedb/pipeline_split_db.go --out-dir gitmap/generated/db/pipelinedb
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

STRUCT_PATTERN = re.compile(r"type\s+([A-Z][A-Za-z0-9_]+)\s+struct\s*\{([^}]+)\}", re.MULTILINE)
FIELD_PATTERN = re.compile(r"^\s*([A-Za-z0-9_]+)\s+([A-Za-z0-9_\[\]*]+)(?:\s+`([^`]+)`)?", re.MULTILINE)
PACKAGE_PATTERN = re.compile(r"package\s+([A-Za-z0-9_]+)")


def parse_structs_from_file(file_path: Path) -> tuple[str, list[dict]]:
    content = file_path.read_text(encoding="utf-8", errors="replace")
    pkg_match = PACKAGE_PATTERN.search(content)
    pkg_name = pkg_match.group(1) if pkg_match else "models"

    structs = []
    for match in STRUCT_PATTERN.finditer(content):
        struct_name = match.group(1)
        if struct_name.endswith(("DbRegistry", "DbRepo", "Repository", "Registry")):
            continue
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


def generate_enums_for_struct(struct_info: dict, pkg_name: str = "") -> str:
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

    pkg_qual = "" if pkg_name == "dbengine" else "dbengine."
    scan_func = f"Scan{s_name}"
    repo_type = f"{s_name}DbRepo"
    qb_type = f"{s_name}QueryBuilder"
    generic_repo_type = f"{s_name}Repository"
    short_name = s_name[:-6] if s_name.endswith("Record") else s_name
    has_alias = short_name != s_name
    short_repo_type = f"{short_name}DbRepo"

    # Scanner function
    lines.append(f"// {scan_func} maps a database row scanner to a {s_name} entity.")
    lines.append(f"func {scan_func}(row {pkg_qual}RowScanner) (*{s_name}, error) {{")
    lines.append(f"\tvar item {s_name}")
    lines.append("\tvar (")
    for f in struct_info["fields"]:
        lines.append(f"\t\traw_{f['name']} any")
    lines.append("\t)")
    lines.append("\terr := row.Scan(")
    for f in struct_info["fields"]:
        lines.append(f"\t\t&raw_{f['name']},")
    lines.append("\t)")
    lines.append("\tif err != nil {")
    lines.append("\t\treturn nil, err")
    lines.append("\t}")
    lines.append("")
    for f in struct_info["fields"]:
        helper = get_scan_helper(f["type"], pkg_qual)
        lines.append(f"\titem.{f['name']} = {helper}(raw_{f['name']})")
    lines.append("\treturn &item, nil")
    lines.append("}")
    lines.append("")

    # Repository struct & methods
    lines.extend([
        f"// {repo_type} provides typed database repository access for {s_name}.",
        f"type {repo_type} struct {{",
        f"\tdb   *{pkg_qual}DbWrapper",
        f"\trepo *{generic_repo_type}",
        "}",
        "",
    ])

    lines.extend([
        f"// New{repo_type} initializes a typed repository for {s_name}.",
        f"func New{repo_type}(db *{pkg_qual}DbWrapper) *{repo_type} {{",
        f"\trepo := {pkg_qual}NewRepository[{s_name}, {enum_type}](",
        "\t\tdb,",
        f"\t\t{s_name}Table,",
        f"\t\t{scan_func},",
        "\t)",
        f"\treturn &{repo_type}{{",
        "\t\tdb:   db,",
        "\t\trepo: repo,",
        "\t}",
        "}",
        "",
    ])

    if has_alias:
        lines.extend([
            f"// New{short_repo_type} is an alias constructor for {repo_type}.",
            f"func New{short_repo_type}(db *{pkg_qual}DbWrapper) *{short_repo_type} {{",
            f"\treturn New{repo_type}(db)",
            "}",
            "",
        ])

    lines.extend([
        f"// Db returns the underlying DbWrapper.",
        f"func (r *{repo_type}) Db() *{pkg_qual}DbWrapper {{",
        "\treturn r.db",
        "}",
        "",
        f"// Repo returns the underlying generic Repository.",
        f"func (r *{repo_type}) Repo() *{generic_repo_type} {{",
        "\treturn r.repo",
        "}",
        "",
        f"// Query returns a fluent QueryBuilder initialized with all standard fields projected.",
        f"func (r *{repo_type}) Query() *{qb_type} {{",
        f"\treturn r.repo.Query().Select({db_var}.All()...)",
        "}",
        "",
        f"// QueryBare returns a fluent QueryBuilder without any pre-selected fields.",
        f"func (r *{repo_type}) QueryBare() *{qb_type} {{",
        "\treturn r.repo.Query()",
        "}",
        "",
        f"// FindAll executes the query selecting all fields and returns a ListResult envelope.",
        f"func (r *{repo_type}) FindAll(ctx context.Context) {pkg_qual}ListResult[{s_name}] {{",
        "\treturn r.Query().FindAll(ctx)",
        "}",
        "",
        f"// First executes the query selecting all fields and returns the first record in an EntityResult envelope.",
        f"func (r *{repo_type}) First(ctx context.Context) {pkg_qual}EntityResult[{s_name}] {{",
        "\treturn r.Query().First(ctx)",
        "}",
        "",
        f"// Count returns the total number of records matching the query.",
        f"func (r *{repo_type}) Count(ctx context.Context) {pkg_qual}Int64Result {{",
        "\treturn r.Query().Count(ctx)",
        "}",
        "",
    ])

    return "\n".join(lines)


def get_scan_helper(field_type: str, pkg_qual: str) -> str:
    ft = field_type.strip("*")
    if ft == "string":
        return f"{pkg_qual}ScanString"
    if ft in ("int", "int32"):
        return f"{pkg_qual}ScanInt"
    if ft == "int64":
        return f"{pkg_qual}ScanInt64"
    if ft in ("uint64", "uint"):
        return f"{pkg_qual}ScanUint64"
    if ft in ("uint32", "uint16", "uint8"):
        return f"{pkg_qual}ScanUint"
    if ft == "bool":
        return f"{pkg_qual}ScanBool"
    if ft in ("float64", "float32"):
        return f"{pkg_qual}ScanFloat64"
    return f"{pkg_qual}ScanString"


def generate_consts_content(pkg_name: str, structs: list[dict]) -> str:
    pkg_qual = "" if pkg_name == "dbengine" else "dbengine."
    import_dbengine = '\t"github.com/alimtvnetwork/gitmap-v28/gitmap/dbengine"\n' if pkg_name != "dbengine" else ""

    lines = [
        "// Code generated by gitmap db generate. DO NOT EDIT.",
        "",
        f"package {pkg_name}",
        "",
    ]
    if import_dbengine:
        lines.extend([
            "import (",
            f"{import_dbengine})",
            "",
        ])

    # 1. Canonical table name constants
    lines.append("// Canonical table name constants.")
    lines.append("const (")
    for s in structs:
        s_name = s["name"]
        short_name = s_name[:-6] if s_name.endswith("Record") else s_name
        lines.append(f'\t{s_name}Table = "{s_name}"')
        if short_name != s_name:
            lines.append(f"\t{short_name}Table = {s_name}Table")
    lines.append(")")
    lines.append("")

    # 2. Dedicated QueryBuilder single data type aliases
    lines.append("// Dedicated QueryBuilder type aliases.")
    lines.append("type (")
    for s in structs:
        s_name = s["name"]
        enum_type = f"{s_name}FieldType"
        qb_type = f"{s_name}QueryBuilder"
        short_name = s_name[:-6] if s_name.endswith("Record") else s_name
        lines.append(f"\t// {qb_type} is the dedicated query builder for {s_name}.")
        lines.append(f"\t{qb_type} = {pkg_qual}QueryBuilder[{s_name}, {enum_type}]")
        if short_name != s_name:
            short_qb_type = f"{short_name}QueryBuilder"
            lines.append(f"\t// {short_qb_type} is an alias to {qb_type} for concise business usage.")
            lines.append(f"\t{short_qb_type} = {qb_type}")
    lines.append(")")
    lines.append("")

    # 3. Dedicated generic Repository type aliases
    lines.append("// Dedicated generic Repository type aliases.")
    lines.append("type (")
    for s in structs:
        s_name = s["name"]
        enum_type = f"{s_name}FieldType"
        generic_repo_type = f"{s_name}Repository"
        short_name = s_name[:-6] if s_name.endswith("Record") else s_name
        lines.append(f"\t// {generic_repo_type} is the dedicated generic repository for {s_name}.")
        lines.append(f"\t{generic_repo_type} = {pkg_qual}Repository[{s_name}, {enum_type}]")
        if short_name != s_name:
            short_generic_repo_type = f"{short_name}Repository"
            lines.append(f"\t// {short_generic_repo_type} is an alias to {generic_repo_type} for concise business usage.")
            lines.append(f"\t{short_generic_repo_type} = {generic_repo_type}")
    lines.append(")")
    lines.append("")

    # 4. Typed DbRepo aliases for concise business usage
    repo_aliases = []
    for s in structs:
        s_name = s["name"]
        short_name = s_name[:-6] if s_name.endswith("Record") else s_name
        if short_name != s_name:
            repo_type = f"{s_name}DbRepo"
            short_repo_type = f"{short_name}DbRepo"
            repo_aliases.append((short_repo_type, repo_type))

    if repo_aliases:
        lines.append("// Typed DbRepo aliases for concise business usage.")
        lines.append("type (")
        for short_repo_type, repo_type in repo_aliases:
            lines.append(f"\t// {short_repo_type} is an alias to {repo_type} for concise business usage.")
            lines.append(f"\t{short_repo_type} = {repo_type}")
        lines.append(")")
        lines.append("")

    return "\n".join(lines)


def to_snake_case(name: str) -> str:
    s = re.sub(r'(.)([A-Z][a-z]+)', r'\1_\2', name)
    return re.sub(r'([a-z0-9])([A-Z])', r'\1_\2', s).lower()


def process_target_file(target_file: Path, out_dir: Path | None = None, dry_run: bool = False) -> bool:
    pkg_name, structs = parse_structs_from_file(target_file)
    if not structs:
        return False

    dest_dir = out_dir if out_dir is not None else target_file.parent
    dest_dir.mkdir(parents=True, exist_ok=True)
    consts_file = dest_dir / "consts.go"

    import_apperror = '\t"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"\n' if pkg_name != "apperror" else ""
    import_dbengine = '\t"github.com/alimtvnetwork/gitmap-v28/gitmap/dbengine"\n' if pkg_name != "dbengine" else ""

    header = f"""// Code generated by gitmap db generate. DO NOT EDIT.

package {pkg_name}

import (
\t"context"
\t"encoding/json"
\t"fmt"
{import_apperror}{import_dbengine})

"""
    consts_content = generate_consts_content(pkg_name, structs)

    if dry_run:
        print(f"[DRY RUN] Would generate: {consts_file} ({len(structs)} structs)")
        return True

    # Remove any legacy _gen.go files in dest_dir
    for legacy_file in dest_dir.glob("*_gen.go"):
        legacy_file.unlink()

    consts_file.write_text(consts_content, encoding="utf-8")
    print(f"  ✔ Generated {consts_file} ({len(structs)} structs)")

    for s in structs:
        s_snake = to_snake_case(s["name"])
        def_file = dest_dir / f"{s_snake}.go"
        enums_body = generate_enums_for_struct(s, pkg_name)
        if def_file.is_file():
            existing = def_file.read_text(encoding="utf-8")
            if f"type {s['name']} struct" in existing:
                # File already has struct definition; ensure enums are appended or updated
                print(f"  ✔ Preserved definition in {def_file}")
                continue
        def_file.write_text(header + enums_body, encoding="utf-8")
        print(f"  ✔ Generated {def_file}")

    return True


def main() -> int:
    parser = argparse.ArgumentParser(description="Generate type-safe column name enums from Go model structs")
    parser.add_argument("--dir", help="Directory containing Go model files", default="")
    parser.add_argument("--file", help="Specific Go model file to parse", default="")
    parser.add_argument("--out-dir", help="Explicit output directory for generated files", default="")
    parser.add_argument("--dry-run", action="store_true", help="Preview output without writing files")

    args = parser.parse_args()
    repo_root = Path(__file__).resolve().parent.parent

    out_dir_path: Path | None = None
    if args.out_dir:
        out_dir_path = Path(args.out_dir)
        if not out_dir_path.is_absolute():
            out_dir_path = repo_root / out_dir_path

    if args.file:
        file_path = Path(args.file)
        if not file_path.is_absolute():
            file_path = repo_root / file_path
        if not file_path.is_file():
            print(f"Error: file not found: {file_path}", file=sys.stderr)
            return 1
        process_target_file(file_path, out_dir_path, args.dry_run)
        return 0

    target_dir = Path(args.dir) if args.dir else repo_root / "gitmap" / "pipelinedb"
    if not target_dir.is_absolute():
        target_dir = repo_root / target_dir

    if not target_dir.is_dir():
        print(f"Error: directory not found: {target_dir}", file=sys.stderr)
        return 1

    # Scan all model files in target_dir, collect all structs for a unified consts.go
    all_structs = []
    pkg_name = "pipelinedb"
    for go_file in target_dir.glob("*.go"):
        if go_file.name.endswith("_test.go") or go_file.name.endswith("_gen.go") or go_file.name == "consts.go":
            continue
        p, structs = parse_structs_from_file(go_file)
        if structs:
            pkg_name = p
            all_structs.extend(structs)

    dest_dir = out_dir_path if out_dir_path is not None else target_dir
    dest_dir.mkdir(parents=True, exist_ok=True)
    consts_file = dest_dir / "consts.go"
    consts_content = generate_consts_content(pkg_name, all_structs)

    if not args.dry_run:
        for legacy_file in dest_dir.glob("*_gen.go"):
            legacy_file.unlink()
        consts_file.write_text(consts_content, encoding="utf-8")
        print(f"  ✔ Generated unified {consts_file} ({len(all_structs)} structs)")
    else:
        print(f"[DRY RUN] Would generate unified {consts_file} ({len(all_structs)} structs)")

    print(f"Completed: processed files in {target_dir} ({len(all_structs)} structs in {consts_file.name}).")
    return 0


if __name__ == "__main__":
    sys.exit(main())

