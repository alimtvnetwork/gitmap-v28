// Package termout — themes.go: design-system theme registry (program 273).
//
// Single source of truth for terminal palettes. Every hex value below is
// taken verbatim from the coding-guidelines-v24 design system —
// 02-spec/07-design-system/16-theme-catalog-and-palettes.md and the
// machine-readable 17-theme-tokens.json. No color value is invented.
//
// A terminal cannot render hex directly, so each role also carries the
// nearest SGR code (16-color or 256-color `38;5;N`), documented per theme
// above its table. Role mapping (spec 02-theme-design.md §2): fg,
// fg-muted, fg-subtle, primary, accent from the catalog's semantic
// tokens; warn/error/ok from the heatmap scale (#F59E0B/#DC2626/#10B981).
// Semantic roles are bound first; literal color codes in help/UI code
// are forbidden (design-system §7). Legacy aliases: `bright` →
// navy-purple, `standard` → vscode-dark-plus, `monochrome`/`mono` →
// SGR strip (unchanged).
package termout

import (
	"os"
	"sort"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)
import "github.com/alimtvnetwork/gitmap-v28/cli/termout/enums"

// RoleType is a semantic terminal color role. Help and UI code references
// these — never raw hex or literal SGR codes.
// RoleType aliases the strongly-typed enum from the enums package.
type RoleType = enums.RoleType

const (
	RoleFG       = enums.RoleFG
	RoleFGMuted  = enums.RoleFGMuted
	RoleFGSubtle = enums.RoleFGSubtle
	RolePrimary  = enums.RolePrimary
	RoleAccent   = enums.RoleAccent
	RoleWarn     = enums.RoleWarn
	RoleError    = enums.RoleError
	RoleOK       = enums.RoleOK
)

// Role is a deprecated legacy alias for RoleType.
type Role = enums.RoleType

// Semantic terminal roles (spec 02-theme-design.md §2).

// Theme slugs exposed via --theme / GITMAP_THEME.
const (
	ThemeNavyPurple     = "navy-purple"
	ThemeTokyoNight     = "tokyo-night"
	ThemeVSCodeDarkPlus = "vscode-dark-plus"
	ThemeMonokaiPro     = "monokai-pro"
)

// Tokyo Night truecolor palette constants (24-bit hex).
const (
	TokyoBase    = "#1A1B26"
	TokyoSurface = "#24283B"
	TokyoPrimary = "#7AA2F7"
	TokyoAccent  = "#BB9AF7"
	TokyoOK      = "#9ECE6A"
	TokyoWarn    = "#E0AF68"
	TokyoError   = "#F7768E"
	TokyoMuted   = "#565F89"
	TokyoFG      = "#C0CAF5"
)

// Palette is one named terminal palette: the design-system hex ground
// plus the nearest SGR parameter bytes per semantic role (e.g.
// "38;5;135" — callers wrap with Role()/RoleSGRBold()).
type Palette struct {
	Slug   string
	Name   string
	Source string
	Hex    map[Role]string
	SGR    map[Role]string
}

// navyPurpleTheme: Family 1 "Midnight Cyber" (catalog §2, tokens
// navy-purple). Recommended for dev tools; the default gitmap theme
// (legacy `bright` alias). Hex → nearest SGR:
// fg #F8FAFC→97, fg-muted #94A3B8→38;5;109, fg-subtle #64748B→38;5;66,
// primary #A855F7→38;5;135, accent #38BDF8→38;5;75, warn #F59E0B→38;5;214,
// error #DC2626→38;5;160, ok #10B981→38;5;36.
var navyPurpleTheme = &Palette{
	Slug:   ThemeNavyPurple,
	Name:   "Navy Blue & Purple",
	Source: "coding-guidelines-v24 16-theme-catalog-and-palettes.md §2",
	Hex: map[Role]string{
		RoleFG: "#F8FAFC", RoleFGMuted: "#94A3B8", RoleFGSubtle: "#64748B",
		RolePrimary: "#A855F7", RoleAccent: "#38BDF8",
		RoleWarn: "#F59E0B", RoleError: "#DC2626", RoleOK: "#10B981",
	},
	SGR: map[Role]string{
		RoleFG: "97", RoleFGMuted: "38;5;109", RoleFGSubtle: "38;5;66",
		RolePrimary: "38;5;135", RoleAccent: "38;5;75",
		RoleWarn: "38;5;214", RoleError: "38;5;160", RoleOK: "38;5;36",
	},
}

// tokyoNightTheme: Family 2 Tokyo Night Storm (catalog §3.2, tokens
// tokyo-night). Recommended for modern CLI tools. Hex → nearest SGR:
// fg #C0CAF5→38;5;153, fg-muted #787C99→38;5;103,
// fg-subtle #565A6E→38;5;59, primary #7AA2F7→38;5;111,
// accent #BB9AF7→38;5;147 (tokens --accent; the spec summary table
// lists --cyan #7DCFFF — the semantic --accent binding wins),
// warn #F59E0B→38;5;214, error #DC2626→38;5;160, ok #10B981→38;5;36.
var tokyoNightTheme = &Palette{
	Slug:   ThemeTokyoNight,
	Name:   "Tokyo Night Storm",
	Source: "coding-guidelines-v24 16-theme-catalog-and-palettes.md §3.2",
	Hex: map[Role]string{
		RoleFG: "#C0CAF5", RoleFGMuted: "#787C99", RoleFGSubtle: "#565A6E",
		RolePrimary: "#7AA2F7", RoleAccent: "#BB9AF7",
		RoleWarn: "#F59E0B", RoleError: "#DC2626", RoleOK: "#10B981",
	},
	SGR: map[Role]string{
		RoleFG: "38;5;153", RoleFGMuted: "38;5;103", RoleFGSubtle: "38;5;59",
		RolePrimary: "38;5;111", RoleAccent: "38;5;147",
		RoleWarn: "38;5;214", RoleError: "38;5;160", RoleOK: "38;5;36",
	},
}

// vscodeDarkPlusTheme: Family 2 VS Code Dark+ (catalog §3.1, tokens
// vscode-dark-plus). Instant familiarity for engineers (legacy
// `standard` alias). Hex → nearest SGR:
// fg #D4D4D4→38;5;252, fg-muted #858585→38;5;245,
// fg-subtle #6E6E6E→38;5;242, primary #007ACC→38;5;32,
// accent #4EC9B0→38;5;79, warn #F59E0B→38;5;214,
// error #DC2626→38;5;160, ok #10B981→38;5;36.
var vscodeDarkPlusTheme = &Palette{
	Slug:   ThemeVSCodeDarkPlus,
	Name:   "VS Code Dark+",
	Source: "coding-guidelines-v24 16-theme-catalog-and-palettes.md §3.1",
	Hex: map[Role]string{
		RoleFG: "#D4D4D4", RoleFGMuted: "#858585", RoleFGSubtle: "#6E6E6E",
		RolePrimary: "#007ACC", RoleAccent: "#4EC9B0",
		RoleWarn: "#F59E0B", RoleError: "#DC2626", RoleOK: "#10B981",
	},
	SGR: map[Role]string{
		RoleFG: "38;5;252", RoleFGMuted: "38;5;245", RoleFGSubtle: "38;5;242",
		RolePrimary: "38;5;32", RoleAccent: "38;5;79",
		RoleWarn: "38;5;214", RoleError: "38;5;160", RoleOK: "38;5;36",
	},
}

// monokaiProTheme: Family 2 Monokai Pro (catalog §3.5, tokens
// monokai-pro). High-impact telemetry/status. Hex → nearest SGR:
// fg #FCFCFA→97, fg-muted #939293→38;5;246, fg-subtle #727072→38;5;243,
// primary #FF6188→38;5;204, accent #FFD866→38;5;221,
// warn #F59E0B→38;5;214, error #DC2626→38;5;160, ok #10B981→38;5;36.
var monokaiProTheme = &Palette{
	Slug:   ThemeMonokaiPro,
	Name:   "Monokai Pro",
	Source: "coding-guidelines-v24 16-theme-catalog-and-palettes.md §3.5",
	Hex: map[Role]string{
		RoleFG: "#FCFCFA", RoleFGMuted: "#939293", RoleFGSubtle: "#727072",
		RolePrimary: "#FF6188", RoleAccent: "#FFD866",
		RoleWarn: "#F59E0B", RoleError: "#DC2626", RoleOK: "#10B981",
	},
	SGR: map[Role]string{
		RoleFG: "97", RoleFGMuted: "38;5;246", RoleFGSubtle: "38;5;243",
		RolePrimary: "38;5;204", RoleAccent: "38;5;221",
		RoleWarn: "38;5;214", RoleError: "38;5;160", RoleOK: "38;5;36",
	},
}

// paletteRegistry is the single theme table keyed by slug.
var paletteRegistry = map[string]*Palette{
	ThemeNavyPurple:     navyPurpleTheme,
	ThemeTokyoNight:     tokyoNightTheme,
	ThemeVSCodeDarkPlus: vscodeDarkPlusTheme,
	ThemeMonokaiPro:     monokaiProTheme,
}

// themeAlias maps legacy labels to registry slugs. monochrome/mono are
// not palettes — they strip SGR and are handled by resolveThemeLabel.
var themeAlias = map[string]string{
	constants.ThemeBright:   ThemeNavyPurple,
	constants.ThemeStandard: ThemeVSCodeDarkPlus,
}

// resolveThemeLabel normalizes a user label to a registry slug.
// Returns strip=true for monochrome/mono (SGR strip, no palette).
func resolveThemeLabel(label string) (slug string, strip bool) {
	norm := strings.ToLower(strings.TrimSpace(label))
	if norm == constants.ThemeMonochrome || norm == constants.ThemeMono {
		return "", true
	}
	if aliased, ok := themeAlias[norm]; ok {
		return aliased, false
	}
	if _, ok := paletteRegistry[norm]; ok {
		return norm, false
	}

	return ThemeNavyPurple, false
}

// ResolveTheme returns the active palette for the current GITMAP_THEME
// value (set by the --theme flag stripper). Unknown/empty labels fall
// back to navy-purple. strip=true means monochrome/mono: no palette,
// callers must emit no SGR at all.
func ResolveTheme() (palette *Palette, strip bool) {
	slug, strip := resolveThemeLabel(os.Getenv(constants.EnvTheme))
	if strip {
		return nil, true
	}

	return paletteRegistry[slug], false
}

// ActivePalette is ResolveTheme for the current process. Convenience
// for renderers that need hex ground (docs UI) as well as SGR.
func ActivePalette() *Palette {
	palette, strip := ResolveTheme()
	if strip {
		return nil
	}

	return palette
}

// ThemeSlugs returns the sorted registry slugs for --theme validation
// messages and help text.
func ThemeSlugs() []string {
	slugs := make([]string, 0, len(paletteRegistry))
	for slug := range paletteRegistry {
		slugs = append(slugs, slug)
	}

	sort.Strings(slugs)

	return slugs
}

// colorAllowed reports whether role colors may be emitted: stdout must
// be a real TTY (dumb-terminal safety — piped output stays plain) and
// the theme must not be the mono strip mode.
func colorAllowed() bool {
	if _, strip := ResolveTheme(); strip {
		return false
	}

	return IsStdoutTTY()
}

// RoleSGR returns the SGR escape for role in the active theme, or "" when
// colors are not allowed (piped output or mono theme). Renderers call
// this instead of embedding literal color codes.
func RoleSGR(role Role) string {
	if !colorAllowed() {
		return ""
	}

	palette, _ := ResolveTheme()
	code, ok := palette.SGR[role]

	if !ok {
		return ""
	}

	return "\033[" + code + "m"
}

// RoleSGRBold is RoleSGR with the bold attribute (headings, banners).
func RoleSGRBold(role Role) string {
	if !colorAllowed() {
		return ""
	}

	palette, _ := ResolveTheme()
	code, ok := palette.SGR[role]

	if !ok {
		return ""
	}

	return "\033[1;" + code + "m"
}

// Reset is the SGR reset sequence. Empty when colors are disallowed so
// piped output never carries stray escapes.
func Reset() string {
	if !colorAllowed() {
		return ""
	}

	return constants.ColorReset
}

// PlasmaScale is the catalog §4.1 sequential scale (hex ground) for
// intensity heatmaps: latency, CPU, request volume, error density.
// SGR mapping is recorded when a heatmap renderer lands; until then
// this table is the single source of truth for the scale values.
var PlasmaScale = []string{
	"#0D0887", "#5A01A5", "#8F0DA4", "#BC3754", "#DE5F39", "#F58C46", "#F0F921",
}

// RAGScale is the catalog §4.2 diverging Red-Amber-Green scale (hex
// ground) for quality/risk/SLA heatmaps.
var RAGScale = map[string]string{
	"critical": "#DC2626", "high": "#EF4444", "warn": "#F59E0B",
	"neutral": "#64748B", "acceptable": "#84CC16", "optimal": "#10B981",
}

// DensityScale is the catalog §4.3 GitHub-style activity density grid
// (dark-mode hex ground) for commit frequency / uptime streaks.
var DensityScale = []string{
	"#161B22", "#0E4429", "#006D32", "#26A641", "#39D353",
}
