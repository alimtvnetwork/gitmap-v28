// Package templates ships curated .gitignore and .gitattributes content
// per language, with a user-overlay system for read-only install paths.
//
// Resolution order for any (kind, lang):
//  1. <UserTemplatesDir>/<kind>/<lang><ext>   (user override)
//  2. embedded assets/<kind>/<lang><ext>      (built-in fallback)
//
// Phase 0: package + embed + resolver + materializer only. CLI wiring
// arrives in Phase 2.
package templates

import "embed"

// FS holds the curated template corpus. Populated in Phase 1.
//
//go:embed all:assets
var FS embed.FS
