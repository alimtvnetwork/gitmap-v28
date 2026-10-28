package cmdlogviewer

import (
	_ "embed"
	"net/http"
)

//go:embed favicon.ico
var FaviconICO []byte

// GitMapLogoSVG is the authentic GitMap brand emblem logo designed with a 100% transparent canvas and crisp vector geometry.
const GitMapLogoSVG = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 32 32" width="24" height="24" fill="none">
  <defs>
    <linearGradient id="gmBranch" x1="0%" y1="0%" x2="100%" y2="100%">
      <stop offset="0%" stop-color="#89b4fa"/>
      <stop offset="100%" stop-color="#b4befe"/>
    </linearGradient>
    <linearGradient id="gmGreen" x1="0%" y1="0%" x2="100%" y2="100%">
      <stop offset="0%" stop-color="#a6e3a1"/>
      <stop offset="100%" stop-color="#94e2d5"/>
    </linearGradient>
  </defs>
  <!-- Main trunk -->
  <line x1="9" y1="7" x2="9" y2="25" stroke="url(#gmBranch)" stroke-width="2.6" stroke-linecap="round"/>
  <!-- Branch fork -->
  <path d="M9 14.5 C9 18.5, 14 20, 22 20 L22 23" fill="none" stroke="url(#gmBranch)" stroke-width="2.6" stroke-linecap="round"/>
  <!-- Root node -->
  <circle cx="9" cy="7" r="3.4" fill="#89b4fa" stroke="#181825" stroke-width="1.2"/>
  <circle cx="9" cy="7" r="1.4" fill="#cdd6f4"/>
  <!-- Base node -->
  <circle cx="9" cy="25" r="3.4" fill="#89b4fa" stroke="#181825" stroke-width="1.2"/>
  <circle cx="9" cy="25" r="1.4" fill="#cdd6f4"/>
  <!-- Merged leaf node -->
  <circle cx="22" cy="23.5" r="3.8" fill="url(#gmGreen)" stroke="#181825" stroke-width="1.2"/>
  <circle cx="22" cy="23.5" r="1.6" fill="#ffffff"/>
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
