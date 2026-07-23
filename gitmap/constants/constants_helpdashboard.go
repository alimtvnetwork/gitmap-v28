package constants

// Help-dashboard help text.
const HelpHelpDashboard = "  help-dashboard (hd) Serve the docs site locally in your browser"

// Help-dashboard flag descriptions.
const (
	FlagDescHDPort = "Port to serve the dashboard on (default: 5173)"
)

// Help-dashboard defaults.
const (
	HDDefaultPort   = 5173
	HDDistDir       = "dist"
	HDDocsDir       = "docs-site"
	DocsSiteArchive = "docs-site.zip"
)

// Help-dashboard terminal messages.
const (
	MsgHDServingStatic  = "  Serving docs from %s on http://localhost:%d\n"
	MsgHDStartingDev    = "  Starting dev server from %s...\n"
	MsgHDRunningNPM     = "  Running npm install...\n"
	MsgHDOpening        = "  Opening http://localhost:%d in browser...\n"
	MsgHDNoDistFallback = "  No pre-built dist/ found, falling back to npm run dev\n"
	MsgHDStopped        = "\n  Server stopped.\n"
)

// Help-dashboard error messages.
const (
	ErrHDNoDocsDir    = "  ✗ Docs site directory not found at %s (operation: resolve, reason: directory does not exist)\n"
	ErrHDNPMInstall   = "  ✗ npm install failed: %v\n"
	ErrHDDevServer    = "  ✗ Dev server failed: %v\n"
	ErrHDServe        = "  ✗ Failed to start server: %v\n"
	ErrHDNPMNotFound  = "  ✗ npm not found — install Node.js to use dev mode\n"
	ErrDocsSiteBundle = "  ✗ Failed to bundle docs-site: %v\n"
)

// Docs-site release messages.
const (
	MsgDocsSiteBundling        = "  Bundling docs-site from %s...\n"
	MsgDocsSiteBundled         = "  ✓ Docs site bundled: %s\n"
	MsgDocsSiteDownload        = "  Downloading docs-site.zip from %s...\n"
	MsgDocsSiteDownloaded      = "  ✓ Downloaded docs-site.zip (%d bytes)\n"
	ErrDocsSiteDownload        = "  ✗ Could not auto-download docs-site.zip (tried %d source(s)). Last error: %v\n    Run `gitmap update` to reinstall, or place docs-site.zip next to the binary at:\n      %s\n"
	MsgHDHostedFallback        = "  ↪ Local docs unavailable — opening hosted docs at %s\n"
	DocsSiteDownloadTimeoutSec = 30
)
