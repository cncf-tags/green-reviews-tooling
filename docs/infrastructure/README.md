# Infrastructure

- [Infrastructure](#infrastructure)
  - [Provisioning Infrastructure \& Cluster](#provisioning-infrastructure--cluster)
  - [Cluster Management](#cluster-management)
    - [Equinix Access](#equinix-access)
    - [Cluster Access](#cluster-access)
  - [Monitoring](#monitoring)

## Provisioning Infrastructure & Cluster

The infrastructure used for this project leverages infrastructure credits that Equinix kindly donated to the CNCF to support open-source projects.

The Green Reviews WG uses Infrastructure as Code (IaC) tools as much as possible. For example, tools used in this project include [OpenTofu](https://opentofu.org/) and [k3s](https://k3s.io/) to provision the infrastructure and cluster. Flux is used to provision the cluster components (see [Cluster Management](#cluster-management) below).

OpenTofu reconciles the files in [infrastructure/](../../infrastructure/):
```bash
./infrastructure
└── equinix-metal
    ├── main.tf
    └── variables.tf
```

OpenTofu provisions the Equinix Metal infrastructure using the [Equinix Terraform provider](https://github.com/equinix/terraform-provider-equinix). `user_data` provisions Kubernetes with k3s. It bootstrap the control plane, worker nodes, Cilium CNI, Flux (see below), etc.

The OpenTofu state is stored in an S3 bucket with AWS credits available through the CNCF. The authentication tokens are stored as Secrets in this repository. There is an S3 bucket in the CNCF AWS account available for the TAG ENV team to use. The bucket name is `tag-env-green-reviews-open-tofu` and the username is `tag-env-technical-user`. There is an ongoing process for creating a 1Password team account to store these secrets, see: https://github.com/cncf/tag-env-sustainability/issues/336

Isolation for the tests is ensured by deploying system-level components on one worker node and CNCF-project-specific components on a separate node (on a project-per-node basis).

## Cluster Management

For provisioning cluster components, the project uses [Flux](https://fluxcd.io/), which is a CNCF Graduated project. Flux manages the minimal set of components that should always be running to support the pipeline.

This GitOps approach was chosen so that it is:
- Clear to all participants which components and versions are running in the cluster
- Easier to contribute to technical tasks by submitting pull requests
- Easier for CNCF Project maintainers to deploy components that they maintain

Flux watches the manifests added to the [clusters/ directory](../clusters/) and applies or reconciles them in the cluster. The `clusters/` directory contains the following subdirecotries:
- [base/](../../clusters/base/) contains system-level applications e.g. [Kepler](https://www.cncf.io/projects/kepler/), [Prometheus](https://www.cncf.io/projects/prometheus/), and [Grafana](https://github.com/grafana/grafana). This is the base of this architectural reference, which is used to surface energy-level metrics to test the CNCF Projects.
- [projects/](../../clusters/projects/) contains project-specific configuration.
  - For example, `projects/falco/` deploys the manifests for Falco. The Falco installation is maintained by the Falco maintainers in the following repository: https://github.com/falcosecurity/cncf-green-review-testing. More info
- [flux-system/](../../clusters/flux-system/) contains the default files needed to operate Flux. This difectory was created when Flux was initially bootstrapped.

```bash
./clusters
├── base
│   ├── kepler-grafana.yaml # deploys the Kepler dashboard for Grafana
│   ├── kepler.yaml
│   ├── kube-prometheus-stack.yaml # deploys Prometheus + Grafana
│   └── monitoring-namespace.yaml
├── flux-system
│   └── <flux files>
└── projects
    └── falco
        ├── falco-namespace.yaml
        ├── falco.yaml
        └── microservices-demo.yaml
```

Please open a Pull Request to [add a component via its Helm Chart](https://fluxcd.io/flux/guides/helmreleases/).

For specific guides on how to deploy some of these components:
- [Install Flux](./install_flux.md)
- [Install Kepler](./install_kepler.md)
- [Microservice demo](./microservices_demo.md)

Flux authenticates with the `cncf-tags/green-reviews-tooling` repository via a fine-grained GitHub Personal Access Token (PAT). The GitHub PAT is a dedicated PAT for this repository alone (see [Flux docs](https://fluxcd.io/flux/installation/bootstrap/github/#github-organization) for GitHub Organizations).

## Access

### Equinix Access

Access to the Equinix Metal console has been granted to TAG & WG Chairs & TLs. Currently this is done via a manual process. TAG ENV leads can request access by following the instructions here: DM Jeffrey Sica in the CNCF or K8s slack (@jeefy) with the email to be granted access.

There is a [Equinix Project API Key](https://deploy.equinix.com/developers/docs/metal/accounts/api-keys/) with read/write permissions in the Equinix Project. These keys allow access to the Equinix Metal API. The API key and Project ID can be shared privately with maintainers who need this to test IaC tooling. Please reach out to one of the TAG/WG leads about this.

### Cluster Access

The project currently has two kubeconfig files: an admin kubeconfig and a read-only kubeconfig. The admin kubeconfig is used by TAG/WG Chairs & TLs. The read-only kubeconfig can be shared with any individual contributor.

If the read-only kubeconfig is needed to help a contributer, please open a Pull Request to make this request. Everyone is welcome to get read-only access. We ask that this request is made in a public channel purely to keep track of who has access to the Kubernetes cluster. Also, the kubeconfig may become redundant if the cluster is created again, so we would like to be able to share it again with you. In the future, we hope to create a 1Password team to store this kubeconfig and add individual contributors to the team.

Steps taken to create the read-only Kubeconfig can be found [here](./read-only-kubeconfig.md).

## Monitoring

One of the aims of this project is to create public data visualizations from the benchmark tests that can then be used for the assessments.

A public Grafana instance is available at http://147.28.134.41/ that can be accessed using the following credentials:
- Username: `admin`
- Password: `prom-operator`

To add a new Grafana dashboard, currently we deploy them as a ConfigMap managed by Flux. For example, in the Kepler dashboard [here](../../clusters/base/kepler-grafana.yaml), `data` contains the Grafana dashboard as a raw JSON object. The [Grafana sidecar for dashboards](https://github.com/grafana/helm-charts/tree/main/charts/grafana#sidecar-for-dashboards) watches all ConfigMaps and looks for the ones with the label (`metadata.labels`) `grafana_dashboard: "1"`.