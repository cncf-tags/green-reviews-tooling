# How to install it manually

1. Get access to the dev-cluster (ask Tag/WG leads)

2. run the following command: 

```
helm upgrade onlineboutique oci://us-docker.pkg.dev/online-boutique-ci/charts/onlineboutique --install
```
3. to uninstall 

```
helm delete onlineboutique
```

# How to install it via flux

We opted for installing microservice demo using flux-oci: example here: https://fluxcd.io/flux/cheatsheets/oci-artifacts/.
A tmp configuration file can be found under clusters/micrservices-demo.yaml


