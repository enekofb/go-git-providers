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
	"encoding/base64"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/google/go-github/v47/github"
	"github.com/gregjones/httpcache"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/fluxcd/go-git-providers/gitprovider"
	"github.com/fluxcd/go-git-providers/gitprovider/testutils"
)

const (
	gitproviderDomain = "dev.azure.com/"

	defaultDescription = "Foo description"
	// TODO: This will change
	defaultBranch = "main"
)

var (
	// customTransportImpl is a shared instance of a customTransport, allowing counting of cache hits.
	customTransportImpl *customTransport
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

func customTransportFactory(transport http.RoundTripper) http.RoundTripper {
	if customTransportImpl != nil {
		panic("didn't expect this function to be called twice")
	}
	customTransportImpl = &customTransport{
		transport:      transport,
		countCacheHits: false,
		cacheHits:      0,
		mux:            &sync.Mutex{},
	}
	return customTransportImpl
}

type customTransport struct {
	transport      http.RoundTripper
	countCacheHits bool
	cacheHits      int
	mux            *sync.Mutex
}

func (t *customTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.mux.Lock()
	defer t.mux.Unlock()

	resp, err := t.transport.RoundTrip(req)
	// If we should count, count all cache hits whenever found
	if t.countCacheHits {
		if _, ok := resp.Header[httpcache.XFromCache]; ok {
			t.cacheHits++
		}
	}
	return resp, err
}

func (t *customTransport) resetCounter() {
	t.mux.Lock()
	defer t.mux.Unlock()

	t.cacheHits = 0
}

func (t *customTransport) setCounter(state bool) {
	t.mux.Lock()
	defer t.mux.Unlock()

	t.countCacheHits = state
}

func (t *customTransport) getCacheHits() int {
	t.mux.Lock()
	defer t.mux.Unlock()

	return t.cacheHits
}

func (t *customTransport) countCacheHitsForFunc(fn func()) int {
	t.setCounter(true)
	t.resetCounter()
	fn()
	t.setCounter(false)
	return t.getCacheHits()
}

var _ = Describe("Azure Devops Provider", func() {
	var (
		ctx context.Context = context.Background()
		c   gitprovider.Client

		testOrgRepoName  string
		testUserRepoName string
		testOrgName      string = "fluxcd-testing"
		testUser         string = "fluxcd-gitprovider-bot"
	)

	BeforeSuite(func() {
		rawToken := os.Getenv("AZURE_DEVOPS_TOKEN")
		if len(rawToken) == 0 {
			Fail("couldn't acquire AZURE_DEVOPS_TOKEN env variable")
		}
		var token string
		if rawToken != "" {
			token = base64.StdEncoding.EncodeToString([]byte(":" + rawToken))
		}

		if orgName := os.Getenv("GIT_PROVIDER_ORGANIZATION"); len(orgName) != 0 {
			testOrgName = orgName
		}

		if gitProviderUser := os.Getenv("GIT_PROVIDER_USER"); len(gitProviderUser) != 0 {
			testUser = gitProviderUser
		}

		var err error
		c, err = NewClient(ClientOptions{
			org:     "efernandezbreis",
			project: "weaveworks",
			token:   token,
		},
		)
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
			return retryOp.IsRetryable(err, fmt.Sprintf("get commits, repository: %s", userRepo.Repository().GetRepository()))
		}, retryOp.Timeout(), retryOp.Interval()).Should(BeTrue())

		latestCommit := commits[0]

		branchName := fmt.Sprintf("test-branch-%03d", rand.Intn(1000))

		err := userRepo.Branches().Create(ctx, branchName, latestCommit.Get().Sha)
		Expect(err).ToNot(HaveOccurred())

		err = userRepo.Branches().Create(ctx, branchName, "wrong-sha")
		Expect(err).To(HaveOccurred())

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

		prs, err := userRepo.PullRequests().List(ctx)
		Expect(len(prs)).To(Equal(1))
		Expect(prs[0].Get().WebURL).To(Equal(pr.Get().WebURL))

		err = userRepo.PullRequests().Merge(ctx, pr.Get().Number, gitprovider.MergeMethodSquash, "squash merged")
		Expect(err).ToNot(HaveOccurred())

		getPR, err := userRepo.PullRequests().Get(ctx, pr.Get().Number)
		Expect(err).ToNot(HaveOccurred())

		Expect(getPR.Get().Merged).To(BeTrue())

		path = "setup/config2.txt"
		content = "yaml content"
		files = []gitprovider.CommitFile{
			{
				Path:    &path,
				Content: &content,
			},
		}

		_, err = userRepo.Commits().Create(ctx, branchName, "added second config file", files)
		Expect(err).ToNot(HaveOccurred())

		pr, err = userRepo.PullRequests().Create(ctx, "Added second config file", branchName, *defaultBranch, "added second config file")
		Expect(err).ToNot(HaveOccurred())
		Expect(pr.Get().WebURL).ToNot(BeEmpty())
		Expect(pr.Get().Merged).To(BeFalse())

		err = userRepo.PullRequests().Merge(ctx, pr.Get().Number, gitprovider.MergeMethodMerge, "merged")
		Expect(err).ToNot(HaveOccurred())

		getPR, err = userRepo.PullRequests().Get(ctx, pr.Get().Number)
		Expect(err).ToNot(HaveOccurred())

		Expect(getPR.Get().Merged).To(BeTrue())
	})

	AfterSuite(func() {
		if os.Getenv("SKIP_CLEANUP") == "1" {
			return
		}
		// Don't do anything more if c wasn't created
		if c == nil {
			return
		}

		// Delete the org test repo used
		orgRepo, err := c.OrgRepositories().Get(ctx, newOrgRepoRef(testOrgName, testOrgRepoName))
		if err != nil && len(os.Getenv("CLEANUP_ALL")) > 0 {
			fmt.Fprintf(os.Stderr, "failed to get repo: %s in org: %s, error: %s\n", testOrgRepoName, testOrgName, err)
			fmt.Fprintf(os.Stderr, "CLEANUP_ALL set so continuing\n")
		} else {
			Expect(err).ToNot(HaveOccurred())
			Expect(orgRepo.Delete(ctx)).ToNot(HaveOccurred())
		}
		// Delete the user test repo used
		userRepo, err := c.UserRepositories().Get(ctx, newUserRepoRef(testUser, testUserRepoName))
		Expect(err).ToNot(HaveOccurred())
		Expect(userRepo.Delete(ctx)).ToNot(HaveOccurred())
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
