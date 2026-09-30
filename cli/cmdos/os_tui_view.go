package cmdos

import (
	"fmt"
	"strings"
)

// View renders the TUI interface per Bubbletea contract.
func (m OSTUIModel) View() string {
	if m.Mode == OSTUIViewModeDoneType {
		return renderDoneView(m)
	}

	return renderNavView(m)
}

func renderNavView(m OSTUIModel) string {
	var b strings.Builder
	totalChecked := countAllCheckedItems(m.Tabs)
	header := fmt.Sprintf("GITMAP OS DASHBOARD  [%d selected across %d tabs]\n", totalChecked, len(m.Tabs))
	b.WriteString(tuiStyleTitle.Render(header))
	b.WriteString(renderTabBar(m))
	b.WriteString("\n" + strings.Repeat("─", 80) + "\n")
	b.WriteString(renderItemList(m))
	b.WriteString(strings.Repeat("─", 80) + "\n")
	b.WriteString(renderFooter(m))

	return b.String()
}

func renderTabBar(m OSTUIModel) string {
	var tabs []string
	for i, t := range m.Tabs {
		label := fmt.Sprintf("[%s] %s", t.ShortcutKey, t.Title)
		if i == m.ActiveTab {
			tabs = append(tabs, tuiStyleTabActive.Render(label))
		} else {
			tabs = append(tabs, tuiStyleTabInactive.Render(label))
		}
	}

	return strings.Join(tabs, "  ")
}

func renderItemList(m OSTUIModel) string {
	if m.ActiveTab >= len(m.Tabs) {
		return "  No items in active tab\n"
	}

	items := m.Tabs[m.ActiveTab].Items
	if len(items) == 0 {
		return "  No items available\n"
	}

	var b strings.Builder
	for i, item := range items {
		isCursor := i == m.Cursor
		b.WriteString(renderItemRow(item, isCursor))
		b.WriteByte('\n')
	}

	return b.String()
}

func renderItemRow(item OSTUIItem, isCursor bool) string {
	prefix := "  "
	if isCursor {
		prefix = tuiStyleCursor.Render("> ")
	}

	mark := "[ ]"
	if item.IsChecked {
		mark = "[x]"
	}

	badge := tuiStyleBadge.Render(fmt.Sprintf("[%s]", item.Badge))
	title := tuiStyleNormal.Render(item.Title)

	return fmt.Sprintf("%s%s %-32s %-12s %s", prefix, mark, title, badge, item.Description)
}

func renderFooter(m OSTUIModel) string {
	hints := "[Tab/←/→] Tabs | [1-6] Jump | [↑/↓] Move | [Space] Toggle | [Enter] Apply | [q] Quit\n"
	status := "Status: Ready"
	if m.StatusMsg != "" {
		status = "Status: " + m.StatusMsg
	}

	return tuiStyleFooter.Render(hints + status + "\n")
}
