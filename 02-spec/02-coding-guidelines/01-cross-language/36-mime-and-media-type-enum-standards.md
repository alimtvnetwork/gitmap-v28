# MIME & Media Type Enum Standards — Coding Guidelines

## 1. Overview & Core Mission
In web, mobile, CLI, desktop, and backend applications interacting with network APIs, file systems, clipboard copy/paste, HTTP headers, Blob constructors, file downloads, content negotiation, and multipart payloads, developers frequently pass raw magic string literals (e.g. `'text/html'`, `'text/plain'`, `'application/json'`, `'image/png'`).

This specification establishes the **MIME & Media Type Enum Standard** as a non-negotiable coding guideline across all supported languages (TypeScript, Go, Python, Rust, PHP). Inline magic string literals for MIME types, media types, and Content-Type headers are strictly banned. All MIME and media types MUST be defined and referenced via strongly-typed, centralized enums or typed constant dictionaries (`MimeType` / `MediaType`).

---

## 2. Anti-Pattern: Magic String MIME Types

### ❌ Anti-Pattern: Browser Blob & Clipboard API Magic Strings
```typescript
// ❌ ANTI-PATTERN: Magic string literals passed to Blob and ClipboardItem constructors
const htmlBlob = new Blob([htmlContent], { type: 'text/html' });
const textBlob = new Blob([mdText], { type: 'text/plain' });

const clipboardItem = new ClipboardItem({
  'text/html': htmlBlob,
  'text/plain': textBlob,
});
await navigator.clipboard.write([clipboardItem]);
```

### ❌ Anti-Pattern: HTTP Headers & API Response Magic Strings
```typescript
// ❌ ANTI-PATTERN: Raw strings in response headers and content-type setters
res.setHeader('Content-Type', 'application/json');
fetch('/api/export', {
  headers: {
    'Accept': 'text/csv',
    'Content-Type': 'application/json',
  },
});
```

```go
// ❌ ANTI-PATTERN: Raw strings in Go HTTP response headers
w.Header().Set("Content-Type", "application/json")
w.Header().Set("Content-Type", "text/html; charset=utf-8")
```

```python
# ❌ ANTI-PATTERN: Raw strings in Python requests or FastAPI responses
response = requests.post(url, headers={"Content-Type": "application/json"})
return Response(content=data, media_type="text/plain")
```

### ❌ Why Magic String MIME Types Fail:
1. **Silent Typo Vulnerability:** Typographical errors like `'text/htm'`, `'text/plian'`, or `'applicaton/json'` compile silently and fail only at runtime or cause browsers/clients to misinterpret payloads.
2. **Scattered Magic Literals:** The same MIME type string is repeated across tens or hundreds of files without a single source of truth.
3. **No Autocompletion or Discovery:** Developers must memorize exact RFC-compliant MIME type strings rather than selecting from compiler-checked enum members.
4. **Refactoring Fragility:** Updating headers or supporting new content encodings requires dangerous global search-and-replace across raw strings.
5. **Violation of Coding Standards:** Violates cross-language rules against magic string literals (`06-constants-and-enums.md` and `26-magic-values-and-immutability.md`).

---

## 3. The Canonical Pattern: Strongly-Typed MimeType Enums

### ✅ TypeScript: Const Object Enum Pattern
```typescript
// ✅ CANONICAL PATTERN: Centralized MimeType enum in src/constants/mimeType.ts
export const MimeType = {
  TextPlain: 'text/plain',
  TextHtml: 'text/html',
  TextMarkdown: 'text/markdown',
  TextCss: 'text/css',
  TextCsv: 'text/csv',
  ApplicationJson: 'application/json',
  ApplicationXml: 'application/xml',
  ApplicationPdf: 'application/pdf',
  ApplicationZip: 'application/zip',
  ApplicationOctetStream: 'application/octet-stream',
  ImagePng: 'image/png',
  ImageJpeg: 'image/jpeg',
  ImageGif: 'image/gif',
  ImageSvg: 'image/svg+xml',
  ImageWebp: 'image/webp',
  AudioMpeg: 'audio/mpeg',
  VideoMp4: 'video/mp4',
  MultipartFormData: 'multipart/form-data',
} as const;

export type MimeType = typeof MimeType[keyof typeof MimeType];
```

#### ✅ TypeScript Usage: Blob & Clipboard APIs
```typescript
// ✅ CANONICAL USAGE: Compiler-verified enum references
import { MimeType } from '@/constants/mimeType';

const htmlBlob = new Blob([htmlContent], { type: MimeType.TextHtml });
const textBlob = new Blob([mdText], { type: MimeType.TextPlain });

const clipboardItem = new ClipboardItem({
  [MimeType.TextHtml]: htmlBlob,
  [MimeType.TextPlain]: textBlob,
});
await navigator.clipboard.write([clipboardItem]);
```

#### ✅ TypeScript Usage: HTTP Headers & Fetch
```typescript
import { MimeType } from '@/constants/mimeType';

const response = await fetch('/api/data', {
  method: 'POST',
  headers: {
    'Accept': MimeType.ApplicationJson,
    'Content-Type': MimeType.ApplicationJson,
  },
  body: JSON.stringify(payload),
});
```

---

### ✅ Go: String-Backed MimeType Enum & Constants
```go
// ✅ CANONICAL PATTERN: Go MimeType enum in pkg/enum/mimetype/mimetype.go
package mimetype

// MimeType represents an RFC-compliant media type string.
type MimeType string

const (
    TextPlain              MimeType = "text/plain"
    TextHTML               MimeType = "text/html"
    TextMarkdown           MimeType = "text/markdown"
    TextCSS                MimeType = "text/css"
    TextCSV                MimeType = "text/csv"
    ApplicationJSON        MimeType = "application/json"
    ApplicationXML         MimeType = "application/xml"
    ApplicationPDF         MimeType = "application/pdf"
    ApplicationZip         MimeType = "application/zip"
    ApplicationOctetStream MimeType = "application/octet-stream"
    ImagePNG               MimeType = "image/png"
    ImageJPEG              MimeType = "image/jpeg"
    ImageGIF               MimeType = "image/gif"
    ImageSVG               MimeType = "image/svg+xml"
    ImageWebP              MimeType = "image/webp"
    MultipartFormData      MimeType = "multipart/form-data"
)

// String returns the underlying string representation.
func (m MimeType) String() string {
    return string(m)
}
```

#### ✅ Go Usage: HTTP Handlers
```go
package handler

import (
    "net/http"
    "pkg/enum/mimetype"
)

func ServeJSON(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", mimetype.ApplicationJSON.String())
    w.WriteHeader(http.StatusOK)
}
```

---

### ✅ Python: StrEnum / Enum Standard
```python
# ✅ CANONICAL PATTERN: Python MimeType enum in app/enums/mime_type.py
from enum import Enum

class MimeType(str, Enum):
    TEXT_PLAIN = "text/plain"
    TEXT_HTML = "text/html"
    TEXT_MARKDOWN = "text/markdown"
    TEXT_CSS = "text/css"
    TEXT_CSV = "text/csv"
    APPLICATION_JSON = "application/json"
    APPLICATION_XML = "application/xml"
    APPLICATION_PDF = "application/pdf"
    APPLICATION_OCTET_STREAM = "application/octet-stream"
    IMAGE_PNG = "image/png"
    IMAGE_JPEG = "image/jpeg"
    IMAGE_GIF = "image/gif"
    IMAGE_SVG = "image/svg+xml"
    IMAGE_WEBP = "image/webp"
```

#### ✅ Python Usage: Requests & FastAPI
```python
import requests
from app.enums.mime_type import MimeType

response = requests.post(
    url,
    headers={"Content-Type": MimeType.APPLICATION_JSON.value},
    json=payload,
)
```

---

### ✅ Rust: Strongly-Typed Enum & Consts
```rust
// ✅ CANONICAL PATTERN: Rust MimeType enum in src/enums/mime_type.rs
#[derive(Debug, Clone, Copy, PartialEq, Eq, Hash)]
pub enum MimeType {
    TextPlain,
    TextHtml,
    TextMarkdown,
    ApplicationJson,
    ApplicationPdf,
    ApplicationOctetStream,
    ImagePng,
    ImageJpeg,
    ImageSvg,
}

impl MimeType {
    pub const fn as_str(&self) -> &'static str {
        match self {
            Self::TextPlain => "text/plain",
            Self::TextHtml => "text/html",
            Self::TextMarkdown => "text/markdown",
            Self::ApplicationJson => "application/json",
            Self::ApplicationPdf => "application/pdf",
            Self::ApplicationOctetStream => "application/octet-stream",
            Self::ImagePng => "image/png",
            Self::ImageJpeg => "image/jpeg",
            Self::ImageSvg => "image/svg+xml",
        }
    }
}
```

---

### ✅ PHP: Backed Enum Standard (PHP 8.1+)
```php
<?php
// ✅ CANONICAL PATTERN: PHP 8.1+ Backed Enum in App\Enums\MimeType.php
namespace App\Enums;

enum MimeType: string
{
    case TextPlain = 'text/plain';
    case TextHtml = 'text/html';
    case TextMarkdown = 'text/markdown';
    case ApplicationJson = 'application/json';
    case ApplicationPdf = 'application/pdf';
    case ApplicationOctetStream = 'application/octet-stream';
    case ImagePng = 'image/png';
    case ImageJpeg = 'image/jpeg';
    case ImageSvg = 'image/svg+xml';
}
```

---

## 4. Extension-to-MIME Mapping Tables (Spec 35 Alignment)

When resolving file extensions to MIME types or vice versa, codebases MUST NOT use cascading `if ... else if` ladders or `switch` statements. Instead, adhere to **Spec 35 (Map Dispatch & Handler Tables)** by using declarative, immutable lookup tables:

```typescript
// ✅ CANONICAL PATTERN: Immutable lookup table mapping extensions to MimeType
import { MimeType } from '@/constants/mimeType';

const EXTENSION_MIME_MAP: Readonly<Record<string, MimeType>> = {
  '.html': MimeType.TextHtml,
  '.htm': MimeType.TextHtml,
  '.txt': MimeType.TextPlain,
  '.md': MimeType.TextMarkdown,
  '.json': MimeType.ApplicationJson,
  '.png': MimeType.ImagePng,
  '.jpg': MimeType.ImageJpeg,
  '.jpeg': MimeType.ImageJpeg,
  '.svg': MimeType.ImageSvg,
  '.pdf': MimeType.ApplicationPdf,
};

export function resolveMimeFromExtension(ext: string): MimeType {
  const normalized = ext.toLowerCase();
  return EXTENSION_MIME_MAP[normalized] ?? MimeType.ApplicationOctetStream;
}
```

---

## 5. Non-Negotiable Invariants & Rules

1. **TOTAL BAN on Magic Strings for MIME Types:** Never hardcode string literals like `'text/html'`, `'application/json'`, or `'image/png'` at call sites.
2. **Blob & Clipboard Mandatory Compliance:** All `new Blob(..., { type })` and `new ClipboardItem({ [mime]: blob })` calls MUST supply the type via `MimeType` enum members.
3. **HTTP Headers Compliance:** `Content-Type`, `Accept`, and response header setters MUST use `MimeType` enum values.
4. **PascalCase Enum Members:** Enum members follow standard PascalCase (`TextHtml`, `ApplicationJson`, `ImagePng`), eliminating inconsistent formatting.
5. **No Loose Helper Functions with Magic Strings:** Any MIME resolution helper must return or accept the strongly-typed `MimeType` enum.
