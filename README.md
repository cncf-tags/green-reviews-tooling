# CNCF Green Reviews WG

The CNCF Green Reviews Working Group (WG) is an open-source, community-led project that is part of the [CNCF Environmental Sustainability Technical Advisory Group (TAG ENV)](https://github.com/cncf/tag-env-sustainability).

The aim of the Green Reviews WG is to set up infrastructure to measure the sustainability footprint of [CNCF Projects](https://www.cncf.io/projects).

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

All the documentation for Green Reviews WG is moved to a dedicated website that's publicly available at [https://docs.green-reviews.tag-env-sustainability.cncf.io/](https://docs.green-reviews.tag-env-sustainability.cncf.io/).
All changes to the documentation must be added to [website/content/docs](./website/content/docs/) folder. Documentation is built on top of the [Doks theme](https://getdoks.org/docs/start-here/getting-started) powered by Hugo. Website is deployed via Netlify.

To run website locally, run the following commands (requires ```npm```):

``` bash
cd website
npm run install
npm run dev
```

More documentation for further development of the website can be found in [Doks theme docs](https://getdoks.org/docs/start-here/getting-started)😁

## Roadmap

The roadmap contains some of the short and long-term goals of the Green Reviews WG. Timelines are estimates and may change.

| Timeline | Goals |
|---|---|
| Q4 '23 / Before & After KubeCon NA  | Announce Green Reviews WG during KubeCon NA Keynotes  |
| | Draft of the implementation for the sustainability footprint measurement workflow (tracking issue) |
| | Cluster up and running on Equinix (tracking issue) |
| Q1 '24 / Before KubeCon EU  | Measure the cloud native sustainability footprint of Falco [milestone](https://github.com/cncf-tags/green-reviews-tooling/milestone/1) |
| | Present architecture as part of the maintainer track at KubeCon EU  |
| Q2 '24 / After KubeCon EU  | Automate cloud native sustainability footprint measuring of the Falco project (alpha version) |
|| First stats published publicly on devstats.cncf.io |
|| Identify the next CNCF Project |
| Q3 '24 / Before KubeCon NA  |  Reproduce workflow with next CNCF Project |
