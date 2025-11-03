package tui

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mathiasdonoso/harborw/internal/api/harbor"
)

type Repository struct {
	Name           string
	Project        string
	ArtifactsCount int
}

func (r Repository) ToColumn() []string {
	decodedOnce, _ := url.PathUnescape(r.Name)
	decodedTwice, _ := url.PathUnescape(decodedOnce)
	columns := []string{
		decodedTwice,
		strconv.Itoa(r.ArtifactsCount),
	}
	return columns
}

type RepositoriesState struct {
	table table.Model
	data  []Repository
}

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
			rowIndex := m.state.repositories.table.Cursor()
			active := m.state.repositories.data[rowIndex]
			m.state.artifacts = m.NewArtifactsState(active.Project, active.Name)
			m = m.SwitchPage(artifactsPage)
			return m, nil
		}
	}
	m.state.repositories.table, cmd = m.state.repositories.table.Update(msg)
	return m, cmd
}

func (m model) NewRepositoriesState(project string) RepositoriesState {
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

	repositories := make([]Repository, len(*r))
	for i, r := range *r {
		nameSections := strings.Split(r.Name, "/")
		if len(nameSections) == 0 {
			// IDK what to do in this scenario
			continue
		}

		name := strings.Join(nameSections[1:], "/")
		// Double encoding needed
		escapedName := url.PathEscape(url.PathEscape(name))
		repository := Repository{
			Name:           escapedName,
			ArtifactsCount: r.ArtifactCount,
			Project:        project,
		}
		repositories[i] = repository
	}

	rows := make([]table.Row, len(*r))
	for i, r := range repositories {
		rows[i] = r.ToColumn()
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(21),
	)

	t.SetStyles(GetTableDefaultStyles())

	state := RepositoriesState{
		table: t,
		data:  repositories,
	}

	slog.Debug("New Repository state created.")

	return state
}
