"""Regenerate the TypeScript UI-remote registry from the Python source.

Single source of truth: `BE/ui_remote/registry.py`. Run from the repo
root:

    python3 03-ai-scripts/generate-ui-remote-registry.py

Writes `src/lib/ui-remote/registry.ts`. Never edit the mirror by hand.
"""

from __future__ import annotations

import json
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
sys.path.insert(0, str(ROOT))

from BE.ui_remote.registry import ELEMENTS  # noqa: E402


def _ts_element(entry: dict) -> str:
    """Render one registry entry as a TypeScript object literal."""

    set_range = "null"

    if entry["set_range"] is not None:
        set_range = f"[{entry['set_range'][0]}, {entry['set_range'][1]}]"

    lines = [
        "  {",
        f"    id: {json.dumps(entry['id'])},",
        f"    name: {json.dumps(entry['name'])},",
        f"    kind: {json.dumps(entry['kind'])},",
        f"    xpath: {json.dumps(entry['xpath'])},",
        f"    routes: {json.dumps(entry['routes'])},",
        f"    actions: {json.dumps(entry['actions'])},",
        f"    setRange: {set_range},",
        "  },",
    ]

    return "\n".join(lines)


def main() -> int:
    """Emit the TypeScript mirror and report the entry count."""

    header = (
        "/** Generated from BE/ui_remote/registry.py — do not edit by hand. */\n"
        "\n"
        "export type UiRemoteAction = \"click\" | \"scroll\" | \"focus\" | \"set\" | \"navigate\";\n"
        "\n"
        "export interface UiRemoteElement {\n"
        "  id: string;\n"
        "  name: string;\n"
        "  kind: \"button\" | \"input\" | \"select\" | \"region\";\n"
        "  xpath: string;\n"
        "  routes: string[];\n"
        "  actions: UiRemoteAction[];\n"
        "  setRange: [number, number] | null;\n"
        "}\n"
        "\n"
        "export interface UiRemotePacket {\n"
        "  packetId: string;\n"
        "  action: UiRemoteAction;\n"
        "  elementId?: string | null;\n"
        "  xpath?: string | null;\n"
        "  route?: string | null;\n"
        "  value?: string | null;\n"
        "  actor: string;\n"
        "  status: string;\n"
        "  createdAt: string;\n"
        "}\n"
        "\n"
        "export interface UiRemoteResult {\n"
        "  packetId: string;\n"
        "  performed: boolean;\n"
        "  resolvedVia?: string | null;\n"
        "  error?: string | null;\n"
        "}\n"
        "\n"
        "export const UI_REMOTE_ELEMENTS: UiRemoteElement[] = [\n"
    )
    body = "\n".join(_ts_element(entry) for entry in ELEMENTS)
    footer = (
        "];\n"
        "\n"
        "export function findUiRemoteElementById(id: string): UiRemoteElement | undefined {\n"
        "\n"
        "  return UI_REMOTE_ELEMENTS.find((entry) => entry.id === id);\n"
        "}\n"
    )
    target = ROOT / "src" / "lib" / "ui-remote" / "registry.ts"
    target.parent.mkdir(parents=True, exist_ok=True)
    target.write_text(header + body + "\n" + footer, encoding="utf-8")
    print(f"Wrote {target.relative_to(ROOT)} ({len(ELEMENTS)} elements)")

    return 0


if __name__ == "__main__":
    raise SystemExit(main())
