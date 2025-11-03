package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type page int

const (
	menuPage page = iota
	projectsPage
	repositoriesPage
	artifactsPage
	statusPage
)

type state struct {
	projects     ProjectsState
	repositories RepositoriesState
	artifacts    ArtifactsState
}

type model struct {
	page     page
	state    state
	renderer *lipgloss.Renderer
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch m.page {
	case menuPage:
		m, cmd = m.menuUpdate(msg)
	case projectsPage:
		m, cmd = m.projectsUpdate(msg)
	case repositoriesPage:
		m, cmd = m.repositoriesUpdate(msg)
	case artifactsPage:
		m, cmd = m.artifactsUpdate(msg)
	case statusPage:
		m, cmd = m.statusUpdate(msg)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	return m, cmd
}

func (m model) View() string {
	header := m.headerView()
	content := m.getContent()
	footer := m.footerView()

	items := []string{}
	items = append(items, header)
	items = append(items, content)
	items = append(items, footer)

	child := lipgloss.JoinVertical(
		lipgloss.Left,
		items...,
	)

	return m.renderer.Place(
		lipgloss.NewStyle().GetMaxWidth(),
		lipgloss.NewStyle().GetMaxHeight(),
		lipgloss.Left,
		lipgloss.Bottom,
		child,
	)
}

func NewModel() (tea.Model, error) {
	m := model{
		page:     projectsPage,
		renderer: &lipgloss.Renderer{},
		state: state{
			projects:     ProjectsState{},
			repositories: RepositoriesState{},
			artifacts:    ArtifactsState{},
		},
	}

	m.state.projects = m.NewProjectsState()

	return m, nil
}

func (m model) SwitchPage(page page) model {
	m.page = page
	return m
}

func (m model) getContent() string {
	page := "unknown"
	switch m.page {
	case menuPage:
		page = m.menuView()
	case projectsPage:
		page = m.projectsView()
	case repositoriesPage:
		page = m.repositoriesView()
	case artifactsPage:
		page = m.artifactsView()
	case statusPage:
		page = m.statusView()
	}
	return page
}
