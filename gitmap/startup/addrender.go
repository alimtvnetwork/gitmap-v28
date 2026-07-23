package startup

// Helpers for Add: .desktop body rendering + atomic write.
// Split from add.go to keep both files under the per-file budget
// and so the renderer is independently testable without touching
// the filesystem.

import (
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// renderDesktop builds the .desktop body. Field order matches what
// `desktop-file-validate` and the freedesktop.org spec recommend:
// Type and Name first (they're the required identity pair), then
// optional descriptors, then Exec, then the gitmap marker LAST so
// it's visible at the bottom of `cat ~/.config/autostart/...` for
// quick eyeballing without scrolling.
//
// We deliberately do NOT emit Categories or MimeType — autostart
// entries are not application-launcher items, those fields are
// noise for this use case.
func renderDesktop(clean string, opts AddOptions) []byte {
	display := opts.DisplayName
	if len(display) == 0 {
		display = clean
	}
	var b strings.Builder
	b.WriteString("[Desktop Entry]\n")
	b.WriteString("Type=Application\n")
	fmt.Fprintf(&b, "Name=%s\n", display)
	if len(opts.Comment) > 0 {
		fmt.Fprintf(&b, "Comment=%s\n", opts.Comment)
	}
	fmt.Fprintf(&b, "Exec=%s\n", opts.Exec)
	// Path= is the XDG-spec field for the working directory the
	// session manager sets before invoking Exec=. Emitted before
	// Terminal= to match the field order recommended by
	// `desktop-file-validate` (identity → exec → path → terminal).
	if len(opts.WorkingDir) > 0 {
		fmt.Fprintf(&b, "Path=%s\n", opts.WorkingDir)
	}
	b.WriteString("Terminal=false\n")
	b.WriteString("X-GNOME-Autostart-enabled=true\n")
	if opts.NoDisplay {
		b.WriteString("NoDisplay=true\n")
	}
	fmt.Fprintf(&b, "%s=%s\n", constants.StartupMarkerKey, constants.StartupMarkerVal)

	return []byte(b.String())
}

// atomicWrite stages the body in a sibling temp file then renames
// it over the final path. Rename is atomic on POSIX filesystems for
// paths within the same directory — this guarantees the autostart
// session never sees a partially-written .desktop file even if
// gitmap is killed mid-write.
//
// The temp file uses a `.gitmap-tmp-` prefix so a crash leaves
// debris under a recognizable name (rather than a random tempname
// the user has to investigate).
func atomicWrite(target string, body []byte) error {
	tmp := tempPathFor(target)
	if err := os.WriteFile(tmp, body, 0o644); err != nil {
		return fmt.Errorf("write temp %s: %w", tmp, err)
	}
	if err := os.Rename(tmp, target); err != nil {
		_ = os.Remove(tmp)
		return fmt.Errorf("rename %s -> %s: %w", tmp, target, err)
	}

	return nil
}

// tempPathFor returns "<dir>/.gitmap-tmp-<basename>" so the rename
// stays inside the same directory (cross-directory rename is NOT
// atomic on Linux when the source/dest are on different filesystems).
func tempPathFor(target string) string {
	idx := strings.LastIndex(target, "/")
	if idx < 0 {
		return ".gitmap-tmp-" + target
	}

	return target[:idx+1] + ".gitmap-tmp-" + target[idx+1:]
}
