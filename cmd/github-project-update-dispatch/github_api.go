package main

import (
	"fmt"
	"net/url"
	"os"
	"strings"
)

func NewGHRepoVarStorage() (*GithubRepository, error) {
	v, ok := os.LookupEnv(EnvGithubPatVarName)
	if !ok || len(v) == 0 {
		return nil, fmt.Errorf("the environment variable for github pat is missing")
	}
	return &GithubRepository{githubToken: v}, nil
}

func generateVariableName(projName string) string {
	return strings.ToLower(projName + "_version")
}

func gitHubRepoEndpoint(
	org, repo string,
	operation githubAPIEndpointType,
	variableName, workflowFileName string) (result string, err error) {

	baseUrl := "https://api.github.com"
	urlPath := []string{"repos", org, repo, "actions"}

	switch operation {
	case GHRepoVariableEndpoint:
		urlPath = append(urlPath, "variables", variableName)
	case GHWorkflowDispatchEndpoint:
		urlPath = append(urlPath, "workflows", workflowFileName, "dispatches")
	case GHLatestRelease:
		urlPath = append(urlPath, "releases", "latest")
	}

	result, err = url.JoinPath(baseUrl, urlPath...)
	return
}

func (obj *GithubRepository) genURLAndHeaders(
	endpointType githubAPIEndpointType,
	variableName string, workflowFileName string,
) (*string, map[string]string, error) {

	url, err := gitHubRepoEndpoint(
		GHRepoOrganizationName,
		GHRepoName,
		endpointType,
		variableName,
		workflowFileName,
	)
	if err != nil {
		return nil, nil, err
	}

	return &url,
		map[string]string{
			"Accept":               "application/vnd.github+json",
			"Authorization":        "Bearer " + obj.githubToken,
			"X-GitHub-Api-Version": "2022-11-28",
		}, nil
}
