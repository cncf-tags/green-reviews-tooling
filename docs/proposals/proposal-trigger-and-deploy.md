# Deploy & trigger GitHub Action workflow from an upstream CNCF project

To trigger our benchmarking task to run when a particular cncf project gets certain kinds of event, let's say its `release`.
Some more info about this proposal is also present in [#83](https://github.com/cncf-tags/green-reviews-tooling/issues/83)

## Authors

- @rossf7
- @dipankardas011

## Status

implementable

<!--
Must be one of provisional, implementable, implemented, deferred,
rejected, withdrawn, or replaced.
-->

## Table of Contents
<!-- toc -->
- [Short, descriptive title](#short-descriptive-title)
  - [Authors](#authors)
  - [Status](#status)
  - [Table of Contents](#table-of-contents)
  - [Summary](#summary)
  - [Motivation](#motivation)
    - [Goals](#goals)
    - [Non-Goals](#non-goals)
    - [Linked Docs](#linked-docs)
  - [Proposal](#proposal)
    - [User Stories (Optional)](#user-stories-optional)
      - [Story 1](#story-1)
      - [Story 2](#story-2)
    - [Notes/Constraints/Caveats (Optional)](#notesconstraintscaveats-optional)
    - [Risks and Mitigations](#risks-and-mitigations)
  - [Design Details](#design-details)
    - [Graduation Criteria (Optional)](#graduation-criteria-optional)
  - [Drawbacks (Optional)](#drawbacks-optional)
  - [Alternatives](#alternatives)
  - [Infrastructure Needed (Optional)](#infrastructure-needed-optional)
<!-- /toc -->

## Summary

This proposal focuses on automating the Green Reviews pipeline for Falco by defining a trigger mechanism, involving the Falco team in the implementation, deploying Falco using Flux, and testing the deployment process. In future the pipeline will support more CNCF projects as they are onboarded.

The proposal also includes considerations for a phased implementation of the automation, starting with manual triggering followed by automation via a webhook.


## Motivation

To automate the trigger of Falco deployment when upstream aka origin repo creates an event.
we will then deploy the benchmarking workfload for the project, in this case its, Falco

### Goals

- Trigger GitHub Action workflow in green-reviews-tooling repo when Falco needs to be tested
- Ask Falco team to implement the trigger
- Deploy correct version of Falco in GitHub Action using Flux
- Test the deployment via the Falco trigger

### Non-Goals

- Creating cluster nodes on demand. [Future Goal Issue #67](https://github.com/cncf-tags/green-reviews-tooling/issues/67)


### Linked Docs

- **Slack Discussion Thread** [Link](https://cloud-native.slack.com/archives/C060EDHN431/p1712765271470189)
- **Triggering GitHub Action**: For triggering the workflow AIUI we could use a webhook to trigger a workflow_dispatch event. [Workflow Dispatch](https://docs.github.com/en/actions/using-workflows/events-that-trigger-workflows#workflow_dispatch). It allows providing custom inputs and as a minimum I think we need the name of the CNCF project and the version to be deployed. [Providing Inputs for event that trigger workflows](https://docs.github.com/en/actions/using-workflows/events-that-trigger-workflows#providing-inputs)


## Proposal

got an example GitHub action workflow file
```yaml
name: TriggerTest

on:
  workflow_dispatch:
    inputs:
      cncf_project:
        description: 'CNCF Project Name'
        required: true
        default: 'falco'
      cncf_project_sub:
        description: 'CNCF Project Subcomponent'
        required: false
        default: 'modern-ebpf'
      version:
        description: 'Version'
        required: true
        default: '0.37.0'

jobs:
  echo-inputs:
    runs-on: ubuntu-latest
    steps:
    - name: Checkout repository
      uses: actions/checkout@v3

    - name: Echo Inputs
      run: |
        echo "Add logic to deploy ${{ github.event.inputs.cncf_project }} ${{ github.event.inputs.cncf_project_sub }}"
        echo "version ${{ github.event.inputs.version }}"
```

for invoking this

```bash
curl -X POST \
     -H "Accept: application/vnd.github.v3+json" \
     -H "Authorization: token $GITHUB_PAT" \
     https://api.github.com/repos/rossf7/green-reviews-tooling/actions/workflows/trigger_test.yaml/dispatches \
     -d '{"ref":"main", "inputs": {"cncf_project": "falco", "cncf_project_sub": "modern-ebpf","version":"0.37.0"}}'
```

> [!NOTE]
> Here fine grained PAT is used
> - Read access to code and metadata
> - Read write access to actions


> [!IMPORTANT]
> We'll need to create these and provide it to the Falco team, aka any future CNCF project we are going to use for benchmarking.


### User Stories (Optional)

<!--
Detail the things that people will be able to do if this proposal is implemented.
Include as much detail as possible so that people can understand the "how" of
the system. The goal here is to make this feel real for users without getting
bogged down.
-->

#### Story 1

#### Story 2

### Notes/Constraints/Caveats (Optional)

<!--
What are the caveats to the proposal?
What are some important details that didn't come across above?
Go in to as much detail as necessary here.
This might be a good place to talk about core concepts and how they relate.
-->

### Risks and Mitigations

<!--
What are the risks of this proposal, and how do we mitigate?
Think broadly.  For example, consider how this will impact or be impacted
by scaling to support more CNCF Projects.

How will this affect the benchmark tests, CNCF Project Maintainers, pipeline maintainers, etc?
-->

## Design Details

> [!TODO]
> Add here

<!--
This section should contain enough information that the specifics of your
change are understandable. This may include manifests or workflow examples
(though not always required) or even code snippets. If there's any ambiguity
about HOW your proposal will be implemented, this is the place to discuss them.
-->

### Graduation Criteria (Optional)

<!--
List criteria which would allow progression from one maturity level to another.
eg. What needs to have been accomplished/demonstrated to move from Alpha to Beta.

If applicable, what is the milestone marker which will allow deprecation of the
replaced capability?
-->

## Drawbacks (Optional)

<!--
What other approaches did you consider, and why did you rule them out? These do
not need to be as detailed as the proposal, but should include enough
information to express the idea and why it was not acceptable.
-->

## Alternatives

<!--
What other approaches did you consider, and why did you rule them out? These do
not need to be as detailed as the proposal (pros and cons are fine),
but should include enough information to express the idea and why it was not acceptable
as well as illustrate why the final solution was selected.
-->

## Infrastructure Needed (Optional)

<!--
Use this section if you need things from the project/SIG. Examples include a
new subproject, repos requested, or GitHub details. Listing these here allows a
SIG to get the process for these resources started right away.
-->

- GitHub actions workflow
- OpenTofu
