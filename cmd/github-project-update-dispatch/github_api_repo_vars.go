package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

func (obj *GithubRepository) UpdateRepoVariable(projName, newVersion string) error {
	variableName := generateVariableName(projName)

	url, header, err := obj.genURLAndHeaders(GHRepoVariableEndpoint, variableName, "")
	if err != nil {
		return err
	}

	newVariableData := struct {
		Name string `json:"name"`
		Val  string `json:"value"`
	}{
		Name: strings.ToUpper(variableName),
		Val:  newVersion,
	}

	var _newVariableData bytes.Buffer

	if err := json.NewEncoder(&_newVariableData).Encode(newVariableData); err != nil {
		return fmt.Errorf("failed to serialize the body: %v", err)
	}

	resp, err := NewHTTPClient(
		http.MethodPatch,
		*url,
		time.Minute,
		&_newVariableData,
		header,
	)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("status code was not 204, got: %v", resp.StatusCode)
	}

	return nil
}

func (obj *GithubRepository) ReadRepoVariable(projName string) (*string, error) {
	variableName := generateVariableName(projName)
	url, header, err := obj.genURLAndHeaders(GHRepoVariableEndpoint, variableName, "")
	if err != nil {
		return nil, err
	}

	resp, err := NewHTTPClient(
		http.MethodGet,
		*url,
		time.Minute,
		nil,
		header,
	)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code was not 200, got: %v", resp.StatusCode)
	}

	var variableData struct {
		VariableValue string `json:"value"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&variableData); err != nil {
		return nil, fmt.Errorf("failed deserialize response body: %v", err)
	}

	return &variableData.VariableValue, nil
}

func fetchLatestRelease(org, proj string) (*string, error) {

	url, err := gitHubRepoEndpoint(
		GHRepoOrganizationName,
		GHRepoName,
		GHLatestRelease,
		"", "",
	)
	if err != nil {
		return nil, err
	}

	resp, err := NewHTTPClient(
		http.MethodGet,
		url,
		time.Minute, nil, nil)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code was not 200, got: %v", resp.StatusCode)
	}

	var release struct {
		TagName string `json:"tag_name"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("failed deserialize response body: %v", err)
	}

	slog.Info("Latest Release", "Proj", proj, "Org", org, "Ver", release.TagName)

	return &release.TagName, nil
}
