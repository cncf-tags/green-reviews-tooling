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

Step 3 in the automated pipeline to report and store the carbon emission results for a CNCF project. 

See also: 
- step 1: [Trigger and Deploy](./proposal-001-trigger-and-deploy.md) 
- step 2: [Run benchmark tests](./proposal-002-run.md).

- Tracking issue: [#95](https://github.com/cncf-tags/green-reviews-tooling/issues/95)
- Implementation issue: TBD

## Authors

- @chrischinchilla
- @AntonioDiTuri

## Status

WIP

## Table of Contents

<!-- toc -->

- [Proposal 003 - Report project benchmark tests from the automated pipeline](#proposal-003---report-project-benchmark-tests-from-the-automated-pipeline)
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
    - [Risks and Mitigations](#risks-and-mitigations)
  - [Design Details](#design-details)
    - [Metrics](#metrics)
    - [Collect](#collect)
    - [Store](#store)
    - [Share](#share)
  - [Drawbacks (Optional)](#drawbacks-optional)
  - [Alternatives](#alternatives)
  - [Infrastructure Needed (Optional)](#infrastructure-needed-optional)

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

- Describe the two kind of metrics that are going to be stored:
  - Project-related metrics: specific metrics that a given project might request
  - Sustainability metrics: the metrics related to the green review
- Describe the steps and infrastructure needed to report and store the results of the pipeline:
  - Collect: the action of getting the metrics from their producers
  - Store: the action of saving the metrics in a state
  - Share: how to expose the metrics for the CNCF project maintainers
- For each step describe how the action should be implemented and why

### Non-Goals

<!--
What is out of scope for this proposal? Listing non-goals helps to focus discussion
and make progress.

It is important to remember that non-goals are still equally important things
which will be dealt with one day but are not things which need to be dealt with immediately
within the scope of this work. This helps make sure everyone is crystal clear on the outcomes.
-->

- Create new metrics from scratch 
- Aggregate existing metrics
- Provide analytic functionalities on top of the raw metrics
- Integration with cncf dev-stat on Grafana

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

**CNCF project maintainer selects project-related metrics** 

If the project produces significant metrics that need to be monitored along with the sustainable metrics, I would like to have them reported.

**Green reviews maintainer adds, modifies or removes sustainability metrics**

As a Green Review maintainer, I would like to change the sustainability metrics over time. 

**CNCF project maintainer is able to check the metrics for their project**

After the pipeline produces the metrics, I would like to see the result of it in an accessible way.


### Risks and Mitigations

<!--
What are the risks of this proposal, and how do we mitigate?
Think broadly.  For example, consider how this will impact or be impacted
by scaling to support more CNCF Projects.

How will this affect the benchmark tests, CNCF Project Maintainers, pipeline maintainers, etc?
-->

As with every design document there are some challanges:
 
- Consistency vs Flexibility: if we change the sustainability metrics overtime it will be difficult to compare different metrics from different green reviews release. However we would rather be flexible in this first phase and possibly change what we store and how if this leads to more correct results. 
- The three sub-steps of the proposal: Collect, Store and Share are co-dependent. How to collect the data depends on how to store it, and how to store the data depends on how to show it. Since we are still in early phases of the working group, an agile approach will be proposed: a first lean solution will be deployed and, most likely, improved in the future.  

## Design Details

This section will have the following subsections:

- Metrics: what metrics to collect?
- Collect: how to collect the metrics?
- Store: how to store them?
- Share: how to share the metrics?

### Metrics

As already mentioned we will have two sets of metrics:

1. Project-related metrics
2. Sustainability metrics

For the project related one Falco has already requested this metrics:

```
rate(container_cpu_usage_seconds_total[5m])
container_memory_rss
container_memory_working_set_bytes
```

For the Sustainbility metrics we will keep this one:

`kepler_container_joules_total`

### Collect

TBD

A prometheus query?
A direct curl to the sources?

Evaluate pros and cons

### Store

TBD

Something simple like a markdown file would do.
Who is writing to the file?
How to organize the file?

### Share

TBD

Grafana dashboard is needed?
Is it enough to show the markdown files?


<!--
This section should contain enough information that the specifics of your
change are understandable. This may include manifests or workflow examples
(though not always required) or even code snippets. If there's any ambiguity
about HOW your proposal will be implemented, this is the place to discuss them.
-->


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
