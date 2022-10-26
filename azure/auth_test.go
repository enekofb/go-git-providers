package azure

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewClient(t *testing.T) {
	client, err := NewClient(ClientOptions{
		org:     "efernandezbreis",
		project: "weaveworks",
	},
	)
	require.NoError(t, err)
	require.NotNil(t, client)
	require.Equal(t, client.SupportedDomain(), "https://dev.azure.com")
	require.NotNil(t, client.Organizations())
}
