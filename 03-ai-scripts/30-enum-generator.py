#!/usr/bin/env python3
"""
30-enum-generator.py — Multi-file type-safe Go enum scaffolder conforming to repo standards.

Scaffolds 4 standard files in a dedicated package folder (04-code/golang/pkg/enum/{name}type/):
  - variant.go: Enum type, constants, predicates, DRY JSON serialization.
  - vars.go: Labels array, baseenumer.CompileMap, All(), Values(), Parse() Result.
  - variant_test.go: Complete unit tests (interfaces, properties, predicates, names, JSON roundtrips).
  - readme.md: Package documentation.

Supports:
  - Backing types: byte (default), int (uint16), string
  - Input modes: CLI comma-separated items, JSON config file (--config), inline JSON (--json)
  - Flags: --dry-run, --overwrite
"""

from __future__ import annotations

import argparse
from importlib import import_module
import json
from pathlib import Path
import sys
from typing import Any

sys.path.insert(0, str(Path(__file__).parent))
engine = import_module("02-shared-engine")

EXIT_SUCCESS = 0
EXIT_FAILURE = 1
EXIT_USAGE_ERROR = 2

write_file_lf = engine.write_file_lf


def parse_cli_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Multi-file type-safe Go enum scaffolder")
    parser.add_argument("--name", "-n", help="Enum base name (e.g. DeliveryMode, ProcessState)")
    parser.add_argument("--type", "-t", choices=["byte", "uint8", "uint16", "int", "string"], default="byte")
    parser.add_argument("--items", "-m", help="Comma-separated member names (e.g. Fast,Standard,Slow)")
    parser.add_argument("--zero-value", "-z", default="Invalid", help="Zero-value label (Invalid or Unknown)")
    parser.add_argument("--package", "-p", help="Target Go package name (defaults to {name.lower()}type)")
    parser.add_argument("--out", "-o", "--target-dir", "-d", dest="target_dir", help="Target output folder")
    parser.add_argument("--config", "-c", help="Path to JSON config file")
    parser.add_argument("--json", help="Inline JSON config string or '-' for stdin")
    parser.add_argument("--overwrite", action="store_true", help="Overwrite existing files if present")
    parser.add_argument("--dry-run", action="store_true", help="Preview generated files without writing")
    return parser.parse_args()


def load_input_source(args: argparse.Namespace) -> dict[str, Any] | None:
    if args.config:
        return json.loads(Path(args.config).read_text(encoding="utf-8"))
    if args.json:
        raw = sys.stdin.read() if args.json == "-" else args.json
        return json.loads(raw)
    return None


def load_raw_config(args: argparse.Namespace) -> dict[str, Any]:
    cfg = load_input_source(args)
    if cfg:
        return cfg
    if not args.name or not args.items:
        raise ValueError("Either --config/--json or --name and --items required")
    return {
        "name": args.name,
        "type": args.type,
        "items": args.items,
        "zero_value": args.zero_value,
        "package": args.package,
        "target_dir": args.target_dir,
    }


def build_items_list(raw_items: Any) -> list[str]:
    if isinstance(raw_items, list):
        return [str(i).strip() for i in raw_items if str(i).strip()]
    if isinstance(raw_items, str):
        return [i.strip() for i in raw_items.split(",") if i.strip()]
    return []


def resolve_type_metadata(backing_type: str, zero_value: str) -> dict[str, Any]:
    norm_type = backing_type.lower()
    is_byte = norm_type in ("byte", "uint8")
    is_str = norm_type == "string"
    is_int = norm_type in ("int", "uint16")
    underlying = "string" if is_str else ("uint16" if is_int else "byte")
    zero = "Unknown" if is_str else zero_value
    return {
        "is_byte": is_byte,
        "is_string": is_str,
        "is_int": is_int,
        "underlying": underlying,
        "zero_value": zero,
    }


def resolve_target_dir(pkg: str, raw_dir: str | None) -> str:
    if raw_dir:
        return raw_dir
    return f"04-code/golang/pkg/enum/{pkg}"


def resolve_pkg_name(raw_pkg: str | None, name: str, raw_dir: str | None = None) -> str:
    if raw_pkg:
        return raw_pkg.strip().lower()
    if raw_dir:
        folder = Path(raw_dir.strip().replace("\\", "/")).name.lower()
        if folder.endswith("type"):
            return folder
    low = name.strip().lower()
    return low if low.endswith("type") else f"{low}type"


def normalize_config(raw: dict[str, Any]) -> dict[str, Any]:
    name = raw["name"].strip()
    items = build_items_list(raw.get("items", []))
    if not items:
        raise ValueError("At least one enum item must be provided")
    pkg = resolve_pkg_name(raw.get("package"), name, raw.get("target_dir"))
    tmeta = resolve_type_metadata(raw.get("type", "byte"), raw.get("zero_value", "Invalid"))
    tdir = resolve_target_dir(pkg, raw.get("target_dir"))
    return {"name": name, "items": items, "package": pkg, "target_dir": tdir, **tmeta}



def build_variant_header(ctx: dict[str, Any]) -> list[str]:
    imports = ['\t"coding-guidelines/common/pkg/baseenumer"']
    return [f"package {ctx['package']}", "", "import ("] + imports + [")", ""]


def build_variant_consts(ctx: dict[str, Any]) -> list[str]:
    if ctx["is_string"]:
        lines = ["const (", f'\t{ctx["zero_value"]} Variant = ""']
        for item in ctx["items"]:
            lines.append(f'\t{item} Variant = "{item}"')
        return lines + [")", ""]
    lines = ["const (", f"\t{ctx['zero_value']} Variant = iota"]
    for item in ctx["items"]:
        lines.append(f"\t{item}")
    return lines + [")", ""]


def resolve_assertion_types(ctx: dict[str, Any]) -> list[str]:
    if ctx["is_string"]:
        return ['\t_ baseenumer.BaseEnumer = Variant("")', '\t_ baseenumer.StringEnumer = Variant("")']
    lines = ["\t_ baseenumer.BaseEnumer = Variant(0)"]
    if ctx["is_byte"]:
        lines.append("\t_ baseenumer.ByteEnumer = Variant(0)")
    lines.append("\t_ baseenumer.NumberEnumer = Variant(0)")
    return lines


def build_variant_assertions(ctx: dict[str, Any]) -> list[str]:
    m_zero = 'Variant("")' if ctx["is_string"] else "Variant(0)"
    lines = ["var ("] + resolve_assertion_types(ctx)
    lines.extend([
        f"\t_ baseenumer.BoundedEnumer[Variant] = {m_zero}",
        ")",
        "",
    ])
    return lines



def build_variant_byte_methods(ctx: dict[str, Any]) -> list[str]:
    if not ctx["is_byte"]:
        return []
    return [
        "func (v Variant) Value() byte {\n\treturn byte(v)\n}\n",
        "func (v Variant) Byte() byte {\n\treturn byte(v)\n}\n",
        "func (v Variant) ValueByte() byte {\n\treturn byte(v)\n}\n",
        "func (v Variant) Bytes() []byte {\n\treturn []byte{byte(v)}\n}\n",
    ]


def build_variant_int_methods(ctx: dict[str, Any]) -> list[str]:
    if ctx["is_string"]:
        return ["func (v Variant) Value() string {\n\treturn string(v)\n}\n"]
    res = []
    if not ctx["is_byte"]:
        res.append(f"func (v Variant) Value() {ctx['underlying']} {{\n\treturn {ctx['underlying']}(v)\n}}\n")
    res.extend([
        "func (v Variant) Int() int {\n\treturn int(v)\n}\n",
        "func (v Variant) Code() uint16 {\n\treturn uint16(v)\n}\n",
    ])
    return res


def build_variant_collection_methods() -> list[str]:
    return [
        "func (v Variant) All() []Variant {\n\treturn All()\n}\n",
        "func (v Variant) Values() []string {\n\treturn Values()\n}\n",
    ]


def build_string_predicates() -> list[str]:
    return [
        "func (v Variant) IsValid() bool {\n\treturn variantRegistry[v]\n}\n",
        "func (v Variant) IsEnum() bool {\n\treturn variantRegistry[v]\n}\n",
    ]


def build_numeric_predicates(first: str, last: str) -> list[str]:
    return [
        f"func (v Variant) IsValid() bool {{\n\treturn baseenumer.IsBetween(v, {first}, {last})\n}}\n",
        f"func (v Variant) IsInvalid() bool {{\n\treturn baseenumer.IsNotBetween(v, {first}, {last})\n}}\n",
        "func (v Variant) IsEnum() bool {\n\treturn v.IsValid()\n}\n",
    ]


def build_variant_predicates(ctx: dict[str, Any]) -> list[str]:
    if ctx["is_string"]:
        return build_string_predicates()
    return build_numeric_predicates(ctx["items"][0], ctx["items"][-1])


def build_variant_boundary_methods() -> list[str]:
    return [
        "func (v Variant) Min() Variant {\n\treturn basicEnum.Min()\n}\n",
        "func (v Variant) Max() Variant {\n\treturn basicEnum.Max()\n}\n",
        "func (v Variant) IsMin() bool {\n\treturn basicEnum.IsMin(v)\n}\n",
        "func (v Variant) IsMax() bool {\n\treturn basicEnum.IsMax(v)\n}\n",
        "func (v Variant) IsInRange(min, max Variant) bool {\n\treturn baseenumer.IsBetween(v, min, max)\n}\n",
    ]


build_variant_body = build_variant_boundary_methods


def build_item_checker(item: str) -> str:
    return f"func (v Variant) Is{item}() bool {{\n\treturn v == {item}\n}}\n"


def build_variant_item_checkers(ctx: dict[str, Any]) -> list[str]:
    lines = []
    if ctx["zero_value"] != "Invalid":
        lines.append(build_item_checker(ctx["zero_value"]))
    for item in ctx["items"]:
        lines.append(build_item_checker(item))
    return lines


def build_string_formatting() -> list[str]:
    return [
        "func (v Variant) Name() string {\n\treturn string(v)\n}\n",
        "func (v Variant) Label() string {\n\treturn v.Name()\n}\n",
        "func (v Variant) String() string {\n\treturn v.Name()\n}\n",
        "func (v Variant) ValueString() string {\n\treturn string(v)\n}\n",
    ]


def build_numeric_formatting(ctx: dict[str, Any]) -> list[str]:
    u, name = ctx["underlying"], ctx["name"]
    return [
        f"func (v Variant) Name() string {{\n\tif int(v) < len(variantLabels) {{\n\t\treturn variantLabels[v]\n\t}}\n\n\treturn baseenumer.FormatNameValue(\"{name}\", {u}(v))\n}}\n",
        "func (v Variant) Label() string {\n\treturn v.Name()\n}\n",
        "func (v Variant) String() string {\n\treturn v.Name()\n}\n",
        f"func (v Variant) ValueString() string {{\n\treturn baseenumer.FormatNameValue(v.Name(), {u}(v))\n}}\n",
    ]


def build_variant_formatting(ctx: dict[str, Any]) -> list[str]:
    if ctx["is_string"]:
        return build_string_formatting()
    return build_numeric_formatting(ctx)


def build_string_json() -> list[str]:
    return [
        "func (v Variant) MarshalJSON() ([]byte, error) {\n\treturn baseenumer.MarshalJSON(string(v))\n}\n",
        "func (v *Variant) UnmarshalJSON(data []byte) error {\n\treturn basicEnum.UnmarshalJSON(data, v)\n}\n",
    ]


def build_numeric_json() -> list[str]:
    return [
        "func (v Variant) MarshalJSON() ([]byte, error) {\n\treturn baseenumer.MarshalJSON(v.Name())\n}\n",
        "func (v *Variant) UnmarshalJSON(data []byte) error {\n\treturn basicEnum.UnmarshalJSON(data, v)\n}\n",
    ]


def build_variant_json(ctx: dict[str, Any]) -> list[str]:
    if ctx["is_string"]:
        return build_string_json()
    return build_numeric_json()


def build_variant_type_section(ctx: dict[str, Any]) -> list[str]:
    return [
        f"type (\n\tVariant {ctx['underlying']}\n\n\t{ctx['name']}Type = Variant\n\n\tVariantPredicate func(v Variant) bool\n)\n"
    ]


def generate_variant_go(ctx: dict[str, Any]) -> str:
    parts = [
        build_variant_header(ctx), build_variant_type_section(ctx),
        build_variant_consts(ctx), build_variant_assertions(ctx),
        build_variant_byte_methods(ctx), build_variant_int_methods(ctx),
        build_variant_collection_methods(),
        build_variant_predicates(ctx), build_variant_boundary_methods(),
        build_variant_item_checkers(ctx), build_variant_formatting(ctx),
        build_variant_json(ctx),
    ]
    return "\n".join(line for chunk in parts for line in chunk)


def build_vars_header(ctx: dict[str, Any]) -> list[str]:
    imports = [
        '\t"coding-guidelines/common/pkg/baseenumer"',
    ]
    return [f"package {ctx['package']}", "", "import ("] + imports + [")", ""]


def build_string_registry(items: list[str]) -> list[str]:
    lines = ["\tvariantRegistry = map[Variant]bool{"]
    lines.extend(f"\t\t{item}: true," for item in items)
    return lines + ["\t}", ""]


def build_string_vars_data(ctx: dict[str, Any]) -> list[str]:
    items_lines = [f"\t\t{item}," for item in ctx["items"]]
    lines = ["var (", "\tallVariants = []Variant{"]
    lines.extend(items_lines)
    lines.extend(["\t}", ""] + build_string_registry(ctx["items"]))
    lines.extend([
        f"\tbasicEnum = baseenumer.NewBasicString(allVariants, {ctx['zero_value']})",
        ")",
        "",
    ])
    return lines


def build_numeric_vars_data(ctx: dict[str, Any]) -> list[str]:
    lines = ["var (", "\tvariantLabels = [...]string{", f'\t\t{ctx["zero_value"]}: "{ctx["zero_value"]}",']
    lines.extend(f'\t\t{item}: "{item}",' for item in ctx["items"])
    lines.extend([
        "\t}",
        "",
        f"\tbasicEnum = baseenumer.NewBasicInteger(variantLabels[:], {ctx['zero_value']})",
        ")",
        "",
    ])
    return lines


def build_vars_data(ctx: dict[str, Any]) -> list[str]:
    if ctx["is_string"]:
        return build_string_vars_data(ctx)
    return build_numeric_vars_data(ctx)


def build_string_vars_helpers() -> list[str]:
    return [
        "func All() []Variant {\n\treturn basicEnum.All()\n}\n",
        "func Values() []string {\n\treturn basicEnum.Values()\n}\n",
        "func Min() Variant {\n\treturn basicEnum.Min()\n}\n",
        "func Max() Variant {\n\treturn basicEnum.Max()\n}\n",
    ]


def build_numeric_vars_helpers() -> list[str]:
    return [
        "func All() []Variant {\n\treturn basicEnum.All()\n}\n",
        "func Values() []string {\n\treturn basicEnum.Values()\n}\n",
        "func Min() Variant {\n\treturn basicEnum.Min()\n}\n",
        "func Max() Variant {\n\treturn basicEnum.Max()\n}\n",
    ]


def build_vars_helpers(ctx: dict[str, Any]) -> list[str]:
    if ctx["is_string"]:
        return build_string_vars_helpers()
    return build_numeric_vars_helpers()


build_vars_body = build_vars_helpers


def build_vars_parser(ctx: dict[str, Any]) -> list[str]:
    zero = ctx["zero_value"]
    lines = [
        "func Parse(s string) (Variant, bool) {\n\treturn basicEnum.Parse(s)\n}\n",
        "func ParseOrZero(s string) Variant {\n\treturn basicEnum.ParseOrZero(s)\n}\n",
        "func ParseOrInvalid(s string) Variant {\n\treturn basicEnum.ParseOrZero(s)\n}\n",
    ]
    if zero != "Invalid":
        lines.append(f"func ParseOr{zero}(s string) Variant {{\n\treturn basicEnum.ParseOrZero(s)\n}}\n")
    return lines


def generate_vars_go(ctx: dict[str, Any]) -> str:
    lines = build_vars_header(ctx)
    lines.extend(build_vars_data(ctx))
    lines.extend(build_vars_helpers(ctx))
    lines.extend(build_vars_parser(ctx))
    return "\n".join(lines)


def build_test_header(ctx: dict[str, Any]) -> list[str]:
    return [
        f"package {ctx['package']}_test",
        "",
        'import (',
        '\t"encoding/json"',
        '\t"testing"',
        "",
        '\t"coding-guidelines/common/pkg/baseenumer"',
        f'\t"coding-guidelines/common/pkg/enum/{ctx["package"]}"',
        ")",
        "",
    ]


def resolve_test_interface_lines(ctx: dict[str, Any], first: str) -> list[str]:
    pkg = ctx["package"]
    if ctx["is_byte"]:
        return [f"\tvar _ baseenumer.ByteEnumer = {pkg}.{first}", f"\tvar _ baseenumer.NumberEnumer = {pkg}.{first}"]
    if ctx["is_string"]:
        return [f"\tvar _ baseenumer.StringEnumer = {pkg}.{first}"]
    return [f"\tvar _ baseenumer.NumberEnumer = {pkg}.{first}"]


def build_test_interfaces(ctx: dict[str, Any]) -> list[str]:
    pkg, first = ctx["package"], ctx["items"][0]
    lines = [f"func Test{ctx['name']}Type_Interfaces(t *testing.T) {{", f"\tvar _ baseenumer.BaseEnumer = {pkg}.{first}"]
    lines.extend(resolve_test_interface_lines(ctx, first))
    lines.extend([
        f"\tvar _ baseenumer.MinMaxer[{pkg}.Variant] = {pkg}.{first}",
        f"\tvar _ baseenumer.BoundedEnumer[{pkg}.Variant] = {pkg}.{first}",
        f"\tvar _ json.Marshaler = {pkg}.{first}",
        f"\tvar _ json.Unmarshaler = (*{pkg}.Variant)(nil)",
        "}",
        "",
    ])
    return lines


def build_numeric_test_properties(pkg: str, first: str, is_byte: bool) -> list[str]:
    lines = []
    if is_byte:
        lines.append(f'\tif {pkg}.{first}.Byte() != 1 || {pkg}.{first}.Value() != 1 {{\n\t\tt.Fatalf("expected 1 from Byte/Value()")\n\t}}')
    else:
        lines.append(f'\tif {pkg}.{first}.Value() != 1 {{\n\t\tt.Fatalf("expected 1 from Value()")\n\t}}')
    lines.append(f'\tif {pkg}.{first}.Int() != 1 || {pkg}.{first}.Code() != 1 {{\n\t\tt.Fatalf("expected 1 from Int/Code")\n\t}}')
    return lines


def build_test_properties(ctx: dict[str, Any]) -> list[str]:
    pkg, first = ctx["package"], ctx["items"][0]
    lines = [f"func Test{ctx['name']}Type_Properties(t *testing.T) {{"]
    if not ctx["is_string"]:
        lines.extend(build_numeric_test_properties(pkg, first, ctx["is_byte"]))
    else:
        lines.append(f'\tif {pkg}.{first}.Value() != "{first}" {{\n\t\tt.Fatalf("expected {first} from Value()")\n\t}}')
    lines.extend([f"\tif !{pkg}.{first}.IsValid() {{\n\t\tt.Fatalf(\"expected {first} to be valid\")\n\t}}", "}", ""])
    return lines


def build_test_boundaries(ctx: dict[str, Any]) -> list[str]:
    p, n = ctx["package"], ctx["name"]
    chk_rec = 'if minVal.Min() != minVal || maxVal.Max() != maxVal {\n\t\tt.Fatalf("receiver Min/Max mismatch")\n\t}'
    chk_pred = 'if !minVal.IsMin() || !maxVal.IsMax() {\n\t\tt.Fatalf("boundary predicates failed")\n\t}'
    chk_rng = 'if !minVal.IsInRange(minVal, maxVal) {\n\t\tt.Fatalf("IsInRange failed")\n\t}'
    body = f"func Test{n}Type_Boundaries(t *testing.T) {{\n\tminVal, maxVal := {p}.Min(), {p}.Max()\n\t{chk_rec}\n\t{chk_pred}\n\t{chk_rng}\n}}"
    return [body, ""]


def build_test_predicates(ctx: dict[str, Any]) -> list[str]:
    pkg, first, zero = ctx["package"], ctx["items"][0], ctx["zero_value"]
    return [
        f"func Test{ctx['name']}Type_Predicates(t *testing.T) {{",
        f'\tif !{pkg}.{first}.Is{first}() {{\n\t\tt.Fatalf("expected Is{first}() to be true")\n\t}}',
        f'\tif {pkg}.{first}.Is{zero}() {{\n\t\tt.Fatalf("expected Is{zero}() to be false for {first}")\n\t}}',
        "}",
        "",
    ]


def build_test_names_and_vars(ctx: dict[str, Any]) -> list[str]:
    pkg, first, zero, n = ctx["package"], ctx["items"][0], ctx["zero_value"], len(ctx["items"])
    return [
        f"func Test{ctx['name']}Type_VarsAndParse(t *testing.T) {{",
        f'\tall := {pkg}.All()\n\tif len(all) != {n} || len({pkg}.{first}.All()) != {n} {{\n\t\tt.Fatalf("expected {n} variants, got %d", len(all))\n\t}}',
        f'\tvals := {pkg}.Values()\n\tif len(vals) != {n} || len({pkg}.{first}.Values()) != {n} {{\n\t\tt.Fatalf("expected {n} values, got %d", len(vals))\n\t}}',
        f'\tv, isOk := {pkg}.Parse("{first}")\n\tif !isOk || v != {pkg}.{first} {{\n\t\tt.Fatalf("expected successful Parse for {first}")\n\t}}',
        f'\tif _, isOk := {pkg}.Parse(""); isOk {{\n\t\tt.Fatalf("expected failure for empty string")\n\t}}',
        f'\tif _, isOk := {pkg}.Parse("invalid_variant_value"); isOk {{\n\t\tt.Fatalf("expected failure for bad variant")\n\t}}',
        f'\tif {pkg}.ParseOrZero("{first}") != {pkg}.{first} {{\n\t\tt.Fatalf("ParseOrZero failed")\n\t}}',
        f'\tif {pkg}.ParseOrInvalid("invalid_variant_value") != {pkg}.{zero} {{\n\t\tt.Fatalf("ParseOrInvalid fallback failed")\n\t}}',
        "}",
        "",
    ]


def build_test_json(ctx: dict[str, Any]) -> list[str]:
    pkg, first, zero = ctx["package"], ctx["items"][0], ctx["zero_value"]
    return [
        f"func Test{ctx['name']}Type_JSON(t *testing.T) {{",
        f'\tdata, err := json.Marshal({pkg}.{first})\n\tif err != nil || string(data) != `"{first}"` {{\n\t\tt.Fatalf("marshal failed: %v", err)\n\t}}',
        f'\tvar v {pkg}.Variant\n\tif err := json.Unmarshal([]byte(`"{first}"`), &v); err != nil || v != {pkg}.{first} {{\n\t\tt.Fatalf("unmarshal failed: %v", err)\n\t}}',
        f'\tif err := json.Unmarshal([]byte(`null`), &v); err != nil || v != {pkg}.{zero} {{\n\t\tt.Fatalf("unmarshal null failed: %v", err)\n\t}}',
        "}",
        "",
    ]


def generate_variant_test_go(ctx: dict[str, Any]) -> str:
    lines = build_test_header(ctx)
    lines.extend(build_test_interfaces(ctx))
    lines.extend(build_test_properties(ctx))
    lines.extend(build_test_boundaries(ctx))
    lines.extend(build_test_predicates(ctx))
    lines.extend(build_test_names_and_vars(ctx))
    lines.extend(build_test_json(ctx))
    return "\n".join(lines)


def generate_readme_md(ctx: dict[str, Any]) -> str:
    members = ", ".join(f"`{item}`" for item in ctx["items"])
    p, u, n = ctx["package"], ctx["underlying"], ctx["name"]
    return (
        f"# `{p}` Package\n\n"
        f"`{p}` provides a type-safe {u}-backed enumeration for {n} variants ({members}).\n\n"
        "## Key Features\n\n"
        "- **Zero Circular Dependencies:** Foundational leaf enum importing only `baseenumer` and standard libraries.\n"
        "- **Zero-Allocation Parsing:** `Parse(s string) (Variant, bool)` directly delegates to `baseenumer.BasicInteger`.\n"
        "- **Safe Fallback Parsers:** `ParseOrZero(s)`, `ParseOrInvalid(s)`, and `ParseOrUnknown(s)` helpers.\n"
        "- **DRY JSON Marshaling:** Implements `json.Marshaler` and `json.Unmarshaler`.\n"
        "- **Boundary Operations:** First-class `Min()`, `Max()`, `IsMin()`, `IsMax()`, and `IsInRange(min, max Variant) bool` conforming to `baseenumer.BoundedEnumer[Variant]`.\n"
    )


def render_bundle(ctx: dict[str, Any]) -> dict[str, str]:
    return {
        "variant.go": generate_variant_go(ctx),
        "vars.go": generate_vars_go(ctx),
        "variant_test.go": generate_variant_test_go(ctx),
        "readme.md": generate_readme_md(ctx),
    }


def execute_dry_run(target_dir: Path, bundle: dict[str, str]) -> int:
    norm_dir = target_dir.as_posix()
    print(f"[DRY-RUN] Target directory: {norm_dir}")
    for filename, content in bundle.items():
        rel_file = f"{norm_dir}/{filename}"
        print(f"[DRY-RUN] Would create {rel_file} ({len(content.splitlines())} lines)")
    return EXIT_SUCCESS


def check_overwrite_conflict(dest: Path, is_overwrite: bool) -> bool:
    if is_overwrite:
        return False
    return dest.exists()


def run_code_formatter(target_dir: Path) -> None:
    formatter = Path(__file__).parent / "26-go-code-formatter.py"
    if not formatter.exists():
        return
    try:
        import subprocess
        subprocess.run(
            [sys.executable, str(formatter), str(target_dir)],
            capture_output=True,
            text=True,
            check=False,
        )
    except Exception:
        pass


def write_bundle(target_dir: Path, bundle: dict[str, str], is_overwrite: bool) -> int:
    target_dir.mkdir(parents=True, exist_ok=True)
    norm_dir = target_dir.as_posix()
    for filename, content in bundle.items():
        dest = target_dir / filename
        has_conflict = check_overwrite_conflict(dest, is_overwrite)
        if has_conflict:
            print(f"Error: {norm_dir}/{filename} already exists. Use --overwrite.", file=sys.stderr)
            return EXIT_USAGE_ERROR
        write_file_lf(dest, content)
        print(f"Created: {norm_dir}/{filename}")
    run_code_formatter(target_dir)
    print(f"✔ Successfully generated {len(bundle)} files in {norm_dir}")
    return EXIT_SUCCESS



def main() -> int:
    args = parse_cli_args()
    try:
        raw_cfg = load_raw_config(args)
        cfg = normalize_config(raw_cfg)
    except Exception as exc:
        print(f"Configuration error: {exc}", file=sys.stderr)
        return EXIT_USAGE_ERROR

    target_dir = Path(cfg["target_dir"])
    bundle = render_bundle(cfg)
    if args.dry_run:
        return execute_dry_run(target_dir, bundle)
    return write_bundle(target_dir, bundle, args.overwrite)


if __name__ == "__main__":
    sys.exit(main())
