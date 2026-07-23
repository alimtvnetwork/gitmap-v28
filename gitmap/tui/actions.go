package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

type actionsModel struct {
	lastResult string
}

func newActionsModel() actionsModel {
	return actionsModel{}
}

func (m actionsModel) Update(msg tea.Msg, selected []model.ScanRecord) (actionsModel, tea.Cmd) {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	count := len(selected)
	if count == 0 {
		m.lastResult = constants.TUINoSelection

		return m, nil
	}

	switch {
	case keys.pull(keyMsg):
		m.lastResult = fmt.Sprintf(constants.TUIActionPull, count)
	case keys.exec(keyMsg):
		m.lastResult = fmt.Sprintf(constants.TUIActionExec, count)
	case keys.status(keyMsg):
		m.lastResult = fmt.Sprintf(constants.TUIActionStatus, count)
	case keys.addToGroup(keyMsg):
		m.lastResult = fmt.Sprintf("Adding %d repo(s) to group...", count)
	}

	return m, nil
}

func (m actionsModel) View() string {
	var b strings.Builder

	b.WriteString(styleHeader.Render("Batch Actions"))
	b.WriteString("\n\n")
	b.WriteString(styleNormalRow.Render("  p — Pull selected repos"))
	b.WriteString("\n")
	b.WriteString(styleNormalRow.Render("  x — Execute git command across selected"))
	b.WriteString("\n")
	b.WriteString(styleNormalRow.Render("  s — Show status for selected"))
	b.WriteString("\n")
	b.WriteString(styleNormalRow.Render("  g — Add selected to a group"))
	b.WriteString("\n")

	if m.lastResult != "" {
		b.WriteString("\n")
		b.WriteString(styleSearch.Render(m.lastResult))
	}

	return b.String()
}
