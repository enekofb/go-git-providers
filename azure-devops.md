# Azure devops support 

## Goals 

- To understand the cost of adding azure devops as an example of extending with other providers
- Same functionality as used by wge


## Current usage

for example 

```go
func (s *GitProviderService) WriteFilesToBranchAndCreatePullRequest(ctx context.Context,
	req WriteFilesToBranchAndCreatePullRequestRequest) (*WriteFilesToBranchAndCreatePullRequestResponse, error) {
	
	...
	
	if err := s.writeFilesToBranch(ctx, writeFilesToBranchRequest{
		Repository:    repo,
		HeadBranch:    req.HeadBranch,
		BaseBranch:    req.BaseBranch,
		CommitMessage: req.CommitMessage,
		Files:         req.Files,
	}); err != nil {
		return nil, fmt.Errorf("unable to write files to branch %q: %w", req.HeadBranch, err)
	}
    ...
	res, err := s.createPullRequest(ctx, createPullRequestRequest{
		Repository:  repo,
		HeadBranch:  req.HeadBranch,
		BaseBranch:  req.BaseBranch,
		Title:       req.Title,
		Description: req.Description,
	})
}
```


```golang
// UserRepository describes a repository owned by an user.
type UserRepository interface {
	// UserRepository and OrgRepository implement the Object interface,
	// allowing access to the underlying object returned from the API.
	Object
	// The repository can be updated.
	Updatable
	// The repository can be reconciled.
	Reconcilable
	// The repository can be deleted.
	Deletable
	// RepositoryBound returns repository reference details.
	RepositoryBound

	// Get returns high-level information about this repository.
	Get() RepositoryInfo
	// Set sets high-level desired state for this repository. In order to apply these changes in
	// the Git provider, run .Update() or .Reconcile().
	Set(RepositoryInfo) error

	// DeployKeys gives access to manipulating deploy keys to access this specific repository.
	DeployKeys() DeployKeyClient

	// Commits gives access to this specific repository commits
	Commits() CommitClient

	// Branches gives access to this specific repository branches
	Branches() BranchClient

	// PullRequests gives access to this specific repository pull requests
	PullRequests() PullRequestClient

	// Files gives access to this specific repository files
	Files() FileClient

	// Trees gives access to this specific repository trees.
	Trees() TreeClient
}
```

- can reference the repo
- manage commits
- manage branches
- mange pull request
- manage repository files

## PoC - Scenarios

We are going to just focus in the PR flow that is the request need

Given azure devops and user for it

1. can create repo
2. can create branch in repo with manifest
3. can create pr for that branch

```
Feature: can support azure devops for wge


Scenario: can create branch on a repo azure devops repo
Given an azure devops user
And repo `my-repo` exists
And file `file-to-add` exists
When created git branch `my-branch` with file `file-to-add`
Then branch is created
 

Scenario: can create pr on a branch
Given an azure devops user
And repo `my-repo` exists
And a branch `my-branch`
When created a pull request for `my-branch` 
Then pull request is created and got its url
```


Some example queries

```bash 
➜  go-git-providers git:(add-azure-devops) ✗ curl -vk -u username:"${AZURE_DEVOPS_TOKEN}" https://dev.azure.com/efernandezbreis/weaveworks/_apis/git/repositories/weaveworks\?api-version=6.0

{"id":"28a51a24-5148-4d1c-920c-4493be96ae25","name":"weaveworks","url":"https://dev.azure.com/efernandezbreis/d603df47-831a-4b35-a24d-2558cdd64f01/_apis/git/repositories/28a51a24-5148-4d1c-920c-4493be96ae25","project":{"id":"d603df47-831a-4b35-a24d-2558cdd64f01","name":"weaveworks","url":"https://dev.azure.com/efernandezbreis/_apis/projects/d603df47-831a-4b35-a24d-2558cdd64f01","state":"wellFormed","revision":11,"visibility":"private","lastUpdateTime":"2022-10-26T18:56:39.583Z"},"size":0,"remoteUrl":"https://efernandezbreis@dev.azure.com/efernandezbreis/weaveworks/_git/weaveworks","sshUrl":"git@ssh.dev.azure.com:v3/efernandezbreis/weaveworks/weaveworks","webUrl":"https://dev.azure.com/efernandezbreis/weaveworks/_git/weaveworks","_links":{"self":{"href":"https://dev.azure.com/efernandezbreis/d603df47-831a-4b35-a24d-2558cdd64f01/_apis/git/repositories/28a51a24-5148-4d1c-920c-4493be96ae25"},"project":{"href":"vstfs:///Classification/TeamProject/d603df47-831a-4b35-a24d-2558cdd64f01"},"web":{"href":"https://dev.azure.com/efernandezbreis/weaveworks/_git/weaveworks"},"ssh":{"href":"git@ssh.dev.azure.com:v3/efernandezbreis/weaveworks/weaveworks"},"commits":{"href":"https://dev.azure.com/efernandezbreis/d603df47-831a-4b35-a24d-2558cdd64f01/_apis/git/repositories/28a51a24-5148-4d1c-920c-4493be96ae25/commits"},"refs":{"href":"https://dev.azure.com/efernandezbreis/d603df47-831a-4b35-a24d-2558cdd64f01/_apis/git/repositories/28a51a24-5148-4d1c-920c-4493be96ae25/refs"},"pullRequests":{"href":"https://dev.azure.com/efernandezbreis/d603df47-831a-4b35-a24d-2558cdd64f01/_apis/git/repositories/28a51a24-5148-4d1c-920c-4493be96ae25/pullRequests"},"items":{"href":"https://dev.azure.com/efernandezbreis/d603df47-831a-4b35-a24d-2558cdd64f01/_apis/git/repositories/28a51a24-5148-4d1c-920c-4493be96ae25/items"},"pushes":{"href":"https://dev.azure.com/efernandezbreis/d603df47-831a-4b35-a24d-2558cdd64f01/_apis/git/repositories/28a51a24-5148-4d1c-920c-4493be96ae25/pushes"}},"isDisabled":false}%

```

# Limitations and furthere info 

- scm support for orgs in general, i believe in the api there is no such a concept  
- 


