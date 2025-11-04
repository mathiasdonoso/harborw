package harbor

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

type CveAllowlist struct {
	CreationTime string `json:"creation_time"`
	Id           int    `json:"id"`
	Items        []any  `json:"items"`
	ProjectId    int    `json:"project_id"`
	UpdateTime   string `json:"update_time"`
}

type Metadata struct {
	AutoScan             string `json:"auto_scan"`
	EnableContentTrust   string `json:"enable_content_trust"`
	PreventVul           string `json:"prevent_vul"`
	Public               string `json:"public"`
	ReuseSysCveAllowlist string `json:"reuse_sys_cve_allowlist"`
	Severity             string `json:"severity"`
}

type ProjectsResult struct {
	ChartCount         int          `json:"chart_count"`
	CreationTime       string       `json:"creation_time"`
	CurrentUserRoleId  int          `json:"current_user_role_id"`
	CurrentUserRoleIds []any        `json:"current_user_role_ids"`
	CveAllowlist       CveAllowlist `json:"cve_allowlist"`
	Metadata           Metadata     `json:"metadata"`
	Name               string       `json:"name"`
	OwnerId            int          `json:"owner_id"`
	OwnerName          string       `json:"owner_name"`
	ProjectId          int          `json:"project_id"`
	RepoCount          int          `json:"repo_count"`
	UpdateTime         string       `json:"update_time"`
}

func (h harborApiClient) FetchProjects() (*[]ProjectsResult, error) {
	url := fmt.Sprintf("%s/api/v2.0/projects", h.baseUrl)
	slog.Debug(fmt.Sprintf("Fetching projects. URL: %s", url))
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

	var projectsResp []ProjectsResult
	if err := json.NewDecoder(resp.Body).Decode(&projectsResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	slog.Debug("Projects fetched", "data", fmt.Sprintf("%+v", projectsResp))

	return &projectsResp, nil
}
