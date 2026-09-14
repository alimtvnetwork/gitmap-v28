package cmdinstall

// StandardIconResolutions defines canonical FreeDesktop icon resolutions.
var StandardIconResolutions = []int{16, 24, 32, 48, 64, 128, 256, 512}

// AntigravityCanonicalIconName defines canonical unextended icon name.
const AntigravityCanonicalIconName = "antigravity"

// HicolorIndexThemeContent contains standard FreeDesktop hicolor index.theme content.
const HicolorIndexThemeContent = `[Icon Theme]
Name=Hicolor
Comment=Fallback icon theme
Hidden=true
Directories=16x16/apps,24x24/apps,32x32/apps,48x48/apps,64x64/apps,128x128/apps,256x256/apps,512x512/apps

[16x16/apps]
Size=16
Type=Threshold

[24x24/apps]
Size=24
Type=Threshold

[32x32/apps]
Size=32
Type=Threshold

[48x48/apps]
Size=48
Type=Threshold

[64x64/apps]
Size=64
Type=Threshold

[128x128/apps]
Size=128
Type=Threshold

[256x256/apps]
Size=256
Type=Threshold

[512x512/apps]
Size=512
Type=Threshold
`

// ResizedIcon contains icon image bytes for a specific square pixel resolution.
type ResizedIcon struct {
	Size int
	Data []byte
}

// IconDeployOptions specifies parameters for deploying multi-resolution desktop icons.
type IconDeployOptions struct {
	IconName   string
	SourceDir  string
	RawBytes   []byte
	IsUserOnly bool
}
