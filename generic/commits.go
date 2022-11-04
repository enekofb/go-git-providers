package generic

import (
	"context"
	"errors"
	"fmt"
	drone "github.com/drone/go-scm/scm"
	"github.com/fluxcd/go-git-providers/gitprovider"
	"net/http"
)

type CommitReference struct {
	reference *drone.Reference
}

func (a CommitReference) APIObject() interface{} {
	return a.reference
}

func (a CommitReference) Get() gitprovider.CommitInfo {
	return gitprovider.CommitInfo{
		Sha: a.reference.Sha,
	}
}

func (a UserRepository) Create(ctx context.Context, branch string, message string, files []gitprovider.CommitFile) (gitprovider.Commit, error) {
	//TODO find a better way to get latest commit sha
	page, err := a.Commits().ListPage(ctx, branch, 10, 1)
	if err != nil {
		return nil, err
	}
	currentCommit := page[0].Get().Sha

	for _, file := range files {
		data := *file.Content
		path := *file.Path

		createParams := drone.ContentParams{
			Message: message,
			Data:    []byte(data),
			Branch:  branch,
			Ref:     currentCommit,
		}
		response, err := a.client.Contents.Create(ctx, a.repositoryId, path, &createParams)
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
	page, err := o.Commits().ListPage(ctx, branch, 10, 1)
	if err != nil {
		return nil, err
	}
	currentCommit := page[0]

	for _, file := range files {
		data := *file.Content
		path := *file.Path

		createParams := drone.ContentParams{
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
		page, err = o.Commits().ListPage(ctx, branch, 10, 1)
		if err != nil {
			return nil, err
		}
		currentCommit = page[0]
	}

	return currentCommit, nil

}
