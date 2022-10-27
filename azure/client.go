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

type AzureCommitClient struct {
	repository UserRepository
}

func (a UserRepository) Commits() gitprovider.CommitClient {
	return a
}

func (a UserRepository) Branches() gitprovider.BranchClient {
	return UserRepositoryBranches{
		repository: a,
	}
}

type AzurePullReqeuestClient struct {
	repository UserRepository
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

func (a UserRepository) PullRequests() gitprovider.PullRequestClient {
	return AzurePullReqeuestClient{
		repository: a,
	}
}

func (a UserRepository) Files() gitprovider.FileClient {
	//TODO implement me
	panic("implement me")
}

func (a UserRepository) Trees() gitprovider.TreeClient {
	//TODO implement me
	panic("implement me")
}
