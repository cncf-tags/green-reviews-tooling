---
title: "Install Flux"
description: "This section provides details around installation of Flux CD tool."
summary: ""
date: 2023-09-07T16:04:48+02:00
lastmod: 2023-09-07T16:04:48+02:00
draft: false
slug: flux-install
weight: 860
toc: true
---

To install Flux manually run the following:

``` sh
export GITHUB_TOKEN=our_path
export GITHUB_USER=cncf-tags
flux bootstrap github --owner=$GITHUB_USER --repository=green-reviews-tooling --path=clusters
```
