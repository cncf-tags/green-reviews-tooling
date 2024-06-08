package main

type Project struct {
	Name          string   `json:"name"`
	Organization  string   `json:"organization"`
	SubComponents []string `json:"sub_components"`
}

type Projects struct {
	Projects []Project `json:"projects"`
}

type GithubRepository struct {
	githubToken string
}

type githubAPIEndpointType uint
