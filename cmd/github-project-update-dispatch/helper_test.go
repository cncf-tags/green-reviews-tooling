package main

import (
	"testing"

	"gotest.tools/assert"
)

func TestConstants(t *testing.T) {
	assert.Equal(t, EnvGithubPatVarName, "GH_TOKEN")
	assert.Equal(t, GHRepoOrganizationName, "cncf-tags")
	assert.Equal(t, GHRepoName, "green-reviews-tooling")

	assert.Equal(t, GHRepoVariableEndpoint, githubAPIEndpointType(0))
	assert.Equal(t, GHWorkflowDispatchEndpoint, githubAPIEndpointType(1))
	assert.Equal(t, GHLatestRelease, githubAPIEndpointType(2))
}
