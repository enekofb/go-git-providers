package generic

import (
	"context"
	"errors"
	"fmt"
	"github.com/fluxcd/go-git-providers/gitprovider"
	"github.com/jenkins-x/go-scm/scm"
	"net/http"
)

type CommitFileReference struct {
	fileEntry *scm.FileEntry
}

type CommitReference struct {
	reference *scm.Reference
}

func (a CommitReference) APIObject() interface{} {
	return a.reference
}

func (a CommitReference) Get() gitprovider.CommitInfo {
	return gitprovider.CommitInfo{
		Sha: a.reference.Sha,
	}
}

func (a CommitFileReference) APIObject() interface{} {
	return a.fileEntry
}

func (a CommitFileReference) Get() gitprovider.CommitInfo {
	return gitprovider.CommitInfo{
		Sha: a.fileEntry.Sha,
	}
}

// TODO: pagination not supported
func (a UserRepository) ListPage(ctx context.Context, branch string, perPage int, page int) ([]gitprovider.Commit, error) {

	driverName := a.client.Driver.String()
	switch driverName {
	case "azure":
		commit, _, err := a.client.Contents.List(ctx, a.repository.ID, "", branch)
		if err != nil {
			return nil, err
		}

		commits := make([]gitprovider.Commit, 0, len(commit))

		for _, fe := range commit {

			commits = append(commits, CommitFileReference{fileEntry: fe})

		}
		return commits, nil

	default:
		//TODO had to do switch as find branch is not supported by azure
		ref, _, err := a.client.Git.FindBranch(ctx, a.repository.FullName, branch)
		if err != nil {
			return nil, err
		}
		commit := CommitReference{reference: ref}
		return []gitprovider.Commit{commit}, nil

	}

	return nil, nil

}

func (o OrgRepository) ListPage(ctx context.Context, branch string, perPage int, page int) ([]gitprovider.Commit, error) {
	//findCommit, _, err := o.client.Git.FindCommit(ctx, o.repository.ID, branch)
	//if err != nil {
	//	return nil, err
	//}
	//
	//commits := make([]gitprovider.Commit, 0, 1)
	//
	//commits = append(commits, CommitFileReference{commit: findCommit})
	//
	//return commits, nil
	return nil, nil
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
