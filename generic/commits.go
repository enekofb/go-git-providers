package generic

import (
	"context"
	"errors"
	"fmt"
	"github.com/fluxcd/go-git-providers/gitprovider"
	"github.com/jenkins-x/go-scm/scm"
	"net/http"
)

type Commmit struct {
	commit *scm.Commit
}

func (a Commmit) APIObject() interface{} {
	return a.commit
}

func (a Commmit) Get() gitprovider.CommitInfo {
	return gitprovider.CommitInfo{
		Sha: a.commit.Sha,
	}
}

// TODO: pagination not supported
func (a UserRepository) ListPage(ctx context.Context, branch string, perPage int, page int) ([]gitprovider.Commit, error) {

	findCommit, _, err := a.client.Git.FindCommit(ctx, a.repository.ID, branch)
	if err != nil {
		return nil, err
	}

	commits := make([]gitprovider.Commit, 0, 1)

	commits = append(commits, Commmit{commit: findCommit})

	return commits, nil
}

func (o OrgRepository) ListPage(ctx context.Context, branch string, perPage int, page int) ([]gitprovider.Commit, error) {
	findCommit, _, err := o.client.Git.FindCommit(ctx, o.repository.ID, branch)
	if err != nil {
		return nil, err
	}

	commits := make([]gitprovider.Commit, 0, 1)

	commits = append(commits, Commmit{commit: findCommit})

	return commits, nil
}

func (a UserRepository) Create(ctx context.Context, branch string, message string, files []gitprovider.CommitFile) (gitprovider.Commit, error) {
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

func (o OrgRepository) Create(ctx context.Context, branch string, message string, files []gitprovider.CommitFile) (gitprovider.Commit, error) {
	//TODO find a better way to get latest commit sha
	page, err := o.ListPage(ctx, branch, 10, 1)
	if err != nil {
		return nil, err
	}
	currentCommit := page[0]

	for _, file := range files {
		data := *file.Content
		path := *file.Path

		createParams := scm.ContentParams{
			Message: message,
			Data:    []byte(data),
			Branch:  branch,
			Ref:     currentCommit.Get().Sha,
		}
		response, err := o.client.Contents.Create(ctx, o.repository.ID, path, &createParams)
		if err != nil {
			return nil, err
		}
		if response.Status != http.StatusCreated {
			return nil, errors.New(fmt.Sprintf("create commit did not get a 200 back %v", response.Status))
		}
		//this aim to update current commit to allow adding more than once
		//TODO: to list the whole thing for getting a sha that we have just created seems expensive, look for a better alternative
		page, err = o.ListPage(ctx, branch, 10, 1)
		if err != nil {
			return nil, err
		}
		currentCommit = page[0]
	}

	return currentCommit, nil

}
