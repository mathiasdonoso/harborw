package harbor

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

type harborApiClient struct {
	client   *http.Client
	baseUrl  string
	username string
	password string
}

func NewHarborApiClient(client *http.Client) (harborApiClient, error) {
	slog.Debug("Creating a new harbor api client")
	username := os.Getenv("LDAP_USERNAME")
	password := os.Getenv("LDAP_PASSWORD")
	baseUrl := os.Getenv("HARBOR_BASEURL")

	if username == "" || password == "" || baseUrl == "" {
		return harborApiClient{}, fmt.Errorf("credentials cannot be empty")
	}

	slog.Debug(fmt.Sprintf("New harbor api client created for %s", baseUrl))

	return harborApiClient{
		client,
		baseUrl,
		username,
		password,
	}, nil
}
