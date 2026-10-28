import type { ClassifiedCell } from "./types";

export function formatByteSize(bytes: number): string {
  if (bytes <= 0) {
    return "0 B";
  }

  if (bytes < 1024) {
    return `${bytes} B`;
  }

  if (bytes < 1024 * 1024) {
    return `${(bytes / 1024).toFixed(1)} KB`;
  }

  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
}

export function isBinaryString(str: string): boolean {
  if (!str || str.length === 0) {
    return false;
  }

  if (str.includes("\uFFFD") || str.includes("\x00")) {
    return true;
  }

  const checkLen = Math.min(str.length, 256);
  for (let i = 0; i < checkLen; i++) {
    const code = str.charCodeAt(i);
    if (code < 32 && code !== 9 && code !== 10 && code !== 13) {
      return true;
    }
  }

  return false;
}

export function classifyCell(val: unknown): ClassifiedCell {
  if (val === null || val === undefined) {
    return { cellType: "null", badgeText: "NULL", rawValue: val, sizeBytes: 0 };
  }

  if (typeof val === "boolean") {
    return { cellType: "bool", badgeText: val ? "[BOOL true]" : "[BOOL false]", rawValue: val, sizeBytes: 1 };
  }

  if (typeof val === "object") {
    return classifyObjectCell(val as object);
  }

  if (typeof val === "string") {
    return classifyStringCell(val);
  }

  return { cellType: "plain", badgeText: String(val), rawValue: val, sizeBytes: String(val).length };
}

function classifyObjectCell(val: object): ClassifiedCell {
  const isArr = Array.isArray(val);
  const count = isArr ? `${(val as unknown[]).length} items` : `${Object.keys(val).length} keys`;
  const jsonStr = JSON.stringify(val);

  return { cellType: "json", badgeText: `[JSON ${count}]`, rawValue: val, sizeBytes: jsonStr.length };
}

function classifyStringCell(val: string): ClassifiedCell {
  if (isBinaryString(val)) {
    return { cellType: "blob", badgeText: `[BLOB ${formatByteSize(val.length)}]`, rawValue: val, sizeBytes: val.length };
  }

  const jsonResult = tryParseJsonString(val);
  if (jsonResult) {
    return jsonResult;
  }

  if (val.includes("\n") || val.includes("\r")) {
    return classifyMultilineString(val);
  }

  if (val.length > 70) {
    return { cellType: "text", badgeText: `${val.slice(0, 68)}...`, rawValue: val, sizeBytes: val.length };
  }

  return { cellType: "plain", badgeText: val, rawValue: val, sizeBytes: val.length };
}

function tryParseJsonString(val: string): ClassifiedCell | null {
  const trimmed = val.trim();
  const isCandidate = (trimmed.startsWith("{") && trimmed.endsWith("}")) || (trimmed.startsWith("[") && trimmed.endsWith("]"));
  if (!isCandidate) {
    return null;
  }

  try {
    const parsed = JSON.parse(trimmed);
    if (typeof parsed === "object" && parsed !== null) {
      const isArr = Array.isArray(parsed);
      const count = isArr ? `${parsed.length} items` : `${Object.keys(parsed).length} keys`;

      return { cellType: "json", badgeText: `[JSON ${count}]`, rawValue: parsed, sizeBytes: val.length };
    }
  } catch {
    return null;
  }

  return null;
}

function classifyMultilineString(val: string): ClassifiedCell {
  const lines = val.split(/\r?\n/);
  const firstLine = lines[0].slice(0, 32);

  return { cellType: "multiline", badgeText: `${firstLine} ↵ [${lines.length} lines]`, rawValue: val, sizeBytes: val.length };
}

export function generateHexDump(raw: unknown): string {
  const str = typeof raw === "string" ? raw : JSON.stringify(raw);
  if (!str) {
    return "(empty)";
  }

  const bytes = extractStringBytes(str);

  return buildHexDumpLines(bytes).join("\n");
}

function extractStringBytes(str: string): number[] {
  const bytes: number[] = [];
  for (let i = 0; i < str.length; i++) {
    bytes.push(str.charCodeAt(i) & 0xff);
  }

  return bytes;
}

function buildHexDumpLines(bytes: number[]): string[] {
  const lines: string[] = [];
  for (let offset = 0; offset < bytes.length; offset += 16) {
    const chunk = bytes.slice(offset, offset + 16);
    lines.push(formatHexLine(offset, chunk));
  }

  return lines;
}

function formatHexLine(offset: number, chunk: number[]): string {
  const hexParts = chunk.map((b) => b.toString(16).padStart(2, "0").toUpperCase());
  while (hexParts.length < 16) {
    hexParts.push("  ");
  }

  const ascii = chunk.map((b) => (b >= 32 && b <= 126 ? String.fromCharCode(b) : ".")).join("");
  const left = hexParts.slice(0, 8).join(" ");
  const right = hexParts.slice(8).join(" ");
  const offStr = offset.toString(16).padStart(8, "0").toUpperCase();

  return `${offStr}  ${left}  ${right}  |${ascii}|`;
}
