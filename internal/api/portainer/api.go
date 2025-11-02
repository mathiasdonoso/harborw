package portainer

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

type portainerApiClient struct {
	client   *http.Client
	baseUrl  string
	username string
	password string
	token    string
}

func NewPortainerApiClient(client *http.Client) (portainerApiClient, error) {
	slog.Debug("Creating a new portainer api client")
	username := os.Getenv("LDAP_USERNAME")
	password := os.Getenv("LDAP_PASSWORD")
	baseUrl := os.Getenv("PORTAINER_BASEURL")

	if username == "" || password == "" || baseUrl == "" {
		return portainerApiClient{}, fmt.Errorf("credentials cannot be empty")
	}

	slog.Debug(fmt.Sprintf("New portainer api client created for %s", baseUrl))

	return portainerApiClient{
		client:   client,
		baseUrl:  baseUrl,
		username: username,
		password: password,
		token:    "",
	}, nil
}
