package azure

import (
	"context"
	"encoding/base64"
	"github.com/fluxcd/go-git-providers/gitprovider"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestGet(t *testing.T) {

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
	require.Nil(t, err)

	org, err := client.Organizations().Get(context.Background(), gitprovider.OrganizationRef{
		Domain:       gitproviderDomain,
		Organization: "efernandezbreis",
	})
	require.Nil(t, err)
	require.NotNil(t, org)
}
