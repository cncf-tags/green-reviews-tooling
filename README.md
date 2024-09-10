# CNCF Green Reviews WG

The CNCF Green Reviews Working Group (WG) is an open-source, community-led project that is part of the [CNCF Environmental Sustainability Technical Advisory Group (TAG ENV)](https://github.com/cncf/tag-env-sustainability).

The aim of the Green Reviews WG is to set up infrastructure to measure the sustainability footprint of [CNCF Projects](https://www.cncf.io/projects). 

Measuring the sustainability footprint of software is not an easy task. Our vision is that the WG will compute the sustainability data for every release of a CNCF project that requests a sustainability footprint assessment. To achieve such a vision, our goal is to develop a workflow that can integrate well with the existing software lifecycle of other CNCF projects.

A good way to practically understand the first version of the workflow that the WG is designing is to take a look at the following architecture diagram:

![green reviews workflow](./docs/images/workflow-vision.png)

The WG’s workflow vision is that every release requesting a sustainability footprint assessment will trigger a Github Action specified in the Green Reviews repo that will start a benchmarking pipeline. The pipeline’s job is to:

1. Spin up the Equinix Metal resources
2. Install Kubernetes and all the needed observability tools
3. Install the software that will be assessed
4. Execute the necessary benchmark test cases
5. Gather sustainability-related metrics
6. Publish sustainability metrics

If you are curious and want to discover more a good:
- You can read the deep dive [article](https://tag-env-sustainability.cncf.io/blog/2024-green-reviews-working-group-measuring-sustainability/)
- You can watch the Kubecon EU '24 [Maintainer Talk](https://www.youtube.com/watch?v=UFa8hxOGKwQ)

## Releases 

| Release | Date | Notes
|---|---|---|
| 0.1.0 | 14.05.24 | [Release Notes](https://github.com/cncf-tags/green-reviews-tooling/releases/tag/0.1.0)
| 0.2.0 | Planned before Kubecon NA '24 | [Tracking issue](https://github.com/cncf-tags/green-reviews-tooling/issues/83)

## Community

### Getting Started

Here are some resources to learn about the project:

- **Charter**: The [Green Reviews WG Charter](https://github.com/cncf/tag-env-sustainability/blob/main/working-groups/green-reviews/charter.md) outlines the WG's motivation, scope, goals, non-goals, and deliverables.
- **Issue Board**: Checkout the project's [Backlog](https://github.com/orgs/cncf/projects/10/views/12) to find something to work on.
- **Design Document** (***Archived - please refer [Documentation section](#documentation) for latest updates***): The [WG's design document](https://docs.google.com/document/d/19fzZW-IMv2kDNatKFHeHh7wqcEN0e2N60wzxvCGZd48/edit?usp=sharing) is a live document created and maintained by the open-source contributors of the WG. Everyone is welcome to contribute ideas, questions, comments, suggestions, and take ownership of the project's implementation.

### Contributing

All contributions are welcome, including code contributions, issues, suggestions, questions, product direction, collaborations, etc.

If you are interested in contributing to the project, head over to the [Contributing Guide](./CONTRIBUTING.md)!

### Documentation

All changes to the documentation must be added to the [docs](./docs/) folder.

## Roadmap

The roadmap contains some of the short and long-term goals of the Green Reviews WG. Timelines are estimates and may change.

| Timeline | Goals |
|---|---|
| Q4 '23 / Before & After KubeCon NA  | Announce Green Reviews WG during KubeCon NA Keynotes  |
|  | Draft of the implementation for the sustainability footprint measurement workflow (tracking issue) |
|  | Cluster up and running on Equinix (tracking issue) |
| Q1 '24 / Before KubeCon EU  | Measure the cloud native sustainability footprint of Falco [milestone](https://github.com/cncf-tags/green-reviews-tooling/milestone/1) |
|  | Present architecture as part of the maintainer track at KubeCon EU  |
| Q2 - Q3 '24 / Before KubeCon NA  | Automate cloud native sustainability footprint measuring of the Falco project (alpha version) - [Tracking issue](https://github.com/cncf-tags/green-reviews-tooling/issues/83) |
| Q4 '24 / After KubeCon NA | TBD: below possible ideas | 
|  | Experiment with a different Energy metrics provider like [Scaphandre](https://github.com/hubblo-org/scaphandre) and compare results
|  | Identify new project |  
|  | Reproduce workflow with next CNCF Project |