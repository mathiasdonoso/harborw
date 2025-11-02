package portainer

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

type Tlsconfig struct {
	Tls           bool   `json:"TLS"`
	Tlsskipverify bool   `json:"TLSSkipVerify"`
	Tlscacert     string `json:"TLSCACert"`
	Tlscert       string `json:"TLSCert"`
	Tlskey        string `json:"TLSKey"`
}

type Azurecredentials struct {
	Applicationid     string `json:"ApplicationID"`
	Tenantid          string `json:"TenantID"`
	Authenticationkey string `json:"AuthenticationKey"`
}

type EndpointsResult struct {
	Id               int              `json:"Id"`
	Name             string           `json:"Name"`
	Type             int              `json:"Type"`
	Url              string           `json:"URL"`
	Groupid          int              `json:"GroupId"`
	Publicurl        string           `json:"PublicURL"`
	Tlsconfig        Tlsconfig        `json:"TLSConfig"`
	Authorizedusers  []interface{}    `json:"AuthorizedUsers"`
	Authorizedteams  []interface{}    `json:"AuthorizedTeams"`
	Extensions       []interface{}    `json:"Extensions"`
	Azurecredentials Azurecredentials `json:"AzureCredentials"`
	Tags             any              `json:"Tags"`
}

func (p *portainerApiClient) GetEndpoints() (*[]EndpointsResult, error) {
	slog.Debug("Fetching endpoints from portainer api")
	url := fmt.Sprintf("%s/api/endpoints", p.baseUrl)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.token))
	req.Header.Set("User-Agent", "harborw/1.0")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code: %d", resp.StatusCode)
	}

	var endpointsResp []EndpointsResult
	if err := json.NewDecoder(resp.Body).Decode(&endpointsResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	slog.Debug("Endpoints from portainer api", fmt.Sprintf("%+v", endpointsResp))

	return &endpointsResp, nil
}
