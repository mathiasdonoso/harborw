package harbor

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

type RepositoriesResult struct {
	ArtifactCount int    `json:"artifact_count"`
	CreationTime  string `json:"creation_time"`
	Id            int    `json:"id"`
	Name          string `json:"name"`
	ProjectId     int    `json:"project_id"`
	PullCount     int    `json:"pull_count"`
	UpdateTime    string `json:"update_time"`
}

func (h harborApiClient) FetchRepositories(project string) (*[]RepositoriesResult, error) {
	slog.Debug(fmt.Sprintf("Fetching repositories from project %s", project))
	url := fmt.Sprintf("%s/api/v2.0/projects/%s/repositories", h.baseUrl, project)
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

	var repositoriesResp []RepositoriesResult
	if err := json.NewDecoder(resp.Body).Decode(&repositoriesResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	slog.Debug("Repositories fetched", fmt.Sprintf("%+v", repositoriesResp))

	return &repositoriesResp, nil
}
