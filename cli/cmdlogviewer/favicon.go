package cmdlogviewer

import (
	_ "embed"
	"net/http"
)

//go:embed favicon.ico
var FaviconICO []byte

// GitMapLogoSVG is the authentic GitMap brand emblem logo.
const GitMapLogoSVG = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 32 32" width="22" height="22">
  <defs>
    <linearGradient id="gmGrad" x1="0%" y1="0%" x2="100%" y2="100%">
      <stop offset="0%" stop-color="#89b4fa"/>
      <stop offset="100%" stop-color="#b4befe"/>
    </linearGradient>
    <linearGradient id="gmBg" x1="0%" y1="0%" x2="100%" y2="100%">
      <stop offset="0%" stop-color="#313244"/>
      <stop offset="100%" stop-color="#181825"/>
    </linearGradient>
  </defs>
  <rect width="32" height="32" rx="8" fill="url(#gmBg)" stroke="#89b4fa" stroke-width="1.5"/>
  <path d="M10 8v16M10 14c3 0 7 2 7 6v4" fill="none" stroke="url(#gmGrad)" stroke-width="2.5" stroke-linecap="round"/>
  <circle cx="10" cy="9" r="2.8" fill="#89b4fa"/>
  <circle cx="10" cy="23" r="2.8" fill="#89b4fa"/>
  <circle cx="17" cy="23" r="2.8" fill="#a6e3a1"/>
</svg>`

func serveFaviconICO(w http.ResponseWriter, r *http.Request) {
	if len(FaviconICO) > 0 {
		w.Header().Set("Content-Type", "image/x-icon")
		w.Header().Set("Cache-Control", "public, max-age=86400")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(FaviconICO)
		return
	}

	serveFaviconSVG(w, r)
}

func serveFaviconSVG(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/svg+xml")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(GitMapLogoSVG))
}

// WebAppManifestJSON provides PWA application metadata for Chrome/Edge app windows.
const WebAppManifestJSON = `{
  "name": "GitMap — Show Log",
  "short_name": "Show Log",
  "start_url": "/",
  "display": "standalone",
  "background_color": "#181825",
  "theme_color": "#181825",
  "icons": [
    {
      "src": "/favicon.ico",
      "sizes": "16x16 24x24 32x32 48x48 64x64",
      "type": "image/x-icon"
    },
    {
      "src": "/favicon.svg",
      "sizes": "any",
      "type": "image/svg+xml"
    }
  ]
}`

func serveManifestJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/manifest+json")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(WebAppManifestJSON))
}
