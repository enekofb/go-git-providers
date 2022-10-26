/*
Copyright 2020 The Flux CD contributors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package azure

import (
	"context"
	"errors"
	"github.com/jenkins-x/go-scm/scm"

	"github.com/fluxcd/go-git-providers/gitprovider"
)

// OrganizationsClient implements the gitprovider.OrganizationsClient interface.
var _ gitprovider.OrganizationsClient = &OrganizationsClient{}

type OrganizationsClient struct {
	client *scm.Client
}

func (c *OrganizationsClient) Get(ctx context.Context, ref gitprovider.OrganizationRef) (gitprovider.Organization, error) {
	return nil, errors.New("not implemented")
}

// TODO: not supported by scm
// https://github.com/jenkins-x/go-scm/blob/main/scm/driver/azure/org.go
func (c *OrganizationsClient) List(ctx context.Context) ([]gitprovider.Organization, error) {
	return nil, errors.New("not implemented")
}

func (c *OrganizationsClient) Children(_ context.Context, _ gitprovider.OrganizationRef) ([]gitprovider.Organization, error) {
	return nil, errors.New("not implemented")
}
