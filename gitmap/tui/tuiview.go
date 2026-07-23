package tui

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

func (m rootModel) View() string {
	if m.quitting {
		return ""
	}

	var b strings.Builder
	b.WriteString(styleTitle.Render(constants.TUITitle))
	b.WriteString("\n")
	b.WriteString(m.renderTabs())
	b.WriteString("\n\n")
	b.WriteString(m.renderContent())
	b.WriteString("\n")
	b.WriteString(m.renderStatusBar())

	return b.String()
}

func (m rootModel) renderTabs() string {
	labels := []string{
		constants.TUIViewBrowser,
		constants.TUIViewActions,
		constants.TUIViewGroups,
		constants.TUIViewDashboard,
		constants.TUIViewReleases,
		constants.TUIViewTempReleases,
		constants.TUIViewZipGroups,
		constants.TUIViewAliases,
		constants.TUIViewLogs,
	}

	var tabs []string
	for i, label := range labels {
		if i == m.activeTab {
			tabs = append(tabs, styleActiveTab.Render(label))
		} else {
			tabs = append(tabs, styleTab.Render(label))
		}
	}

	return strings.Join(tabs, " ")
}

func (m rootModel) renderContent() string {
	switch m.activeTab {
	case viewBrowser:
		return m.browser.View()
	case viewActions:
		return m.actions.View()
	case viewGroups:
		return m.groupsMgr.View()
	case viewDashboard:
		return m.dashboard.View()
	case viewReleases:
		return m.releases.View()
	case viewTempReleases:
		return m.tempReleases.View()
	case viewZipGroups:
		return m.zipGroups.View()
	case viewAliases:
		return m.aliases.View()
	case viewLogs:
		return m.logs.View()
	}

	return ""
}

func (m rootModel) renderStatusBar() string {
	hints := []string{constants.TUIQuitHint, constants.TUITabHint}

	switch m.activeTab {
	case viewBrowser:
		hints = append(hints, constants.TUISelectHint)
	case viewActions:
		hints = append(hints, constants.TUIBatchHint)
	case viewGroups:
		hints = append(hints, constants.TUIGroupHint)
	case viewDashboard:
		hints = append(hints, constants.TUIDashHint)
	case viewReleases:
		hints = append(hints, constants.TUIRelHint)
	case viewTempReleases:
		hints = append(hints, constants.TUITRHint)
	case viewZipGroups:
		hints = append(hints, constants.TUIZGHint)
	case viewAliases:
		hints = append(hints, constants.TUIAliasHint)
	case viewLogs:
		hints = append(hints, constants.TUILogHint)
	}

	return styleStatusBar.Render(strings.Join(hints, "  │  "))
}
