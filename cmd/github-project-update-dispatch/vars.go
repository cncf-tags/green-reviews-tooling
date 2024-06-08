package main

const (
	EnvGithubPatVarName    string = "GH_TOKEN"
	GHRepoOrganizationName string = "cncf-tags"
	GHRepoName             string = "green-reviews-tooling"
)

const (
	GHRepoVariableEndpoint     githubAPIEndpointType = iota
	GHWorkflowDispatchEndpoint githubAPIEndpointType = iota
	GHLatestRelease            githubAPIEndpointType = iota
)
