package tui

import tea "github.com/charmbracelet/bubbletea"

func (m model) menuView() string {
	return "menu-view"
}

func (m model) menuUpdate(msg tea.Msg) (model, tea.Cmd) {
	return m, nil
}
