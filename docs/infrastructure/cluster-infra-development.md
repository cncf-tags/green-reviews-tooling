# Cluster infra development

## Overview

Our cluster runs on physical servers from Equinix Metal contributed to CNCF.

- We use OpenTofu to manage the cluster infra and an S3 bucket from AWS to store the tofu state.
- To manage components running in the cluster we use Flux.

## Pre-Requisites

- Install the [tofu](https://opentofu.org/docs/intro/install/) CLI
- Fork the tooling repo https://github.com/cncf-tags/green-reviews-tooling
- Get access to the `Green Reviews Cluster Dev` vault in the TAG ENV 1Password account (please post in the #tag-env-wg-green-reviews channel in CNCF Slack so we have tracking for your request) 

## Setup

- In a local copy of your fork change to the working directory

```sh
cd infrastructure/equinix-metal
```

- Set env vars for tofu from 1Password with AWS and Equinix credentials

```sh
export AWS_ACCESS_KEY_ID="*****"
export AWS_SECRET_ACCESS_KEY="*****"
export TF_VAR_equinix_auth_token="*****"
export TF_VAR_equinix_project_id="*****"
export TF_VAR_k3s_token="*****"
```

- Set env vars for your fork including a GitHub PAT for bootstrapping Flux

```sh
export TF_VAR_cluster_name="green-reviews-dev"
export TF_VAR_flux_github_user="*** Your GitHub user ***"
export TF_VAR_flux_github_repo="green-reviews-tooling"
export TF_VAR_flux_github_token="*** Your GitHub PAT ***"
export TF_VAR_flux_github_branch="*** Your branch ***"
export TF_VAR_ssh_public_key="*** Your SSH public key ***"
```

- Check tofu workspaces

```sh
tofu workspace list

  default
* dev
```

- If the dev workspace doesn't exist create it.

```sh
tofu workspace new dev
```

- Ensure you are using the **dev** workspace

```sh
tofu workspace select dev
```

## Making changes

Follow the usual tofu workflow. See [core workflow](https://opentofu.org/docs/intro/core-workflow/)
for more details.

- Make changes
- Initialize tofu

```sh
tofu init
```

- Review changes

```sh
tofu plan
```

- Test changes - **make sure you are using dev workspace**

```sh
tofu apply
```

- Cleanup - **make sure you are using dev workspace**

```sh
tofu destroy
```
