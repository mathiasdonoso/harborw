package tui

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/mathiasdonoso/harborw/internal/api/harbor"
)

type Project struct {
	// Selected bool
	Columns []string
}

type ProjectsState struct {
	table  table.Model
	data   []Project
	status string
	filter string
}

var projectsLoaded = false

func (m model) projectsView() string {
	var s string

	if m.state.projects.status == "filtering" {
		s += fmt.Sprintf("Filter: %s\n", m.state.projects.filter)
	}

	s += m.state.projects.table.View()
	return s
}

func (m model) projectsUpdate(msg tea.Msg) (model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "/":
			fmt.Printf("Press /\n")
			m.state.projects.status = "filtering"
			m.state.projects.filter = ""
			return m, nil
		case "enter":
			fmt.Printf("Press enter\n")
			m.state.repositories.table = m.NewRepositoriesState("onboarding")
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

func (m model) NewProjectsState() table.Model {
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
	projectsLoaded = true

	projects := make([]Project, len(*r))

	for i, p := range *r {
		projects[i] = Project{
			Columns: []string{p.Name, strconv.Itoa(p.RepoCount)},
		}
	}

	rows := make([]table.Row, len(*r))
	for i, a := range projects {
		rows[i] = a.Columns
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(21),
	)

	t.SetStyles(GetTableDefaultStyles())

	return t
}
