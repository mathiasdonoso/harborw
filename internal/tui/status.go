package tui

import tea "github.com/charmbracelet/bubbletea"

type StatusState struct {
}

func (m model) statusView() string {
	return ""
}

func (m model) statusUpdate(msg tea.Msg) (model, tea.Cmd) {
	return m, nil
}
