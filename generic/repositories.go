package generic

import (
	"context"
	"errors"
	"github.com/fluxcd/go-git-providers/gitprovider"
	"github.com/jenkins-x/go-scm/scm"
)

// UserRepositories operates on repositories the user has access to.
type UserRepositories struct {
	client *scm.Client
}

type OrgRepositories struct {
	client *scm.Client
}

// Get returns the repository at the given path.
//
// ErrNotFound is returned if the resource does not exist.
func (c *UserRepositories) Get(ctx context.Context, ref gitprovider.UserRepositoryRef) (gitprovider.UserRepository, error) {

	repository, _, err := c.client.Repositories.Find(ctx, ref.GetRepository())
	if err != nil {
		return nil, err
	}

	return UserRepository{
		repository: repository,
		client:     c.client,
	}, nil
}

// List all repositories in the given organization.
//
// List returns all available repositories, using multiple paginated requests if needed.
func (c *UserRepositories) List(ctx context.Context, ref gitprovider.UserRef) ([]gitprovider.UserRepository, error) {
	return nil, errors.New("not implemented")
}

// Create creates a repository for the given organization, with the data and options
//
// ErrAlreadyExists will be returned if the resource already exists.
func (c *UserRepositories) Create(ctx context.Context,
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
func (c *UserRepositories) Reconcile(ctx context.Context, ref gitprovider.UserRepositoryRef, req gitprovider.RepositoryInfo, opts ...gitprovider.RepositoryReconcileOption) (gitprovider.UserRepository, bool, error) {
	return nil, false, errors.New("not implemented")
}

func (ca *OrgRepositories) Get(ctx context.Context, r gitprovider.OrgRepositoryRef) (gitprovider.OrgRepository, error) {

	repository, _, err := ca.client.Repositories.Find(ctx, r.GetRepository())
	if err != nil {
		return nil, err
	}

	return OrgRepository{
		repository: repository,
		client:     ca.client,
	}, nil
}

// List all repositories in the given organization.
//
// List returns all available repositories, using multiple paginated requests if needed.
func (c *OrgRepositories) List(ctx context.Context, o gitprovider.OrganizationRef) ([]gitprovider.OrgRepository, error) {
	return nil, errors.New("not implemented")
}

// Create creates a repository for the given organization, with the data and options.
//
// ErrAlreadyExists will be returned if the resource already exists.
func (c *OrgRepositories) Create(ctx context.Context, r gitprovider.OrgRepositoryRef, req gitprovider.RepositoryInfo, opts ...gitprovider.RepositoryCreateOption) (gitprovider.OrgRepository, error) {
	return nil, errors.New("not implemented")
}

// Reconcile makes sure the given desired state (req) becomes the actual state in the backing Git provider.
//
// If req doesn't exist under the hood, it is created (actionTaken == true).
// If req doesn't equal the actual state, the resource will be updated (actionTaken == true).
// If req is already the actual state, this is a no-op (actionTaken == false).
func (c *OrgRepositories) Reconcile(ctx context.Context, r gitprovider.OrgRepositoryRef, req gitprovider.RepositoryInfo, opts ...gitprovider.RepositoryReconcileOption) (resp gitprovider.OrgRepository, actionTaken bool, err error) {
	return nil, false, errors.New("not implemented")
}

type UserRepository struct {
	repository *scm.Repository
	client     *scm.Client
}

type OrgRepository struct {
	repository *scm.Repository
	client     *scm.Client
}

func (o OrgRepository) APIObject() interface{} {
	return o.repository
}

func (o OrgRepository) Update(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (o OrgRepository) Reconcile(ctx context.Context) (actionTaken bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (o OrgRepository) Delete(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (o OrgRepository) Repository() gitprovider.RepositoryRef {
	//TODO implement me
	panic("implement me")
}

func (o OrgRepository) Get() gitprovider.RepositoryInfo {
	return gitprovider.RepositoryInfo{
		Description:   &(o.repository.FullName),
		DefaultBranch: &(o.repository.Branch),
		Visibility:    nil,
	}
}

func (o OrgRepository) Set(info gitprovider.RepositoryInfo) error {
	//TODO implement me
	panic("implement me")
}

func (o OrgRepository) DeployKeys() gitprovider.DeployKeyClient {
	//TODO implement me
	panic("implement me")
}

func (o OrgRepository) Commits() gitprovider.CommitClient {
	return o
}

func (o OrgRepository) Branches() gitprovider.BranchClient {
	return OrgRepositoryBranches{
		repository: o,
	}
}

func (o OrgRepository) PullRequests() gitprovider.PullRequestClient {
	return OrgRepositoryPullRequests{
		repository: o,
	}
}

func (o OrgRepository) Files() gitprovider.FileClient {
	//TODO implement me
	panic("implement me")
}

func (o OrgRepository) Trees() gitprovider.TreeClient {
	//TODO implement me
	panic("implement me")
}

func (o OrgRepository) TeamAccess() gitprovider.TeamAccessClient {
	//TODO implement me
	panic("implement me")
}

func (a UserRepository) APIObject() interface{} {

	//TODO implement me
	panic("implement me")
}

func (a UserRepository) Update(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (a UserRepository) Reconcile(ctx context.Context) (actionTaken bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (a UserRepository) Delete(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (a UserRepository) Repository() gitprovider.RepositoryRef {
	//TODO implement me
	panic("implement me")
}

func (a UserRepository) Get() gitprovider.RepositoryInfo {
	return gitprovider.RepositoryInfo{
		Description:   &(a.repository.FullName),
		DefaultBranch: &(a.repository.Branch),
		Visibility:    nil,
	}
}

func (a UserRepository) Set(info gitprovider.RepositoryInfo) error {
	//TODO implement me
	panic("implement me")
}

func (a UserRepository) DeployKeys() gitprovider.DeployKeyClient {
	//TODO implement me
	panic("implement me")
}
