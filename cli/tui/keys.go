package tui

import tea "github.com/charmbracelet/bubbletea"

// keyMap defines all TUI key bindings.
type keyMap struct{}

func (k keyMap) quit(msg tea.KeyMsg) bool {
	return msg.String() == "q" || msg.String() == "esc"
}

func (k keyMap) tab(msg tea.KeyMsg) bool {
	return msg.String() == "tab"
}

func (k keyMap) up(msg tea.KeyMsg) bool {
	return msg.String() == "up" || msg.String() == "k"
}

func (k keyMap) down(msg tea.KeyMsg) bool {
	return msg.String() == "down" || msg.String() == "j"
}

func (k keyMap) selectItem(msg tea.KeyMsg) bool {
	return msg.String() == " "
}

func (k keyMap) enter(msg tea.KeyMsg) bool {
	return msg.String() == "enter"
}

func (k keyMap) search(msg tea.KeyMsg) bool {
	return msg.String() == "/"
}

func (k keyMap) selectAll(msg tea.KeyMsg) bool {
	return msg.String() == "a"
}

func (k keyMap) pull(msg tea.KeyMsg) bool {
	return msg.String() == "p"
}

func (k keyMap) exec(msg tea.KeyMsg) bool {
	return msg.String() == "x"
}

func (k keyMap) status(msg tea.KeyMsg) bool {
	return msg.String() == "s"
}

func (k keyMap) addToGroup(msg tea.KeyMsg) bool {
	return msg.String() == "g"
}

func (k keyMap) create(msg tea.KeyMsg) bool {
	return msg.String() == "c"
}

func (k keyMap) delete(msg tea.KeyMsg) bool {
	return msg.String() == "d"
}

func (k keyMap) refresh(msg tea.KeyMsg) bool {
	return msg.String() == "r"
}

var keys = keyMap{}
