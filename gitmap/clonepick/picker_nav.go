package clonepick

// picker_nav.go: cursor + scroll math for the bubbletea picker.
// Pure functions / private helpers split out of picker.go so it
// stays under the strict 200-line cap.

// applyCursorMove updates the cursor for navigation keys. Page keys
// jump by viewportHeight; g / G are vim-style home / end. Toggle
// keys (space, a, n) fall through unchanged -- handleNavKey owns
// the selection mutation.
func applyCursorMove(m pickerModel, key string) pickerModel {
	switch key {
	case "up", "k":
		m.cursor = maxInt(0, m.cursor-1)
	case "down", "j":
		m.cursor = minInt(len(m.paths)-1, m.cursor+1)
	case "pgup", "ctrl+b":
		m.cursor = maxInt(0, m.cursor-m.viewportHeight)
	case "pgdown", "ctrl+f":
		m.cursor = minInt(len(m.paths)-1, m.cursor+m.viewportHeight)
	case "g", "home":
		m.cursor = 0
	case "G", "end":
		m.cursor = maxInt(0, len(m.paths)-1)
	}

	return m
}

// clampScroll keeps the cursor in the visible window. Returns the
// new offset such that cursor in [offset, offset+height). Pure --
// no model mutation -- so the caller decides when to commit it.
func clampScroll(cursor, offset, height, total int) int {
	if height < 1 || total == 0 {
		return 0
	}
	if cursor < offset {
		return cursor
	}
	if cursor >= offset+height {
		return cursor - height + 1
	}
	maxOffset := total - height
	if maxOffset < 0 {
		maxOffset = 0
	}
	if offset > maxOffset {
		return maxOffset
	}

	return offset
}

// minInt / maxInt: tiny stdlib-free helpers so picker.go has no
// external dep beyond bubbletea. Kept private to the package.
func minInt(a, b int) int {
	if a < b {
		return a
	}

	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}

	return b
}
