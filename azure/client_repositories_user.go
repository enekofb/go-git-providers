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

// UserRepositoriesClient implements the gitprovider.UserRepositoriesClient interface.
var _ gitprovider.UserRepositoriesClient = &UserRepositoriesClient{}

// UserRepositoriesClient operates on repositories the user has access to.
type UserRepositoriesClient struct {
	client *scm.Client
}

type AzureUserRepository struct {
	repository *scm.Repository
}

func (a AzureUserRepository) APIObject() interface{} {

	//TODO implement me
	panic("implement me")
}

func (a AzureUserRepository) Update(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (a AzureUserRepository) Reconcile(ctx context.Context) (actionTaken bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (a AzureUserRepository) Delete(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (a AzureUserRepository) Repository() gitprovider.RepositoryRef {
	//TODO implement me
	panic("implement me")
}

func (a AzureUserRepository) Get() gitprovider.RepositoryInfo {
	//TODO implement me
	panic("implement me")
}

func (a AzureUserRepository) Set(info gitprovider.RepositoryInfo) error {
	//TODO implement me
	panic("implement me")
}

func (a AzureUserRepository) DeployKeys() gitprovider.DeployKeyClient {
	//TODO implement me
	panic("implement me")
}

func (a AzureUserRepository) Commits() gitprovider.CommitClient {
	//TODO implement me
	panic("implement me")
}

func (a AzureUserRepository) Branches() gitprovider.BranchClient {
	//TODO implement me
	panic("implement me")
}

func (a AzureUserRepository) PullRequests() gitprovider.PullRequestClient {
	//TODO implement me
	panic("implement me")
}

func (a AzureUserRepository) Files() gitprovider.FileClient {
	//TODO implement me
	panic("implement me")
}

func (a AzureUserRepository) Trees() gitprovider.TreeClient {
	//TODO implement me
	panic("implement me")
}

// Get returns the repository at the given path.
//
// ErrNotFound is returned if the resource does not exist.
func (c *UserRepositoriesClient) Get(ctx context.Context, ref gitprovider.UserRepositoryRef) (gitprovider.UserRepository, error) {

	repository, _, err := c.client.Repositories.Find(ctx, ref.GetRepository())
	if err != nil {
		return nil, err
	}

	return AzureUserRepository{repository: repository}, nil
}

// List all repositories in the given organization.
//
// List returns all available repositories, using multiple paginated requests if needed.
func (c *UserRepositoriesClient) List(ctx context.Context, ref gitprovider.UserRef) ([]gitprovider.UserRepository, error) {
	return nil, errors.New("not implemented")
}

// Create creates a repository for the given organization, with the data and options
//
// ErrAlreadyExists will be returned if the resource already exists.
func (c *UserRepositoriesClient) Create(ctx context.Context,
	ref gitprovider.UserRepositoryRef,
	req gitprovider.RepositoryInfo,
	opts ...gitprovider.RepositoryCreateOption,
) (gitprovider.UserRepository, error) {
	return nil, errors.New("not implemented")
}

// Reconcile makes sure the given desired state (req) becomes the actual state in the backing Git provider.
//
// If req doesn't exist under the hood, it is created (actionTaken == true).
// If req doesn't equal the actual state, the resource will be updated (actionTaken == true).
// If req is already the actual state, this is a no-op (actionTaken == false).
func (c *UserRepositoriesClient) Reconcile(ctx context.Context, ref gitprovider.UserRepositoryRef, req gitprovider.RepositoryInfo, opts ...gitprovider.RepositoryReconcileOption) (gitprovider.UserRepository, bool, error) {
	return nil, false, errors.New("not implemented")
}
