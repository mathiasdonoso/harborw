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

func (p Project) ToColumn() []string {
	columns := []string{p.Name, strconv.Itoa(p.RepoCount)}
	return columns
}

type ProjectsState struct {
	table  table.Model
	data   []Project
	status string
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
		default:
			if m.state.projects.status == "filtering" {
			}
		}
	}

	m.state.projects.table, cmd = m.state.projects.table.Update(msg)
	return m, cmd
}

func (m model) NewProjectsState() ProjectsState {
	const PROJECT_COLUMN_NAME = "Project"
	const REPOSITORIES_COUNT_COLUMN_NAME = "Repositores count"

	columns := []table.Column{
		{Title: PROJECT_COLUMN_NAME, Width: 40},
		{Title: REPOSITORIES_COUNT_COLUMN_NAME, Width: len(REPOSITORIES_COUNT_COLUMN_NAME)},
	}

	harborClient, err := harbor.NewHarborApiClient(http.DefaultClient)
	if err != nil {
		fmt.Printf("Error creating harbor client: %v\n", err)
	}
	r, err := harborClient.FetchProjects()
	if err != nil {
		fmt.Printf("Error fetching projects: %v\n", err)
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
		rows[i] = a.ToColumn()
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(21),
	)

	t.SetStyles(GetTableDefaultStyles())

	state := ProjectsState{
		table:  t,
		data:   projects,
		status: "",
	}

	slog.Debug("New Project state created.")

	return state
}
