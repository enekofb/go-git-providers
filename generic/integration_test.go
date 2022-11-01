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

package generic

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/google/go-github/v47/github"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/fluxcd/go-git-providers/gitprovider"
	"github.com/fluxcd/go-git-providers/gitprovider/testutils"
)

const (
	gitproviderDomain = "generic"

	defaultDescription = "Foo description"
	// TODO: This will change
	defaultBranch = "main"
)

func init() {
	// Call testing.Init() prior to tests.NewParams(), as otherwise -test.* will not be recognised. See also: https://golang.org/doc/go1.13#testing
	testing.Init()
	rand.Seed(time.Now().UnixNano())
}

func TestProvider(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GitHub Provider Suite")
}

var _ = Describe("Gitea Provider", func() {
	var (
		ctx context.Context = context.Background()
		c   gitprovider.Client
	)

	BeforeSuite(func() {
		var err error
		c, err = NewClientFromEnvironment()
		Expect(err).ToNot(HaveOccurred())
	})

	It("should be possible to create a pr for a user repository", func() {

		userRepoRef := newUserRepoRef("efernandezbreis", "weaveworks")

		var userRepo gitprovider.UserRepository
		retryOp := testutils.RetryOp{}
		Eventually(func() bool {
			var err error
			userRepo, err = c.UserRepositories().Get(ctx, userRepoRef)
			return retryOp.IsRetryable(err, fmt.Sprintf("get user repository: %s", userRepoRef.RepositoryName))
		}, retryOp.Timeout(), retryOp.Interval()).Should(BeTrue())

		defaultBranch := userRepo.Get().DefaultBranch

		var commits []gitprovider.Commit = []gitprovider.Commit{}
		Eventually(func() bool {
			var err error
			commits, err = userRepo.Commits().ListPage(ctx, *defaultBranch, 1, 0)
			if err == nil && len(commits) == 0 {
				err = errors.New("empty commits list")
			}
			return err == nil && len(commits) > 0
		}, retryOp.Timeout(), retryOp.Interval()).Should(BeTrue())

		latestCommit := commits[0]

		branchName := fmt.Sprintf("test-branch-%03d", rand.Intn(1000))

		err := userRepo.Branches().Create(ctx, branchName, latestCommit.Get().Sha)
		Expect(err).ToNot(HaveOccurred())

		path := "setup/config.txt"
		content := "yaml content"
		files := []gitprovider.CommitFile{
			{
				Path:    &path,
				Content: &content,
			},
		}

		_, err = userRepo.Commits().Create(ctx, branchName, "added config file", files)
		Expect(err).ToNot(HaveOccurred())

		pr, err := userRepo.PullRequests().Create(ctx, "Added config file", branchName, *defaultBranch, "added config file")
		Expect(err).ToNot(HaveOccurred())
		Expect(pr.Get().WebURL).ToNot(BeEmpty())
		Expect(pr.Get().Merged).To(BeFalse())

	})

	It("should be possible to create a pr for an org repository", func() {

		// get org
		orgRef := newOrgRef("efernandezbreis")
		var _ gitprovider.Organization
		retryOp := testutils.RetryOp{}
		Eventually(func() bool {
			var err error
			_, err = c.Organizations().Get(ctx, orgRef)
			return retryOp.IsRetryable(err, fmt.Sprintf("get org: %s", orgRef.Organization))
		}, retryOp.Timeout(), retryOp.Interval()).Should(BeTrue())

		// get org repository

		orgRepoRef := newOrgRepoRef(orgRef.Organization, "weaveworks")
		orgRepo, err := c.OrgRepositories().Get(ctx, orgRepoRef)

		defaultBranch := orgRepo.Get().DefaultBranch

		var commits []gitprovider.Commit = []gitprovider.Commit{}
		Eventually(func() bool {
			var err error
			commits, err = orgRepo.Commits().ListPage(ctx, *defaultBranch, 1, 0)
			if err == nil && len(commits) == 0 {
				err = errors.New("empty commits list")
			}
			return err == nil && len(commits) > 0
		}, retryOp.Timeout(), retryOp.Interval()).Should(BeTrue())

		latestCommit := commits[0]

		branchName := fmt.Sprintf("test-branch-%03d", rand.Intn(1000))

		err = orgRepo.Branches().Create(ctx, branchName, latestCommit.Get().Sha)
		Expect(err).ToNot(HaveOccurred())

		path := "setup/config.txt"
		content := "yaml content"
		files := []gitprovider.CommitFile{
			{
				Path:    &path,
				Content: &content,
			},
		}

		_, err = orgRepo.Commits().Create(ctx, branchName, "added config file", files)
		Expect(err).ToNot(HaveOccurred())

		pr, err := orgRepo.PullRequests().Create(ctx, "Added config file", branchName, *defaultBranch, "added config file")
		Expect(err).ToNot(HaveOccurred())
		Expect(pr.Get().WebURL).ToNot(BeEmpty())
		Expect(pr.Get().Merged).To(BeFalse())

	})
})

func newOrgRef(organizationName string) gitprovider.OrganizationRef {
	return gitprovider.OrganizationRef{
		Domain:       gitproviderDomain,
		Organization: organizationName,
	}
}

func newOrgRepoRef(organizationName, repoName string) gitprovider.OrgRepositoryRef {
	return gitprovider.OrgRepositoryRef{
		OrganizationRef: newOrgRef(organizationName),
		RepositoryName:  repoName,
	}
}

func newUserRef(userLogin string) gitprovider.UserRef {
	return gitprovider.UserRef{
		Domain:    gitproviderDomain,
		UserLogin: userLogin,
	}
}

func newUserRepoRef(userLogin, repoName string) gitprovider.UserRepositoryRef {
	return gitprovider.UserRepositoryRef{
		UserRef:        newUserRef(userLogin),
		RepositoryName: repoName,
	}
}

func findUserRepo(repos []gitprovider.UserRepository, name string) gitprovider.UserRepository {
	if name == "" {
		return nil
	}
	for _, repo := range repos {
		if repo.Repository().GetRepository() == name {
			return repo
		}
	}
	return nil
}

func findOrgRepo(repos []gitprovider.OrgRepository, name string) gitprovider.OrgRepository {
	if name == "" {
		return nil
	}
	for _, repo := range repos {
		if repo.Repository().GetRepository() == name {
			return repo
		}
	}
	return nil
}

func validateRepo(repo gitprovider.OrgRepository, expectedRepoRef gitprovider.RepositoryRef) {
	info := repo.Get()
	// Expect certain fields to be set
	Expect(repo.Repository()).To(Equal(expectedRepoRef))
	Expect(*info.Description).To(Equal(defaultDescription))
	Expect(*info.Visibility).To(Equal(gitprovider.RepositoryVisibilityPrivate))
	Expect(*info.DefaultBranch).To(Equal(defaultBranch))
	// Expect high-level fields to match their underlying data
	internal := repo.APIObject().(*github.Repository)
	Expect(repo.Repository().GetRepository()).To(Equal(*internal.Name))
	Expect(repo.Repository().GetIdentity()).To(Equal(internal.Owner.GetLogin()))
	Expect(*info.Description).To(Equal(*internal.Description))
	Expect(string(*info.Visibility)).To(Equal(*internal.Visibility))
	Expect(*info.DefaultBranch).To(Equal(*internal.DefaultBranch))
}
