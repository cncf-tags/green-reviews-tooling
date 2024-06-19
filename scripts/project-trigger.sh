#!/bin/bash


json_file="../projects/projects.json"
gh_token=$GH_TOKEN
git_ref="main"
workflow_organization_name="cncf-tags"
workflow_project_name="green-reviews-tooling"
workflow_dispatcher_file_name="dispatch.yaml"

# if [ -z "$gh_token" ]; then
#     echo "[FATL] GH_TOKEN not set"
#     exit 20
# fi

jq -c '.projects[]' "$json_file" | while read -r project; do
    proj_name=$(echo "$project" | jq -r '.name')
    proj_organization=$(echo "$project" | jq -r '.organization')
    sub_components=$(echo "$project" | jq -r '.sub_components')

    echo "[DBG] Project Name: $proj_name"
    echo "[DBG] Organization: $proj_organization"
    echo "[DBG] SubComponents: $sub_components"

    release_url="https://api.github.com/repos/${proj_organization}/${proj_name}/releases/latest"
    
    
    response=$(curl --fail-with-body -sSL -X GET $release_url)
    status_code=$?
    if [ $status_code -ne 0 ]; then
        echo "[ERR] fetching latest release for ${proj_name} from ${release_url}. Status code: $status_code"
        continue
    fi

    latest_proj_version=$(echo "$response" | jq -r '.tag_name')
    if [ -z "$latest_proj_version" ]; then
        echo "[ERR] Could not find the latest version for ${proj_name}"
        continue
    fi
    echo "[DBG] Version: $latest_proj_version"

    workflow_dispatch=$(curl --fail-with-body -sSL -X POST \
        -H "Accept: application/vnd.github+json" \
        -H "Authorization: Bearer $gh_token" \
        -H "X-GitHub-Api-Version: 2022-11-28" \
        "https://api.github.com/repos/$workflow_organization_name/$workflow_project_name/actions/workflows/$workflow_dispatcher_file_name/dispatches" \
        -d "{\"ref\":\"${git_ref}\",\"inputs\":{\"cncf_project\":\"${proj_name}\",\"cncf_project_sub\":${sub_components},\"version\":\"${latest_proj_version}\"}}")

    status_code=$?
    if [ $status_code -ne 0 ]; then
        echo "[ERR] dispatching workflow for ${proj_name}. Status code: $status_code"
        echo "[DBG] Response: $workflow_dispatch_response"
        continue
    fi

    echo "[INF] workflow_call event [proj: $proj_name]=> $workflow_dispatch"
    echo "-----------------------------"
done
