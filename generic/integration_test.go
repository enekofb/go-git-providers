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
	"fmt"
	"github.com/stretchr/testify/require"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/google/go-github/v47/github"
	. "github.com/onsi/gomega"

	"github.com/fluxcd/go-git-providers/gitprovider"
)

const (
	gitproviderDomain = "dev.azure.com/"

	defaultDescription = "Foo description"
	// TODO: This will change
	defaultBranch = "main"
)

func init() {
	// Call testing.Init() prior to tests.NewParams(), as otherwise -test.* will not be recognised. See also: https://golang.org/doc/go1.13#testing
	testing.Init()
	rand.Seed(time.Now().UnixNano())
}

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

func TestCreatePR(t *testing.T) {

	var gitProviders = []struct {
		kind        string
		server      string
		tokenEnvVar string
		user        string
		repo        string
	}{
		//{"azure", "https://dev.azure.com", "AZURE_DEVOPS_TOKEN", "efernandezbreis", "weaveworks"},
		{"gitea", "http://localhost:3000", "GITEA_TOKEN", "gitea", "gitea/weaveworks"},
		//{"bitbucketcloud", "", "GITEA_TOKEN", "enekoww", "enekoww/test"},
	}

	for _, gitProvider := range gitProviders {
		t.Run(gitProvider.kind+" should be possible to create a pr for a user repository", func(t *testing.T) {
			// 1 - create client

			os.Setenv("GIT_KIND", gitProvider.kind)
			os.Setenv("GIT_SERVER", gitProvider.server)
			if gitProvider.tokenEnvVar != "" {
				os.Setenv("GIT_TOKEN", os.Getenv(gitProvider.tokenEnvVar))
			}
			os.Setenv("GIT_USER", gitProvider.user)

			var c gitprovider.Client
			var err error

			//TODO inconsistent api for creating client between azure and generic
			switch gitProvider.kind {
			case "azure":
				c, err = createAzureClient(gitProvider)
			default:
				c, err = NewClientFromEnvironment()
			}

			// 2- get repo

			ctx := context.Background()
			userRepoRef := newUserRepoRef(gitProvider.user, gitProvider.repo)
			var userRepo gitprovider.UserRepository
			userRepo, err = c.UserRepositories().Get(ctx, userRepoRef)
			require.NoError(t, err)

			// 3 - get latest commit for default branch

			defaultBranch := userRepo.Get().DefaultBranch

			var commits []gitprovider.Commit = []gitprovider.Commit{}

			commits, err = userRepo.Commits().ListPage(ctx, *defaultBranch, 1, 0)
			require.NoError(t, err)
			if err == nil && len(commits) == 0 {
				t.Errorf("empty commits list")
			}

			latestCommit := commits[0]
			require.NotEmpty(t, latestCommit)

			branchName := fmt.Sprintf("test-branch-%03d", rand.Intn(1000))

			// 4 create branch out of it
			switch gitProvider.kind {
			//TODO gitea does not support creating references https://github.com/jenkins-x/go-scm/blob/main/scm/driver/gitea/git.go#L39
			//it would require it to extend it via https://pkg.go.dev/code.gitea.io/sdk/gitea#Client.CreateBranch
			case "gitea":
				//using a fixed one instead to move one
				branchName = "test"
			default:
				err = userRepo.Branches().Create(ctx, branchName, latestCommit.Get().Sha)
				require.NoError(t, err)
			}

			// 5 create commit on branch
			path := "setup/config.txt"
			content := "yaml content"
			files := []gitprovider.CommitFile{
				{
					Path:    &path,
					Content: &content,
				},
			}
			_, err = userRepo.Commits().Create(ctx, branchName, "added config file", files)
			//require.NoError(t, err)

			// 6 create pr

			pr, err := userRepo.PullRequests().Create(ctx, "Added config file", branchName, *defaultBranch, "added config file")
			require.NoError(t, err)
			require.NotEmpty(t, pr.Get().WebURL)

		})
	}
}
