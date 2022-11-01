package generic

import (
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
	}{
		{"gitea", "http://localhost:3000", "GITEA_TOKEN", "gitea"},
		//{"azure", "https://dev.azure.com", "AZURE_DEVOPS_TOKEN", "efernandezbreis"},
		{"bitbucketcloud", "", "", "enekoww"},
	}

	for _, gitProvider := range gitProviders {
		t.Run(gitProvider.kind, func(t *testing.T) {
			os.Setenv("GIT_KIND", gitProvider.kind)
			os.Setenv("GIT_SERVER", gitProvider.server)
			if gitProvider.tokenEnvVar != "" {
				os.Setenv("GIT_TOKEN", os.Getenv(gitProvider.tokenEnvVar))
			}
			os.Setenv("GIT_USER", gitProvider.user)

			client, err := NewClientFromEnvironment()

			require.NoError(t, err)
			require.NotNil(t, client)
			require.NotNil(t, client.SupportedDomain())
			require.NotNil(t, client.Organizations())
			require.NotNil(t, client.UserRepositories())
			require.NotNil(t, client.Organizations())
		})
	}
}
