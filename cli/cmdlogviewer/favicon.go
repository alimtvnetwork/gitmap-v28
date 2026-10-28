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
      <stop offset="0%" stop-color="#1e1e2e"/>
      <stop offset="100%" stop-color="#11111b"/>
    </linearGradient>
    <filter id="gmShadow" x="-20%" y="-20%" width="140%" height="140%">
      <feDropShadow dx="0" dy="1" stdDeviation="1.5" flood-color="#000000" flood-opacity="0.4"/>
    </filter>
  </defs>
  <rect x="1" y="1" width="30" height="30" rx="7" fill="url(#gmBg)" stroke="#89b4fa" stroke-width="1.6" filter="url(#gmShadow)"/>
  <rect x="2.5" y="2.5" width="27" height="27" rx="5.5" fill="none" stroke="#313244" stroke-width="0.8" opacity="0.6"/>
  <path d="M11 7.5v17M11 14c3.5 0 10 2.5 10 9.5" fill="none" stroke="url(#gmGrad)" stroke-width="2.2" stroke-linecap="round"/>
  <circle cx="11" cy="7.5" r="3" fill="#89b4fa" stroke="#cdd6f4" stroke-width="0.8"/>
  <circle cx="11" cy="7.5" r="1.2" fill="#1e1e2e"/>
  <circle cx="11" cy="24.5" r="3" fill="#89b4fa" stroke="#cdd6f4" stroke-width="0.8"/>
  <circle cx="11" cy="24.5" r="1.2" fill="#1e1e2e"/>
  <circle cx="21" cy="23.5" r="3.2" fill="#a6e3a1" stroke="#ffffff" stroke-width="0.8"/>
  <circle cx="21" cy="23.5" r="1.3" fill="#1e1e2e"/>
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
