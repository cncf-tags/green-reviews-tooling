# How to install it manually

1. Get access to the dev-cluster (ask Tag/WG leads)

2. run the following command: 

```
helm upgrade online-boutique oci://us-docker.pkg.dev/online-boutique-ci/charts/online-boutique --install
```
3. to uninstall 

```
helm delete online-boutique
```

# How to install it via flux

We opted for installing microservice demo using flux-oci: example here: https://fluxcd.io/flux/cheatsheets/oci-artifacts/.
A tmp configuration file can be found under clusters/micrservices-demo.yaml


