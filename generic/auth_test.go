package generic

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestNewClient(t *testing.T) {

	os.Setenv("GIT_KIND", "gitea")
	os.Setenv("GIT_SERVER", "http://localhost:3000/")
	//os.Setenv("GIT_TOKEN", "12321")
	os.Setenv("GIT_USER", "gitea")

	client, err := NewClientFromEnvironment()

	require.NoError(t, err)
	require.NotNil(t, client)
	require.NotNil(t, client.SupportedDomain())
	require.NotNil(t, client.Organizations())
	require.NotNil(t, client.UserRepositories())
	require.NotNil(t, client.Organizations())
}
