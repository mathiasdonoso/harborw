package harbor

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

type BuildHistory struct {
	Absolute bool   `json:"absolute"`
	Href     string `json:"href"`
}

type AdditionLinks struct {
	BuildHistory BuildHistory `json:"build_history"`
}

type N80Tcp struct {
}

type N8080Tcp struct {
}

type N8443Tcp struct {
}

type N9001Tcp struct {
}

type Exposedports struct {
	N80Tcp   N80Tcp   `json:"80/tcp"`
	N8080Tcp N8080Tcp `json:"8080/tcp"`
	N8443Tcp N8443Tcp `json:"8443/tcp"`
	N9001Tcp N9001Tcp `json:"9001/tcp"`
}

type Labels struct {
	DevopsService string `json:"devops-service"`
	Maintainer    string `json:"maintainer"`
	Version       string `json:"version"`
}

type Config struct {
	Cmd          []any        `json:"Cmd"`
	Entrypoint   []any        `json:"Entrypoint"`
	Env          []any        `json:"Env"`
	Exposedports Exposedports `json:"ExposedPorts"`
	Labels       Labels       `json:"Labels"`
	Stopsignal   string       `json:"StopSignal"`
	User         string       `json:"User"`
	Workingdir   string       `json:"WorkingDir"`
}

type ExtraAttrs struct {
	Architecture string `json:"architecture"`
	Author       string `json:"author"`
	Config       Config `json:"config"`
	Created      string `json:"created"`
	Os           string `json:"os"`
}

type Tag struct {
	ArtifactId   int    `json:"artifact_id"`
	Id           int    `json:"id"`
	Immutable    bool   `json:"immutable"`
	Name         string `json:"name"`
	PullTime     string `json:"pull_time"`
	PushTime     string `json:"push_time"`
	RepositoryId int    `json:"repository_id"`
	Signed       bool   `json:"signed"`
}

type ArtifactsResult struct {
	AdditionLinks     AdditionLinks `json:"addition_links"`
	Digest            string        `json:"digest"`
	ExtraAttrs        ExtraAttrs    `json:"extra_attrs"`
	Icon              string        `json:"icon"`
	Id                int           `json:"id"`
	Labels            any           `json:"labels"`
	ManifestMediaType string        `json:"manifest_media_type"`
	MediaType         string        `json:"media_type"`
	ProjectId         int           `json:"project_id"`
	PullTime          string        `json:"pull_time"`
	PushTime          string        `json:"push_time"`
	References        any           `json:"references"`
	RepositoryId      int           `json:"repository_id"`
	Size              int           `json:"size"`
	Tags              []Tag         `json:"tags"`
	Type              string        `json:"type"`
}

func (h harborApiClient) FetchArtifacts(project string, repository string) (*[]ArtifactsResult, error) {
	url := fmt.Sprintf("%s/api/v2.0/projects/%s/repositories/%s/artifacts", h.baseUrl, project, repository)
	slog.Debug(fmt.Sprintf("Fetching artifacts. URL: %s", url))
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", h.credentials))
	req.Header.Set("User-Agent", "harborw/1.0")

	q := req.URL.Query()
	q.Add("page", "1")
	q.Add("page_size", "100")
	req.URL.RawQuery = q.Encode()

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code: %d", resp.StatusCode)
	}

	var artifactsResp []ArtifactsResult
	if err := json.NewDecoder(resp.Body).Decode(&artifactsResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	slog.Debug("Artifacts fetched", fmt.Sprintf("%+v", artifactsResp))

	return &artifactsResp, nil
}

func (h harborApiClient) DeleteArtifact(project string, repository string, artifactHashOrTag string) error {
	url := fmt.Sprintf("%s/api/v2.0/projects/%s/repositories/%s/artifacts/%s", h.baseUrl, project, repository, artifactHashOrTag)
	slog.Debug(fmt.Sprintf("Deleting artifact. URL: %s", url))
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", h.credentials))
	req.Header.Set("User-Agent", "harborw/1.0")

	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("status code: %d", resp.StatusCode)
	}

	slog.Debug(fmt.Sprintf("Artifact with hash: %s deleted", artifactHashOrTag[:10]))

	return nil
}
