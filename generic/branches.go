package generic

import (
	"context"
	"errors"
	"fmt"
	drone "github.com/drone/go-scm/scm"
	"net/http"
)

type UserRepositoryBranches struct {
	repository UserRepository
}

func (a UserRepositoryBranches) Create(ctx context.Context, branch, sha string) error {

	input := &drone.CreateBranch{
		Name: branch,
		Sha:  sha,
	}
	response, err := a.repository.client.Git.CreateBranch(ctx, a.repository.repository.ID, input)
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

	input := &drone.CreateBranch{
		Name: branch,
		Sha:  sha,
	}
	response, err := o.repository.client.Git.CreateBranch(ctx, o.repository.repository.ID, input)
	if err != nil {
		return err
	}

	if response.Status != http.StatusOK {
		return errors.New(fmt.Sprintf("CreateBranch did not get a 200 back %v", response.Status))
	}

	return nil
}
