package azure

import (
	"fmt"
	"github.com/fluxcd/go-git-providers/generic"
	"github.com/fluxcd/go-git-providers/gitprovider"
	"github.com/jenkins-x/go-scm/scm/driver/azure"
	"github.com/jenkins-x/go-scm/scm/transport"
	"net/http"
)

const (
	domain = "https://dev.azure.com"
)

type ClientOptions struct {
	Org     string
	Project string
	Token   string
}

func NewClient(clientOptions ClientOptions) (gitprovider.Client, error) {

	project := clientOptions.Project
	org := clientOptions.Org

	c, _ := azure.New(domain, org, project)

	c.Client = &http.Client{
		Transport: &transport.Custom{
			Before: func(r *http.Request) {
				r.Header.Set("Authorization", fmt.Sprintf("Basic %s", clientOptions.Token))
			},
		},
	}

	return generic.NewClientFromScm(c)
}
