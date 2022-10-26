package azure

import (
	"context"
	"github.com/fluxcd/go-git-providers/gitprovider"
	"github.com/jenkins-x/go-scm/scm"
	"github.com/jenkins-x/go-scm/scm/driver/azure"
)

const domain = "https://dev.azure.com"

type ClientOptions struct {
	org     string
	project string
}

type wrapper struct {
	client *scm.Client
}

func NewClient(clientOptions ClientOptions) (gitprovider.Client, error) {

	project := clientOptions.project
	org := clientOptions.org

	c, err := azure.New(domain, org, project)

	return wrapper{client: c}, err
}

// SupportedDomain returns the domain endpoint for this client, e.g. "github.com", "enterprise.github.com" or
// "my-custom-git-server.com:6443". This allows a higher-level user to know what Client to use for
// what endpoints.
// This field is set at client creation time, and can't be changed.
func (c wrapper) SupportedDomain() string {
	return domain
}

// ProviderID returns the provider ID "github".
// This field is set at client creation time, and can't be changed.
func (c wrapper) ProviderID() gitprovider.ProviderID {
	return ""
}

// Raw returns the Go GitHub client (github.com/google/go-github/v47/github *Client)
// used under the hood for accessing GitHub.
func (c wrapper) Raw() interface{} {
	return nil
}

// Organizations returns the OrganizationsClient handling sets of organizations.
func (c wrapper) Organizations() gitprovider.OrganizationsClient {
	return nil
}

// OrgRepositories returns the OrgRepositoriesClient handling sets of repositories in an organization.
func (c wrapper) OrgRepositories() gitprovider.OrgRepositoriesClient {
	return nil
}

// UserRepositories returns the UserRepositoriesClient handling sets of repositories for a user.
func (c wrapper) UserRepositories() gitprovider.UserRepositoriesClient {
	return nil
}

// HasTokenPermission returns true if the given token has the given permissions.
func (c wrapper) HasTokenPermission(ctx context.Context, permission gitprovider.TokenPermission) (bool, error) {
	return false, nil
}
