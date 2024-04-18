---
title: "Install Microservice Demo App"
description: "This section provides details around installation of Google microservice demo application."
summary: ""
date: 2023-09-07T16:04:48+02:00
lastmod: 2023-09-07T16:04:48+02:00
draft: false
slug: microservice-demo-install
weight: 880
toc: true
seo:
  title: "" # custom title (optional)
  description: "" # custom description (recommended)
  canonical: "" # custom canonical URL (optional)
  noindex: false # false (default) or true
---

## Manual installation

1. Get access to the dev-cluster (ask Tag/WG leads)

2. Run the following command: ```helm upgrade online-boutique oci://us-docker.pkg.dev/online-boutique-ci/charts/online-boutique --install```

3. To uninstall, run the following command: ```helm delete online-boutique```

# Install via Flux

We opted in for installing microservice demo using flux-oci. Example [here](https://fluxcd.io/flux/cheatsheets/oci-artifacts).
A temporaryconfiguration file can be found under [microservices-demo.yaml](https://github.com/cncf-tags/green-reviews-tooling/blob/main/clusters/projects/falco/microservices-demo.yaml).
