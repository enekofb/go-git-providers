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
	"github.com/fluxcd/go-git-providers/gitprovider"
)

type AzureCommitClient struct {
	repository UserRepository
}

func (a UserRepository) Commits() gitprovider.CommitClient {
	return a
}

func (a UserRepository) Branches() gitprovider.BranchClient {
	return UserRepositoryBranches{
		repository: a,
	}
}
