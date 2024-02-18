# Install Flux

To install flux manually run the following: 

```sh
export GITHUB_TOKEN=our_path
export GITHUB_USER=cncf-tags
flux bootstrap github --owner=$GITHUB_USER --repository=green-reviews-tooling --path=clusters
```

# SOPS encryption

Mozilla SOPS is used to encrypt secrets using GPG.

> Refer: [Link](https://fluxcd.io/flux/guides/mozilla-sops/)

## Import public key

Import GPG public key:

```sh
gpg --import ./clusters/.sops.pub.asc
```

## Encrypt secret

Generate manifest:

```sh
kubectl -n default create secret generic basic-auth \
--from-literal=user=admin \
--from-literal=password=change-me \
--dry-run=client \
-o yaml > clusters/basic-auth.yaml
```

Encrypt secret:

```sh
cd clusters
sops --encrypt --in-place basic-auth.yaml
```
