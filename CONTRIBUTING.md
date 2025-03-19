# Contributing Guide

Thank you for contributing to this project!
We really appreciate your time and willingness to contribute. The [Green Reviews
Working Group](https://github.com/cncf/tag-env-sustainability/tree/main/working-groups/green-reviews)
is an entirely volunteer-led open-source project where we innovate
in the open in the hopes of making this architectural reference available to
everyone who operates in a cloud-native environment.

Read on to learn more about the project and how to contribute. As always, please
don't hesitate to ask if anything is unclear. We value questions, guidance, and
suggestions - they help us build the right thing.

## Learn through our docs

Checkout our [docs](./docs/README.md) if you are interested in contributing and/or using the pipeline.

## Finding Work

Here are some suggestions for how to find work to contribute to the project.

First, make sure you're up to speed with the project by looking through the resources in the [Getting Started](./README.md#getting-started) guide.

Then, familiarise yourself with the current work and priorities, or find someone who would be willing to pair with you:

- Join the twice monthly [WG meetings](https://github.com/cncf/tag-env-sustainability/tree/main/working-groups/green-reviews#meetings-and-contact).
- [async] Leave a comment in any of the [issues](https://github.com/orgs/cncf/projects/10/views/12).
- [async] Contribute to discussions in the CNCF Slack channel [#wg-green-reviews](https://cloud-native.slack.com/archives/C060EDHN431).

When you feel ready to contribute:

- Check the Backlog in the [issue board](https://github.com/orgs/cncf-tags/projects/5). More pressing issues are labeled with `good first issue` or `help wanted`.
- If you would like to make a feature request or raise a bug, feel free to [open an issue](https://github.com/cncf-tags/green-reviews-tooling/issues/new/choose).

We encourage all communication to remain public by going through the communications channels listed above so that everyone can stay informed!

### Opening a Pull Request

Recommendations for a faster Pull Request review:

- Please ensure that your commits are **signed** and **verified**.
  - **Signed**: Include a `Signed-off-by: Author Name <authoremail@example.com>` in all commits by doing `git rebase HEAD~2 --signoff` (replace `~2` with the number of commits to sign) and then `git push -f`. More info [here](https://github.com/cncf-tags/green-reviews-tooling/pull/53/checks?check_run_id=21563565653).
  - **Verified**: [Learn more about commit signature verification with GPG.](https://docs.github.com/en/authentication/managing-commit-signature-verification/about-commit-signature-verification#gpg-commit-signature-verification)
- The Kubernetes Best Practices for faster Reviews is a great resource for PR best practices: [Best practices for faster reviews](https://git.k8s.io/community/contributors/guide/pull-requests.md#best-practices-for-faster-reviews)

### Proposals

For larger feature requests, please submit a design proposal in [docs/proposals/](./docs/proposals/). This is similar to a [Kubernetes Enhancement Proposal (KEP)](https://github.com/kubernetes/enhancements) or a [Architecture Decision Record (ADR)](https://adr.github.io/).

First, create a copy of the template found in the proposal directory, [docs/proposals/proposal-000-template.md](./docs/proposals/proposal-000-template.md). Rename the file to the next number in the sequence and add a name for the proposal e.g. `proposal-001-my-feature.md`. Fill in the required fields, then open a PR for review.

The initial PR can be a barebone PR with the goals/non-goals sections clarified that can be merged quickly and iterated on.

Initial merging of the PR does not mean that the proposal is approved. The status of the PR is defined in the `Status` section. Any KEP marked as `provisional` is a working document and subject to change.

A proposal that is accepted is a living document. Even the most well-planned ideas may change at some point.
