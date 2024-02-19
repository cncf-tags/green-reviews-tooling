# Collaboration with CNCF Projects

This project aims to provide an architectural reference for how to assess the sustainability footprint of a cloud-native application (including the SCI) using cloud-native tooling. The community cluster aims to provide a playground for emerging metrics related to environmental sustainability, such as the [Software Carbon Intensity (SCI) specification](https://sci-guide.greensoftware.foundation/). The Green Reviews pipeline is also inspired by the process established by [TAG Security’s Security Assessment (TSSA)](https://github.com/cncf/tag-security/tree/main/assessments) for CNCF Projects.

## CNCF Project Maintainer Responsibilities

There are different roles and responsibilities involved in the review of the CNCF Projects. What do CNCF Project Maintainers provide? What does the Green Reviews WG provide? These questions are being answered through the current active collaboration with the maintainers of Falco.

There are certain key differences around the deployment and implementation that vary for each project. CNCF Project Maintainers can help with the following:
- Provide a way to install the CNCF Project in the community cluster and c
- Share any requirements for deploying the CNCF Project or running the benchmarks
- Contribute test scenarios for each of the specific CNCF Project, which should roughly be equivalent to the [functional unit](https://sci-guide.greensoftware.foundation/R) of the tool

CNCF Projects are welcome to flag their interest to collaborate with the Green Reviews WG by leaving a comment in [this issue](https://github.com/cncf/tag-env-sustainability/issues/223).

## Project 1: Falco

Falco is the first project to go through the TAG ENV Green Review pipeline - more info [here](./falco/README.md).