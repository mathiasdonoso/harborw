package tui

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mathiasdonoso/harborw/internal/api/harbor"
	"github.com/mathiasdonoso/harborw/internal/api/portainer"
)

type Artifact struct {
	Selected   bool
	Project    string
	Repository string
	Name       string
	Hash       string
	Size       float64
	PullTime   string
	PushTime   string
}

func (a Artifact) ToColumn() []string {
	checked := "[ ]"

	if a.Selected {
		checked = "[x]"
	}

	size := float64(a.Size) / 1024 / 1024

	return []string{
		checked,
		a.Name,
		a.Hash,
		"",
		fmt.Sprintf("%.2f MiB", size),
		a.PullTime,
		a.PushTime,
	}
}

type DeletionProcess struct {
	artifacts Artifact
	status    string
	deleted   bool
	err       error
}

type ArtifactsState struct {
	table table.Model
	data  []Artifact
}

type processDeleteArtifactMsg struct {
	Artifact  Artifact
	canDelete bool
}

func ImageIsInUse(hash string) (bool, error) {
	portainerClient, err := portainer.NewPortainerApiClient(http.DefaultClient)
	if err != nil {
		return false, err
	}

	slog.Debug(fmt.Sprintf("Searching for usage for image with hash: %s", hash))
	endpoints, err := portainerClient.GetEndpoints()
	if err != nil {
		return false, err
	}

	for i, e := range *endpoints {
		fmt.Printf("i: %d\n", i)
		slog.Debug(fmt.Sprintf("Searching for usage of image with hash %s inside endpoint %s", hash, e.Name))

		// containerInfo, err := portainerClient.GetContainersJson(e.Id)
		containerInfo, err := portainerClient.GetContainersJson(1009)
		if err != nil {
			return false, err
		}

		for _, c := range *containerInfo {
			slog.Debug(fmt.Sprintf("Searching for usage of image with hash: %s inside container: %s", hash, c.Image))

			// "Image": "hub.fif.tech/omnichannel/privatesite:bf-co-executive-eta@sha256:b4e2f5ad6ce67c3033317119a1044044642425b2e7c9619372eab6ae226c5e75",
			imageSections := strings.Split(c.Image, "@")
			if len(imageSections) == 2 {
				imageHash := imageSections[1]

				if hash == imageHash {
					return true, nil
				}
			}
		}
	}

	return false, nil
}

type ArtifactDeleteMsg struct {
	err error
}

func deleteArtifact(artifact Artifact) tea.Cmd {
	var err error
	harborClient, err := harbor.NewHarborApiClient(http.DefaultClient)
	if err != nil {
		return func() tea.Msg {
			return ArtifactDeleteMsg{
				err,
			}
		}
	}

	err = harborClient.DeleteArtifact(artifact.Project, artifact.Repository, artifact.Hash)
	return func() tea.Msg {
		return ArtifactDeleteMsg{
			err,
		}
	}
}

func processDeleteArtifact(artifact Artifact) tea.Cmd {
	canDelete := false

	result, err := ImageIsInUse(artifact.Hash)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
	}

	if result {
		fmt.Println("Artifact is in USE")
		canDelete = true
	} else {
		fmt.Println("Artifact is NOT in USE")
	}

	return func() tea.Msg {
		return processDeleteArtifactMsg{
			Artifact:  artifact,
			canDelete: canDelete,
		}
	}
}

func getSelectedArtifacts(artifacsState ArtifactsState) []Artifact {
	selected := make([]Artifact, 0)

	for _, a := range artifacsState.data {
		if a.Selected {
			selected = append(selected, a)
		}
	}

	return selected
}

func (m model) artifactsView() string {
	return m.state.artifacts.table.View()
}

func (m model) artifactsUpdate(msg tea.Msg) (model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case ArtifactDeleteMsg:
		fmt.Printf("DELETED")
	case processDeleteArtifactMsg:
		if msg.canDelete {
			fmt.Println("processDoneMsg!!!")
		} else {
			fmt.Println("processErrorMsg!!!")
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "c":
			// Clear selected artifacts
			rows := make([]table.Row, len(m.state.artifacts.data))

			for i := range m.state.artifacts.data {
				m.state.artifacts.data[i].Selected = false
				rows[i] = m.state.artifacts.data[i].ToColumn()
			}

			m.state.artifacts.table.SetRows(rows)
			return m, nil
		case "d":
			// Delete selected artifacts
			selected := getSelectedArtifacts(m.state.artifacts)
			cmds := []tea.Cmd{}

			for _, s := range selected {
				// cmds = append(cmds, processDeleteArtifact(s))
				cmds = append(cmds, deleteArtifact(s))
			}

			return m, tea.Batch(cmds...)
		case "-":
			// Go back
			m = m.SwitchPage(repositoriesPage)
			return m, nil
		case " ":
			// Check artifact for deletion
			rowIndex := m.state.artifacts.table.Cursor()
			m.state.artifacts.data[rowIndex].Selected = !m.state.artifacts.data[rowIndex].Selected

			rows := make([]table.Row, len(m.state.artifacts.data))
			for i, a := range m.state.artifacts.data {
				rows[i] = a.ToColumn()
			}

			m.state.artifacts.table.SetRows(rows)

			return m, nil
		}
	}
	m.state.artifacts.table, cmd = m.state.artifacts.table.Update(msg)
	return m, cmd
}

func (m model) NewArtifactsState(project string, repository string) ArtifactsState {
	columns := []table.Column{
		{Title: "Select", Width: 6},
		{Title: "Repository", Width: 25},
		{Title: "sha256", Width: 15},
		{Title: "Labels", Width: 20},
		{Title: "Size (MiB)", Width: 10},
		{Title: "Pull time", Width: 25},
		{Title: "Push time", Width: 25},
	}

	harborClient, err := harbor.NewHarborApiClient(http.DefaultClient)
	if err != nil {
		fmt.Printf("Error creating harbor client: %v\n", err)
	}

	r, err := harborClient.FetchArtifacts(project, repository)
	if err != nil {
		fmt.Printf("Error fetching projects: %v\n", err)
	}
	repositoriesLoaded = true

	artifacts := make([]Artifact, len(*r))

	for i, a := range *r {
		tag := a.Tags[0]
		// size := float64(a.Size) / 1024 / 1024
		artifact := Artifact{
			Selected:   false,
			Name:       tag.Name,
			Hash:       a.Digest,
			Size:       float64(a.Size),
			PullTime:   a.PullTime,
			PushTime:   a.PushTime,
			Project:    project,
			Repository: repository,
		}

		artifacts[i] = artifact
	}

	rows := make([]table.Row, len(*r))
	for i, a := range artifacts {
		rows[i] = a.ToColumn()
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(41),
	)

	t.SetStyles(GetTableDefaultStyles())

	state := ArtifactsState{
		table: t,
		data:  artifacts,
	}

	return state
}
