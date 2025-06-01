# Infrastructure

The infrastructure used for this project leverages credits that Oracle Cloud kindly donated to the CNCF to support open-source projects.

Until June 2025 this project used credits donated by Equinix Metal. See [legacy](https://github.com/cncf-tags/green-reviews-tooling/blob/main/docs/infrastructure/legacy/) directory.

The Green Reviews WG uses Infrastructure as Code (IaC) tools as much as possible. For example, tools used in this project include [OpenTofu](https://opentofu.org) and [k3s](https://k3s.io) to provision the infrastructure and cluster.

The [benchmark pipeline](https://github.com/cncf-tags/green-reviews-tooling/blob/main/.github/workflows/benchmark-pipeline.yaml) workflow uses OpenTofu and the [Oracle Cloud Infrastructure (OCI) Provider](https://docs.oracle.com/en-us/iaas/Content/dev/terraform/home.htm) to provision the resources in the `infrastructure/oracle-cloud` directory.

```bash
./infrastructure
└── oracle-cloud
    ├── main.tf
    └── variables.tf
```

For each pipeline run we provision a bare metal instance plus required networking resources. The bare metal instance runs Ubuntu 24.04 and k3s. The opentofu stack outputs the k3s kubeconfig that is used by dagger to run the benchmark pipeline.

The [clean_iptables.sh](https://github.com/cncf-tags/green-reviews-tooling/blob/main/infrastructure/oracle-cloud/clean_iptables.sh) script is required to remove iptables rules that ship with the Oracle Ubuntu image that break cluster networking.

At the end of the pipeline run on success or failure we delete all resources via the opentofu stack.
