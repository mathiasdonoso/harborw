package tui

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mathiasdonoso/harborw/internal/api/harbor"
)

type RepositoriesState struct {
	table table.Model
}

var repositoriesLoaded = false

func (m model) repositoriesView() string {
	return m.state.repositories.table.View()
}

func (m model) repositoriesUpdate(msg tea.Msg) (model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "-":
			m = m.SwitchPage(projectsPage)
			return m, nil
		case "enter":
			fmt.Printf("Press enter\n")
			m.state.artifacts = m.NewArtifactsState("onboarding", "ng-ui-mx")
			m = m.SwitchPage(artifactsPage)
			return m, nil
		}
	}
	m.state.repositories.table, cmd = m.state.repositories.table.Update(msg)
	return m, cmd
}

func (m model) NewRepositoriesState(project string) table.Model {
	const REPOSITORY_COLUMN_NAME = "Repository"
	const ARTIFACTS_COUNT_COLUMN_NAME = "Artifacts count"

	columns := []table.Column{
		{Title: REPOSITORY_COLUMN_NAME, Width: 40},
		{Title: ARTIFACTS_COUNT_COLUMN_NAME, Width: len(ARTIFACTS_COUNT_COLUMN_NAME)},
	}

	harborClient, err := harbor.NewHarborApiClient(http.DefaultClient)
	if err != nil {
		fmt.Printf("Error creating harbor client: %v\n", err)
	}

	r, err := harborClient.FetchRepositories(project)
	if err != nil {
		fmt.Printf("Error fetching projects: %v\n", err)
	}
	repositoriesLoaded = true

	rows := []table.Row{}
	for _, v := range *r {
		rows = append(rows, table.Row{v.Name, strconv.Itoa(v.ArtifactCount)})
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
