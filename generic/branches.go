package generic

import (
	"context"
	"errors"
	"fmt"
	"github.com/jenkins-x/go-scm/scm"
	"net/http"
)

type UserRepositoryBranches struct {
	repository UserRepository
}

func (a UserRepositoryBranches) Create(ctx context.Context, branch, sha string) error {

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

type OrgRepositoryBranches struct {
	repository OrgRepository
}

func (o OrgRepositoryBranches) Create(ctx context.Context, branch, sha string) error {

	input := &scm.ReferenceInput{
		Name: branch,
		Sha:  sha,
	}
	_, response, err := o.repository.client.Git.CreateRef(ctx, o.repository.repository.ID, input.Name, input.Sha)
	if err != nil {
		return err
	}

	if response.Status != http.StatusOK {
		return errors.New(fmt.Sprintf("CreateBranch did not get a 200 back %v", response.Status))
	}

	return nil
}
