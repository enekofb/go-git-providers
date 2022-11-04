package generic

import (
	"context"
	"errors"
	drone "github.com/drone/go-scm/scm"
	"github.com/fluxcd/go-git-providers/gitprovider"
)

type Organisations struct {
	client *drone.Client
}

type Organization struct {
	organisation gitprovider.OrganizationRef
}

func (o Organization) APIObject() interface{} {
	return o.organisation
}

func (o Organization) Organization() gitprovider.OrganizationRef {
	//TODO implement me
	panic("implement me")
}

func (o Organization) Get() gitprovider.OrganizationInfo {
	//TODO implement me
	panic("implement me")
}

func (o Organization) Teams() gitprovider.TeamsClient {
	//TODO implement me
	panic("implement me")
}

// TODO org is not supported by azure scm
func (c *Organisations) Get(ctx context.Context, ref gitprovider.OrganizationRef) (gitprovider.Organization, error) {

	return Organization{organisation: ref}, nil
}

// TODO: not supported by scm
// https://github.com/jenkins-x/go-scm/bpalob/main/scm/driver/azure/org.go
func (c *Organisations) List(ctx context.Context) ([]gitprovider.Organization, error) {
	return nil, errors.New("not implemented")
}

func (c *Organisations) Children(_ context.Context, _ gitprovider.OrganizationRef) ([]gitprovider.Organization, error) {
	return nil, errors.New("not implemented")
}
