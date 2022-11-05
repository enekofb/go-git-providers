package generic

import (
	"context"
	"github.com/drone/go-scm/scm"
	"github.com/fluxcd/go-git-providers/gitprovider"
	"github.com/stretchr/testify/require"
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
		//{"gitea", "http://localhost:3000", "GITEA_TOKEN", "gitea", "gitea/weaveworks"},
		//{"bitbucketcloud", "", "GITEA_TOKEN", "enekoww", "enekoww/test"},
	}

	for _, gitProvider := range gitProviders {
		t.Run(gitProvider.kind, func(t *testing.T) {
			os.Setenv("GIT_KIND", gitProvider.kind)
			os.Setenv("GIT_SERVER", gitProvider.server)
			if gitProvider.tokenEnvVar != "" {
				os.Setenv("GIT_TOKEN", os.Getenv(gitProvider.tokenEnvVar))
			}
			os.Setenv("GIT_USER", gitProvider.user)

			os.Setenv("GIT_REPO", gitProvider.repo)

			var c gitprovider.Client
			var err error

			c, err = NewClientFromEnvironmentDrone()

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
			require.Equal(t, gitProvider.repo, repository.Name)

		})
	}
}
