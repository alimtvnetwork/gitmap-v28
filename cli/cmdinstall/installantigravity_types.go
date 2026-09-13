package cmdinstall

// AntigravityPlatformInfo encapsulates platform metadata and target locations for installation.
type AntigravityPlatformInfo struct {
	OS           string
	Arch         string
	DistroName   string
	PlatformID   string
	ArtifactName string
	DownloadURL  string
	InstallDir   string
	BinDir       string
}

const (
	AntigravityDefaultVersion = "2.13.0"
	AntigravityDefaultBuildID = "6362815968182272"
	AntigravityBaseURL        = "https://storage.googleapis.com/antigravity-public/antigravity-hub"
	ErrUnsupportedPlatform    = "E_UNSUPPORTED_PLATFORM"
	ErrPrerequisiteFailed     = "E_PREREQUISITE_FAILED"
	ErrDownloadFailed         = "E_DOWNLOAD_FAILED"
)

// HasValidPlatform reports whether PlatformID and ArtifactName are populated.
func (p AntigravityPlatformInfo) HasValidPlatform() bool {
	return p.PlatformID != "" && p.ArtifactName != ""
}

// HasDownloadURL reports whether DownloadURL is populated.
func (p AntigravityPlatformInfo) HasDownloadURL() bool {
	return p.DownloadURL != ""
}

// IsDarwin reports whether the platform is macOS.
func (p AntigravityPlatformInfo) IsDarwin() bool {
	return p.OS == "darwin" || p.OS == "macos"
}

// IsWindows reports whether the platform is Windows.
func (p AntigravityPlatformInfo) IsWindows() bool {
	return p.OS == "windows"
}

// IsLinux reports whether the platform is Linux.
func (p AntigravityPlatformInfo) IsLinux() bool {
	return p.OS == "linux"
}
