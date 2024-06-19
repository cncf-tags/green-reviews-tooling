#!/bin/bash


json_file="../projects/projects.json"
gh_token=$GH_TOKEN
workflow_organization_name="cncf-tags"
workflow_project_name="green-reviews-tooling"
workflow_dispatcher_file_name="dispatch.yaml"


jq -c '.projects[]' "$json_file" | while read -r project; do
    proj_name=$(echo "$project" | jq -r '.name')
    proj_organization=$(echo "$project" | jq -r '.organization')
    sub_components=$(echo "$project" | jq -r '.sub_components')

    echo "Project Name: $proj_name"
    echo "Organization: $proj_organization"
    echo "SubComponents: $sub_components"

    releaseUrl="https://api.github.com/repos/${proj_organization}/${proj_name}/releases/latest"
    
    
    latest_proj_version=$(curl -fsSL -X GET $releaseUrl | jq -r '.tag_name')
    echo "Version: $latest_proj_version"


    workflow_dispatch=$(curl -fsSL -X POST \
      -H "Accept: application/vnd.github+json" \
      -H "Authorization: Bearer $gh_token" \
      -H "X-GitHub-Api-Version: 2022-11-28" \
      "https://api.github.com/repos/$workflow_organization_name/$workflow_project_name/actions/workflows/$workflow_dispatcher_file_name/dispatches" \
      -d "{\"ref\":\"main\",\"inputs\":{\"cncf_project\":\"${proj_name}\",\"cncf_project_sub\":\"${sub_components}\",\"version\":\"${latest_proj_version}\"}}")

    echo "workflow_call event [proj: $proj_name]=> $workflow_dispatch"
    echo "-----------------------------"
done
