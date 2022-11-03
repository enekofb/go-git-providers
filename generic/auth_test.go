package generic

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/fluxcd/go-git-providers/gitprovider"
	"github.com/jenkins-x/go-scm/scm"
	"github.com/jenkins-x/go-scm/scm/driver/azure"
	"github.com/jenkins-x/go-scm/scm/transport"
	"github.com/stretchr/testify/require"
	"net/http"
	"os"
	"testing"
)

func TestNewClient(t *testing.T) {
	var gitProviders = []struct {
		kind        string
		server      string
		tokenEnvVar string
		user        string
		repo        string
	}{
		{"azure", "https://dev.azure.com", "AZURE_DEVOPS_TOKEN", "efernandezbreis", "weaveworks"},
		{"gitea", "http://localhost:3000", "GITEA_TOKEN", "gitea", "gitea/weaveworks"},
		{"bitbucketcloud", "", "GITEA_TOKEN", "enekoww", "enekoww/test"},
	}

	for _, gitProvider := range gitProviders {
		t.Run(gitProvider.kind, func(t *testing.T) {
			os.Setenv("GIT_KIND", gitProvider.kind)
			os.Setenv("GIT_SERVER", gitProvider.server)
			if gitProvider.tokenEnvVar != "" {
				os.Setenv("GIT_TOKEN", os.Getenv(gitProvider.tokenEnvVar))
			}
			os.Setenv("GIT_USER", gitProvider.user)

			var c gitprovider.Client
			var err error

			//TODO inconsistent api for creating client between azure and generic
			switch gitProvider.kind {
			case "azure":
				c, err = createAzureClient(gitProvider)
			default:
				c, err = NewClientFromEnvironment()
			}

			require.NoError(t, err)
			require.NotNil(t, c)
			require.NotNil(t, c.SupportedDomain())
			require.NotNil(t, c.Organizations())
			require.NotNil(t, c.UserRepositories())
			require.NotNil(t, c.Organizations())

			// get repo
			ctx := context.Background()
			userRepoRef := newUserRepoRef(gitProvider.user, gitProvider.repo)
			var userRepo gitprovider.UserRepository
			userRepo, err = c.UserRepositories().Get(ctx, userRepoRef)
			require.NoError(t, err)

			//var repositoryName string

			//TODO had to move to scm due to limitations on git providers getting a consistent name or
			//implementation limitation
			object := userRepo.APIObject()
			repository := object.(*scm.Repository)

			//TODO inconsistent api for naming between generic and azure
			switch gitProvider.kind {
			case "azure":
				require.Equal(t, gitProvider.repo, repository.Name)
			default:
				require.Equal(t, gitProvider.repo, repository.FullName)
			}
		})
	}
}

func createAzureClient(provider struct {
	kind        string
	server      string
	tokenEnvVar string
	user        string
	repo        string
}) (gitprovider.Client, error) {
	rawToken := os.Getenv(provider.tokenEnvVar)
	if len(rawToken) == 0 {
		return nil, fmt.Errorf("couldn't acquire AZURE_DEVOPS_TOKEN env variable")
	}
	var token string
	if rawToken != "" {
		token = base64.StdEncoding.EncodeToString([]byte(":" + rawToken))
	}

	c, _ := azure.New("https://dev.azure.com", provider.user, provider.repo)

	c.Client = &http.Client{
		Transport: &transport.Custom{
			Before: func(r *http.Request) {
				r.Header.Set("Authorization", fmt.Sprintf("Basic %s", token))
			},
		},
	}

	return NewClientFromScm(c)
}
