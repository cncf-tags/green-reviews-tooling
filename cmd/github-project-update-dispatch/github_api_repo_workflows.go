package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

func (obj *GithubRepository) DispatchCall(projData Project, version string) error {
	projToWorkflowTable := map[string]string{
		"falco": "falco.yaml",
	}

	workflowFileName, found := projToWorkflowTable[projData.Name]
	if !found {
		return fmt.Errorf("failed to find the workflow file related to project: %s", projData.Name+"/"+projData.Organization)
	}

	url, header, err := obj.genURLAndHeaders(GHWorkflowDispatchEndpoint, "", workflowFileName)
	if err != nil {
		return err
	}

	type WorkflowData struct {
		Ref    string `json:"ref"`
		Inputs struct {
			CncfProject    string `json:"cncf_project"`
			Version        string `json:"version"`
			CncfProjectSub string `json:"cncf_project_sub"`
		} `json:"inputs"`
	}

	for i := range projData.SubComponents {
		workflowData := WorkflowData{}
		workflowData.Inputs.CncfProject = projData.Name
		workflowData.Ref = "main" // TODO: need to figure it out
		workflowData.Inputs.Version = version
		workflowData.Inputs.CncfProjectSub = projData.SubComponents[i]

		var _newWorkflowData bytes.Buffer

		if err := json.NewEncoder(&_newWorkflowData).Encode(workflowData); err != nil {
			return fmt.Errorf("failed to serialize the body: %v", err)
		}

		resp, err := NewHTTPClient(
			http.MethodPost,
			*url,
			time.Minute,
			&_newWorkflowData,
			header,
		)
		if err != nil {
			return err
		}

		defer resp.Body.Close()
		if resp.StatusCode != http.StatusNoContent {
			return fmt.Errorf(
				"failed in workflowDispatch, status code was not 204, got: %v",
				resp.StatusCode)
		}

		slog.Info("Successfully workflow_dispatch",
			"Proj", projData.Name,
			"Org", projData.Organization,
			"Ver", version,
			"Subcomponent", projData.SubComponents[i])
	}

	return nil
}
