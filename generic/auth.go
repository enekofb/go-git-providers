package generic

import (
	"context"
	"encoding/base64"
	"fmt"
	drone "github.com/drone/go-scm/scm"
	droneAzure "github.com/drone/go-scm/scm/driver/azure"
	"github.com/drone/go-scm/scm/transport"
	"net/http"

	"github.com/fluxcd/go-git-providers/gitprovider"
	"os"
)

type ClientOptions struct {
	Token string
	Uri   string
}

type wrapperDrone struct {
	client *drone.Client
}

func NewClientFromScm(c *drone.Client) (gitprovider.Client, error) {

	return wrapperDrone{client: c}, nil
}

func NewClientFromEnvironmentDrone() (wrapperDrone, error) {
	var c *drone.Client
	driver := os.Getenv("GIT_KIND")
	//serverURL := os.Getenv("GIT_SERVER")
	token := os.Getenv("GIT_TOKEN")
	username := os.Getenv("GIT_USER")
	repo := os.Getenv("GIT_REPO")

	//clientID := os.Getenv("BB_OAUTH_CLIENT_ID")
	//clientSecret := os.Getenv("BB_OAUTH_CLIENT_SECRET")
	switch driver {
	case "azure":
		c = droneAzure.NewDefault(username, repo)

		var encodedToken string
		if token != "" {
			encodedToken = base64.StdEncoding.EncodeToString([]byte(":" + token))
		}

		c.Client = &http.Client{
			Transport: &transport.Custom{
				Before: func(r *http.Request) {
					r.Header.Set("Authorization", fmt.Sprintf("Basic %s", encodedToken))
				},
			},
		}

	default:
		return wrapperDrone{}, fmt.Errorf("Unsupported GIT_KIND value: %s", driver)
	}

	return wrapperDrone{client: c}, nil

}

// SupportedDomain returns the domain endpoint for this client, e.g. "github.com", "enterprise.github.com" or
// "my-custom-git-server.com:6443". This allows a higher-level user to know what Client to use for
// what endpoints.
// This field is set at client creation time, and can't be changed.
func (c wrapperDrone) SupportedDomain() string {
	return ""
}

// ProviderID returns the provider ID "github".
// This field is set at client creation time, and can't be changed.
func (c wrapperDrone) ProviderID() gitprovider.ProviderID {
	return ""
}

// Raw returns the Go GitHub client (github.com/google/go-github/v47/github *Client)
// used under the hood for accessing GitHub.
func (c wrapperDrone) Raw() interface{} {
	return nil
}

// Organizations returns the Organisations handling sets of organizations.
func (c wrapperDrone) Organizations() gitprovider.OrganizationsClient {
	return &Organisations{client: c.client}
}

// OrgRepositories returns the OrgRepositoriesClient handling sets of repositories in an organization.
func (c wrapperDrone) OrgRepositories() gitprovider.OrgRepositoriesClient {
	return &OrgRepositories{client: c.client}
}

// UserRepositories returns the UserRepositories handling sets of repositories for a user.
func (c wrapperDrone) UserRepositories() gitprovider.UserRepositoriesClient {
	return &UserRepositories{client: c.client}
}

// HasTokenPermission returns true if the given Token has the given permissions.
func (c wrapperDrone) HasTokenPermission(ctx context.Context, permission gitprovider.TokenPermission) (bool, error) {
	return false, nil
}
