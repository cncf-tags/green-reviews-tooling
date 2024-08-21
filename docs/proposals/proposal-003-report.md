---
title: "003 - Report"
description: "This is a change proposal template for the Green Reviews WG repository."
summary: ""
date: 2024-08-21T13:30:48+02:00
lastmod: 2024-08-21T13:30:48+02:00
draft: false
slug: 003-proposal-report
weight: 940
toc: true
---

<!--
How to use this template:

- Make a copy of this file in the docs/proposals/ directory
- Set the name of the file to contain the next logical number and the name of the feature
- Fill out at least the Status, Motivation and Goals/Non-Goals fields.
- Open a PR to green-reviews-tooling
- Merge early and iterate

For more tips see the Contributing docs: https://github.com/cncf-tags/green-reviews-tooling/blob/main/CONTRIBUTING.md#proposals
-->

# Proposal 003 - Report project benchmark tests from the automated pipeline

<!--
Keep the title short, simple, and descriptive. A good
title can help communicate what the proposal is and should be
considered as part of any review.
-->

Step 3 in the automated pipeline to report and store the carbon emission results for a CNCF project. See also step 1: [Trigger and Deploy](./proposal-001-trigger-and-deploy.md) and step 2: [Run benchmark tests](./proposal-002-run.md).

- Tracking issue: [#95](https://github.com/cncf-tags/green-reviews-tooling/issues/95)
- Implementation issue: TBD

## Authors

- @chrischinchilla
- @AntonioDiTuri

## Status

Provisional

<!--
The headings here are just starting points, add more as makes sense for what you
are proposing.
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

<!--
A good summary is at least a paragraph in length and should be written with a wide audience
in mind.

It should encompass the entire document, and serve as both future documentation
and as a quick reference for people coming by to learn the proposal's purpose
without reading the entire thing.

Both in this section and below, follow the guidelines of the [documentation
style guide]. In particular, wrap lines to a reasonable length, to make it
easier for reviewers to cite specific portions, and to minimize diff churn on
updates.

[documentation style guide]: https://github.com/kubernetes/community/blob/master/contributors/guide/style-guide.md
-->

## Motivation

<!--
This section is for explicitly listing the motivation, goals and non-goals of
this proposal. Describe why the change is important, how it fits into the project's
goals and the benefits to users.

It is helpful to frame this to answer the question: "What is the problem this proposal
is trying to solve?"
-->

This proposal is part of the pipeline automation of the Green Reviews tooling for Falco (and new CNCF projects in the future). It builds upon the previous steps of the pipeline to record, report, and store the results of the pipeline. It records, reports, and stores SRE metrics specified by the project as well as standard metrics for carbon emissions based on energy usage.

### Goals

<!--
List the specific goals of the proposal. What is it trying to achieve? How will we
know that this has succeeded?
-->

- Describe the steps and infrastructure needed to report and store the results of the pipeline.
- Export and store the reported metrics in an accessible format.

### Non-Goals

<!--
What is out of scope for this proposal? Listing non-goals helps to focus discussion
and make progress.

It is important to remember that non-goals are still equally important things
which will be dealt with one day but are not things which need to be dealt with immediately
within the scope of this work. This helps make sure everyone is crystal clear on the outcomes.
-->

- Creating new metrics

### Linked Docs

<!--
Provide links to previous discussions, Slack threads, motivation issues or any other document
with context. It is really helpful to provide a "source of truth" for the work
so that people aren't searching all over the place for lost context.
-->

## Proposal

<!--
This is where we get down to the specifics of what the proposal actually is:
outlining your solution to the problem described in the Motivation section.
This should have enough detail that reviewers can understand exactly what
you're proposing, but should not include things like API designs or
implementation. The "Design Details" section below is for the real
nitty-gritty.
-->

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

The main risks are that the metrics captured and recorded aren't useful or don't show much.

### Risks and Mitigations

<!--
What are the risks of this proposal, and how do we mitigate?
Think broadly.  For example, consider how this will impact or be impacted
by scaling to support more CNCF Projects.

How will this affect the benchmark tests, CNCF Project Maintainers, pipeline maintainers, etc?
-->

## Design Details

<!--
This section should contain enough information that the specifics of your
change are understandable. This may include manifests or workflow examples
(though not always required) or even code snippets. If there's any ambiguity
about HOW your proposal will be implemented, this is the place to discuss them.
-->

### Setup

1. Start Kubernetes
2. [Install and start Prometheus](https://sustainable-computing.io/installation/kepler/#deploy-the-prometheus-operator) 2. `cd kube-prometheus` 3. `kubectl apply --server-side -f manifests/setup` 4. `kubectl apply -f manifests/` 5. Wait… 6. `kubectl -n monitoring port-forward svc/grafana 3000` 7. Open dashboard _localhost:3000_
3. Install metrics server
   1. `kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml`
   2. Patch

```shell
kubectl patch -n kube-system deployment metrics-server --type=json \
-p '[{"op":"add","path":"/spec/template/spec/containers/0/args/-","value":"--kubelet-insecure-tls"}]'
```

4. Install and start Kepler
   1. Open a new terminal
   2. `git clone --depth 1 git@github.com:sustainable-computing-io/kepler.git`
5. `cd kepler`
6. `make build-manifest OPTS="PROMETHEUS_DEPLOY"`
7. `kubectl apply -f _output/generated-manifest/deployment.yaml`
8. Add [dashboard](https://raw.githubusercontent.com/sustainable-computing-io/kepler/main/grafana-dashboards/Kepler-Exporter.json) to Grafana.
9. Install and start Falco
   1. Open a new terminal
   2. [Install Helm](https://helm.sh/docs/intro/install/)
   3. `helm repo add falcosecurity https://falcosecurity.github.io/charts`
   4. `helm repo update`
   5. `helm install falco falcosecurity/falco --namespace falco --create-namespace --set driver.kind=modern-bpf  --set falco.grpc.enabled=true --set falco.grpc_output.enabled=true`
   6. `helm install falco-exporter falcosecurity/falco-exporter`
10. Run Falco tests
    1. https://github.com/falcosecurity/cncf-green-review-testing/tree/main/benchmark-tests
11. May need to remove `nodeSelector`
12. Write out metrics to JSON
13. Thinking about https://github.com/prometheus/prom2json
14. These metrics:

    ```
    rate(container_cpu_usage_seconds_total[5m])
    container_memory_rss
    container_memory_working_set_bytes
    kepler_container_joules_total
    ```

## Drawbacks (Optional)

JSON storage isn't so viable in the long term, this is an MVP solution and in the long term, will store metrics in a Prometheus-compatible storage solution.

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
