package cmdos

import tea "github.com/charmbracelet/bubbletea"

func (m OSTUIModel) handleNavKey(k tea.KeyMsg) (tea.Model, tea.Cmd) {
	keyStr := k.String()
	switch keyStr {
	case "q", "ctrl+c":
		m.IsDone = true
		return m, tea.Quit
	case "tab", "l", "right":
		return handleTabSwitch(m, 1), nil
	case "shift+tab", "h", "left":
		return handleTabSwitch(m, -1), nil
	case "up", "k":
		return handleCursorMove(m, -1), nil
	case "down", "j":
		return handleCursorMove(m, 1), nil
	case " ":
		return toggleCurrentItem(m), nil
	case "a":
		return setAllItemsInTab(m, true), nil
	case "n":
		return setAllItemsInTab(m, false), nil
	case "enter":
		return m.triggerExecution(), nil
	default:
		return m.handleNumberKey(keyStr), nil
	}
}

func (m OSTUIModel) handleNumberKey(keyStr string) OSTUIModel {
	if len(keyStr) == 1 && keyStr[0] >= '1' && keyStr[0] <= '6' {
		idx := int(keyStr[0] - '1')
		return handleJumpTab(m, idx)
	}

	return m
}

func (m OSTUIModel) triggerExecution() OSTUIModel {
	count := countAllCheckedItems(m.Tabs)
	if count == 0 {
		m.StatusMsg = "Select at least one item (Space) before executing"
		return m
	}

	m.Results = executeCheckedItems(m.Tabs, m.IsDryRun)
	m.Mode = OSTUIViewModeDoneType

	return m
}
