package portainer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

type AuthResult struct {
	Jwt string `json:"jwt"`
}

type AuthRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (p *portainerApiClient) PostAuth() error {
	slog.Debug("Authenticating using portainer api")
	url := fmt.Sprintf("%s/api/auth", p.baseUrl)

	reqBody := AuthRequestBody{
		Username: p.username,
		Password: p.password,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "harborw/1.0")

	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	var authResp AuthResult
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	p.token = authResp.Jwt
	slog.Debug("Authenticated from portainer api")

	return nil
}
