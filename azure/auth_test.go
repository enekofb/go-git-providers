package azure

import (
	"encoding/base64"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestNewClient(t *testing.T) {
	var token string
	rawToken := os.Getenv("AZURE_DEVOPS_TOKEN")
	if rawToken != "" {
		token = base64.StdEncoding.EncodeToString([]byte(":" + rawToken))
	}

	client, err := NewClient(ClientOptions{
		Org:     "efernandezbreis",
		Project: "weaveworks",
		Token:   token,
	},
	)
	require.NoError(t, err)
	require.NotNil(t, client)
	require.Equal(t, client.SupportedDomain(), "https://dev.azure.com")
	require.NotNil(t, client.Organizations())
	require.NotNil(t, client.UserRepositories())
	require.NotNil(t, client.Organizations())
}
