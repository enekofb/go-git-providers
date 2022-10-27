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
	"fmt"
	"github.com/jenkins-x/go-scm/scm"
	"net/http"

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
	client     *scm.Client
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
	return gitprovider.RepositoryInfo{
		Description:   &(a.repository.FullName),
		DefaultBranch: &(a.repository.Branch),
		Visibility:    nil,
	}
}

func (a AzureUserRepository) Set(info gitprovider.RepositoryInfo) error {
	//TODO implement me
	panic("implement me")
}

func (a AzureUserRepository) DeployKeys() gitprovider.DeployKeyClient {
	//TODO implement me
	panic("implement me")
}

type AzureCommitClient struct {
	repository AzureUserRepository
}

type AzureBranchClient struct {
	repository AzureUserRepository
}

type AzureCommit struct {
	fileEntry *scm.FileEntry
}

func (a AzureCommit) APIObject() interface{} {
	return a.fileEntry
}

func (a AzureCommit) Get() gitprovider.CommitInfo {
	return gitprovider.CommitInfo{
		Sha: a.fileEntry.Sha,
	}
}

// TODO: pagination not supported
func (a AzureUserRepository) ListPage(ctx context.Context, branch string, perPage int, page int) ([]gitprovider.Commit, error) {
	commit, _, err := a.client.Contents.List(ctx, a.repository.ID, "", branch)
	if err != nil {
		return nil, err
	}

	commits := make([]gitprovider.Commit, 0, len(commit))

	for _, fe := range commit {

		commits = append(commits, AzureCommit{fileEntry: fe})

	}
	return commits, nil
}

func (a AzureUserRepository) Create(ctx context.Context, branch string, message string, files []gitprovider.CommitFile) (gitprovider.Commit, error) {
	//TODO find a better way to get latest commit sha
	page, err := a.ListPage(ctx, branch, 10, 1)
	if err != nil {
		return nil, err
	}
	currentCommit := page[0].Get().Sha

	for _, file := range files {
		data := *file.Content
		path := *file.Path

		createParams := scm.ContentParams{
			Message: message,
			Data:    []byte(data),
			Branch:  branch,
			Ref:     currentCommit,
		}
		response, err := a.client.Contents.Create(ctx, a.repository.ID, path, &createParams)
		if err != nil {
			return nil, err
		}
		if response.Status != http.StatusCreated {
			return nil, errors.New(fmt.Sprintf("create commit did not get a 200 back %v", response.Status))
		}

	}

	return nil, nil
}

func (a AzureUserRepository) Commits() gitprovider.CommitClient {
	return a
}

func (a AzureUserRepository) Branches() gitprovider.BranchClient {
	return AzureBranchClient{
		repository: a,
	}
}

func (a AzureBranchClient) Create(ctx context.Context, branch, sha string) error {

	input := &scm.ReferenceInput{
		Name: branch,
		Sha:  sha,
	}
	_, response, err := a.repository.client.Git.CreateRef(ctx, a.repository.repository.ID, input.Name, input.Sha)
	if err != nil {
		return err
	}

	if response.Status != http.StatusOK {
		return errors.New(fmt.Sprintf("CreateBranch did not get a 200 back %v", response.Status))
	}

	return nil
}

type AzurePullReqeuestClient struct {
	repository AzureUserRepository
}

func (a AzurePullReqeuestClient) List(ctx context.Context) ([]gitprovider.PullRequest, error) {
	//TODO implement me
	panic("implement me")
}

type AzurePullRequest struct {
	pullRequest *scm.PullRequest
}

func (a AzurePullRequest) APIObject() interface{} {
	return a.pullRequest
}

func (a AzurePullRequest) Get() gitprovider.PullRequestInfo {

	return gitprovider.PullRequestInfo{
		Number: a.pullRequest.Number,
		WebURL: a.pullRequest.Link,
		Merged: a.pullRequest.Merged,
	}
}

func (a AzurePullReqeuestClient) Create(ctx context.Context, title, branch, baseBranch, description string) (gitprovider.PullRequest, error) {

	input := &scm.PullRequestInput{
		Title: title,
		Body:  description,
		Head:  branch,
		Base:  baseBranch,
	}

	outputPR, response, err := a.repository.client.PullRequests.Create(context.Background(), a.repository.repository.ID, input)
	if err != nil {
		return nil, err
	}
	if response.Status != http.StatusCreated {
		return nil, errors.New(fmt.Sprintf("PullRequests.Create did not get a 201 back %v", response.Status))
	}

	return AzurePullRequest{pullRequest: outputPR}, nil
}

func (a AzurePullReqeuestClient) Get(ctx context.Context, number int) (gitprovider.PullRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (a AzurePullReqeuestClient) Merge(ctx context.Context, number int, mergeMethod gitprovider.MergeMethod, message string) error {
	//TODO implement me
	panic("implement me")
}

func (a AzureUserRepository) PullRequests() gitprovider.PullRequestClient {
	return AzurePullReqeuestClient{
		repository: a,
	}
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

	return AzureUserRepository{
		repository: repository,
		client:     c.client,
	}, nil
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
