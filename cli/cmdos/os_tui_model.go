package cmdos

import tea "github.com/charmbracelet/bubbletea"

// OSTUIModel encapsulates state and update lifecycle for the OS TUI.
type OSTUIModel struct {
	Tabs      []OSTUITab
	ActiveTab int
	Cursor    int
	Mode      OSTUIViewModeType
	Results   []OSTUIActionResult
	StatusMsg string
	IsDryRun  bool
	IsDone    bool
}

// NewOSTUIModel instantiates a default dashboard model.
func NewOSTUIModel(isDryRun bool) OSTUIModel {
	return OSTUIModel{
		Tabs:     buildOSTUITabs(),
		Mode:     OSTUIViewModeNavType,
		IsDryRun: isDryRun,
	}
}

// Init satisfies the tea.Model interface.
func (m OSTUIModel) Init() tea.Cmd {
	return nil
}

// Update routes messages to key or resize handlers.
func (m OSTUIModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch t := msg.(type) {
	case tea.KeyMsg:
		return m.handleKey(t)
	}

	return m, nil
}

func (m OSTUIModel) handleKey(k tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.Mode == OSTUIViewModeDoneType {
		return m.handleDoneKey(k)
	}

	return m.handleNavKey(k)
}

func (m OSTUIModel) handleDoneKey(k tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch k.String() {
	case "q", "ctrl+c":
		m.IsDone = true
		return m, tea.Quit
	case "esc", "enter", "backspace":
		m.Mode = OSTUIViewModeNavType
		m.StatusMsg = "Ready"
		return m, nil
	}

	return m, nil
}
