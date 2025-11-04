package tui

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/mathiasdonoso/harborw/internal/api/harbor"
)

type Project struct {
	Name      string
	RepoCount int
}

func (p Project) ToRow() []string {
	columns := []string{p.Name, strconv.Itoa(p.RepoCount)}
	return columns
}

type ProjectsState struct {
	table table.Model
	data  []Project
}

func (m model) projectsView() string {
	return m.state.projects.table.View()
}

func (m model) projectsUpdate(msg tea.Msg) (model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			rowIndex := m.state.projects.table.Cursor()
			active := m.state.projects.data[rowIndex]
			m.state.repositories = m.NewRepositoriesState(active.Name)
			slog.Debug(fmt.Sprintf("Selecting project: %s", active.Name))
			m = m.SwitchPage(repositoriesPage)
			return m, nil
		}
	}

	m.state.projects.table, cmd = m.state.projects.table.Update(msg)
	return m, cmd
}

var PROJECTS_COLUMNS = []table.Column{
	{Title: "Project", Width: 40},
	{Title: "Repositories count", Width: 18},
}

func newEmptyProjectsState() ProjectsState {
	t := table.New(
		table.WithColumns(PROJECTS_COLUMNS),
		table.WithRows([]table.Row{{"No data available", ""}}),
		table.WithFocused(true),
		table.WithHeight(2),
	)

	t.SetStyles(GetTableDefaultStyles())

	state := ProjectsState{
		table: t,
		data:  []Project{},
	}

	return state
}

func (m model) NewProjectsState() ProjectsState {
	harborClient, err := harbor.NewHarborApiClient(http.DefaultClient)
	if err != nil {
		slog.Error("Error creating harbor client", "err", err)
		return newEmptyProjectsState()
	}
	r, err := harborClient.FetchProjects()
	if err != nil {
		slog.Error("Error fetching projects", "err", err)
		return newEmptyProjectsState()
	}

	projects := make([]Project, len(*r))
	for i, p := range *r {
		project := Project{
			Name:      p.Name,
			RepoCount: p.RepoCount,
		}
		projects[i] = project
	}

	rows := make([]table.Row, len(*r))
	for i, a := range projects {
		rows[i] = a.ToRow()
	}

	t := table.New(
		table.WithColumns(PROJECTS_COLUMNS),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(21),
	)

	t.SetStyles(GetTableDefaultStyles())

	state := ProjectsState{
		table: t,
		data:  projects,
	}

	slog.Debug("New Project state created.")

	return state
}
