package azure

import (
	"context"
	"errors"
	"fmt"
	"github.com/fluxcd/go-git-providers/gitprovider"
	"github.com/jenkins-x/go-scm/scm"
	"net/http"
)

type UserRepositoryPullRequests struct {
	repository UserRepository
}

func (a UserRepositoryPullRequests) List(ctx context.Context) ([]gitprovider.PullRequest, error) {
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

func (a UserRepositoryPullRequests) Create(ctx context.Context, title, branch, baseBranch, description string) (gitprovider.PullRequest, error) {

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

func (a UserRepositoryPullRequests) Get(ctx context.Context, number int) (gitprovider.PullRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (a UserRepositoryPullRequests) Merge(ctx context.Context, number int, mergeMethod gitprovider.MergeMethod, message string) error {
	//TODO implement me
	panic("implement me")
}

func (a UserRepository) PullRequests() gitprovider.PullRequestClient {
	return UserRepositoryPullRequests{
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

type OrgRepositoryPullRequests struct {
	repository OrgRepository
}

func (o OrgRepositoryPullRequests) List(ctx context.Context) ([]gitprovider.PullRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (o OrgRepositoryPullRequests) Create(ctx context.Context, title, branch, baseBranch, description string) (gitprovider.PullRequest, error) {
	input := &scm.PullRequestInput{
		Title: title,
		Body:  description,
		Head:  branch,
		Base:  baseBranch,
	}

	outputPR, response, err := o.repository.client.PullRequests.Create(context.Background(), o.repository.repository.ID, input)
	if err != nil {
		return nil, err
	}
	if response.Status != http.StatusCreated {
		return nil, errors.New(fmt.Sprintf("PullRequests.Create did not get a 201 back %v", response.Status))
	}

	return AzurePullRequest{pullRequest: outputPR}, nil
}

func (o OrgRepositoryPullRequests) Get(ctx context.Context, number int) (gitprovider.PullRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (o OrgRepositoryPullRequests) Merge(ctx context.Context, number int, mergeMethod gitprovider.MergeMethod, message string) error {
	//TODO implement me
	panic("implement me")
}
