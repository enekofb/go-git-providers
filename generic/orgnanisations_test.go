package generic

import (
	"context"
	"github.com/fluxcd/go-git-providers/gitprovider"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGet(t *testing.T) {

	client, err := NewClientFromEnvironment()
	require.Nil(t, err)

	org, err := client.Organizations().Get(context.Background(), gitprovider.OrganizationRef{
		Domain:       "test",
		Organization: "efernandezbreis",
	})
	require.Nil(t, err)
	require.NotNil(t, org)
}
